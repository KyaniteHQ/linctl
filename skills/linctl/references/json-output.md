# linctl `--json` output shapes

Pass `--json` to any command to get one 2-space-indented, newline-terminated JSON object.
Add `--compact` for a single-line JSON object. Add `--fields key,nested.key` to project JSON output
to just the requested keys; for list commands, projection applies to each item in `issues`, `projects`,
or `members`.
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

## Usage

`usage` · `issue usage` · `project usage` → `{ "topic": string, "text": string }`
