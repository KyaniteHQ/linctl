---
name: linctl
description: Use linctl as the Linear control surface for Application, AgentActivity, AgentSkill, AuditEntry, organization, rate-limit, Notification, ReleasePipeline, ReleaseStage, Release, ReleaseNote, EntityExternalLink, SemanticSearch, issue, IssueRelation, comment, project, ProjectUpdate, ProjectStatus, ProjectLabel, ProjectRelation, Cycle, ProjectMilestone, document, label, team, TeamMembership, user, Draft, WorkflowState, TimeSchedule, TriageResponsibility, SLA configuration, Template, initiative, InitiativeRelation, InitiativeToProject, InitiativeUpdate, Roadmap, RoadmapToProject, CustomView, Customer, CustomerNeed, CustomerStatus, CustomerTier, Favorite, Emoji, and Attachment work: reads, guarded writes, branch lookup, next-work preview, doctor checks, and live smoke. Prefer it over Linear MCP, ad hoc API calls, or hand-written GraphQL when linctl covers the operation.
---

# linctl

`linctl` is the Linear control surface. Reads are cheap. Writes are guarded: the CLI resolves the active token, compares it to the pinned target, and fails closed on mismatch.

## Resolve

Pick one command prefix and use it for every Linear operation:

1. `command -v linctl` -> use `linctl`.
2. Else, inside a checkout containing `cmd/linctl/main.go` -> use `go run ./cmd/linctl`.
3. Else stop and report that `linctl` is unavailable.

Helper:

```bash
prefix="$(bash skills/linctl/scripts/linctl-resolve.sh)" || exit 1
$prefix doctor --json
```

Completion criterion: every Linear operation in the task runs through `linctl`, never through MCP, raw GraphQL, or an ad hoc script.

## Discover

Before a write:

1. Read `.linctl.toml` if present. It overlays `~/.config/linctl/config.toml`.
2. Run `linctl doctor --json` or `linctl target --json`.
3. Run `linctl usage`, plus `linctl issue usage`, `linctl project usage`, or `linctl cycle usage` before an unfamiliar domain write.
4. Use `--json` when another tool or agent will parse output. See `references/json-output.md` for exact fields.

Completion criterion: command, target, and output format are known before mutation.

## Surface

Global flags:

```bash
--json --compact --fields identifier,title,state
--id-only --quiet --fail-on-empty
--sort title --order asc
--format minimal|compact|full
--profile NAME --org ORG_ID --team TEAM --project PROJECT_ID
--timeout 30s
```

Read commands:

