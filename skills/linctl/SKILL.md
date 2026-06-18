---
name: linctl
description: Use the linctl CLI for schema-aligned Linear reads and fail-closed writes from agent sessions. Load when a task asks to inspect, create, update, comment on, close, archive, or coordinate Linear issues or projects through linctl.
---

# linctl

Use `linctl` as the Linear control surface when an agent needs a compact, target-pinned CLI instead of an MCP connector.

## Discover First

1. Read the repo's `.linctl.toml` if it exists. Confirm `[target]` has `org_id`, `team_key`, `team_id`, and optional `project_id`.
2. Run `linctl target --json` before any write. Treat a target mismatch as a hard stop.
3. Run `linctl usage`, then the domain-specific help such as `linctl issue usage` or `linctl project usage`.
4. Prefer `--json` whenever another tool or agent will parse the result.

## Act Safely

- Reads may run without a pinned target.
- Writes require pinned `org_id` plus team identity and fail closed on mismatch.
- Team-scoped writes create a new resource and compare org/team only.
- Resource-scoped writes resolve the resource first and compare the pinned `project_id` when one is configured.
- Use namespaced throwaway resources for tests, then archive them.
- Never paste or print Linear tokens. Inject them through `LINCTL_TOKEN` or the local secret manager.

## Common Commands

```bash
linctl target --json
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

## AGENTS.md Snippet

```markdown
## Linear via linctl

- Use `linctl target --json` before Linear writes and stop on target mismatch.
- Use `linctl usage` and `<domain> usage` before unfamiliar commands.
- Use `--json` for agent-readable output.
- Never print Linear tokens; inject credentials via the configured secret manager.
- Writes must stay pinned to `.linctl.toml` `[target]`; do not add bypass flags.
- Test resources must use a namespaced prefix and be archived after verification.
```
