package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientWriteFailureScenarios_wrap_graphql_operation_errors(t *testing.T) {
	operationErr := errors.New("linear unavailable")

	_, err := CreateIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"IssueCreate": "",
	}).withError(operationErr), matchingTarget(), IssueCreateRequest{Title: "title"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "create issue")

	_, err = UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-40",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"IssueUpdate": "",
	}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-40", Title: "title"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "update issue LIT-40")

	_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-40",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"IssueCommentCreate": "",
	}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-40", Body: "body"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "comment on issue LIT-40")

	_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-40",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"StartedWorkflowStates": `{"workflowStates":{"nodes":[{"id":"started-state","name":"Started","type":"started","position":1}]}}`,
		"IssueUpdate":           "",
	}).withError(operationErr), matchingTarget(), "LIT-40")
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "start issue LIT-40")

	_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-40",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"CompletedWorkflowStates": `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`,
		"IssueClose":              "",
	}).withError(operationErr), matchingTarget(), "LIT-40")
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "close issue LIT-40")

	_, err = CreateProject(context.Background(), projectWriteFakeClient(map[string]string{
		"ProjectCreate": "",
	}).withError(operationErr), matchingTarget(), ProjectCreateRequest{Name: "name"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "create project")

	_, err = UpdateProject(context.Background(), projectWriteFakeClient(map[string]string{
		"project":       `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"ProjectUpdate": "",
	}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "update project project-id")

	_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
		"project":        `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"ProjectArchive": "",
	}).withError(operationErr), matchingTarget(), "project-id")
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "archive project project-id")

	_, err = CreateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
		"project":                `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"ProjectMilestoneCreate": "",
	}).withError(operationErr), matchingTarget(), ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "create project milestone")

	_, err = UpdateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
		"ProjectMilestoneUpdate": "",
	}).withError(operationErr), matchingTarget(), ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "update project milestone project-milestone-id")

	_, err = CreateCycle(context.Background(), projectWriteFakeClient(map[string]string{
		"CycleCreate": "",
	}).withError(operationErr), matchingTarget(), CycleCreateRequest{
		StartsAt: "2026-07-01T00:00:00Z",
		EndsAt:   "2026-07-15T00:00:00Z",
	})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "create cycle")

	_, err = UpdateCycle(context.Background(), projectWriteFakeClient(map[string]string{
		"cycle":       `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
		"CycleUpdate": "",
	}).withError(operationErr), matchingTarget(), CycleUpdateRequest{ID: "cycle-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "update cycle cycle-id")

	_, err = ArchiveCycle(context.Background(), projectWriteFakeClient(map[string]string{
		"cycle":        `{"cycle":` + cycleJSON("Planning cycle", "team-id", "LIT") + `}`,
		"CycleArchive": "",
	}).withError(operationErr), matchingTarget(), "cycle-id")
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "archive cycle cycle-id")
}

func Test_ClientWriteFailureScenarios_return_guard_read_errors(t *testing.T) {
	operationErr := errors.New("guard read failed")

	_, err := UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": "",
	}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-50", Title: "title"})
	require.ErrorIs(t, err, operationErr)

	_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": "",
	}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-50", Body: "body"})
	require.ErrorIs(t, err, operationErr)

	_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": "",
	}).withError(operationErr), matchingTarget(), "LIT-50")
	require.ErrorIs(t, err, operationErr)

	_, err = StartIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-51",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"StartedWorkflowStates": "",
	}).withError(operationErr), matchingTarget(), "LIT-51")
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "list started workflow states")

	_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": "",
	}).withError(operationErr), matchingTarget(), "LIT-50")
	require.ErrorIs(t, err, operationErr)

	_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-51",
			Title:      "target",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"CompletedWorkflowStates": "",
	}).withError(operationErr), matchingTarget(), "LIT-51")
	require.ErrorIs(t, err, operationErr)
	require.Contains(t, err.Error(), "list completed workflow states")

	_, err = UpdateProject(context.Background(), projectWriteFakeClient(map[string]string{
		"project": "",
	}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)

	_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
		"project": "",
	}).withError(operationErr), matchingTarget(), "project-id")
	require.ErrorIs(t, err, operationErr)

	_, err = CreateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
		"project": "",
	}).withError(operationErr), matchingTarget(), ProjectMilestoneCreateRequest{ProjectID: "project-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)

	_, err = UpdateProjectMilestone(context.Background(), projectWriteFakeClient(map[string]string{
		"projectMilestone": "",
	}).withError(operationErr), matchingTarget(), ProjectMilestoneUpdateRequest{ID: "project-milestone-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)

	_, err = UpdateCycle(context.Background(), projectWriteFakeClient(map[string]string{
		"cycle": "",
	}).withError(operationErr), matchingTarget(), CycleUpdateRequest{ID: "cycle-id", Name: "name"})
	require.ErrorIs(t, err, operationErr)

	_, err = ArchiveCycle(context.Background(), projectWriteFakeClient(map[string]string{
		"cycle": "",
	}).withError(operationErr), matchingTarget(), "cycle-id")
	require.ErrorIs(t, err, operationErr)
}
