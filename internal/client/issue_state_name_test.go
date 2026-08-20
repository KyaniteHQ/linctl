package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UpdateIssue_selects_exact_started_state_name(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	after := before
	after.StateID = "in-review-state"
	after.State = "In Review"
	after.StateType = "started"
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(before) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(after) + `}}`,
	})}

	issue, err := UpdateIssue(
		context.Background(),
		withIssueAfterWrite(recorder, after),
		matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"stateId": "in-review-state"}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_refuses_when_readback_state_does_not_match(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	wrong := before
	wrong.StateID = "in-progress-state"
	wrong.State = "In Progress"
	wrong.StateType = "started"
	graphqlClient := withIssueAfterWrite(issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(before) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(wrong) + `}}`,
	}), wrong)

	_, err := UpdateIssue(
		context.Background(), graphqlClient, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.ErrorIs(t, err, ErrStateMismatch)
	require.ErrorContains(t, err, "in-review-state")
}

func Test_UpdateIssue_skips_write_when_issue_already_has_selected_state(t *testing.T) {
	current := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "in-review-state", State: "In Review", StateType: "started",
	}
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(current) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
	})}

	issue, err := UpdateIssue(
		context.Background(), recorder, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
	require.Zero(t, recorder.countOf("IssueUpdate"))
}

func Test_UpdateIssue_reconciles_ambiguous_state_write_without_replay(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	after := before
	after.StateID = "in-review-state"
	after.State = "In Review"
	after.StateType = "started"
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(before) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          "",
	}).withError(errors.New("timeout"))}
	graphqlClient := withIssueAfterWrite(recorder, after)

	issue, err := UpdateIssue(
		context.Background(), graphqlClient, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
	require.Equal(t, 1, recorder.countOf("IssueUpdate"))
}

func Test_CreateIssue_reconciles_ambiguous_create_without_replay(t *testing.T) {
	after := issueFixture{
		ID: "issue-id", Identifier: "LIT-3", Title: "typed", ProjectID: "project-id", Project: "fixture",
		StateID: "in-review-state", State: "In Review", StateType: "started",
	}
	inner := issueWriteFakeClient(map[string]string{
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssuesByTeamFiltered": `{"issues":{"nodes":[` + issueJSON(after) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	})
	sequenced := newSequentialOpClient(inner)
	sequenced.failAt["IssueCreate"] = 1
	graphqlClient := withIssueAfterWrite(sequenced, after)

	issue, err := CreateIssue(
		context.Background(), graphqlClient, matchingTarget(),
		IssueCreateRequest{Title: "typed", StateSelector: "In Review"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
	require.Equal(t, 1, sequenced.calls["IssueCreate"])
}

func Test_CreateIssue_reconciles_without_a_state_selector(t *testing.T) {
	after := issueFixture{
		ID: "issue-id", Identifier: "LIT-3", Title: "typed", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	sequenced := newSequentialOpClient(issueWriteFakeClient(map[string]string{
		"IssuesByTeamFiltered": `{"issues":{"nodes":[` + issueJSON(after) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}))
	sequenced.failAt["IssueCreate"] = 1

	issue, err := CreateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueCreateRequest{Title: "typed"},
	)

	require.NoError(t, err)
	require.Equal(t, "LIT-3", issue.Identifier)
	require.Equal(t, 1, sequenced.calls["IssueCreate"])
}

func Test_CreateIssue_returns_mutation_error_when_payload_has_issue_without_success(t *testing.T) {
	after := issueFixture{
		ID: "issue-id", Identifier: "LIT-3", Title: "typed", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}

	_, err := CreateIssue(
		context.Background(),
		issueWriteFakeClient(map[string]string{
			"IssueCreate": `{"issueCreate":{"success":false,"issue":` + issueJSON(after) + `}}`,
		}),
		matchingTarget(),
		IssueCreateRequest{Title: "typed"},
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_CreateIssue_returns_write_error_when_created_issue_is_missing(t *testing.T) {
	sequenced := newSequentialOpClient(issueWriteFakeClient(map[string]string{
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssuesByTeamFiltered": `{"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}))
	sequenced.failAt["IssueCreate"] = 1

	_, err := CreateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueCreateRequest{Title: "typed", StateSelector: "In Review"},
	)

	require.ErrorContains(t, err, "injected IssueCreate failure")
}

func Test_CreateIssue_fails_closed_when_created_title_is_ambiguous(t *testing.T) {
	first := issueJSON(issueFixture{
		ID: "issue-a", Identifier: "LIT-3", Title: "typed", ProjectID: "project-id", Project: "fixture",
		StateID: "in-review-state", State: "In Review", StateType: "started",
	})
	second := issueJSON(issueFixture{
		ID: "issue-b", Identifier: "LIT-4", Title: "typed", ProjectID: "project-id", Project: "fixture",
		StateID: "in-review-state", State: "In Review", StateType: "started",
	})
	sequenced := newSequentialOpClient(issueWriteFakeClient(map[string]string{
		"IssuesByTeamFiltered": `{"issues":{"nodes":[` + first + `,` + second +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}))
	sequenced.failAt["IssueCreate"] = 1

	_, err := CreateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueCreateRequest{Title: "typed"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "ambiguous")
}

func Test_CreateIssue_fails_closed_when_created_title_page_is_truncated(t *testing.T) {
	sequenced := newSequentialOpClient(issueWriteFakeClient(map[string]string{
		"IssuesByTeamFiltered": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			ID: "issue-id", Identifier: "LIT-3", Title: "typed", ProjectID: "project-id", Project: "fixture",
			StateID: "todo-state", State: "Todo", StateType: "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"c"}}}`,
	}))
	sequenced.failAt["IssueCreate"] = 1

	_, err := CreateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueCreateRequest{Title: "typed"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "ambiguous")
}

func Test_CreateIssue_uses_state_selector_over_type(t *testing.T) {
	after := issueFixture{
		Identifier: "LIT-3", Title: "typed", ProjectID: "project-id", Project: "fixture",
		StateID: "in-review-state", State: "In Review", StateType: "started",
	}
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueCreate":          `{"issueCreate":{"success":true,"issue":` + issueJSON(after) + `}}`,
	})}

	issue, err := CreateIssue(
		context.Background(),
		withIssueAfterWrite(recorder, after),
		matchingTarget(),
		IssueCreateRequest{Title: "typed", StateSelector: "In Review", StateType: "started"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
}

func Test_UpdateIssue_returns_readback_error_when_issue_vanishes(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	inner := issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(before) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(before) + `}}`,
	})
	sequenced := newSequentialOpClient(inner)
	sequenced.failAt["issue"] = 2

	_, err := UpdateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.Error(t, err)
	require.ErrorContains(t, err, "injected issue failure")
}

func Test_UpdateIssue_prefers_mutation_error_when_readback_also_fails(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	inner := issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(before) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          "",
	}).withError(errors.New("timeout"))
	sequenced := newSequentialOpClient(inner)
	sequenced.failAt["issue"] = 2

	_, err := UpdateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.ErrorContains(t, err, "timeout")
}

