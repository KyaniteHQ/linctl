package client

import (
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CreateIssue_returns_created_issue_when_target_matches(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-2",
			Title:      "created",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}}`,
	})

	// When
	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:       "created",
		Description: "body",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-2", issue.Identifier)
	require.Equal(t, "project-id", issue.ProjectID)
}

func Test_CreateIssue_refuses_label_from_another_team(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "other", "other-team", "OTHER") + `}`,
	})}

	_, err := CreateIssue(context.Background(), recorder, matchingTarget(), IssueCreateRequest{
		Title:    "created",
		LabelIDs: []string{"label-id"},
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueCreate"))
}

func Test_UpdateIssue_refuses_label_from_another_team(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "existing",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "other", "other-team", "OTHER") + `}`,
	})}

	_, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:       "LIT-1",
		LabelIDs: []string{"label-id"},
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueUpdate"))
}

func Test_UpdateIssue_refuses_when_issue_project_does_not_match_pinned_project(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-3",
			Title:      "wrong project",
			ProjectID:  "other-project",
			Project:    "other",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
	})

	// When
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:    "LIT-3",
		Title: "new title",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

// Test_UpdateIssue_succeeds_without_pinned_project_when_team_matches is the
// success-path counterpart to Test_UpdateIssue_refuses_when_issue_project_does_not_match_pinned_project:
// with no project pinned, requireIssueDetail skips the project comparison but
// still enforces the team match.
func Test_UpdateIssue_succeeds_without_pinned_project_when_team_matches(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-3",
			Title:      "any project",
			ProjectID:  "other-project",
			Project:    "other",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-3",
			Title:      "new title",
			ProjectID:  "other-project",
			Project:    "other",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}}`,
	})

	// When
	issue, err := UpdateIssue(context.Background(), graphqlClient, teamOnlyTarget(), IssueUpdateRequest{
		ID:    "LIT-3",
		Title: "new title",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "new title", issue.Title)
}

func Test_StartIssue_assigns_viewer_and_moves_to_started_state_when_target_matches(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-5",
			Title:      "start",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo-state",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"StartedWorkflowStates": `{"workflowStates":{"nodes":[
			{"id":"later-started-state","name":"Later","type":"started","position":2},
			{"id":"started-state","name":"Started","type":"started","position":1}
		]}}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-5",
			Title:      "start",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "started-state",
			State:      "Started",
			StateType:  "started",
		}, "Omer") + `}}`,
	})

	// When
	issue, err := StartIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-5")

	// Then
	require.NoError(t, err)
	require.Equal(t, "started", issue.StateType)
	require.Equal(t, "started-state", issue.StateID)
	require.Equal(t, "Omer", issue.Assignee)
}

