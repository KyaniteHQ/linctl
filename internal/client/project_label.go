package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ProjectLabelSummary is the compact project label model used by read-only commands.
type ProjectLabelSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Color         string `json:"color"`
	IsGroup       bool   `json:"is_group"`
	OrgID         string `json:"org_id,omitempty"`
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

//nolint:lll
type projectLabelsNode = gql.XProjectLabelsProjectLabelsProjectLabelConnectionNodesProjectLabel

//nolint:lll
type projectLabelChildrenNode = gql.XProjectLabel_childrenProjectLabelChildrenProjectLabelConnectionNodesProjectLabel

//nolint:lll
type projectLabelProjectsNode = gql.XProjectLabel_projectsProjectLabelProjectsProjectConnectionNodesProject

type projectLabelsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type projectLabelScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	labelID       string
	labelName     string
}

// ListProjectLabels returns visible Linear project labels.
func ListProjectLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectLabelList, error) {
	query := projectLabelsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list project labels", limit, defaultListPageSize,
		query.page,
		projectLabelsNodeSummary,
	)
	if err != nil {
		return ProjectLabelList{}, err
	}

	return ProjectLabelList{ProjectLabels: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetProjectLabelByID returns one Linear project label by id.
func GetProjectLabelByID(ctx context.Context, graphqlClient graphql.Client, id string) (ProjectLabelSummary, error) {
	result, err := gql.XProjectLabel(ctx, graphqlClient, id)
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
	query := &projectLabelScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list project label children "+id, limit, defaultListPageSize,
		query.children,
		projectLabelChildrenNodeSummary,
	)
	if err != nil {
		return ProjectLabelChildrenList{}, err
	}

	return ProjectLabelChildrenList{
		ProjectLabelID:   query.labelID,
		ProjectLabelName: query.labelName,
		ProjectLabels:    page.Items,
		HasNextPage:      page.HasNextPage,
		EndCursor:        page.EndCursor,
	}, nil
}

// ListProjectLabelProjects returns projects associated with one Linear project label.
func ListProjectLabelProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectLabelProjectsList, error) {
	query := &projectLabelScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list project label projects "+id, limit, defaultListPageSize,
		query.projects,
		projectLabelProjectsNodeSummary,
	)
	if err != nil {
		return ProjectLabelProjectsList{}, err
	}

	return ProjectLabelProjectsList{
		ProjectLabelID:   query.labelID,
		ProjectLabelName: query.labelName,
		Projects:         page.Items,
		HasNextPage:      page.HasNextPage,
		EndCursor:        page.EndCursor,
	}, nil
}

func (query projectLabelsQuery) page(pageSize int, after *string) ([]projectLabelsNode, bool, *string, error) {
	result, err := gql.XProjectLabels(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ProjectLabels.Nodes,
		result.ProjectLabels.PageInfo.HasNextPage,
		result.ProjectLabels.PageInfo.EndCursor,
		nil
}

func (query *projectLabelScopedQuery) children(
	pageSize int,
	after *string,
) ([]projectLabelChildrenNode, bool, *string, error) {
	result, err := gql.XProjectLabel_children(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.labelID = result.ProjectLabel.Id
	query.labelName = result.ProjectLabel.Name

	return result.ProjectLabel.Children.Nodes,
		result.ProjectLabel.Children.PageInfo.HasNextPage,
		result.ProjectLabel.Children.PageInfo.EndCursor,
		nil
}

func (query *projectLabelScopedQuery) projects(
	pageSize int,
	after *string,
) ([]projectLabelProjectsNode, bool, *string, error) {
	result, err := gql.XProjectLabel_projects(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.labelID = result.ProjectLabel.Id
	query.labelName = result.ProjectLabel.Name

	return result.ProjectLabel.Projects.Nodes,
		result.ProjectLabel.Projects.PageInfo.HasNextPage,
		result.ProjectLabel.Projects.PageInfo.EndCursor,
		nil
}

func projectLabelsNodeSummary(label projectLabelsNode) ProjectLabelSummary {
	return projectLabelSummary(label.ProjectLabelSummaryFields)
}

func projectLabelChildrenNodeSummary(label projectLabelChildrenNode) ProjectLabelSummary {
	return projectLabelSummary(label.ProjectLabelSummaryFields)
}

func projectLabelProjectsNodeSummary(project projectLabelProjectsNode) ProjectSummary {
	return projectSummaryFromFields(project.ProjectSummaryFields)
}

func projectLabelSummary(fields gql.ProjectLabelSummaryFields) ProjectLabelSummary {
	label := ProjectLabelSummary{
		ID:            fields.Id,
		Name:          fields.Name,
		Description:   stringValue(fields.Description),
		Color:         fields.Color,
		IsGroup:       fields.IsGroup,
		OrgID:         fields.Organization.Id,
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
