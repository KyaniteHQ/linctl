package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListProjectStatuses returns visible Linear project statuses.
func ListProjectStatuses(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectStatusList, error) {
	result, err := projectStatuses(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectStatusList{}, fmt.Errorf("list project statuses: %w", err)
	}

	statuses := make([]ProjectStatusSummary, 0, len(result.ProjectStatuses.Nodes))
	for _, status := range result.ProjectStatuses.Nodes {
		statuses = append(statuses, projectStatusSummary(status.ProjectStatusSummaryFields))
	}

	return ProjectStatusList{
		ProjectStatuses: statuses,
		HasNextPage:     result.ProjectStatuses.PageInfo.HasNextPage,
		EndCursor:       result.ProjectStatuses.PageInfo.EndCursor,
	}, nil
}

// GetProjectStatusByID returns one Linear project status by id.
func GetProjectStatusByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectStatusSummary, error) {
	result, err := projectStatus(ctx, graphqlClient, id)
	if err != nil {
		return ProjectStatusSummary{}, fmt.Errorf("get project status %s: %w", id, err)
	}

	return projectStatusSummary(result.ProjectStatus.ProjectStatusSummaryFields), nil
}

func projectStatusSummary(fields ProjectStatusSummaryFields) ProjectStatusSummary {
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
