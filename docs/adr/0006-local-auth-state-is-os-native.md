# Local auth state is OS-native

linctl will store OAuth client and token material as machine-local auth state outside repo configuration, using the operating system's user-specific config or state location with strict permissions where the platform supports them. Repo `.linctl.toml` remains for the pinned target and other shareable project settings, not OAuth secrets.

**Status**: accepted
