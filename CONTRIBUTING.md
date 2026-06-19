# Contributing

`linctl` is a schema-aligned Go CLI for Linear. Keep changes small, typed, and backed by the generated
GraphQL schema.

## Local Checks

```bash
go generate ./...
git diff --exit-code -- internal/client/generated.go
go build ./...
go vet ./...
go test -race -shuffle=on -count=1 ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --timeout 5m ./...
go run github.com/go-task/task/v3/cmd/task@latest ci
```

Run live integration tests only with a disposable Linear API token:

```bash
LINCTL_TEST_TOKEN=<token> go test -count=1 -tags=integration ./internal/client
```

The full live smoke harness is:

```bash
task live-smoke
```

Never run write tests against real project data. Test resources must use a `linctl-it-<runid>` prefix and
be archived during cleanup.

## Schema Changes

Refresh the vendored Linear schema before adding or changing GraphQL operations:

```bash
./scripts/refresh-schema.sh
go generate ./...
```

Generated code must be committed with the operation that requires it.

## Releases

Release builds are produced by GoReleaser from `v*` tags:

```bash
goreleaser release --snapshot --clean
```

The release workflow publishes GitHub artifacts and updates the `KyaniteHQ/homebrew-linctl` tap cask.
The tap token must be provided as `HOMEBREW_TAP_GITHUB_TOKEN`.
