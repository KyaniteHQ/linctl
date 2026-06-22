#!/usr/bin/env bash
# Token-free linctl smoke. Proves the binary builds and runs in a headless
# environment with NO Linear credentials and NO network access.
#
# Complements linctl-smoke.sh (read-only, needs a token) and `task live-smoke`
# (full harness with disposable writes). Use this first on a clean machine to
# confirm linctl is wired up before any credential or target work.
#
# Every command here avoids the command runtime, so none of them resolves a
# token or contacts Linear: --version, --help, usage guidance, and completion.
# A non-zero exit from any one fails the smoke.
set -uo pipefail

here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
prefix="$(bash "$here/linctl-resolve.sh")" || exit 1

status=0
run() {
  local label="$1"; shift
  # shellcheck disable=SC2086 # $prefix may be "go run ./cmd/linctl" and must word-split
  if $prefix "$@" >/dev/null 2>&1; then
    printf 'PASS  %s\n' "$label"
  else
    printf 'FAIL  %s (exit %d)\n' "$label" "$?"
    status=1
  fi
}

run "--version"        --version
run "--help"           --help
run "usage"            usage
run "usage overview"   usage overview
run "issue --help"     issue --help
run "project --help"   project --help
run "cycle --help"     cycle --help
run "completion bash"  completion bash

if [ "$status" -ne 0 ]; then
  printf '\noffline smoke failed; the binary is not runnable in this environment\n' >&2
fi
exit "$status"
