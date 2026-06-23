# Thermo-Nuclear Code Quality Review

Date: 2026-06-23
Repo: `/home/oruc/Desktop/workspace/linctl`
Branch: `master`

## Scope

Full repo review across Go source, tests, scripts, CI, docs, generated-code boundaries, and the bundled `linctl` skill. Generated GraphQL output and schema were treated as integration boundaries rather than hand-written code to simplify.

The repo is healthy enough to pass its main gates, but it has one dangerous structural weakness: command and safety truth is scattered across Cobra registration, command docs, output projection, tests, live smoke, domain docs, and the Linear coverage generator. That duplication makes the codebase look more complete and safer than the actual authority model can prove.

## Verification

Passed:

- `bash scripts/go-packages.sh`
- `go test -count=1 $(bash scripts/go-packages.sh)`
- `go test -race -shuffle=on -count=1 ./internal/cli ./internal/client`
- `go vet $(bash scripts/go-packages.sh)`
- `bash scripts/coverage.sh` - reported `100.0%` hand-written coverage
- `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --timeout 5m $(bash scripts/go-packages.sh)`
- `shellcheck scripts/*.sh`
- `SHELLCHECK_OPTS=--exclude=SC2046 go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12`
- `go tool govulncheck $(bash scripts/go-packages.sh)`

Useful failing probes:

- `go list ./...` fails because ignored local skill examples under `skills/cc-skills-golang` are still inside the Go module and import packages that are not in this module.
- `go tool gofumpt -l cmd internal scripts` fails because `go.sum` has only `mvdan.cc/gofumpt v0.10.0/go.mod`, not the module checksum needed by `go tool gofumpt`.
- Raw `actionlint` fails on `SC2046`; the repo intentionally passes when using `SHELLCHECK_OPTS=--exclude=SC2046`, matching `Taskfile.yml`.

Not run:

- `task ci`, because `generate-check` can write generated artifacts and the worktree already had unrelated dirty files.
- Live smoke/write integrations, because they require real Linear credentials and disposable-resource authority.

## Findings

### 1. Critical - Coverage ledger misclassifies mutating API methods as safe read candidates

Evidence:

- `docs/linear-api-coverage.md:60`, `:68`, `:70`, `:143`, and `:318-324` mark methods such as `attachmentLinkFront`, `attachmentLinkURL`, `attachmentSyncToSlack`, `customerMerge`, `notificationArchiveAll`, `notificationMarkReadAll`, and `notificationSnoozeAll` as `safe_candidate`.
- The same names are classified later as blocked mutations at `docs/linear-api-coverage.md:671`, `:679`, `:681`, `:699`, and `:873-882`.
- `scripts/linear_api_coverage.go:1072-1115` falls through to `safe_candidate` unless the loose name heuristic catches a danger word.

Why it matters:

The ledger is supposed to be the safety map for what should or should not become CLI surface. It currently says some mutating methods are safe future reads. This is worse than missing coverage because it creates false confidence.

Fix:

Stop using loose string heuristics as the final authority. Classify SDK methods through generated GraphQL root-field kind first. Unknown methods should default to `blocked_needs_design` or `unclassified`, and CI should fail until each is explicitly classified.

### 2. High - There is no canonical command/control-surface model

Evidence:

- Cobra wiring and command behavior are embedded directly in large files such as `internal/cli/issue.go` and `internal/cli/project.go`.
- Output projection keeps a separate global collection registry in `internal/cli/output.go:255`.
- Command coverage is hard-coded again in `scripts/linear_api_coverage.go:645`.
- Command docs are generated from Cobra through `scripts/gen-skill.go`, while the coverage ledger uses its own registry.
- Live smoke has another command inventory in `scripts/live-smoke.sh`.
- Flow tests duplicate large parts of the command model in `internal/cli/command_flow_test.go`.

Why it matters:

Adding a command currently means updating several different control planes and hoping they stay aligned. The repo has green tests because those tests mirror the duplication, not because one authoritative model is being checked from multiple angles.

Fix:

Introduce a typed command metadata layer that records command path, entity, resolver, parent target, GraphQL operation/root, output collection key, write/read safety, and doc category. Generate or verify Cobra docs, coverage ledger rows, projection support, and smoke scenarios from that one model.

### 3. High - Write safety is enforced by convention instead of a single guarded mutation boundary

