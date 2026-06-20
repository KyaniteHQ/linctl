#!/usr/bin/env bash
# Resolve the linctl command prefix for this environment and print it to stdout.
# Order: an installed `linctl` binary, else `go run ./cmd/linctl` inside a linctl checkout.
# Exits non-zero with a clear message when neither is available.
#
# Usage (note: $prefix is left unquoted so `go run ./cmd/linctl` splits into words):
#   prefix="$(bash skills/linctl/scripts/linctl-resolve.sh)" || exit 1
#   $prefix target --json
set -euo pipefail

if command -v linctl >/dev/null 2>&1; then
  printf 'linctl\n'
  exit 0
fi

if [ -f cmd/linctl/main.go ]; then
  printf 'go run ./cmd/linctl\n'
  exit 0
fi

printf 'linctl is unavailable: no installed binary and no cmd/linctl/main.go in this checkout\n' >&2
exit 1