func Test_CloseIssue_moves_issue_to_completed_state_when_target_matches(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-4",
			Title:      "close",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo-state",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"CompletedWorkflowStates": `{"workflowStates":{"nodes":[
			{"id":"done-state","name":"Done","type":"completed","position":2},
			{"id":"complete-state","name":"Complete","type":"completed","position":1}
		]}}`,
		"IssueClose": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-4",
			Title:      "close",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "complete-state",
			State:      "Complete",
			StateType:  "completed",
		}) + `}}`,
	})

	// When
	issue, err := CloseIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-4")

	// Then
	require.NoError(t, err)
	require.Equal(t, "completed", issue.StateType)
	require.Equal(t, "complete-state", issue.StateID)
}

type issueWriteFakeClient map[string]string

func (client issueWriteFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	return fakeGraphQLClient(client.withTargetResponses()).MakeRequest(ctx, request, response)
}

func (client issueWriteFakeClient) withTargetResponses() map[string]string {
	responses := map[string]string{
		"Viewer": `{
			"viewer": {
				"id": "user-id",
				"name": "Omer",
				"displayName": "Omer",
				"email": "omer@example.com",
				"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
			}
		}`,
		"Teams": `{
			"teams": {
				"nodes": [{
					"id": "team-id",
					"key": "LIT",
					"name": "linctl-it",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}],
				"pageInfo": {"hasNextPage": false, "endCursor": null}
			}
		}`,
		"TargetProject": `{
			"project": {
				"id": "project-id",
				"name": "fixture",
				"teams": {
					"nodes": [{
						"id": "team-id",
						"key": "LIT",
						"name": "linctl-it",
						"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
					}]
				}
			}
		}`,
	}
	for operation, response := range client {
		responses[operation] = response
	}

	return responses
}

func matchingTarget() config.Target {
	return config.Target{
		OrgID:     "org-id",
		TeamKey:   "LIT",
		TeamID:    "team-id",
		ProjectID: "project-id",
	}
}

// teamOnlyTarget is a pinned target without a project, exercising the
// project-unpinned guard mode (CONTEXT.md: project_id optional).
func teamOnlyTarget() config.Target {
	return config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	}
}

type issueFixture struct {
	Identifier string
	Title      string
	ProjectID  string
	Project    string
	StateID    string
	State      string
	StateType  string
}

func Test_CreateIssue_resolves_state_type_and_priority_when_provided(t *testing.T) {
	// Given
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStatesByType": `{"workflowStates":{"nodes":[
			{"id":"todo-state","name":"Todo","type":"unstarted","position":2},
			{"id":"backlog-state","name":"Backlog","type":"unstarted","position":1}
		]}}`,
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-3",
			Title:      "typed",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "backlog-state",
			State:      "Backlog",
			StateType:  "unstarted",
		}) + `}}`,
	})}

	// When
	issue, err := CreateIssue(context.Background(), recorder, matchingTarget(), IssueCreateRequest{
		Title:     "typed",
		StateType: "unstarted",
		Priority:  "2",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-3", issue.Identifier)
	require.JSONEq(t, `{
		"input": {
			"title": "typed",
			"teamId": "team-id",
			"projectId": "project-id",
			"stateId": "backlog-state",
			"priority": 2
		}
	}`, string(recorder.variablesFor(t, "IssueCreate")))
}

func Test_CreateIssue_returns_error_when_state_type_has_no_workflow_states(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"WorkflowStatesByType": `{"workflowStates":{"nodes":[]}}`,
	})

	// When
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:     "typed",
		StateType: "unstarted",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_returns_error_when_state_type_has_no_workflow_states(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "existing",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo-state",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"WorkflowStatesByType": `{"workflowStates":{"nodes":[]}}`,
	})

	// When
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:        "LIT-1",
		StateType: "completed",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_resolves_state_type_and_priority_when_provided(t *testing.T) {
	// Given
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "existing",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo-state",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
		"WorkflowStatesByType": `{"workflowStates":{"nodes":[
			{"id":"done-state","name":"Done","type":"completed","position":1}
		]}}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "existing",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "done-state",
			State:      "Done",
			StateType:  "completed",
		}) + `}}`,
	})}

	// When
	issue, err := UpdateIssue(context.Background(), recorder, matchingTarget(), IssueUpdateRequest{
		ID:        "LIT-1",
		StateType: "completed",
		Priority:  "1",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "completed", issue.StateType)
	require.JSONEq(t, `{
		"id": "LIT-1",
		"input": {
			"stateId": "done-state",
			"priority": 1
		}
	}`, string(recorder.variablesFor(t, "IssueUpdate")))
}

func Test_UpdateIssue_returns_error_when_all_fields_empty(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})

	// When
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID: "LIT-1",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssue_returns_error_for_invalid_priority_string(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})

	// When
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:    "typed",
		Priority: "not-a-number",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_returns_error_for_invalid_priority_string(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "existing",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo-state",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}`,
	})

	// When
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:       "LIT-1",
		Priority: "not-a-number",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_firstStateIDOfType_returns_error_on_graphql_failure(t *testing.T) {
	// Given - empty fake client with no WorkflowStatesByType response triggers error
	graphqlClient := fakeGraphQLClient(map[string]string{})

	// When
	_, err := firstStateIDOfType(context.Background(), graphqlClient, "team-id", "started")

	// Then
	require.ErrorContains(t, err, "list started workflow states")
}

func Test_parsePriority_returns_nil_for_empty_string(t *testing.T) {
	result, err := parsePriority("")

	require.NoError(t, err)
	require.Nil(t, result)
}

