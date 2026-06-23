# Linear CLI Feature Leech

This doc lists ideas worth borrowing from comparable Linear/agent/API CLIs into linctl —
in plain words, with why each is worth taking and what it brings — plus the ideas we should
*not* copy. It was built from a primary-source read of 47 comparable repos, and every
"already have / missing / partial" call below was checked against linctl's actual code.

## linctl's moat (everything borrowed must preserve this)

linctl's value is **target-pinned writes**: before any mutation it re-resolves the active
credential's org, team, and optional project and compares them to `.linctl.toml`. On a
mismatch the write fails immediately — no prompt, no `--force`, no soft warning. Reads are
free; writes fail closed. Borrowed features must keep this, keep operations schema-aligned
(genqlient, no raw query strings on the main surface), prefer archive over hard delete, and
never print token values.

## Comparable repositories

| Repo | Stack | Best idea |
|---|---|---|
| `schpet/linear-cli` | Deno/TS | Git/jj-aware `issue start` with branch creation; generated SKILL.md kept in sync with the binary |
| `linearis-oss/linearis` | Node/TS | Two-tier `usage` discovery (top-level <200 tokens, domain <500); explicit "CLI instead of MCP" agent framing |
| `joa23/linear-cli` | Go | Skill installer + task export for Claude Code; cycle analytics; token-efficient ASCII output |
| `Finesssee/linear-cli` | Rust | Broad output flags (`--fail-on-empty`, `--id-only`, `--dry-run`, `--fields`); dynamic completions querying live Linear |
| `flipbit03/lineark` | Rust | Usage-first agent setup snippet; human-readable names instead of raw UUIDs |
| `frr149/lql` | Rust | Tolerant alias contract: accepts common agent mistakes, then reports the normalization |
| `Securiteru/linear-cli` | Go | JSON-lines batch create; quiet flag that suppresses all output on success |
| `juanbermudez/linear-agent-cli` | Deno/TS | Smart resource resolution (name/identifier/ID interchangeable); structured JSON errors with codes |
| `choam2426/Linear-CLI` | Python | Single-file agent CLI; JSON errors on stderr; idempotent `save-*` (create-or-update) |
| `Valian/linear-cli-skill` | JS | Minimal Claude skill bundled with the CLI so invocation stays in sync |
| `dabblewriter/linear-cli` | Node | `--unblocked` + `next` worktree flow; append-to-description update mode |
| `nikpietanze/linear-cli` | Go | Template-driven issue creation with cache + dry-run preview |
| `evangodon/linear-cli` | Node/oclif | Multiple output formats (csv/json/yaml); column selection; list sort/filter |
| `rubyists/linear-cli` | Ruby | Short command aliases (`lcls`, `lcomment`, `lclose`); semantic PR title workflow |
| `AdiKsOnDev/linear-cli` | Python | Project health update history as a first-class read surface |
| `danielrearden/linear-cli` | Node | Cached org metadata with an explicit `cache clear` |
| `dorkitude/linctl` | Go | Dynamic MCP namespace; doctor-style smoke command |
| `tdwells90/lncli` | Rust | Structured agent error envelopes with a machine-readable `error_code` |
| `0x80/pick-linear-ticket` | TS | Ranked next-work picker (unblock count, priority, age) → `branchName` JSON |
| `alleneubank/linear-cli` | Zig | Config hygiene: `0600` config, redacted auth show, separate `auth test` |
| `aliou/linear-cli` | TS | Multi-profile auth with per-directory profile selection (explicit, not implicit) |
| `TrevorS/linear-cli` | Rust | Markdown/frontmatter issue creation; browser open; dry-run shows parsed fields |
| `wiseiodev/linear-cli` | TS | `bulk-update --dry-run --json` to plan batch mutations |
| `anoncam/linear-cli` | JS | AI-assisted label cleanup proposed before any write |
| `rinadelph/clinear` | Python | Typed validation (Pydantic/Typer) that rejects malformed input before any API call |
| `ChuckMayo/linear-multi-workspace-cli` | TS | Multi-workspace surface — a useful *contrast* to our pinned single-target model |

## What linctl already has (so we don't re-propose it)

- **Output controls**: `--json`, `--compact`, `--fields`, `--id-only`, `--quiet`,
  `--fail-on-empty`, `--sort`/`--order`, `--format minimal|compact|full`, `--debug`,
  `--profile`, `--org`/`--team`/`--project` overrides, `--timeout`.
