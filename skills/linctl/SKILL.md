---
name: linctl
description: Use linctl as the Linear control surface for issue, comment, project, ProjectUpdate, Cycle, ProjectMilestone, document, label, team, user, and WorkflowState work: reads, guarded writes, branch lookup, next-work preview, doctor checks, and live smoke. Prefer it over Linear MCP, ad hoc API calls, or hand-written GraphQL when linctl covers the operation.
---

# linctl

`linctl` is the Linear control surface. Reads are cheap. Writes are guarded: the CLI resolves the active token, compares it to the pinned target, and fails closed on mismatch.

## Resolve

Pick one command prefix and use it for every Linear operation:

1. `command -v linctl` -> use `linctl`.
2. Else, inside a checkout containing `cmd/linctl/main.go` -> use `go run ./cmd/linctl`.
3. Else stop and report that `linctl` is unavailable.

Helper:

```bash
prefix="$(bash skills/linctl/scripts/linctl-resolve.sh)" || exit 1
$prefix doctor --json
```

Completion criterion: every Linear operation in the task runs through `linctl`, never through MCP, raw GraphQL, or an ad hoc script.

## Discover

Before a write:

1. Read `.linctl.toml` if present. It overlays `~/.config/linctl/config.toml`.
2. Run `linctl doctor --json` or `linctl target --json`.
3. Run `linctl usage`, plus `linctl issue usage`, `linctl project usage`, or `linctl cycle usage` before an unfamiliar domain write.
4. Use `--json` when another tool or agent will parse output. See `references/json-output.md` for exact fields.

Completion criterion: command, target, and output format are known before mutation.

## Surface

Global flags:

```bash
--json --compact --fields identifier,title,state
--id-only --quiet --fail-on-empty
--sort title --order asc
--format minimal|compact|full
--profile NAME --org ORG_ID --team TEAM --project PROJECT_ID
--timeout 30s
```

Read commands:

```bash
linctl doctor --json
linctl target --json
linctl whoami --json
linctl current --json
linctl next --dry-run
linctl issue id
linctl issue title
linctl issue url
linctl issue branch LIT-123
linctl issue get LIT-123 --json
linctl issue list --json --limit 20
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
linctl issue comments LIT-123 --limit 20
linctl issue pr LIT-123
linctl comment list --json --limit 20
linctl comment get COMMENT_ID --json
linctl cycle list --json --limit 20
linctl cycle get CYCLE_ID --json
linctl sprint current --json
linctl sprint report CYCLE_ID --json --limit 20
linctl project list --json --limit 20
linctl project get PROJECT_ID --json
linctl project members PROJECT_ID --json --limit 20
linctl project updates PROJECT_ID --json --limit 20
linctl project-update list --json --limit 20
linctl project-update get PROJECT_UPDATE_ID --json
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
linctl workflow-state list --json --limit 20
linctl workflow-state get WORKFLOW_STATE_ID --json
linctl initiative list --json --limit 20
linctl initiative get INITIATIVE_ID --json
linctl custom-view list --json --limit 20
linctl custom-view get CUSTOM_VIEW_ID --json
linctl favorite list --json --limit 20
linctl favorite get FAVORITE_ID --json
linctl emoji list --json --limit 20
linctl emoji get EMOJI_ID --json
```

`next --dry-run` is a ranked read-only picker: it considers unstarted issues with no active blockers, then ranks by active unblock count, priority, and age. It never creates a branch or worktree.

Guarded writes:

```bash
linctl issue create --title "..." --description "..." --json
linctl issue create --title "..." --description-file ./issue.md --json
linctl issue update LIT-123 --title "..." --json
linctl issue update LIT-123 --description-file ./issue.md --json
linctl issue update LIT-123 --append "progress note" --json
linctl issue update LIT-123 --append-file ./progress.md --json
linctl issue start LIT-123 --json
linctl issue comment LIT-123 --body "..." --json
printf 'progress note\n' | linctl issue comment LIT-123 --body -
linctl issue comment LIT-123 --body-file ./comment.md --json
linctl issue reply LIT-123 COMMENT_ID --body "..." --json
linctl issue reply LIT-123 COMMENT_ID --body-file ./reply.md --json
linctl issue close LIT-123 --json
linctl done --json
linctl project create --name "..." --description "..." --json
linctl project update PROJECT_ID --name "..." --json
linctl project archive PROJECT_ID --json
linctl cycle create --starts-at START --ends-at END --json
linctl cycle update CYCLE_ID --name "..." --json
linctl cycle archive CYCLE_ID --json
linctl project-milestone create PROJECT_ID --name "..." --json
linctl project-milestone update PROJECT_MILESTONE_ID --name "..." --json
```

