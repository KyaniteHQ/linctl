# Quickstart

This page takes you from nothing to a write that lands, in about five minutes.

You need a Linear account, and the permission to create an OAuth application in it.

## 1. Install linctl

```bash
# macOS
brew install --cask KyaniteHQ/linctl/linctl

# anywhere with Go
go install github.com/KyaniteHQ/linctl/cmd/linctl@latest
```

Check that linctl runs. These two commands need no auth and no config.

```bash
linctl --version
linctl usage
```

`linctl usage` prints a compact orientation summary. Give it to an agent that has not used
linctl before.

## 2. Create a Linear OAuth application

In Linear, open **Settings → API → OAuth applications → Create new**.

- Set the **redirect URI** to `http://127.0.0.1:8765/callback`.
- Copy the **client ID**. Also copy the **client secret** if Linear shows one.

linctl authenticates with OAuth, not with a personal API key. This is deliberate. An OAuth token
carries scopes and an actor identity that linctl can verify before it permits a write.

## 3. Configure the auth

The auth material lives outside your repo, in the config directory of your operating system. It
never goes into `.linctl.toml`, and it never goes into git.

```bash
export LINCTL_OAUTH_CLIENT_ID=<your-client-id>

linctl auth configure \
  --client-id "$LINCTL_OAUTH_CLIENT_ID" \
  --redirect-uri "http://127.0.0.1:8765/callback" \
  --scopes read,write,issues:create,comments:create
```

If your application has a secret, add `--client-secret "$LINCTL_OAUTH_CLIENT_SECRET"`.

Now log in. This command opens a browser.

```bash
linctl auth login
```

Confirm that the login worked.

```bash
linctl auth status
```

The output shows the actor, the granted scopes, and an expiry. linctl checks this against Linear
live, so a good status means that the credential really works. It does not only mean that a file
exists on disk.

linctl never prints a secret. The auth output reports each secret as `set` or `missing`.

## 4. Pin the target

Run this in the repo where your agent works. `linctl init` writes a target-only pin from the
active credential. Auth material never goes into this file.

```bash
# one visible team: this writes .linctl.toml in the current directory
linctl init

# more than one team: select one. Use linctl team list --json to find the key.
linctl init --team LIT

# optional: also pin a project, after you confirm that it belongs to that team
linctl init --team LIT --project PROJECT_ID
```

`init` refuses to overwrite an existing `.linctl.toml`, and linctl has no `--force` flag. Edit
the file or remove it instead. The file carries the whole safety story: *a write from this repo
goes here and nowhere else.* The file holds no secret, so you can commit it.

<details>
<summary>Manual pin, when you already know the ids</summary>

```bash
linctl user me --json          # the organization
linctl team list --json        # the teams that the credential can see
```

```toml
[target]
org_id   = "your-linear-org-id"
team_key = "LIT"
team_id  = "your-linear-team-id"

# Optional. This makes a write narrower, down to a single project.
# project_id = "your-linear-project-id"
```

</details>

## 5. Confirm the target

Confirm that the credential really resolves to the pinned target.

```bash
linctl target --json
```

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
true` means that the two agree, so a write works. If the output shows something else, run
`linctl doctor`. It tells you which of the two is wrong.

Without `--json` you get the same facts on one line.

```
org your-linear-org-id team LIT/your-linear-team-id project  confirmed true
```

## 6. Read something

A read needs no pin and no confirmation. A read works as soon as you authenticate.

```bash
linctl issue list --mine
linctl issue list --state started --limit 5
linctl issue get LIT-123 --json
```

## 7. Write something

Now the important part.

```bash
linctl issue create --title "Hello from linctl"
```

```
LIT-42  Hello from linctl  [Backlog]
```

That issue is real. It landed in the pinned team, because the live credential resolved to the
pinned team.

## 8. See linctl refuse a write

Run the same write, but aim it at a team that your credential cannot reach. An override on the
command line needs **both** flags, a key and an id, because a key alone does not identify a team.

```bash
linctl issue create --title "Hello from linctl" \
  --team STG --team-id 00000000-0000-0000-0000-000000000000
