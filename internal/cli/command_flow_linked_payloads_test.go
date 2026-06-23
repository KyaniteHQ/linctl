package cli

func commandFlowStateAndCommentPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "workflowStates":
		return `{"workflowStates":{"nodes":[` + commandWorkflowStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "workflowState":
		return `{"workflowState":` + commandWorkflowStateJSON() + `}`, true
	case "workflowState_issues":
		return `{"workflowState":{"id":"workflow-state-id","name":"Started","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "state-id", "Todo", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	}

	return commandFlowInitiativePayload(operation, fake)
}

func commandFlowFilePayload(operation string) (string, bool) {
	if operation != "fileUpload" {
		return "", false
	}

	return `{"fileUpload":{"success":true,"uploadFile":{` +
		`"filename":"upload.txt","contentType":"text/plain","size":11,` +
		`"uploadUrl":"https://uploads.example/put","assetUrl":"https://assets.example/file.txt",` +
		`"headers":[{"key":"x-test","value":"1"}]}}}`, true
}

func commandFlowCommentPayload(operation string) (string, bool) {
	switch operation {
	case "comments":
		return `{"comments":{"nodes":[` + commandTopLevelCommentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "comment":
		return `{"comment":` + commandTopLevelCommentJSON() + `}`, true
	case "comment_botActor":
		return `{"comment":{"id":"comment-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "comment_children":
		return `{"comment":{"id":"comment-id","children":{"nodes":[` +
			commandCommentMetadataJSONWithID("child-comment-id", "comment-id", "", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "comment_createdIssues":
		return `{"comment":{"id":"comment-id","createdIssues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "CommentUpdate":
		return `{"commentUpdate":{"success":true,"comment":` + commandTopLevelCommentJSON() + `}}`, true
	case "CommentDelete":
		return `{"commentDelete":{"success":true,"entityId":"comment-id"}}`, true
	default:
		return "", false
	}
}

func commandFlowInitiativePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "initiatives":
		return `{"initiatives":{"nodes":[` + commandInitiativeJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiative":
		return `{"initiative":` + commandInitiativeJSON() + `}`, true
	case "initiative_history":
		return `{"initiative":{"history":{"nodes":[` + commandInitiativeHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_links":
		return `{"initiative":{"links":{"nodes":[` + commandEntityExternalLinkJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_subInitiatives":
		return `{"initiative":{"subInitiatives":{"nodes":[` + commandSubInitiativeJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_initiativeUpdates":
		return `{"initiative":{"initiativeUpdates":{"nodes":[` + commandInitiativeUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_documents":
		return `{"initiative":{"documents":{"nodes":[` + commandDocumentJSON(
			"Spec",
			`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
		) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_projects":
		return `{"initiative":{"projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiativeRelations":
		return `{"initiativeRelations":{"nodes":[` + commandInitiativeRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiativeRelation":
		return `{"initiativeRelation":` + commandInitiativeRelationJSON() + `}`, true
	}

	return commandFlowInitiativeUpdatePayload(operation, fake)
}

func commandFlowInitiativeUpdatePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "initiativeToProjects":
		return `{"initiativeToProjects":{"nodes":[` + commandInitiativeToProjectJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiativeToProject":
		return `{"initiativeToProject":` + commandInitiativeToProjectJSON() + `}`, true
	case "roadmapToProjects":
		return `{"roadmapToProjects":{"nodes":[` + commandRoadmapToProjectJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "roadmapToProject":
		return `{"roadmapToProject":` + commandRoadmapToProjectJSON() + `}`, true
	case "initiativeUpdates":
		return `{"initiativeUpdates":{"nodes":[` + commandInitiativeUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiativeUpdate":
		return `{"initiativeUpdate":` + commandInitiativeUpdateJSON() + `}`, true
	case "initiativeUpdate_comments":
		return `{"initiativeUpdate":{"id":"initiative-update-id","comments":{"nodes":[` +
			commandCommentMetadataJSON("", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	}

	return commandFlowExtraReadPayload(operation, fake)
}

//nolint:gocyclo // The table-driven command-flow fake is intentionally centralized by operation name.
func commandFlowExtraReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "auditEntryTypes":
		return `{"auditEntryTypes":[{"type":"user_login","description":"User logged in"}]}`, true
	case "notifications":
		return `{"notifications":{"nodes":[` + commandNotificationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "notification":
		return `{"notification":` + commandNotificationJSON() + `}`, true
	case "notificationSubscriptions":
		return `{"notificationSubscriptions":{"nodes":[` + commandNotificationSubscriptionJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "notificationSubscription":
		return `{"notificationSubscription":` + commandNotificationSubscriptionJSON() + `}`, true
	case "triageResponsibilities":
		return `{"triageResponsibilities":{"nodes":[` + commandTriageResponsibilityJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "triageResponsibility":
		return `{"triageResponsibility":` + commandTriageResponsibilityJSON() + `}`, true
	case "triageResponsibility_manualSelection":
		return `{"triageResponsibility":{"id":"triage-responsibility-id","manualSelection":{"userIds":["user-id","other-user-id"]}}}`, true
	case "slaConfigurations":
		if fake.emptySLAConfigurations {
			return `{"slaConfigurations":[]}`, true
		}
		return `{"slaConfigurations":[` + commandSLAConfigurationJSON() + `]}`, true
	case "semanticSearch":
		if fake.emptySemanticSearch {
			return `{"semanticSearch":{"results":[]}}`, true
		}
		return `{"semanticSearch":{"results":[` + commandSemanticSearchResultJSON() + `]}}`, true
	case "searchDocuments":
		if fake.emptySearchDocuments {
			return `{"searchDocuments":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":0}}`, true
		}
		return `{"searchDocuments":{"nodes":[` + commandSearchDocumentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":1}}`, true
	case "searchIssues":
		if fake.emptySearchIssues {
			return `{"searchIssues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":0}}`, true
		}
		return `{"searchIssues":{"nodes":[` + commandSearchIssueJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":1}}`, true
	case "searchProjects":
		if fake.emptySearchProjects {
			return `{"searchProjects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":0}}`, true
		}
		return `{"searchProjects":{"nodes":[` + commandSearchProjectJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":1}}`, true
	case "releasePipelines":
		return `{"releasePipelines":{"nodes":[` + commandReleasePipelineJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releasePipeline":
		return `{"releasePipeline":` + commandReleasePipelineJSON() + `}`, true
	case "releasePipeline_releases":
		return `{"releasePipeline":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releasePipeline_stages":
		return `{"releasePipeline":{"stages":{"nodes":[` + commandReleaseStageJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releasePipeline_teams":
		return `{"releasePipeline":{"teams":{"nodes":[` + commandTeamJSON(true) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releaseStages":
		return `{"releaseStages":{"nodes":[` + commandReleaseStageJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releaseStage":
		return `{"releaseStage":` + commandReleaseStageJSON() + `}`, true
	case "releaseStage_releases":
		return `{"releaseStage":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releases":
		return `{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releaseSearch":
		if fake.emptyReleaseSearch {
			return `{"releaseSearch":[]}`, true
		}
		return `{"releaseSearch":[` + commandReleaseJSON() + `]}`, true
	case "release":
		return `{"release":` + commandReleaseJSON() + `}`, true
	case "release_history":
		return `{"release":{"history":{"nodes":[` + commandReleaseHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "release_documents":
		return `{"release":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "release_issues":
		return `{"release":{"issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "release_links":
		return `{"release":{"links":{"nodes":[` + commandEntityExternalLinkJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "entityExternalLink":
		return `{"entityExternalLink":` + commandEntityExternalLinkJSON() + `}`, true
	case "releaseNotes":
		return `{"releaseNotes":{"nodes":[` + commandReleaseNoteJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releaseNote":
		return `{"releaseNote":` + commandReleaseNoteJSON() + `}`, true
	case "issueToReleases":
		return `{"issueToReleases":{"nodes":[` + commandIssueToReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueToRelease":
		return `{"issueToRelease":` + commandIssueToReleaseJSON() + `}`, true
	case "timeSchedules":
		return `{"timeSchedules":{"nodes":[` + commandTimeScheduleJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "timeSchedule":
		return `{"timeSchedule":` + commandTimeScheduleJSON() + `}`, true
	case "templates":
		return `{"templates":[` + commandTemplateJSON() + `]}`, true
	case "template":
		return `{"template":` + commandTemplateJSON() + `}`, true
	case "templateContent":
		return `{"template":{"id":"template-id","name":"Bug report","templateData":` +
			`{"title":"Template title","description":"## Steps\n\nReproduce here"}}}`, true
	case "roadmaps":
		return `{"roadmaps":{"nodes":[` + commandRoadmapJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "roadmap":
		return `{"roadmap":` + commandRoadmapJSON() + `}`, true
	case "roadmap_projects":
		return `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customViews":
		return `{"customViews":{"nodes":[` + commandCustomViewJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customViewHasSubscribers":
		return `{"customViewHasSubscribers":{"hasSubscribers":true}}`, true
	case "customView":
		return `{"customView":` + commandCustomViewJSON() + `}`, true
	case "customView_initiatives":
		return `{"customView":{"initiatives":{"nodes":[` + commandInitiativeJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customView_issues":
		return `{"customView":{"issues":{"nodes":[` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customView_organizationViewPreferences":
		return `{"customView":{"organizationViewPreferences":` + commandCustomViewPreferencesJSON("priority", "list") + `}}`, true
	case "customView_organizationViewPreferences_preferences":
		return `{"customView":{"organizationViewPreferences":{"preferences":` + commandCustomViewPreferenceValuesJSON("priority", "list") + `}}}`, true
	case "customView_projects":
		return `{"customView":{"projects":{"nodes":[` + commandProjectJSON("Pinned project", "Backlog", "backlog") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customView_userViewPreferences":
		return `{"customView":{"userViewPreferences":` + commandCustomViewScopedPreferencesJSON("user", "updatedAt", "board") + `}}`, true
	case "customView_userViewPreferences_preferences":
		return `{"customView":{"userViewPreferences":{"preferences":` + commandCustomViewPreferenceValuesJSON("updatedAt", "board") + `}}}`, true
	case "customView_viewPreferencesValues":
		return `{"customView":{"viewPreferencesValues":` + commandCustomViewPreferenceValuesJSON("updatedAt", "board") + `}}`, true
	case "customers":
		return `{"customers":{"nodes":[` + commandCustomerJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customer":
		return `{"customer":` + commandCustomerJSON() + `}`, true
	case "customerNeeds":
		return `{"customerNeeds":{"nodes":[` + commandCustomerNeedJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customerNeed":
		return `{"customerNeed":` + commandCustomerNeedJSON() + `}`, true
	case "customerNeed_projectAttachment":
		if fake.missingCustomerNeedAttachment {
			return `{"customerNeed":{"id":"customer-need-id","projectAttachment":null}}`, true
		}
		return `{"customerNeed":{"id":"customer-need-id","projectAttachment":` + commandAttachmentJSON() + `}}`, true
	case "customerStatuses":
		return `{"customerStatuses":{"nodes":[` + commandCustomerStatusJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customerStatus":
		return `{"customerStatus":` + commandCustomerStatusJSON() + `}`, true
	case "customerTiers":
		return `{"customerTiers":{"nodes":[` + commandCustomerTierJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customerTier":
		return `{"customerTier":` + commandCustomerTierJSON() + `}`, true
	case "favorites":
		return `{"favorites":{"nodes":[` + commandFavoriteJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "favorite_children":
		return `{"favorite":{"children":{"nodes":[` + commandFavoriteChildJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "favorite":
		return `{"favorite":` + commandFavoriteJSON() + `}`, true
	case "emojis":
		return `{"emojis":{"nodes":[` + commandEmojiJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "emoji":
		return `{"emoji":` + commandEmojiJSON() + `}`, true
	case "attachments":
		return `{"attachments":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "attachmentsForURL":
		return `{"attachmentsForURL":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "attachment":
		return `{"attachment":` + commandAttachmentJSON() + `}`, true
	default:
		return "", false
	}
}

func commandFlowLabelChildPayload(operation string) (string, bool) {
	switch operation {
	case "issueLabel_children":
		return `{"issueLabel":{"id":"label-id","name":"Bug","children":{"nodes":[` +
			commandNamedLabelJSON("child-label-id", "Mobile", "#56ccf2", "child label") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueLabel_issues":
		return `{"issueLabel":{"id":"label-id","name":"Bug","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}
