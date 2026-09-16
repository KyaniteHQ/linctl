package client

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UpdateIssue_sets_state_and_delegate_in_one_write(t *testing.T) {
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
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review", DelegateID: "agent-user-id"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"stateId": "in-review-state", "delegateId": "agent-user-id"}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_sets_delegate_alone(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:         "LIT-1",
		DelegateID: "agent-user-id",
	})

	require.NoError(t, err)
	require.Equal(t, 1, recorder.countOf("IssueUpdate"))
}

func Test_UpdateIssue_clears_delegate_when_requested(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:            "LIT-1",
		ClearDelegate: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"delegateId": null}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_rejects_delegate_with_clear_delegate(t *testing.T) {
	_, err := UpdateIssue(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", DelegateID: "agent-user-id", ClearDelegate: true},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_issueSummaryFromFields_carries_delegate_and_updated_at(t *testing.T) {
	body := strings.Replace(
		issueJSON(b1IssueFixture("LIT-1")), `"assignee":null,`,
		`"assignee":null,"delegate":{"id":"agent-user-id"},"updatedAt":"2026-09-16T14:26:12.000Z",`, 1,
	)
	graphqlClient := issueWriteFakeClient(map[string]string{"issue": `{"issue":` + body + `}`})

	detail, err := GetIssueDetail(context.Background(), graphqlClient, "LIT-1")

	require.NoError(t, err)
	require.Equal(t, "agent-user-id", detail.Summary.DelegateID)
	require.Equal(t, "2026-09-16T14:26:12.000Z", detail.Summary.UpdatedAt)
}

func Test_UpdateIssue_trusts_the_mutation_state_when_the_read_trails_the_write(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	after := before
	after.StateID = "in-review-state"
	after.State = "In Review"
	after.StateType = "started"
	sequenced := newSequentialOpClient(issueWriteFakeClient(map[string]string{
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(after) + `}}`,
	}))
	sequenced.payloads["issue"] = []string{
		`{"issue":` + issueJSON(before) + `}`,
		`{"issue":` + issueJSON(before) + `}`,
	}

	issue, err := UpdateIssue(
		context.Background(), sequenced, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.NoError(t, err)
	require.Equal(t, "in-review-state", issue.StateID)
	require.Equal(t, 2, sequenced.calls["issue"])
}

func Test_UpdateIssue_reports_a_mismatch_when_the_mutation_state_differs(t *testing.T) {
	before := issueFixture{
		Identifier: "LIT-1", Title: "job", ProjectID: "project-id", Project: "fixture",
		StateID: "todo-state", State: "Todo", StateType: "unstarted",
	}
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":                `{"issue":` + issueJSON(before) + `}`,
		"WorkflowStatesByTeam": multipleStartedStatesJSON(),
		"IssueUpdate":          `{"issueUpdate":{"success":true,"issue":` + issueJSON(before) + `}}`,
	})

	_, err := UpdateIssue(
		context.Background(), graphqlClient, matchingTarget(),
		IssueUpdateRequest{ID: "LIT-1", StateSelector: "In Review"},
	)

	require.ErrorIs(t, err, ErrStateMismatch)
}
