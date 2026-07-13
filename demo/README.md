# Demo assets

Two separate things live here.

| Directory | What it is |
| --- | --- |
| this one (`demo.tape`, `render*.sh`) | the **terminal recording** in the README: real command output, captured with vhs |
| `readme-hero/` | the **explainer video**: a designed HyperFrames composition, 4K60 |
| `launch/` | the 17s deadpan launch spot, kept as-is |

## The terminal recording

Reads are free, a guarded write into the pinned target lands, and the same write aimed at a team
the credential cannot reach is refused.

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

## The explainer video (`readme-hero/`)

A [HyperFrames](https://hyperframes.heygen.com) composition: five scenes covering the problem,
the standing context cost of an MCP server, a write that lands, the same write refused, and the
honest limit. Deadpan, monospace, no narration, no music.

Only the source is tracked (`index.html`, `meta.json`, `hyperframes.json`, `package.json`).
Renders and snapshots are generated.

```bash
cd readme-hero
npx hyperframes check      # lint, runtime, layout, motion, contrast
npx hyperframes preview    # watch it in the browser while editing
npx hyperframes render --resolution landscape-4k --fps 60 --quality high \
  --output renders/linctl-4k60.mp4
```

The composition is authored at 1920x1080; `--resolution landscape-4k` renders it at 3840x2160 by
raising Chrome's device pixel ratio, so there is one set of coordinates to reason about.

**Every number and string on screen is load-bearing.** The token counts come from
`scripts/mcp-token-measure.sh`, the terminal lines match the real binary's output shapes, and the
closing "confused agent, not a hostile one" is the guard's actual limit. If any of those change,
change the composition too.
