# linctl domain map

This map is derived from the vendored Linear schema at `internal/client/schema.graphql`.
Command names below are either implemented CLI surface or intentionally deferred surface. Implementation slices must use GraphQL operations backed by these schema fields.

## Core target

| CLI surface | Schema backing | Notes |
| --- | --- | --- |
| `whoami` | `Query.viewer`, `User` | Reads the authenticated user. |
| `target` | `Query.organization`, `Query.teams`, `Query.team`, `Query.projects`, `Query.project` | Resolves the active token's organization, team, and optional project. |
| `doctor` | `Query.viewer`, `Query.teams`, optional `Query.project` | Read-only health check for config load, token presence, and pinned-target confirmation. Does not print token values. |
| `organization exists` | `Query.organizationExists` | Read-only URL-key existence check for workspace lookup. |
| `rate-limit status` | `Query.rateLimitStatus` | Read-only quota status for the authenticated Linear client. |

The target vocabulary is `org_id`, `team_key`, `team_id`, and optional `project_id`. Do not introduce `workspace` as a flag or JSON key synonym.

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

- Types: `Release`, `ReleasePipeline`, `ReleaseStage`, `ReleaseNote`, `IssueToRelease`
- Reads: `Query.releasePipelines`, `Query.releasePipeline`, `Query.releaseStages`, `Query.releaseStage`, `Query.releases`, `Query.release`, `Query.releaseSearch`, `Query.releaseNotes`, `Query.releaseNote`
- Deferred reads: nested release documents/issues/history/links and access-key release reads
- Writes: `Mutation.releasePipelineCreate`, `Mutation.releasePipelineUpdate`, `Mutation.releasePipelineArchive`, `Mutation.releasePipelineDelete`, `Mutation.releaseStageCreate`, `Mutation.releaseStageUpdate`, `Mutation.releaseStageArchive`, `Mutation.releaseStageUnarchive`, plus Release/ReleaseNote/IssueToRelease create/update/archive/delete/sync/complete mutations
- Relevant fields: `Release.id`, `Release.name`, `Release.slugId`, `Release.version`, `Release.pipeline`, `Release.stage`, `Release.issueCount`, `ReleaseNote.id`, `ReleaseNote.title`, `ReleaseNote.slugId`, `ReleaseNote.pipeline`, `ReleaseNote.releaseCount`, `ReleasePipeline.id`, `ReleasePipeline.name`, `ReleasePipeline.slugId`, `ReleasePipeline.type`, `ReleasePipeline.isProduction`, `ReleasePipeline.approximateReleaseCount`, `ReleaseStage.id`, `ReleaseStage.name`, `ReleaseStage.type`, `ReleaseStage.pipeline`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `release-pipeline list` | `Query.releasePipelines` | Read-only |
| `release-pipeline get` | `Query.releasePipeline` | Read-only |
| `release-stage list` | `Query.releaseStages` | Read-only |
| `release-stage get` | `Query.releaseStage` | Read-only |
| `release list` | `Query.releases` | Read-only |
| `release search` | `Query.releaseSearch` | Read-only |
| `release get` | `Query.release` | Read-only |
| `release-note list` | `Query.releaseNotes` | Read-only |
| `release-note get` | `Query.releaseNote` | Read-only |
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

Release, ReleasePipeline, ReleaseStage, and ReleaseNote read commands are implemented in the current CLI. IssueToRelease, sync, complete, access-key, and association commands remain deferred until their control-surface shape and guard model are explicit.

## Issue

Schema backing:

