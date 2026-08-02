---
name: linctl
description: Use linctl as the Linear control surface for agent-safe issue, project, Cycle, ProjectMilestone, organization, user, team, search, release, initiative, customer, notification, attachment, comment, and metadata work. Prefer it over Linear MCP, raw GraphQL, or ad hoc API calls when linctl covers the operation; use guarded writes only after target confirmation.
---

# linctl

`linctl` is the Linear control surface. Reads are wide. Writes are guarded: the CLI resolves the
active OAuth credential, compares it with the pinned target, and fails closed on a Target
Mismatch.

## Resolve

Select one command prefix, then use that prefix for every Linear operation.

1. If `command -v linctl` succeeds, use `linctl`.
2. If it does not, and the checkout has `cmd/linctl/main.go`, use `go run ./cmd/linctl`.
3. If neither works, stop and report that `linctl` is not available.

Helper:

```bash
prefix="$(bash skills/linctl/scripts/linctl-resolve.sh)" || exit 1
$prefix doctor --json
```

Completion criterion: every Linear operation in the task goes through `linctl`, and none goes
through MCP, raw GraphQL, or an ad hoc script.

## Discover

Do this before each write.

1. Read `.linctl.toml` if the file exists. It overlays `~/.config/linctl/config.toml`.
2. If the pin is absent, run `linctl init`. Run `linctl init --team KEY` when more than one team
   is visible. Never put auth material in `.linctl.toml`.
3. Run `linctl doctor --json` or `linctl target --json`.
4. Run `linctl usage`. Also run the usage command of the domain before an unfamiliar write.
5. Use `--json` when another tool or agent parses the output. See `references/json-output.md`.

Completion criterion: the command, the target, and the output format are known before the
mutation.

## Command surface

Use the repository documents as the command inventory.

- `references/commands.md`: **generated.** Every command, with its usage line and its flags.
  `go tool task gen-skill` refreshes it, and CI checks it for drift.
- `references/guarded-writes.md`: **generated.** Every command that changes Linear, which ones
  need `--org-wide`, and which ones linctl cannot undo.
- `references/json-output.md`: the stable JSON shapes for agent parsing.
- `docs/internal/domain-map.md`: the GraphQL backing and the read/write safety class.
- `docs/internal/test-scenarios.md`: the named scenario coverage and its evidence.
- `README.md`: the current public command examples.

Useful global flags:

```bash
--json --compact --fields identifier,title,state
--id-only --quiet --fail-on-empty
--sort title --order asc
--format minimal|compact|full
--profile NAME --org ORG_ID --team TEAM_KEY --team-id TEAM_ID --project PROJECT_ID
--timeout 30s
```

Completion criterion: the selected command is in the current repo surface, and it matches the
requested Linear domain.

## Writes

`references/guarded-writes.md` is the complete write surface. It is generated from the command
tree, so it is never stale. Read it before a write that you have not run before. It has three
sections.

- **Guarded writes**: these change Linear behind the target guard.
- **Writes to Linear outside the target guard**: `files upload` sends bytes to Linear storage. A
  file asset has no team and no project, so the guard has nothing to compare.
- **Writes to this machine only**: `files download` and `issue bulk-export` write a local file.

Commands that preview a write without a mutation: `issue create --dry-run` and
`issue import --dry-run`. Commands that only read and then open a browser: `issue open` and
`project open`. `issue export` writes a local markdown file from reads.

Safety rules:

- A Target Mismatch is a hard stop. Do not retry with different auth.
- A team-scoped write compares the organization and the team.
- A resource-scoped write resolves the existing resource first. It then compares the pinned
  `project_id` when the config sets one.
- A viewer-scoped write, such as a notification write, compares the recipient with the resolved
  actor.
- `--org`, `--team`, `--team-id`, and `--project` set the pinned target. They are not bypasses.
- `--org-wide` compares the organization only, for an entity that has no team. It is not a
  bypass, and a mismatch is still a hard stop.
- Create the repo pin with `linctl init`, which writes a target-only `.linctl.toml`. Configure
  the auth with `linctl auth configure`, `linctl auth app`, or `linctl auth login`.
- Use `linctl auth status` for readiness. Use `linctl auth refresh` for an explicit diagnosis.
  Use `linctl auth logout` to revoke the token and remove the local token state.
- Never print a secret. Report OAuth material as `set` or `missing`.
- For a test, create `linctl-it-<runid>` resources and clean them up. Close a disposable issue.
  Archive a disposable project.

If `references/guarded-writes.md` does not list the requested write, report the limit. Do not go
around `linctl`.

