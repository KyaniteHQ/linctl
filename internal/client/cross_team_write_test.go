package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func destinationTeamJSON(teamID string, teamKey string, orgID string) string {
	return `{
		"id": "` + teamID + `",
		"key": "` + teamKey + `",
		"name": "` + teamKey + `",
		"description": "",
		"archivedAt": null,
		"organization": {"id": "` + orgID + `", "name": "Kyanite", "urlKey": "kyanite"}
	}`
}

func teamsListWithDestinationJSON() string {
	return `{
		"teams": {
			"nodes": [
				{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"description": "",
					"archivedAt": null,
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				},
				` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `
			],
			"pageInfo": {"hasNextPage": false, "endCursor": null}
		}
	}`
}

func projectJSONWithTwoTeams(project projectFixture) string {
	return `{
		"id":"` + project.ID + `",
		"name":"` + project.Name + `",
		"description":"description",
		"archivedAt":null,
		"slugId":"` + project.Name + `",
		"url":"https://linear.app/kyanite/project/` + project.ID + `",
		"priority":0,
		"status":{"id":"status-id","name":"` + project.Status + `","type":"backlog"},
		"lead":null,
		"teams":{
			"nodes":[
				{"id":"team-id","key":"LIT","name":"linctl-it"},
				{"id":"ops-team-id","key":"OPS","name":"OPS"}
			],
			"pageInfo":{"hasNextPage":false}
		}
	}`
}

func issueJSONWithTeam(issue issueFixture, teamID string, teamKey string) string {
	project := `null`
	if issue.ProjectID != "" {
		project = `{"id":"` + issue.ProjectID + `","name":"` + issue.Project + `"}`
	}

	return `{
		"id":"issue-id",
		"identifier":"` + issue.Identifier + `",
		"title":"` + issue.Title + `",
		"url":"https://linear.app/kyanite/issue/` + issue.Identifier + `",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"` + teamID + `","key":"` + teamKey + `","name":"` + teamKey + `"},
		"state":{"id":"` + issue.StateID + `","name":"` + issue.State + `","type":"` + issue.StateType + `"},
		"assignee":null,
		"project":` + project + `
	}`
}

func Test_AddProjectTeam_merges_destination_without_dropping_existing(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSONWithTwoTeams(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
		) + `}}`,
	})}

	project, err := AddProjectTeam(context.Background(), recorder, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamID:    "ops-team-id",
	})

	require.NoError(t, err)
	require.True(t, projectHasTeam(project, "team-id", "LIT"))
	require.True(t, projectHasTeam(project, "ops-team-id", "OPS"))
	require.JSONEq(t, `{
		"id": "project-id",
		"input": {"teamIds": ["team-id", "ops-team-id"]}
	}`, string(recorder.variablesFor(t, "ProjectUpdate")))
}

func Test_AddProjectTeam_is_idempotent_when_team_already_attached(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTwoTeams(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
		) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
	})}

	project, err := AddProjectTeam(context.Background(), recorder, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamID:    "ops-team-id",
		TeamKey:   "OPS",
	})

	require.NoError(t, err)
	require.True(t, projectHasTeam(project, "ops-team-id", "OPS"))
	require.False(t, recorder.sentOperation("ProjectUpdate"))
}

func Test_AddProjectTeam_refuses_destination_in_another_organization(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "other-org") + `}`,
	})}

	_, err := AddProjectTeam(context.Background(), recorder, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamID:    "ops-team-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectUpdate"))
}

func Test_AddProjectTeam_refuses_when_project_is_not_pinned(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "other-project", Name: "Other", Status: "Backlog",
		}) + `}`,
	})}

	_, err := AddProjectTeam(context.Background(), recorder, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "other-project",
		TeamID:    "ops-team-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("ProjectUpdate"))
}