- Types: `Issue`, `IssueConnection`
- Reads: `Query.issues`, `Query.issue`
- Writes: `Mutation.issueCreate`, `Mutation.issueUpdate`, `Mutation.issueArchive`, `Mutation.commentCreate`
- Inputs: `IssueCreateInput`, `IssueUpdateInput`
- Relevant fields: `Issue.id`, `Issue.identifier`, `Issue.number`, `Issue.title`, `Issue.team`, `Issue.cycle`, `Issue.project`, `Issue.projectMilestone`, `Issue.assignee`, `Issue.state`, `Issue.documents`, `Issue.comments`, `Issue.url`, `Issue.branchName`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `issue list` | `Query.issues`, optionally filtered by `Issue.team.id`, `Issue.state.type`, `Issue.project.id`, `Issue.assignee.id`, `Issue.labels.some.id`, `Issue.cycle.id`, `Issue.createdAt.gte` (`--created-after` / `--created-since`), `Issue.createdAt.lte`, `Issue.hasBlockedByRelations.eq`, or `Issue.hasBlockingRelations.eq`; `--blocked-by ISSUE` traverses `Issue.relations` with `IssueRelation.type == "blocks"` and returns matching `IssueRelation.relatedIssue`; `--all-teams` omits the team filter | Read-only |
| `issue search` | `Query.issues`, filtered by `Issue.searchableContent` | Read-only |
| `issue get` | `Query.issue` | Read-only |
| `issue deps` | `Query.issue`, `Issue.parent`, `Issue.children`, `Issue.relations`, `Issue.inverseRelations`; `IssueRelation.type == "blocks"` separates blocked issues from blockers | Read-only |
| `issue id` | Current checkout issue identifier from git/jj context | Read-only |
| `issue title` | `Query.issue` after current checkout or explicit issue resolution | Read-only |
| `issue url` | `Query.issue` after current checkout or explicit issue resolution | Read-only |
| `issue branch` | `Query.issue`, `Issue.branchName` | Read-only |
| `issue pr` | `Query.issue`; emits a local `gh pr create` title/body plan without calling GitHub | Read-only |
| `next --dry-run` | `Query.issues`, filtered by `Issue.team.id`, `Issue.state.type == "unstarted"`, and `Issue.hasBlockedByRelations.eq == false`; fetches `Issue.relations`, `Issue.priority`, and `Issue.createdAt`, then ranks by active unblock count, priority, and age before printing one candidate without checkout/worktree creation | Read-only |
| `done` | Current checkout issue identifier, then `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved |
| `issue create` | `Mutation.issueCreate` with `IssueCreateInput.teamId`, optional `projectId`; `--description-file` is resolved locally before mutation | Team-scoped unless `projectId` is set |
| `issue update` | `Mutation.issueUpdate` with `IssueUpdateInput`; `--description-file` replaces description, while `--append` or `--append-file` first reads `Issue.description` and appends text before sending `description` | Resource-scoped when a project target is involved |
| `issue start` | `Query.viewer`, `Query.workflowStates` filtered to `started`, then `Mutation.issueUpdate` with `IssueUpdateInput.assigneeId` and `stateId` | Resource-scoped when a project target is involved |
| `issue comment` | `Mutation.commentCreate`; `--body -` reads stdin and `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project |
| `issue reply` | `Mutation.commentCreate` with `CommentCreateInput.parentId`; `--body-file` reads a local file before mutation | Resource-scoped to the issue's resolved team/project |
| `issue close` | `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved |
| `issue comments` | `Issue.comments` via `Query.issue` | Read-only |

## Comment

Schema backing:

- Types: `Comment`, `CommentConnection`
- Reads: `Query.comments`, `Query.comment`, `Issue.comments`
- Writes: `Mutation.commentCreate`, `Mutation.commentResolve`, `Mutation.commentUnresolve`
- Inputs: `CommentCreateInput`
- Relevant fields: `Comment.id`, `Comment.body`, `Comment.url`, `Comment.createdAt`, `Comment.updatedAt`, `Comment.parentId`, `Comment.issueId`, `Comment.projectId`, `Comment.projectUpdateId`, `Comment.initiativeId`, `Comment.initiativeUpdateId`, `Comment.documentContentId`, `Comment.user`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `comment list` | `Query.comments` | Read-only |
| `comment get` | `Query.comment` | Read-only |
| `comment resolve` | `Mutation.commentResolve` | Blocked: resolving must first identify and compare the parent issue/project/update/document scope |
| `comment unresolve` | `Mutation.commentUnresolve` | Blocked: unresolving must first identify and compare the parent issue/project/update/document scope |

Only `comment list` and `comment get` are implemented in the current CLI. Issue-scoped comment creation and replies remain under the guarded `issue comment` and `issue reply` commands.

## Project

Schema backing:

- Types: `Project`, `ProjectConnection`
- Reads: `Query.projects`, `Query.project`
- Writes: `Mutation.projectCreate`, `Mutation.projectUpdate`, `Mutation.projectArchive`
- Inputs: `ProjectCreateInput`, `ProjectUpdateInput`
- Relevant fields: `Project.id`, `Project.name`, `Project.description`, `Project.status`, `Project.lead`, `Project.url`, `Project.teams`, `Project.members`, `Project.documents`, `Project.projectMilestones`, `Project.issues`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project list` | `Query.projects` | Read-only |
| `project get` | `Query.project` | Read-only |
| `project create` | `Mutation.projectCreate` with `ProjectCreateInput.teamIds` | Team-scoped |
| `project update` | `Mutation.projectUpdate` with `ProjectUpdateInput` | Resource-scoped, compare `project_id` |
| `project archive` | `Mutation.projectArchive` | Resource-scoped, compare `project_id` |
| `project members` | `Project.members` plus `Mutation.projectUpdate` with `ProjectUpdateInput.memberIds` | Read-only for list, resource-scoped for writes |
| `project updates` | `Project.projectUpdates` | Read-only |

