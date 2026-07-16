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

func Test_collectNodePages_rejects_invalid_limit(t *testing.T) {
	_, err := collectNodePages(
		"list projects", 0, projectListPageSize,
		func(int, *string) (nodePage[ProjectSummary], error) {
			t.Fatal("fetch must not run")
			return nodePage[ProjectSummary]{}, nil
		},
	)

	require.ErrorContains(t, err, "list projects: limit must be positive")
}

func Test_collectNodePages_rejects_missing_next_cursor(t *testing.T) {
	_, err := collectNodePages(
		"list projects", 100, projectListPageSize,
		func(int, *string) (nodePage[ProjectSummary], error) {
			return nodePage[ProjectSummary]{HasNextPage: true}, nil
		},
	)

	require.ErrorContains(t, err, "list projects: next page has no end cursor")
}
