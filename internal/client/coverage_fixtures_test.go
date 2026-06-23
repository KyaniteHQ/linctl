package client

import (
	"strconv"
	"strings"
)

func issueJSONWithAssignee(issue issueFixture, assignee string) string {
	return strings.ReplaceAll(issueJSON(issue), `"assignee":null`, `"assignee":{"id":"user-id","name":"omer","displayName":"`+assignee+`"}`)
}

func issueJSONWithDescription(issue issueFixture, description string) string {
	return strings.Replace(issueJSON(issue), `"id":"issue-id",`, `"id":"issue-id","description":"`+description+`",`, 1)
}

func nextIssueJSON(
	identifier string,
	title string,
	priority int,
	priorityLabel string,
	createdAt string,
	relations []string,
) string {
	return strings.TrimSuffix(issueJSON(issueFixture{
		Identifier: identifier,
		Title:      title,
		StateID:    "todo",
		State:      "Todo",
		StateType:  "unstarted",
	}), "\n\t}") +
		`,
		"priority":` + strconv.Itoa(priority) + `,
		"priorityLabel":"` + priorityLabel + `",
		"createdAt":"` + createdAt + `",
		"relations":{"nodes":[` + strings.Join(relations, ",") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}
	}`
}

func issueIdentifiers(issues []IssueSummary) []string {
	identifiers := make([]string, 0, len(issues))
	for _, issue := range issues {
		identifiers = append(identifiers, issue.Identifier)
	}

	return identifiers
}

func projectJSONWithLead(project projectFixture, lead string) string {
	return strings.Replace(projectJSON(project), `"lead":null`, `"lead":{"id":"user-id","name":"omer","displayName":"`+lead+`"}`, 1)
}

func notificationJSON() string {
	return `{
		"__typename":"IssueNotification",
		"id":"notification-id",
		"type":"issueMention",
		"category":"mentions",
		"title":"Mentioned you",
		"subtitle":"LIT-1",
		"url":"https://linear.app/kyanite/issue/LIT-1",
		"inboxUrl":"https://linear.app/kyanite/inbox/notification-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"readAt":null,
		"emailedAt":null,
		"snoozedUntilAt":null,
		"unsnoozedAt":null,
		"user":{"id":"user-id","displayName":"Omer"},
		"actor":{"id":"actor-id","displayName":"Ada"},
		"externalUserActor":{"id":"external-user-id"}
	}`
}

func notificationSubscriptionJSON() string {
	return `{
		"__typename":"ProjectNotificationSubscription",
		"id":"notification-subscription-id",
		"active":true,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"contextViewType":null,
		"userContextViewType":null,
		"subscriber":{"id":"user-id","displayName":"Omer"},
		"customer":null,
		"customView":null,
		"cycle":null,
		"initiative":null,
		"label":null,
		"project":{"id":"project-id","name":"Roadmap"},
		"team":null,
		"user":null
	}`
}

func draftJSON(parentType string) string {
	payload := `{
		"id":"draft-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"issue":null,
		"project":null,
		"projectUpdate":null,
		"initiative":null,
		"initiativeUpdate":null,
		"parentComment":null,
		"customerNeed":null,
		"team":null
	}`
	switch parentType {
	case "issue":
		return strings.Replace(payload, `"issue":null`, `"issue":{"id":"issue-id","identifier":"LIT-3","title":"Draft issue"}`, 1)
	case "project":
		return strings.Replace(payload, `"project":null`, `"project":{"id":"project-id","name":"Draft project"}`, 1)
	case "project_update":
		return strings.Replace(payload, `"projectUpdate":null`, `"projectUpdate":{"id":"project-update-id"}`, 1)
	case "initiative":
		return strings.Replace(payload, `"initiative":null`, `"initiative":{"id":"initiative-id","name":"Draft initiative"}`, 1)
	case "initiative_update":
		return strings.Replace(payload, `"initiativeUpdate":null`, `"initiativeUpdate":{"id":"initiative-update-id"}`, 1)
	case "comment":
		return strings.Replace(payload, `"parentComment":null`, `"parentComment":{"id":"comment-id"}`, 1)
	case "customer_need":
		return strings.Replace(payload, `"customerNeed":null`, `"customerNeed":{"id":"customer-need-id"}`, 1)
	case "team":
		return strings.Replace(payload, `"team":null`, `"team":{"id":"team-id","key":"LIT","name":"linctl"}`, 1)
	default:
		return payload
	}
}

func triageResponsibilityJSON() string {
	return `{
		"id":"triage-responsibility-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"action":"notify",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"timeSchedule":{"id":"time-schedule-id","name":"Primary rotation"},
		"currentUser":{"id":"user-id","displayName":"Omer"},
		"manualSelection":{"userIds":["user-id","other-user-id"]}
	}`
}

func notificationSubscriptionTargetJSON(
	typename string,
	targetField string,
	targetPayload string,
	withContextView bool,
	withUserContextView bool,
) string {
	contextViewType := "null"
	if withContextView {
		contextViewType = `"backlog"`
	}
	userContextViewType := "null"
	if withUserContextView {
		userContextViewType = `"assigned"`
	}

	payload := `{
		"__typename":"` + typename + `",
		"id":"notification-subscription-id",
		"active":true,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"contextViewType":` + contextViewType + `,
		"userContextViewType":` + userContextViewType + `,
		"subscriber":{"id":"user-id","displayName":"Omer"},
		"customer":null,
		"customView":null,
		"cycle":null,
		"initiative":null,
		"label":null,
		"project":null,
		"team":null,
		"user":null
	}`

	return strings.Replace(payload, `"`+targetField+`":null`, `"`+targetField+`":`+targetPayload, 1)
}

func subInitiativeJSON() string {
	return `{
		"id":"child-initiative-id",
		"name":"Child platform",
		"description":"Child initiative",
		"status":"Planned",
		"priority":1,
		"targetDate":"2026-11-30",
		"slugId":"child-platform",
		"url":"https://linear.app/kyanite/initiative/child-platform"
	}`
}

func initiativeHistoryJSON() string {
	return `{
		"id":"initiative-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Planned","to":"Active"}],
		"initiative":{"id":"initiative-id"}
	}`
}

func initiativeUpdateJSON() string {
	return `{
		"id":"initiative-update-id",
		"body":"First initiative update",
		"health":"onTrack",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"url":"https://linear.app/initiative-update/initiative-update-id",
		"slugId":"initiative-update-slug",
		"commentCount":1,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func customViewPreferencesJSON(ordering string, layout string) string {
	return customViewScopedPreferencesJSON("organization", ordering, layout)
}

func customViewScopedPreferencesJSON(scope string, ordering string, layout string) string {
	return `{
		"id":"view-preferences-id",
		"createdAt":"2026-06-01T12:00:00Z",
		"updatedAt":"2026-06-01T12:01:00Z",
		"archivedAt":null,
		"type":"` + scope + `",
		"viewType":"customView",
		"preferences":` + customViewPreferenceValuesJSON(ordering, layout) + `
	}`
}

func customViewPreferenceValuesJSON(ordering string, layout string) string {
	return `{
		"layout":"` + layout + `",
		"viewOrdering":"` + ordering + `",
		"viewOrderingDirection":"Descending",
		"issueGrouping":"status",
		"issueSubGrouping":"priority",
		"showCompletedIssues":"all",
		"showArchivedItems":true,
		"showEmptyGroups":true,
		"hiddenColumns":["column-id"],
		"hiddenRows":["row-id"],
		"hiddenGroupsList":["group-id"],
		"columnOrderBoard":["board-column-id"],
		"columnOrderList":["list-column-id"],
		"projectLayout":"timeline",
		"projectViewOrdering":"priority",
		"projectGrouping":"status",
		"projectSubGrouping":"lead",
		"projectShowEmptyGroups":"all",
		"projectShowEmptySubGroups":"all"
	}`
}

func templateJSON() string {
	return `{
		"id":"template-id",
		"name":"Bug report",
		"type":"issue",
		"description":"Bug report template",
		"icon":"bug",
		"color":"#ff0000",
		"templateData":{"title":"Bug: "},
		"sortOrder":1,
		"lastAppliedAt":"2026-06-19T12:00:00Z",
		"createdAt":"2026-06-18T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"pipeline":{"id":"pipeline-id"},
		"creator":{"id":"creator-id"},
		"lastUpdatedBy":{"id":"updated-by-id"},
		"inheritedFrom":{"id":"parent-template-id"}
	}`
}

func agentSkillJSON() string {
	return `{
		"id":"agent-skill-id",
		"title":"Triage Helper",
		"body":"Use this skill for triage.",
		"description":"Helps triage issues",
		"slugId":"triage-helper",
		"teamId":"team-id",
		"shared":true,
		"icon":"sparkles",
		"color":"#5e6ad2",
		"recentUsageCount":3,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"lastUsedAt":"2026-06-20T12:00:00Z",
		"owner":{"id":"owner-id"},
		"creator":{"id":"creator-id"},
		"lastUpdatedBy":{"id":"updater-id"}
	}`
}

func externalUserJSON() string {
	return `{
		"id":"external-user-id",
		"name":"External User",
		"displayName":"@external",
		"avatarUrl":"https://example.com/avatar.png",
		"lastSeen":"2026-06-19T12:00:00Z",
		"createdAt":"2026-06-18T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null
	}`
}

func agentActivityJSON(contentType string) string {
	content := map[string]string{
		"action": `{
			"__typename":"AgentActivityActionContent",
			"type":"action",
			"action":"read_file",
			"parameter":"README.md",
			"result":"Read file"
		}`,
		"elicitation": `{
			"__typename":"AgentActivityElicitationContent",
			"type":"elicitation",
			"body":"Need a choice"
		}`,
		"error": `{
			"__typename":"AgentActivityErrorContent",
			"type":"error",
			"body":"Tool failed",
			"reasonCode":"tool_error"
		}`,
		"prompt": `{
			"__typename":"AgentActivityPromptContent",
			"type":"prompt",
			"body":"Please continue"
		}`,
		"response": `{
			"__typename":"AgentActivityResponseContent",
			"type":"response",
			"body":"Done"
		}`,
		"thought": `{
			"__typename":"AgentActivityThoughtContent",
			"type":"thought",
			"body":"Thinking"
		}`,
	}[contentType]

	return `{
		"id":"agent-activity-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"signal":"continue",
		"ephemeral":false,
		"agentSession":{"id":"agent-session-id"},
		"sourceComment":{"id":"comment-id"},
		"user":{"id":"user-id"},
		"content":` + content + `
	}`
}

func slaConfigurationJSON() string {
	return `{
		"id":"sla-configuration-id",
		"name":"First response",
		"conditions":{"priority":{"eq":1}},
		"sla":3600000,
		"slaType":"all",
		"removesSla":false
	}`
}

func semanticSearchResultJSON(resultType string) string {
	switch resultType {
	case "issue":
		return `{
			"id":"issue-id",
			"type":"issue",
			"issue":{"id":"issue-id","identifier":"LIT-3","title":"Search result","url":"https://linear.app/kyanite/issue/LIT-3"},
			"project":null,
			"initiative":null,
			"document":null
		}`
	case "project":
		return `{
			"id":"project-id",
			"type":"project",
			"issue":null,
			"project":{"id":"project-id","name":"Search project","url":"https://linear.app/kyanite/project/search-project"},
			"initiative":null,
			"document":null
		}`
	case "initiative":
		return `{
			"id":"initiative-id",
			"type":"initiative",
			"issue":null,
			"project":null,
			"initiative":{"id":"initiative-id","name":"Search initiative","url":"https://linear.app/kyanite/initiative/search-initiative"},
			"document":null
		}`
	case "document":
		return `{
			"id":"document-id",
			"type":"document",
			"issue":null,
			"project":null,
			"initiative":null,
			"document":{"id":"document-id","title":"Search document","url":"https://linear.app/kyanite/document/search-document"}
		}`
	default:
		return `{
			"id":"unknown-id",
			"type":"issue",
			"issue":null,
			"project":null,
			"initiative":null,
			"document":null
		}`
	}
}

func searchDocumentJSON() string {
	return searchDocumentJSONWithParent(
		"search-document-id",
		"Search spec",
		`"project":null,"initiative":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"release":null,"cycle":null`,
	)
}

func searchDocumentJSONWithParent(id string, title string, parentFields string) string {
	return `{
		"id":"` + id + `",
		"title":"` + title + `",
		"slugId":"` + id + `",
		"url":"https://linear.app/kyanite/document/` + id + `",
		` + parentFields + `
	}`
}

func searchIssueJSON() string {
	return `{
		"id":"search-issue-id",
		"identifier":"LIT-30",
		"title":"Search issue",
		"url":"https://linear.app/kyanite/issue/LIT-30",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"state":{"id":"state-id","name":"Todo","type":"unstarted"},
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func searchProjectJSON() string {
	return `{
		"id":"search-project-id",
		"name":"Search project",
		"slugId":"search-project",
		"url":"https://linear.app/kyanite/project/search-project",
		"status":{"id":"status-id","name":"Backlog","type":"backlog"},
		"lead":{"id":"user-id","name":"omer","displayName":"Omer"},
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}]}
	}`
}

func releasePipelineJSON() string {
	return `{
		"id":"release-pipeline-id",
		"name":"Production",
		"slugId":"production",
		"type":"scheduled",
		"isProduction":true,
		"autoGenerateReleaseNotesOnCompletion":true,
		"includePathPatterns":["services/api/**"],
		"approximateReleaseCount":4,
		"trashed":null,
		"releaseNoteTemplate":{"id":"template-id"},
		"latestReleaseNote":{"id":"release-note-id"},
		"url":"https://linear.app/kyanite/releases/production",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null
	}`
}

func trashedReleasePipelineJSON() string {
	return strings.Replace(releasePipelineJSON(), `"trashed":null`, `"trashed":false`, 1)
}

func releaseStageJSON() string {
	return `{
		"id":"release-stage-id",
		"name":"Started",
		"color":"#00ff00",
		"type":"started",
		"position":2,
		"frozen":false,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"pipeline":{"id":"release-pipeline-id","name":"Production","slugId":"production"}
	}`
}

func releaseJSON() string {
	return `{
		"id":"release-id",
		"name":"Mobile 1.2.3",
		"slugId":"mobile-1-2-3",
		"version":"v1.2.3",
		"description":"Release body",
		"commitSha":"abc123",
		"issueCount":3,
		"trashed":null,
		"url":"https://linear.app/kyanite/release/mobile-1-2-3",
		"startDate":"2026-06-20",
		"targetDate":"2026-06-30",
		"startedAt":"2026-06-20T12:00:00Z",
		"completedAt":null,
		"canceledAt":null,
		"autoArchivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-20T12:00:00Z",
		"archivedAt":null,
		"pipeline":{"id":"release-pipeline-id","name":"Production","slugId":"production"},
		"stage":{"id":"release-stage-id","name":"Started","type":"started"},
		"releaseNotes":[{"id":"release-note-id","title":"Launch notes","slugId":"launch-notes"}],
		"creator":{"id":"user-id","displayName":"Omer"}
	}`
}

func releaseHistoryJSON() string {
	return `{
		"id":"release-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"stage","from":"planned","to":"started"}],
		"release":{"id":"release-id"}
	}`
}

func projectHistoryJSON() string {
	return `{
		"id":"project-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Backlog","to":"Started"}],
		"project":{"id":"project-id"}
	}`
}

func issueHistoryJSON() string {
	return `{
		"id":"issue-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"actorId":"user-id",
		"updatedDescription":true,
		"issue":{"id":"issue-id"}
	}`
}

func issueStateSpanJSON() string {
	return `{
		"id":"issue-state-span-id",
		"stateId":"started-state",
		"startedAt":"2026-06-19T12:00:00Z",
		"endedAt":null,
		"state":{"id":"started-state","name":"Started","type":"started"}
	}`
}

func actorBotJSON() string {
	return `{
		"id":"bot-actor-id",
		"type":"github",
		"subType":"issue",
		"name":"GitHub",
		"userDisplayName":"octocat",
		"avatarUrl":"https://example.com/github.png"
	}`
}

func commentMetadataJSON(projectID string, projectUpdateID string, userID string) string {
	return commentMetadataJSONWithID("comment-id", "", projectID, projectUpdateID, userID)
}

func commentMetadataJSONWithID(
	id string,
	parentID string,
	projectID string,
	projectUpdateID string,
	userID string,
) string {
	user := `null`
	if userID != "" {
		user = `{"id":"` + userID + `","name":"omer","displayName":"Omer"}`
	}

	return `{
		"id":"` + id + `",
		"url":"https://linear.app/comment/` + id + `",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:02:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":` + nullableStringJSON(parentID) + `,
		"issueId":null,
		"projectId":` + nullableStringJSON(projectID) + `,
		"projectUpdateId":` + nullableStringJSON(projectUpdateID) + `,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":` + user + `
	}`
}

func nullableStringJSON(value string) string {
	if value == "" {
		return `null`
	}

	return `"` + value + `"`
}

func projectAttachmentJSON() string {
	return `{
		"id":"project-attachment-id",
		"title":"Project link",
		"subtitle":"overview",
		"url":"https://example.com/project-link",
		"sourceType":"github"
	}`
}

func gitAutomationStateJSON() string {
	return `{
		"id":"git-automation-state-id",
		"event":"review",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"state":{"id":"workflow-state-id","name":"Started","type":"started"},
		"targetBranch":{"id":"target-branch-id","branchPattern":"main","isRegex":false}
	}`
}

func issueRelationJSON() string {
	return `{
		"id":"issue-relation-id",
		"type":"blocks",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},
		"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}
	}`
}

func projectLabelJSON(id string, name string) string {
	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":"Project label",
		"color":"#f2c94c",
		"isGroup":false,
		"lastAppliedAt":"2026-06-19T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"parent":null
	}`
}

