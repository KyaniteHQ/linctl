package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListCyclesByTeam_returns_cycle_page(t *testing.T) {
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"cycles": `{"cycles":{"nodes":[{"id":"cycle-id","number":12,"name":null,"description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"}},{"id":"named-cycle-id","number":13,"name":"Named cycle","description":null,"startsAt":"2026-01-01T00:00:00Z","endsAt":"2026-02-01T00:00:00Z","completedAt":"2026-01-30T00:00:00Z","progress":1,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
	}

	cycles, err := ListCyclesByTeam(context.Background(), graphqlClient, "team-id", 1)

	require.NoError(t, err)
	require.True(t, cycles.HasNextPage)
	require.Equal(t, &endCursor, cycles.EndCursor)
	require.Equal(t, "cycle-id", cycles.Cycles[0].ID)
	require.Equal(t, "Cycle 12", cycles.Cycles[0].Name)
	require.Equal(t, "cycle body", cycles.Cycles[0].Description)
	require.Equal(t, "active", cycles.Cycles[0].Status)
	require.Equal(t, "LIT", cycles.Cycles[0].TeamKey)
	require.Equal(t, "Named cycle", cycles.Cycles[1].Name)
	require.Empty(t, cycles.Cycles[1].Description)
	require.Equal(t, "2026-01-30T00:00:00Z", cycles.Cycles[1].CompletedAt)
	require.Equal(t, "completed", cycles.Cycles[1].Status)
}

func Test_ListCyclesByTeam_wraps_graphql_errors(t *testing.T) {
	_, err := ListCyclesByTeam(context.Background(), errorGraphQLClient{err: errors.New("network down")}, "team-id", 1)

	require.Error(t, err)
	require.Contains(t, err.Error(), "list cycles")
}

func Test_GetCycleByID_returns_cycle(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"cycle": `{"cycle":{"id":"cycle-id","number":12,"name":"Named cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"}}}`,
	}

	cycle, err := GetCycleByID(context.Background(), graphqlClient, "cycle-id")

	require.NoError(t, err)
	require.Equal(t, "cycle-id", cycle.ID)
	require.Equal(t, "Named cycle", cycle.Name)
	require.Equal(t, "cycle body", cycle.Description)
	require.Equal(t, "active", cycle.Status)
	require.Equal(t, "LIT", cycle.TeamKey)
}

func Test_GetCycleByID_wraps_graphql_errors(t *testing.T) {
	_, err := GetCycleByID(context.Background(), errorGraphQLClient{err: errors.New("network down")}, "cycle-id")

	require.Error(t, err)
	require.Contains(t, err.Error(), "get cycle cycle-id")
}

func Test_CurrentCycleByTeam_returns_active_cycle(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"cycles": `{"cycles":{"nodes":[{"id":"future-cycle-id","number":13,"name":"Future cycle","description":null,"startsAt":"2099-01-01T00:00:00Z","endsAt":"2099-02-01T00:00:00Z","completedAt":null,"progress":0,"team":{"id":"team-id","key":"LIT","name":"linctl"}},{"id":"cycle-id","number":12,"name":"Current sprint","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	cycle, err := CurrentCycleByTeam(context.Background(), graphqlClient, "team-id")

	require.NoError(t, err)
	require.Equal(t, "cycle-id", cycle.ID)
	require.Equal(t, "Current sprint", cycle.Name)
	require.Equal(t, "active", cycle.Status)
}

func Test_CurrentCycleByTeam_reports_empty_active_cycle(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"cycles": `{"cycles":{"nodes":[{"id":"future-cycle-id","number":13,"name":"Future cycle","description":null,"startsAt":"2099-01-01T00:00:00Z","endsAt":"2099-02-01T00:00:00Z","completedAt":null,"progress":0,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	_, err := CurrentCycleByTeam(context.Background(), graphqlClient, "team-id")

	require.Error(t, err)
	require.Contains(t, err.Error(), "current sprint: no active Cycle")
}

func Test_CurrentCycleByTeam_wraps_graphql_errors(t *testing.T) {
	_, err := CurrentCycleByTeam(context.Background(), errorGraphQLClient{err: errors.New("network down")}, "team-id")

	require.Error(t, err)
	require.Contains(t, err.Error(), "current sprint: list cycles")
}

func Test_GetSprintReport_returns_cycle_and_issues(t *testing.T) {
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"CycleReport": `{"cycle":{"id":"cycle-id","number":12,"name":"Current sprint","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "Ship report",
			StateID:    "started",
			State:      "Started",
			StateType:  "started",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
	}

	report, err := GetSprintReport(context.Background(), graphqlClient, "cycle-id", 1)

	require.NoError(t, err)
	require.Equal(t, "cycle-id", report.Cycle.ID)
	require.Equal(t, "Current sprint", report.Cycle.Name)
	require.Equal(t, "LIT-1", report.Issues[0].Identifier)
	require.Equal(t, "Ship report", report.Issues[0].Title)
	require.True(t, report.HasNextPage)
	require.Equal(t, &endCursor, report.EndCursor)
}

func Test_GetSprintReport_wraps_graphql_errors(t *testing.T) {
	_, err := GetSprintReport(context.Background(), errorGraphQLClient{err: errors.New("network down")}, "cycle-id", 1)

	require.Error(t, err)
	require.Contains(t, err.Error(), "sprint report cycle-id")
}

func Test_ListCycleIssues_returns_issue_page(t *testing.T) {
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"cycle_issues": `{"cycle":{"id":"cycle-id","number":12,"name":"Current cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "Ship cycle issue",
			StateID:    "started",
			State:      "Started",
			StateType:  "started",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
	}

	issues, err := ListCycleIssues(context.Background(), graphqlClient, "cycle-id", 1)

	require.NoError(t, err)
	require.Equal(t, "cycle-id", issues.Cycle.ID)
	require.Equal(t, "Current cycle", issues.Cycle.Name)
	require.Equal(t, "LIT-1", issues.Issues[0].Identifier)
	require.Equal(t, "Ship cycle issue", issues.Issues[0].Title)
	require.True(t, issues.HasNextPage)
	require.Equal(t, &endCursor, issues.EndCursor)
}

func Test_ListCycleIssues_wraps_graphql_errors(t *testing.T) {
	_, err := ListCycleIssues(context.Background(), errorGraphQLClient{err: errors.New("network down")}, "cycle-id", 1)

	require.Error(t, err)
	require.Contains(t, err.Error(), "list cycle issues cycle-id")
}

func Test_ListCycleUncompletedIssuesUponClose_returns_issue_page(t *testing.T) {
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"cycle_uncompletedIssuesUponClose": `{"cycle":{"id":"cycle-id","number":12,"name":"Closed cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2026-01-15T00:00:00Z","completedAt":"2026-01-15T00:00:00Z","progress":0.75,"team":{"id":"team-id","key":"LIT","name":"linctl"},"uncompletedIssuesUponClose":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-2",
			Title:      "Carry issue forward",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
	}

	issues, err := ListCycleUncompletedIssuesUponClose(context.Background(), graphqlClient, "cycle-id", 1)

	require.NoError(t, err)
	require.Equal(t, "cycle-id", issues.Cycle.ID)
	require.Equal(t, "Closed cycle", issues.Cycle.Name)
	require.Equal(t, "LIT-2", issues.Issues[0].Identifier)
	require.Equal(t, "Carry issue forward", issues.Issues[0].Title)
	require.True(t, issues.HasNextPage)
	require.Equal(t, &endCursor, issues.EndCursor)
}

func Test_ListCycleUncompletedIssuesUponClose_wraps_graphql_errors(t *testing.T) {
	_, err := ListCycleUncompletedIssuesUponClose(
		context.Background(),
		errorGraphQLClient{err: errors.New("network down")},
		"cycle-id",
		1,
	)

	require.Error(t, err)
	require.Contains(t, err.Error(), "list cycle uncompleted issues cycle-id")
}

func Test_CycleStatus_describes_completion_and_date_edges(t *testing.T) {
	require.Equal(t, "completed", cycleStatus("2026-01-01T00:00:00Z", "2099-01-01T00:00:00Z", "2026-01-02T00:00:00Z"))
	require.Equal(t, "future", cycleStatus("2099-01-01T00:00:00Z", "2099-02-01T00:00:00Z", ""))
	require.Equal(t, "past", cycleStatus("2000-01-01T00:00:00Z", "2000-02-01T00:00:00Z", ""))
	require.Equal(t, "unknown", cycleStatus("not-a-date", "2099-01-01T00:00:00Z", ""))
}
