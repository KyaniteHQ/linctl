package cli

func commandFlowProjectPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowProjectStatusPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowProjectLabelPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowProjectRelationPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowProjectReadPayload(operation, fake); ok {
		return payload, true
	}

	return commandFlowProjectWritePayload(operation)
}

//nolint:gocyclo // The fake payload switch mirrors the project command operation surface.
func commandFlowProjectReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "Projects":
		if fake.emptyProjectList {
			return `{"team":{"projects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"team":{"projects":{"nodes":[` + commandProjectJSON("Listed project", "Backlog", "backlog") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projects":
		if fake.emptyProjectList {
			return `{"projects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "project":
		return `{"project":` + commandProjectJSON("Detail project", "Backlog", "backlog") + `}`, true
	case "project_attachments":
		return `{"project":{"id":"project-id","name":"Detail project","attachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_documents":
		return `{"project":{"id":"project-id","name":"Detail project","documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_externalLinks":
		return `{"project":{"id":"project-id","name":"Detail project","externalLinks":{"nodes":[` +
			commandEntityExternalLinkJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_history":
		return `{"project":{"id":"project-id","name":"Detail project","history":{"nodes":[` +
			commandProjectHistoryJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_initiativeToProjects":
		return `{"project":{"id":"project-id","name":"Detail project","initiativeToProjects":{"nodes":[` +
			commandInitiativeToProjectJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_initiatives":
		return `{"project":{"id":"project-id","name":"Detail project","initiatives":{"nodes":[` +
			commandInitiativeJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_inverseRelations":
		return `{"project":{"id":"project-id","name":"Detail project","inverseRelations":{"nodes":[` +
			commandProjectRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_issues":
		return `{"project":{"id":"project-id","name":"Detail project","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_comments":
		return `{"project":{"id":"project-id","name":"Detail project","comments":{"nodes":[` +
			commandCommentMetadataJSON("project-id", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_labels":
		return `{"project":{"id":"project-id","name":"Detail project","labels":{"nodes":[` +
			commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_members":
		if fake.emptyProjectMembers {
			return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_needs":
		return `{"project":{"id":"project-id","name":"Detail project","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_relations":
		return `{"project":{"id":"project-id","name":"Detail project","relations":{"nodes":[` +
			commandProjectRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_teams":
		return `{"project":{"id":"project-id","name":"Detail project","teams":{"nodes":[` +
			commandTeamJSON(true) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_projectUpdates":
		if fake.emptyProjectUpdates {
			return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[{"id":"project-update-id","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectFilterSuggestion":
		return `{"projectFilterSuggestion":{"filter":{"status":{"type":{"eq":"started"}}},"logId":"filter-log-id"}}`, true
	case "projectUpdates":
		if fake.emptyProjectUpdates {
			return `{"projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projectUpdates":{"nodes":[` + commandProjectUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectUpdate":
		return `{"projectUpdate":` + commandProjectUpdateJSON() + `}`, true
	case "projectUpdate_comments":
		return `{"projectUpdate":{"id":"project-update-id","comments":{"nodes":[` +
			commandCommentMetadataJSON("", "project-update-id") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_projectMilestones":
		if fake.emptyProjectMilestones {
			return `{"project":{"id":"project-id","name":"Detail project","projectMilestones":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectMilestones":
		if fake.emptyProjectMilestones {
			return `{"projectMilestones":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectMilestone":
		return `{"projectMilestone":` + commandProjectMilestoneJSON("Launch milestone", "next") + `}`, true
	case "projectMilestone_issues":
		return `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","issues":{"nodes":[` +
			commandIssueJSON("LIT-2", "Milestone issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectStatusPayload(operation string) (string, bool) {
	switch operation {
	case "projectStatuses":
		return `{"projectStatuses":{"nodes":[` +
			commandProjectStatusJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectStatus":
		return `{"projectStatus":` + commandProjectStatusJSON() + `}`, true
	case "projectStatusProjectCount":
		return `{"projectStatusProjectCount":{"count":12,"privateCount":2,"archivedTeamCount":1}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectLabelPayload(operation string) (string, bool) {
	switch operation {
	case "projectLabels":
		return `{"projectLabels":{"nodes":[` +
			commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectLabel":
		return `{"projectLabel":` + commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") + `}`, true
	case "projectLabel_children":
		return `{"projectLabel":{"id":"project-label-id","name":"Roadmap","children":{"nodes":[` +
			commandProjectLabelJSON("child-project-label-id", "Mobile", "#56ccf2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectLabel_projects":
		return `{"projectLabel":{"id":"project-label-id","name":"Roadmap","projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectRelationPayload(operation string) (string, bool) {
	switch operation {
	case "projectRelations":
		return `{"projectRelations":{"nodes":[` +
			commandProjectRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectRelation":
		return `{"projectRelation":` + commandProjectRelationJSON() + `}`, true
	default:
		return "", false
	}
}

func commandFlowProjectWritePayload(operation string) (string, bool) {
	switch operation {
	case "ProjectMilestoneCreate":
		return `{"projectMilestoneCreate":{"success":true,"projectMilestone":` + commandProjectMilestoneJSON("Created milestone", "next") + `}}`, true
	case "ProjectMilestoneUpdate":
		return `{"projectMilestoneUpdate":{"success":true,"projectMilestone":` + commandProjectMilestoneJSON("Updated milestone", "done") + `}}`, true
	case "ProjectCreate":
		return `{"projectCreate":{"success":true,"project":` + commandProjectJSON("Created project", "Backlog", "backlog") + `}}`, true
	case "ProjectUpdate":
		return `{"projectUpdate":{"success":true,"project":` + commandProjectJSON("Updated project", "Started", "started") + `}}`, true
	case "ProjectArchive":
		return `{"projectArchive":{"success":true,"entity":` + commandProjectJSON("Archived project", "Canceled", "canceled") + `}}`, true
	default:
		return "", false
	}
}
