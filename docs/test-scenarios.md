# Test Scenarios

This file defines the repeatable scenario set for the coverage, logging, and live-smoke goal.

## Method

N is 5 for the local regression streak. Each local scenario runs under the
same unit-test conditions: fake GraphQL responses, no live Linear writes, and no
secret material in inputs or logs.

Success is pass/fail:

- The expected behavior is asserted by an automated test.
- Important failure paths produce actionable errors or diagnostic logs.
- Logs must not include Linear tokens, request bodies, response bodies, or fixture user data.
- The default suite remains fast enough for local iteration.
- Live smoke uses the same pass/fail rule, but runs only when a disposable
  Linear token is available.

## Scenarios

1. Target-pinned issue write
   - Success: creates, updates, comments, and closes only after target resolution and project/team checks.
   - Evidence: `go test ./internal/client`, `Test_ClientWriteScenarios_guard_writes_and_report_results`.

2. Target-pinned project write
   - Success: creates, updates, and archives only after target resolution and project/team checks.
   - Evidence: `go test ./internal/client`, `Test_ClientWriteScenarios_guard_writes_and_report_results`.

3. Read-only issue/project inspection
   - Success: list/get/member commands map generated GraphQL responses into compact models with pagination data.
   - Evidence: `go test ./internal/client`, `Test_ClientReadScenarios_return_compact_lists_details_and_members`.

4. Transport retry and diagnostics
   - Success: 429 responses retry, terminal failures remain errors, and diagnostics include attempt/status without secrets or bodies.
   - Evidence: `go test ./internal/client`, `Test_Transport_retries_429_with_retry_after_when_present`,
     `Test_Transport_logs_decode_failures_without_response_body`.

5. Production-record classification boundary
   - Success: no repo-owned production-record dataset or classification logic is present; generated Linear schema references are external API surface, not linctl-owned production records.
   - Evidence: `rg -n "production records|classification|classify|record|allowed definition|production" README.md CONTEXT.md docs internal scripts test -S`.

## Current Outcome

All five local scenarios pass under the method above. The complete local suite also passes with `go test ./...`.

Coverage is enforced with `task coverage`, which runs uncached tests and excludes generated GraphQL code and the thin process entrypoint from the hand-written behavior metric. The enforced hand-written statement coverage target is 100.0%.

## Live Smoke

Run the complete live smoke suite with:

```bash
task live-smoke
```

The command requires a disposable Linear token in `LINCTL_TEST_TOKEN`, `LINCTL_TOKEN`, or `LINEAR_API_KEY`.
It builds a temporary `linctl` binary, smoke-tests read-only CLI commands, then runs the integration-tagged
client round trips. Write smoke tests create `linctl-it-<runid>` resources and archive them during cleanup.

Current credential state: live smoke is blocked unless one of the token variables above is injected. The
local readiness check is `env -u LINCTL_TEST_TOKEN -u LINCTL_TOKEN -u LINEAR_API_KEY bash scripts/live-smoke.sh`,
which must exit 2 with a missing-token message and without printing secret values.
