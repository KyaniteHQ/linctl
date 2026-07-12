package cli

// Every list page type has one pageWithItems setter here: it puts the sorted
// items back on the page envelope for JSON output in the shared list pipeline.

import (
	"github.com/KyaniteHQ/linctl/internal/client"
)

func agentActivityPageWithItems(
	page client.AgentActivityList,
	activities []client.AgentActivitySummary,
) client.AgentActivityList {
	page.AgentActivities = activities
	return page
}

func agentSessionPageWithItems(
	page client.AgentSessionList,
	sessions []client.AgentSessionSummary,
) client.AgentSessionList {
	page.AgentSessions = sessions
	return page
}

func agentSkillPageWithItems(
	page client.AgentSkillList,
	skills []client.AgentSkillSummary,
) client.AgentSkillList {
	page.AgentSkills = skills
	return page
}

func attachmentPageWithItems(
	page client.AttachmentList,
	attachments []client.AttachmentSummary,
) client.AttachmentList {
	page.Attachments = attachments
	return page
}

func commentPageWithItems(page client.CommentList, comments []client.CommentSummary) client.CommentList {
	page.Comments = comments
	return page
}

func commentChildPageWithItems(
	page client.CommentChildList,
	comments []client.CommentMetadataSummary,
) client.CommentChildList {
	page.Comments = comments
	return page
}

func customerPageWithItems(
	page client.CustomerList,
	customers []client.CustomerSummary,
) client.CustomerList {
	page.Customers = customers
	return page
}

func customerNeedPageWithItems(
	page client.CustomerNeedList,
	needs []client.CustomerNeedSummary,
) client.CustomerNeedList {
	page.Needs = needs
	return page
}

func customerStatusPageWithItems(
	page client.CustomerStatusList,
	statuses []client.CustomerStatusSummary,
) client.CustomerStatusList {
	page.Statuses = statuses
	return page
}

func customerTierPageWithItems(
	page client.CustomerTierList,
	tiers []client.CustomerTierSummary,
) client.CustomerTierList {
	page.Tiers = tiers
	return page
}

func customViewPageWithItems(
	page client.CustomViewList,
	views []client.CustomViewSummary,
) client.CustomViewList {
	page.CustomViews = views
	return page
}

func documentPageWithItems(page client.DocumentList, documents []client.DocumentSummary) client.DocumentList {
	page.Documents = documents
	return page
}

func documentCommentPageWithItems(
	page client.DocumentCommentList,
	comments []client.CommentMetadataSummary,
) client.DocumentCommentList {
	page.Comments = comments
	return page
}

func emojiPageWithItems(
	page client.EmojiList,
	emojis []client.EmojiSummary,
) client.EmojiList {
	page.Emojis = emojis
	return page
}

func externalUserPageWithItems(
	page client.ExternalUserList,
	users []client.ExternalUserSummary,
) client.ExternalUserList {
	page.ExternalUsers = users
	return page
}

func favoritePageWithItems(
	page client.FavoriteList,
	favorites []client.FavoriteSummary,
) client.FavoriteList {
	page.Favorites = favorites
	return page
}

func initiativePageWithItems(
	page client.InitiativeList,
	initiatives []client.InitiativeSummary,
) client.InitiativeList {
	page.Initiatives = initiatives
	return page
}

func initiativeHistoryPageWithItems(
	page client.InitiativeHistoryList,
	history []client.InitiativeHistorySummary,
) client.InitiativeHistoryList {
	page.History = history
	return page
}

func initiativeRelationPageWithItems(
	page client.InitiativeRelationList,
	relations []client.InitiativeRelationSummary,
) client.InitiativeRelationList {
	page.Relations = relations
	return page
}

func initiativeToProjectPageWithItems(
	page client.InitiativeToProjectList,
	associations []client.InitiativeToProjectSummary,
) client.InitiativeToProjectList {
	page.Associations = associations
	return page
}

