package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_CliRenderHelpers_write_text_and_json_output(t *testing.T) {
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

	textOut := bytes.Buffer{}
	textCommand := &cobra.Command{}
	textCommand.SetOut(&textOut)
	textOptions := rootOptions{}

	require.NoError(t, writeIssue(textCommand, &textOptions, issue))
	require.NoError(t, writeIssueBotActor(textCommand, &textOptions, issueBotActor))
	require.NoError(t, writeIssueBotActor(textCommand, &textOptions, client.IssueBotActor{IssueID: "plain-issue-id"}))
	require.NoError(t, writeIssueStateSpan(textCommand, &textOptions, issueStateSpan))
	require.NoError(t, writeCycle(textCommand, &textOptions, cycle))
	require.NoError(t, writeProject(textCommand, &textOptions, project))
	require.NoError(t, writeProjectUpdate(textCommand, &textOptions, projectUpdate))
	require.NoError(t, writeProjectMilestone(textCommand, &textOptions, milestone))
	require.NoError(t, writeProjectStatus(textCommand, &textOptions, projectStatus))
	require.NoError(t, writeProjectStatusProjectCount(textCommand, &textOptions, projectStatusProjectCount))
	require.NoError(t, writeProjectLabel(textCommand, &textOptions, projectLabel))
	require.NoError(t, writeProjectRelation(textCommand, &textOptions, projectRelation))
	require.NoError(t, writeIssueRelation(textCommand, &textOptions, issueRelation))
	require.NoError(t, writeIssueToRelease(textCommand, &textOptions, issueToRelease))
	require.NoError(t, writeDocument(textCommand, &textOptions, document))
	require.NoError(t, writeLabel(textCommand, &textOptions, label))
	require.NoError(t, writeTeam(textCommand, &textOptions, team))
	require.NoError(t, writeTeamMembership(textCommand, &textOptions, teamMembership))
	require.NoError(t, writeGitAutomationState(textCommand, &textOptions, gitAutomationState))
	require.NoError(t, writeUser(textCommand, &textOptions, user))
	require.NoError(t, writeDraft(textCommand, &textOptions, draft))
	require.NoError(t, writeComment(textCommand, &textOptions, comment))
	require.NoError(t, writeCommentMetadata(textCommand, &textOptions, commentMetadata))
	require.NoError(t, writeCommentBotActor(textCommand, &textOptions, commentBotActor))
	require.NoError(t, writeCommentBotActor(textCommand, &textOptions, client.CommentBotActor{CommentID: "plain-comment-id"}))
	require.NoError(t, writeWorkflowState(textCommand, &textOptions, workflowState))
	require.NoError(t, writeTimeSchedule(textCommand, &textOptions, timeSchedule))
	require.NoError(t, writeTemplate(textCommand, &textOptions, template))
	require.NoError(t, writeInitiative(textCommand, &textOptions, initiative))
	require.NoError(t, writeInitiativeHistory(textCommand, &textOptions, initiativeHistory))
	require.NoError(t, writeInitiativeRelation(textCommand, &textOptions, initiativeRelation))
	require.NoError(t, writeInitiativeToProject(textCommand, &textOptions, initiativeToProject))
	require.NoError(t, writeRoadmapToProject(textCommand, &textOptions, roadmapToProject))
	require.NoError(t, writeInitiativeUpdate(textCommand, &textOptions, initiativeUpdate))
	require.NoError(t, writeRoadmap(textCommand, &textOptions, roadmap))
	require.NoError(t, writeCustomView(textCommand, &textOptions, customView))
	require.NoError(t, writeCustomViewSubscriberStatus(textCommand, &textOptions, customViewSubscriberStatus))
	require.NoError(t, writeCustomViewPreferences(textCommand, &textOptions, customViewPreferences))
	require.NoError(t, writeCustomViewPreferenceValues(textCommand, &textOptions, customViewPreferenceValues))
	require.NoError(t, writeCustomViewPreferences(textCommand, &textOptions, client.CustomViewPreferences{CustomViewID: "empty-custom-view-id"}))
	require.NoError(t, writeSLAConfiguration(textCommand, &textOptions, slaConfiguration))
	require.NoError(t, writeSLAConfiguration(textCommand, &textOptions, client.SLAConfigurationSummary{ID: "sla-remove-id", Name: "Remove SLA", RemovesSLA: true}))
	require.NoError(t, writeCustomer(textCommand, &textOptions, customer))
	require.NoError(t, writeCustomerNeed(textCommand, &textOptions, customerNeed))
	require.NoError(t, writeCustomerStatus(textCommand, &textOptions, customerStatus))
	require.NoError(t, writeCustomerTier(textCommand, &textOptions, customerTier))
	require.NoError(t, writeApplicationInfo(textCommand, &textOptions, application))
	require.NoError(t, writeAgentActivity(textCommand, &textOptions, agentActivity))
	require.NoError(t, writeAgentSkill(textCommand, &textOptions, agentSkill))
	require.NoError(t, writeExternalUser(textCommand, &textOptions, externalUser))
	require.NoError(t, writeAuditEntryType(textCommand, &textOptions, auditEntryType))
	require.NoError(t, writeOrganizationExists(textCommand, &textOptions, organizationExistsStatus))
	require.NoError(t, writeRateLimitStatus(textCommand, &textOptions, rateLimitStatus))
	require.NoError(t, writeFavorite(textCommand, &textOptions, favorite))
	require.NoError(t, writeEmoji(textCommand, &textOptions, emoji))
	require.NoError(t, writeAttachment(textCommand, &textOptions, attachment))
	require.NoError(t, writeNotification(textCommand, &textOptions, notification))
	require.NoError(t, writeNotificationSubscription(textCommand, &textOptions, notificationSubscription))
	require.NoError(t, writeTriageResponsibility(textCommand, &textOptions, triageResponsibility))
	require.NoError(t, writeTriageResponsibilityManualSelection(textCommand, &textOptions, triageManualSelection))
	require.NoError(t, writeReleasePipeline(textCommand, &textOptions, releasePipeline))
	require.NoError(t, writeReleaseStage(textCommand, &textOptions, releaseStage))
	require.NoError(t, writeRelease(textCommand, &textOptions, release))
	require.NoError(t, writeReleaseHistory(textCommand, &textOptions, releaseHistory))
	require.NoError(t, writeIssueHistory(textCommand, &textOptions, issueHistory))
	require.NoError(t, writeEntityExternalLink(textCommand, &textOptions, releaseLink))
	require.NoError(t, writeReleaseNote(textCommand, &textOptions, releaseNote))
	require.Equal(
		t,
		"LIT-1 Ship coverage [Todo]\nissue-id bot bot-actor-id GitHub [github]\n"+
			"plain-issue-id bot -\nissue-state-span-id Started started 2026-06-19T12:00:00Z -> -\n"+
			"cycle-id Planning cycle [active]\n"+
			"project-id Coverage [Backlog]\nproject-update-id onTrack Omer First update\n"+
			"project-milestone-id Launch milestone [next]\n"+
			"project-status-id Backlog [backlog] #bec2c8\n"+
			"project-status-id count 12 private 2 archived_team 1\n"+
			"project-label-id Roadmap #f2c94c\n"+
			"project-relation-id blocks Pinned project -> Related project\n"+
			"issue-relation-id blocks LIT-1 -> LIT-2\n"+
			"issue-to-release-id issue issue-id -> release release-id\n"+
			"document-id Spec [project]\nlabel-id Bug #ff0000\nteam-id LIT linctl\n"+
			"team-membership-id LIT Omer owner true order 1.50\n"+
			"git-automation-state-id review state Started target main\n"+
			"user-id Omer <omer@example.com>\ndraft-id issue LIT-3 Draft issue\n"+
			"comment-id Omer First comment\ncomment-id Omer 2026-06-19T12:00:00Z\n"+
			"comment-id bot bot-actor-id GitHub [github]\nplain-comment-id bot -\n"+
			"workflow-state-id Started [started]\n"+
			"time-schedule-id Primary on-call entries 1\n"+
			"template-id Bug report [issue] team LIT\n"+
			"initiative-id Platform [Active]\ninitiative-history-id initiative initiative-id entries 1\n"+
			"initiative-relation-id Platform -> Child initiative order 1.50\n"+
			"initiative-to-project-id Platform -> Pinned project order 1\n"+
			"roadmap-to-project-id Platform roadmap -> Pinned project order 1 [legacy]\n"+
			"initiative-update-id onTrack Omer First initiative update\n"+
			"roadmap-id Platform roadmap platform-roadmap [legacy]\n"+
			"custom-view-id My issues [Issue]\n"+
			"custom-view-id has_subscribers true\n"+
			"custom-view-id organization preferences organization customView layout list\n"+
			"custom-view-id preference values layout board ordering updatedAt\n"+
			"empty-custom-view-id organization preferences -\n"+
			"sla-configuration-id First response sla 3600000 type all removes false\n"+
			"sla-remove-id Remove SLA sla - type - removes true\n"+
			"customer-id Acme [Active] needs 3\n"+
			"customer-need-id Acme LIT-1 priority 1\n"+
			"customer-status-id Active #00ff00 1\n"+
			"customer-tier-id Enterprise #0000ff 2\n"+
			"app-id Demo App by Kyanite\n"+
			"agent-activity-id session agent-session-id [action] signal continue\n"+
			"agent-skill-id Triage Helper shared true recent 3\n"+
			"external-user-id External User @external last_seen 2026-06-19T12:00:00Z\n"+
			"user_login User logged in\n"+
			"kyanite exists true success true\n"+
			"api api-key\ncomplexity remaining 900/1000 reset 1720000000000\n"+
			"favorite-id [issue] https://linear.app/kyanite/issue/LIT-1\nemoji-id party [custom]\n"+
			"attachment-id Linked PR [github]\n"+
			"notification-id issueMention [mentions] Mentioned you\n"+
			"notification-subscription-id project Roadmap active true\n"+
			"triage-responsibility-id team LIT action notify current Omer\n"+
			"triage-responsibility-id manual users user-id,other-user-id\n"+
			"release-pipeline-id Production production releases 4\n"+
			"release-stage-id Started [started] pipeline Production\n"+
			"release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3\n"+
			"release-history-id release release-id entries 1\n"+
			"issue-history-id issue issue-id updated_description true\n"+
			"release-link-id Runbook https://example.com/runbook order 1.5\n"+
			"release-note-id Launch notes pipeline Production releases 2\n",
		textOut.String(),
	)

	jsonOut := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOut)
	jsonOptions := rootOptions{json: true}

	require.NoError(t, writeIssue(jsonCommand, &jsonOptions, issue))
	require.NoError(t, writeIssueBotActor(jsonCommand, &jsonOptions, issueBotActor))
	require.NoError(t, writeIssueStateSpan(jsonCommand, &jsonOptions, issueStateSpan))
	require.NoError(t, writeCycle(jsonCommand, &jsonOptions, cycle))
	require.NoError(t, writeProject(jsonCommand, &jsonOptions, project))
	require.NoError(t, writeProjectUpdate(jsonCommand, &jsonOptions, projectUpdate))
	require.NoError(t, writeProjectMilestone(jsonCommand, &jsonOptions, milestone))
	require.NoError(t, writeProjectStatus(jsonCommand, &jsonOptions, projectStatus))
	require.NoError(t, writeProjectStatusProjectCount(jsonCommand, &jsonOptions, projectStatusProjectCount))
	require.NoError(t, writeProjectLabel(jsonCommand, &jsonOptions, projectLabel))
	require.NoError(t, writeProjectRelation(jsonCommand, &jsonOptions, projectRelation))
	require.NoError(t, writeIssueRelation(jsonCommand, &jsonOptions, issueRelation))
	require.NoError(t, writeIssueToRelease(jsonCommand, &jsonOptions, issueToRelease))
	require.NoError(t, writeDocument(jsonCommand, &jsonOptions, document))
	require.NoError(t, writeLabel(jsonCommand, &jsonOptions, label))
	require.NoError(t, writeTeam(jsonCommand, &jsonOptions, team))
	require.NoError(t, writeGitAutomationState(jsonCommand, &jsonOptions, gitAutomationState))
	require.NoError(t, writeUser(jsonCommand, &jsonOptions, user))
	require.NoError(t, writeDraft(jsonCommand, &jsonOptions, draft))
	require.NoError(t, writeComment(jsonCommand, &jsonOptions, comment))
	require.NoError(t, writeCommentMetadata(jsonCommand, &jsonOptions, commentMetadata))
	require.NoError(t, writeCommentBotActor(jsonCommand, &jsonOptions, commentBotActor))
	require.NoError(t, writeWorkflowState(jsonCommand, &jsonOptions, workflowState))
	require.NoError(t, writeTimeSchedule(jsonCommand, &jsonOptions, timeSchedule))
	require.NoError(t, writeTemplate(jsonCommand, &jsonOptions, template))
	require.NoError(t, writeInitiative(jsonCommand, &jsonOptions, initiative))
	require.NoError(t, writeInitiativeHistory(jsonCommand, &jsonOptions, initiativeHistory))
	require.NoError(t, writeInitiativeRelation(jsonCommand, &jsonOptions, initiativeRelation))
	require.NoError(t, writeInitiativeToProject(jsonCommand, &jsonOptions, initiativeToProject))
	require.NoError(t, writeRoadmapToProject(jsonCommand, &jsonOptions, roadmapToProject))
	require.NoError(t, writeInitiativeUpdate(jsonCommand, &jsonOptions, initiativeUpdate))
	require.NoError(t, writeRoadmap(jsonCommand, &jsonOptions, roadmap))
	require.NoError(t, writeCustomView(jsonCommand, &jsonOptions, customView))
	require.NoError(t, writeCustomViewSubscriberStatus(jsonCommand, &jsonOptions, customViewSubscriberStatus))
	require.NoError(t, writeCustomViewPreferences(jsonCommand, &jsonOptions, customViewPreferences))
	require.NoError(t, writeCustomViewPreferenceValues(jsonCommand, &jsonOptions, customViewPreferenceValues))
	require.NoError(t, writeSLAConfiguration(jsonCommand, &jsonOptions, slaConfiguration))
	require.NoError(t, writeCustomer(jsonCommand, &jsonOptions, customer))
	require.NoError(t, writeCustomerNeed(jsonCommand, &jsonOptions, customerNeed))
	require.NoError(t, writeCustomerStatus(jsonCommand, &jsonOptions, customerStatus))
	require.NoError(t, writeCustomerTier(jsonCommand, &jsonOptions, customerTier))
	require.NoError(t, writeApplicationInfo(jsonCommand, &jsonOptions, application))
	require.NoError(t, writeAgentActivity(jsonCommand, &jsonOptions, agentActivity))
	require.NoError(t, writeAgentSkill(jsonCommand, &jsonOptions, agentSkill))
	require.NoError(t, writeExternalUser(jsonCommand, &jsonOptions, externalUser))
	require.NoError(t, writeAuditEntryType(jsonCommand, &jsonOptions, auditEntryType))
	require.NoError(t, writeOrganizationExists(jsonCommand, &jsonOptions, organizationExistsStatus))
	require.NoError(t, writeRateLimitStatus(jsonCommand, &jsonOptions, rateLimitStatus))
	require.NoError(t, writeFavorite(jsonCommand, &jsonOptions, favorite))
	require.NoError(t, writeEmoji(jsonCommand, &jsonOptions, emoji))
	require.NoError(t, writeAttachment(jsonCommand, &jsonOptions, attachment))
	require.NoError(t, writeNotification(jsonCommand, &jsonOptions, notification))
	require.NoError(t, writeNotificationSubscription(jsonCommand, &jsonOptions, notificationSubscription))
	require.NoError(t, writeTriageResponsibility(jsonCommand, &jsonOptions, triageResponsibility))
	require.NoError(t, writeTriageResponsibilityManualSelection(jsonCommand, &jsonOptions, triageManualSelection))
	require.NoError(t, writeSemanticSearchResult(jsonCommand, &jsonOptions, semanticSearchResult))
	require.NoError(t, writeReleasePipeline(jsonCommand, &jsonOptions, releasePipeline))
	require.NoError(t, writeReleaseStage(jsonCommand, &jsonOptions, releaseStage))
	require.NoError(t, writeRelease(jsonCommand, &jsonOptions, release))
	require.NoError(t, writeReleaseHistory(jsonCommand, &jsonOptions, releaseHistory))
	require.NoError(t, writeIssueHistory(jsonCommand, &jsonOptions, issueHistory))
	require.NoError(t, writeEntityExternalLink(jsonCommand, &jsonOptions, releaseLink))
	require.NoError(t, writeReleaseNote(jsonCommand, &jsonOptions, releaseNote))
	require.Contains(t, jsonOut.String(), `"identifier": "LIT-1"`)
	require.Contains(t, jsonOut.String(), `"issue_id": "issue-id"`)
	require.Contains(t, jsonOut.String(), `"state_name": "Started"`)
	require.Contains(t, jsonOut.String(), `"name": "Planning cycle"`)
	require.Contains(t, jsonOut.String(), `"name": "Coverage"`)
	require.Contains(t, jsonOut.String(), `"body": "First update"`)
	require.Contains(t, jsonOut.String(), `"name": "Launch milestone"`)
	require.Contains(t, jsonOut.String(), `"id": "project-status-id"`)
	require.Contains(t, jsonOut.String(), `"id": "project-label-id"`)
	require.Contains(t, jsonOut.String(), `"key": "LIT-3"`)
	require.Contains(t, jsonOut.String(), `"title": "Spec"`)
	require.Contains(t, jsonOut.String(), `"color": "#ff0000"`)
	require.Contains(t, jsonOut.String(), `"key": "LIT"`)
	require.Contains(t, jsonOut.String(), `"target_branch_pattern": "main"`)
	require.Contains(t, jsonOut.String(), `"email": "omer@example.com"`)
	require.Contains(t, jsonOut.String(), `"parent_key": "LIT-3"`)
	require.Contains(t, jsonOut.String(), `"body": "First comment"`)
	require.Contains(t, jsonOut.String(), `"comment_id": "comment-id"`)
	require.Contains(t, jsonOut.String(), `"type": "started"`)
	require.Contains(t, jsonOut.String(), `"entry_count": 1`)
	require.Contains(t, jsonOut.String(), `"team_key": "LIT"`)
	require.Contains(t, jsonOut.String(), `"status": "Active"`)
	require.Contains(t, jsonOut.String(), `"related_initiative_name": "Child initiative"`)
	require.Contains(t, jsonOut.String(), `"project_name": "Pinned project"`)
	require.Contains(t, jsonOut.String(), `"body": "First initiative update"`)
	require.Contains(t, jsonOut.String(), `"slug_id": "platform-roadmap"`)
	require.Contains(t, jsonOut.String(), `"model_name": "Issue"`)
	require.Contains(t, jsonOut.String(), `"has_subscribers": true`)
	require.Contains(t, jsonOut.String(), `"view_ordering": "updatedAt"`)
	require.Contains(t, jsonOut.String(), `"hidden_columns": [`)
	require.Contains(t, jsonOut.String(), `"sla_type": "all"`)
	require.Contains(t, jsonOut.String(), `"approximate_need_count": 3`)
	require.Contains(t, jsonOut.String(), `"customer_name": "Acme"`)
	require.Contains(t, jsonOut.String(), `"display_name": "Active"`)
	require.Contains(t, jsonOut.String(), `"display_name": "Enterprise"`)
	require.Contains(t, jsonOut.String(), `"client_id": "app-client-id"`)
	require.Contains(t, jsonOut.String(), `"content_type": "action"`)
	require.Contains(t, jsonOut.String(), `"title": "Triage Helper"`)
	require.Contains(t, jsonOut.String(), `"display_name": "@external"`)
	require.Contains(t, jsonOut.String(), `"description": "User logged in"`)
	require.Contains(t, jsonOut.String(), `"url_key": "kyanite"`)
	require.Contains(t, jsonOut.String(), `"remaining_amount": 900`)
	require.Contains(t, jsonOut.String(), `"type": "issue"`)
	require.Contains(t, jsonOut.String(), `"source": "custom"`)
	require.Contains(t, jsonOut.String(), `"source_type": "github"`)
	require.Contains(t, jsonOut.String(), `"category": "mentions"`)
	require.Contains(t, jsonOut.String(), `"target_type": "project"`)
	require.Contains(t, jsonOut.String(), `"team_key": "LIT"`)
	require.Contains(t, jsonOut.String(), `"user_ids": [`)
	require.Contains(t, jsonOut.String(), `"slug_id": "production"`)
	require.Contains(t, jsonOut.String(), `"pipeline_name": "Production"`)
	require.Contains(t, jsonOut.String(), `"version": "v1.2.3"`)
	require.Contains(t, jsonOut.String(), `"release_count": 2`)
}
