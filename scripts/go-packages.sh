#!/usr/bin/env bash
set -euo pipefail

module_root="$(pwd -P)"

go list -f '{{.Dir}}' ./... 2>/dev/null |
  grep -v '/skills/' |
  sed "s#^${module_root}#.#"
