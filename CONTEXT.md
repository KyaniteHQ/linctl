# linctl

linctl is a Linear control surface for agent-safe coordination. Its language separates free reads from guarded writes so agents can inspect Linear broadly while mutating only an explicitly pinned target.

## Language

**Linear Control Surface**:
A compact command interface for reading and changing Linear entities from an agent session. It exists to make Linear coordination predictable and safe.
_Avoid_: Linear wrapper, API helper, generic CLI

**Pinned Target**:
The Linear organization, team, and optional project that a repo declares as the only allowed destination for writes. A pinned target is expressed as `org_id`, `team_key`, `team_id`, and optional `project_id`.
_Avoid_: Workspace, account, default project

**Resolved Target**:
The Linear organization, team, and optional project proven from the active credential at command time. A resolved target must match the pinned target before guarded writes proceed.
_Avoid_: Current workspace, active account

**OAuth App Credential**:
A Linear OAuth application credential used by linctl as the product authentication model for agent-safe coordination.
_Avoid_: API key, personal token, user token

**App Actor**:
The Linear application identity that performs changes when linctl authenticates through OAuth app authorization.
_Avoid_: Human actor, bot user, service account

**User Actor**:
The Linear user identity that performs changes when linctl is explicitly authorized to act as that user.
_Avoid_: Default actor, app actor

**User Authorization Flow**:
The OAuth authorization-code flow where a Linear user approves linctl in the browser and linctl stores renewable OAuth tokens for later commands.
_Avoid_: Device login, API-key setup, manual token paste

**Manual Callback Fallback**:
The browser login fallback where linctl prints an OAuth authorization URL and accepts the returned callback or authorization code from the user.
_Avoid_: Device flow, token paste, API-key fallback

**PKCE Login**:
The authorization-code login behavior that always uses a generated PKCE verifier and S256 challenge.
_Avoid_: Plain challenge, secret-only login, unbound code exchange

**Client Secret Requirement**:
The rule that OAuth client secret is optional for PKCE browser login when Linear permits it, but required for headless app authorization.
_Avoid_: Always-secret login, secretless app auth

**Local Auth State**:
Machine-local OAuth client and token material stored outside repo configuration using the operating system's user-specific config or state location.
_Avoid_: Repo auth config, checked-in credential, project token

**Auth Profile**:
The named local OAuth auth state selected by linctl's active profile.
_Avoid_: Separate login slot, workspace switcher, hidden account

**Auth Override**:
Non-persistent OAuth credential material supplied by the process environment for automation, without changing local auth state.
_Avoid_: Primary config, repo credential, saved token

**Redacted Auth Output**:
Auth-related command, JSON, debug, and test output that reports presence and metadata without printing OAuth secret values.
_Avoid_: Token output, secret echo, raw auth debug

**OAuth Logout**:
The command behavior that revokes Linear OAuth tokens when possible and removes local token state while preserving OAuth app configuration unless explicitly forgotten.
_Avoid_: Config reset, app deletion, token expiry

**Auth Command Surface**:
The `linctl auth` command family for configuring OAuth app settings, authorizing browser login, authorizing headless app access, inspecting auth status, refreshing tokens, and logging out.
_Avoid_: Token command, API-key setup, credential helper

**Live Auth Status**:
The auth status behavior that refreshes expired OAuth tokens when possible and verifies saved or overridden OAuth state against Linear before reporting readiness.
_Avoid_: Local-only status, cached readiness

**Auth Readiness**:
The live proof that OAuth credentials can authenticate to Linear with the expected actor, scopes, and resolved target before linctl reports auth as usable.
_Avoid_: Token saved, local config present, assumed readiness

**Token Refresh**:
The runtime behavior that renews expired OAuth access tokens, persists rotated token state, and retries the original Linear request once.
_Avoid_: Manual refresh step, endless retry, command-specific auth

**App Token Reacquisition**:
The runtime behavior that reuses a client-credentials access token until expiry or authorization failure, then requests a new app actor token once.
_Avoid_: Refresh token, per-command token fetch, endless retry