```bash
linctl doctor --json
linctl target --json
linctl whoami --json
linctl application info CLIENT_ID --json
linctl agent-activity list --json --limit 20
linctl agent-activity get AGENT_ACTIVITY_ID --json
linctl agent-skill list --json --limit 20
linctl agent-skill get AGENT_SKILL_ID --json
linctl audit-entry types --json
linctl organization exists URL_KEY --json
linctl organization labels --json --limit 20
linctl organization project-labels --json --limit 20
linctl organization teams --json --limit 20
linctl organization templates --json --limit 20
linctl organization users --json --limit 20
linctl rate-limit status --json
linctl notification list --json --limit 20
linctl notification get NOTIFICATION_ID --json
linctl notification subscription list --json --limit 20
linctl notification subscription get NOTIFICATION_SUBSCRIPTION_ID --json
linctl triage-responsibility list --json --limit 20
linctl triage-responsibility get TRIAGE_RESPONSIBILITY_ID --json
linctl triage-responsibility manual-selection TRIAGE_RESPONSIBILITY_ID --json
linctl sla-configuration list TEAM_ID_OR_KEY --json
linctl semantic-search QUERY --json --limit 20
linctl search documents QUERY --json --limit 20
linctl search issues QUERY --json --limit 20
linctl search projects QUERY --json --limit 20
linctl release-pipeline list --json --limit 20
linctl release-pipeline get RELEASE_PIPELINE_ID --json
linctl release-pipeline releases RELEASE_PIPELINE_ID --json --limit 20
linctl release-pipeline stages RELEASE_PIPELINE_ID --json --limit 20
linctl release-pipeline teams RELEASE_PIPELINE_ID --json --limit 20
linctl release-stage list --json --limit 20
linctl release-stage get RELEASE_STAGE_ID --json
linctl release-stage releases RELEASE_STAGE_ID --json --limit 20
linctl release list --json --limit 20
linctl release search TERM --json --limit 20
linctl release get RELEASE_ID --json
linctl release history RELEASE_ID --json --limit 20
linctl release documents RELEASE_ID --json --limit 20
linctl release issues RELEASE_ID --json --limit 20
linctl release links RELEASE_ID --json --limit 20
linctl external-link get EXTERNAL_LINK_ID --json
linctl release-note list --json --limit 20
linctl release-note get RELEASE_NOTE_ID --json
linctl issue-to-release list --json --limit 20
linctl issue-to-release get ISSUE_TO_RELEASE_ID --json
linctl external-user list --json --limit 20
linctl external-user get EXTERNAL_USER_ID --json
linctl current --json
linctl next --dry-run
linctl issue id
linctl issue title
linctl issue url
linctl issue branch LIT-123
linctl issue get LIT-123 --json
linctl issue list --json --limit 20
linctl issue list --state started --limit 20
linctl issue list --project PROJECT_ID --limit 20
linctl issue list --mine --limit 20
linctl issue list --assignee USER_ID --limit 20
linctl issue list --label LABEL_ID --limit 20
linctl issue list --cycle CYCLE_ID --limit 20
linctl issue list --created-after 2026-06-01 --limit 20
linctl issue list --created-since 2026-06-01 --limit 20
linctl issue list --created-before 2026-06-30 --limit 20
linctl issue list --has-blockers --limit 20
linctl issue list --blocks --limit 20
linctl issue list --blocked-by LIT-123 --limit 20
linctl issue list --all-teams --limit 20
linctl issue search "needle" --limit 20
linctl issue deps LIT-123 --limit 20
linctl issue attachments LIT-123 --json --limit 20
linctl issue children LIT-123 --json --limit 20
linctl issue documents LIT-123 --json --limit 20
linctl issue former-attachments LIT-123 --json --limit 20
linctl issue history LIT-123 --json --limit 20
linctl issue inverse-relations LIT-123 --json --limit 20
linctl issue labels LIT-123 --json --limit 20
linctl issue relations LIT-123 --json --limit 20
linctl issue releases LIT-123 --json --limit 20
linctl issue comments LIT-123 --limit 20
linctl issue-relation list --json --limit 20
linctl issue-relation get ISSUE_RELATION_ID --json
linctl issue pr LIT-123
linctl comment list --json --limit 20
linctl comment get COMMENT_ID --json
linctl comment bot-actor COMMENT_ID --json
linctl comment children COMMENT_ID --json --limit 20
linctl comment created-issues COMMENT_ID --json --limit 20
linctl cycle list --json --limit 20
linctl cycle get CYCLE_ID --json
linctl cycle issues CYCLE_ID --json --limit 20
linctl cycle uncompleted-issues CYCLE_ID --json --limit 20
linctl sprint current --json
linctl sprint report CYCLE_ID --json --limit 20
linctl project list --json --limit 20
linctl project get PROJECT_ID --json
linctl project attachments PROJECT_ID --json --limit 20
linctl project documents PROJECT_ID --json --limit 20
linctl project external-links PROJECT_ID --json --limit 20
linctl project history PROJECT_ID --json --limit 20
linctl project initiative-links PROJECT_ID --json --limit 20
linctl project initiatives PROJECT_ID --json --limit 20
linctl project inverse-relations PROJECT_ID --json --limit 20
linctl project issues PROJECT_ID --json --limit 20
linctl project comments PROJECT_ID --json --limit 20
linctl project labels PROJECT_ID --json --limit 20
linctl project members PROJECT_ID --json --limit 20
linctl project needs PROJECT_ID --json --limit 20
linctl project relations PROJECT_ID --json --limit 20
linctl project teams PROJECT_ID --json --limit 20
linctl project updates PROJECT_ID --json --limit 20
linctl project-update list --json --limit 20
linctl project-update get PROJECT_UPDATE_ID --json
linctl project-update comments PROJECT_UPDATE_ID --json --limit 20
linctl project-milestone list PROJECT_ID --json --limit 20
linctl project-milestone get PROJECT_MILESTONE_ID --json
linctl project-milestone issues PROJECT_MILESTONE_ID --json --limit 20
linctl project-status list --json --limit 20
linctl project-status get PROJECT_STATUS_ID --json
linctl project-label list --json --limit 20
linctl project-label get PROJECT_LABEL_ID --json
linctl project-label children PROJECT_LABEL_ID --json --limit 20
linctl project-label projects PROJECT_LABEL_ID --json --limit 20
linctl project-relation list --json --limit 20
linctl project-relation get PROJECT_RELATION_ID --json
linctl document list --json --limit 20
linctl document get DOCUMENT_ID --json
linctl label list --json --limit 20
linctl label get LABEL_ID --json
linctl label children LABEL_ID --json --limit 20
linctl label issues LABEL_ID --json --limit 20
linctl team list --json --limit 20
linctl team get TEAM_ID --json
linctl team cycles TEAM_ID --json --limit 20
linctl team issues TEAM_ID --json --limit 20
linctl team labels TEAM_ID --json --limit 20
linctl team members TEAM_ID --json --limit 20
linctl team memberships TEAM_ID --json --limit 20
linctl team projects TEAM_ID --json --limit 20
linctl team release-pipelines TEAM_ID --json --limit 20
linctl team states TEAM_ID --json --limit 20
linctl team templates TEAM_ID --json --limit 20
linctl team-membership list --json --limit 20
linctl team-membership get TEAM_MEMBERSHIP_ID --json
linctl user list --json --limit 20
linctl user get USER_ID --json
linctl user me --json
linctl user drafts --json --limit 20
linctl user assigned-issues USER_ID --json --limit 20
linctl user created-issues USER_ID --json --limit 20
linctl user delegated-issues USER_ID --json --limit 20
linctl user team-memberships USER_ID --json --limit 20
linctl user teams USER_ID --json --limit 20
linctl user my-assigned-issues --json --limit 20
linctl user my-created-issues --json --limit 20
linctl user my-delegated-issues --json --limit 20
linctl user my-team-memberships --json --limit 20
linctl user my-teams --json --limit 20
linctl workflow-state list --json --limit 20
linctl workflow-state get WORKFLOW_STATE_ID --json
linctl time-schedule list --json --limit 20
linctl time-schedule get TIME_SCHEDULE_ID --json
linctl template list --json --limit 20
linctl template get TEMPLATE_ID --json
linctl initiative list --json --limit 20
linctl initiative get INITIATIVE_ID --json
linctl initiative history INITIATIVE_ID --json --limit 20
linctl initiative links INITIATIVE_ID --json --limit 20
linctl initiative sub-initiatives INITIATIVE_ID --json --limit 20
linctl initiative updates INITIATIVE_ID --json --limit 20
linctl initiative documents INITIATIVE_ID --json --limit 20
linctl initiative projects INITIATIVE_ID --json --limit 20
linctl initiative-relation list --json --limit 20
linctl initiative-relation get INITIATIVE_RELATION_ID --json
linctl initiative-to-project list --json --limit 20
linctl initiative-to-project get INITIATIVE_TO_PROJECT_ID --json
linctl initiative-update list --json --limit 20
linctl initiative-update get INITIATIVE_UPDATE_ID --json
linctl roadmap list --json --limit 20
linctl roadmap get ROADMAP_ID --json
linctl roadmap-to-project list --json --limit 20
linctl roadmap-to-project get ROADMAP_TO_PROJECT_ID --json
linctl custom-view list --json --limit 20
linctl custom-view subscribers CUSTOM_VIEW_ID --json
linctl custom-view get CUSTOM_VIEW_ID --json
linctl custom-view initiatives CUSTOM_VIEW_ID --json --limit 20
linctl custom-view issues CUSTOM_VIEW_ID --json --limit 20
linctl custom-view organization-preferences CUSTOM_VIEW_ID --json
linctl custom-view organization-preferences values CUSTOM_VIEW_ID --json
linctl custom-view projects CUSTOM_VIEW_ID --json --limit 20
linctl custom-view user-preferences CUSTOM_VIEW_ID --json
linctl custom-view user-preferences values CUSTOM_VIEW_ID --json
linctl custom-view preference-values CUSTOM_VIEW_ID --json
linctl customer list --json --limit 20
linctl customer get CUSTOMER_ID --json
linctl customer-need list --json --limit 20
linctl customer-need get CUSTOMER_NEED_ID --json
linctl customer-status list --json --limit 20
linctl customer-status get CUSTOMER_STATUS_ID --json
linctl customer-tier list --json --limit 20
linctl customer-tier get CUSTOMER_TIER_ID --json
linctl favorite list --json --limit 20
linctl favorite children FAVORITE_ID --json --limit 20
linctl favorite get FAVORITE_ID --json
linctl emoji list --json --limit 20
linctl emoji get EMOJI_ID --json
linctl attachment list --json --limit 20
linctl attachment url URL --json --limit 20
linctl attachment get ATTACHMENT_ID --json
```

