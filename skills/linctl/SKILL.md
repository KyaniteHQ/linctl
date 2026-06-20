---
name: linctl
description: Use the linctl Linear CLI as the control surface for Linear issue and project work from agent sessions — reads, target-pinned writes, branch-issue lookup, and live smoke checks. Load this whenever a task touches Linear issues or projects: listing or getting issues/projects, resolving the current branch's issue (e.g. an identifier like LIT-123), creating/updating/commenting/closing issues, creating/updating/archiving projects, or running a live Linear smoke check — even when the user only says "check the ticket", "update the Linear issue", or names an issue key without mentioning linctl. Prefer it over the Linear MCP server, ad hoc Linear API calls, or hand-written GraphQL for these operations.
---

# linctl

Use `linctl` as the Linear control surface for issue and project reads, target-pinned writes, branch issue lookup, and live smoke verification. Run only the commands a task needs — this is an à la carte surface, not a battery to run end to end. Do not hand-write Linear API or GraphQL calls, and do not reach for the Linear MCP server, when `linctl` already covers the operation.

## Resolve the CLI

Resolve one command prefix and route every Linear operation through it:

1. `command -v linctl` → use `linctl`.
2. Else, inside a linctl checkout (`cmd/linctl/main.go` present) → use `go run ./cmd/linctl`.
3. Else stop and report that `linctl` is unavailable.

`scripts/linctl-resolve.sh` does exactly this and prints the prefix (left unquoted so `go run ./cmd/linctl` splits into words):

```bash
prefix="$(bash skills/linctl/scripts/linctl-resolve.sh)" || exit 1
$prefix target --json
```

Completion criterion: every Linear operation runs through a known `linctl` binary, never an ad hoc API script.

## Discover first

1. Read `.linctl.toml` in the repo (it overlays the global `~/.config/linctl/config.toml`). Confirm `[target]` carries `org_id`, `team_key`, `team_id`, and optional `project_id`.
2. Run `linctl usage` (and `linctl issue usage` / `linctl project usage`) before an unfamiliar write.
3. Run `linctl target --json` before any write — a target mismatch is a hard stop, not a warning.
4. Pass `--json` whenever another tool, agent, or test will parse the output. See `references/json-output.md` for the exact field names of every command's JSON.

Completion criterion: target, command, and output format are settled before the first write.

## CLI surface

Global flags: `--json` (structured output), `--compact` (one-line JSON), `--fields identifier,title,state` (JSON projection), `--id-only` (emit ids for chaining), `--quiet` (suppress successful output), `--fail-on-empty` (non-zero empty list), `--sort field --order asc|desc` (deterministic list output), `--format minimal|compact|full` (human output), `--profile` (named config profile), `--org` / `--team` / `--project` (explicit pinned-target overrides), `--timeout` (request bound).

Reads — no pinned target required:

```bash
linctl target --json                          # resolved org/team/project for the active token
linctl whoami --json                          # authenticated user
linctl current --json                         # issue keyed by the current branch / jj description
linctl next --dry-run                         # preview first unstarted issue with no blockers
linctl done                                   # close the current checkout issue
linctl issue id                               # current branch issue identifier, e.g. LIT-123
linctl issue title                            # current branch issue title
linctl issue url                              # current branch Linear URL
linctl issue branch LIT-123                   # Linear's suggested branch name, no checkout
linctl issue deps LIT-123 --limit 20          # parent, children, blocked issues, and blockers
linctl issue pr LIT-123                       # print gh pr create title/body plan
linctl issue comments LIT-123 --limit 20      # issue discussion comments
linctl issue list --json --limit 20
linctl issue list --state started --limit 20   # workflow state type
linctl issue list --project PROJECT_ID --limit 20
linctl issue list --mine --limit 20             # assigned to the authenticated user
linctl issue list --assignee USER_ID --limit 20 # assigned to a Linear user id
linctl issue list --label LABEL_ID --limit 20   # issues with a Linear label id
linctl issue list --cycle CYCLE_ID --limit 20   # issues attached to a Linear Cycle id
linctl issue list --created-after 2026-06-01 --limit 20 # created on or after a date
linctl issue list --created-since 2026-06-01 --limit 20 # alias for created-after
linctl issue list --created-before 2026-06-30 --limit 20 # created on or before a date
linctl issue list --has-blockers --limit 20     # issues blocked by another issue
linctl issue list --blocks --limit 20           # issues blocking another issue
linctl issue list --blocked-by LIT-123 --limit 20 # issues blocked by that issue
linctl issue list --all-teams --limit 20        # broad read-only issue inspection
linctl issue search "needle" --limit 20        # resolved-team text search
linctl issue deps LIT-123 --limit 20           # dependency graph for one issue
linctl issue get LIT-123 --json
linctl --json --compact --fields identifier,title,state issue list --limit 20
linctl --id-only issue get LIT-123
linctl --fail-on-empty --sort title --order asc issue list --limit 20
linctl cycle list --json --limit 20
linctl cycle get CYCLE_ID --json
linctl sprint current --json
linctl sprint report CYCLE_ID --json --limit 20
linctl project list --json --limit 20
linctl project get PROJECT_ID --json
linctl project members PROJECT_ID --json --limit 20
linctl project updates PROJECT_ID --json --limit 20
linctl project-milestone list PROJECT_ID --json --limit 20
linctl project-milestone get PROJECT_MILESTONE_ID --json
linctl document list --json --limit 20
linctl document get DOCUMENT_ID --json
linctl label list --json --limit 20
linctl label get LABEL_ID --json
linctl team list --json --limit 20
linctl team get TEAM_ID --json
linctl team members TEAM_ID --json --limit 20
linctl user list --json --limit 20
linctl user get USER_ID --json
linctl user me --json
```

