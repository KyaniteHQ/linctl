package cli

import (
	"encoding/json"
	"errors"

	"github.com/Khan/genqlient/graphql"
)

func commandGitAutomationStateJSON() string {
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

func commandNotificationJSON() string {
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
		"externalUserActor":null
	}`
}

func commandNotificationSubscriptionJSON() string {
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

func commandTriageResponsibilityJSON() string {
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

func commandSLAConfigurationJSON() string {
	return `{
		"id":"sla-configuration-id",
		"name":"First response",
		"conditions":{"priority":{"eq":1}},
		"sla":3600000,
		"slaType":"all",
		"removesSla":false
	}`
}

func commandSemanticSearchResultJSON() string {
	return `{
		"id":"issue-id",
		"type":"issue",
		"issue":{"id":"issue-id","identifier":"LIT-3","title":"Search result","url":"https://linear.app/kyanite/issue/LIT-3"},
		"project":null,
		"initiative":null,
		"document":null
	}`
}

func commandSearchDocumentJSON() string {
	return `{
		"id":"search-document-id",
		"title":"Search spec",
		"slugId":"search-spec",
		"url":"https://linear.app/kyanite/document/search-spec",
		"project":null,
		"initiative":null,
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"issue":null,
		"release":null,
		"cycle":null
	}`
}

func commandSearchIssueJSON() string {
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

func commandSearchProjectJSON() string {
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

func commandReleasePipelineJSON() string {
	return `{
		"id":"release-pipeline-id",
		"name":"Production",
		"slugId":"production",
		"type":"scheduled",
		"isProduction":true,
		"autoGenerateReleaseNotesOnCompletion":true,
		"approximateReleaseCount":4,
		"url":"https://linear.app/kyanite/releases/production",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"trashed":null,
		"includePathPatterns":["services/api/**"],
		"releaseNoteTemplate":{"id":"template-id"},
		"latestReleaseNote":{"id":"release-note-id"}
	}`
}

func commandReleaseStageJSON() string {
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

func commandReleaseJSON() string {
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

func commandReleaseHistoryJSON() string {
	return `{
		"id":"release-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"stage","from":"planned","to":"started"}],
		"release":{"id":"release-id"}
	}`
}

func commandEntityExternalLinkJSON() string {
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

func commandReleaseNoteJSON() string {
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

func commandTopLevelCommentJSON() string {
	return `{
		"id":"comment-id",
		"body":"First comment",
		"url":"https://linear.app/comment/comment-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":null,
		"issueId":"issue-id",
		"projectId":null,
		"projectUpdateId":null,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandCommentMetadataJSON(projectID string, projectUpdateID string) string {
	return commandCommentMetadataJSONWithID("comment-id", "", projectID, projectUpdateID)
}

func commandCommentMetadataJSONWithID(
	id string,
	parentID string,
	projectID string,
	projectUpdateID string,
) string {
	return `{
		"id":"` + id + `",
		"url":"https://linear.app/comment/` + id + `",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":` + commandNullableStringJSON(parentID) + `,
		"issueId":null,
		"projectId":` + commandNullableStringJSON(projectID) + `,
		"projectUpdateId":` + commandNullableStringJSON(projectUpdateID) + `,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandActorBotJSON() string {
	return `{
		"id":"bot-actor-id",
		"type":"github",
		"subType":"issue",
		"name":"GitHub",
		"userDisplayName":"octocat",
		"avatarUrl":"https://example.com/github.png"
	}`
}

func commandUserSettingsJSON() string {
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
		"notificationCategoryPreferences":` + commandNotificationCategoriesJSON() + `,
		"notificationChannelPreferences":` + commandNotificationChannelJSON() + `,
		"notificationDeliveryPreferences":` + commandNotificationDeliveryPreferencesJSON() + `
	}`
}

func commandNotificationCategoriesJSON() string {
	channel := commandNotificationChannelJSON()
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

func commandNotificationChannelJSON() string {
	return `{"desktop":true,"email":false,"mobile":true,"slack":false}`
}

func commandNotificationDeliveryPreferencesJSON() string {
	return `{"mobile":` + commandNotificationDeliveryChannelJSON() + `}`
}

func commandNotificationDeliveryChannelJSON() string {
	return `{"notificationsDisabled":false,"schedule":` + commandNotificationDeliveryScheduleJSON() + `}`
}

func commandNotificationDeliveryScheduleJSON() string {
	day := commandNotificationDeliveryDayJSON()
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

func commandNotificationDeliveryDayJSON() string {
	return `{"start":"09:00","end":"18:00"}`
}

func commandUserSettingsThemeJSON(includeCustom bool) string {
	custom := "null"
	if includeCustom {
		custom = commandUserSettingsCustomThemeJSON(true)
	}

	return `{"preset":"custom","custom":` + custom + `}`
}

func commandUserSettingsCustomThemeJSON(includeSidebar bool) string {
	sidebar := "null"
	if includeSidebar {
		sidebar = commandUserSettingsCustomSidebarThemeJSON()
	}

	return `{"accent":[50.5,20.5,10.5],"base":[90.5,0,0],"contrast":50,"sidebar":` + sidebar + `}`
}

func commandUserSettingsCustomSidebarThemeJSON() string {
	return `{"accent":[60.5,20.5,10.5],"base":[20.5,0,0],"contrast":70}`
}

func commandNullableStringJSON(value string) string {
	if value == "" {
		return `null`
	}

	return `"` + value + `"`
}

var _ graphql.Client = commandFlowFakeClient{}

func requestVariable[T comparable](request *graphql.Request, keys ...string) (T, error) {
	var zero T
	payload, err := json.Marshal(request.Variables)
	if err != nil {
		return zero, err
	}
	var variables map[string]any
	if err := json.Unmarshal(payload, &variables); err != nil {
		return zero, err
	}
	current := any(variables)
	for _, key := range keys {
		object, ok := current.(map[string]any)
		if !ok {
			return zero, errors.New("request variable is not an object")
		}
		value, ok := object[key]
		if !ok {
			return zero, errors.New("request variable missing " + key)
		}
		current = value
	}
	value, ok := current.(T)
	if !ok {
		return zero, errors.New("request variable has unexpected type")
	}

	return value, nil
}