func initiativeUpdatePageWithItems(
	page client.InitiativeUpdateList,
	updates []client.InitiativeUpdateSummary,
) client.InitiativeUpdateList {
	page.Updates = updates
	return page
}

func initiativeUpdateCommentPageWithItems(
	page client.InitiativeUpdateCommentList,
	comments []client.CommentMetadataSummary,
) client.InitiativeUpdateCommentList {
	page.Comments = comments
	return page
}

func commentMetadataPageWithItems(
	page client.IssueCommentMetadataList,
	comments []client.CommentMetadataSummary,
) client.IssueCommentMetadataList {
	page.Comments = comments
	return page
}

func customerNeedMetadataPageWithItems(
	page client.IssueCustomerNeedMetadataList,
	needs []client.CustomerNeedMetadataSummary,
) client.IssueCustomerNeedMetadataList {
	page.Needs = needs
	return page
}

func issueHistoryPageWithItems(
	page client.IssueHistoryList,
	history []client.IssueHistorySummary,
) client.IssueHistoryList {
	page.History = history
	return page
}

func issueStateSpanPageWithItems(
	page client.IssueStateHistoryList,
	spans []client.IssueStateSpanSummary,
) client.IssueStateHistoryList {
	page.Spans = spans
	return page
}

func issueRelationPageWithItems(
	page client.IssueRelationList,
	relations []client.IssueRelationSummary,
) client.IssueRelationList {
	page.Relations = relations
	return page
}

func issueToReleasePageWithItems(
	page client.IssueToReleaseList,
	associations []client.IssueToReleaseSummary,
) client.IssueToReleaseList {
	page.Associations = associations
	return page
}

func labelPageWithItems(page client.LabelList, labels []client.LabelSummary) client.LabelList {
	page.Labels = labels
	return page
}

func labelChildrenPageWithItems(page client.LabelChildList, labels []client.LabelSummary) client.LabelChildList {
	page.Labels = labels
	return page
}

func labelIssuesPageWithItems(page client.LabelIssueList, issues []client.IssueSummary) client.LabelIssueList {
	page.Issues = issues
	return page
}

func notificationPageWithItems(
	page client.NotificationList,
	notifications []client.NotificationSummary,
) client.NotificationList {
	page.Notifications = notifications
	return page
}

func notificationSubscriptionPageWithItems(
	page client.NotificationSubscriptionList,
	subscriptions []client.NotificationSubscriptionSummary,
) client.NotificationSubscriptionList {
	page.Subscriptions = subscriptions
	return page
}

func projectPageWithItems(page client.ProjectList, projects []client.ProjectSummary) client.ProjectList {
	page.Projects = projects
	return page
}

func projectLabelPageWithItems(
	page client.ProjectLabelList,
	labels []client.ProjectLabelSummary,
) client.ProjectLabelList {
	page.ProjectLabels = labels
	return page
}

func projectLabelChildrenPageWithItems(
	page client.ProjectLabelChildrenList,
	labels []client.ProjectLabelSummary,
) client.ProjectLabelChildrenList {
	page.ProjectLabels = labels
	return page
}

func projectLabelProjectsPageWithItems(
	page client.ProjectLabelProjectsList,
	projects []client.ProjectSummary,
) client.ProjectLabelProjectsList {
	page.Projects = projects
	return page
}

func projectMilestonePageWithItems(
	page client.ProjectMilestoneList,
	milestones []client.ProjectMilestoneSummary,
) client.ProjectMilestoneList {
	page.Milestones = milestones
	return page
}

func projectMilestoneIssuePageWithItems(
	page client.ProjectMilestoneIssueList,
	issues []client.IssueSummary,
) client.ProjectMilestoneIssueList {
	page.Issues = issues
	return page
}

func projectRelationPageWithItems(
	page client.ProjectRelationList,
	relations []client.ProjectRelationSummary,
) client.ProjectRelationList {
	page.Relations = relations
	return page
}

func projectStatusPageWithItems(
	page client.ProjectStatusList,
	statuses []client.ProjectStatusSummary,
) client.ProjectStatusList {
	page.ProjectStatuses = statuses
	return page
}

