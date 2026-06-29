# OAuth live gate is explicit

linctl will add a separate `task live-oauth` gate for OAuth fixture-backed live tests, and `task live-smoke` may call it only when the OAuth fixture environment is present. The fixture environment is `LINCTL_OAUTH_CLIENT_ID`, `LINCTL_OAUTH_CLIENT_SECRET`, `LINCTL_OAUTH_REDIRECT_URI`, `LINCTL_OAUTH_SCOPES`, and `LINCTL_OAUTH_EXPECTED_ACTOR=app`; `LINCTL_OAUTH_CLIENT_SECRET` is required for the initial live gate because the automatic path verifies client-credentials auth. This keeps OAuth failures diagnosable and secret-aware while still maximizing live regression detection in prepared environments.

**Status**: accepted
