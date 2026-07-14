#!/usr/bin/env bash
# shellcheck shell=bash
set -euo pipefail

module_root="$(pwd -P)"

# Collect the unique directories that hold Go files. This avoids the bash 4
# builtin 'mapfile' and the GNU-only 'xargs -r', so the helper also runs on the
# macOS system bash (3.2) and BSD userland, not just Linux.
#
# The find|sort pipeline runs in a command substitution, not a process
# substitution, so 'set -euo pipefail' sees its exit status. A process
# substitution lets a partial list survive a failed producer, and every gate
# fed by this helper -- test, vet, build, lint, govulncheck, coverage -- would
# then pass while silently skipping the packages that were never discovered.
#
# scripts/ is intentionally excluded: its only Go file is a standalone
# maintenance tool tagged '//go:build ignore', so 'go list ./scripts' would
# fail with "build constraints exclude all Go files".
package_dir_output="$(
  find ./cmd ./internal -type f -name '*.go' -exec dirname {} \; |
    sort -u
)"

package_dirs=()
while IFS= read -r dir; do
  if [[ -n "$dir" ]]; then
    package_dirs+=("$dir")
  fi
done <<< "$package_dir_output"

# Discovering nothing is a failure, not an empty success: every caller word-
# splits this list into a gate command, so returning zero packages with status 0
# would hand each gate an empty scope to declare green.
if ((${#package_dirs[@]} == 0)); then
  printf 'go package discovery found no Go directories under ./cmd or ./internal\n' >&2
  exit 1
fi

go list -f '{{.Dir}}' "${package_dirs[@]}" |
  sed "s#^${module_root}#.#"
