package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
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

// eoirIssueFixture is the destination-side issue: the same issue after it landed
// in the destination project.
func eoirIssueFixture(identifier string) issueFixture {
	issue := b1IssueFixture(identifier)
	issue.ProjectID = "eoir-project-id"
	issue.Project = "EOIR Case Scraper"

	return issue
}

func destinationProjectFixture() projectFixture {
	return projectFixture{ID: "eoir-project-id", Name: "EOIR Case Scraper", Status: "Backlog"}
}

func Test_MoveIssueProject_moves_issue_off_the_pinned_project(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue":   `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"project": `{"project":` + projectJSON(destinationProjectFixture()) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(
			eoirIssueFixture("LIT-1"),
		) + `}}`,
	})}

	issue, err := MoveIssueProject(context.Background(), recorder, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.NoError(t, err)
	require.Equal(t, "eoir-project-id", issue.ProjectID)
	require.Equal(t, "EOIR Case Scraper", issue.Project)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {"projectId": "eoir-project-id"}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_MoveIssueProject_requires_issue_id(t *testing.T) {
	_, err := MoveIssueProject(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), IssueMoveProjectRequest{
		ProjectID: "eoir-project-id",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_MoveIssueProject_requires_destination(t *testing.T) {
	_, err := MoveIssueProject(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), IssueMoveProjectRequest{
		IssueID: "LIT-1",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_MoveIssueProject_refuses_when_target_unresolved(t *testing.T) {
	_, err := MoveIssueProject(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, IssueMoveProjectRequest{IssueID: "LIT-1", ProjectID: "eoir-project-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_MoveIssueProject_refuses_issue_outside_the_pin(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(eoirIssueFixture("LIT-1")) + `}`,
	})}

	_, err := MoveIssueProject(context.Background(), recorder, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueProject_reports_unreadable_destination(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
	})}

	_, err := MoveIssueProject(context.Background(), recorder, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.Error(t, err)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueProject_refuses_destination_outside_the_pinned_team(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"project": `{"project":` + projectJSONWithTeam(
			destinationProjectFixture(), "ops-team-id", "OPS",
		) + `}`,
	})}

	_, err := MoveIssueProject(context.Background(), recorder, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueProject_refuses_when_issue_is_already_on_destination(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "fixture", Status: "Backlog",
		}) + `}`,
	})}

	_, err := MoveIssueProject(context.Background(), recorder, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "project-id",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_MoveIssueProject_reports_failed_mutation(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"project":     `{"project":` + projectJSON(destinationProjectFixture()) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":false,"issue":null}}`,
	})

	_, err := MoveIssueProject(context.Background(), graphqlClient, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_MoveIssueProject_reports_transport_failure(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":   `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"project": `{"project":` + projectJSON(destinationProjectFixture()) + `}`,
	})

	_, err := MoveIssueProject(context.Background(), graphqlClient, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.Error(t, err)
}

