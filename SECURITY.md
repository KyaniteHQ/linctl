# Security Policy

## Supported versions

linctl is distributed as a rolling release. Only the latest tagged release
receives security fixes.

| Version        | Supported          |
| -------------- | ------------------ |
| latest release | :white_check_mark: |
| older releases | :x:                |

## Reporting a vulnerability

Please report security issues privately through GitHub's
[private vulnerability reporting](https://github.com/KyaniteHQ/linctl/security/advisories/new)
(repository **Security** tab → **Report a vulnerability**).

Do **not** open a public issue for a suspected vulnerability.

We aim to acknowledge a report within 72 hours and to provide a remediation
timeline after triage. Because linctl performs guarded, target-pinned writes
against Linear, please include the resolved/pinned target details and the exact
command invocation (with any token value redacted) so the report can be
reproduced safely.
