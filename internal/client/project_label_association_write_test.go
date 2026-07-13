package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func projectLabelWithOrgJSON(orgID string) string {
	return `{
		"id":"label-id",
		"name":"priority",
		"description":"",
		"color":"#000000",
		"isGroup":false,
		"lastAppliedAt":null,
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-06-20T00:00:00Z",
		"updatedAt":"2026-06-20T00:00:00Z",
		"organization":{"id":"` + orgID + `"},
		"parent":null
	}`
}

func Test_AddProjectLabel_attaches_label_when_project_and_org_match(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project":      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectAddLabel": `{"projectAddLabel":{"success":true,"project":` +
			projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}}`,
	})

	project, err := AddProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "project-id", project.ID)
}

func Test_AddProjectLabel_refuses_when_label_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"project":      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := AddProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectAddLabel"))
}

func Test_AddProjectLabel_refuses_when_project_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "other-project", Name: "other", Status: "Backlog"}) + `}`,
	})}

	_, err := AddProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "other-project", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectAddLabel"))
}

func Test_AddProjectLabel_requires_project_and_label_ids(t *testing.T) {
	_, err := AddProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelAssociationRequest{ProjectID: "project-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectLabel_refuses_when_target_unresolved(t *testing.T) {
	_, err := AddProjectLabel(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, ProjectLabelAssociationRequest{ProjectID: "project-id", LabelID: "label-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectLabel_wraps_label_resolution_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
	})

	_, err := AddProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project":      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := AddProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project":         `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel":    `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectAddLabel": `{"projectAddLabel":{"success":false,"project":null}}`,
	})

	_, err := AddProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RemoveProjectLabel_detaches_label_when_project_and_org_match(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project":      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectRemoveLabel": `{"projectRemoveLabel":{"success":true,"project":` +
			projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}}`,
	})

	project, err := RemoveProjectLabel(
		context.Background(), graphqlClient, matchingTarget(),
		ProjectLabelAssociationRequest{ProjectID: "project-id", LabelID: "label-id"},
	)

	require.NoError(t, err)
	require.Equal(t, "project-id", project.ID)
}

func Test_RemoveProjectLabel_wraps_project_resolution_error(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := RemoveProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectRemoveLabel"))
}

func Test_RemoveProjectLabel_requires_project_and_label_ids(t *testing.T) {
	_, err := RemoveProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelAssociationRequest{LabelID: "label-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RemoveProjectLabel_refuses_when_label_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"project":      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := RemoveProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectRemoveLabel"))
}

func Test_RemoveProjectLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project":      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := RemoveProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RemoveProjectLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"project":            `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"projectLabel":       `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectRemoveLabel": `{"projectRemoveLabel":{"success":false,"project":null}}`,
	})

	_, err := RemoveProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelAssociationRequest{
		ProjectID: "project-id", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}
