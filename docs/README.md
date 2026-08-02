# linctl documentation

## Start here

| | |
| --- | --- |
| [Quickstart](quickstart.md) | Install linctl, authenticate, pin a target, and land a write. About five minutes. |
| [linctl and the Linear MCP server](why-not-mcp.md) | The token cost, the write boundary, and what MCP does better. |
| [`../README.md`](../README.md) | What linctl is, the command surface, and the output flags. |
| [`../CONTEXT.md`](../CONTEXT.md) | The vocabulary. Read this before you name anything in the codebase. |

## Reference

| | |
| --- | --- |
| [Architecture baseline (ADR 0001)](adr/0001-linctl-architecture-baseline.md) | Why the guard works this way, and which alternatives the project rejected. |
| [`../skills/linctl/SKILL.md`](../skills/linctl/SKILL.md) | How to drive linctl from an AI agent. |
| [`../skills/linctl/references/commands.md`](../skills/linctl/references/commands.md) | **Generated.** Every command, with its usage line and its flags. |
| [`../skills/linctl/references/guarded-writes.md`](../skills/linctl/references/guarded-writes.md) | **Generated.** Every command that changes Linear, and which ones need `--org-wide`. |
| [`../skills/linctl/references/json-output.md`](../skills/linctl/references/json-output.md) | The stable JSON shapes, for any tool that parses linctl output. |
| [`../CONTRIBUTING.md`](../CONTRIBUTING.md) | The local gate, schema changes, and releases. |
| [`../SECURITY.md`](../SECURITY.md) | How to report a vulnerability. |

## Internal engineering ledgers

These are working documents for maintainers. They are large and exhaustive, and nobody reads them
from start to end. You do not need them to use linctl.

| | |
| --- | --- |
| [`internal/domain-map.md`](internal/domain-map.md) | Each command mapped to the GraphQL operation behind it. |
| [`internal/linear-api-coverage.md`](internal/linear-api-coverage.md) | **Generated.** The coverage of the Linear API surface. `scripts/linear_api_coverage.go` regenerates it, and CI fails on drift. Never edit it by hand. |
| [`internal/test-scenarios.md`](internal/test-scenarios.md) | The named test flows, tied to the Go test functions that prove them. |