func Test_MoveIssueProject_reports_issue_landing_on_another_project(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":   `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"project": `{"project":` + projectJSON(destinationProjectFixture()) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(
			b1IssueFixture("LIT-1"),
		) + `}}`,
	})

	_, err := MoveIssueProject(context.Background(), graphqlClient, matchingTarget(), IssueMoveProjectRequest{
		IssueID:   "LIT-1",
		ProjectID: "eoir-project-id",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_AddProjectTeam_requires_project_id(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), ProjectAddTeamRequest{
		TeamID: "ops-team-id",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_refuses_when_target_unresolved(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectTeam_wraps_project_lookup_error(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamID:    "ops-team-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_refuses_when_project_lacks_pinned_team(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeam(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}, "other-team", "OTHER") + `}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectTeam_wraps_destination_lookup_error(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_wraps_mutation_error(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectTeam_fails_when_mutation_reports_no_success(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team":          `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":false,"project":null}}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_AddProjectTeam_fails_when_destination_missing_after_update(t *testing.T) {
	// Linear returned success without the destination on the project: that is not a
	// successful attach, even though the mutation payload claimed success.
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_AddProjectTeam_fails_when_update_drops_the_pinned_team(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSONWithTeam(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
			"ops-team-id", "OPS",
		) + `}}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectTeam_refuses_when_team_id_and_key_disagree(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
	}), matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamID:    "ops-team-id",
		TeamKey:   "WRONG",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_fails_when_destination_key_is_unknown(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"teams_list": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"description": "",
					"archivedAt": null,
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamKey: "MISSING"})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_wraps_team_list_error_when_resolving_by_key(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamKey: "OPS"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_pages_team_list_when_destination_is_not_on_first_page(t *testing.T) {
	// requestKey in the fake GraphQL client appends ":cursor" for after=cursor pages.
	client := projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"teams_list": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"description": "",
					"archivedAt": null,
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": true, "endCursor": "team-cursor-1"}
			}
		}`,
		"teams_list:team-cursor-1": teamsListWithDestinationJSON(),
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSONWithTwoTeams(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
		) + `}}`,
	})

	project, err := AddProjectTeam(context.Background(), client, matchingTarget(), ProjectAddTeamRequest{
		ProjectID: "project-id",
		TeamKey:   "OPS",
	})

	require.NoError(t, err)
	require.True(t, projectHasTeam(project, "ops-team-id", "OPS"))
}

func Test_AddProjectTeam_fails_when_team_list_page_has_no_end_cursor(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"teams_list": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"description": "",
					"archivedAt": null,
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": true, "endCursor": null}
			}
		}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamKey: "OPS"})

	require.Error(t, err)
	require.Contains(t, err.Error(), "end cursor")
}

func Test_AddProjectTeam_refuses_key_lookup_team_in_another_organization(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}) + `}`,
		"teams_list": `{
			"teams": {
				"nodes": [` + destinationTeamJSON("ops-team-id", "OPS", "other-org") + `],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamKey: "OPS"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectTeam_wraps_project_teams_list_error_when_truncated(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeamPage(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}, "team-id", "LIT", true) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddProjectTeam_fails_when_truncated_teams_page_is_empty(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeamPage(projectFixture{
			ID: "project-id", Name: "Harness", Status: "Backlog",
		}, "team-id", "LIT", true) + `}`,
		"team": `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"project_teams": `{
			"project": {
				"id": "project-id",
				"name": "Harness",
				"teams": {
					"nodes": [],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}
		}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_AddProjectTeam_fails_when_project_teams_page_has_no_end_cursor(t *testing.T) {
	_, err := AddProjectTeam(context.Background(), projectWriteFakeClient(map[string]string{
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
						{"id":"team-id","key":"LIT","name":"linctl-it","description":"","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}
					],
					"pageInfo": {"hasNextPage": true, "endCursor": null}
				}
			}
		}`,
	}), matchingTarget(), ProjectAddTeamRequest{ProjectID: "project-id", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddProjectTeam_continues_project_teams_pagination(t *testing.T) {
	client := projectWriteFakeClient(map[string]string{
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
						{"id":"team-id","key":"LIT","name":"linctl-it","description":"","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}
					],
					"pageInfo": {"hasNextPage": true, "endCursor": "teams-cursor-1"}
				}
			}
		}`,
		"project_teams:teams-cursor-1": `{
			"project": {
				"id": "project-id",
				"name": "Harness",
				"teams": {
					"nodes": [
						{"id":"extra-team-id","key":"EXT","name":"Extra","description":"","archivedAt":null,"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}
					],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}
		}`,
		"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSONWithTwoTeams(
			projectFixture{ID: "project-id", Name: "Harness", Status: "Backlog"},
		) + `}}`,
	})
	recorder := &recordingGraphQLClient{inner: client}

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

func Test_MoveIssueTeam_requires_issue_id(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), IssueMoveTeamRequest{
		TeamID: "ops-team-id",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_MoveIssueTeam_refuses_when_target_unresolved(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, IssueMoveTeamRequest{IssueID: "LIT-1", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_MoveIssueTeam_wraps_destination_lookup_error(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
	}), matchingTarget(), IssueMoveTeamRequest{IssueID: "LIT-1", TeamID: "ops-team-id"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrWriteInvalid)
}

func Test_MoveIssueTeam_wraps_mutation_error(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"team":  `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
	}), matchingTarget(), IssueMoveTeamRequest{IssueID: "LIT-1", TeamID: "ops-team-id"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_MoveIssueTeam_fails_when_mutation_reports_no_success(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"team":        `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":false,"issue":null}}`,
	}), matchingTarget(), IssueMoveTeamRequest{IssueID: "LIT-1", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_MoveIssueTeam_fails_when_issue_lands_on_the_wrong_team(t *testing.T) {
	_, err := MoveIssueTeam(context.Background(), issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"team":        `{"team":` + destinationTeamJSON("ops-team-id", "OPS", "org-id") + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	}), matchingTarget(), IssueMoveTeamRequest{IssueID: "LIT-1", TeamID: "ops-team-id"})

	require.ErrorIs(t, err, ErrMutationFailed)
}
