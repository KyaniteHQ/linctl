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
		if fake.orgWideLabel {
			return `{"issueLabel":{"id":"label-id","name":"Bug","description":null,"color":"#ff0000","isGroup":false,"team":null}}`, true
		}

		return `{"issueLabel":` + commandLabelJSON("") + `}`, true
	case "IssueLabelCreate":
		return `{"issueLabelCreate":{"success":true,"issueLabel":` +
			commandNamedLabelJSON("label-id", "Created label", "#ff0000", "") + `}}`, true
	case "IssueLabelUpdate":
		return `{"issueLabelUpdate":{"success":true,"issueLabel":` +
			commandNamedLabelJSON("label-id", "Updated label", "#ff0000", "") + `}}`, true
	case "IssueLabelRetire":
		return `{"issueLabelRetire":{"success":true,"issueLabel":` +
			commandNamedLabelJSON("label-id", "Retired label", "#ff0000", "") + `}}`, true
	case "IssueLabelRestore":
		return `{"issueLabelRestore":{"success":true,"issueLabel":` +
			commandNamedLabelJSON("label-id", "Restored label", "#ff0000", "") + `}}`, true
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
	case "userSettings_notificationChannelPreferences":
		return `{"userSettings":{"notificationChannelPreferences":` + commandNotificationChannelJSON() + `}}`, true
	case "userSettings_notificationDeliveryPreferences":
		return `{"userSettings":{"notificationDeliveryPreferences":` + commandNotificationDeliveryPreferencesJSON() + `}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile":
		return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":` + commandNotificationDeliveryChannelJSON() + `}}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule":
		return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":` + commandNotificationDeliveryScheduleJSON() + `}}}}`, true
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
