#!/usr/bin/env bash
set -euo pipefail

token="${LINCTL_TEST_TOKEN:-${LINCTL_TOKEN:-${LINEAR_API_KEY:-}}}"
if [[ -z "$token" ]]; then
  printf 'missing disposable Linear token: set LINCTL_TEST_TOKEN, LINCTL_TOKEN, or LINEAR_API_KEY\n' >&2
  exit 2
fi

export LINCTL_TEST_TOKEN="$token"
export LINCTL_TOKEN="$token"

binary="$(mktemp -t linctl-live-smoke.XXXXXX)"
smoke_dir="$(mktemp -d -t linctl-live-smoke.XXXXXX)"
trap 'rm -f "$binary"; rm -rf "$smoke_dir"' EXIT

go build -trimpath -o "$binary" ./cmd/linctl
python3 - test/integration-config.json "$smoke_dir/.linctl.toml" <<'PY'
import json
import sys

input_path = sys.argv[1]
output_path = sys.argv[2]
with open(input_path, "r", encoding="utf-8") as input_file:
    config = json.load(input_file)
with open(output_path, "w", encoding="utf-8") as output:
    output.write("[target]\n")
    output.write(f'org_id = "{config["org_id"]}"\n')
    output.write(f'team_key = "{config["team_key"]}"\n')
    output.write(f'team_id = "{config["team_id"]}"\n')
    output.write(f'project_id = "{config["project_id"]}"\n')
PY

(
  cd "$smoke_dir"
  "$binary" usage >/dev/null
  target_json="$("$binary" target --json)"
  org_url_key="$(python3 -c 'import json, sys; print(json.load(sys.stdin)["org"]["url_key"])' <<<"$target_json")"
  "$binary" organization exists "$org_url_key" --json >/dev/null
  "$binary" rate-limit status --json >/dev/null
  "$binary" whoami --json >/dev/null
  "$binary" issue usage >/dev/null
  "$binary" issue list --json --limit 5 >/dev/null
  "$binary" project usage >/dev/null
  "$binary" project list --json --limit 5 >/dev/null
  "$binary" notification list --json --limit 5 >/dev/null
  "$binary" notification subscription list --json --limit 5 >/dev/null
  release_pipeline_json="$("$binary" release-pipeline list --json --limit 5)"
  release_pipeline_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("release_pipelines", []); print(items[0]["id"] if items else "")' <<<"$release_pipeline_json")"
  if [[ -n "$release_pipeline_id" ]]; then
    "$binary" release-pipeline releases "$release_pipeline_id" --json --limit 5 >/dev/null
    "$binary" release-pipeline stages "$release_pipeline_id" --json --limit 5 >/dev/null
  fi
  release_stage_json="$("$binary" release-stage list --json --limit 5)"
  release_stage_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("release_stages", []); print(items[0]["id"] if items else "")' <<<"$release_stage_json")"
  if [[ -n "$release_stage_id" ]]; then
    "$binary" release-stage releases "$release_stage_id" --json --limit 5 >/dev/null
  fi
  release_json="$("$binary" release list --json --limit 5)"
  release_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("releases", []); print(items[0]["id"] if items else "")' <<<"$release_json")"
  if [[ -n "$release_id" ]]; then
    "$binary" release history "$release_id" --json --limit 5 >/dev/null
    "$binary" release links "$release_id" --json --limit 5 >/dev/null
  fi
  "$binary" release-note list --json --limit 5 >/dev/null
  "$binary" time-schedule list --json --limit 5 >/dev/null
  "$binary" initiative-relation list --json --limit 5 >/dev/null
  "$binary" initiative-to-project list --json --limit 5 >/dev/null
  "$binary" initiative-update list --json --limit 5 >/dev/null
  "$binary" roadmap list --json --limit 5 >/dev/null
  "$binary" customer list --json --limit 5 >/dev/null
  "$binary" customer-need list --json --limit 5 >/dev/null
  "$binary" customer-status list --json --limit 5 >/dev/null
  "$binary" customer-tier list --json --limit 5 >/dev/null
)

go test -count=1 -tags=integration ./internal/client
