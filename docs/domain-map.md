# linctl domain map

This map is derived from the vendored Linear schema at `internal/client/schema.graphql`.
Command names below are either implemented CLI surface or intentionally deferred surface. Implementation slices must use GraphQL operations backed by these schema fields.

## Core target

| CLI surface | Schema backing | Notes |
| --- | --- | --- |
| `whoami` | `Query.viewer`, `User` | Reads the authenticated user. |
| `target` | `Query.organization`, `Query.teams`, `Query.team`, `Query.projects`, `Query.project` | Resolves the active token's organization, team, and optional project. |
| `doctor` | `Query.viewer`, `Query.teams`, `TargetProject` (`Query.project`) when `project_id` is pinned | Read-only health check for config load, token presence, and pinned-target confirmation. Does not print token values. |
| `application info` | `Query.applicationInfo` | Read-only public OAuth application metadata by client id. |
| `organization exists` | `Query.organizationExists` | Read-only URL-key existence check for organization lookup. |
| `organization labels` | `Organization.labels` via `Query.organization` | Read-only organization-level issue labels. |
| `organization project-labels` | `Organization.projectLabels` via `Query.organization` | Read-only organization-level project labels. |
| `organization teams` | `Organization.teams` via `Query.organization` | Read-only teams visible to the authenticated user. |
| `organization templates` | `Organization.templates` via `Query.organization` | Read-only organization-level templates. |
| `organization users` | `Organization.users` via `Query.organization` | Read-only active users visible to the authenticated user. |
| `rate-limit status` | `Query.rateLimitStatus` | Read-only quota status for the authenticated Linear client. |

The target vocabulary is `org_id`, `team_key`, `team_id`, and optional `project_id`. Do not introduce `workspace` as a flag or JSON key synonym.

## AgentActivity

Schema backing:

- Types: `AgentActivity`, `AgentActivityConnection`, `AgentActivityContent`
- Reads: `Query.agentActivities`, `Query.agentActivity`
- Writes: `Mutation.agentActivityCreate`, `Mutation.agentActivityUpdate`, `Mutation.agentActivityArchive`, prompt-specific activity mutations
- Relevant fields: `AgentActivity.id`, `AgentActivity.agentSession`, `AgentActivity.content`, `AgentActivity.signal`, `AgentActivity.ephemeral`, `AgentActivity.sourceComment`, `AgentActivity.user`, `AgentActivity.createdAt`, `AgentActivity.updatedAt`

Command coverage:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `agent-activity list` | `Query.agentActivities` | Read-only |
| `agent-activity get` | `Query.agentActivity` | Read-only |
| `agent-activity create` | `Mutation.agentActivityCreate` | Blocked: create writes into an agent session and needs explicit session/comment guard semantics |
| `agent-activity update` | `Mutation.agentActivityUpdate` | Blocked: update must resolve the agent session and activity scope before mutation |
| `agent-activity archive` | `Mutation.agentActivityArchive` | Blocked: destructive command needs explicit AgentActivity safety semantics |

Only `agent-activity list` and `agent-activity get` are implemented in the current CLI. AgentActivity writes remain deferred until their session and comment guard model is explicit.

## AgentSkill

Schema backing:

- Types: `AgentSkill`, `AgentSkillConnection`
- Reads: `Query.agentSkills`, `Query.agentSkill`
- Writes: `Mutation.agentSkillCreate`, `Mutation.agentSkillUpdate`, `Mutation.agentSkillArchive`
- Relevant fields: `AgentSkill.id`, `AgentSkill.title`, `AgentSkill.body`, `AgentSkill.description`, `AgentSkill.slugId`, `AgentSkill.teamId`, `AgentSkill.shared`, `AgentSkill.recentUsageCount`, `AgentSkill.owner`, `AgentSkill.creator`, `AgentSkill.lastUpdatedBy`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `agent-skill list` | `Query.agentSkills` | Read-only |
| `agent-skill get` | `Query.agentSkill` | Read-only |
| `agent-skill create` | `Mutation.agentSkillCreate` | Blocked: create can expose reusable agent instructions and needs explicit team/owner guard semantics |
| `agent-skill update` | `Mutation.agentSkillUpdate` | Blocked: update must resolve the AgentSkill's team and ownership scope before mutation |
| `agent-skill archive` | `Mutation.agentSkillArchive` | Blocked: destructive command needs explicit AgentSkill safety semantics |

Only `agent-skill list` and `agent-skill get` are implemented in the current CLI. AgentSkill writes remain deferred until their guard model is explicit.

## ExternalUser

Schema backing:

- Types: `ExternalUser`, `ExternalUserConnection`
- Reads: `Query.externalUsers`, `Query.externalUser`
- Writes: none exposed directly; `Mutation.userExternalUserDisconnect` is tracked with the User write surface.
- Relevant fields: `ExternalUser.id`, `ExternalUser.name`, `ExternalUser.displayName`, `ExternalUser.avatarUrl`, `ExternalUser.lastSeen`, `ExternalUser.createdAt`, `ExternalUser.updatedAt`, `ExternalUser.archivedAt`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `external-user list` | `Query.externalUsers` | Read-only |
| `external-user get` | `Query.externalUser` | Read-only |

Only `external-user list` and `external-user get` are implemented in the current CLI. `ExternalUser.email` is intentionally omitted from local GraphQL selections and default output.

## AuditEntry

Schema backing:

- Types: `AuditEntryType`
- Reads: `Query.auditEntryTypes`
- Deferred reads: `Query.auditEntries`
- Relevant fields: `AuditEntryType.type`, `AuditEntryType.description`

Command coverage:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `audit-entry types` | `Query.auditEntryTypes` | Read-only |
| `audit-entry list` | `Query.auditEntries` | Blocked: audit log entries can expose actor, IP, country, and request metadata; needs an explicit admin/security output model |

Only `audit-entry types` is implemented in the current CLI. Audit entry listing remains deferred until the security-sensitive output model is explicit.

## Notification

Schema backing:

- Types: `Notification`, `NotificationConnection`, `NotificationSubscription`, `NotificationSubscriptionConnection`
- Reads: `Query.notifications`, `Query.notification`, `Query.notificationSubscriptions`, `Query.notificationSubscription`
- Writes: `Mutation.notificationArchive`, `Mutation.notificationArchiveAll`, `Mutation.notificationUpdate`, `Mutation.notificationMarkReadAll`, `Mutation.notificationMarkUnreadAll`, `Mutation.notificationSnoozeAll`, `Mutation.notificationUnsnoozeAll`, `Mutation.notificationCategoryChannelSubscriptionUpdate`, `Mutation.notificationSubscriptionCreate`, `Mutation.notificationSubscriptionUpdate`, `Mutation.notificationSubscriptionDelete`
- Relevant fields: `Notification.id`, `Notification.type`, `Notification.category`, `Notification.title`, `Notification.subtitle`, `Notification.url`, `Notification.inboxUrl`, `Notification.user`, `Notification.actor`, `NotificationSubscription.id`, `NotificationSubscription.active`, `NotificationSubscription.subscriber`, target entity fields

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `notification list` | `Query.notifications` | Read-only |
| `notification get` | `Query.notification` | Read-only |
| `notification subscription list` | `Query.notificationSubscriptions` | Read-only |
| `notification subscription get` | `Query.notificationSubscription` | Read-only |
| `notification archive` | `Mutation.notificationArchive` | Blocked: mutates the authenticated user's inbox state; needs an explicit viewer-state safety model |
| `notification archive all` | `Mutation.notificationArchiveAll` | Blocked: bulk inbox mutation needs explicit safety semantics |
| `notification update` | `Mutation.notificationUpdate` | Blocked: direct inbox-state mutation needs an explicit viewer-state safety model |
| `notification mark read all` | `Mutation.notificationMarkReadAll` | Blocked: bulk inbox mutation needs explicit safety semantics |
| `notification mark unread all` | `Mutation.notificationMarkUnreadAll` | Blocked: bulk inbox mutation needs explicit safety semantics |
| `notification snooze all` | `Mutation.notificationSnoozeAll` | Blocked: bulk inbox mutation needs explicit safety semantics |
| `notification unsnooze all` | `Mutation.notificationUnsnoozeAll` | Blocked: bulk inbox mutation needs explicit safety semantics |
| `notification category channel subscription update` | `Mutation.notificationCategoryChannelSubscriptionUpdate` | Blocked: viewer notification preference mutation needs an explicit consent model |
| `notification subscription create` | `Mutation.notificationSubscriptionCreate` | Blocked: subscription writes can target several entity types and need explicit target-resolution semantics |
| `notification subscription update` | `Mutation.notificationSubscriptionUpdate` | Blocked: update must resolve the subscription target before mutation |
| `notification subscription delete` | `Mutation.notificationSubscriptionDelete` | Blocked: destructive viewer preference command needs explicit safety semantics |

Only `notification list`, `notification get`, `notification subscription list`, and `notification subscription get` are implemented in the current CLI. Notification writes are deferred as viewer-state and preference surface.

## Release

Schema backing:

