# Demo assets

Source for the terminal demo shown in the project README: reads are free, a guarded write into
the pinned target lands, and the same write aimed at a team the credential cannot reach is
refused.

There are two ways to render it. Both produce `../docs/assets/demo.{gif,mp4}`.

## Files

| File | Purpose |
| --- | --- |
| `demo.tape` | [vhs](https://github.com/charmbracelet/vhs) script for the recording |
| `render-fixture.sh` | renders from `fixture/linctl`. **No Linear account, no credentials.** |
| `fixture/linctl` | replays the real binary's output shapes against invented ids |
| `render.sh` | builds and records the **real binary** against a disposable Linear org |
| `.linctl.toml` | pins the demo org/team for `render.sh` (you provide it; gitignored) |

## Render from the fixture (default)

Only needs [vhs](https://github.com/charmbracelet/vhs):

```bash
./render-fixture.sh
```

The ids are invented, so nothing is created in Linear and nobody needs a throwaway organization
to refresh the README. The output shapes are copied from the real binary. **If you change how
`issue list`, `issue create`, `target`, or the error envelope print, update `fixture/linctl` to
match**, or re-record with `render.sh` below.

## Render against real Linear

Prerequisites: vhs, Go, and a **throwaway** Linear org/team. Pin that target in
`demo/.linctl.toml` (`org_id`, `team_key`, `team_id`), then:

```bash
LINCTL_OAUTH_ACCESS_TOKEN=<demo-oauth-access-token> ./render.sh
```

The OAuth access token is read from the environment and never printed. The successful `issue
create` lands a real issue in the demo target, so use a disposable one.

## Storyboard

1. `issue list` — reads need no pinned target.
2. `cat .linctl.toml` + `linctl target` — the pin, then the live auth re-resolved against it.
3. `issue create` — a guarded write into the pinned target succeeds.
4. `issue create --team STG` — the same write aimed at a team the credential cannot reach is
   refused with `{"error_code":"TARGET_MISMATCH"}` and a non-zero exit. `--team` sets the
   pinned target; it does not relax the guard. There is no bypass flag.
