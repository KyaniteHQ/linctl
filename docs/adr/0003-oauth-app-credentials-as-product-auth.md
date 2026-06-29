# OAuth app credentials as product auth

linctl will make Linear OAuth app credentials the product authentication model and remove personal API keys from the product direction. Because this auth surface has not been released, there is no compatibility period for `LINCTL_TOKEN`, `LINEAR_API_KEY`, or config-file API keys. This trades the previous local/CI simplicity of API keys for app identity, installable integration semantics, and clearer agent attribution while keeping target-pinned guarded writes as the write-safety boundary.

**Status**: accepted
