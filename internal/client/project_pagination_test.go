package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListProjectsByTeam_pages_in_API_safe_chunks(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{
			operation: "Projects",
			variables: `{"teamId":"team-id","first":50,"after":null,"includeArchived":true}`,
			payload: `{"team":{"projects":{"nodes":[` + projectJSON(projectFixture{
				ID: "project-1", Name: "First", Status: "Backlog",
			}) + `],"pageInfo":{
				"hasNextPage":true,"endCursor":"project-cursor-1"
			}}}}`,
		},
		{
			operation: "Projects",
			variables: `{"teamId":"team-id","first":50,"after":"project-cursor-1","includeArchived":true}`,
			payload: `{"team":{"projects":{"nodes":[` + projectJSON(projectFixture{
				ID: "project-2", Name: "Second", Status: "Backlog",
			}) + `],"pageInfo":{
				"hasNextPage":false,"endCursor":null
			}}}}`,
		},
	}}

	projects, err := ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"project-1", "project-2"}, []string{
		projects.Projects[0].ID,
		projects.Projects[1].ID,
	})
	require.False(t, projects.HasNextPage)
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_ListProjects_pages_in_API_safe_chunks(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{fixtures: []strictTargetFixture{
		{
			operation: "projects",
			variables: `{"first":50,"after":null,"includeArchived":true}`,
			payload: `{"projects":{"nodes":[` + projectJSON(projectFixture{
				ID: "project-1", Name: "First", Status: "Backlog",
			}) + `],"pageInfo":{
				"hasNextPage":true,"endCursor":"project-cursor-1"
			}}}`,
		},
		{
			operation: "projects",
			variables: `{"first":50,"after":"project-cursor-1","includeArchived":true}`,
			payload: `{"projects":{"nodes":[` + projectJSON(projectFixture{
				ID: "project-2", Name: "Second", Status: "Backlog",
			}) + `],"pageInfo":{
				"hasNextPage":false,"endCursor":null
			}}}`,
		},
	}}

	projects, err := ListProjects(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"project-1", "project-2"}, []string{
		projects.Projects[0].ID,
		projects.Projects[1].ID,
	})
	require.False(t, projects.HasNextPage)
	require.Equal(t, len(graphqlClient.fixtures), graphqlClient.next)
}

func Test_ListProjects_reports_ErrGraphQL_when_next_page_has_no_end_cursor(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"projects": `{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID: "project-1", Name: "First", Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":null}}}`,
	}

	_, err := ListProjects(context.Background(), graphqlClient, 100)

	require.ErrorIs(t, err, ErrGraphQL)
	require.ErrorContains(t, err, "list projects: next page has no end cursor")
}

func Test_ListProjects_rejects_non_positive_limit_before_any_request(t *testing.T) {
	graphqlClient := &strictTargetGraphQLClient{}

	_, err := ListProjects(context.Background(), graphqlClient, 0)

	require.ErrorContains(t, err, "list projects: limit must be positive")
	require.Zero(t, graphqlClient.next)
}
