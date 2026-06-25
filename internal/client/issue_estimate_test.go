package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func teamEstimateConfigJSON(estimationType string, allowZero bool) string {
	zero := "false"
	if allowZero {
		zero = "true"
	}

	return `{"team":{"id":"team-id","issueEstimationType":"` + estimationType + `","issueEstimationAllowZero":` + zero + `}}`
}

func Test_CreateIssue_sets_estimate_when_team_scale_accepts_it(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"teamEstimateConfig": teamEstimateConfigJSON("fibonacci", false),
		"IssueCreate":        `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-7")) + `}}`,
	})

	estimate := 3
	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b3",
		Estimate: &estimate,
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-7", issue.Identifier)
}

func Test_CreateIssue_accepts_zero_estimate_when_team_allows_zero(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"teamEstimateConfig": teamEstimateConfigJSON("fibonacci", true),
		"IssueCreate":        `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-7")) + `}}`,
	})

	estimate := 0
	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b3",
		Estimate: &estimate,
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-7", issue.Identifier)
}

func Test_CreateIssue_rejects_estimate_when_team_estimates_disabled(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"teamEstimateConfig": teamEstimateConfigJSON("notUsed", false),
	})

	estimate := 3
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b3",
		Estimate: &estimate,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssue_rejects_zero_estimate_when_team_disallows_zero(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"teamEstimateConfig": teamEstimateConfigJSON("fibonacci", false),
	})

	estimate := 0
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b3",
		Estimate: &estimate,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssue_rejects_negative_estimate(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	estimate := -1
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b3",
		Estimate: &estimate,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssue_wraps_team_estimate_config_read_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	estimate := 3
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "b3",
		Estimate: &estimate,
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrWriteInvalid)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateIssue_sets_estimate_when_team_scale_accepts_it(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":              `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"teamEstimateConfig": teamEstimateConfigJSON("fibonacci", false),
		"IssueUpdate":        `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})

	estimate := 5
	issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:       "LIT-1",
		Estimate: &estimate,
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_UpdateIssue_clears_estimate_when_requested(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})

	issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:            "LIT-1",
		ClearEstimate: true,
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_UpdateIssue_rejects_estimate_with_clear_estimate(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	estimate := 3
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:            "LIT-1",
		Estimate:      &estimate,
		ClearEstimate: true,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_rejects_estimate_when_team_estimates_disabled(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":              `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"teamEstimateConfig": teamEstimateConfigJSON("notUsed", false),
	})

	estimate := 3
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:       "LIT-1",
		Estimate: &estimate,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}