Writes — guarded, require a pinned target, fail closed on mismatch:

```bash
linctl issue create --title "..." --description "..." --json    # team-scoped
linctl issue update LIT-123 --title "..." --json                # resource-scoped
linctl issue update LIT-123 --append "progress note" --json     # append to description
linctl issue start LIT-123 --json                               # assign to viewer and move to started
linctl done --json                                              # close the current checkout issue
linctl issue comment LIT-123 --body "..." --json                # resource-scoped
linctl issue reply LIT-123 comment-id --body "..." --json       # resource-scoped
printf 'progress note\n' | linctl issue comment LIT-123 --body - # stdin body
linctl issue close LIT-123 --json                               # resource-scoped
linctl project create --name "..." --description "..." --json   # team-scoped
linctl project update PROJECT_ID --name "..." --json            # resource-scoped
linctl project archive PROJECT_ID --json                        # resource-scoped (cleanup path)
linctl cycle create --starts-at START --ends-at END --json      # team-scoped
linctl cycle update CYCLE_ID --name "..." --json                # team-scoped
linctl cycle archive CYCLE_ID --json                            # team-scoped
linctl project-milestone create PROJECT_ID --name "..." --json  # resource-scoped project write
linctl project-milestone update PROJECT_MILESTONE_ID --name "..." --json # resource-scoped
```

That is the whole implemented surface. ProjectMilestone delete and Document, label, team, or user writes are unsupported by `linctl`; confirm with `linctl --help` or `linctl usage`, then report the limit. Do not silently fall back to GraphQL or the Linear MCP server.

Completion criterion: the chosen command exists in the surface above and matches the requested domain.

## Act safely

- Reads may run without a pinned target. Writes require `org_id` plus team identity (in `.linctl.toml` or via `--org`/`--team`/`--project`) and fail closed on mismatch.
- Team-scoped writes (`issue create`, `project create`) compare org + team only — the resource does not exist yet.
- Resource-scoped writes resolve the existing resource first and compare the pinned `project_id` when one is configured, so a same-team wrong-project write is refused.
- On a target mismatch, inspect the expected vs resolved ids (see the `target --json` shape in `references/json-output.md`). Do not retry blindly with a different token.
- Never paste or print a Linear token. The CLI reads credentials from `LINCTL_TOKEN`, then `LINEAR_API_KEY`, then a config `token` — values stay out of chat and logs.
- Use namespaced throwaway resources for any write check (`linctl-it-<runid>`), then clean up: close disposable issues, archive disposable projects. Archive is the only cleanup path — there is no hard delete.

Completion criterion: every write has a pinned target and a cleanup path, or the agent stops before writing.

## Live smoke

Use this when a task asks for live smoke, production verification, or confidence against the real Linear service.

- Inside this checkout, prefer the full harness: `task live-smoke`. It builds a temporary binary, runs read-only checks, then disposable write round-trips with cleanup. It takes a disposable token via `LINCTL_TEST_TOKEN` (preferred) `>` `LINCTL_TOKEN` `>` `LINEAR_API_KEY` and forwards it to the CLI — `LINCTL_TEST_TOKEN` is a harness input, not a credential the CLI itself reads.
- Anywhere else (installed binary, no checkout), run the portable read-only smoke first:

```bash
bash skills/linctl/scripts/linctl-smoke.sh    # target, whoami, issue list, project list — read-only
```

- For a write smoke, create one namespaced resource, verify it with a read, then clean it up. If the repo target pins a different fixture project, archive a newly created disposable project with `--project <created-id>` — that still goes through target comparison, so it is not a bypass.
- Record command names, exit status, and redacted evidence only.

Completion criterion: live smoke passed with redacted evidence, or is explicitly blocked on missing credentials/target.

## Command patterns

Branch-driven issue lookup:

```bash
linctl target --json
linctl current --json        # LIT-123 from the branch or jj description, then the issue
linctl done                  # close the current branch or jj issue
linctl issue id              # just LIT-123 from the current branch or jj description
linctl issue title           # just the current issue title
linctl issue url             # just the current issue URL
linctl issue branch LIT-123  # branch slug from Linear, without checkout
linctl issue deps LIT-123 --limit 20
linctl issue pr              # gh pr create plan from the current issue
linctl next --dry-run        # preview first unstarted issue with no blockers
```

Single issue write:

