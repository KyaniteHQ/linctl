# linctl upstream coverage control

## Allowed scope

- Maximize `linctl` coverage of the current upstream Linear SDK/API while preserving the agent-safe Linear control surface.
- Implement accepted repo-planned PM domains before broad upstream expansion.
- Prefer safe read/list/get/search commands first; add guarded writes only when Pinned Target / Resolved Target / Target Mismatch semantics are clear.
- Keep generated GraphQL code generated through `go generate ./...`; do not hand-edit `internal/client/generated.go`.

## Protected files and behavior

- Preserve unrelated dirty work in the current `master` checkout.
- Do not relax `internal/client/write_guard.go` or add bypass flags.
- Do not expose destructive, admin, auth, internal, or alpha operations as ordinary commands without explicit user approval.
- Do not print Linear tokens or live secret values.

## Current phase

Phase 5: final verification and objective audit.

## Latest coverage counts

- Upstream Linear SDK root methods: 458 total, 30 implemented/root-backed, 458 classified.
- Upstream GraphQL Query root fields: 158 total, 18 implemented/root-backed, 158 classified.
- Upstream GraphQL Mutation root fields: 364 total, 12 implemented/root-backed, 364 classified.
- Local generated Go operations: 54 total, 54 implemented, 54 classified.
- Repo domain-map commands: 58 total, 48 implemented, 58 classified.

## Final verification evidence

- `go generate ./...` idempotence check for `internal/client/generated.go`: passed.
- `go run github.com/go-task/task/v3/cmd/task@latest test`: passed.
- `go run github.com/go-task/task/v3/cmd/task@latest lint`: passed with 0 issues.
- `go run github.com/go-task/task/v3/cmd/task@latest coverage`: passed with total hand-written statement coverage 100.0%.
- `go vet ./...`: passed.
- `go build ./...`: passed.
- `go run github.com/go-task/task/v3/cmd/task@latest live-smoke`: passed.
- Tech debt audit artifact: `/tmp/linctl-tech-debt-audit-2026-06-20.md`.
- Thermo-nuclear review artifact: `/tmp/linctl-thermo-review-2026-06-20.md`.

## Accepted exclusion categories

- Unsafe destructive operation without a safe archive/cleanup model.
- Admin/auth/internal/alpha operation without a user-approved command contract.
- Upstream root field that is not useful from an agent CLI or cannot be scoped safely.
- Mutation whose target comparison semantics are ambiguous.
- Product deferral documented with evidence and a future command shape.

## Decision gates

Require explicit Omer approval before:

- Exposing destructive/admin/auth/internal/alpha operations.
- Adding a new dependency.
- Introducing public command names that conflict with `CONTEXT.md` or `docs/domain-map.md`.
- Making a broad architecture pivot away from the current Cobra plus genqlient module shape.