- **Discovery / health**: `usage [overview|issue|project|cycle]`, `target --json`,
  `doctor` (reports config/token/target without printing the token), `whoami`.
- **Issue reads**: `issue list` (state, project, assignee, label, cycle, created-after/
  since/before, has-blockers, blocks, blocked-by, all-teams, mine), `search`, `get`,
  `deps`, `pr`, `id`/`title`/`url`/`branch`, `vcs-branch-search`, AI helpers
  (filter/title suggestion, figma key, priority-values).
- **Issue writes**: `create` (+`--description-file`), `update` (+`--append`/`--append-file`/
  `--description-file`), `start`, `comment` (+`--body-file`, `--body -`), `reply`, `close`,
  template-backed create, guarded import, relation writes, and comment update/delete.
- **Current-branch flow**: `current`, `done`, `next --dry-run` (ranked by unblock count,
  priority, age).
- **Projects**: list/get/members, create/update/archive, `updates` (read), plus
  `project-update` create/read and `project-status`/`project-label` reads.
- **ProjectMilestone**: all/list/get/create/update (no delete). **Cycle**: list/get/create/
  update/archive, issues, uncompleted-issues. **Sprint**: current/report (read-only aliases).
- **v0.3.0 coordination surface**: guarded document create/update, project-update create,
  issue relate/unrelate, comment update/delete, issue template dry-run, `next --checkout`,
  file upload/download, browser open, issue export, and CSV/JSON import/export.
- **Reads across the schema**: documents, templates, initiatives, roadmaps (legacy),
  search + semantic-search, labels/teams/users/workflow-states, releases & pipelines,
  comments, attachments, notifications, custom views, organization, rate-limit, customers,
  and more.
- **Packaging is live for `v0.3.0`**: `.goreleaser.yaml` +
  `.github/workflows/release.yml` build cross-platform binaries and generate the Homebrew
  cask (`dist/homebrew/Casks/linctl.rb`) from `v*` tags.
- **Quality**: 100% statement-coverage gate, ~30 linters, build-tag integration tests,
  fuzz tests, genqlient/skill drift checks, and `actionlint` in the local `task ci` gate.

The linctl agent skill lives at `skills/linctl/SKILL.md`; its generated command reference is
drift-checked from the Cobra tree.

---

## P1 — high value, safe, start now

### 1. Tolerant alias normalization for agent mistakes
*(status: done — `issue create`/`update`/`list` normalize `--status`→`--state`, human state
names → schema state types, and priority words → Linear ints before target comparison, printing
one stderr note per changed field and rejecting unknown values. Sources: `frr149/lql`,
`linearis-oss/linearis`.)*

- **What it is**: when an agent passes a near-miss — `--status` for `--state`, `high`
  instead of Linear's priority integer, a human state name — linctl maps it to the correct
  value, runs the command, and prints **one line on stderr** saying what it normalized.
  Truly ambiguous or destructive input is still rejected.
- **Why leech it**: agents constantly emit field names from training data that differ from
  Linear's schema. lql showed this kills a whole class of failed runs without loosening
  safety, because normalization happens *before* target comparison.
- **What it brings**: fewer dead-end runs on valid intent; agents stop needing to memorize
  Linear's enum internals.
- **Moat fit**: pure pre-write transform on flag values. Doesn't touch `write_guard.go`;
  the normalized value resolves against the pinned target exactly as a raw value would.

### 2. Generated command reference with a CI drift check
*(status: done — `scripts/gen-skill.go` generates `skills/linctl/references/commands.md`
from the cobra tree; `task gen-skill` refreshes it and `generate-check` + CI fail on drift.
The curated `SKILL.md` is hand-maintained and links to the generated reference rather than
being overwritten. Sources: `schpet/linear-cli`, `joa23/linear-cli`.)*

- **What it is**: a generator reads every command's `--help` and writes
  `skills/linctl/references/commands.md` (the curated `SKILL.md` links to it). CI re-runs it
  and fails if the committed file differs from the binary's surface.
- **Why leech it**: agents read the skill doc to learn the commands. Without a drift check
  it silently goes stale as commands are added, so agents call things that don't exist or
  miss things that do.
