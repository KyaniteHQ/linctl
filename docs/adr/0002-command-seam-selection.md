# ADR 0002: command seam selection

**Status**: accepted

## Context

The CLI accumulated consumer-owned Command Ports whose production adapters forwarded one method to one exported `internal/client` function. Those interfaces did not isolate command decisions or external effects. They repeated the client API, added adapter code, and made tests prove forwarding instead of command behavior.

Some commands do need a seam. Issue creation may resolve optional template defaults before creating an issue. Issue listing and search choose a client operation after resolving command context. `next` lists, ranks, and may start an issue. Batch creation, OAuth, HTTP, and temporary-file handling also isolate behavior or effects that should be replaceable in focused tests.

Guarded-write safety is separate from command seam selection. The fail-closed target check belongs to `internal/client` so every caller crosses the same safety boundary.

## Decision

One-call commands call exported `internal/client` functions directly. Reads pass the runtime GraphQL client. Guarded writes pass the runtime GraphQL client and pinned target to the exported client function, which resolves and compares the target before mutation.

The CLI defines a narrow consumer-owned Command Port only when it isolates one of these:

- A multi-call decision or workflow owned by the command.
- An external effect whose replacement makes focused tests deterministic.

The retained issue workflow seams cover:

- Template resolution for issue creation.
- List dispatch between all-team and resolved-team reads.
- Search dispatch after target resolution.
- Next-issue selection and guarded start.

The separate batch `CreateIssues`, OAuth, HTTP, and download temporary-file seams remain. Their boundaries isolate batching or external effects rather than one-method forwarding.

Ports return command or domain types, not generated GraphQL response types. Tests use a port fake when exercising the decision or effect behind that port. One-call command behavior is tested through the command flow and client boundary instead of a forwarding-adapter test.

## Consequences

- One-call commands have no forwarding interface or production adapter to maintain.
- Multi-call decisions and external effects remain independently testable.
- Guarded writes remain fail-closed regardless of whether a command uses a port.
- Adding a port requires a concrete command decision or replaceable external effect, not anticipated future complexity.

## Rejected Alternatives

- **A Command Port for every command**: rejected because one-method ports duplicate exported client functions without isolating behavior.
- **No Command Ports**: rejected because multi-call decisions and external effects still need focused deterministic tests.
- **Write guarding in each command or adapter**: rejected because duplicated checks can drift or be skipped.
- **Generated GraphQL types in CLI ports**: rejected because generated transport shapes should not become command contracts.

## Code Alignment

- Guarded target resolution and mutation boundary: `internal/client/target.go`, `internal/client/write_guard.go`.
- Runtime GraphQL client and pinned target: `internal/cli/runtime.go`.
- Issue workflow seams: `internal/cli/issue_port.go`.
- Batch issue creation seam: `internal/cli/bulk.go`.
- OAuth seam and implementation: `internal/cli/auth.go`, `internal/oauth/`.
- HTTP and download temporary-file seams: `internal/cli/files.go`.
