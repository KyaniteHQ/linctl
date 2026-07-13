package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CreateProjectLabel_creates_label_when_org_wide_passed(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"ProjectLabelCreate": `{"projectLabelCreate":{"success":true,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := CreateProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelCreateRequest{
		Name: "priority", OrgWide: true,
	})

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
}

func Test_CreateProjectLabel_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelCreateRequest{
		Name: "priority",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("ProjectLabelCreate"))
}

func Test_CreateProjectLabel_requires_name(t *testing.T) {
	_, err := CreateProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelCreateRequest{OrgWide: true},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateProjectLabel_refuses_when_target_unresolved(t *testing.T) {
	_, err := CreateProjectLabel(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, ProjectLabelCreateRequest{Name: "priority", OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateProjectLabel_wraps_mutation_error(t *testing.T) {
	_, err := CreateProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelCreateRequest{Name: "priority", OrgWide: true},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateProjectLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	_, err := CreateProjectLabel(context.Background(), projectWriteFakeClient(map[string]string{
		"ProjectLabelCreate": `{"projectLabelCreate":{"success":false,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	}), matchingTarget(), ProjectLabelCreateRequest{Name: "priority", OrgWide: true})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_CreateProjectLabel_proceeds_when_pinned_project_present(t *testing.T) {
	// matchingTarget() pins project-id; ProjectLabel taxonomy writes compare
	// organization only, so a pinned project must not block this write (103
	// decision memo ruling).
	graphqlClient := projectWriteFakeClient(map[string]string{
		"ProjectLabelCreate": `{"projectLabelCreate":{"success":true,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	_, err := CreateProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelCreateRequest{
		Name: "priority", OrgWide: true,
	})

	require.NoError(t, err)
}

func Test_UpdateProjectLabel_updates_label_when_organization_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectLabelUpdate": `{"projectLabelUpdate":{"success":true,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := UpdateProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelUpdateRequest{
		ID: "label-id", Name: "renamed", OrgWide: true,
	})

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
}

func Test_UpdateProjectLabel_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := UpdateProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("ProjectLabelUpdate"))
}

func Test_UpdateProjectLabel_refuses_when_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := UpdateProjectLabel(context.Background(), recorder, matchingTarget(), ProjectLabelUpdateRequest{
		ID: "label-id", Name: "renamed", OrgWide: true,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectLabelUpdate"))
}

func Test_UpdateProjectLabel_requires_id(t *testing.T) {
	_, err := UpdateProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelUpdateRequest{Name: "renamed", OrgWide: true},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateProjectLabel_requires_a_field(t *testing.T) {
	_, err := UpdateProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelUpdateRequest{ID: "label-id", OrgWide: true},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateProjectLabel_wraps_resolution_error(t *testing.T) {
	_, err := UpdateProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		ProjectLabelUpdateRequest{ID: "label-id", Name: "renamed", OrgWide: true},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateProjectLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := UpdateProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelUpdateRequest{
		ID: "label-id", Name: "renamed", OrgWide: true,
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateProjectLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectLabelUpdate": `{"projectLabelUpdate":{"success":false,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	_, err := UpdateProjectLabel(context.Background(), graphqlClient, matchingTarget(), ProjectLabelUpdateRequest{
		ID: "label-id", Name: "renamed", OrgWide: true,
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RetireProjectLabel_retires_label_when_organization_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectLabelRetire": `{"projectLabelRetire":{"success":true,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := RetireProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
}

func Test_RetireProjectLabel_requires_id(t *testing.T) {
	_, err := RetireProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "", true,
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RetireProjectLabel_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := RetireProjectLabel(context.Background(), recorder, matchingTarget(), "label-id", false)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("ProjectLabelRetire"))
}

func Test_RetireProjectLabel_refuses_when_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := RetireProjectLabel(context.Background(), recorder, matchingTarget(), "label-id", true)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectLabelRetire"))
}

func Test_RetireProjectLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := RetireProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RetireProjectLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectLabelRetire": `{"projectLabelRetire":{"success":false,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	_, err := RetireProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RestoreProjectLabel_restores_label_when_organization_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectLabelRestore": `{"projectLabelRestore":{"success":true,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := RestoreProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
}

func Test_RestoreProjectLabel_requires_id(t *testing.T) {
	_, err := RestoreProjectLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "", true,
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RestoreProjectLabel_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := RestoreProjectLabel(context.Background(), recorder, matchingTarget(), "label-id", false)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("ProjectLabelRestore"))
}

func Test_RestoreProjectLabel_refuses_when_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := RestoreProjectLabel(context.Background(), recorder, matchingTarget(), "label-id", true)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectLabelRestore"))
}

func Test_RestoreProjectLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := RestoreProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RestoreProjectLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectLabel": `{"projectLabel":` + projectLabelWithOrgJSON("org-id") + `}`,
		"ProjectLabelRestore": `{"projectLabelRestore":{"success":false,"projectLabel":` +
			projectLabelWithOrgJSON("org-id") + `}}`,
	})

	_, err := RestoreProjectLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.ErrorIs(t, err, ErrMutationFailed)
}
