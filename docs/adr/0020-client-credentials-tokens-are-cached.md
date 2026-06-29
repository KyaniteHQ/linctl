# Client-credentials tokens are cached

`linctl auth app` will save the client-credentials access token with its expiry and reuse it until expiry or an authorization failure. The runtime may request a new app actor token once and retry the original Linear request, avoiding per-command token churn while still recovering from expired or invalidated app tokens.

**Status**: accepted
