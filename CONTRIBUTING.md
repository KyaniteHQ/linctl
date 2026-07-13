# Contributing

`linctl` is a schema-aligned Go CLI for Linear. Keep changes small, typed, and backed by the generated
GraphQL schema.

## Local Checks

`go tool task ci` is the single local review gate. Run it before opening a PR:

```bash
go tool task ci
```

It runs, in order: `deps-check` (module checksums and tidy state), `fmt-check`,
`generate-check` (generated GraphQL client, command reference, and the
upstream schema/coverage-ledger checks), `domain-language-check`,
`browser-login-smoke-check`, `vet`, `test`, a build, `smoke-run`, `lint`,
`shellcheck`, `actionlint`, and `vuln`. None of these modify source files.

`go tool task coverage` is separate from `ci` and enforces 100% hand-written
statement coverage. Run it on its own for product code changes:

```bash
go tool task coverage
```

### Network & offline

`go tool task ci` validates local GraphQL operations and
`docs/linear-api-coverage.md` against the upstream `linear/linear` schema. The
shared source contract is:

- Remote: `https://github.com/linear/linear.git`
- Default checkout: run-local temporary checkout managed by Taskfile
- Default ref: `master`
- Reusable checkout: `LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear`
- Override ref: `LINCTL_LINEAR_SDK_REF=<branch-or-tag>`
- Skip the refresh fetch against an existing reusable checkout: `LINCTL_LINEAR_SDK_OFFLINE=1`
  (requires `LINCTL_LINEAR_SDK_UPSTREAM` to already point at a checkout; fails loud otherwise)

Prepare or refresh an explicit reusable checkout with:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear \
go tool task linear-sdk-upstream-checkout
```

Run the gate offline once that checkout exists:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear \
LINCTL_LINEAR_SDK_OFFLINE=1 \
go tool task ci
```

To skip Taskfile and run the helper scripts directly, pass `-upstream`:

```bash
go run scripts/linear_graphql_operation_check.go -upstream /path/to/linear
go run scripts/linear_api_coverage*.go -upstream /path/to/linear
```

GitHub-only checks remain separate from `go tool task ci`: dependency review runs only
on pull requests, coverage stays in `go tool task coverage`, and live OAuth/integration
checks require disposable fixture credentials.

Run live integration tests only with a disposable OAuth app fixture:

```bash
LINCTL_OAUTH_CLIENT_ID=<client-id> \
LINCTL_OAUTH_CLIENT_SECRET=<client-secret> \
LINCTL_OAUTH_REDIRECT_URI=http://127.0.0.1:8765/callback \
LINCTL_OAUTH_SCOPES=read,write,issues:create,comments:create \
LINCTL_OAUTH_EXPECTED_ACTOR=app \
go test -count=1 -tags=integration ./internal/client
```

The full live smoke harness is:

```bash
go tool task live-smoke
```

For the project Infisical setup, fixture secrets live under `/linctl`, not the
root secret path. Use the pinned aliases so the folder is not easy to forget:

```bash
go tool task live-oauth-infisical
go tool task live-smoke-infisical
```

Never run write tests against real project data. Test resources must use a `linctl-it-<runid>` prefix and
be archived during cleanup.

## Schema Changes

A nightly `schema-drift` job (`.github/workflows/integration.yml`) reports semantic drift
between the vendored `internal/client/schema.graphql` and the upstream `linear/linear` SDK
schema: types, fields, and enum values added or removed, and field type changes on types
present in both. It is read-only and not part of `go tool task ci`, so it never blocks a PR;
it exists to flag when the vendored snapshot has fallen behind. Run it locally with:

```bash
go tool task schema-drift-check
```

When the nightly job fails, refresh the schema and land the refresh as its own commit:

```bash
npm ci
LINCTL_OAUTH_ACCESS_TOKEN=<token> bash scripts/refresh-schema.sh
go generate ./...
go tool task ci
```

Review the generated diff before committing. Drift confined to schema areas linctl doesn't
use is safe to batch into a routine refresh; drift touching fields or types the current
operations depend on is urgent, though `go tool task graphql-operation-check` would already
be failing in that case.

Refresh the vendored Linear schema before adding or changing GraphQL operations:

```bash
npm ci
./scripts/refresh-schema.sh
go generate ./...
```

`scripts/refresh-schema.sh` uses the repo-managed `graphql` dependency from
`package-lock.json` and requires Node 22 or newer. Set
`LINCTL_OAUTH_ACCESS_TOKEN` for the command, but never print or paste the token
value into logs. Generated code must be committed with the operation that
requires it.

## Releases

A release is triggered by pushing a `v*` tag. The release workflow then runs GoReleaser to
publish the GitHub artifacts (archives, SBOMs, `checksums.txt`, and a keyless cosign sigstore
bundle) and update the `KyaniteHQ/homebrew-linctl` tap cask. The tap token must be provided as
`HOMEBREW_TAP_GITHUB_TOKEN`.

Before creating or pushing the tag, run the local non-publishing release preflight:

```bash
go tool task release-preflight
```

The preflight runs the local CI gate, statement coverage, and `goreleaser check`.
It does not create a tag, push to Git, publish a release, or require release
secrets. If you want the heavier local artifact build as a final manual check,
run the snapshot task. It also does not publish anything:

```bash
go tool task release-snapshot
```

Only create and push the release tag after the preflight passes:

```bash
git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z
```
