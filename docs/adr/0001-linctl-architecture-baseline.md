# ADR 0001: linctl architecture baseline

**Status**: accepted baseline

**Consolidates**: the previous ADR set from `0001` through `0025`. This file is the only current ADR. It keeps the target-pinned write invariant from the former `0001` and absorbs the later collection-projection and OAuth decisions into one codebase-aligned baseline.

## Context

linctl is an agent-first Linear control surface. Its core product promise is narrow:

- Reads should be broad, cheap, and easy to run from an agent shell.
- Writes should fail closed unless the active OAuth credential proves it can reach the repo's Pinned Target.
- Output should be deterministic enough for shell pipelines and LLM agents.
- Auth should use Linear OAuth app credentials and renewable OAuth token state, not checked-in or repo-local credential material.

That makes the write guard, OAuth runtime, command output policy, and command-to-client seams part of one design. Splitting those decisions across many tiny ADRs made the record hard to use, especially after OAuth moved from proposal to implementation.

## Decision

### 1. Reads Are Broad, Writes Are Guarded

Read commands may inspect Linear data visible to the active OAuth credential without requiring a Pinned Target first. Guarded Writes must resolve the active credential to a Resolved Target at command time and compare it to the repo's Pinned Target before mutating Linear.

The Pinned Target is `org_id`, `team_key`, `team_id`, and optional `project_id` from config. The Resolved Target is the organization, team, and optional project proven through Linear GraphQL using the active OAuth access token.

Target Mismatch is a hard stop. It is not a warning, confirmation prompt, or recoverable branch. `--org`, `--team`, `--team-id`, and `--project` set the target used for comparison; they do not relax the guard.

Team-scoped creates compare organization and team because the entity does not exist yet. Resource-scoped writes resolve the existing resource first and compare its team, plus `project_id` when one is pinned.

### 2. OAuth App Credentials Are Product Auth

The product auth model is Linear OAuth. linctl supports OAuth app configuration, app-actor client credentials, and browser authorization-code login. Legacy personal-token material is not a product path, and process env overrides must not reintroduce it indirectly.

Environment OAuth variables are non-persistent automation overrides. They can supply OAuth app material or an OAuth access token for a process, but saved local auth state remains the default interactive path.

### 3. Auth State Is Local, Profile-Scoped, And Split By Role

Repo config stays shareable. It owns the Pinned Target, not OAuth secrets or token state.

OAuth app configuration and OAuth token state live in OS-native user config/state locations. They are profile-scoped with the same profile selection used by config loading. App configuration and token state remain logically separate so commands can clear one without rewriting the other.

`auth logout` removes token state and keeps app configuration by default. `--forget-app` is the explicit operation that removes saved app configuration too.

### 4. Auth Commands Are The Public Auth Interface

The stable auth command surface is:

- `linctl auth configure`: save OAuth app configuration.
- `linctl auth login`: run browser authorization-code login.
- `linctl auth app`: authorize non-interactively with client credentials as the app actor.
- `linctl auth status`: inspect token, actor, scopes, expiry, and target readiness.
- `linctl auth refresh`: explicitly refresh or reacquire token state.
- `linctl auth logout`: revoke tokens when Linear accepts revocation, then clear local token state.

The app actor is the default OAuth actor. Browser login sends `actor=app` unless the user passes `--actor user`. User attribution remains available, but it is explicit.

Setup and status commands must prove readiness before reporting success. Saving token bytes is not enough. linctl must verify the expected actor when available from the readiness path, required scopes from token state, and the Resolved Target before treating auth as usable.

### 5. Browser Login Uses PKCE And Manual Callback Fallback

Authorization-code login always uses PKCE with an S256 challenge. The CLI prints or returns a Linear authorization URL and then accepts either a callback URL or raw authorization code through `--callback`; `--callback -` reads the callback from stdin.

Ordinary commands do not auto-launch browser reauthorization when scopes are missing. They fail with a structured missing-scope error and show the reauthorization path. Scope escalation stays deliberate and scriptable.

### 6. Token Recovery Is Bounded

The runtime may recover OAuth token state once for the original GraphQL request:

- Authorization-code tokens refresh through their refresh token.
- Client-credentials app tokens are reacquired through the app credential.
- Rotated token state is persisted only for persistent local auth sessions.
- A request rejected after recovery returns an auth error instead of looping.

Client-credentials tokens are cached until expiry or authorization failure. linctl avoids per-command token churn while still recovering from expired or invalidated app tokens.

### 7. Scopes And Error Codes Are Stable

The default OAuth scope set is:

```text
read,write,issues:create,comments:create
```

`admin` is not requested by default. Commands that need more permission fail with `MISSING_SCOPE` and tell the user how to reauthorize.

OAuth failures use stable structured codes:

```text
AUTH_NOT_CONFIGURED
AUTH_TOKEN_EXPIRED
AUTH_REFRESH_FAILED
AUTH_REAUTH_REQUIRED
MISSING_SCOPE
AUTH_ACTOR_MISMATCH
AUTH_TARGET_MISMATCH
```

Guarded write failures keep the write-guard code `TARGET_MISMATCH`. Auth readiness mismatch uses `AUTH_TARGET_MISMATCH` so agents can tell readiness failure from a guarded mutation refusal.

### 8. Output Is For Agents First

stdout is command output. stderr is diagnostics and structured error envelopes. Successful command output must stay pipeable.

The global output controls are stable product surface:

