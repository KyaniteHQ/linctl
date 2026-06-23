package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientReadFailureScenarios_wrap_graphql_errors(t *testing.T) {
	graphqlClient := errorGraphQLClient{err: errors.New("network down")}

	_, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssues(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{StateType: "started"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{ProjectID: "project-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{AssigneeID: "user-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{LabelID: "label-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{CycleID: "cycle-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{CreatedAfter: "2026-06-01"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{CreatedBefore: "2026-06-30"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{HasBlockers: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{Blocks: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list next issues")

	_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{BlockedBy: "LIT-1"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issues")

	_, err = SearchIssuesByTeam(context.Background(), graphqlClient, "team-id", "needle", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "search issues")

	_, err = SearchIssuesByFigmaFileKey(context.Background(), graphqlClient, "figma-key", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "search issues by Figma file key")

	_, err = ListIssuePriorityValues(context.Background(), graphqlClient)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue priority values")

	_, err = GetIssueFilterSuggestion(context.Background(), graphqlClient, "started issues", "team-id", "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue filter suggestion")

	_, err = GetIssueTitleSuggestionFromCustomerRequest(
		context.Background(),
		graphqlClient,
		"Customer asks for faster exports",
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue title suggestion")

	_, err = GetIssueByID(context.Background(), graphqlClient, "LIT-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue LIT-1")

	_, err = GetIssueDependencies(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue dependencies LIT-1")

	_, err = ListIssueAttachments(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue attachments LIT-1")

	_, err = GetIssueBotActor(context.Background(), graphqlClient, "LIT-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue bot actor LIT-1")

	_, err = ListIssueChildren(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue children LIT-1")

	_, err = ListIssueDocuments(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue documents LIT-1")

	_, err = ListIssueFormerAttachments(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue former attachments LIT-1")

	_, err = ListIssueHistory(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue history LIT-1")

	_, err = ListIssueInverseRelations(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue inverse relations LIT-1")

	_, err = ListIssueLabels(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue labels LIT-1")

	_, err = ListIssueRelationsForIssue(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue relations LIT-1")

	_, err = ListIssueReleases(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue releases LIT-1")

	_, err = ListIssueStateHistory(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue state history LIT-1")

	_, err = ListIssueSubscribers(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue subscribers LIT-1")

	_, err = ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list projects")

	_, err = GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project project-id")

	_, err = ListProjectMembers(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project members project-id")

	_, err = ListProjectAttachments(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project attachments project-id")

	_, err = ListProjectDocuments(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project documents project-id")

	_, err = ListProjectExternalLinks(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project external links project-id")

	_, err = ListProjectHistory(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project history project-id")

	_, err = ListProjectInitiativeToProjects(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project initiative associations project-id")

	_, err = ListProjectInitiatives(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project initiatives project-id")

	_, err = ListProjectInverseRelations(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project inverse relations project-id")

	_, err = ListProjectIssues(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project issues project-id")

	_, err = ListLabelsForProject(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project labels project-id")

	_, err = ListProjectNeeds(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project customer needs project-id")

	_, err = ListProjectRelationsForProject(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project relations project-id")

	_, err = ListProjectTeams(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project teams project-id")

	_, err = ListProjectComments(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project comments project-id")

	_, err = ListProjectUpdates(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project updates project-id")

	_, err = GetProjectFilterSuggestion(context.Background(), graphqlClient, "started projects", "team-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project filter suggestion")

	_, err = ListAllProjectUpdates(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project updates")

	_, err = GetProjectUpdateByID(context.Background(), graphqlClient, "project-update-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project update project-update-id")

	_, err = ListProjectUpdateComments(context.Background(), graphqlClient, "project-update-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project update comments project-update-id")

	_, err = ListProjectMilestones(context.Background(), graphqlClient, "project-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project milestones project-id")

	_, err = ListAllProjectMilestones(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project milestones")

	_, err = GetProjectMilestoneByID(context.Background(), graphqlClient, "project-milestone-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project milestone project-milestone-id")

	_, err = ListProjectMilestoneIssues(context.Background(), graphqlClient, "project-milestone-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project milestone issues project-milestone-id")

	_, err = ListProjects(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list projects")

	_, err = ListProjectStatuses(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project statuses")

	_, err = GetProjectStatusByID(context.Background(), graphqlClient, "project-status-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project status project-status-id")

	_, err = GetProjectStatusProjectCount(context.Background(), graphqlClient, "project-status-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project status project count project-status-id")

	_, err = ListProjectLabels(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project labels")

	_, err = GetProjectLabelByID(context.Background(), graphqlClient, "project-label-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project label project-label-id")

	_, err = ListProjectLabelChildren(context.Background(), graphqlClient, "project-label-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project label children project-label-id")

	_, err = ListProjectLabelProjects(context.Background(), graphqlClient, "project-label-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project label projects project-label-id")

	_, err = ListProjectRelations(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list project relations")

	_, err = GetProjectRelationByID(context.Background(), graphqlClient, "project-relation-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get project relation project-relation-id")

	_, err = ListIssueRelations(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue relations")

	_, err = GetIssueRelationByID(context.Background(), graphqlClient, "issue-relation-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue relation issue-relation-id")

	_, err = ListIssueToReleases(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue to releases")

	_, err = GetIssueToReleaseByID(context.Background(), graphqlClient, "issue-to-release-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue to release issue-to-release-id")

	_, err = ListTeamMemberships(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team memberships")

	_, err = GetTeamMembershipByID(context.Background(), graphqlClient, "team-membership-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get team membership team-membership-id")

	_, err = GetApplicationInfo(context.Background(), graphqlClient, "app-client-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get application info app-client-id")

	_, err = ListAgentActivities(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list agent activities")

	_, err = GetAgentActivityByID(context.Background(), graphqlClient, "agent-activity-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get agent activity agent-activity-id")

	_, err = ListAgentSkills(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list agent skills")

	_, err = GetAgentSkillByID(context.Background(), graphqlClient, "agent-skill-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get agent skill agent-skill-id")

	_, err = ListExternalUsers(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list external users")

	_, err = GetExternalUserByID(context.Background(), graphqlClient, "external-user-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get external user external-user-id")

	_, err = ListAuditEntryTypes(context.Background(), graphqlClient)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list audit entry types")

	_, err = ListIssueComments(context.Background(), graphqlClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue comments LIT-1")

	_, err = ListComments(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list comments")

	_, err = GetCommentByID(context.Background(), graphqlClient, "comment-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get comment comment-id")

	_, err = GetCommentBotActor(context.Background(), graphqlClient, "comment-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get comment bot actor comment-id")

	_, err = ListCommentChildren(context.Background(), graphqlClient, "comment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list comment children comment-id")

	_, err = ListCommentCreatedIssues(context.Background(), graphqlClient, "comment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list comment created issues comment-id")

	_, err = ListDocuments(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list documents")

	_, err = GetDocumentByID(context.Background(), graphqlClient, "document-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get document document-id")

	_, err = ListDocumentComments(context.Background(), graphqlClient, "document-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list document comments document-id")

	_, err = ListLabels(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list labels")

	_, err = GetLabelByID(context.Background(), graphqlClient, "label-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get label label-id")

	_, err = ListLabelChildren(context.Background(), graphqlClient, "label-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list label children label-id")

	_, err = ListLabelIssues(context.Background(), graphqlClient, "label-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list label issues label-id")

	_, err = ListTeams(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list teams")

	_, err = GetTeamByID(context.Background(), graphqlClient, "team-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get team team-id")

	_, err = ListTeamMembers(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team members team-id")

	_, err = ListTeamCycles(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team cycles team-id")

	_, err = ListTeamIssues(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team issues team-id")

	_, err = ListTeamLabels(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team labels team-id")

	_, err = ListTeamMembershipsForTeam(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team memberships team-id")

	_, err = ListTeamProjects(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team projects team-id")

	_, err = ListTeamReleasePipelines(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team release pipelines team-id")

	_, err = ListTeamWorkflowStates(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team states team-id")

	_, err = ListTeamGitAutomationStates(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team git automation states team-id")

	_, err = ListTeamTemplates(context.Background(), graphqlClient, "team-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list team templates team-id")

	_, err = ListUsers(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list users")

	_, err = GetUserByID(context.Background(), graphqlClient, "user-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get user user-id")

	_, err = GetViewerUser(context.Background(), graphqlClient)
	require.Error(t, err)
	require.Contains(t, err.Error(), "get viewer user")

	_, err = ListViewerDrafts(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list viewer drafts")

	_, err = ListUserAssignedIssues(context.Background(), graphqlClient, "user-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list user assigned issues user-id")

	_, err = ListUserCreatedIssues(context.Background(), graphqlClient, "user-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list user created issues user-id")

	_, err = ListUserDelegatedIssues(context.Background(), graphqlClient, "user-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list user delegated issues user-id")

	_, err = ListUserTeamMemberships(context.Background(), graphqlClient, "user-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list user team memberships user-id")

	_, err = ListUserTeams(context.Background(), graphqlClient, "user-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list user teams user-id")

	_, err = ListViewerAssignedIssues(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list viewer assigned issues")

	_, err = ListViewerCreatedIssues(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list viewer created issues")

	_, err = ListViewerDelegatedIssues(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list viewer delegated issues")

	_, err = ListViewerTeamMemberships(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list viewer team memberships")

	_, err = ListViewerTeams(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list viewer teams")

	_, err = ListWorkflowStates(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list workflow states")

	_, err = GetWorkflowStateByID(context.Background(), graphqlClient, "workflow-state-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get workflow state workflow-state-id")

	_, err = ListWorkflowStateIssues(context.Background(), graphqlClient, "workflow-state-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list workflow state issues workflow-state-id")

	_, err = ListTimeSchedules(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list time schedules")

	_, err = GetTimeScheduleByID(context.Background(), graphqlClient, "time-schedule-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get time schedule time-schedule-id")

	_, err = ListOrganizationLabels(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list organization labels")

	_, err = ListOrganizationProjectLabels(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list organization project labels")

	_, err = ListOrganizationTeams(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list organization teams")

	_, err = ListOrganizationTemplates(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list organization templates")

	_, err = ListOrganizationUsers(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list organization users")

	_, err = ListTemplates(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list templates")

	_, err = GetTemplateByID(context.Background(), graphqlClient, "template-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get template template-id")

	_, err = ListInitiatives(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiatives")

	_, err = GetInitiativeByID(context.Background(), graphqlClient, "initiative-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get initiative initiative-id")

	_, err = ListInitiativeHistory(context.Background(), graphqlClient, "initiative-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative history initiative-id")

	_, err = ListInitiativeLinks(context.Background(), graphqlClient, "initiative-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative links initiative-id")

	_, err = ListSubInitiatives(context.Background(), graphqlClient, "initiative-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative sub-initiatives initiative-id")

	_, err = ListInitiativeUpdatesForInitiative(context.Background(), graphqlClient, "initiative-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative updates initiative-id")

	_, err = ListInitiativeDocuments(context.Background(), graphqlClient, "initiative-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative documents initiative-id")

	_, err = ListInitiativeProjects(context.Background(), graphqlClient, "initiative-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative projects initiative-id")

	_, err = ListInitiativeRelations(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative relations")

	_, err = GetInitiativeRelationByID(context.Background(), graphqlClient, "initiative-relation-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get initiative relation initiative-relation-id")

	_, err = ListInitiativeToProjects(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative to projects")

	_, err = GetInitiativeToProjectByID(context.Background(), graphqlClient, "initiative-to-project-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get initiative to project initiative-to-project-id")

	_, err = ListRoadmapToProjects(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list roadmap to projects")

	_, err = GetRoadmapToProjectByID(context.Background(), graphqlClient, "roadmap-to-project-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get roadmap to project roadmap-to-project-id")

	_, err = ListInitiativeUpdates(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative updates")

	_, err = GetInitiativeUpdateByID(context.Background(), graphqlClient, "initiative-update-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get initiative update initiative-update-id")

	_, err = ListInitiativeUpdateComments(context.Background(), graphqlClient, "initiative-update-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list initiative update comments initiative-update-id")

	_, err = ListRoadmaps(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list roadmaps")

	_, err = GetRoadmapByID(context.Background(), graphqlClient, "roadmap-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get roadmap roadmap-id")

	_, err = ListRoadmapProjects(context.Background(), graphqlClient, "roadmap-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list roadmap projects roadmap-id")

	_, err = ListCustomViews(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list custom views")

	_, err = GetCustomViewSubscriberStatus(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view subscribers custom-view-id")

	_, err = GetCustomViewByID(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view custom-view-id")

	_, err = ListCustomViewInitiatives(context.Background(), graphqlClient, "custom-view-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list custom view initiatives custom-view-id")

	_, err = ListCustomViewIssues(context.Background(), graphqlClient, "custom-view-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list custom view issues custom-view-id")

	_, err = GetCustomViewOrganizationPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view organization preferences custom-view-id")

	_, err = GetCustomViewOrganizationPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view organization preference values custom-view-id")

	_, err = ListCustomViewProjects(context.Background(), graphqlClient, "custom-view-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list custom view projects custom-view-id")

	_, err = GetCustomViewUserPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view user preferences custom-view-id")

	_, err = GetCustomViewUserPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view user preference values custom-view-id")

	_, err = GetCustomViewPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get custom view preference values custom-view-id")

	_, err = ListCustomers(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list customers")

	_, err = GetCustomerByID(context.Background(), graphqlClient, "customer-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get customer customer-id")

	_, err = ListCustomerNeeds(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list customer needs")

	_, err = GetCustomerNeedByID(context.Background(), graphqlClient, "customer-need-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get customer need customer-need-id")

	_, err = GetCustomerNeedProjectAttachment(context.Background(), graphqlClient, "customer-need-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get customer need project attachment customer-need-id")

	_, err = ListCustomerStatuses(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list customer statuses")

	_, err = GetCustomerStatusByID(context.Background(), graphqlClient, "customer-status-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get customer status customer-status-id")

	_, err = ListCustomerTiers(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list customer tiers")

	_, err = GetCustomerTierByID(context.Background(), graphqlClient, "customer-tier-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get customer tier customer-tier-id")

	_, err = ListFavorites(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list favorites")

	_, err = ListFavoriteChildren(context.Background(), graphqlClient, "favorite-folder-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list favorite children favorite-folder-id")

	_, err = GetFavoriteByID(context.Background(), graphqlClient, "favorite-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get favorite favorite-id")

	_, err = ListEmojis(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list emojis")

	_, err = GetEmojiByID(context.Background(), graphqlClient, "emoji-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get emoji emoji-id")

	_, err = ListNotifications(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list notifications")

	_, err = GetNotificationByID(context.Background(), graphqlClient, "notification-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get notification notification-id")

	_, err = ListNotificationSubscriptions(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list notification subscriptions")

	_, err = GetNotificationSubscriptionByID(context.Background(), graphqlClient, "notification-subscription-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get notification subscription notification-subscription-id")

	_, err = ListTriageResponsibilities(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list triage responsibilities")

	_, err = GetTriageResponsibilityByID(context.Background(), graphqlClient, "triage-responsibility-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get triage responsibility triage-responsibility-id")

	_, err = GetTriageResponsibilityManualSelection(context.Background(), graphqlClient, "triage-responsibility-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get triage responsibility manual selection triage-responsibility-id")

	_, err = ListSLAConfigurations(context.Background(), graphqlClient, "team-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "list SLA configurations team-id")

	_, err = SearchSemantic(context.Background(), graphqlClient, "agent search", 2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "semantic search")

	_, err = SearchDocuments(context.Background(), graphqlClient, "agent search", 2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "search documents")

	_, err = SearchIssues(context.Background(), graphqlClient, "agent search", 2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "search issues")

	_, err = SearchProjects(context.Background(), graphqlClient, "agent search", 2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "search projects")

	_, err = ListReleasePipelines(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release pipelines")

	_, err = GetReleasePipelineByID(context.Background(), graphqlClient, "release-pipeline-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get release pipeline release-pipeline-id")

	_, err = ListReleasePipelineReleases(context.Background(), graphqlClient, "release-pipeline-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release pipeline releases release-pipeline-id")

	_, err = ListReleasePipelineStages(context.Background(), graphqlClient, "release-pipeline-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release pipeline stages release-pipeline-id")

	_, err = ListReleasePipelineTeams(context.Background(), graphqlClient, "release-pipeline-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release pipeline teams release-pipeline-id")

	_, err = ListReleaseStages(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release stages")

	_, err = GetReleaseStageByID(context.Background(), graphqlClient, "release-stage-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get release stage release-stage-id")

	_, err = ListReleaseStageReleases(context.Background(), graphqlClient, "release-stage-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release stage releases release-stage-id")

	_, err = ListReleases(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list releases")

	_, err = GetReleaseByID(context.Background(), graphqlClient, "release-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get release release-id")

	_, err = ListReleaseHistory(context.Background(), graphqlClient, "release-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release history release-id")

	_, err = ListReleaseDocuments(context.Background(), graphqlClient, "release-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release documents release-id")

	_, err = ListReleaseIssues(context.Background(), graphqlClient, "release-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release issues release-id")

	_, err = ListReleaseLinks(context.Background(), graphqlClient, "release-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release links release-id")

	_, err = GetEntityExternalLinkByID(context.Background(), graphqlClient, "release-link-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get external link release-link-id")

	_, err = SearchReleases(context.Background(), graphqlClient, "mobile", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "search releases")

	_, err = ListReleaseNotes(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list release notes")

	_, err = GetReleaseNoteByID(context.Background(), graphqlClient, "release-note-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get release note release-note-id")

	_, err = ListAttachments(context.Background(), graphqlClient, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachments")

	_, err = ListAttachmentsForURL(context.Background(), graphqlClient, "https://example.com/spec", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachments for url https://example.com/spec")

	_, err = GetAttachmentByID(context.Background(), graphqlClient, "attachment-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get attachment attachment-id")

	_, err = GetIssueByVCSBranch(context.Background(), graphqlClient, "omer/branch")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue by vcs branch omer/branch")

	_, err = ListIssueVCSBranchAttachments(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch attachments omer/branch")

	_, err = GetIssueVCSBranchBotActor(context.Background(), graphqlClient, "omer/branch")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue vcs branch bot actor omer/branch")

	_, err = ListIssueVCSBranchChildren(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch children omer/branch")

	_, err = ListIssueVCSBranchDocuments(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch documents omer/branch")

	_, err = ListIssueVCSBranchFormerAttachments(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch former attachments omer/branch")

	_, err = ListIssueVCSBranchHistory(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch history omer/branch")

	_, err = ListIssueVCSBranchInverseRelations(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch inverse relations omer/branch")

	_, err = ListIssueVCSBranchLabels(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch labels omer/branch")

	_, err = ListIssueVCSBranchRelations(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch relations omer/branch")

	_, err = ListIssueVCSBranchReleases(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch releases omer/branch")

	_, err = ListIssueVCSBranchStateHistory(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch state history omer/branch")

	_, err = ListIssueVCSBranchSubscribers(context.Background(), graphqlClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch subscribers omer/branch")

	_, err = GetAttachmentIssue(context.Background(), graphqlClient, "attachment-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get attachment issue attachment-id")

	_, err = ListAttachmentIssueAttachments(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue attachments attachment-id")

	_, err = GetAttachmentIssueBotActor(context.Background(), graphqlClient, "attachment-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get attachment issue bot actor attachment-id")

	_, err = ListAttachmentIssueChildren(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue children attachment-id")

	_, err = ListAttachmentIssueDocuments(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue documents attachment-id")

	_, err = ListAttachmentIssueFormerAttachments(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue former attachments attachment-id")

	_, err = ListAttachmentIssueHistory(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue history attachment-id")

	_, err = ListAttachmentIssueInverseRelations(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue inverse relations attachment-id")

	_, err = ListAttachmentIssueLabels(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue labels attachment-id")

	_, err = ListAttachmentIssueRelations(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue relations attachment-id")

	_, err = ListAttachmentIssueReleases(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue releases attachment-id")

	_, err = ListAttachmentIssueStateHistory(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue state history attachment-id")

	_, err = ListAttachmentIssueSubscribers(context.Background(), graphqlClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue subscribers attachment-id")
}
