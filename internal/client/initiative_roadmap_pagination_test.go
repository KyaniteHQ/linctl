package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListInitiatives_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"initiatives": `{"initiatives":{"nodes":[` +
			initiativePageJSON("initiative-1", "First") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"initiative-cursor-1"}}}`,
		"initiatives:initiative-cursor-1": `{"initiatives":{"nodes":[` +
			initiativePageJSON("initiative-2", "Second") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	initiatives, err := ListInitiatives(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"initiative-1", "initiative-2"}, []string{
		initiatives.Initiatives[0].ID,
		initiatives.Initiatives[1].ID,
	})
	require.False(t, initiatives.HasNextPage)
}

func Test_ListInitiativeHistory_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"initiative_history": `{"initiative":{"history":{"nodes":[` +
			initiativeHistoryPageJSON("history-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"history-cursor-1"}}}}`,
		"initiative_history:history-cursor-1": `{"initiative":{"history":{"nodes":[` +
			initiativeHistoryPageJSON("history-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	history, err := ListInitiativeHistory(context.Background(), graphqlClient, "initiative-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"history-1", "history-2"}, []string{
		history.History[0].ID,
		history.History[1].ID,
	})
	require.False(t, history.HasNextPage)
}

func Test_ListInitiativeProjects_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"initiative_projects": `{"initiative":{"projects":{"nodes":[` +
			projectJSON(projectFixture{ID: "project-1", Name: "First", Status: "Started"}) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"project-cursor-1"}}}}`,
		"initiative_projects:project-cursor-1": `{"initiative":{"projects":{"nodes":[` +
			projectJSON(projectFixture{ID: "project-2", Name: "Second", Status: "Started"}) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	projects, err := ListInitiativeProjects(context.Background(), graphqlClient, "initiative-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"project-1", "project-2"}, []string{
		projects.Projects[0].ID,
		projects.Projects[1].ID,
	})
	require.False(t, projects.HasNextPage)
}

func Test_ListInitiativeToProjects_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"initiativeToProjects": `{"initiativeToProjects":{"nodes":[` +
			initiativeToProjectPageJSON("association-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"association-cursor-1"}}}`,
		"initiativeToProjects:association-cursor-1": `{"initiativeToProjects":{"nodes":[` +
			initiativeToProjectPageJSON("association-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	associations, err := ListInitiativeToProjects(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"association-1", "association-2"}, []string{
		associations.Associations[0].ID,
		associations.Associations[1].ID,
	})
	require.False(t, associations.HasNextPage)
}

func Test_ListInitiativeUpdateComments_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"initiativeUpdate_comments": `{"initiativeUpdate":{"id":"initiative-update-id","comments":{"nodes":[` +
			commentMetadataJSONWithID("comment-1", "", "", "", "user-id") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"comment-cursor-1"}}}}`,
		"initiativeUpdate_comments:comment-cursor-1": `{"initiativeUpdate":{"id":"initiative-update-id","comments":{"nodes":[` +
			commentMetadataJSONWithID("comment-2", "", "", "", "user-id") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	comments, err := ListInitiativeUpdateComments(context.Background(), graphqlClient, "initiative-update-id", 100)

	require.NoError(t, err)
	require.Equal(t, "initiative-update-id", comments.InitiativeUpdateID)
	require.Equal(t, []string{"comment-1", "comment-2"}, []string{
		comments.Comments[0].ID,
		comments.Comments[1].ID,
	})
	require.False(t, comments.HasNextPage)
}

func Test_ListRoadmaps_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"roadmaps": `{"roadmaps":{"nodes":[` +
			roadmapPageJSON("roadmap-1", "First") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"roadmap-cursor-1"}}}`,
		"roadmaps:roadmap-cursor-1": `{"roadmaps":{"nodes":[` +
			roadmapPageJSON("roadmap-2", "Second") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	roadmaps, err := ListRoadmaps(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"roadmap-1", "roadmap-2"}, []string{
		roadmaps.Roadmaps[0].ID,
		roadmaps.Roadmaps[1].ID,
	})
	require.False(t, roadmaps.HasNextPage)
}

func Test_ListRoadmapProjects_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"roadmap_projects": `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","projects":{"nodes":[` +
			projectJSON(projectFixture{ID: "project-1", Name: "First", Status: "Backlog"}) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"project-cursor-1"}}}}`,
		"roadmap_projects:project-cursor-1": `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","projects":{"nodes":[` +
			projectJSON(projectFixture{ID: "project-2", Name: "Second", Status: "Backlog"}) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	projects, err := ListRoadmapProjects(context.Background(), graphqlClient, "roadmap-id", 100)

	require.NoError(t, err)
	require.Equal(t, "roadmap-id", projects.RoadmapID)
	require.Equal(t, "Platform roadmap", projects.RoadmapName)
	require.Equal(t, []string{"project-1", "project-2"}, []string{
		projects.Projects[0].ID,
		projects.Projects[1].ID,
	})
	require.False(t, projects.HasNextPage)
}

func initiativePageJSON(id string, name string) string {
	return `{"id":"` + id + `","name":"` + name +
		`","description":"Platform initiative","status":"Active","priority":2,` +
		`"targetDate":"2026-12-31","slugId":"` + id +
		`","url":"https://linear.app/kyanite/initiative/` + id +
		`","organization":{"id":"org-id"}}`
}

func initiativeHistoryPageJSON(id string) string {
	return `{"id":"` + id + `","createdAt":"2026-06-03T12:00:00Z",` +
		`"updatedAt":"2026-06-03T12:01:00Z","archivedAt":null,` +
		`"entries":[{"type":"status","from":"Planned","to":"Active"}],` +
		`"initiative":{"id":"initiative-id"}}`
}

func initiativeToProjectPageJSON(id string) string {
	return `{"id":"` + id + `","sortOrder":"1","createdAt":"2026-06-19T12:00:00Z",` +
		`"updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,` +
		`"initiative":{"id":"initiative-id","name":"Platform"},` +
		`"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project",` +
		`"url":"https://linear.app/project/project-id"}}`
}

func roadmapPageJSON(id string, name string) string {
	return `{"id":"` + id + `","name":"` + name +
		`","description":"Roadmap body","color":"#5e6ad2","slugId":"` + id +
		`","sortOrder":1,"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z",` +
		`"updatedAt":"2026-06-19T12:01:00Z",` +
		`"url":"https://linear.app/kyanite/roadmap/` + id +
		`","creator":{"id":"user-id","displayName":"Omer"},` +
		`"owner":{"id":"owner-id","displayName":"Owner"}}`
}
