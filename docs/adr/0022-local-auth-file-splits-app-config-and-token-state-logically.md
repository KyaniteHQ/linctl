# Local auth file splits app config and token state logically

linctl may store OAuth app configuration and token state in one physical local auth file for the MVP, but the model keeps them logically separate. This keeps local setup simple while preserving command semantics such as `auth logout` deleting token state without forgetting OAuth app configuration.

**Status**: accepted