```bash
linctl target --json
linctl issue list --json --limit 20    # confirm the visible queue
linctl issue list --state started --limit 20
linctl issue list --project PROJECT_ID --limit 20
linctl issue list --mine --limit 20
linctl issue list --assignee USER_ID --limit 20
linctl issue list --label LABEL_ID --limit 20
linctl issue list --cycle CYCLE_ID --limit 20
linctl issue list --created-after 2026-06-01 --limit 20
linctl issue list --created-since 2026-06-01 --limit 20
linctl issue list --created-before 2026-06-30 --limit 20
linctl issue list --has-blockers --limit 20
linctl issue list --blocks --limit 20
linctl issue list --blocked-by LIT-123 --limit 20
linctl issue list --all-teams --limit 20
linctl issue search "needle" --limit 20
linctl issue deps LIT-123 --limit 20
linctl issue pr LIT-123
linctl issue comments LIT-123 --limit 20
linctl next --dry-run
linctl cycle get CYCLE_ID --json
linctl sprint current --json
linctl sprint report CYCLE_ID --json --limit 20
linctl document list --json --limit 20
linctl label list --json --limit 20
linctl team members TEAM_ID --json --limit 20
linctl user me --json
linctl issue update LIT-123 --append "progress note" --json
linctl issue start LIT-123 --json
linctl done --json
linctl issue comment LIT-123 --body "concise update" --json
linctl issue reply LIT-123 comment-id --body "thread reply" --json
printf 'longer update\n' | linctl issue comment LIT-123 --body -
```

Machine-readable issue queue:

```bash
linctl --json --compact --fields identifier,title,state issue list --limit 20
linctl issue list --state started --limit 20
linctl issue list --project PROJECT_ID --limit 20
linctl issue list --mine --limit 20
linctl issue list --assignee USER_ID --limit 20
linctl issue list --label LABEL_ID --limit 20
linctl issue list --cycle CYCLE_ID --limit 20
linctl issue list --created-after 2026-06-01 --limit 20
linctl issue list --created-since 2026-06-01 --limit 20
linctl issue list --created-before 2026-06-30 --limit 20
linctl issue list --has-blockers --limit 20
linctl issue list --blocks --limit 20
linctl issue list --blocked-by LIT-123 --limit 20
linctl issue list --all-teams --limit 20
linctl issue search "needle" --limit 20
linctl issue deps LIT-123 --limit 20
linctl issue pr LIT-123
linctl next --dry-run
linctl cycle list --json --limit 20
linctl document list --json --limit 20
linctl label list --json --limit 20
linctl team list --json --limit 20
linctl user list --json --limit 20
linctl --fail-on-empty --sort title --order asc issue list --limit 20
```

Disposable project smoke with cleanup:

```bash
linctl target --json
linctl project create --name "linctl-it-<runid>" --description "disposable smoke" --json
linctl project get <created-id> --json
linctl --project <created-id> project archive <created-id> --json
```

Project status history:

```bash
linctl target --json
linctl project updates PROJECT_ID --json --limit 20
linctl --json --fields id,health,display_name project updates PROJECT_ID --limit 20
```

Project milestone list:

```bash
linctl target --json
linctl project-milestone list PROJECT_ID --json --limit 20
linctl --json --fields id,name,status project-milestone list PROJECT_ID --limit 20
linctl project-milestone create PROJECT_ID --name "linctl-it-<runid>" --json
linctl project-milestone update PROJECT_MILESTONE_ID --name "renamed" --json
```

## Evidence pattern

Report in three sections so the reader sees state, change, and gaps at a glance:

```markdown
Known:
- `linctl target --json` confirmed team `TEAM`.
- `linctl issue list --json --limit 5` returned parseable JSON.

Changed:
- Created `LIT-123` (title prefix `linctl-it-2026-06-19`), then closed it after verification.

Blocked:
- Live write smoke not run: `LINCTL_TOKEN` / `LINEAR_API_KEY` unset.
```

## AGENTS.md snippet

Paste into a project's AGENTS.md to keep agents on the CLI:

```markdown
## Linear via linctl

- Use the `linctl` Linear CLI for Linear issue/project work instead of the Linear MCP server, ad hoc API calls, or hand-written GraphQL.
- Resolve the command with `command -v linctl`; inside the linctl checkout use `go run ./cmd/linctl`.
- Run `linctl target --json` before any write and stop on target mismatch.
- Run `linctl usage` / `<domain> usage` before unfamiliar commands; pass `--json` for agent-readable output.
- Run `linctl current --json` when the branch carries a Linear issue key.
- Never print Linear tokens. The CLI reads `LINCTL_TOKEN` > `LINEAR_API_KEY` > config `token`; the live-smoke harness also accepts `LINCTL_TEST_TOKEN`.
- Keep writes pinned to `.linctl.toml` `[target]`; do not add bypass flags.
- Name test resources `linctl-it-<runid>` and close or archive them after verification.
- For live smoke prefer `task live-smoke`; otherwise run the read-only `linctl-smoke.sh` first, then disposable write checks with cleanup.
- If a Linear domain is unsupported by `linctl`, report that limit instead of calling Linear directly.
```
