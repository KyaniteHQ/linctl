package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectMilestoneSummary is one milestone within a project.
type ProjectMilestoneSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	TargetDate  string  `json:"target_date,omitempty"`
	Status      string  `json:"status"`
	Progress    float64 `json:"progress"`
	SortOrder   float64 `json:"sort_order"`
}

// ProjectMilestoneList is a page of project milestones.
type ProjectMilestoneList struct {
	ProjectID   string                    `json:"project_id"`
	ProjectName string                    `json:"project_name"`
	Milestones  []ProjectMilestoneSummary `json:"milestones"`
	HasNextPage bool                      `json:"has_next_page"`
	EndCursor   *string                   `json:"end_cursor,omitempty"`
}

// ProjectMilestoneDetail is a ProjectMilestone with its parent project.
type ProjectMilestoneDetail struct {
	Summary ProjectMilestoneSummary `json:"summary"`
	Project ProjectSummary          `json:"project"`
}

// ProjectMilestoneIssueList is a page of issues associated with one ProjectMilestone.
type ProjectMilestoneIssueList struct {
	ProjectMilestoneID   string         `json:"project_milestone_id"`
	ProjectMilestoneName string         `json:"project_milestone_name"`
	Issues               []IssueSummary `json:"issues"`
	HasNextPage          bool           `json:"has_next_page"`
	EndCursor            *string        `json:"end_cursor,omitempty"`
}

// ProjectMilestoneCreateRequest describes a guarded ProjectMilestone create.
type ProjectMilestoneCreateRequest struct {
	ProjectID   string
	Name        string
	Description string
	TargetDate  string
}

// ProjectMilestoneUpdateRequest describes a guarded ProjectMilestone update.
type ProjectMilestoneUpdateRequest struct {
	ID          string
	Name        string
	Description string
	TargetDate  string
}

func projectMilestoneSummary(milestone gql.ProjectMilestoneSummaryFields) ProjectMilestoneSummary {
	return ProjectMilestoneSummary{
		ID:          milestone.Id,
		Name:        milestone.Name,
		Description: stringValue(milestone.Description),
		TargetDate:  stringValue(milestone.TargetDate),
		Status:      string(milestone.Status),
		Progress:    milestone.Progress,
		SortOrder:   milestone.SortOrder,
	}
}

//nolint:lll
type projectMilestonesNode = gql.XProjectMilestonesProjectMilestonesProjectMilestoneConnectionNodesProjectMilestone

//nolint:lll
type projectScopedMilestonesNode = gql.XProject_projectMilestonesProjectProjectMilestonesProjectMilestoneConnectionNodesProjectMilestone

//nolint:lll
type projectMilestoneIssuesNode = gql.XProjectMilestone_issuesProjectMilestoneIssuesIssueConnectionNodesIssue

type projectMilestonesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type projectScopedMilestonesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	projectID     string
	projectName   string
}

type projectMilestoneIssuesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	milestoneID   string
	milestoneName string
}

