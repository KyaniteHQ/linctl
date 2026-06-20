package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// TemplateSummary is the compact Linear template model used by read-only commands.
type TemplateSummary struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Description     string  `json:"description,omitempty"`
	Icon            string  `json:"icon,omitempty"`
	Color           string  `json:"color,omitempty"`
	SortOrder       float64 `json:"sort_order"`
	LastAppliedAt   string  `json:"last_applied_at,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	ArchivedAt      string  `json:"archived_at,omitempty"`
	TeamID          string  `json:"team_id,omitempty"`
	TeamKey         string  `json:"team_key,omitempty"`
	TeamName        string  `json:"team_name,omitempty"`
	PipelineID      string  `json:"pipeline_id,omitempty"`
	CreatorID       string  `json:"creator_id,omitempty"`
	LastUpdatedByID string  `json:"last_updated_by_id,omitempty"`
	InheritedFromID string  `json:"inherited_from_id,omitempty"`
}

// TemplateList is a local page of Linear templates.
type TemplateList struct {
	Templates  []TemplateSummary `json:"templates"`
	TotalCount int               `json:"total_count"`
}

// ListTemplates returns visible Linear templates.
func ListTemplates(ctx context.Context, graphqlClient graphql.Client, limit int) (TemplateList, error) {
	result, err := templates(ctx, graphqlClient)
	if err != nil {
		return TemplateList{}, fmt.Errorf("list templates: %w", err)
	}

	summaries := make([]TemplateSummary, 0, len(result.Templates))
	for index, node := range result.Templates {
		if index >= limit {
			break
		}
		summaries = append(summaries, templateSummary(node.TemplateSummaryFields))
	}

	return TemplateList{
		Templates:  summaries,
		TotalCount: len(result.Templates),
	}, nil
}

// GetTemplateByID returns one Linear template by id.
func GetTemplateByID(ctx context.Context, graphqlClient graphql.Client, id string) (TemplateSummary, error) {
	result, err := template(ctx, graphqlClient, id)
	if err != nil {
		return TemplateSummary{}, fmt.Errorf("get template %s: %w", id, err)
	}

	return templateSummary(result.Template.TemplateSummaryFields), nil
}

func templateSummary(fields TemplateSummaryFields) TemplateSummary {
	summary := TemplateSummary{
		ID:            fields.Id,
		Name:          fields.Name,
		Type:          fields.Type,
		Description:   stringValue(fields.Description),
		Icon:          stringValue(fields.Icon),
		Color:         stringValue(fields.Color),
		SortOrder:     fields.SortOrder,
		LastAppliedAt: stringValue(fields.LastAppliedAt),
		CreatedAt:     fields.CreatedAt,
		UpdatedAt:     fields.UpdatedAt,
		ArchivedAt:    stringValue(fields.ArchivedAt),
	}
	if fields.Team != nil {
		summary.TeamID = fields.Team.Id
		summary.TeamKey = fields.Team.Key
		summary.TeamName = fields.Team.Name
	}
	if fields.Pipeline != nil {
		summary.PipelineID = fields.Pipeline.Id
	}
	if fields.Creator != nil {
		summary.CreatorID = fields.Creator.Id
	}
	if fields.LastUpdatedBy != nil {
		summary.LastUpdatedByID = fields.LastUpdatedBy.Id
	}
	if fields.InheritedFrom != nil {
		summary.InheritedFromID = fields.InheritedFrom.Id
	}

	return summary
}
