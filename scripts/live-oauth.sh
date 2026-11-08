#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

if ! command -v python3 >/dev/null 2>&1; then
  printf 'python3 is required to run the live OAuth harness\n' >&2
  exit 1
fi

required_env=(
  LINCTL_OAUTH_CLIENT_ID
  LINCTL_OAUTH_CLIENT_SECRET
  LINCTL_OAUTH_REDIRECT_URI
  LINCTL_OAUTH_SCOPES
  LINCTL_OAUTH_EXPECTED_ACTOR
)
missing=()
for key in "${required_env[@]}"; do
  if [[ -z "${!key:-}" ]]; then
    missing+=("$key")
  fi
done
if ((${#missing[@]} > 0)); then
  printf 'missing OAuth fixture env: set %s\n' "${missing[*]}" >&2
  exit 2
fi
if [[ "$LINCTL_OAUTH_EXPECTED_ACTOR" != "app" ]]; then
  printf 'LINCTL_OAUTH_EXPECTED_ACTOR must be app for the live OAuth fixture\n' >&2
  exit 2
fi

binary="${LINCTL_BINARY:-${1:-./bin/linctl}}"
if [[ ! -x "$binary" ]]; then
  printf 'linctl binary is not executable: %s\n' "$binary" >&2
  exit 2
fi

app_output="$(mktemp -t linctl-live-oauth-app.XXXXXX)"
status_output="$(mktemp -t linctl-live-oauth-status.XXXXXX)"
trap 'rm -f "$app_output" "$status_output"' EXIT

validate_auth_json() {
  local path="$1"
  local expected_actor="$2"
  local mode="$3"
  python3 - "$path" "$expected_actor" "$mode" <<'PY'
import json
import sys

path = sys.argv[1]
expected_actor = sys.argv[2]
mode = sys.argv[3]

with open(path, "r", encoding="utf-8") as input_file:
    payload = json.load(input_file)


def fail(message):
    sys.stderr.write(message + "\n")
    sys.exit(1)


prefix = "live OAuth status " if mode == "status" else "live OAuth "

if payload.get("actor") != expected_actor:
    fail(prefix + "actor mismatch")
if payload.get("target", {}).get("status") != "ready":
    fail(prefix + "target is not ready")
if payload.get("token", {}).get("status") != "set":
    token_message = "token is not set" if mode == "status" else "token was not saved"
    fail(prefix + token_message)
app = payload.get("app", {})
if app.get("client_id") != "set" or app.get("client_secret") != "set":
    fail(prefix + "app material was not redacted")
PY
}

"$binary" auth configure \
  --client-id "$LINCTL_OAUTH_CLIENT_ID" \
  --client-secret "$LINCTL_OAUTH_CLIENT_SECRET" \
  --redirect-uri "$LINCTL_OAUTH_REDIRECT_URI" \
  --scopes "$LINCTL_OAUTH_SCOPES" \
  --quiet

"$binary" auth app --json >"$app_output"
validate_auth_json "$app_output" "$LINCTL_OAUTH_EXPECTED_ACTOR" app

"$binary" auth status --json >"$status_output"
validate_auth_json "$status_output" "$LINCTL_OAUTH_EXPECTED_ACTOR" status

if [[ "${LINCTL_LIVE_OAUTH_RECOVERY:-}" == "1" ]]; then
  other_state="$(mktemp -d -t linctl-live-oauth-other.XXXXXX)"
  recovered_output="$(mktemp -t linctl-live-oauth-recovered.XXXXXX)"
  injected_error="$(mktemp -t linctl-live-oauth-injected.XXXXXX)"
  trap 'rm -f "$app_output" "$status_output" "$recovered_output" "$injected_error"; rm -rf "${other_state:?}"' EXIT

  env -u LINCTL_OAUTH_SCOPES XDG_CONFIG_HOME="$other_state/config" XDG_STATE_HOME="$other_state/state" \
    "$binary" auth app --scopes read --quiet
  "$binary" auth status --json >"$recovered_output"
  validate_auth_json "$recovered_output" "$LINCTL_OAUTH_EXPECTED_ACTOR" status
  python3 - "$status_output" "$recovered_output" <<'PY'
import json
import sys

with open(sys.argv[1], "r", encoding="utf-8") as before_file:
    before = json.load(before_file)
with open(sys.argv[2], "r", encoding="utf-8") as after_file:
    after = json.load(after_file)
if before.get("expires_at") == after.get("expires_at"):
    sys.stderr.write("live OAuth status did not replace the token Linear rejected\n")
    sys.exit(1)
PY

  if LINCTL_OAUTH_ACCESS_TOKEN=lin_oauth_linctl_live_rejected "$binary" auth status --json >/dev/null 2>"$injected_error"; then
    printf 'live OAuth status accepted an invalid injected token\n' >&2
    exit 1
  fi
  if ! grep -q '"error_code":"AUTH_REAUTH_REQUIRED"' "$injected_error"; then
    printf 'live OAuth status did not report AUTH_REAUTH_REQUIRED for an invalid injected token\n' >&2
    exit 1
  fi
  printf 'live OAuth recovery ok: rejected token replaced, invalid injected token refused\n'
fi

printf 'live OAuth ok: actor=app target=ready token=set scopes=set\n'
