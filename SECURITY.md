# Security Policy

## Supported versions

linctl has a rolling release. Only the latest tagged release gets a security fix.

| Version        | Supported          |
| -------------- | ------------------ |
| latest release | :white_check_mark: |
| older releases | :x:                |

## Report a vulnerability

Report a security problem privately through
[private vulnerability reporting](https://github.com/KyaniteHQ/linctl/security/advisories/new)
on GitHub. Open the **Security** tab of the repository, then select **Report a vulnerability**.

Do **not** open a public issue for a suspected vulnerability.

We try to acknowledge a report in 72 hours. We give a remediation schedule after the triage.

linctl makes guarded, target-pinned writes against Linear. Add the resolved target and the pinned
target to your report. Add the exact command that you ran, and redact each token value. We can
then reproduce the problem safely.
