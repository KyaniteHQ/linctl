package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// issueFilterFakeClient serves one final page so each filter case records
// exactly one IssuesByTeamFiltered request.
func issueFilterFakeClient() fakeGraphQLClient {
	return fakeGraphQLClient{
		"IssuesByTeamFiltered": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-10",
			Title:      "listed issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}
}

func Test_ListIssuesByTeam_composes_combined_filters_into_one_operation(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueFilterFakeClient()}

	issues, err := ListIssuesByTeam(context.Background(), recorder, "team-id", 2, IssueListFilters{
		StateType:  "started",
		AssigneeID: "user-id",
	})

	require.NoError(t, err)
	require.NotEmpty(t, issues.Issues)
	require.JSONEq(t, `{
		"filter": {
			"team": {"id": {"eq": "team-id"}},
			"state": {"type": {"eq": "started"}},
			"assignee": {"id": {"eq": "user-id"}}
		},
		"first": 2,
		"after": null,
		"includeArchived": true
	}`, string(recorder.variablesFor(t, "IssuesByTeamFiltered")))
}

func Test_ListIssuesByTeam_composes_each_filter_clause(t *testing.T) {
	team := `"team": {"id": {"eq": "team-id"}}`
	tests := []struct {
		name    string
		filters IssueListFilters
		want    string
	}{
		{"unfiltered", IssueListFilters{}, `{` + team + `}`},
		{"state", IssueListFilters{StateType: "started"}, `{` + team + `, "state": {"type": {"eq": "started"}}}`},
		{"project", IssueListFilters{ProjectID: "project-id"}, `{` + team + `, "project": {"id": {"eq": "project-id"}}}`},
		{"assignee", IssueListFilters{AssigneeID: "user-id"}, `{` + team + `, "assignee": {"id": {"eq": "user-id"}}}`},
		{"label", IssueListFilters{LabelID: "label-id"}, `{` + team + `, "labels": {"some": {"id": {"eq": "label-id"}}}}`},
		{"cycle", IssueListFilters{CycleID: "cycle-id"}, `{` + team + `, "cycle": {"id": {"eq": "cycle-id"}}}`},
		{
			"created range",
			IssueListFilters{CreatedAfter: "2026-06-01", CreatedBefore: "2026-06-30"},
			`{` + team + `, "createdAt": {"gte": "2026-06-01", "lte": "2026-06-30"}}`,
		},
		{
			"updated range",
			IssueListFilters{UpdatedAfter: "2026-07-01", UpdatedBefore: "2026-07-30"},
			`{` + team + `, "updatedAt": {"gte": "2026-07-01", "lte": "2026-07-30"}}`,
		},
		{
			"updated after only",
			IssueListFilters{UpdatedAfter: "2026-07-01"},
			`{` + team + `, "updatedAt": {"gte": "2026-07-01"}}`,
		},
		{
			"state plus updated after",
			IssueListFilters{StateType: "started", UpdatedAfter: "2026-07-01"},
			`{` + team + `, "state": {"type": {"eq": "started"}}, "updatedAt": {"gte": "2026-07-01"}}`,
		},
		{"has blockers", IssueListFilters{HasBlockers: true}, `{` + team + `, "hasBlockedByRelations": {"eq": true}}`},
		{"blocks", IssueListFilters{Blocks: true}, `{` + team + `, "hasBlockingRelations": {"eq": true}}`},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := &recordingGraphQLClient{inner: issueFilterFakeClient()}

			_, err := ListIssuesByTeam(context.Background(), recorder, "team-id", 2, testCase.filters)

			require.NoError(t, err)
			require.JSONEq(t, `{
				"filter": `+testCase.want+`,
				"first": 2,
				"after": null,
				"includeArchived": true
			}`, string(recorder.variablesFor(t, "IssuesByTeamFiltered")))
		})
	}
}
