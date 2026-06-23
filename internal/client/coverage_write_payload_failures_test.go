package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientWriteFailureScenarios_fail_when_issue_payload_omits_entity(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"IssueCreate": `{"issueCreate":{"success":false,"issue":null}}`,
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-20",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"IssueUpdate":        `{"issueUpdate":{"success":false,"issue":null}}`,
		"IssueCommentCreate": `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"body","url":"url","issue":null}}}`,
		"CompletedWorkflowStates": `{"workflowStates":{"nodes":[
			{"id":"done-state","name":"Done","type":"completed","position":1}
		]}}`,
		"StartedWorkflowStates": `{"workflowStates":{"nodes":[
			{"id":"started-state","name":"Started","type":"started","position":1}
		]}}`,
		"IssueClose": `{"issueUpdate":{"success":false,"issue":null}}`,
	})

	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{Title: "title"})
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{ID: "LIT-20", Title: "title"})
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{ID: "LIT-20", Body: "body"})
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = StartIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-20")
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = CloseIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-20")
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_ClientWriteFailureScenarios_fail_when_project_payload_omits_entity(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"ProjectCreate": `{"projectCreate":{"success":false,"project":null}}`,
		"project": `{"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "fixture",
			Status: "Backlog",
		}) + `}`,
		"ProjectUpdate":  `{"projectUpdate":{"success":false,"project":null}}`,
		"ProjectArchive": `{"projectArchive":{"success":false,"entity":null}}`,
	})

	_, err := CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{Name: "name"})
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = ArchiveProject(context.Background(), graphqlClient, matchingTarget(), "project-id")
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = CreateProjectMilestone(
		context.Background(),
		projectWriteFakeClient(map[string]string{
			"project":                `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectMilestoneCreate": `{"projectMilestoneCreate":{"success":false,"projectMilestone":null}}`,
		}),
		matchingTarget(),
		ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"},
	)
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = UpdateProjectMilestone(
		context.Background(),
		projectWriteFakeClient(map[string]string{
			"projectMilestone": `{"projectMilestone":` +
				projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
			"ProjectMilestoneUpdate": `{"projectMilestoneUpdate":{"success":false,"projectMilestone":null}}`,
		}),
		matchingTarget(),
		ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"},
	)
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = CreateCycle(
		context.Background(),
		projectWriteFakeClient(map[string]string{
			"CycleCreate": `{"cycleCreate":{"success":false,"cycle":null}}`,
		}),
		matchingTarget(),
		CycleCreateRequest{StartsAt: "2026-07-01T00:00:00Z", EndsAt: "2026-07-15T00:00:00Z"},
	)
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = UpdateCycle(
		context.Background(),
		projectWriteFakeClient(map[string]string{
			"cycle":       `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
			"CycleUpdate": `{"cycleUpdate":{"success":false,"cycle":null}}`,
		}),
		matchingTarget(),
		CycleUpdateRequest{ID: "cycle-id", Name: "name"},
	)
	require.ErrorIs(t, err, ErrMutationFailed)

	_, err = ArchiveCycle(
		context.Background(),
		projectWriteFakeClient(map[string]string{
			"cycle":        `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
			"CycleArchive": `{"cycleArchive":{"success":false,"entity":null}}`,
		}),
		matchingTarget(),
		"cycle-id",
	)
	require.ErrorIs(t, err, ErrMutationFailed)
}
