package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func agentSessionJSON(issue string) string {
	return `{
		"id": "agent-session-id",
		"slugId": "session-slug",
		"status": "active",
		"summary": null,
		"url": null,
		"startedAt": null,
		"endedAt": null,
		"createdAt": "2026-06-19T11:00:00Z",
		"updatedAt": "2026-06-19T12:30:00Z",
		"archivedAt": null,
		"creator": null,
		"appUser": {"id": "app-user-id"},
		"issue": ` + issue + `
	}`
}

func createdAgentActivityJSON() string {
	return `{
		"id": "created-activity-id",
		"createdAt": "2026-06-19T12:00:00Z",
		"updatedAt": "2026-06-19T12:01:00Z",
		"archivedAt": null,
		"signal": null,
		"ephemeral": false,
		"agentSession": {"id": "agent-session-id"},
		"sourceComment": null,
		"user": {"id": "user-id"},
		"content": {"__typename": "AgentActivityThoughtContent", "type": "thought", "body": "Reading the plan"}
	}`
}

func pinnedIssueJSON() string {
	return `{"issue":` + issueJSON(issueFixture{
		Identifier: "LIT-1", Title: "Pinned issue", ProjectID: "project-id", Project: "fixture",
		StateID: "state-id", State: "Todo", StateType: "unstarted",
	}) + `}`
}

func agentActivityCreatePayloads() map[string]string {
	return map[string]string{
		"agentSession":        `{"agentSession":` + agentSessionJSON(`{"id":"issue-id","identifier":"LIT-1"}`) + `}`,
		"issue":               pinnedIssueJSON(),
		"AgentActivityCreate": `{"agentActivityCreate":{"success":true,"agentActivity":` + createdAgentActivityJSON() + `}}`,
	}
}

func thoughtRequest() AgentActivityCreateRequest {
	return AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Type: "thought", Body: "Reading the plan"}
}

func Test_CreateAgentActivity_emits_a_thought_into_a_pinned_session(t *testing.T) {
	activity, err := CreateAgentActivity(
		context.Background(), issueWriteFakeClient(agentActivityCreatePayloads()), matchingTarget(), thoughtRequest(),
	)

	require.NoError(t, err)
	require.Equal(t, "created-activity-id", activity.ID)
	require.Equal(t, "thought", activity.ContentType)
}

func Test_CreateAgentActivity_sends_every_field_it_was_given(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(agentActivityCreatePayloads())}

	_, err := CreateAgentActivity(context.Background(), recorder, matchingTarget(), AgentActivityCreateRequest{
		AgentSessionID: "agent-session-id", Type: "action", Action: "read_file", Parameter: "README.md",
		Result: "Read file", Signal: "continue", Ephemeral: true,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"input": {
			"agentSessionId": "agent-session-id",
			"content": {"type": "action", "action": "read_file", "parameter": "README.md", "result": "Read file"},
			"signal": "continue",
			"ephemeral": true
		}
	}`, string(recorder.variablesFor(t, "AgentActivityCreate")))
}

func Test_CreateAgentActivity_omits_the_optional_fields_it_was_not_given(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(agentActivityCreatePayloads())}

	_, err := CreateAgentActivity(context.Background(), recorder, matchingTarget(), thoughtRequest())

	require.NoError(t, err)
	require.JSONEq(t, `{
		"input": {
			"agentSessionId": "agent-session-id",
			"content": {"type": "thought", "body": "Reading the plan"}
		}
	}`, string(recorder.variablesFor(t, "AgentActivityCreate")))
}

func Test_CreateAgentActivity_validates_the_request_before_any_request_is_sent(t *testing.T) {
	tests := []struct {
		name     string
		request  AgentActivityCreateRequest
		contains string
	}{
		{name: "missing session", request: AgentActivityCreateRequest{Type: "thought", Body: "x"}, contains: "agent session id is required"},
		{name: "missing type", request: AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Body: "x"}, contains: "type is required"},
		{name: "unknown type", request: AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Type: "prompt", Body: "x"}, contains: "type must be one of thought, elicitation, response, error, action"},
		{name: "unknown signal", request: AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Type: "thought", Body: "x", Signal: "pause"}, contains: "signal must be one of auth, continue, select, stop"},
		{name: "missing body", request: AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Type: "response"}, contains: "body is required"},
		{name: "missing action", request: AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Type: "action", Parameter: "x"}, contains: "action is required"},
		{name: "missing parameter", request: AgentActivityCreateRequest{AgentSessionID: "agent-session-id", Type: "action", Action: "x"}, contains: "parameter is required"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

			_, err := CreateAgentActivity(context.Background(), recorder, matchingTarget(), test.request)

			require.ErrorIs(t, err, ErrWriteInvalid)
			require.ErrorContains(t, err, test.contains)
			require.False(t, recorder.sentOperation("AgentActivityCreate"))
		})
	}
}

func Test_CreateAgentActivity_refuses_when_target_unresolved(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateAgentActivity(context.Background(), recorder, config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, thoughtRequest())

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("AgentActivityCreate"))
}

func Test_CreateAgentActivity_wraps_session_lookup_error(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateAgentActivity(context.Background(), recorder, matchingTarget(), thoughtRequest())

	require.ErrorContains(t, err, "get agent session agent-session-id")
	require.False(t, recorder.sentOperation("AgentActivityCreate"))
}

func Test_CreateAgentActivity_refuses_a_session_without_an_issue(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"agentSession": `{"agentSession":` + agentSessionJSON("null") + `}`,
	})}

	_, err := CreateAgentActivity(context.Background(), recorder, matchingTarget(), thoughtRequest())

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "not attached to an issue")
	require.False(t, recorder.sentOperation("AgentActivityCreate"))
}

func Test_CreateAgentActivity_refuses_a_session_on_another_teams_issue(t *testing.T) {
	payloads := agentActivityCreatePayloads()
	payloads["issue"] = `{"issue":` + issueJSONWithTeam(issueFixture{
		Identifier: "OPS-1", Title: "Other issue", StateID: "state-id", State: "Todo", StateType: "unstarted",
	}, "other-team-id", "OTHER") + `}`
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(payloads)}

	_, err := CreateAgentActivity(context.Background(), recorder, matchingTarget(), thoughtRequest())

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("AgentActivityCreate"))
}

func Test_CreateAgentActivity_wraps_mutation_error(t *testing.T) {
	payloads := agentActivityCreatePayloads()
	delete(payloads, "AgentActivityCreate")

	_, err := CreateAgentActivity(
		context.Background(), issueWriteFakeClient(payloads), matchingTarget(), thoughtRequest(),
	)

	require.ErrorContains(t, err, "create agent activity")
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateAgentActivity_fails_when_mutation_reports_no_success(t *testing.T) {
	payloads := agentActivityCreatePayloads()
	payloads["AgentActivityCreate"] = `{"agentActivityCreate":{"success":false,"agentActivity":` + createdAgentActivityJSON() + `}}`

	_, err := CreateAgentActivity(
		context.Background(), issueWriteFakeClient(payloads), matchingTarget(), thoughtRequest(),
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}
