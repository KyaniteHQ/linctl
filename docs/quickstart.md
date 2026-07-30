# Quickstart

From nothing to a write that lands, in about five minutes.

You will need a Linear account and permission to create an OAuth application in it.

## 1. Install

```bash
# macOS
brew install --cask KyaniteHQ/linctl/linctl

# anywhere with Go
go install github.com/KyaniteHQ/linctl/cmd/linctl@latest
```

Check it works. This needs no auth and no config:

```bash
linctl --version
linctl usage
```

`linctl usage` prints a compact orientation summary. It is the thing to hand an agent that has
never seen linctl before.

## 2. Create a Linear OAuth application

In Linear: **Settings → API → OAuth applications → Create new**.

- **Redirect URI**: `http://127.0.0.1:8765/callback`
- Copy the **client ID**. Copy the **client secret** too if Linear gives you one.

linctl authenticates with OAuth, not personal API keys. That is deliberate: OAuth tokens carry
scopes and an actor identity that linctl can verify before it lets a write through.

## 3. Configure auth

Auth material lives outside your repo, in your OS config directory. It never goes in
`.linctl.toml` and never goes in git.

```bash
export LINCTL_OAUTH_CLIENT_ID=<your-client-id>

linctl auth configure \
  --client-id "$LINCTL_OAUTH_CLIENT_ID" \
  --redirect-uri "http://127.0.0.1:8765/callback" \
  --scopes read,write,issues:create,comments:create
```

If your app has a secret, add `--client-secret "$LINCTL_OAUTH_CLIENT_SECRET"`.

Now log in. This opens a browser:

```bash
linctl auth login
```

Confirm it worked:

```bash
linctl auth status
```

You should see the actor, the granted scopes, and an expiry. linctl checks this against Linear
live, so a green status means the credential really works, not just that a file exists on disk.

Secrets are never printed. Auth output reports them as `set` or `missing`.

## 4. Pin the target

From the repo where your agent will run, scaffold a target-only pin from the active
credential. Auth material never goes in this file.

```bash
# one visible team: writes .linctl.toml in the current directory
linctl init

# multiple teams: pick one (discover with linctl team list --json)
linctl init --team LIT
# optional: also pin a project after verifying it belongs to that team
linctl init --team LIT --project PROJECT_ID
```

`init` refuses to overwrite an existing `.linctl.toml` (edit or remove it; there is no
`--force`). The file is the whole safety story: *writes from this repo go here and
nowhere else.* It is safe to commit. It holds no secrets.

Confirm the credential actually resolves to it:

```bash
linctl target --json
```

<details>
<summary>Manual pin (when you already know the ids)</summary>

```bash
linctl user me --json          # organization
linctl team list --json        # teams the credential can see
```

```toml
[target]
org_id   = "your-linear-org-id"
team_key = "LIT"
team_id  = "your-linear-team-id"

# Optional. Narrows writes further, to a single project.
# project_id = "your-linear-project-id"
```

</details>

```json
{
  "viewer": { "id": "...", "name": "You", "email": "you@example.com" },
  "org": { "id": "your-linear-org-id", "name": "Your Org", "url_key": "your-org" },
  "team": { "id": "your-linear-team-id", "key": "LIT", "name": "Litmus" },
  "expected": { "org_id": "...", "team_key": "LIT", "team_id": "..." },
  "resolved": { "org_id": "...", "team_key": "LIT", "team_id": "..." },
  "confirmed": true
}
```

`expected` is what you pinned. `resolved` is what the live credential proves. `"confirmed":
true` means they agree, so writes will work. If it says otherwise, `linctl doctor` will tell
you which of the two is wrong.

Without `--json` you get the same thing on one line:

```
org your-linear-org-id team LIT/your-linear-team-id project  confirmed true
```

## 6. Read something

Reads need no pin and no confirmation. They work as soon as you are authenticated.

```bash
linctl issue list --mine
linctl issue list --state started --limit 5
linctl issue get LIT-123 --json
```

## 7. Write something

Now the part that matters.

