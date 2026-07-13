#!/usr/bin/env bash
# Render the README terminal demo (docs/assets/demo.{gif,mp4}) from the fixture in
# demo/fixture, so no Linear organization and no credentials are needed.
#
# The fixture replays the real binary's output shapes against invented ids. To
# record the actual binary against a throwaway Linear organization instead, use
# ./render.sh.
#
# Requires: vhs (https://github.com/charmbracelet/vhs)
set -euo pipefail

cd "$(dirname "$0")"

command -v vhs >/dev/null || {
  echo "vhs not found on PATH" >&2
  exit 1
}

work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

install -m 0755 fixture/linctl "$work/linctl"

cat >"$work/.linctl.toml" <<'TOML'
[target]
org_id     = "c5bc56eb-7e79-4d90-969b-65d938103e0d"
team_key   = "LIT"
team_id    = "9c870067-bd25-412e-9e78-62c56a6d2716"
project_id = "151e4a94-678f-4cf6-9fce-ccdc803fe7b1"
TOML

# vhs resolves Output paths relative to its own working directory, so render into
# the temp dir and copy the results back.
sed \
  -e 's#^Output \.\./docs/assets/#Output #' \
  demo.tape >"$work/demo.tape"

(cd "$work" && PATH="$work:$PATH" vhs demo.tape)

cp "$work/demo.gif" ../docs/assets/demo.gif
cp "$work/demo.mp4" ../docs/assets/demo.mp4

echo "rendered docs/assets/demo.gif and docs/assets/demo.mp4"
