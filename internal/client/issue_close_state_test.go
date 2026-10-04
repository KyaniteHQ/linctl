package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func closeStateTarget(name string) config.Target {
	target := matchingTarget()
	target.States = config.States{Close: name}

	return target
}

// closeStatePayloads is a team whose workflow has two completed states, with
// "Merged" at the lower position, the state a bare close would pick.
func closeStatePayloads() map[string]string {
	payloads := transitionPayloads()
	payloads["WorkflowStatesByTeam"] = workflowStatesByTeamJSON(`
		{"id":"todo-state","name":"Todo","type":"unstarted","position":0},
		{"id":"merged-state","name":"Merged","type":"completed","position":1},
		{"id":"done-state","name":"Done","type":"completed","position":2}
	`)
	payloads["IssueClose"] = `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
		Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
		StateID: "done-state", State: "Done", StateType: "completed",
	}) + `}}`

	return payloads
}

func Test_CloseIssue_moves_to_the_configured_state_over_a_lower_position_completed_state(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(closeStatePayloads())}
	graphqlClient := withIssueAfterWrite(recorder, issueFixture{
		Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
		StateID: "done-state", State: "Done", StateType: "completed",
	})

	issue, err := CloseIssue(context.Background(), graphqlClient, closeStateTarget("done"), "LIT-1")

	require.NoError(t, err)
	require.Equal(t, "done-state", issue.StateID)
	require.JSONEq(t, `{"id": "LIT-1", "input": {"stateId": "done-state"}}`,
		string(recorder.variablesFor(t, "IssueClose")))
}

func Test_CloseIssue_picks_the_lowest_position_completed_state_when_no_close_state_is_set(t *testing.T) {
	merged := issueFixture{
		Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
		StateID: "merged-state", State: "Merged", StateType: "completed",
	}
	payloads := closeStatePayloads()
	payloads["IssueClose"] = `{"issueUpdate":{"success":true,"issue":` + issueJSON(merged) + `}}`
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(payloads)}
	graphqlClient := withIssueAfterWrite(recorder, merged)

	_, err := CloseIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-1")

	require.NoError(t, err)
	require.JSONEq(t, `{"id": "LIT-1", "input": {"stateId": "merged-state"}}`,
		string(recorder.variablesFor(t, "IssueClose")))
}

func Test_CloseIssue_refuses_a_configured_state_missing_from_the_team(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(closeStatePayloads())}

	_, err := CloseIssue(context.Background(), recorder, closeStateTarget("Shipped"), "LIT-1")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "Shipped")
	require.False(t, recorder.sentOperation("IssueClose"))
}

func Test_CloseIssue_refuses_a_configured_state_that_is_not_completed(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(closeStatePayloads())}

	_, err := CloseIssue(context.Background(), recorder, closeStateTarget("Todo"), "LIT-1")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "unstarted")
	require.False(t, recorder.sentOperation("IssueClose"))
}

func Test_CloseIssue_applies_the_transitions_allowlist_to_the_configured_state(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(closeStatePayloads())}
	target := closeStateTarget("Done")
	target.Transitions = config.Transitions{"Todo": {"Merged"}}

	_, err := CloseIssue(context.Background(), recorder, target, "LIT-1")

	require.ErrorIs(t, err, ErrTransitionDenied)
	require.False(t, recorder.sentOperation("IssueClose"))
}

func Test_CloseIssue_wraps_a_failed_workflow_state_listing_for_a_configured_state(t *testing.T) {
	operationErr := errors.New("workflow states unavailable")
	payloads := closeStatePayloads()
	payloads["WorkflowStatesByTeam"] = ""
	graphqlClient := issueWriteFakeClient(payloads).withError(operationErr)

	_, err := CloseIssue(context.Background(), graphqlClient, closeStateTarget("Done"), "LIT-1")

	require.ErrorIs(t, err, operationErr)
	require.ErrorContains(t, err, "list workflow states")
}