func Test_AddProjectTeam_requires_destination(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_resolves_destination_by_key(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"teams_list": teamsListWithDestinationJSON(),
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSONWithTwoTeams(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
		) + `}}`,
	})}

	project, err := AddProjectTeam(context.Background(), recorder, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamKey:   "OPS",
	})

	require.NoError(t, err)
	require.True(t, projectHasTeam(project, "ops-team-id", "OPS"))
	require.JSONEq(t, `{
		"id": "project-id",
		"input": {"teamIds": ["team-id", "ops-team-id"]}
	}`, string(recorder.variablesFor(t, "ProjectUpdate")))
}

func Test_AddProjectTeam_paginates_when_teams_page_is_truncated(t *testing.T) {
	// A truncated first page must be fully walked before merge: otherwise an
	// unlisted team would be dropped by projectUpdate's full-set teamIds field.
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeamPage(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}, "team-id", "LIT", true) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"project_teams": `{
			"project": {
				"id": "project-id",
				"name": "Harness",
				"teams": {
					"nodes": [
						{"id":"team-id","key":"LIT","name":"linctl-it","description":"","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}},
						{"id":"extra-team-id","key":"EXT","name":"Extra","description":"","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}
					],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}
		}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSONWithTwoTeams(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
		) + `}}`,
	})}

	_, err := AddProjectTeam(context.Background(), recorder, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamID:    "ops-team-id",
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"id": "project-id",
		"input": {"teamIds": ["team-id", "extra-team-id", "ops-team-id"]}
	}`, string(recorder.variablesFor(t, "ProjectUpdate")))
}

func Test_MoveIssueTeam_moves_issue_off_the_pinned_team(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"team":  `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSONWithTeam(
			b1IssueFixture("OPS-1"), "ops-team-id", "OPS",
		) + `}}`,
	})}

	issue, err := MoveIssueTeam(context.Background(), recorder, matchingTarget(), IssueMoveTeamRequest{
		IssueID: "LIT-1",
		TeamID:  "ops-team-id",
	})

	require.NoError(t, err)
	require.Equal(t, "ops-team-id", issue.TeamID)
	require.Equal(t, "OPS", issue.Team)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"teamId": "ops-team-id"}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_MoveIssueTeam_refuses_issue_outside_the_pin(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSONWithTeam(b1IssueFixture("OPS-1"), "ops-team-id", "OPS") + `}`,
	})}

	_, err := MoveIssueTeam(context.Background(), recorder, matchingTarget(), IssueMoveTeamRequest{
		IssueID: "OPS-1",
		TeamID:  "ops-team-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueTeam_refuses_destination_in_another_organization(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"team":  `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "other-org") + `}`,
	})}

	_, err := MoveIssueTeam(context.Background(), recorder, matchingTarget(), IssueMoveTeamRequest{
		IssueID: "LIT-1",
		TeamID:  "ops-team-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueTeam_refuses_when_issue_is_already_on_destination(t *testing.T) {
	// Same team id/key as the pin, so the issue matches requireIssueDetail but
	// the destination is not a move.
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"team":  `{"team":` + destinationTeamJSON("team-id", "LIT", "org-id") + `}`,
	})}

	_, err := MoveIssueTeam(context.Background(), recorder, matchingTarget(), IssueMoveTeamRequest{
		IssueID: "LIT-1",
		TeamID:  "team-id",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueTeam_requires_destination(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), IssueMoveTeamRequest{
		IssueID: "LIT-1",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_MoveIssueTeam_resolves_destination_by_key(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":      `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"teams_list": teamsListWithDestinationJSON(),
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSONWithTeam(
			b1IssueFixture("OPS-1"), "ops-team-id", "OPS",
		) + `}}`,
	})}

	issue, err := MoveIssueTeam(context.Background(), recorder, matchingTarget(), IssueMoveTeamRequest{
		IssueID: "LIT-1",
		TeamKey: "OPS",
	})

	require.NoError(t, err)
	require.Equal(t, "OPS", issue.Team)
}
