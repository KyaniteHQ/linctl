# Linear CLI Feature Leech Notes

This file is a repo-grounded feature extraction from comparable GitHub Linear CLI
projects cloned into `/tmp/linctl-linear-cli-research` on 2026-06-19.

`linctl`'s existing differentiator is target-pinned writes: writes resolve the
active Linear org, team, and optional project before mutating data. Borrow
features below only when they preserve that fail-closed model.

## Cloned Repositories

| Repo | Local clone | Stack | Good features to leech |
| --- | --- | --- | --- |
| `schpet/linear-cli` | `/tmp/linctl-linear-cli-research/schpet__linear-cli` | Deno/TypeScript | Git and jj aware `issue start`, current issue views, PR creation, generated agent skill docs, documents, milestones, comments, browser/app open commands, project-local config. |
| `linearis-oss/linearis` | `/tmp/linctl-linear-cli-research/linearis-oss__linearis` | Node/TypeScript | Two-tier `usage` discovery, all-JSON output, smart ID resolution, issue discussion/reply model, file upload/download guidance, explicit "CLI instead of MCP" agent prompt. |
| `joa23/linear-cli` | `/tmp/linctl-linear-cli-research/joa23__linear-cli` | Go | Token-efficient ASCII formats, OAuth user vs agent/app modes, dependency graph search, task export for Claude Code, issue attachment/export, cycle analytics, skill installer, `.linear.yaml` defaults. |
| `Finesssee/linear-cli` | `/tmp/linctl-linear-cli-research/Finesssee__linear-cli` | Rust | Very broad command surface, output flags (`--fields`, `--id-only`, `--fail-on-empty`, `--dry-run`), import/export, dynamic completions, watch mode, webhooks, project updates, saved views, status/metrics/triage. |
| `flipbit03/lineark` | `/tmp/linctl-linear-cli-research/flipbit03__lineark` | Rust | Compact `usage` for agents, human-readable names instead of UUIDs, generated SDK/codegen split, file embeds, profiles, usage-first agent setup snippet. |
| `frr149/lql` | `/tmp/linctl-linear-cli-research/frr149__lql` | Rust | Adversarial CLI design, tolerant aliases for common agent mistakes, TOON compact output, credential helper support, cwd context map, retired-team hints, schema-first tests from real fixtures. |
| `Securiteru/linear-cli` | `/tmp/linctl-linear-cli-research/Securiteru__linear-cli` | Go | Small static-binary feel, psst secret-injection workflow, JSON lines batch create, quiet output, broad but simple command list. |
| `juanbermudez/linear-agent-cli` | `/tmp/linctl-linear-cli-research/juanbermudez__linear-agent-cli` | Deno/TypeScript | Agent-first fork of schpet-style CLI, VCS integration, cross-entity operations, JSON/error-code discipline, smart resource resolution, caching. |
| `choam2426/Linear-CLI` | `/tmp/linctl-linear-cli-research/choam2426__Linear-CLI` | Python | Single-file no-dependency agent CLI, compact JSON, simple command table, `save-*` create-or-update pattern, upload/download for Linear file URLs, JSON errors on stderr. |
| `Valian/linear-cli-skill` | `/tmp/linctl-linear-cli-research/Valian__linear-cli-skill` | JavaScript | Minimal Claude skill plus CLI, `resource action` command shape, `@me`, sub-issues, project assignment, JSON flag everywhere. |
| `dabblewriter/linear-cli` | `/tmp/linctl-linear-cli-research/dabblewriter__linear-cli` | Node.js | Zero-dependency Node CLI, `--unblocked`, `next` worktree flow, `done` cleanup flow, issue branch creation, per-project `.linear` config, append-to-description update mode. |
| `nikpietanze/linear-cli` | `/tmp/linctl-linear-cli-research/nikpietanze__linear-cli` | Go | AI-oriented one-command issue creation from Linear templates, automatic template sync/cache, no-delete production safety, CI issue creation examples. |
| `evangodon/linear-cli` | `/tmp/linctl-linear-cli-research/evangodon__linear-cli` | Node/oclif | Older but useful generated command docs, cache show/refresh, workspace switch, output formats (`csv`, `json`, `yaml`), column selection, list sorting/filtering. |
| `rubyists/linear-cli` | `/tmp/linctl-linear-cli-research/rubyists__linear-cli` | Ruby | OCI container install path, wrapper alias script, strong aliases (`lcls`, `lcomment`, `lclose`), multi-issue update/comment/close, semantic PR title workflow. |
| `AdiKsOnDev/linear-cli` | `/tmp/linctl-linear-cli-research/AdiKsOnDev__linear-cli` | Python | OAuth plus API-key auth, keyring storage, YAML output, project health updates, project update history, advanced search, shell completion. |
| `danielrearden/linear-cli` | `/tmp/linctl-linear-cli-research/danielrearden__linear-cli` | Node | Simple interactive create/status flows, keychain storage, branch checkout by issue, cached org-specific metadata with clear command. |