**Default OAuth Scope Set**:
The OAuth scopes linctl requests for ordinary read commands and guarded writes without asking for admin-level access.
_Avoid_: Admin scope, all scopes, per-command prompt

**Scope Escalation**:
An explicit reauthorization step the user runs after linctl reports that the current OAuth authorization is missing a required scope.
_Avoid_: Automatic browser launch, silent permission upgrade

**Auth Error Code**:
A stable machine-readable OAuth failure reason emitted by linctl so agents can recover without scraping prose.
_Avoid_: Free-form auth failure, panic text, debug-only reason

**Target Mismatch**:
A refusal state where the resolved target does not match the pinned target. A target mismatch is a hard stop for guarded writes.
_Avoid_: Soft warning, confirmation prompt

**Read Command**:
A command that inspects Linear without changing it. Read commands may operate before a pinned target is confirmed.
_Avoid_: Safe write, dry run

**Guarded Write**:
A command that changes Linear only after comparing the pinned target with the resolved target. Guarded writes fail closed on mismatch.
_Avoid_: Mutation, update call

**Team-Scoped Write**:
A guarded write that creates a new Linear entity inside the resolved team. It compares organization and team because the entity does not exist yet.
_Avoid_: Unscoped create, workspace write

**Resource-Scoped Write**:
A guarded write against an existing Linear entity. It resolves the entity first and compares the pinned project when a project is configured.
_Avoid_: Direct update, blind write

**Command Port**:
The narrow, domain-typed interface a Read Command or Guarded Write depends on to reach Linear, decoupled from the GraphQL transport. A command port is defined by its consumer, returns domain summaries rather than GraphQL responses, and is satisfied in production by a thin adapter over the client and in tests by an in-memory fake. It makes the command's interface the test surface.
_Avoid_: client, gateway, service, mock

**Current Issue**:
The Linear issue referenced by the current checkout context. It comes from an issue identifier in the branch name or checkout metadata.
_Avoid_: Active ticket, selected issue

**Issue Identifier**:
The human-readable Linear issue key such as `LIT-123`. It is distinct from the Linear issue id.
_Avoid_: Issue id, ticket number

**ProjectMilestone**:
The Linear project milestone entity. Use the full schema name when discussing code or command surface.
_Avoid_: Milestone

**Cycle**:
Linear's time-boxed planning entity. Use Cycle for mutations and schema-aligned command names.
_Avoid_: Sprint

**Sprint**:
A report-facing alias over Cycle. Sprint is not a Linear schema entity and should not own mutations.
_Avoid_: Cycle mutation

**Initiative**:
Linear's current strategic planning entity for grouping projects toward a goal. Use Initiative for new planning workflows.
_Avoid_: Roadmap for new planning

**Roadmap**:
Linear's deprecated project grouping surface. Keep Roadmap reads for legacy compatibility only; do not add Roadmap writes without an explicit guard model.
_Avoid_: Current planning entity, Initiative replacement

**Namespaced Throwaway Resource**:
A temporary Linear entity created for verification with a recognizable namespace. It must be cleaned up through the supported cleanup path.
_Avoid_: Test fixture, dummy data

**OAuth App Fixture**:
A dedicated pre-created Linear OAuth application used by live tests to verify OAuth behavior without dynamically creating or rotating OAuth apps.
_Avoid_: Throwaway OAuth app, admin-created test app

**Browser Login Smoke**:
A manual or semi-automated live verification of the OAuth browser authorization flow.
_Avoid_: Fully unattended consent test, mocked login

**OAuth Live Gate**:
The explicit live test task that verifies OAuth behavior with a dedicated OAuth app fixture and may be included by broader live smoke when the fixture is configured.
_Avoid_: Unit-only auth proof, always-on browser CI

**Cleanup**:
The safe removal path for temporary Linear entities. For projects, cleanup means archive rather than hard delete.
_Avoid_: Delete, purge
