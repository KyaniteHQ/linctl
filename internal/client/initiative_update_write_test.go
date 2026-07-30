package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func initiativeWriteJSON(health string, body string) string {
	return `{
		"id":"initiative-update-id",
		"body":"` + body + `",
		"health":"` + health + `",
		"createdAt":"2026-06-20T00:00:00Z",
		"updatedAt":"2026-06-20T00:00:00Z",
		"url":"https://linear.app/kyanite/initiative-update/initiative-update-id",
		"slugId":"initiative-update-slug",
		"commentCount":0,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"user":{"id":"user-id","name":"Omer","displayName":"Omer"}
	}`
}

func initiativeOrgJSON(orgID string) string {
	return `{
		"id":"initiative-id",
		"name":"Platform",
		"description":"Platform initiative",
		"status":"Active",
		"priority":2,
		"targetDate":"2026-12-31",
		"slugId":"platform-init",
		"url":"https://linear.app/kyanite/initiative/platform-init",
		"organization":{"id":"` + orgID + `"}
	}`
}

func Test_CreateInitiativeUpdate_returns_summary_when_org_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"initiative": `{"initiative":` + initiativeOrgJSON("org-id") + `}`,
		"InitiativeUpdateCreate": `{"initiativeUpdateCreate":{"success":true,"initiativeUpdate":` +
			initiativeWriteJSON("onTrack", "All good") + `}}`,
	})}

	update, err := CreateInitiativeUpdate(context.Background(), recorder, matchingTarget(), InitiativeUpdateCreateRequest{
		InitiativeID: "initiative-id",
		Body:         "All good",
		Health:       "onTrack",
	})

	require.NoError(t, err)
	require.Equal(t, "initiative-update-id", update.ID)
	require.Equal(t, "onTrack", update.Health)
	require.Equal(t, "All good", update.Body)
	require.Equal(t, "initiative-id", update.InitiativeID)
	require.JSONEq(t, `{
		"input": {
			"initiativeId": "initiative-id",
			"body": "All good",
			"health": "onTrack"
		}
	}`, string(recorder.variablesFor(t, "InitiativeUpdateCreate")))
}

func Test_CreateInitiativeUpdate_refuses_wrong_org(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"initiative": `{"initiative":` + initiativeOrgJSON("other-org") + `}`,
	})}

	_, err := CreateInitiativeUpdate(context.Background(), recorder, matchingTarget(), InitiativeUpdateCreateRequest{
		InitiativeID: "initiative-id",
		Health:       "onTrack",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("InitiativeUpdateCreate"))
}

func Test_CreateInitiativeUpdate_requires_body_or_health(t *testing.T) {
	_, err := CreateInitiativeUpdate(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		InitiativeUpdateCreateRequest{InitiativeID: "initiative-id"},
	)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateInitiativeUpdate_requires_initiative_id(t *testing.T) {
	_, err := CreateInitiativeUpdate(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(),
		InitiativeUpdateCreateRequest{Health: "onTrack"},
	)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateInitiativeUpdate_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiative": `{"initiative":` + initiativeOrgJSON("org-id") + `}`,
		"InitiativeUpdateCreate": `{"initiativeUpdateCreate":{"success":false,"initiativeUpdate":` +
			initiativeWriteJSON("onTrack", "x") + `}}`,
	})

	_, err := CreateInitiativeUpdate(context.Background(), graphqlClient, matchingTarget(), InitiativeUpdateCreateRequest{
		InitiativeID: "initiative-id",
		Health:       "onTrack",
	})
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_CreateInitiativeUpdate_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"initiative": `{"initiative":` + initiativeOrgJSON("org-id") + `}`,
	})

	_, err := CreateInitiativeUpdate(context.Background(), graphqlClient, matchingTarget(), InitiativeUpdateCreateRequest{
		InitiativeID: "initiative-id",
		Health:       "onTrack",
	})
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateInitiativeUpdate_refuses_when_target_unresolved(t *testing.T) {
	_, err := CreateInitiativeUpdate(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, InitiativeUpdateCreateRequest{InitiativeID: "initiative-id", Health: "onTrack"})
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateInitiativeUpdate_wraps_initiative_read_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{})
	_, err := CreateInitiativeUpdate(context.Background(), graphqlClient, matchingTarget(), InitiativeUpdateCreateRequest{
		InitiativeID: "initiative-id",
		Health:       "onTrack",
	})
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}
