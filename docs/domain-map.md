# linctl domain map

This map is derived from the vendored Linear schema at `internal/client/schema.graphql`.
Command names below are either implemented CLI surface or intentionally deferred surface. Implementation slices must use GraphQL operations backed by these schema fields.

## Core target

| CLI surface | Schema backing | Notes |
| --- | --- | --- |
| `whoami` | `Query.viewer`, `User` | Reads the authenticated user. |
| `target` | `Query.organization`, `Query.teams`, `Query.team`, `Query.projects`, `Query.project` | Resolves the active token's organization, team, and optional project. |
| `doctor` | `Query.viewer`, `Query.teams`, optional `Query.project` | Read-only health check for config load, token presence, and pinned-target confirmation. Does not print token values. |

The target vocabulary is `org_id`, `team_key`, `team_id`, and optional `project_id`. Do not introduce `workspace` as a flag or JSON key synonym.

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

## CustomView

Use the schema name `CustomView` in code and docs. It is Linear's saved view over issues, projects, or initiatives.

Schema backing:

- Types: `CustomView`, `CustomViewConnection`
- Reads: `Query.customViews`, `Query.customView`
- Writes: `Mutation.createCustomView`, `Mutation.updateCustomView`, `Mutation.deleteCustomView`
- Inputs: `CustomViewCreateInput`, `CustomViewUpdateInput`
- Relevant fields: `CustomView.id`, `CustomView.name`, `CustomView.description`, `CustomView.modelName`, `CustomView.shared`, `CustomView.color`, `CustomView.slugId`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `custom-view list` | `Query.customViews` | Read-only |
| `custom-view get` | `Query.customView` | Read-only |
| `custom-view create` | `Mutation.createCustomView` | Blocked: custom view create needs an explicit organization-scoped safety model |
| `custom-view update` | `Mutation.updateCustomView` | Blocked: update must resolve and compare the owning organization before mutation |
| `custom-view delete` | `Mutation.deleteCustomView` | Blocked: destructive command needs explicit safety semantics |

Only `custom-view list` and `custom-view get` are implemented in the current CLI. CustomView writes are deferred as organization-scoped view configuration surface.

## Favorite

Use the schema name `Favorite` in code and docs. It is the authenticated user's bookmarked entity in the Linear sidebar.

Schema backing:

- Types: `Favorite`, `FavoriteConnection`
- Reads: `Query.favorites`, `Query.favorite`
- Writes: `Mutation.createFavorite`, `Mutation.updateFavorite`, `Mutation.deleteFavorite`
- Inputs: `FavoriteCreateInput`, `FavoriteUpdateInput`
- Relevant fields: `Favorite.id`, `Favorite.type`, `Favorite.folderName`, `Favorite.url`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `favorite list` | `Query.favorites` | Read-only |
| `favorite get` | `Query.favorite` | Read-only |
| `favorite create` | `Mutation.createFavorite` | Blocked: favorite create needs an explicit viewer-scoped safety model |
| `favorite update` | `Mutation.updateFavorite` | Blocked: update must resolve and compare the owning viewer before mutation |
| `favorite delete` | `Mutation.deleteFavorite` | Blocked: destructive command needs explicit safety semantics |

Only `favorite list` and `favorite get` are implemented in the current CLI. Favorite writes are deferred as viewer-scoped personalization surface.

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
