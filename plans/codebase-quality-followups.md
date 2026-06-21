# Codebase Quality Follow-ups

These items came from the June 2026 quality loop. They are intentionally not mixed into the small
correctness and setup fixes in that loop because each needs its own before/after verification.

## 1. Split Broad Test Fixtures

Problem: `internal/cli/command_flow_test.go` and `internal/client/coverage_test.go` are high-value but
large, which makes small command changes harder to review.

Success criteria:

- Preserve the current scenario coverage and 100.0% hand-written statement coverage.
- Split by command family or domain with shared fixture builders.
- Keep command-flow evidence names stable enough for `docs/test-scenarios.md`.

Verification:

```bash
go run github.com/go-task/task/v3/cmd/task@latest ci
go run github.com/go-task/task/v3/cmd/task@latest coverage
```

## 2. Deepen Read Command Registration

Problem: read-only command wiring still repeats domain callbacks and has many `//nolint:dupl` suppressions.

Success criteria:

- Map one repeated change path before adding a seam.
- Remove more code than the new abstraction adds.
- Keep Cobra help, JSON output, and test scenario names unchanged.

Verification:

```bash
go test -race -shuffle=on -count=1 ./internal/cli ./internal/client
go run github.com/go-task/task/v3/cmd/task@latest ci
```

## 3. Mark Deprecated Roadmap Surface As Legacy

Problem: Linear marks Roadmap as deprecated in favor of Initiative, but the CLI exposes Roadmap reads as
ordinary commands.

Success criteria:

- Keep existing Roadmap read compatibility.
- Make help/docs bias new planning workflows toward Initiative.
- Do not add Roadmap writes without an explicit guard model.

Verification:

```bash
go test -race -shuffle=on -count=1 ./internal/cli
go run github.com/go-task/task/v3/cmd/task@latest ci
```
