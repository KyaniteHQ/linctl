package client

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
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
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-11",
			Title:      "detail issue",
			StateID:    "done",
			State:      "Done",
			StateType:  "completed",
		}) + `}`,
		"issue_attachments":       `{"issue":{"attachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_children":          `{"issue":{"children":{"nodes":[` + issueJSON(issueFixture{Identifier: "LIT-31", Title: "child issue", StateID: "todo", State: "Todo", StateType: "unstarted"}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_documents":         `{"issue":{"documents":{"nodes":[{"id":"issue-document-id","title":"Issue spec","slugId":"issue-spec","archivedAt":null,"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_formerAttachments": `{"issue":{"formerAttachments":{"nodes":[` + projectAttachmentJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_history":           `{"issue":{"history":{"nodes":[` + issueHistoryJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_inverseRelations":  `{"issue":{"inverseRelations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_labels":            `{"issue":{"labels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_relations":         `{"issue":{"relations":{"nodes":[` + issueRelationJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"issue_releases":          `{"issue":{"releases":{"nodes":[` + releaseJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"Projects": `{"team":{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
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
		"ProjectUpdates":            `{"project":{"id":"project-id","name":"detail","projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectUpdates":            `{"projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","project":{"id":"project-id","name":"detail"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectUpdate":             `{"projectUpdate":{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","project":{"id":"project-id","name":"detail"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"projectUpdate_comments":    `{"projectUpdate":{"id":"project-update-id","comments":{"nodes":[` + commentMetadataJSON("", "project-update-id", "user-id") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"project_projectMilestones": `{"project":{"id":"project-id","name":"detail","projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectMilestone":          `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}}`,
		"projectMilestone_issues": `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-52",
			Title:      "Milestone issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectStatuses":       `{"projectStatuses":{"nodes":[{"id":"project-status-id","name":"Backlog","description":"Ready for planning","type":"backlog","color":"#bec2c8","position":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectStatus":         `{"projectStatus":{"id":"project-status-id","name":"Backlog","description":"Ready for planning","type":"backlog","color":"#bec2c8","position":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z"}}`,
		"projectLabels":         `{"projectLabels":{"nodes":[` + projectLabelJSON("project-label-id", "Roadmap") + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectLabel":          `{"projectLabel":{"id":"project-label-id","name":"Roadmap","description":"Project label","color":"#f2c94c","isGroup":false,"lastAppliedAt":"2026-06-19T12:00:00Z","retiredAt":null,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","parent":{"id":"parent-project-label-id","name":"Parent","color":"#828282"}}}`,
		"projectLabel_children": `{"projectLabel":{"id":"project-label-id","name":"Roadmap","children":{"nodes":[{"id":"child-project-label-id","name":"Mobile","description":"Child project label","color":"#56ccf2","isGroup":false,"lastAppliedAt":null,"retiredAt":null,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","parent":{"id":"project-label-id","name":"Roadmap","color":"#f2c94c"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectLabel_projects": `{"projectLabel":{"id":"project-label-id","name":"Roadmap","projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectRelations":    `{"projectRelations":{"nodes":[{"id":"project-relation-id","type":"blocks","anchorType":"project","relatedAnchorType":"project","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":null,"relatedProject":{"id":"related-project-id","name":"Related project"},"relatedProjectMilestone":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"project-relation-no-user","type":"relates","anchorType":"project","relatedAnchorType":"project","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone"},"relatedProject":{"id":"other-related-project-id","name":"Other related"},"relatedProjectMilestone":{"id":"related-project-milestone-id","name":"Related milestone"},"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectRelation":     `{"projectRelation":{"id":"project-relation-id","type":"blocks","anchorType":"project","relatedAnchorType":"project","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":null,"relatedProject":{"id":"related-project-id","name":"Related project"},"relatedProjectMilestone":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"issueRelations":      `{"issueRelations":{"nodes":[{"id":"issue-relation-id","type":"blocks","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issueRelation":       `{"issueRelation":{"id":"issue-relation-id","type":"blocks","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}}`,
		"issueToReleases":     `{"issueToReleases":{"nodes":[{"id":"issue-to-release-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id"},"release":{"id":"release-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issueToRelease":      `{"issueToRelease":{"id":"issue-to-release-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"issue":{"id":"issue-id"},"release":{"id":"release-id"}}}`,
		"applicationInfo":     `{"applicationInfo":{"id":"app-id","clientId":"app-client-id","name":"Demo App","description":"Demo authorization app","developer":"Kyanite","developerUrl":"https://example.com","imageUrl":"https://example.com/app.png"}}`,
		"issue_comments":      `{"issue":{"id":"issue-id","identifier":"LIT-12","comments":{"nodes":[{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":"parent-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"bot-comment-id","body":"bot note","url":"https://linear.app/comment/bot-comment-id","createdAt":"2026-06-19T12:01:00Z","parentId":null,"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"comments":            `{"comments":{"nodes":[{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:02:00Z","editedAt":"2026-06-19T12:02:00Z","resolvedAt":null,"parentId":"parent-id","issueId":"issue-id","projectId":null,"projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"bot-comment-id","body":"bot note","url":"https://linear.app/comment/bot-comment-id","createdAt":"2026-06-19T12:01:00Z","updatedAt":"2026-06-19T12:01:00Z","editedAt":null,"resolvedAt":null,"parentId":null,"issueId":null,"projectId":"project-id","projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"comment":             `{"comment":{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:02:00Z","editedAt":"2026-06-19T12:02:00Z","resolvedAt":null,"parentId":"parent-id","issueId":"issue-id","projectId":null,"projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"Documents":           `{"documents":{"nodes":[{"id":"document-id","title":"Spec","slugId":"spec","archivedAt":null,"project":{"id":"project-id","name":"fixture"},"team":null,"issue":null,"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"document":            `{"document":{"id":"document-id","title":"Team note","slugId":"team-note","archivedAt":null,"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null}}`,
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
		"team_templates":        `{"team":{"id":"team-id","key":"LIT","name":"linctl","templates":{"nodes":[` + templateJSON() + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
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
		"viewer_teamMemberships":       `{"viewer":{"teamMemberships":{"nodes":[{"id":"team-membership-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"owner":true,"sortOrder":1.5,"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"viewer_teams":                 `{"viewer":{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"workflowStates":               `{"workflowStates":{"nodes":[{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"workflowState":                `{"workflowState":{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}}`,
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
		"initiativeRelations":      `{"initiativeRelations":{"nodes":[{"id":"initiative-relation-id","sortOrder":1.5,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"relatedInitiative":{"id":"child-initiative-id","name":"Child initiative"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"initiative-relation-no-user","sortOrder":2,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"relatedInitiative":{"id":"other-child-initiative-id","name":"Other child"},"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeRelation":       `{"initiativeRelation":{"id":"initiative-relation-id","sortOrder":1.5,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"relatedInitiative":{"id":"child-initiative-id","name":"Child initiative"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"initiativeToProjects":     `{"initiativeToProjects":{"nodes":[{"id":"initiative-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeToProject":      `{"initiativeToProject":{"id":"initiative-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"initiative":{"id":"initiative-id","name":"Platform"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}}`,
		"roadmapToProjects":        `{"roadmapToProjects":{"nodes":[{"id":"roadmap-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"roadmap":{"id":"roadmap-id","name":"Platform roadmap"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"roadmapToProject":         `{"roadmapToProject":{"id":"roadmap-to-project-id","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,"roadmap":{"id":"roadmap-id","name":"Platform roadmap"},"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}}}`,
		"initiativeUpdates":        `{"initiativeUpdates":{"nodes":[{"id":"initiative-update-id","body":"First initiative update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/initiative-update/initiative-update-id","slugId":"initiative-update-slug","commentCount":1,"initiative":{"id":"initiative-id","name":"Platform"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"initiativeUpdate":         `{"initiativeUpdate":{"id":"initiative-update-id","body":"First initiative update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/initiative-update/initiative-update-id","slugId":"initiative-update-slug","commentCount":1,"initiative":{"id":"initiative-id","name":"Platform"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"roadmaps":                 `{"roadmaps":{"nodes":[{"id":"roadmap-id","name":"Platform roadmap","description":"Roadmap body","color":"#5e6ad2","slugId":"platform-roadmap","sortOrder":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","url":"https://linear.app/kyanite/roadmap/platform-roadmap","creator":{"id":"user-id","displayName":"Omer"},"owner":{"id":"owner-id","displayName":"Owner"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"roadmap":                  `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","description":"Roadmap body","color":"#5e6ad2","slugId":"platform-roadmap","sortOrder":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","url":"https://linear.app/kyanite/roadmap/platform-roadmap","creator":{"id":"user-id","displayName":"Omer"},"owner":{"id":"owner-id","displayName":"Owner"}}}`,
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
	issue, err := GetIssueByID(context.Background(), graphqlClient, "LIT-11")
	require.NoError(t, err)
	issueAttachments, err := ListIssueAttachments(context.Background(), graphqlClient, "LIT-1", 2)
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
	allProjectUpdates, err := ListAllProjectUpdates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectUpdate, err := GetProjectUpdateByID(context.Background(), graphqlClient, "project-update-id")
	require.NoError(t, err)
	projectUpdateComments, err := ListProjectUpdateComments(context.Background(), graphqlClient, "project-update-id", 2)
	require.NoError(t, err)
	projectMilestones, err := ListProjectMilestones(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectMilestone, err := GetProjectMilestoneByID(context.Background(), graphqlClient, "project-milestone-id")
	require.NoError(t, err)
	projectMilestoneIssues, err := ListProjectMilestoneIssues(context.Background(), graphqlClient, "project-milestone-id", 2)
	require.NoError(t, err)
	projectStatuses, err := ListProjectStatuses(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectStatus, err := GetProjectStatusByID(context.Background(), graphqlClient, "project-status-id")
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
	documents, err := ListDocuments(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	document, err := GetDocumentByID(context.Background(), graphqlClient, "document-id")
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
	roadmaps, err := ListRoadmaps(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	roadmap, err := GetRoadmapByID(context.Background(), graphqlClient, "roadmap-id")
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
	require.Equal(t, "LIT-11", issue.Identifier)
	require.True(t, issueAttachments.HasNextPage)
	require.Equal(t, &endCursor, issueAttachments.EndCursor)
	require.Equal(t, "project-attachment-id", issueAttachments.Attachments[0].ID)
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
	require.Equal(t, &endCursor, projectUpdates.EndCursor)
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
	require.Equal(t, "project-milestone-id", projectMilestone.ID)
	require.Equal(t, "Launch milestone", projectMilestone.Name)
	require.Equal(t, "next", projectMilestone.Status)
	require.True(t, projectMilestoneIssues.HasNextPage)
	require.Equal(t, "project-milestone-id", projectMilestoneIssues.ProjectMilestoneID)
	require.Equal(t, "LIT-52", projectMilestoneIssues.Issues[0].Identifier)
	require.True(t, projectStatuses.HasNextPage)
	require.Equal(t, &endCursor, projectStatuses.EndCursor)
	require.Equal(t, "project-status-id", projectStatuses.ProjectStatuses[0].ID)
	require.Equal(t, "Backlog", projectStatuses.ProjectStatuses[0].Name)
	require.Equal(t, "backlog", projectStatuses.ProjectStatuses[0].Type)
	require.Equal(t, "#bec2c8", projectStatuses.ProjectStatuses[0].Color)
	require.Equal(t, "project-status-id", projectStatus.ID)
	require.Equal(t, "Ready for planning", projectStatus.Description)
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
	require.True(t, documents.HasNextPage)
	require.Equal(t, "project", documents.Documents[0].ParentType)
	require.Equal(t, "Team note", document.Title)
	require.Equal(t, "team", document.ParentType)
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
	require.True(t, roadmaps.HasNextPage)
	require.Equal(t, &endCursor, roadmaps.EndCursor)
	require.Equal(t, "Platform roadmap", roadmaps.Roadmaps[0].Name)
	require.Equal(t, "roadmap-id", roadmap.ID)
	require.Equal(t, "Owner", roadmap.OwnerDisplayName)
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

func Test_ClientWriteScenarios_guard_writes_and_report_results(t *testing.T) {
	// Given
	t.Run("invalid requests fail before network", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{})

		_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{ID: "LIT-1"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:          "LIT-1",
			Description: "description",
			Append:      "append",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{Title: "missing id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{ID: "LIT-1"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{Body: "body"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = StartIssue(context.Background(), graphqlClient, matchingTarget(), "")
		require.Error(t, err)

		_, err = CloseIssue(context.Background(), graphqlClient, matchingTarget(), "")
		require.Error(t, err)

		_, err = CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{ID: "project-id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateProjectMilestone(context.Background(), graphqlClient, matchingTarget(), ProjectMilestoneCreateRequest{
			Name: "name",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateProjectMilestone(context.Background(), graphqlClient, matchingTarget(), ProjectMilestoneCreateRequest{
			ProjectID: "project-id",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateProjectMilestone(
			context.Background(),
			graphqlClient,
			matchingTarget(),
			ProjectMilestoneUpdateRequest{ID: "project-milestone-id"},
		)
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateCycle(context.Background(), graphqlClient, matchingTarget(), CycleCreateRequest{
			EndsAt: "2026-07-15T00:00:00Z",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CreateCycle(context.Background(), graphqlClient, matchingTarget(), CycleCreateRequest{
			StartsAt: "2026-07-01T00:00:00Z",
		})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateCycle(context.Background(), graphqlClient, matchingTarget(), CycleUpdateRequest{ID: "cycle-id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = ArchiveCycle(context.Background(), graphqlClient, matchingTarget(), "")
		require.ErrorIs(t, err, ErrWriteInvalid)
	})

	t.Run("issue comment succeeds", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-12",
				Title:      "comment target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueCommentCreate": `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","issue":` + issueJSON(issueFixture{
				Identifier: "LIT-12",
				Title:      "comment target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}}`,
		})

		comment, err := CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{
			ID:   "LIT-12",
			Body: "hello",
		})

		require.NoError(t, err)
		require.Equal(t, "comment-id", comment.ID)
		require.Equal(t, "LIT-12", comment.Issue.Identifier)
	})

	t.Run("issue update succeeds", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-21",
				Title:      "update target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-21",
				Title:      "updated",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}`,
		})

		issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:    "LIT-21",
			Title: "updated",
		})

		require.NoError(t, err)
		require.Equal(t, "updated", issue.Title)
	})

	t.Run("issue update appends to description", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSONWithDescription(issueFixture{
				Identifier: "LIT-22",
				Title:      "append target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}, "Existing description") + `}`,
			"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-22",
				Title:      "append target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}`,
		})

		issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:     "LIT-22",
			Append: "Progress note",
		})

		require.NoError(t, err)
		require.Equal(t, "append target", issue.Title)
		require.Equal(t, "Progress note", appendIssueDescription("", "Progress note"))
		require.Equal(t, "Existing description\n\nProgress note", appendIssueDescription("Existing description\n", "Progress note"))
	})

	t.Run("project update and archive succeed", func(t *testing.T) {
		graphqlClient := projectWriteFakeClient(map[string]string{
			"project": `{"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "fixture",
				Status: "Backlog",
			}) + `}`,
			"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "updated",
				Status: "Started",
			}) + `}}`,
			"ProjectArchive": `{"projectArchive":{"success":true,"entity":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "updated",
				Status: "Canceled",
			}) + `}}`,
		})

		project, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
			ID:   "project-id",
			Name: "updated",
		})
		require.NoError(t, err)
		require.Equal(t, "updated", project.Name)

		project, err = ArchiveProject(context.Background(), graphqlClient, matchingTarget(), "project-id")
		require.NoError(t, err)
		require.Equal(t, "Canceled", project.Status.Name)
	})
}

func Test_SummaryMappingScenarios_preserve_optional_people(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issue": `{"issue":` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-30",
			Title:      "assigned",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}, "Omer") + `}`,
		"project": `{"project":` + projectJSONWithLead(projectFixture{
			ID:     "project-id",
			Name:   "led",
			Status: "Backlog",
		}, "Omer") + `}`,
	}

	issue, err := GetIssueByID(context.Background(), graphqlClient, "LIT-30")
	require.NoError(t, err)
	require.Equal(t, "Omer", issue.Assignee)

	project, err := GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.NoError(t, err)
	require.Equal(t, "Omer", project.Lead)
}

func Test_SummaryMappingScenarios_preserve_reference_domain_variants(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"document": `{"document":{
			"id":"document-id",
			"title":"Cycle spec",
			"slugId":"cycle-spec",
			"archivedAt":"2026-06-19T12:00:00Z",
			"project":{"id":"project-id","name":"Pinned project"},
			"team":{"id":"team-id","key":"LIT","name":"linctl"},
			"issue":{"id":"issue-id","identifier":"LIT-1","title":"Issue"},
			"cycle":{"id":"cycle-id","number":7,"name":"Planning"}
		}}`,
		"team": `{"team":{
			"id":"team-id",
			"key":"LIT",
			"name":"linctl",
			"description":null,
			"archivedAt":"2026-06-19T12:00:00Z",
			"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
		}}`,
	}

	document, err := GetDocumentByID(context.Background(), graphqlClient, "document-id")
	require.NoError(t, err)
	require.Equal(t, "2026-06-19T12:00:00Z", document.ArchivedAt)
	require.Equal(t, "cycle", document.ParentType)
	require.Equal(t, "Planning", document.ParentName)

	team, err := GetTeamByID(context.Background(), graphqlClient, "team-id")
	require.NoError(t, err)
	require.Empty(t, team.Description)
	require.Equal(t, "2026-06-19T12:00:00Z", team.ArchivedAt)
}

func Test_SummaryMappingScenarios_preserve_release_note_without_generation_status(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"releaseNote": `{"releaseNote":` + strings.Replace(
			releaseNoteJSON(),
			`"generationStatus":"completed"`,
			`"generationStatus":null`,
			1,
		) + `}`,
	}

	note, err := GetReleaseNoteByID(context.Background(), graphqlClient, "release-note-id")

	require.NoError(t, err)
	require.Empty(t, note.GenerationStatus)
	require.Equal(t, "Launch notes", note.Title)
}

func Test_ClientFailureScenarios_wrap_read_and_mutation_errors(t *testing.T) {
	t.Run("read operations wrap graphql errors", func(t *testing.T) {
		graphqlClient := errorGraphQLClient{err: errors.New("network down")}

		_, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssues(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{StateType: "started"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{ProjectID: "project-id"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{AssigneeID: "user-id"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{LabelID: "label-id"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{CycleID: "cycle-id"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{CreatedAfter: "2026-06-01"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{CreatedBefore: "2026-06-30"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{HasBlockers: true})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{Blocks: true})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list next issues")

		_, err = ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1, IssueListFilters{BlockedBy: "LIT-1"})
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = SearchIssuesByTeam(context.Background(), graphqlClient, "team-id", "needle", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "search issues")

		_, err = GetIssueByID(context.Background(), graphqlClient, "LIT-1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")

		_, err = GetIssueDependencies(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue dependencies LIT-1")

		_, err = ListIssueAttachments(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue attachments LIT-1")

		_, err = ListIssueChildren(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue children LIT-1")

		_, err = ListIssueDocuments(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue documents LIT-1")

		_, err = ListIssueFormerAttachments(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue former attachments LIT-1")

		_, err = ListIssueHistory(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue history LIT-1")

		_, err = ListIssueInverseRelations(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue inverse relations LIT-1")

		_, err = ListIssueLabels(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue labels LIT-1")

		_, err = ListIssueRelationsForIssue(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue relations LIT-1")

		_, err = ListIssueReleases(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue releases LIT-1")

		_, err = ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list projects")

		_, err = GetProjectByID(context.Background(), graphqlClient, "project-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project project-id")

		_, err = ListProjectMembers(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project members project-id")

		_, err = ListProjectAttachments(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project attachments project-id")

		_, err = ListProjectDocuments(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project documents project-id")

		_, err = ListProjectExternalLinks(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project external links project-id")

		_, err = ListProjectHistory(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project history project-id")

		_, err = ListProjectInitiativeToProjects(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project initiative associations project-id")

		_, err = ListProjectInitiatives(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project initiatives project-id")

		_, err = ListProjectInverseRelations(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project inverse relations project-id")

		_, err = ListProjectIssues(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project issues project-id")

		_, err = ListLabelsForProject(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project labels project-id")

		_, err = ListProjectNeeds(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project customer needs project-id")

		_, err = ListProjectRelationsForProject(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project relations project-id")

		_, err = ListProjectTeams(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project teams project-id")

		_, err = ListProjectComments(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project comments project-id")

		_, err = ListProjectUpdates(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project updates project-id")

		_, err = ListAllProjectUpdates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project updates")

		_, err = GetProjectUpdateByID(context.Background(), graphqlClient, "project-update-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project update project-update-id")

		_, err = ListProjectUpdateComments(context.Background(), graphqlClient, "project-update-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project update comments project-update-id")

		_, err = ListProjectMilestones(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project milestones project-id")

		_, err = GetProjectMilestoneByID(context.Background(), graphqlClient, "project-milestone-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project milestone project-milestone-id")

		_, err = ListProjectMilestoneIssues(context.Background(), graphqlClient, "project-milestone-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project milestone issues project-milestone-id")

		_, err = ListProjectStatuses(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project statuses")

		_, err = GetProjectStatusByID(context.Background(), graphqlClient, "project-status-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project status project-status-id")

		_, err = ListProjectLabels(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project labels")

		_, err = GetProjectLabelByID(context.Background(), graphqlClient, "project-label-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project label project-label-id")

		_, err = ListProjectLabelChildren(context.Background(), graphqlClient, "project-label-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project label children project-label-id")

		_, err = ListProjectLabelProjects(context.Background(), graphqlClient, "project-label-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project label projects project-label-id")

		_, err = ListProjectRelations(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project relations")

		_, err = GetProjectRelationByID(context.Background(), graphqlClient, "project-relation-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project relation project-relation-id")

		_, err = ListIssueRelations(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue relations")

		_, err = GetIssueRelationByID(context.Background(), graphqlClient, "issue-relation-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue relation issue-relation-id")

		_, err = ListIssueToReleases(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue to releases")

		_, err = GetIssueToReleaseByID(context.Background(), graphqlClient, "issue-to-release-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue to release issue-to-release-id")

		_, err = ListTeamMemberships(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team memberships")

		_, err = GetTeamMembershipByID(context.Background(), graphqlClient, "team-membership-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get team membership team-membership-id")

		_, err = GetApplicationInfo(context.Background(), graphqlClient, "app-client-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get application info app-client-id")

		_, err = ListAgentActivities(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list agent activities")

		_, err = GetAgentActivityByID(context.Background(), graphqlClient, "agent-activity-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get agent activity agent-activity-id")

		_, err = ListAgentSkills(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list agent skills")

		_, err = GetAgentSkillByID(context.Background(), graphqlClient, "agent-skill-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get agent skill agent-skill-id")

		_, err = ListExternalUsers(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list external users")

		_, err = GetExternalUserByID(context.Background(), graphqlClient, "external-user-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get external user external-user-id")

		_, err = ListAuditEntryTypes(context.Background(), graphqlClient)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list audit entry types")

		_, err = ListIssueComments(context.Background(), graphqlClient, "LIT-1", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue comments LIT-1")

		_, err = ListComments(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list comments")

		_, err = GetCommentByID(context.Background(), graphqlClient, "comment-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get comment comment-id")

		_, err = ListDocuments(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list documents")

		_, err = GetDocumentByID(context.Background(), graphqlClient, "document-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get document document-id")

		_, err = ListLabels(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list labels")

		_, err = GetLabelByID(context.Background(), graphqlClient, "label-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get label label-id")

		_, err = ListLabelChildren(context.Background(), graphqlClient, "label-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list label children label-id")

		_, err = ListLabelIssues(context.Background(), graphqlClient, "label-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list label issues label-id")

		_, err = ListTeams(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list teams")

		_, err = GetTeamByID(context.Background(), graphqlClient, "team-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get team team-id")

		_, err = ListTeamMembers(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team members team-id")

		_, err = ListTeamCycles(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team cycles team-id")

		_, err = ListTeamIssues(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team issues team-id")

		_, err = ListTeamLabels(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team labels team-id")

		_, err = ListTeamMembershipsForTeam(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team memberships team-id")

		_, err = ListTeamProjects(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team projects team-id")

		_, err = ListTeamReleasePipelines(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team release pipelines team-id")

		_, err = ListTeamWorkflowStates(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team states team-id")

		_, err = ListTeamTemplates(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team templates team-id")

		_, err = ListUsers(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list users")

		_, err = GetUserByID(context.Background(), graphqlClient, "user-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get user user-id")

		_, err = GetViewerUser(context.Background(), graphqlClient)
		require.Error(t, err)
		require.Contains(t, err.Error(), "get viewer user")

		_, err = ListViewerDrafts(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list viewer drafts")

		_, err = ListUserAssignedIssues(context.Background(), graphqlClient, "user-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list user assigned issues user-id")

		_, err = ListUserCreatedIssues(context.Background(), graphqlClient, "user-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list user created issues user-id")

		_, err = ListUserDelegatedIssues(context.Background(), graphqlClient, "user-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list user delegated issues user-id")

		_, err = ListUserTeamMemberships(context.Background(), graphqlClient, "user-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list user team memberships user-id")

		_, err = ListUserTeams(context.Background(), graphqlClient, "user-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list user teams user-id")

		_, err = ListViewerAssignedIssues(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list viewer assigned issues")

		_, err = ListViewerCreatedIssues(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list viewer created issues")

		_, err = ListViewerDelegatedIssues(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list viewer delegated issues")

		_, err = ListViewerTeamMemberships(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list viewer team memberships")

		_, err = ListViewerTeams(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list viewer teams")

		_, err = ListWorkflowStates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list workflow states")

		_, err = GetWorkflowStateByID(context.Background(), graphqlClient, "workflow-state-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get workflow state workflow-state-id")

		_, err = ListTimeSchedules(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list time schedules")

		_, err = GetTimeScheduleByID(context.Background(), graphqlClient, "time-schedule-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get time schedule time-schedule-id")

		_, err = ListOrganizationLabels(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list organization labels")

		_, err = ListOrganizationProjectLabels(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list organization project labels")

		_, err = ListOrganizationTeams(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list organization teams")

		_, err = ListOrganizationTemplates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list organization templates")

		_, err = ListOrganizationUsers(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list organization users")

		_, err = ListTemplates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list templates")

		_, err = GetTemplateByID(context.Background(), graphqlClient, "template-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get template template-id")

		_, err = ListInitiatives(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiatives")

		_, err = GetInitiativeByID(context.Background(), graphqlClient, "initiative-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get initiative initiative-id")

		_, err = ListInitiativeHistory(context.Background(), graphqlClient, "initiative-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative history initiative-id")

		_, err = ListInitiativeLinks(context.Background(), graphqlClient, "initiative-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative links initiative-id")

		_, err = ListSubInitiatives(context.Background(), graphqlClient, "initiative-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative sub-initiatives initiative-id")

		_, err = ListInitiativeUpdatesForInitiative(context.Background(), graphqlClient, "initiative-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative updates initiative-id")

		_, err = ListInitiativeDocuments(context.Background(), graphqlClient, "initiative-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative documents initiative-id")

		_, err = ListInitiativeProjects(context.Background(), graphqlClient, "initiative-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative projects initiative-id")

		_, err = ListInitiativeRelations(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative relations")

		_, err = GetInitiativeRelationByID(context.Background(), graphqlClient, "initiative-relation-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get initiative relation initiative-relation-id")

		_, err = ListInitiativeToProjects(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative to projects")

		_, err = GetInitiativeToProjectByID(context.Background(), graphqlClient, "initiative-to-project-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get initiative to project initiative-to-project-id")

		_, err = ListRoadmapToProjects(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list roadmap to projects")

		_, err = GetRoadmapToProjectByID(context.Background(), graphqlClient, "roadmap-to-project-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get roadmap to project roadmap-to-project-id")

		_, err = ListInitiativeUpdates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list initiative updates")

		_, err = GetInitiativeUpdateByID(context.Background(), graphqlClient, "initiative-update-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get initiative update initiative-update-id")

		_, err = ListRoadmaps(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list roadmaps")

		_, err = GetRoadmapByID(context.Background(), graphqlClient, "roadmap-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get roadmap roadmap-id")

		_, err = ListCustomViews(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list custom views")

		_, err = GetCustomViewSubscriberStatus(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view subscribers custom-view-id")

		_, err = GetCustomViewByID(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view custom-view-id")

		_, err = ListCustomViewInitiatives(context.Background(), graphqlClient, "custom-view-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list custom view initiatives custom-view-id")

		_, err = ListCustomViewIssues(context.Background(), graphqlClient, "custom-view-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list custom view issues custom-view-id")

		_, err = GetCustomViewOrganizationPreferences(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view organization preferences custom-view-id")

		_, err = GetCustomViewOrganizationPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view organization preference values custom-view-id")

		_, err = ListCustomViewProjects(context.Background(), graphqlClient, "custom-view-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list custom view projects custom-view-id")

		_, err = GetCustomViewUserPreferences(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view user preferences custom-view-id")

		_, err = GetCustomViewUserPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view user preference values custom-view-id")

		_, err = GetCustomViewPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get custom view preference values custom-view-id")

		_, err = ListCustomers(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list customers")

		_, err = GetCustomerByID(context.Background(), graphqlClient, "customer-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get customer customer-id")

		_, err = ListCustomerNeeds(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list customer needs")

		_, err = GetCustomerNeedByID(context.Background(), graphqlClient, "customer-need-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get customer need customer-need-id")

		_, err = ListCustomerStatuses(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list customer statuses")

		_, err = GetCustomerStatusByID(context.Background(), graphqlClient, "customer-status-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get customer status customer-status-id")

		_, err = ListCustomerTiers(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list customer tiers")

		_, err = GetCustomerTierByID(context.Background(), graphqlClient, "customer-tier-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get customer tier customer-tier-id")

		_, err = ListFavorites(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list favorites")

		_, err = ListFavoriteChildren(context.Background(), graphqlClient, "favorite-folder-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list favorite children favorite-folder-id")

		_, err = GetFavoriteByID(context.Background(), graphqlClient, "favorite-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get favorite favorite-id")

		_, err = ListEmojis(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list emojis")

		_, err = GetEmojiByID(context.Background(), graphqlClient, "emoji-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get emoji emoji-id")

		_, err = ListNotifications(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list notifications")

		_, err = GetNotificationByID(context.Background(), graphqlClient, "notification-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get notification notification-id")

		_, err = ListNotificationSubscriptions(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list notification subscriptions")

		_, err = GetNotificationSubscriptionByID(context.Background(), graphqlClient, "notification-subscription-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get notification subscription notification-subscription-id")

		_, err = ListTriageResponsibilities(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list triage responsibilities")

		_, err = GetTriageResponsibilityByID(context.Background(), graphqlClient, "triage-responsibility-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get triage responsibility triage-responsibility-id")

		_, err = GetTriageResponsibilityManualSelection(context.Background(), graphqlClient, "triage-responsibility-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get triage responsibility manual selection triage-responsibility-id")

		_, err = ListSLAConfigurations(context.Background(), graphqlClient, "team-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "list SLA configurations team-id")

		_, err = SearchSemantic(context.Background(), graphqlClient, "agent search", 2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "semantic search")

		_, err = ListReleasePipelines(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release pipelines")

		_, err = GetReleasePipelineByID(context.Background(), graphqlClient, "release-pipeline-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get release pipeline release-pipeline-id")

		_, err = ListReleasePipelineReleases(context.Background(), graphqlClient, "release-pipeline-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release pipeline releases release-pipeline-id")

		_, err = ListReleasePipelineStages(context.Background(), graphqlClient, "release-pipeline-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release pipeline stages release-pipeline-id")

		_, err = ListReleasePipelineTeams(context.Background(), graphqlClient, "release-pipeline-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release pipeline teams release-pipeline-id")

		_, err = ListReleaseStages(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release stages")

		_, err = GetReleaseStageByID(context.Background(), graphqlClient, "release-stage-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get release stage release-stage-id")

		_, err = ListReleaseStageReleases(context.Background(), graphqlClient, "release-stage-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release stage releases release-stage-id")

		_, err = ListReleases(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list releases")

		_, err = GetReleaseByID(context.Background(), graphqlClient, "release-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get release release-id")

		_, err = ListReleaseHistory(context.Background(), graphqlClient, "release-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release history release-id")

		_, err = ListReleaseDocuments(context.Background(), graphqlClient, "release-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release documents release-id")

		_, err = ListReleaseIssues(context.Background(), graphqlClient, "release-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release issues release-id")

		_, err = ListReleaseLinks(context.Background(), graphqlClient, "release-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release links release-id")

		_, err = GetEntityExternalLinkByID(context.Background(), graphqlClient, "release-link-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get external link release-link-id")

		_, err = SearchReleases(context.Background(), graphqlClient, "mobile", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "search releases")

		_, err = ListReleaseNotes(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list release notes")

		_, err = GetReleaseNoteByID(context.Background(), graphqlClient, "release-note-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get release note release-note-id")

		_, err = ListAttachments(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list attachments")

		_, err = ListAttachmentsForURL(context.Background(), graphqlClient, "https://example.com/spec", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list attachments for url https://example.com/spec")

		_, err = GetAttachmentByID(context.Background(), graphqlClient, "attachment-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get attachment attachment-id")
	})

	t.Run("issue mutations fail when payload omits entity", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"IssueCreate": `{"issueCreate":{"success":false,"issue":null}}`,
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-20",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate":        `{"issueUpdate":{"success":false,"issue":null}}`,
			"IssueCommentCreate": `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"body","url":"url","issue":null}}}`,
			"CompletedWorkflowStates": `{"workflowStates":{"nodes":[
				{"id":"done-state","name":"Done","type":"completed","position":1}
			]}}`,
			"StartedWorkflowStates": `{"workflowStates":{"nodes":[
				{"id":"started-state","name":"Started","type":"started","position":1}
			]}}`,
			"IssueClose": `{"issueUpdate":{"success":false,"issue":null}}`,
		})

		_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{Title: "title"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{ID: "LIT-20", Title: "title"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{ID: "LIT-20", Body: "body"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = StartIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-20")
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = CloseIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-20")
		require.ErrorIs(t, err, ErrMutationFailed)
	})

	t.Run("project mutations fail when payload omits entity", func(t *testing.T) {
		graphqlClient := projectWriteFakeClient(map[string]string{
			"ProjectCreate": `{"projectCreate":{"success":false,"project":null}}`,
			"project": `{"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "fixture",
				Status: "Backlog",
			}) + `}`,
			"ProjectUpdate":  `{"projectUpdate":{"success":false,"project":null}}`,
			"ProjectArchive": `{"projectArchive":{"success":false,"entity":null}}`,
		})

		_, err := CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{Name: "name"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = ArchiveProject(context.Background(), graphqlClient, matchingTarget(), "project-id")
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = CreateProjectMilestone(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"project":                `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
				"ProjectMilestoneCreate": `{"projectMilestoneCreate":{"success":false,"projectMilestone":null}}`,
			}),
			matchingTarget(),
			ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateProjectMilestone(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"projectMilestone": `{"projectMilestone":` +
					projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
				"ProjectMilestoneUpdate": `{"projectMilestoneUpdate":{"success":false,"projectMilestone":null}}`,
			}),
			matchingTarget(),
			ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = CreateCycle(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"CycleCreate": `{"cycleCreate":{"success":false,"cycle":null}}`,
			}),
			matchingTarget(),
			CycleCreateRequest{StartsAt: "2026-07-01T00:00:00Z", EndsAt: "2026-07-15T00:00:00Z"},
		)
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateCycle(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"cycle":       `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
				"CycleUpdate": `{"cycleUpdate":{"success":false,"cycle":null}}`,
			}),
			matchingTarget(),
			CycleUpdateRequest{ID: "cycle-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = ArchiveCycle(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"cycle":        `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
				"CycleArchive": `{"cycleArchive":{"success":false,"entity":null}}`,
			}),
			matchingTarget(),
			"cycle-id",
		)
		require.ErrorIs(t, err, ErrMutationFailed)
	})

	t.Run("write operations wrap graphql operation errors", func(t *testing.T) {
		operationErr := errors.New("linear unavailable")

		_, err := CreateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueCreate": "",
		}).withError(operationErr), matchingTarget(), IssueCreateRequest{Title: "title"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create issue")

		_, err = UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate": "",
		}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-40", Title: "title"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update issue LIT-40")

		_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueCommentCreate": "",
		}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-40", Body: "body"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "comment on issue LIT-40")

		_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"StartedWorkflowStates": `{"workflowStates":{"nodes":[{"id":"started-state","name":"Started","type":"started","position":1}]}}`,
			"IssueUpdate":           "",
		}).withError(operationErr), matchingTarget(), "LIT-40")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "start issue LIT-40")

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"CompletedWorkflowStates": `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`,
			"IssueClose":              "",
		}).withError(operationErr), matchingTarget(), "LIT-40")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "close issue LIT-40")

		_, err = CreateProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectCreate": "",
		}).withError(operationErr), matchingTarget(), ProjectCreateRequest{Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create project")

		_, err = UpdateProject(context.Background(), projectWriteFakeClient(map[string]string{
			"project":       `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectUpdate": "",
		}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update project project-id")

		_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
			"project":        `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectArchive": "",
		}).withError(operationErr), matchingTarget(), "project-id")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "archive project project-id")

		_, err = CreateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"project":                `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectMilestoneCreate": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create project milestone")

		_, err = UpdateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"projectMilestone": `{"projectMilestone":` +
				projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
			"ProjectMilestoneUpdate": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update project milestone project-milestone-id")

		_, err = CreateCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"CycleCreate": "",
		}).withError(operationErr), matchingTarget(), CycleCreateRequest{
			StartsAt: "2026-07-01T00:00:00Z",
			EndsAt:   "2026-07-15T00:00:00Z",
		})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create cycle")

		_, err = UpdateCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"cycle":       `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
			"CycleUpdate": "",
		}).withError(operationErr), matchingTarget(), CycleUpdateRequest{ID: "cycle-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update cycle cycle-id")

		_, err = ArchiveCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"cycle":        `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
			"CycleArchive": "",
		}).withError(operationErr), matchingTarget(), "cycle-id")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "archive cycle cycle-id")
	})

	t.Run("write operations return guard read errors", func(t *testing.T) {
		operationErr := errors.New("guard read failed")

		_, err := UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": "",
		}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-50", Title: "title"})
		require.ErrorIs(t, err, operationErr)

		_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": "",
		}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-50", Body: "body"})
		require.ErrorIs(t, err, operationErr)

		_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": "",
		}).withError(operationErr), matchingTarget(), "LIT-50")
		require.ErrorIs(t, err, operationErr)

		_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-51",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"StartedWorkflowStates": "",
		}).withError(operationErr), matchingTarget(), "LIT-51")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "list started workflow states")

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": "",
		}).withError(operationErr), matchingTarget(), "LIT-50")
		require.ErrorIs(t, err, operationErr)

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-51",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"CompletedWorkflowStates": "",
		}).withError(operationErr), matchingTarget(), "LIT-51")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "list completed workflow states")

		_, err = UpdateProject(context.Background(), projectWriteFakeClient(map[string]string{
			"project": "",
		}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
			"project": "",
		}).withError(operationErr), matchingTarget(), "project-id")
		require.ErrorIs(t, err, operationErr)

		_, err = CreateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"project": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = UpdateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"projectMilestone": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = UpdateCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"cycle": "",
		}).withError(operationErr), matchingTarget(), CycleUpdateRequest{ID: "cycle-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = ArchiveCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"cycle": "",
		}).withError(operationErr), matchingTarget(), "cycle-id")
		require.ErrorIs(t, err, operationErr)
	})

	t.Run("write operations refuse unpinned targets", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{})
		emptyTarget := config.Target{}

		_, err := CreateIssue(context.Background(), graphqlClient, emptyTarget, IssueCreateRequest{Title: "title"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = UpdateIssue(context.Background(), graphqlClient, emptyTarget, IssueUpdateRequest{ID: "LIT-1", Title: "title"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CommentOnIssue(context.Background(), graphqlClient, emptyTarget, IssueCommentRequest{ID: "LIT-1", Body: "body"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = StartIssue(context.Background(), graphqlClient, emptyTarget, "LIT-1")
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CloseIssue(context.Background(), graphqlClient, emptyTarget, "LIT-1")
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CreateProject(context.Background(), graphqlClient, emptyTarget, ProjectCreateRequest{Name: "name"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = UpdateProject(context.Background(), graphqlClient, emptyTarget, ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = ArchiveProject(context.Background(), graphqlClient, emptyTarget, "project-id")
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CreateProjectMilestone(
			context.Background(),
			graphqlClient,
			emptyTarget,
			ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = UpdateProjectMilestone(
			context.Background(),
			graphqlClient,
			emptyTarget,
			ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CreateCycle(
			context.Background(),
			graphqlClient,
			emptyTarget,
			CycleCreateRequest{StartsAt: "2026-07-01T00:00:00Z", EndsAt: "2026-07-15T00:00:00Z"},
		)
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = UpdateCycle(
			context.Background(),
			graphqlClient,
			emptyTarget,
			CycleUpdateRequest{ID: "cycle-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = ArchiveCycle(context.Background(), graphqlClient, emptyTarget, "cycle-id")
		require.ErrorIs(t, err, ErrTargetMismatch)
	})
}

func Test_TargetFailureScenarios_refuse_unpinned_or_mismatched_targets(t *testing.T) {
	_, err := ResolveTarget(context.Background(), fakeGraphQLClient{}, config.Target{})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer": `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":  "",
	}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve teams")

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"other-org","name":"Other","urlKey":"other"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": "",
	}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve project")

	graphqlClient := fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"other-team","key":"ABC","name":"other","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}

	_, err = ResolveTarget(context.Background(), graphqlClient, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = firstCompletedStateID(context.Background(), fakeGraphQLClient{
		"CompletedWorkflowStates": `{"workflowStates":{"nodes":[]}}`,
	}, "team-id")
	require.ErrorIs(t, err, ErrWriteInvalid)

	_, err = firstStartedStateID(context.Background(), fakeGraphQLClient{
		"StartedWorkflowStates": `{"workflowStates":{"nodes":[]}}`,
	}, "team-id")
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = requireTargetMatch(config.Target{OrgID: "org-id", TeamID: "team-id", TeamKey: "LIT"}, config.Target{
		OrgID:   "other-org",
		TeamID:  "team-id",
		TeamKey: "LIT",
	})
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_TransportScenarios_return_actionable_errors(t *testing.T) {
	require.Equal(t, "fallback", firstNonEmpty("", "fallback"))
	require.Equal(t, "primary", firstNonEmpty("primary", "fallback"))
	require.Equal(t, 3*time.Second, defaultDuration(3*time.Second, time.Second))
	require.Equal(t, time.Second, defaultDuration(0, time.Second))
	require.Equal(t, 200*time.Millisecond, retryDelay("", 1))
	require.Equal(t, 100*time.Millisecond, retryDelay("not-a-number", 0))
	require.Equal(t, 2*time.Second, retryDelay("2", 0))

	response := graphql.Response{}
	err := decodeGraphQLResponse([]byte("not json"), http.StatusOK, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode graphql response")

	err = decodeGraphQLResponse([]byte("server down"), http.StatusBadGateway, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "graphql http status 502")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = waitForRateLimitRetry(ctx, http.StatusTooManyRequests, http.Header{}, 0, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wait for retry")
}

func Test_CustomViewPreferenceReads_return_empty_values_when_organization_defaults_are_absent(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient{
		"customView_organizationViewPreferences":             `{"customView":{"organizationViewPreferences":null}}`,
		"customView_organizationViewPreferences_preferences": `{"customView":{"organizationViewPreferences":null}}`,
	}

	// When
	preferences, err := GetCustomViewOrganizationPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	values, err := GetCustomViewOrganizationPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)

	// Then
	require.Equal(t, "custom-view-id", preferences.CustomViewID)
	require.Empty(t, preferences.ID)
	require.Equal(t, "custom-view-id", values.CustomViewID)
	require.False(t, values.HasOrganizationPreferences)
}

func Test_CustomViewPreferenceReads_return_empty_values_when_user_preferences_are_absent(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient{
		"customView_userViewPreferences":             `{"customView":{"userViewPreferences":null}}`,
		"customView_userViewPreferences_preferences": `{"customView":{"userViewPreferences":null}}`,
	}

	// When
	preferences, err := GetCustomViewUserPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	values, err := GetCustomViewUserPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)

	// Then
	require.Equal(t, "custom-view-id", preferences.CustomViewID)
	require.Empty(t, preferences.ID)
	require.Equal(t, "custom-view-id", values.CustomViewID)
	require.False(t, values.HasUserPreferences)
}

type errorGraphQLClient struct {
	err error
}

func (client errorGraphQLClient) MakeRequest(
	_ context.Context,
	_ *graphql.Request,
	_ *graphql.Response,
) error {
	return client.err
}

type operationErrorFakeClient struct {
	responses map[string]string
	err       error
}

func (client operationErrorFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if client.responses[request.OpName] == "" {
		return client.err
	}

	return fakeGraphQLClient(client.responses).MakeRequest(ctx, request, response)
}

func (client issueWriteFakeClient) withError(err error) operationErrorFakeClient {
	return operationErrorFakeClient{
		responses: client.withTargetResponses(),
		err:       err,
	}
}

func (client projectWriteFakeClient) withError(err error) operationErrorFakeClient {
	return operationErrorFakeClient{
		responses: client.withTargetResponses(),
		err:       err,
	}
}

func Test_WriteGuardScenarios_refuse_mismatched_resources(t *testing.T) {
	guard := writeGuard{
		target: ResolvedTarget{
			Team: TargetTeam{ID: "team-id", Key: "LIT"},
		},
	}
	graphqlClient := fakeGraphQLClient{
		"issue": `{"issue":` + strings.ReplaceAll(issueJSON(issueFixture{
			Identifier: "ABC-1",
			Title:      "wrong team",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
		"project": `{"project":` + strings.ReplaceAll(projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "wrong-team",
			Status: "Backlog",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
	}

	_, err := guard.requireIssue(context.Background(), graphqlClient, "ABC-1")
	require.ErrorIs(t, err, ErrTargetMismatch)

	err = guard.requireProject(context.Background(), graphqlClient, "project-id")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = newWriteGuard(context.Background(), errorGraphQLClient{err: errors.New("resolve failed")}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve failed")

	_, err = guard.requireIssue(context.Background(), errorGraphQLClient{err: errors.New("read issue failed")}, "LIT-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read issue failed")

	err = guard.requireProject(context.Background(), errorGraphQLClient{err: errors.New("read project failed")}, "project-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read project failed")
}

func Test_FakeGraphQLClient_respects_context_and_missing_operations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := fakeGraphQLClient{}.MakeRequest(ctx, &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)

	err = fakeGraphQLClient{}.MakeRequest(context.Background(), &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)
	require.True(t, errors.Is(err, errors.New("missing fake response for Viewer")) || strings.Contains(err.Error(), "missing fake response"))
}

func Test_TargetScenarios_allow_unpinned_project_and_matching_team(t *testing.T) {
	require.Nil(t, optionalString(""))
	require.Equal(t, "value", *optionalString("value"))
	require.Equal(t, "value", *stringPtr("value"))
	require.Equal(t, 7, *intPtr(7))
	require.True(t, *boolPtr(true))
	require.Nil(t, issueDependencyParent(nil))
	require.True(t, projectHasTeam(ProjectSummary{Teams: []ProjectTeam{{ID: "team-id", Key: "LIT"}}}, "team-id", "LIT"))
	require.False(t, projectHasTeam(ProjectSummary{Teams: []ProjectTeam{{ID: "team-id", Key: "ABC"}}}, "team-id", "LIT"))

	guard, err := newWriteGuard(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.NoError(t, err)
	require.Nil(t, guard.target.Project)

	err = validateProjectUpdateRequest(ProjectUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = validateProjectMilestoneUpdateRequest(ProjectMilestoneUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = validateCycleUpdateRequest(CycleUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func issueJSONWithAssignee(issue issueFixture, assignee string) string {
	return strings.ReplaceAll(issueJSON(issue), `"assignee":null`, `"assignee":{"id":"user-id","name":"omer","displayName":"`+assignee+`"}`)
}

func issueJSONWithDescription(issue issueFixture, description string) string {
	return strings.Replace(issueJSON(issue), `"id":"issue-id",`, `"id":"issue-id","description":"`+description+`",`, 1)
}

func nextIssueJSON(
	identifier string,
	title string,
	priority int,
	priorityLabel string,
	createdAt string,
	relations []string,
) string {
	return strings.TrimSuffix(issueJSON(issueFixture{
		Identifier: identifier,
		Title:      title,
		StateID:    "todo",
		State:      "Todo",
		StateType:  "unstarted",
	}), "\n\t}") +
		`,
		"priority":` + strconv.Itoa(priority) + `,
		"priorityLabel":"` + priorityLabel + `",
		"createdAt":"` + createdAt + `",
		"relations":{"nodes":[` + strings.Join(relations, ",") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}
	}`
}

func issueIdentifiers(issues []IssueSummary) []string {
	identifiers := make([]string, 0, len(issues))
	for _, issue := range issues {
		identifiers = append(identifiers, issue.Identifier)
	}

	return identifiers
}

func projectJSONWithLead(project projectFixture, lead string) string {
	return strings.Replace(projectJSON(project), `"lead":null`, `"lead":{"id":"user-id","name":"omer","displayName":"`+lead+`"}`, 1)
}

func notificationJSON() string {
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
		"externalUserActor":{"id":"external-user-id"}
	}`
}

func notificationSubscriptionJSON() string {
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

func draftJSON(parentType string) string {
	payload := `{
		"id":"draft-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"issue":null,
		"project":null,
		"projectUpdate":null,
		"initiative":null,
		"initiativeUpdate":null,
		"parentComment":null,
		"customerNeed":null,
		"team":null
	}`
	switch parentType {
	case "issue":
		return strings.Replace(payload, `"issue":null`, `"issue":{"id":"issue-id","identifier":"LIT-3","title":"Draft issue"}`, 1)
	case "project":
		return strings.Replace(payload, `"project":null`, `"project":{"id":"project-id","name":"Draft project"}`, 1)
	case "project_update":
		return strings.Replace(payload, `"projectUpdate":null`, `"projectUpdate":{"id":"project-update-id"}`, 1)
	case "initiative":
		return strings.Replace(payload, `"initiative":null`, `"initiative":{"id":"initiative-id","name":"Draft initiative"}`, 1)
	case "initiative_update":
		return strings.Replace(payload, `"initiativeUpdate":null`, `"initiativeUpdate":{"id":"initiative-update-id"}`, 1)
	case "comment":
		return strings.Replace(payload, `"parentComment":null`, `"parentComment":{"id":"comment-id"}`, 1)
	case "customer_need":
		return strings.Replace(payload, `"customerNeed":null`, `"customerNeed":{"id":"customer-need-id"}`, 1)
	case "team":
		return strings.Replace(payload, `"team":null`, `"team":{"id":"team-id","key":"LIT","name":"linctl"}`, 1)
	default:
		return payload
	}
}

func triageResponsibilityJSON() string {
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

func notificationSubscriptionTargetJSON(
	typename string,
	targetField string,
	targetPayload string,
	withContextView bool,
	withUserContextView bool,
) string {
	contextViewType := "null"
	if withContextView {
		contextViewType = `"backlog"`
	}
	userContextViewType := "null"
	if withUserContextView {
		userContextViewType = `"assigned"`
	}

	payload := `{
		"__typename":"` + typename + `",
		"id":"notification-subscription-id",
		"active":true,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"contextViewType":` + contextViewType + `,
		"userContextViewType":` + userContextViewType + `,
		"subscriber":{"id":"user-id","displayName":"Omer"},
		"customer":null,
		"customView":null,
		"cycle":null,
		"initiative":null,
		"label":null,
		"project":null,
		"team":null,
		"user":null
	}`

	return strings.Replace(payload, `"`+targetField+`":null`, `"`+targetField+`":`+targetPayload, 1)
}

func subInitiativeJSON() string {
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

func initiativeHistoryJSON() string {
	return `{
		"id":"initiative-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Planned","to":"Active"}],
		"initiative":{"id":"initiative-id"}
	}`
}

func initiativeUpdateJSON() string {
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

func customViewPreferencesJSON(ordering string, layout string) string {
	return customViewScopedPreferencesJSON("organization", ordering, layout)
}

func customViewScopedPreferencesJSON(scope string, ordering string, layout string) string {
	return `{
		"id":"view-preferences-id",
		"createdAt":"2026-06-01T12:00:00Z",
		"updatedAt":"2026-06-01T12:01:00Z",
		"archivedAt":null,
		"type":"` + scope + `",
		"viewType":"customView",
		"preferences":` + customViewPreferenceValuesJSON(ordering, layout) + `
	}`
}

func customViewPreferenceValuesJSON(ordering string, layout string) string {
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

func templateJSON() string {
	return `{
		"id":"template-id",
		"name":"Bug report",
		"type":"issue",
		"description":"Bug report template",
		"icon":"bug",
		"color":"#ff0000",
		"templateData":{"title":"Bug: "},
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

func agentSkillJSON() string {
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

func externalUserJSON() string {
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

func agentActivityJSON(contentType string) string {
	content := map[string]string{
		"action": `{
			"__typename":"AgentActivityActionContent",
			"type":"action",
			"action":"read_file",
			"parameter":"README.md",
			"result":"Read file"
		}`,
		"elicitation": `{
			"__typename":"AgentActivityElicitationContent",
			"type":"elicitation",
			"body":"Need a choice"
		}`,
		"error": `{
			"__typename":"AgentActivityErrorContent",
			"type":"error",
			"body":"Tool failed",
			"reasonCode":"tool_error"
		}`,
		"prompt": `{
			"__typename":"AgentActivityPromptContent",
			"type":"prompt",
			"body":"Please continue"
		}`,
		"response": `{
			"__typename":"AgentActivityResponseContent",
			"type":"response",
			"body":"Done"
		}`,
		"thought": `{
			"__typename":"AgentActivityThoughtContent",
			"type":"thought",
			"body":"Thinking"
		}`,
	}[contentType]

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
		"content":` + content + `
	}`
}

func slaConfigurationJSON() string {
	return `{
		"id":"sla-configuration-id",
		"name":"First response",
		"conditions":{"priority":{"eq":1}},
		"sla":3600000,
		"slaType":"all",
		"removesSla":false
	}`
}

func semanticSearchResultJSON(resultType string) string {
	switch resultType {
	case "issue":
		return `{
			"id":"issue-id",
			"type":"issue",
			"issue":{"id":"issue-id","identifier":"LIT-3","title":"Search result","url":"https://linear.app/kyanite/issue/LIT-3"},
			"project":null,
			"initiative":null,
			"document":null
		}`
	case "project":
		return `{
			"id":"project-id",
			"type":"project",
			"issue":null,
			"project":{"id":"project-id","name":"Search project","url":"https://linear.app/kyanite/project/search-project"},
			"initiative":null,
			"document":null
		}`
	case "initiative":
		return `{
			"id":"initiative-id",
			"type":"initiative",
			"issue":null,
			"project":null,
			"initiative":{"id":"initiative-id","name":"Search initiative","url":"https://linear.app/kyanite/initiative/search-initiative"},
			"document":null
		}`
	case "document":
		return `{
			"id":"document-id",
			"type":"document",
			"issue":null,
			"project":null,
			"initiative":null,
			"document":{"id":"document-id","title":"Search document","url":"https://linear.app/kyanite/document/search-document"}
		}`
	default:
		return `{
			"id":"unknown-id",
			"type":"issue",
			"issue":null,
			"project":null,
			"initiative":null,
			"document":null
		}`
	}
}

func releasePipelineJSON() string {
	return `{
		"id":"release-pipeline-id",
		"name":"Production",
		"slugId":"production",
		"type":"scheduled",
		"isProduction":true,
		"autoGenerateReleaseNotesOnCompletion":true,
		"includePathPatterns":["services/api/**"],
		"approximateReleaseCount":4,
		"trashed":null,
		"releaseNoteTemplate":{"id":"template-id"},
		"latestReleaseNote":{"id":"release-note-id"},
		"url":"https://linear.app/kyanite/releases/production",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null
	}`
}

func trashedReleasePipelineJSON() string {
	return strings.Replace(releasePipelineJSON(), `"trashed":null`, `"trashed":false`, 1)
}

func releaseStageJSON() string {
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

func releaseJSON() string {
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

func releaseHistoryJSON() string {
	return `{
		"id":"release-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"stage","from":"planned","to":"started"}],
		"release":{"id":"release-id"}
	}`
}

func projectHistoryJSON() string {
	return `{
		"id":"project-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Backlog","to":"Started"}],
		"project":{"id":"project-id"}
	}`
}

func issueHistoryJSON() string {
	return `{
		"id":"issue-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"actorId":"user-id",
		"updatedDescription":true,
		"issue":{"id":"issue-id"}
	}`
}

func commentMetadataJSON(projectID string, projectUpdateID string, userID string) string {
	user := `null`
	if userID != "" {
		user = `{"id":"` + userID + `","name":"omer","displayName":"Omer"}`
	}

	return `{
		"id":"comment-id",
		"url":"https://linear.app/comment/comment-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:02:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":null,
		"issueId":null,
		"projectId":` + nullableStringJSON(projectID) + `,
		"projectUpdateId":` + nullableStringJSON(projectUpdateID) + `,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":` + user + `
	}`
}

func nullableStringJSON(value string) string {
	if value == "" {
		return `null`
	}

	return `"` + value + `"`
}

func projectAttachmentJSON() string {
	return `{
		"id":"project-attachment-id",
		"title":"Project link",
		"subtitle":"overview",
		"url":"https://example.com/project-link",
		"sourceType":"github"
	}`
}

func issueRelationJSON() string {
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

func projectLabelJSON(id string, name string) string {
	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":"Project label",
		"color":"#f2c94c",
		"isGroup":false,
		"lastAppliedAt":"2026-06-19T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"parent":null
	}`
}

func projectRelationJSON() string {
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

func customerNeedJSON() string {
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

func entityExternalLinkJSON() string {
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

func entityExternalLinkWithParentsJSON() string {
	return `{
		"id":"release-link-parent-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"url":"https://example.com/plan",
		"label":"Plan",
		"sortOrder":2,
		"creator":null,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func releaseNoteJSON() string {
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
