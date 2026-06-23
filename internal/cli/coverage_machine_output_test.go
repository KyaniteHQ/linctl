package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CliOutputHelpers_cover_machine_output_edges(t *testing.T) {
	command := &cobra.Command{}
	output := bytes.Buffer{}
	command.SetOut(&output)
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
	organizationExistsStatus := client.OrganizationExistsStatus{
		URLKey:  "kyanite",
		Success: true,
		Exists:  true,
	}
	rateLimitStatus := client.RateLimitStatus{
		Identifier: "api-key",
		Kind:       "api",
		Limits: []client.RateLimit{
			{Type: "complexity", AllowedAmount: 1000, RemainingAmount: 900, Reset: 1720000000000},
		},
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
	issueBotActor := client.IssueBotActor{
		IssueID: "issue-id",
		Bot: &client.ActorBotSummary{
			ID:   "bot-actor-id",
			Type: "github",
			Name: "GitHub",
		},
	}
	commentBotActor := client.CommentBotActor{
		CommentID: "comment-id",
		Bot: &client.ActorBotSummary{
			ID:   "bot-actor-id",
			Type: "github",
			Name: "GitHub",
		},
	}

	require.NoError(t, writeIssue(command, &rootOptions{format: "full"}, issue))
	require.NoError(t, writeIssue(command, &rootOptions{idOnly: true}, issue))
	require.NoError(t, writeIssueStateSpan(command, &rootOptions{idOnly: true}, issueStateSpan))
	require.NoError(t, writeCycle(command, &rootOptions{format: "minimal"}, cycle))
	require.NoError(t, writeCycle(command, &rootOptions{format: "full"}, cycle))
	require.NoError(t, writeCycle(command, &rootOptions{idOnly: true}, cycle))
	require.NoError(t, writeProject(command, &rootOptions{format: "minimal"}, project))
	require.NoError(t, writeProject(command, &rootOptions{format: "full"}, project))
	require.NoError(t, writeProject(command, &rootOptions{idOnly: true}, project))
	require.NoError(t, writeProjectUpdate(command, &rootOptions{idOnly: true}, projectUpdate))
	require.NoError(t, writeProjectMilestone(command, &rootOptions{format: "minimal"}, milestone))
	require.NoError(t, writeProjectMilestone(command, &rootOptions{format: "full"}, milestone))
	require.NoError(t, writeProjectMilestone(command, &rootOptions{idOnly: true}, milestone))
	require.NoError(t, writeProjectStatus(command, &rootOptions{idOnly: true}, projectStatus))
	require.NoError(t, writeProjectLabel(command, &rootOptions{format: "minimal"}, projectLabel))
	require.NoError(t, writeProjectLabel(command, &rootOptions{format: "full"}, projectLabel))
	require.NoError(t, writeProjectLabel(command, &rootOptions{idOnly: true}, projectLabel))
	require.Error(t, writeProjectLabel(command, &rootOptions{format: "wide"}, projectLabel))
	require.NoError(t, writeProjectRelation(command, &rootOptions{idOnly: true}, projectRelation))
	require.NoError(t, writeIssueRelation(command, &rootOptions{idOnly: true}, issueRelation))
	require.NoError(t, writeIssueToRelease(command, &rootOptions{idOnly: true}, issueToRelease))
	require.NoError(t, writeTeamMembership(command, &rootOptions{idOnly: true}, teamMembership))
	require.NoError(t, writeDocument(command, &rootOptions{idOnly: true}, document))
	require.NoError(t, writeLabel(command, &rootOptions{idOnly: true}, label))
	require.NoError(t, writeTeam(command, &rootOptions{idOnly: true}, team))
	require.NoError(t, writeUser(command, &rootOptions{idOnly: true}, user))
	require.NoError(t, writeComment(command, &rootOptions{idOnly: true}, comment))
	require.NoError(t, writeCommentMetadata(command, &rootOptions{idOnly: true}, commentMetadata))
	require.NoError(t, writeWorkflowState(command, &rootOptions{idOnly: true}, workflowState))
	require.NoError(t, writeTimeSchedule(command, &rootOptions{idOnly: true}, timeSchedule))
	require.NoError(t, writeTemplate(command, &rootOptions{idOnly: true}, template))
	require.NoError(t, writeInitiative(command, &rootOptions{idOnly: true}, initiative))
	require.NoError(t, writeInitiativeHistory(command, &rootOptions{idOnly: true}, initiativeHistory))
	require.NoError(t, writeInitiativeRelation(command, &rootOptions{idOnly: true}, initiativeRelation))
	require.NoError(t, writeInitiativeToProject(command, &rootOptions{idOnly: true}, initiativeToProject))
	require.NoError(t, writeRoadmapToProject(command, &rootOptions{idOnly: true}, roadmapToProject))
	require.NoError(t, writeInitiativeUpdate(command, &rootOptions{idOnly: true}, initiativeUpdate))
	require.NoError(t, writeRoadmap(command, &rootOptions{idOnly: true}, roadmap))
	require.NoError(t, writeCustomView(command, &rootOptions{idOnly: true}, customView))
	require.NoError(t, writeCustomViewSubscriberStatus(command, &rootOptions{idOnly: true}, customViewSubscriberStatus))
	require.NoError(t, writeCustomViewPreferences(command, &rootOptions{idOnly: true}, customViewPreferences))
	require.NoError(t, writeCustomViewPreferenceValues(command, &rootOptions{idOnly: true}, customViewPreferenceValues))
	require.NoError(t, writeSLAConfiguration(command, &rootOptions{idOnly: true}, slaConfiguration))
	require.NoError(t, writeCustomer(command, &rootOptions{idOnly: true}, customer))
	require.NoError(t, writeCustomerNeed(command, &rootOptions{idOnly: true}, customerNeed))
	require.NoError(t, writeCustomerStatus(command, &rootOptions{idOnly: true}, customerStatus))
	require.NoError(t, writeCustomerTier(command, &rootOptions{idOnly: true}, customerTier))
	require.NoError(t, writeApplicationInfo(command, &rootOptions{idOnly: true}, application))
	require.NoError(t, writeAgentActivity(command, &rootOptions{idOnly: true}, agentActivity))
	require.NoError(t, writeAgentSkill(command, &rootOptions{idOnly: true}, agentSkill))
	require.NoError(t, writeExternalUser(command, &rootOptions{idOnly: true}, externalUser))
	require.NoError(t, writeAuditEntryType(command, &rootOptions{idOnly: true}, auditEntryType))
	require.NoError(t, writeFavorite(command, &rootOptions{idOnly: true}, favorite))
	require.NoError(t, writeEmoji(command, &rootOptions{idOnly: true}, emoji))
	require.NoError(t, writeAttachment(command, &rootOptions{idOnly: true}, attachment))
	require.NoError(t, writeNotification(command, &rootOptions{idOnly: true}, notification))
	require.NoError(t, writeNotificationSubscription(command, &rootOptions{idOnly: true}, notificationSubscription))
	require.NoError(t, writeTriageResponsibility(command, &rootOptions{idOnly: true}, triageResponsibility))
	require.NoError(t, writeTriageResponsibilityManualSelection(command, &rootOptions{idOnly: true}, triageManualSelection))
	require.NoError(t, writeReleasePipeline(command, &rootOptions{idOnly: true}, releasePipeline))
	require.NoError(t, writeReleaseStage(command, &rootOptions{idOnly: true}, releaseStage))
	require.NoError(t, writeRelease(command, &rootOptions{idOnly: true}, release))
	require.NoError(t, writeReleaseHistory(command, &rootOptions{idOnly: true}, releaseHistory))
	require.NoError(t, writeEntityExternalLink(command, &rootOptions{idOnly: true}, releaseLink))
	require.NoError(t, writeReleaseNote(command, &rootOptions{idOnly: true}, releaseNote))
	require.Contains(t, output.String(), "project=Pinned project")
	require.Contains(t, output.String(), "issue-id")
	require.Contains(t, output.String(), "starts_at=2026-07-01T00:00:00Z")
	require.Contains(t, output.String(), "cycle-id")
	require.Contains(t, output.String(), "project-id")
	require.Contains(t, output.String(), "project-update-id")
	require.Contains(t, output.String(), "target_date=2026-06-30")
	require.Contains(t, output.String(), "project-milestone-id")
	require.Contains(t, output.String(), "project-status-id")
	require.Contains(t, output.String(), "project-label-id")
	require.Contains(t, output.String(), "issue-relation-id")
	require.Contains(t, output.String(), "document-id")
	require.Contains(t, output.String(), "label-id")
	require.Contains(t, output.String(), "team-id")
	require.Contains(t, output.String(), "user-id")
	require.Contains(t, output.String(), "comment-id")
	require.Contains(t, output.String(), "workflow-state-id")
	require.Contains(t, output.String(), "time-schedule-id")
	require.Contains(t, output.String(), "template-id")
	require.Contains(t, output.String(), "initiative-id")
	require.Contains(t, output.String(), "initiative-history-id")
	require.Contains(t, output.String(), "initiative-relation-id")
	require.Contains(t, output.String(), "initiative-to-project-id")
	require.Contains(t, output.String(), "roadmap-to-project-id")
	require.Contains(t, output.String(), "initiative-update-id")
	require.Contains(t, output.String(), "roadmap-id")
	require.Contains(t, output.String(), "custom-view-id")
	require.Contains(t, output.String(), "sla-configuration-id")
	require.Contains(t, output.String(), "customer-id")
	require.Contains(t, output.String(), "customer-need-id")
	require.Contains(t, output.String(), "customer-status-id")
	require.Contains(t, output.String(), "customer-tier-id")
	require.Contains(t, output.String(), "app-id")
	require.Contains(t, output.String(), "agent-activity-id")
	require.Contains(t, output.String(), "agent-skill-id")
	require.Contains(t, output.String(), "user_login")
	require.Contains(t, output.String(), "favorite-id")
	require.Contains(t, output.String(), "emoji-id")
	require.Contains(t, output.String(), "attachment-id")
	require.Contains(t, output.String(), "notification-id")
	require.Contains(t, output.String(), "notification-subscription-id")
	require.Contains(t, output.String(), "triage-responsibility-id")
	require.Contains(t, output.String(), "release-pipeline-id")
	require.Contains(t, output.String(), "release-stage-id")
	require.Contains(t, output.String(), "release-id")
	require.Contains(t, output.String(), "release-history-id")
	require.Contains(t, output.String(), "release-link-id")
	require.Contains(t, output.String(), "release-note-id")
	require.Equal(t, "-", emptyDash(""))

	quietOutput := bytes.Buffer{}
	quietCommand := &cobra.Command{}
	quietCommand.SetOut(&quietOutput)
	require.NoError(t, writeJSONValue(quietCommand, &rootOptions{quiet: true}, issue))
	require.NoError(t, writeIssueBotActor(quietCommand, &rootOptions{quiet: true}, issueBotActor))
	require.NoError(t, writeIssueStateSpan(quietCommand, &rootOptions{quiet: true}, issueStateSpan))
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
	require.NoError(t, writeCycle(quietCommand, &rootOptions{quiet: true}, cycle))
	require.NoError(t, writeProject(quietCommand, &rootOptions{quiet: true}, project))
	require.NoError(t, writeProjectUpdate(quietCommand, &rootOptions{quiet: true}, projectUpdate))
	require.NoError(t, writeProjectMilestone(quietCommand, &rootOptions{quiet: true}, milestone))
	require.NoError(t, writeProjectStatus(quietCommand, &rootOptions{quiet: true}, projectStatus))
	require.NoError(t, writeProjectStatusProjectCount(
		quietCommand,
		&rootOptions{quiet: true},
		client.ProjectStatusProjectCount{ProjectStatusID: "project-status-id", Count: 12},
	))
	require.NoError(t, writeProjectLabel(quietCommand, &rootOptions{quiet: true}, projectLabel))
	require.NoError(t, writeProjectRelation(quietCommand, &rootOptions{quiet: true}, projectRelation))
	require.NoError(t, writeIssueRelation(quietCommand, &rootOptions{quiet: true}, issueRelation))
	require.NoError(t, writeIssueToRelease(quietCommand, &rootOptions{quiet: true}, issueToRelease))
	require.NoError(t, writeTeamMembership(quietCommand, &rootOptions{quiet: true}, teamMembership))
	require.NoError(t, writeDocument(quietCommand, &rootOptions{quiet: true}, document))
	require.NoError(t, writeLabel(quietCommand, &rootOptions{quiet: true}, label))
	require.NoError(t, writeTeam(quietCommand, &rootOptions{quiet: true}, team))
	require.NoError(t, writeUser(quietCommand, &rootOptions{quiet: true}, user))
	require.NoError(t, writeComment(quietCommand, &rootOptions{quiet: true}, comment))
	require.NoError(t, writeCommentMetadata(quietCommand, &rootOptions{quiet: true}, commentMetadata))
	require.NoError(t, writeCommentBotActor(quietCommand, &rootOptions{quiet: true}, commentBotActor))
	require.NoError(t, writeWorkflowState(quietCommand, &rootOptions{quiet: true}, workflowState))
	require.NoError(t, writeTimeSchedule(quietCommand, &rootOptions{quiet: true}, timeSchedule))
	require.NoError(t, writeTemplate(quietCommand, &rootOptions{quiet: true}, template))
	require.NoError(t, writeInitiative(quietCommand, &rootOptions{quiet: true}, initiative))
	require.NoError(t, writeInitiativeHistory(quietCommand, &rootOptions{quiet: true}, initiativeHistory))
	require.NoError(t, writeInitiativeRelation(quietCommand, &rootOptions{quiet: true}, initiativeRelation))
	require.NoError(t, writeInitiativeToProject(quietCommand, &rootOptions{quiet: true}, initiativeToProject))
	require.NoError(t, writeRoadmapToProject(quietCommand, &rootOptions{quiet: true}, roadmapToProject))
	require.NoError(t, writeInitiativeUpdate(quietCommand, &rootOptions{quiet: true}, initiativeUpdate))
	require.NoError(t, writeRoadmap(quietCommand, &rootOptions{quiet: true}, roadmap))
	require.NoError(t, writeCustomView(quietCommand, &rootOptions{quiet: true}, customView))
	require.NoError(t, writeCustomViewSubscriberStatus(quietCommand, &rootOptions{quiet: true}, customViewSubscriberStatus))
	require.NoError(t, writeCustomViewPreferences(quietCommand, &rootOptions{quiet: true}, customViewPreferences))
	require.NoError(t, writeCustomViewPreferenceValues(quietCommand, &rootOptions{quiet: true}, customViewPreferenceValues))
	require.NoError(t, writeSLAConfiguration(quietCommand, &rootOptions{quiet: true}, slaConfiguration))
	require.NoError(t, writeCustomer(quietCommand, &rootOptions{quiet: true}, customer))
	require.NoError(t, writeCustomerNeed(quietCommand, &rootOptions{quiet: true}, customerNeed))
	require.NoError(t, writeCustomerStatus(quietCommand, &rootOptions{quiet: true}, customerStatus))
	require.NoError(t, writeCustomerTier(quietCommand, &rootOptions{quiet: true}, customerTier))
	require.NoError(t, writeApplicationInfo(quietCommand, &rootOptions{quiet: true}, application))
	require.NoError(t, writeAgentActivity(quietCommand, &rootOptions{quiet: true}, agentActivity))
	require.NoError(t, writeAgentSkill(quietCommand, &rootOptions{quiet: true}, agentSkill))
	require.NoError(t, writeExternalUser(quietCommand, &rootOptions{quiet: true}, externalUser))
	require.NoError(t, writeAuditEntryType(quietCommand, &rootOptions{quiet: true}, auditEntryType))
	require.NoError(t, writeAuditEntryTypes(
		quietCommand,
		&rootOptions{quiet: true},
		client.AuditEntryTypeList{AuditEntryTypes: []client.AuditEntryTypeSummary{auditEntryType}},
	))
	require.NoError(t, writeOrganizationExists(quietCommand, &rootOptions{quiet: true}, organizationExistsStatus))
	require.NoError(t, writeRateLimitStatus(quietCommand, &rootOptions{quiet: true}, rateLimitStatus))
	require.NoError(t, writeFavorite(quietCommand, &rootOptions{quiet: true}, favorite))
	require.NoError(t, writeEmoji(quietCommand, &rootOptions{quiet: true}, emoji))
	require.NoError(t, writeAttachment(quietCommand, &rootOptions{quiet: true}, attachment))
	require.NoError(t, writeNotification(quietCommand, &rootOptions{quiet: true}, notification))
	require.NoError(t, writeNotificationSubscription(quietCommand, &rootOptions{quiet: true}, notificationSubscription))
	require.NoError(t, writeTriageResponsibility(quietCommand, &rootOptions{quiet: true}, triageResponsibility))
	require.NoError(t, writeTriageResponsibilityManualSelection(quietCommand, &rootOptions{quiet: true}, triageManualSelection))
	require.NoError(t, writeReleasePipeline(quietCommand, &rootOptions{quiet: true}, releasePipeline))
	require.NoError(t, writeReleaseStage(quietCommand, &rootOptions{quiet: true}, releaseStage))
	require.NoError(t, writeRelease(quietCommand, &rootOptions{quiet: true}, release))
	require.NoError(t, writeReleaseHistory(quietCommand, &rootOptions{quiet: true}, releaseHistory))
	require.NoError(t, writeEntityExternalLink(quietCommand, &rootOptions{quiet: true}, releaseLink))
	require.NoError(t, writeReleaseNote(quietCommand, &rootOptions{quiet: true}, releaseNote))
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
	require.Contains(t, err.Error(), "id is empty")

	require.NoError(t, ensureNonEmpty(&rootOptions{}, 0))
	require.Error(t, writeIssue(command, &rootOptions{format: "wide"}, issue))
	require.Error(t, writeCycle(command, &rootOptions{format: "wide"}, cycle))
	require.Error(t, writeProject(command, &rootOptions{format: "wide"}, project))
	require.Error(t, writeProjectMilestone(command, &rootOptions{format: "wide"}, milestone))
	_, err = normalizedHumanFormat(&rootOptions{format: "wide"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid format")

	err = writeJSONValue(command, &rootOptions{json: true, fields: "missing"}, issue)
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"missing\" is not present")
}
