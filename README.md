# linctl

`linctl` is a schema-aligned Go CLI for Linear.

It is built for agent-safe daily coordination: reads are lightweight, writes re-resolve the active Linear
token and fail closed unless the resolved org/team/project matches the pinned target.

## Install

### Clean Linux Machine

These commands start from a fresh Ubuntu 24.04 environment with no project tools installed:

```bash
apt-get update
apt-get install -y ca-certificates curl git tar

curl -fsSL https://go.dev/dl/go1.26.4.linux-amd64.tar.gz -o /tmp/go.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf /tmp/go.tar.gz
export PATH="/usr/local/go/bin:$PATH"

git clone https://github.com/KyaniteHQ/linctl.git
cd linctl

go run ./cmd/linctl usage
go run ./cmd/linctl --version
```

From source:

```bash
go install github.com/KyaniteHQ/linctl/cmd/linctl@latest
```

After the first tagged release:

```bash
brew install --cask KyaniteHQ/linctl/linctl
```

## Configure

Create `.linctl.toml` in a repo:

```toml
[target]
org_id = "linear-org-id"
team_key = "LIT"
team_id = "linear-team-id"
project_id = "optional-linear-project-id"
```

Inject credentials with `LINCTL_TOKEN` or `LINEAR_API_KEY`; do not commit tokens.

## Usage

```bash
linctl usage
linctl target --json
linctl current --json
linctl issue usage
linctl project usage
```

Issue and project writes require a pinned target. Team-scoped creates compare org/team; resource-scoped
updates and archives resolve the resource first and compare the pinned project when configured.

## Development

After following the clean-machine setup above:

```bash
go generate ./...
git diff --exit-code -- internal/client/generated.go
go build ./...
go vet ./...
go test -race -shuffle=on -count=1 ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --timeout 5m ./...
```

The temporary integration fixture is configured in `test/integration-config.json`. Inject
`LINCTL_TEST_TOKEN` from secret storage only when running live integration tests:

```bash
LINCTL_TEST_TOKEN=<token> go test -count=1 -tags=integration ./internal/client
```
