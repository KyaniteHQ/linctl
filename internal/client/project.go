package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ProjectSummary is the compact project model used by project commands.
type ProjectSummary struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	SlugID      string        `json:"slug_id"`
	URL         string        `json:"url"`
	Priority    int           `json:"priority"`
	Status      ProjectStatus `json:"status"`
	Lead        string        `json:"lead,omitempty"`
	Teams       []ProjectTeam `json:"teams"`
}

// ProjectStatus is the compact project lifecycle status.
type ProjectStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ProjectTeam is a project-associated team.
type ProjectTeam struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ProjectList is a page of projects scoped to a team.
type ProjectList struct {
	Projects    []ProjectSummary `json:"projects"`
	HasNextPage bool             `json:"has_next_page"`
	EndCursor   *string          `json:"end_cursor,omitempty"`
}

// ProjectMember is a project member.
type ProjectMember struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// ProjectMemberList is a page of project members.
type ProjectMemberList struct {
	ProjectID   string          `json:"project_id"`
	ProjectName string          `json:"project_name"`
	Members     []ProjectMember `json:"members"`
	HasNextPage bool            `json:"has_next_page"`
	EndCursor   *string         `json:"end_cursor,omitempty"`
}

// ProjectUpdateSummary is one project status update.
type ProjectUpdateSummary struct {
	ID          string `json:"id"`
	Body        string `json:"body"`
	Health      string `json:"health"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	URL         string `json:"url"`
	ProjectID   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// ProjectUpdateList is a page of project status updates.
type ProjectUpdateList struct {
	ProjectID   string                 `json:"project_id"`
	ProjectName string                 `json:"project_name"`
	Updates     []ProjectUpdateSummary `json:"updates"`
	HasNextPage bool                   `json:"has_next_page"`
	EndCursor   *string                `json:"end_cursor,omitempty"`
}

// ListProjectsByTeam returns projects scoped to a resolved team.
func ListProjectsByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (ProjectList, error) {
	projects, err := Projects(ctx, graphqlClient, teamID, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list projects: %w", err)
	}

	summaries := make([]ProjectSummary, 0, len(projects.Team.Projects.Nodes))
	for _, project := range projects.Team.Projects.Nodes {
		summaries = append(summaries, projectSummaryFromFields(project.ProjectSummaryFields))
	}

	return ProjectList{
		Projects:    summaries,
		HasNextPage: projects.Team.Projects.PageInfo.HasNextPage,
		EndCursor:   projects.Team.Projects.PageInfo.EndCursor,
	}, nil
}

// GetProjectByID returns a project by Linear id or slug.
func GetProjectByID(ctx context.Context, graphqlClient graphql.Client, id string) (ProjectSummary, error) {
	projectResult, err := project(ctx, graphqlClient, id)
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("get project %s: %w", id, err)
	}

	return projectSummaryFromFields(projectResult.Project.ProjectSummaryFields), nil
}

// ListProjectMembers returns members for one project.
func ListProjectMembers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMemberList, error) {
	project, err := project_members(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectMemberList{}, fmt.Errorf("list project members %s: %w", id, err)
	}

	members := make([]ProjectMember, 0, len(project.Project.Members.Nodes))
	for _, member := range project.Project.Members.Nodes {
		members = append(members, ProjectMember{
			ID:          member.Id,
			Name:        member.Name,
			DisplayName: member.DisplayName,
			Email:       member.Email,
		})
	}

	return ProjectMemberList{
		ProjectID:   project.Project.Id,
		ProjectName: project.Project.Name,
		Members:     members,
		HasNextPage: project.Project.Members.PageInfo.HasNextPage,
		EndCursor:   project.Project.Members.PageInfo.EndCursor,
	}, nil
}

// ListProjectUpdates returns status updates for one project.
func ListProjectUpdates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectUpdateList, error) {
	project, err := ProjectUpdates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateList{}, fmt.Errorf("list project updates %s: %w", id, err)
	}

	updates := make([]ProjectUpdateSummary, 0, len(project.Project.ProjectUpdates.Nodes))
	for _, update := range project.Project.ProjectUpdates.Nodes {
		updates = append(updates, ProjectUpdateSummary{
			ID:          update.Id,
			Body:        update.Body,
			Health:      string(update.Health),
			CreatedAt:   update.CreatedAt,
			UpdatedAt:   update.UpdatedAt,
			URL:         update.Url,
			UserID:      update.User.Id,
			Name:        update.User.Name,
			DisplayName: update.User.DisplayName,
		})
	}

	return ProjectUpdateList{
		ProjectID:   project.Project.Id,
		ProjectName: project.Project.Name,
		Updates:     updates,
		HasNextPage: project.Project.ProjectUpdates.PageInfo.HasNextPage,
		EndCursor:   project.Project.ProjectUpdates.PageInfo.EndCursor,
	}, nil
}

// ListAllProjectUpdates returns visible project status updates across projects.
func ListAllProjectUpdates(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectUpdateList, error) {
	updatesResponse, err := projectUpdates(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateList{}, fmt.Errorf("list project updates: %w", err)
	}

	updates := make([]ProjectUpdateSummary, 0, len(updatesResponse.ProjectUpdates.Nodes))
	for _, update := range updatesResponse.ProjectUpdates.Nodes {
		updates = append(updates, projectUpdateSummary(update.TopLevelProjectUpdateSummaryFields))
	}

	return ProjectUpdateList{
		Updates:     updates,
		HasNextPage: updatesResponse.ProjectUpdates.PageInfo.HasNextPage,
		EndCursor:   updatesResponse.ProjectUpdates.PageInfo.EndCursor,
	}, nil
}

// GetProjectUpdateByID returns one project update by Linear id.
func GetProjectUpdateByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectUpdateSummary, error) {
	update, err := projectUpdate(ctx, graphqlClient, id)
	if err != nil {
		return ProjectUpdateSummary{}, fmt.Errorf("get project update %s: %w", id, err)
	}

	return projectUpdateSummary(update.ProjectUpdate.TopLevelProjectUpdateSummaryFields), nil
}

func projectUpdateSummary(update TopLevelProjectUpdateSummaryFields) ProjectUpdateSummary {
	return ProjectUpdateSummary{
		ID:          update.Id,
		Body:        update.Body,
		Health:      string(update.Health),
		CreatedAt:   update.CreatedAt,
		UpdatedAt:   update.UpdatedAt,
		URL:         update.Url,
		ProjectID:   update.Project.Id,
		ProjectName: update.Project.Name,
		UserID:      update.User.Id,
		Name:        update.User.Name,
		DisplayName: update.User.DisplayName,
	}
}

func projectSummaryFromFields(project ProjectSummaryFields) ProjectSummary {
	lead := ""
	if project.Lead != nil {
		lead = project.Lead.DisplayName
	}

	teams := make([]ProjectTeam, 0, len(project.Teams.Nodes))
	for _, team := range project.Teams.Nodes {
		teams = append(teams, ProjectTeam{
			ID:   team.Id,
			Key:  team.Key,
			Name: team.Name,
		})
	}

	return ProjectSummary{
		ID:          project.Id,
		Name:        project.Name,
		Description: project.Description,
		SlugID:      project.SlugId,
		URL:         project.Url,
		Priority:    project.Priority,
		Status: ProjectStatus{
			ID:   project.Status.Id,
			Name: project.Status.Name,
			Type: string(project.Status.Type),
		},
		Lead:  lead,
		Teams: teams,
	}
}
