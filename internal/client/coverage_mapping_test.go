package client

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SummaryMappingScenarios_preserve_optional_people(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issue": `{"issue":` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-30",
			Title:      "assigned",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}, "Omer") + `}`,
		"project": `{"project":` + projectJSONWithLead(projectFixture{
			ID:     "project-id",
			Name:   "led",
			Status: "Backlog",
		}, "Omer") + `}`,
	}

	issue, err := GetIssueByID(context.Background(), graphqlClient, "LIT-30")
	require.NoError(t, err)
	require.Equal(t, "Omer", issue.Assignee)

	project, err := GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.NoError(t, err)
	require.Equal(t, "Omer", project.Lead)
}

func Test_SummaryMappingScenarios_preserve_reference_domain_variants(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"document": `{"document":{
			"id":"document-id",
			"title":"Cycle spec",
			"slugId":"cycle-spec",
			"archivedAt":"2026-06-19T12:00:00Z",
			"project":{"id":"project-id","name":"Pinned project"},
			"team":{"id":"team-id","key":"LIT","name":"linctl"},
			"issue":{"id":"issue-id","identifier":"LIT-1","title":"Issue"},
			"cycle":{"id":"cycle-id","number":7,"name":"Planning"}
		}}`,
		"team": `{"team":{
			"id":"team-id",
			"key":"LIT",
			"name":"linctl",
			"description":null,
			"archivedAt":"2026-06-19T12:00:00Z",
			"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
		}}`,
	}

	document, err := GetDocumentByID(context.Background(), graphqlClient, "document-id")
	require.NoError(t, err)
	require.Equal(t, "2026-06-19T12:00:00Z", document.ArchivedAt)
	require.Equal(t, "cycle", document.ParentType)
	require.Equal(t, "Planning", document.ParentName)

	team, err := GetTeamByID(context.Background(), graphqlClient, "team-id")
	require.NoError(t, err)
	require.Empty(t, team.Description)
	require.Equal(t, "2026-06-19T12:00:00Z", team.ArchivedAt)
}

func Test_SummaryMappingScenarios_preserve_release_note_without_generation_status(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"releaseNote": `{"releaseNote":` + strings.Replace(
			releaseNoteJSON(),
			`"generationStatus":"completed"`,
			`"generationStatus":null`,
			1,
		) + `}`,
	}

	note, err := GetReleaseNoteByID(context.Background(), graphqlClient, "release-note-id")

	require.NoError(t, err)
	require.Empty(t, note.GenerationStatus)
	require.Equal(t, "Launch notes", note.Title)
}
