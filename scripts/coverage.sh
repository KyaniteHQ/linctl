#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

profile="${1:-/tmp/linctl.cover}"
filtered="${profile%.cover}.handwritten.cover"

# Portable replacement for 'mapfile' (absent on the macOS system bash 3.2), so
# 'task coverage' enforces the gate on every platform instead of silently
# passing zero packages to go test.
package_output="$(bash scripts/go-packages.sh)"
packages=()
while IFS= read -r package_dir; do
  if [[ -n "$package_dir" ]]; then
    packages+=("$package_dir")
  fi
done <<< "$package_output"
if (( ${#packages[@]} == 0 )); then
  printf 'coverage package discovery returned no packages\n' >&2
  exit 1
fi

# -race here doubles as a second race pass alongside the test job at no extra
# wall-clock cost (CI jobs run in parallel). Generated transport and export
# wrappers plus the thin main entrypoint are excluded from the hand-written
# coverage gate; scripts/ is already excluded by go-packages.sh.
go test -race -count=1 -coverprofile="$profile" "${packages[@]}"
# Any nonzero filter status is fatal. A valid profile always keeps its 'mode:'
# header through both filters, so grep cannot legitimately report "nothing
# selected"; tolerating status 1 would let a failed redirect leave a stale
# $filtered behind and have the gate validate the previous run's profile.
filter_status=0
grep -Fv '/internal/client/internal/gql/generated.go:' "$profile" |
  grep -Fv '/internal/client/internal/gql/exports.go:' |
  grep -Fv '/cmd/linctl/main.go:' > "$filtered" || filter_status=$?
if (( filter_status != 0 )); then
  printf 'failed to filter coverage profile: status %d\n' "$filter_status" >&2
  exit "$filter_status"
fi

# Field 2 is the block's statement count. Go emits zero-statement blocks, and a
# profile made only of those measures nothing while still satisfying every check
# below, so require at least one block that actually carries statements.
if ! awk '$1 != "mode:" && $2 > 0 { found = 1 } END { exit !found }' "$filtered"; then
  printf 'hand-written coverage profile has no statement blocks\n' >&2
  exit 1
fi

coverage_output="$(go tool cover -func="$filtered")"
printf '%s\n' "$coverage_output"

zero_count_blocks="$(awk '$1 != "mode:" && $2 > 0 && $3 == 0 {
  location = $1
  sub(/,.*/, "", location)
  sub(/\.[0-9]+$/, "", location)
  print location
}' "$filtered")"
if [[ -n "$zero_count_blocks" ]]; then
  printf 'hand-written coverage has zero-count statement blocks:\n%s\n' "$zero_count_blocks" >&2
  exit 1
fi