Project is the first implemented PM domain; later domains should reuse its target-comparison vocabulary.

## ProjectUpdate

Use `ProjectUpdate` for Linear project status updates. Avoid calling these generic comments or notes.

Schema backing:

- Types: `ProjectUpdate`, `ProjectUpdateConnection`
- Reads: `Query.projectUpdates`, `Query.projectUpdate`, `Project.projectUpdates`
- Writes: `Mutation.projectUpdateCreate`, `Mutation.projectUpdateUpdate`, `Mutation.projectUpdateArchive`
- Relevant fields: `ProjectUpdate.id`, `ProjectUpdate.body`, `ProjectUpdate.health`, `ProjectUpdate.createdAt`, `ProjectUpdate.updatedAt`, `ProjectUpdate.url`, `ProjectUpdate.project`, `ProjectUpdate.user`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-update list` | `Query.projectUpdates` | Read-only |
| `project-update get` | `Query.projectUpdate` | Read-only |
| `project-update create` | `Mutation.projectUpdateCreate` | Blocked: create must resolve and compare the target project before posting |
| `project-update update` | `Mutation.projectUpdateUpdate` | Blocked: update must resolve and compare the owning project before mutation |
| `project-update archive` | `Mutation.projectUpdateArchive` | Blocked: destructive command needs explicit safety semantics |

Only `project-update list` and `project-update get` are implemented in the current top-level CLI. `project updates PROJECT_ID` remains the project-scoped history view.

## Cycle

Schema backing:

- Types: `Cycle`, `CycleConnection`
- Reads: `Query.cycles`, `Query.cycle`, `Team.cycles`
- Writes: `Mutation.cycleCreate`, `Mutation.cycleUpdate`, `Mutation.cycleArchive`, `Mutation.cycleShiftAll`, `Mutation.cycleStartUpcomingCycleToday`
- Relevant fields: `Cycle.id`, `Cycle.number`, `Cycle.name`, `Cycle.startsAt`, `Cycle.endsAt`, `Cycle.completedAt`, `Cycle.team`, `Cycle.issues`, `Cycle.progress`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `cycle list` | `Query.cycles` | Read-only |
| `cycle get` | `Query.cycle` | Read-only |
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
- Reads: `Query.projectMilestones`, `Query.projectMilestone`, `Project.projectMilestones`
- Writes: `Mutation.projectMilestoneCreate`, `Mutation.projectMilestoneUpdate`, `Mutation.projectMilestoneDelete`
- Inputs: `ProjectMilestoneCreateInput`, `ProjectMilestoneUpdateInput`
- Relevant fields: `ProjectMilestone.id`, `ProjectMilestone.name`, `ProjectMilestone.description`, `ProjectMilestone.targetDate`, `ProjectMilestone.status`, `ProjectMilestone.project`, `ProjectMilestone.sortOrder`, `ProjectMilestone.issues`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-milestone list` | `Project.projectMilestones` via `Query.project` | Read-only |
| `project-milestone get` | `Query.projectMilestone` | Read-only |
| `project-milestone create` | `Mutation.projectMilestoneCreate` with `projectId` | Resource-scoped, compare `project_id` |
| `project-milestone update` | `Mutation.projectMilestoneUpdate` | Resource-scoped, compare resolved project |
| `project-milestone delete` | `Mutation.projectMilestoneDelete` | Resource-scoped, compare resolved project |

## Document

Schema backing:

