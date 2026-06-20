# linctl upstream coverage notes

## Domain language

- Use Linear Control Surface, Pinned Target, Resolved Target, Target Mismatch, Read Command, Guarded Write, Team-Scoped Write, Resource-Scoped Write, Current Issue, Issue Identifier, ProjectMilestone, Cycle, Sprint, Namespaced Throwaway Resource, and Cleanup exactly as defined in `CONTEXT.md`.
- Avoid "workspace" for target language.
- Sprint is a read-only report alias over Cycle, not a mutation-owning domain.

## Safety model

- Reads can be broad.
- Writes must compare the active credential's resolved organization/team/project against the pinned target.
- Target Mismatch is a hard stop.
- Project cleanup means archive, not hard delete.

## Early architecture observations

- `internal/cli/command_flow_test.go` has crossed 1000 handwritten lines.
- `internal/client/coverage_test.go` has crossed 1000 handwritten lines.
- The new goal should treat test decomposition as an architecture/maintainability item, not as incidental cleanup.

## Baseline

- Upstream Linear SDK commit: `df20561`.
- Initial ledger path: `docs/linear-api-coverage.md`.
- The ledger classifies every upstream SDK root method, upstream Query root field, upstream Mutation root field, local generated operation, and domain-map command.
- Domain-map implementation started at 28/58 commands and ends at 48/58 commands.
- Remaining 10 domain-map commands are intentionally blocked: ProjectMilestone delete, Document create/update/delete, Label create/update/delete, Team create/update/delete.
- Document, Label, Team, and User reads are implemented; Team membership is read-only.

## Architecture report

- Report path: `/tmp/architecture-review-linctl-20260620-021942.html`.
- Top recommendation: split the GraphQL operation module by domain to improve locality while preserving the existing genqlient seam.
- Implemented by replacing `internal/client/operations/viewer.graphql` with domain modules under `internal/client/operations/`.

## Final review artifacts

- Tech debt audit: `/tmp/linctl-tech-debt-audit-2026-06-20.md`.
- Thermo-nuclear review: `/tmp/linctl-thermo-review-2026-06-20.md`.
- Structural blocker found during final lint/review: duplicated list scaffolding. Fixed with `internal/cli/list.go`.
