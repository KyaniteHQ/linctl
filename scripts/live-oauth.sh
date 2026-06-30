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

"$binary" auth configure \
  --client-id "$LINCTL_OAUTH_CLIENT_ID" \
  --client-secret "$LINCTL_OAUTH_CLIENT_SECRET" \
  --redirect-uri "$LINCTL_OAUTH_REDIRECT_URI" \
  --scopes "$LINCTL_OAUTH_SCOPES" \
  --quiet

"$binary" auth app --json >"$app_output"
python3 - "$app_output" "$LINCTL_OAUTH_EXPECTED_ACTOR" <<'PY'
import json
import sys

path = sys.argv[1]
expected_actor = sys.argv[2]

with open(path, "r", encoding="utf-8") as input_file:
    payload = json.load(input_file)


def fail(message):
    sys.stderr.write(message + "\n")
    sys.exit(1)


if payload.get("actor") != expected_actor:
    fail("live OAuth actor mismatch")
if payload.get("target", {}).get("status") != "ready":
    fail("live OAuth target is not ready")
if payload.get("token", {}).get("status") != "set":
    fail("live OAuth token was not saved")
app = payload.get("app", {})
if app.get("client_id") != "set" or app.get("client_secret") != "set":
    fail("live OAuth app material was not redacted")
PY

"$binary" auth status --json >"$status_output"
python3 - "$status_output" "$LINCTL_OAUTH_EXPECTED_ACTOR" <<'PY'
import json
import sys

path = sys.argv[1]
expected_actor = sys.argv[2]

with open(path, "r", encoding="utf-8") as input_file:
    payload = json.load(input_file)


def fail(message):
    sys.stderr.write(message + "\n")
    sys.exit(1)


if payload.get("actor") != expected_actor:
    fail("live OAuth status actor mismatch")
if payload.get("target", {}).get("status") != "ready":
    fail("live OAuth status target is not ready")
if payload.get("token", {}).get("status") != "set":
    fail("live OAuth status token is not set")
app = payload.get("app", {})
if app.get("client_id") != "set" or app.get("client_secret") != "set":
    fail("live OAuth status app material was not redacted")
PY

printf 'live OAuth ok: actor=app target=ready token=set scopes=set\n'
