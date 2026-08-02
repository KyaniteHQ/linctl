#!/usr/bin/env bash
# shellcheck shell=bash
# Refresh the vendored Linear schema from the upstream linear/linear SDK schema,
# bump the .linear-sdk-ref pin, regenerate the genqlient client, and rewrite the
# coverage ledger. No OAuth token required: the source of truth is the same SDK
# schema file that schema-drift-check compares against.
#
# Env:
#   LINCTL_LINEAR_SDK_REF       upstream git ref (default: master)
#   LINCTL_LINEAR_SDK_UPSTREAM  optional existing Linear checkout path
#   LINCTL_LINEAR_SDK_OFFLINE   if set, reuse LINCTL_LINEAR_SDK_UPSTREAM without fetch
#
# On success prints:
#   UPSTREAM_SHA=<full sha>
#   REFRESHED=1|0   (1 when any tracked file changed)
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

ref="${LINCTL_LINEAR_SDK_REF:-master}"
upstream="${LINCTL_LINEAR_SDK_UPSTREAM:-}"
upstream_root=""
cleanup() {
  if [[ -n "$upstream_root" ]]; then
    rm -rf "$upstream_root"
  fi
}
trap cleanup EXIT

if [[ -z "$upstream" ]]; then
  upstream_root="$(mktemp -d)"
  upstream="$upstream_root/linear"
fi

LINCTL_LINEAR_SDK_UPSTREAM="$upstream" LINCTL_LINEAR_SDK_REF="$ref" \
  go tool task linear-sdk-upstream-checkout

upstream_schema="$upstream/packages/sdk/src/schema.graphql"
if [[ ! -f "$upstream_schema" ]]; then
  printf 'upstream schema missing at %s\n' "$upstream_schema" >&2
  exit 1
fi

sha="$(git -C "$upstream" rev-parse HEAD)"
if [[ ! "$sha" =~ ^[0-9a-f]{40}$ ]]; then
  printf 'unexpected upstream SHA: %s\n' "$sha" >&2
  exit 1
fi

cp "$upstream_schema" "$repo_root/internal/client/schema.graphql"

printf '%s\n' "$sha" > "$repo_root/.linear-sdk-ref"
printf 'pinned the Linear SDK ref to %s\n' "$sha" >&2

go generate ./...

LINCTL_LINEAR_SDK_UPSTREAM="$upstream" LINCTL_LINEAR_SDK_OFFLINE=1 LINCTL_LINEAR_SDK_REF="$sha" \
  go run scripts/linear_api_coverage*.go --upstream "$upstream" --output docs/internal/linear-api-coverage.md

# Confirm the vendored schema now matches the upstream we just synced from.
LINCTL_LINEAR_SDK_UPSTREAM="$upstream" LINCTL_LINEAR_SDK_OFFLINE=1 LINCTL_LINEAR_SDK_REF="$sha" \
  go tool task schema-drift-check >/dev/null

changed=0
if ! git diff --quiet -- \
  internal/client/schema.graphql \
  internal/client/internal/gql/generated.go \
  docs/internal/linear-api-coverage.md \
  .linear-sdk-ref; then
  changed=1
fi

printf 'UPSTREAM_SHA=%s\n' "$sha"
printf 'REFRESHED=%s\n' "$changed"