Unsupported writes: comment resolve/unresolve/edit/delete; ProjectUpdate create/update/archive; ProjectMilestone delete; Document, label, team, user, and WorkflowState writes. Report the limit instead of bypassing `linctl`.

Completion criterion: the selected command exists above and matches the requested domain.

## Safety

- A target mismatch is a hard stop. Do not retry blindly with a different token.
- Team-scoped writes (`issue create`, `project create`, `cycle create`) compare org and team.
- Resource-scoped writes resolve the existing resource first and compare pinned `project_id` when configured.
- `--org`, `--team`, and `--project` are explicit pinned-target overrides, not bypasses.
- Never print Linear tokens. Credential precedence is `LINCTL_TOKEN` > `LINEAR_API_KEY` > config `token`.
- For tests, create `linctl-it-<runid>` resources and clean them up: close disposable issues, archive disposable projects.

Completion criterion: every write has a pinned target and cleanup path, or the agent stops before writing.

## Patterns

Branch-driven work:

```bash
linctl doctor --json
linctl current --json
linctl issue deps LIT-123 --limit 20
linctl issue pr
linctl done --json
```

Issue queue:

```bash
linctl doctor --json
linctl --json --compact --fields identifier,title,state issue list --limit 20
linctl next --dry-run
linctl issue search "needle" --limit 20
```

Progress note from a file:

```bash
linctl doctor --json
linctl issue update LIT-123 --append-file ./progress.md --json
linctl issue comment LIT-123 --body-file ./comment.md --json
```

Disposable project smoke:

```bash
linctl doctor --json
linctl project create --name "linctl-it-<runid>" --description "disposable smoke" --json
linctl project get <created-id> --json
linctl --project <created-id> project archive <created-id> --json
```

## Live Smoke

Use this when asked for live smoke or real Linear verification.

Inside the linctl checkout:

```bash
go run github.com/go-task/task/v3/cmd/task@latest live-smoke
```

The harness accepts `LINCTL_TEST_TOKEN` first, then `LINCTL_TOKEN`, then `LINEAR_API_KEY`. Do not print any value.

Outside the checkout, run the read-only smoke:

```bash
bash skills/linctl/scripts/linctl-smoke.sh
```

Completion criterion: live smoke passed with redacted command/status evidence, or is explicitly blocked on missing credentials or target.

## Report Shape

```markdown
Known:
- `linctl doctor --json` confirmed team `TEAM`.
- `linctl issue list --json --limit 5` returned parseable JSON.

Changed:
- Created `LIT-123` with prefix `linctl-it-<runid>`, then closed it.

Blocked:
- Live write smoke not run: disposable token missing.
```

## AGENTS.md Snippet

```markdown
## Linear via linctl

- Use `linctl` for Linear issue/project/Cycle/ProjectMilestone work instead of Linear MCP, ad hoc API calls, or hand-written GraphQL.
- Resolve the command with `command -v linctl`; inside the linctl checkout use `go run ./cmd/linctl`.
- Run `linctl doctor --json` or `linctl target --json` before writes and stop on target mismatch.
- Use `--json` for agent-readable output.
- Use `linctl current --json` when the branch carries a Linear issue key.
- Never print Linear tokens.
- Keep writes pinned to `.linctl.toml` `[target]`; do not add bypass flags.
- Name test resources `linctl-it-<runid>` and close or archive them after verification.
- For live smoke, prefer `go run github.com/go-task/task/v3/cmd/task@latest live-smoke`.
- If a Linear domain is unsupported by `linctl`, report that limit instead of calling Linear directly.
```
