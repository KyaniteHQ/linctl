# Scope escalation is explicit

linctl will not automatically start browser authorization from ordinary commands when the current OAuth authorization lacks a required scope. Commands fail with a structured missing-scope error and show the exact reauthorization command, making permission expansion deliberate and automation-friendly.

**Status**: accepted