func projectUpdatePageWithItems(
	page client.ProjectUpdateList,
	updates []client.ProjectUpdateSummary,
) client.ProjectUpdateList {
	page.Updates = updates
	return page
}

func projectUpdateCommentPageWithItems(
	page client.ProjectUpdateCommentList,
	comments []client.CommentMetadataSummary,
) client.ProjectUpdateCommentList {
	page.Comments = comments
	return page
}

func releasePageWithItems(page client.ReleaseList, releases []client.ReleaseSummary) client.ReleaseList {
	page.Releases = releases
	return page
}

func releaseHistoryPageWithItems(
	page client.ReleaseHistoryList,
	history []client.ReleaseHistorySummary,
) client.ReleaseHistoryList {
	page.History = history
	return page
}

func releaseLinksPageWithItems(
	page client.EntityExternalLinkList,
	links []client.EntityExternalLinkSummary,
) client.EntityExternalLinkList {
	page.Links = links
	return page
}

func releaseNotePageWithItems(
	page client.ReleaseNoteList,
	notes []client.ReleaseNoteSummary,
) client.ReleaseNoteList {
	page.ReleaseNotes = notes
	return page
}

func releasePipelinePageWithItems(
	page client.ReleasePipelineList,
	pipelines []client.ReleasePipelineSummary,
) client.ReleasePipelineList {
	page.ReleasePipelines = pipelines
	return page
}

func releaseStagePageWithItems(
	page client.ReleaseStageList,
	stages []client.ReleaseStageSummary,
) client.ReleaseStageList {
	page.ReleaseStages = stages
	return page
}

func roadmapPageWithItems(
	page client.RoadmapList,
	roadmaps []client.RoadmapSummary,
) client.RoadmapList {
	page.Roadmaps = roadmaps
	return page
}

func roadmapProjectPageWithItems(
	page client.RoadmapProjectList,
	projects []client.ProjectSummary,
) client.RoadmapProjectList {
	page.Projects = projects
	return page
}

func roadmapToProjectPageWithItems(
	page client.RoadmapToProjectList,
	associations []client.RoadmapToProjectSummary,
) client.RoadmapToProjectList {
	page.Associations = associations
	return page
}

func searchDocumentPageWithItems(
	page client.SearchDocumentList,
	documents []client.SearchDocumentSummary,
) client.SearchDocumentList {
	page.Documents = documents
	return page
}

func searchIssuePageWithItems(
	page client.SearchIssueList,
	issues []client.SearchIssueSummary,
) client.SearchIssueList {
	page.Issues = issues
	return page
}

func searchProjectPageWithItems(
	page client.SearchProjectList,
	projects []client.SearchProjectSummary,
) client.SearchProjectList {
	page.Projects = projects
	return page
}

func semanticSearchPageWithItems(
	page client.SemanticSearchList,
	results []client.SemanticSearchResultSummary,
) client.SemanticSearchList {
	page.Results = results
	return page
}

func teamPageWithItems(page client.TeamList, teams []client.TeamSummary) client.TeamList {
	page.Teams = teams
	return page
}

func teamMemberPageWithItems(page client.TeamMemberList, members []client.UserSummary) client.TeamMemberList {
	page.Members = members
	return page
}

func cyclePageWithItems(page client.CycleList, cycles []client.CycleSummary) client.CycleList {
	page.Cycles = cycles
	return page
}

func issuePageWithItems(page client.IssueList, issues []client.IssueSummary) client.IssueList {
	page.Issues = issues
	return page
}

func gitAutomationStatePageWithItems(
	page client.GitAutomationStateList,
	states []client.GitAutomationStateSummary,
) client.GitAutomationStateList {
	page.States = states
	return page
}

func teamMembershipPageWithItems(
	page client.TeamMembershipList,
	memberships []client.TeamMembershipSummary,
) client.TeamMembershipList {
	page.Memberships = memberships
	return page
}