- **What it brings**: the agent-facing doc always matches the running binary; CI blocks
  drift before merge.
- **Moat fit**: pure docs generation at build time. No writes, no tokens, no schema change.

### 3. Dynamic shell completions
*(status: done — `internal/cli/completion.go` adds live flag completion for `--team`,
`--project`, and `--state`, plus `ValidArgsFunction` for `team get`/`project get`, all via
existing read paths and degrading silently without a token; static completions documented in
`skills/linctl/references/completions.md`. Sources: `Finesssee/linear-cli`, `aliou/linear-cli`,
`AdiKsOnDev/linear-cli`.)*

- **What it is**: wire `ValidArgsFunction`/flag-completion so the shell can complete *live*
  values — team keys, project IDs, workflow-state names — by calling read commands. Also
  document and test the static completions we already get for free.
- **Why leech it**: Cobra makes this cheap, and value-completion is the difference between
  "completion exists" and "completion is useful." Several peers ship it.
- **What it brings**: operators stop hand-looking-up IDs; far fewer typos at the prompt.
- **Moat fit**: read-only. Completion callbacks use existing read paths; the write guard is
  never involved.

### 4. Structured error envelope on stderr
*(status: done — `execute` in `internal/cli/root.go` emits one JSON `{error_code, message}`
line to stderr on any failure; `errorCode` in `output.go` maps the `errors.Is` sentinels
(`TARGET_MISMATCH`, `RATE_LIMITED`, `MUTATION_FAILED`, `INVALID_WRITE`, `GRAPHQL_ERROR`),
not-found, and `INTERNAL`. Human text still follows. Documented in
`references/json-output.md`. Sources: `tdwells90/lncli`, `choam2426/Linear-CLI`,
`juanbermudez/linear-agent-cli`.)*

- **What it is**: on failure, emit a JSON object to stderr with a machine `error_code`
  (`TARGET_MISMATCH`, `NOT_FOUND`, `RATE_LIMITED`, …), a human `message`, and optional
  context. Human-readable text stays available too.
- **Why leech it**: agents that shell out to linctl need to branch on *what* failed without
  parsing English. `TARGET_MISMATCH` must never look like `NOT_FOUND`.
- **What it brings**: reliable agent retry/escalation — retry only on `RATE_LIMITED`,
  inspect config on `TARGET_MISMATCH`, give up on `NOT_FOUND`.
- **Moat fit**: output-only and additive. Maps the existing `errors.Is` sentinels to codes;
  doesn't touch the guard.

---

## P2 — high coordination value, after P1

### 5. Document writes (create / update)
*(status: done — `document create`/`document update` in `internal/cli/document.go` backed by
`CreateDocument`/`UpdateDocument` (`internal/client/document_write.go`): create is team-scoped
(+ pinned project), update resolves the document and fails closed unless its team and pinned
project match. `--content`, `--content-file`, and `--content -` (stdin) supported. Sources:
`schpet`, `linearis`.)*
- **What/why/brings**: `document create --title T --content-file F` and
  `document update ID --content-file F` (both `--body -` for stdin) let agents write plans,
  notes, and spec drafts into Linear Documents mid-workflow instead of bouncing to the web
  UI/MCP.
- **Moat fit**: same guarded-write pattern as issue writes (resolve → compare team/optional
  project → mutate).

### 6. Project status-update writes
*(status: done — `project-update create PROJECT --health --body` in
`internal/cli/project_update.go` backed by `CreateProjectUpdate`
(`internal/client/project_update_write.go`): resource-scoped via `requireProject`, fails
closed unless the resolved project matches the pinned `project_id` and team. `--health`
accepts tolerant aliases (`on-track`/`onTrack`/...), `--body`/`--body-file`/`--body -`
(stdin) supported. Sources: `AdiKsOnDev`, `schpet`, `Finesssee`.)*
- **What/why/brings**: `project-update create PROJECT --health onTrack|atRisk|offTrack
  --body T` lets an agent post automated health-check results back to Linear, keeping
  dashboards current without a human.
- **Moat fit**: resource-scoped write — resolve project, compare to pinned `project_id` if
  set, then create.

