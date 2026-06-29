# Auth status is live

`linctl auth status` will refresh expired OAuth tokens when possible and verify the current OAuth state against Linear before reporting readiness. This makes the command slower than a local-only inspection, but it keeps status output truthful about token validity, authorization, actor state, scopes, and target readiness.

**Status**: accepted
