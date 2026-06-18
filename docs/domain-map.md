# linctl domain map

This map is derived from the vendored Linear schema at `internal/client/schema.graphql`.
Command names below are planned names only; implementation slices must use GraphQL operations backed by these schema fields.

## Core target

| CLI surface | Schema backing | Notes |
| --- | --- | --- |
| `whoami` | `Query.viewer`, `User` | Reads the authenticated user. |
| `target` | `Query.organization`, `Query.teams`, `Query.team`, `Query.projects`, `Query.project` | Resolves the active token's organization, team, and optional project. |

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
| `issue list` | `Query.issues` | Read-only |
| `issue get` | `Query.issue` | Read-only |
| `issue create` | `Mutation.issueCreate` with `IssueCreateInput.teamId`, optional `projectId` | Team-scoped unless `projectId` is set |
| `issue update` | `Mutation.issueUpdate` with `IssueUpdateInput` | Resource-scoped when a project target is involved |
| `issue comment` | `Mutation.commentCreate` | Resource-scoped to the issue's resolved team/project |
| `issue close` | `Mutation.issueUpdate` state change | Resource-scoped when a project target is involved |

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

S6a implements this domain first as the template for later domains.

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
- Relevant fields: `ProjectMilestone.id`, `ProjectMilestone.name`, `ProjectMilestone.description`, `ProjectMilestone.targetDate`, `ProjectMilestone.project`, `ProjectMilestone.sortOrder`, `ProjectMilestone.issues`

Planned commands:

| Command | Operation backing | Write scope |
| --- | --- | --- |
| `project-milestone list` | `Query.projectMilestones` | Read-only |
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
| `document create` | `Mutation.documentCreate` with optional `projectId`, `teamId`, `issueId`, `cycleId` | Team-scoped or resource-scoped by target fields |
| `document update` | `Mutation.documentUpdate` | Resource-scoped by resolved parent |
| `document delete` | `Mutation.documentDelete` | Resource-scoped by resolved parent |

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
| `label create` | `Mutation.issueLabelCreate` with optional `teamId` | Team-scoped |
| `label update` | `Mutation.issueLabelUpdate` | Team-scoped |
| `label delete` | `Mutation.issueLabelDelete` | Team-scoped |

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
| `team create` | `Mutation.teamCreate` | Org-scoped, compare `org_id` |
| `team update` | `Mutation.teamUpdate` | Team-scoped |
| `team delete` | `Mutation.teamDelete` | Team-scoped |
| `team members` | `Team.members` plus team membership mutations | Read-only for list, team-scoped for writes |

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
