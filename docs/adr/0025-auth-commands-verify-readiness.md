# Auth commands verify readiness

`linctl auth app` and other auth setup/status commands will verify OAuth readiness against Linear before reporting success. Saving local auth state is not enough; linctl must prove the expected actor, scopes, and resolved target before telling the user the authorization is usable. If readiness verification fails, the token is not saved as active auth state.

**Status**: accepted
