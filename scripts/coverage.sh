#!/usr/bin/env bash
set -euo pipefail

profile="${1:-/tmp/linctl.cover}"
filtered="${profile%.cover}.handwritten.cover"
mapfile -t packages < <(bash scripts/go-packages.sh)

go test -count=1 -coverprofile="$profile" "${packages[@]}"
grep -v '/internal/client/generated.go:' "$profile" |
  grep -v '/cmd/linctl/main.go:' |
  grep -v '/scripts/' > "$filtered"

coverage_output="$(go tool cover -func="$filtered")"
printf '%s\n' "$coverage_output"

total="$(printf '%s\n' "$coverage_output" | awk '/^total:/ {print $3}')"
if [[ "$total" != "100.0%" ]]; then
  printf 'hand-written coverage must be 100.0%%, got %s\n' "$total" >&2
  exit 1
fi
