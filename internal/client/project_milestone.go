package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

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

// LinearProjectMilestoneCreateInput is the sparse Linear projectMilestoneCreate payload linctl supports.
type LinearProjectMilestoneCreateInput struct {
	ProjectID   string  `json:"projectId"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	TargetDate  *string `json:"targetDate,omitempty"`
}

// LinearProjectMilestoneUpdateInput is the sparse Linear projectMilestoneUpdate payload linctl supports.
type LinearProjectMilestoneUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	TargetDate  *string `json:"targetDate,omitempty"`
}

func projectMilestoneSummary(milestone ProjectMilestoneSummaryFields) ProjectMilestoneSummary {
	description := ""
	if milestone.Description != nil {
		description = *milestone.Description
	}
	targetDate := ""
	if milestone.TargetDate != nil {
		targetDate = *milestone.TargetDate
	}

	return ProjectMilestoneSummary{
		ID:          milestone.Id,
		Name:        milestone.Name,
		Description: description,
		TargetDate:  targetDate,
		Status:      string(milestone.Status),
		Progress:    milestone.Progress,
		SortOrder:   milestone.SortOrder,
	}
}

// ListProjectMilestones returns milestones for one project.
func ListProjectMilestones(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMilestoneList, error) {
	project, err := project_projectMilestones(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectMilestoneList{}, fmt.Errorf("list project milestones %s: %w", id, err)
	}

	milestones := make([]ProjectMilestoneSummary, 0, len(project.Project.ProjectMilestones.Nodes))
	for _, milestone := range project.Project.ProjectMilestones.Nodes {
		milestones = append(milestones, projectMilestoneSummary(milestone.ProjectMilestoneSummaryFields))
	}

	return ProjectMilestoneList{
		ProjectID:   project.Project.Id,
		ProjectName: project.Project.Name,
		Milestones:  milestones,
		HasNextPage: project.Project.ProjectMilestones.PageInfo.HasNextPage,
		EndCursor:   project.Project.ProjectMilestones.PageInfo.EndCursor,
	}, nil
}

// ListAllProjectMilestones returns visible ProjectMilestones across the workspace.
func ListAllProjectMilestones(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (ProjectMilestoneList, error) {
	result, err := projectMilestones(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectMilestoneList{}, fmt.Errorf("list project milestones: %w", err)
	}

	milestones := make([]ProjectMilestoneSummary, 0, len(result.ProjectMilestones.Nodes))
	for _, milestone := range result.ProjectMilestones.Nodes {
		milestones = append(milestones, projectMilestoneSummary(milestone.ProjectMilestoneSummaryFields))
	}

	return ProjectMilestoneList{
		Milestones:  milestones,
		HasNextPage: result.ProjectMilestones.PageInfo.HasNextPage,
		EndCursor:   result.ProjectMilestones.PageInfo.EndCursor,
	}, nil
}

// ListProjectMilestoneIssues returns issues associated with one ProjectMilestone.
func ListProjectMilestoneIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMilestoneIssueList, error) {
	result, err := projectMilestone_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectMilestoneIssueList{}, fmt.Errorf("list project milestone issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.ProjectMilestone.Issues.Nodes))
	for _, node := range result.ProjectMilestone.Issues.Nodes {
		issues = append(issues, issueSummaryFromFields(node.IssueSummaryFields))
	}

	return ProjectMilestoneIssueList{
		ProjectMilestoneID:   result.ProjectMilestone.Id,
		ProjectMilestoneName: result.ProjectMilestone.Name,
		Issues:               issues,
		HasNextPage:          result.ProjectMilestone.Issues.PageInfo.HasNextPage,
		EndCursor:            result.ProjectMilestone.Issues.PageInfo.EndCursor,
	}, nil
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
	milestone, err := projectMilestone(ctx, graphqlClient, id)
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
		return ProjectMilestoneSummary{}, fmt.Errorf("%w: project id is required", ErrWriteInvalid)
	}
	if request.Name == "" {
		return ProjectMilestoneSummary{}, fmt.Errorf("%w: name is required", ErrWriteInvalid)
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectMilestoneSummary{}, err
	}
	if err := guard.requireProject(ctx, graphqlClient, request.ProjectID); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	created, err := ProjectMilestoneCreate(ctx, graphqlClient, LinearProjectMilestoneCreateInput{
		ProjectID:   request.ProjectID,
		Name:        request.Name,
		Description: optionalString(request.Description),
		TargetDate:  optionalString(request.TargetDate),
	})
	if err != nil {
		return ProjectMilestoneSummary{}, fmt.Errorf("create project milestone: %w", err)
	}
	if !created.ProjectMilestoneCreate.Success {
		return ProjectMilestoneSummary{}, fmt.Errorf("%w: projectMilestoneCreate failed", ErrMutationFailed)
	}

	return projectMilestoneSummary(
		created.ProjectMilestoneCreate.ProjectMilestone.ProjectMilestoneSummaryFields,
	), nil
}

// UpdateProjectMilestone updates a ProjectMilestone after resolving and comparing its project.
func UpdateProjectMilestone(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectMilestoneUpdateRequest,
) (ProjectMilestoneSummary, error) {
	if err := validateProjectMilestoneUpdateRequest(request); err != nil {
		return ProjectMilestoneSummary{}, err
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectMilestoneSummary{}, err
	}
	if _, err := guard.requireProjectMilestone(ctx, graphqlClient, request.ID); err != nil {
		return ProjectMilestoneSummary{}, err
	}

	updated, err := ProjectMilestoneUpdate(ctx, graphqlClient, request.ID, LinearProjectMilestoneUpdateInput{
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		TargetDate:  optionalString(request.TargetDate),
	})
	if err != nil {
		return ProjectMilestoneSummary{}, fmt.Errorf("update project milestone %s: %w", request.ID, err)
	}
	if !updated.ProjectMilestoneUpdate.Success {
		return ProjectMilestoneSummary{}, fmt.Errorf("%w: projectMilestoneUpdate failed", ErrMutationFailed)
	}

	return projectMilestoneSummary(
		updated.ProjectMilestoneUpdate.ProjectMilestone.ProjectMilestoneSummaryFields,
	), nil
}

func validateProjectMilestoneUpdateRequest(request ProjectMilestoneUpdateRequest) error {
	if request.ID == "" {
		return fmt.Errorf("%w: project milestone id is required", ErrWriteInvalid)
	}
	if request.Name == "" && request.Description == "" && request.TargetDate == "" {
		return fmt.Errorf("%w: name, description, or target date is required", ErrWriteInvalid)
	}

	return nil
}
