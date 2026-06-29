# Browser login supports manual callback fallback

`linctl auth login` will prefer a localhost callback for browser authorization, but it will also support printing the authorization URL and accepting the returned callback or authorization code manually. This keeps browser login usable from SSH sessions, containers, agent shells, and desktops where opening a browser or binding a local callback port is unavailable.

**Status**: accepted
