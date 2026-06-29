# Live OAuth tests use a pre-created app fixture

linctl live OAuth tests will use a dedicated pre-created Linear OAuth application supplied through secrets instead of dynamically creating, updating, or rotating OAuth apps during the test run. This keeps live regression coverage focused on OAuth behavior while avoiding brittle admin-surface cleanup and secret-management risk.

**Status**: accepted
