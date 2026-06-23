package client

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientReadScenarios_return_compact_lists_details_and_members(t *testing.T) {
	// Given
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"issues": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-20",
			Title:      "all-team issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeam": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-10",
			Title:      "listed issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"IssuesByTeamState": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-12",
			Title:      "started issue",
			StateID:    "started",
			State:      "Started",
			StateType:  "started",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamProject": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-14",
			Title:      "project issue",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamAssignee": `{"issues":{"nodes":[` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-15",
			Title:      "mine issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}, "Omer") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamLabel": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-16",
			Title:      "labeled issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamCycle": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-17",
			Title:      "cycle issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamCreatedAfter": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-18",
			Title:      "recent issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamCreatedBefore": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-19",
			Title:      "older issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamHasBlockers": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-21",
			Title:      "blocked issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssuesByTeamBlocks": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-22",
			Title:      "blocking issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"NextIssuesByTeam": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-27",
			Title:      "next issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssueBlockedIssues": `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[{"id":"relation-id","type":"blocks","relatedIssue":` + issueJSON(issueFixture{
			Identifier: "LIT-23",
			Title:      "blocked by issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
		"IssueDependencies": `{"issue":{
			"id":"issue-id",
			"identifier":"LIT-1",
			"parent":` + issueJSON(issueFixture{
			Identifier: "LIT-25",
			Title:      "parent issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `,
			"children":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-26",
			Title:      "child issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}},
			"relations":{"nodes":[{"id":"blocks-relation","type":"blocks","relatedIssue":` + issueJSON(issueFixture{
			Identifier: "LIT-23",
			Title:      "blocked issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}],"pageInfo":{"hasNextPage":false,"endCursor":null}},
			"inverseRelations":{"nodes":[{"id":"blocked-by-relation","type":"blocks","issue":` + issueJSON(issueFixture{
			Identifier: "LIT-24",
			Title:      "blocker issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}],"pageInfo":{"hasNextPage":false,"endCursor":null}}
		}}`,
		"issueSearch": `{"issueSearch":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-13",
			Title:      "search result",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"issueFigmaFileKeySearch": `{"issueFigmaFileKeySearch":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-41",
			Title:      "Figma issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issuePriorityValues":                     `{"issuePriorityValues":[{"priority":1,"label":"Urgent"},{"priority":0,"label":"No priority"}]}`,
		"issueFilterSuggestion":                   `{"issueFilterSuggestion":{"filter":{"state":{"type":{"eq":"started"}}},"logId":"issue-filter-log-id"}}`,
		"issueTitleSuggestionFromCustomerRequest": `{"issueTitleSuggestionFromCustomerRequest":{"title":"Improve exports","logId":"title-log-id"}}`,
		"searchDocuments": `{"searchDocuments":{"nodes":[` + strings.Join([]string{
			searchDocumentJSON(),
			searchDocumentJSONWithParent(
				"search-project-document-id",
				"Project search spec",
				`"project":{"id":"project-id","name":"Pinned project"},"initiative":null,"team":null,"issue":null,"release":null,"cycle":null`,
			),
			searchDocumentJSONWithParent(
				"search-initiative-document-id",
				"Initiative search spec",
				`"project":null,"initiative":{"id":"initiative-id","name":"Platform"},"team":null,"issue":null,"release":null,"cycle":null`,
			),
			searchDocumentJSONWithParent(
				"search-issue-document-id",
				"Issue search spec",
				`"project":null,"initiative":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-30","title":"Search issue"},"release":null,"cycle":null`,
			),
			searchDocumentJSONWithParent(
				"search-release-document-id",
				"Release search spec",
				`"project":null,"initiative":null,"team":null,"issue":null,"release":{"id":"release-id","name":"Mobile"},"cycle":null`,
			),
			searchDocumentJSONWithParent(
				"search-cycle-document-id",
				"Cycle search spec",
				`"project":null,"initiative":null,"team":null,"issue":null,"release":null,"cycle":{"id":"cycle-id","number":12,"name":"Planning cycle"}`,
			),
		}, ",") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"},"totalCount":6}}`,
		"searchIssues":   `{"searchIssues":{"nodes":[` + searchIssueJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"},"totalCount":1}}`,
		"searchProjects": `{"searchProjects":{"nodes":[` + searchProjectJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"},"totalCount":1}}`,
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-11",
			Title:      "detail issue",
			StateID:    "done",
			State:      "Done",
			StateType:  "completed",
		}) + `}`,
		"issue_attachments":       `{"issue":{"attachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_botActor":          `{"issue":{"id":"issue-id","botActor":` + actorBotJSON() + `}}`,
		"issue_children":          `{"issue":{"children":{"nodes":[` + issueJSON(issueFixture{Identifier: "LIT-31", Title: "child issue", StateID: "todo", State: "Todo", StateType: "unstarted"}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_documents":         `{"issue":{"documents":{"nodes":[{"id":"issue-document-id","title":"Issue spec","slugId":"issue-spec","archivedAt":null,"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_formerAttachments": `{"issue":{"formerAttachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_history":           `{"issue":{"history":{"nodes":[` + issueHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_inverseRelations":  `{"issue":{"inverseRelations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_labels":            `{"issue":{"labels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_relations":         `{"issue":{"relations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_releases":          `{"issue":{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_stateHistory":      `{"issue":{"id":"issue-id","stateHistory":{"nodes":[` + issueStateSpanJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_subscribers":       `{"issue":{"id":"issue-id","subscribers":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch": `{"issueVcsBranchSearch":` + issueJSON(issueFixture{
			Identifier: "LIT-40",
			Title:      "branch issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"issueVcsBranchSearch_attachments":       `{"issueVcsBranchSearch":{"attachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_botActor":          `{"issueVcsBranchSearch":{"id":"issue-id","botActor":` + actorBotJSON() + `}}`,
		"issueVcsBranchSearch_children":          `{"issueVcsBranchSearch":{"children":{"nodes":[` + issueJSON(issueFixture{Identifier: "LIT-43", Title: "branch child issue", StateID: "todo", State: "Todo", StateType: "unstarted"}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_documents":         `{"issueVcsBranchSearch":{"documents":{"nodes":[{"id":"branch-issue-document-id","title":"Branch issue spec","slugId":"branch-issue-spec","archivedAt":null,"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_formerAttachments": `{"issueVcsBranchSearch":{"formerAttachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_history":           `{"issueVcsBranchSearch":{"history":{"nodes":[` + issueHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_inverseRelations":  `{"issueVcsBranchSearch":{"inverseRelations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_labels":            `{"issueVcsBranchSearch":{"labels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_relations":         `{"issueVcsBranchSearch":{"relations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_releases":          `{"issueVcsBranchSearch":{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_stateHistory":      `{"issueVcsBranchSearch":{"id":"issue-id","stateHistory":{"nodes":[` + issueStateSpanJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueVcsBranchSearch_subscribers":       `{"issueVcsBranchSearch":{"id":"issue-id","subscribers":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"Projects": `{"team":{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projects": `{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "workspace-project-id",
			Name:   "workspace listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"project": `{"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "detail",
			Status: "Started",
		}) + `}`,
		"project_members":              `{"project":{"id":"project-id","name":"detail","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_attachments":          `{"project":{"id":"project-id","name":"detail","attachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_documents":            `{"project":{"id":"project-id","name":"detail","documents":{"nodes":[{"id":"project-document-id","title":"Project spec","slugId":"project-spec","archivedAt":null,"project":{"id":"project-id","name":"detail"},"team":null,"issue":null,"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_externalLinks":        `{"project":{"id":"project-id","name":"detail","externalLinks":{"nodes":[` + entityExternalLinkJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_history":              `{"project":{"id":"project-id","name":"detail","history":{"nodes":[` + projectHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_initiativeToProjects": `{"project":{"id":"project-id","name":"detail","initiativeToProjects":{"nodes":[{"id":"initiative-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_initiatives":          `{"project":{"id":"project-id","name":"detail","initiatives":{"nodes":[{"id":"initiative-id","name":"Platform","description":"Platform initiative","status":"Active","priority":2,"targetDate":"2026-12-31","slugId":"platform-init","url":"https://linear.app/kyanite/initiative/platform-init"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_inverseRelations":     `{"project":{"id":"project-id","name":"detail","inverseRelations":{"nodes":[` + projectRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_issues": `{"project":{"id":"project-id","name":"detail","issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-47",
			Title:      "Project issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_comments":          `{"project":{"id":"project-id","name":"detail","comments":{"nodes":[` + commentMetadataJSON("project-id", "", "user-id") + `,` + commentMetadataJSON("project-id", "", "") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_labels":            `{"project":{"id":"project-id","name":"detail","labels":{"nodes":[` + projectLabelJSON("project-label-id", "Roadmap") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_needs":             `{"project":{"id":"project-id","name":"detail","needs":{"nodes":[` + customerNeedJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_relations":         `{"project":{"id":"project-id","name":"detail","relations":{"nodes":[` + projectRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_teams":             `{"project":{"id":"project-id","name":"detail","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_projectUpdates":    `{"project":{"id":"project-id","name":"detail","projectUpdates":{"nodes":[{"id":"project-update-id","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectFilterSuggestion":   `{"projectFilterSuggestion":{"filter":{"status":{"type":{"eq":"started"}}},"logId":"filter-log-id"}}`,
		"projectUpdates":            `{"projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","project":{"id":"project-id","name":"detail"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectUpdate":             `{"projectUpdate":{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","project":{"id":"project-id","name":"detail"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"projectUpdate_comments":    `{"projectUpdate":{"id":"project-update-id","comments":{"nodes":[` + commentMetadataJSON("", "project-update-id", "user-id") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_projectMilestones": `{"project":{"id":"project-id","name":"detail","projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectMilestones":         `{"projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectMilestone":          `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}}`,
		"projectMilestone_issues": `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-52",
			Title:      "Milestone issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectStatuses":           `{"projectStatuses":{"nodes":[{"id":"project-status-id","name":"Backlog","description":"Ready for planning","type":"backlog","color":"#bec2c8","position":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectStatus":             `{"projectStatus":{"id":"project-status-id","name":"Backlog","description":"Ready for planning","type":"backlog","color":"#bec2c8","position":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z"}}`,
		"projectStatusProjectCount": `{"projectStatusProjectCount":{"count":12,"privateCount":2,"archivedTeamCount":1}}`,
		"projectLabels":             `{"projectLabels":{"nodes":[` + projectLabelJSON("project-label-id", "Roadmap") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectLabel":              `{"projectLabel":{"id":"project-label-id","name":"Roadmap","description":"Project label","color":"#f2c94c","isGroup":false,"lastAppliedAt":"2026-06-19T12:00:00Z","retiredAt":null,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","parent":{"id":"parent-project-label-id","name":"Parent","color":"#828282"}}}`,
		"projectLabel_children":     `{"projectLabel":{"id":"project-label-id","name":"Roadmap","children":{"nodes":[{"id":"child-project-label-id","name":"Mobile","description":"Child project label","color":"#56ccf2","isGroup":false,"lastAppliedAt":null,"retiredAt":null,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","parent":{"id":"project-label-id","name":"Roadmap","color":"#f2c94c"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectLabel_projects": `{"projectLabel":{"id":"project-label-id","name":"Roadmap","projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectRelations": `{"projectRelations":{"nodes":[{"id":"project-relation-id","type":"blocks","anchorType":"project","relatedAnchorType":"project","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":null,"relatedProject":{"id":"related-project-id","name":"Related project"},"relatedProjectMilestone":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"project-relation-no-user","type":"relates","anchorType":"project","relatedAnchorType":"project","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone"},"relatedProject":{"id":"other-related-project-id","name":"Other related"},"relatedProjectMilestone":{"id":"related-project-milestone-id","name":"Related milestone"},"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectRelation":  `{"projectRelation":{"id":"project-relation-id","type":"blocks","anchorType":"project","relatedAnchorType":"project","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":null,"relatedProject":{"id":"related-project-id","name":"Related project"},"relatedProjectMilestone":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"issueRelations":   `{"issueRelations":{"nodes":[{"id":"issue-relation-id","type":"blocks","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issueRelation":    `{"issueRelation":{"id":"issue-relation-id","type":"blocks","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}}`,
		"issueToReleases":  `{"issueToReleases":{"nodes":[{"id":"issue-to-release-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id"},"release":{"id":"release-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issueToRelease":   `{"issueToRelease":{"id":"issue-to-release-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id"},"release":{"id":"release-id"}}}`,
		"applicationInfo":  `{"applicationInfo":{"id":"app-id","clientId":"app-client-id","name":"Demo App","description":"Demo authorization app","developer":"Kyanite","developerUrl":"https://example.com","imageUrl":"https://example.com/app.png"}}`,
		"issue_comments":   `{"issue":{"id":"issue-id","identifier":"LIT-12","comments":{"nodes":[{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":"parent-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"bot-comment-id","body":"bot note","url":"https://linear.app/comment/bot-comment-id","createdAt":"2026-06-19T12:01:00Z","parentId":null,"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"comments":         `{"comments":{"nodes":[{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:02:00Z","editedAt":"2026-06-19T12:02:00Z","resolvedAt":null,"parentId":"parent-id","issueId":"issue-id","projectId":null,"projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"bot-comment-id","body":"bot note","url":"https://linear.app/comment/bot-comment-id","createdAt":"2026-06-19T12:01:00Z","updatedAt":"2026-06-19T12:01:00Z","editedAt":null,"resolvedAt":null,"parentId":null,"issueId":null,"projectId":"project-id","projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"comment":          `{"comment":{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:02:00Z","editedAt":"2026-06-19T12:02:00Z","resolvedAt":null,"parentId":"parent-id","issueId":"issue-id","projectId":null,"projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"comment_botActor": `{"comment":{"id":"comment-id","botActor":{"id":"bot-actor-id","type":"github","subType":"issue","name":"GitHub","userDisplayName":"octocat","avatarUrl":"https://example.com/github.png"}}}`,
		"comment_children": `{"comment":{"id":"comment-id","children":{"nodes":[` + commentMetadataJSONWithID("child-comment-id", "comment-id", "", "", "user-id") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"comment_createdIssues": `{"comment":{"id":"comment-id","createdIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-54",
			Title:      "Created from comment",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"Documents":           `{"documents":{"nodes":[{"id":"document-id","title":"Spec","slugId":"spec","archivedAt":null,"project":{"id":"project-id","name":"fixture"},"team":null,"issue":null,"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"document":            `{"document":{"id":"document-id","title":"Team note","slugId":"team-note","archivedAt":null,"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null}}`,
		"document_comments":   `{"document":{"id":"document-id","comments":{"nodes":[` + commentMetadataJSON("", "", "user-id") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"IssueLabels":         `{"issueLabels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issueLabel":          `{"issueLabel":{"id":"label-id","name":"Bug","description":null,"color":"#ff0000","isGroup":false,"team":null}}`,
		"issueLabel_children": `{"issueLabel":{"id":"label-id","name":"Bug","children":{"nodes":[{"id":"child-label-id","name":"Mobile","description":"child label","color":"#56ccf2","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issueLabel_issues": `{"issueLabel":{"id":"label-id","name":"Bug","issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-53",
			Title:      "Label issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"Teams":           `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"teamMemberships": `{"teamMemberships":{"nodes":[{"id":"team-membership-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"owner":true,"sortOrder":1.5,"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"teamMembership":  `{"teamMembership":{"id":"team-membership-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"owner":true,"sortOrder":1.5,"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},"team":{"id":"team-id","key":"LIT","name":"linctl"}}}`,
		"team_cycles":     `{"team":{"id":"team-id","key":"LIT","name":"linctl","cycles":{"nodes":[` + cycleJSON("Planning cycle", "team-id", "LIT") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_issues": `{"team":{"id":"team-id","key":"LIT","name":"linctl","issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "Team issue",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "backlog",
			Project:    "Team project",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_labels":           `{"team":{"id":"team-id","key":"LIT","name":"linctl","labels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_memberships":      `{"team":{"id":"team-id","key":"LIT","name":"linctl","memberships":{"nodes":[{"id":"team-membership-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"owner":true,"sortOrder":1.5,"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_projects":         `{"team":{"id":"team-id","key":"LIT","name":"linctl","projects":{"nodes":[` + projectJSON(projectFixture{ID: "project-id", Name: "Team project", Status: "Backlog"}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_releasePipelines": `{"team":{"id":"team-id","key":"LIT","name":"linctl","releasePipelines":{"nodes":[` + releasePipelineJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_states":           `{"team":{"id":"team-id","key":"LIT","name":"linctl","states":{"nodes":[{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_gitAutomationStates": `{"team":{"id":"team-id","key":"LIT","name":"linctl","gitAutomationStates":{"nodes":[` +
			gitAutomationStateJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"team_templates": `{"team":{"id":"team-id","key":"LIT","name":"linctl","templates":{"nodes":[` + templateJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"agentActivities": `{"agentActivities":{"nodes":[` + strings.Join([]string{
			agentActivityJSON("action"),
			agentActivityJSON("elicitation"),
			agentActivityJSON("error"),
			agentActivityJSON("prompt"),
			agentActivityJSON("response"),
			agentActivityJSON("thought"),
		}, ",") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"agentActivity":              `{"agentActivity":` + agentActivityJSON("action") + `}`,
		"agentSkills":                `{"agentSkills":{"nodes":[` + agentSkillJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"agentSkill":                 `{"agentSkill":` + agentSkillJSON() + `}`,
		"externalUsers":              `{"externalUsers":{"nodes":[` + externalUserJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"externalUser":               `{"externalUser":` + externalUserJSON() + `}`,
		"auditEntryTypes":            `{"auditEntryTypes":[{"type":"user_login","description":"User logged in"}]}`,
		"organizationExists":         `{"organizationExists":{"success":true,"exists":true}}`,
		"organization_labels":        `{"organization":{"labels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"organization_projectLabels": `{"organization":{"projectLabels":{"nodes":[{"id":"project-label-id","name":"Roadmap","description":"Project label","color":"#f2c94c","isGroup":false,"lastAppliedAt":"2026-06-19T12:00:00Z","retiredAt":null,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","parent":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"organization_teams":         `{"organization":{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","description":null,"archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"organization_templates":     `{"organization":{"templates":{"nodes":[` + templateJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"organization_users":         `{"organization":{"users":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"rateLimitStatus":            `{"rateLimitStatus":{"identifier":"api-key","kind":"api","limits":[{"type":"complexity","requestedAmount":1,"allowedAmount":1000,"period":60000,"remainingAmount":900,"reset":1720000000000}]}}`,
		"notifications":              `{"notifications":{"nodes":[` + notificationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"notification":               `{"notification":` + notificationJSON() + `}`,
		"notificationSubscriptions": `{"notificationSubscriptions":{"nodes":[` + strings.Join([]string{
			notificationSubscriptionJSON(),
			notificationSubscriptionTargetJSON("CustomerNotificationSubscription", "customer", `{"id":"customer-id","name":"Acme"}`, false, false),
			notificationSubscriptionTargetJSON("CustomViewNotificationSubscription", "customView", `{"id":"custom-view-id","name":"My issues"}`, false, false),
			notificationSubscriptionTargetJSON("CycleNotificationSubscription", "cycle", `{"id":"cycle-id","name":"Cycle 7"}`, false, false),
			notificationSubscriptionTargetJSON("InitiativeNotificationSubscription", "initiative", `{"id":"initiative-id","name":"Platform"}`, false, false),
			notificationSubscriptionTargetJSON("LabelNotificationSubscription", "label", `{"id":"label-id","name":"Bug"}`, false, false),
			notificationSubscriptionTargetJSON("TeamNotificationSubscription", "team", `{"id":"team-id","key":"LIT","name":"linctl"}`, true, false),
			notificationSubscriptionTargetJSON("UserNotificationSubscription", "user", `{"id":"target-user-id","displayName":"Ada"}`, false, true),
		}, ",") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"notificationSubscription":             `{"notificationSubscription":` + notificationSubscriptionJSON() + `}`,
		"triageResponsibilities":               `{"triageResponsibilities":{"nodes":[` + triageResponsibilityJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"triageResponsibility":                 `{"triageResponsibility":` + triageResponsibilityJSON() + `}`,
		"triageResponsibility_manualSelection": `{"triageResponsibility":{"id":"triage-responsibility-id","manualSelection":{"userIds":["user-id","other-user-id"]}}}`,
		"slaConfigurations":                    `{"slaConfigurations":[` + slaConfigurationJSON() + `]}`,
		"semanticSearch": `{"semanticSearch":{"results":[` + strings.Join([]string{
			semanticSearchResultJSON("issue"),
			semanticSearchResultJSON("project"),
			semanticSearchResultJSON("initiative"),
			semanticSearchResultJSON("document"),
			semanticSearchResultJSON("unknown"),
		}, ",") + `]}}`,
		"releasePipelines":         `{"releasePipelines":{"nodes":[` + releasePipelineJSON() + `,` + trashedReleasePipelineJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"releasePipeline":          `{"releasePipeline":` + releasePipelineJSON() + `}`,
		"releasePipeline_releases": `{"releasePipeline":{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"releasePipeline_stages":   `{"releasePipeline":{"stages":{"nodes":[` + releaseStageJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"releasePipeline_teams":    `{"releasePipeline":{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"releaseStages":            `{"releaseStages":{"nodes":[` + releaseStageJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"releaseStage":             `{"releaseStage":` + releaseStageJSON() + `}`,
		"releaseStage_releases":    `{"releaseStage":{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"releases":                 `{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"release":                  `{"release":` + releaseJSON() + `}`,
		"releaseSearch":            `{"releaseSearch":[` + releaseJSON() + `]}`,
		"release_history":          `{"release":{"history":{"nodes":[` + releaseHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"release_documents":        `{"release":{"documents":{"nodes":[{"id":"release-document-id","title":"Release spec","slugId":"release-spec","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"release_issues": `{"release":{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-48",
			Title:      "Release issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"release_links":      `{"release":{"links":{"nodes":[` + entityExternalLinkJSON() + `,` + entityExternalLinkWithParentsJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"entityExternalLink": `{"entityExternalLink":` + entityExternalLinkJSON() + `}`,
		"releaseNotes":       `{"releaseNotes":{"nodes":[` + releaseNoteJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"releaseNote":        `{"releaseNote":` + releaseNoteJSON() + `}`,
		"team":               `{"team":{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"team_members":       `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"users":              `{"users":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"user":               `{"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}}`,
		"viewer":             `{"viewer":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}}`,
		"viewer_drafts": `{"viewer":{"drafts":{"nodes":[` + strings.Join([]string{
			draftJSON("issue"),
			draftJSON("project"),
			draftJSON("project_update"),
			draftJSON("initiative"),
			draftJSON("initiative_update"),
			draftJSON("comment"),
			draftJSON("customer_need"),
			draftJSON("team"),
			draftJSON("unknown"),
		}, ",") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"user_assignedIssues": `{"user":{"assignedIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-41",
			Title:      "User assigned issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"user_createdIssues": `{"user":{"createdIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-42",
			Title:      "User created issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"user_delegatedIssues": `{"user":{"delegatedIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-43",
			Title:      "User delegated issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"user_teamMemberships": `{"user":{"teamMemberships":{"nodes":[{"id":"team-membership-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"owner":true,"sortOrder":1.5,"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"user_teams":           `{"user":{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"viewer_assignedIssues": `{"viewer":{"assignedIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-44",
			Title:      "Viewer assigned issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"viewer_createdIssues": `{"viewer":{"createdIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-45",
			Title:      "Viewer created issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"viewer_delegatedIssues": `{"viewer":{"delegatedIssues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-46",
			Title:      "Viewer delegated issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"viewer_teamMemberships": `{"viewer":{"teamMemberships":{"nodes":[{"id":"team-membership-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"owner":true,"sortOrder":1.5,"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"viewer_teams":           `{"viewer":{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"workflowStates":         `{"workflowStates":{"nodes":[{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"workflowState":          `{"workflowState":{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}}`,
		"workflowState_issues": `{"workflowState":{"id":"workflow-state-id","name":"Started","issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-55",
			Title:      "Workflow state issue",
			StateID:    "workflow-state-id",
			State:      "Started",
			StateType:  "started",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"timeSchedules":                `{"timeSchedules":{"nodes":[{"id":"time-schedule-id","name":"Primary on-call","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,"externalId":"pd-primary","externalUrl":"https://example.com/schedule","integration":{"id":"integration-id"},"entries":[{"startsAt":"2026-06-20T00:00:00Z","endsAt":"2026-06-21T00:00:00Z","userId":"user-id","userEmail":"omer@example.com"}]}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"timeSchedule":                 `{"timeSchedule":{"id":"time-schedule-id","name":"Primary on-call","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,"externalId":"pd-primary","externalUrl":"https://example.com/schedule","integration":{"id":"integration-id"},"entries":[{"startsAt":"2026-06-20T00:00:00Z","endsAt":"2026-06-21T00:00:00Z","userId":"user-id","userEmail":"omer@example.com"}]}}`,
		"templates":                    `{"templates":[` + templateJSON() + `,` + strings.Replace(templateJSON(), "template-id", "template-two-id", 1) + `]}`,
		"template":                     `{"template":` + templateJSON() + `}`,
		"initiatives":                  `{"initiatives":{"nodes":[{"id":"initiative-id","name":"Platform","description":"Platform initiative","status":"Active","priority":2,"targetDate":"2026-12-31","slugId":"platform-init","url":"https://linear.app/kyanite/initiative/platform-init"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiative":                   `{"initiative":{"id":"initiative-id","name":"Platform","description":"Platform initiative","status":"Active","priority":2,"targetDate":"2026-12-31","slugId":"platform-init","url":"https://linear.app/kyanite/initiative/platform-init"}}`,
		"initiative_history":           `{"initiative":{"history":{"nodes":[` + initiativeHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"initiative_links":             `{"initiative":{"links":{"nodes":[` + entityExternalLinkJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"initiative_subInitiatives":    `{"initiative":{"subInitiatives":{"nodes":[` + subInitiativeJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"initiative_initiativeUpdates": `{"initiative":{"initiativeUpdates":{"nodes":[` + initiativeUpdateJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"initiative_documents":         `{"initiative":{"documents":{"nodes":[{"id":"initiative-document-id","title":"Initiative spec","slugId":"initiative-spec","archivedAt":null,"project":null,"team":null,"issue":null,"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"initiative_projects": `{"initiative":{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "initiative-project-id",
			Name:   "Initiative project",
			Status: "Started",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"initiativeRelations":  `{"initiativeRelations":{"nodes":[{"id":"initiative-relation-id","sortOrder":1.5,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"relatedInitiative":{"id":"child-initiative-id","name":"Child initiative"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"initiative-relation-no-user","sortOrder":2,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"relatedInitiative":{"id":"other-child-initiative-id","name":"Other child"},"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeRelation":   `{"initiativeRelation":{"id":"initiative-relation-id","sortOrder":1.5,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"relatedInitiative":{"id":"child-initiative-id","name":"Child initiative"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"initiativeToProjects": `{"initiativeToProjects":{"nodes":[{"id":"initiative-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeToProject":  `{"initiativeToProject":{"id":"initiative-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}}`,
		"roadmapToProjects":    `{"roadmapToProjects":{"nodes":[{"id":"roadmap-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"roadmap":{"id":"roadmap-id","name":"Platform roadmap"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"roadmapToProject":     `{"roadmapToProject":{"id":"roadmap-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"roadmap":{"id":"roadmap-id","name":"Platform roadmap"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}}`,
		"initiativeUpdates":    `{"initiativeUpdates":{"nodes":[{"id":"initiative-update-id","body":"First initiative update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/initiative-update/initiative-update-id","slugId":"initiative-update-slug","commentCount":1,"initiative":{"id":"initiative-id","name":"Platform"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeUpdate":     `{"initiativeUpdate":{"id":"initiative-update-id","body":"First initiative update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/initiative-update/initiative-update-id","slugId":"initiative-update-slug","commentCount":1,"initiative":{"id":"initiative-id","name":"Platform"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"initiativeUpdate_comments": `{"initiativeUpdate":{"id":"initiative-update-id","comments":{"nodes":[` +
			commentMetadataJSON("", "", "user-id") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"roadmaps":                 `{"roadmaps":{"nodes":[{"id":"roadmap-id","name":"Platform roadmap","description":"Roadmap body","color":"#5e6ad2","slugId":"platform-roadmap","sortOrder":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","url":"https://linear.app/kyanite/roadmap/platform-roadmap","creator":{"id":"user-id","displayName":"Omer"},"owner":{"id":"owner-id","displayName":"Owner"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"roadmap":                  `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","description":"Roadmap body","color":"#5e6ad2","slugId":"platform-roadmap","sortOrder":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","url":"https://linear.app/kyanite/roadmap/platform-roadmap","creator":{"id":"user-id","displayName":"Omer"},"owner":{"id":"owner-id","displayName":"Owner"}}}`,
		"roadmap_projects":         `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","projects":{"nodes":[` + projectJSON(projectFixture{ID: "roadmap-project-id", Name: "Roadmap project", Status: "Backlog"}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"customViews":              `{"customViews":{"nodes":[{"id":"custom-view-id","name":"My issues","description":"Saved issue view","modelName":"Issue","shared":true,"color":"#5e6ad2","slugId":"my-issues"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"customViewHasSubscribers": `{"customViewHasSubscribers":{"hasSubscribers":true}}`,
		"customView":               `{"customView":{"id":"custom-view-id","name":"My issues","description":"Saved issue view","modelName":"Issue","shared":true,"color":"#5e6ad2","slugId":"my-issues"}}`,
		"customView_initiatives":   `{"customView":{"initiatives":{"nodes":[{"id":"initiative-id","name":"Platform","description":"Platform initiative","status":"Active","priority":2,"targetDate":"2026-12-31","slugId":"platform-init","url":"https://linear.app/kyanite/initiative/platform-init"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"customView_issues": `{"customView":{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "Custom view issue",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "backlog",
			Project:    "Custom view project",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"customView_organizationViewPreferences":             `{"customView":{"organizationViewPreferences":` + customViewPreferencesJSON("priority", "list") + `}}`,
		"customView_organizationViewPreferences_preferences": `{"customView":{"organizationViewPreferences":{"preferences":` + customViewPreferenceValuesJSON("priority", "list") + `}}}`,
		"customView_projects": `{"customView":{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "Custom view project",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"customView_userViewPreferences":             `{"customView":{"userViewPreferences":` + customViewScopedPreferencesJSON("user", "updatedAt", "board") + `}}`,
		"customView_userViewPreferences_preferences": `{"customView":{"userViewPreferences":{"preferences":` + customViewPreferenceValuesJSON("updatedAt", "board") + `}}}`,
		"customView_viewPreferencesValues":           `{"customView":{"viewPreferencesValues":` + customViewPreferenceValuesJSON("updatedAt", "board") + `}}`,
		"customers":                                  `{"customers":{"nodes":[{"id":"customer-id","name":"Acme","domains":["acme.example"],"externalIds":["crm-acme"],"slackChannelId":"slack-channel-id","status":{"id":"status-id","name":"Active"},"tier":{"id":"tier-id","name":"Enterprise"},"owner":{"id":"user-id","displayName":"Omer"},"revenue":120000,"size":42,"approximateNeedCount":3,"slugId":"acme","url":"https://linear.app/kyanite/customer/acme"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"customer":                                   `{"customer":{"id":"customer-id","name":"Acme","domains":["acme.example"],"externalIds":["crm-acme"],"slackChannelId":"slack-channel-id","status":{"id":"status-id","name":"Active"},"tier":{"id":"tier-id","name":"Enterprise"},"owner":{"id":"user-id","displayName":"Omer"},"revenue":120000,"size":42,"approximateNeedCount":3,"slugId":"acme","url":"https://linear.app/kyanite/customer/acme"}}`,
		"customerNeeds":                              `{"customerNeeds":{"nodes":[{"id":"customer-need-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,"priority":1,"body":"Need body","content":"Need content","url":"https://example.com/need","customer":{"id":"customer-id","name":"Acme"},"issue":{"id":"issue-id","identifier":"LIT-1","title":"Need issue"},"project":{"id":"project-id","name":"Customer project"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"customerNeed":                               `{"customerNeed":{"id":"customer-need-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,"priority":1,"body":"Need body","content":"Need content","url":"https://example.com/need","customer":{"id":"customer-id","name":"Acme"},"issue":{"id":"issue-id","identifier":"LIT-1","title":"Need issue"},"project":{"id":"project-id","name":"Customer project"}}}`,
		"customerNeed_projectAttachment":             `{"customerNeed":{"id":"customer-need-id","projectAttachment":` + projectAttachmentJSON() + `}}`,
		"customerStatuses":                           `{"customerStatuses":{"nodes":[{"id":"customer-status-id","name":"active","displayName":"Active","color":"#00ff00","description":"Active customers","position":1,"archivedAt":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"customerStatus":                             `{"customerStatus":{"id":"customer-status-id","name":"active","displayName":"Active","color":"#00ff00","description":"Active customers","position":1,"archivedAt":null}}`,
		"customerTiers":                              `{"customerTiers":{"nodes":[{"id":"customer-tier-id","name":"enterprise","displayName":"Enterprise","color":"#0000ff","description":"Enterprise customers","position":2,"archivedAt":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"customerTier":                               `{"customerTier":{"id":"customer-tier-id","name":"enterprise","displayName":"Enterprise","color":"#0000ff","description":"Enterprise customers","position":2,"archivedAt":null}}`,
		"favorites":                                  `{"favorites":{"nodes":[{"id":"favorite-id","type":"issue","folderName":null,"url":"https://linear.app/kyanite/issue/LIT-1"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"favorite_children":                          `{"favorite":{"children":{"nodes":[{"id":"favorite-child-id","type":"project","folderName":null,"url":"https://linear.app/kyanite/project/project-id"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"favorite":                                   `{"favorite":{"id":"favorite-id","type":"issue","folderName":null,"url":"https://linear.app/kyanite/issue/LIT-1"}}`,
		"emojis":                                     `{"emojis":{"nodes":[{"id":"emoji-id","name":"party","url":"https://linear.app/kyanite/emoji/party.png","source":"custom"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"emoji":                                      `{"emoji":{"id":"emoji-id","name":"party","url":"https://linear.app/kyanite/emoji/party.png","source":"custom"}}`,
		"attachments":                                `{"attachments":{"nodes":[{"id":"attachment-id","title":"Linked PR","subtitle":"feat: add thing","url":"https://github.com/kyanite/linctl/pull/1","sourceType":"github"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"attachmentsForURL":                          `{"attachmentsForURL":{"nodes":[{"id":"attachment-url-id","title":"Linked URL","subtitle":"url source","url":"https://example.com/spec","sourceType":"url"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"attachment":                                 `{"attachment":{"id":"attachment-id","title":"Linked PR","subtitle":"feat: add thing","url":"https://github.com/kyanite/linctl/pull/1","sourceType":"github"}}`,
		"attachmentIssue": `{"attachmentIssue":` + issueJSON(issueFixture{
			Identifier: "LIT-41",
			Title:      "attachment issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"attachmentIssue_attachments":       `{"attachmentIssue":{"attachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_botActor":          `{"attachmentIssue":{"id":"issue-id","botActor":` + actorBotJSON() + `}}`,
		"attachmentIssue_children":          `{"attachmentIssue":{"children":{"nodes":[` + issueJSON(issueFixture{Identifier: "LIT-42", Title: "attachment child issue", StateID: "todo", State: "Todo", StateType: "unstarted"}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_documents":         `{"attachmentIssue":{"documents":{"nodes":[{"id":"attachment-issue-document-id","title":"Attachment issue spec","slugId":"attachment-issue-spec","archivedAt":null,"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_formerAttachments": `{"attachmentIssue":{"formerAttachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_history":           `{"attachmentIssue":{"history":{"nodes":[` + issueHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_inverseRelations":  `{"attachmentIssue":{"inverseRelations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_labels":            `{"attachmentIssue":{"labels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_relations":         `{"attachmentIssue":{"relations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_releases":          `{"attachmentIssue":{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_stateHistory":      `{"attachmentIssue":{"id":"issue-id","stateHistory":{"nodes":[` + issueStateSpanJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"attachmentIssue_subscribers":       `{"attachmentIssue":{"id":"issue-id","subscribers":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
	}

	// When
	allTeamIssues, err := ListIssues(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	issues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{})
	require.NoError(t, err)
	startedIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{StateType: "started"})
	require.NoError(t, err)
	projectIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{ProjectID: "project-id"})
	require.NoError(t, err)
	mineIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{AssigneeID: "user-id"})
	require.NoError(t, err)
	labelIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{LabelID: "label-id"})
	require.NoError(t, err)
	cycleIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{CycleID: "cycle-id"})
	require.NoError(t, err)
	recentIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{CreatedAfter: "2026-06-01"})
	require.NoError(t, err)
	olderIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{CreatedBefore: "2026-06-30"})
	require.NoError(t, err)
	blockedIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{HasBlockers: true})
	require.NoError(t, err)
	blockingIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{Blocks: true})
	require.NoError(t, err)
	nextIssues, err := ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	blockedByIssues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2, IssueListFilters{BlockedBy: "LIT-1"})
	require.NoError(t, err)
	dependencies, err := GetIssueDependencies(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	searchIssues, err := SearchIssuesByTeam(context.Background(), graphqlClient, "team-id", "needle", 2)
	require.NoError(t, err)
	figmaIssues, err := SearchIssuesByFigmaFileKey(context.Background(), graphqlClient, "figma-key", 2)
	require.NoError(t, err)
	issuePriorityValues, err := ListIssuePriorityValues(context.Background(), graphqlClient)
	require.NoError(t, err)
	issueFilterSuggestion, err := GetIssueFilterSuggestion(
		context.Background(),
		graphqlClient,
		"started issues",
		"team-id",
		"",
	)
	require.NoError(t, err)
	issueTitleSuggestion, err := GetIssueTitleSuggestionFromCustomerRequest(
		context.Background(),
		graphqlClient,
		"Customer asks for faster exports",
	)
	require.NoError(t, err)
	issue, err := GetIssueByID(context.Background(), graphqlClient, "LIT-11")
	require.NoError(t, err)
	issueAttachments, err := ListIssueAttachments(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueBotActor, err := GetIssueBotActor(context.Background(), graphqlClient, "LIT-1")
	require.NoError(t, err)
	issueChildren, err := ListIssueChildren(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueDocuments, err := ListIssueDocuments(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueFormerAttachments, err := ListIssueFormerAttachments(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueHistory, err := ListIssueHistory(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueInverseRelations, err := ListIssueInverseRelations(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueLabels, err := ListIssueLabels(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueScopedRelations, err := ListIssueRelationsForIssue(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueReleases, err := ListIssueReleases(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueStateHistory, err := ListIssueStateHistory(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	issueSubscribers, err := ListIssueSubscribers(context.Background(), graphqlClient, "LIT-1", 2)
	require.NoError(t, err)
	branchIssue, err := GetIssueByVCSBranch(context.Background(), graphqlClient, "omer/branch")
	require.NoError(t, err)
	branchIssueAttachments, err := ListIssueVCSBranchAttachments(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueBotActor, err := GetIssueVCSBranchBotActor(context.Background(), graphqlClient, "omer/branch")
	require.NoError(t, err)
	branchIssueChildren, err := ListIssueVCSBranchChildren(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueDocuments, err := ListIssueVCSBranchDocuments(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueFormerAttachments, err := ListIssueVCSBranchFormerAttachments(
		context.Background(),
		graphqlClient,
		"omer/branch",
		2,
	)
	require.NoError(t, err)
	branchIssueHistory, err := ListIssueVCSBranchHistory(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueInverseRelations, err := ListIssueVCSBranchInverseRelations(
		context.Background(),
		graphqlClient,
		"omer/branch",
		2,
	)
	require.NoError(t, err)
	branchIssueLabels, err := ListIssueVCSBranchLabels(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueRelations, err := ListIssueVCSBranchRelations(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueReleases, err := ListIssueVCSBranchReleases(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueStateHistory, err := ListIssueVCSBranchStateHistory(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	branchIssueSubscribers, err := ListIssueVCSBranchSubscribers(context.Background(), graphqlClient, "omer/branch", 2)
	require.NoError(t, err)
	projects, err := ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	project, err := GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.NoError(t, err)
	members, err := ListProjectMembers(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectAttachments, err := ListProjectAttachments(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectDocuments, err := ListProjectDocuments(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectExternalLinks, err := ListProjectExternalLinks(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectHistory, err := ListProjectHistory(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectInitiativeAssociations, err := ListProjectInitiativeToProjects(
		context.Background(),
		graphqlClient,
		"project-id",
		2,
	)
	require.NoError(t, err)
	projectInitiatives, err := ListProjectInitiatives(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectInverseRelations, err := ListProjectInverseRelations(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	scopedProjectIssues, err := ListProjectIssues(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectComments, err := ListProjectComments(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	labelsForProject, err := ListLabelsForProject(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectNeeds, err := ListProjectNeeds(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectScopedRelations, err := ListProjectRelationsForProject(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectTeams, err := ListProjectTeams(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectUpdates, err := ListProjectUpdates(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectFilterSuggestion, err := GetProjectFilterSuggestion(context.Background(), graphqlClient, "started projects", "team-id")
	require.NoError(t, err)
	allProjectUpdates, err := ListAllProjectUpdates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectUpdate, err := GetProjectUpdateByID(context.Background(), graphqlClient, "project-update-id")
	require.NoError(t, err)
	projectUpdateComments, err := ListProjectUpdateComments(context.Background(), graphqlClient, "project-update-id", 2)
	require.NoError(t, err)
	projectMilestones, err := ListProjectMilestones(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	allProjectMilestones, err := ListAllProjectMilestones(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectMilestone, err := GetProjectMilestoneByID(context.Background(), graphqlClient, "project-milestone-id")
	require.NoError(t, err)
	projectMilestoneIssues, err := ListProjectMilestoneIssues(context.Background(), graphqlClient, "project-milestone-id", 2)
	require.NoError(t, err)
	workspaceProjects, err := ListProjects(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectStatuses, err := ListProjectStatuses(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectStatus, err := GetProjectStatusByID(context.Background(), graphqlClient, "project-status-id")
	require.NoError(t, err)
	projectStatusProjectCount, err := GetProjectStatusProjectCount(context.Background(), graphqlClient, "project-status-id")
	require.NoError(t, err)
	projectLabels, err := ListProjectLabels(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectLabel, err := GetProjectLabelByID(context.Background(), graphqlClient, "project-label-id")
	require.NoError(t, err)
	projectLabelChildren, err := ListProjectLabelChildren(context.Background(), graphqlClient, "project-label-id", 2)
	require.NoError(t, err)
	projectLabelProjects, err := ListProjectLabelProjects(context.Background(), graphqlClient, "project-label-id", 2)
	require.NoError(t, err)
	projectRelations, err := ListProjectRelations(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectRelation, err := GetProjectRelationByID(context.Background(), graphqlClient, "project-relation-id")
	require.NoError(t, err)
	issueRelations, err := ListIssueRelations(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	issueRelation, err := GetIssueRelationByID(context.Background(), graphqlClient, "issue-relation-id")
	require.NoError(t, err)
	issueToReleases, err := ListIssueToReleases(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	issueToRelease, err := GetIssueToReleaseByID(context.Background(), graphqlClient, "issue-to-release-id")
	require.NoError(t, err)
	application, err := GetApplicationInfo(context.Background(), graphqlClient, "app-client-id")
	require.NoError(t, err)
	agentActivities, err := ListAgentActivities(context.Background(), graphqlClient, 6)
	require.NoError(t, err)
	agentActivity, err := GetAgentActivityByID(context.Background(), graphqlClient, "agent-activity-id")
	require.NoError(t, err)
	agentSkills, err := ListAgentSkills(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	agentSkill, err := GetAgentSkillByID(context.Background(), graphqlClient, "agent-skill-id")
	require.NoError(t, err)
	externalUsers, err := ListExternalUsers(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	externalUser, err := GetExternalUserByID(context.Background(), graphqlClient, "external-user-id")
	require.NoError(t, err)
	auditEntryTypes, err := ListAuditEntryTypes(context.Background(), graphqlClient)
	require.NoError(t, err)
	comments, err := ListIssueComments(context.Background(), graphqlClient, "LIT-12", 2)
	require.NoError(t, err)
	topLevelComments, err := ListComments(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	topLevelComment, err := GetCommentByID(context.Background(), graphqlClient, "comment-id")
	require.NoError(t, err)
	commentBotActor, err := GetCommentBotActor(context.Background(), graphqlClient, "comment-id")
	require.NoError(t, err)
	commentChildren, err := ListCommentChildren(context.Background(), graphqlClient, "comment-id", 1)
	require.NoError(t, err)
	commentCreatedIssues, err := ListCommentCreatedIssues(context.Background(), graphqlClient, "comment-id", 1)
	require.NoError(t, err)
	documents, err := ListDocuments(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	document, err := GetDocumentByID(context.Background(), graphqlClient, "document-id")
	require.NoError(t, err)
	documentComments, err := ListDocumentComments(context.Background(), graphqlClient, "document-id", 2)
	require.NoError(t, err)
	labels, err := ListLabels(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	label, err := GetLabelByID(context.Background(), graphqlClient, "label-id")
	require.NoError(t, err)
	labelChildren, err := ListLabelChildren(context.Background(), graphqlClient, "label-id", 2)
	require.NoError(t, err)
	labelScopedIssues, err := ListLabelIssues(context.Background(), graphqlClient, "label-id", 2)
	require.NoError(t, err)
	teams, err := ListTeams(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	teamMemberships, err := ListTeamMemberships(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	teamMembership, err := GetTeamMembershipByID(context.Background(), graphqlClient, "team-membership-id")
	require.NoError(t, err)
	organizationExists, err := CheckOrganizationExists(context.Background(), graphqlClient, "kyanite")
	require.NoError(t, err)
	rateLimitStatus, err := GetRateLimitStatus(context.Background(), graphqlClient)
	require.NoError(t, err)
	notifications, err := ListNotifications(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	notification, err := GetNotificationByID(context.Background(), graphqlClient, "notification-id")
	require.NoError(t, err)
	notificationSubscriptions, err := ListNotificationSubscriptions(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	notificationSubscription, err := GetNotificationSubscriptionByID(
		context.Background(),
		graphqlClient,
		"notification-subscription-id",
	)
	require.NoError(t, err)
	triageResponsibilities, err := ListTriageResponsibilities(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	triageResponsibility, err := GetTriageResponsibilityByID(
		context.Background(),
		graphqlClient,
		"triage-responsibility-id",
	)
	require.NoError(t, err)
	triageManualSelection, err := GetTriageResponsibilityManualSelection(
		context.Background(),
		graphqlClient,
		"triage-responsibility-id",
	)
	require.NoError(t, err)
	slaConfigurations, err := ListSLAConfigurations(context.Background(), graphqlClient, "team-id")
	require.NoError(t, err)
	semanticSearch, err := SearchSemantic(context.Background(), graphqlClient, "agent search", 2)
	require.NoError(t, err)
	documentSearch, err := SearchDocuments(context.Background(), graphqlClient, "agent search", 2)
	require.NoError(t, err)
	issueSearch, err := SearchIssues(context.Background(), graphqlClient, "agent search", 2)
	require.NoError(t, err)
	projectSearch, err := SearchProjects(context.Background(), graphqlClient, "agent search", 2)
	require.NoError(t, err)
	releasePipelines, err := ListReleasePipelines(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	releasePipeline, err := GetReleasePipelineByID(context.Background(), graphqlClient, "release-pipeline-id")
	require.NoError(t, err)
	releasePipelineReleases, err := ListReleasePipelineReleases(context.Background(), graphqlClient, "release-pipeline-id", 2)
	require.NoError(t, err)
	releasePipelineStages, err := ListReleasePipelineStages(context.Background(), graphqlClient, "release-pipeline-id", 2)
	require.NoError(t, err)
	releasePipelineTeams, err := ListReleasePipelineTeams(context.Background(), graphqlClient, "release-pipeline-id", 2)
	require.NoError(t, err)
	releaseStages, err := ListReleaseStages(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	releaseStage, err := GetReleaseStageByID(context.Background(), graphqlClient, "release-stage-id")
	require.NoError(t, err)
	releaseStageReleases, err := ListReleaseStageReleases(context.Background(), graphqlClient, "release-stage-id", 2)
	require.NoError(t, err)
	releases, err := ListReleases(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	release, err := GetReleaseByID(context.Background(), graphqlClient, "release-id")
	require.NoError(t, err)
	releaseDocuments, err := ListReleaseDocuments(context.Background(), graphqlClient, "release-id", 2)
	require.NoError(t, err)
	releaseIssues, err := ListReleaseIssues(context.Background(), graphqlClient, "release-id", 2)
	require.NoError(t, err)
	releaseHistory, err := ListReleaseHistory(context.Background(), graphqlClient, "release-id", 2)
	require.NoError(t, err)
	releaseLinks, err := ListReleaseLinks(context.Background(), graphqlClient, "release-id", 2)
	require.NoError(t, err)
	externalLink, err := GetEntityExternalLinkByID(context.Background(), graphqlClient, "release-link-id")
	require.NoError(t, err)
	releaseSearch, err := SearchReleases(context.Background(), graphqlClient, "mobile", 2)
	require.NoError(t, err)
	releaseNotes, err := ListReleaseNotes(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	releaseNote, err := GetReleaseNoteByID(context.Background(), graphqlClient, "release-note-id")
	require.NoError(t, err)
	team, err := GetTeamByID(context.Background(), graphqlClient, "team-id")
	require.NoError(t, err)
	teamMembers, err := ListTeamMembers(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamCycles, err := ListTeamCycles(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamIssues, err := ListTeamIssues(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamLabels, err := ListTeamLabels(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamScopedMemberships, err := ListTeamMembershipsForTeam(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamProjects, err := ListTeamProjects(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamReleasePipelines, err := ListTeamReleasePipelines(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamWorkflowStates, err := ListTeamWorkflowStates(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamGitAutomationStates, err := ListTeamGitAutomationStates(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	teamTemplates, err := ListTeamTemplates(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	users, err := ListUsers(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	user, err := GetUserByID(context.Background(), graphqlClient, "user-id")
	require.NoError(t, err)
	viewerUser, err := GetViewerUser(context.Background(), graphqlClient)
	require.NoError(t, err)
	drafts, err := ListViewerDrafts(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	userAssignedIssues, err := ListUserAssignedIssues(context.Background(), graphqlClient, "user-id", 2)
	require.NoError(t, err)
	userCreatedIssues, err := ListUserCreatedIssues(context.Background(), graphqlClient, "user-id", 2)
	require.NoError(t, err)
	userDelegatedIssues, err := ListUserDelegatedIssues(context.Background(), graphqlClient, "user-id", 2)
	require.NoError(t, err)
	userTeamMemberships, err := ListUserTeamMemberships(context.Background(), graphqlClient, "user-id", 2)
	require.NoError(t, err)
	userTeams, err := ListUserTeams(context.Background(), graphqlClient, "user-id", 2)
	require.NoError(t, err)
	viewerAssignedIssues, err := ListViewerAssignedIssues(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	viewerCreatedIssues, err := ListViewerCreatedIssues(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	viewerDelegatedIssues, err := ListViewerDelegatedIssues(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	viewerTeamMemberships, err := ListViewerTeamMemberships(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	viewerTeams, err := ListViewerTeams(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	workflowStates, err := ListWorkflowStates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	workflowState, err := GetWorkflowStateByID(context.Background(), graphqlClient, "workflow-state-id")
	require.NoError(t, err)
	workflowStateIssues, err := ListWorkflowStateIssues(context.Background(), graphqlClient, "workflow-state-id", 2)
	require.NoError(t, err)
	timeSchedules, err := ListTimeSchedules(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	timeSchedule, err := GetTimeScheduleByID(context.Background(), graphqlClient, "time-schedule-id")
	require.NoError(t, err)
	organizationLabels, err := ListOrganizationLabels(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	organizationProjectLabels, err := ListOrganizationProjectLabels(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	organizationTeams, err := ListOrganizationTeams(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	organizationTemplates, err := ListOrganizationTemplates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	organizationUsers, err := ListOrganizationUsers(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	templates, err := ListTemplates(context.Background(), graphqlClient, 1)
	require.NoError(t, err)
	template, err := GetTemplateByID(context.Background(), graphqlClient, "template-id")
	require.NoError(t, err)
	initiatives, err := ListInitiatives(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	initiative, err := GetInitiativeByID(context.Background(), graphqlClient, "initiative-id")
	require.NoError(t, err)
	initiativeHistory, err := ListInitiativeHistory(context.Background(), graphqlClient, "initiative-id", 2)
	require.NoError(t, err)
	initiativeLinks, err := ListInitiativeLinks(context.Background(), graphqlClient, "initiative-id", 2)
	require.NoError(t, err)
	subInitiatives, err := ListSubInitiatives(context.Background(), graphqlClient, "initiative-id", 2)
	require.NoError(t, err)
	initiativeScopedUpdates, err := ListInitiativeUpdatesForInitiative(
		context.Background(),
		graphqlClient,
		"initiative-id",
		2,
	)
	require.NoError(t, err)
	initiativeDocuments, err := ListInitiativeDocuments(context.Background(), graphqlClient, "initiative-id", 2)
	require.NoError(t, err)
	initiativeProjects, err := ListInitiativeProjects(context.Background(), graphqlClient, "initiative-id", 2)
	require.NoError(t, err)
	initiativeRelations, err := ListInitiativeRelations(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	initiativeRelation, err := GetInitiativeRelationByID(context.Background(), graphqlClient, "initiative-relation-id")
	require.NoError(t, err)
	initiativeToProjects, err := ListInitiativeToProjects(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	initiativeToProject, err := GetInitiativeToProjectByID(context.Background(), graphqlClient, "initiative-to-project-id")
	require.NoError(t, err)
	roadmapToProjects, err := ListRoadmapToProjects(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	roadmapToProject, err := GetRoadmapToProjectByID(context.Background(), graphqlClient, "roadmap-to-project-id")
	require.NoError(t, err)
	initiativeUpdates, err := ListInitiativeUpdates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	initiativeUpdate, err := GetInitiativeUpdateByID(context.Background(), graphqlClient, "initiative-update-id")
	require.NoError(t, err)
	initiativeUpdateComments, err := ListInitiativeUpdateComments(
		context.Background(),
		graphqlClient,
		"initiative-update-id",
		2,
	)
	require.NoError(t, err)
	roadmaps, err := ListRoadmaps(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	roadmap, err := GetRoadmapByID(context.Background(), graphqlClient, "roadmap-id")
	require.NoError(t, err)
	roadmapProjects, err := ListRoadmapProjects(context.Background(), graphqlClient, "roadmap-id", 2)
	require.NoError(t, err)
	customViews, err := ListCustomViews(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	customViewSubscribers, err := GetCustomViewSubscriberStatus(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	customView, err := GetCustomViewByID(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	customViewInitiatives, err := ListCustomViewInitiatives(context.Background(), graphqlClient, "custom-view-id", 2)
	require.NoError(t, err)
	customViewIssues, err := ListCustomViewIssues(context.Background(), graphqlClient, "custom-view-id", 2)
	require.NoError(t, err)
	customViewOrganizationPreferences, err := GetCustomViewOrganizationPreferences(
		context.Background(),
		graphqlClient,
		"custom-view-id",
	)
	require.NoError(t, err)
	customViewOrganizationPreferenceValues, err := GetCustomViewOrganizationPreferenceValues(
		context.Background(),
		graphqlClient,
		"custom-view-id",
	)
	require.NoError(t, err)
	customViewProjects, err := ListCustomViewProjects(context.Background(), graphqlClient, "custom-view-id", 2)
	require.NoError(t, err)
	customViewUserPreferences, err := GetCustomViewUserPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	customViewUserPreferenceValues, err := GetCustomViewUserPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	customViewPreferenceValues, err := GetCustomViewPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	customers, err := ListCustomers(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	customer, err := GetCustomerByID(context.Background(), graphqlClient, "customer-id")
	require.NoError(t, err)
	customerNeeds, err := ListCustomerNeeds(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	customerNeed, err := GetCustomerNeedByID(context.Background(), graphqlClient, "customer-need-id")
	require.NoError(t, err)
	customerNeedProjectAttachment, err := GetCustomerNeedProjectAttachment(
		context.Background(),
		graphqlClient,
		"customer-need-id",
	)
	require.NoError(t, err)
	customerStatuses, err := ListCustomerStatuses(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	customerStatus, err := GetCustomerStatusByID(context.Background(), graphqlClient, "customer-status-id")
	require.NoError(t, err)
	customerTiers, err := ListCustomerTiers(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	customerTier, err := GetCustomerTierByID(context.Background(), graphqlClient, "customer-tier-id")
	require.NoError(t, err)
	favorites, err := ListFavorites(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	favoriteChildren, err := ListFavoriteChildren(context.Background(), graphqlClient, "favorite-folder-id", 2)
	require.NoError(t, err)
	favorite, err := GetFavoriteByID(context.Background(), graphqlClient, "favorite-id")
	require.NoError(t, err)
	emojis, err := ListEmojis(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	emoji, err := GetEmojiByID(context.Background(), graphqlClient, "emoji-id")
	require.NoError(t, err)
	attachments, err := ListAttachments(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	attachmentsForURL, err := ListAttachmentsForURL(context.Background(), graphqlClient, "https://example.com/spec", 2)
	require.NoError(t, err)
	attachment, err := GetAttachmentByID(context.Background(), graphqlClient, "attachment-id")
	require.NoError(t, err)
	attachmentIssue, err := GetAttachmentIssue(context.Background(), graphqlClient, "attachment-id")
	require.NoError(t, err)
	attachmentIssueAttachments, err := ListAttachmentIssueAttachments(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueBotActor, err := GetAttachmentIssueBotActor(context.Background(), graphqlClient, "attachment-id")
	require.NoError(t, err)
	attachmentIssueChildren, err := ListAttachmentIssueChildren(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueDocuments, err := ListAttachmentIssueDocuments(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueFormerAttachments, err := ListAttachmentIssueFormerAttachments(
		context.Background(),
		graphqlClient,
		"attachment-id",
		2,
	)
	require.NoError(t, err)
	attachmentIssueHistory, err := ListAttachmentIssueHistory(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueInverseRelations, err := ListAttachmentIssueInverseRelations(
		context.Background(),
		graphqlClient,
		"attachment-id",
		2,
	)
	require.NoError(t, err)
	attachmentIssueLabels, err := ListAttachmentIssueLabels(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueRelations, err := ListAttachmentIssueRelations(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueReleases, err := ListAttachmentIssueReleases(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueStateHistory, err := ListAttachmentIssueStateHistory(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)
	attachmentIssueSubscribers, err := ListAttachmentIssueSubscribers(context.Background(), graphqlClient, "attachment-id", 2)
	require.NoError(t, err)

	// Then
	require.Equal(t, "LIT-20", allTeamIssues.Issues[0].Identifier)
	require.True(t, issues.HasNextPage)
	require.Equal(t, "LIT-10", issues.Issues[0].Identifier)
	require.False(t, startedIssues.HasNextPage)
	require.Equal(t, "started", startedIssues.Issues[0].StateType)
	require.Equal(t, "project-id", projectIssues.Issues[0].ProjectID)
	require.Equal(t, "Omer", mineIssues.Issues[0].Assignee)
	require.Equal(t, "LIT-16", labelIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-17", cycleIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-18", recentIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-19", olderIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-21", blockedIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-22", blockingIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-27", nextIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-23", blockedByIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-1", dependencies.Identifier)
	require.Equal(t, "LIT-25", dependencies.Parent.Identifier)
	require.Equal(t, "LIT-26", dependencies.Children[0].Identifier)
	require.Equal(t, "LIT-23", dependencies.Blocks[0].Identifier)
	require.Equal(t, "LIT-24", dependencies.BlockedBy[0].Identifier)
	require.True(t, dependencies.HasNextPage)
	require.Equal(t, "LIT-13", searchIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-41", figmaIssues.Issues[0].Identifier)
	require.True(t, figmaIssues.HasNextPage)
	require.Equal(t, &endCursor, figmaIssues.EndCursor)
	require.Equal(t, 1, issuePriorityValues[0].Priority)
	require.Equal(t, "Urgent", issuePriorityValues[0].Label)
	require.JSONEq(t, `{"state":{"type":{"eq":"started"}}}`, string(issueFilterSuggestion.Filter))
	require.Equal(t, "issue-filter-log-id", issueFilterSuggestion.LogID)
	require.Equal(t, "Improve exports", issueTitleSuggestion.Title)
	require.Equal(t, "title-log-id", issueTitleSuggestion.LogID)
	require.Equal(t, "LIT-11", issue.Identifier)
	require.True(t, issueAttachments.HasNextPage)
	require.Equal(t, &endCursor, issueAttachments.EndCursor)
	require.Equal(t, "project-attachment-id", issueAttachments.Attachments[0].ID)
	require.Equal(t, "issue-id", issueBotActor.IssueID)
	require.Equal(t, "bot-actor-id", issueBotActor.Bot.ID)
	require.Equal(t, "github", issueBotActor.Bot.Type)
	require.Equal(t, "LIT-31", issueChildren.Issues[0].Identifier)
	require.Equal(t, "Issue spec", issueDocuments.Documents[0].Title)
	require.Equal(t, "issue", issueDocuments.Documents[0].ParentType)
	require.Equal(t, "project-attachment-id", issueFormerAttachments.Attachments[0].ID)
	require.Equal(t, "issue-history-id", issueHistory.History[0].ID)
	require.Equal(t, "issue-id", issueHistory.History[0].IssueID)
	require.Equal(t, "user-id", issueHistory.History[0].ActorID)
	require.True(t, issueHistory.History[0].UpdatedDescription)
	require.Equal(t, "issue-relation-id", issueInverseRelations.Relations[0].ID)
	require.Equal(t, "Bug", issueLabels.Labels[0].Name)
	require.Equal(t, "issue-relation-id", issueScopedRelations.Relations[0].ID)
	require.Equal(t, "release-id", issueReleases.Releases[0].ID)
	require.True(t, issueStateHistory.HasNextPage)
	require.Equal(t, "issue-id", issueStateHistory.IssueID)
	require.Equal(t, "issue-state-span-id", issueStateHistory.Spans[0].ID)
	require.Equal(t, "Started", issueStateHistory.Spans[0].StateName)
	require.Equal(t, "started", issueStateHistory.Spans[0].StateType)
	require.Equal(t, "Omer", issueSubscribers.Users[0].DisplayName)
	require.True(t, issueSubscribers.HasNextPage)
	require.Equal(t, "LIT-40", branchIssue.Identifier)
	require.Equal(t, "project-attachment-id", branchIssueAttachments.Attachments[0].ID)
	require.Equal(t, "issue-id", branchIssueBotActor.IssueID)
	require.Equal(t, "bot-actor-id", branchIssueBotActor.Bot.ID)
	require.Equal(t, "LIT-43", branchIssueChildren.Issues[0].Identifier)
	require.Equal(t, "Branch issue spec", branchIssueDocuments.Documents[0].Title)
	require.Equal(t, "project-attachment-id", branchIssueFormerAttachments.Attachments[0].ID)
	require.Equal(t, "issue-history-id", branchIssueHistory.History[0].ID)
	require.Equal(t, "issue-relation-id", branchIssueInverseRelations.Relations[0].ID)
	require.Equal(t, "Bug", branchIssueLabels.Labels[0].Name)
	require.Equal(t, "issue-relation-id", branchIssueRelations.Relations[0].ID)
	require.Equal(t, "release-id", branchIssueReleases.Releases[0].ID)
	require.Equal(t, "issue-id", branchIssueStateHistory.IssueID)
	require.Equal(t, "issue-state-span-id", branchIssueStateHistory.Spans[0].ID)
	require.Equal(t, "Omer", branchIssueSubscribers.Users[0].DisplayName)
	require.True(t, projects.HasNextPage)
	require.Equal(t, "listed", projects.Projects[0].Name)
	require.Equal(t, "detail", project.Name)
	require.Equal(t, "Omer", members.Members[0].DisplayName)
	require.Equal(t, &endCursor, members.EndCursor)
	require.Equal(t, "project-attachment-id", projectAttachments.Attachments[0].ID)
	require.Equal(t, "github", projectAttachments.Attachments[0].SourceType)
	require.Equal(t, "Project spec", projectDocuments.Documents[0].Title)
	require.Equal(t, "project", projectDocuments.Documents[0].ParentType)
	require.Equal(t, "release-link-id", projectExternalLinks.Links[0].ID)
	require.Equal(t, "Runbook", projectExternalLinks.Links[0].Label)
	require.Equal(t, "project-history-id", projectHistory.History[0].ID)
	require.Equal(t, 1, projectHistory.History[0].EntryCount)
	require.JSONEq(
		t,
		`[{"from":"Backlog","to":"Started","type":"status"}]`,
		string(projectHistory.History[0].Entries),
	)
	require.Equal(t, "initiative-to-project-id", projectInitiativeAssociations.Associations[0].ID)
	require.Equal(t, "Platform", projectInitiativeAssociations.Associations[0].InitiativeName)
	require.Equal(t, "initiative-id", projectInitiatives.Initiatives[0].ID)
	require.Equal(t, "Active", projectInitiatives.Initiatives[0].Status)
	require.Equal(t, "project-relation-id", projectInverseRelations.Relations[0].ID)
	require.Equal(t, "Related project", projectInverseRelations.Relations[0].RelatedProjectName)
	require.Equal(t, "LIT-47", scopedProjectIssues.Issues[0].Identifier)
	require.True(t, projectComments.HasNextPage)
	require.Equal(t, "project-id", projectComments.Comments[0].ProjectID)
	require.Equal(t, "Omer", projectComments.Comments[0].DisplayName)
	require.Empty(t, projectComments.Comments[1].UserID)
	require.Equal(t, &endCursor, projectComments.EndCursor)
	require.Equal(t, "project-label-id", labelsForProject.ProjectLabels[0].ID)
	require.Equal(t, "customer-need-id", projectNeeds.Needs[0].ID)
	require.Equal(t, "Acme", projectNeeds.Needs[0].CustomerName)
	require.Equal(t, "project-relation-id", projectScopedRelations.Relations[0].ID)
	require.Equal(t, "linctl", projectTeams.Teams[0].Name)
	require.True(t, projectUpdates.HasNextPage)
	require.Equal(t, "project-update-id", projectUpdates.Updates[0].ID)
	require.Equal(t, "onTrack", projectUpdates.Updates[0].Health)
	require.Equal(t, "Omer", projectUpdates.Updates[0].DisplayName)
	require.Empty(t, projectUpdates.Updates[0].Body)
	require.Equal(t, &endCursor, projectUpdates.EndCursor)
	require.JSONEq(t, `{"status":{"type":{"eq":"started"}}}`, string(projectFilterSuggestion.Filter))
	require.Equal(t, "filter-log-id", projectFilterSuggestion.LogID)
	require.True(t, allProjectUpdates.HasNextPage)
	require.Equal(t, &endCursor, allProjectUpdates.EndCursor)
	require.Equal(t, "project-id", allProjectUpdates.Updates[0].ProjectID)
	require.Equal(t, "detail", allProjectUpdates.Updates[0].ProjectName)
	require.Equal(t, "project-update-id", projectUpdate.ID)
	require.Equal(t, "detail", projectUpdate.ProjectName)
	require.True(t, projectUpdateComments.HasNextPage)
	require.Equal(t, "project-update-id", projectUpdateComments.ProjectUpdateID)
	require.Equal(t, "project-update-id", projectUpdateComments.Comments[0].ProjectUpdateID)
	require.True(t, projectMilestones.HasNextPage)
	require.Equal(t, "project-milestone-id", projectMilestones.Milestones[0].ID)
	require.Equal(t, "Launch milestone", projectMilestones.Milestones[0].Name)
	require.Equal(t, "milestone body", projectMilestones.Milestones[0].Description)
	require.Equal(t, "2026-06-30", projectMilestones.Milestones[0].TargetDate)
	require.Equal(t, "next", projectMilestones.Milestones[0].Status)
	require.Equal(t, &endCursor, projectMilestones.EndCursor)
	require.True(t, allProjectMilestones.HasNextPage)
	require.Equal(t, "project-milestone-id", allProjectMilestones.Milestones[0].ID)
	require.Equal(t, "project-milestone-id", projectMilestone.ID)
	require.Equal(t, "Launch milestone", projectMilestone.Name)
	require.Equal(t, "next", projectMilestone.Status)
	require.True(t, projectMilestoneIssues.HasNextPage)
	require.Equal(t, "project-milestone-id", projectMilestoneIssues.ProjectMilestoneID)
	require.Equal(t, "LIT-52", projectMilestoneIssues.Issues[0].Identifier)
	require.Equal(t, "workspace-project-id", workspaceProjects.Projects[0].ID)
	require.Equal(t, &endCursor, workspaceProjects.EndCursor)
	require.True(t, projectStatuses.HasNextPage)
	require.Equal(t, &endCursor, projectStatuses.EndCursor)
	require.Equal(t, "project-status-id", projectStatuses.ProjectStatuses[0].ID)
	require.Equal(t, "Backlog", projectStatuses.ProjectStatuses[0].Name)
	require.Equal(t, "backlog", projectStatuses.ProjectStatuses[0].Type)
	require.Equal(t, "#bec2c8", projectStatuses.ProjectStatuses[0].Color)
	require.Equal(t, "project-status-id", projectStatus.ID)
	require.Equal(t, "Ready for planning", projectStatus.Description)
	require.Equal(t, "project-status-id", projectStatusProjectCount.ProjectStatusID)
	require.InEpsilon(t, float64(12), projectStatusProjectCount.Count, 0.001)
	require.InEpsilon(t, float64(2), projectStatusProjectCount.PrivateCount, 0.001)
	require.InEpsilon(t, float64(1), projectStatusProjectCount.ArchivedTeamCount, 0.001)
	require.True(t, projectLabels.HasNextPage)
	require.Equal(t, &endCursor, projectLabels.EndCursor)
	require.Equal(t, "project-label-id", projectLabels.ProjectLabels[0].ID)
	require.Equal(t, "Roadmap", projectLabels.ProjectLabels[0].Name)
	require.Equal(t, "#f2c94c", projectLabels.ProjectLabels[0].Color)
	require.Equal(t, "project-label-id", projectLabel.ID)
	require.Equal(t, "Parent", projectLabel.ParentName)
	require.Equal(t, "project-label-id", projectLabelChildren.ProjectLabelID)
	require.Equal(t, "child-project-label-id", projectLabelChildren.ProjectLabels[0].ID)
	require.Equal(t, "project-id", projectLabelProjects.Projects[0].ID)
	require.True(t, projectRelations.HasNextPage)
	require.Equal(t, &endCursor, projectRelations.EndCursor)
	require.Equal(t, "project-relation-id", projectRelations.Relations[0].ID)
	require.Equal(t, "blocks", projectRelations.Relations[0].Type)
	require.Equal(t, "Pinned project", projectRelations.Relations[0].ProjectName)
	require.Equal(t, "Related project", projectRelations.Relations[0].RelatedProjectName)
	require.Equal(t, "Omer", projectRelations.Relations[0].DisplayName)
	require.Equal(t, "project-milestone-id", projectRelations.Relations[1].ProjectMilestoneID)
	require.Equal(t, "related-project-milestone-id", projectRelations.Relations[1].RelatedProjectMilestoneID)
	require.Empty(t, projectRelations.Relations[1].DisplayName)
	require.Equal(t, "project-relation-id", projectRelation.ID)
	require.Equal(t, "blocks", projectRelation.Type)
	require.True(t, issueRelations.HasNextPage)
	require.Equal(t, &endCursor, issueRelations.EndCursor)
	require.Equal(t, "issue-relation-id", issueRelations.Relations[0].ID)
	require.Equal(t, "blocks", issueRelations.Relations[0].Type)
	require.Equal(t, "LIT-1", issueRelations.Relations[0].IssueIdentifier)
	require.Equal(t, "LIT-2", issueRelations.Relations[0].RelatedIssueIdentifier)
	require.Equal(t, "issue-relation-id", issueRelation.ID)
	require.Equal(t, "related-issue-id", issueRelation.RelatedIssueID)
	require.True(t, issueToReleases.HasNextPage)
	require.Equal(t, &endCursor, issueToReleases.EndCursor)
	require.Equal(t, "issue-to-release-id", issueToReleases.Associations[0].ID)
	require.Equal(t, "issue-id", issueToReleases.Associations[0].IssueID)
	require.Equal(t, "release-id", issueToReleases.Associations[0].ReleaseID)
	require.Equal(t, "issue-to-release-id", issueToRelease.ID)
	require.Equal(t, "release-id", issueToRelease.ReleaseID)
	require.Equal(t, "app-id", application.ID)
	require.Equal(t, "app-client-id", application.ClientID)
	require.Equal(t, "Demo App", application.Name)
	require.Equal(t, "Kyanite", application.Developer)
	require.True(t, agentActivities.HasNextPage)
	require.Equal(t, &endCursor, agentActivities.EndCursor)
	require.Equal(t, "action", agentActivities.AgentActivities[0].ContentType)
	require.Equal(t, "elicitation", agentActivities.AgentActivities[1].ContentType)
	require.Equal(t, "error", agentActivities.AgentActivities[2].ContentType)
	require.Equal(t, "prompt", agentActivities.AgentActivities[3].ContentType)
	require.Equal(t, "response", agentActivities.AgentActivities[4].ContentType)
	require.Equal(t, "thought", agentActivities.AgentActivities[5].ContentType)
	require.Equal(t, "agent-activity-id", agentActivity.ID)
	require.Equal(t, "agent-session-id", agentActivity.AgentSessionID)
	require.Equal(t, "comment-id", agentActivity.SourceCommentID)
	require.Equal(t, "read_file", agentActivity.Content.Action)
	require.Empty(t, agentActivityContentSummary(nil).Type)
	require.True(t, agentSkills.HasNextPage)
	require.Equal(t, &endCursor, agentSkills.EndCursor)
	require.Equal(t, "agent-skill-id", agentSkills.AgentSkills[0].ID)
	require.Equal(t, "Triage Helper", agentSkills.AgentSkills[0].Title)
	require.Equal(t, "agent-skill-id", agentSkill.ID)
	require.Equal(t, "updater-id", agentSkill.LastUpdatedByID)
	require.True(t, externalUsers.HasNextPage)
	require.Equal(t, &endCursor, externalUsers.EndCursor)
	require.Equal(t, "external-user-id", externalUsers.ExternalUsers[0].ID)
	require.Equal(t, "@external", externalUsers.ExternalUsers[0].DisplayName)
	require.Equal(t, "https://example.com/avatar.png", externalUsers.ExternalUsers[0].AvatarURL)
	require.Equal(t, "external-user-id", externalUser.ID)
	require.Equal(t, "2026-06-19T12:00:00Z", externalUser.LastSeen)
	require.Equal(t, "user_login", auditEntryTypes.AuditEntryTypes[0].Type)
	require.Equal(t, "User logged in", auditEntryTypes.AuditEntryTypes[0].Description)
	require.Equal(t, "LIT-12", comments.Identifier)
	require.True(t, comments.HasNextPage)
	require.Equal(t, &endCursor, comments.EndCursor)
	require.Equal(t, "parent-id", comments.Comments[0].ParentID)
	require.Equal(t, "Omer", comments.Comments[0].DisplayName)
	require.Empty(t, comments.Comments[1].UserID)
	require.Empty(t, comments.Comments[1].ParentID)
	require.True(t, topLevelComments.HasNextPage)
	require.Equal(t, &endCursor, topLevelComments.EndCursor)
	require.Equal(t, "parent-id", topLevelComments.Comments[0].ParentID)
	require.Equal(t, "issue-id", topLevelComments.Comments[0].IssueID)
	require.Equal(t, "Omer", topLevelComments.Comments[0].DisplayName)
	require.Equal(t, "project-id", topLevelComments.Comments[1].ProjectID)
	require.Empty(t, topLevelComments.Comments[1].UserID)
	require.Equal(t, "comment-id", topLevelComment.ID)
	require.Equal(t, "issue-id", topLevelComment.IssueID)
	require.Equal(t, "Omer", topLevelComment.DisplayName)
	require.Equal(t, "comment-id", commentBotActor.CommentID)
	require.Equal(t, "bot-actor-id", commentBotActor.Bot.ID)
	require.Equal(t, "github", commentBotActor.Bot.Type)
	require.Equal(t, "GitHub", commentBotActor.Bot.Name)
	require.True(t, commentChildren.HasNextPage)
	require.Equal(t, "comment-id", commentChildren.CommentID)
	require.Equal(t, "child-comment-id", commentChildren.Comments[0].ID)
	require.Equal(t, "comment-id", commentChildren.Comments[0].ParentID)
	require.True(t, commentCreatedIssues.HasNextPage)
	require.Equal(t, "LIT-54", commentCreatedIssues.Issues[0].Identifier)
	require.True(t, documents.HasNextPage)
	require.Equal(t, "project", documents.Documents[0].ParentType)
	require.Equal(t, "Team note", document.Title)
	require.Equal(t, "team", document.ParentType)
	require.Equal(t, "document-id", documentComments.DocumentID)
	require.Equal(t, "comment-id", documentComments.Comments[0].ID)
	require.Equal(t, &endCursor, documentComments.EndCursor)
	require.True(t, labels.HasNextPage)
	require.Equal(t, "Bug", labels.Labels[0].Name)
	require.Equal(t, "LIT", labels.Labels[0].TeamKey)
	require.Equal(t, "label-id", label.ID)
	require.Empty(t, label.Description)
	require.Equal(t, "child-label-id", labelChildren.Labels[0].ID)
	require.Equal(t, "Mobile", labelChildren.Labels[0].Name)
	require.Equal(t, &endCursor, labelChildren.EndCursor)
	require.Equal(t, "LIT-53", labelScopedIssues.Issues[0].Identifier)
	require.Equal(t, "Label issue", labelScopedIssues.Issues[0].Title)
	require.Equal(t, &endCursor, labelScopedIssues.EndCursor)
	require.True(t, teams.HasNextPage)
	require.True(t, teamMemberships.HasNextPage)
	require.Equal(t, &endCursor, teamMemberships.EndCursor)
	require.Equal(t, "team-membership-id", teamMemberships.Memberships[0].ID)
	require.Equal(t, "team-id", teamMemberships.Memberships[0].TeamID)
	require.Equal(t, "LIT", teamMemberships.Memberships[0].TeamKey)
	require.Equal(t, "user-id", teamMemberships.Memberships[0].UserID)
	require.Equal(t, "Omer", teamMemberships.Memberships[0].DisplayName)
	require.True(t, teamMemberships.Memberships[0].Owner)
	require.InEpsilon(t, 1.5, teamMemberships.Memberships[0].SortOrder, 0)
	require.True(t, teamCycles.HasNextPage)
	require.Equal(t, &endCursor, teamCycles.EndCursor)
	require.Equal(t, "cycle-id", teamCycles.Cycles[0].ID)
	require.Equal(t, "LIT-1", teamIssues.Issues[0].Identifier)
	require.Equal(t, "label-id", teamLabels.Labels[0].ID)
	require.Equal(t, "team-membership-id", teamScopedMemberships.Memberships[0].ID)
	require.Equal(t, "Team project", teamProjects.Projects[0].Name)
	require.Equal(t, "release-pipeline-id", teamReleasePipelines.ReleasePipelines[0].ID)
	require.Equal(t, "workflow-state-id", teamWorkflowStates.WorkflowStates[0].ID)
	require.Equal(t, "git-automation-state-id", teamGitAutomationStates.States[0].ID)
	require.Equal(t, "review", teamGitAutomationStates.States[0].Event)
	require.Equal(t, "Started", teamGitAutomationStates.States[0].StateName)
	require.Equal(t, "main", teamGitAutomationStates.States[0].TargetBranchPattern)
	require.Equal(t, "template-id", teamTemplates.Templates[0].ID)
	require.Equal(t, "team-membership-id", teamMembership.ID)
	require.Equal(t, "omer@example.com", teamMembership.Email)
	require.Equal(t, "LIT", teams.Teams[0].Key)
	require.Equal(t, "kyanite", organizationExists.URLKey)
	require.True(t, organizationExists.Success)
	require.True(t, organizationExists.Exists)
	require.Equal(t, "api-key", rateLimitStatus.Identifier)
	require.Equal(t, "api", rateLimitStatus.Kind)
	require.Equal(t, "complexity", rateLimitStatus.Limits[0].Type)
	require.InDelta(t, 900, rateLimitStatus.Limits[0].RemainingAmount, 0)
	require.True(t, notifications.HasNextPage)
	require.Equal(t, &endCursor, notifications.EndCursor)
	require.Equal(t, "Mentioned you", notifications.Notifications[0].Title)
	require.Equal(t, "mentions", notifications.Notifications[0].Category)
	require.Equal(t, "actor-id", notifications.Notifications[0].ActorID)
	require.Equal(t, "external-user-id", notification.ExternalUserActorID)
	require.True(t, notificationSubscriptions.HasNextPage)
	require.Equal(t, &endCursor, notificationSubscriptions.EndCursor)
	require.Equal(t, "project", notificationSubscriptions.Subscriptions[0].TargetType)
	require.Equal(t, "Roadmap", notificationSubscriptions.Subscriptions[0].TargetName)
	require.Equal(t, "customer", notificationSubscriptions.Subscriptions[1].TargetType)
	require.Equal(t, "custom_view", notificationSubscriptions.Subscriptions[2].TargetType)
	require.Equal(t, "Cycle 7", notificationSubscriptions.Subscriptions[3].TargetName)
	require.Equal(t, "initiative", notificationSubscriptions.Subscriptions[4].TargetType)
	require.Equal(t, "label", notificationSubscriptions.Subscriptions[5].TargetType)
	require.Equal(t, "backlog", notificationSubscriptions.Subscriptions[6].ContextViewType)
	require.Equal(t, "LIT", notificationSubscriptions.Subscriptions[6].TargetName)
	require.Equal(t, "assigned", notificationSubscriptions.Subscriptions[7].UserContextViewType)
	require.Equal(t, "Ada", notificationSubscriptions.Subscriptions[7].TargetName)
	require.Equal(t, "project-id", notificationSubscription.TargetID)
	require.True(t, triageResponsibilities.HasNextPage)
	require.Equal(t, &endCursor, triageResponsibilities.EndCursor)
	require.Equal(t, "notify", triageResponsibilities.TriageResponsibilities[0].Action)
	require.Equal(t, "LIT", triageResponsibilities.TriageResponsibilities[0].TeamKey)
	require.Equal(t, "Primary rotation", triageResponsibilities.TriageResponsibilities[0].TimeScheduleName)
	require.Equal(t, "Omer", triageResponsibilities.TriageResponsibilities[0].CurrentUserName)
	require.Equal(t, []string{"user-id", "other-user-id"}, triageResponsibility.ManualUserIDs)
	require.Equal(t, []string{"user-id", "other-user-id"}, triageManualSelection.UserIDs)
	require.Equal(t, "team-id", slaConfigurations.TeamIDOrKey)
	require.Equal(t, "First response", slaConfigurations.SLAConfigurations[0].Name)
	require.InDelta(t, 3600000, slaConfigurations.SLAConfigurations[0].SLA, 0)
	require.Equal(t, "all", slaConfigurations.SLAConfigurations[0].SLAType)
	require.False(t, slaConfigurations.SLAConfigurations[0].RemovesSLA)
	require.Len(t, semanticSearch.Results, 5)
	require.Equal(t, "issue", semanticSearch.Results[0].Type)
	require.Equal(t, "issue-id", semanticSearch.Results[0].ID)
	require.Equal(t, "LIT-3", semanticSearch.Results[0].Key)
	require.Equal(t, "Search result", semanticSearch.Results[0].Title)
	require.Equal(t, "project", semanticSearch.Results[1].Type)
	require.Equal(t, "Search project", semanticSearch.Results[1].Title)
	require.Equal(t, "initiative", semanticSearch.Results[2].Type)
	require.Equal(t, "Search initiative", semanticSearch.Results[2].Title)
	require.Equal(t, "document", semanticSearch.Results[3].Type)
	require.Equal(t, "Search document", semanticSearch.Results[3].Title)
	require.Equal(t, "unknown-id", semanticSearch.Results[4].ID)
	require.Empty(t, semanticSearch.Results[4].Title)
	require.True(t, documentSearch.HasNextPage)
	require.Equal(t, &endCursor, documentSearch.EndCursor)
	require.Len(t, documentSearch.Documents, 6)
	require.Equal(t, "search-document-id", documentSearch.Documents[0].ID)
	require.Equal(t, "team", documentSearch.Documents[0].ParentType)
	require.Equal(t, "linctl", documentSearch.Documents[0].ParentName)
	require.Equal(t, "project", documentSearch.Documents[1].ParentType)
	require.Equal(t, "initiative", documentSearch.Documents[2].ParentType)
	require.Equal(t, "issue", documentSearch.Documents[3].ParentType)
	require.Equal(t, "release", documentSearch.Documents[4].ParentType)
	require.Equal(t, "cycle", documentSearch.Documents[5].ParentType)
	require.True(t, issueSearch.HasNextPage)
	require.Equal(t, &endCursor, issueSearch.EndCursor)
	require.Equal(t, "LIT-30", issueSearch.Issues[0].Identifier)
	require.Equal(t, "Pinned project", issueSearch.Issues[0].ProjectName)
	require.True(t, projectSearch.HasNextPage)
	require.Equal(t, &endCursor, projectSearch.EndCursor)
	require.Equal(t, "search-project-id", projectSearch.Projects[0].ID)
	require.Equal(t, "Search project", projectSearch.Projects[0].Name)
	require.Equal(t, "Omer", projectSearch.Projects[0].Lead)
	require.True(t, releasePipelines.HasNextPage)
	require.Equal(t, &endCursor, releasePipelines.EndCursor)
	require.Equal(t, "Production", releasePipelines.ReleasePipelines[0].Name)
	require.Equal(t, "scheduled", releasePipelines.ReleasePipelines[0].Type)
	require.True(t, releasePipeline.IsProduction)
	require.Equal(t, "template-id", releasePipeline.ReleaseNoteTemplateID)
	require.True(t, releasePipelineReleases.HasNextPage)
	require.Equal(t, &endCursor, releasePipelineReleases.EndCursor)
	require.Equal(t, "release-id", releasePipelineReleases.Releases[0].ID)
	require.Equal(t, "Production", releasePipelineReleases.Releases[0].PipelineName)
	require.True(t, releasePipelineStages.HasNextPage)
	require.Equal(t, &endCursor, releasePipelineStages.EndCursor)
	require.Equal(t, "release-stage-id", releasePipelineStages.ReleaseStages[0].ID)
	require.Equal(t, "Production", releasePipelineStages.ReleaseStages[0].PipelineName)
	require.True(t, releasePipelineTeams.HasNextPage)
	require.Equal(t, &endCursor, releasePipelineTeams.EndCursor)
	require.Equal(t, "LIT", releasePipelineTeams.Teams[0].Key)
	require.Equal(t, "Kyanite", releasePipelineTeams.Teams[0].OrgName)
	require.True(t, releaseStages.HasNextPage)
	require.Equal(t, &endCursor, releaseStages.EndCursor)
	require.Equal(t, "Started", releaseStages.ReleaseStages[0].Name)
	require.Equal(t, "started", releaseStages.ReleaseStages[0].Type)
	require.Equal(t, "Production", releaseStage.PipelineName)
	require.True(t, releaseStageReleases.HasNextPage)
	require.Equal(t, &endCursor, releaseStageReleases.EndCursor)
	require.Equal(t, "release-id", releaseStageReleases.Releases[0].ID)
	require.Equal(t, "release-stage-id", releaseStageReleases.Releases[0].StageID)
	require.True(t, releases.HasNextPage)
	require.Equal(t, &endCursor, releases.EndCursor)
	require.Equal(t, "Mobile 1.2.3", releases.Releases[0].Name)
	require.Equal(t, "Started", releases.Releases[0].StageName)
	require.Equal(t, "v1.2.3", release.Version)
	require.Equal(t, "Omer", release.CreatorName)
	require.True(t, releaseDocuments.HasNextPage)
	require.Equal(t, "Release spec", releaseDocuments.Documents[0].Title)
	require.Equal(t, "project", releaseDocuments.Documents[0].ParentType)
	require.True(t, releaseIssues.HasNextPage)
	require.Equal(t, "LIT-48", releaseIssues.Issues[0].Identifier)
	require.Equal(t, "Release issue", releaseIssues.Issues[0].Title)
	require.Equal(t, 1, release.ReleaseNoteCount)
	require.True(t, releaseHistory.HasNextPage)
	require.Equal(t, &endCursor, releaseHistory.EndCursor)
	require.Equal(t, "release-history-id", releaseHistory.History[0].ID)
	require.Equal(t, "release-id", releaseHistory.History[0].ReleaseID)
	require.Equal(t, 1, releaseHistory.History[0].EntryCount)
	require.JSONEq(t, `[{"from":"planned","to":"started","type":"stage"}]`, string(releaseHistory.History[0].Entries))
	require.True(t, releaseLinks.HasNextPage)
	require.Equal(t, &endCursor, releaseLinks.EndCursor)
	require.Equal(t, "release-link-id", releaseLinks.Links[0].ID)
	require.Equal(t, "Runbook", releaseLinks.Links[0].Label)
	require.Equal(t, "https://example.com/runbook", releaseLinks.Links[0].URL)
	require.Equal(t, "user-id", releaseLinks.Links[0].CreatorID)
	require.Equal(t, "initiative-id", releaseLinks.Links[1].InitiativeID)
	require.Equal(t, "Platform", releaseLinks.Links[1].InitiativeName)
	require.Equal(t, "project-id", releaseLinks.Links[1].ProjectID)
	require.Equal(t, "Pinned project", releaseLinks.Links[1].ProjectName)
	require.Equal(t, "release-link-id", externalLink.ID)
	require.Equal(t, "Runbook", externalLink.Label)
	require.Equal(t, "https://example.com/runbook", externalLink.URL)
	require.Equal(t, "release-id", releaseSearch.Releases[0].ID)
	require.True(t, releaseNotes.HasNextPage)
	require.Equal(t, &endCursor, releaseNotes.EndCursor)
	require.Equal(t, "Launch notes", releaseNotes.ReleaseNotes[0].Title)
	require.Equal(t, "completed", releaseNote.GenerationStatus)
	require.Equal(t, "Mobile 1.2.3", releaseNote.LastReleaseName)
	require.Equal(t, "team body", team.Description)
	require.Equal(t, "Omer", teamMembers.Members[0].DisplayName)
	require.Equal(t, &endCursor, teamMembers.EndCursor)
	require.True(t, users.HasNextPage)
	require.True(t, users.Users[0].Admin)
	require.Equal(t, "Omer", user.DisplayName)
	require.Equal(t, "Omer", viewerUser.DisplayName)
	require.True(t, drafts.HasNextPage)
	require.Equal(t, &endCursor, drafts.EndCursor)
	require.Equal(t, "issue", drafts.Drafts[0].ParentType)
	require.Equal(t, "LIT-3", drafts.Drafts[0].ParentKey)
	require.Equal(t, "project", drafts.Drafts[1].ParentType)
	require.Equal(t, "project_update", drafts.Drafts[2].ParentType)
	require.Equal(t, "initiative", drafts.Drafts[3].ParentType)
	require.Equal(t, "initiative_update", drafts.Drafts[4].ParentType)
	require.Equal(t, "comment", drafts.Drafts[5].ParentType)
	require.Equal(t, "customer_need", drafts.Drafts[6].ParentType)
	require.Equal(t, "team", drafts.Drafts[7].ParentType)
	require.Equal(t, "unknown", drafts.Drafts[8].ParentType)
	require.Equal(t, "LIT-41", userAssignedIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-42", userCreatedIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-43", userDelegatedIssues.Issues[0].Identifier)
	require.Equal(t, "team-membership-id", userTeamMemberships.Memberships[0].ID)
	require.Equal(t, "LIT", userTeams.Teams[0].Key)
	require.Equal(t, "LIT-44", viewerAssignedIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-45", viewerCreatedIssues.Issues[0].Identifier)
	require.Equal(t, "LIT-46", viewerDelegatedIssues.Issues[0].Identifier)
	require.Equal(t, "team-membership-id", viewerTeamMemberships.Memberships[0].ID)
	require.Equal(t, "LIT", viewerTeams.Teams[0].Key)
	require.True(t, workflowStates.HasNextPage)
	require.Equal(t, &endCursor, workflowStates.EndCursor)
	require.Equal(t, "Started", workflowStates.WorkflowStates[0].Name)
	require.Equal(t, "LIT", workflowStates.WorkflowStates[0].TeamKey)
	require.Equal(t, "started", workflowState.Type)
	require.Equal(t, "linctl", workflowState.TeamName)
	require.True(t, workflowStateIssues.HasNextPage)
	require.Equal(t, "workflow-state-id", workflowStateIssues.WorkflowStateID)
	require.Equal(t, "LIT-55", workflowStateIssues.Issues[0].Identifier)
	require.True(t, timeSchedules.HasNextPage)
	require.Equal(t, &endCursor, timeSchedules.EndCursor)
	require.Equal(t, "Primary on-call", timeSchedules.TimeSchedules[0].Name)
	require.Equal(t, 1, timeSchedule.EntryCount)
	require.Equal(t, "integration-id", timeSchedule.IntegrationID)
	require.Equal(t, "omer@example.com", timeSchedule.Entries[0].UserEmail)
	require.True(t, organizationLabels.HasNextPage)
	require.Equal(t, &endCursor, organizationLabels.EndCursor)
	require.Equal(t, "Bug", organizationLabels.Labels[0].Name)
	require.True(t, organizationProjectLabels.HasNextPage)
	require.Equal(t, "Roadmap", organizationProjectLabels.ProjectLabels[0].Name)
	require.True(t, organizationTeams.HasNextPage)
	require.Equal(t, "linctl", organizationTeams.Teams[0].Name)
	require.True(t, organizationTemplates.HasNextPage)
	require.Equal(t, &endCursor, organizationTemplates.EndCursor)
	require.Equal(t, "Bug report", organizationTemplates.Templates[0].Name)
	require.True(t, organizationUsers.HasNextPage)
	require.Equal(t, "Omer", organizationUsers.Users[0].DisplayName)
	require.Equal(t, 2, templates.TotalCount)
	require.Len(t, templates.Templates, 1)
	require.Equal(t, "Bug report", templates.Templates[0].Name)
	require.Equal(t, "issue", templates.Templates[0].Type)
	require.Equal(t, "LIT", templates.Templates[0].TeamKey)
	require.Equal(t, "template-id", template.ID)
	require.Equal(t, "pipeline-id", template.PipelineID)
	require.Equal(t, "creator-id", template.CreatorID)
	require.True(t, initiatives.HasNextPage)
	require.Equal(t, &endCursor, initiatives.EndCursor)
	require.Equal(t, "Platform", initiatives.Initiatives[0].Name)
	require.Equal(t, "Active", initiatives.Initiatives[0].Status)
	require.Equal(t, "2026-12-31", initiatives.Initiatives[0].TargetDate)
	require.Equal(t, "initiative-id", initiative.ID)
	require.Equal(t, "Platform initiative", initiative.Description)
	require.True(t, initiativeHistory.HasNextPage)
	require.Equal(t, &endCursor, initiativeHistory.EndCursor)
	require.Equal(t, "initiative-history-id", initiativeHistory.History[0].ID)
	require.Equal(t, "initiative-id", initiativeHistory.History[0].InitiativeID)
	require.Equal(t, 1, initiativeHistory.History[0].EntryCount)
	require.JSONEq(
		t,
		`[{"from":"Planned","to":"Active","type":"status"}]`,
		string(initiativeHistory.History[0].Entries),
	)
	require.True(t, initiativeLinks.HasNextPage)
	require.Equal(t, &endCursor, initiativeLinks.EndCursor)
	require.Equal(t, "release-link-id", initiativeLinks.Links[0].ID)
	require.Equal(t, "Runbook", initiativeLinks.Links[0].Label)
	require.True(t, subInitiatives.HasNextPage)
	require.Equal(t, "child-initiative-id", subInitiatives.Initiatives[0].ID)
	require.Equal(t, "Child platform", subInitiatives.Initiatives[0].Name)
	require.Equal(t, "Planned", subInitiatives.Initiatives[0].Status)
	require.True(t, initiativeScopedUpdates.HasNextPage)
	require.Equal(t, &endCursor, initiativeScopedUpdates.EndCursor)
	require.Equal(t, "initiative-update-id", initiativeScopedUpdates.Updates[0].ID)
	require.Equal(t, "initiative-id", initiativeScopedUpdates.Updates[0].InitiativeID)
	require.True(t, initiativeDocuments.HasNextPage)
	require.Equal(t, &endCursor, initiativeDocuments.EndCursor)
	require.Equal(t, "initiative-document-id", initiativeDocuments.Documents[0].ID)
	require.Equal(t, "Initiative spec", initiativeDocuments.Documents[0].Title)
	require.True(t, initiativeProjects.HasNextPage)
	require.Equal(t, &endCursor, initiativeProjects.EndCursor)
	require.Equal(t, "initiative-project-id", initiativeProjects.Projects[0].ID)
	require.Equal(t, "Initiative project", initiativeProjects.Projects[0].Name)
	require.True(t, initiativeRelations.HasNextPage)
	require.Equal(t, &endCursor, initiativeRelations.EndCursor)
	require.Equal(t, "initiative-relation-id", initiativeRelations.Relations[0].ID)
	require.Equal(t, "Platform", initiativeRelations.Relations[0].ParentInitiativeName)
	require.Equal(t, "Child initiative", initiativeRelations.Relations[0].RelatedInitiativeName)
	require.Equal(t, "Omer", initiativeRelations.Relations[0].DisplayName)
	require.Empty(t, initiativeRelations.Relations[1].DisplayName)
	require.Equal(t, "initiative-relation-id", initiativeRelation.ID)
	require.InEpsilon(t, 1.5, initiativeRelation.SortOrder, 0)
	require.True(t, initiativeToProjects.HasNextPage)
	require.Equal(t, &endCursor, initiativeToProjects.EndCursor)
	require.Equal(t, "initiative-to-project-id", initiativeToProjects.Associations[0].ID)
	require.Equal(t, "Platform", initiativeToProjects.Associations[0].InitiativeName)
	require.Equal(t, "Pinned project", initiativeToProjects.Associations[0].ProjectName)
	require.Equal(t, "initiative-to-project-id", initiativeToProject.ID)
	require.Equal(t, "project-id", initiativeToProject.ProjectID)
	require.True(t, roadmapToProjects.HasNextPage)
	require.Equal(t, &endCursor, roadmapToProjects.EndCursor)
	require.Equal(t, "roadmap-to-project-id", roadmapToProjects.Associations[0].ID)
	require.Equal(t, "Platform roadmap", roadmapToProjects.Associations[0].RoadmapName)
	require.Equal(t, "Pinned project", roadmapToProjects.Associations[0].ProjectName)
	require.Equal(t, "roadmap-to-project-id", roadmapToProject.ID)
	require.Equal(t, "project-id", roadmapToProject.ProjectID)
	require.True(t, initiativeUpdates.HasNextPage)
	require.Equal(t, &endCursor, initiativeUpdates.EndCursor)
	require.Equal(t, "initiative-update-id", initiativeUpdates.Updates[0].ID)
	require.Equal(t, "Platform", initiativeUpdates.Updates[0].InitiativeName)
	require.Equal(t, 1, initiativeUpdates.Updates[0].CommentCount)
	require.Equal(t, "initiative-update-id", initiativeUpdate.ID)
	require.Equal(t, "First initiative update", initiativeUpdate.Body)
	require.Equal(t, "initiative-update-id", initiativeUpdateComments.InitiativeUpdateID)
	require.Equal(t, "comment-id", initiativeUpdateComments.Comments[0].ID)
	require.Equal(t, &endCursor, initiativeUpdateComments.EndCursor)
	require.True(t, roadmaps.HasNextPage)
	require.Equal(t, &endCursor, roadmaps.EndCursor)
	require.Equal(t, "Platform roadmap", roadmaps.Roadmaps[0].Name)
	require.Equal(t, "roadmap-id", roadmap.ID)
	require.Equal(t, "Owner", roadmap.OwnerDisplayName)
	require.True(t, roadmapProjects.HasNextPage)
	require.Equal(t, "roadmap-id", roadmapProjects.RoadmapID)
	require.Equal(t, "Roadmap project", roadmapProjects.Projects[0].Name)
	require.True(t, customViews.HasNextPage)
	require.Equal(t, &endCursor, customViews.EndCursor)
	require.Equal(t, "My issues", customViews.CustomViews[0].Name)
	require.Equal(t, "Issue", customViews.CustomViews[0].ModelName)
	require.True(t, customViews.CustomViews[0].Shared)
	require.Equal(t, "custom-view-id", customViewSubscribers.ID)
	require.True(t, customViewSubscribers.HasSubscribers)
	require.Equal(t, "custom-view-id", customView.ID)
	require.Equal(t, "Saved issue view", customView.Description)
	require.True(t, customViewInitiatives.HasNextPage)
	require.Equal(t, "Platform", customViewInitiatives.Initiatives[0].Name)
	require.True(t, customViewIssues.HasNextPage)
	require.Equal(t, "LIT-1", customViewIssues.Issues[0].Identifier)
	require.Equal(t, "view-preferences-id", customViewOrganizationPreferences.ID)
	require.Equal(t, "organization", customViewOrganizationPreferences.Type)
	require.Equal(t, "list", customViewOrganizationPreferences.Values.Layout)
	require.True(t, customViewOrganizationPreferenceValues.HasOrganizationPreferences)
	require.Equal(t, "priority", customViewOrganizationPreferenceValues.ViewOrdering)
	require.True(t, customViewOrganizationPreferenceValues.ShowArchivedItems)
	require.Equal(t, []string{"column-id"}, customViewOrganizationPreferenceValues.HiddenColumns)
	require.True(t, customViewProjects.HasNextPage)
	require.Equal(t, "Custom view project", customViewProjects.Projects[0].Name)
	require.Equal(t, "view-preferences-id", customViewUserPreferences.ID)
	require.Equal(t, "user", customViewUserPreferences.Type)
	require.True(t, customViewUserPreferenceValues.HasUserPreferences)
	require.Equal(t, "updatedAt", customViewUserPreferenceValues.ViewOrdering)
	require.True(t, customViewPreferenceValues.HasEffectivePreferenceValue)
	require.Equal(t, "board", customViewPreferenceValues.Layout)
	require.Equal(t, "updatedAt", customViewPreferenceValues.ViewOrdering)
	require.True(t, customers.HasNextPage)
	require.Equal(t, &endCursor, customers.EndCursor)
	require.Equal(t, "Acme", customers.Customers[0].Name)
	require.Equal(t, "Active", customers.Customers[0].StatusName)
	require.InDelta(t, 3, customers.Customers[0].ApproximateNeedCount, 0)
	require.Equal(t, "customer-id", customer.ID)
	require.Equal(t, "Enterprise", customer.TierName)
	require.Equal(t, "Omer", customer.OwnerDisplayName)
	require.True(t, customerNeeds.HasNextPage)
	require.Equal(t, &endCursor, customerNeeds.EndCursor)
	require.Equal(t, "Acme", customerNeeds.Needs[0].CustomerName)
	require.Equal(t, "LIT-1", customerNeed.Issue)
	require.Equal(t, "Need content", customerNeed.Content)
	require.Equal(t, "customer-need-id", customerNeedProjectAttachment.CustomerNeedID)
	require.NotNil(t, customerNeedProjectAttachment.Attachment)
	require.Equal(t, "project-attachment-id", customerNeedProjectAttachment.Attachment.ID)
	require.True(t, customerStatuses.HasNextPage)
	require.Equal(t, &endCursor, customerStatuses.EndCursor)
	require.Equal(t, "Active", customerStatuses.Statuses[0].DisplayName)
	require.Equal(t, "customer-status-id", customerStatus.ID)
	require.Equal(t, "#00ff00", customerStatus.Color)
	require.True(t, customerTiers.HasNextPage)
	require.Equal(t, &endCursor, customerTiers.EndCursor)
	require.Equal(t, "Enterprise", customerTiers.Tiers[0].DisplayName)
	require.Equal(t, "customer-tier-id", customerTier.ID)
	require.Equal(t, "#0000ff", customerTier.Color)
	require.True(t, favorites.HasNextPage)
	require.Equal(t, &endCursor, favorites.EndCursor)
	require.Equal(t, "issue", favorites.Favorites[0].Type)
	require.True(t, favoriteChildren.HasNextPage)
	require.Equal(t, &endCursor, favoriteChildren.EndCursor)
	require.Equal(t, "favorite-child-id", favoriteChildren.Favorites[0].ID)
	require.Equal(t, "project", favoriteChildren.Favorites[0].Type)
	require.Equal(t, "favorite-id", favorite.ID)
	require.Equal(t, "https://linear.app/kyanite/issue/LIT-1", favorite.URL)
	require.True(t, emojis.HasNextPage)
	require.Equal(t, &endCursor, emojis.EndCursor)
	require.Equal(t, "party", emojis.Emojis[0].Name)
	require.Equal(t, "custom", emojis.Emojis[0].Source)
	require.Equal(t, "emoji-id", emoji.ID)
	require.Equal(t, "party", emoji.Name)
	require.True(t, attachments.HasNextPage)
	require.Equal(t, &endCursor, attachments.EndCursor)
	require.Equal(t, "Linked PR", attachments.Attachments[0].Title)
	require.Equal(t, "github", attachments.Attachments[0].SourceType)
	require.True(t, attachmentsForURL.HasNextPage)
	require.Equal(t, &endCursor, attachmentsForURL.EndCursor)
	require.Equal(t, "Linked URL", attachmentsForURL.Attachments[0].Title)
	require.Equal(t, "url", attachmentsForURL.Attachments[0].SourceType)
	require.Equal(t, "attachment-id", attachment.ID)
	require.Equal(t, "feat: add thing", attachment.Subtitle)
	require.Equal(t, "LIT-41", attachmentIssue.Identifier)
	require.Equal(t, "project-attachment-id", attachmentIssueAttachments.Attachments[0].ID)
	require.Equal(t, "issue-id", attachmentIssueBotActor.IssueID)
	require.Equal(t, "bot-actor-id", attachmentIssueBotActor.Bot.ID)
	require.Equal(t, "LIT-42", attachmentIssueChildren.Issues[0].Identifier)
	require.Equal(t, "Attachment issue spec", attachmentIssueDocuments.Documents[0].Title)
	require.Equal(t, "project-attachment-id", attachmentIssueFormerAttachments.Attachments[0].ID)
	require.Equal(t, "issue-history-id", attachmentIssueHistory.History[0].ID)
	require.Equal(t, "issue-relation-id", attachmentIssueInverseRelations.Relations[0].ID)
	require.Equal(t, "Bug", attachmentIssueLabels.Labels[0].Name)
	require.Equal(t, "issue-relation-id", attachmentIssueRelations.Relations[0].ID)
	require.Equal(t, "release-id", attachmentIssueReleases.Releases[0].ID)
	require.Equal(t, "issue-id", attachmentIssueStateHistory.IssueID)
	require.Equal(t, "issue-state-span-id", attachmentIssueStateHistory.Spans[0].ID)
	require.Equal(t, "Omer", attachmentIssueSubscribers.Users[0].DisplayName)
}

func Test_CheckOrganizationExists_returns_operation_errors(t *testing.T) {
	_, err := CheckOrganizationExists(context.Background(), fakeGraphQLClient{}, "missing")

	require.Error(t, err)
	require.Contains(t, err.Error(), "missing fake response for organizationExists")
}

func Test_GetRateLimitStatus_returns_operation_errors(t *testing.T) {
	_, err := GetRateLimitStatus(context.Background(), fakeGraphQLClient{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "missing fake response for rateLimitStatus")
}

func Test_ClientReadHelpers_cover_nil_actor_bot_summaries(t *testing.T) {
	require.Nil(t, actorBotSummary(nil))
	require.Nil(t, commentActorBotSummary(nil))
	require.Nil(t, issueActorBotSummary(nil))
	require.Nil(t, issueVCSBranchActorBotSummary(nil))
	require.Nil(t, attachmentIssueActorBotSummary(nil))
}

func Test_ClientReadScenarios_return_not_found_for_null_vcs_branch_issue(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issueVcsBranchSearch":                   `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_attachments":       `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_botActor":          `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_children":          `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_documents":         `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_formerAttachments": `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_history":           `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_inverseRelations":  `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_labels":            `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_relations":         `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_releases":          `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_stateHistory":      `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_subscribers":       `{"issueVcsBranchSearch":null}`,
	}

	_, err := GetIssueByVCSBranch(context.Background(), graphqlClient, "missing/branch")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchAttachments(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = GetIssueVCSBranchBotActor(context.Background(), graphqlClient, "missing/branch")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchChildren(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchDocuments(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchFormerAttachments(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchHistory(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchInverseRelations(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchLabels(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchRelations(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchReleases(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchStateHistory(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
	_, err = ListIssueVCSBranchSubscribers(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
}

func Test_ClientReadScenarios_rank_next_issues(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"NextIssuesByTeam": `{"issues":{"nodes":[` +
			nextIssueJSON("LIT-31", "Low priority standalone", 4, "Low", "2026-01-01T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-32", "Urgent standalone", 1, "Urgent", "2026-02-01T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-33", "Older high standalone", 2, "High", "2026-01-15T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-34", "Newer high standalone", 2, "High", "2026-02-15T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-35", "No priority standalone", 0, "No priority", "2026-01-01T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-36", "Unblocks active work", 3, "Normal", "2026-03-01T00:00:00Z", []string{
				`{"type":"blocks","relatedIssue":{"id":"active-1","state":{"type":"started"}}}`,
				`{"type":"blocks","relatedIssue":{"id":"done-1","state":{"type":"completed"}}}`,
				`{"type":"relates","relatedIssue":{"id":"active-2","state":{"type":"unstarted"}}}`,
			}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	issues, err := ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 6)

	require.NoError(t, err)
	require.Equal(t, []string{"LIT-36", "LIT-32", "LIT-33", "LIT-34", "LIT-31", "LIT-35"}, issueIdentifiers(issues.Issues))
	require.Equal(t, 1, issues.Issues[0].UnblocksCount)
}
