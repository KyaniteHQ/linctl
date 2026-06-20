package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListAgentActivities returns AgentActivities visible to the authenticated user.
func ListAgentActivities(ctx context.Context, graphqlClient graphql.Client, limit int) (AgentActivityList, error) {
	result, err := agentActivities(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AgentActivityList{}, fmt.Errorf("list agent activities: %w", err)
	}

	summaries := make([]AgentActivitySummary, 0, len(result.AgentActivities.Nodes))
	for _, node := range result.AgentActivities.Nodes {
		summaries = append(summaries, agentActivitySummary(node.AgentActivitySummaryFields))
	}

	return AgentActivityList{
		AgentActivities: summaries,
		HasNextPage:     result.AgentActivities.PageInfo.HasNextPage,
		EndCursor:       result.AgentActivities.PageInfo.EndCursor,
	}, nil
}

// GetAgentActivityByID returns one AgentActivity by id.
func GetAgentActivityByID(ctx context.Context, graphqlClient graphql.Client, id string) (AgentActivitySummary, error) {
	result, err := agentActivity(ctx, graphqlClient, id)
	if err != nil {
		return AgentActivitySummary{}, fmt.Errorf("get agent activity %s: %w", id, err)
	}

	return agentActivitySummary(result.AgentActivity.AgentActivitySummaryFields), nil
}

func agentActivitySummary(fields AgentActivitySummaryFields) AgentActivitySummary {
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
	content AgentActivitySummaryFieldsContentAgentActivityContent,
) AgentActivityContentSummary {
	switch value := content.(type) {
	case *AgentActivitySummaryFieldsContentAgentActivityActionContent:
		return AgentActivityContentSummary{
			Type:      string(value.Type),
			Action:    value.Action,
			Parameter: value.Parameter,
			Result:    stringValue(value.Result),
		}
	case *AgentActivitySummaryFieldsContentAgentActivityElicitationContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	case *AgentActivitySummaryFieldsContentAgentActivityErrorContent:
		return AgentActivityContentSummary{
			Type:       string(value.Type),
			Body:       value.Body,
			ReasonCode: stringValue(value.ReasonCode),
		}
	case *AgentActivitySummaryFieldsContentAgentActivityPromptContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	case *AgentActivitySummaryFieldsContentAgentActivityResponseContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	case *AgentActivitySummaryFieldsContentAgentActivityThoughtContent:
		return AgentActivityContentSummary{Type: string(value.Type), Body: value.Body}
	default:
		return AgentActivityContentSummary{}
	}
}