```

```json
{"error_code":"TARGET_MISMATCH","message":"target mismatch: expected team_id=00000000-0000-0000-0000-000000000000 team_key=STG"}
```

The exit code is 1. linctl created nothing. The credential cannot reach that team, so linctl
stopped before it sent a mutation.

That is the guard. `--team` and `--team-id` **set** the pinned target. They did not relax the
check. The write still had to prove itself against the live credential, and it could not. No
flag makes this write succeed.

This is what protects you when an agent runs with an old token, in the wrong checkout, or from a
config that it misread.

> `--team` on its own gives you `TARGET_NOT_CONFIGURED` instead, because a target override
> without a team id is incomplete. Pass both flags, or pass neither.

## 9. Give linctl to your agent

```bash
linctl issue list --json --compact --fields identifier,title,state
id=$(linctl --id-only issue create --title "Investigate export timeouts")
linctl issue start "$id"
linctl issue comment "$id" --body "Reproduced on staging."
```

Give your agent [`skills/linctl/SKILL.md`](../skills/linctl/SKILL.md). It teaches the command
surface, the JSON contracts, and the guard. It also has an `AGENTS.md` block that you can copy
into the repo.

## Run linctl without a browser (CI, containers, agents)

`linctl auth login` opens a browser, which does not work in CI. Two other paths do work.

**Authorize as the application.** This needs a client secret. It gives you an app actor instead
of a user actor, so Linear attributes each change to the application.

```bash
linctl auth configure --client-id "$LINCTL_OAUTH_CLIENT_ID" \
                      --client-secret "$LINCTL_OAUTH_CLIENT_SECRET" \
                      --scopes read,write,issues:create,comments:create
linctl auth app
```

**Or give linctl a token directly.** linctl reads these variables from the environment, and
linctl never writes them to disk. That makes them the correct fit for a CI secret.

| Variable | Purpose |
| --- | --- |
| `LINCTL_OAUTH_ACCESS_TOKEN` | use this token, and skip the local auth state |
| `LINCTL_OAUTH_CLIENT_ID` | the client id of the OAuth application |
| `LINCTL_OAUTH_CLIENT_SECRET` | the client secret of the OAuth application |
| `LINCTL_OAUTH_REDIRECT_URI` | the redirect URI |
| `LINCTL_OAUTH_SCOPES` | the requested scopes |

These overrides are not persistent. They do not change your saved auth state, so a CI run cannot
quietly rewrite the login of your laptop.

The guard behaves the same way on both paths. A CI job with a token for the wrong organization
gets `TARGET_MISMATCH` and a non-zero exit, which is what a pipeline must do.

## Where to go next

- [linctl compared to the Linear MCP server](why-not-mcp.md): the token cost and the write boundary
- [All documentation](README.md)
- `linctl <group> --help`: every subcommand of every group
- `linctl doctor`: run this first when something behaves in an unexpected way

## Troubleshooting

**`TARGET_NOT_CONFIGURED`**: there is no `.linctl.toml`, or the file has no `org_id`, no
`team_key`, or no `team_id`. A read still works. A write does not work until you pin a target.

**`TARGET_MISMATCH` when you did not expect it**: the credential does not resolve to what you
pinned. Run `linctl target --json` and `linctl auth status`. Usually the token belongs to a
different organization than the file names, or the token is old. This is the guard at work, so
read the message before you assume a defect.

**The task runs from another checkout**: pass `--config /absolute/repo/.linctl.toml`.
Run `linctl --config /absolute/repo/.linctl.toml doctor --json` before a write.

**The auth expired**: run `linctl auth refresh`. linctl also refreshes the token inside a
command and retries one time, so you rarely need this.

**Anything else**: `linctl doctor` checks the config, the auth, and the target together, and it
tells you which layer is broken.
