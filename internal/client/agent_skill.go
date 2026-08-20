package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Page
}

//nolint:lll
type agentSkillsNode = gql.XAgentSkillsAgentSkillsAgentSkillConnectionNodesAgentSkill

type agentSkillsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListAgentSkills returns AgentSkills visible to the authenticated user.
func ListAgentSkills(ctx context.Context, graphqlClient graphql.Client, limit int) (AgentSkillList, error) {
	query := agentSkillsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list agent skills", limit, defaultListPageSize,
		query.page,
		agentSkillsNodeSummary,
	)
	if err != nil {
		return AgentSkillList{}, err
	}

	return AgentSkillList{AgentSkills: page.Items, Page: page.Page}, nil
}

// GetAgentSkillByID returns one AgentSkill by id.
func GetAgentSkillByID(ctx context.Context, graphqlClient graphql.Client, id string) (AgentSkillSummary, error) {
	result, err := gql.XAgentSkill(ctx, graphqlClient, id)
	if err != nil {
		return AgentSkillSummary{}, fmt.Errorf("get agent skill %s: %w", id, err)
	}

	return agentSkillSummary(result.AgentSkill.AgentSkillSummaryFields), nil
}

func (query agentSkillsQuery) page(pageSize int, after *string) ([]agentSkillsNode, bool, *string, error) {
	result, err := gql.XAgentSkills(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.AgentSkills.Nodes,
		result.AgentSkills.PageInfo.HasNextPage,
		result.AgentSkills.PageInfo.EndCursor,
		nil
}

func agentSkillsNodeSummary(node agentSkillsNode) AgentSkillSummary {
	return agentSkillSummary(node.AgentSkillSummaryFields)
}

func agentSkillSummary(fields gql.AgentSkillSummaryFields) AgentSkillSummary {
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