// ListProjectMilestones returns milestones for one project.
func ListProjectMilestones(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMilestoneList, error) {
	query := &projectScopedMilestonesQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list project milestones "+id, limit, defaultListPageSize,
		query.page,
		projectScopedMilestoneNodeSummary,
	)
	if err != nil {
		return ProjectMilestoneList{}, err
	}

	return ProjectMilestoneList{
		ProjectID:   query.projectID,
		ProjectName: query.projectName,
		Milestones:  page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListAllProjectMilestones returns project milestones visible to the authenticated user.
func ListAllProjectMilestones(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (ProjectMilestoneList, error) {
	query := projectMilestonesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list project milestones", limit, defaultListPageSize,
		query.page,
		projectMilestoneNodeSummary,
	)
	if err != nil {
		return ProjectMilestoneList{}, err
	}

	return ProjectMilestoneList{Milestones: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListProjectMilestoneIssues returns issues associated with one ProjectMilestone.
func ListProjectMilestoneIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMilestoneIssueList, error) {
	query := &projectMilestoneIssuesQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list project milestone issues "+id, limit, defaultListPageSize,
		query.page,
		projectMilestoneIssueNodeSummary,
	)
	if err != nil {
		return ProjectMilestoneIssueList{}, err
	}

	return ProjectMilestoneIssueList{
		ProjectMilestoneID:   query.milestoneID,
		ProjectMilestoneName: query.milestoneName,
		Issues:               page.Items,
		HasNextPage:          page.HasNextPage,
		EndCursor:            page.EndCursor,
	}, nil
}

func (query *projectScopedMilestonesQuery) page(
	pageSize int,
	after *string,
) ([]projectScopedMilestonesNode, bool, *string, error) {
	result, err := gql.XProject_projectMilestones(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.projectID = result.Project.Id
	query.projectName = result.Project.Name

	return result.Project.ProjectMilestones.Nodes,
		result.Project.ProjectMilestones.PageInfo.HasNextPage,
		result.Project.ProjectMilestones.PageInfo.EndCursor,
		nil
}

func (query projectMilestonesQuery) page(
	pageSize int,
	after *string,
) ([]projectMilestonesNode, bool, *string, error) {
	result, err := gql.XProjectMilestones(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ProjectMilestones.Nodes,
		result.ProjectMilestones.PageInfo.HasNextPage,
		result.ProjectMilestones.PageInfo.EndCursor,
		nil
}

func (query *projectMilestoneIssuesQuery) page(
	pageSize int,
	after *string,
) ([]projectMilestoneIssuesNode, bool, *string, error) {
	result, err := gql.XProjectMilestone_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.milestoneID = result.ProjectMilestone.Id
	query.milestoneName = result.ProjectMilestone.Name

	return result.ProjectMilestone.Issues.Nodes,
		result.ProjectMilestone.Issues.PageInfo.HasNextPage,
		result.ProjectMilestone.Issues.PageInfo.EndCursor,
		nil
}

func projectScopedMilestoneNodeSummary(milestone projectScopedMilestonesNode) ProjectMilestoneSummary {
	return projectMilestoneSummary(milestone.ProjectMilestoneSummaryFields)
}

func projectMilestoneNodeSummary(milestone projectMilestonesNode) ProjectMilestoneSummary {
	return projectMilestoneSummary(milestone.ProjectMilestoneSummaryFields)
}

func projectMilestoneIssueNodeSummary(node projectMilestoneIssuesNode) IssueSummary {
	return issueSummaryFromFields(node.IssueSummaryFields)
}

// GetProjectMilestoneByID returns one ProjectMilestone by Linear id.
func GetProjectMilestoneByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectMilestoneSummary, error) {
	detail, err := GetProjectMilestoneDetail(ctx, graphqlClient, id)
	if err != nil {
		return ProjectMilestoneSummary{}, fmt.Errorf("get project milestone %s: %w", id, err)
	}

	return detail.Summary, nil
}

// GetProjectMilestoneDetail returns one ProjectMilestone and its parent project.
func GetProjectMilestoneDetail(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectMilestoneDetail, error) {
	milestone, err := gql.XProjectMilestone(ctx, graphqlClient, id)
	if err != nil {
		return ProjectMilestoneDetail{}, fmt.Errorf("get project milestone %s: %w", id, err)
	}

	return ProjectMilestoneDetail{
		Summary: projectMilestoneSummary(milestone.ProjectMilestone.ProjectMilestoneSummaryFields),
		Project: projectSummaryFromFields(milestone.ProjectMilestone.Project.ProjectSummaryFields),
	}, nil
}

// CreateProjectMilestone creates a ProjectMilestone after resolving and comparing its project.
func CreateProjectMilestone(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectMilestoneCreateRequest,
) (ProjectMilestoneSummary, error) {
	if request.ProjectID == "" {
		return ProjectMilestoneSummary{}, requiredFieldError("project id")
	}
	if request.Name == "" {
		return ProjectMilestoneSummary{}, requiredFieldError("name")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectMilestoneSummary{}, err
	}

	return guard.createProjectMilestone(ctx, request)
}

func (guard *guardedClient) createProjectMilestone(
	ctx context.Context,
	request ProjectMilestoneCreateRequest,
) (ProjectMilestoneSummary, error) {
	if err := guard.requireProject(ctx, request.ProjectID); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	created, err := gql.ProjectMilestoneCreate(ctx, guard.graphqlClient, LinearProjectMilestoneCreateInput{
		ProjectID:   request.ProjectID,
		Name:        request.Name,
		Description: optionalString(request.Description),
		TargetDate:  optionalString(request.TargetDate),
	})
	if err != nil {
		return ProjectMilestoneSummary{}, fmt.Errorf("create project milestone: %w", err)
	}
	if err := mutationSuccess(created.ProjectMilestoneCreate.Success, "projectMilestoneCreate"); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	return projectMilestoneSummary(
		created.ProjectMilestoneCreate.ProjectMilestone.ProjectMilestoneSummaryFields,
	), nil
}

// UpdateProjectMilestone updates a ProjectMilestone after resolving and comparing its project.
//
//nolint:dupl // Mirrors UpdateProjectLabel's resolve-then-mutate shape; the guard target and mutation differ.
func UpdateProjectMilestone(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectMilestoneUpdateRequest,
) (ProjectMilestoneSummary, error) {
	if err := validateProjectMilestoneUpdateRequest(request); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectMilestoneSummary{}, err
	}

	return guard.updateProjectMilestone(ctx, request)
}

func (guard *guardedClient) updateProjectMilestone(
	ctx context.Context,
	request ProjectMilestoneUpdateRequest,
) (ProjectMilestoneSummary, error) {
	if err := guard.requireProjectMilestone(ctx, request.ID); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	updated, err := gql.ProjectMilestoneUpdate(ctx, guard.graphqlClient, request.ID, LinearProjectMilestoneUpdateInput{
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		TargetDate:  optionalString(request.TargetDate),
	})
	if err != nil {
		return ProjectMilestoneSummary{}, fmt.Errorf("update project milestone %s: %w", request.ID, err)
	}
	if err := mutationSuccess(updated.ProjectMilestoneUpdate.Success, "projectMilestoneUpdate"); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	return projectMilestoneSummary(
		updated.ProjectMilestoneUpdate.ProjectMilestone.ProjectMilestoneSummaryFields,
	), nil
}

// DeleteProjectMilestone hard deletes a ProjectMilestone after resolving and
// comparing its project. This is linctl's one approved irreversible write:
// there is no restore path via linctl.
func DeleteProjectMilestone(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	projectMilestoneID string,
) (string, error) {
	if projectMilestoneID == "" {
		return "", requiredFieldError("project milestone id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}

	return guard.deleteProjectMilestone(ctx, projectMilestoneID)
}

func (guard *guardedClient) deleteProjectMilestone(ctx context.Context, projectMilestoneID string) (string, error) {
	if err := guard.requireProjectMilestone(ctx, projectMilestoneID); err != nil {
		return "", err
	}

	deleted, err := gql.ProjectMilestoneDelete(ctx, guard.graphqlClient, projectMilestoneID)
	if err != nil {
		return "", fmt.Errorf("delete project milestone %s: %w", projectMilestoneID, err)
	}
	if err := mutationSuccess(deleted.ProjectMilestoneDelete.Success, "projectMilestoneDelete"); err != nil {
		return "", err
	}

	return deleted.ProjectMilestoneDelete.EntityId, nil
}

func validateProjectMilestoneUpdateRequest(request ProjectMilestoneUpdateRequest) error {
	if request.ID == "" {
		return requiredFieldError("project milestone id")
	}
	if request.Name == "" && request.Description == "" && request.TargetDate == "" {
		return requiredFieldError("name, description, or target date")
	}

	return nil
}
