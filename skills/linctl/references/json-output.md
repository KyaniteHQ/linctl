# linctl `--json` output shapes

Pass `--json` to get one 2-space-indented, newline-terminated JSON value.
Most commands emit an object; a few commands with value-list semantics, such as
`issue priority-values`, emit an array. Add `--compact` for a single-line JSON value.
Add `--fields key,nested.key` to project JSON output
to just the requested keys.

For list-page commands, projection applies to each item in the command's collection array rather than
to the scalar pagination fields. The current projection collection keys come from
`internal/cli.CollectionKeys()`:

`issues`, `associations`, `cycles`, `projects`, `members`, `comments`, `updates`, `milestones`,
`documents`, `labels`, `teams`, `users`, `memberships`, `drafts`, `initiatives`, `notifications`,
`notification_subscriptions`, `release_pipelines`, `release_stages`, `releases`, `history`,
`links`, `release_notes`, `customers`, `customer_needs`, `customer_statuses`, `customer_tiers`,
`relations`, `roadmaps`, `time_schedules`, `triage_responsibilities`, `sla_configurations`,
`results`, `templates`, `workflow_states`, `agent_activities`, `agent_skills`, `agent_sessions`,
`external_users`, `audit_entry_types`, `favorites`, `emojis`, `attachments`, `custom_views`,
`project_labels`, `project_statuses`, `spans`, and `git_automation_states`.

These are the exact keys (from `internal/client/*.go`). Fields marked *optional* are omitted
when empty.

## Issue

`issue get` · `issue create` · `issue update` · `issue start` · `issue close` · `current` · `done` → **IssueSummary**

| key | type | notes |
| --- | --- | --- |
| `id` | string | Linear UUID |
| `identifier` | string | human key, e.g. `LIT-123` |
| `title` | string | |
| `branch_name` | string | Linear's suggested git branch name |
| `url` | string | |
| `priority` | number | 0–4 |
| `priority_label` | string | e.g. `Medium` |
| `team_id` | string | |
| `team` | string | team key |
| `state_id` | string | |
| `state` | string | workflow state name |
| `state_type` | string | e.g. `started`, `completed` |
| `assignee` | string | *optional* — display name |
| `project_id` | string | *optional* |
| `project` | string | *optional* — project name |

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
| `lead` | string | *optional* — display name |
| `teams` | array | `[{ "id", "key", "name" }]` |

`project list` → **ProjectList**:
`{ "projects": [ProjectSummary], "has_next_page": bool, "end_cursor": string|absent }`

`project members` → **ProjectMemberList**:
`{ "project_id", "project_name", "members": [{ "id", "name", "display_name", "email" }], "has_next_page": bool, "end_cursor": string|absent }`

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

Two things to know when parsing `target --json`:

- `project` is omitted when no `project_id` is pinned.
- `expected` and `resolved` use Go-default capitalized keys (`OrgID`, `TeamKey`, `TeamID`,
  `ProjectID`), not the snake_case used elsewhere — they mirror the config struct. Compare them
  field by field to explain a target mismatch.

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

Auth readiness succeeds only after linctl proves the token actor, token scopes, and pinned target.
App config and token material are reported as `set` or `missing`; secret values are never printed.

## Usage

`usage` · `issue usage` · `project usage` → `{ "topic": string, "text": string }`

## Error envelope

On any failure linctl writes one JSON line to **stderr** (in addition to the human-readable
error), so an agent can branch on a stable code instead of parsing prose:

```json
{ "error_code": "TARGET_MISMATCH", "message": "target mismatch: expected team_id=... resolved ..." }
```

`error_code` is one of:

- `TARGET_MISMATCH` — resolved target does not match the pinned target (hard stop; do not retry blindly).
- `TARGET_NOT_CONFIGURED` — no pinned target found; set org_id, team_key, and team_id in .linctl.toml.
- `RATE_LIMITED` — Linear returned a rate-limit response; back off and retry.
- `MUTATION_FAILED` — the mutation ran but Linear reported no success/entity.
- `INVALID_WRITE` — the write request was rejected before any API call (missing/!valid input).
- `GRAPHQL_ERROR` — the GraphQL request itself failed.
- `NOT_FOUND` — the referenced entity was not found.
- `AUTH_NOT_CONFIGURED` — required OAuth app or token state is missing.
- `AUTH_TOKEN_EXPIRED` — the saved OAuth access token is expired and cannot be used directly.
- `AUTH_REFRESH_FAILED` — Linear rejected or failed an OAuth refresh request.
- `AUTH_REAUTH_REQUIRED` — the saved OAuth state cannot be refreshed without a new login.
- `MISSING_SCOPE` — OAuth token state does not include every configured required scope.
- `AUTH_ACTOR_MISMATCH` — OAuth readiness could not prove the expected actor.
- `AUTH_TARGET_MISMATCH` — OAuth readiness could not prove the token can access the pinned target.
- `INTERNAL` — any other error (config, unknown command, decode, etc.).

Read the JSON line from stderr; the human-readable line follows it on stderr too.
