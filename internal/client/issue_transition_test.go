package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func transitionTarget(allowlist config.Transitions) config.Target {
	target := matchingTarget()
	target.Transitions = allowlist

	return target
}

func transitionPayloads() map[string]string {
	return map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
			StateID: "todo-state", State: "Todo", StateType: "unstarted",
		}) + `}`,
		"WorkflowStatesByTeam": workflowStatesByTeamJSON(`
			{"id":"todo-state","name":"Todo","type":"unstarted","position":0},
			{"id":"in-progress-state","name":"In Progress","type":"started","position":1},
			{"id":"done-state","name":"Done","type":"completed","position":2}
		`),
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
			StateID: "in-progress-state", State: "In Progress", StateType: "started",
		}) + `}}`,
	}
}

func Test_UpdateIssue_moves_state_when_the_transition_is_listed(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(transitionPayloads())}
	graphqlClient := withIssueAfterWrite(recorder, issueFixture{
		Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
		StateID: "in-progress-state", State: "In Progress", StateType: "started",
	})

	issue, err := UpdateIssue(context.Background(), graphqlClient,
		transitionTarget(config.Transitions{"todo": {"in progress"}}),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Progress"})

	require.NoError(t, err)
	require.Equal(t, "in-progress-state", issue.StateID)
	require.JSONEq(t, `{"id": "LIT-1", "input": {"stateId": "in-progress-state"}}`,
		string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_refuses_a_transition_the_allowlist_does_not_list(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(transitionPayloads())}

	_, err := UpdateIssue(context.Background(), recorder,
		transitionTarget(config.Transitions{"Todo": {"Canceled"}, "In Progress": {"Done"}}),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Progress"})

	require.ErrorIs(t, err, ErrTransitionDenied)
	require.ErrorContains(t, err, `"Todo" -> "In Progress"`)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_UpdateIssue_lets_a_non_state_edit_through_the_allowlist(t *testing.T) {
	// A title edit is not a transition, so a strict allowlist must not block it.
	payloads := transitionPayloads()
	payloads["IssueUpdate"] = `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
		Identifier: "LIT-1", Title: "renamed", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}) + `}}`
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(payloads)}

	issue, err := UpdateIssue(context.Background(), recorder,
		transitionTarget(config.Transitions{"Done": {"Todo"}}),
		IssueUpdateRequest{ID: "LIT-1", Title: "renamed"})

	require.NoError(t, err)
	require.Equal(t, "renamed", issue.Title)
	require.True(t, recorder.sentOperation("IssueUpdate"))
}

func Test_UpdateIssue_lets_a_same_state_edit_through_the_allowlist(t *testing.T) {
	// Re-selecting the current state alongside another field is not a transition.
	payloads := transitionPayloads()
	payloads["IssueUpdate"] = `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
		Identifier: "LIT-1", Title: "renamed", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}) + `}}`
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(payloads)}

	_, err := UpdateIssue(context.Background(), recorder,
		transitionTarget(config.Transitions{"Done": {"Todo"}}),
		IssueUpdateRequest{ID: "LIT-1", Title: "renamed", StateSelector: "Todo"})

	require.NoError(t, err)
	require.True(t, recorder.sentOperation("IssueUpdate"))
}

func Test_StartIssue_refuses_when_the_allowlist_omits_the_started_state(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(transitionPayloads())}

	_, err := StartIssue(context.Background(), recorder,
		transitionTarget(config.Transitions{"Todo": {"Done"}}), "LIT-1")

	require.ErrorIs(t, err, ErrTransitionDenied)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_CloseIssue_refuses_when_the_allowlist_omits_the_completed_state(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(transitionPayloads())}

	_, err := CloseIssue(context.Background(), recorder,
		transitionTarget(config.Transitions{"Todo": {"In Progress"}}), "LIT-1")

	require.ErrorIs(t, err, ErrTransitionDenied)
	require.False(t, recorder.sentOperation("IssueClose"))
}

func Test_CloseIssue_moves_state_when_the_transition_is_listed(t *testing.T) {
	payloads := transitionPayloads()
	payloads["IssueClose"] = `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
		Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
		StateID: "done-state", State: "Done", StateType: "completed",
	}) + `}}`
	graphqlClient := withIssueAfterWrite(issueWriteFakeClient(payloads), issueFixture{
		Identifier: "LIT-1", Title: "existing", ProjectID: "project-id", Project: "fixture",
		StateID: "done-state", State: "Done", StateType: "completed",
	})

	issue, err := CloseIssue(context.Background(), graphqlClient,
		transitionTarget(config.Transitions{"Todo": {"Done"}}), "LIT-1")

	require.NoError(t, err)
	require.Equal(t, "done-state", issue.StateID)
}

func Test_requireTransition_refuses_a_state_missing_from_the_cache(t *testing.T) {
	// The selected state always comes from the same cached list, so a miss is a
	// programming error surfaced as an invalid write rather than a nil deref.
	guard := &guardedClient{
		target:   ResolvedTarget{Expected: transitionTarget(config.Transitions{"Todo": {"Done"}})},
		stateIDs: &stateIDCache{lists: map[string][]workflowStateCandidate{}},
	}

	err := guard.requireTransition(IssueSummary{TeamID: "team-id", Team: "LIT", StateID: "todo-state", State: "Todo"}, "done-state")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "not a cached workflow state")
}