## Highest-Value Features For `linctl`

### 1. Agent-sized command discovery

`linctl usage` already exists. Make it more like `linearis usage`,
`linearis <domain> usage`, and `lineark usage`:

- Keep top-level `usage` under roughly 200 tokens.
- Keep domain usage under roughly 500 tokens.
- Include examples that are copy-pasteable and safe.
- Generate `skills/linctl/SKILL.md` from command help so CLI docs and agent docs
  cannot drift.
- Add a CI check that generated skill docs are current.

Useful evidence:
- `linearis-oss__linearis/README.md`
- `flipbit03__lineark/README.md`
- `schpet__linear-cli/skills/linear-cli/scripts/generate-docs.ts`

### 2. Machine output controls

Keep `--json`, then add script-friendly output flags from `Finesssee` and
`joa23`:

- `--fields identifier,title,state.name` for JSON field projection.
- `--compact` for no pretty-printing.
- `--id-only` for create/update chaining.
- `--quiet` for scripts.
- `--fail-on-empty` for monitors and automations.
- `--sort field --order asc|desc` for deterministic list output.
- `--format minimal|compact|full` for human/agent text output.

Do not add lossy output as the only surface. JSON must remain stable.

Implemented slice:
- `--compact` keeps JSON stable but removes pretty-print whitespace.
- `--fields` projects JSON keys, including list item projection for `issues`, `projects`, and `members`.
- `--id-only`, `--quiet`, `--fail-on-empty`, `--sort field --order asc|desc`, and
  `--format minimal|compact|full` are global CLI flags.
- Evidence: `go test ./internal/cli`, especially the `Test_CommandFlows_*` output-control tests.

Useful evidence:
- `Finesssee__linear-cli/README.md`
- `joa23__linear-cli/README.md`
- `evangodon__linear-cli/README.md`

### 3. Tolerant but explicit agent ergonomics

Borrow `lql`'s tolerance contract:

- Accept common non-destructive aliases, then report the normalization.
- Examples: `--status` -> `--state`, priority names -> Linear priority ints,
  human state names -> schema state types.
- Reject destructive or ambiguous input instead of guessing.
- Suggest the right command for wrong flags, especially `--comment` on update.

This fits `linctl` if normalized values are still resolved against the pinned
target before writes.

Useful evidence:
- `frr149__lql/README.md`

### 4. Branch and work-start flows

`linctl current` already derives an issue from git/jj context. Expand around it:

- `issue id`, `issue title`, `issue url` aliases for current branch.
- `issue start ISSUE` to assign/start and optionally create/check out a branch.
- `issue branch ISSUE` to print the canonical branch slug without checkout.
- `issue pr ISSUE` or `current pr` to generate a `gh pr create` title/body.
- `done` to close the current branch issue after target comparison.
- Optional later: extend `next` beyond dry-run to create/check out a branch or worktree.

