# OAuth failures have structured error codes

linctl will expose OAuth runtime failures through stable structured error codes so agents can branch on recovery paths without scraping prose. The initial OAuth error set includes `AUTH_NOT_CONFIGURED`, `AUTH_TOKEN_EXPIRED`, `AUTH_REFRESH_FAILED`, `AUTH_REAUTH_REQUIRED`, `MISSING_SCOPE`, `AUTH_ACTOR_MISMATCH`, and `AUTH_TARGET_MISMATCH`. `AUTH_TARGET_MISMATCH` is for auth readiness checks; guarded writes keep the existing target mismatch semantics.

**Status**: accepted
