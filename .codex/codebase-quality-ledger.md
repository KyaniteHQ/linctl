# Codebase Quality Ledger

## Current State

- Repo: `/home/oruc/Desktop/workspace/linctl`.
- Branch: `master`.
- Baseline commit: `c389900`.
- Latest local commit: `71240f0` ported document writes through narrow Command Ports.
- Rolling report: `/tmp/codex-codebase-quality/linctl-quality-report.html`.
- Worktree exception at run start: `.gitignore` modified with local ignore rules and untracked `.directory`; both are treated as pre-existing Omer changes and must remain unstaged.
- Likely next action: after the first slice is committed, re-discover whether a similarly small Command Port locality slice remains.

## Validation Surface

- Focused tests: `go test ./internal/cli -run 'Test_runIssue.*Relation|Test_issueClientAdapter_forwards_to_client'`.
- Required broad gates after a completed slice: `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
- Live smoke: only if credentials are available safely and the slice touches live Linear behavior. The first planned slice is behavior-preserving command-local refactoring, so live smoke is not required unless later evidence changes that.
- Commit preflight: `git status --short` must show only scoped changes plus the allowed `.gitignore` / `.directory` exception; `git diff --cached --check` must pass before each commit.

## Completed Slices

- 2026-06-26: Port-level issue relation writes.
  - Files: `internal/cli/issue_port.go`, `internal/cli/issue_relation_write.go`, `internal/cli/issue_port_test.go`.
  - Behavior impact: no public CLI behavior change; `issue relate` and `issue unrelate` still build the same requests, call the same guarded client adapter, and render through the same writers.
  - Quality impact: moved issue relation create/delete command bodies behind narrow Command Ports and focused run functions, matching adjacent issue write seams.
  - Validation: `go test ./internal/cli -run 'Test_runIssue.*Relation|Test_issueClientAdapter_forwards_to_client'`; `go test ./internal/cli -cover`; `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
  - Notes: `task ci` skipped coverage-ledger drift because `/tmp/linctl-upstream-linear` is unavailable; all other CI steps passed.
  - Commit: this commit.
- 2026-06-26: Port-level issue link write.
  - Files: `internal/cli/issue_port.go`, `internal/cli/issue_write.go`, `internal/cli/issue_port_test.go`.
  - Behavior impact: no public CLI behavior change; `issue link` still accepts the same positional args and `--title` / `--subtitle` flags, calls the same guarded client adapter, and renders the same attachment-link output.
  - Quality impact: moved attachment-link command execution behind a narrow Command Port and focused run function, so request assembly is characterized without transport payloads.
  - Validation: `go test ./internal/cli -run 'Test_runIssueLink|Test_issueClientAdapter_forwards_to_client'`; `go test ./internal/cli -cover`; `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
  - Notes: `task ci` skipped coverage-ledger drift because `/tmp/linctl-upstream-linear` is unavailable; all other CI steps passed.
  - Commit: this commit.
- 2026-06-26: Port-level comment writes.
  - Files: `internal/cli/comment.go`, `internal/cli/comment_port.go`, `internal/cli/comment_port_test.go`.
  - Behavior impact: no public CLI behavior change; `comment update` still resolves `--body`, stdin, and `--body-file` before the same guarded client write, and `comment delete` renders the same deletion output.
  - Quality impact: moved comment update/delete execution behind small Command Ports and focused run functions, so command request/body handling is covered without GraphQL payloads.
  - Validation: `go test ./internal/cli -run 'Test_runComment|Test_commentClientAdapter'`; `go test ./internal/cli -cover`; `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
  - Notes: `task ci` skipped coverage-ledger drift because `/tmp/linctl-upstream-linear` is unavailable; all other CI steps passed.
  - Commit: this commit.
- 2026-06-26: Port-level document writes.
  - Files: `internal/cli/document.go`, `internal/cli/document_port.go`, `internal/cli/document_port_test.go`.
  - Behavior impact: no public CLI behavior change; `document create` and `document update` still resolve `--content`, stdin, and `--content-file` before the same guarded client writes and render the same document output.
  - Quality impact: moved document create/update execution behind small Command Ports and focused run functions, so content/request handling is covered without GraphQL payloads.
  - Validation: `go test ./internal/cli -run 'Test_runDocument|Test_documentClientAdapter'`; `go test ./internal/cli -cover`; `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
  - Notes: `task ci` skipped coverage-ledger drift because `/tmp/linctl-upstream-linear` is unavailable; all other CI steps passed.
  - Commit: this commit.
- 2026-06-26: Port-level project writes.
  - Files: `internal/cli/project_write.go`, `internal/cli/project_port.go`, `internal/cli/project_port_test.go`.
  - Behavior impact: no public CLI behavior change; `project create`, `project update`, and `project archive` still build the same requests, call the same guarded client writes, and render the same project output.
  - Quality impact: replaced closure-based direct client calls in project command wiring with small Command Ports and focused run functions.
  - Validation: `go test ./internal/cli -run 'Test_runProject|Test_projectClientAdapter'`; `go test ./internal/cli -cover`; `go generate ./...`; `go run github.com/go-task/task/v3/cmd/task@latest ci`; `go run github.com/go-task/task/v3/cmd/task@latest coverage`.
  - Notes: first `task ci` attempt failed on a 122-character adapter method line; the line was wrapped and all gates passed afterward. `task ci` skipped coverage-ledger drift because `/tmp/linctl-upstream-linear` is unavailable; all other CI steps passed.
  - Commit: this commit.

## Deferred Needs Omer

- First seen: 2026-06-26. Area: public CLI expansion and write-surface additions. Blocked reason: this loop is behavior-preserving and must not add public CLI contracts. Unblock action: Omer explicitly requests a public command or write model.
- First seen: 2026-06-26. Area: generated Linear coverage expansion. Blocked reason: generated files and upstream coverage changes are out of scope unless tied to the current slice. Unblock action: run a dedicated coverage expansion loop from current upstream truth.

## Candidate Signals

- Candidate: `issueClientAdapter` now satisfies issue, bulk issue import, and project-update Command Ports; a later naming/locality cleanup may make sense if it stays small.
- Candidate: `issue start` remains a simple one-id guarded write; it is lower leverage than request-assembly ports unless a future refactor touches start semantics.
- Candidate: simple guarded-write wrappers may benefit from one more characterization test if a future refactor touches `runGuardedWrite`.
- Candidate: `project-milestone` and `Cycle` writes still use generic closure wrappers; they are safe only if the next slice keeps domain naming precise and avoids public CLI changes.
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
- 2026-06-26T00:05:00+03:00: Completed and validated the issue relation Command Port slice; ready to commit after staged diff checks.
- 2026-06-26T00:07:00+03:00: Completed and validated the issue link Command Port slice; ready to commit after staged diff checks.
- 2026-06-26T00:10:00+03:00: Completed and validated the comment write Command Port slice; ready to commit after staged diff checks.
- 2026-06-26T00:13:00+03:00: Completed and validated the document write Command Port slice; ready to commit after staged diff checks.
- 2026-06-26T00:16:00+03:00: Completed and validated the project write Command Port slice after fixing one line-length lint issue; ready to commit after staged diff checks.
