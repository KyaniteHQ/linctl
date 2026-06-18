#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
schema_path="$repo_root/internal/client/schema.graphql"
token="${LINCTL_TEST_TOKEN:-${LINCTL_TOKEN:-${LINEAR_API_KEY:-}}}"

if [[ -z "$token" ]]; then
  printf 'missing Linear API token: set LINCTL_TEST_TOKEN, LINCTL_TOKEN, or LINEAR_API_KEY\n' >&2
  exit 1
fi

mkdir -p "$(dirname "$schema_path")"
tmp_schema="$(mktemp)"
trap 'rm -f "$tmp_schema"' EXIT

npx --yes get-graphql-schema@2.1.2 \
  --header "Authorization=$token" \
  https://api.linear.app/graphql >"$tmp_schema"

mv "$tmp_schema" "$schema_path"
