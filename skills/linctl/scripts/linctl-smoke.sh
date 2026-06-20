#!/usr/bin/env bash
# Portable, read-only linctl smoke. Runs a few read commands with --json and reports
# pass/fail per command without printing tokens or full payloads.
#
# Complements the in-repo `task live-smoke` (full harness with disposable writes); use this
# when you only have an installed binary, or want a quick read-only confidence check outside
# the checkout. Read-only by design: it never creates, updates, or archives Linear resources.
#
# Credentials: the CLI reads LINCTL_TOKEN, then LINEAR_API_KEY, then a config token. This
# script never sets, echoes, or logs any token value. On a missing token or target mismatch
# the underlying command fails closed and this script reports FAIL with the first error line.
set -uo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
prefix="$(bash "$here/linctl-resolve.sh")" || exit 1

status=0
run() {
  local label="$1"; shift
  local out
  # shellcheck disable=SC2086 # $prefix may be "go run ./cmd/linctl" and must word-split
  if out="$($prefix "$@" --json 2>&1)"; then
    printf 'PASS  %s\n' "$label"
  else
    printf 'FAIL  %s: %s\n' "$label" "$(printf '%s' "$out" | head -n 1)"
    status=1
  fi
}

run "target"       target
run "whoami"       whoami
run "issue list"   issue list --limit 5
run "project list" project list --limit 5

if [ "$status" -ne 0 ]; then
  printf '\nread-only smoke failed; check token and pinned target (no values printed)\n' >&2
fi
exit "$status"
