#!/usr/bin/env bash
# Render the linctl terminal demo (docs/assets/demo.gif + demo.mp4).
#
# The tape runs the real linctl binary against a Linear workspace, so you need:
#   - vhs (https://github.com/charmbracelet/vhs) on PATH
#   - a Linear personal API key for a THROWAWAY demo workspace in LINEAR_API_KEY
#   - a .linctl.toml in this directory pinning that workspace's org/team
#
# The guarded write lands in the pinned target; the wrong-team write is refused.
# Usage: LINEAR_API_KEY=<token> ./render.sh
set -euo pipefail
cd "$(dirname "$0")"

: "${LINEAR_API_KEY:?set LINEAR_API_KEY to a Linear token for the demo workspace}"
command -v vhs >/dev/null || { echo "vhs not found on PATH" >&2; exit 1; }
[ -f .linctl.toml ] || { echo ".linctl.toml missing — pin a demo org/team here" >&2; exit 1; }

# Build a fresh binary into this dir so the tape's bare `linctl` resolves to it.
go build -o ./linctl ../cmd/linctl
export PATH="$PWD:$PATH"

exec vhs demo.tape
