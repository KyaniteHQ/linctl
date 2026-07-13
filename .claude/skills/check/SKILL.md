---
name: check
description: Run the full linctl pre-PR local gate (generate-check, vet, test, 100% coverage, lint) and report pass/fail per gate. Use before opening a PR or when asked to verify linctl changes locally.
---

# check — linctl pre-PR gate

Run the complete local gate for the linctl repo and report each step's result. `go tool task ci` and `go tool task coverage` are the gate; run them before any PR.

`task` is a pinned Go tool dependency — always invoke it as `go tool task`, never fetch an unpinned `task` binary at run time.

Run these in order, capturing pass/fail and key output for each. Do not stop at the first failure — run all, then report.

1. **`go tool task ci`** — deps-check, fmt-check, generate-check (generated client, command reference, upstream schema/coverage-ledger checks), domain-language-check, browser-login-smoke-check, vet, test, build, smoke-run, lint, shellcheck, actionlint, vuln.
2. **`go tool task coverage`** (separate from `ci`) — enforces 100% hand-written statement coverage; generated.go, cmd/linctl/main.go, and scripts/ are excluded. Fails with the uncovered file:line.

The upstream schema/coverage-ledger checks in `ci` need network access to `github.com/linear/linear` unless a reusable checkout is already prepared. Use `LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear` to point at one, and `LINCTL_LINEAR_SDK_OFFLINE=1` to skip the refresh fetch against it (fails loud if the checkout doesn't exist).

Do NOT run `task live-smoke` here — it needs a disposable Linear token and hits the live API. Mention it separately if the change touches transport/client behavior.

Report in this shape:

```
Gate results:
- ci: PASS
- coverage: FAIL — internal/cli/foo.go:42 uncovered
```

Fix failures at the root, re-run the failed gate, and only report green once every gate passes.
