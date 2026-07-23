#!/usr/bin/env bash
# shellcheck shell=bash
# Open or update the chore/schema-refresh PR after scripts/schema-refresh-from-upstream.sh.
# Intended for GitHub Actions (needs gh + push credentials). Safe to re-run: force-pushes
# the dedicated bot branch and edits an existing open PR when present.
#
# Required env:
#   BRANCH          bot branch name (e.g. chore/schema-refresh)
#   BASE_BRANCH     PR base (e.g. master)
#   UPSTREAM_SHA    full linear/linear SHA the refresh targets
#   SHORT_SHA       short SHA for titles
#   DRIFT_REPORT    pre-refresh drift report text
#   BREAKING        true|false
#   GH_TOKEN        token for gh + git push (set by the workflow)
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

: "${BRANCH:?BRANCH is required}"
: "${BASE_BRANCH:?BASE_BRANCH is required}"
: "${UPSTREAM_SHA:?UPSTREAM_SHA is required}"
: "${SHORT_SHA:?SHORT_SHA is required}"
: "${DRIFT_REPORT:?DRIFT_REPORT is required}"
: "${BREAKING:?BREAKING is required}"

git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

git add \
  internal/client/schema.graphql \
  internal/client/internal/gql/generated.go \
  docs/internal/linear-api-coverage.md \
  Taskfile.yml \
  .github/workflows/ci.yml

if git diff --cached --quiet; then
  printf 'refresh produced no staged changes; nothing to open\n'
  exit 0
fi

git commit -m "$(cat <<EOF
Refresh vendored Linear schema and SDK ref pin from upstream

Sync to linear/linear ${SHORT_SHA} so the nightly schema-drift job stays green.
EOF
)"

# Dedicated bot branch: replace tip each run so the PR tracks base + latest schema.
git push --force origin "HEAD:refs/heads/${BRANCH}"

title="Refresh vendored Linear schema (${SHORT_SHA})"

warning=""
if [[ "$BREAKING" == "true" ]]; then
  warning="$(cat <<'EOF'
## Review carefully

This drift includes removals or field type changes. Additive-only refreshes are usually safe; this one may need operation updates before merge.

EOF
)"
fi

body="$(cat <<EOF
## What changed

Automated refresh of the vendored Linear GraphQL schema against \`linear/linear\` \`${SHORT_SHA}\` (\`${UPSTREAM_SHA}\`).

- \`internal/client/schema.graphql\` synced from \`packages/sdk/src/schema.graphql\`
- \`LINEAR_SDK_REF\` / \`LINCTL_LINEAR_SDK_REF\` pins bumped
- genqlient client regenerated
- \`docs/internal/linear-api-coverage.md\` regenerated

${warning}## Drift report (before refresh)

\`\`\`
${DRIFT_REPORT}
\`\`\`

## Safety

- [x] Reads stay free; no write-path changes in this refresh
- [x] Operations remain schema-aligned (client regenerated from the new snapshot)
- [x] No token value is printed or logged

## Checks

- [ ] \`task ci\` green on this PR
- [ ] Human review of removals / type changes if called out above

Opened by the \`schema-refresh\` workflow. Merge closes the nightly schema-drift failure loop.
EOF
)"

existing="$(
  gh pr list \
    --head "$BRANCH" \
    --base "$BASE_BRANCH" \
    --state open \
    --json number \
    --jq '.[0].number // empty'
)"
if [[ -n "$existing" ]]; then
  gh pr edit "$existing" --title "$title" --body "$body"
  url="$(gh pr view "$existing" --json url --jq .url)"
  printf 'updated existing PR: %s\n' "$url"
else
  url="$(
    gh pr create \
      --base "$BASE_BRANCH" \
      --head "$BRANCH" \
      --title "$title" \
      --body "$body"
  )"
  printf 'opened PR: %s\n' "$url"
fi
