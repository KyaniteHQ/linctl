package client

import (
	"context"
	"strconv"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

func Test_CreateProject_returns_created_project_when_target_matches(t *testing.T) {
	// Given
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"ProjectCreate": `{"projectCreate":{"success":true,"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "created",
			Status: "Backlog",
		}) + `}}`,
	})}

	// When
	project, err := CreateProject(context.Background(), recorder, matchingTarget(), ProjectCreateRequest{
		Name:        "created",
		Description: "body",
		Content:     "# heading",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "project-id", project.ID)
	require.Equal(t, "created", project.Name)
	require.True(t, projectHasTeam(project, "team-id", "LIT"))
	require.JSONEq(t, `{
		"input": {
			"name": "created",
			"description": "body",
			"content": "# heading",
			"teamIds": ["team-id"]
		}
	}`, string(recorder.variablesFor(t, "ProjectCreate")))
}

func Test_UpdateProject_refuses_when_pinned_project_differs(t *testing.T) {
	// Given
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
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

func Test_UpdateProject_refuses_when_project_lacks_pinned_team(t *testing.T) {
	// Given: the resolved project id matches the pinned project, but the project
	// belongs to a different team. This isolates requireProject's projectHasTeam
	// branch (write_guard.go) from the project-id mismatch branch above it.
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeam(projectFixture{
			ID:     "project-id",
			Name:   "fixture",
			Status: "Backlog",
		}, "other-team", "OTHER") + `}`,
	})

	// When
	_, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
		ID:   "project-id",
		Name: "updated",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

// Test_UpdateProject_succeeds_when_pinned_team_matches_on_truncated_page proves
// a truncated teams page (hasNextPage=true) does not block a legitimate write
// when the pinned team is found within the fetched page.
func Test_UpdateProject_succeeds_when_pinned_team_matches_on_truncated_page(t *testing.T) {
	// Given
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeamPage(projectFixture{
			ID:     "project-id",
			Name:   "fixture",
			Status: "Backlog",
		}, "team-id", "LIT", true) + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "updated",
			Status: "Backlog",
		}) + `}}`,
	})

	// When
	project, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
		ID:   "project-id",
		Name: "updated",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "updated", project.Name)
}

// Test_UpdateProject_refuses_fail_closed_when_teams_page_is_truncated proves an
// unmatched pinned team on a truncated teams page fails closed with a distinct
// message instead of silently trusting the first 50 teams.
func Test_UpdateProject_refuses_fail_closed_when_teams_page_is_truncated(t *testing.T) {
	// Given
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeamPage(projectFixture{
			ID:     "project-id",
			Name:   "fixture",
			Status: "Backlog",
		}, "other-team", "OTHER", true) + `}`,
	})

	// When
	_, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
		ID:   "project-id",
		Name: "updated",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Contains(t, err.Error(), "more than 50 teams")
}

func Test_UpdateProject_returns_updated_project_when_target_matches(t *testing.T) {
	// Given
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "fixture",
			Status: "Backlog",
		}) + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "updated",
			Status: "Backlog",
		}) + `}}`,
	})}

	// When
	project, err := UpdateProject(context.Background(), recorder, matchingTarget(), ProjectUpdateRequest{
		ID:      "project-id",
		Name:    "updated",
		Content: "# updated content",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "updated", project.Name)
	require.JSONEq(t, `{
		"id": "project-id",
		"input": {
			"name": "updated",
			"content": "# updated content"
		}
	}`, string(recorder.variablesFor(t, "ProjectUpdate")))
}

// Test_UpdateProject_succeeds_without_pinned_project_when_team_matches is the
// success-path counterpart to Test_UpdateProject_refuses_when_pinned_project_differs:
// with no project pinned, requireProject skips the project-id comparison but
// still enforces the pinned team (Test_UpdateProject_refuses_when_project_lacks_pinned_team).
func Test_UpdateProject_succeeds_without_pinned_project_when_team_matches(t *testing.T) {
	// Given
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID:     "any-project",
			Name:   "fixture",
			Status: "Backlog",
		}) + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
			ID:     "any-project",
			Name:   "updated",
			Status: "Backlog",
		}) + `}}`,
	})

	// When
	project, err := UpdateProject(context.Background(), graphqlClient, teamOnlyTarget(), ProjectUpdateRequest{
		ID:   "any-project",
		Name: "updated",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "updated", project.Name)
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
	return projectJSONWithTeam(project, "team-id", "LIT")
}

func projectJSONWithTeam(project projectFixture, teamID string, teamKey string) string {
	return projectJSONWithTeamPage(project, teamID, teamKey, false)
}

func projectJSONWithTeamPage(project projectFixture, teamID string, teamKey string, teamsHasNextPage bool) string {
	return `{
		"id":"` + project.ID + `",
		"name":"` + project.Name + `",
		"description":"description",
		"slugId":"` + project.Name + `",
		"url":"https://linear.app/kyanite/project/` + project.ID + `",
		"priority":0,
		"status":{"id":"status-id","name":"` + project.Status + `","type":"backlog"},
		"lead":null,
		"teams":{
			"nodes":[{"id":"` + teamID + `","key":"` + teamKey + `","name":"linctl-it"}],
			"pageInfo":{"hasNextPage":` + strconv.FormatBool(teamsHasNextPage) + `}
		}
	}`
}
