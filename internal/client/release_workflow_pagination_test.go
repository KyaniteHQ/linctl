package client

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListReleases_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"releases": `{"releases":{"nodes":[` +
			releasePageJSON("release-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"release-cursor-1"}}}`,
		"releases:release-cursor-1": `{"releases":{"nodes":[` +
			releasePageJSON("release-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	releases, err := ListReleases(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"release-1", "release-2"}, []string{
		releases.Releases[0].ID,
		releases.Releases[1].ID,
	})
	require.False(t, releases.HasNextPage)
}

func Test_ListReleaseHistory_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"release_history": `{"release":{"history":{"nodes":[` +
			releaseHistoryPageJSON("history-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"history-cursor-1"}}}}`,
		"release_history:history-cursor-1": `{"release":{"history":{"nodes":[` +
			releaseHistoryPageJSON("history-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	history, err := ListReleaseHistory(context.Background(), graphqlClient, "release-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"history-1", "history-2"}, []string{
		history.History[0].ID,
		history.History[1].ID,
	})
	require.False(t, history.HasNextPage)
}

func Test_ListReleasePipelineReleases_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"releasePipeline_releases": `{"releasePipeline":{"releases":{"nodes":[` +
			releasePageJSON("release-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"pipeline-release-cursor-1"}}}}`,
		"releasePipeline_releases:pipeline-release-cursor-1": `{"releasePipeline":{"releases":{"nodes":[` +
			releasePageJSON("release-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	releases, err := ListReleasePipelineReleases(context.Background(), graphqlClient, "release-pipeline-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"release-1", "release-2"}, []string{
		releases.Releases[0].ID,
		releases.Releases[1].ID,
	})
	require.False(t, releases.HasNextPage)
}

func Test_ListReleaseNotes_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"releaseNotes": `{"releaseNotes":{"nodes":[` +
			releaseNotePageJSON("release-note-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"note-cursor-1"}}}`,
		"releaseNotes:note-cursor-1": `{"releaseNotes":{"nodes":[` +
			releaseNotePageJSON("release-note-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	notes, err := ListReleaseNotes(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"release-note-1", "release-note-2"}, []string{
		notes.ReleaseNotes[0].ID,
		notes.ReleaseNotes[1].ID,
	})
	require.False(t, notes.HasNextPage)
}

func Test_ListWorkflowStateIssues_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"workflowState_issues": `{"workflowState":{"id":"workflow-state-id","name":"Started","issues":{"nodes":[` +
			issueJSON(issueFixture{
				ID:         "issue-1",
				Identifier: "LIT-1",
				Title:      "First",
				StateID:    "workflow-state-id",
				State:      "Started",
				StateType:  "started",
			}) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"issue-cursor-1"}}}}`,
		"workflowState_issues:issue-cursor-1": `{"workflowState":{"id":"workflow-state-id","name":"Started","issues":{"nodes":[` +
			issueJSON(issueFixture{
				ID:         "issue-2",
				Identifier: "LIT-2",
				Title:      "Second",
				StateID:    "workflow-state-id",
				State:      "Started",
				StateType:  "started",
			}) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	issues, err := ListWorkflowStateIssues(context.Background(), graphqlClient, "workflow-state-id", 100)

	require.NoError(t, err)
	require.Equal(t, "workflow-state-id", issues.WorkflowStateID)
	require.Equal(t, "Started", issues.WorkflowStateName)
	require.Equal(t, []string{"issue-1", "issue-2"}, []string{
		issues.Issues[0].ID,
		issues.Issues[1].ID,
	})
	require.False(t, issues.HasNextPage)
}

func Test_ListTimeSchedules_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"timeSchedules": `{"timeSchedules":{"nodes":[` +
			timeSchedulePageJSON("time-schedule-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"schedule-cursor-1"}}}`,
		"timeSchedules:schedule-cursor-1": `{"timeSchedules":{"nodes":[` +
			timeSchedulePageJSON("time-schedule-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	schedules, err := ListTimeSchedules(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"time-schedule-1", "time-schedule-2"}, []string{
		schedules.TimeSchedules[0].ID,
		schedules.TimeSchedules[1].ID,
	})
	require.False(t, schedules.HasNextPage)
}

func Test_ListTriageResponsibilities_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"triageResponsibilities": `{"triageResponsibilities":{"nodes":[` +
			triageResponsibilityPageJSON("triage-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"triage-cursor-1"}}}`,
		"triageResponsibilities:triage-cursor-1": `{"triageResponsibilities":{"nodes":[` +
			triageResponsibilityPageJSON("triage-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	responsibilities, err := ListTriageResponsibilities(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"triage-1", "triage-2"}, []string{
		responsibilities.TriageResponsibilities[0].ID,
		responsibilities.TriageResponsibilities[1].ID,
	})
	require.False(t, responsibilities.HasNextPage)
}

func releasePageJSON(id string) string {
	return strings.Replace(releaseJSON(), `"id":"release-id"`, `"id":"`+id+`"`, 1)
}

func releaseHistoryPageJSON(id string) string {
	return strings.Replace(releaseHistoryJSON(), `"id":"release-history-id"`, `"id":"`+id+`"`, 1)
}

func releaseNotePageJSON(id string) string {
	return strings.Replace(releaseNoteJSON(), `"id":"release-note-id"`, `"id":"`+id+`"`, 1)
}

func timeSchedulePageJSON(id string) string {
	return `{"id":"` + id + `","name":"Primary on-call","createdAt":"2026-06-19T12:00:00Z",` +
		`"updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,"externalId":"pd-primary",` +
		`"externalUrl":"https://example.com/schedule","integration":{"id":"integration-id"},` +
		`"entries":[{"startsAt":"2026-06-20T00:00:00Z","endsAt":"2026-06-21T00:00:00Z",` +
		`"userId":"user-id","userEmail":"omer@example.com"}]}`
}

func triageResponsibilityPageJSON(id string) string {
	return strings.Replace(triageResponsibilityJSON(), `"id":"triage-responsibility-id"`, `"id":"`+id+`"`, 1)
}
