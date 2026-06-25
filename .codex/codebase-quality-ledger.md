# Codebase Quality Ledger

## Current State

- Repo: `/home/oruc/Desktop/workspace/linctl`.
- Branch: `master`.
- Baseline commit: `c389900`.
- Rolling report: `/tmp/codex-codebase-quality/linctl-quality-report.html`.
- Worktree exception at run start: `.gitignore` modified with local ignore rules and untracked `.directory`; both are treated as pre-existing Omer changes and must remain unstaged.
- Likely next action: deepen Command Port coverage for the `issue relate` and `issue unrelate` guarded writes without changing CLI behavior.

## Validation Surface

- Focused tests: `go test ./internal/cli -run 'Test_runIssue.*Relation|Test_issueClientAdapter_forwards_to_client'`.
- Required broad gates after a completed slice: `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
- Live smoke: only if credentials are available safely and the slice touches live Linear behavior. The first planned slice is behavior-preserving command-local refactoring, so live smoke is not required unless later evidence changes that.
- Commit preflight: `git status --short` must show only scoped changes plus the allowed `.gitignore` / `.directory` exception; `git diff --cached --check` must pass before each commit.

## Completed Slices

None yet.

## Deferred Needs Omer

- First seen: 2026-06-26. Area: public CLI expansion and write-surface additions. Blocked reason: this loop is behavior-preserving and must not add public CLI contracts. Unblock action: Omer explicitly requests a public command or write model.
- First seen: 2026-06-26. Area: generated Linear coverage expansion. Blocked reason: generated files and upstream coverage changes are out of scope unless tied to the current slice. Unblock action: run a dedicated coverage expansion loop from current upstream truth.

## Candidate Signals

- Selected first: `issue relate` and `issue unrelate` still assemble guarded-write request and output dispatch directly inside Cobra `RunE`; adjacent issue writes use small `runIssue...` functions with fake Command Ports.
- Candidate: `issueClientAdapter` now satisfies issue, bulk issue import, and project-update Command Ports; a later naming/locality cleanup may make sense if it stays small.
- Candidate: simple guarded-write wrappers may benefit from one more characterization test if a future refactor touches `runGuardedWrite`.
- Deferred for now: docs/test scenario cleanup unless tied to verified behavior from a code slice.

## Recently Failed

None yet.

## Assumptions To Re-check

- `task` is invoked through `go run github.com/go-task/task/v3/cmd/task@latest`.
- `scripts/coverage.sh` gates rounded 100.0% and has known pre-existing uncovered `command_inventory.go` branches; new statements in this run must be covered.
- The allowed dirty worktree exception remains `.gitignore` and `.directory` only.

## History

- 2026-06-26T00:04:00+03:00: Started autonomous quality loop from `c389900`; read `CLAUDE.md`, `CONTEXT.md`, recent Command Port commits, and quality-loop instructions.
- 2026-06-26T00:04:00+03:00: Ranked first slice as Command Port coverage/locality for `issue relate` and `issue unrelate`; no product code changed yet.
