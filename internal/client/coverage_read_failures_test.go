package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// Every client read wraps a transport failure with its operation context and
// keeps the injected error in the chain for errors.Is-based recovery.
func Test_ClientReadFailureScenarios_wrap_graphql_errors(t *testing.T) {
	sentinel := errors.New("network down")
	graphqlClient := errorGraphQLClient{err: sentinel}
	ctx := context.Background()

	tests := []struct {
		name string
		call func() error
		want string
	}{
		{"ListIssuesByTeam", func() error {
			_, err := ListIssuesByTeam(ctx, graphqlClient, "team-id", 1, IssueListFilters{})
			return err
		}, "list issues"},
		{"ListIssues", func() error { _, err := ListIssues(ctx, graphqlClient, 1); return err }, "list issues"},
		{"ListIssuesByTeam filtered", func() error {
			_, err := ListIssuesByTeam(ctx, graphqlClient, "team-id", 1, IssueListFilters{StateType: "started"})
			return err
		}, "list issues"},
		{"ListNextIssuesByTeam", func() error { _, err := ListNextIssuesByTeam(ctx, graphqlClient, "team-id", 1); return err }, "list next issues"},
		{"ListIssuesByTeam blocked-by", func() error {
			_, err := ListIssuesByTeam(ctx, graphqlClient, "team-id", 1, IssueListFilters{BlockedBy: "LIT-1"})
			return err
		}, "list issues"},
		{"SearchIssuesByTeam", func() error { _, err := SearchIssuesByTeam(ctx, graphqlClient, "team-id", "needle", 1); return err }, "search issues"},
		{"SearchIssuesByFigmaFileKey", func() error { _, err := SearchIssuesByFigmaFileKey(ctx, graphqlClient, "figma-key", 1); return err }, "search issues by Figma file key"},
		{"ListIssuePriorityValues", func() error { _, err := ListIssuePriorityValues(ctx, graphqlClient); return err }, "list issue priority values"},
		{"GetIssueFilterSuggestion", func() error {
			_, err := GetIssueFilterSuggestion(ctx, graphqlClient, "started issues", "team-id", "")
			return err
		}, "get issue filter suggestion"},
		{"GetIssueTitleSuggestionFromCustomerRequest", func() error {
			_, err := GetIssueTitleSuggestionFromCustomerRequest(ctx, graphqlClient, "Customer asks for faster exports")
			return err
		}, "get issue title suggestion"},
		{"GetIssueByID", func() error { _, err := GetIssueByID(ctx, graphqlClient, "LIT-1"); return err }, "get issue LIT-1"},
		{"GetIssueDependencies", func() error { _, err := GetIssueDependencies(ctx, graphqlClient, "LIT-1", 1); return err }, "get issue dependencies LIT-1"},
		{"ListIssueAttachments", func() error { _, err := ListIssueAttachments(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue attachments LIT-1"},
		{"GetIssueBotActor", func() error { _, err := GetIssueBotActor(ctx, graphqlClient, "LIT-1"); return err }, "get issue bot actor LIT-1"},
		{"ListIssueChildren", func() error { _, err := ListIssueChildren(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue children LIT-1"},
		{"ListIssueDocuments", func() error { _, err := ListIssueDocuments(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue documents LIT-1"},
		{"ListIssueFormerAttachments", func() error { _, err := ListIssueFormerAttachments(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue former attachments LIT-1"},
		{"ListIssueHistory", func() error { _, err := ListIssueHistory(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue history LIT-1"},
		{"ListIssueInverseRelations", func() error { _, err := ListIssueInverseRelations(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue inverse relations LIT-1"},
		{"ListIssueLabels", func() error { _, err := ListIssueLabels(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue labels LIT-1"},
		{"ListIssueRelationsForIssue", func() error { _, err := ListIssueRelationsForIssue(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue relations LIT-1"},
		{"ListIssueReleases", func() error { _, err := ListIssueReleases(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue releases LIT-1"},
		{"ListIssueStateHistory", func() error { _, err := ListIssueStateHistory(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue state history LIT-1"},
		{"ListIssueSubscribers", func() error { _, err := ListIssueSubscribers(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue subscribers LIT-1"},
		{"ListProjectsByTeam", func() error { _, err := ListProjectsByTeam(ctx, graphqlClient, "team-id", 1); return err }, "list projects"},
		{"GetProjectByID", func() error { _, err := GetProjectByID(ctx, graphqlClient, "project-id"); return err }, "get project project-id"},
		{"ListProjectMembers", func() error { _, err := ListProjectMembers(ctx, graphqlClient, "project-id", 1); return err }, "list project members project-id"},
		{"ListProjectAttachments", func() error { _, err := ListProjectAttachments(ctx, graphqlClient, "project-id", 1); return err }, "list project attachments project-id"},
		{"ListProjectDocuments", func() error { _, err := ListProjectDocuments(ctx, graphqlClient, "project-id", 1); return err }, "list project documents project-id"},
		{"ListProjectExternalLinks", func() error { _, err := ListProjectExternalLinks(ctx, graphqlClient, "project-id", 1); return err }, "list project external links project-id"},
		{"ListProjectHistory", func() error { _, err := ListProjectHistory(ctx, graphqlClient, "project-id", 1); return err }, "list project history project-id"},
		{"ListProjectInitiativeToProjects", func() error {
			_, err := ListProjectInitiativeToProjects(ctx, graphqlClient, "project-id", 1)
			return err
		}, "list project initiative associations project-id"},
		{"ListProjectInitiatives", func() error { _, err := ListProjectInitiatives(ctx, graphqlClient, "project-id", 1); return err }, "list project initiatives project-id"},
		{"ListProjectInverseRelations", func() error { _, err := ListProjectInverseRelations(ctx, graphqlClient, "project-id", 1); return err }, "list project inverse relations project-id"},
		{"ListProjectIssues", func() error { _, err := ListProjectIssues(ctx, graphqlClient, "project-id", 1); return err }, "list project issues project-id"},
		{"ListLabelsForProject", func() error { _, err := ListLabelsForProject(ctx, graphqlClient, "project-id", 1); return err }, "list project labels project-id"},
		{"ListProjectNeeds", func() error { _, err := ListProjectNeeds(ctx, graphqlClient, "project-id", 1); return err }, "list project customer needs project-id"},
		{"ListProjectRelationsForProject", func() error {
			_, err := ListProjectRelationsForProject(ctx, graphqlClient, "project-id", 1)
			return err
		}, "list project relations project-id"},
		{"ListProjectTeams", func() error { _, err := ListProjectTeams(ctx, graphqlClient, "project-id", 1); return err }, "list project teams project-id"},
		{"ListProjectComments", func() error { _, err := ListProjectComments(ctx, graphqlClient, "project-id", 1); return err }, "list project comments project-id"},
		{"ListProjectUpdates", func() error { _, err := ListProjectUpdates(ctx, graphqlClient, "project-id", 1); return err }, "list project updates project-id"},
		{"GetProjectFilterSuggestion", func() error {
			_, err := GetProjectFilterSuggestion(ctx, graphqlClient, "started projects", "team-id")
			return err
		}, "get project filter suggestion"},
		{"ListAllProjectUpdates", func() error { _, err := ListAllProjectUpdates(ctx, graphqlClient, 1); return err }, "list project updates"},
		{"GetProjectUpdateByID", func() error { _, err := GetProjectUpdateByID(ctx, graphqlClient, "project-update-id"); return err }, "get project update project-update-id"},
		{"ListProjectUpdateComments", func() error {
			_, err := ListProjectUpdateComments(ctx, graphqlClient, "project-update-id", 1)
			return err
		}, "list project update comments project-update-id"},
		{"ListProjectMilestones", func() error { _, err := ListProjectMilestones(ctx, graphqlClient, "project-id", 1); return err }, "list project milestones project-id"},
		{"ListAllProjectMilestones", func() error { _, err := ListAllProjectMilestones(ctx, graphqlClient, 1); return err }, "list project milestones"},
		{"GetProjectMilestoneByID", func() error {
			_, err := GetProjectMilestoneByID(ctx, graphqlClient, "project-milestone-id")
			return err
		}, "get project milestone project-milestone-id"},
		{"ListProjectMilestoneIssues", func() error {
			_, err := ListProjectMilestoneIssues(ctx, graphqlClient, "project-milestone-id", 1)
			return err
		}, "list project milestone issues project-milestone-id"},
		{"ListProjects", func() error { _, err := ListProjects(ctx, graphqlClient, 1); return err }, "list projects"},
		{"ListProjectStatuses", func() error { _, err := ListProjectStatuses(ctx, graphqlClient, 1); return err }, "list project statuses"},
		{"GetProjectStatusByID", func() error { _, err := GetProjectStatusByID(ctx, graphqlClient, "project-status-id"); return err }, "get project status project-status-id"},
		{"GetProjectStatusProjectCount", func() error {
			_, err := GetProjectStatusProjectCount(ctx, graphqlClient, "project-status-id")
			return err
		}, "get project status project count project-status-id"},
		{"ListProjectLabels", func() error { _, err := ListProjectLabels(ctx, graphqlClient, 1); return err }, "list project labels"},
		{"GetProjectLabelByID", func() error { _, err := GetProjectLabelByID(ctx, graphqlClient, "project-label-id"); return err }, "get project label project-label-id"},
		{"ListProjectLabelChildren", func() error {
			_, err := ListProjectLabelChildren(ctx, graphqlClient, "project-label-id", 1)
			return err
		}, "list project label children project-label-id"},
		{"ListProjectLabelProjects", func() error {
			_, err := ListProjectLabelProjects(ctx, graphqlClient, "project-label-id", 1)
			return err
		}, "list project label projects project-label-id"},
		{"ListProjectRelations", func() error { _, err := ListProjectRelations(ctx, graphqlClient, 1); return err }, "list project relations"},
		{"GetProjectRelationByID", func() error { _, err := GetProjectRelationByID(ctx, graphqlClient, "project-relation-id"); return err }, "get project relation project-relation-id"},
		{"ListIssueRelations", func() error { _, err := ListIssueRelations(ctx, graphqlClient, 1); return err }, "list issue relations"},
		{"GetIssueRelationByID", func() error { _, err := GetIssueRelationByID(ctx, graphqlClient, "issue-relation-id"); return err }, "get issue relation issue-relation-id"},
		{"ListIssueToReleases", func() error { _, err := ListIssueToReleases(ctx, graphqlClient, 1); return err }, "list issue to releases"},
		{"GetIssueToReleaseByID", func() error { _, err := GetIssueToReleaseByID(ctx, graphqlClient, "issue-to-release-id"); return err }, "get issue to release issue-to-release-id"},
		{"ListTeamMemberships", func() error { _, err := ListTeamMemberships(ctx, graphqlClient, 1); return err }, "list team memberships"},
		{"GetTeamMembershipByID", func() error { _, err := GetTeamMembershipByID(ctx, graphqlClient, "team-membership-id"); return err }, "get team membership team-membership-id"},
		{"GetApplicationInfo", func() error { _, err := GetApplicationInfo(ctx, graphqlClient, "app-client-id"); return err }, "get application info app-client-id"},
		{"ListAgentActivities", func() error { _, err := ListAgentActivities(ctx, graphqlClient, 1); return err }, "list agent activities"},
		{"GetAgentActivityByID", func() error { _, err := GetAgentActivityByID(ctx, graphqlClient, "agent-activity-id"); return err }, "get agent activity agent-activity-id"},
		{"ListAgentSkills", func() error { _, err := ListAgentSkills(ctx, graphqlClient, 1); return err }, "list agent skills"},
		{"GetAgentSkillByID", func() error { _, err := GetAgentSkillByID(ctx, graphqlClient, "agent-skill-id"); return err }, "get agent skill agent-skill-id"},
		{"ListExternalUsers", func() error { _, err := ListExternalUsers(ctx, graphqlClient, 1); return err }, "list external users"},
		{"GetExternalUserByID", func() error { _, err := GetExternalUserByID(ctx, graphqlClient, "external-user-id"); return err }, "get external user external-user-id"},
		{"ListAuditEntryTypes", func() error { _, err := ListAuditEntryTypes(ctx, graphqlClient); return err }, "list audit entry types"},
		{"ListIssueComments", func() error { _, err := ListIssueComments(ctx, graphqlClient, "LIT-1", 1); return err }, "list issue comments LIT-1"},
		{"ListComments", func() error { _, err := ListComments(ctx, graphqlClient, 1); return err }, "list comments"},
		{"GetCommentByID", func() error { _, err := GetCommentByID(ctx, graphqlClient, "comment-id"); return err }, "get comment comment-id"},
		{"GetCommentBotActor", func() error { _, err := GetCommentBotActor(ctx, graphqlClient, "comment-id"); return err }, "get comment bot actor comment-id"},
		{"ListCommentChildren", func() error { _, err := ListCommentChildren(ctx, graphqlClient, "comment-id", 1); return err }, "list comment children comment-id"},
		{"ListCommentCreatedIssues", func() error { _, err := ListCommentCreatedIssues(ctx, graphqlClient, "comment-id", 1); return err }, "list comment created issues comment-id"},
		{"ListDocuments", func() error { _, err := ListDocuments(ctx, graphqlClient, 1); return err }, "list documents"},
		{"GetDocumentByID", func() error { _, err := GetDocumentByID(ctx, graphqlClient, "document-id"); return err }, "get document document-id"},
		{"ListDocumentComments", func() error { _, err := ListDocumentComments(ctx, graphqlClient, "document-id", 1); return err }, "list document comments document-id"},
		{"ListLabels", func() error { _, err := ListLabels(ctx, graphqlClient, 1); return err }, "list labels"},
		{"GetLabelByID", func() error { _, err := GetLabelByID(ctx, graphqlClient, "label-id"); return err }, "get label label-id"},
		{"ListLabelChildren", func() error { _, err := ListLabelChildren(ctx, graphqlClient, "label-id", 1); return err }, "list label children label-id"},
		{"ListLabelIssues", func() error { _, err := ListLabelIssues(ctx, graphqlClient, "label-id", 1); return err }, "list label issues label-id"},
		{"ListTeams", func() error { _, err := ListTeams(ctx, graphqlClient, 1); return err }, "list teams"},
		{"GetTeamByID", func() error { _, err := GetTeamByID(ctx, graphqlClient, "team-id"); return err }, "get team team-id"},
		{"ListTeamMembers", func() error { _, err := ListTeamMembers(ctx, graphqlClient, "team-id", 1); return err }, "list team members team-id"},
		{"ListTeamCycles", func() error { _, err := ListTeamCycles(ctx, graphqlClient, "team-id", 1); return err }, "list team cycles team-id"},
		{"ListTeamIssues", func() error { _, err := ListTeamIssues(ctx, graphqlClient, "team-id", 1); return err }, "list team issues team-id"},
		{"ListTeamLabels", func() error { _, err := ListTeamLabels(ctx, graphqlClient, "team-id", 1); return err }, "list team labels team-id"},
		{"ListTeamMembershipsForTeam", func() error { _, err := ListTeamMembershipsForTeam(ctx, graphqlClient, "team-id", 1); return err }, "list team memberships team-id"},
		{"ListTeamProjects", func() error { _, err := ListTeamProjects(ctx, graphqlClient, "team-id", 1); return err }, "list team projects team-id"},
		{"ListTeamReleasePipelines", func() error { _, err := ListTeamReleasePipelines(ctx, graphqlClient, "team-id", 1); return err }, "list team release pipelines team-id"},
		{"ListTeamWorkflowStates", func() error { _, err := ListTeamWorkflowStates(ctx, graphqlClient, "team-id", 1); return err }, "list team states team-id"},
		{"ListTeamGitAutomationStates", func() error { _, err := ListTeamGitAutomationStates(ctx, graphqlClient, "team-id", 1); return err }, "list team git automation states team-id"},
		{"ListTeamTemplates", func() error { _, err := ListTeamTemplates(ctx, graphqlClient, "team-id", 1); return err }, "list team templates team-id"},
		{"ListUsers", func() error { _, err := ListUsers(ctx, graphqlClient, 1); return err }, "list users"},
		{"GetUserByID", func() error { _, err := GetUserByID(ctx, graphqlClient, "user-id"); return err }, "get user user-id"},
		{"GetViewerUser", func() error { _, err := GetViewerUser(ctx, graphqlClient); return err }, "get viewer user"},
		{"ListViewerDrafts", func() error { _, err := ListViewerDrafts(ctx, graphqlClient, 1); return err }, "list viewer drafts"},
		{"ListUserAssignedIssues", func() error { _, err := ListUserAssignedIssues(ctx, graphqlClient, "user-id", 1); return err }, "list user assigned issues user-id"},
		{"ListUserCreatedIssues", func() error { _, err := ListUserCreatedIssues(ctx, graphqlClient, "user-id", 1); return err }, "list user created issues user-id"},
		{"ListUserDelegatedIssues", func() error { _, err := ListUserDelegatedIssues(ctx, graphqlClient, "user-id", 1); return err }, "list user delegated issues user-id"},
		{"ListUserTeamMemberships", func() error { _, err := ListUserTeamMemberships(ctx, graphqlClient, "user-id", 1); return err }, "list user team memberships user-id"},
		{"ListUserTeams", func() error { _, err := ListUserTeams(ctx, graphqlClient, "user-id", 1); return err }, "list user teams user-id"},
		{"ListViewerAssignedIssues", func() error { _, err := ListViewerAssignedIssues(ctx, graphqlClient, 1); return err }, "list viewer assigned issues"},
		{"ListViewerCreatedIssues", func() error { _, err := ListViewerCreatedIssues(ctx, graphqlClient, 1); return err }, "list viewer created issues"},
		{"ListViewerDelegatedIssues", func() error { _, err := ListViewerDelegatedIssues(ctx, graphqlClient, 1); return err }, "list viewer delegated issues"},
		{"ListViewerTeamMemberships", func() error { _, err := ListViewerTeamMemberships(ctx, graphqlClient, 1); return err }, "list viewer team memberships"},
		{"ListViewerTeams", func() error { _, err := ListViewerTeams(ctx, graphqlClient, 1); return err }, "list viewer teams"},
		{"ListWorkflowStates", func() error { _, err := ListWorkflowStates(ctx, graphqlClient, 1); return err }, "list workflow states"},
		{"GetWorkflowStateByID", func() error { _, err := GetWorkflowStateByID(ctx, graphqlClient, "workflow-state-id"); return err }, "get workflow state workflow-state-id"},
		{"ListWorkflowStateIssues", func() error {
			_, err := ListWorkflowStateIssues(ctx, graphqlClient, "workflow-state-id", 1)
			return err
		}, "list workflow state issues workflow-state-id"},
		{"ListTimeSchedules", func() error { _, err := ListTimeSchedules(ctx, graphqlClient, 1); return err }, "list time schedules"},
		{"GetTimeScheduleByID", func() error { _, err := GetTimeScheduleByID(ctx, graphqlClient, "time-schedule-id"); return err }, "get time schedule time-schedule-id"},
		{"ListOrganizationLabels", func() error { _, err := ListOrganizationLabels(ctx, graphqlClient, 1); return err }, "list organization labels"},
		{"ListOrganizationProjectLabels", func() error { _, err := ListOrganizationProjectLabels(ctx, graphqlClient, 1); return err }, "list organization project labels"},
		{"ListOrganizationTeams", func() error { _, err := ListOrganizationTeams(ctx, graphqlClient, 1); return err }, "list organization teams"},
		{"ListOrganizationTemplates", func() error { _, err := ListOrganizationTemplates(ctx, graphqlClient, 1); return err }, "list organization templates"},
		{"ListOrganizationUsers", func() error { _, err := ListOrganizationUsers(ctx, graphqlClient, 1); return err }, "list organization users"},
		{"ListTemplates", func() error { _, err := ListTemplates(ctx, graphqlClient, 1); return err }, "list templates"},
		{"GetTemplateByID", func() error { _, err := GetTemplateByID(ctx, graphqlClient, "template-id"); return err }, "get template template-id"},
		{"ListInitiatives", func() error { _, err := ListInitiatives(ctx, graphqlClient, 1); return err }, "list initiatives"},
		{"GetInitiativeByID", func() error { _, err := GetInitiativeByID(ctx, graphqlClient, "initiative-id"); return err }, "get initiative initiative-id"},
		{"ListInitiativeHistory", func() error { _, err := ListInitiativeHistory(ctx, graphqlClient, "initiative-id", 1); return err }, "list initiative history initiative-id"},
		{"ListInitiativeLinks", func() error { _, err := ListInitiativeLinks(ctx, graphqlClient, "initiative-id", 1); return err }, "list initiative links initiative-id"},
		{"ListSubInitiatives", func() error { _, err := ListSubInitiatives(ctx, graphqlClient, "initiative-id", 1); return err }, "list initiative sub-initiatives initiative-id"},
		{"ListInitiativeUpdatesForInitiative", func() error {
			_, err := ListInitiativeUpdatesForInitiative(ctx, graphqlClient, "initiative-id", 1)
			return err
		}, "list initiative updates initiative-id"},
		{"ListInitiativeDocuments", func() error { _, err := ListInitiativeDocuments(ctx, graphqlClient, "initiative-id", 1); return err }, "list initiative documents initiative-id"},
		{"ListInitiativeProjects", func() error { _, err := ListInitiativeProjects(ctx, graphqlClient, "initiative-id", 1); return err }, "list initiative projects initiative-id"},
		{"ListInitiativeRelations", func() error { _, err := ListInitiativeRelations(ctx, graphqlClient, 1); return err }, "list initiative relations"},
		{"GetInitiativeRelationByID", func() error {
			_, err := GetInitiativeRelationByID(ctx, graphqlClient, "initiative-relation-id")
			return err
		}, "get initiative relation initiative-relation-id"},
		{"ListInitiativeToProjects", func() error { _, err := ListInitiativeToProjects(ctx, graphqlClient, 1); return err }, "list initiative to projects"},
		{"GetInitiativeToProjectByID", func() error {
			_, err := GetInitiativeToProjectByID(ctx, graphqlClient, "initiative-to-project-id")
			return err
		}, "get initiative to project initiative-to-project-id"},
		{"ListRoadmapToProjects", func() error { _, err := ListRoadmapToProjects(ctx, graphqlClient, 1); return err }, "list roadmap to projects"},
		{"GetRoadmapToProjectByID", func() error {
			_, err := GetRoadmapToProjectByID(ctx, graphqlClient, "roadmap-to-project-id")
			return err
		}, "get roadmap to project roadmap-to-project-id"},
		{"ListInitiativeUpdates", func() error { _, err := ListInitiativeUpdates(ctx, graphqlClient, 1); return err }, "list initiative updates"},
		{"GetInitiativeUpdateByID", func() error {
			_, err := GetInitiativeUpdateByID(ctx, graphqlClient, "initiative-update-id")
			return err
		}, "get initiative update initiative-update-id"},
		{"ListInitiativeUpdateComments", func() error {
			_, err := ListInitiativeUpdateComments(ctx, graphqlClient, "initiative-update-id", 1)
			return err
		}, "list initiative update comments initiative-update-id"},
		{"ListRoadmaps", func() error { _, err := ListRoadmaps(ctx, graphqlClient, 1); return err }, "list roadmaps"},
		{"GetRoadmapByID", func() error { _, err := GetRoadmapByID(ctx, graphqlClient, "roadmap-id"); return err }, "get roadmap roadmap-id"},
		{"ListRoadmapProjects", func() error { _, err := ListRoadmapProjects(ctx, graphqlClient, "roadmap-id", 1); return err }, "list roadmap projects roadmap-id"},
		{"ListCustomViews", func() error { _, err := ListCustomViews(ctx, graphqlClient, 1); return err }, "list custom views"},
		{"GetCustomViewSubscriberStatus", func() error {
			_, err := GetCustomViewSubscriberStatus(ctx, graphqlClient, "custom-view-id")
			return err
		}, "get custom view subscribers custom-view-id"},
		{"GetCustomViewByID", func() error { _, err := GetCustomViewByID(ctx, graphqlClient, "custom-view-id"); return err }, "get custom view custom-view-id"},
		{"ListCustomViewInitiatives", func() error { _, err := ListCustomViewInitiatives(ctx, graphqlClient, "custom-view-id", 1); return err }, "list custom view initiatives custom-view-id"},
		{"ListCustomViewIssues", func() error { _, err := ListCustomViewIssues(ctx, graphqlClient, "custom-view-id", 1); return err }, "list custom view issues custom-view-id"},
		{"GetCustomViewOrganizationPreferences", func() error {
			_, err := GetCustomViewOrganizationPreferences(ctx, graphqlClient, "custom-view-id")
			return err
		}, "get custom view organization preferences custom-view-id"},
		{"GetCustomViewOrganizationPreferenceValues", func() error {
			_, err := GetCustomViewOrganizationPreferenceValues(ctx, graphqlClient, "custom-view-id")
			return err
		}, "get custom view organization preference values custom-view-id"},
		{"ListCustomViewProjects", func() error { _, err := ListCustomViewProjects(ctx, graphqlClient, "custom-view-id", 1); return err }, "list custom view projects custom-view-id"},
		{"GetCustomViewUserPreferences", func() error { _, err := GetCustomViewUserPreferences(ctx, graphqlClient, "custom-view-id"); return err }, "get custom view user preferences custom-view-id"},
		{"GetCustomViewUserPreferenceValues", func() error {
			_, err := GetCustomViewUserPreferenceValues(ctx, graphqlClient, "custom-view-id")
			return err
		}, "get custom view user preference values custom-view-id"},
		{"GetCustomViewPreferenceValues", func() error {
			_, err := GetCustomViewPreferenceValues(ctx, graphqlClient, "custom-view-id")
			return err
		}, "get custom view preference values custom-view-id"},
		{"ListCustomers", func() error { _, err := ListCustomers(ctx, graphqlClient, 1); return err }, "list customers"},
		{"GetCustomerByID", func() error { _, err := GetCustomerByID(ctx, graphqlClient, "customer-id"); return err }, "get customer customer-id"},
		{"ListCustomerNeeds", func() error { _, err := ListCustomerNeeds(ctx, graphqlClient, 1); return err }, "list customer needs"},
		{"GetCustomerNeedByID", func() error { _, err := GetCustomerNeedByID(ctx, graphqlClient, "customer-need-id"); return err }, "get customer need customer-need-id"},
		{"GetCustomerNeedProjectAttachment", func() error {
			_, err := GetCustomerNeedProjectAttachment(ctx, graphqlClient, "customer-need-id")
			return err
		}, "get customer need project attachment customer-need-id"},
		{"ListCustomerStatuses", func() error { _, err := ListCustomerStatuses(ctx, graphqlClient, 1); return err }, "list customer statuses"},
		{"GetCustomerStatusByID", func() error { _, err := GetCustomerStatusByID(ctx, graphqlClient, "customer-status-id"); return err }, "get customer status customer-status-id"},
		{"ListCustomerTiers", func() error { _, err := ListCustomerTiers(ctx, graphqlClient, 1); return err }, "list customer tiers"},
		{"GetCustomerTierByID", func() error { _, err := GetCustomerTierByID(ctx, graphqlClient, "customer-tier-id"); return err }, "get customer tier customer-tier-id"},
		{"ListFavorites", func() error { _, err := ListFavorites(ctx, graphqlClient, 1); return err }, "list favorites"},
		{"ListFavoriteChildren", func() error { _, err := ListFavoriteChildren(ctx, graphqlClient, "favorite-folder-id", 1); return err }, "list favorite children favorite-folder-id"},
		{"GetFavoriteByID", func() error { _, err := GetFavoriteByID(ctx, graphqlClient, "favorite-id"); return err }, "get favorite favorite-id"},
		{"ListEmojis", func() error { _, err := ListEmojis(ctx, graphqlClient, 1); return err }, "list emojis"},
		{"GetEmojiByID", func() error { _, err := GetEmojiByID(ctx, graphqlClient, "emoji-id"); return err }, "get emoji emoji-id"},
		{"ListNotifications", func() error { _, err := ListNotifications(ctx, graphqlClient, 1); return err }, "list notifications"},
		{"GetNotificationByID", func() error { _, err := GetNotificationByID(ctx, graphqlClient, "notification-id"); return err }, "get notification notification-id"},
		{"ListNotificationSubscriptions", func() error { _, err := ListNotificationSubscriptions(ctx, graphqlClient, 1); return err }, "list notification subscriptions"},
		{"GetNotificationSubscriptionByID", func() error {
			_, err := GetNotificationSubscriptionByID(ctx, graphqlClient, "notification-subscription-id")
			return err
		}, "get notification subscription notification-subscription-id"},
		{"ListTriageResponsibilities", func() error { _, err := ListTriageResponsibilities(ctx, graphqlClient, 1); return err }, "list triage responsibilities"},
		{"GetTriageResponsibilityByID", func() error {
			_, err := GetTriageResponsibilityByID(ctx, graphqlClient, "triage-responsibility-id")
			return err
		}, "get triage responsibility triage-responsibility-id"},
		{"GetTriageResponsibilityManualSelection", func() error {
			_, err := GetTriageResponsibilityManualSelection(ctx, graphqlClient, "triage-responsibility-id")
			return err
		}, "get triage responsibility manual selection triage-responsibility-id"},
		{"ListSLAConfigurations", func() error { _, err := ListSLAConfigurations(ctx, graphqlClient, "team-id"); return err }, "list SLA configurations team-id"},
		{"SearchSemantic", func() error { _, err := SearchSemantic(ctx, graphqlClient, "agent search", 2); return err }, "semantic search"},
		{"SearchDocuments", func() error { _, err := SearchDocuments(ctx, graphqlClient, "agent search", 2); return err }, "search documents"},
		{"SearchIssues", func() error { _, err := SearchIssues(ctx, graphqlClient, "agent search", 2); return err }, "search issues"},
		{"SearchProjects", func() error { _, err := SearchProjects(ctx, graphqlClient, "agent search", 2); return err }, "search projects"},
		{"ListReleasePipelines", func() error { _, err := ListReleasePipelines(ctx, graphqlClient, 1); return err }, "list release pipelines"},
		{"GetReleasePipelineByID", func() error { _, err := GetReleasePipelineByID(ctx, graphqlClient, "release-pipeline-id"); return err }, "get release pipeline release-pipeline-id"},
		{"ListReleasePipelineReleases", func() error {
			_, err := ListReleasePipelineReleases(ctx, graphqlClient, "release-pipeline-id", 1)
			return err
		}, "list release pipeline releases release-pipeline-id"},
		{"ListReleasePipelineStages", func() error {
			_, err := ListReleasePipelineStages(ctx, graphqlClient, "release-pipeline-id", 1)
			return err
		}, "list release pipeline stages release-pipeline-id"},
		{"ListReleasePipelineTeams", func() error {
			_, err := ListReleasePipelineTeams(ctx, graphqlClient, "release-pipeline-id", 1)
			return err
		}, "list release pipeline teams release-pipeline-id"},
		{"ListReleaseStages", func() error { _, err := ListReleaseStages(ctx, graphqlClient, 1); return err }, "list release stages"},
		{"GetReleaseStageByID", func() error { _, err := GetReleaseStageByID(ctx, graphqlClient, "release-stage-id"); return err }, "get release stage release-stage-id"},
		{"ListReleaseStageReleases", func() error {
			_, err := ListReleaseStageReleases(ctx, graphqlClient, "release-stage-id", 1)
			return err
		}, "list release stage releases release-stage-id"},
		{"ListReleases", func() error { _, err := ListReleases(ctx, graphqlClient, 1); return err }, "list releases"},
		{"GetReleaseByID", func() error { _, err := GetReleaseByID(ctx, graphqlClient, "release-id"); return err }, "get release release-id"},
		{"ListReleaseHistory", func() error { _, err := ListReleaseHistory(ctx, graphqlClient, "release-id", 1); return err }, "list release history release-id"},
		{"ListReleaseDocuments", func() error { _, err := ListReleaseDocuments(ctx, graphqlClient, "release-id", 1); return err }, "list release documents release-id"},
		{"ListReleaseIssues", func() error { _, err := ListReleaseIssues(ctx, graphqlClient, "release-id", 1); return err }, "list release issues release-id"},
		{"ListReleaseLinks", func() error { _, err := ListReleaseLinks(ctx, graphqlClient, "release-id", 1); return err }, "list release links release-id"},
		{"GetEntityExternalLinkByID", func() error { _, err := GetEntityExternalLinkByID(ctx, graphqlClient, "release-link-id"); return err }, "get external link release-link-id"},
		{"SearchReleases", func() error { _, err := SearchReleases(ctx, graphqlClient, "mobile", 1); return err }, "search releases"},
		{"ListReleaseNotes", func() error { _, err := ListReleaseNotes(ctx, graphqlClient, 1); return err }, "list release notes"},
		{"GetReleaseNoteByID", func() error { _, err := GetReleaseNoteByID(ctx, graphqlClient, "release-note-id"); return err }, "get release note release-note-id"},
		{"ListAttachments", func() error { _, err := ListAttachments(ctx, graphqlClient, 1); return err }, "list attachments"},
		{"ListAttachmentsForURL", func() error {
			_, err := ListAttachmentsForURL(ctx, graphqlClient, "https://example.com/spec", 1)
			return err
		}, "list attachments for url https://example.com/spec"},
		{"GetAttachmentByID", func() error { _, err := GetAttachmentByID(ctx, graphqlClient, "attachment-id"); return err }, "get attachment attachment-id"},
		{"GetIssueByVCSBranch", func() error { _, err := GetIssueByVCSBranch(ctx, graphqlClient, "omer/branch"); return err }, "get issue by vcs branch omer/branch"},
		{"ListIssueVCSBranchAttachments", func() error {
			_, err := ListIssueVCSBranchAttachments(ctx, graphqlClient, "omer/branch", 1)
			return err
		}, "list issue vcs branch attachments omer/branch"},
		{"GetIssueVCSBranchBotActor", func() error { _, err := GetIssueVCSBranchBotActor(ctx, graphqlClient, "omer/branch"); return err }, "get issue vcs branch bot actor omer/branch"},
		{"ListIssueVCSBranchChildren", func() error { _, err := ListIssueVCSBranchChildren(ctx, graphqlClient, "omer/branch", 1); return err }, "list issue vcs branch children omer/branch"},
		{"ListIssueVCSBranchDocuments", func() error { _, err := ListIssueVCSBranchDocuments(ctx, graphqlClient, "omer/branch", 1); return err }, "list issue vcs branch documents omer/branch"},
		{"ListIssueVCSBranchFormerAttachments", func() error {
			_, err := ListIssueVCSBranchFormerAttachments(ctx, graphqlClient, "omer/branch", 1)
			return err
		}, "list issue vcs branch former attachments omer/branch"},
		{"ListIssueVCSBranchHistory", func() error { _, err := ListIssueVCSBranchHistory(ctx, graphqlClient, "omer/branch", 1); return err }, "list issue vcs branch history omer/branch"},
		{"ListIssueVCSBranchInverseRelations", func() error {
			_, err := ListIssueVCSBranchInverseRelations(ctx, graphqlClient, "omer/branch", 1)
			return err
		}, "list issue vcs branch inverse relations omer/branch"},
		{"ListIssueVCSBranchLabels", func() error { _, err := ListIssueVCSBranchLabels(ctx, graphqlClient, "omer/branch", 1); return err }, "list issue vcs branch labels omer/branch"},
		{"ListIssueVCSBranchRelations", func() error { _, err := ListIssueVCSBranchRelations(ctx, graphqlClient, "omer/branch", 1); return err }, "list issue vcs branch relations omer/branch"},
		{"ListIssueVCSBranchReleases", func() error { _, err := ListIssueVCSBranchReleases(ctx, graphqlClient, "omer/branch", 1); return err }, "list issue vcs branch releases omer/branch"},
		{"ListIssueVCSBranchStateHistory", func() error {
			_, err := ListIssueVCSBranchStateHistory(ctx, graphqlClient, "omer/branch", 1)
			return err
		}, "list issue vcs branch state history omer/branch"},
		{"ListIssueVCSBranchSubscribers", func() error {
			_, err := ListIssueVCSBranchSubscribers(ctx, graphqlClient, "omer/branch", 1)
			return err
		}, "list issue vcs branch subscribers omer/branch"},
		{"GetAttachmentIssue", func() error { _, err := GetAttachmentIssue(ctx, graphqlClient, "attachment-id"); return err }, "get attachment issue attachment-id"},
		{"ListAttachmentIssueAttachments", func() error {
			_, err := ListAttachmentIssueAttachments(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue attachments attachment-id"},
		{"GetAttachmentIssueBotActor", func() error { _, err := GetAttachmentIssueBotActor(ctx, graphqlClient, "attachment-id"); return err }, "get attachment issue bot actor attachment-id"},
		{"ListAttachmentIssueChildren", func() error {
			_, err := ListAttachmentIssueChildren(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue children attachment-id"},
		{"ListAttachmentIssueDocuments", func() error {
			_, err := ListAttachmentIssueDocuments(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue documents attachment-id"},
		{"ListAttachmentIssueFormerAttachments", func() error {
			_, err := ListAttachmentIssueFormerAttachments(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue former attachments attachment-id"},
		{"ListAttachmentIssueHistory", func() error { _, err := ListAttachmentIssueHistory(ctx, graphqlClient, "attachment-id", 1); return err }, "list attachment issue history attachment-id"},
		{"ListAttachmentIssueInverseRelations", func() error {
			_, err := ListAttachmentIssueInverseRelations(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue inverse relations attachment-id"},
		{"ListAttachmentIssueLabels", func() error { _, err := ListAttachmentIssueLabels(ctx, graphqlClient, "attachment-id", 1); return err }, "list attachment issue labels attachment-id"},
		{"ListAttachmentIssueRelations", func() error {
			_, err := ListAttachmentIssueRelations(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue relations attachment-id"},
		{"ListAttachmentIssueReleases", func() error {
			_, err := ListAttachmentIssueReleases(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue releases attachment-id"},
		{"ListAttachmentIssueStateHistory", func() error {
			_, err := ListAttachmentIssueStateHistory(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue state history attachment-id"},
		{"ListAttachmentIssueSubscribers", func() error {
			_, err := ListAttachmentIssueSubscribers(ctx, graphqlClient, "attachment-id", 1)
			return err
		}, "list attachment issue subscribers attachment-id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.call()
			require.ErrorIs(t, err, sentinel)
			require.ErrorContains(t, err, test.want)
		})
	}
}