Implemented slice:
- `issue id` prints the Current Issue identifier from git/jj context.
- `issue title [ISSUE]` and `issue url [ISSUE]` print one scalar from either an explicit issue or the Current Issue.
- `issue branch ISSUE` prints Linear's `branchName` without checkout.
- `issue start ISSUE` assigns the issue to the authenticated viewer and moves it to the team's first started workflow state after target comparison; branch checkout is intentionally still future work.
- `done` closes the Current Issue from the checkout after the same target comparison as `issue close`.
- `issue pr [ISSUE]` reads either an explicit issue or the Current Issue and prints a local `gh pr create` title/body plan without calling GitHub.
- `next --dry-run` resolves the target, reads unstarted issues with no blocking relations, and prints the first candidate without creating a branch or worktree.
- Evidence: `go test ./internal/cli`, especially the `Test_CommandFlows_print_current_issue_*`
  tests, `Test_CommandFlows_print_issue_branch_from_issue_branch`, and
  `Test_CommandFlows_execute_read_and_write_commands/issue_start`, plus
  `Test_CommandFlows_close_current_issue_from_done`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_pr`, and
  `Test_CommandFlows_print_issue_pr_from_current_branch`, plus
  `Test_CommandFlows_execute_read_and_write_commands/next_dry_run`.

Useful evidence:
- `schpet__linear-cli/README.md`
- `dabblewriter__linear-cli/README.md`
- `rubyists__linear-cli/Readme.adoc`
- `Finesssee__linear-cli/README.md`

### 5. Search and list filters

Current `issue list` now includes the highest-value read-only filters from the comparable CLIs.

Keep writes target-pinned; broad reads can stay less restrictive.

Implemented slice:
- `issue list --state TYPE` filters by Linear workflow state type, such as `started` or `completed`.
- `issue list --project PROJECT_ID` filters by Linear project id inside the resolved team.
- `issue list --mine` filters by authenticated viewer id inside the resolved team.
- `issue list --assignee USER_ID` filters by Linear assignee user id inside the resolved team.
- `issue list --label LABEL_ID` filters by Linear issue label id inside the resolved team.
- `issue list --cycle CYCLE_ID` filters by Linear Cycle id inside the resolved team.
- `issue list --created-after DATE` filters by `Issue.createdAt.gte` inside the resolved team.
- `issue list --created-since DATE` is an alias for the `Issue.createdAt.gte` filter inside the resolved team.
- `issue list --created-before DATE` filters by `Issue.createdAt.lte` inside the resolved team.
- `issue list --has-blockers` filters by `Issue.hasBlockedByRelations.eq` inside the resolved team.
- `issue list --blocks` filters by `Issue.hasBlockingRelations.eq` inside the resolved team.
- `issue list --blocked-by ISSUE` traverses that issue's `blocks` relations and returns blocked issues inside the resolved team.
- `issue list --all-teams` lists issues across every visible Linear team as a read-only broad inspection command.
- `issue search QUERY` searches issue content within the resolved team.
- Broader dependency graph commands remain future work.
- Evidence: `go test ./internal/cli`, especially `Test_CommandFlows_execute_read_and_write_commands/issue_list_state_filter`
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_project_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_mine_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_assignee_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_label_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_cycle_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_created-after_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_created-since_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_created-before_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_has_blockers_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_blocks_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_blocked_by_filter`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_list_all_teams`, and
  `Test_CommandFlows_execute_read_and_write_commands/issue_search`.

Useful evidence:
- `joa23__linear-cli/README.md`
- `schpet__linear-cli/README.md`
- `Finesssee__linear-cli/README.md`
- `evangodon__linear-cli/README.md`

### 6. Dependency graph commands

This is high leverage for agents because it turns Linear into an execution
state machine:

- `issue deps ISSUE` prints blockers, blocked issues, parent, and children.
- A broader `deps --team TEAM --project PROJECT` graph view remains future work.
- Detect circular dependencies and exit non-zero.
- `issue relate ISSUE blocks OTHER` and `issue unrelate ...`.
- Use dependency filters in `issue search`.

Implemented slice:
- `issue deps ISSUE --limit N` reads one issue's parent, children, outgoing `blocks` relations, and incoming `blocks` relations.
- Text output groups related issues under `parent`, `children`, `blocks`, and `blocked_by`; JSON output returns the same graph as `IssueDependencyGraph`.
- Team/project graph views, circular dependency detection, relation writes, and dependency filters in `issue search` remain future work.
- Evidence: `go test ./internal/cli`, especially `Test_CommandFlows_execute_read_and_write_commands/issue_deps`.

Useful evidence:
- `joa23__linear-cli/README.md`
- `Finesssee__linear-cli/README.md`
- `frr149__lql/README.md`

### 7. Comments, discussions, and append flows

`linctl issue comment` exists. Improve it:

- List comments/discussions.
- Reply to a root discussion thread.
- Edit/delete only if target comparison proves the issue/project scope.
- `--body -` to read from stdin.
- `issue update --append FILE_OR_TEXT` for progress notes without replacing the
  description.
- Prefer issue discussion commands over a separate top-level comments domain.

Implemented slice:
- `issue comment ISSUE --body -` reads the full command stdin as the comment body before running the guarded write.
- `issue comments ISSUE --limit N` lists issue discussion comments.
- `issue reply ISSUE COMMENT --body TEXT` creates a threaded reply with `CommentCreateInput.parentId` after the same issue target comparison.
- `issue update ISSUE --append TEXT` reads the existing description, appends text with a blank-line separator, and writes the combined description through the guarded issue update path.
- Edit/delete comment flows and file-backed append remain future work.
- Evidence: `go test ./internal/cli`, especially `Test_CommandFlows_read_issue_comment_body_from_stdin`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_comments`,
  `Test_CommandFlows_execute_read_and_write_commands/issue_reply`, and
  `Test_CommandFlows_execute_read_and_write_commands/issue_update_append`.

