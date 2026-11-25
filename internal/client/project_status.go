package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ProjectStatusSummary is the compact project status model used by read-only commands.
type ProjectStatusSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Position    string `json:"position"`
	ArchivedAt  string `json:"archived_at,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ProjectStatusList is a page of Linear project statuses.
type ProjectStatusList struct {
	ProjectStatuses []ProjectStatusSummary `json:"project_statuses"`
	Page
}

//nolint:lll
type projectStatusesNode = gql.XProjectStatusesProjectStatusesProjectStatusConnectionNodesProjectStatus

type projectStatusesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListProjectStatuses returns visible Linear project statuses.
func ListProjectStatuses(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectStatusList, error) {
	query := projectStatusesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list project statuses", limit, defaultListPageSize,
		query.page,
		projectStatusNodeSummary,
	)
	if err != nil {
		return ProjectStatusList{}, err
	}

	return ProjectStatusList{ProjectStatuses: page.Items, Page: page.Page}, nil
}

func (query projectStatusesQuery) page(
	pageSize int,
	after *string,
) ([]projectStatusesNode, bool, *string, error) {
	result, err := gql.XProjectStatuses(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ProjectStatuses.Nodes,
		result.ProjectStatuses.PageInfo.HasNextPage,
		result.ProjectStatuses.PageInfo.EndCursor,
		nil
}

func projectStatusNodeSummary(status projectStatusesNode) ProjectStatusSummary {
	return projectStatusSummary(status.ProjectStatusSummaryFields)
}

// GetProjectStatusByID returns one Linear project status by id.
func GetProjectStatusByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectStatusSummary, error) {
	result, err := gql.XProjectStatus(ctx, graphqlClient, id)
	if err != nil {
		return ProjectStatusSummary{}, fmt.Errorf("get project status %s: %w", id, err)
	}

	return projectStatusSummary(result.ProjectStatus.ProjectStatusSummaryFields), nil
}

func projectStatusSummary(fields gql.ProjectStatusSummaryFields) ProjectStatusSummary {
	return ProjectStatusSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		Description: stringValue(fields.Description),
		Type:        string(fields.Type),
		Color:       fields.Color,
		Position:    fmt.Sprintf("%.2f", fields.Position),
		ArchivedAt:  stringValue(fields.ArchivedAt),
		CreatedAt:   fields.CreatedAt,
		UpdatedAt:   fields.UpdatedAt,
	}
}
