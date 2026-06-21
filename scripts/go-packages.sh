#!/usr/bin/env bash
set -euo pipefail

module_root="$(pwd -P)"

mapfile -t package_dirs < <(
  find ./cmd ./internal ./scripts -type f -name '*.go' -print |
    xargs -r -n1 dirname |
    sort -u
)

if ((${#package_dirs[@]} == 0)); then
  exit 0
fi

go list -f '{{.Dir}}' "${package_dirs[@]}" |
  sed "s#^${module_root}#.#"
