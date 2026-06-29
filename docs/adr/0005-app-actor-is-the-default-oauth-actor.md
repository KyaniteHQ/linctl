# App actor is the default OAuth actor

linctl will default OAuth authorization to the Linear app actor. Browser login builds authorization URLs with `actor=app` unless the user explicitly passes `--actor user`, and client-credentials auth relies on Linear's app actor token behavior. This favors agent-safe attribution and service behavior over personal attribution, while preserving an intentional path for commands that must appear as the approving Linear user.

**Status**: accepted
