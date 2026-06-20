package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CliRenderHelpers_write_text_and_json_output(t *testing.T) {
	issue := client.IssueSummary{
		Identifier: "LIT-1",
		Title:      "Ship coverage",
		State:      "Todo",
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
	require.NoError(t, writeCycle(textCommand, &textOptions, cycle))
	require.NoError(t, writeProject(textCommand, &textOptions, project))
	require.NoError(t, writeProjectUpdate(textCommand, &textOptions, projectUpdate))
	require.NoError(t, writeProjectMilestone(textCommand, &textOptions, milestone))
	require.NoError(t, writeProjectStatus(textCommand, &textOptions, projectStatus))
	require.NoError(t, writeProjectLabel(textCommand, &textOptions, projectLabel))
	require.NoError(t, writeProjectRelation(textCommand, &textOptions, projectRelation))
	require.NoError(t, writeDocument(textCommand, &textOptions, document))
	require.NoError(t, writeLabel(textCommand, &textOptions, label))
	require.NoError(t, writeTeam(textCommand, &textOptions, team))
	require.NoError(t, writeTeamMembership(textCommand, &textOptions, teamMembership))
	require.NoError(t, writeUser(textCommand, &textOptions, user))
	require.NoError(t, writeDraft(textCommand, &textOptions, draft))
	require.NoError(t, writeComment(textCommand, &textOptions, comment))
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
	require.NoError(t, writeEntityExternalLink(textCommand, &textOptions, releaseLink))
	require.NoError(t, writeReleaseNote(textCommand, &textOptions, releaseNote))
	require.Equal(
		t,
		"LIT-1 Ship coverage [Todo]\ncycle-id Planning cycle [active]\n"+
			"project-id Coverage [Backlog]\nproject-update-id onTrack Omer First update\n"+
			"project-milestone-id Launch milestone [next]\n"+
			"project-status-id Backlog [backlog] #bec2c8\n"+
			"project-label-id Roadmap #f2c94c\n"+
			"project-relation-id blocks Pinned project -> Related project\n"+
			"document-id Spec [project]\nlabel-id Bug #ff0000\nteam-id LIT linctl\n"+
			"team-membership-id LIT Omer owner true order 1.50\n"+
			"user-id Omer <omer@example.com>\ndraft-id issue LIT-3 Draft issue\n"+
			"comment-id Omer First comment\nworkflow-state-id Started [started]\n"+
			"time-schedule-id Primary on-call entries 1\n"+
			"template-id Bug report [issue] team LIT\n"+
			"initiative-id Platform [Active]\ninitiative-history-id initiative initiative-id entries 1\n"+
			"initiative-relation-id Platform -> Child initiative order 1.50\n"+
			"initiative-to-project-id Platform -> Pinned project order 1\n"+
			"roadmap-to-project-id Platform roadmap -> Pinned project order 1\n"+
			"initiative-update-id onTrack Omer First initiative update\n"+
			"roadmap-id Platform roadmap platform-roadmap\n"+
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
			"release-link-id Runbook https://example.com/runbook order 1.5\n"+
			"release-note-id Launch notes pipeline Production releases 2\n",
		textOut.String(),
	)

	jsonOut := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOut)
	jsonOptions := rootOptions{json: true}

	require.NoError(t, writeIssue(jsonCommand, &jsonOptions, issue))
	require.NoError(t, writeCycle(jsonCommand, &jsonOptions, cycle))
	require.NoError(t, writeProject(jsonCommand, &jsonOptions, project))
	require.NoError(t, writeProjectUpdate(jsonCommand, &jsonOptions, projectUpdate))
	require.NoError(t, writeProjectMilestone(jsonCommand, &jsonOptions, milestone))
	require.NoError(t, writeProjectStatus(jsonCommand, &jsonOptions, projectStatus))
	require.NoError(t, writeProjectLabel(jsonCommand, &jsonOptions, projectLabel))
	require.NoError(t, writeProjectRelation(jsonCommand, &jsonOptions, projectRelation))
	require.NoError(t, writeDocument(jsonCommand, &jsonOptions, document))
	require.NoError(t, writeLabel(jsonCommand, &jsonOptions, label))
	require.NoError(t, writeTeam(jsonCommand, &jsonOptions, team))
	require.NoError(t, writeUser(jsonCommand, &jsonOptions, user))
	require.NoError(t, writeDraft(jsonCommand, &jsonOptions, draft))
	require.NoError(t, writeComment(jsonCommand, &jsonOptions, comment))
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
	require.NoError(t, writeEntityExternalLink(jsonCommand, &jsonOptions, releaseLink))
	require.NoError(t, writeReleaseNote(jsonCommand, &jsonOptions, releaseNote))
	require.Contains(t, jsonOut.String(), `"identifier": "LIT-1"`)
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
	require.Contains(t, jsonOut.String(), `"email": "omer@example.com"`)
	require.Contains(t, jsonOut.String(), `"parent_key": "LIT-3"`)
	require.Contains(t, jsonOut.String(), `"body": "First comment"`)
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

	require.NoError(t, writeIssue(command, &rootOptions{format: "full"}, issue))
	require.NoError(t, writeIssue(command, &rootOptions{idOnly: true}, issue))
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
	require.NoError(t, writeTeamMembership(command, &rootOptions{idOnly: true}, teamMembership))
	require.NoError(t, writeDocument(command, &rootOptions{idOnly: true}, document))
	require.NoError(t, writeLabel(command, &rootOptions{idOnly: true}, label))
	require.NoError(t, writeTeam(command, &rootOptions{idOnly: true}, team))
	require.NoError(t, writeUser(command, &rootOptions{idOnly: true}, user))
	require.NoError(t, writeComment(command, &rootOptions{idOnly: true}, comment))
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
	require.NoError(t, writeCycle(quietCommand, &rootOptions{quiet: true}, cycle))
	require.NoError(t, writeProject(quietCommand, &rootOptions{quiet: true}, project))
	require.NoError(t, writeProjectUpdate(quietCommand, &rootOptions{quiet: true}, projectUpdate))
	require.NoError(t, writeProjectMilestone(quietCommand, &rootOptions{quiet: true}, milestone))
	require.NoError(t, writeProjectStatus(quietCommand, &rootOptions{quiet: true}, projectStatus))
	require.NoError(t, writeProjectLabel(quietCommand, &rootOptions{quiet: true}, projectLabel))
	require.NoError(t, writeProjectRelation(quietCommand, &rootOptions{quiet: true}, projectRelation))
	require.NoError(t, writeTeamMembership(quietCommand, &rootOptions{quiet: true}, teamMembership))
	require.NoError(t, writeDocument(quietCommand, &rootOptions{quiet: true}, document))
	require.NoError(t, writeLabel(quietCommand, &rootOptions{quiet: true}, label))
	require.NoError(t, writeTeam(quietCommand, &rootOptions{quiet: true}, team))
	require.NoError(t, writeUser(quietCommand, &rootOptions{quiet: true}, user))
	require.NoError(t, writeComment(quietCommand, &rootOptions{quiet: true}, comment))
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

func Test_CliOutputHelpers_cover_json_projection_and_sort_edges(t *testing.T) {
	projected, err := projectJSONFields(
		map[string]any{"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}}},
		"identifier,state.name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"projects": []any{map[string]any{"id": "project-id", "status": map[string]any{"name": "Backlog"}}}},
		"id,status.name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"projects": []any{map[string]any{"id": "project-id", "status": map[string]any{"name": "Backlog"}}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"members": []any{map[string]any{"id": "user-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"members": []any{map[string]any{"id": "user-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customers": []any{map[string]any{"id": "customer-id", "status_name": "Active"}}},
		"id,status_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customers": []any{map[string]any{"id": "customer-id", "status_name": "Active"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"roadmaps": []any{map[string]any{"id": "roadmap-id", "slug_id": "platform-roadmap"}}},
		"id,slug_id",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"roadmaps": []any{map[string]any{"id": "roadmap-id", "slug_id": "platform-roadmap"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"time_schedules": []any{map[string]any{"id": "time-schedule-id", "entry_count": float64(1)}}},
		"id,entry_count",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"time_schedules": []any{map[string]any{"id": "time-schedule-id", "entry_count": float64(1)}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"notifications": []any{map[string]any{"id": "notification-id", "category": "mentions"}}},
		"id,category",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"notifications": []any{map[string]any{"id": "notification-id", "category": "mentions"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{
			"notification_subscriptions": []any{
				map[string]any{"id": "notification-subscription-id", "target_type": "project"},
			},
		},
		"id,target_type",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"notification_subscriptions": []any{
			map[string]any{"id": "notification-subscription-id", "target_type": "project"},
		},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"triage_responsibilities": []any{
			map[string]any{"id": "triage-responsibility-id", "team_key": "LIT", "action": "notify"},
		}},
		"id,team_key,action",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"triage_responsibilities": []any{
			map[string]any{"id": "triage-responsibility-id", "team_key": "LIT", "action": "notify"},
		},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"sla_configurations": []any{
			map[string]any{"id": "sla-configuration-id", "name": "First response", "sla_type": "all"},
		}},
		"id,name,sla_type",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"sla_configurations": []any{
			map[string]any{"id": "sla-configuration-id", "name": "First response", "sla_type": "all"},
		},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"audit_entry_types": []any{map[string]any{"type": "user_login", "description": "User logged in"}}},
		"type,description",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"audit_entry_types": []any{map[string]any{"type": "user_login", "description": "User logged in"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"release_pipelines": []any{map[string]any{"id": "release-pipeline-id", "slug_id": "production"}}},
		"id,slug_id",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"release_pipelines": []any{map[string]any{"id": "release-pipeline-id", "slug_id": "production"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"release_stages": []any{map[string]any{"id": "release-stage-id", "pipeline_name": "Production"}}},
		"id,pipeline_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"release_stages": []any{map[string]any{"id": "release-stage-id", "pipeline_name": "Production"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customer_needs": []any{map[string]any{"id": "customer-need-id", "customer_name": "Acme"}}},
		"id,customer_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customer_needs": []any{map[string]any{"id": "customer-need-id", "customer_name": "Acme"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customer_statuses": []any{map[string]any{"id": "customer-status-id", "display_name": "Active"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customer_statuses": []any{map[string]any{"id": "customer-status-id", "display_name": "Active"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"customer_tiers": []any{map[string]any{"id": "customer-tier-id", "display_name": "Enterprise"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"customer_tiers": []any{map[string]any{"id": "customer-tier-id", "display_name": "Enterprise"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"comments": []any{map[string]any{"id": "comment-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"comments": []any{map[string]any{"id": "comment-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"documents": []any{map[string]any{"id": "document-id", "title": "Spec"}}},
		"id,title",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"documents": []any{map[string]any{"id": "document-id", "title": "Spec"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"labels": []any{map[string]any{"id": "label-id", "color": "#ff0000"}}},
		"id,color",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"labels": []any{map[string]any{"id": "label-id", "color": "#ff0000"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"teams": []any{map[string]any{"id": "team-id", "key": "LIT"}}},
		"id,key",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"teams": []any{map[string]any{"id": "team-id", "key": "LIT"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"users": []any{map[string]any{"id": "user-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"users": []any{map[string]any{"id": "user-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"drafts": []any{map[string]any{"id": "draft-id", "parent_key": "LIT-3"}}},
		"id,parent_key",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"drafts": []any{map[string]any{"id": "draft-id", "parent_key": "LIT-3"}},
	}, projected)

	projected, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "identifier")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	projected, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "identifier,, ")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	_, err = projectJSONFields(map[string]any{"bad": func() {}}, "bad")
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal output")

	_, err = projectJSONFields([]string{"not-an-object"}, "id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode output")

	_, err = projectJSONFields(map[string]any{"issues": []any{"bad-item"}}, "identifier")
	require.Error(t, err)
	require.Contains(t, err.Error(), "item is not an object")

	_, err = projectJSONFields(map[string]any{"issues": []any{map[string]any{"title": "Missing id"}}}, "identifier")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"identifier\" is not present")

	_, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "missing")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"missing\" is not present")

	_, err = projectJSONFields(map[string]any{"state": "Todo"}, "state.name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"state\" is not an object")

	items := []client.IssueSummary{
		{Identifier: "LIT-2", Title: "Zebra"},
		{Identifier: "LIT-1", Title: "Alpha"},
	}
	sortedItems, err := sortByJSONField(items, "", "asc")
	require.NoError(t, err)
	require.Equal(t, items, sortedItems)

	sortedItems, err = sortByJSONField(items, "title", "asc")
	require.NoError(t, err)
	require.Equal(t, "Alpha", sortedItems[0].Title)

	_, err = sortByJSONField(items, "title", "sideways")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid sort order")

	_, err = sortByJSONField(items, "missing", "asc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "sort field \"missing\" is not present")

	_, err = sortByJSONField([]map[string]any{{"state": "Todo"}}, "state.name", "asc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not an object path")

	_, err = jsonFieldValue(map[string]any{"bad": func() {}}, "bad")
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal output")

	destination := map[string]any{}
	require.NoError(t, copyJSONPath(map[string]any{"id": "issue-id"}, destination, nil))
	require.Empty(t, destination)
}

func Test_CommandFlows_cover_output_error_and_quiet_branches(t *testing.T) {
	quietCommands := [][]string{
		{"--quiet", "target"},
		{"--quiet", "whoami"},
		{"--quiet", "issue", "deps", "LIT-1"},
		{"--quiet", "issue", "pr", "LIT-1"},
		{"--quiet", "usage"},
	}
	for _, args := range quietCommands {
		t.Run("quiet "+args[len(args)-1], func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Empty(t, output.String())
		})
	}

	errorCommands := [][]string{
		{"--sort", "missing", "issue", "list"},
		{"--sort", "missing", "issue", "list", "--project", "project-id"},
		{"--sort", "missing", "issue", "list", "--mine"},
		{"--sort", "missing", "issue", "list", "--assignee", "assignee-id"},
		{"--sort", "missing", "issue", "list", "--label", "label-id"},
		{"--sort", "missing", "issue", "list", "--cycle", "cycle-id"},
		{"--sort", "missing", "issue", "list", "--created-after", "2026-06-01"},
		{"--sort", "missing", "issue", "list", "--created-since", "2026-06-01"},
		{"--sort", "missing", "issue", "list", "--created-before", "2026-06-30"},
		{"--sort", "missing", "issue", "list", "--has-blockers"},
		{"--sort", "missing", "issue", "list", "--blocks"},
		{"--sort", "missing", "issue", "list", "--blocked-by", "LIT-1"},
		{"--sort", "missing", "issue", "list", "--all-teams"},
		{"--sort", "missing", "issue", "comments", "LIT-1"},
		{"--sort", "missing", "issue", "search", "needle"},
		{"--sort", "missing", "project", "list"},
		{"--sort", "missing", "project", "members", "project-id"},
	}
	for _, args := range errorCommands {
		t.Run("sort error "+args[len(args)-1], func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "sort field")
		})
	}

	emptyCommands := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{name: "issue list project", args: []string{"--fail-on-empty", "issue", "list", "--project", "project-id"}, fake: commandFlowFakeClient{emptyIssueProject: true}},
		{name: "issue list mine", args: []string{"--fail-on-empty", "issue", "list", "--mine"}, fake: commandFlowFakeClient{emptyIssueMine: true}},
		{name: "issue list assignee", args: []string{"--fail-on-empty", "issue", "list", "--assignee", "assignee-id"}, fake: commandFlowFakeClient{emptyIssueMine: true}},
		{name: "issue list label", args: []string{"--fail-on-empty", "issue", "list", "--label", "label-id"}, fake: commandFlowFakeClient{emptyIssueLabel: true}},
		{name: "issue list cycle", args: []string{"--fail-on-empty", "issue", "list", "--cycle", "cycle-id"}, fake: commandFlowFakeClient{emptyIssueCycle: true}},
		{name: "issue list created-after", args: []string{"--fail-on-empty", "issue", "list", "--created-after", "2026-06-01"}, fake: commandFlowFakeClient{emptyIssueCreatedAfter: true}},
		{name: "issue list created-since", args: []string{"--fail-on-empty", "issue", "list", "--created-since", "2026-06-01"}, fake: commandFlowFakeClient{emptyIssueCreatedAfter: true}},
		{name: "issue list created-before", args: []string{"--fail-on-empty", "issue", "list", "--created-before", "2026-06-30"}, fake: commandFlowFakeClient{emptyIssueCreatedBefore: true}},
		{name: "issue list has blockers", args: []string{"--fail-on-empty", "issue", "list", "--has-blockers"}, fake: commandFlowFakeClient{emptyIssueHasBlockers: true}},
		{name: "issue list blocks", args: []string{"--fail-on-empty", "issue", "list", "--blocks"}, fake: commandFlowFakeClient{emptyIssueBlocks: true}},
		{name: "issue list blocked by", args: []string{"--fail-on-empty", "issue", "list", "--blocked-by", "LIT-1"}, fake: commandFlowFakeClient{emptyIssueBlockedBy: true}},
		{name: "issue list all teams", args: []string{"--fail-on-empty", "issue", "list", "--all-teams"}, fake: commandFlowFakeClient{emptyIssueAllTeams: true}},
		{name: "issue comments", args: []string{"--fail-on-empty", "issue", "comments", "LIT-1"}, fake: commandFlowFakeClient{emptyIssueComments: true}},
		{name: "issue search", args: []string{"--fail-on-empty", "issue", "search", "needle"}, fake: commandFlowFakeClient{emptyIssueSearch: true}},
		{name: "project list", args: []string{"--fail-on-empty", "project", "list"}, fake: commandFlowFakeClient{emptyProjectList: true}},
		{name: "project members", args: []string{"--fail-on-empty", "project", "members", "project-id"}, fake: commandFlowFakeClient{emptyProjectMembers: true}},
	}
	for _, test := range emptyCommands {
		t.Run("empty "+test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "empty result")
		})
	}
}

func Test_CommandFlows_cover_issue_list_filter_validation(t *testing.T) {
	tests := [][]string{
		{"issue", "list", "--state", "started", "--project", "project-id"},
		{"issue", "list", "--state", "started", "--mine"},
		{"issue", "list", "--state", "started", "--assignee", "assignee-id"},
		{"issue", "list", "--state", "started", "--label", "label-id"},
		{"issue", "list", "--state", "started", "--cycle", "cycle-id"},
		{"issue", "list", "--state", "started", "--created-after", "2026-06-01"},
		{"issue", "list", "--state", "started", "--created-since", "2026-06-01"},
		{"issue", "list", "--created-after", "2026-06-01", "--created-since", "2026-06-01"},
		{"issue", "list", "--state", "started", "--created-before", "2026-06-30"},
		{"issue", "list", "--state", "started", "--has-blockers"},
		{"issue", "list", "--has-blockers", "--blocks"},
		{"issue", "list", "--blocks", "--blocked-by", "LIT-1"},
		{"issue", "list", "--state", "started", "--all-teams"},
	}
	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "use only one")
		})
	}
}

func Test_CommandFlows_cover_issue_current_error_branches(t *testing.T) {
	t.Run("id missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("title missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "title"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("title runtime error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "title"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			},
		)

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("url lookup error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "url"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return testCommandRuntime(commandFlowFakeClient{failOperation: "issue"}), nil
			},
		)

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})

	t.Run("branch argument lookup error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "branch", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})
}

func Test_CommandFlows_cover_issue_comment_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_issue_reply_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "reply", "LIT-1", "comment-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_issue_comments_error_branches(t *testing.T) {
	t.Run("operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_comments"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "comments", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue comments LIT-1")
	})

	t.Run("writer error", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssueComments(command, []client.IssueCommentSummary{{ID: "comment-id", DisplayName: "Omer", Body: "body"}})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_issue_deps_writer_error(t *testing.T) {
	t.Run("issue header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("section header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssueDependencySection(command, &rootOptions{}, "children", nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("parent issue", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 3})
		parent := client.IssueSummary{Identifier: "LIT-2", Title: "Parent", State: "Todo"}
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1", Parent: &parent}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("children section", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 2})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("blocks section", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 3})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_rate_limit_writer_errors(t *testing.T) {
	status := client.RateLimitStatus{
		Identifier: "api-key",
		Kind:       "api",
		Limits: []client.RateLimit{
			{Type: "complexity", AllowedAmount: 1000, RemainingAmount: 900, Reset: 1720000000000},
		},
	}

	t.Run("header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeRateLimitStatus(command, &rootOptions{}, status)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("limit", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 2})

		err := writeRateLimitStatus(command, &rootOptions{}, status)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_audit_entry_type_writer_errors(t *testing.T) {
	command := &cobra.Command{}
	command.SetOut(commandFailingWriter{})
	types := client.AuditEntryTypeList{
		AuditEntryTypes: []client.AuditEntryTypeSummary{
			{Type: "user_login", Description: "User logged in"},
		},
	}

	err := writeAuditEntryTypes(command, &rootOptions{}, types)

	require.Error(t, err)
	require.Contains(t, err.Error(), "write line")
}

type countFailingWriter struct {
	failAt int
	writes int
}

func (writer *countFailingWriter) Write(content []byte) (int, error) {
	writer.writes++
	if writer.writes == writer.failAt {
		return 0, errors.New("write failed")
	}

	return len(content), nil
}

func Test_CliHelpers_resolve_target_overrides_and_project_ids(t *testing.T) {
	options := rootOptions{
		orgID:   "org-id",
		team:    "LIT",
		project: "project-id",
	}

	target := targetOverride(&options)

	require.Equal(t, "org-id", target.OrgID)
	require.Equal(t, "LIT", target.TeamKey)
	require.Equal(t, "project-id", target.ProjectID)
	require.Empty(t, projectID(nil))
	require.Equal(t, "project-id", projectID(&client.ResolvedProject{ID: "project-id"}))
	require.NotEmpty(t, defaultGlobalConfigPath())
}

type commandFailingReader struct{}

func (reader commandFailingReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func Test_CommandRuntime_loads_config_and_requires_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LINCTL_TOKEN", "")
	t.Setenv("LINEAR_API_KEY", "")
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing Linear token")

	t.Setenv("LINCTL_TOKEN", "test-token")
	runtime, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.NoError(t, err)
	require.Equal(t, "project-id", runtime.config.Target.ProjectID)
	require.NotNil(t, runtime.graphqlClient)
}

func Test_CommandRuntime_reports_config_load_errors(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	require.NoError(t, os.WriteFile(".linctl.toml", []byte("[target\n"), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "parse config")
}

func Test_DefaultGlobalConfigPath_returns_empty_when_home_is_unset(t *testing.T) {
	t.Setenv("HOME", "")

	require.Empty(t, defaultGlobalConfigPath())
}

func Test_WriteUsage_reports_unknown_topics(t *testing.T) {
	command := &cobra.Command{}

	err := writeUsage(command, &rootOptions{}, "missing")

	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown usage topic "missing"`)
}
