package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Templates   []TemplateSummary `json:"templates"`
	TotalCount  int               `json:"total_count"`
	HasNextPage bool              `json:"has_next_page,omitempty"`
	EndCursor   *string           `json:"end_cursor,omitempty"`
}

// HasMore reports whether more templates exist. TemplateList keeps its own
// page fields instead of embedding Page because its has_next_page tag carries
// omitempty, and embedding Page would start emitting the key when false.
func (list TemplateList) HasMore() bool {
	return list.HasNextPage
}

//nolint:lll
type organizationTemplatesNode = gql.XOrganization_templatesOrganizationTemplatesTemplateConnectionNodesTemplate

type organizationTemplatesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListOrganizationTemplates returns organization-wide Linear templates.
func ListOrganizationTemplates(ctx context.Context, graphqlClient graphql.Client, limit int) (TemplateList, error) {
	query := organizationTemplatesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list organization templates", limit, defaultListPageSize,
		query.page,
		organizationTemplatesNodeSummary,
	)
	if err != nil {
		return TemplateList{}, err
	}

	return TemplateList{
		Templates:   page.Items,
		TotalCount:  len(page.Items),
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListTemplates returns visible Linear templates.
// It cannot paginate: gql.XTemplates takes no arguments, fetches everything, and truncates client-side.
func ListTemplates(ctx context.Context, graphqlClient graphql.Client, limit int) (TemplateList, error) {
	result, err := gql.XTemplates(ctx, graphqlClient)
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
	result, err := gql.XTemplate(ctx, graphqlClient, id)
	if err != nil {
		return TemplateSummary{}, fmt.Errorf("get template %s: %w", id, err)
	}

	return templateSummary(result.Template.TemplateSummaryFields), nil
}

// TemplateDetail is the exact template model used by template content and
// guarded write receipts. Ordinary list and get commands stay on TemplateSummary.
type TemplateDetail struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Data       json.RawMessage `json:"data"`
	TeamID     string          `json:"team_id"`
	TeamKey    string          `json:"team_key"`
	TeamName   string          `json:"team_name"`
	PipelineID string          `json:"pipeline_id,omitempty"`
}

// GetTemplateDetail returns exact template data and scope by id.
func GetTemplateDetail(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (TemplateDetail, error) {
	if id == "" {
		return TemplateDetail{}, requiredFieldError("template id")
	}
	result, err := gql.XTemplateContent(ctx, graphqlClient, id)
	if err != nil {
		return TemplateDetail{}, fmt.Errorf("get template content %s: %w", id, err)
	}
	teamID, teamKey, teamName := "", "", ""
	if result.Template.Team != nil {
		teamID = result.Template.Team.Id
		teamKey = result.Template.Team.Key
		teamName = result.Template.Team.Name
	}
	pipelineID := ""
	if result.Template.Pipeline != nil {
		pipelineID = result.Template.Pipeline.Id
	}

	return newTemplateDetail(
		id,
		result.Template.Id,
		result.Template.Name,
		result.Template.Type,
		result.Template.TemplateData,
		teamID,
		teamKey,
		teamName,
		pipelineID,
	)
}

func newTemplateDetail(
	requestedID string,
	id string,
	name string,
	templateType string,
	data json.RawMessage,
	teamID string,
	teamKey string,
	teamName string,
	pipelineID string,
) (TemplateDetail, error) {
	if id == "" {
		return TemplateDetail{}, notFoundError("template %s", requestedID)
	}
	canonical, err := canonicalTemplateReadback(data)
	if err != nil {
		return TemplateDetail{}, fmt.Errorf("decode template %s data: %w", requestedID, err)
	}

	return TemplateDetail{
		ID:         id,
		Name:       name,
		Type:       templateType,
		Data:       canonical,
		TeamID:     teamID,
		TeamKey:    teamKey,
		TeamName:   teamName,
		PipelineID: pipelineID,
	}, nil
}

func (query organizationTemplatesQuery) page(
	pageSize int,
	after *string,
) ([]organizationTemplatesNode, bool, *string, error) {
	result, err := gql.XOrganization_templates(
		query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Organization.Templates.Nodes,
		result.Organization.Templates.PageInfo.HasNextPage,
		result.Organization.Templates.PageInfo.EndCursor,
		nil
}

func organizationTemplatesNodeSummary(node organizationTemplatesNode) TemplateSummary {
	return templateSummary(node.TemplateSummaryFields)
}

func templateSummary(fields gql.TemplateSummaryFields) TemplateSummary {
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
