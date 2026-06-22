#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

profile="${1:-/tmp/linctl.cover}"
filtered="${profile%.cover}.handwritten.cover"

# Portable replacement for 'mapfile' (absent on the macOS system bash 3.2), so
# 'task coverage' enforces the gate on every platform instead of silently
# passing zero packages to go test.
packages=()
while IFS= read -r package_dir; do
  packages+=("$package_dir")
done < <(bash scripts/go-packages.sh)

# -race here doubles as a second race pass alongside the test job at no extra
# wall-clock cost (CI jobs run in parallel). The generated client and the thin
# main entrypoint are excluded from the hand-written coverage gate; scripts/ is
# already excluded by go-packages.sh.
go test -race -count=1 -coverprofile="$profile" "${packages[@]}"
grep -v '/internal/client/generated.go:' "$profile" |
  grep -v '/cmd/linctl/main.go:' > "$filtered"

coverage_output="$(go tool cover -func="$filtered")"
printf '%s\n' "$coverage_output"

total="$(printf '%s\n' "$coverage_output" | awk '/^total:/ {print $3}')"
if [[ "$total" != "100.0%" ]]; then
  printf 'hand-written coverage must be 100.0%%, got %s\n' "$total" >&2
  exit 1
fi