func Test_resolveStateID_requires_selector(t *testing.T) {
	guard := &guardedClient{stateIDs: &stateIDCache{lists: map[string][]workflowStateCandidate{}}}

	_, err := guard.resolveStateID(context.Background(), "team-id", "  ")

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_resolveStateID_returns_error_when_team_states_cannot_be_listed(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}

	_, err := UpdateIssue(
		context.Background(),
		issueWriteFakeClient(map[string]string{
			"issue": `{"issue":` + issueJSON(before) + `}`,
		}),
		matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.ErrorContains(t, err, "list workflow states")
}

func Test_StartIssue_selects_lowest_position_started_type_not_name_started(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-5", Title: "start", ProjectID: "project-id", Project: "fixture",
			StateID: "todo-state", State: "Todo", StateType: "unstarted",
		}) + `}`,
		"WorkflowStatesByTeam": workflowStatesByTeamJSON(`
			{"id":"started-state","name":"Started","type":"started","position":2},
			{"id":"in-progress-state","name":"In Progress","type":"started","position":0}
		`),
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-5", Title: "start", ProjectID: "project-id", Project: "fixture",
			StateID: "in-progress-state", State: "In Progress", StateType: "started",
		}, "Ada") + `}}`,
	})}
	graphqlClient := withIssueAfterWriteJSON(recorder, issueJSONWithAssignee(issueFixture{
		Identifier: "LIT-5", Title: "start", ProjectID: "project-id", Project: "fixture",
		StateID: "in-progress-state", State: "In Progress", StateType: "started",
	}, "Ada"))

	issue, err := StartIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-5")

	require.NoError(t, err)
	require.Equal(t, "in-progress-state", issue.StateID)
	require.JSONEq(t, `{
		"id": "LIT-5",
		"input": {"assigneeId": "user-id", "stateId": "in-progress-state"}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_selectWorkflowStateID_returns_error_when_name_is_ambiguous(t *testing.T) {
	states := []workflowStateCandidate{
		{ID: "a", Name: "In Review", Type: "started", Position: 1},
		{ID: "b", Name: "in review", Type: "started", Position: 2},
	}

	_, err := selectWorkflowStateID(states, "In Review")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "ambiguous")
}