- Types: `Release`, `ReleasePipeline`, `ReleaseStage`, `ReleaseNote`, `EntityExternalLink`, `IssueToRelease`
- Reads: `Query.releasePipelines`, `Query.releasePipeline`, `ReleasePipeline.releases`, `ReleasePipeline.stages`, `ReleasePipeline.teams`, `Query.releaseStages`, `Query.releaseStage`, `ReleaseStage.releases`, `Query.releases`, `Query.release`, `Release.history`, `Release.documents`, `Release.issues`, `Release.links`, `Query.entityExternalLink`, `Query.releaseSearch`, `Query.releaseNotes`, `Query.releaseNote`, `Query.issueToReleases`, `Query.issueToRelease`
- Deferred reads: access-key release reads and release document-content reads
- Writes: `Mutation.releasePipelineCreate`, `Mutation.releasePipelineUpdate`, `Mutation.releasePipelineArchive`, `Mutation.releasePipelineDelete`, `Mutation.releaseStageCreate`, `Mutation.releaseStageUpdate`, `Mutation.releaseStageArchive`, `Mutation.releaseStageUnarchive`, plus Release/ReleaseNote/IssueToRelease create/update/archive/delete/sync/complete mutations
- Relevant fields: `Release.id`, `Release.name`, `Release.slugId`, `Release.version`, `Release.pipeline`, `Release.stage`, `Release.issueCount`, `ReleaseNote.id`, `ReleaseNote.title`, `ReleaseNote.slugId`, `ReleaseNote.pipeline`, `ReleaseNote.releaseCount`, `ReleasePipeline.id`, `ReleasePipeline.name`, `ReleasePipeline.slugId`, `ReleasePipeline.type`, `ReleasePipeline.isProduction`, `ReleasePipeline.approximateReleaseCount`, `ReleaseStage.id`, `ReleaseStage.name`, `ReleaseStage.type`, `ReleaseStage.pipeline`, `EntityExternalLink.id`, `EntityExternalLink.label`, `EntityExternalLink.url`, `EntityExternalLink.sortOrder`, `EntityExternalLink.creator`, `EntityExternalLink.initiative`, `EntityExternalLink.project`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `release-pipeline list` | `Query.releasePipelines` | Read-only |
| `release-pipeline get` | `Query.releasePipeline` | Read-only |
| `release-pipeline releases` | `ReleasePipeline.releases` via `Query.releasePipeline` | Read-only |
| `release-pipeline stages` | `ReleasePipeline.stages` via `Query.releasePipeline` | Read-only |
| `release-pipeline teams` | `ReleasePipeline.teams` via `Query.releasePipeline` | Read-only |
| `release-stage list` | `Query.releaseStages` | Read-only |
| `release-stage get` | `Query.releaseStage` | Read-only |
| `release-stage releases` | `ReleaseStage.releases` via `Query.releaseStage` | Read-only |
| `release list` | `Query.releases` | Read-only |
| `release search` | `Query.releaseSearch` | Read-only |
| `release get` | `Query.release` | Read-only |
| `release history` | `Release.history` via `Query.release` | Read-only |
| `release documents` | `Release.documents` via `Query.release` | Read-only |
| `release issues` | `Release.issues` via `Query.release` | Read-only |
| `release links` | `Release.links` via `Query.release` | Read-only |
| `external-link get` | `Query.entityExternalLink` | Read-only |
| `release-note list` | `Query.releaseNotes` | Read-only |
| `release-note get` | `Query.releaseNote` | Read-only |
| `issue-to-release list` | `Query.issueToReleases` | Read-only |
| `issue-to-release get` | `Query.issueToRelease` | Read-only |
| `release-pipeline create` | `Mutation.releasePipelineCreate` | Blocked: pipeline configuration is team/admin release surface and needs explicit guard semantics |
| `release-pipeline update` | `Mutation.releasePipelineUpdate` | Blocked: update must resolve and compare associated teams before mutation |
| `release-pipeline archive` | `Mutation.releasePipelineArchive` | Blocked: destructive command needs explicit safety semantics |
| `release-pipeline unarchive` | `Mutation.releasePipelineUnarchive` | Blocked: restore command needs explicit safety semantics |
| `release-pipeline delete` | `Mutation.releasePipelineDelete` | Blocked: destructive command needs explicit safety semantics |
| `release-stage create` | `Mutation.releaseStageCreate` | Blocked: release workflow configuration needs explicit pipeline/team guard semantics |
| `release-stage update` | `Mutation.releaseStageUpdate` | Blocked: update must resolve the stage's pipeline and teams before mutation |
| `release-stage archive` | `Mutation.releaseStageArchive` | Blocked: destructive command needs explicit safety semantics |
| `release-stage unarchive` | `Mutation.releaseStageUnarchive` | Blocked: restore command needs explicit safety semantics |
| `release create` | `Mutation.releaseCreate` | Blocked: create must resolve pipeline/team guard semantics before mutation |
| `release update` | `Mutation.releaseUpdate` | Blocked: update must resolve the release pipeline/stage and associated teams before mutation |
| `release archive` | `Mutation.releaseArchive` | Blocked: destructive command needs explicit safety semantics |
| `release unarchive` | `Mutation.releaseUnarchive` | Blocked: restore command needs explicit safety semantics |
| `release delete` | `Mutation.releaseDelete` | Blocked: destructive command needs explicit safety semantics |
| `release complete` | `Mutation.releaseComplete`, `Mutation.releaseCompleteByAccessKey` | Blocked: lifecycle transition and access-key behavior need explicit guard semantics |
| `release sync` | `Mutation.releaseSync`, `Mutation.releaseSyncByAccessKey` | Blocked: sync mutates release associations and needs explicit guard semantics |
| `release-note create` | `Mutation.releaseNoteCreate` | Blocked: create must resolve release pipeline and release range semantics before mutation |
| `release-note update` | `Mutation.releaseNoteUpdate` | Blocked: update must resolve covered releases and pipeline before mutation |
| `release-note archive` | `Mutation.releaseNoteArchive` | Blocked: destructive command needs explicit safety semantics |
| `release-note delete` | `Mutation.releaseNoteDelete` | Blocked: destructive command needs explicit safety semantics |
| `issue-to-release create` | `Mutation.issueToReleaseCreate` | Blocked: association write must compare issue and release scope before mutation |
| `issue-to-release update` | `Mutation.issueToReleaseUpdate` | Blocked: association update must compare issue and release scope before mutation |
| `issue-to-release delete` | `Mutation.issueToReleaseDelete` | Blocked: destructive association command needs explicit safety semantics |

Release, ReleasePipeline, ReleaseStage, ReleaseNote, EntityExternalLink, and IssueToRelease read commands are implemented in the current CLI. IssueToRelease writes, sync, complete, access-key, and broader association commands remain deferred until their control-surface shape and guard model are explicit.

## Issue

Schema backing:

