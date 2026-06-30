# Brag Plan: linctl

## What is this app?
A Go CLI for Linear whose one trick is refusing to write to the wrong place: every guarded write re-resolves active OAuth auth against the org/team live and fails closed on mismatch, with no bypass flag. Reads are free.

## The angle
The whole video is one flat, unbothered demonstration: the same `issue create` command, run twice. Once it lands. Once it is refused, because the auth does not resolve to the team the command aimed at. No drama, no music swell at the "fail" — the joke of deadpan is that the guard treats a catastrophe (an agent writing to the wrong team) as a shrug. The refusal IS the product.

## Hook (first 2-3 seconds)
A single dry line on black: "An agent with the wrong auth writes to the wrong team." Beat. A dim second line: "Usually, nothing stops it." That sets up the only thing linctl does.

## Key moments (the middle)
- The pin appears (`.linctl.toml` → team LIT) and `linctl target` prints `confirmed true` — the live auth re-resolved against the pin.
- A guarded `issue create` lands in the pinned target: `LIT-15 ... [Backlog]`, calm green check.
- The same command aimed at `--team STG` returns `{"error_code":"TARGET_MISMATCH"}` in red, exit 1. The line "the write does not happen" stays on screen, unbothered.

## Outro / punchline
Tagline holds alone: **"reads everywhere, writes fail closed."** Then the install line: `go install github.com/KyaniteHQ/linctl/cmd/linctl@latest`.

## User flow worth showing
The CLI flow itself, recreated as a styled terminal card (not a raw recording — the authentic recording is the README GIF):
entry `linctl target` (confirmed) → action `linctl issue create` (lands) → same action at the wrong team (refused). Entry → key action → result, where the result is a refusal.

## Tone
- Preset: deadpan
- Creative direction: a security guard who has seen everything and is not impressed; the wrong-team write is handled like a misdelivered package.
- Interpretation: long holds, big empty space, one idea per scene, dry sparse sound. The refusal gets no sting — it is just another line of output. Restraint everywhere.

## Format: landscape — 1920x1080
## Duration: 17s

## Visual identity (from the project)
- Background: #0b0e14
- Accent: #e6b450 (amber cursor), #90e1c6 (cyan labels)
- Text: #bfbdb6; success green #7fd962; refusal red #ea6c73
- Display font: a monospace (JetBrains Mono / Berkeley Mono / system monospace fallback)
- Body font: same monospace — this is a CLI, everything is mono
- Strongest visual element: the `{"error_code":"TARGET_MISMATCH"}` line in red, and the green `confirmed true` / `LIT-15` lines. The styled terminal card is the recurring frame.

## Share copy (draft)
An AI agent with stale Linear auth can quietly create issues in the wrong team. linctl re-resolves auth on every write and fails closed. No bypass flag. `go install github.com/KyaniteHQ/linctl/cmd/linctl@latest`

## Audio direction
- Role: sparse professional accents over a very low bed
- Music: `happy-beats-business-moves-vol-12-by-ende-dot-app.mp3` (steady, clean), volume 0.15, deadpan-low
- Music treatment: start at 0, hold flat, fade out under the outro tagline; do not swell on the refusal (the point is it is unremarkable)
- Music cue guidance: vol-12 preset (109.96 BPM). Optional locks: the refusal landing near a strong beat (~10.93s or ~13.11s) and the outro tagline near ~17.47s. Deadpan: keep locks rare and quiet; ignore any cue that rushes a readable line.
- Audio-reactive treatment: none (deadpan stillness)
- SFX posture: very sparse — at most a soft keypress run on the typed command, one dry confirm on the success line, one dry low thud on the refusal. Nothing comedic, no error-buzz.
- Audio-coupled moments: typed command (subtle keypress), success line landing (soft drop), refusal line landing (low dry thud)
- Restraint rule: no music swell, no celebratory sting, no glitch/error sounds. Silence is allowed to sit.

## Storyboard

### Scene 1 — The problem — 3.5s
Black (#0b0e14). Line 1 fades in fast, holds: "An agent with the wrong auth writes to the wrong team." After ~1.4s, dim line 2 below: "Usually, nothing stops it." Hold.
Sequential/interaction: yes — line 2 appears after line 1 settles (single staggered reveal, each held to its reading floor).
Audio intent: near silence; music barely present; one soft drop as line 2 lands.
Audio-coupled idea: soft `interface/drop` on line 2.
Music: vol-12 entering at 0.15, almost subliminal.
Transition mood: slow crossfade → Scene 2

### Scene 2 — The pin — 4.5s
Styled terminal card slides up (rounded window, three dots, dark). Prompt `❯ cat .linctl.toml` types; a compact pin shows `[target]` with `team_key = "LIT"` highlighted (org/ids dimmed). Then `❯ linctl target` types; output line `... team LIT/... confirmed true` with `confirmed true` in green. Caption beside, small/cyan: "reads are free. writes are pinned." A dim parenthetical, smaller: "(no 13k-token MCP schema to load first.)"
Sequential/interaction: yes — two commands type in sequence; hold each output to its floor. Keypress ticks subtle.
Audio intent: calm, procedural.
Audio-coupled idea: subtle `keyboard/keypress-*` on typing; soft `interface/drop` when `confirmed true` lands.
Transition mood: clean → Scene 3

### Scene 3 — Same command, twice — 5.5s
Same terminal card. `❯ linctl issue create --title "Add rate-limit backoff to sync"` types and runs → green `LIT-15  Add rate-limit backoff to sync  [Backlog]`. Hold ~1s. Then the SAME line re-types with ` --team STG` appended (the only change), runs → red `{"error_code":"TARGET_MISMATCH","message":"target mismatch: ... team_key=STG"}`, then dim `exit 1`. A calm caption fades in: "the write does not happen."
Sequential/interaction: yes — success first, then the wrong-team retry; the `--team STG` suffix types on character by character so the viewer sees the only difference.
Audio intent: the success is soft and positive; the refusal gets a single dry low thud, no buzz — matter-of-fact.
Audio-coupled idea: `keyboard/keypress-*` on the `--team STG` suffix; soft confirm (`impact/impactSoft_medium`) on LIT-15; one dry `impact/impactSoft_heavy` on the TARGET_MISMATCH line.
Music: hold flat; do NOT swell on the refusal.
Transition mood: long hold, then slow crossfade → Scene 4

### Scene 4 — Outro — 3.5s
Card recedes. Centered tagline holds alone: "reads everywhere, writes fail closed." After ~1.2s, below it the install line in mono: `go install github.com/KyaniteHQ/linctl/cmd/linctl@latest`. Small dim footer: `github.com/KyaniteHQ/linctl`. Music fades to 0 under the tagline.
Sequential/interaction: yes — tagline, then install line.
Audio intent: settle to silence; let the tagline sit.
Audio-coupled idea: one soft `interface/drop` as the install line lands; nothing after.
Transition mood: fade to hold (end)

**Music mood for this video:** deadpan
**Audio summary:** A near-subliminal steady bed under dry, sparse keypress and drop accents; the refusal is deliberately unscored; everything fades to silence on the closing tagline.
