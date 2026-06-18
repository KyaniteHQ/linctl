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
	project, err := ProjectByID(ctx, graphqlClient, id)
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("get project %s: %w", id, err)
	}

	return projectSummaryFromFields(project.Project.ProjectSummaryFields), nil
}

// ListProjectMembers returns members for one project.
func ListProjectMembers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMemberList, error) {
	project, err := ProjectMembers(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
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
