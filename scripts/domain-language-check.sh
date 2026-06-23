#!/usr/bin/env bash
set -euo pipefail

if rg -n -i '\bworkspace(-level)?\b' \
  README.md \
  CONTRIBUTING.md \
  internal/cli \
  docs \
  skills/linctl/references/commands.md \
  -g '!docs/domain-map.md' \
  -g '!docs/linear-api-coverage.md' \
  -g '!docs/linear-cli-feature-leech.md' \
  -g '!docs/setup-gap-log.md'; then
  cat <<'MSG'
Avoid "workspace" in user-facing linctl docs/help. Prefer organization, team,
or visible-to-authenticated-user language. Generated schema/API snapshots and
explicit domain-rule docs are excluded from this check.
MSG
  exit 1
fi
