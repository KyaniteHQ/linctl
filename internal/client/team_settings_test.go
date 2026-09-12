package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func teamStateJSON(teamID string, teamKey string) string {
	return `{"workflowState":{"id":"triage-state","name":"Triage","type":"triage","color":"#bec2c8","position":0,` +
		`"team":{"id":"` + teamID + `","key":"` + teamKey + `","name":"` + teamKey + `"}}}`
}

func teamSettingsPayloads() map[string]string {
	return map[string]string{
		"team":               `{"team":` + otherTeamJSON("ops-team-id", "org-id") + `}`,
		"workflowState":      teamStateJSON("ops-team-id", "OPS"),
		"TeamSettingsUpdate": `{"teamUpdate":{"success":true,"team":` + otherTeamJSON("ops-team-id", "org-id") + `}}`,
	}
}

func Test_UpdateTeamSettings_sends_every_setting_it_was_given(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(teamSettingsPayloads())}

	team, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(), TeamSettingsRequest{
		ID: "ops-team-id", Triage: boolPtr(true), DefaultStateID: "triage-state", Inherit: boolPtr(false),
		OrgWide: true,
	})

	require.NoError(t, err)
	require.Equal(t, "ops-team-id", team.ID)
	require.JSONEq(t, `{"id": "ops-team-id", "input": {
		"triageEnabled": true, "defaultIssueStateId": "triage-state", "inheritWorkflowStatuses": false
	}}`, string(recorder.variablesFor(t, "TeamSettingsUpdate")))
}

func Test_UpdateTeamSettings_omits_the_settings_it_was_not_given(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(teamSettingsPayloads())}

	_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(), TeamSettingsRequest{
		ID: "ops-team-id", Triage: boolPtr(false), OrgWide: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{"id": "ops-team-id", "input": {"triageEnabled": false}}`,
		string(recorder.variablesFor(t, "TeamSettingsUpdate")))
}

func Test_UpdateTeamSettings_validates_before_any_request(t *testing.T) {
	tests := []struct {
		name    string
		request TeamSettingsRequest
		want    error
	}{
		{name: "missing id", request: TeamSettingsRequest{Triage: boolPtr(true), OrgWide: true}, want: ErrWriteInvalid},
		{name: "no setting", request: TeamSettingsRequest{ID: "ops-team-id", OrgWide: true}, want: ErrWriteInvalid},
		{name: "no org-wide", request: TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true)}, want: ErrTargetMismatch},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

			_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(), test.request)

			require.ErrorIs(t, err, test.want)
			require.False(t, recorder.sentOperation("TeamSettingsUpdate"))
		})
	}
}

func Test_UpdateTeamSettings_refuses_when_target_unresolved(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := UpdateTeamSettings(context.Background(), recorder, config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true), OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TeamSettingsUpdate"))
}

func Test_UpdateTeamSettings_wraps_team_lookup_error(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{})}

	_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true), OrgWide: true})

	require.ErrorContains(t, err, "get team ops-team-id")
	require.False(t, recorder.sentOperation("TeamSettingsUpdate"))
}

func Test_UpdateTeamSettings_refuses_a_team_in_another_organization(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"team": `{"team":` + otherTeamJSON("ops-team-id", "other-org-id") + `}`,
	})}

	_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true), OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TeamSettingsUpdate"))
}

func Test_UpdateTeamSettings_wraps_default_state_lookup_error(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"team": `{"team":` + otherTeamJSON("ops-team-id", "org-id") + `}`,
	})}

	_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", DefaultStateID: "triage-state", OrgWide: true})

	require.ErrorContains(t, err, "get workflow state triage-state")
	require.False(t, recorder.sentOperation("TeamSettingsUpdate"))
}

func Test_UpdateTeamSettings_refuses_a_default_state_from_another_team(t *testing.T) {
	payloads := teamSettingsPayloads()
	payloads["workflowState"] = teamStateJSON("team-id", "LIT")
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(payloads)}

	_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", DefaultStateID: "triage-state", OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorContains(t, err, "belongs to team LIT")
	require.False(t, recorder.sentOperation("TeamSettingsUpdate"))
}

func Test_UpdateTeamSettings_accepts_a_parent_state_for_a_sub_team(t *testing.T) {
	payloads := teamSettingsPayloads()
	payloads["team"] = `{"team":` + subTeamJSON("team-id", "org-id") + `}`
	payloads["workflowState"] = teamStateJSON("team-id", "LIT")
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(payloads)}

	_, err := UpdateTeamSettings(context.Background(), recorder, matchingTarget(),
		TeamSettingsRequest{ID: "created-team-id", DefaultStateID: "triage-state", OrgWide: true})

	require.NoError(t, err)
	require.True(t, recorder.sentOperation("TeamSettingsUpdate"))
}

func Test_UpdateTeamSettings_wraps_mutation_error(t *testing.T) {
	payloads := teamSettingsPayloads()
	delete(payloads, "TeamSettingsUpdate")

	_, err := UpdateTeamSettings(context.Background(), projectWriteFakeClient(payloads), matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true), OrgWide: true})

	require.ErrorContains(t, err, "update team settings ops-team-id")
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateTeamSettings_fails_when_mutation_returns_no_team(t *testing.T) {
	payloads := teamSettingsPayloads()
	payloads["TeamSettingsUpdate"] = `{"teamUpdate":{"success":true,"team":null}}`

	_, err := UpdateTeamSettings(context.Background(), projectWriteFakeClient(payloads), matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true), OrgWide: true})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UpdateTeamSettings_refuses_a_team_returned_in_another_organization(t *testing.T) {
	payloads := teamSettingsPayloads()
	payloads["TeamSettingsUpdate"] = `{"teamUpdate":{"success":true,"team":` + otherTeamJSON("ops-team-id", "other-org-id") + `}}`

	_, err := UpdateTeamSettings(context.Background(), projectWriteFakeClient(payloads), matchingTarget(),
		TeamSettingsRequest{ID: "ops-team-id", Triage: boolPtr(true), OrgWide: true})

	require.ErrorIs(t, err, ErrTargetMismatch)
}
