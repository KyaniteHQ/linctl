package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// AgentActivityContentSummary is the compact content model for an AgentActivity.
type AgentActivityContentSummary struct {
	Type       string `json:"type"`
	Body       string `json:"body,omitempty"`
	Action     string `json:"action,omitempty"`
	Parameter  string `json:"parameter,omitempty"`
	Result     string `json:"result,omitempty"`
	ReasonCode string `json:"reason_code,omitempty"`
}

// AgentActivitySummary is the compact AgentActivity model used by read-only commands.
type AgentActivitySummary struct {
	ID              string                      `json:"id"`
	AgentSessionID  string                      `json:"agent_session_id"`
	Content         AgentActivityContentSummary `json:"content"`
	ContentType     string                      `json:"content_type"`
	Signal          string                      `json:"signal,omitempty"`
	Ephemeral       bool                        `json:"ephemeral"`
	SourceCommentID string                      `json:"source_comment_id,omitempty"`
	UserID          string                      `json:"user_id"`
	CreatedAt       string                      `json:"created_at"`
	UpdatedAt       string                      `json:"updated_at"`
	ArchivedAt      string                      `json:"archived_at,omitempty"`
}

// AgentActivityList is a page of AgentActivities.
type AgentActivityList struct {
	AgentActivities []AgentActivitySummary `json:"agent_activities"`
	HasNextPage     bool                   `json:"has_next_page"`
	EndCursor       *string                `json:"end_cursor,omitempty"`
}

//nolint:lll
type agentActivitiesNode = gql.XAgentActivitiesAgentActivitiesAgentActivityConnectionNodesAgentActivity

type agentActivitiesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListAgentActivities returns AgentActivities visible to the authenticated user.
func ListAgentActivities(ctx context.Context, graphqlClient graphql.Client, limit int) (AgentActivityList, error) {
	query := agentActivitiesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list agent activities", limit, defaultListPageSize,
		query.page,
		agentActivitiesNodeSummary,
	)
	if err != nil {
		return AgentActivityList{}, err
	}

	return AgentActivityList{
		AgentActivities: page.Items,
		HasNextPage:     page.HasNextPage,
		EndCursor:       page.EndCursor,
	}, nil
}

// GetAgentActivityByID returns one AgentActivity by id.
func GetAgentActivityByID(ctx context.Context, graphqlClient graphql.Client, id string) (AgentActivitySummary, error) {
	result, err := gql.XAgentActivity(ctx, graphqlClient, id)
	if err != nil {
		return AgentActivitySummary{}, fmt.Errorf("get agent activity %s: %w", id, err)
	}

	return agentActivitySummary(result.AgentActivity.AgentActivitySummaryFields), nil
}

func (query agentActivitiesQuery) page(pageSize int, after *string) ([]agentActivitiesNode, bool, *string, error) {
	result, err := gql.XAgentActivities(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.AgentActivities.Nodes,
		result.AgentActivities.PageInfo.HasNextPage,
		result.AgentActivities.PageInfo.EndCursor,
		nil
}

func agentActivitiesNodeSummary(node agentActivitiesNode) AgentActivitySummary {
	return agentActivitySummary(node.AgentActivitySummaryFields)
}

func agentActivitySummary(fields gql.AgentActivitySummaryFields) AgentActivitySummary {
	summary := AgentActivitySummary{
		ID:             fields.Id,
		AgentSessionID: fields.AgentSession.Id,
		Content:        agentActivityContentSummary(fields.Content),
		Ephemeral:      fields.Ephemeral,
		UserID:         fields.User.Id,
		CreatedAt:      fields.CreatedAt,
		UpdatedAt:      fields.UpdatedAt,
		ArchivedAt:     stringValue(fields.ArchivedAt),
	}
	summary.ContentType = summary.Content.Type
	if fields.Signal != nil {
		summary.Signal = string(*fields.Signal)
	}
	if fields.SourceComment != nil {
		summary.SourceCommentID = fields.SourceComment.Id
	}

	return summary
}

func agentActivityContentSummary(
	content gql.AgentActivitySummaryFieldsContentAgentActivityContent,
) AgentActivityContentSummary {
	switch value := content.(type) {
	case *gql.AgentActivitySummaryFieldsContentAgentActivityActionContent:
		return AgentActivityContentSummary{
			Type:      string(value.Type),
			Action:    value.Action,
			Parameter: value.Parameter,
			Result:    stringValue(value.Result),
		}
	case *gql.AgentActivitySummaryFieldsContentAgentActivityElicitationContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	case *gql.AgentActivitySummaryFieldsContentAgentActivityErrorContent:
		return AgentActivityContentSummary{
			Type:       string(value.Type),
			Body:       value.Body,
			ReasonCode: stringValue(value.ReasonCode),
		}
	case *gql.AgentActivitySummaryFieldsContentAgentActivityPromptContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	case *gql.AgentActivitySummaryFieldsContentAgentActivityResponseContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	case *gql.AgentActivitySummaryFieldsContentAgentActivityThoughtContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	default:
		return AgentActivityContentSummary{}
	}
}