`next --dry-run` is a ranked read-only picker: it considers unstarted issues with no active blockers, then ranks by active unblock count, priority, and age. It never creates a branch or worktree.

Guarded writes:

```bash
linctl issue create --title "..." --description "..." --json
linctl issue create --title "..." --description-file ./issue.md --json
linctl issue update LIT-123 --title "..." --json
linctl issue update LIT-123 --description-file ./issue.md --json
linctl issue update LIT-123 --append "progress note" --json
linctl issue update LIT-123 --append-file ./progress.md --json
linctl issue start LIT-123 --json
linctl issue comment LIT-123 --body "..." --json
printf 'progress note\n' | linctl issue comment LIT-123 --body -
linctl issue comment LIT-123 --body-file ./comment.md --json
linctl issue reply LIT-123 COMMENT_ID --body "..." --json
linctl issue reply LIT-123 COMMENT_ID --body-file ./reply.md --json
linctl issue close LIT-123 --json
linctl done --json
linctl project create --name "..." --description "..." --json
linctl project update PROJECT_ID --name "..." --json
linctl project archive PROJECT_ID --json
linctl cycle create --starts-at START --ends-at END --json
linctl cycle update CYCLE_ID --name "..." --json
linctl cycle archive CYCLE_ID --json
linctl project-milestone create PROJECT_ID --name "..." --json
linctl project-milestone update PROJECT_MILESTONE_ID --name "..." --json
```