Evidence:

- `internal/client/write_guard.go` defines guard primitives, but callers still repeat guard construction and mutation flow manually.
- Similar guarded-write patterns appear across issue, cycle, document, project, comment, and related client write files.
- `internal/client/issue_write.go:377-386` parses any integer priority even though the error message promises `0-4`; CLI normalization catches ordinary command input, but the client request layer does not enforce its own contract.
- Several tests assert operation names or fake response shape more strongly than payload semantics.

Why it matters:

The repo's core safety promise is target-pinned guarded writes. That promise should be impossible to bypass accidentally from any write path. Today it depends on every write implementation repeating the pattern correctly.

Fix:

Create one `guardedMutation` helper that takes target identity, expected key/name, operation name, variables, and response decoder. Put validation such as priority range at the client request boundary too. Update write tests to assert target variables and guarded mutation payloads, not only operation names.

### 4. High - File upload/download bypasses the CLI timeout model and buffers entire files

Evidence:

- `internal/cli/runtime.go` applies timeout to the GraphQL runtime, but `internal/cli/files.go:26` uses `http.DefaultClient` for storage transfers.
- `internal/cli/files.go:84` reads the full upload into memory before preparing the upload.
- `internal/cli/files.go:110-126` performs the PUT through the package-level HTTP client.
- `internal/cli/files.go:156-176` downloads through the same client and `io.ReadAll` before writing the file.

Why it matters:

A user can set CLI timeout expectations and still hang or consume unbounded memory on file transfer paths. This is a real behavior gap, not just style.

Fix:

Pass a timeout-aware HTTP client through runtime or create a file-transfer runtime from the same root options. Stream upload/download bodies, enforce maximum expected sizes where Linear provides them, and write downloads through a temporary file before rename.

### 5. High - Structured error codes depend on a string suffix for not-found

Evidence:

- `internal/cli/output.go:46-49` maps not-found errors by checking whether `err.Error()` ends with `"not found"`.
- Client code emits many separate not-found messages from different files.

Why it matters:

The CLI promises machine-readable error codes. A prose edit can silently change `NOT_FOUND` into `INTERNAL`, breaking callers and agents.

Fix:

Add `client.ErrNotFound` and wrap all entity not-found errors with it. Keep the CLI mapping entirely sentinel-based.

### 6. High - The test suite protects coverage numbers more than refactor safety

Evidence:

- `internal/cli/command_flow_test.go` is over 5,000 lines and contains a large fake command world.
- `internal/client/coverage_test.go` is almost 5,000 lines and pins many private helper branches.
- `internal/cli/coverage_test.go` is over 2,700 lines.
- `scripts/coverage.sh` enforces 100% coverage over hand-written code.
- `internal/cli/runtime.go` uses mutable runtime globals, and `internal/cli/command_flow_test.go` explicitly avoids parallelism because of shared global state.

Why it matters:

The tests are expensive to change and make design improvements look riskier than they are. They also encourage adding coverage for helper branches instead of validating user-visible contracts.

Fix:

Keep the coverage gate, but move the valuable assertions up to contract-level tests: command metadata inventory, renderer behavior, target guard invariants, and live-safe read smoke. Then split the mega tests by domain and remove private-helper pinning where a public behavior test exists.

### 7. Medium - Generated and generated-adjacent tooling can drift

Evidence:

- `Taskfile.yml:33-39` drift-checks `internal/client/generated.go` and `skills/linctl/references/commands.md`.
- `docs/linear-api-coverage.md` is generated by `scripts/linear_api_coverage.go`, but is not part of the same drift check.
- `scripts/linear_api_coverage.go` is build-ignored and excluded from `scripts/go-packages.sh`, even though it is the safety ledger generator.
- Release preflight checks generated client drift but does not run the same full `generate-check` command reference gate.

Why it matters:

The repo treats the coverage ledger and generated skill docs as authoritative, but only some generated artifacts are checked in every path. The riskiest generator is outside normal Go package quality gates.

Fix:

Add a pinned coverage-ledger check or rename it as an explicit snapshot with source/date metadata. Add generator compile/lint coverage for `scripts/linear_api_coverage.go`. Run `task generate-check` or equivalent in release preflight.

### 8. Medium - Live smoke can accidentally use a normal production token

Evidence:

