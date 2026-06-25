package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func agentSessionFieldsJSON(withCreatorAndIssue bool) string {
	creator := "null"
	issue := "null"
	if withCreatorAndIssue {
		creator = `{"id":"creator-id"}`
		issue = `{"id":"issue-id","identifier":"LIT-1"}`
	}

	return `{
		"id":"agent-session-id",
		"slugId":"session-slug",
		"status":"active",
		"summary":"did work",
		"url":"https://linear.app/kyanite/agent/session-slug",
		"startedAt":"2026-06-19T12:00:00Z",
		"endedAt":null,
		"createdAt":"2026-06-19T11:00:00Z",
		"updatedAt":"2026-06-19T12:30:00Z",
		"archivedAt":null,
		"creator":` + creator + `,
		"appUser":{"id":"app-user-id"},
		"issue":` + issue + `
	}`
}

func Test_ListAgentSessions_returns_compact_sessions(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"agentSessions": `{"agentSessions":{"nodes":[` + agentSessionFieldsJSON(true) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"cursor-1"}}}`,
	})

	sessions, err := ListAgentSessions(context.Background(), graphqlClient, 2)

	require.NoError(t, err)
	require.True(t, sessions.HasNextPage)
	require.Len(t, sessions.AgentSessions, 1)
	require.Equal(t, "agent-session-id", sessions.AgentSessions[0].ID)
	require.Equal(t, "active", sessions.AgentSessions[0].Status)
	require.Equal(t, "creator-id", sessions.AgentSessions[0].CreatorID)
	require.Equal(t, "LIT-1", sessions.AgentSessions[0].IssueIdentifier)
	require.Equal(t, "app-user-id", sessions.AgentSessions[0].AppUserID)
}

func Test_GetAgentSessionByID_returns_session_without_creator_or_issue(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"agentSession": `{"agentSession":` + agentSessionFieldsJSON(false) + `}`,
	})

	session, err := GetAgentSessionByID(context.Background(), graphqlClient, "agent-session-id")

	require.NoError(t, err)
	require.Equal(t, "agent-session-id", session.ID)
	require.Empty(t, session.CreatorID)
	require.Empty(t, session.IssueIdentifier)
}

func Test_ListAgentSessions_wraps_read_error(t *testing.T) {
	_, err := ListAgentSessions(context.Background(), fakeGraphQLClient(map[string]string{}), 1)

	require.Error(t, err)
}

func Test_GetAgentSessionByID_wraps_read_error(t *testing.T) {
	_, err := GetAgentSessionByID(context.Background(), fakeGraphQLClient(map[string]string{}), "agent-session-id")

	require.Error(t, err)
}
