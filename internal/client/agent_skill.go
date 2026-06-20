package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// AgentSkillSummary is the compact AgentSkill model used by read-only commands.
type AgentSkillSummary struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Body             string  `json:"body"`
	Description      string  `json:"description,omitempty"`
	SlugID           string  `json:"slug_id"`
	TeamID           string  `json:"team_id,omitempty"`
	Shared           bool    `json:"shared"`
	Icon             string  `json:"icon,omitempty"`
	Color            string  `json:"color,omitempty"`
	RecentUsageCount float64 `json:"recent_usage_count"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	ArchivedAt       string  `json:"archived_at,omitempty"`
	LastUsedAt       string  `json:"last_used_at,omitempty"`
	OwnerID          string  `json:"owner_id"`
	CreatorID        string  `json:"creator_id"`
	LastUpdatedByID  string  `json:"last_updated_by_id,omitempty"`
}

// AgentSkillList is a page of AgentSkills.
type AgentSkillList struct {
	AgentSkills []AgentSkillSummary `json:"agent_skills"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ListAgentSkills returns AgentSkills visible to the authenticated user.
func ListAgentSkills(ctx context.Context, graphqlClient graphql.Client, limit int) (AgentSkillList, error) {
	result, err := agentSkills(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AgentSkillList{}, fmt.Errorf("list agent skills: %w", err)
	}

	summaries := make([]AgentSkillSummary, 0, len(result.AgentSkills.Nodes))
	for _, node := range result.AgentSkills.Nodes {
		summaries = append(summaries, agentSkillSummary(node.AgentSkillSummaryFields))
	}

	return AgentSkillList{
		AgentSkills: summaries,
		HasNextPage: result.AgentSkills.PageInfo.HasNextPage,
		EndCursor:   result.AgentSkills.PageInfo.EndCursor,
	}, nil
}

// GetAgentSkillByID returns one AgentSkill by id.
func GetAgentSkillByID(ctx context.Context, graphqlClient graphql.Client, id string) (AgentSkillSummary, error) {
	result, err := agentSkill(ctx, graphqlClient, id)
	if err != nil {
		return AgentSkillSummary{}, fmt.Errorf("get agent skill %s: %w", id, err)
	}

	return agentSkillSummary(result.AgentSkill.AgentSkillSummaryFields), nil
}

func agentSkillSummary(fields AgentSkillSummaryFields) AgentSkillSummary {
	summary := AgentSkillSummary{
		ID:               fields.Id,
		Title:            fields.Title,
		Body:             fields.Body,
		Description:      stringValue(fields.Description),
		SlugID:           fields.SlugId,
		TeamID:           stringValue(fields.TeamId),
		Shared:           fields.Shared,
		Icon:             stringValue(fields.Icon),
		Color:            stringValue(fields.Color),
		RecentUsageCount: fields.RecentUsageCount,
		CreatedAt:        fields.CreatedAt,
		UpdatedAt:        fields.UpdatedAt,
		ArchivedAt:       stringValue(fields.ArchivedAt),
		LastUsedAt:       stringValue(fields.LastUsedAt),
		OwnerID:          fields.Owner.Id,
		CreatorID:        fields.Creator.Id,
	}
	if fields.LastUpdatedBy != nil {
		summary.LastUpdatedByID = fields.LastUpdatedBy.Id
	}

	return summary
}
