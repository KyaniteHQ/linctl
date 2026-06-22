package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func projectUpdateJSON(health string, body string) string {
	return `{
		"id":"project-update-id",
		"body":"` + body + `",
		"health":"` + health + `",
		"createdAt":"2026-06-20T00:00:00Z",
		"updatedAt":"2026-06-20T00:00:00Z",
		"url":"https://linear.app/kyanite/project-update/project-update-id",
		"project":{"id":"project-id","name":"fixture"},
		"user":{"id":"user-id","name":"Omer","displayName":"Omer"}
	}`
}

func Test_CreateProjectUpdate_returns_summary_when_target_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"ProjectUpdateCreate": `{"projectUpdateCreate":{"success":true,"projectUpdate":` +
			projectUpdateJSON("onTrack", "All good") + `}}`,
	})

	update, err := CreateProjectUpdate(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateCreateRequest{
		ProjectID: "project-id",
		Body:      "All good",
		Health:    "onTrack",
	})

	require.NoError(t, err)
	require.Equal(t, "project-update-id", update.ID)
	require.Equal(t, "onTrack", update.Health)
	require.Equal(t, "All good", update.Body)
	require.Equal(t, "project-id", update.ProjectID)
}

func Test_CreateProjectUpdate_refuses_when_pinned_project_differs(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "other-project", Name: "other", Status: "Backlog"}) + `}`,
	})

	_, err := CreateProjectUpdate(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateCreateRequest{
		ProjectID: "other-project",
		Health:    "onTrack",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateProjectUpdate_refuses_when_target_unresolved(t *testing.T) {
	_, err := CreateProjectUpdate(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, ProjectUpdateCreateRequest{ProjectID: "project-id", Health: "onTrack"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateProjectUpdate_requires_project_id(t *testing.T) {
	_, err := CreateProjectUpdate(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectUpdateCreateRequest{Health: "onTrack"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateProjectUpdate_requires_body_or_health(t *testing.T) {
	_, err := CreateProjectUpdate(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectUpdateCreateRequest{ProjectID: "project-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateProjectUpdate_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
	})

	_, err := CreateProjectUpdate(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateCreateRequest{
		ProjectID: "project-id",
		Health:    "onTrack",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateProjectUpdate_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"ProjectUpdateCreate": `{"projectUpdateCreate":{"success":false,"projectUpdate":` +
			projectUpdateJSON("onTrack", "x") + `}}`,
	})

	_, err := CreateProjectUpdate(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateCreateRequest{
		ProjectID: "project-id",
		Health:    "onTrack",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}
