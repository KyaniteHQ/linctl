package client

import (
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

func Test_CreateProject_returns_created_project_when_target_matches(t *testing.T) {
	// Given
	graphqlClient := projectWriteFakeClient(map[string]string{
		"ProjectCreate": `{"projectCreate":{"success":true,"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "created",
			Status: "Backlog",
		}) + `}}`,
	})

	// When
	project, err := CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{
		Name:        "created",
		Description: "body",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "project-id", project.ID)
	require.Equal(t, "created", project.Name)
	require.True(t, projectHasTeam(project, "team-id", "LIT"))
}

func Test_UpdateProject_refuses_when_pinned_project_differs(t *testing.T) {
	// Given
	graphqlClient := projectWriteFakeClient(map[string]string{
		"ProjectByID": `{"project":` + projectJSON(projectFixture{
			ID:     "other-project",
			Name:   "other",
			Status: "Backlog",
		}) + `}`,
	})

	// When
	_, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
		ID:   "other-project",
		Name: "updated",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

type projectWriteFakeClient map[string]string

func (client projectWriteFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	return fakeGraphQLClient(client.withTargetResponses()).MakeRequest(ctx, request, response)
}

func (client projectWriteFakeClient) withTargetResponses() map[string]string {
	responses := map[string]string{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
		"TargetProject": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [{
						"id": "team-id",
						"key": "LIT",
						"name": "linctl-it",
						"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
					}]
				}
			}
		}`,
	}
	for operation, response := range client {
		responses[operation] = response
	}

	return responses
}

type projectFixture struct {
	ID     string
	Name   string
	Status string
}

func projectJSON(project projectFixture) string {
	return `{
		"id":"` + project.ID + `",
		"name":"` + project.Name + `",
		"description":"description",
		"slugId":"` + project.Name + `",
		"url":"https://linear.app/kyanite/project/` + project.ID + `",
		"priority":0,
		"status":{"id":"status-id","name":"` + project.Status + `","type":"backlog"},
		"lead":null,
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl-it"}]}
	}`
}
