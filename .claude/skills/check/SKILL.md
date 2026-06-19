---
name: check
description: Run the full linctl pre-PR local gate (generate-check, vet, test, 100% coverage, lint) and report pass/fail per gate. Use before opening a PR or when asked to verify linctl changes locally.
---

# check — linctl pre-PR gate

Run the complete local gate for the linctl repo and report each step's result. This is the gate CI plus the coverage script enforce; run it before any PR.

`task` may not be installed — use `go run github.com/go-task/task/v3/cmd/task@latest` as the `task` command, or call the underlying commands directly (shown below).

Run these in order, capturing pass/fail and key output for each. Do not stop at the first failure — run all, then report.

1. **Generated client is committed** — `go generate ./... && git diff --exit-code -- internal/client/generated.go`. Fails if generated code drifted; regenerate and commit.
2. **Vet** — `go vet ./...`
3. **Tests** — `go test -race -shuffle=on -count=1 ./...`
4. **Coverage (100% hand-written)** — `bash scripts/coverage.sh`. Fails if any hand-written statement is uncovered; generated.go and cmd/linctl/main.go are excluded.
5. **Lint** — `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --timeout 5m ./...`

Do NOT run `task live-smoke` here — it needs a disposable Linear token and hits the live API. Mention it separately if the change touches transport/client behavior.

Report in this shape:

```
Gate results:
- generate-check: PASS
- vet: PASS
- test: PASS
- coverage: FAIL — internal/cli/foo.go:42 uncovered
- lint: PASS
```

Fix failures at the root, re-run the failed gate, and only report green once every gate passes.
