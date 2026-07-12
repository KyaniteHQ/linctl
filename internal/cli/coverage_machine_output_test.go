package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type machineOutputCase struct {
	name    string
	options rootOptions
	write   func(*cobra.Command, *rootOptions) error
	want    string
	wantErr bool
}

//nolint:maintidx // One row per writer/output-mode pair; the table is long but flat.
func machineOutputCases() []machineOutputCase {
	issue := client.IssueSummary{
		ID:         "issue-id",
		Identifier: "LIT-1",
		Title:      "Ship coverage",
		State:      "Todo",
		Project:    "Pinned project",
		URL:        "https://linear.app/issue/LIT-1",
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
		Status:     "next",
		TargetDate: "2026-06-30",
		Progress:   0.5,
	}
	projectStatus := client.ProjectStatusSummary{
		ID:    "project-status-id",
		Name:  "Backlog",
		Type:  "backlog",
		Color: "#bec2c8",
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
	user := client.UserSummary{
		ID:          "user-id",
		DisplayName: "Omer",
		Email:       "omer@example.com",
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
			CustomViewID: "custom-view-id",
			Layout:       "list",
		},
	}
	customViewPreferenceValues := client.CustomViewPreferencesValues{
		CustomViewID: "custom-view-id",
		Layout:       "board",
	}
	slaConfiguration := client.SLAConfigurationSummary{
		ID:      "sla-configuration-id",
		Name:    "First response",
		SLA:     3600000,
		SLAType: "all",
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
	issueStateSpan := client.IssueStateSpanSummary{
		ID:        "issue-state-span-id",
		StateName: "Started",
		StateType: "started",
		StartedAt: "2026-06-19T12:00:00Z",
	}

	return []machineOutputCase{
		{
			name:    "Issue format=full",
			options: rootOptions{format: "full"},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeIssue(command, options, issue) },
			want:    "LIT-1 Ship coverage [Todo] project=Pinned project url=https://linear.app/issue/LIT-1\n",
		},
		{
			name:    "Issue idOnly=true",
			options: rootOptions{idOnly: true},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeIssue(command, options, issue) },
			want:    "issue-id\n",
		},
		{
			name:    "IssueStateSpan idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueStateSpan(command, options, issueStateSpan)
			},
			want: "issue-state-span-id\n",
		},
		{
			name:    "Cycle format=minimal",
			options: rootOptions{format: "minimal"},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeCycle(command, options, cycle) },
			want:    "cycle-id\n",
		},
		{
			name:    "Cycle format=full",
			options: rootOptions{format: "full"},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeCycle(command, options, cycle) },
			want:    "cycle-id Planning cycle [active] starts_at=2026-07-01T00:00:00Z ends_at=2026-07-15T00:00:00Z progress=0.50\n",
		},
		{
			name:    "Cycle idOnly=true",
			options: rootOptions{idOnly: true},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeCycle(command, options, cycle) },
			want:    "cycle-id\n",
		},
		{
			name:    "Project format=minimal",
			options: rootOptions{format: "minimal"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProject(command, options, project)
			},
			want: "project-id\n",
		},
		{
			name:    "Project format=full",
			options: rootOptions{format: "full"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProject(command, options, project)
			},
			want: "project-id Coverage [Backlog] url=https://linear.app/project/project-id\n",
		},
		{
			name:    "Project idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProject(command, options, project)
			},
			want: "project-id\n",
		},
		{
			name:    "ProjectUpdate idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectUpdate(command, options, projectUpdate)
			},
			want: "project-update-id\n",
		},
		{
			name:    "ProjectMilestone format=minimal",
			options: rootOptions{format: "minimal"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectMilestone(command, options, milestone)
			},
			want: "project-milestone-id\n",
		},
		{
			name:    "ProjectMilestone format=full",
			options: rootOptions{format: "full"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectMilestone(command, options, milestone)
			},
			want: "project-milestone-id Launch milestone [next] target_date=2026-06-30 progress=0.50\n",
		},
		{
			name:    "ProjectMilestone idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectMilestone(command, options, milestone)
			},
			want: "project-milestone-id\n",
		},
		{
			name:    "ProjectStatus idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectStatus(command, options, projectStatus)
			},
			want: "project-status-id\n",
		},
		{
			name:    "ProjectLabel format=minimal",
			options: rootOptions{format: "minimal"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectLabel(command, options, projectLabel)
			},
			want: "project-label-id\n",
		},
		{
			name:    "ProjectLabel format=full",
			options: rootOptions{format: "full"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectLabel(command, options, projectLabel)
			},
			want: "project-label-id Roadmap #f2c94c group=false parent=Parent\n",
		},
		{
			name:    "ProjectLabel idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectLabel(command, options, projectLabel)
			},
			want: "project-label-id\n",
		},
		{
			name:    "ProjectLabel format=wide",
			options: rootOptions{format: "wide"},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectLabel(command, options, projectLabel)
			},
			wantErr: true,
		},
		{
			name:    "ProjectRelation idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeProjectRelation(command, options, projectRelation)
			},
			want: "project-relation-id\n",
		},
		{
			name:    "IssueRelation idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueRelation(command, options, issueRelation)
			},
			want: "issue-relation-id\n",
		},
		{
			name:    "IssueToRelease idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeIssueToRelease(command, options, issueToRelease)
			},
			want: "issue-to-release-id\n",
		},
		{
			name:    "TeamMembership idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTeamMembership(command, options, teamMembership)
			},
			want: "team-membership-id\n",
		},
		{
			name:    "Document idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeDocument(command, options, document)
			},
			want: "document-id\n",
		},
		{
			name:    "Label idOnly=true",
			options: rootOptions{idOnly: true},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeLabel(command, options, label) },
			want:    "label-id\n",
		},
		{
			name:    "Team idOnly=true",
			options: rootOptions{idOnly: true},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeTeam(command, options, team) },
			want:    "team-id\n",
		},
		{
			name:    "User idOnly=true",
			options: rootOptions{idOnly: true},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeUser(command, options, user) },
			want:    "user-id\n",
		},
		{
			name:    "Comment idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeComment(command, options, comment)
			},
			want: "comment-id\n",
		},
		{
			name:    "CommentMetadata idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCommentMetadata(command, options, commentMetadata)
			},
			want: "comment-id\n",
		},
		{
			name:    "WorkflowState idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeWorkflowState(command, options, workflowState)
			},
			want: "workflow-state-id\n",
		},
		{
			name:    "TimeSchedule idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTimeSchedule(command, options, timeSchedule)
			},
			want: "time-schedule-id\n",
		},
		{
			name:    "Template idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTemplate(command, options, template)
			},
			want: "template-id\n",
		},
		{
			name:    "Initiative idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiative(command, options, initiative)
			},
			want: "initiative-id\n",
		},
		{
			name:    "InitiativeHistory idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeHistory(command, options, initiativeHistory)
			},
			want: "initiative-history-id\n",
		},
		{
			name:    "InitiativeRelation idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeRelation(command, options, initiativeRelation)
			},
			want: "initiative-relation-id\n",
		},
		{
			name:    "InitiativeToProject idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeToProject(command, options, initiativeToProject)
			},
			want: "initiative-to-project-id\n",
		},
		{
			name:    "RoadmapToProject idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRoadmapToProject(command, options, roadmapToProject)
			},
			want: "roadmap-to-project-id\n",
		},
		{
			name:    "InitiativeUpdate idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeInitiativeUpdate(command, options, initiativeUpdate)
			},
			want: "initiative-update-id\n",
		},
		{
			name:    "Roadmap idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRoadmap(command, options, roadmap)
			},
			want: "roadmap-id\n",
		},
		{
			name:    "CustomView idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomView(command, options, customView)
			},
			want: "custom-view-id\n",
		},
		{
			name:    "CustomViewSubscriberStatus idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewSubscriberStatus(command, options, customViewSubscriberStatus)
			},
			want: "custom-view-id\n",
		},
		{
			name:    "CustomViewPreferences idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewPreferences(command, options, customViewPreferences)
			},
			want: "view-preferences-id\n",
		},
		{
			name:    "CustomViewPreferenceValues idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomViewPreferenceValues(command, options, customViewPreferenceValues)
			},
			want: "",
		},
		{
			name:    "SLAConfiguration idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSLAConfiguration(command, options, slaConfiguration)
			},
			want: "sla-configuration-id\n",
		},
		{
			name:    "Customer idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomer(command, options, customer)
			},
			want: "customer-id\n",
		},
		{
			name:    "CustomerNeed idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomerNeed(command, options, customerNeed)
			},
			want: "customer-need-id\n",
		},
		{
			name:    "CustomerStatus idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomerStatus(command, options, customerStatus)
			},
			want: "customer-status-id\n",
		},
		{
			name:    "CustomerTier idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeCustomerTier(command, options, customerTier)
			},
			want: "customer-tier-id\n",
		},
		{
			name:    "ApplicationInfo idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeApplicationInfo(command, options, application)
			},
			want: "app-id\n",
		},
		{
			name:    "AgentActivity idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAgentActivity(command, options, agentActivity)
			},
			want: "agent-activity-id\n",
		},
		{
			name:    "AgentSkill idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAgentSkill(command, options, agentSkill)
			},
			want: "agent-skill-id\n",
		},
		{
			name:    "ExternalUser idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeExternalUser(command, options, externalUser)
			},
			want: "external-user-id\n",
		},
		{
			name:    "AuditEntryType idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAuditEntryType(command, options, auditEntryType)
			},
			want: "user_login\n",
		},
		{
			name:    "Favorite idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeFavorite(command, options, favorite)
			},
			want: "favorite-id\n",
		},
		{
			name:    "Emoji idOnly=true",
			options: rootOptions{idOnly: true},
			write:   func(command *cobra.Command, options *rootOptions) error { return writeEmoji(command, options, emoji) },
			want:    "emoji-id\n",
		},
		{
			name:    "Attachment idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeAttachment(command, options, attachment)
			},
			want: "attachment-id\n",
		},
		{
			name:    "Notification idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeNotification(command, options, notification)
			},
			want: "notification-id\n",
		},
		{
			name:    "NotificationSubscription idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeNotificationSubscription(command, options, notificationSubscription)
			},
			want: "notification-subscription-id\n",
		},
		{
			name:    "TriageResponsibility idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTriageResponsibility(command, options, triageResponsibility)
			},
			want: "triage-responsibility-id\n",
		},
		{
			name:    "TriageResponsibilityManualSelection idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeTriageResponsibilityManualSelection(command, options, triageManualSelection)
			},
			want: "triage-responsibility-id\n",
		},
		{
			name:    "ReleasePipeline idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleasePipeline(command, options, releasePipeline)
			},
			want: "release-pipeline-id\n",
		},
		{
			name:    "ReleaseStage idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleaseStage(command, options, releaseStage)
			},
			want: "release-stage-id\n",
		},
		{
			name:    "Release idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeRelease(command, options, release)
			},
			want: "release-id\n",
		},
		{
			name:    "ReleaseHistory idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleaseHistory(command, options, releaseHistory)
			},
			want: "release-history-id\n",
		},
		{
			name:    "EntityExternalLink idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeEntityExternalLink(command, options, releaseLink)
			},
			want: "release-link-id\n",
		},
		{
			name:    "ReleaseNote idOnly=true",
			options: rootOptions{idOnly: true},
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeReleaseNote(command, options, releaseNote)
			},
			want: "release-note-id\n",
		},
	}
}