```bash
linctl issue create --title "Hello from linctl"
```

```
LIT-42  Hello from linctl  [Backlog]
```

That issue is real. It landed in the pinned team, because the live credential resolved to the
pinned team.

## 8. Watch it refuse

Run the same write, but aim it at a team your credential cannot reach. Overriding the target on
the command line takes **both** flags, a key and an id, because a key alone is not enough to
identify a team:

```bash
linctl issue create --title "Hello from linctl" \
  --team STG --team-id 00000000-0000-0000-0000-000000000000
```

```json
{"error_code":"TARGET_MISMATCH","message":"target mismatch: expected team_id=00000000-0000-0000-0000-000000000000 team_key=STG"}
```

Exit code 1. Nothing was created. The credential cannot reach that team, so linctl stopped
before sending any mutation.

That is the guard. `--team` and `--team-id` **set** the pinned target. They did not relax the
check. The write still had to prove itself against the live credential, and it could not. No
flag makes this succeed.

This is what protects you when an agent runs with a stale token, in the wrong checkout, or from
a config it misread.

> Passing `--team` on its own gives you `TARGET_NOT_CONFIGURED` instead, because a target
> override with no team id is incomplete. Pass both, or neither.

## 9. Point your agent at it

```bash
linctl issue list --json --compact --fields identifier,title,state
id=$(linctl --id-only issue create --title "Investigate export timeouts")
linctl issue start "$id"
linctl issue comment "$id" --body "Reproduced on staging."
```

Give your agent [`skills/linctl/SKILL.md`](../skills/linctl/SKILL.md). It teaches the command
surface, the JSON contracts, and the guard, and it includes an `AGENTS.md` snippet you can drop
into the repo.

## Running headless (CI, containers, agents with no browser)

`linctl auth login` opens a browser, which is useless in CI. Two options that do not:

**Authorize as the app itself.** This needs a client secret and gives you an app actor rather
than a user actor, so changes are attributed to the application:

```bash
linctl auth configure --client-id "$LINCTL_OAUTH_CLIENT_ID" \
                      --client-secret "$LINCTL_OAUTH_CLIENT_SECRET" \
                      --scopes read,write,issues:create,comments:create
linctl auth app
```

**Or hand it a token directly.** linctl reads these from the environment and never writes them
to disk, which makes them the right fit for CI secrets:

| Variable | Purpose |
| --- | --- |
| `LINCTL_OAUTH_ACCESS_TOKEN` | use this token, skip local auth state entirely |
| `LINCTL_OAUTH_CLIENT_ID` | OAuth app client id |
| `LINCTL_OAUTH_CLIENT_SECRET` | OAuth app client secret |
| `LINCTL_OAUTH_REDIRECT_URI` | redirect URI |
| `LINCTL_OAUTH_SCOPES` | requested scopes |

These are non-persistent overrides. They do not change your saved auth state, so a CI run cannot
quietly rewrite what your laptop is logged in as.

The guard behaves identically either way. A CI job with a token for the wrong organization gets
`TARGET_MISMATCH` and a non-zero exit, which is exactly what you want a pipeline to do.

## Where to go next

- [Why not the Linear MCP server?](why-not-mcp.md) — the token cost and the write boundary
- [All documentation](README.md)
- `linctl <group> --help` — every subcommand of every group
- `linctl doctor` — run this first whenever something is not behaving

## Troubleshooting

**`TARGET_NOT_CONFIGURED`** — no `.linctl.toml`, or it is missing `org_id`, `team_key`, or
`team_id`. Reads still work. Writes will not, until you pin a target.

**`TARGET_MISMATCH` when you did not expect it** — the credential does not resolve to what you
pinned. Run `linctl target --json` and `linctl auth status`. Usually the token belongs to a
different organization than the one in the file, or the token is stale. This is the guard doing
its job, so read the message before assuming it is a bug.

**Auth expired** — `linctl auth refresh`. linctl also refreshes on its own mid-command and
retries once, so you should rarely need this.

**Anything else** — `linctl doctor` checks config, auth, and target together and tells you which
layer is broken.
