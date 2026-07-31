package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func initiativeLabelWithOrgJSON(orgID string) string {
	return `{
		"id":"initiative-label-id",
		"name":"Strategy",
		"description":"Strategic theme",
		"color":"#5e6ad2",
		"isGroup":false,
		"lastAppliedAt":"2026-07-10T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-07-01T12:00:00Z",
		"updatedAt":"2026-07-10T12:00:00Z",
		"organization":{"id":"` + orgID + `"},
		"parent":null
	}`
}

func initiativeLabelWithParentJSON(orgID string) string {
	return `{
		"id":"initiative-label-id",
		"name":"Strategy",
		"description":"Strategic theme",
		"color":"#5e6ad2",
		"isGroup":false,
		"lastAppliedAt":"2026-07-10T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-07-01T12:00:00Z",
		"updatedAt":"2026-07-10T12:00:00Z",
		"organization":{"id":"` + orgID + `"},
		"parent":{"id":"initiative-label-group-id","name":"Themes","color":"#8a8f98"}
	}`
}

func Test_RetireInitiativeLabel_retires_label_when_organization_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
		"InitiativeLabelRetire": `{"initiativeLabelRetire":{"success":true,"initiativeLabel":` +
			initiativeLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := RetireInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.NoError(t, err)
	require.Equal(t, "initiative-label-id", label.ID)
	require.Equal(t, "org-id", label.OrgID)
}

func Test_RetireInitiativeLabel_requires_id(t *testing.T) {
	_, err := RetireInitiativeLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "", true,
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RetireInitiativeLabel_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := RetireInitiativeLabel(context.Background(), recorder, matchingTarget(), "initiative-label-id", false)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("InitiativeLabelRetire"))
}

func Test_RetireInitiativeLabel_refuses_when_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := RetireInitiativeLabel(context.Background(), recorder, matchingTarget(), "initiative-label-id", true)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("InitiativeLabelRetire"))
}

func Test_RetireInitiativeLabel_refuses_when_target_unresolved(t *testing.T) {
	_, err := RetireInitiativeLabel(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, "initiative-label-id", true)

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_RetireInitiativeLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := RetireInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RetireInitiativeLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
		"InitiativeLabelRetire": `{"initiativeLabelRetire":{"success":false,"initiativeLabel":` +
			initiativeLabelWithOrgJSON("org-id") + `}}`,
	})

	_, err := RetireInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RetireInitiativeLabel_proceeds_when_pinned_project_present(t *testing.T) {
	// matchingTarget() pins project-id; InitiativeLabel taxonomy writes compare
	// organization only, so a pinned project must not block this write.
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
		"InitiativeLabelRetire": `{"initiativeLabelRetire":{"success":true,"initiativeLabel":` +
			initiativeLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := RetireInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.NoError(t, err)
	require.Equal(t, "initiative-label-id", label.ID)
}

func Test_RestoreInitiativeLabel_restores_label_when_organization_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
		"InitiativeLabelRestore": `{"initiativeLabelRestore":{"success":true,"initiativeLabel":` +
			initiativeLabelWithOrgJSON("org-id") + `}}`,
	})

	label, err := RestoreInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.NoError(t, err)
	require.Equal(t, "initiative-label-id", label.ID)
}

func Test_RestoreInitiativeLabel_requires_id(t *testing.T) {
	_, err := RestoreInitiativeLabel(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "", true,
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RestoreInitiativeLabel_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := RestoreInitiativeLabel(context.Background(), recorder, matchingTarget(), "initiative-label-id", false)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("InitiativeLabelRestore"))
}

func Test_RestoreInitiativeLabel_refuses_when_organization_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("other-org") + `}`,
	})}

	_, err := RestoreInitiativeLabel(context.Background(), recorder, matchingTarget(), "initiative-label-id", true)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("InitiativeLabelRestore"))
}

func Test_RestoreInitiativeLabel_refuses_when_target_unresolved(t *testing.T) {
	_, err := RestoreInitiativeLabel(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, "initiative-label-id", true)

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_RestoreInitiativeLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
	})

	_, err := RestoreInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RestoreInitiativeLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiativeLabel": `{"initiativeLabel":` + initiativeLabelWithOrgJSON("org-id") + `}`,
		"InitiativeLabelRestore": `{"initiativeLabelRestore":{"success":false,"initiativeLabel":` +
			initiativeLabelWithOrgJSON("org-id") + `}}`,
	})

	_, err := RestoreInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RestoreInitiativeLabel_wraps_resolve_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{})

	_, err := RestoreInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}

func Test_RetireInitiativeLabel_wraps_resolve_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{})

	_, err := RetireInitiativeLabel(context.Background(), graphqlClient, matchingTarget(), "initiative-label-id", true)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}
