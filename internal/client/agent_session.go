package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// AgentSessionSummary is the compact AgentSession model used by read-only commands.
type AgentSessionSummary struct {
	ID              string `json:"id"`
	SlugID          string `json:"slug_id"`
	Status          string `json:"status"`
	Summary         string `json:"summary,omitempty"`
	URL             string `json:"url,omitempty"`
	StartedAt       string `json:"started_at,omitempty"`
	EndedAt         string `json:"ended_at,omitempty"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	ArchivedAt      string `json:"archived_at,omitempty"`
	CreatorID       string `json:"creator_id,omitempty"`
	AppUserID       string `json:"app_user_id"`
	IssueID         string `json:"issue_id,omitempty"`
	IssueIdentifier string `json:"issue_identifier,omitempty"`
}

// AgentSessionList is a page of AgentSessions.
type AgentSessionList struct {
	AgentSessions []AgentSessionSummary `json:"agent_sessions"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

//nolint:lll
type agentSessionsNode = gql.XAgentSessionsAgentSessionsAgentSessionConnectionNodesAgentSession

type agentSessionsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListAgentSessions returns AgentSessions visible to the authenticated user.
func ListAgentSessions(ctx context.Context, graphqlClient graphql.Client, limit int) (AgentSessionList, error) {
	query := agentSessionsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list agent sessions", limit, defaultListPageSize,
		query.page,
		agentSessionsNodeSummary,
	)
	if err != nil {
		return AgentSessionList{}, err
	}

	return AgentSessionList{
		AgentSessions: page.Items,
		HasNextPage:   page.HasNextPage,
		EndCursor:     page.EndCursor,
	}, nil
}

// GetAgentSessionByID returns one AgentSession by id.
func GetAgentSessionByID(ctx context.Context, graphqlClient graphql.Client, id string) (AgentSessionSummary, error) {
	result, err := gql.XAgentSession(ctx, graphqlClient, id)
	if err != nil {
		return AgentSessionSummary{}, fmt.Errorf("get agent session %s: %w", id, err)
	}

	return agentSessionSummary(result.AgentSession.AgentSessionSummaryFields), nil
}

func (query agentSessionsQuery) page(pageSize int, after *string) ([]agentSessionsNode, bool, *string, error) {
	result, err := gql.XAgentSessions(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.AgentSessions.Nodes,
		result.AgentSessions.PageInfo.HasNextPage,
		result.AgentSessions.PageInfo.EndCursor,
		nil
}

func agentSessionsNodeSummary(node agentSessionsNode) AgentSessionSummary {
	return agentSessionSummary(node.AgentSessionSummaryFields)
}

func agentSessionSummary(fields gql.AgentSessionSummaryFields) AgentSessionSummary {
	summary := AgentSessionSummary{
		ID:         fields.Id,
		SlugID:     fields.SlugId,
		Status:     string(fields.Status),
		Summary:    stringValue(fields.Summary),
		URL:        stringValue(fields.Url),
		StartedAt:  stringValue(fields.StartedAt),
		EndedAt:    stringValue(fields.EndedAt),
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
		AppUserID:  fields.AppUser.Id,
	}
	if fields.Creator != nil {
		summary.CreatorID = fields.Creator.Id
	}
	if fields.Issue != nil {
		summary.IssueID = fields.Issue.Id
		summary.IssueIdentifier = fields.Issue.Identifier
	}

	return summary
}
