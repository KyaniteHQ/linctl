#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

if ! command -v python3 >/dev/null 2>&1; then
  printf 'python3 is required to run the live sweep\n' >&2
  exit 1
fi

required_env=(
  LINCTL_OAUTH_CLIENT_ID
  LINCTL_OAUTH_CLIENT_SECRET
  LINCTL_OAUTH_REDIRECT_URI
  LINCTL_OAUTH_SCOPES
  LINCTL_OAUTH_EXPECTED_ACTOR
  LINCTL_TEST_ORG_ID
  LINCTL_TEST_TEAM_KEY
  LINCTL_TEST_TEAM_ID
)
missing_env=()
for key in "${required_env[@]}"; do
  if [[ -z "${!key:-}" ]]; then
    missing_env+=("$key")
  fi
done
if ((${#missing_env[@]} > 0)); then
  printf 'missing fixture env for live sweep: set %s\n' "${missing_env[*]}" >&2
  exit 2
fi

namespace_prefix="linctl-it-"

binary="${LINCTL_BINARY:-}"
cleanup_binary=0
if [[ -z "$binary" ]]; then
  binary="$(mktemp -t linctl-live-sweep.XXXXXX)"
  cleanup_binary=1
  go build -trimpath -o "$binary" ./cmd/linctl
fi

sweep_dir="$(mktemp -d -t linctl-live-sweep.XXXXXX)"
repo_root="$(pwd -P)"
cleanup() {
  rm -rf "$sweep_dir"
  if ((cleanup_binary)); then
    rm -f "$binary"
  fi
}
trap cleanup EXIT

{
  printf '[target]\n'
  printf 'org_id = "%s"\n' "$LINCTL_TEST_ORG_ID"
  printf 'team_key = "%s"\n' "$LINCTL_TEST_TEAM_KEY"
  printf 'team_id = "%s"\n' "$LINCTL_TEST_TEAM_ID"
  printf 'project_id = "%s"\n' "${LINCTL_TEST_PROJECT_ID:-}"
} >"$sweep_dir/.linctl.toml"

filter_ids_by_prefix() {
  local collection="$1"
  local id_field="$2"
  local name_field="$3"
  local archived_field="${4:-}"
  local status_field="${5:-}"
  python3 -c '
import json, sys

collection, id_field, name_field, prefix, archived_field, status_field = sys.argv[1:7]
data = json.load(sys.stdin)
for item in data.get(collection, []):
    archived = archived_field and item.get(archived_field)
    terminal = status_field and item.get(status_field, {}).get("type") == "canceled"
    active = not archived and not terminal
    if active and item.get(name_field, "").startswith(prefix):
        print(item[id_field])
' "$collection" "$id_field" "$name_field" "$namespace_prefix" "$archived_field" "$status_field"
}

(
  export XDG_CONFIG_HOME="$sweep_dir/config"
  export XDG_STATE_HOME="$sweep_dir/state"
  export LINCTL_BINARY="$binary"

  cd "$sweep_dir"
  (cd "$repo_root" && go tool task --taskfile "$repo_root/Taskfile.yml" --dir "$sweep_dir" live-oauth) >/dev/null

  swept=0

  issue_json="$("$binary" issue list --json --limit 250)"
  while IFS= read -r issue_id; do
    [[ -z "$issue_id" ]] && continue
    "$binary" issue close "$issue_id" >/dev/null
    swept=$((swept + 1))
  done < <(filter_ids_by_prefix issues id title <<<"$issue_json")

  project_json="$("$binary" project list --json --limit 250)"
  while IFS= read -r project_id; do
    [[ -z "$project_id" ]] && continue
    "$binary" project archive "$project_id" >/dev/null
    swept=$((swept + 1))
  done < <(filter_ids_by_prefix projects id name archived_at status <<<"$project_json")

  printf 'live-sweep: closed/archived %d namespaced (%s) resource(s)\n' "$swept" "$namespace_prefix"
)
