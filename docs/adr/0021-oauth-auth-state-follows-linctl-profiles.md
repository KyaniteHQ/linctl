# OAuth auth state follows linctl profiles

linctl will optimize OAuth auth state for one default profile while keeping named profile selection through the existing `--profile` and config mechanisms. OAuth does not introduce a separate profile system, preserving the existing relationship between profile, pinned target, and runtime behavior.

**Status**: accepted
