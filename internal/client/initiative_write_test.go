package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func initiativeJSON(orgID string) string {
	return `{
		"id": "created-initiative-id",
		"name": "Platform",
		"description": "platform",
		"status": "Planned",
		"priority": 0,
		"targetDate": null,
		"slugId": "platform",
		"url": "https://linear.app/kyanite/initiative/platform",
		"organization": {"id": "` + orgID + `"}
	}`
}

func initiativeCreatePayload(orgID string) map[string]string {
	return map[string]string{
		"InitiativeCreate": `{"initiativeCreate":{"success":true,"initiative":` + initiativeJSON(orgID) + `}}`,
	}
}

func Test_CreateInitiative_creates_initiative_when_org_wide_passed(t *testing.T) {
	initiative, err := CreateInitiative(
		context.Background(), projectWriteFakeClient(initiativeCreatePayload("org-id")), matchingTarget(),
		InitiativeCreateRequest{Name: "Platform", Description: "platform", Status: "Planned", OrgWide: true},
	)

	require.NoError(t, err)
	require.Equal(t, "created-initiative-id", initiative.ID)
	require.Equal(t, "Planned", initiative.Status)
	require.Equal(t, "org-id", initiative.OrgID)
}

func Test_CreateInitiative_refuses_without_org_wide(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateInitiative(context.Background(), recorder, matchingTarget(), InitiativeCreateRequest{
		Name: "Platform",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("InitiativeCreate"))
}

func Test_CreateInitiative_requires_name(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateInitiative(context.Background(), recorder, matchingTarget(), InitiativeCreateRequest{
		OrgWide: true,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("InitiativeCreate"))
}

func Test_CreateInitiative_rejects_an_unknown_status(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateInitiative(context.Background(), recorder, matchingTarget(), InitiativeCreateRequest{
		Name: "Platform", Status: "Started", OrgWide: true,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "Proposed, Planned, Active, Completed, Canceled")
	require.False(t, recorder.sentOperation("InitiativeCreate"))
}

func Test_CreateInitiative_refuses_when_target_unresolved(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := CreateInitiative(context.Background(), recorder, config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, InitiativeCreateRequest{Name: "Platform", OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("InitiativeCreate"))
}

func Test_CreateInitiative_wraps_mutation_error(t *testing.T) {
	_, err := CreateInitiative(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		InitiativeCreateRequest{Name: "Platform", OrgWide: true},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateInitiative_fails_when_mutation_reports_no_success(t *testing.T) {
	_, err := CreateInitiative(context.Background(), projectWriteFakeClient(map[string]string{
		"InitiativeCreate": `{"initiativeCreate":{"success":false,"initiative":` + initiativeJSON("org-id") + `}}`,
	}), matchingTarget(), InitiativeCreateRequest{Name: "Platform", OrgWide: true})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_CreateInitiative_refuses_an_initiative_created_in_another_organization(t *testing.T) {
	// The Org-Scoped Write comparison is made against what Linear returned, not
	// against target resolution alone, so an initiative landing elsewhere is a hard stop.
	_, err := CreateInitiative(
		context.Background(), projectWriteFakeClient(initiativeCreatePayload("other-org-id")), matchingTarget(),
		InitiativeCreateRequest{Name: "Platform", OrgWide: true},
	)

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateInitiative_sends_every_field_it_was_given(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(initiativeCreatePayload("org-id"))}

	_, err := CreateInitiative(context.Background(), recorder, matchingTarget(), InitiativeCreateRequest{
		Name: "Platform", Description: "platform", Status: "Planned", OrgWide: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"input": {
			"name": "Platform",
			"description": "platform",
			"status": "Planned"
		}
	}`, string(recorder.variablesFor(t, "InitiativeCreate")))
}

func Test_CreateInitiative_omits_the_optional_fields_it_was_not_given(t *testing.T) {
	// Linear picks the default status itself, and an explicit null would override
	// that default rather than defer to it, so an unset field must not be sent.
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(initiativeCreatePayload("org-id"))}

	_, err := CreateInitiative(context.Background(), recorder, matchingTarget(), InitiativeCreateRequest{
		Name: "Platform", OrgWide: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{"input": {"name": "Platform"}}`, string(recorder.variablesFor(t, "InitiativeCreate")))
}
