# Setup Gap Log

This log records clean-clone setup assumptions found while replaying the README from disposable
environments.

## Attempt 1

Environment: fresh `ubuntu:24.04` container.

Command:

```bash
git clone https://github.com/KyaniteHQ/linctl.git /tmp/linctl
cd /tmp/linctl
go generate ./...
```

Gap:

- `git` was not installed, and the README did not state clean-machine prerequisites.
- The README also did not state the required Go version before development commands.

Fix:

- Added `README.md` clean Linux setup commands for `ca-certificates`, `curl`, `git`, `tar`, Go `1.26.4`,
  cloning, and running `linctl` from source.
- Added `go build ./...` to the documented development gate.

## Attempt 2

Environment: fresh `ubuntu:24.04` container.

Command:

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
go generate ./...
git diff --exit-code -- internal/client/generated.go
go build ./...
go vet ./...
go test -race -shuffle=on -count=1 ./...
```

Gap:

- `go test -race` failed with `-race requires cgo`; the README did not state the C toolchain required
  by Go's race detector.

Fix:

- Added `build-essential` to the clean Linux prerequisite install command.

## Attempt 3

Environment: fresh `ubuntu:24.04` container.

Command:

```bash
apt-get update
apt-get install -y build-essential ca-certificates curl git tar
curl -fsSL https://go.dev/dl/go1.26.4.linux-amd64.tar.gz -o /tmp/go.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf /tmp/go.tar.gz
export PATH="/usr/local/go/bin:$PATH"
git clone https://github.com/KyaniteHQ/linctl.git
cd linctl
go run ./cmd/linctl usage
go run ./cmd/linctl --version
go generate ./...
git diff --exit-code -- internal/client/generated.go
go build ./...
go vet ./...
go test -race -shuffle=on -count=1 ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run --timeout 5m ./...
```

Gap:

- The clean run passed, but the lint command used `@latest`, which quietly made the documented gate depend
  on whatever golangci-lint version is current at execution time.

Fix:

- Pinned the README lint command and CI lint action to golangci-lint `v2.12.2`.