func templatePageWithItems(page client.TemplateList, templates []client.TemplateSummary) client.TemplateList {
	page.Templates = templates
	return page
}

func timeSchedulePageWithItems(
	page client.TimeScheduleList,
	schedules []client.TimeScheduleSummary,
) client.TimeScheduleList {
	page.TimeSchedules = schedules
	return page
}

func triageResponsibilityPageWithItems(
	page client.TriageResponsibilityList,
	responsibilities []client.TriageResponsibilitySummary,
) client.TriageResponsibilityList {
	page.TriageResponsibilities = responsibilities
	return page
}

func userPageWithItems(page client.UserList, users []client.UserSummary) client.UserList {
	page.Users = users
	return page
}

func draftPageWithItems(page client.DraftList, drafts []client.DraftSummary) client.DraftList {
	page.Drafts = drafts
	return page
}

func workflowStatePageWithItems(
	page client.WorkflowStateList,
	states []client.WorkflowStateSummary,
) client.WorkflowStateList {
	page.WorkflowStates = states
	return page
}

func workflowStateIssuePageWithItems(
	page client.WorkflowStateIssueList,
	issues []client.IssueSummary,
) client.WorkflowStateIssueList {
	page.Issues = issues
	return page
}

func cycleIssuePageWithItems(page client.CycleIssueList, issues []client.IssueSummary) client.CycleIssueList {
	page.Issues = issues
	return page
}

func projectAttachmentPageWithItems(
	page client.ProjectAttachmentList,
	attachments []client.AttachmentSummary,
) client.ProjectAttachmentList {
	page.Attachments = attachments
	return page
}

func projectDocumentPageWithItems(
	page client.ProjectDocumentList,
	documents []client.DocumentSummary,
) client.ProjectDocumentList {
	page.Documents = documents
	return page
}

func projectExternalLinkPageWithItems(
	page client.ProjectExternalLinkList,
	links []client.EntityExternalLinkSummary,
) client.ProjectExternalLinkList {
	page.Links = links
	return page
}

func projectHistoryPageWithItems(
	page client.ProjectHistoryList,
	history []client.ProjectHistorySummary,
) client.ProjectHistoryList {
	page.History = history
	return page
}

func projectInitiativeToProjectPageWithItems(
	page client.ProjectInitiativeToProjectList,
	associations []client.InitiativeToProjectSummary,
) client.ProjectInitiativeToProjectList {
	page.Associations = associations
	return page
}

func projectInitiativePageWithItems(
	page client.ProjectInitiativeList,
	initiatives []client.InitiativeSummary,
) client.ProjectInitiativeList {
	page.Initiatives = initiatives
	return page
}

func projectIssuePageWithItems(page client.ProjectIssueList, issues []client.IssueSummary) client.ProjectIssueList {
	page.Issues = issues
	return page
}

func projectCommentPageWithItems(
	page client.ProjectCommentList,
	comments []client.CommentMetadataSummary,
) client.ProjectCommentList {
	page.Comments = comments
	return page
}

func projectProjectLabelPageWithItems(
	page client.ProjectProjectLabelList,
	labels []client.ProjectLabelSummary,
) client.ProjectProjectLabelList {
	page.ProjectLabels = labels
	return page
}

func projectMemberPageWithItems(
	page client.ProjectMemberList,
	members []client.ProjectMember,
) client.ProjectMemberList {
	page.Members = members
	return page
}

func projectCustomerNeedPageWithItems(
	page client.ProjectCustomerNeedList,
	needs []client.CustomerNeedSummary,
) client.ProjectCustomerNeedList {
	page.Needs = needs
	return page
}

func projectTeamPageWithItems(page client.ProjectTeamList, teams []client.TeamSummary) client.ProjectTeamList {
	page.Teams = teams
	return page
}

func projectProjectRelationPageWithItems(
	page client.ProjectProjectRelationList,
	relations []client.ProjectRelationSummary,
) client.ProjectProjectRelationList {
	page.Relations = relations
	return page
}
