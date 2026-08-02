#!/usr/bin/env bash
# shellcheck shell=bash
# Run the Linear upstream-schema checks behind one Linear SDK checkout.
#
# Usage: scripts/upstream-check.sh operations|ledger|drift|all
#
#   operations  validate local GraphQL operations against the upstream schema
#   ledger      regenerate the coverage ledger and diff it against the committed one
#   drift       report semantic drift between the vendored and upstream schemas
#   all         operations, then ledger
#
# Env:
#   LINCTL_LINEAR_SDK_UPSTREAM  reuse an existing Linear checkout at this path.
#                               Unset creates a run-local checkout and removes it on exit.
set -euo pipefail

check="${1:-all}"
case "$check" in
  operations | ledger | drift | all) ;;
  *)
    printf 'usage: %s operations|ledger|drift|all\n' "$0" >&2
    exit 2
    ;;
esac

upstream="${LINCTL_LINEAR_SDK_UPSTREAM:-}"
upstream_root=""
ledger_tmp=""
cleanup() {
  if [[ -n "$ledger_tmp" ]]; then
    rm -f "$ledger_tmp"
  fi
  if [[ -n "$upstream_root" ]]; then
    rm -rf "$upstream_root"
  fi
}
trap cleanup EXIT

if [[ -z "$upstream" ]]; then
  upstream_root="$(mktemp -d)"
  upstream="$upstream_root/linear"
fi
LINCTL_LINEAR_SDK_UPSTREAM="$upstream" go tool task linear-sdk-upstream-checkout

run_operations() {
  go run scripts/linear_graphql_operation_check.go --upstream "$upstream"
}

run_ledger() {
  ledger_tmp="$(mktemp)"
  go run scripts/linear_api_coverage*.go --upstream "$upstream" --output "$ledger_tmp"
  diff -u docs/internal/linear-api-coverage.md "$ledger_tmp"
}

run_drift() {
  go run scripts/linear_schema_drift_check.go --upstream "$upstream"
}

case "$check" in
  operations) run_operations ;;
  ledger) run_ledger ;;
  drift) run_drift ;;
  all)
    run_operations
    run_ledger
    ;;
esac
