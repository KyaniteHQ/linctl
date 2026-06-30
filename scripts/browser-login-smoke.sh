#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

if ! command -v python3 >/dev/null 2>&1; then
  printf 'python3 is required to run the browser login smoke\n' >&2
  exit 1
fi

actor="${1:-${LINCTL_BROWSER_LOGIN_ACTOR:-user}}"
if (($# > 1)); then
  printf 'usage: scripts/browser-login-smoke.sh [user|app]\n' >&2
  exit 2
fi
case "$actor" in
  app | user) ;;
  *)
    printf 'browser login actor must be app or user\n' >&2
    exit 2
    ;;
esac

required_env=(
  LINCTL_OAUTH_CLIENT_ID
  LINCTL_OAUTH_REDIRECT_URI
)
missing=()
for key in "${required_env[@]}"; do
  if [[ -z "${!key:-}" ]]; then
    missing+=("$key")
  fi
done
if ((${#missing[@]} > 0)); then
  printf 'missing browser login fixture env: set %s\n' "${missing[*]}" >&2
  exit 2
fi

scopes="${LINCTL_OAUTH_SCOPES:-read,write,issues:create,comments:create}"
callback_timeout="${LINCTL_BROWSER_LOGIN_CALLBACK_TIMEOUT:-300}"
if ! [[ "$callback_timeout" =~ ^[0-9]+$ ]] || ((callback_timeout < 1)); then
  printf 'LINCTL_BROWSER_LOGIN_CALLBACK_TIMEOUT must be a positive integer number of seconds\n' >&2
  exit 2
fi
binary="${LINCTL_BINARY:-}"
built_binary=""
smoke_dir="$(mktemp -d -t linctl-browser-login-smoke.XXXXXX)"

cleanup() {
  if [[ -n "$built_binary" ]]; then
    rm -f "$built_binary"
  fi
  rm -rf "$smoke_dir"
}
trap cleanup EXIT

if [[ -z "$binary" ]]; then
  built_binary="$(mktemp -t linctl-browser-login.XXXXXX)"
  go build -trimpath -o "$built_binary" ./cmd/linctl
  binary="$built_binary"
elif [[ ! -x "$binary" ]]; then
  printf 'linctl binary is not executable: %s\n' "$binary" >&2
  exit 2
fi

# Resolve the pinned target from env vars first, falling back to the same local
# untracked integration config used by the full live smoke.
python3 - "${LINCTL_TEST_CONFIG:-test/integration-config.json}" "$smoke_dir/.linctl.toml" <<'PY'
import json
import os
import sys

input_path = sys.argv[1]
output_path = sys.argv[2]

env_keys = {
    "org_id": "LINCTL_TEST_ORG_ID",
    "team_key": "LINCTL_TEST_TEAM_KEY",
    "team_id": "LINCTL_TEST_TEAM_ID",
    "project_id": "LINCTL_TEST_PROJECT_ID",
}
config = {key: os.environ.get(env, "") for key, env in env_keys.items()}

if not (config["org_id"] and config["team_key"] and config["team_id"]):
    if not os.path.exists(input_path):
        sys.stderr.write(
            "missing integration config: set LINCTL_TEST_ORG_ID, "
            "LINCTL_TEST_TEAM_KEY, LINCTL_TEST_TEAM_ID (and optional "
            f"LINCTL_TEST_PROJECT_ID), or provide {input_path}\n"
        )
        sys.exit(2)
    with open(input_path, "r", encoding="utf-8") as input_file:
        config = json.load(input_file)

with open(output_path, "w", encoding="utf-8") as output:
    output.write("[target]\n")
    output.write(f'org_id = "{config["org_id"]}"\n')
    output.write(f'team_key = "{config["team_key"]}"\n')
    output.write(f'team_id = "{config["team_id"]}"\n')
    output.write(f'project_id = "{config.get("project_id", "")}"\n')
PY

validate_auth_json() {
  local path="$1"
  local expected_actor="$2"
  local mode="$3"
  local client_secret_expected="$4"
  python3 - "$path" "$expected_actor" "$mode" "$client_secret_expected" <<'PY'
import json
import sys

path = sys.argv[1]
expected_actor = sys.argv[2]
mode = sys.argv[3]
client_secret_expected = sys.argv[4] == "yes"

with open(path, "r", encoding="utf-8") as input_file:
    payload = json.load(input_file)


def fail(message):
    sys.stderr.write(message + "\n")
    sys.exit(1)


prefix = "browser login status " if mode == "status" else "browser login "

if payload.get("actor") != expected_actor:
    fail(prefix + "actor mismatch")
if payload.get("target", {}).get("status") != "ready":
    fail(prefix + "target is not ready")
if payload.get("token", {}).get("status") != "set":
    token_message = "token is not set" if mode == "status" else "token was not saved"
    fail(prefix + token_message)

app = payload.get("app", {})
if app.get("client_id") != "set":
    fail(prefix + "client id was not redacted")
expected_secret_status = "set" if client_secret_expected else "missing"
if app.get("client_secret") != expected_secret_status:
    fail(prefix + "client secret status mismatch")
PY
}

assert_not_printed() {
  local needle="$1"
  shift
  if [[ -z "$needle" ]]; then
    return 0
  fi
  local path
  for path in "$@"; do
    if grep -Fq "$needle" "$path"; then
      printf 'browser login smoke printed secret material in %s\n' "$path" >&2
      exit 1
    fi
  done
}

wait_for_authorize_url() {
  local path="$1"
  local pid="$2"
  local authorize_url=""
  for _ in {1..200}; do
    authorize_url="$(grep -m1 '^https://linear.app/oauth/authorize' "$path" || true)"
    if [[ -n "$authorize_url" ]]; then
      printf '%s\n' "$authorize_url"
      return 0
    fi
    if ! kill -0 "$pid" 2>/dev/null; then
      break
    fi
    sleep 0.1
  done

  return 1
}

start_callback_listener() {
  local redirect_uri="$1"
  local callback_path="$2"
  local ready_path="$3"
  local error_path="$4"
  local timeout_seconds="$5"

  python3 - "$redirect_uri" "$callback_path" "$ready_path" "$error_path" "$timeout_seconds" <<'PY' &
import html
import http.server
import os
import socket
import sys
import time
import urllib.parse

redirect_uri = sys.argv[1]
callback_path = sys.argv[2]
ready_path = sys.argv[3]
error_path = sys.argv[4]
timeout_seconds = int(sys.argv[5])


def fail(message):
    with open(error_path, "w", encoding="utf-8") as output:
        output.write(message + "\n")
    sys.exit(1)


parsed_redirect = urllib.parse.urlparse(redirect_uri)
host = parsed_redirect.hostname
port = parsed_redirect.port
path = parsed_redirect.path or "/"
if (
    parsed_redirect.scheme != "http"
    or host not in {"127.0.0.1", "localhost", "::1"}
    or port is None
):
    fail("redirect URI is not a supported localhost HTTP callback")


class CallbackHandler(http.server.BaseHTTPRequestHandler):
    server_version = "linctl-browser-login-smoke"

    def log_message(self, _format, *_args):
        return

    def do_GET(self):
        request = urllib.parse.urlparse(self.path)
        if request.path != path:
            self.send_response(404)
            self.send_header("Content-Type", "text/plain; charset=utf-8")
            self.end_headers()
            self.wfile.write(b"linctl browser login callback is waiting on a different path.\n")
            return

        query = urllib.parse.parse_qs(request.query)
        if "code" not in query or "state" not in query:
            self.send_response(400)
            self.send_header("Content-Type", "text/plain; charset=utf-8")
            self.end_headers()
            self.wfile.write(b"Linear did not send an OAuth code and state.\n")
            return

        callback_url = urllib.parse.urlunparse(
            (
                parsed_redirect.scheme,
                parsed_redirect.netloc,
                request.path,
                "",
                request.query,
                "",
            )
        )
        with open(callback_path, "w", encoding="utf-8") as output:
            output.write(callback_url + "\n")

        body = """<!doctype html>
<html lang="en">
<meta charset="utf-8">
<title>linctl browser login complete</title>
<style>
body {
  color: #1f2937;
  font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  margin: 4rem auto;
  max-width: 42rem;
  line-height: 1.5;
}
.box {
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 1.25rem 1.5rem;
}
h1 {
  font-size: 1.35rem;
  margin: 0 0 .5rem;
}
p {
  margin: .35rem 0 0;
}
</style>
<main class="box">
  <h1>linctl browser login complete</h1>
  <p>You can close this tab and return to the terminal.</p>
  <p>The callback code was captured locally and was not printed.</p>
</main>
</html>
"""
        encoded = body.encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "text/html; charset=utf-8")
        self.send_header("Content-Length", str(len(encoded)))
        self.end_headers()
        self.wfile.write(encoded)


address_family = socket.AF_INET6 if host == "::1" else socket.AF_INET


class CallbackServer(http.server.HTTPServer):
    allow_reuse_address = True


CallbackServer.address_family = address_family
try:
    server = CallbackServer((host, port), CallbackHandler)
except OSError as error:
    fail(f"could not listen on {host}:{port}: {error}")

with open(ready_path, "w", encoding="utf-8") as output:
    output.write("ready\n")

deadline = time.monotonic() + timeout_seconds
server.timeout = 0.5
try:
    while not os.path.exists(callback_path):
        if time.monotonic() > deadline:
            fail("timed out waiting for Linear OAuth callback")
        server.handle_request()
finally:
    server.server_close()
PY
}

wait_for_callback_listener() {
  local ready_path="$1"
  local error_path="$2"
  local pid="$3"

  for _ in {1..100}; do
    if [[ -s "$ready_path" ]]; then
      return 0
    fi
    if [[ -s "$error_path" ]]; then
      cat "$error_path" >&2
      return 1
    fi
    if ! kill -0 "$pid" 2>/dev/null; then
      break
    fi
    sleep 0.1
  done

  if [[ -s "$error_path" ]]; then
    cat "$error_path" >&2
  else
    printf 'callback listener did not start\n' >&2
  fi
  return 1
}

wait_for_callback_file() {
  local callback_path="$1"
  local error_path="$2"
  local timeout_seconds="$3"
  local attempts=$((timeout_seconds * 10))

  for ((attempt = 0; attempt < attempts; attempt++)); do
    if [[ -s "$callback_path" ]]; then
      sed -n '1p' "$callback_path"
      return 0
    fi
    if [[ -s "$error_path" ]]; then
      cat "$error_path" >&2
      return 1
    fi
    sleep 0.1
  done

  printf 'timed out waiting for Linear OAuth callback\n' >&2
  return 1
}

read_manual_callback() {
  if [[ ! -t 0 ]]; then
    printf 'manual callback mode requires an interactive terminal\n' >&2
    return 2
  fi

  printf 'After authorization, your browser may show "This site can'\''t be reached".\n' >&2
  printf 'That is expected in manual mode. Copy the full callback URL from the address bar.\n' >&2
  printf 'Paste callback URL here (hidden because it contains a one-time OAuth code): ' >&2
  local callback
  IFS= read -r -s callback
  printf '\n' >&2
  if [[ -z "$callback" ]]; then
    printf 'missing OAuth callback URL\n' >&2
    return 2
  fi

  printf '%s\n' "$callback"
}

(
  export XDG_CONFIG_HOME="$smoke_dir/config"
  export XDG_STATE_HOME="$smoke_dir/state"

  callback_listener_pid=""
  # shellcheck disable=SC2329 # Invoked by the EXIT trap below.
  cleanup_listener() {
    if [[ -n "$callback_listener_pid" ]]; then
      kill "$callback_listener_pid" 2>/dev/null || true
      wait "$callback_listener_pid" 2>/dev/null || true
    fi
  }
  trap cleanup_listener EXIT

  cd "$smoke_dir"

  configure_args=(
    "$binary"
    auth configure
    --client-id "$LINCTL_OAUTH_CLIENT_ID"
    --redirect-uri "$LINCTL_OAUTH_REDIRECT_URI"
    --scopes "$scopes"
    --quiet
  )
  if [[ -n "${LINCTL_OAUTH_CLIENT_SECRET:-}" ]]; then
    configure_args+=(--client-secret "$LINCTL_OAUTH_CLIENT_SECRET")
  fi
  "${configure_args[@]}"

  login_output="$smoke_dir/login.json"
  login_error="$smoke_dir/login.err"
  status_output="$smoke_dir/status.json"
  status_error="$smoke_dir/status.err"
  callback_file="$smoke_dir/callback.url"
  listener_ready="$smoke_dir/callback.ready"
  listener_error="$smoke_dir/callback.err"
  callback_mode="${LINCTL_BROWSER_LOGIN_CALLBACK_MODE:-auto}"
  case "$callback_mode" in
    auto | manual) ;;
    *)
      printf 'LINCTL_BROWSER_LOGIN_CALLBACK_MODE must be auto or manual\n' >&2
      exit 2
      ;;
  esac

  if [[ "$callback_mode" == "auto" ]]; then
    start_callback_listener \
      "$LINCTL_OAUTH_REDIRECT_URI" \
      "$callback_file" \
      "$listener_ready" \
      "$listener_error" \
      "$callback_timeout"
    callback_listener_pid="$!"
    if ! wait_for_callback_listener "$listener_ready" "$listener_error" "$callback_listener_pid"; then
      exit 1
    fi
  fi

  coproc AUTH_LOGIN {
    "$binary" --json auth login --actor "$actor" --callback - >"$login_output" 2>"$login_error"
  }
  auth_login_stdin="${AUTH_LOGIN[1]}"

  if ! authorize_url="$(wait_for_authorize_url "$login_error" "$AUTH_LOGIN_PID")"; then
    wait "$AUTH_LOGIN_PID" || true
    cat "$login_error" >&2
    exit 1
  fi

  printf 'Open this Linear OAuth URL:\n%s\n\n' "$authorize_url" >&2
  if [[ "$callback_mode" == "auto" ]]; then
    printf 'Waiting for Linear to redirect back to %s ...\n' "$LINCTL_OAUTH_REDIRECT_URI" >&2
    printf 'The browser will show a linctl success page. No callback copy/paste is needed.\n\n' >&2
    if ! callback="$(wait_for_callback_file "$callback_file" "$listener_error" "$callback_timeout")"; then
      kill "$AUTH_LOGIN_PID" 2>/dev/null || true
      wait "$AUTH_LOGIN_PID" || true
      exit 1
    fi
  else
    if ! callback="$(read_manual_callback)"; then
      kill "$AUTH_LOGIN_PID" 2>/dev/null || true
      wait "$AUTH_LOGIN_PID" || true
      exit 2
    fi
  fi

  printf '%s\n' "$callback" >&"$auth_login_stdin"
  exec {auth_login_stdin}>&-

  if ! wait "$AUTH_LOGIN_PID"; then
    cat "$login_error" >&2
    exit 1
  fi

  client_secret_expected="no"
  if [[ -n "${LINCTL_OAUTH_CLIENT_SECRET:-}" ]]; then
    client_secret_expected="yes"
  fi
  validate_auth_json "$login_output" "$actor" login "$client_secret_expected"
  assert_not_printed "$callback" "$login_output" "$login_error"
  assert_not_printed "${LINCTL_OAUTH_CLIENT_SECRET:-}" "$login_output" "$login_error"

  "$binary" auth status --json >"$status_output" 2>"$status_error"
  validate_auth_json "$status_output" "$actor" status "$client_secret_expected"
  assert_not_printed "${LINCTL_OAUTH_CLIENT_SECRET:-}" "$status_output" "$status_error"

  "$binary" auth logout --forget-app --quiet

  printf 'browser login smoke ok: actor=%s target=ready token=set scopes=set\n' "$actor"
)
