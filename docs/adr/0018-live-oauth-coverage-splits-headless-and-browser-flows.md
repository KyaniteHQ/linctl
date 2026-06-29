# Live OAuth coverage splits headless and browser flows

linctl will automatically live-test the client-credentials OAuth path and initially verify browser authorization-code login through a manual or semi-automated browser smoke. This maximizes early live regression detection for unattended automation while avoiding brittle CI dependence on an interactive consent screen.

**Status**: accepted
