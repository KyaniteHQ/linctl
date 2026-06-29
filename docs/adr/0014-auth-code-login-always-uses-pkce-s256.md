# Auth-code login always uses PKCE S256

`linctl auth login` will always use PKCE with an S256 code challenge for OAuth authorization-code login. This makes the browser and manual callback paths safer for a CLI without depending only on a stored client secret to bind the authorization code exchange.

**Status**: accepted