- `--json`: machine-readable output.
- `--compact`: single-line JSON when used with `--json`.
- `--fields`: JSON field projection.
- `--id-only`: print only an id when the command returns one.
- `--quiet`: suppress successful output.
- `--fail-on-empty`: make empty list results fail intentionally.
- `--sort` and `--order`: deterministic list ordering by JSON field.
- `--format`: human output density for non-JSON output.

Auth commands, JSON output, debug logs, diagnostics, tests, and errors must not print OAuth access tokens, refresh tokens, or client secret values. They may report presence as `set` or `missing` and may report non-secret metadata such as actor, scopes, expiry, token type, and target readiness.

### 9. Command Ports Are The Preferred Command Seam

The GraphQL client package owns generated GraphQL operations, transport behavior, target resolution, and the write guard. CLI command logic should not grow around generated response shapes.

For commands with meaningful behavior, the CLI package should define a narrow consumer-owned Command Port. A Command Port returns domain summaries or command request/result types, not generated GraphQL responses. Production adapters should be thin forwarding adapters over the client package. Command tests should prefer in-memory fakes at the Command Port seam.

This is the baseline for new and changed command logic. The migration is intentionally incremental. Existing broad read commands may still call the runtime and client helpers directly until they gain enough behavior to justify a Command Port.

### 10. Collection Projection Stays Type-Derived With A Curated Fallback

`--fields` projection over list pages should derive the collection key from the list page type when the shape has exactly one exported slice field with a JSON name. That keeps common list pages type-driven instead of hand-registered.

The explicit collection-key allowlist remains for ambiguous cases and for detail shapes that contain incidental arrays. Generic "top-level array" detection is rejected because it would project detail objects incorrectly.

### 11. Live Verification Is Explicit And Fixture-Based

OAuth live coverage uses a pre-created Linear OAuth app fixture supplied by environment variables. `task live-oauth` verifies client-credentials auth against the pinned target, requires app actor readiness, and emits only redacted auth status.

`task live-smoke` may bootstrap through the same OAuth harness when fixture env is present. Scheduled and manual integration runs use the same live gate before live integration and smoke checks. Browser authorization-code login remains a manual smoke path because interactive consent is not reliable unattended CI surface.

## Consequences

- Guarded writes cost an extra target-resolution step, but stale or wrong auth fails before mutation.
- Auth setup is more explicit than token-paste workflows, but readiness output is truthful.
- Missing scopes do not trigger surprise browser windows from ordinary commands.
- Local auth state is machine-local by design, so moving a repo does not move auth.
- Live OAuth coverage requires prepared secrets and a matching pinned target.
- Command Port migration proceeds by behavior pressure, not by a broad rewrite.
- The old per-topic ADR files are removed from the working tree. Git history remains the provenance for their exact original wording.

## Rejected Alternatives

- **Bypass flag for guarded writes**: rejected because it would undermine the only hard safety boundary.
- **Personal-token product auth**: rejected because linctl's product direction is OAuth app credentials, renewable token state, and explicit actor attribution.
- **Admin scope by default**: rejected because ordinary reads and guarded writes do not need admin-level permission.
- **Automatic browser reauthorization from ordinary commands**: rejected because it is hostile to automation and makes scope escalation implicit.
- **Generic array-based field projection**: rejected because detail responses can contain arrays that are not the primary collection.
- **Creating or rotating OAuth apps during live tests**: rejected because it expands the fixture surface into admin behavior and cleanup risk.
- **Immediate Command Port rewrite for every command**: rejected because thin read commands do not all justify a new seam yet.

## Code Alignment

- Domain language: `CONTEXT.md`.
- Config and profiles: `internal/config/load.go`.
- Target resolution and guarded writes: `internal/client/target.go`, `internal/client/write_guard.go`.
- Transport, retries, and auth failure sentinel: `internal/client/transport.go`.
- Auth state and profile selection: `internal/auth/state.go`, `internal/auth/session.go`, `internal/auth/token.go`.
- OAuth error codes: `internal/auth/token_error.go`.
- OAuth token client and PKCE: `internal/oauth/client.go`, `internal/oauth/pkce.go`.
- Auth command surface: `internal/cli/auth.go`, `internal/cli/auth_login.go`.
- Runtime token recovery: `internal/cli/runtime.go`.
- Output and error envelope policy: `internal/cli/output.go`, `internal/cli/root.go`.
- Command inventory and collection projection: `internal/cli/command_inventory.go`.
- Command Port examples: `internal/cli/issue_port.go`, `internal/cli/bulk.go`, `internal/cli/comment_port.go`, `internal/cli/cycle_port.go`, `internal/cli/document_port.go`, `internal/cli/project_update_port.go`.
- Live OAuth gate: `scripts/live-oauth.sh`, `Taskfile.yml`, `.github/workflows/ci.yml`, `.github/workflows/integration.yml`.

## Fitness Checks

Use these checks when changing the baseline behavior:

```bash
go test -count=1 ./internal/client ./internal/cli ./internal/auth ./internal/oauth ./internal/config
go run github.com/go-task/task/v3/cmd/task@latest ci
go run github.com/go-task/task/v3/cmd/task@latest coverage
go run github.com/go-task/task/v3/cmd/task@latest live-oauth
```

`live-oauth` requires fixture env and should be skipped when the fixture is absent. Browser login remains a manual smoke check around `linctl auth login --callback -`.