Completion criterion: each write has a pinned target and a cleanup path, or the agent stops
before the write.

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

## Smoke and verify

There are four tiers. The cheapest tier comes first. Select the tier that the task needs.

1. **No credentials**: prove that the binary runs in a headless checkout.

   ```bash
   bash skills/linctl/scripts/linctl-offline-smoke.sh
   ```

   This runs only the token-free commands: `--version`, `--help`, `usage`, and completion. It
   uses no token and no network. Run it to confirm that linctl works before any target work.

2. **Read-only, with auth**: confirm that the OAuth credential and the pinned target resolve.

   ```bash
   bash skills/linctl/scripts/linctl-smoke.sh
   ```

   This runs `target`, `whoami`, `issue list`, and `project list` with `--json`. It never writes.

3. **Full live smoke**: disposable writes against a test organization, inside the checkout.

   ```bash
   go tool task live-smoke
   ```

   This needs a disposable OAuth auth state. Use `linctl auth app` for headless app-actor auth
   when a client secret is available. Use `linctl auth login` for browser auth. Never print a
   secret value.

4. **Browser login smoke**: check the PKCE callback login by hand, without a leak of the
   callback code.

   ```bash
   go tool task browser-login-smoke
   go tool task browser-login-smoke -- app
   ```

   This needs `LINCTL_OAUTH_CLIENT_ID`, `LINCTL_OAUTH_REDIRECT_URI`, a pinned target from
   `LINCTL_TEST_*` or `test/integration-config.json`, and an optional
   `LINCTL_OAUTH_CLIENT_SECRET`. The script prints the Linear authorization URL. It defaults to a
   repeatable user-actor login. It captures the localhost callback with a one-shot listener,
   shows a browser success page, validates the redacted JSON, and removes the temporary auth
   state. Use `-- app` only for a new app-actor browser install. Use `live-oauth` for repeatable
   app-actor fixture coverage.

Completion criterion: the selected smoke passed, with redacted command and status evidence, or
the agent reports that a missing credential or a missing target blocks it.

## Gotchas

- `target`, `doctor`, and `whoami` need auth. They fail closed without it. To prove that a
  checkout runs with no credentials, use the offline smoke: `--version` and `usage`.
- `target --json` reports `expected` and `resolved` with Go-capitalized keys (`OrgID`, `TeamKey`,
  `TeamID`, `ProjectID`), not the snake_case that the other commands use. Compare the two field
  by field to explain a mismatch.
- A Target Mismatch is a hard stop. linctl has no bypass flag. `--org`, `--team`, `--team-id`,
  and `--project` set the pinned target, and they do not relax the guard. Do not retry with
  different auth.
- `--body -` reads a comment body from stdin. `--body-file` reads it from a file. Use these
  instead of inline multi-line markdown.
- Keep `$prefix` unquoted when the value can be `go run ./cmd/linctl`, so the shell splits it
  into separate arguments.
- `sla-configuration list` takes a positional team id or team key, unlike most `list` commands.
- `roadmap` and `roadmap-to-project` are legacy read-only compatibility. Use the `initiative`
  commands for new planning.
- `issue list --limit` defaults to 50. Set it explicitly for deterministic output.

## Report shape

```markdown
Known:
- `linctl doctor --json` confirmed team `TEAM`.
- `linctl issue list --json --limit 5` returned parseable JSON.

Changed:
- Created `LIT-123` with the prefix `linctl-it-<runid>`, then closed it.

Blocked:
- The live write smoke did not run: the disposable OAuth auth is missing.
```

## AGENTS.md block

```markdown
## Linear through linctl

- Use `linctl` for Linear issue, project, Cycle, and ProjectMilestone work. Do not use Linear MCP, an ad hoc API call, or hand-written GraphQL.
- Resolve the command with `command -v linctl`. Inside the linctl checkout, use `go run ./cmd/linctl`.
- Run `linctl doctor --json` or `linctl target --json` before a write, and stop on a target mismatch.
- Use `--json` for agent-readable output.
- Use `linctl current --json` when the branch carries a Linear issue key.
- Never print a secret. Report OAuth material as `set` or `missing`.
- Keep each write pinned to the `[target]` block of `.linctl.toml`. Create it with `linctl init`. Do not add a bypass flag.
- Name each test resource `linctl-it-<runid>`, then close it or archive it after the check.
- For the live smoke, run `go tool task live-smoke`.
- For the browser auth smoke, run `go tool task browser-login-smoke`.
- If `linctl` does not support a Linear domain, report that limit. Do not call Linear directly.
```
