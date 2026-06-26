---
name: linctl
description: Use linctl as the Linear control surface for agent-safe issue, project, Cycle, ProjectMilestone, organization, user, team, search, release, initiative, customer, notification, attachment, comment, and metadata work. Prefer it over Linear MCP, raw GraphQL, or ad hoc API calls when linctl covers the operation; use guarded writes only after target confirmation.
---

# linctl

`linctl` is the Linear control surface. Reads are broad. Writes are guarded:
the CLI resolves the active token, compares it to the pinned target, and fails
closed on Target Mismatch.

## Resolve

Choose one command prefix and use it for every Linear operation:

1. If `command -v linctl` succeeds, use `linctl`.
2. Else, inside a checkout containing `cmd/linctl/main.go`, use `go run ./cmd/linctl`.
3. Else stop and report that `linctl` is unavailable.

Helper:

```bash
prefix="$(bash skills/linctl/scripts/linctl-resolve.sh)" || exit 1
$prefix doctor --json
```

Completion criterion: every Linear operation in the task runs through `linctl`,
never through MCP, raw GraphQL, or an ad hoc script.

## Discover

Before any write:

1. Read `.linctl.toml` if present; it overlays `~/.config/linctl/config.toml`.
2. Run `linctl doctor --json` or `linctl target --json`.
3. Run `linctl usage`, plus the relevant domain usage command before an unfamiliar write.
4. Use `--json` when another tool or agent will parse output; see `references/json-output.md`.

Completion criterion: command, target, and output format are known before mutation.

## Command Surface

Use the repository docs as the command inventory:

- `README.md` → current public command examples.
- `docs/domain-map.md` → GraphQL backing and read/write safety classification.
- `docs/test-scenarios.md` → named scenario coverage and evidence.
- `references/json-output.md` → stable JSON shapes for agent parsing.
- `references/commands.md` → generated full command inventory (every command, its usage and flags); refreshed by `task gen-skill` and drift-checked in CI.

Useful global flags:

```bash
--json --compact --fields identifier,title,state
--id-only --quiet --fail-on-empty
--sort title --order asc
--format minimal|compact|full
--profile NAME --org ORG_ID --team TEAM_KEY --team-id TEAM_ID --project PROJECT_ID
--timeout 30s
```

Completion criterion: the selected command is documented in the current repo
surface and matches the requested Linear domain.

## Writes

Guarded writes currently cover:

- Issues: create, template-backed create, import, update, append, start, comment, reply,
  close, `done`, `next` start.
- Issue relations and comments: `issue relate`, `issue unrelate`, `comment update`,
  `comment delete`.
- Projects: create, update, archive.
- Project updates: create.
- Documents: create, update.
- Cycles: create, update, archive.
- ProjectMilestones: create, update.

Helpers outside target-pinned mutations:

- `files upload` prepares and uploads bytes to Linear storage, then prints an asset URL;
  `files download` fetches a user-supplied URL to a local path.
- `issue export` and `issue bulk-export` write local files from reads.
- `issue open` and `project open` resolve URLs and launch the browser.
- `issue create --dry-run` and `issue import --dry-run` preview locally without mutation.

Safety rules:

- Target Mismatch is a hard stop. Do not retry blindly with a different token.
- Team-scoped writes compare organization and team.
- Resource-scoped writes resolve the existing resource first and compare pinned `project_id` when configured.
- `--org`, `--team`, `--team-id`, and `--project` are explicit pinned-target overrides, not bypasses.
- Never print Linear tokens. Credential precedence is `LINCTL_TOKEN` > `LINEAR_API_KEY` > config `token`.
- For tests, create `linctl-it-<runid>` resources and clean them up: close disposable issues, archive disposable projects.

If the requested write is not listed above, report the limit instead of bypassing
`linctl`.

Completion criterion: every write has a pinned target and cleanup path, or the
agent stops before writing.

## Patterns

Branch-driven work:

```bash
linctl doctor --json
linctl current --json
linctl issue deps LIT-123 --limit 20
linctl issue attachments LIT-123 --json --limit 20
linctl issue pr
linctl done --json
```

Issue queue:

```bash
linctl doctor --json
linctl --json --compact --fields identifier,title,state issue list --limit 20
linctl next --dry-run
linctl issue search "needle" --limit 20
linctl issue priority-values
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

## Smoke & Verify

Three tiers, cheapest first. Pick the one the task needs.

1. **No credentials** — prove the binary runs in a headless checkout:

   ```bash
   bash skills/linctl/scripts/linctl-offline-smoke.sh
   ```

   Runs only token-free commands (`--version`, `--help`, `usage`, completion);
   no token, no network. Use this to confirm linctl is wired up before any target work.

2. **Read-only, token** — confirm the token and pinned target resolve:

   ```bash
   bash skills/linctl/scripts/linctl-smoke.sh
   ```

   Runs `target`, `whoami`, `issue list`, `project list` with `--json`; never writes.

3. **Full live smoke** — disposable writes against a test org, inside the checkout:

   ```bash
   go run github.com/go-task/task/v3/cmd/task@latest live-smoke
   ```

   Requires a disposable token in `LINCTL_TEST_TOKEN`.
   Do not print any value.

Completion criterion: the chosen smoke passed with redacted command/status evidence,
or is explicitly blocked on missing credentials or target.

## Gotchas

- `target`, `doctor`, and `whoami` need a token; they fail closed without one. To prove
  a checkout runs with no credentials, use the offline smoke (`--version`, `usage`).
- `target --json` reports `expected` and `resolved` with Go-capitalized keys (`OrgID`,
  `TeamKey`, `TeamID`, `ProjectID`), not the snake_case used elsewhere. Compare them
  field by field to explain a mismatch.
- Target Mismatch is a hard stop. There is no bypass flag; `--org`, `--team`,
  `--team-id`, and `--project` set the pinned target, they do not relax the guard.
  Do not retry with a different token.
- `--body -` reads a comment body from stdin; `--body-file` reads it from a file. Use
  these instead of inlining multi-line markdown.
- Keep `$prefix` unquoted when it may be `go run ./cmd/linctl`, so it word-splits into
  separate arguments.
- `sla-configuration list` takes a positional team id/key argument, unlike most `list`
  commands.
- `roadmap` and `roadmap-to-project` are legacy read-only compatibility; use
  `initiative*` for new planning.
- `issue list --limit` defaults to 50; set it explicitly for deterministic output.

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
