package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func commandIssueJSON(identifier string, title string, stateID string, state string, stateType string) string {
	return `{
		"id":"issue-id",
		"description":"Existing description",
		"identifier":"` + identifier + `",
		"title":"` + title + `",
		"branchName":"` + strings.ToLower(identifier) + `-` + strings.ToLower(strings.ReplaceAll(title, " ", "-")) + `",
		"url":"https://linear.app/kyanite/issue/` + identifier + `",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"state":{"id":"` + stateID + `","name":"` + state + `","type":"` + stateType + `"},
		"assignee":null,
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func commandIssueRelationJSON() string {
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

func commandIssueToReleaseJSON() string {
	return `{
		"id":"issue-to-release-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id"},
		"release":{"id":"release-id"}
	}`
}

func commandIssueHistoryJSON() string {
	return `{
		"id":"issue-history-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"actorId":"user-id",
		"updatedDescription":true,
		"issue":{"id":"issue-id"}
	}`
}

func commandIssueStateSpanJSON() string {
	return `{
		"id":"issue-state-span-id",
		"stateId":"started-state",
		"startedAt":"2026-06-19T12:00:00Z",
		"endedAt":null,
		"state":{"id":"started-state","name":"Started","type":"started"}
	}`
}

func commandCycleJSON() string {
	return `{
		"id":"cycle-id",
		"number":7,
		"name":"Planning cycle",
		"description":"Cycle body",
		"startsAt":"2026-06-01T00:00:00Z",
		"endsAt":"2026-07-15T00:00:00Z",
		"completedAt":null,
		"progress":0.5,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandIssueWithNextRankJSON(
	identifier string,
	title string,
	priority int,
	priorityLabel string,
	createdAt string,
	unblocksCount int,
) string {
	return strings.TrimSuffix(commandIssueJSON(identifier, title, "todo-state", "Todo", "unstarted"), "\n\t}") +
		`,
		"priority":` + strconv.Itoa(priority) + `,
		"priorityLabel":"` + priorityLabel + `",
		"createdAt":"` + createdAt + `",
		"relations":{"nodes":[` + commandBlockingRelationsJSON(unblocksCount) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}
	}`
}

func commandBlockingRelationsJSON(count int) string {
	relations := make([]string, 0, count)
	for i := range count {
		relations = append(relations, fmt.Sprintf(`{"type":"blocks","relatedIssue":{"id":"blocked-%d","state":{"type":"unstarted"}}}`, i))
	}

	return strings.Join(relations, ",")
}

func commandProjectJSON(name string, status string, statusType string) string {
	return `{
		"id":"project-id",
		"name":"` + name + `",
		"description":"description",
		"slugId":"` + name + `",
		"url":"https://linear.app/kyanite/project/project-id",
		"priority":0,
		"status":{"id":"status-id","name":"` + status + `","type":"` + statusType + `"},
		"lead":null,
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}]}
	}`
}

func commandProjectUpdateJSON() string {
	return `{
		"id":"project-update-id",
		"body":"First update",
		"health":"onTrack",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"url":"https://linear.app/project-update/project-update-id",
		"project":{"id":"project-id","name":"Pinned project"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandProjectStatusJSON() string {
	return `{
		"id":"project-status-id",
		"name":"Backlog",
		"description":"Ready for planning",
		"type":"backlog",
		"color":"#bec2c8",
		"position":1,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z"
	}`
}

func commandProjectLabelJSON(id string, name string, color string) string {
	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":"Project label",
		"color":"` + color + `",
		"isGroup":false,
		"lastAppliedAt":"2026-06-19T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"parent":null
	}`
}

func commandProjectRelationJSON() string {
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

func commandProjectHistoryJSON() string {
	return `{
		"id":"project-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Backlog","to":"Started"}],
		"project":{"id":"project-id"}
	}`
}

func commandInitiativeUpdateJSON() string {
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

func commandInitiativeRelationJSON() string {
	return `{
		"id":"initiative-relation-id",
		"sortOrder":1.5,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"relatedInitiative":{"id":"child-initiative-id","name":"Child initiative"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandInitiativeToProjectJSON() string {
	return `{
		"id":"initiative-to-project-id",
		"sortOrder":"1",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}
	}`
}

func commandRoadmapToProjectJSON() string {
	return `{
		"id":"roadmap-to-project-id",
		"sortOrder":"1",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"roadmap":{"id":"roadmap-id","name":"Platform roadmap"},
		"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}
	}`
}

func commandProjectMilestoneJSON(name string, status string) string {
	return `{
		"id":"project-milestone-id",
		"name":"` + name + `",
		"description":"milestone body",
		"targetDate":"2026-06-30",
		"status":"` + status + `",
		"progress":0.5,
		"sortOrder":1,
		"project":` + commandProjectJSON("Pinned project", "Backlog", "backlog") + `
	}`
}

func commandDocumentJSON(title string, parents string) string {
	return `{
		"id":"document-id",
		"title":"` + title + `",
		"slugId":"document-slug",
		"archivedAt":null,
		` + parents + `
	}`
}

func commandLabelJSON(description string) string {
	return commandNamedLabelJSON("label-id", "Bug", "#ff0000", description)
}

func commandNamedLabelJSON(id string, name string, color string, description string) string {
	descriptionPayload := "null"
	if description != "" {
		descriptionPayload = `"` + description + `"`
	}

	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":` + descriptionPayload + `,
		"color":"` + color + `",
		"isGroup":false,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandTeamJSON(includeDescription bool) string {
	descriptionPayload := "null"
	if includeDescription {
		descriptionPayload = `"team body"`
	}

	return `{
		"id":"team-id",
		"key":"LIT",
		"name":"linctl",
		"description":` + descriptionPayload + `,
		"archivedAt":null,
		"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
	}`
}

func commandUserJSON() string {
	return `{
		"id":"user-id",
		"name":"omer",
		"displayName":"Omer",
		"email":"omer@example.com",
		"active":true,
		"guest":false,
		"admin":true
	}`
}

func commandDraftJSON() string {
	return `{
		"id":"draft-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id","identifier":"LIT-3","title":"Draft issue"},
		"project":null,
		"projectUpdate":null,
		"initiative":null,
		"initiativeUpdate":null,
		"parentComment":null,
		"customerNeed":null,
		"team":null
	}`
}

func commandFlowUserIssueListPayload(parent string, field string) string {
	return `{"` + parent + `":{"` + field + `":{"nodes":[` +
		commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
		`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`
}