- Types: `Document`, `DocumentConnection`
- Reads: `Query.documents`, `Query.document`, `Project.documents`, `Team.documents`, `Issue.documents`, `Cycle.documents`
- Writes: `Mutation.documentCreate`, `Mutation.documentUpdate`, `Mutation.documentDelete`
- Inputs: `DocumentCreateInput`, `DocumentUpdateInput`
- Relevant fields: `Document.id`, `Document.title`, `Document.slugId`, `Document.archivedAt`, `Document.project`, `Document.team`, `Document.issue`, `Document.cycle`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `document list` | `Query.documents` | Read-only |
| `document get` | `Query.document` | Read-only |
| `document create` | `Mutation.documentCreate` with optional `projectId`, `teamId`, `issueId`, `cycleId` | Blocked: parent can be project, team, issue, or cycle; write guard needs explicit parent-resolution semantics |
| `document update` | `Mutation.documentUpdate` | Blocked: update must resolve and compare the existing parent before changing content |
| `document delete` | `Mutation.documentDelete` | Blocked: destructive command needs explicit safety semantics |

Only `document list` and `document get` are implemented in the current CLI. Document writes are deferred until the parent-resolution guard is designed.

## Label

CLI name `label` maps to Linear schema type `IssueLabel`.

Schema backing:

- Types: `IssueLabel`, `IssueLabelConnection`
- Reads: `Query.issueLabels`, `Query.issueLabel`, `Team.labels`
- Writes: `Mutation.issueLabelCreate`, `Mutation.issueLabelUpdate`, `Mutation.issueLabelDelete`
- Inputs: `IssueLabelCreateInput`, `IssueLabelUpdateInput`
- Relevant fields: `IssueLabel.id`, `IssueLabel.name`, `IssueLabel.description`, `IssueLabel.color`, `IssueLabel.team`, `IssueLabel.issues`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `label list` | `Query.issueLabels` | Read-only |
| `label get` | `Query.issueLabel` | Read-only |
| `label create` | `Mutation.issueLabelCreate` with optional `teamId` | Blocked: optional team scope needs explicit org/team target behavior before writes |
| `label update` | `Mutation.issueLabelUpdate` | Blocked: update must resolve and compare the label's owning team before mutation |
| `label delete` | `Mutation.issueLabelDelete` | Blocked: destructive command needs explicit safety semantics |

Only `label list` and `label get` are implemented in the current CLI. Label writes are deferred until the team-scope guard is designed.

## Team

Schema backing:

- Types: `Team`, `TeamConnection`, `TeamMembership`
- Reads: `Query.teams`, `Query.team`, `Team.members`, `Team.issues`, `Team.cycles`, `Team.projects`
- Writes: `Mutation.teamCreate`, `Mutation.teamUpdate`, `Mutation.teamDelete`, `Mutation.teamMembershipCreate`, `Mutation.teamMembershipUpdate`, `Mutation.teamMembershipDelete`
- Inputs: `TeamCreateInput`, `TeamUpdateInput`
- Relevant fields: `Team.id`, `Team.name`, `Team.key`, `Team.description`, `Team.archivedAt`, `Team.issues`, `Team.cycles`, `Team.members`, `Team.projects`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `team list` | `Query.teams` | Read-only |
| `team get` | `Query.team` | Read-only |
| `team create` | `Mutation.teamCreate` | Blocked: organization administration surface needs an explicit admin safety model |
| `team update` | `Mutation.teamUpdate` | Blocked: team metadata writes need stronger authority checks than ordinary target comparison |
| `team delete` | `Mutation.teamDelete` | Blocked: destructive command needs explicit safety semantics |
| `team members` | `Team.members` | Read-only |

Only `team list`, `team get`, and `team members` are implemented in the current CLI. Team creation, metadata mutation, and membership writes are deferred as organization/admin surface.

## User

Schema backing:

- Types: `User`, `UserConnection`
- Reads: `Query.users`, `Query.user`, `Query.viewer`, `Team.members`, `Project.members`
- Relevant fields: `User.id`, `User.name`, `User.displayName`, `User.email`, `User.active`, `User.guest`, `User.admin`, `User.url`, `User.assignedIssues`, `User.teams`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `user list` | `Query.users` | Read-only |
| `user get` | `Query.user` | Read-only |
| `user me` | `Query.viewer` | Read-only |

User writes are not part of the v1 PM command surface until a later slice proves the exact Linear mutation and safety semantics.

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
| `workflow-state create` | `Mutation.workflowStateCreate` | Blocked: team workflow configuration needs an explicit admin safety model |
| `workflow-state update` | `Mutation.workflowStateUpdate` | Blocked: update must resolve and compare the owning team before mutation |
| `workflow-state archive` | `Mutation.workflowStateArchive` | Blocked: destructive command needs explicit safety semantics |

Only `workflow-state list` and `workflow-state get` are implemented in the current CLI. WorkflowState writes are deferred as team/admin configuration surface.

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

