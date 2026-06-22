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
	graphqlClient := issueWriteFakeClient(map[string]string{
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
	})

	// When
	issue, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{
		Title:     "typed",
		StateType: "unstarted",
		Priority:  "2",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-3", issue.Identifier)
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
	})

	// When
	issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
		ID:        "LIT-1",
		StateType: "completed",
		Priority:  "1",
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "completed", issue.StateType)
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
	require.Error(t, err)
	require.Contains(t, err.Error(), "list started workflow states")
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
