package cli

func commandFlowBasePayload(operation string) (string, bool) {
	switch operation {
	case "Viewer":
		return `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`, true
	case "Teams":
		return `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "TargetProject":
		return `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`, true
	case "applicationInfo":
		return commandApplicationInfoPayload(), true
	case "agentActivities":
		return `{"agentActivities":{"nodes":[` + commandAgentActivityJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "agentActivity":
		return `{"agentActivity":` + commandAgentActivityJSON() + `}`, true
	case "agentSkills":
		return `{"agentSkills":{"nodes":[` + commandAgentSkillJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "agentSkill":
		return `{"agentSkill":` + commandAgentSkillJSON() + `}`, true
	case "externalUsers":
		return `{"externalUsers":{"nodes":[` + commandExternalUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "externalUser":
		return `{"externalUser":` + commandExternalUserJSON() + `}`, true
	case "rateLimitStatus":
		return commandRateLimitStatusPayload(), true
	default:
		return "", false
	}
}

func commandFlowOrganizationPayload(operation string) (string, bool) {
	switch operation {
	case "organizationExists":
		return `{"organizationExists":{"success":true,"exists":true}}`, true
	case "organization_labels":
		return `{"organization":{"labels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_projectLabels":
		return `{"organization":{"projectLabels":{"nodes":[` + commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_teams":
		return `{"organization":{"teams":{"nodes":[` + commandTeamJSON(false) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_templates":
		return `{"organization":{"templates":{"nodes":[` + commandTemplateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_users":
		return `{"organization":{"users":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowTeamMembershipPayload(operation string) (string, bool) {
	switch operation {
	case "teamMemberships":
		return `{"teamMemberships":{"nodes":[` + commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "teamMembership":
		return `{"teamMembership":` + commandTeamMembershipJSON() + `}`, true
	default:
		return "", false
	}
}

func commandRateLimitStatusPayload() string {
	return `{"rateLimitStatus":{"identifier":"api-key","kind":"api","limits":[{"type":"complexity","requestedAmount":1,"allowedAmount":1000,"period":60000,"remainingAmount":900,"reset":1720000000000}]}}`
}

func commandApplicationInfoPayload() string {
	return `{"applicationInfo":{"id":"app-id","clientId":"app-client-id","name":"Demo App","description":"Demo authorization app","developer":"Kyanite","developerUrl":"https://example.com","imageUrl":"https://example.com/app.png"}}`
}

func commandTeamMembershipJSON() string {
	return `{
		"id":"team-membership-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"owner":true,
		"sortOrder":1.5,
		"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandAgentActivityJSON() string {
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
		"content":{
			"__typename":"AgentActivityActionContent",
			"type":"action",
			"action":"read_file",
			"parameter":"README.md",
			"result":"Read file"
		}
	}`
}

func commandAgentSkillJSON() string {
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

func commandExternalUserJSON() string {
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