Useful evidence:
- `linearis-oss__linearis/README.md`
- `dabblewriter__linear-cli/README.md`
- `rubyists__linear-cli/Readme.adoc`
- `joa23__linear-cli/README.md`

### 8. Projects, milestones, and project updates

`linctl project` already has list/get/members/create/update/archive. Next
project-domain features:

- `project update-status PROJECT --body ... --health onTrack|atRisk|offTrack`
- `project-milestone get|create|update|delete`
- `project issues PROJECT`
- `project labels add|remove|set`
- `project open PROJECT` and `project url PROJECT`

Every write should compare the resolved project against `[target].project_id`
when configured.

Implemented slice:
- `project updates PROJECT --limit N` lists project status updates with health, author, body, timestamps, and URL.
- `project-milestone list PROJECT --limit N` lists a project's milestones with status, target date, progress, and sort order.
- Project update creation/edit/archive and project milestone get/create/update/delete remain future work.
- Evidence: `go test ./internal/cli`, especially
  `Test_CommandFlows_execute_read_and_write_commands/project_updates` and
  `Test_CommandFlows_execute_read_and_write_commands/project_milestone_list`, and
  `go test ./internal/client`, especially
  `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

Useful evidence:
- `AdiKsOnDev__linear-cli/README.md`
- `Finesssee__linear-cli/README.md`
- `schpet__linear-cli/README.md`
- `choam2426__Linear-CLI/plugins/claude/skills/linear-cli/SKILL.md`

### 9. Documents, files, and attachments

These are valuable because agents often need to move plans, screenshots, and
artifacts into Linear:

- `document list|get|create|update`
- `document create --content-file FILE`
- `document create` from stdin.
- `files upload PATH` returning an asset URL.
- `files download URL --output PATH`
- `issue export ISSUE DIR` to write description, comments, and assets locally.

Target comparison should follow the attached issue/project/team.

Useful evidence:
- `schpet__linear-cli/README.md`
- `linearis-oss__linearis/README.md`
- `choam2426__Linear-CLI/plugins/claude/skills/linear-cli/SKILL.md`
- `joa23__linear-cli/README.md`

### 10. Template-driven issue creation

Useful if Omer wants agents to create high-quality Linear tasks:

- `template list --team TEAM`
- `template get TEMPLATE`
- `issue create --template TEMPLATE --section Name=Value`
- Cache templates for offline/low-latency use.
- Add `--dry-run` to preview rendered description.

This should be additive. Do not make basic `issue create --title` depend on
templates.

Useful evidence:
- `nikpietanze__linear-cli/README.md`
- `Finesssee__linear-cli/README.md`

### 11. Auth, config, and diagnostics

Do not copy other CLIs' unsafe token flows directly. Use the ideas while
preserving `linctl`'s secret rules:

- `doctor` checks config, token presence, Linear reachability, and target match.
- `target --json` should stay the first proof command.
- Support credential helper commands later, but never print tokens.
- Consider named profiles only as explicit overrides, never implicit workspace
  guessing.
- Add `cache status|clear` if schema/name resolution caching arrives.
- Keep `.linctl.toml` scoped to org/team/project, not loose workspace defaults.

Useful evidence:
- `Finesssee__linear-cli/README.md`
- `frr149__lql/README.md`
- `schpet__linear-cli/docs/authentication.md`
- `danielrearden__linear-cli/README.md`

### 12. Packaging and install paths

Leech the distribution ergonomics:

- Homebrew tap/cask after tagged release.
- `go install github.com/KyaniteHQ/linctl/cmd/linctl@latest`.
- GitHub release binaries for Linux/macOS/Windows.
- Static shell completions.
- Optional dynamic completions that query teams/projects/statuses.
- OCI image or wrapper script only if there is actual demand.

Useful evidence:
- `Finesssee__linear-cli/README.md`
- `rubyists__linear-cli/Readme.adoc`
- `joa23__linear-cli/README.md`

## Suggested `linctl` Feature Backlog

### P0: Preserve the moat

- Keep target-pinned writes mandatory.
- Keep schema-aligned GraphQL operations, not ad hoc raw query strings.
- Add tests for every new write proving org/team/project mismatch fails closed.
- Prefer archive over hard delete.

### P1: Agent daily-driver slice

1. Improve `usage` and generated `skills/linctl/SKILL.md`.
2. Add `--fields`, `--compact`, `--id-only`, `--quiet`, `--fail-on-empty`.
3. Add `issue search` plus `issue list` filters: `--state`, `--project`.
4. Add branch helpers around existing `current`: `issue id`, `issue title`,
   `issue url`, `issue branch`, `issue start`, `done`.
5. Add `issue comments`, `issue reply`, and stdin body support.

### P2: Coordination layer

1. Add dependency graph read commands and relation writes.
2. Add project updates and project milestones.
3. Add cycle current/report analytics.
4. Add document create/read/update and issue export.
5. Add template-driven issue create with dry-run.

### P3: Power-user automation

1. Add import/export for CSV/JSON with dry-run.
2. Add watch/monitor commands.
3. Add dynamic completions.
4. Add file upload/download.
5. Extend `next` from dry-run preview into a branch/worktree flow only after the read-only picker stays solid.

## Things Not To Leech Blindly

- Hard delete commands. Prefer archive or require a separate, explicit design.
- Raw GraphQL escape hatches. They undermine schema-aligned safety unless
  isolated under a clearly unsafe debug namespace.
- Implicit workspace switching. `linctl` should keep explicit pinned targets.
- Printing tokens or using `.env` examples that expose secret values.
- Large all-in-one command surfaces before the smaller safe slices are proven.

## Verification Evidence

The cloned repos and primary evidence files:

- `/tmp/linctl-linear-cli-research/schpet__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/schpet__linear-cli/docs/authentication.md`
- `/tmp/linctl-linear-cli-research/schpet__linear-cli/skills/linear-cli/`
- `/tmp/linctl-linear-cli-research/linearis-oss__linearis/README.md`
- `/tmp/linctl-linear-cli-research/joa23__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/Finesssee__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/flipbit03__lineark/README.md`
- `/tmp/linctl-linear-cli-research/frr149__lql/README.md`
- `/tmp/linctl-linear-cli-research/Securiteru__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/juanbermudez__linear-agent-cli/README.md`
- `/tmp/linctl-linear-cli-research/choam2426__Linear-CLI/README.md`
- `/tmp/linctl-linear-cli-research/choam2426__Linear-CLI/plugins/claude/skills/linear-cli/SKILL.md`
- `/tmp/linctl-linear-cli-research/Valian__linear-cli-skill/README.md`
- `/tmp/linctl-linear-cli-research/dabblewriter__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/nikpietanze__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/evangodon__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/rubyists__linear-cli/Readme.adoc`
- `/tmp/linctl-linear-cli-research/AdiKsOnDev__linear-cli/README.md`
- `/tmp/linctl-linear-cli-research/danielrearden__linear-cli/README.md`
