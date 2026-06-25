package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateIssue_creates_sub_issue_when_parent_is_in_target(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-7")) + `}}`,
	})

	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b4",
		ParentID: "LIT-1",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-7", issue.Identifier)
}

func Test_CreateIssue_refuses_parent_in_a_different_team(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueReadWrongTeam(),
	})

	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b4",
		ParentID: "LIT-1",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssue_wraps_parent_read_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b4",
		ParentID: "LIT-1",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}
