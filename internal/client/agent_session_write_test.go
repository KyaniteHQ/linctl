package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func agentSessionCreatePayloads() map[string]string {
	return map[string]string{
		"issue": pinnedIssueJSON(),
		"AgentSessionCreateOnIssue": `{"agentSessionCreateOnIssue":{"success":true,"agentSession":` +
			agentSessionJSON(`{"id":"issue-id","identifier":"LIT-1"}`) + `}}`,
	}
}

func Test_CreateAgentSessionOnIssue_sends_the_resolved_issue_id(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(agentSessionCreatePayloads())}

	session, err := CreateAgentSessionOnIssue(context.Background(), recorder, matchingTarget(), "LIT-1")

	require.NoError(t, err)
	require.Equal(t, "agent-session-id", session.ID)
	require.JSONEq(t, `{"input":{"issueId":"issue-id"}}`,
		string(recorder.variablesFor(t, "AgentSessionCreateOnIssue")))
}

func Test_CreateAgentSessionOnIssue_requires_an_issue(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateAgentSessionOnIssue(context.Background(), recorder, matchingTarget(), "")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "issue is required")
	require.False(t, recorder.sentOperation("AgentSessionCreateOnIssue"))
}

func Test_CreateAgentSessionOnIssue_refuses_when_target_unresolved(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateAgentSessionOnIssue(context.Background(), recorder, config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, "LIT-1")

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("AgentSessionCreateOnIssue"))
}

func Test_CreateAgentSessionOnIssue_refuses_an_issue_outside_the_pinned_target(t *testing.T) {
	otherTeam := issueJSONWithTeam(issueFixture{
		Identifier: "OPS-1", Title: "Other issue", StateID: "state-id", State: "Todo", StateType: "unstarted",
	}, "other-team-id", "OTHER")
	otherProject := issueJSON(issueFixture{
		Identifier: "LIT-2", Title: "Other project", ProjectID: "other-project-id", Project: "other",
		StateID: "state-id", State: "Todo", StateType: "unstarted",
	})
	noProject := issueJSON(issueFixture{
		Identifier: "LIT-3", Title: "No project", StateID: "state-id", State: "Todo", StateType: "unstarted",
	})
	tests := []struct {
		name  string
		issue string
	}{
		{name: "another team", issue: otherTeam},
		{name: "same team, another project", issue: otherProject},
		{name: "same team, no project", issue: noProject},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payloads := agentSessionCreatePayloads()
			payloads["issue"] = `{"issue":` + test.issue + `}`
			recorder := &mutationRecordingClient{inner: issueWriteFakeClient(payloads)}

			_, err := CreateAgentSessionOnIssue(context.Background(), recorder, matchingTarget(), "ANY-1")

			require.ErrorIs(t, err, ErrTargetMismatch)
			require.False(t, recorder.sentOperation("AgentSessionCreateOnIssue"))
		})
	}
}

func Test_CreateAgentSessionOnIssue_wraps_mutation_error(t *testing.T) {
	payloads := agentSessionCreatePayloads()
	delete(payloads, "AgentSessionCreateOnIssue")

	_, err := CreateAgentSessionOnIssue(
		context.Background(), issueWriteFakeClient(payloads), matchingTarget(), "LIT-1",
	)

	require.ErrorContains(t, err, "create agent session on issue LIT-1")
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateAgentSessionOnIssue_fails_when_mutation_reports_no_success(t *testing.T) {
	payloads := agentSessionCreatePayloads()
	payloads["AgentSessionCreateOnIssue"] = `{"agentSessionCreateOnIssue":{"success":false,"agentSession":` +
		agentSessionJSON(`{"id":"issue-id","identifier":"LIT-1"}`) + `}}`

	_, err := CreateAgentSessionOnIssue(
		context.Background(), issueWriteFakeClient(payloads), matchingTarget(), "LIT-1",
	)

	require.ErrorIs(t, err, ErrMutationFailed)
}
