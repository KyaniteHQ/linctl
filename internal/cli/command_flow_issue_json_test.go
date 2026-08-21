package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func commandIssueExportIdentifier(fake commandFlowFakeClient) string {
	if fake.invalidExportLeaf {
		return "../escape"
	}

	return "LIT-1"
}

func commandFlowWorkflowStatesByTeamJSON() string {
	return `{"workflowStates":{"nodes":[
		{"id":"todo-state","name":"Todo","type":"unstarted","position":1},
		{"id":"type-state-id","name":"Started","type":"started","position":0},
		{"id":"in-review-state","name":"In Review","type":"started","position":2},
		{"id":"done-state","name":"Done","type":"completed","position":1}
	],"pageInfo":{"hasNextPage":false}}}`
}

func commandFlowStateByID(stateID string) (name string, stateType string) {
	switch stateID {
	case "in-review-state":
		return "In Review", "started"
	case "done-state":
		return "Done", "completed"
	case "type-state-id", "started-state":
		return "Started", "started"
	default:
		return "Todo", "unstarted"
	}
}

func commandIssueJSON(identifier string, title string, stateID string, state string, stateType string) string {
	return commandIssueJSONWithTeam(identifier, title, stateID, state, stateType, "team-id", "LIT")
}

func commandIssueJSONWithID(
	id string,
	identifier string,
	title string,
	stateID string,
	state string,
	stateType string,
) string {
	return strings.Replace(
		commandIssueJSON(identifier, title, stateID, state, stateType),
		`"id":"issue-id"`,
		`"id":"`+id+`"`,
		1,
	)
}

func commandIssueJSONWithTeam(
	identifier string,
	title string,
	stateID string,
	state string,
	stateType string,
	teamID string,
	teamKey string,
) string {
	return `{
		"id":"issue-id",
		"description":"Existing description",
		"identifier":"` + identifier + `",
		"title":"` + title + `",
		"branchName":"` + strings.ToLower(identifier) + `-` + strings.ToLower(strings.ReplaceAll(title, " ", "-")) + `",
		"url":"https://linear.app/kyanite/issue/` + identifier + `",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"` + teamID + `","key":"` + teamKey + `","name":"` + teamKey + `","organization":{"id":"org-id"}},
		"state":{"id":"` + stateID + `","name":"` + state + `","type":"` + stateType + `"},
		"assignee":null,
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

// commandDestinationProjectID names the move-project destination in command flow
// tests. It differs from the pinned "project-id" so the move is a real move.
const commandDestinationProjectID = "eoir-project-id"

// commandDestinationProjectJSON is the move-project destination: another project
// on the pinned LIT team, so requireProjectTeam passes.
func commandDestinationProjectJSON() string {
	return strings.ReplaceAll(
		commandProjectJSON("EOIR Case Scraper", "Backlog", "backlog"),
		`"project-id"`,
		`"`+commandDestinationProjectID+`"`,
	)
}

func commandIssueJSONWithProject(
	identifier string,
	title string,
	projectID string,
	projectName string,
) string {
	return strings.Replace(
		commandIssueJSONWithTeam(identifier, title, "todo-state", "Todo", "unstarted", "team-id", "LIT"),
		`"project":{"id":"project-id","name":"Pinned project"}`,
		`"project":{"id":"`+projectID+`","name":"`+projectName+`"}`,
		1,
	)
}

func commandDestinationTeamJSON(teamID string, teamKey string) string {
	return `{
		"id":"` + teamID + `",
		"key":"` + teamKey + `",
		"name":"` + teamKey + `",
		"description":"destination team",
		"archivedAt":null,
		"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
	}`
}

type commandProjectTeam struct {
	ID   string
	Key  string
	Name string
}

func commandProjectJSONWithTeams(
	name string,
	status string,
	statusType string,
	teams []commandProjectTeam,
) string {
	nodes := make([]string, 0, len(teams))
	for _, team := range teams {
		nodes = append(nodes, `{"id":"`+team.ID+`","key":"`+team.Key+`","name":"`+team.Name+`"}`)
	}

	return `{
		"id":"project-id",
		"name":"` + name + `",
		"description":"description",
		"slugId":"` + name + `",
		"url":"https://linear.app/kyanite/project/project-id",
		"priority":0,
		"status":{"id":"status-id","name":"` + status + `","type":"` + statusType + `"},
		"lead":null,
		"teams":{"nodes":[` + strings.Join(nodes, ",") + `]}
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
		"endsAt":"2099-01-01T00:00:00Z",
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
	return commandProjectJSONWithTeams(name, status, statusType, []commandProjectTeam{
		{ID: "team-id", Key: "LIT", Name: "linctl"},
	})
}

func commandProjectDetailJSON(name string, status string, statusType string, content string) string {
	return strings.Replace(
		commandProjectJSON(name, status, statusType),
		`"description":"description",`,
		`"description":"description","content":"`+content+`",`,
		1,
	)
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
		"organization":{"id":"org-id"},
		"parent":null
	}`
}

func commandInitiativeLabelJSON(name string) string {
	return commandInitiativeLabelJSONWithOrg(name, "org-id")
}

func commandInitiativeLabelJSONWithOrg(name string, orgID string) string {
	return `{
		"id":"initiative-label-id",
		"name":"` + name + `",
		"description":"Strategic theme",
		"color":"#5e6ad2",
		"isGroup":false,
		"lastAppliedAt":"2026-07-10T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-07-01T12:00:00Z",
		"updatedAt":"2026-07-10T12:00:00Z",
		"organization":{"id":"` + orgID + `"},
		"parent":{"id":"initiative-label-group-id","name":"Themes","color":"#8a8f98"}
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
