package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UpdateIssue_sets_milestone_when_target_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})}

	issue, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:                 "LIT-1",
		ProjectMilestoneID: "project-milestone-id",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"projectMilestoneId": "project-milestone-id"}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_clears_milestone_when_requested(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})}

	issue, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:             "LIT-1",
		ClearMilestone: true,
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"projectMilestoneId": null}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_rejects_milestone_with_clear_milestone(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:                 "LIT-1",
		ProjectMilestoneID: "project-milestone-id",
		ClearMilestone:     true,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_refuses_milestone_from_a_different_project(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Wrong project milestone", "next", "other-project") + `}`,
	})

	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:                 "LIT-1",
		ProjectMilestoneID: "project-milestone-id",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}
