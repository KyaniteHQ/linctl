# linctl and the Linear MCP server

Both give an AI agent access to Linear. They are built for different setups, and the honest
comparison matters more than a win.

**Use the Linear MCP server** if your agent lives in a hosted client (Claude Desktop, an IDE
panel, a chat product) and cannot run a shell. It installs in seconds, needs no binary, and
tool-calling is what those clients are designed around. It is a good product.

**Use linctl** if your agent has a shell and works in a repo, and you care about two things:
what it costs you on every session, and what it is allowed to change.

## Standing context cost

Be precise about whose fault this is. Nothing in the MCP protocol *requires* a client to inject
every tool definition up front, and the official guidance actually points at progressive
discovery. But most agent clients today do fetch the whole tool list at session start and put it
in context, so in practice that is what you pay.

Measured on 2026-07-13 against `https://mcp.linear.app/mcp`:

| | |
| --- | --- |
| Tools exposed | **47** |
| `tools/list` JSON, compact | **13,620 tokens** |
| `tools/list` JSON, pretty-printed | **20,256 tokens** |
Counted with tiktoken `o200k_base`, which is an OpenAI tokenizer. Your model may count somewhat
differently, and whether you pay the compact or the pretty-printed figure depends on how your
client serializes the schema. The honest way to state it is the range: **roughly 14k to 20k
tokens**, every session.

Reproduce it in about a minute:

```bash
LINCTL_OAUTH_ACCESS_TOKEN=<token> bash scripts/mcp-token-measure.sh
```

Re-run it rather than trusting this page, because the number moves fast. In late June 2026 the
same server was 38 tools and 10,200 compact tokens. Two weeks later: 47 tools and 13,620. The
tool count grew about a quarter; the standing token cost grew about a third. None of that is
Linear's fault. It is what happens to any MCP server that keeps shipping features.

**linctl has no tool schema to load.** It is a binary behind a Bash tool, so nothing about
linctl enters the context until the agent runs a command and reads its output.

Be fair about the caveat: an agent still has to *know how* to drive linctl, and if you hand it
[`skills/linctl/SKILL.md`](../skills/linctl/SKILL.md) or an `AGENTS.md` snippet, that costs
tokens too. The difference is that you choose when to spend it and how much, and the cheap path
really is cheap:

```bash
linctl usage           # ~400 tokens, on demand
linctl issue --help    # the rest, one group at a time
```

So the comparison is not "17k versus 400". It is "a tool schema in context every session,
whether or not the agent touches Linear" versus "as much or as little guidance as you decide to
give it".

## The write boundary

This is the actual reason linctl exists. To be accurate: the MCP *protocol* does not forbid
this. The official Linear MCP server simply does not do it.

That server gives an agent whatever its token can reach. If the token writes to nine teams, the
agent writes to nine teams. Your controls are the token's permissions, which are usually broader
than the job in front of it, and the model's judgment, which is not a control.

linctl adds a second condition. Before any write it checks the **live** credential:

1. Does the credential's own organization match the `org_id` pinned in `.linctl.toml`?
2. Is the pinned team one the credential can actually reach?

Both must hold. On failure linctl exits non-zero and prints
`{"error_code":"TARGET_MISMATCH"}`. Not a confirmation prompt, not a warning it logs and moves
past. There is no bypass flag, and adding one was
[explicitly rejected](adr/0001-linctl-architecture-baseline.md).

The flags that look like escape hatches are not. `--org`, `--team`, `--team-id`, and `--project`
*set* the pinned target. Whatever you set still gets checked against the credential.

### Be precise about what that buys you

The guard **narrows a credential**. It cannot be stronger than that credential. Specifically:

- It catches a token from the **wrong organization**, a **stale** token that lost access, and a
  pinned team the credential **cannot reach**. These are the everyday failures.
- It does **not** catch an agent that edits `.linctl.toml` to name a **sibling team its own
  token already reaches**. If that matters to you, scope the token; a config file cannot stop a
  determined agent, and anyone claiming otherwise is selling something.
- It is a **team boundary, not a correctness check**. The wrong issue inside the right team goes
  through.
- **Reads are not guarded**, on purpose.
- **Label, project-label, and initiative-label retire and restore** can be `--org-wide`. Those
  entities have no team, so the check is organization-only. Explicit flag, still fail-closed,
  but not a team-level guarantee.
- The check and the write are **separate API calls**. If something moves in between, the guard
  is acting on what it saw a moment ago. The window is small, but it is not zero.

What you get is not an unbreakable boundary. It is that the boring, common, expensive failure,
an agent with the wrong auth writing somewhere real, becomes an exit code instead of an
incident.

## Everything else

| | Linear MCP server | linctl |
| --- | --- | --- |
| Standing context cost | ~14k to 20k tokens per session | no tool schema; you choose what guidance to give |
| Install | none | one binary |
| Needs a shell | no | yes |
| Write boundary | the token's permissions | the token's permissions, narrowed to a pinned target |
| Bypass | not applicable | none exists |
| Output | tool-call results | JSON, with `--fields`, `--id-only`, `--compact` for piping |
| Composability | inside the client | pipes, `$(...)`, exit codes, shell scripts, CI |
| Surface | 47 tools | 62 command groups |

"47 tools" and "62 command groups" are not the same unit, so do not read that last row as a
coverage score. It only tells you both surfaces are broad.

They are not exclusive, either. Plenty of setups run the MCP server for a chat client and linctl
for the coding agent in the repo. That is a reasonable place to land.
