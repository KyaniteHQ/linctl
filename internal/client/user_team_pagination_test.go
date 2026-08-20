package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListUsers_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"users": `{"users":{"nodes":[` +
			userPageJSON("user-1", "omer") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"user-cursor-1"}}}`,
		"users:user-cursor-1": `{"users":{"nodes":[` +
			userPageJSON("user-2", "ada") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	users, err := ListUsers(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"user-1", "user-2"}, []string{
		users.Users[0].ID,
		users.Users[1].ID,
	})
	require.False(t, users.HasNextPage)
}

func Test_ListUserAssignedIssues_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"user_assignedIssues": `{"user":{"assignedIssues":{"nodes":[` +
			issueJSON(issueFixture{
				Identifier: "LIT-41",
				Title:      "First assigned",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"assigned-cursor-1"}}}}`,
		"user_assignedIssues:assigned-cursor-1": `{"user":{"assignedIssues":{"nodes":[` +
			issueJSON(issueFixture{
				Identifier: "LIT-42",
				Title:      "Second assigned",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	issues, err := ListUserAssignedIssues(context.Background(), graphqlClient, "user-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"LIT-41", "LIT-42"}, []string{
		issues.Issues[0].Identifier,
		issues.Issues[1].Identifier,
	})
	require.False(t, issues.HasNextPage)
}

func Test_ListViewerDrafts_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"viewer_drafts": `{"viewer":{"drafts":{"nodes":[` +
			draftPageJSON("draft-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"draft-cursor-1"}}}}`,
		"viewer_drafts:draft-cursor-1": `{"viewer":{"drafts":{"nodes":[` +
			draftPageJSON("draft-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	drafts, err := ListViewerDrafts(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"draft-1", "draft-2"}, []string{
		drafts.Drafts[0].ID,
		drafts.Drafts[1].ID,
	})
	require.False(t, drafts.HasNextPage)
}

func Test_ListTeamMembers_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"team_members": `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[` +
			userPageJSON("user-1", "omer") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"member-cursor-1"}}}}`,
		"team_members:member-cursor-1": `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[` +
			userPageJSON("user-2", "ada") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	members, err := ListTeamMembers(context.Background(), graphqlClient, "team-id", 100)

	require.NoError(t, err)
	require.Equal(t, "team-id", members.TeamID)
	require.Equal(t, "LIT", members.TeamKey)
	require.Equal(t, "linctl", members.TeamName)
	require.Equal(t, []string{"user-1", "user-2"}, []string{
		members.Members[0].ID,
		members.Members[1].ID,
	})
	require.False(t, members.HasNextPage)
}

func Test_ListTeamMemberships_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"teamMemberships": `{"teamMemberships":{"nodes":[` +
			teamMembershipPageJSON("membership-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"membership-cursor-1"}}}`,
		"teamMemberships:membership-cursor-1": `{"teamMemberships":{"nodes":[` +
			teamMembershipPageJSON("membership-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	memberships, err := ListTeamMemberships(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"membership-1", "membership-2"}, []string{
		memberships.Memberships[0].ID,
		memberships.Memberships[1].ID,
	})
	require.False(t, memberships.HasNextPage)
}

func Test_ListTeamGitAutomationStates_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"team_gitAutomationStates": `{"team":{"id":"team-id","key":"LIT","name":"linctl","gitAutomationStates":{"nodes":[` +
			gitAutomationStatePageJSON("git-state-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"git-cursor-1"}}}}`,
		"team_gitAutomationStates:git-cursor-1": `{"team":{"id":"team-id","key":"LIT","name":"linctl","gitAutomationStates":{"nodes":[` +
			gitAutomationStatePageJSON("git-state-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	states, err := ListTeamGitAutomationStates(context.Background(), graphqlClient, "team-id", 100)

	require.NoError(t, err)
	require.Equal(t, "team-id", states.TeamID)
	require.Equal(t, "LIT", states.TeamKey)
	require.Equal(t, "linctl", states.TeamName)
	require.Equal(t, []string{"git-state-1", "git-state-2"}, []string{
		states.States[0].ID,
		states.States[1].ID,
	})
	require.False(t, states.HasNextPage)
}

func userPageJSON(id string, name string) string {
	return `{"id":"` + id + `","name":"` + name + `","displayName":"` + name +
		`","email":"` + name + `@example.com","active":true,"guest":false,"admin":false}`
}

func draftPageJSON(id string) string {
	return `{"id":"` + id + `","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,` +
		`"issue":{"id":"issue-id","identifier":"LIT-3","title":"Draft issue"},` +
		`"project":null,"projectUpdate":null,"initiative":null,"initiativeUpdate":null,` +
		`"parentComment":null,"customerNeed":null,"team":null}`
}

func teamMembershipPageJSON(id string) string {
	return `{"id":"` + id + `","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,` +
		`"owner":true,"sortOrder":1.5,` +
		`"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},` +
		`"team":{"id":"team-id","key":"LIT","name":"linctl"}}`
}

func gitAutomationStatePageJSON(id string) string {
	return `{"id":"` + id + `","event":"review","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,` +
		`"state":{"id":"workflow-state-id","name":"Started","type":"started"},` +
		`"targetBranch":{"id":"target-branch-id","branchPattern":"main","isRegex":false}}`
}
