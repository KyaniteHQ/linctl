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

	stateID, err := selectWorkflowStateID(states, "team-id", "In Review")

	require.NoError(t, err)
	require.Equal(t, "in-review-state", stateID)
}

func Test_selectWorkflowStateID_matches_state_name_without_regard_to_case(t *testing.T) {
	states := []workflowStateCandidate{
		{ID: "in-review-state", Name: "In Review", Type: "started", Position: 1},
		{ID: "in-progress-state", Name: "In Progress", Type: "started", Position: 0},
	}

	stateID, err := selectWorkflowStateID(states, "team-id", "in review")

	require.NoError(t, err)
	require.Equal(t, "in-review-state", stateID)
}

func Test_selectWorkflowStateID_uses_lowest_position_for_a_type(t *testing.T) {
	states := []workflowStateCandidate{
		{ID: "second-state", Name: "Second", Type: "started", Position: 2},
		{ID: "first-state", Name: "First", Type: "started", Position: 1},
	}

	stateID, err := selectWorkflowStateID(states, "team-id", "started")

	require.NoError(t, err)
	require.Equal(t, "first-state", stateID)
}

func Test_selectWorkflowStateID_returns_error_when_type_has_no_states(t *testing.T) {
	_, err := selectWorkflowStateID(nil, "team-id", "started")

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_selectWorkflowStateID_returns_error_for_unknown_selector(t *testing.T) {
	_, err := selectWorkflowStateID(nil, "team-id", "sprinting")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "unknown workflow state")
}

func Test_listTeamWorkflowStates_returns_error_on_graphql_failure(t *testing.T) {
	_, err := listTeamWorkflowStates(context.Background(), fakeGraphQLClient(map[string]string{}), "team-id")

	require.ErrorContains(t, err, "list workflow states")
}

func Test_listTeamWorkflowStates_fails_closed_when_page_is_truncated(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"WorkflowStatesByTeam": `{"workflowStates":{"nodes":[
			{"id":"started-state","name":"Started","type":"started","position":1}
		],"pageInfo":{"hasNextPage":true}}}`,
	})

	_, err := listTeamWorkflowStates(context.Background(), graphqlClient, "team-id")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "more than 50 workflow states")
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
