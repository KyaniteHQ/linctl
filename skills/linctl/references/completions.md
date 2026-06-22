# Shell completions

linctl ships shell completions on two levels.

## Static (built in by Cobra)

The root `completion` command emits a completion script for your shell. It needs
no token and no network:

```bash
linctl completion bash   > /etc/bash_completion.d/linctl   # bash
linctl completion zsh    > "${fpath[1]}/_linctl"           # zsh
linctl completion fish   > ~/.config/fish/completions/linctl.fish
linctl completion powershell | Out-String | Invoke-Expression
```

These complete command names, subcommands, and flag names from the binary's
command tree, so they stay correct as commands are added.

## Dynamic (live value completion)

When the static script is installed, these values complete from live Linear
reads (they degrade silently to no suggestions when no token or target is
configured — completion never errors):

- `--team` → team keys (via the team list read)
- `--project` → project ids, annotated with the project name (via the project list read)
- `--state` on `issue list`, `issue create`, and `issue update` → workflow state
  types (via the workflow-state list read)
- positional argument of `team get` → team keys
- positional argument of `project get` → project ids

Dynamic completion is read-only and never touches the write guard.
