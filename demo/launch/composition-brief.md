# Hyperframes Composition Brief: linctl

## Objective
A ~17s deadpan launch/social video for linctl. Companion to the authentic terminal GIF
(`docs/assets/demo.gif`) — this one is branded and animated, not a raw recording.

## Output
- Composition directory: `demo/launch/composition/`
- Rendered video: copied to `docs/assets/launch.mp4`
- Format: landscape — 1920x1080
- Duration: 17s

## Source material
- Product: linctl (Go CLI for Linear; github.com/KyaniteHQ/linctl)
- Strongest claim: writes fail closed on org/team mismatch, no bypass flag; reads are free.
- Key visual moment to recreate: a styled terminal card running the SAME `issue create`
  command twice — once it lands (green `LIT-15 ... [Backlog]`), once it is refused
  (`{"error_code":"TARGET_MISMATCH"}` in red, exit 1).
- Copy that must appear verbatim:
  - `{"error_code":"TARGET_MISMATCH"}`
  - `go install github.com/KyaniteHQ/linctl/cmd/linctl@latest`
  - reads everywhere, writes fail closed.

## Creative direction
- Tone preset: deadpan
- Interpretation: long holds, big empty space, dry sparse sound; the refusal is NOT
  scored as a sting — it is treated as routine. Restraint everywhere.
- Angle + storyboard: see `brag-plan.md` (the creative contract). 4 scenes:
  1) the problem (3.5s), 2) the pin + `target` confirmed (4.5s), 3) same command twice:
  success then wrong-team refusal (5.5s), 4) tagline + install (3.5s).
- Hook: "An agent with the wrong auth writes to the wrong team." / "Usually, nothing stops it."
- Outro: "reads everywhere, writes fail closed." + the go install line.
- Avoid: generic SaaS language, abstract filler, error-buzz/comedic sounds, music swell on the refusal.

## Visual identity
- Background: #0b0e14
- Text: #bfbdb6
- Success green: #7fd962 ; refusal red: #ea6c73 ; amber cursor/accent: #e6b450 ; cyan labels: #90e1c6
- Terminal card surface: ~#11161f with a subtle #2a3342 border, rounded ~12px, three window dots.
- Display + body font: a monospace (JetBrains Mono if built-in; else a mono fallback). Everything is mono.

## Storyboard
Use `demo/launch/brag-plan.md` as the creative contract. Scene summary:
1. The problem — 3.5s — two flat lines on black.
2. The pin — 4.5s — terminal card: `cat .linctl.toml` (team_key "LIT" highlighted) then `linctl target` → `confirmed true` (green). Caption: "reads are free. writes are pinned." + dim aside "(no 13k-token MCP schema to load first.)"
3. Same command, twice — 5.5s — `issue create` → green `LIT-15 ... [Backlog]`; then the same line with ` --team STG` appended → red `{"error_code":"TARGET_MISMATCH", ...}` + dim `exit 1`; caption "the write does not happen."
4. Outro — 3.5s — "reads everywhere, writes fail closed." then the go install line, dim footer `github.com/KyaniteHQ/linctl`.

## Audio
- Music: `assets/music/happy-beats-business-moves-vol-12-by-ende-dot-app.mp3`, volume 0.15, flat, fade out under the outro tagline. Never swells on the refusal.
- Music cue source: vol-12 preset (`assets/music/cues/...vol-12...music-cues.json`), 109.96 BPM. Optional locks: refusal landing near a strong beat (~10.93s or ~13.11s), outro tagline near ~17.47s. Keep locks rare/quiet (deadpan).
- SFX (very sparse, CC0 Kenney): subtle `keyboard/keypress-*` on typed commands; one soft `impact/impactSoft_medium_*` when `LIT-15` lands; one dry `impact/impactSoft_heavy_*` on the TARGET_MISMATCH line; one soft `interface/drop_*` on the outro install line. No error-buzz, no glitch.
- Audio-reactive: none (deadpan stillness).
- Copy chosen music + only the selected SFX into `composition/assets/`.

## Hyperframes instructions
Use the current hyperframes skill + CLI workflow. Standalone `index.html` (no `<template>`),
`data-composition-id` on the root div, one GSAP timeline registered on `window.__timelines`.
Transitions between all scenes (slow crossfades for deadpan); entrance-only animations except the
final scene. Keep every line legible at 1080p and held to its reading floor. Run lint + validate +
inspect before render. Render landscape 1920x1080, then a 1200x630 social-card still is extracted
from the refusal frame in delivery.
