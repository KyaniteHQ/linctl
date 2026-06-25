<!-- Keep changes small, typed, and schema-aligned. See CONTRIBUTING.md. -->

## What changed

<!-- One or two sentences. For a new command, name the command and its GraphQL backing. -->

## Why

<!-- The problem this solves. Link any issue with "Closes #NNN". -->

## Safety

<!-- If this adds or touches a write, confirm it goes through the target guard. -->

- [ ] Reads stay free; any new write is target-pinned and fails closed on mismatch (no bypass flag)
- [ ] Operations are schema-aligned (genqlient); `internal/client/generated.go` regenerated and committed if operations changed
- [ ] No token value is printed or logged

## Checks

- [ ] `task ci` passes (generate-check → vet → test → build → lint → actionlint → vuln)
- [ ] `task coverage` is 100.0% for hand-written product code
- [ ] Docs/skill reference regenerated if the command surface changed
