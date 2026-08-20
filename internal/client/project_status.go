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
	HasNextPage     bool                   `json:"has_next_page"`
	EndCursor       *string                `json:"end_cursor,omitempty"`
}

// ProjectStatusProjectCount summarizes projects using one project status.
type ProjectStatusProjectCount struct {
	ProjectStatusID   string  `json:"project_status_id"`
	Count             float64 `json:"count"`
	PrivateCount      float64 `json:"private_count"`
	ArchivedTeamCount float64 `json:"archived_team_count"`
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

	return ProjectStatusList{ProjectStatuses: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
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

// GetProjectStatusProjectCount returns project counts for one Linear project status.
func GetProjectStatusProjectCount(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectStatusProjectCount, error) {
	result, err := gql.XProjectStatusProjectCount(ctx, graphqlClient, id)
	if err != nil {
		return ProjectStatusProjectCount{}, fmt.Errorf("get project status project count %s: %w", id, err)
	}

	count := result.ProjectStatusProjectCount
	return ProjectStatusProjectCount{
		ProjectStatusID:   id,
		Count:             count.Count,
		PrivateCount:      count.PrivateCount,
		ArchivedTeamCount: count.ArchivedTeamCount,
	}, nil
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
