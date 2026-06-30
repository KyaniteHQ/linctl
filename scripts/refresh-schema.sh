#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
schema_path="$repo_root/internal/client/schema.graphql"
token="${LINCTL_OAUTH_ACCESS_TOKEN:-}"

if [[ -z "$token" ]]; then
  printf 'missing Linear OAuth access token: set LINCTL_OAUTH_ACCESS_TOKEN\n' >&2
  exit 1
fi

if ! command -v npx >/dev/null 2>&1; then
  printf 'npx (Node.js) is required to introspect the Linear schema\n' >&2
  exit 1
fi

mkdir -p "$(dirname "$schema_path")"
tmp_schema="$(mktemp)"
trap 'rm -f "$tmp_schema"' EXIT

npx --yes --package graphql@16.14.2 node >"$tmp_schema" <<'NODE'
const { buildClientSchema, getIntrospectionQuery, printSchema } = require("graphql");

async function main() {
  const token = process.env.LINCTL_OAUTH_ACCESS_TOKEN;
  const response = await fetch("https://api.linear.app/graphql", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ query: getIntrospectionQuery() }),
  });
  const responseBody = await response.text();
  if (!response.ok) {
    throw new Error(`Linear schema request failed (${response.status}): ${responseBody}`);
  }
  const payload = JSON.parse(responseBody);
  if (payload.errors && payload.errors.length > 0) {
    throw new Error(`Linear schema introspection failed: ${JSON.stringify(payload.errors)}`);
  }
  process.stdout.write(printSchema(buildClientSchema(payload.data)));
  process.stdout.write("\n");
}

main().catch((error) => {
  console.error(error.message);
  process.exit(1);
});
NODE

mv "$tmp_schema" "$schema_path"
