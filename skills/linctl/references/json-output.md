# linctl `--json` output shapes

Pass `--json` to get one JSON value with 2-space indentation and a final newline. Most commands
emit an object. A few commands have value-list semantics, such as `issue priority-values`, and
they emit an array. Add `--compact` for a JSON value on one line. Add `--fields key,nested.key`
to keep only the requested keys.

For a list-page command, the projection applies to each item of the collection array, not to the
scalar pagination fields. Each list command declares its collection key at registration. For
example, `issue list` declares `issues`, and `release search` declares `releases`. Every other
command projects the whole object and fails on a field that is not present. This includes a
single-entity `get` command with an incidental array inside the payload, such as the `teams`
array of a project.

These are the exact keys, from `internal/client/*.go`. linctl omits a field marked *optional*
when the value is empty.

## Issue

`issue get` · `issue create` · `issue update` · `issue start` · `issue close` · `current` · `done` → **IssueSummary**

| key | type | notes |
| --- | --- | --- |
| `id` | string | Linear UUID |
| `identifier` | string | the human key, for example `LIT-123` |
| `title` | string | |
| `branch_name` | string | the git branch name that Linear suggests |
| `url` | string | |
| `priority` | number | 0–4 |
| `priority_label` | string | for example `Medium` |
| `team_id` | string | |
| `team` | string | team key |
| `state_id` | string | |
| `state` | string | workflow state name |
| `state_type` | string | for example `started` or `completed` |
| `assignee` | string | *optional*: the display name |
| `project_id` | string | *optional* |
| `project` | string | *optional*: the project name |

`issue list` → **IssueList**:
`{ "issues": [IssueSummary], "has_next_page": bool, "end_cursor": string|absent }`

`issue comment` · `issue reply` → **IssueCommentResult**:
`{ "id": string, "body": string, "url": string, "issue": IssueSummary }`

`issue comments` → **IssueCommentList**:
`{ "issue_id": string, "identifier": string, "comments": [IssueCommentSummary], "has_next_page": bool, "end_cursor": string|absent }`

**IssueCommentSummary** keys:
`id`, `body`, `url`, `created_at`, optional `parent_id`, optional `user_id`, optional `user_name`, optional `display_name`.

`issue deps` → **IssueDependencyGraph**:
`{ "id": string, "identifier": string, "parent": IssueSummary|absent, "children": [IssueSummary], "blocks": [IssueSummary], "blocked_by": [IssueSummary], "has_next_page": bool }`

`issue pr` → **PullRequestPlan**:
`{ "title": string, "body": string, "command": ["gh", "pr", "create", "--title", title, "--body", body] }`

`issue import FILE` → **issueImportResult**:
`{ "count": number, "issues": [IssueSummary], "failures": [{ "row": number, "title": string, "error": string }]|absent }`

`count` and `issues` cover only the rows that linctl created. `failures` is present only when at
least one row failed. An import with no failure omits the key, so the success shape does not
change. Each failure names the row number counted from 1, the title of that row, and the error.
A row before or after a failed row still creates normally, so linctl never discards a partial
import in silence. When any row fails, the command exits non-zero and writes the error
envelope.

## Project

`project get` · `project create` · `project update` · `project archive` → **ProjectSummary**

| key | type | notes |
| --- | --- | --- |
| `id` | string | |
| `name` | string | |
| `description` | string | |
| `slug_id` | string | |
| `url` | string | |
| `priority` | number | |
| `status` | object | `{ "id", "name", "type" }` |
| `lead` | string | *optional*: the display name |
| `teams` | array | `[{ "id", "key", "name" }]` |

`project list` → **ProjectList**:
`{ "projects": [ProjectSummary], "has_next_page": bool, "end_cursor": string|absent }`

`project members` → **ProjectMemberList**:
`{ "project_id", "project_name", "members": [{ "id", "name", "display_name", "email" }], "has_next_page": bool, "end_cursor": string|absent }`

`project export PROJECT_ID DIR` → **projectExportResult**:
`{ "path": string, "slug_id": string, "attachments": number, "truncated": bool|absent }`

`path` is the written markdown file. `truncated` is present only when more attachment pages exist beyond the export cap. The file body is the project `content` field. Strip the header before `## Content` and the `## Attachments` section when you write that body back with `project update --content-file`.

