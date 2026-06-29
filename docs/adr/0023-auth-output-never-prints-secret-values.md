# Auth output never prints secret values

linctl auth commands, JSON output, debug logs, diagnostics, and tests will never print OAuth token or client secret values. Auth output may report `set`, `missing`, expiry, scopes, actor, token type, and other non-secret metadata, with secret-bearing fields redacted before any output path.

**Status**: accepted
