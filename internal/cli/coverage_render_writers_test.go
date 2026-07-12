package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type renderWriterCase struct {
	name  string
	item  any
	write func(*cobra.Command, *rootOptions) error
	text  string
}

//nolint:maintidx // One row per render writer; the table is long but flat.
func renderWriterCases() []renderWriterCase {
	issue := client.IssueSummary{
		Identifier: "LIT-1",
		Title:      "Ship coverage",
		State:      "Todo",
		URL:        "https://linear.app/issue/LIT-1",
	}
	issueBotActor := client.IssueBotActor{
		IssueID: "issue-id",
		Bot: &client.ActorBotSummary{
			ID:   "bot-actor-id",
			Type: "github",
			Name: "GitHub",
		},
	}
	issueStateSpan := client.IssueStateSpanSummary{
		ID:        "issue-state-span-id",
		StateName: "Started",
		StateType: "started",
		StartedAt: "2026-06-19T12:00:00Z",
	}
	project := client.ProjectSummary{
		ID:   "project-id",
		Name: "Coverage",
		URL:  "https://linear.app/project/project-id",
		Status: client.ProjectStatus{
			Name: "Backlog",
		},
	}
	projectUpdate := client.ProjectUpdateSummary{
		ID:          "project-update-id",
		Body:        "First update",
		Health:      "onTrack",
		DisplayName: "Omer",
	}
	cycle := client.CycleSummary{
		ID:       "cycle-id",
		Name:     "Planning cycle",
		Status:   "active",
		StartsAt: "2026-07-01T00:00:00Z",
		EndsAt:   "2026-07-15T00:00:00Z",
		Progress: 0.5,
	}
	milestone := client.ProjectMilestoneSummary{
		ID:         "project-milestone-id",
		Name:       "Launch milestone",
		TargetDate: "2026-06-30",
		Status:     "next",
		Progress:   0.5,
	}
	projectStatus := client.ProjectStatusSummary{
		ID:    "project-status-id",
		Name:  "Backlog",
		Type:  "backlog",
		Color: "#bec2c8",
	}
	projectStatusProjectCount := client.ProjectStatusProjectCount{
		ProjectStatusID:   "project-status-id",
		Count:             12,
		PrivateCount:      2,
		ArchivedTeamCount: 1,
	}
	projectLabel := client.ProjectLabelSummary{
		ID:         "project-label-id",
		Name:       "Roadmap",
		Color:      "#f2c94c",
		ParentName: "Parent",
	}
	projectRelation := client.ProjectRelationSummary{
		ID:                 "project-relation-id",
		Type:               "blocks",
		ProjectName:        "Pinned project",
		RelatedProjectName: "Related project",
	}
	issueRelation := client.IssueRelationSummary{
		ID:                     "issue-relation-id",
		Type:                   "blocks",
		IssueIdentifier:        "LIT-1",
		RelatedIssueIdentifier: "LIT-2",
	}
	issueToRelease := client.IssueToReleaseSummary{
		ID:        "issue-to-release-id",
		IssueID:   "issue-id",
		ReleaseID: "release-id",
	}
	document := client.DocumentSummary{
		ID:         "document-id",
		Title:      "Spec",
		ParentType: "project",
	}
	label := client.LabelSummary{
		ID:    "label-id",
		Name:  "Bug",
		Color: "#ff0000",
	}
	team := client.TeamSummary{
		ID:   "team-id",
		Key:  "LIT",
		Name: "linctl",
	}
	teamMembership := client.TeamMembershipSummary{
		ID:          "team-membership-id",
		UserID:      "user-id",
		DisplayName: "Omer",
		TeamKey:     "LIT",
		Owner:       true,
		SortOrder:   1.5,
	}
	gitAutomationState := client.GitAutomationStateSummary{
		ID:                  "git-automation-state-id",
		Event:               "review",
		StateName:           "Started",
		TargetBranchPattern: "main",
	}
	user := client.UserSummary{
		ID:          "user-id",
		DisplayName: "Omer",
		Email:       "omer@example.com",
	}
	draft := client.DraftSummary{
		ID:          "draft-id",
		ParentType:  "issue",
		ParentKey:   "LIT-3",
		ParentTitle: "Draft issue",
	}
	comment := client.CommentSummary{
		ID:          "comment-id",
		Body:        "First comment",
		DisplayName: "Omer",
	}
	commentMetadata := client.CommentMetadataSummary{
		ID:          "comment-id",
		DisplayName: "Omer",
		CreatedAt:   "2026-06-19T12:00:00Z",
		ProjectID:   "project-id",
	}
	commentBotActor := client.CommentBotActor{
		CommentID: "comment-id",
		Bot: &client.ActorBotSummary{
			ID:   "bot-actor-id",
			Type: "github",
			Name: "GitHub",
		},
	}
	workflowState := client.WorkflowStateSummary{
		ID:   "workflow-state-id",
		Name: "Started",
		Type: "started",
	}
	timeSchedule := client.TimeScheduleSummary{
		ID:         "time-schedule-id",
		Name:       "Primary on-call",
		EntryCount: 1,
	}
	template := client.TemplateSummary{
		ID:      "template-id",
		Name:    "Bug report",
		Type:    "issue",
		TeamKey: "LIT",
	}
	initiative := client.InitiativeSummary{
		ID:     "initiative-id",
		Name:   "Platform",
		Status: "Active",
	}
	initiativeHistory := client.InitiativeHistorySummary{
		ID:           "initiative-history-id",
		InitiativeID: "initiative-id",
		EntryCount:   1,
		Entries:      json.RawMessage(`[{"type":"status"}]`),
	}
	initiativeRelation := client.InitiativeRelationSummary{
		ID:                    "initiative-relation-id",
		ParentInitiativeName:  "Platform",
		RelatedInitiativeName: "Child initiative",
		SortOrder:             1.5,
	}
	initiativeToProject := client.InitiativeToProjectSummary{
		ID:             "initiative-to-project-id",
		InitiativeName: "Platform",
		ProjectName:    "Pinned project",
		SortOrder:      "1",
	}
	roadmapToProject := client.RoadmapToProjectSummary{
		ID:          "roadmap-to-project-id",
		RoadmapName: "Platform roadmap",
		ProjectName: "Pinned project",
		SortOrder:   "1",
	}
	initiativeUpdate := client.InitiativeUpdateSummary{
		ID:          "initiative-update-id",
		Body:        "First initiative update",
		Health:      "onTrack",
		DisplayName: "Omer",
	}
	roadmap := client.RoadmapSummary{
		ID:     "roadmap-id",
		Name:   "Platform roadmap",
		SlugID: "platform-roadmap",
	}
	customView := client.CustomViewSummary{
		ID:        "custom-view-id",
		Name:      "My issues",
		ModelName: "Issue",
	}
	customViewSubscriberStatus := client.CustomViewSubscriberStatus{
		ID:             "custom-view-id",
		HasSubscribers: true,
	}
	customViewPreferences := client.CustomViewPreferences{
		CustomViewID: "custom-view-id",
		ID:           "view-preferences-id",
		Type:         "organization",
		ViewType:     "customView",
		Values: client.CustomViewPreferencesValues{
			CustomViewID:  "custom-view-id",
			Layout:        "list",
			ViewOrdering:  "priority",
			HiddenColumns: []string{"column-id"},
		},
	}
	customViewPreferenceValues := client.CustomViewPreferencesValues{
		CustomViewID:  "custom-view-id",
		Layout:        "board",
		ViewOrdering:  "updatedAt",
		HiddenColumns: []string{"column-id"},
	}
	slaConfiguration := client.SLAConfigurationSummary{
		ID:         "sla-configuration-id",
		Name:       "First response",
		SLA:        3600000,
		SLAType:    "all",
		RemovesSLA: false,
	}
	semanticSearchResult := client.SemanticSearchResultSummary{
		Type:  "issue",
		ID:    "issue-id",
		Key:   "LIT-3",
		Title: "Search result",
		URL:   "https://linear.app/kyanite/issue/LIT-3",
	}
	customer := client.CustomerSummary{
		ID:                   "customer-id",
		Name:                 "Acme",
		StatusName:           "Active",
		ApproximateNeedCount: 3,
	}
	customerNeed := client.CustomerNeedSummary{
		ID:           "customer-need-id",
		CustomerName: "Acme",
		Issue:        "LIT-1",
		Priority:     1,
	}
	customerStatus := client.CustomerStatusSummary{
		ID:          "customer-status-id",
		DisplayName: "Active",
		Color:       "#00ff00",
		Position:    1,
	}
	customerTier := client.CustomerTierSummary{
		ID:          "customer-tier-id",
		DisplayName: "Enterprise",
		Color:       "#0000ff",
		Position:    2,
	}
	organizationExistsStatus := client.OrganizationExistsStatus{
		URLKey:  "kyanite",
		Success: true,
		Exists:  true,
	}
	application := client.ApplicationInfo{
		ID:           "app-id",
		ClientID:     "app-client-id",
		Name:         "Demo App",
		Developer:    "Kyanite",
		DeveloperURL: "https://example.com",
	}
	agentSkill := client.AgentSkillSummary{
		ID:               "agent-skill-id",
		Title:            "Triage Helper",
		Body:             "Use this skill for triage.",
		Shared:           true,
		RecentUsageCount: 3,
	}
	externalUser := client.ExternalUserSummary{
		ID:          "external-user-id",
		Name:        "External User",
		DisplayName: "@external",
		LastSeen:    "2026-06-19T12:00:00Z",
	}
	auditEntryType := client.AuditEntryTypeSummary{
		Type:        "user_login",
		Description: "User logged in",
	}
	agentActivity := client.AgentActivitySummary{
		ID:             "agent-activity-id",
		AgentSessionID: "agent-session-id",
		ContentType:    "action",
		Content: client.AgentActivityContentSummary{
			Type:      "action",
			Action:    "read_file",
			Parameter: "README.md",
		},
		Signal: "continue",
		UserID: "user-id",
	}
	rateLimitStatus := client.RateLimitStatus{
		Identifier: "api-key",
		Kind:       "api",
		Limits: []client.RateLimit{
			{
				Type:            "complexity",
				RequestedAmount: 1,
				AllowedAmount:   1000,
				Period:          60000,
				RemainingAmount: 900,
				Reset:           1720000000000,
			},
		},
	}
	favorite := client.FavoriteSummary{
		ID:   "favorite-id",
		Type: "issue",
		URL:  "https://linear.app/kyanite/issue/LIT-1",
	}
	emoji := client.EmojiSummary{
		ID:     "emoji-id",
		Name:   "party",
		Source: "custom",
	}
	attachment := client.AttachmentSummary{
		ID:         "attachment-id",
		Title:      "Linked PR",
		SourceType: "github",
	}
	notification := client.NotificationSummary{
		ID:       "notification-id",
		Type:     "issueMention",
		Category: "mentions",
		Title:    "Mentioned you",
	}
	notificationSubscription := client.NotificationSubscriptionSummary{
		ID:         "notification-subscription-id",
		Active:     true,
		TargetType: "project",
		TargetName: "Roadmap",
	}
	triageResponsibility := client.TriageResponsibilitySummary{
		ID:              "triage-responsibility-id",
		Action:          "notify",
		TeamKey:         "LIT",
		CurrentUserName: "Omer",
	}
	triageManualSelection := client.TriageResponsibilityManualSelection{
		ID:      "triage-responsibility-id",
		UserIDs: []string{"user-id", "other-user-id"},
	}
	releasePipeline := client.ReleasePipelineSummary{
		ID:                      "release-pipeline-id",
		Name:                    "Production",
		SlugID:                  "production",
		ApproximateReleaseCount: 4,
	}
	releaseStage := client.ReleaseStageSummary{
		ID:           "release-stage-id",
		Name:         "Started",
		Type:         "started",
		PipelineName: "Production",
	}
	release := client.ReleaseSummary{
		ID:           "release-id",
		Name:         "Mobile 1.2.3",
		Version:      "v1.2.3",
		PipelineName: "Production",
		StageName:    "Started",
		IssueCount:   3,
	}
	releaseHistory := client.ReleaseHistorySummary{
		ID:         "release-history-id",
		ReleaseID:  "release-id",
		EntryCount: 1,
		Entries:    json.RawMessage(`[{"type":"stage"}]`),
	}
	issueHistory := client.IssueHistorySummary{
		ID:                 "issue-history-id",
		IssueID:            "issue-id",
		ActorID:            "user-id",
		UpdatedDescription: true,
	}
	releaseLink := client.EntityExternalLinkSummary{
		ID:        "release-link-id",
		Label:     "Runbook",
		URL:       "https://example.com/runbook",
		SortOrder: 1.5,
	}
	releaseNote := client.ReleaseNoteSummary{
		ID:           "release-note-id",
		Title:        "Launch notes",
		PipelineName: "Production",
		ReleaseCount: 2,
	}

	return []renderWriterCase{
		{
			name:  "Issue",
			item:  issue,
			write: func(command *cobra.Command, options *rootOptions) error { return writeIssue(command, options, issue) },
			text:  "LIT-1 Ship coverage [Todo]\n",
		},
		{
			name: "IssueBotActor",
			item: issueBotActor,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueBotActor(command, options, issueBotActor)
			},
			text: "issue-id bot bot-actor-id GitHub [github]\n",
		},
		{
			name: "IssueBotActor 2",
			item: client.IssueBotActor{IssueID: "plain-issue-id"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueBotActor(command, options, client.IssueBotActor{IssueID: "plain-issue-id"})
			},
			text: "plain-issue-id bot -\n",
		},
		{
			name: "IssueStateSpan",
			item: issueStateSpan,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueStateSpan(command, options, issueStateSpan)
			},
			text: "issue-state-span-id Started started 2026-06-19T12:00:00Z -> -\n",
		},
		{
			name:  "Cycle",
			item:  cycle,
			write: func(command *cobra.Command, options *rootOptions) error { return writeCycle(command, options, cycle) },
			text:  "cycle-id Planning cycle [active]\n",
		},
		{
			name: "Project",
			item: project,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProject(command, options, project)
			},
			text: "project-id Coverage [Backlog]\n",
		},
		{
			name: "ProjectUpdate",
			item: projectUpdate,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectUpdate(command, options, projectUpdate)
			},
			text: "project-update-id onTrack Omer First update\n",
		},
		{
			name: "ProjectMilestone",
			item: milestone,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectMilestone(command, options, milestone)
			},
			text: "project-milestone-id Launch milestone [next]\n",
		},
		{
			name: "ProjectStatus",
			item: projectStatus,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectStatus(command, options, projectStatus)
			},
			text: "project-status-id Backlog [backlog] #bec2c8\n",
		},
		{
			name: "ProjectStatusProjectCount",
			item: projectStatusProjectCount,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectStatusProjectCount(command, options, projectStatusProjectCount)
			},
			text: "project-status-id count 12 private 2 archived_team 1\n",
		},
		{
			name: "ProjectLabel",
			item: projectLabel,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectLabel(command, options, projectLabel)
			},
			text: "project-label-id Roadmap #f2c94c\n",
		},
		{
			name: "ProjectRelation",
			item: projectRelation,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectRelation(command, options, projectRelation)
			},
			text: "project-relation-id blocks Pinned project -> Related project\n",
		},
		{
			name: "IssueRelation",
			item: issueRelation,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueRelation(command, options, issueRelation)
			},
			text: "issue-relation-id blocks LIT-1 -> LIT-2\n",
		},
		{
			name: "IssueToRelease",
			item: issueToRelease,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueToRelease(command, options, issueToRelease)
			},
			text: "issue-to-release-id issue issue-id -> release release-id\n",
		},
		{
			name: "Document",
			item: document,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeDocument(command, options, document)
			},
			text: "document-id Spec [project]\n",
		},
		{
			name:  "Label",
			item:  label,
			write: func(command *cobra.Command, options *rootOptions) error { return writeLabel(command, options, label) },
			text:  "label-id Bug #ff0000\n",
		},
		{
			name:  "Team",
			item:  team,
			write: func(command *cobra.Command, options *rootOptions) error { return writeTeam(command, options, team) },
			text:  "team-id LIT linctl\n",
		},
		{
			name: "TeamMembership",
			item: teamMembership,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTeamMembership(command, options, teamMembership)
			},
			text: "team-membership-id LIT Omer owner true order 1.50\n",
		},
		{
			name: "GitAutomationState",
			item: gitAutomationState,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeGitAutomationState(command, options, gitAutomationState)
			},
			text: "git-automation-state-id review state Started target main\n",
		},
		{
			name:  "User",
			item:  user,
			write: func(command *cobra.Command, options *rootOptions) error { return writeUser(command, options, user) },
			text:  "user-id Omer <omer@example.com>\n",
		},
		{
			name:  "Draft",
			item:  draft,
			write: func(command *cobra.Command, options *rootOptions) error { return writeDraft(command, options, draft) },
			text:  "draft-id issue LIT-3 Draft issue\n",
		},
		{
			name: "Comment",
			item: comment,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeComment(command, options, comment)
			},
			text: "comment-id Omer First comment\n",
		},
		{
			name: "CommentMetadata",
			item: commentMetadata,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCommentMetadata(command, options, commentMetadata)
			},
			text: "comment-id Omer 2026-06-19T12:00:00Z\n",
		},
		{
			name: "CommentBotActor",
			item: commentBotActor,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCommentBotActor(command, options, commentBotActor)
			},
			text: "comment-id bot bot-actor-id GitHub [github]\n",
		},
		{
			name: "CommentBotActor 2",
			item: client.CommentBotActor{CommentID: "plain-comment-id"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCommentBotActor(command, options, client.CommentBotActor{CommentID: "plain-comment-id"})
			},
			text: "plain-comment-id bot -\n",
		},
		{
			name: "WorkflowState",
			item: workflowState,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeWorkflowState(command, options, workflowState)
			},
			text: "workflow-state-id Started [started]\n",
		},
		{
			name: "TimeSchedule",
			item: timeSchedule,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTimeSchedule(command, options, timeSchedule)
			},
			text: "time-schedule-id Primary on-call entries 1\n",
		},
		{
			name: "Template",
			item: template,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTemplate(command, options, template)
			},
			text: "template-id Bug report [issue] team LIT\n",
		},
		{
			name: "Initiative",
			item: initiative,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiative(command, options, initiative)
			},
			text: "initiative-id Platform [Active]\n",
		},
		{
			name: "InitiativeHistory",
			item: initiativeHistory,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeHistory(command, options, initiativeHistory)
			},
			text: "initiative-history-id initiative initiative-id entries 1\n",
		},
		{
			name: "InitiativeRelation",
			item: initiativeRelation,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeRelation(command, options, initiativeRelation)
			},
			text: "initiative-relation-id Platform -> Child initiative order 1.50\n",
		},
		{
			name: "InitiativeToProject",
			item: initiativeToProject,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeToProject(command, options, initiativeToProject)
			},
			text: "initiative-to-project-id Platform -> Pinned project order 1\n",
		},
		{
			name: "RoadmapToProject",
			item: roadmapToProject,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRoadmapToProject(command, options, roadmapToProject)
			},
			text: "roadmap-to-project-id Platform roadmap -> Pinned project order 1 [legacy]\n",
		},
		{
			name: "InitiativeUpdate",
			item: initiativeUpdate,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeUpdate(command, options, initiativeUpdate)
			},
			text: "initiative-update-id onTrack Omer First initiative update\n",
		},
		{
			name: "Roadmap",
			item: roadmap,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRoadmap(command, options, roadmap)
			},
			text: "roadmap-id Platform roadmap platform-roadmap [legacy]\n",
		},
		{
			name: "CustomView",
			item: customView,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomView(command, options, customView)
			},
			text: "custom-view-id My issues [Issue]\n",
		},
		{
			name: "CustomViewSubscriberStatus",
			item: customViewSubscriberStatus,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewSubscriberStatus(command, options, customViewSubscriberStatus)
			},
			text: "custom-view-id has_subscribers true\n",
		},
		{
			name: "CustomViewPreferences",
			item: customViewPreferences,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewPreferences(command, options, customViewPreferences)
			},
			text: "custom-view-id organization preferences organization customView layout list\n",
		},
		{
			name: "CustomViewPreferenceValues",
			item: customViewPreferenceValues,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewPreferenceValues(command, options, customViewPreferenceValues)
			},
			text: "custom-view-id preference values layout board ordering updatedAt\n",
		},
		{
			name: "CustomViewPreferences 2",
			item: client.CustomViewPreferences{CustomViewID: "empty-custom-view-id"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewPreferences(command, options, client.CustomViewPreferences{CustomViewID: "empty-custom-view-id"})
			},
			text: "empty-custom-view-id organization preferences -\n",
		},
		{
			name: "SLAConfiguration",
			item: slaConfiguration,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSLAConfiguration(command, options, slaConfiguration)
			},
			text: "sla-configuration-id First response sla 3600000 type all removes false\n",
		},
		{
			name: "SLAConfiguration 2",
			item: client.SLAConfigurationSummary{ID: "sla-remove-id", Name: "Remove SLA", RemovesSLA: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSLAConfiguration(command, options, client.SLAConfigurationSummary{ID: "sla-remove-id", Name: "Remove SLA", RemovesSLA: true})
			},
			text: "sla-remove-id Remove SLA sla - type - removes true\n",
		},
		{
			name: "Customer",
			item: customer,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomer(command, options, customer)
			},
			text: "customer-id Acme [Active] needs 3\n",
		},
		{
			name: "CustomerNeed",
			item: customerNeed,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomerNeed(command, options, customerNeed)
			},
			text: "customer-need-id Acme LIT-1 priority 1\n",
		},
		{
			name: "CustomerStatus",
			item: customerStatus,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomerStatus(command, options, customerStatus)
			},
			text: "customer-status-id Active #00ff00 1\n",
		},
		{
			name: "CustomerTier",
			item: customerTier,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomerTier(command, options, customerTier)
			},
			text: "customer-tier-id Enterprise #0000ff 2\n",
		},
		{
			name: "ApplicationInfo",
			item: application,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeApplicationInfo(command, options, application)
			},
			text: "app-id Demo App by Kyanite\n",
		},
		{
			name: "AgentActivity",
			item: agentActivity,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAgentActivity(command, options, agentActivity)
			},
			text: "agent-activity-id session agent-session-id [action] signal continue\n",
		},
		{
			name: "AgentSkill",
			item: agentSkill,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAgentSkill(command, options, agentSkill)
			},
			text: "agent-skill-id Triage Helper shared true recent 3\n",
		},
		{
			name: "ExternalUser",
			item: externalUser,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeExternalUser(command, options, externalUser)
			},
			text: "external-user-id External User @external last_seen 2026-06-19T12:00:00Z\n",
		},
		{
			name: "AuditEntryType",
			item: auditEntryType,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAuditEntryType(command, options, auditEntryType)
			},
			text: "user_login User logged in\n",
		},
		{
			name: "OrganizationExists",
			item: organizationExistsStatus,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeOrganizationExists(command, options, organizationExistsStatus)
			},
			text: "kyanite exists true success true\n",
		},
		{
			name: "RateLimitStatus",
			item: rateLimitStatus,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRateLimitStatus(command, options, rateLimitStatus)
			},
			text: "api api-key\ncomplexity remaining 900/1000 reset 1720000000000\n",
		},
		{
			name: "Favorite",
			item: favorite,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeFavorite(command, options, favorite)
			},
			text: "favorite-id [issue] https://linear.app/kyanite/issue/LIT-1\n",
		},
		{
			name:  "Emoji",
			item:  emoji,
			write: func(command *cobra.Command, options *rootOptions) error { return writeEmoji(command, options, emoji) },
			text:  "emoji-id party [custom]\n",
		},
		{
			name: "Attachment",
			item: attachment,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAttachment(command, options, attachment)
			},
			text: "attachment-id Linked PR [github]\n",
		},
		{
			name: "Notification",
			item: notification,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeNotification(command, options, notification)
			},
			text: "notification-id issueMention [mentions] Mentioned you\n",
		},
		{
			name: "NotificationSubscription",
			item: notificationSubscription,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeNotificationSubscription(command, options, notificationSubscription)
			},
			text: "notification-subscription-id project Roadmap active true\n",
		},
		{
			name: "TriageResponsibility",
			item: triageResponsibility,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTriageResponsibility(command, options, triageResponsibility)
			},
			text: "triage-responsibility-id team LIT action notify current Omer\n",
		},
		{
			name: "TriageResponsibilityManualSelection",
			item: triageManualSelection,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTriageResponsibilityManualSelection(command, options, triageManualSelection)
			},
			text: "triage-responsibility-id manual users user-id,other-user-id\n",
		},
		{
			name: "ReleasePipeline",
			item: releasePipeline,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleasePipeline(command, options, releasePipeline)
			},
			text: "release-pipeline-id Production production releases 4\n",
		},
		{
			name: "ReleaseStage",
			item: releaseStage,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleaseStage(command, options, releaseStage)
			},
			text: "release-stage-id Started [started] pipeline Production\n",
		},
		{
			name: "Release",
			item: release,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRelease(command, options, release)
			},
			text: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3\n",
		},
		{
			name: "ReleaseHistory",
			item: releaseHistory,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleaseHistory(command, options, releaseHistory)
			},
			text: "release-history-id release release-id entries 1\n",
		},
		{
			name: "IssueHistory",
			item: issueHistory,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueHistory(command, options, issueHistory)
			},
			text: "issue-history-id issue issue-id updated_description true\n",
		},
		{
			name: "EntityExternalLink",
			item: releaseLink,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeEntityExternalLink(command, options, releaseLink)
			},
			text: "release-link-id Runbook https://example.com/runbook order 1.5\n",
		},
		{
			name: "ReleaseNote",
			item: releaseNote,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleaseNote(command, options, releaseNote)
			},
			text: "release-note-id Launch notes pipeline Production releases 2\n",
		},
		{
			name: "SemanticSearchResult json only",
			item: semanticSearchResult,
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSemanticSearchResult(command, options, semanticSearchResult)
			},
		},
	}
}

// Each render writer emits its documented human line in text mode and the
// item's JSON encoding in JSON mode.
func Test_CliRenderHelpers_write_text_and_json_output(t *testing.T) {
	for _, test := range renderWriterCases() {
		if test.text != "" {
			t.Run(test.name+" text", func(t *testing.T) {
				output := bytes.Buffer{}
				command := &cobra.Command{}
				command.SetOut(&output)
				require.NoError(t, test.write(command, &rootOptions{}))
				require.Equal(t, test.text, output.String())
			})
		}
		t.Run(test.name+" json", func(t *testing.T) {
			output := bytes.Buffer{}
			command := &cobra.Command{}
			command.SetOut(&output)
			require.NoError(t, test.write(command, &rootOptions{json: true}))
			expected, err := json.Marshal(test.item)
			require.NoError(t, err)
			require.JSONEq(t, string(expected), output.String())
		})
	}
}