### 7. Relation writes (`issue relate` / `unrelate`)
*(status: done — `issue relate ISSUE RELATED --type` and `issue unrelate RELATION_ID` in
`internal/cli/issue_relation_write.go` backed by `CreateIssueRelation`/`DeleteIssueRelation`
(`internal/client/issue_relation_write.go`): both endpoints resolve through `requireIssue`
and fail closed unless each issue belongs to the resolved team. `--type blocks` is rejected
when the related issue already blocks the source issue (cycle pre-check via the existing deps
read). Sources: `Finesssee`, `joa23`.)*
- **What/why/brings**: `issue relate LIT-123 blocks LIT-456` / `unrelate …` lets agents
  maintain the dependency graph as work progresses. With existing `issue deps` reads, this
  turns Linear into a live execution state machine for agent-driven work.
- **Moat fit**: resolve both issues; guard compares the pinned team against both. Add
  circular-dependency detection as a pre-write check using the existing deps read.

### 8. Comment edit / delete
*(status: done — `comment update ID --body` and `comment delete ID` in
`internal/cli/comment.go` backed by `UpdateComment`/`DeleteComment`
(`internal/client/comment_write.go`): both resolve the comment and fail closed unless its
parent issue belongs to the resolved team (`guardCommentTarget`); a comment not attached to
an issue is refused. `comment delete` is the one approved delete. `--body`/`--body-file`/
`--body -` (stdin) supported on update. Sources: `linearis`, `rubyists`.)*
- **What/why/brings**: `comment update ID --body T` / `comment delete ID` let an agent fix
  or remove a progress note it posted, instead of stacking stale comments.
- **Moat fit**: resource-scoped write — resolve the comment's parent issue, compare to
  pinned target, then mutate.

### 9. `issue create --template` with `--dry-run`
*(status: done — `issue create --template ID --section NAME=VALUE --dry-run` in
`internal/cli/issue_write.go`/`issue_template.go`: `--template` reads `Template.templateData`
(free read via `GetIssueTemplateContent`, tolerant of object or JSON-encoded-string data) and
fills title/description defaults that explicit flags override; `--section` fills/appends a
markdown section locally; `--dry-run` renders the assembled draft and performs no mutation.
The real write reuses the guarded `issue create` (`CreateIssue`). Sources: `nikpietanze`,
`Finesssee`.)*
- **What/why/brings**: apply a Linear template by ID with `--section Name=Value`;
  `--dry-run` prints the rendered description without writing. Consistent, high-quality
  issues from agents, previewable before commit.
- **Moat fit**: write path reuses the guarded `issue create`; dry-run is local and never
  hits the API.

### 10. `next` with branch checkout
*(status: done — `next` in `internal/cli/next.go` now starts the picked issue through the
guarded `StartIssue`; `--dry-run` keeps the read-only preview, and `--checkout` runs
`git checkout -b <branchName>` (injectable `checkoutBranch`) before starting. The pick still
goes through target comparison and `StartIssue` re-runs the write guard. Sources: `schpet`,
`dabblewriter`, `0x80/pick-linear-ticket`.)*
- **What/why/brings**: after picking the top unblocked issue, optionally
  `git checkout -b <branchName>` so pick-and-start is one command — closing the window where
  the agent might act on a different issue between steps.
- **Moat fit**: the pick still goes through target comparison; branch creation is a local
  git op; the existing guarded `issue start` can run on the same issue.

### 11. File upload / download for Linear assets
*(status: done — `files upload PATH` and `files download URL --output PATH` in
`internal/cli/files.go` backed by `PrepareFileUpload` (`internal/client/file.go`): upload calls
`fileUpload` for a pre-signed target, PUTs the bytes with the returned headers (injectable
`fileHTTPClient`), and prints the asset URL for a later guarded attachment write; download is a
plain unauthenticated GET so a user-supplied URL never receives the Linear token. Sources:
`linearis`, `choam2426`, `joa23`.)*
- **What/why/brings**: `files upload PATH` → asset URL; `files download URL --output PATH`.
  Lets agents attach CI artifacts, screenshots, and diagrams to issues without a browser.
- **Moat fit**: upload returns a URL; attaching it uses the existing attachment write path.
  Download is read-only.

---

## P3 — nice to have / lower urgency

