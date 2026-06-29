# Client secret is flow-specific

linctl will allow PKCE browser login without a client secret when Linear permits it, but `linctl auth app` will require a client secret for the client-credentials flow. This keeps user authorization flexible for public-client style CLI login while preserving explicit app authentication for headless OAuth app authorization.

**Status**: accepted
