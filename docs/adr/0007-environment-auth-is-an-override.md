# Environment auth is an override

linctl will prefer saved local auth state for OAuth app credentials and tokens, while allowing environment-supplied OAuth material as a non-persistent automation override. Environment auth accepts OAuth app material or OAuth token state only; it does not accept personal API keys. This keeps CI, containers, and secret managers workable without making process environment variables the product's primary auth configuration surface or reintroducing API-key auth through the side door.

**Status**: accepted
