# Target-pinned Linear writes

linctl allows broad Linear reads but requires writes to compare the active credential's resolved organization, team, and optional project against the repo's pinned target. This trades a little command friction for agent safety: a stale token, wrong team, or wrong project fails closed instead of mutating unrelated Linear data.
