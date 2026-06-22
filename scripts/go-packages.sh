#!/usr/bin/env bash
set -euo pipefail

module_root="$(pwd -P)"

# Collect the unique directories that hold Go files. This avoids the bash 4
# builtin 'mapfile' and the GNU-only 'xargs -r', so the helper also runs on the
# macOS system bash (3.2) and BSD userland, not just Linux.
package_dirs=()
while IFS= read -r dir; do
  package_dirs+=("$dir")
done < <(
  find ./cmd ./internal ./scripts -type f -name '*.go' -exec dirname {} \; |
    sort -u
)

if ((${#package_dirs[@]} == 0)); then
  exit 0
fi

go list -f '{{.Dir}}' "${package_dirs[@]}" |
  sed "s#^${module_root}#.#"