## Initiative

Use the schema name `Initiative` in code and docs. It is Linear's strategic grouping of projects toward a goal.

Schema backing:

- Types: `Initiative`, `InitiativeConnection`
- Reads: `Query.initiatives`, `Query.initiative`
- Writes: `Mutation.createInitiative`, `Mutation.updateInitiative`, `Mutation.archiveInitiative`, `Mutation.deleteInitiative`
- Inputs: `InitiativeCreateInput`, `InitiativeUpdateInput`
- Relevant fields: `Initiative.id`, `Initiative.name`, `Initiative.description`, `Initiative.status`, `Initiative.priority`, `Initiative.targetDate`, `Initiative.slugId`, `Initiative.url`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `initiative list` | `Query.initiatives` | Read-only |
| `initiative get` | `Query.initiative` | Read-only |
| `initiative create` | `Mutation.createInitiative` | Blocked: initiative create needs an explicit organization-scoped safety model |
| `initiative update` | `Mutation.updateInitiative` | Blocked: update must resolve and compare the owning organization before mutation |
| `initiative archive` | `Mutation.archiveInitiative` | Blocked: destructive command needs explicit safety semantics |

Only `initiative list` and `initiative get` are implemented in the current CLI. Initiative writes are deferred as organization-scoped planning surface.

## InitiativeUpdate

Use `InitiativeUpdate` for Linear initiative status updates. Avoid calling these generic comments or notes.

Schema backing:

- Types: `InitiativeUpdate`, `InitiativeUpdateConnection`
- Reads: `Query.initiativeUpdates`, `Query.initiativeUpdate`
- Writes: `Mutation.initiativeUpdateCreate`, `Mutation.initiativeUpdateUpdate`, `Mutation.initiativeUpdateArchive`, `Mutation.initiativeUpdateUnarchive`
- Inputs: `InitiativeUpdateCreateInput`, `InitiativeUpdateUpdateInput`
- Relevant fields: `InitiativeUpdate.id`, `InitiativeUpdate.body`, `InitiativeUpdate.health`, `InitiativeUpdate.createdAt`, `InitiativeUpdate.updatedAt`, `InitiativeUpdate.url`, `InitiativeUpdate.slugId`, `InitiativeUpdate.commentCount`, `InitiativeUpdate.initiative`, `InitiativeUpdate.user`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `initiative-update list` | `Query.initiativeUpdates` | Read-only |
| `initiative-update get` | `Query.initiativeUpdate` | Read-only |
| `initiative-update create` | `Mutation.initiativeUpdateCreate` | Blocked: create must resolve and compare the owning Initiative before posting |
| `initiative-update update` | `Mutation.initiativeUpdateUpdate` | Blocked: update must resolve and compare the owning Initiative before mutation |
| `initiative-update archive` | `Mutation.initiativeUpdateArchive` | Blocked: destructive command needs explicit safety semantics |
| `initiative-update unarchive` | `Mutation.initiativeUpdateUnarchive` | Blocked: unarchive needs explicit lifecycle and target semantics |

Only `initiative-update list` and `initiative-update get` are implemented in the current CLI. InitiativeUpdate writes and reminders are deferred until their guard model is explicit.

## Roadmap

Use the schema name `Roadmap` in code and docs. It is Linear's deprecated roadmap grouping for projects; prefer `Initiative` for new planning workflows.

Schema backing:

- Types: `Roadmap`, `RoadmapConnection`
- Reads: `Query.roadmaps`, `Query.roadmap`
- Writes: `Mutation.roadmapCreate`, `Mutation.roadmapUpdate`, `Mutation.roadmapArchive`, `Mutation.roadmapDelete`
- Inputs: `RoadmapCreateInput`, `RoadmapUpdateInput`
- Relevant fields: `Roadmap.id`, `Roadmap.name`, `Roadmap.description`, `Roadmap.color`, `Roadmap.slugId`, `Roadmap.sortOrder`, `Roadmap.url`, `Roadmap.creator`, `Roadmap.owner`

Command status:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `roadmap list` | `Query.roadmaps` | Read-only |
| `roadmap get` | `Query.roadmap` | Read-only |
| `roadmap create` | `Mutation.roadmapCreate` | Blocked: deprecated organization-scoped planning surface needs an explicit safety model |
| `roadmap update` | `Mutation.roadmapUpdate` | Blocked: update must resolve and compare the owning organization before mutation |
| `roadmap archive` | `Mutation.roadmapArchive` | Blocked: destructive command needs explicit safety semantics |
| `roadmap delete` | `Mutation.roadmapDelete` | Blocked: destructive command needs explicit safety semantics |

