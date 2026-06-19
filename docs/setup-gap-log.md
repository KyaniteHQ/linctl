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
