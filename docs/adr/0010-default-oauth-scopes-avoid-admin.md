# Default OAuth scopes avoid admin

linctl will request `read`, `write`, `issues:create`, and `comments:create` as the initial default OAuth scope set, and will not request `admin` by default. This keeps the first OAuth switch broad enough for existing read commands and guarded writes while preserving the safety boundary around admin-level Linear operations.

**Status**: accepted