Only `roadmap list` and `roadmap get` are implemented in the current CLI. Roadmap writes and roadmap-project associations are deferred; prefer Initiative commands for current Linear planning workflows.

## CustomView

Use the schema name `CustomView` in code and docs. It is Linear's saved view over issues, projects, or initiatives.

Schema backing:

- Types: `CustomView`, `CustomViewConnection`
- Reads: `Query.customViews`, `Query.customView`, `Query.customViewHasSubscribers`
- Writes: `Mutation.createCustomView`, `Mutation.updateCustomView`, `Mutation.deleteCustomView`
- Inputs: `CustomViewCreateInput`, `CustomViewUpdateInput`
- Relevant fields: `CustomView.id`, `CustomView.name`, `CustomView.description`, `CustomView.modelName`, `CustomView.shared`, `CustomView.color`, `CustomView.slugId`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `custom-view list` | `Query.customViews` | Read-only |
| `custom-view subscribers` | `Query.customViewHasSubscribers` | Read-only |
| `custom-view get` | `Query.customView` | Read-only |
| `custom-view create` | `Mutation.createCustomView` | Blocked: custom view create needs an explicit organization-scoped safety model |
| `custom-view update` | `Mutation.updateCustomView` | Blocked: update must resolve and compare the owning organization before mutation |
| `custom-view delete` | `Mutation.deleteCustomView` | Blocked: destructive command needs explicit safety semantics |

Only `custom-view list`, `custom-view subscribers`, and `custom-view get` are implemented in the current CLI. CustomView writes are deferred as organization-scoped view configuration surface.

## Customer

Use the schema name `Customer` in code and docs. It is Linear's customer organization record for customer requests and feedback.

Schema backing:

- Types: `Customer`, `CustomerConnection`, `CustomerNeed`, `CustomerNeedConnection`, `CustomerStatus`, `CustomerStatusConnection`, `CustomerTier`, `CustomerTierConnection`
- Reads: `Query.customers`, `Query.customer`, `Query.customerNeeds`, `Query.customerNeed`, `Query.customerStatuses`, `Query.customerStatus`, `Query.customerTiers`, `Query.customerTier`
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
| `customer-status create` | `Mutation.customerStatusCreate` | Blocked: workspace lifecycle configuration needs an explicit admin safety model |
| `customer-status update` | `Mutation.customerStatusUpdate` | Blocked: workspace lifecycle configuration needs an explicit admin safety model |
| `customer-status delete` | `Mutation.customerStatusDelete` | Blocked: destructive admin command needs explicit safety semantics |
| `customer-tier create` | `Mutation.customerTierCreate` | Blocked: workspace tier configuration needs an explicit admin safety model |
| `customer-tier update` | `Mutation.customerTierUpdate` | Blocked: workspace tier configuration needs an explicit admin safety model |
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

Use the schema name `Emoji` in code and docs. It is a workspace custom emoji.

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

## Attachment

Use the schema name `Attachment` in code and docs. It is an external resource linked to a Linear issue.

Schema backing:

- Types: `Attachment`, `AttachmentConnection`
- Reads: `Query.attachments`, `Query.attachment`, `Query.attachmentsForURL`
- Writes: `Mutation.attachmentCreate`, `Mutation.attachmentUpdate`, `Mutation.attachmentDelete`, `Mutation.attachmentLinkURL`
- Inputs: `AttachmentCreateInput`, `AttachmentUpdateInput`
- Relevant fields: `Attachment.id`, `Attachment.title`, `Attachment.subtitle`, `Attachment.url`, `Attachment.sourceType`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `attachment list` | `Query.attachments` | Read-only |
| `attachment url` | `Query.attachmentsForURL` | Read-only |
| `attachment get` | `Query.attachment` | Read-only |
| `attachment create` | `Mutation.attachmentCreate` | Blocked: attachment create must resolve and compare the owning issue's team before mutation |
| `attachment update` | `Mutation.attachmentUpdate` | Blocked: update must resolve and compare the owning issue before mutation |
| `attachment delete` | `Mutation.attachmentDelete` | Blocked: destructive command needs explicit safety semantics |

Only `attachment list`, `attachment url`, and `attachment get` are implemented in the current CLI. Attachment writes are deferred until the owning-issue guard model is explicit.
