package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

const projectListPageSize = 50

// ProjectSummary is the compact project model used by project commands.
type ProjectSummary struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ArchivedAt  string        `json:"archived_at,omitempty"`
	SlugID      string        `json:"slug_id"`
	URL         string        `json:"url"`
	Priority    int           `json:"priority"`
	Status      ProjectStatus `json:"status"`
	Lead        string        `json:"lead,omitempty"`
	Teams       []ProjectTeam `json:"teams"`
	// TeamsTruncated is true when the project has more teams than the fetched
	// page, so an unmatched pinned team cannot be ruled out from this page alone.
	TeamsTruncated bool `json:"teams_truncated,omitempty"`
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
	Body        string `json:"body,omitempty"`
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

// ProjectFilterSuggestion is an AI-generated project filter suggestion.
type ProjectFilterSuggestion struct {
	Filter json.RawMessage `json:"filter,omitempty"`
	LogID  string          `json:"log_id,omitempty"`
}

// ProjectUpdateCommentList is a page of body-free Comments associated with one ProjectUpdate.
type ProjectUpdateCommentList struct {
	ProjectUpdateID string                   `json:"project_update_id"`
	Comments        []CommentMetadataSummary `json:"comments"`
	HasNextPage     bool                     `json:"has_next_page"`
	EndCursor       *string                  `json:"end_cursor,omitempty"`
}

// ProjectAttachmentList is a page of Attachments associated with one Project.
type ProjectAttachmentList struct {
	ProjectID   string              `json:"project_id"`
	ProjectName string              `json:"project_name"`
	Attachments []AttachmentSummary `json:"attachments"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ProjectDocumentList is a page of Documents associated with one Project.
type ProjectDocumentList struct {
	ProjectID   string            `json:"project_id"`
	ProjectName string            `json:"project_name"`
	Documents   []DocumentSummary `json:"documents"`
	HasNextPage bool              `json:"has_next_page"`
	EndCursor   *string           `json:"end_cursor,omitempty"`
}

// ProjectExternalLinkList is a page of external links associated with one Project.
type ProjectExternalLinkList struct {
	ProjectID   string                      `json:"project_id"`
	ProjectName string                      `json:"project_name"`
	Links       []EntityExternalLinkSummary `json:"links"`
	HasNextPage bool                        `json:"has_next_page"`
	EndCursor   *string                     `json:"end_cursor,omitempty"`
}

// ProjectHistorySummary is the compact project history model used by read-only commands.
type ProjectHistorySummary struct {
	ID         string          `json:"id"`
	ProjectID  string          `json:"project_id"`
	EntryCount int             `json:"entry_count"`
	Entries    json.RawMessage `json:"entries"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
	ArchivedAt string          `json:"archived_at,omitempty"`
}

