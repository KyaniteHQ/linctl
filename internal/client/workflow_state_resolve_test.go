package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_selectWorkflowStateID_uses_exact_name_among_started_states(t *testing.T) {
	states := []workflowStateCandidate{
		{ID: "in-progress-state", Name: "In Progress", Type: "started", Position: 0},
		{ID: "in-review-state", Name: "In Review", Type: "started", Position: 1},
		{ID: "started-state", Name: "Started", Type: "started", Position: 2},
	}

	stateID, err := selectWorkflowStateID(states, "In Review")

	require.NoError(t, err)
	require.Equal(t, "in-review-state", stateID)
}

func Test_UpdateIssue_matches_state_name_without_regard_to_case(t *testing.T) {
	after := b1IssueFixture("LIT-1")
	after.StateID = "in-review-state"
	after.State = "In Review"
	after.StateType = "started"
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(after) + `}}`,
	})}

	_, err := UpdateIssue(context.Background(), withIssueAfterWrite(recorder, after), matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "in review"})

	require.NoError(t, err)
	require.JSONEq(t, `{"id": "LIT-1", "input": {"stateId": "in-review-state"}}`,
		string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_state_selector_does_not_fall_back_to_type(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"WorkflowStatesByTeam": workflowStatesByTeamJSON(`
			{"id":"in-review-state","name":"In Review","type":"started","position":1},
			{"id":"in-progress-state","name":"In Progress","type":"started","position":0}
		`),
	})}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "started"})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, `unknown workflow state "started"`)
	require.Zero(t, recorder.countOf("IssueUpdate"))
}

func Test_listTeamWorkflowStates_returns_error_on_graphql_failure(t *testing.T) {
	_, err := listTeamWorkflowStates(context.Background(), fakeGraphQLClient(map[string]string{}), "team-id")

	require.ErrorContains(t, err, "list workflow states")
}

func Test_CreateIssue_fails_closed_when_workflow_states_page_is_truncated(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStatesByTeam": `{"workflowStates":{"nodes":[
			{"id":"started-state","name":"Started","type":"started","position":1}
		],"pageInfo":{"hasNextPage":true}}}`,
	})}

	_, err := CreateIssue(context.Background(), recorder, matchingTarget(),
		IssueCreateRequest{Title: "typed", StateType: "started"})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "more than 50 workflow states")
	require.Zero(t, recorder.countOf("IssueCreate"))
}

func Test_IssueStateWriteRetryClass_is_reconcile_before_retry(t *testing.T) {
	require.Equal(t, MutationRetryReconcile, IssueStateWriteRetryClass())
}

func Test_IssueRelationCreateRetryClass_is_reconcile_before_retry(t *testing.T) {
	require.Equal(t, MutationRetryReconcile, IssueRelationCreateRetryClass())
}

func Test_CanonicalWorkflowStateType_maps_aliases(t *testing.T) {
	canonical, ok := CanonicalWorkflowStateType("in progress")

	require.True(t, ok)
	require.Equal(t, "started", canonical)
}
