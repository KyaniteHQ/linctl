# OAuth logout preserves app configuration

`linctl auth logout` will revoke Linear OAuth tokens when possible and remove local token state, but it will preserve OAuth app client configuration unless the user passes an explicit forget option. This separates ending an authorization from removing the local app setup needed to authorize again.

**Status**: accepted
