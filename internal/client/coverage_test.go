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
		"AllTeamIssues": `{"issues":{"nodes":[` + issueJSON(issueFixture{
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
		"IssueSearch": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-13",
			Title:      "search result",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"IssueByID": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-11",
			Title:      "detail issue",
			StateID:    "done",
			State:      "Done",
			StateType:  "completed",
		}) + `}`,
		"Projects": `{"team":{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"ProjectByID": `{"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "detail",
			Status: "Started",
		}) + `}`,
		"ProjectMembers":       `{"project":{"id":"project-id","name":"detail","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"ProjectUpdates":       `{"project":{"id":"project-id","name":"detail","projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"projectUpdates":       `{"projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","project":{"id":"project-id","name":"detail"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"projectUpdate":        `{"projectUpdate":{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","project":{"id":"project-id","name":"detail"},"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"ProjectMilestones":    `{"project":{"id":"project-id","name":"detail","projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"ProjectMilestoneByID": `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}}`,
		"IssueComments":        `{"issue":{"id":"issue-id","identifier":"LIT-12","comments":{"nodes":[{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":"parent-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"bot-comment-id","body":"bot note","url":"https://linear.app/comment/bot-comment-id","createdAt":"2026-06-19T12:01:00Z","parentId":null,"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"comments":             `{"comments":{"nodes":[{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:02:00Z","editedAt":"2026-06-19T12:02:00Z","resolvedAt":null,"parentId":"parent-id","issueId":"issue-id","projectId":null,"projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}},{"id":"bot-comment-id","body":"bot note","url":"https://linear.app/comment/bot-comment-id","createdAt":"2026-06-19T12:01:00Z","updatedAt":"2026-06-19T12:01:00Z","editedAt":null,"resolvedAt":null,"parentId":null,"issueId":null,"projectId":"project-id","projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"comment":              `{"comment":{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:02:00Z","editedAt":"2026-06-19T12:02:00Z","resolvedAt":null,"parentId":"parent-id","issueId":"issue-id","projectId":null,"projectUpdateId":null,"initiativeId":null,"initiativeUpdateId":null,"documentContentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}}`,
		"Documents":            `{"documents":{"nodes":[{"id":"document-id","title":"Spec","slugId":"spec","archivedAt":null,"project":{"id":"project-id","name":"fixture"},"team":null,"issue":null,"cycle":null}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"document":             `{"document":{"id":"document-id","title":"Team note","slugId":"team-note","archivedAt":null,"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null}}`,
		"IssueLabels":          `{"issueLabels":{"nodes":[{"id":"label-id","name":"Bug","description":"label body","color":"#ff0000","isGroup":false,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"issueLabel":           `{"issueLabel":{"id":"label-id","name":"Bug","description":null,"color":"#ff0000","isGroup":false,"team":null}}`,
		"Teams":                `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"TeamByID":             `{"team":{"id":"team-id","key":"LIT","name":"linctl","description":"team body","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"TeamMembers":          `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"Users":                `{"users":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"UserByID":             `{"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}}`,
		"ViewerUser":           `{"viewer":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":true}}`,
		"workflowStates":       `{"workflowStates":{"nodes":[{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"workflowState":        `{"workflowState":{"id":"workflow-state-id","name":"Started","type":"started","color":"#f2c94c","position":2,"team":{"id":"team-id","key":"LIT","name":"linctl"}}}`,
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
	projects, err := ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	project, err := GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.NoError(t, err)
	members, err := ListProjectMembers(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectUpdates, err := ListProjectUpdates(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	allProjectUpdates, err := ListAllProjectUpdates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	projectUpdate, err := GetProjectUpdateByID(context.Background(), graphqlClient, "project-update-id")
	require.NoError(t, err)
	projectMilestones, err := ListProjectMilestones(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)
	projectMilestone, err := GetProjectMilestoneByID(context.Background(), graphqlClient, "project-milestone-id")
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
	teams, err := ListTeams(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	team, err := GetTeamByID(context.Background(), graphqlClient, "team-id")
	require.NoError(t, err)
	teamMembers, err := ListTeamMembers(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	users, err := ListUsers(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	user, err := GetUserByID(context.Background(), graphqlClient, "user-id")
	require.NoError(t, err)
	viewerUser, err := GetViewerUser(context.Background(), graphqlClient)
	require.NoError(t, err)
	workflowStates, err := ListWorkflowStates(context.Background(), graphqlClient, 2)
	require.NoError(t, err)
	workflowState, err := GetWorkflowStateByID(context.Background(), graphqlClient, "workflow-state-id")
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
	require.True(t, projects.HasNextPage)
	require.Equal(t, "listed", projects.Projects[0].Name)
	require.Equal(t, "detail", project.Name)
	require.Equal(t, "Omer", members.Members[0].DisplayName)
	require.Equal(t, &endCursor, members.EndCursor)
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
	require.True(t, teams.HasNextPage)
	require.Equal(t, "LIT", teams.Teams[0].Key)
	require.Equal(t, "team body", team.Description)
	require.Equal(t, "Omer", teamMembers.Members[0].DisplayName)
	require.Equal(t, &endCursor, teamMembers.EndCursor)
	require.True(t, users.HasNextPage)
	require.True(t, users.Users[0].Admin)
	require.Equal(t, "Omer", user.DisplayName)
	require.Equal(t, "Omer", viewerUser.DisplayName)
	require.True(t, workflowStates.HasNextPage)
	require.Equal(t, &endCursor, workflowStates.EndCursor)
	require.Equal(t, "Started", workflowStates.WorkflowStates[0].Name)
	require.Equal(t, "LIT", workflowStates.WorkflowStates[0].TeamKey)
	require.Equal(t, "started", workflowState.Type)
	require.Equal(t, "linctl", workflowState.TeamName)
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
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"IssueByID": `{"issue":` + issueJSONWithDescription(issueFixture{
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
			"ProjectByID": `{"project":` + projectJSON(projectFixture{
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
		"IssueByID": `{"issue":` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-30",
			Title:      "assigned",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}, "Omer") + `}`,
		"ProjectByID": `{"project":` + projectJSONWithLead(projectFixture{
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
		"TeamByID": `{"team":{
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

		_, err = ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list projects")

		_, err = GetProjectByID(context.Background(), graphqlClient, "project-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project project-id")

		_, err = ListProjectMembers(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project members project-id")

		_, err = ListProjectUpdates(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project updates project-id")

		_, err = ListAllProjectUpdates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project updates")

		_, err = GetProjectUpdateByID(context.Background(), graphqlClient, "project-update-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project update project-update-id")

		_, err = ListProjectMilestones(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project milestones project-id")

		_, err = GetProjectMilestoneByID(context.Background(), graphqlClient, "project-milestone-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project milestone project-milestone-id")

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

		_, err = ListTeams(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list teams")

		_, err = GetTeamByID(context.Background(), graphqlClient, "team-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get team team-id")

		_, err = ListTeamMembers(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list team members team-id")

		_, err = ListUsers(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list users")

		_, err = GetUserByID(context.Background(), graphqlClient, "user-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get user user-id")

		_, err = GetViewerUser(context.Background(), graphqlClient)
		require.Error(t, err)
		require.Contains(t, err.Error(), "get viewer user")

		_, err = ListWorkflowStates(context.Background(), graphqlClient, 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list workflow states")

		_, err = GetWorkflowStateByID(context.Background(), graphqlClient, "workflow-state-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get workflow state workflow-state-id")
	})

	t.Run("issue mutations fail when payload omits entity", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"IssueCreate": `{"issueCreate":{"success":false,"issue":null}}`,
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"ProjectByID": `{"project":` + projectJSON(projectFixture{
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
				"ProjectByID":            `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
				"ProjectMilestoneCreate": `{"projectMilestoneCreate":{"success":false,"projectMilestone":null}}`,
			}),
			matchingTarget(),
			ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateProjectMilestone(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"ProjectMilestoneByID": `{"projectMilestone":` +
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
				"CycleByID":   `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
				"CycleUpdate": `{"cycleUpdate":{"success":false,"cycle":null}}`,
			}),
			matchingTarget(),
			CycleUpdateRequest{ID: "cycle-id", Name: "name"},
		)
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = ArchiveCycle(
			context.Background(),
			projectWriteFakeClient(map[string]string{
				"CycleByID":    `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
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
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"ProjectByID":   `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectUpdate": "",
		}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update project project-id")

		_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID":    `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectArchive": "",
		}).withError(operationErr), matchingTarget(), "project-id")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "archive project project-id")

		_, err = CreateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID":            `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectMilestoneCreate": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create project milestone")

		_, err = UpdateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectMilestoneByID": `{"projectMilestone":` +
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
			"CycleByID":   `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
			"CycleUpdate": "",
		}).withError(operationErr), matchingTarget(), CycleUpdateRequest{ID: "cycle-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update cycle cycle-id")

		_, err = ArchiveCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"CycleByID":    `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
			"CycleArchive": "",
		}).withError(operationErr), matchingTarget(), "cycle-id")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "archive cycle cycle-id")
	})

	t.Run("write operations return guard read errors", func(t *testing.T) {
		operationErr := errors.New("guard read failed")

		_, err := UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-50", Title: "title"})
		require.ErrorIs(t, err, operationErr)

		_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-50", Body: "body"})
		require.ErrorIs(t, err, operationErr)

		_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), "LIT-50")
		require.ErrorIs(t, err, operationErr)

		_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), "LIT-50")
		require.ErrorIs(t, err, operationErr)

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
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
			"ProjectByID": "",
		}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID": "",
		}).withError(operationErr), matchingTarget(), "project-id")
		require.ErrorIs(t, err, operationErr)

		_, err = CreateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = UpdateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectMilestoneByID": "",
		}).withError(operationErr), matchingTarget(), ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = UpdateCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"CycleByID": "",
		}).withError(operationErr), matchingTarget(), CycleUpdateRequest{ID: "cycle-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = ArchiveCycle(context.Background(), projectWriteFakeClient(map[string]string{
			"CycleByID": "",
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
		"IssueByID": `{"issue":` + strings.ReplaceAll(issueJSON(issueFixture{
			Identifier: "ABC-1",
			Title:      "wrong team",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
		"ProjectByID": `{"project":` + strings.ReplaceAll(projectJSON(projectFixture{
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