// idOnly and format modes emit the documented machine line per writer;
// invalid formats fail loudly.
func Test_CliOutputHelpers_cover_machine_output_edges(t *testing.T) {
	for _, test := range machineOutputCases() {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			command := &cobra.Command{}
			command.SetOut(&output)
			options := test.options
			err := test.write(command, &options)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, output.String())
		})
	}

	t.Run("quiet mode writes nothing for every render writer", func(t *testing.T) {
		output := bytes.Buffer{}
		command := &cobra.Command{}
		command.SetOut(&output)
		for _, test := range renderWriterCases() {
			require.NoError(t, test.write(command, &rootOptions{quiet: true}))
		}
		require.Empty(t, output.String())
	})

	t.Run("scalar and edge helpers", func(t *testing.T) {
		command := &cobra.Command{}
		output := bytes.Buffer{}
		command.SetOut(&output)
		issue := client.IssueSummary{Identifier: "LIT-1", Title: "Ship coverage", State: "Todo"}
		cycle := client.CycleSummary{ID: "cycle-id"}
		project := client.ProjectSummary{ID: "project-id"}
		milestone := client.ProjectMilestoneSummary{ID: "project-milestone-id"}
		auditEntryType := client.AuditEntryTypeSummary{Type: "user_login", Description: "User logged in"}
		require.Equal(t, "-", emptyDash(""))

		quietOutput := bytes.Buffer{}
		quietCommand := &cobra.Command{}
		quietCommand.SetOut(&quietOutput)
		require.NoError(t, writeJSONValue(quietCommand, &rootOptions{quiet: true}, issue))
		require.NoError(t, writeIssuePriorityValues(
			quietCommand,
			&rootOptions{quiet: true},
			[]client.IssuePriorityValue{{Priority: 1, Label: "Urgent"}},
		))
		require.NoError(t, writeIssueFilterSuggestion(
			quietCommand,
			&rootOptions{quiet: true},
			client.IssueFilterSuggestion{Filter: json.RawMessage(`{"state":{"type":{"eq":"started"}}}`), LogID: "log-id"},
		))
		require.NoError(t, writeIssueTitleSuggestion(
			quietCommand,
			&rootOptions{quiet: true},
			client.IssueTitleSuggestion{Title: "Improve exports", LogID: "log-id"},
		))
		require.NoError(t, writeAuditEntryTypes(
			quietCommand,
			&rootOptions{quiet: true},
			client.AuditEntryTypeList{AuditEntryTypes: []client.AuditEntryTypeSummary{auditEntryType}},
		))
		require.NoError(t, writeScalar(quietCommand, &rootOptions{quiet: true}, "title", "quiet"))
		wrote, err := writeIDOnly(quietCommand, &rootOptions{idOnly: true, quiet: true}, "issue-id")
		require.NoError(t, err)
		require.True(t, wrote)
		require.Empty(t, quietOutput.String())

		scalarJSONOutput := bytes.Buffer{}
		scalarJSONCommand := &cobra.Command{}
		scalarJSONCommand.SetOut(&scalarJSONOutput)
		require.NoError(t, writeScalar(scalarJSONCommand, &rootOptions{json: true}, "title", "Ship coverage"))
		require.Contains(t, scalarJSONOutput.String(), `"title": "Ship coverage"`)

		wrote, err = writeIDOnly(command, &rootOptions{idOnly: true}, "")
		require.Error(t, err)
		require.True(t, wrote)
		require.ErrorContains(t, err, "id is empty")

		require.NoError(t, ensureNonEmpty(&rootOptions{}, 0))
		require.Error(t, writeIssue(command, &rootOptions{format: "wide"}, issue))
		require.Error(t, writeCycle(command, &rootOptions{format: "wide"}, cycle))
		require.Error(t, writeProject(command, &rootOptions{format: "wide"}, project))
		require.Error(t, writeProjectMilestone(command, &rootOptions{format: "wide"}, milestone))
		_, err = normalizedHumanFormat(&rootOptions{format: "wide"})
		require.ErrorContains(t, err, "invalid format")

		err = writeJSONValue(command, &rootOptions{json: true, fields: "missing"}, issue)
		require.ErrorContains(t, err, "field \"missing\" is not present")
	})
}
