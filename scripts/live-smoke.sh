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
  team_id="$(python3 -c 'import json, sys; print(json.load(sys.stdin)["team"]["id"])' <<<"$target_json")"
  project_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); print(data.get("project", {}).get("id", ""))' <<<"$target_json")"
  if [[ -n "${LINCTL_APPLICATION_CLIENT_ID:-}" ]]; then
    "$binary" application info "$LINCTL_APPLICATION_CLIENT_ID" --json >/dev/null
  fi
  "$binary" agent-activity list --json --limit 5 >/dev/null
  "$binary" agent-skill list --json --limit 5 >/dev/null
  "$binary" audit-entry types --json >/dev/null
  "$binary" organization exists "$org_url_key" --json >/dev/null
  "$binary" organization labels --json --limit 5 >/dev/null
  "$binary" organization project-labels --json --limit 5 >/dev/null
  "$binary" organization teams --json --limit 5 >/dev/null
  "$binary" organization templates --json --limit 5 >/dev/null
  "$binary" organization users --json --limit 5 >/dev/null
  "$binary" rate-limit status --json >/dev/null
  viewer_json="$("$binary" whoami --json)"
  user_id="$(python3 -c 'import json, sys; print(json.load(sys.stdin)["id"])' <<<"$viewer_json")"
  "$binary" user drafts --json --limit 5 >/dev/null
  "$binary" user assigned-issues "$user_id" --json --limit 5 >/dev/null
  "$binary" user created-issues "$user_id" --json --limit 5 >/dev/null
  "$binary" user delegated-issues "$user_id" --json --limit 5 >/dev/null
  "$binary" user team-memberships "$user_id" --json --limit 5 >/dev/null
  "$binary" user teams "$user_id" --json --limit 5 >/dev/null
  "$binary" user my-assigned-issues --json --limit 5 >/dev/null
  "$binary" user my-created-issues --json --limit 5 >/dev/null
  "$binary" user my-delegated-issues --json --limit 5 >/dev/null
  "$binary" user my-team-memberships --json --limit 5 >/dev/null
  "$binary" user my-teams --json --limit 5 >/dev/null
  "$binary" issue usage >/dev/null
  issue_json="$("$binary" issue list --json --limit 5)"
  issue_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("issues", []); print(items[0]["id"] if items else "")' <<<"$issue_json")"
  if [[ -n "$issue_id" ]]; then
    "$binary" issue attachments "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue bot-actor "$issue_id" --json >/dev/null
    "$binary" issue children "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue documents "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue former-attachments "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue history "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue inverse-relations "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue labels "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue relations "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue releases "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue state-history "$issue_id" --json --limit 5 >/dev/null
    "$binary" issue subscribers "$issue_id" --json --limit 5 >/dev/null
  fi
  comment_json="$("$binary" comment list --json --limit 5)"
  comment_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("comments", []); print(items[0]["id"] if items else "")' <<<"$comment_json")"
  if [[ -n "$comment_id" ]]; then
    "$binary" comment bot-actor "$comment_id" --json >/dev/null
    "$binary" comment children "$comment_id" --json --limit 5 >/dev/null
    "$binary" comment created-issues "$comment_id" --json --limit 5 >/dev/null
  fi
  "$binary" issue-relation list --json --limit 5 >/dev/null
  "$binary" issue-to-release list --json --limit 5 >/dev/null
  "$binary" project usage >/dev/null
  "$binary" project list --json --limit 5 >/dev/null
  if [[ -n "$project_id" ]]; then
    "$binary" project comments "$project_id" --json --limit 5 >/dev/null
    project_milestone_json="$("$binary" project-milestone list "$project_id" --json --limit 5)"
    project_milestone_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("milestones", []); print(items[0]["id"] if items else "")' <<<"$project_milestone_json")"
    if [[ -n "$project_milestone_id" ]]; then
      "$binary" project-milestone issues "$project_milestone_id" --json --limit 5 >/dev/null
    fi
  fi
  project_update_json="$("$binary" project-update list --json --limit 5)"
  project_update_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("updates", []); print(items[0]["id"] if items else "")' <<<"$project_update_json")"
  if [[ -n "$project_update_id" ]]; then
    "$binary" project-update comments "$project_update_id" --json --limit 5 >/dev/null
  fi
  "$binary" project-status list --json --limit 5 >/dev/null
  "$binary" project-label list --json --limit 5 >/dev/null
  "$binary" project-relation list --json --limit 5 >/dev/null
  label_json="$("$binary" label list --json --limit 5)"
  label_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("labels", []); print(items[0]["id"] if items else "")' <<<"$label_json")"
  if [[ -n "$label_id" ]]; then
    "$binary" label children "$label_id" --json --limit 5 >/dev/null
    "$binary" label issues "$label_id" --json --limit 5 >/dev/null
  fi
  cycle_json="$("$binary" cycle list --json --limit 5)"
  cycle_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("cycles", []); print(items[0]["id"] if items else "")' <<<"$cycle_json")"
  if [[ -n "$cycle_id" ]]; then
    "$binary" cycle issues "$cycle_id" --json --limit 5 >/dev/null
    "$binary" cycle uncompleted-issues "$cycle_id" --json --limit 5 >/dev/null
  fi
  "$binary" team cycles "$team_id" --json --limit 5 >/dev/null
  "$binary" team issues "$team_id" --json --limit 5 >/dev/null
  "$binary" team labels "$team_id" --json --limit 5 >/dev/null
  "$binary" team memberships "$team_id" --json --limit 5 >/dev/null
  "$binary" team projects "$team_id" --json --limit 5 >/dev/null
  "$binary" team release-pipelines "$team_id" --json --limit 5 >/dev/null
  "$binary" team states "$team_id" --json --limit 5 >/dev/null
  "$binary" team templates "$team_id" --json --limit 5 >/dev/null
  "$binary" team-membership list --json --limit 5 >/dev/null
  "$binary" notification list --json --limit 5 >/dev/null
  "$binary" notification subscription list --json --limit 5 >/dev/null
  triage_responsibility_json="$("$binary" triage-responsibility list --json --limit 5)"
  triage_responsibility_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("triage_responsibilities", []); print(items[0]["id"] if items else "")' <<<"$triage_responsibility_json")"
  if [[ -n "$triage_responsibility_id" ]]; then
    "$binary" triage-responsibility get "$triage_responsibility_id" --json >/dev/null
    "$binary" triage-responsibility manual-selection "$triage_responsibility_id" --json >/dev/null
  fi
  "$binary" sla-configuration list "$team_id" --json >/dev/null
  "$binary" semantic-search linear --json --limit 1 >/dev/null
  "$binary" search documents linear --json --limit 1 >/dev/null
  "$binary" search issues linear --json --limit 1 >/dev/null
  "$binary" search projects linear --json --limit 1 >/dev/null
  release_pipeline_json="$("$binary" release-pipeline list --json --limit 5)"
  release_pipeline_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("release_pipelines", []); print(items[0]["id"] if items else "")' <<<"$release_pipeline_json")"
  if [[ -n "$release_pipeline_id" ]]; then
    "$binary" release-pipeline releases "$release_pipeline_id" --json --limit 5 >/dev/null
    "$binary" release-pipeline stages "$release_pipeline_id" --json --limit 5 >/dev/null
    "$binary" release-pipeline teams "$release_pipeline_id" --json --limit 5 >/dev/null
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
    "$binary" release documents "$release_id" --json --limit 5 >/dev/null
    "$binary" release issues "$release_id" --json --limit 5 >/dev/null
    release_links_json="$("$binary" release links "$release_id" --json --limit 5)"
    external_link_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("links", []); print(items[0]["id"] if items else "")' <<<"$release_links_json")"
    if [[ -n "$external_link_id" ]]; then
      "$binary" external-link get "$external_link_id" --json >/dev/null
    fi
  fi
  "$binary" release-note list --json --limit 5 >/dev/null
  external_user_json="$("$binary" external-user list --json --limit 5)"
  external_user_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("external_users", []); print(items[0]["id"] if items else "")' <<<"$external_user_json")"
  if [[ -n "$external_user_id" ]]; then
    "$binary" external-user get "$external_user_id" --json >/dev/null
  fi
  "$binary" time-schedule list --json --limit 5 >/dev/null
  template_json="$("$binary" template list --json --limit 5)"
  template_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("templates", []); print(items[0]["id"] if items else "")' <<<"$template_json")"
  if [[ -n "$template_id" ]]; then
    "$binary" template get "$template_id" --json >/dev/null
  fi
  initiative_json="$("$binary" initiative list --json --limit 5)"
  initiative_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("initiatives", []); print(items[0]["id"] if items else "")' <<<"$initiative_json")"
  if [[ -n "$initiative_id" ]]; then
    "$binary" initiative history "$initiative_id" --json --limit 5 >/dev/null
    "$binary" initiative links "$initiative_id" --json --limit 5 >/dev/null
    "$binary" initiative sub-initiatives "$initiative_id" --json --limit 5 >/dev/null
    "$binary" initiative updates "$initiative_id" --json --limit 5 >/dev/null
    "$binary" initiative documents "$initiative_id" --json --limit 5 >/dev/null
    "$binary" initiative projects "$initiative_id" --json --limit 5 >/dev/null
  fi
  "$binary" initiative-relation list --json --limit 5 >/dev/null
  "$binary" initiative-to-project list --json --limit 5 >/dev/null
  "$binary" initiative-update list --json --limit 5 >/dev/null
  "$binary" roadmap list --json --limit 5 >/dev/null
  "$binary" roadmap-to-project list --json --limit 5 >/dev/null
  custom_view_json="$("$binary" custom-view list --json --limit 5)"
  custom_view_id="$(python3 -c 'import json, sys; data=json.load(sys.stdin); items=data.get("custom_views", []); print(items[0]["id"] if items else "")' <<<"$custom_view_json")"
  if [[ -n "$custom_view_id" ]]; then
    "$binary" custom-view get "$custom_view_id" --json >/dev/null
    "$binary" custom-view subscribers "$custom_view_id" --json >/dev/null
    "$binary" custom-view initiatives "$custom_view_id" --json --limit 5 >/dev/null
    "$binary" custom-view issues "$custom_view_id" --json --limit 5 >/dev/null
    "$binary" custom-view organization-preferences "$custom_view_id" --json >/dev/null
    "$binary" custom-view organization-preferences values "$custom_view_id" --json >/dev/null
    "$binary" custom-view projects "$custom_view_id" --json --limit 5 >/dev/null
    "$binary" custom-view user-preferences "$custom_view_id" --json >/dev/null
    "$binary" custom-view user-preferences values "$custom_view_id" --json >/dev/null
    "$binary" custom-view preference-values "$custom_view_id" --json >/dev/null
  fi
  "$binary" customer list --json --limit 5 >/dev/null
  "$binary" customer-need list --json --limit 5 >/dev/null
  "$binary" customer-status list --json --limit 5 >/dev/null
  "$binary" customer-tier list --json --limit 5 >/dev/null
)

go test -count=1 -tags=integration ./internal/client
