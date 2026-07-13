#!/usr/bin/env bash
# shellcheck shell=bash
# Requires curl and uv. uv supplies tiktoken without changing project manifests.
set -euo pipefail

endpoint="https://mcp.linear.app/mcp"
protocol_version="2025-06-18"
token="${LINCTL_OAUTH_ACCESS_TOKEN:-}"

if [[ -z "$token" ]]; then
  printf 'missing Linear OAuth access token: set LINCTL_OAUTH_ACCESS_TOKEN\n' >&2
  exit 1
fi

if ! command -v curl >/dev/null 2>&1; then
  printf 'curl is required to query the Linear MCP server\n' >&2
  exit 1
fi

if ! command -v uv >/dev/null 2>&1; then
  printf 'uv is required to count MCP tool-definition tokens\n' >&2
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

short_body_summary() {
  local body_file="$1"

  uv run --no-project python /dev/fd/3 "$body_file" 3<<'PY'
import os
import re
import sys
from pathlib import Path

body = Path(sys.argv[1]).read_text(encoding="utf-8", errors="replace")
token = os.environ.get("LINCTL_OAUTH_ACCESS_TOKEN", "")
if token:
    body = body.replace(token, "[REDACTED]")
body = re.sub(r"(?i)bearer\s+\S+", "Bearer [REDACTED]", body)
body = "".join(character if character.isprintable() else " " for character in body)
body = " ".join(body.split())
print((body[:240] if body else "(empty body)"))
PY
}

mcp_post() {
  local label="$1"
  local payload="$2"
  local response_file="$3"
  local headers_file="$4"
  local session_id="${5:-}"
  local include_protocol_version="${6:-}"
  local http_status
  local -a headers=(
    --header "Authorization: Bearer $token"
    --header "Accept: application/json, text/event-stream"
    --header "Content-Type: application/json"
  )

  if [[ -n "$session_id" ]]; then
    headers+=(--header "Mcp-Session-Id: $session_id")
  fi
  if [[ -n "$include_protocol_version" ]]; then
    headers+=(--header "MCP-Protocol-Version: $protocol_version")
  fi

  if ! http_status="$(curl \
    "${headers[@]}" \
    --silent \
    --show-error \
    --request POST \
    --url "$endpoint" \
    --data-binary "$payload" \
    --dump-header "$headers_file" \
    --output "$response_file" \
    --write-out '%{http_code}')"; then
    printf 'Linear MCP %s request failed before receiving an HTTP response\n' "$label" >&2
    return 1
  fi

  if [[ ! "$http_status" =~ ^2[0-9][0-9]$ ]]; then
    printf 'Linear MCP %s request failed with HTTP %s: %s\n' \
      "$label" "$http_status" "$(short_body_summary "$response_file")" >&2
    return 1
  fi
}

initialize_payload="$(printf \
  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"%s","capabilities":{},"clientInfo":{"name":"linctl-mcp-token-measure","version":"1.0"}}}' \
  "$protocol_version")"
initialized_payload='{"jsonrpc":"2.0","method":"notifications/initialized"}'
tools_list_payload='{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'

initialize_response="$tmp_dir/initialize-response"
initialize_headers="$tmp_dir/initialize-headers"
mcp_post "initialize" "$initialize_payload" "$initialize_response" "$initialize_headers"

session_id="$(awk '
  tolower($1) == "mcp-session-id:" {
    sub(/\r$/, "", $2)
    print $2
    exit
  }
' "$initialize_headers")"

initialized_response="$tmp_dir/initialized-response"
initialized_headers="$tmp_dir/initialized-headers"
mcp_post \
  "initialized notification" \
  "$initialized_payload" \
  "$initialized_response" \
  "$initialized_headers" \
  "$session_id" \
  "include"

tools_response="$tmp_dir/tools-response"
tools_headers="$tmp_dir/tools-headers"
mcp_post \
  "tools/list" \
  "$tools_list_payload" \
  "$tools_response" \
  "$tools_headers" \
  "$session_id" \
  "include"

uv run --no-project --with tiktoken python /dev/fd/3 "$initialize_response" 3<<'PY' <"$tools_response"
import json
import sys
from pathlib import Path

import tiktoken


def parse_messages(raw):
    try:
        payload = json.loads(raw)
    except json.JSONDecodeError:
        payloads = []
        data_lines = []
        for line in raw.splitlines() + [""]:
            if line.startswith("data:"):
                data_lines.append(line[5:].lstrip())
            elif not line and data_lines:
                data = "\n".join(data_lines)
                data_lines = []
                if data != "[DONE]":
                    payloads.append(json.loads(data))
        if not payloads:
            raise ValueError("response was neither JSON nor a populated SSE stream")
    else:
        payloads = payload if isinstance(payload, list) else [payload]
    return payloads


def response_for_id(raw, response_id, label):
    try:
        messages = parse_messages(raw)
    except (json.JSONDecodeError, ValueError) as error:
        raise ValueError(f"Linear MCP {label} response could not be parsed: {error}") from error

    for message in messages:
        if isinstance(message, dict) and message.get("id") == response_id:
            if "error" in message:
                code = message["error"].get("code", "unknown")
                raise ValueError(f"Linear MCP {label} returned JSON-RPC error {code}")
            return message.get("result")
    raise ValueError(f"Linear MCP {label} response did not contain JSON-RPC id {response_id}")


try:
    initialize_raw = Path(sys.argv[1]).read_text(encoding="utf-8", errors="replace")
    initialize_result = response_for_id(initialize_raw, 1, "initialize")
    if not isinstance(initialize_result, dict):
        raise ValueError("Linear MCP initialize response did not contain a result object")

    tools_result = response_for_id(sys.stdin.read(), 2, "tools/list")
    if not isinstance(tools_result, dict) or not isinstance(tools_result.get("tools"), list):
        raise ValueError("Linear MCP tools/list response did not contain a tools array")

    tools = tools_result["tools"]
    compact = json.dumps(tools, ensure_ascii=False, separators=(",", ":"))
    pretty = json.dumps(tools, ensure_ascii=False, indent=2)
    encoding = tiktoken.get_encoding("o200k_base")

    print(f"tool count: {len(tools)}")
    print(f"compact tokens (o200k_base): {len(encoding.encode(compact))}")
    print(f"pretty tokens (o200k_base): {len(encoding.encode(pretty))}")
except ValueError as error:
    print(error, file=sys.stderr)
    raise SystemExit(1) from error
PY
