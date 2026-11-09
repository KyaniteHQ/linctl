package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// CreateAgentSessionOnIssue starts an AgentSession on an issue after resolving
// the issue and comparing it against the pinned target. The session belongs to
// the authenticated app user, so Linear sends the app its created event.
func CreateAgentSessionOnIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (AgentSessionSummary, error) {
	if issueID == "" {
		return AgentSessionSummary{}, requiredFieldError("issue")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return AgentSessionSummary{}, err
	}
	issue, err := guard.requireIssue(ctx, issueID)
	if err != nil {
		return AgentSessionSummary{}, err
	}

	created, err := gql.AgentSessionCreateOnIssue(
		ctx, guard.graphqlClient, LinearAgentSessionCreateOnIssueInput{IssueID: issue.ID},
	)
	if err != nil {
		return AgentSessionSummary{}, fmt.Errorf("create agent session on issue %s: %w", issueID, err)
	}
	if err := mutationSuccess(created.AgentSessionCreateOnIssue.Success, "agentSessionCreateOnIssue"); err != nil {
		return AgentSessionSummary{}, err
	}

	return agentSessionSummary(created.AgentSessionCreateOnIssue.AgentSession.AgentSessionSummaryFields), nil
}
