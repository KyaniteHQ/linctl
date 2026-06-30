# Contributing

`linctl` is a schema-aligned Go CLI for Linear. Keep changes small, typed, and backed by the generated
GraphQL schema.

## Local Checks

```bash
go generate ./...
git diff --exit-code -- internal/client/generated.go
go build $(bash scripts/go-packages.sh)
go vet $(bash scripts/go-packages.sh)
go test -race -shuffle=on -count=1 $(bash scripts/go-packages.sh)
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --timeout 5m $(bash scripts/go-packages.sh)
go tool govulncheck $(bash scripts/go-packages.sh)
go run github.com/go-task/task/v3/cmd/task@latest actionlint
go run github.com/go-task/task/v3/cmd/task@latest ci
go run github.com/go-task/task/v3/cmd/task@latest coverage
```

`task ci` also validates local GraphQL operations and `docs/linear-api-coverage.md`
against the upstream Linear SDK checkout. The shared source contract is:

- Remote: `https://github.com/linear/linear.git`
- Default checkout: `/tmp/linctl-upstream-linear`
- Default ref: `master`
- Override path: `LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear`
- Override ref: `LINCTL_LINEAR_SDK_REF=<branch-or-tag>`

Prepare or refresh the default checkout with:

```bash
go run github.com/go-task/task/v3/cmd/task@latest linear-sdk-upstream-checkout
```

If the default path is unavailable, use an override:

```bash
LINCTL_LINEAR_SDK_UPSTREAM=/path/to/linear \
go run github.com/go-task/task/v3/cmd/task@latest coverage-ledger-check
```

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
go run github.com/go-task/task/v3/cmd/task@latest live-smoke
```

For the project Infisical setup, fixture secrets live under `/linctl`, not the
root secret path. Use the pinned aliases so the folder is not easy to forget:

```bash
go run github.com/go-task/task/v3/cmd/task@latest live-oauth-infisical
go run github.com/go-task/task/v3/cmd/task@latest live-smoke-infisical
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

```bash
git tag -a vX.Y.Z -m "vX.Y.Z" && git push origin vX.Y.Z
```

To dry-run the build locally without publishing anything (no tag, no upload):

```bash
goreleaser release --snapshot --clean
```
