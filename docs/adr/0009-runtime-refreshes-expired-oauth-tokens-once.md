# Runtime refreshes expired OAuth tokens once

linctl will automatically refresh an expired OAuth access token, persist the rotated token state, and retry the original Linear request once. If refresh fails, the command returns a structured authentication error that points the user back to `linctl auth login` or `linctl auth app` instead of making normal commands require a manual refresh step.

**Status**: accepted