- Types: `Issue`, `IssueConnection`
- Reads: `Query.issues`, `Query.issue`, `Issue.botActor`, `Issue.stateHistory`, `Issue.subscribers`
- Writes: `Mutation.issueCreate`, `Mutation.issueUpdate`, `Mutation.issueArchive`, `Mutation.commentCreate`
- Inputs: `IssueCreateInput`, `IssueUpdateInput`
- Relevant fields: `Issue.id`, `Issue.identifier`, `Issue.number`, `Issue.title`, `Issue.team`, `Issue.cycle`, `Issue.project`, `Issue.projectMilestone`, `Issue.assignee`, `Issue.state`, `Issue.documents`, `Issue.comments`, `Issue.url`, `Issue.branchName`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `issue list` | `Query.issues`, optionally filtered by `Issue.team.id`, `Issue.state.type` (`--state`, with `--status` as an alias; human state names are normalized to the schema state type before filtering), `Issue.project.id`, `Issue.assignee.id`, `Issue.labels.some.id`, `Issue.cycle.id`, `Issue.createdAt.gte` (`--created-after` / `--created-since`), `Issue.createdAt.lte`, `Issue.hasBlockedByRelations.eq`, or `Issue.hasBlockingRelations.eq`; `--blocked-by ISSUE` traverses `Issue.relations` with `IssueRelation.type == "blocks"` and returns matching `IssueRelation.relatedIssue`; `--all-teams` omits the team filter | Read-only |
| `issue search` | `Query.issues`, filtered by `Issue.searchableContent` | Read-only |
| `issue figma-file-key-search` | `Query.issueFigmaFileKeySearch`; returns compact issue summaries for a Figma file key | Read-only |
| `issue priority-values` | `Query.issuePriorityValues` | Read-only |
| `issue filter-suggestion` | `Query.issueFilterSuggestion`; returns the suggested filter JSON plus log id only | Read-only |
| `issue title-suggestion` | `Query.issueTitleSuggestionFromCustomerRequest`; returns the suggested title plus log id only | Read-only |
| `issue get` | `Query.issue` | Read-only |
| `issue deps` | `Query.issue`, `Issue.parent`, `Issue.children`, `Issue.relations`, `Issue.inverseRelations`; `IssueRelation.type == "blocks"` separates blocked issues from blockers | Read-only |
| `issue attachments` | `Issue.attachments` via `Query.issue` | Read-only |
| `issue bot-actor` | `Issue.botActor` via `Query.issue` | Read-only, bot metadata only |
| `issue children` | `Issue.children` via `Query.issue` | Read-only |
| `issue documents` | `Issue.documents` via `Query.issue` | Read-only |
| `issue former-attachments` | `Issue.formerAttachments` via `Query.issue` | Read-only |
| `issue former-needs` | `Issue.formerNeeds` via `Query.issue`; returns customer-need metadata without body/content | Read-only |
| `issue history` | `Issue.history` via `Query.issue`; returns compact metadata only, not raw change payloads or content fields | Read-only |
| `issue inverse-relations` | `Issue.inverseRelations` via `Query.issue` | Read-only |
| `issue labels` | `Issue.labels` via `Query.issue` | Read-only |
| `issue needs` | `Issue.needs` via `Query.issue`; returns customer-need metadata without body/content | Read-only |
| `issue relations` | `Issue.relations` via `Query.issue` | Read-only |
| `issue releases` | `Issue.releases` via `Query.issue` | Read-only |
| `issue shared-access` | `Issue.sharedAccess` via `Query.issue`; omits shared user details and exposes only flags/counts/disallowed fields | Read-only |
| `issue state-history` | `Issue.stateHistory` via `Query.issue` | Read-only, workflow-state span metadata |
| `issue subscribers` | `Issue.subscribers` via `Query.issue` | Read-only |
| `issue vcs-branch-search get` | `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search attachments` | `Issue.attachments` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search bot-actor` | `Issue.botActor` via `Query.issueVcsBranchSearch` | Read-only, bot metadata only |
| `issue vcs-branch-search children` | `Issue.children` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search documents` | `Issue.documents` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search former-attachments` | `Issue.formerAttachments` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search comments` | `Issue.comments` via `Query.issueVcsBranchSearch`; returns comment metadata without body | Read-only |
| `issue vcs-branch-search former-needs` | `Issue.formerNeeds` via `Query.issueVcsBranchSearch`; returns customer-need metadata without body/content | Read-only |
| `issue vcs-branch-search history` | `Issue.history` via `Query.issueVcsBranchSearch`; returns compact metadata only, not raw change payloads or content fields | Read-only |
| `issue vcs-branch-search inverse-relations` | `Issue.inverseRelations` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search labels` | `Issue.labels` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search needs` | `Issue.needs` via `Query.issueVcsBranchSearch`; returns customer-need metadata without body/content | Read-only |
| `issue vcs-branch-search relations` | `Issue.relations` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search releases` | `Issue.releases` via `Query.issueVcsBranchSearch` | Read-only |
| `issue vcs-branch-search shared-access` | `Issue.sharedAccess` via `Query.issueVcsBranchSearch`; omits shared user details and exposes only flags/counts/disallowed fields | Read-only |
| `issue vcs-branch-search state-history` | `Issue.stateHistory` via `Query.issueVcsBranchSearch` | Read-only, workflow-state span metadata |
| `issue vcs-branch-search subscribers` | `Issue.subscribers` via `Query.issueVcsBranchSearch` | Read-only |
| `issue id` | Current checkout issue identifier from git/jj context | Read-only |
| `issue title` | `Query.issue` after current checkout or explicit issue resolution | Read-only |
| `issue url` | `Query.issue` after current checkout or explicit issue resolution | Read-only |
| `issue open` | `Query.issue` resolves `Issue.url`, then the platform opener (`xdg-open`/`open`/`rundll32`) launches it with the URL as a discrete argv argument | Read-only |
| `issue export` | `Query.issue` (`GetIssueDetail`), `Issue.comments`, and `Issue.attachments` are assembled into a single markdown file (`<DIR>/<identifier>.md`) holding the metadata header, description, comments, and attachment URLs; capped at 250 comments/attachments with a stderr note when more pages exist | Read-only, writes only local files |
| `issue import` | Reads a CSV or JSON file (format from extension), normalizes each row's state/priority, rejects any row whose `team` key ≠ the pinned `team_key`, then creates each issue via guarded `Mutation.issueCreate` (`CreateIssue`); `--dry-run` renders the normalized rows locally and performs no mutation | Team-scoped per row; each create re-runs the write guard; `--dry-run` writes nothing |
| `issue bulk-export` | `Query.team`/`Team.issues` (`ListIssuesByTeam`) for the resolved team are written to a CSV or JSON file (format from extension), capped by `--limit` (default 250) | Read-only, writes only the local file |
| `issue branch` | `Query.issue`, `Issue.branchName` | Read-only |
| `issue pr` | `Query.issue`; emits a local `gh pr create` title/body plan without calling GitHub | Read-only |
| `next` | `Query.issues`, filtered by `Issue.team.id`, `Issue.state.type == "unstarted"`, and `Issue.hasBlockedByRelations.eq == false`; fetches `Issue.relations`, `Issue.priority`, and `Issue.createdAt`, then ranks by active unblock count, priority, and age. `--dry-run` prints the top candidate and writes nothing; without it the top candidate is started via guarded `Mutation.issueUpdate` (`StartIssue`); `--checkout` runs `git checkout -b <Issue.branchName>` before starting | `--dry-run` read-only; otherwise resource-scoped start of the picked issue |
| `done` | Current checkout issue identifier, then `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved |
| `issue create` | `Mutation.issueCreate` with `IssueCreateInput.teamId`, optional `projectId`; `--description-file` is resolved locally before mutation; `--template` reads `Template.templateData` via `Query.template` (free read) and fills title/description defaults that explicit flags override; `--section NAME=VALUE` fills a markdown section locally before mutation; `--dry-run` renders the assembled draft locally and performs no mutation; `--state` (alias `--status`) normalizes a human state name to a schema state type and resolves `IssueCreateInput.stateId` via `Query.workflowStates` filtered by team + type; `--priority` normalizes human words (`urgent`/`high`/`medium`/`low`/`none`) or `0-4` to `IssueCreateInput.priority` | Team-scoped unless `projectId` is set; `--dry-run` writes nothing |
| `issue update` | `Mutation.issueUpdate` with `IssueUpdateInput`; `--description-file` replaces description, while `--append` or `--append-file` first reads `Issue.description` and appends text before sending `description`; `--state` (alias `--status`) and `--priority` are normalized the same way as on `issue create`, with `stateId` resolved via `Query.workflowStates` filtered by the issue's team + type | Resource-scoped when a project target is involved |
| `issue start` | `Query.viewer`, `Query.workflowStates` filtered to `started`, then `Mutation.issueUpdate` with `IssueUpdateInput.assigneeId` and `stateId` | Resource-scoped when a project target is involved |
| `issue comment` | `Mutation.commentCreate`; `--body -` reads stdin and `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project |
| `issue reply` | `Mutation.commentCreate` with `CommentCreateInput.parentId`; `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project |
| `issue close` | `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved |
| `issue comments` | `Issue.comments` via `Query.issue` | Read-only |

Issue customer-need child reads use a metadata-only projection and intentionally omit `body` and `content`. Shared-access reads omit shared user details and expose only booleans, counts, and disallowed field names.

## IssueRelation

Use `IssueRelation` for Linear dependency and similarity relations between issues. It is root issue graph metadata; `issue deps` remains the focused per-issue dependency view.

Schema backing:

- Types: `IssueRelation`, `IssueRelationConnection`
- Reads: `Query.issueRelations`, `Query.issueRelation`
- Writes: `Mutation.issueRelationCreate`, `Mutation.issueRelationUpdate`, `Mutation.issueRelationDelete`
- Inputs: `IssueRelationCreateInput`, `IssueRelationUpdateInput`
- Relevant fields: `IssueRelation.id`, `IssueRelation.type`, `IssueRelation.issue`, `IssueRelation.relatedIssue`, `IssueRelation.createdAt`, `IssueRelation.updatedAt`, `IssueRelation.archivedAt`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `issue-relation list` | `Query.issueRelations` | Read-only |
| `issue-relation get` | `Query.issueRelation` | Read-only |
| `issue relate` | `Mutation.issueRelationCreate` with `IssueRelationCreateInput` | Team-scoped on both endpoints: resolve each issue and compare the pinned team before linking; `--type blocks` is refused when it would close a direct cycle |
| `issue unrelate` | `Mutation.issueRelationDelete` | Resolve the relation, then compare the pinned team for both linked issues before deleting |
| `issue-relation update` | `Mutation.issueRelationUpdate` | Blocked: update must resolve and compare both issue endpoints before mutation |

`issue-relation list` and `issue-relation get` (reads) plus `issue relate ISSUE RELATED --type` and `issue unrelate RELATION_ID` (writes) are implemented in the current CLI. `issue relate` resolves both endpoints through `requireIssue` so a relation lands only when both issues belong to the resolved team, and a `blocks` relation is rejected when the related issue already blocks the source issue. `issue unrelate` is the one approved relation delete; it resolves the relation and confirms both linked issues before removing it. `issue-relation update` stays deferred until its endpoint guard model is explicit.

## Comment

Schema backing:

- Types: `Comment`, `CommentConnection`
- Reads: `Query.comments`, `Query.comment`, `Issue.comments`, `Comment.botActor`, `Comment.children`, `Comment.createdIssues`
- Writes: `Mutation.commentCreate`, `Mutation.commentUpdate`, `Mutation.commentDelete`, `Mutation.commentResolve`, `Mutation.commentUnresolve`
- Inputs: `CommentCreateInput`, `CommentUpdateInput`
- Relevant fields: `Comment.id`, `Comment.body`, `Comment.url`, `Comment.createdAt`, `Comment.updatedAt`, `Comment.parentId`, `Comment.issueId`, `Comment.projectId`, `Comment.projectUpdateId`, `Comment.initiativeId`, `Comment.initiativeUpdateId`, `Comment.documentContentId`, `Comment.user`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `comment list` | `Query.comments` | Read-only |
| `comment get` | `Query.comment` | Read-only |
| `comment bot-actor` | `Comment.botActor` via `Query.comment` | Read-only, bot metadata only |
| `comment children` | `Comment.children` via `Query.comment` | Read-only, body-free metadata |
| `comment created-issues` | `Comment.createdIssues` via `Query.comment` | Read-only |
| `comment update` | `Mutation.commentUpdate` with `CommentUpdateInput` | Resolve the comment, then compare the pinned team through its parent issue; non-issue comments are refused |
| `comment delete` | `Mutation.commentDelete` | Resolve the comment, then compare the pinned team through its parent issue before deleting; non-issue comments are refused |
| `comment resolve` | `Mutation.commentResolve` | Blocked: resolving must first identify and compare the parent issue/project/update/document scope |
| `comment unresolve` | `Mutation.commentUnresolve` | Blocked: unresolving must first identify and compare the parent issue/project/update/document scope |

`comment list`, `comment get`, `comment bot-actor`, `comment children`, and `comment created-issues` (reads) plus `comment update COMMENT_ID --body` and `comment delete COMMENT_ID` (writes) are implemented in the current CLI. Both writes resolve the comment, then compare the pinned team through the comment's parent issue (`guardCommentTarget`): a comment not attached to an issue is refused because the issue guard cannot prove its target. `comment delete` is the one approved delete. Comment child reads omit comment body content by default. Document content and external thread reads remain out of the default surface because they expose content/thread payloads. Issue-scoped comment creation and replies remain under the guarded `issue comment` and `issue reply` commands.

## Project

Schema backing:

- Types: `Project`, `ProjectConnection`
- Reads: `Query.projects`, `Query.project`, `Project.attachments`, `Project.documents`, `Project.externalLinks`, `Project.history`, `Project.initiativeToProjects`, `Project.initiatives`, `Project.inverseRelations`, `Project.issues`, `Project.labels`, `Project.members`, `Project.needs`, `Project.relations`, `Project.teams`, `Project.projectUpdates`
- Writes: `Mutation.projectCreate`, `Mutation.projectUpdate`, `Mutation.projectArchive`
- Inputs: `ProjectCreateInput`, `ProjectUpdateInput`
- Relevant fields: `Project.id`, `Project.name`, `Project.description`, `Project.status`, `Project.lead`, `Project.url`, `Project.teams`, `Project.members`, `Project.documents`, `Project.projectMilestones`, `Project.issues`, `Project.comments`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project list` | `Query.team`, `Team.projects` | Read-only, resolved-team scoped |
| `project all` | `Query.projects` | Read-only |
| `project get` | `Query.project` | Read-only |
| `project open` | `Query.project` resolves `Project.url`, then the platform opener (`xdg-open`/`open`/`rundll32`) launches it with the URL as a discrete argv argument | Read-only |
| `project attachments` | `Project.attachments` | Read-only |
| `project documents` | `Project.documents` | Read-only |
| `project external-links` | `Project.externalLinks` | Read-only |
| `project history` | `Project.history` | Read-only |
| `project initiative-links` | `Project.initiativeToProjects` | Read-only |
| `project initiatives` | `Project.initiatives` | Read-only |
| `project inverse-relations` | `Project.inverseRelations` | Read-only |
| `project issues` | `Project.issues` | Read-only |
| `project comments` | `Project.comments` | Read-only, body-free metadata |
| `project labels` | `Project.labels` | Read-only |
| `project create` | `Mutation.projectCreate` with `ProjectCreateInput.teamIds` | Team-scoped |
| `project update` | `Mutation.projectUpdate` with `ProjectUpdateInput` | Resource-scoped, compare `project_id` |
| `project archive` | `Mutation.projectArchive` | Resource-scoped, compare `project_id` |
| `project members` | `Project.members` plus `Mutation.projectUpdate` with `ProjectUpdateInput.memberIds` | Read-only for list, resource-scoped for writes |
| `project needs` | `Project.needs` | Read-only |
| `project relations` | `Project.relations` | Read-only |
| `project teams` | `Project.teams` | Read-only |
| `project updates` | `Project.projectUpdates` | Read-only, body-free metadata |
| `project filter-suggestion` | `Query.projectFilterSuggestion` | Read-only suggestion payload |