Unsupported writes: Notification archive/update/read-state/snooze/subscription/preference changes; ReleasePipeline and ReleaseStage configuration writes; Release, ReleaseNote, EntityExternalLink, IssueToRelease, release sync, and release complete writes; comment resolve/unresolve/edit/delete; IssueRelation create/update/delete; ProjectUpdate create/update/archive; ProjectMilestone delete; ProjectLabel create/update/delete/retire/restore; ProjectRelation create/update/delete; Document, label, team, TeamMembership, user, WorkflowState, TimeSchedule, SLA configuration, initiative, InitiativeRelation, InitiativeToProject, InitiativeUpdate, Roadmap, RoadmapToProject, CustomView, Favorite, Emoji, and Attachment writes. Report the limit instead of bypassing `linctl`.

Completion criterion: the selected command exists above and matches the requested domain.

## Safety

- A target mismatch is a hard stop. Do not retry blindly with a different token.
- Team-scoped writes (`issue create`, `project create`, `cycle create`) compare org and team.
- Resource-scoped writes resolve the existing resource first and compare pinned `project_id` when configured.
- `--org`, `--team`, and `--project` are explicit pinned-target overrides, not bypasses.
- Never print Linear tokens. Credential precedence is `LINCTL_TOKEN` > `LINEAR_API_KEY` > config `token`.
- For tests, create `linctl-it-<runid>` resources and clean them up: close disposable issues, archive disposable projects.

Completion criterion: every write has a pinned target and cleanup path, or the agent stops before writing.

## Patterns

Branch-driven work:

```bash
linctl doctor --json
linctl current --json
linctl issue deps LIT-123 --limit 20
linctl issue attachments LIT-123 --json --limit 20
linctl issue children LIT-123 --json --limit 20
linctl issue pr
linctl done --json
```

Issue queue:

```bash
linctl doctor --json
linctl --json --compact --fields identifier,title,state issue list --limit 20
linctl next --dry-run
linctl issue search "needle" --limit 20
```

Progress note from a file:

```bash
linctl doctor --json
linctl issue update LIT-123 --append-file ./progress.md --json
linctl issue comment LIT-123 --body-file ./comment.md --json
```

Disposable project smoke:

```bash
linctl doctor --json
linctl project create --name "linctl-it-<runid>" --description "disposable smoke" --json
linctl project get <created-id> --json
linctl --project <created-id> project archive <created-id> --json
```

## Live Smoke

Use this when asked for live smoke or real Linear verification.

Inside the linctl checkout:

```bash
go run github.com/go-task/task/v3/cmd/task@latest live-smoke
```

The harness accepts `LINCTL_TEST_TOKEN` first, then `LINCTL_TOKEN`, then `LINEAR_API_KEY`. Do not print any value.

Outside the checkout, run the read-only smoke:

```bash
bash skills/linctl/scripts/linctl-smoke.sh
```

Completion criterion: live smoke passed with redacted command/status evidence, or is explicitly blocked on missing credentials or target.

## Report Shape

```markdown
Known:
- `linctl doctor --json` confirmed team `TEAM`.
- `linctl issue list --json --limit 5` returned parseable JSON.

Changed:
- Created `LIT-123` with prefix `linctl-it-<runid>`, then closed it.

Blocked:
- Live write smoke not run: disposable token missing.
```

## AGENTS.md Snippet

```markdown
## Linear via linctl

- Use `linctl` for Linear issue/project/Cycle/ProjectMilestone work instead of Linear MCP, ad hoc API calls, or hand-written GraphQL.
- Resolve the command with `command -v linctl`; inside the linctl checkout use `go run ./cmd/linctl`.
- Run `linctl doctor --json` or `linctl target --json` before writes and stop on target mismatch.
- Use `--json` for agent-readable output.
- Use `linctl current --json` when the branch carries a Linear issue key.
- Never print Linear tokens.
- Keep writes pinned to `.linctl.toml` `[target]`; do not add bypass flags.
- Name test resources `linctl-it-<runid>` and close or archive them after verification.
- For live smoke, prefer `go run github.com/go-task/task/v3/cmd/task@latest live-smoke`.
- If a Linear domain is unsupported by `linctl`, report that limit instead of calling Linear directly.
```
