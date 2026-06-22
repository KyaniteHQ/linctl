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

**Cleanup**:
The safe removal path for temporary Linear entities. For projects, cleanup means archive rather than hard delete.
_Avoid_: Delete, purge
