package cli

func commandFlowTeamChildPayload(operation string) (string, bool) {
	switch operation {
	case "team_cycles":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","cycles":{"nodes":[` +
			commandCycleJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_issues":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "state-id", "Todo", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_labels":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","labels":{"nodes":[` +
			commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_memberships":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","memberships":{"nodes":[` +
			commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_projects":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_releasePipelines":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","releasePipelines":{"nodes":[` +
			commandReleasePipelineJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_states":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","states":{"nodes":[` +
			commandWorkflowStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_gitAutomationStates":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","gitAutomationStates":{"nodes":[` +
			commandGitAutomationStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_templates":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","templates":{"nodes":[` +
			commandTemplateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
func commandFlowPeopleAndReferencePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowTeamChildPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowUserChildPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowLabelChildPayload(operation); ok {
		return payload, true
	}

	switch operation {
	case "Documents":
		return `{"documents":{"nodes":[` + commandDocumentJSON(
			"Spec",
			`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
		) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "document":
		return `{"document":` + commandDocumentJSON(
			"Team note",
			`"project":{"id":"project-id","name":"Pinned project"},`+
				`"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}`, true
	case "DocumentCreate":
		return `{"documentCreate":{"success":true,"document":` + commandDocumentJSON(
			"Created doc",
			`"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}}`, true
	case "DocumentUpdate":
		return `{"documentUpdate":{"success":true,"document":` + commandDocumentJSON(
			"Updated doc",
			`"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}}`, true
	case "ProjectUpdateCreate":
		return `{"projectUpdateCreate":{"success":true,"projectUpdate":` + commandProjectUpdateJSON() + `}}`, true
	case "document_comments":
		return `{"document":{"id":"document-id","comments":{"nodes":[` +
			commandCommentMetadataJSON("", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "IssueLabels":
		return `{"issueLabels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueLabel":
		return `{"issueLabel":` + commandLabelJSON("") + `}`, true
	case "team":
		return `{"team":` + commandTeamJSON(true) + `}`, true
	case "team_members":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "users":
		return `{"users":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "user":
		return `{"user":` + commandUserJSON() + `}`, true
	case "viewer":
		return `{"viewer":` + commandUserJSON() + `}`, true
	case "viewer_drafts":
		if fake.emptyViewerDrafts {
			return `{"viewer":{"drafts":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"viewer":{"drafts":{"nodes":[` + commandDraftJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	}
	if payload, ok := commandFlowUserSettingsPayload(operation); ok {
		return payload, true
	}

	return commandFlowStateAndCommentPayload(operation, fake)
}

func commandFlowUserChildPayload(operation string) (string, bool) {
	switch operation {
	case "user_assignedIssues":
		return commandFlowUserIssueListPayload("user", "assignedIssues"), true
	case "user_createdIssues":
		return commandFlowUserIssueListPayload("user", "createdIssues"), true
	case "user_delegatedIssues":
		return commandFlowUserIssueListPayload("user", "delegatedIssues"), true
	case "user_teamMemberships":
		return `{"user":{"teamMemberships":{"nodes":[` + commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "user_teams":
		return `{"user":{"teams":{"nodes":[` + commandTeamJSON(false) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "viewer_assignedIssues":
		return commandFlowUserIssueListPayload("viewer", "assignedIssues"), true
	case "viewer_createdIssues":
		return commandFlowUserIssueListPayload("viewer", "createdIssues"), true
	case "viewer_delegatedIssues":
		return commandFlowUserIssueListPayload("viewer", "delegatedIssues"), true
	case "viewer_teamMemberships":
		return `{"viewer":{"teamMemberships":{"nodes":[` + commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "viewer_teams":
		return `{"viewer":{"teams":{"nodes":[` + commandTeamJSON(false) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

//nolint:gocyclo // Each branch mirrors a distinct official UserSettings operation.
func commandFlowUserSettingsPayload(operation string) (string, bool) {
	switch operation {
	case "userSettings":
		return `{"userSettings":` + commandUserSettingsJSON() + `}`, true
	case "userSettings_notificationCategoryPreferences":
		return `{"userSettings":{"notificationCategoryPreferences":` + commandNotificationCategoriesJSON() + `}}`, true
	case "userSettings_notificationCategoryPreferences_appsAndIntegrations":
		return commandUserSettingsCategoryPayload("appsAndIntegrations"), true
	case "userSettings_notificationCategoryPreferences_assignments":
		return commandUserSettingsCategoryPayload("assignments"), true
	case "userSettings_notificationCategoryPreferences_billing":
		return commandUserSettingsCategoryPayload("billing"), true
	case "userSettings_notificationCategoryPreferences_commentsAndReplies":
		return commandUserSettingsCategoryPayload("commentsAndReplies"), true
	case "userSettings_notificationCategoryPreferences_customers":
		return commandUserSettingsCategoryPayload("customers"), true
	case "userSettings_notificationCategoryPreferences_documentChanges":
		return commandUserSettingsCategoryPayload("documentChanges"), true
	case "userSettings_notificationCategoryPreferences_feed":
		return commandUserSettingsCategoryPayload("feed"), true
	case "userSettings_notificationCategoryPreferences_mentions":
		return commandUserSettingsCategoryPayload("mentions"), true
	case "userSettings_notificationCategoryPreferences_postsAndUpdates":
		return commandUserSettingsCategoryPayload("postsAndUpdates"), true
	case "userSettings_notificationCategoryPreferences_reactions":
		return commandUserSettingsCategoryPayload("reactions"), true
	case "userSettings_notificationCategoryPreferences_reminders":
		return commandUserSettingsCategoryPayload("reminders"), true
	case "userSettings_notificationCategoryPreferences_reviews":
		return commandUserSettingsCategoryPayload("reviews"), true
	case "userSettings_notificationCategoryPreferences_statusChanges":
		return commandUserSettingsCategoryPayload("statusChanges"), true
	case "userSettings_notificationCategoryPreferences_subscriptions":
		return commandUserSettingsCategoryPayload("subscriptions"), true
	case "userSettings_notificationCategoryPreferences_system":
		return commandUserSettingsCategoryPayload("system"), true
	case "userSettings_notificationCategoryPreferences_triage":
		return commandUserSettingsCategoryPayload("triage"), true
	case "userSettings_notificationChannelPreferences":
		return `{"userSettings":{"notificationChannelPreferences":` + commandNotificationChannelJSON() + `}}`, true
	case "userSettings_notificationDeliveryPreferences":
		return `{"userSettings":{"notificationDeliveryPreferences":` + commandNotificationDeliveryPreferencesJSON() + `}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile":
		return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":` + commandNotificationDeliveryChannelJSON() + `}}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule":
		return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":` + commandNotificationDeliveryScheduleJSON() + `}}}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_friday":
		return commandUserSettingsScheduleDayPayload("friday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_monday":
		return commandUserSettingsScheduleDayPayload("monday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_saturday":
		return commandUserSettingsScheduleDayPayload("saturday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_sunday":
		return commandUserSettingsScheduleDayPayload("sunday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_thursday":
		return commandUserSettingsScheduleDayPayload("thursday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_tuesday":
		return commandUserSettingsScheduleDayPayload("tuesday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_wednesday":
		return commandUserSettingsScheduleDayPayload("wednesday"), true
	case "userSettings_theme":
		return `{"userSettings":{"theme":` + commandUserSettingsThemeJSON(true) + `}}`, true
	case "userSettings_theme_custom":
		return `{"userSettings":{"theme":{"custom":` + commandUserSettingsCustomThemeJSON(true) + `}}}`, true
	case "userSettings_theme_custom_sidebar":
		return `{"userSettings":{"theme":{"custom":{"sidebar":` + commandUserSettingsCustomSidebarThemeJSON() + `}}}}`, true
	default:
		return "", false
	}
}

func commandUserSettingsCategoryPayload(category string) string {
	return `{"userSettings":{"notificationCategoryPreferences":{"` + category + `":` + commandNotificationChannelJSON() + `}}}`
}

func commandUserSettingsScheduleDayPayload(day string) string {
	return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":{"` + day + `":` +
		commandNotificationDeliveryDayJSON() + `}}}}}`
}