### 12. `issue open` / `project open` (browser)
*(status: done — `issue open ISSUE_ID` and `project open PROJECT_ID` in
`internal/cli/open.go` resolve the entity's existing `url` via the free `GetIssueByID` /
`GetProjectByID` reads, then launch the platform opener (`xdg-open`/`open`/`rundll32`,
injectable `openExecutor`) with the URL as a discrete argv argument — no shell. Read-only,
no guard; respects `--id-only`/`--json`/`--quiet`. Sources: `schpet`, `TrevorS`.)* Open the
entity's existing URL via `xdg-open`. Read-only, no guard. Smooths CLI→web handoff for human
operators.

### 13. `issue export ISSUE DIR`
*(status: done — `issue export ISSUE_ID DIR` in `internal/cli/export.go` assembles
`GetIssueDetail` (description + metadata), `ListIssueComments`, and `ListIssueAttachments`
into `<DIR>/<identifier>.md` (metadata header → description → comments → attachment URLs).
All three are free reads; the only write is the local markdown file. Comments/attachments
are capped at 250 with a stderr note when more pages exist. Respects
`--id-only`/`--json`/`--quiet`. Sources: `joa23`, `choam2426`.)* Dump description + comments
+ attachment URLs to local files for retrospectives or feeding a local LLM. Read-only
assembly of existing reads.

### 14. CSV/JSON import-export with `--dry-run`
*(status: done — `issue import FILE` and `issue bulk-export FILE` in
`internal/cli/bulk.go`. Import parses CSV or JSON (format from the extension), normalizes
each row's state/priority, rejects any row whose `team` key ≠ the pinned `team_key`, then
creates each issue through the guarded `CreateIssue` (each create re-runs the write guard);
`--dry-run` renders the normalized rows and writes nothing. `bulk-export` writes the resolved
team's issues (`ListIssuesByTeam`, `--limit` default 250) to a CSV or JSON file — the only
write is the local file. Both respect `--id-only`/`--json`/`--quiet`. Sources: `Finesssee`,
`wiseiodev`.)* Bulk create from a file with a preview step; export the team's issues to a
file. Each create goes through the guarded `CreateIssue`; import rejects rows whose team key
≠ pinned target.

### 15. Finish the release (packaging is already built)
*(status: done — `v0.3.0` is tagged at `ee4b3b9` after the feature-leech merge and
the local `actionlint` gate.)*
goreleaser (`.goreleaser.yaml`), the tag-triggered release workflow
(`.github/workflows/release.yml`, `on: push: tags: v*`), and the Homebrew cask
(`homebrew_casks` → `KyaniteHQ/homebrew-linctl`) are all wired: the workflow runs the full
preflight gate, then `goreleaser release --clean` publishes the GitHub release, archives,
checksums, SBOM, cosign signature, and the cask, and a final job downloads and verifies the
published archive. The binary version is injected from the tag via
`-X main.version={{ .Version }}`. The local repo state now has `HEAD`, `origin/master`, and
tag `v0.3.0` on the actionlint-gated release commit; no follow-up tag command remains for
this branch state.

---

## Do not leech

| Thing | Why not |
|---|---|
| Raw GraphQL escape hatch (`--query` for arbitrary mutations) | Bypasses schema-aligned ops and the write guard; an agent could mutate anything regardless of pinned target. Predictable, constrained writes are the whole point. |
| Implicit workspace switching (auto-pick org from repo) | The pinned target is the guarantee that writes land in the right org/team. Implicit switching silently breaks it and creates hard-to-diagnose mismatches. |
| Hard delete for issues/projects | Archive is the safe cleanup path. Hard delete is irreversible; it needs a separate explicit design that doesn't exist, and the data-loss risk outweighs the benefit. |
| Printing/logging token values | `doctor` deliberately reports `set`, not the value. Any path that could leak a token breaks the never-print-tokens rule. |
| OAuth device flow as the *primary* auth | For an agent-first CLI, API-key auth is simpler and CI-stable (no browser redirect). OAuth can be a secondary path later if asked. |
| Watch/monitor polling loops | Adds process-lifecycle complexity; agents are better served calling explicit poll commands on their own schedule. Add only on real demand. |
| Interactive TUI | Incompatible with agent use, which needs deterministic scriptable output; pulls in deps that complicate the output contract. |
| OS keychain/keyring token storage | Adds platform-specific deps and complicates CI token flow. `LINCTL_TOKEN` env + `.linctl.toml` is simpler and more portable. |
