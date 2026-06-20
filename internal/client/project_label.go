package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ProjectLabelSummary is the compact project label model used by read-only commands.
type ProjectLabelSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Color         string `json:"color"`
	IsGroup       bool   `json:"is_group"`
	ParentID      string `json:"parent_id,omitempty"`
	ParentName    string `json:"parent_name,omitempty"`
	ParentColor   string `json:"parent_color,omitempty"`
	LastAppliedAt string `json:"last_applied_at,omitempty"`
	RetiredAt     string `json:"retired_at,omitempty"`
	ArchivedAt    string `json:"archived_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// ProjectLabelList is a page of Linear project labels.
type ProjectLabelList struct {
	ProjectLabels []ProjectLabelSummary `json:"project_labels"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

// ProjectLabelChildrenList is a page of child labels for one ProjectLabel.
type ProjectLabelChildrenList struct {
	ProjectLabelID   string                `json:"project_label_id"`
	ProjectLabelName string                `json:"project_label_name"`
	ProjectLabels    []ProjectLabelSummary `json:"project_labels"`
	HasNextPage      bool                  `json:"has_next_page"`
	EndCursor        *string               `json:"end_cursor,omitempty"`
}

// ProjectLabelProjectsList is a page of projects associated with one ProjectLabel.
type ProjectLabelProjectsList struct {
	ProjectLabelID   string           `json:"project_label_id"`
	ProjectLabelName string           `json:"project_label_name"`
	Projects         []ProjectSummary `json:"projects"`
	HasNextPage      bool             `json:"has_next_page"`
	EndCursor        *string          `json:"end_cursor,omitempty"`
}

// ListProjectLabels returns visible Linear project labels.
func ListProjectLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectLabelList, error) {
	result, err := projectLabels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectLabelList{}, fmt.Errorf("list project labels: %w", err)
	}

	labels := make([]ProjectLabelSummary, 0, len(result.ProjectLabels.Nodes))
	for _, label := range result.ProjectLabels.Nodes {
		labels = append(labels, projectLabelSummary(label.ProjectLabelSummaryFields))
	}

	return ProjectLabelList{
		ProjectLabels: labels,
		HasNextPage:   result.ProjectLabels.PageInfo.HasNextPage,
		EndCursor:     result.ProjectLabels.PageInfo.EndCursor,
	}, nil
}

// GetProjectLabelByID returns one Linear project label by id.
func GetProjectLabelByID(ctx context.Context, graphqlClient graphql.Client, id string) (ProjectLabelSummary, error) {
	result, err := projectLabel(ctx, graphqlClient, id)
	if err != nil {
		return ProjectLabelSummary{}, fmt.Errorf("get project label %s: %w", id, err)
	}

	return projectLabelSummary(result.ProjectLabel.ProjectLabelSummaryFields), nil
}

// ListProjectLabelChildren returns children for one Linear project label.
func ListProjectLabelChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectLabelChildrenList, error) {
	result, err := projectLabel_children(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectLabelChildrenList{}, fmt.Errorf("list project label children %s: %w", id, err)
	}

	labels := make([]ProjectLabelSummary, 0, len(result.ProjectLabel.Children.Nodes))
	for _, label := range result.ProjectLabel.Children.Nodes {
		labels = append(labels, projectLabelSummary(label.ProjectLabelSummaryFields))
	}

	return ProjectLabelChildrenList{
		ProjectLabelID:   result.ProjectLabel.Id,
		ProjectLabelName: result.ProjectLabel.Name,
		ProjectLabels:    labels,
		HasNextPage:      result.ProjectLabel.Children.PageInfo.HasNextPage,
		EndCursor:        result.ProjectLabel.Children.PageInfo.EndCursor,
	}, nil
}

// ListProjectLabelProjects returns projects associated with one Linear project label.
func ListProjectLabelProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectLabelProjectsList, error) {
	result, err := projectLabel_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectLabelProjectsList{}, fmt.Errorf("list project label projects %s: %w", id, err)
	}

	projects := make([]ProjectSummary, 0, len(result.ProjectLabel.Projects.Nodes))
	for _, project := range result.ProjectLabel.Projects.Nodes {
		projects = append(projects, projectSummaryFromFields(project.ProjectSummaryFields))
	}

	return ProjectLabelProjectsList{
		ProjectLabelID:   result.ProjectLabel.Id,
		ProjectLabelName: result.ProjectLabel.Name,
		Projects:         projects,
		HasNextPage:      result.ProjectLabel.Projects.PageInfo.HasNextPage,
		EndCursor:        result.ProjectLabel.Projects.PageInfo.EndCursor,
	}, nil
}

func projectLabelSummary(fields ProjectLabelSummaryFields) ProjectLabelSummary {
	label := ProjectLabelSummary{
		ID:            fields.Id,
		Name:          fields.Name,
		Description:   stringValue(fields.Description),
		Color:         fields.Color,
		IsGroup:       fields.IsGroup,
		LastAppliedAt: stringValue(fields.LastAppliedAt),
		RetiredAt:     stringValue(fields.RetiredAt),
		ArchivedAt:    stringValue(fields.ArchivedAt),
		CreatedAt:     fields.CreatedAt,
		UpdatedAt:     fields.UpdatedAt,
	}
	if fields.Parent != nil {
		label.ParentID = fields.Parent.Id
		label.ParentName = fields.Parent.Name
		label.ParentColor = fields.Parent.Color
	}

	return label
}