func Test_parsePriority_returns_error_for_non_numeric_string(t *testing.T) {
	_, err := parsePriority("high")

	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_parsePriority_rejects_out_of_range_values(t *testing.T) {
	for _, raw := range []string{"-1", "5"} {
		t.Run(raw, func(t *testing.T) {
			_, err := parsePriority(raw)

			require.Error(t, err)
			require.ErrorIs(t, err, ErrWriteInvalid)
		})
	}
}

func Test_firstStateIDOfType_returns_state_with_lowest_position(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient(map[string]string{
		"WorkflowStatesByType": `{"workflowStates":{"nodes":[
			{"id":"second-state","name":"Second","type":"started","position":2},
			{"id":"first-state","name":"First","type":"started","position":1}
		]}}`,
	})

	// When
	stateID, err := firstStateIDOfType(context.Background(), graphqlClient, "team-id", "started")

	// Then
	require.NoError(t, err)
	require.Equal(t, "first-state", stateID)
}

func Test_firstStateIDOfType_returns_error_when_no_states(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient(map[string]string{
		"WorkflowStatesByType": `{"workflowStates":{"nodes":[]}}`,
	})

	// When
	_, err := firstStateIDOfType(context.Background(), graphqlClient, "team-id", "started")

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func issueJSON(issue issueFixture) string {
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
		"team":{"id":"team-id","key":"LIT","name":"linctl-it"},
		"state":{"id":"` + issue.StateID + `","name":"` + issue.State + `","type":"` + issue.StateType + `"},
		"assignee":null,
		"project":` + project + `
	}`
}

func b1IssueFixture(identifier string) issueFixture {
	return issueFixture{
		Identifier: identifier,
		Title:      "b1",
		ProjectID:  "project-id",
		Project:    "fixture",
		StateID:    "todo-state",
		State:      "Todo",
		StateType:  "unstarted",
	}
}

func Test_CreateIssue_sets_assignee_label_and_due_date(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel":  `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-7")) + `}}`,
	})

	// When
	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:      "b1",
		AssigneeID: "user-id",
		LabelIDs:   []string{"label-id"},
		DueDate:    "2026-07-01",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-7", issue.Identifier)
}

func Test_CreateIssue_rejects_invalid_due_date(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})

	// When
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:   "b1",
		DueDate: "July 1",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_sets_assignee_label_and_due_date(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"issueLabel":  `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})

	// When
	issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:         "LIT-1",
		AssigneeID: "user-id",
		LabelIDs:   []string{"label-id"},
		DueDate:    "2026-07-01",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_UpdateIssue_clears_due_date_when_requested(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":       `{"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}`,
		"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})

	// When
	issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:           "LIT-1",
		ClearDueDate: true,
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_UpdateIssue_rejects_due_date_with_clear_due_date(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})

	// When
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:           "LIT-1",
		DueDate:      "2026-07-01",
		ClearDueDate: true,
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateIssue_rejects_invalid_due_date(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})

	// When
	_, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:      "LIT-1",
		DueDate: "not-a-date",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssue_assigns_milestone_after_pinned_target_check(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Launch milestone", "next", "project-id") + `}`,
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-9",
			Title:      "with milestone",
			ProjectID:  "project-id",
			Project:    "fixture",
		}) + `}}`,
	})

	// When
	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:              "with milestone",
		ProjectMilestoneID: "project-milestone-id",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-9", issue.Identifier)
}

