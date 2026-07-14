# Codebase Quality Ledger

## Current State

- Repo: `/home/oruc/Desktop/workspace/linctl`.
- Branch: `master`.
- Original quality-loop baseline: `c389900`.
- Latest completed program: the 2026-07-13 thermo-nuclear audit and behavior-preserving fix loop.
- Run evidence: `/tmp/codex-runs/linctl-thermo-fix-20260713/manifest.json` and its journal.
- No standing dirty-worktree exception applies. Preserve unrelated user changes and stage only the intended release set.

## Validation Surface

- Required broad gates: `go tool task ci` and `go tool task coverage`.
- Generated boundary: `go tool task generated-client-check` requires both generated outputs to be tracked, rejects the legacy generated-client path, and checks regeneration drift.
- Integration compatibility: `go test -run '^$' -tags=integration ./internal/client` and `go vet -tags=integration ./internal/client`.
- Live behavior: run credentialed integration and smoke tests only with prepared disposable Linear fixtures.
- Commit preflight: inspect `git status --short`, stage only intended files, and require `git diff --cached --check`.

## Completed Programs

- 2026-06-26: introduced focused command seams for issue selection and other command workflows. The later architecture review retained only seams that isolate a multi-call decision, batching, or an external effect.
- 2026-07-13: made auth token persistence transactional across processes, made injected credential provenance explicit, and closed target-proof gaps.
- 2026-07-13: moved generated GraphQL transport code behind `internal/client/internal/gql`, with stable wire operation names and reproducible output.
- 2026-07-13: replaced one-call forwarding command seams with direct exported client calls while keeping fail-closed write safety in the client package.
- 2026-07-13: deleted per-page JSON callbacks and consolidated repeated issue child projections with byte-identical command and GraphQL contract proofs.
- 2026-07-13: hardened package discovery and coverage so partial discovery, stale profiles, empty profiles, and any zero-count hand-written statement block fail the gate.

## Deferred Needs Omer

- Credentialed live integration remains pending when the required disposable Linear fixture is unavailable.
- Public CLI expansion, update-field tri-state semantics, and pagination policy remain separate product decisions.

## Assumptions To Re-check

- Repository tasks run through the pinned tool dependency as `go tool task <task>`.
- `scripts/coverage.sh` enforces the underlying invariant directly: every hand-written statement block must have a nonzero execution count. The rounded `go tool cover` total is display output, not the gate.
- Ordinary schema and coverage checks use Linear SDK commit `202a5e0fbce142cd1f0cd5c5265c3a310f629702`; nightly schema drift explicitly follows `master`.
- Generated GraphQL outputs are `internal/client/internal/gql/generated.go` and `internal/client/internal/gql/exports.go`. `internal/client/generated.go` must remain absent and untracked.

## History

- 2026-06-26: completed the original command-seam quality loop from baseline `c389900` with focused tests and broad gates.
- 2026-07-13: completed the thermo-nuclear audit remediation program; CI, exact coverage, generated reproducibility, integration-tag compile/vet, and independent reviews passed. Credentialed live integration was unavailable.
- 2026-07-14: aligned release checks with the generated package boundary, pinned ordinary upstream checks, and refreshed this ledger to current architecture and gate behavior.
