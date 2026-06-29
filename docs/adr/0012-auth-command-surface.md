# Auth command surface

linctl will expose OAuth through `linctl auth configure`, `linctl auth login`, `linctl auth app`, `linctl auth status`, `linctl auth refresh`, and `linctl auth logout`. The command surface separates OAuth app setup, browser authorization, headless app authorization, operational inspection, explicit refresh, and logout so each behavior can be tested and documented as a stable public interface.

**Status**: accepted
