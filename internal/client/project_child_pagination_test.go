package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListProjectMembers_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"project_members": `{"project":{"id":"project-id","name":"detail","members":{"nodes":[` +
			projectMemberJSON("user-1", "omer") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"member-cursor-1"}}}}`,
		"project_members:member-cursor-1": `{"project":{"id":"project-id","name":"detail","members":{"nodes":[` +
			projectMemberJSON("user-2", "ada") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	members, err := ListProjectMembers(context.Background(), graphqlClient, "project-id", 100)

	require.NoError(t, err)
	require.Equal(t, "project-id", members.ProjectID)
	require.Equal(t, "detail", members.ProjectName)
	require.Equal(t, []string{"user-1", "user-2"}, []string{
		members.Members[0].ID,
		members.Members[1].ID,
	})
	require.False(t, members.HasNextPage)
}

func Test_ListProjectLabels_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"projectLabels": `{"projectLabels":{"nodes":[` +
			projectLabelJSON("label-1", "First") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"label-cursor-1"}}}`,
		"projectLabels:label-cursor-1": `{"projectLabels":{"nodes":[` +
			projectLabelJSON("label-2", "Second") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	labels, err := ListProjectLabels(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"label-1", "label-2"}, []string{
		labels.ProjectLabels[0].ID,
		labels.ProjectLabels[1].ID,
	})
	require.False(t, labels.HasNextPage)
}

func Test_ListProjectLabelChildren_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"projectLabel_children": `{"projectLabel":{"id":"project-label-id","name":"Roadmap","children":{"nodes":[` +
			projectLabelJSON("child-1", "Mobile") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"child-cursor-1"}}}}`,
		"projectLabel_children:child-cursor-1": `{"projectLabel":{"id":"project-label-id","name":"Roadmap","children":{"nodes":[` +
			projectLabelJSON("child-2", "Web") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	children, err := ListProjectLabelChildren(context.Background(), graphqlClient, "project-label-id", 100)

	require.NoError(t, err)
	require.Equal(t, "project-label-id", children.ProjectLabelID)
	require.Equal(t, "Roadmap", children.ProjectLabelName)
	require.Equal(t, []string{"child-1", "child-2"}, []string{
		children.ProjectLabels[0].ID,
		children.ProjectLabels[1].ID,
	})
	require.False(t, children.HasNextPage)
}

func Test_ListProjectMilestones_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"project_projectMilestones": `{"project":{"id":"project-id","name":"detail","projectMilestones":{"nodes":[` +
			projectMilestonePageJSON("milestone-1", "First") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"milestone-cursor-1"}}}}`,
		"project_projectMilestones:milestone-cursor-1": `{"project":{"id":"project-id","name":"detail","projectMilestones":{"nodes":[` +
			projectMilestonePageJSON("milestone-2", "Second") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	milestones, err := ListProjectMilestones(context.Background(), graphqlClient, "project-id", 100)

	require.NoError(t, err)
	require.Equal(t, "project-id", milestones.ProjectID)
	require.Equal(t, "detail", milestones.ProjectName)
	require.Equal(t, []string{"milestone-1", "milestone-2"}, []string{
		milestones.Milestones[0].ID,
		milestones.Milestones[1].ID,
	})
	require.False(t, milestones.HasNextPage)
}

func Test_ListProjectUpdates_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"project_projectUpdates": `{"project":{"id":"project-id","name":"detail","projectUpdates":{"nodes":[` +
			projectUpdatePageJSON("update-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"update-cursor-1"}}}}`,
		"project_projectUpdates:update-cursor-1": `{"project":{"id":"project-id","name":"detail","projectUpdates":{"nodes":[` +
			projectUpdatePageJSON("update-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	updates, err := ListProjectUpdates(context.Background(), graphqlClient, "project-id", 100)

	require.NoError(t, err)
	require.Equal(t, "project-id", updates.ProjectID)
	require.Equal(t, "detail", updates.ProjectName)
	require.Equal(t, []string{"update-1", "update-2"}, []string{
		updates.Updates[0].ID,
		updates.Updates[1].ID,
	})
	require.False(t, updates.HasNextPage)
}

func Test_ListProjectRelations_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"projectRelations": `{"projectRelations":{"nodes":[` +
			projectRelationJSONWithID("relation-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"relation-cursor-1"}}}`,
		"projectRelations:relation-cursor-1": `{"projectRelations":{"nodes":[` +
			projectRelationJSONWithID("relation-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	relations, err := ListProjectRelations(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"relation-1", "relation-2"}, []string{
		relations.Relations[0].ID,
		relations.Relations[1].ID,
	})
	require.False(t, relations.HasNextPage)
}

func Test_ListProjectStatuses_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"projectStatuses": `{"projectStatuses":{"nodes":[` +
			projectStatusPageJSON("status-1", "Backlog") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"status-cursor-1"}}}`,
		"projectStatuses:status-cursor-1": `{"projectStatuses":{"nodes":[` +
			projectStatusPageJSON("status-2", "Started") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	statuses, err := ListProjectStatuses(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"status-1", "status-2"}, []string{
		statuses.ProjectStatuses[0].ID,
		statuses.ProjectStatuses[1].ID,
	})
	require.False(t, statuses.HasNextPage)
}

func projectMemberJSON(id string, name string) string {
	return `{"id":"` + id + `","name":"` + name + `","displayName":"` + name + `","email":"` + name + `@example.com"}`
}

func projectMilestonePageJSON(id string, name string) string {
	return `{"id":"` + id + `","name":"` + name +
		`","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}`
}

func projectUpdatePageJSON(id string) string {
	return `{"id":"` + id +
		`","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z",` +
		`"url":"https://linear.app/project-update/` + id +
		`","user":{"id":"user-id","name":"omer","displayName":"Omer"}}`
}

func projectRelationJSONWithID(id string) string {
	return `{"id":"` + id + `","type":"blocks","anchorType":"project","relatedAnchorType":"project",` +
		`"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,` +
		`"project":{"id":"project-id","name":"Pinned project"},"projectMilestone":null,` +
		`"relatedProject":{"id":"related-project-id","name":"Related project"},` +
		`"relatedProjectMilestone":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}`
}

func projectStatusPageJSON(id string, name string) string {
	return `{"id":"` + id + `","name":"` + name +
		`","description":"Ready for planning","type":"backlog","color":"#bec2c8","position":1,` +
		`"archivedAt":null,"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z"}`
}
