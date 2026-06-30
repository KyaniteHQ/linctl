# Contributing

`linctl` is a schema-aligned Go CLI for Linear. Keep changes small, typed, and backed by the generated
GraphQL schema.

## Local Checks

```bash
go generate ./...
git diff --exit-code -- internal/client/generated.go
go mod download
go mod verify
go mod tidy -diff
go tool task fmt-check
go build $(bash scripts/go-packages.sh)
go vet $(bash scripts/go-packages.sh)
go test -race -shuffle=on -count=1 $(bash scripts/go-packages.sh)
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --timeout 5m $(bash scripts/go-packages.sh)
shellcheck scripts/*.sh
go tool govulncheck $(bash scripts/go-packages.sh)
go tool task actionlint
go tool task ci
go tool task coverage
```

`go tool task ci` is the local review gate. It verifies module checksums and
tidy state, formatting, generated artifacts, domain language, vet, tests, build,
CLI smoke output, golangci-lint, ShellCheck, actionlint, and govulncheck without
modifying source files. It also validates local GraphQL operations and
`docs/linear-api-coverage.md` against the upstream Linear SDK checkout. The
shared source contract is:

- Remote: `https://github.com/linear/linear.git`
- Default checkout: run-local temporary checkout managed by Taskfile
- Default ref: `master`
- Reusable checkout: `LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear`
- Override ref: `LINCTL_LINEAR_SDK_REF=<branch-or-tag>`

Prepare or refresh an explicit reusable checkout with:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear \
go tool task linear-sdk-upstream-checkout
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

Refresh the vendored Linear schema before adding or changing GraphQL operations:

```bash
npm ci
./scripts/refresh-schema.sh
go generate ./...
```

`scripts/refresh-schema.sh` uses the repo-managed `graphql` dependency from
`package-lock.json`. Set `LINCTL_OAUTH_ACCESS_TOKEN` for the command, but never
print or paste the token value into logs. Generated code must be committed with
the operation that requires it.

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
