package client

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CreateProjectMilestone_returns_created_milestone_when_target_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "fixture",
			Status: "Backlog",
		}) + `}`,
		"ProjectMilestoneCreate": `{"projectMilestoneCreate":{"success":true,"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}}`,
	})}

	milestone, err := CreateProjectMilestone(
		context.Background(),
		recorder,
		matchingTarget(),
		ProjectMilestoneCreateRequest{
			ProjectID:   "project-id",
			Name:        "Launch milestone",
			Description: "milestone body",
			TargetDate:  "2026-06-30",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "project-milestone-id", milestone.ID)
	require.Equal(t, "Launch milestone", milestone.Name)
	require.Equal(t, "2026-06-30", milestone.TargetDate)
	require.JSONEq(t, `{
		"input": {
			"projectId": "project-id",
			"name": "Launch milestone",
			"description": "milestone body",
			"targetDate": "2026-06-30"
		}
	}`, string(recorder.variablesFor(t, "ProjectMilestoneCreate")))
}

func Test_UpdateProjectMilestone_returns_updated_milestone_when_target_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
		"ProjectMilestoneUpdate": `{"projectMilestoneUpdate":{"success":true,"projectMilestone":` +
			projectMilestoneJSON("Updated milestone", "done", "project-id") + `}}`,
	})

	milestone, err := UpdateProjectMilestone(
		context.Background(),
		graphqlClient,
		matchingTarget(),
		ProjectMilestoneUpdateRequest{
			ID:         "project-milestone-id",
			Name:       "Updated milestone",
			TargetDate: "2026-07-01",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "project-milestone-id", milestone.ID)
	require.Equal(t, "Updated milestone", milestone.Name)
	require.Equal(t, "done", milestone.Status)
}

func Test_UpdateProjectMilestone_refuses_when_pinned_project_differs(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Wrong project milestone", "next", "other-project") + `}`,
	})

	_, err := UpdateProjectMilestone(
		context.Background(),
		graphqlClient,
		matchingTarget(),
		ProjectMilestoneUpdateRequest{
			ID:   "project-milestone-id",
			Name: "Updated milestone",
		},
	)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateProjectMilestone_refuses_when_project_team_differs(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSONWithTeam("Wrong team milestone", "next", "other-project", "other-team", "OTHER") + `}`,
	})

	_, err := UpdateProjectMilestone(
		context.Background(),
		graphqlClient,
		config.Target{
			OrgID:   "org-id",
			TeamKey: "LIT",
			TeamID:  "team-id",
		},
		ProjectMilestoneUpdateRequest{
			ID:   "project-milestone-id",
			Name: "Updated milestone",
		},
	)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteProjectMilestone_removes_milestone_when_target_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
		"ProjectMilestoneDelete": `{"projectMilestoneDelete":{"success":true,"entityId":"project-milestone-id"}}`,
	})

	id, err := DeleteProjectMilestone(context.Background(), graphqlClient, matchingTarget(), "project-milestone-id")

	require.NoError(t, err)
	require.Equal(t, "project-milestone-id", id)
}

func Test_DeleteProjectMilestone_requires_id(t *testing.T) {
	_, err := DeleteProjectMilestone(
		context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "",
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_DeleteProjectMilestone_refuses_when_target_unresolved(t *testing.T) {
	_, err := DeleteProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, "project-milestone-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteProjectMilestone_refuses_when_pinned_project_differs(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Wrong project milestone", "next", "other-project") + `}`,
	})

	_, err := DeleteProjectMilestone(context.Background(), graphqlClient, matchingTarget(), "project-milestone-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteProjectMilestone_refuses_when_project_team_differs(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSONWithTeam("Wrong team milestone", "next", "other-project", "other-team", "OTHER") + `}`,
	})

	_, err := DeleteProjectMilestone(context.Background(), graphqlClient, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	}, "project-milestone-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteProjectMilestone_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
	})

	_, err := DeleteProjectMilestone(context.Background(), graphqlClient, matchingTarget(), "project-milestone-id")

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteProjectMilestone_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
		"ProjectMilestoneDelete": `{"projectMilestoneDelete":{"success":false,"entityId":"project-milestone-id"}}`,
	})

	_, err := DeleteProjectMilestone(context.Background(), graphqlClient, matchingTarget(), "project-milestone-id")

	require.ErrorIs(t, err, ErrMutationFailed)
}

// Test_UpdateProjectMilestone_succeeds_without_pinned_project_when_team_matches
// is the success-path counterpart to Test_UpdateProjectMilestone_refuses_when_pinned_project_differs:
// with no project pinned, requireProjectMilestone skips the project-id
// comparison but still enforces the pinned team
// (Test_UpdateProjectMilestone_refuses_when_project_team_differs).
func Test_UpdateProjectMilestone_succeeds_without_pinned_project_when_team_matches(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Any project milestone", "next", "any-project") + `}`,
		"ProjectMilestoneUpdate": `{"projectMilestoneUpdate":{"success":true,"projectMilestone":` +
			projectMilestoneJSON("Updated milestone", "done", "any-project") + `}}`,
	})

	milestone, err := UpdateProjectMilestone(
		context.Background(),
		graphqlClient,
		teamOnlyTarget(),
		ProjectMilestoneUpdateRequest{
			ID:   "project-milestone-id",
			Name: "Updated milestone",
		},
	)

	require.NoError(t, err)
	require.Equal(t, "Updated milestone", milestone.Name)
}

func projectMilestoneJSON(name string, status string, projectID string) string {
	return projectMilestoneJSONWithTeam(name, status, projectID, "team-id", "LIT")
}

func projectMilestoneJSONWithTeam(name string, status string, projectID string, teamID string, teamKey string) string {
	project := projectJSON(projectFixture{
		ID:     projectID,
		Name:   "fixture",
		Status: "Backlog",
	})
	project = strings.ReplaceAll(project, `"id":"team-id","key":"LIT"`, `"id":"`+teamID+`","key":"`+teamKey+`"`)

	return `{
		"id":"project-milestone-id",
		"name":"` + name + `",
		"description":"milestone body",
		"targetDate":"2026-06-30",
		"status":"` + status + `",
		"progress":0.5,
		"sortOrder":1,
		"project":` + project + `
	}`
}
