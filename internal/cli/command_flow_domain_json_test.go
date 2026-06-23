package cli

func commandWorkflowStateJSON() string {
	return `{
		"id":"workflow-state-id",
		"name":"Started",
		"type":"started",
		"color":"#f2c94c",
		"position":2,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandInitiativeJSON() string {
	return `{
		"id":"initiative-id",
		"name":"Platform",
		"description":"Platform initiative",
		"status":"Active",
		"priority":2,
		"targetDate":"2026-12-31",
		"slugId":"platform-init",
		"url":"https://linear.app/kyanite/initiative/platform-init"
	}`
}

func commandSubInitiativeJSON() string {
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

func commandInitiativeHistoryJSON() string {
	return `{
		"id":"initiative-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Planned","to":"Active"}],
		"initiative":{"id":"initiative-id"}
	}`
}

func commandCustomViewJSON() string {
	return `{
		"id":"custom-view-id",
		"name":"My issues",
		"description":"Saved issue view",
		"modelName":"Issue",
		"shared":true,
		"color":"#5e6ad2",
		"slugId":"my-issues"
	}`
}

func commandCustomViewPreferencesJSON(ordering string, layout string) string {
	return commandCustomViewScopedPreferencesJSON("organization", ordering, layout)
}

func commandCustomViewScopedPreferencesJSON(scope string, ordering string, layout string) string {
	return `{
		"id":"view-preferences-id",
		"createdAt":"2026-06-01T12:00:00Z",
		"updatedAt":"2026-06-01T12:01:00Z",
		"archivedAt":null,
		"type":"` + scope + `",
		"viewType":"customView",
		"preferences":` + commandCustomViewPreferenceValuesJSON(ordering, layout) + `
	}`
}

func commandCustomViewPreferenceValuesJSON(ordering string, layout string) string {
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

func commandCustomerJSON() string {
	return `{
		"id":"customer-id",
		"name":"Acme",
		"domains":["acme.example"],
		"externalIds":["crm-acme"],
		"slackChannelId":"slack-channel-id",
		"status":{"id":"status-id","name":"Active"},
		"tier":{"id":"tier-id","name":"Enterprise"},
		"owner":{"id":"user-id","displayName":"Omer"},
		"revenue":120000,
		"size":42,
		"approximateNeedCount":3,
		"slugId":"acme",
		"url":"https://linear.app/kyanite/customer/acme"
	}`
}

func commandRoadmapJSON() string {
	return `{
		"id":"roadmap-id",
		"name":"Platform roadmap",
		"description":"Roadmap body",
		"color":"#5e6ad2",
		"slugId":"platform-roadmap",
		"sortOrder":1,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"url":"https://linear.app/kyanite/roadmap/platform-roadmap",
		"creator":{"id":"user-id","displayName":"Omer"},
		"owner":{"id":"owner-id","displayName":"Owner"}
	}`
}

func commandTimeScheduleJSON() string {
	return `{
		"id":"time-schedule-id",
		"name":"Primary on-call",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"externalId":"pd-primary",
		"externalUrl":"https://example.com/schedule",
		"integration":{"id":"integration-id"},
		"entries":[{"startsAt":"2026-06-20T00:00:00Z","endsAt":"2026-06-21T00:00:00Z","userId":"user-id","userEmail":"omer@example.com"}]
	}`
}

func commandTemplateJSON() string {
	return `{
		"id":"template-id",
		"name":"Bug report",
		"type":"issue",
		"description":"Bug report template",
		"icon":"bug",
		"color":"#ff0000",
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

func commandCustomerNeedJSON() string {
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

func commandIssueSharedAccessJSON() string {
	return `{
		"isShared":true,
		"viewerHasOnlySharedAccess":false,
		"sharedWithCount":2,
		"disallowedIssueFields":["description","priority"]
	}`
}

func commandCustomerStatusJSON() string {
	return `{
		"id":"customer-status-id",
		"name":"active",
		"displayName":"Active",
		"color":"#00ff00",
		"description":"Active customers",
		"position":1,
		"archivedAt":null
	}`
}

func commandCustomerTierJSON() string {
	return `{
		"id":"customer-tier-id",
		"name":"enterprise",
		"displayName":"Enterprise",
		"color":"#0000ff",
		"description":"Enterprise customers",
		"position":2,
		"archivedAt":null
	}`
}

func commandFavoriteJSON() string {
	return `{
		"id":"favorite-id",
		"type":"issue",
		"folderName":null,
		"url":"https://linear.app/kyanite/issue/LIT-1"
	}`
}

func commandFavoriteChildJSON() string {
	return `{
		"id":"favorite-child-id",
		"type":"project",
		"folderName":null,
		"url":"https://linear.app/kyanite/project/project-id"
	}`
}

func commandEmojiJSON() string {
	return `{
		"id":"emoji-id",
		"name":"party",
		"url":"https://linear.app/kyanite/emoji/party.png",
		"source":"custom"
	}`
}

func commandAttachmentJSON() string {
	return `{
		"id":"attachment-id",
		"title":"Linked PR",
		"subtitle":"feat: add thing",
		"url":"https://github.com/kyanite/linctl/pull/1",
		"sourceType":"github"
	}`
}
