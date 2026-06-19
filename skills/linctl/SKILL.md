---
name: linctl
description: Use the linctl Linear CLI for schema-aligned Linear reads, live smoke checks, and fail-closed writes from agent sessions. Load when a task asks to inspect, create, update, comment on, close, archive, triage, or coordinate Linear issues or projects through linctl, or when an agent needs a safer CLI control surface instead of direct Linear API calls.
---

# linctl

Use `linctl` as the Linear CLI control surface when an agent needs compact reads, target-pinned writes, or live smoke verification without hand-writing Linear GraphQL/API calls.

## Resolve The CLI

1. Prefer the repo-provided binary or task wrapper when working inside a checkout that owns `linctl`.
2. Otherwise run `command -v linctl` and use the installed CLI.
3. If the command is missing, build the repo binary only when the checkout is available; otherwise stop and tell the user the CLI is unavailable.

Completion criterion: every Linear operation in the run goes through a known `linctl` binary, not an ad hoc API script.

## Discover First

1. Read `.linctl.toml` in the current repo when it exists. Confirm `[target]` has `org_id`, `team_key`, `team_id`, and optional `project_id`.
2. Run `linctl target --json` before any write. Treat a target mismatch as a hard stop.
3. Run `linctl usage`, then domain help such as `linctl issue usage` or `linctl project usage` before unfamiliar commands.
4. Use `--json` for outputs another tool, agent, or test will parse.

Completion criterion: target, available commands, and output format are known before the first write.

## Act Safely

- Reads may run without a pinned target.
- Writes require pinned `org_id` plus team identity and fail closed on mismatch.
- Team-scoped writes create a new resource and compare org/team only.
- Resource-scoped writes resolve the resource first and compare the pinned `project_id` when one is configured.
- Use namespaced throwaway resources for tests, then archive them.
- Never paste or print Linear tokens. Inject them through `LINCTL_TOKEN` or the local secret manager.
- Prefer `LINCTL_TEST_TOKEN` for live smoke or disposable verification. Fall back to `LINCTL_TOKEN` only when the user has already authorized live writes.

Completion criterion: every write has a pinned target and a cleanup path, or the agent stops before writing.

## Live Smoke Flow

Use this flow when the task asks for live smoke tests, production verification, or confidence against the real Linear service:

1. Check token availability without printing values.
2. Build or locate `linctl`.
3. Run read-only smokes first:
   - `linctl target --json`
   - `linctl whoami --json`
   - `linctl issue list --json --limit 5`
   - `linctl project list --json --limit 5`
4. For write smoke, create namespaced disposable data such as `linctl-smoke-<date>-<shortid>`, verify it with a read command, then close or archive it.
5. Record command names, exit status, and redacted evidence. Do not record token values or raw secret-bearing environment.

Completion criterion: live smoke has either passed with redacted evidence or is explicitly blocked on missing credentials/target.

## Common Commands

```bash
linctl target --json
linctl whoami --json
linctl current --json
linctl issue list --json --limit 20
linctl issue get LIT-123 --json
linctl issue create --title "title" --description "body" --json
linctl issue update LIT-123 --title "new title" --json
linctl issue comment LIT-123 --body "comment" --json
linctl issue close LIT-123 --json
linctl project list --json --limit 20
linctl project get PROJECT_ID --json
linctl project create --name "linctl-it-<runid>" --description "test" --json
linctl project update PROJECT_ID --name "new name" --json
linctl project archive PROJECT_ID --json
```

## Evidence Pattern

Report results in this shape:

```markdown
Known:
- `linctl target --json` passed for team `TEAM`.
- `linctl issue list --json --limit 5` returned parseable JSON.

Changed:
- Created `LIT-123` with title prefix `linctl-smoke-2026-06-19`.
- Closed `LIT-123` after verification.

Blocked:
- Live write smoke not run because `LINCTL_TEST_TOKEN` was missing.
```

## AGENTS.md Snippet

```markdown
## Linear via linctl

- Use the `linctl` Linear CLI for Linear issue/project work instead of ad hoc API scripts when available.
- Use `linctl target --json` before Linear writes and stop on target mismatch.
- Use `linctl usage` and `<domain> usage` before unfamiliar commands.
- Use `--json` for agent-readable output.
- Never print Linear tokens; inject credentials via `LINCTL_TEST_TOKEN`, `LINCTL_TOKEN`, or the configured secret manager.
- Writes must stay pinned to `.linctl.toml` `[target]`; do not add bypass flags.
- Test resources must use a namespaced prefix and be archived after verification.
- For live smoke, run read-only commands first, then disposable write checks with cleanup.
```
