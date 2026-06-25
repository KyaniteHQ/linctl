# linctl — guarded-writes terminal demo (the README hero)
# Render with: LINEAR_API_KEY=<demo-workspace-token> ./render.sh
# Reads are free; the guarded write is refused when the token does not
# resolve to the pinned target. See README.md in this directory.

Output ../docs/assets/demo.gif
Output ../docs/assets/demo.mp4

Require linctl

Set Shell bash
Set FontSize 21
Set Width 1240
Set Height 760
Set Padding 26
Set TypingSpeed 30ms
Set Theme { "name": "linctl", "background": "#0b0e14", "foreground": "#bfbdb6", "black": "#0b0e14", "red": "#ea6c73", "green": "#7fd962", "yellow": "#f9af4f", "blue": "#53bdfa", "magenta": "#cda1fa", "cyan": "#90e1c6", "white": "#c7c7c7", "brightBlack": "#686868", "brightRed": "#f07178", "brightGreen": "#aad94c", "brightYellow": "#ffb454", "brightBlue": "#59c2ff", "brightMagenta": "#d2a6ff", "brightWhite": "#ffffff", "cursor": "#e6b450", "selection": "#1c2733" }

# --- clean prompt, no recording ---
Hide
Type "PS1='❯ ' && clear" Enter
Sleep 400ms
Show

# --- reads are free: no pinned target required ---
Type "# reads are free" Sleep 150ms Enter
Sleep 250ms
Type "linctl issue list --state started" Enter
Sleep 1500ms

# --- the pin, then re-resolve the live token against it ---
Type "cat .linctl.toml" Enter
Sleep 1200ms
Type "linctl target" Enter
Sleep 1600ms

# --- a guarded write into the pin: succeeds ---
Type 'linctl issue create --title "Add rate-limit backoff to sync"' Enter
Sleep 1700ms

# --- the money shot: aim the same write at a team this token can't reach ---
Type "# same write, wrong team:" Sleep 150ms Enter
Sleep 250ms
Type 'linctl issue create --title "Add rate-limit backoff to sync" --team STG' Enter
Sleep 2000ms
Type "# refused — token re-resolved, no bypass flag" Sleep 150ms Enter
Sleep 2400ms
