# Demo assets

Source for the terminal demo shown in the project README. The tape runs the **real
`linctl` binary** against a Linear workspace, so the output is genuine — reads are free,
and the guarded write is refused when the active token does not resolve to the pinned target.

## Files

| File | Purpose |
| --- | --- |
| `demo.tape` | [vhs](https://github.com/charmbracelet/vhs) script for the recording |
| `render.sh` | builds `linctl` and runs vhs → `../docs/assets/demo.{gif,mp4}` |
| `.linctl.toml` | pins the demo workspace (you provide this; it is gitignored) |

## Render

Prerequisites: [vhs](https://github.com/charmbracelet/vhs), Go, and a **throwaway** Linear
workspace. Pin that workspace in `demo/.linctl.toml` (`org_id`, `team_key`, `team_id`), then:

```bash
LINEAR_API_KEY=<demo-workspace-token> ./render.sh
```

The token is read from the environment and never printed. The successful `issue create`
lands a real issue in the demo workspace, so use a disposable one.

## Storyboard

1. `issue list` — reads need no pinned target.
2. `cat .linctl.toml` + `linctl target` — the pin, then the live token re-resolved against it.
3. `issue create` — a guarded write into the pinned target succeeds.
4. `issue create --team STG` — the same write aimed at a team the token does not own is
   refused with `{"error_code":"TARGET_MISMATCH"}` and a non-zero exit. The flags set the
   pinned target; they do not relax the guard. There is no bypass flag.