func Test_CreateIssue_rejects_milestone_without_pinned_project(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})
	target := matchingTarget()
	target.ProjectID = ""

	// When
	_, err := CreateIssue(context.Background(), graphqlClient, target, IssueCreateRequest{
		Title:              "with milestone",
		ProjectMilestoneID: "project-milestone-id",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssue_refuses_milestone_outside_pinned_project(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"projectMilestone": `{"projectMilestone":` +
			projectMilestoneJSON("Wrong project milestone", "next", "other-project") + `}`,
	})

	// When
	_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:              "with milestone",
		ProjectMilestoneID: "project-milestone-id",
	})

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssues_resolves_target_once(t *testing.T) {
	// Given
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})}
	requests := []IssueCreateRequest{{Title: "First"}, {Title: "Second"}, {Title: "Third"}}

	// When
	outcomes, err := CreateIssues(context.Background(), recorder, matchingTarget(), requests, 3)

	// Then
	require.NoError(t, err)
	require.Len(t, outcomes, 3)
	for _, outcome := range outcomes {
		require.NoError(t, outcome.Err)
		require.Equal(t, "LIT-1", outcome.Issue.Identifier)
	}
	require.Equal(t, 1, recorder.countOf("Viewer"), "viewer should resolve once, not per row")
	// The fixture serves no "team" fixture, so the direct-lookup fast path errors
	// and falls through to the "Teams" scan; both run exactly once per
	// resolveTarget call, not once per row.
	require.Equal(t, 1, recorder.countOf("team"), "the direct team lookup should run once, not per row")
	require.Equal(t, 1, recorder.countOf("Teams"), "the team scan should run once, not per row")
}

func Test_CreateIssues_validates_each_distinct_label_once(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel":  `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})}
	requests := []IssueCreateRequest{
		{Title: "First", LabelIDs: []string{"label-id"}},
		{Title: "Second", LabelIDs: []string{"label-id"}},
		{Title: "Third", LabelIDs: []string{"label-id"}},
	}

	outcomes, err := CreateIssues(context.Background(), recorder, matchingTarget(), requests, 3)

	require.NoError(t, err)
	for _, outcome := range outcomes {
		require.NoError(t, outcome.Err)
	}
	require.Equal(t, 1, recorder.countOf("issueLabel"))
}

func Test_CreateIssues_reports_row_outcomes_in_order(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-9")) + `}}`,
	})
	requests := []IssueCreateRequest{{Title: "First"}, {Title: ""}, {Title: "Third"}}

	// When
	outcomes, err := CreateIssues(context.Background(), graphqlClient, matchingTarget(), requests, 3)

	// Then
	require.NoError(t, err)
	require.Len(t, outcomes, 3)
	require.Equal(t, 0, outcomes[0].Index)
	require.NoError(t, outcomes[0].Err)
	require.Equal(t, "LIT-9", outcomes[0].Issue.Identifier)
	require.Equal(t, 1, outcomes[1].Index)
	require.ErrorIs(t, outcomes[1].Err, ErrWriteInvalid)
	require.Equal(t, 2, outcomes[2].Index)
	require.NoError(t, outcomes[2].Err)
	require.Equal(t, "LIT-9", outcomes[2].Issue.Identifier)
}

func Test_CreateIssues_reports_row_error_for_invalid_due_date(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})
	requests := []IssueCreateRequest{{Title: "First", DueDate: "not-a-date"}}

	// When
	outcomes, err := CreateIssues(context.Background(), graphqlClient, matchingTarget(), requests, 1)

	// Then
	require.NoError(t, err)
	require.Len(t, outcomes, 1)
	require.ErrorIs(t, outcomes[0].Err, ErrWriteInvalid)
}

func Test_CreateIssues_defaults_and_clamps_concurrency(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{
		"IssueCreate": `{"issueCreate":{"success":true,"issue":` + issueJSON(b1IssueFixture("LIT-1")) + `}}`,
	})
	requests := []IssueCreateRequest{{Title: "First"}}

	// When / Then
	for _, concurrency := range []int{0, -1, 100} {
		outcomes, err := CreateIssues(context.Background(), graphqlClient, matchingTarget(), requests, concurrency)

		require.NoError(t, err)
		require.Len(t, outcomes, 1)
		require.NoError(t, outcomes[0].Err)
	}
}

func Test_CreateIssues_returns_error_when_target_resolution_fails(t *testing.T) {
	// Given
	graphqlClient := issueWriteFakeClient(map[string]string{})
	mismatchedTarget := config.Target{OrgID: "org-id", TeamKey: "OTHER", TeamID: "other-team-id"}

	// When
	outcomes, err := CreateIssues(
		context.Background(), graphqlClient, mismatchedTarget, []IssueCreateRequest{{Title: "First"}}, 1,
	)

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Nil(t, outcomes)
}