func projectRelationJSON() string {
	return `{
		"id":"project-relation-id",
		"type":"blocks",
		"anchorType":"project",
		"relatedAnchorType":"project",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"project":{"id":"project-id","name":"Pinned project"},
		"projectMilestone":null,
		"relatedProject":{"id":"related-project-id","name":"Related project"},
		"relatedProjectMilestone":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func customerNeedJSON() string {
	return `{
		"id":"customer-need-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"priority":1,
		"body":"Need body",
		"content":"Need content",
		"url":"https://example.com/need",
		"customer":{"id":"customer-id","name":"Acme"},
		"issue":{"id":"issue-id","identifier":"LIT-1","title":"Need issue"},
		"project":{"id":"project-id","name":"Customer project"}
	}`
}

func entityExternalLinkJSON() string {
	return `{
		"id":"release-link-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"url":"https://example.com/runbook",
		"label":"Runbook",
		"sortOrder":1.5,
		"creator":{"id":"user-id","displayName":"Omer"},
		"initiative":null,
		"project":null
	}`
}

func entityExternalLinkWithParentsJSON() string {
	return `{
		"id":"release-link-parent-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"url":"https://example.com/plan",
		"label":"Plan",
		"sortOrder":2,
		"creator":null,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func userSettingsJSON() string {
	return `{
		"id":"settings-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"autoAssignToSelf":true,
		"feedLastSeenTime":null,
		"feedSummarySchedule":"daily",
		"showFullUserNames":false,
		"subscribedToChangelog":true,
		"subscribedToDPA":false,
		"subscribedToInviteAccepted":true,
		"subscribedToPrivacyLegalUpdates":true,
		"user":{"id":"user-id"},
		"notificationCategoryPreferences":` + notificationCategoriesJSON() + `,
		"notificationChannelPreferences":` + notificationChannelJSON() + `,
		"notificationDeliveryPreferences":` + notificationDeliveryPreferencesJSON() + `
	}`
}

func notificationCategoriesJSON() string {
	channel := notificationChannelJSON()
	return `{
		"appsAndIntegrations":` + channel + `,
		"assignments":` + channel + `,
		"billing":` + channel + `,
		"commentsAndReplies":` + channel + `,
		"customers":` + channel + `,
		"documentChanges":` + channel + `,
		"feed":` + channel + `,
		"mentions":` + channel + `,
		"postsAndUpdates":` + channel + `,
		"reactions":` + channel + `,
		"reminders":` + channel + `,
		"reviews":` + channel + `,
		"statusChanges":` + channel + `,
		"subscriptions":` + channel + `,
		"system":` + channel + `,
		"triage":` + channel + `
	}`
}

func userSettingsCategoryJSON(category string) string {
	return `{"userSettings":{"notificationCategoryPreferences":{"` + category + `":` + notificationChannelJSON() + `}}}`
}

func notificationChannelJSON() string {
	return `{"desktop":true,"email":false,"mobile":true,"slack":true}`
}

func notificationDeliveryPreferencesJSON() string {
	return `{"mobile":` + notificationDeliveryChannelJSON() + `}`
}

func notificationDeliveryChannelJSON() string {
	return `{"notificationsDisabled":false,"schedule":` + notificationDeliveryScheduleJSON() + `}`
}

func notificationDeliveryScheduleJSON() string {
	day := notificationDeliveryDayJSON()
	return `{
		"disabled":false,
		"friday":` + day + `,
		"monday":` + day + `,
		"saturday":` + day + `,
		"sunday":` + day + `,
		"thursday":` + day + `,
		"tuesday":` + day + `,
		"wednesday":` + day + `
	}`
}

func notificationDeliveryDayJSON() string {
	return `{"start":"09:00","end":"18:00"}`
}

func userSettingsScheduleDayJSON(day string) string {
	return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":{"` + day + `":` +
		notificationDeliveryDayJSON() + `}}}}}`
}

func userSettingsThemeJSON(includeCustom bool) string {
	custom := "null"
	if includeCustom {
		custom = userSettingsCustomThemeJSON(true)
	}

	return `{"preset":"custom","custom":` + custom + `}`
}

func userSettingsCustomThemeJSON(includeSidebar bool) string {
	sidebar := "null"
	if includeSidebar {
		sidebar = userSettingsCustomSidebarThemeJSON()
	}

	return `{"accent":[50.5,20.5,10.5],"base":[90.5,0,0],"contrast":50,"sidebar":` + sidebar + `}`
}

func userSettingsCustomSidebarThemeJSON() string {
	return `{"accent":[60.5,20.5,10.5],"base":[20.5,0,0],"contrast":70}`
}

func releaseNoteJSON() string {
	return `{
		"id":"release-note-id",
		"title":"Launch notes",
		"slugId":"launch-notes",
		"generationStatus":"completed",
		"releaseCount":2,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-20T12:00:00Z",
		"archivedAt":null,
		"pipeline":{"id":"release-pipeline-id","name":"Production","slugId":"production"},
		"firstRelease":{"id":"release-id","name":"Mobile 1.2.2","version":"v1.2.2"},
		"lastRelease":{"id":"release-id","name":"Mobile 1.2.3","version":"v1.2.3"}
	}`
}