// ProjectHistoryList is a page of Linear project history records.
type ProjectHistoryList struct {
	ProjectID   string                  `json:"project_id"`
	ProjectName string                  `json:"project_name"`
	History     []ProjectHistorySummary `json:"history"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// ProjectInitiativeToProjectList is a page of initiative associations for one Project.
type ProjectInitiativeToProjectList struct {
	ProjectID    string                       `json:"project_id"`
	ProjectName  string                       `json:"project_name"`
	Associations []InitiativeToProjectSummary `json:"associations"`
	HasNextPage  bool                         `json:"has_next_page"`
	EndCursor    *string                      `json:"end_cursor,omitempty"`
}

// ProjectInitiativeList is a page of Initiatives associated with one Project.
type ProjectInitiativeList struct {
	ProjectID   string              `json:"project_id"`
	ProjectName string              `json:"project_name"`
	Initiatives []InitiativeSummary `json:"initiatives"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ProjectIssueList is a page of Issues associated with one Project.
type ProjectIssueList struct {
	ProjectID   string         `json:"project_id"`
	ProjectName string         `json:"project_name"`
	Issues      []IssueSummary `json:"issues"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// ProjectCommentList is a page of body-free Comments associated with one Project.
type ProjectCommentList struct {
	ProjectID   string                   `json:"project_id"`
	ProjectName string                   `json:"project_name"`
	Comments    []CommentMetadataSummary `json:"comments"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ProjectProjectLabelList is a page of ProjectLabels associated with one Project.
type ProjectProjectLabelList struct {
	ProjectID     string                `json:"project_id"`
	ProjectName   string                `json:"project_name"`
	ProjectLabels []ProjectLabelSummary `json:"project_labels"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

// ProjectCustomerNeedList is a page of customer needs associated with one Project.
type ProjectCustomerNeedList struct {
	ProjectID   string                `json:"project_id"`
	ProjectName string                `json:"project_name"`
	Needs       []CustomerNeedSummary `json:"customer_needs"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// ProjectProjectRelationList is a page of project relations associated with one Project.
type ProjectProjectRelationList struct {
	ProjectID   string                   `json:"project_id"`
	ProjectName string                   `json:"project_name"`
	Relations   []ProjectRelationSummary `json:"relations"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ProjectTeamList is a page of Teams associated with one Project.
type ProjectTeamList struct {
	ProjectID   string        `json:"project_id"`
	ProjectName string        `json:"project_name"`
	Teams       []TeamSummary `json:"teams"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// ListProjectsByTeam returns projects scoped to a resolved team.
func ListProjectsByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (ProjectList, error) {
	page, err := collectNodePages(
		"list projects", limit, projectListPageSize,
		func(pageSize int, after *string) (nodePage[ProjectSummary], error) {
			projects, err := gql.Projects(ctx, graphqlClient, teamID, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nodePage[ProjectSummary]{}, err
			}

			result := projects.Team.Projects
			return nodePage[ProjectSummary]{
				Items: mapNodes(result.Nodes, func(
					project gql.ProjectsTeamProjectsProjectConnectionNodesProject,
				) ProjectSummary {
					return projectSummaryFromFields(project.ProjectSummaryFields)
				}),
				HasNextPage: result.PageInfo.HasNextPage,
				EndCursor:   result.PageInfo.EndCursor,
			}, nil
		},
	)
	if err != nil {
		return ProjectList{}, err
	}

	return ProjectList{Projects: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListProjects returns Linear projects visible to the authenticated user.
func ListProjects(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectList, error) {
	page, err := collectNodePages(
		"list projects", limit, projectListPageSize,
		func(pageSize int, after *string) (nodePage[ProjectSummary], error) {
			result, err := gql.XProjects(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nodePage[ProjectSummary]{}, err
			}

			projects := result.Projects
			return nodePage[ProjectSummary]{
				Items: mapNodes(projects.Nodes, func(
					project gql.XProjectsProjectsProjectConnectionNodesProject,
				) ProjectSummary {
					return projectSummaryFromFields(project.ProjectSummaryFields)
				}),
				HasNextPage: projects.PageInfo.HasNextPage,
				EndCursor:   projects.PageInfo.EndCursor,
			}, nil
		},
	)
	if err != nil {
		return ProjectList{}, err
	}

	return ProjectList{Projects: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetProjectByID returns a project by Linear id or slug.
func GetProjectByID(ctx context.Context, graphqlClient graphql.Client, id string) (ProjectSummary, error) {
	projectResult, err := gql.XProject(ctx, graphqlClient, id)
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("get project %s: %w", id, err)
	}

	return projectSummaryFromFields(projectResult.Project.ProjectSummaryFields), nil
}

// GetProjectFilterSuggestion returns a JSON project filter suggestion for a prompt.
func GetProjectFilterSuggestion(
	ctx context.Context,
	graphqlClient graphql.Client,
	prompt string,
	teamID string,
) (ProjectFilterSuggestion, error) {
	suggestion, err := gql.XProjectFilterSuggestion(ctx, graphqlClient, prompt, optionalString(teamID))
	if err != nil {
		return ProjectFilterSuggestion{}, fmt.Errorf("get project filter suggestion: %w", err)
	}

	filter := json.RawMessage(nil)
	if suggestion.ProjectFilterSuggestion.Filter != nil {
		filter = *suggestion.ProjectFilterSuggestion.Filter
	}

	return ProjectFilterSuggestion{
		Filter: filter,
		LogID:  stringValue(suggestion.ProjectFilterSuggestion.LogId),
	}, nil
}

func projectSummaryFromFields(project gql.ProjectSummaryFields) ProjectSummary {
	lead := ""
	if project.Lead != nil {
		lead = project.Lead.DisplayName
	}

	teams := mapNodes(project.Teams.Nodes, func(team gql.ProjectSummaryFieldsTeamsTeamConnectionNodesTeam) ProjectTeam {
		return ProjectTeam{
			ID:   team.Id,
			Key:  team.Key,
			Name: team.Name,
		}
	})

	return ProjectSummary{
		ID:          project.Id,
		Name:        project.Name,
		Description: project.Description,
		ArchivedAt:  stringValue(project.ArchivedAt),
		SlugID:      project.SlugId,
		URL:         project.Url,
		Priority:    project.Priority,
		Status: ProjectStatus{
			ID:   project.Status.Id,
			Name: project.Status.Name,
			Type: string(project.Status.Type),
		},
		Lead:           lead,
		Teams:          teams,
		TeamsTruncated: project.Teams.PageInfo.HasNextPage,
	}
}