Project is the first implemented PM domain; later domains should reuse its target-comparison vocabulary.

## ProjectUpdate

Use `ProjectUpdate` for Linear project status updates. Avoid calling these generic comments or notes.

Schema backing:

- Types: `ProjectUpdate`, `ProjectUpdateConnection`
- Reads: `Query.projectUpdates`, `Query.projectUpdate`, `Project.projectUpdates`, `ProjectUpdate.comments`
- Writes: `Mutation.projectUpdateCreate`, `Mutation.projectUpdateUpdate`, `Mutation.projectUpdateArchive`
- Relevant fields: `ProjectUpdate.id`, `ProjectUpdate.body`, `ProjectUpdate.health`, `ProjectUpdate.createdAt`, `ProjectUpdate.updatedAt`, `ProjectUpdate.url`, `ProjectUpdate.project`, `ProjectUpdate.user`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-update list` | `Query.projectUpdates` | Read-only |
| `project-update get` | `Query.projectUpdate` | Read-only |
| `project-update comments` | `ProjectUpdate.comments` | Read-only, body-free metadata |
| `project-update create` | `Mutation.projectUpdateCreate` with `ProjectUpdateCreateInput` | Resource-scoped, compare `project_id` (pinned project) and team ownership |
| `project-update update` | `Mutation.projectUpdateUpdate` | Blocked: update must resolve and compare the owning project before mutation |
| `project-update archive` | `Mutation.projectUpdateArchive` | Blocked: destructive command needs explicit safety semantics |

`project-update list`, `project-update get`, `project-update comments`, and `project-update create` are implemented in the current top-level CLI. `project-update create PROJECT --health --body` resolves the pinned project through `requireProject` before posting, so a status update lands only when the resolved project matches the pinned target. `project updates PROJECT_ID` remains the project-scoped history view and omits update body content by default.

## ProjectStatus

Use `ProjectStatus` for Linear project lifecycle status configuration. Do not confuse it with `ProjectUpdate`, which is the user-authored project status update feed.

Schema backing:

- Types: `ProjectStatus`, `ProjectStatusConnection`
- Reads: `Query.projectStatuses`, `Query.projectStatus`
- Writes: `Mutation.projectStatusCreate`, `Mutation.projectStatusUpdate`, `Mutation.projectStatusArchive`, `Mutation.projectStatusUnarchive`
- Relevant fields: `ProjectStatus.id`, `ProjectStatus.name`, `ProjectStatus.description`, `ProjectStatus.type`, `ProjectStatus.color`, `ProjectStatus.position`, `ProjectStatus.archivedAt`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-status list` | `Query.projectStatuses` | Read-only |
| `project-status get` | `Query.projectStatus` | Read-only |
| `project-status project-count` | `Query.projectStatusProjectCount` | Read-only count payload |
| `project-status create` | `Mutation.projectStatusCreate` | Blocked: organization project status configuration needs an explicit admin safety model |
| `project-status update` | `Mutation.projectStatusUpdate` | Blocked: update must resolve and compare the owning organization before mutation |
| `project-status archive` | `Mutation.projectStatusArchive` | Blocked: destructive command needs explicit safety semantics |
| `project-status unarchive` | `Mutation.projectStatusUnarchive` | Blocked: restore semantics need an explicit admin safety model |

Only `project-status list`, `project-status get`, and `project-status project-count` are implemented in the current CLI. ProjectStatus writes are deferred as organization/admin configuration surface.

## ProjectLabel

Use `ProjectLabel` for Linear project labels. Do not confuse it with issue labels.

Schema backing:

- Types: `ProjectLabel`, `ProjectLabelConnection`
- Reads: `Query.projectLabels`, `Query.projectLabel`, `ProjectLabel.children`, `ProjectLabel.projects`
- Writes: `Mutation.projectLabelCreate`, `Mutation.projectLabelUpdate`, `Mutation.projectLabelDelete`, `Mutation.projectLabelRetire`, `Mutation.projectLabelRestore`
- Relevant fields: `ProjectLabel.id`, `ProjectLabel.name`, `ProjectLabel.description`, `ProjectLabel.color`, `ProjectLabel.isGroup`, `ProjectLabel.parent`, `ProjectLabel.retiredAt`, `ProjectLabel.archivedAt`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-label list` | `Query.projectLabels` | Read-only |
| `project-label get` | `Query.projectLabel` | Read-only |
| `project-label children` | `ProjectLabel.children` via `Query.projectLabel` | Read-only |
| `project-label projects` | `ProjectLabel.projects` via `Query.projectLabel` | Read-only |
| `project-label create` | `Mutation.projectLabelCreate` | Blocked: organization label configuration needs an explicit admin safety model |
| `project-label update` | `Mutation.projectLabelUpdate` | Blocked: update must resolve and compare the owning organization before mutation |
| `project-label delete` | `Mutation.projectLabelDelete` | Blocked: destructive command needs explicit safety semantics |
| `project-label retire` | `Mutation.projectLabelRetire` | Blocked: lifecycle command needs explicit admin safety semantics |
| `project-label restore` | `Mutation.projectLabelRestore` | Blocked: restore semantics need an explicit admin safety model |

Only `project-label list`, `project-label get`, `project-label children`, and `project-label projects` are implemented in the current CLI. ProjectLabel writes are deferred as organization/admin configuration surface.

## ProjectRelation

Use `ProjectRelation` for Linear dependency relations between Projects. It is project graph metadata, not issue dependency state.

Schema backing:

- Types: `ProjectRelation`, `ProjectRelationConnection`
- Reads: `Query.projectRelations`, `Query.projectRelation`
- Writes: `Mutation.projectRelationCreate`, `Mutation.projectRelationUpdate`, `Mutation.projectRelationDelete`
- Inputs: `ProjectRelationCreateInput`, `ProjectRelationUpdateInput`
- Relevant fields: `ProjectRelation.id`, `ProjectRelation.type`, `ProjectRelation.project`, `ProjectRelation.projectMilestone`, `ProjectRelation.anchorType`, `ProjectRelation.relatedProject`, `ProjectRelation.relatedProjectMilestone`, `ProjectRelation.relatedAnchorType`, `ProjectRelation.createdAt`, `ProjectRelation.updatedAt`, `ProjectRelation.archivedAt`, `ProjectRelation.user`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-relation list` | `Query.projectRelations` | Read-only |
| `project-relation get` | `Query.projectRelation` | Read-only |
| `project-relation create` | `Mutation.projectRelationCreate` | Blocked: create must resolve and compare both project dependency endpoints before mutation |
| `project-relation update` | `Mutation.projectRelationUpdate` | Blocked: update must resolve and compare both project dependency endpoints before mutation |
| `project-relation delete` | `Mutation.projectRelationDelete` | Blocked: destructive command needs explicit project dependency safety semantics |

Only `project-relation list` and `project-relation get` are implemented in the current CLI. ProjectRelation writes are deferred until their endpoint guard model is explicit.

## Cycle

Schema backing:

- Types: `Cycle`, `CycleConnection`
- Reads: `Query.cycles`, `Query.cycle`, `Cycle.issues`, `Cycle.uncompletedIssuesUponClose`, `Team.cycles`
- Writes: `Mutation.cycleCreate`, `Mutation.cycleUpdate`, `Mutation.cycleArchive`, `Mutation.cycleShiftAll`, `Mutation.cycleStartUpcomingCycleToday`
- Relevant fields: `Cycle.id`, `Cycle.number`, `Cycle.name`, `Cycle.startsAt`, `Cycle.endsAt`, `Cycle.completedAt`, `Cycle.team`, `Cycle.issues`, `Cycle.uncompletedIssuesUponClose`, `Cycle.progress`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `cycle list` | `Query.cycles` | Read-only |
| `cycle get` | `Query.cycle` | Read-only |
| `cycle issues` | `Cycle.issues` | Read-only |
| `cycle uncompleted-issues` | `Cycle.uncompletedIssuesUponClose` | Read-only |
| `cycle create` | `Mutation.cycleCreate` | Team-scoped |
| `cycle update` | `Mutation.cycleUpdate` | Team-scoped |
| `cycle archive` | `Mutation.cycleArchive` | Team-scoped |

## Sprint

`sprint` is not a Linear schema type. It is a report alias over `Cycle` only.

Schema backing:

- Types: `Cycle`, `Issue`
- Reads: `Query.cycles`, `Query.cycle`, `Cycle.issues`, `Cycle.progress`
- Relevant fields: `Cycle.number`, `Cycle.name`, `Cycle.startsAt`, `Cycle.endsAt`, `Cycle.completedAt`, `Cycle.progress`, `Issue.identifier`, `Issue.title`, `Issue.state`, `Issue.assignee`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `sprint current` | `Query.cycles` filtered to active/current cycles | Read-only |
| `sprint report` | `Query.cycle` plus `Cycle.issues` | Read-only |

No `sprint create`, `sprint update`, or `sprint archive` command exists. Use `cycle` for Cycle mutations.

## ProjectMilestone

Use the schema name `ProjectMilestone` in code and docs. Avoid the loose name `milestone`.

Schema backing:

- Types: `ProjectMilestone`, `ProjectMilestoneConnection`
- Reads: `Query.projectMilestones`, `Query.projectMilestone`, `Project.projectMilestones`, `ProjectMilestone.issues`
- Writes: `Mutation.projectMilestoneCreate`, `Mutation.projectMilestoneUpdate`, `Mutation.projectMilestoneDelete`
- Inputs: `ProjectMilestoneCreateInput`, `ProjectMilestoneUpdateInput`
- Relevant fields: `ProjectMilestone.id`, `ProjectMilestone.name`, `ProjectMilestone.description`, `ProjectMilestone.targetDate`, `ProjectMilestone.status`, `ProjectMilestone.project`, `ProjectMilestone.sortOrder`, `ProjectMilestone.issues`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-milestone all` | `Query.projectMilestones` | Read-only |
| `project-milestone list` | `Project.projectMilestones` via `Query.project` | Read-only |
| `project-milestone get` | `Query.projectMilestone` | Read-only |
| `project-milestone issues` | `ProjectMilestone.issues` | Read-only |
| `project-milestone create` | `Mutation.projectMilestoneCreate` with `projectId` | Resource-scoped, compare `project_id` |
| `project-milestone update` | `Mutation.projectMilestoneUpdate` | Resource-scoped, compare resolved project |
| `project-milestone delete` | `Mutation.projectMilestoneDelete` | Resource-scoped, compare resolved project |

## Document

Schema backing:

- Types: `Document`, `DocumentConnection`
- Reads: `Query.documents`, `Query.document`, `Document.comments`, `Project.documents`, `Team.documents`, `Issue.documents`, `Cycle.documents`
- Writes: `Mutation.documentCreate`, `Mutation.documentUpdate`, `Mutation.documentDelete`
- Inputs: `DocumentCreateInput`, `DocumentUpdateInput`
- Relevant fields: `Document.id`, `Document.title`, `Document.slugId`, `Document.archivedAt`, `Document.project`, `Document.team`, `Document.issue`, `Document.cycle`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `document list` | `Query.documents` | Read-only |
| `document get` | `Query.document` | Read-only |
| `document comments` | `Document.comments` | Read-only, body-free metadata |
| `document create` | `Mutation.documentCreate` with `DocumentCreateInput.teamId` from the resolved team and optional `projectId` from the pinned project; `--content` (or `--content-file`, or `--content -` for stdin) | Team-scoped unless a `project_id` is pinned |
| `document update` | `Mutation.documentUpdate`; resolves the existing document via `Query.document` and compares its `team` (and pinned `project`) before mutating | Resource-scoped, compare team and pinned project |
| `document delete` | `Mutation.documentDelete` | Blocked: destructive command needs explicit safety semantics |

`document list`, `document get`, `document comments`, `document create`, and `document update` are implemented in the current CLI. Document comment reads omit comment body content by default. `document create` is a team-scoped guarded write (carrying the pinned project when set); `document update` resolves the existing document and fails closed unless its team — and the pinned project, when configured — match. Document delete stays deferred.

## Label

CLI name `label` maps to Linear schema type `IssueLabel`.

Schema backing:

- Types: `IssueLabel`, `IssueLabelConnection`
- Reads: `Query.issueLabels`, `Query.issueLabel`, `IssueLabel.children`, `IssueLabel.issues`, `Team.labels`
- Writes: `Mutation.issueLabelCreate`, `Mutation.issueLabelUpdate`, `Mutation.issueLabelDelete`
- Inputs: `IssueLabelCreateInput`, `IssueLabelUpdateInput`
- Relevant fields: `IssueLabel.id`, `IssueLabel.name`, `IssueLabel.description`, `IssueLabel.color`, `IssueLabel.team`, `IssueLabel.children`, `IssueLabel.issues`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `label list` | `Query.issueLabels` | Read-only |
| `label get` | `Query.issueLabel` | Read-only |
| `label children` | `IssueLabel.children` | Read-only |
| `label issues` | `IssueLabel.issues` | Read-only |
| `label create` | `Mutation.issueLabelCreate` with optional `teamId` | Blocked: optional team scope needs explicit org/team target behavior before writes |
| `label update` | `Mutation.issueLabelUpdate` | Blocked: update must resolve and compare the label's owning team before mutation |
| `label delete` | `Mutation.issueLabelDelete` | Blocked: destructive command needs explicit safety semantics |

Only read-only label commands are implemented in the current CLI. Label writes are deferred until the team-scope guard is designed.

## Team

Schema backing:

- Types: `Team`, `TeamConnection`, `TeamMembership`, `TeamMembershipConnection`
- Reads: `Query.teams`, `Query.team`, `Team.cycles`, `Team.issues`, `Team.labels`, `Team.members`, `Team.memberships`, `Team.projects`, `Team.releasePipelines`, `Team.states`, `Team.templates`, `Query.teamMemberships`, `Query.teamMembership`
- Writes: `Mutation.teamCreate`, `Mutation.teamUpdate`, `Mutation.teamDelete`, `Mutation.teamMembershipCreate`, `Mutation.teamMembershipUpdate`, `Mutation.teamMembershipDelete`
- Inputs: `TeamCreateInput`, `TeamUpdateInput`
- Relevant fields: `Team.id`, `Team.name`, `Team.key`, `Team.description`, `Team.archivedAt`, `Team.issues`, `Team.cycles`, `Team.members`, `Team.projects`, `TeamMembership.id`, `TeamMembership.user`, `TeamMembership.team`, `TeamMembership.owner`, `TeamMembership.sortOrder`, `TeamMembership.archivedAt`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `team list` | `Query.teams` | Read-only |
| `team get` | `Query.team` | Read-only |
| `team create` | `Mutation.teamCreate` | Blocked: organization administration surface needs an explicit admin safety model |
| `team update` | `Mutation.teamUpdate` | Blocked: team metadata writes need stronger authority checks than ordinary target comparison |
| `team delete` | `Mutation.teamDelete` | Blocked: destructive command needs explicit safety semantics |
| `team cycles` | `Team.cycles` | Read-only |
| `team issues` | `Team.issues` | Read-only |
| `team labels` | `Team.labels` | Read-only |
| `team members` | `Team.members` | Read-only |
| `team memberships` | `Team.memberships` | Read-only |
| `team projects` | `Team.projects` | Read-only |
| `team release-pipelines` | `Team.releasePipelines` | Read-only |
| `team states` | `Team.states` | Read-only |
| `team git-automation-states` | `Team.gitAutomationStates` | Read-only, rule/state/target-branch metadata only |
| `team templates` | `Team.templates` | Read-only |
| `team-membership list` | `Query.teamMemberships` | Read-only |
| `team-membership get` | `Query.teamMembership` | Read-only |
| `team-membership create` | `Mutation.teamMembershipCreate` | Blocked: organization membership administration needs an explicit admin safety model |
| `team-membership update` | `Mutation.teamMembershipUpdate` | Blocked: update must resolve and compare the membership's team and organization before mutation |
| `team-membership delete` | `Mutation.teamMembershipDelete` | Blocked: destructive membership command needs explicit admin safety semantics |

Team list/get, the read-only Team child-list commands above, and `team-membership list/get` are implemented in the current CLI. `team git-automation-states` exposes rule/state/target-branch metadata without write controls. Team creation, metadata mutation, and membership writes are deferred as organization/admin surface.

## User

Schema backing:

- Types: `User`, `UserConnection`, `UserSettings`
- Reads: `Query.users`, `Query.user`, `Query.viewer`, `Query.userSettings`, `User.assignedIssues`, `User.createdIssues`, `User.delegatedIssues`, `User.teamMemberships`, `User.teams`, `User.drafts`, `Team.members`, `Project.members`, `UserSettings.notificationCategoryPreferences`, `UserSettings.notificationChannelPreferences`, `UserSettings.notificationDeliveryPreferences`, `UserSettings.theme`
- Relevant fields: `User.id`, `User.name`, `User.displayName`, `User.email`, `User.active`, `User.guest`, `User.admin`, `User.url`, `User.assignedIssues`, `User.teams`, `Draft.id`, `Draft.issue`, `Draft.project`, `Draft.projectUpdate`, `Draft.initiative`, `Draft.initiativeUpdate`, `Draft.parentComment`, `Draft.customerNeed`, `Draft.team`, `UserSettings.id`, `UserSettings.user.id`, notification channel booleans, mobile delivery windows, theme preset and custom color values

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `user list` | `Query.users` | Read-only |
| `user get` | `Query.user` | Read-only |
| `user me` | `Query.viewer` | Read-only |
| `user drafts` | `User.drafts` via `Query.viewer` | Read-only |
| `user settings get` | `Query.userSettings` | Read-only |
| `user settings notification-categories` | `Query.userSettings.notificationCategoryPreferences` | Read-only |
| `user settings notification-category CATEGORY` | `Query.userSettings.notificationCategoryPreferences.<category>` | Read-only |
| `user settings notification-channels` | `Query.userSettings.notificationChannelPreferences` | Read-only |
| `user settings notification-delivery` | `Query.userSettings.notificationDeliveryPreferences` | Read-only |
| `user settings mobile-delivery` | `Query.userSettings.notificationDeliveryPreferences.mobile` | Read-only |
| `user settings mobile-schedule` | `Query.userSettings.notificationDeliveryPreferences.mobile.schedule` | Read-only |
| `user settings mobile-schedule-day DAY` | `Query.userSettings.notificationDeliveryPreferences.mobile.schedule.<day>` | Read-only |
| `user settings theme` | `Query.userSettings.theme` | Read-only |
| `user settings custom-theme` | `Query.userSettings.theme.custom` | Read-only |
| `user settings custom-sidebar-theme` | `Query.userSettings.theme.custom.sidebar` | Read-only |
| `user assigned-issues` | `User.assignedIssues` | Read-only |
| `user created-issues` | `User.createdIssues` | Read-only |
| `user delegated-issues` | `User.delegatedIssues` | Read-only |
| `user team-memberships` | `User.teamMemberships` | Read-only |
| `user teams` | `User.teams` | Read-only |
| `user my-assigned-issues` | `User.assignedIssues` via `Query.viewer` | Read-only |
| `user my-created-issues` | `User.createdIssues` via `Query.viewer` | Read-only |
| `user my-delegated-issues` | `User.delegatedIssues` via `Query.viewer` | Read-only |
| `user my-team-memberships` | `User.teamMemberships` via `Query.viewer` | Read-only |
| `user my-teams` | `User.teams` via `Query.viewer` | Read-only |

The read-only user commands are implemented in the current CLI. Draft reads intentionally omit draft body/data and return parent metadata only. User settings reads omit calendar hashes, raw unsubscribe arrays, and user email by default; they expose compact preference booleans, delivery windows, and theme colors. User issue-list commands use compact issue summaries and do not include body content. User writes are not part of the v1 PM command surface until a later slice proves the exact Linear mutation and safety semantics.

## WorkflowState

Use the schema name `WorkflowState` in code and docs. It is Linear's issue status entity.

Schema backing:

- Types: `WorkflowState`, `WorkflowStateConnection`
- Reads: `Query.workflowStates`, `Query.workflowState`, `Team.states`
- Writes: `Mutation.workflowStateCreate`, `Mutation.workflowStateUpdate`, `Mutation.workflowStateArchive`
- Inputs: `WorkflowStateCreateInput`, `WorkflowStateUpdateInput`
- Relevant fields: `WorkflowState.id`, `WorkflowState.name`, `WorkflowState.type`, `WorkflowState.color`, `WorkflowState.position`, `WorkflowState.team`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `workflow-state list` | `Query.workflowStates` | Read-only |
| `workflow-state get` | `Query.workflowState` | Read-only |
| `workflow-state issues` | `WorkflowState.issues` via `Query.workflowState` | Read-only |
| `workflow-state create` | `Mutation.workflowStateCreate` | Blocked: team workflow configuration needs an explicit admin safety model |
| `workflow-state update` | `Mutation.workflowStateUpdate` | Blocked: update must resolve and compare the owning team before mutation |
| `workflow-state archive` | `Mutation.workflowStateArchive` | Blocked: destructive command needs explicit safety semantics |

`workflow-state list`, `workflow-state get`, and `workflow-state issues` are implemented in the current CLI. WorkflowState writes are deferred as team/admin configuration surface.

## TimeSchedule

Use the schema name `TimeSchedule` in code and docs. It is Linear's on-call or availability schedule used by triage responsibilities.

Schema backing:

- Types: `TimeSchedule`, `TimeScheduleConnection`, `TimeScheduleEntry`
- Reads: `Query.timeSchedules`, `Query.timeSchedule`
- Writes: `Mutation.timeScheduleCreate`, `Mutation.timeScheduleUpdate`, `Mutation.timeScheduleDelete`, `Mutation.timeScheduleUpsertExternal`
- Inputs: `TimeScheduleCreateInput`, `TimeScheduleUpdateInput`, `TimeScheduleEntryInput`
- Relevant fields: `TimeSchedule.id`, `TimeSchedule.name`, `TimeSchedule.externalId`, `TimeSchedule.externalUrl`, `TimeSchedule.integration`, `TimeSchedule.entries`, `TimeSchedule.createdAt`, `TimeSchedule.updatedAt`, `TimeSchedule.archivedAt`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `time-schedule list` | `Query.timeSchedules` | Read-only |
| `time-schedule get` | `Query.timeSchedule` | Read-only |
| `time-schedule create` | `Mutation.timeScheduleCreate` | Blocked: schedule create needs explicit owner/admin safety semantics |
| `time-schedule update` | `Mutation.timeScheduleUpdate` | Blocked: update must resolve schedule scope before mutation |
| `time-schedule delete` | `Mutation.timeScheduleDelete` | Blocked: destructive command needs explicit safety semantics |
| `time-schedule upsert-external` | `Mutation.timeScheduleUpsertExternal` | Blocked: external integration sync surface is not an ordinary agent workflow |

Only `time-schedule list` and `time-schedule get` are implemented in the current CLI. TimeSchedule writes and external upserts are deferred as triage/admin configuration surface.

## TriageResponsibility

Use the schema name `TriageResponsibility` in code and docs. It is Linear's team triage assignment or notification responsibility configuration.

Schema backing:

- Types: `TriageResponsibility`, `TriageResponsibilityConnection`, `TriageResponsibilityManualSelection`
- Reads: `Query.triageResponsibilities`, `Query.triageResponsibility`, `TriageResponsibility.manualSelection`
- Writes: `Mutation.triageResponsibilityCreate`, `Mutation.triageResponsibilityUpdate`, `Mutation.triageResponsibilityDelete`
- Inputs: `TriageResponsibilityCreateInput`, `TriageResponsibilityUpdateInput`
- Relevant fields: `TriageResponsibility.id`, `TriageResponsibility.action`, `TriageResponsibility.team`, `TriageResponsibility.timeSchedule`, `TriageResponsibility.currentUser`, `TriageResponsibility.manualSelection`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `triage-responsibility list` | `Query.triageResponsibilities` | Read-only |
| `triage-responsibility get` | `Query.triageResponsibility` | Read-only |
| `triage-responsibility manual-selection` | `TriageResponsibility.manualSelection` via `Query.triageResponsibility` | Read-only |
| `triage-responsibility create` | `Mutation.triageResponsibilityCreate` | Blocked: team triage configuration needs an explicit admin safety model |
| `triage-responsibility update` | `Mutation.triageResponsibilityUpdate` | Blocked: update must resolve and compare the owning team before mutation |
| `triage-responsibility delete` | `Mutation.triageResponsibilityDelete` | Blocked: destructive team triage configuration command needs explicit safety semantics |

Only `triage-responsibility list`, `triage-responsibility get`, and `triage-responsibility manual-selection` are implemented in the current CLI. TriageResponsibility writes are deferred as team/admin configuration surface.

## SLA Configuration

Use the command name `sla-configuration` for Linear's `SlaConfiguration` schema type. It is an active SLA rule that can apply to a team.

Schema backing:

- Types: `SlaConfiguration`
- Reads: `Query.slaConfigurations`
- Writes: no direct write operation is implemented in linctl
- Relevant fields: `SlaConfiguration.id`, `SlaConfiguration.name`, `SlaConfiguration.conditions`, `SlaConfiguration.sla`, `SlaConfiguration.slaType`, `SlaConfiguration.removesSla`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `sla-configuration list` | `Query.slaConfigurations` | Read-only |

Only `sla-configuration list` is implemented in the current CLI. SLA rule changes remain part of team/admin workflow configuration and do not have a guarded write surface.

## SemanticSearch

Use the command name `semantic-search` for Linear's semantic search query. It searches visible issues, projects, initiatives, and documents and returns compact references only.

Schema backing:

- Types: `SemanticSearchPayload`, `SemanticSearchResult`, `SemanticSearchResultType`
- Reads: `Query.semanticSearch`
- Writes: no write operation exists
- Relevant fields: `SemanticSearchResult.type`, `SemanticSearchResult.id`, and compact reference fields from `Issue`, `Project`, `Initiative`, and `Document`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `semantic-search` | `Query.semanticSearch` | Read-only |

Only `semantic-search` is implemented in the current CLI. Results intentionally omit body/content snippets so the command stays a compact reference lookup.

## Search

Use `search` for Linear's typed full-text/vector search roots. These commands return compact result summaries and intentionally omit archive payloads, metadata, comments, descriptions, document content, and project content by default.

Schema backing:

- Types: `DocumentSearchPayload`, `IssueSearchPayload`, `ProjectSearchPayload`, `DocumentSearchResult`, `IssueSearchResult`, `ProjectSearchResult`
- Reads: `Query.searchDocuments`, `Query.searchIssues`, `Query.searchProjects`
- Writes: none
- Inputs: required `term`, optional pagination in the client layer
- Relevant fields: compact identity, URL, status/team/project/parent fields only

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `search documents` | `Query.searchDocuments` | Read-only |
| `search issues` | `Query.searchIssues` | Read-only |
| `search projects` | `Query.searchProjects` | Read-only |

Only typed read commands are implemented. Archive payload and metadata variants are intentionally not exposed as default workflow.

## Template

Use the schema name `Template` in code and docs. It is Linear's reusable issue, project, document, and release-note template entity.

Schema backing:

- Types: `Template`
- Reads: `Query.templates`, `Query.template`
- Writes: `Mutation.templateCreate`, `Mutation.templateUpdate`, `Mutation.templateDelete`
- Inputs: `TemplateCreateInput`, `TemplateUpdateInput`
- Relevant fields: `Template.id`, `Template.name`, `Template.type`, `Template.description`, `Template.icon`, `Template.color`, `Template.sortOrder`, `Template.lastAppliedAt`, `Template.team`, `Template.pipeline`, `Template.creator`, `Template.lastUpdatedBy`, `Template.inheritedFrom`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `template list` | `Query.templates` | Read-only |
| `template get` | `Query.template` | Read-only |
| `template create` | `Mutation.templateCreate` | Blocked: create can be organization-, team-, or pipeline-scoped and needs explicit guard semantics |
| `template update` | `Mutation.templateUpdate` | Blocked: update must resolve and compare the template's organization, team, or pipeline scope before mutation |
| `template delete` | `Mutation.templateDelete` | Blocked: destructive command needs explicit template-scope safety semantics |

Only `template list` and `template get` are implemented in the current CLI. Template writes are deferred until their organization, team, and pipeline guard model is explicit.

## Initiative

Use the schema name `Initiative` in code and docs. It is Linear's current strategic grouping of projects toward a goal. Use Initiative, InitiativeToProject, and InitiativeUpdate for new planning workflows.

Schema backing:

- Types: `Initiative`, `InitiativeConnection`, `InitiativeHistory`, `InitiativeUpdate`, `EntityExternalLink`, `Document`, `Project`
- Reads: `Query.initiatives`, `Query.initiative`, `Initiative.history`, `Initiative.links`, `Initiative.subInitiatives`, `Initiative.initiativeUpdates`, `Initiative.documents`, `Initiative.projects`
- Writes: `Mutation.createInitiative`, `Mutation.updateInitiative`, `Mutation.archiveInitiative`, `Mutation.deleteInitiative`
- Inputs: `InitiativeCreateInput`, `InitiativeUpdateInput`
- Relevant fields: `Initiative.id`, `Initiative.name`, `Initiative.description`, `Initiative.status`, `Initiative.priority`, `Initiative.targetDate`, `Initiative.slugId`, `Initiative.url`, compact `Document` identity and parent fields, compact `Project` identity, status, lead, and team fields

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `initiative list` | `Query.initiatives` | Read-only |
| `initiative get` | `Query.initiative` | Read-only |
| `initiative history` | `Initiative.history` via `Query.initiative` | Read-only |
| `initiative links` | `Initiative.links` via `Query.initiative` | Read-only |
| `initiative sub-initiatives` | `Initiative.subInitiatives` via `Query.initiative` | Read-only |
| `initiative updates` | `Initiative.initiativeUpdates` via `Query.initiative` | Read-only |
| `initiative documents` | `Initiative.documents` via `Query.initiative` | Read-only |
| `initiative projects` | `Initiative.projects` via `Query.initiative` | Read-only direct projects |
| `initiative create` | `Mutation.createInitiative` | Blocked: initiative create needs an explicit organization-scoped safety model |
| `initiative update` | `Mutation.updateInitiative` | Blocked: update must resolve and compare the owning organization before mutation |
| `initiative archive` | `Mutation.archiveInitiative` | Blocked: destructive command needs explicit safety semantics |

Initiative writes are deferred as organization-scoped planning surface.

## InitiativeRelation

Use `InitiativeRelation` for Linear parent-child Initiative hierarchy edges.

Schema backing:

- Types: `InitiativeRelation`, `InitiativeRelationConnection`
- Reads: `Query.initiativeRelations`, `Query.initiativeRelation`
- Writes: `Mutation.initiativeRelationCreate`, `Mutation.initiativeRelationUpdate`, `Mutation.initiativeRelationDelete`
- Inputs: `InitiativeRelationCreateInput`, `InitiativeRelationUpdateInput`
- Relevant fields: `InitiativeRelation.id`, `InitiativeRelation.initiative`, `InitiativeRelation.relatedInitiative`, `InitiativeRelation.sortOrder`, `InitiativeRelation.createdAt`, `InitiativeRelation.updatedAt`, `InitiativeRelation.archivedAt`, `InitiativeRelation.user`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `initiative-relation list` | `Query.initiativeRelations` | Read-only |
| `initiative-relation get` | `Query.initiativeRelation` | Read-only |
| `initiative-relation create` | `Mutation.initiativeRelationCreate` | Blocked: create must resolve and compare both Initiative hierarchy endpoints before mutation |
| `initiative-relation update` | `Mutation.initiativeRelationUpdate` | Blocked: update must resolve and compare both Initiative hierarchy endpoints before mutation |
| `initiative-relation delete` | `Mutation.initiativeRelationDelete` | Blocked: destructive command needs explicit hierarchy safety semantics |

Only `initiative-relation list` and `initiative-relation get` are implemented in the current CLI. InitiativeRelation writes are deferred until their hierarchy guard model is explicit.

## InitiativeToProject

Use `InitiativeToProject` for Linear associations between Initiatives and Projects.

Schema backing:

- Types: `InitiativeToProject`, `InitiativeToProjectConnection`
- Reads: `Query.initiativeToProjects`, `Query.initiativeToProject`
- Writes: `Mutation.initiativeToProjectCreate`, `Mutation.initiativeToProjectUpdate`, `Mutation.initiativeToProjectDelete`
- Inputs: `InitiativeToProjectCreateInput`, `InitiativeToProjectUpdateInput`
- Relevant fields: `InitiativeToProject.id`, `InitiativeToProject.initiative`, `InitiativeToProject.project`, `InitiativeToProject.sortOrder`, `InitiativeToProject.createdAt`, `InitiativeToProject.updatedAt`, `InitiativeToProject.archivedAt`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `initiative-to-project list` | `Query.initiativeToProjects` | Read-only |
| `initiative-to-project get` | `Query.initiativeToProject` | Read-only |
| `initiative-to-project create` | `Mutation.initiativeToProjectCreate` | Blocked: create must resolve and compare both Initiative and Project endpoints before mutation |
| `initiative-to-project update` | `Mutation.initiativeToProjectUpdate` | Blocked: update must resolve and compare both Initiative and Project endpoints before mutation |
| `initiative-to-project delete` | `Mutation.initiativeToProjectDelete` | Blocked: destructive command needs explicit association safety semantics |

Only `initiative-to-project list` and `initiative-to-project get` are implemented in the current CLI. InitiativeToProject writes are deferred until their endpoint guard model is explicit.

## RoadmapToProject

`RoadmapToProject` is the deprecated Linear association between Roadmaps and Projects. Keep the read commands for compatibility; prefer `InitiativeToProject` for new workflows when Linear offers both surfaces.

Schema backing:

- Types: `RoadmapToProject`, `RoadmapToProjectConnection`
- Reads: `Query.roadmapToProjects`, `Query.roadmapToProject`
- Writes: `Mutation.roadmapToProjectCreate`, `Mutation.roadmapToProjectUpdate`, `Mutation.roadmapToProjectDelete`
- Inputs: `RoadmapToProjectCreateInput`, `RoadmapToProjectUpdateInput`
- Relevant fields: `RoadmapToProject.id`, `RoadmapToProject.roadmap`, `RoadmapToProject.project`, `RoadmapToProject.sortOrder`, `RoadmapToProject.createdAt`, `RoadmapToProject.updatedAt`, `RoadmapToProject.archivedAt`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `roadmap-to-project list` | `Query.roadmapToProjects` | Legacy read-only |
| `roadmap-to-project get` | `Query.roadmapToProject` | Legacy read-only |
| `roadmap-to-project create` | `Mutation.roadmapToProjectCreate` | Blocked: deprecated create must resolve and compare both Roadmap and Project endpoints before mutation |
| `roadmap-to-project update` | `Mutation.roadmapToProjectUpdate` | Blocked: deprecated update must resolve and compare both Roadmap and Project endpoints before mutation |
| `roadmap-to-project delete` | `Mutation.roadmapToProjectDelete` | Blocked: destructive deprecated association command needs explicit safety semantics |

Only `roadmap-to-project list` and `roadmap-to-project get` are implemented as legacy reads in the current CLI. RoadmapToProject writes are deferred until their endpoint guard model is explicit.

## InitiativeUpdate

Use `InitiativeUpdate` for Linear initiative status updates. Avoid calling these generic comments or notes.

Schema backing:

- Types: `InitiativeUpdate`, `InitiativeUpdateConnection`
- Reads: `Query.initiativeUpdates`, `Query.initiativeUpdate`, `InitiativeUpdate.comments`
- Writes: `Mutation.initiativeUpdateCreate`, `Mutation.initiativeUpdateUpdate`, `Mutation.initiativeUpdateArchive`, `Mutation.initiativeUpdateUnarchive`
- Inputs: `InitiativeUpdateCreateInput`, `InitiativeUpdateUpdateInput`
- Relevant fields: `InitiativeUpdate.id`, `InitiativeUpdate.body`, `InitiativeUpdate.health`, `InitiativeUpdate.createdAt`, `InitiativeUpdate.updatedAt`, `InitiativeUpdate.url`, `InitiativeUpdate.slugId`, `InitiativeUpdate.commentCount`, `InitiativeUpdate.initiative`, `InitiativeUpdate.user`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `initiative-update list` | `Query.initiativeUpdates` | Read-only |
| `initiative-update get` | `Query.initiativeUpdate` | Read-only |
| `initiative-update comments` | `InitiativeUpdate.comments` | Read-only, body-free metadata |
| `initiative-update create` | `Mutation.initiativeUpdateCreate` | Blocked: create must resolve and compare the owning Initiative before posting |
| `initiative-update update` | `Mutation.initiativeUpdateUpdate` | Blocked: update must resolve and compare the owning Initiative before mutation |
| `initiative-update archive` | `Mutation.initiativeUpdateArchive` | Blocked: destructive command needs explicit safety semantics |
| `initiative-update unarchive` | `Mutation.initiativeUpdateUnarchive` | Blocked: unarchive needs explicit lifecycle and target semantics |

`initiative-update list`, `initiative-update get`, and `initiative-update comments` are implemented in the current CLI. InitiativeUpdate comment reads omit comment body content by default. InitiativeUpdate writes and reminders are deferred until their guard model is explicit.

## Roadmap

Use the schema name `Roadmap` in code and docs for legacy compatibility only. Roadmap is Linear's deprecated grouping for projects; use `Initiative` for new planning workflows.

Schema backing:

- Types: `Roadmap`, `RoadmapConnection`
- Reads: `Query.roadmaps`, `Query.roadmap`, `Roadmap.projects`
- Writes: `Mutation.roadmapCreate`, `Mutation.roadmapUpdate`, `Mutation.roadmapArchive`, `Mutation.roadmapDelete`
- Inputs: `RoadmapCreateInput`, `RoadmapUpdateInput`
- Relevant fields: `Roadmap.id`, `Roadmap.name`, `Roadmap.description`, `Roadmap.color`, `Roadmap.slugId`, `Roadmap.sortOrder`, `Roadmap.url`, `Roadmap.creator`, `Roadmap.owner`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `roadmap list` | `Query.roadmaps` | Legacy read-only |
| `roadmap get` | `Query.roadmap` | Legacy read-only |
| `roadmap projects` | `Roadmap.projects` via `Query.roadmap` | Legacy read-only |
| `roadmap create` | `Mutation.roadmapCreate` | Blocked: deprecated organization-scoped planning surface needs an explicit safety model |
| `roadmap update` | `Mutation.roadmapUpdate` | Blocked: update must resolve and compare the owning organization before mutation |
| `roadmap archive` | `Mutation.roadmapArchive` | Blocked: destructive command needs explicit safety semantics |
| `roadmap delete` | `Mutation.roadmapDelete` | Blocked: destructive command needs explicit safety semantics |

`roadmap list`, `roadmap get`, and `roadmap projects` are implemented as legacy reads in the current CLI. Roadmap writes and roadmap-project association writes are deferred; prefer Initiative commands for current Linear planning workflows.

## CustomView

Use the schema name `CustomView` in code and docs. It is Linear's saved view over issues, projects, or initiatives.

Schema backing:

- Types: `CustomView`, `CustomViewConnection`
- Reads: `Query.customViews`, `Query.customView`, `Query.customViewHasSubscribers`, `Query.customView_initiatives`,
  `Query.customView_organizationViewPreferences`, `Query.customView_organizationViewPreferences_preferences`,
  `Query.customView_viewPreferencesValues`
- Writes: `Mutation.createCustomView`, `Mutation.updateCustomView`, `Mutation.deleteCustomView`
- Inputs: `CustomViewCreateInput`, `CustomViewUpdateInput`
- Relevant fields: `CustomView.id`, `CustomView.name`, `CustomView.description`, `CustomView.modelName`, `CustomView.shared`, `CustomView.color`, `CustomView.slugId`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `custom-view list` | `Query.customViews` | Read-only |
| `custom-view subscribers` | `Query.customViewHasSubscribers` | Read-only |
| `custom-view get` | `Query.customView` | Read-only |
| `custom-view initiatives` | `Query.customView_initiatives` | Read-only |
| `custom-view issues` | `Query.customView_issues` | Read-only |
| `custom-view organization-preferences` | `Query.customView_organizationViewPreferences` | Read-only |
| `custom-view organization-preferences values` | `Query.customView_organizationViewPreferences_preferences` | Read-only |
| `custom-view projects` | `Query.customView_projects` | Read-only |
| `custom-view user-preferences` | `Query.customView_userViewPreferences` | Read-only |
| `custom-view user-preferences values` | `Query.customView_userViewPreferences_preferences` | Read-only |
| `custom-view preference-values` | `Query.customView_viewPreferencesValues` | Read-only |
| `custom-view create` | `Mutation.createCustomView` | Blocked: custom view create needs an explicit organization-scoped safety model |
| `custom-view update` | `Mutation.updateCustomView` | Blocked: update must resolve and compare the owning organization before mutation |
| `custom-view delete` | `Mutation.deleteCustomView` | Blocked: destructive command needs explicit safety semantics |

Only CustomView reads are implemented in the current CLI. CustomView writes are deferred as organization-scoped view configuration surface.

## Customer

Use the schema name `Customer` in code and docs. It is Linear's customer organization record for customer requests and feedback.

Schema backing:

- Types: `Customer`, `CustomerConnection`, `CustomerNeed`, `CustomerNeedConnection`, `CustomerStatus`, `CustomerStatusConnection`, `CustomerTier`, `CustomerTierConnection`
- Reads: `Query.customers`, `Query.customer`, `Query.customerNeeds`, `Query.customerNeed`, `CustomerNeed.projectAttachment`, `Query.customerStatuses`, `Query.customerStatus`, `Query.customerTiers`, `Query.customerTier`
- Writes: `Mutation.customerCreate`, `Mutation.customerUpdate`, `Mutation.customerArchive`, `Mutation.customerNeedCreate`, `Mutation.customerNeedUpdate`, `Mutation.customerNeedArchive`, `Mutation.customerNeedDelete`, `Mutation.customerStatusCreate`, `Mutation.customerStatusUpdate`, `Mutation.customerStatusDelete`, `Mutation.customerTierCreate`, `Mutation.customerTierUpdate`, `Mutation.customerTierDelete`
- Inputs: `CustomerCreateInput`, `CustomerUpdateInput`, `CustomerNeedCreateInput`, `CustomerNeedUpdateInput`, `CustomerStatusCreateInput`, `CustomerStatusUpdateInput`, `CustomerTierCreateInput`, `CustomerTierUpdateInput`
- Relevant fields: `Customer.id`, `Customer.name`, `Customer.domains`, `Customer.externalIds`, `Customer.status`, `Customer.tier`, `Customer.owner`, `Customer.approximateNeedCount`, `Customer.slugId`, `Customer.url`, `CustomerNeed.id`, `CustomerNeed.customer`, `CustomerNeed.issue`, `CustomerNeed.project`, `CustomerNeed.priority`, `CustomerNeed.content`, `CustomerStatus.id`, `CustomerStatus.displayName`, `CustomerStatus.color`, `CustomerStatus.position`, `CustomerTier.id`, `CustomerTier.displayName`, `CustomerTier.color`, `CustomerTier.position`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `customer list` | `Query.customers` | Read-only |
| `customer get` | `Query.customer` | Read-only |
| `customer-need list` | `Query.customerNeeds` | Read-only |
| `customer-need get` | `Query.customerNeed` | Read-only |
| `customer-need project-attachment` | `CustomerNeed.projectAttachment` via `Query.customerNeed` | Read-only, metadata-only projection |
| `customer-status list` | `Query.customerStatuses` | Read-only |
| `customer-status get` | `Query.customerStatus` | Read-only |
| `customer-tier list` | `Query.customerTiers` | Read-only |
| `customer-tier get` | `Query.customerTier` | Read-only |
| `customer create` | `Mutation.customerCreate` | Blocked: customer create needs an explicit organization-scoped safety model |
| `customer update` | `Mutation.customerUpdate` | Blocked: update must resolve and compare the owning organization before mutation |
| `customer archive` | `Mutation.customerArchive` | Blocked: destructive command needs explicit safety semantics |
| `customer-need create` | `Mutation.customerNeedCreate` | Blocked: need creation must prove the linked issue, project, or customer target before mutation |
| `customer-need update` | `Mutation.customerNeedUpdate` | Blocked: update must resolve the need and compare the linked issue or project target before mutation |
| `customer-need archive` | `Mutation.customerNeedArchive` | Blocked: destructive command needs explicit safety semantics |
| `customer-need delete` | `Mutation.customerNeedDelete` | Blocked: destructive command needs explicit safety semantics |
| `customer-status create` | `Mutation.customerStatusCreate` | Blocked: organization lifecycle configuration needs an explicit admin safety model |
| `customer-status update` | `Mutation.customerStatusUpdate` | Blocked: organization lifecycle configuration needs an explicit admin safety model |
| `customer-status delete` | `Mutation.customerStatusDelete` | Blocked: destructive admin command needs explicit safety semantics |
| `customer-tier create` | `Mutation.customerTierCreate` | Blocked: organization tier configuration needs an explicit admin safety model |
| `customer-tier update` | `Mutation.customerTierUpdate` | Blocked: organization tier configuration needs an explicit admin safety model |
| `customer-tier delete` | `Mutation.customerTierDelete` | Blocked: destructive admin command needs explicit safety semantics |

Only Customer read commands are implemented in the current CLI. Customer, CustomerNeed, CustomerStatus, and CustomerTier writes are deferred until they have explicit target or admin safety models.

## Favorite

Use the schema name `Favorite` in code and docs. It is the authenticated user's bookmarked entity in the Linear sidebar.

Schema backing:

- Types: `Favorite`, `FavoriteConnection`
- Reads: `Query.favorites`, `Query.favorite`, `Favorite.children`
- Writes: `Mutation.createFavorite`, `Mutation.updateFavorite`, `Mutation.deleteFavorite`
- Inputs: `FavoriteCreateInput`, `FavoriteUpdateInput`
- Relevant fields: `Favorite.id`, `Favorite.type`, `Favorite.folderName`, `Favorite.url`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `favorite list` | `Query.favorites` | Read-only |
| `favorite children` | `Favorite.children` via `Query.favorite` | Read-only |
| `favorite get` | `Query.favorite` | Read-only |
| `favorite create` | `Mutation.createFavorite` | Blocked: favorite create needs an explicit viewer-scoped safety model |
| `favorite update` | `Mutation.updateFavorite` | Blocked: update must resolve and compare the owning viewer before mutation |
| `favorite delete` | `Mutation.deleteFavorite` | Blocked: destructive command needs explicit safety semantics |

Only `favorite list`, `favorite children`, and `favorite get` are implemented in the current CLI. Favorite writes are deferred as viewer-scoped personalization surface.

## Emoji

Use the schema name `Emoji` in code and docs. It is an organization custom emoji.

Schema backing:

- Types: `Emoji`, `EmojiConnection`
- Reads: `Query.emojis`, `Query.emoji`
- Writes: `Mutation.createEmoji`, `Mutation.deleteEmoji`
- Inputs: `EmojiCreateInput`
- Relevant fields: `Emoji.id`, `Emoji.name`, `Emoji.url`, `Emoji.source`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `emoji list` | `Query.emojis` | Read-only |
| `emoji get` | `Query.emoji` | Read-only |
| `emoji create` | `Mutation.createEmoji` | Blocked: emoji create needs an explicit organization-scoped safety model |
| `emoji delete` | `Mutation.deleteEmoji` | Blocked: destructive command needs explicit safety semantics |

Only `emoji list` and `emoji get` are implemented in the current CLI. Emoji writes are deferred as organization-scoped asset surface.

## File

Use `File` for raw asset upload/download. A file upload produces a workspace asset URL; it is not a target-pinned write because a raw asset has no team or project. The asset URL is attached to an issue or project through the guarded attachment commands.

Schema backing:

- Types: `UploadPayload`, `UploadFile`, `UploadFileHeader`
- Writes: `Mutation.fileUpload` (prepares a pre-signed upload target; the bytes are then PUT to storage out of band)
- Relevant fields: `UploadFile.uploadUrl`, `UploadFile.assetUrl`, `UploadFile.headers`, `UploadFile.contentType`, `UploadFile.size`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `files upload` | `Mutation.fileUpload` then an HTTP PUT of the bytes to the pre-signed URL | Workspace asset, not target-pinned; prints the asset URL for a later guarded attachment write |
| `files download` | Plain HTTP GET of the asset URL to a local path | Read-only, no API; no auth header is attached so a user-supplied URL never receives the Linear token |

`files upload PATH` infers the content type from the file extension (overridable with `--content-type`), calls `fileUpload` for a pre-signed target, PUTs the bytes with the returned headers, and prints `UploadFile.assetUrl`. `files download URL --output PATH` performs an unauthenticated GET and writes the body to the path; it is meant for public or signed asset URLs.

## Attachment

Use the schema name `Attachment` in code and docs. It is an external resource linked to a Linear issue.

Schema backing:

- Types: `Attachment`, `AttachmentConnection`
- Reads: `Query.attachments`, `Query.attachment`, `Query.attachmentsForURL`, `Query.attachmentIssue`, `Issue.attachments`, `Issue.botActor`, `Issue.children`, `Issue.documents`, `Issue.formerAttachments`, `Issue.history`, `Issue.inverseRelations`, `Issue.labels`, `Issue.relations`, `Issue.releases`, `Issue.stateHistory`, `Issue.subscribers` via `Query.attachmentIssue`
- Writes: `Mutation.attachmentCreate`, `Mutation.attachmentUpdate`, `Mutation.attachmentDelete`, `Mutation.attachmentLinkURL`
- Inputs: `AttachmentCreateInput`, `AttachmentUpdateInput`
- Relevant fields: `Attachment.id`, `Attachment.title`, `Attachment.subtitle`, `Attachment.url`, `Attachment.sourceType`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `attachment list` | `Query.attachments` | Read-only |
| `attachment url` | `Query.attachmentsForURL` | Read-only |
| `attachment get` | `Query.attachment` | Read-only |
| `attachment issue get` | `Query.attachmentIssue` | Read-only |
| `attachment issue attachments` | `Issue.attachments` via `Query.attachmentIssue` | Read-only |
| `attachment issue bot-actor` | `Issue.botActor` via `Query.attachmentIssue` | Read-only, bot metadata only |
| `attachment issue children` | `Issue.children` via `Query.attachmentIssue` | Read-only |
| `attachment issue comments` | `Issue.comments` via `Query.attachmentIssue`; returns comment metadata without body | Read-only |
| `attachment issue documents` | `Issue.documents` via `Query.attachmentIssue` | Read-only |
| `attachment issue former-attachments` | `Issue.formerAttachments` via `Query.attachmentIssue` | Read-only |
| `attachment issue former-needs` | `Issue.formerNeeds` via `Query.attachmentIssue`; returns customer-need metadata without body/content | Read-only |
| `attachment issue history` | `Issue.history` via `Query.attachmentIssue` | Read-only, compact metadata only |
| `attachment issue inverse-relations` | `Issue.inverseRelations` via `Query.attachmentIssue` | Read-only |
| `attachment issue labels` | `Issue.labels` via `Query.attachmentIssue` | Read-only |
| `attachment issue needs` | `Issue.needs` via `Query.attachmentIssue`; returns customer-need metadata without body/content | Read-only |
| `attachment issue relations` | `Issue.relations` via `Query.attachmentIssue` | Read-only |
| `attachment issue releases` | `Issue.releases` via `Query.attachmentIssue` | Read-only |
| `attachment issue shared-access` | `Issue.sharedAccess` via `Query.attachmentIssue`; omits shared user details and exposes only flags/counts/disallowed fields | Read-only |
| `attachment issue state-history` | `Issue.stateHistory` via `Query.attachmentIssue` | Read-only, workflow-state span metadata |
| `attachment issue subscribers` | `Issue.subscribers` via `Query.attachmentIssue` | Read-only |
| `attachment create` | `Mutation.attachmentCreate` | Blocked: attachment create must resolve and compare the owning issue's team before mutation |
| `attachment update` | `Mutation.attachmentUpdate` | Blocked: update must resolve and compare the owning issue before mutation |
| `attachment delete` | `Mutation.attachmentDelete` | Blocked: destructive command needs explicit safety semantics |

Only read-only attachment commands are implemented in the current CLI. Attachment writes are deferred until the owning-issue guard model is explicit.