## Target

`whoami` → **TargetViewer**: `{ "id", "name", "display_name", "email" }`

`target` → **ResolvedTarget**:

```json
{
  "viewer":   { "id": "...", "name": "...", "display_name": "...", "email": "..." },
  "org":      { "id": "...", "name": "...", "url_key": "..." },
  "team":     { "id": "...", "key": "LIT", "name": "..." },
  "project":  { "id": "...", "name": "..." },
  "expected": { "OrgID": "...", "TeamKey": "LIT", "TeamID": "...", "ProjectID": "..." },
  "resolved": { "OrgID": "...", "TeamKey": "LIT", "TeamID": "...", "ProjectID": "..." },
  "confirmed": true
}
```

Know two facts before you parse `target --json`.

- linctl omits `project` when the config pins no `project_id`.
- `expected` and `resolved` use Go-default capitalized keys (`OrgID`, `TeamKey`, `TeamID`,
  `ProjectID`), not the snake_case of the other commands. These keys mirror the config struct.
  Compare the two objects field by field to explain a target mismatch.

## Auth

`auth app` · `auth login --callback ...` · `auth status` · `auth refresh` → **AuthStatus**:

```json
{
  "app": { "client_id": "set", "client_secret": "set", "redirect_uri": "...", "scopes": ["read"] },
  "token": { "status": "set", "type": "Bearer", "expires_at": "...", "scopes": ["read"] },
  "actor": "app",
  "scopes": ["read"],
  "expires_at": "...",
  "token_type": "Bearer",
  "target": {
    "status": "ready",
    "expected": { "org_id": "...", "team_key": "LIT", "team_id": "...", "project_id": "..." },
    "resolved": { "org_id": "...", "team_key": "LIT", "team_id": "...", "project_id": "..." }
  }
}
```

Auth readiness succeeds only after linctl proves the token actor, the token scopes, and the
pinned target. linctl reports the app config and the token material as `set` or `missing`, and
linctl never prints a secret value.

## Usage

`usage` · `issue usage` · `project usage` → `{ "topic": string, "text": string }`

## Error envelope

On any failure, linctl writes one JSON line to **stderr**, next to the human-readable error. An
agent can then branch on a stable code instead of a parse of the prose.

```json
{ "error_code": "TARGET_MISMATCH", "message": "target mismatch: expected team_id=... resolved ..." }
```

`error_code` is one of:

- `TARGET_MISMATCH`: the resolved target does not match the pinned target. This is a hard stop. Do not retry with different auth.
- `TARGET_NOT_CONFIGURED`: there is no pinned target. Set `org_id`, `team_key`, and `team_id` in `.linctl.toml`.
- `RATE_LIMITED`: Linear returned a rate-limit response. Wait, then retry.
- `MUTATION_FAILED`: the mutation ran, but Linear reported no success and no entity.
- `INVALID_WRITE`: linctl rejected the write request before any API call, because an input was absent or not valid.
- `GRAPHQL_ERROR`: the GraphQL request failed.
- `NOT_FOUND`: linctl did not find the referenced entity.
- `AUTH_NOT_CONFIGURED`: a required OAuth app state or token state is absent.
- `AUTH_TOKEN_EXPIRED`: the saved OAuth access token expired, so linctl cannot use it directly.
- `AUTH_REFRESH_FAILED`: an OAuth refresh request failed, or Linear rejected it.
- `AUTH_REAUTH_REQUIRED`: linctl cannot refresh the saved OAuth state without a new login.
- `MISSING_SCOPE`: the OAuth token state does not have every required scope.
- `AUTH_ACTOR_MISMATCH`: OAuth readiness could not prove the expected actor.
- `AUTH_TARGET_MISMATCH`: OAuth readiness could not prove that the token reaches the pinned target.
- `AUTH_TARGET_NOT_CONFIGURED`: there is no pinned target. Set `org_id`, `team_key`, and `team_id` in `.linctl.toml`.
- `INTERNAL`: any other error, such as a config error, an unknown command, or a decode error.

Read the JSON line from stderr. The human-readable line follows it, also on stderr.
