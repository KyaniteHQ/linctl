# linctl and the Linear MCP server

Both tools give an AI agent access to Linear. They are built for different setups, and an honest
comparison is more useful than a winner.

**Use the Linear MCP server** if your agent lives in a hosted client (Claude Desktop, an IDE
panel, a chat product) and cannot run a shell. It installs in seconds, it needs no binary, and
those clients are designed around tool calls. It is a good product.

**Use linctl** if your agent has a shell and works in a repo, and you care about two things:
the cost of each session, and what the agent can change.

## The standing context cost

Be precise about the cause. Nothing in the MCP protocol *requires* a client to inject every tool
definition at the start, and the official guidance points at progressive discovery. But most
agent clients today read the whole tool list at the start of a session and put it in the context.
That is what you pay in practice.

The measurement on 2026-07-13 against `https://mcp.linear.app/mcp`:

| | |
| --- | --- |
| Tools exposed | **47** |
| `tools/list` JSON, compact | **13,620 tokens** |
| `tools/list` JSON, pretty-printed | **20,256 tokens** |

The count uses tiktoken `o200k_base`, which is an OpenAI tokenizer. Your model can count in a
different way. Whether you pay the compact figure or the pretty-printed figure depends on how
your client serializes the schema. The honest statement is the range: **about 14,000 to 20,000
tokens**, in each session.

Reproduce the measurement in about one minute.

```bash
LINCTL_OAUTH_ACCESS_TOKEN=<token> bash scripts/mcp-token-measure.sh
```

Run it again instead of trusting this page, because the number moves quickly. In late June 2026
the same server had 38 tools and 10,200 compact tokens. Two weeks later it had 47 tools and
13,620 compact tokens. The tool count grew about one quarter, and the standing token cost grew
about one third. None of this is a fault of Linear. It is what happens to any MCP server that
keeps adding features.

**linctl has no tool schema to load.** It is a binary behind a Bash tool. Nothing about linctl
enters the context until the agent runs a command and reads the output.

State the caveat as well. An agent still has to know how to drive linctl. If you give it
[`skills/linctl/SKILL.md`](../skills/linctl/SKILL.md) or an `AGENTS.md` block, that costs tokens
too. The difference is that you choose when to spend the tokens and how many, and the cheap path
really is cheap.

```bash
linctl usage           # about 400 tokens, on demand
linctl issue --help    # the rest, one group at a time
```

So the comparison is not "17,000 against 400". On one side, a tool schema sits in the context in
each session, whether or not the agent touches Linear. On the other side, you decide how much
guidance to give.

## The write boundary

This is the real reason that linctl exists. To be accurate: the MCP *protocol* does not forbid a
write boundary. The official Linear MCP server simply does not have one.

That server gives an agent everything that its token can reach. If the token writes to nine
teams, the agent writes to nine teams. Your controls are the permissions of the token, which are
usually wider than the task, and the judgement of the model, which is not a control.

linctl adds a second condition. Before each write, linctl checks the **live** credential.

1. Is the organization of the credential the `org_id` pinned in `.linctl.toml`?
2. Can the credential reach the pinned team?

Both facts must hold. On a failure linctl exits non-zero and prints
`{"error_code":"TARGET_MISMATCH"}`. This is not a confirmation prompt, and it is not a warning
that linctl logs before it continues. linctl has no bypass flag, and
[ADR 0001](adr/0001-linctl-architecture-baseline.md) rejects one.

The flags that look like escape hatches are not escape hatches. `--org`, `--team`, `--team-id`,
and `--project` *set* the pinned target. linctl still checks whatever you set against the
credential.

### Be precise about what the guard buys you

The guard **makes a credential narrow**. The guard cannot be stronger than that credential.

- The guard catches a token from the **wrong organization**, an **old** token that lost its
  access, and a pinned team that the credential **cannot reach**. These are the everyday
  failures.
- The guard does **not** catch an agent that edits `.linctl.toml` to name a **sibling team that
  its own token already reaches**. If that matters to you, make the token scope small. A config
  file cannot stop a determined agent.
- The guard is a **team boundary, not a correctness check**. The wrong issue inside the correct
  team goes through.
- **A read is not guarded.** This is deliberate.
- **Some label and team commands are organization-wide.** A label, a project label, and an
  initiative label have no team, so `--org-wide` compares the organization only. `team create`
  needs the same flag, because a new team is always outside the pinned team. The flag is
  explicit, the write still fails closed, but the check is not a team check. The
  [generated guarded-write list](../skills/linctl/references/guarded-writes.md) marks every
  command that takes the flag.
- The check and the write are **two API calls**. If the state changes between the two calls, the
  guard acted on the older state. The window is small. It is not zero.

What you get is not an unbreakable boundary. One failure is common and expensive: an agent with
the wrong auth writes into a real team. The guard turns that failure into an exit code.

## Everything else

| | Linear MCP server | linctl |
| --- | --- | --- |
| Standing context cost | about 14,000 to 20,000 tokens per session | no tool schema; you choose the guidance |
| Install | none | one binary |
| Needs a shell | no | yes |
| Write boundary | the permissions of the token | the permissions of the token, made narrow by a pinned target |
| Bypass | not applicable | none exists |
| Output | tool-call results | JSON, with `--fields`, `--id-only`, and `--compact` for pipes |
| Composability | inside the client | pipes, `$(...)`, exit codes, shell scripts, CI |
| Surface | 47 tools | more than 60 top-level commands |

"47 tools" and "more than 60 commands" are not the same unit, so do not read the last row as a
coverage score. It only tells you that both surfaces are wide.

The two tools are not exclusive. Many setups run the MCP server for a chat client and linctl for
the coding agent in the repo. That is a reasonable place to stop.
