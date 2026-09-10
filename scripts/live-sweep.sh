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
# Matches defaultRetries in internal/client/transport.go: 3 retries after the first try.
max_transient_retries=3

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
    terminal = status_field and item.get(status_field, {}).get("type") in {"canceled", "completed"}
    active = not archived and not terminal
    if active and item.get(name_field, "").startswith(prefix):
        print(item[id_field])
' "$collection" "$id_field" "$name_field" "$namespace_prefix" "$archived_field" "$status_field"
}

is_unarchivable_project() {
  local error_path="$1"
  python3 - "$error_path" <<'PY'
import json
import sys

for line in open(sys.argv[1], encoding="utf-8"):
    try:
        error = json.loads(line)
    except json.JSONDecodeError:
        continue
    if (
        error.get("error_code") == "GRAPHQL_ERROR"
        and "projectArchive cannot delete project (INPUT_ERROR)" in error.get("message", "")
    ):
        sys.exit(0)
sys.exit(1)
PY
}

is_transient_http_error() {
  local error_path="$1"
  python3 - "$error_path" <<'PY'
import json
import re
import sys

pattern = re.compile(r"graphql http status 50[23]\b")
for line in open(sys.argv[1], encoding="utf-8"):
    try:
        error = json.loads(line)
    except json.JSONDecodeError:
        if pattern.search(line):
            sys.exit(0)
        continue
    if pattern.search(str(error.get("message", ""))):
        sys.exit(0)
sys.exit(1)
PY
}

transient_retry_delay() {
  local attempt="$1"
  local error_path="$2"
  python3 - "$attempt" "$error_path" <<'PY'
import json
import sys

attempt = int(sys.argv[1])
max_delay = 30.0
retry_after = None
for line in open(sys.argv[2], encoding="utf-8"):
    try:
        error = json.loads(line)
    except json.JSONDecodeError:
        continue
    for key in ("retry_after", "retryAfter"):
        value = error.get(key)
        if isinstance(value, int) and value > 0:
            retry_after = float(value)
        elif isinstance(value, str) and value.isdigit() and int(value) > 0:
            retry_after = float(value)
if retry_after is not None:
    print(min(retry_after, max_delay))
else:
    print(min((attempt + 1) * 0.1, max_delay))
PY
}

run_write_with_retry() {
  local error_path="$1"
  shift
  local attempt=0
  while true; do
    if "$@" >/dev/null 2>"$error_path"; then
      return 0
    fi
    if is_transient_http_error "$error_path" && ((attempt < max_transient_retries)); then
      sleep "$(transient_retry_delay "$attempt" "$error_path")"
      attempt=$((attempt + 1))
      continue
    fi
    return 1
  done
}

report_write_failure() {
  local kind="$1"
  local resource_id="$2"
  local error_path="$3"
  printf 'live-sweep: failed to %s %s\n' "$kind" "$resource_id" >&2
  cat "$error_path" >&2
}

(
  export XDG_CONFIG_HOME="$sweep_dir/config"
  export XDG_STATE_HOME="$sweep_dir/state"
  export LINCTL_BINARY="$binary"

  cd "$sweep_dir"
  (cd "$repo_root" && go tool task --taskfile "$repo_root/Taskfile.yml" --dir "$sweep_dir" live-oauth) >/dev/null

  swept=0
  failed=0
  skipped=0

  issue_json="$("$binary" issue list --json --limit 250)"
  while IFS= read -r issue_id; do
    [[ -z "$issue_id" ]] && continue
    issue_error="$sweep_dir/issue-close.err"
    if run_write_with_retry "$issue_error" "$binary" issue close "$issue_id"; then
      swept=$((swept + 1))
    else
      report_write_failure "close issue" "$issue_id" "$issue_error"
      failed=$((failed + 1))
    fi
  done < <(filter_ids_by_prefix issues id title <<<"$issue_json")

  project_json="$("$binary" project list --json --limit 250)"
  while IFS= read -r project_id; do
    [[ -z "$project_id" ]] && continue
    project_error="$sweep_dir/project-archive.err"
    if run_write_with_retry "$project_error" "$binary" project archive "$project_id"; then
      swept=$((swept + 1))
    elif is_unarchivable_project "$project_error"; then
      skipped=$((skipped + 1))
    else
      report_write_failure "archive project" "$project_id" "$project_error"
      failed=$((failed + 1))
    fi
  done < <(filter_ids_by_prefix projects id name archived_at status <<<"$project_json")

  printf 'live-sweep: closed/archived %d namespaced (%s) resource(s)\n' "$swept" "$namespace_prefix"
  if ((skipped)); then
    printf 'live-sweep: skipped %d namespaced project(s) rejected as unarchivable by Linear\n' "$skipped" >&2
  fi
  if ((failed)); then
    printf 'live-sweep: %d namespaced resource(s) still failed after retries\n' "$failed" >&2
    exit 1
  fi
)
