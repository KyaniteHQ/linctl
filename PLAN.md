# linctl upstream coverage plan

## Objective

Maximize `linctl` coverage of the current upstream Linear SDK/API while preserving the agent-safe control-surface model.

## Strategy

1. Recompute baselines from current upstream Linear SDK and local repo state.
2. Write `docs/linear-api-coverage.md` with classified SDK root methods, GraphQL root fields, local generated Go operations, and domain-map commands.
3. Run architecture and quality scans before broad command expansion.
4. Implement repo-planned PM domains first with vertical TDD tracer bullets.
5. Expand safe upstream read coverage by product domain.
6. Add guarded writes only when target semantics are explicit.
7. Keep docs, skill, usage, tests, generated code, and live verification current after each domain.

## Completed slices

- Fresh upstream baseline from Linear commit `df20561`.
- Initial `docs/linear-api-coverage.md` with no unclassified rows.
- Architecture report generated in `/tmp`.
- Top architecture recommendation implemented: domain-split GraphQL operation modules behind the unchanged `operations/*.graphql` seam.
- Cycle read expansion: `cycle list`, `cycle get`.
- ProjectMilestone read expansion: `project-milestone list`, `project-milestone get`.
- Sprint read aliases: `sprint current`, `sprint report`.
- ProjectMilestone guarded writes: `project-milestone create`, `project-milestone update`.
- Cycle guarded writes: `cycle create`, `cycle update`, `cycle archive`.
- Document reads: `document list`, `document get`; writes documented as blocked pending parent-resolution guard design.
- Label reads: `label list`, `label get`; writes documented as blocked pending team-scope guard design.
- Team reads: `team list`, `team get`, `team members`; writes documented as blocked organization/admin surface.
- User reads: `user list`, `user get`, `user me`; user writes documented outside the v1 PM command surface.
- Final `docs/linear-api-coverage.md` regenerated from upstream Linear commit `df20561`.
- Tech-debt and thermo-nuclear review artifacts written under `/tmp`.
- Final objective audit completed for the upstream coverage push; the authoritative current coverage
  ledger is `docs/linear-api-coverage.md`.

## Current slice

Maintenance mode: keep generated coverage, command docs, test scenarios, and `skills/linctl` aligned whenever
the Linear schema or CLI command surface changes.

## Next verification

- Run `go run github.com/go-task/task/v3/cmd/task@latest ci`.
- Run `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
- Refresh the coverage ledger when `internal/client/schema.graphql`, generated operations, or public command
  registrations change.