- `scripts/live-smoke.sh` falls back from `LINCTL_TEST_TOKEN` to `LINCTL_TOKEN` or `LINEAR_API_KEY`.
- The same script can run integration tests, and write integrations create Linear resources when `LINCTL_TEST_ENABLE_WRITES=1`.
- The bundled skill documents the same fallback behavior.

Why it matters:

Write-capable smoke should require a clearly named test credential. Falling back to ordinary token names makes accidental writes more plausible.

Fix:

Require `LINCTL_TEST_TOKEN` for `live-smoke`. If fallback is still desired for read-only smoke, split the command or require an explicit opt-in flag such as `LINCTL_ALLOW_NON_TEST_TOKEN_FOR_LIVE_SMOKE=1`.

### 9. Medium - Domain language drift leaks into help and generated skill docs

Evidence:

- `CONTEXT.md` and `docs/domain-map.md` tell the repo not to introduce a "workspace" model.
- Cobra help strings and generated skill docs still use wording like "workspace-level" in several commands.
- The generated command reference therefore teaches users and agents a domain model the repo explicitly rejects.

Why it matters:

This CLI is meant for agent-safe Linear operations. Help text is not just user copy; it is part of the control surface agents consume.

Fix:

Replace "workspace" with organization/team/visible-to-authenticated-user language as appropriate, regenerate skill docs, and add a small `rg`-based guard over external docs and help text with intentional exceptions.

### 10. Medium - High-complexity domain files are hiding missing concepts

Evidence:

- `internal/cli/issue.go` is about 1,500 lines.
- `internal/client/issue.go` is about 1,500 lines.
- `internal/client/user.go` is over 1,100 lines.
- `internal/client/user.go:260` has `gocognit`/`gocyclo` complexity around `33`.
- `internal/cli/user.go:114` has complexity around `32`.

Why it matters:

The problem is not file size alone. The repeated issue subcommands, user settings switches, and list/filter branches show that the code is missing domain-level vocabulary for child collections, settings categories, and query modes.

Fix:

Extract narrow, typed concepts only where repetition already exists: child collection specs, user setting category metadata, issue list mode/query builders, and output collection metadata. Avoid generic frameworks.

### 11. Medium - Clean install documentation is too destructive for a broad Linux claim

Evidence:

- `README.md` presents a broad macOS/Linux/Windows install path.
- The clean Linux path assumes `apt-get`, downloads a Go tarball, and removes `/usr/local/go`.

Why it matters:

This is surprising on non-Debian Linux and risky for users with a managed Go installation.

Fix:

Make `go install` the primary source install path. Move tarball install instructions to a distro-specific note with checksum verification and explicit permission guidance.

### 12. Medium - Tooling shape conflicts with the repo layout

Evidence:

- `go list ./...` walks ignored skill examples under `skills/cc-skills-golang`, causing missing-package failures.
- The repo works around this with `scripts/go-packages.sh`, which only lists product packages.
- `task fmt` depends on `go tool gofumpt`, but the current `go.sum` is missing the module checksum required to run it cleanly.

Why it matters:

New contributors and agents naturally try `go list ./...`, `go test ./...`, or `task fmt`. Those are common Go muscle-memory commands, and two of them fail or require a workaround.

Fix:

Move non-module examples outside the root module, put them behind a nested module, or make `skills/cc-skills-golang` impossible for Go package discovery to enter. Refresh tool checksums so `task fmt` works from a clean checkout.

## Suggested Fix Order

1. Fix the coverage classifier and ledger drift check first. That is the only finding that can directly mislabel unsafe future work as safe.
2. Build the command metadata model and generate/check the repeated surfaces from it.
3. Centralize guarded writes and typed not-found errors.
4. Repair file transfer timeout/streaming behavior.
5. Split the mega tests around the new contracts instead of preserving the current fake world.
6. Clean docs/help terminology and install instructions.
7. Fix Go module/tooling ergonomics for `go list ./...` and `task fmt`.

## Existing Dirty Worktree

The review started with unrelated modifications already present in:

- `CLAUDE.md`
- `CONTRIBUTING.md`
- `README.md`
- `docs/linear-cli-feature-leech.md`
- `docs/test-scenarios.md`
- `internal/client/generated_integration_test.go`
- `scripts/live-smoke.sh`
- `skills/linctl/SKILL.md`

This report intentionally does not revert or alter those files.
