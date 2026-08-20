//nolint:dupl // Minimal association read glue is intentionally uniform across project-association domains.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// InitiativeToProjectSummary is one project association under an Initiative.
type InitiativeToProjectSummary struct {
	ID             string `json:"id"`
	InitiativeID   string `json:"initiative_id"`
	InitiativeName string `json:"initiative_name"`
	ProjectID      string `json:"project_id"`
	ProjectName    string `json:"project_name"`
	ProjectSlugID  string `json:"project_slug_id"`
	ProjectURL     string `json:"project_url"`
	SortOrder      string `json:"sort_order"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ArchivedAt     string `json:"archived_at,omitempty"`
}

// InitiativeToProjectList is a page of Initiative-to-Project associations.
type InitiativeToProjectList struct {
	Associations []InitiativeToProjectSummary `json:"associations"`
	HasNextPage  bool                         `json:"has_next_page"`
	EndCursor    *string                      `json:"end_cursor,omitempty"`
}

//nolint:lll
type initiativeToProjectsNode = gql.XInitiativeToProjectsInitiativeToProjectsInitiativeToProjectConnectionNodesInitiativeToProject

// ListInitiativeToProjects returns visible Initiative-to-Project associations.
func ListInitiativeToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (InitiativeToProjectList, error) {
	page, err := listConnection(
		"list initiative to projects", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]initiativeToProjectsNode, bool, *string, error) {
			result, err := gql.XInitiativeToProjects(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.InitiativeToProjects.Nodes,
				result.InitiativeToProjects.PageInfo.HasNextPage,
				result.InitiativeToProjects.PageInfo.EndCursor,
				nil
		},
		initiativeToProjectsNodeSummary,
	)
	if err != nil {
		return InitiativeToProjectList{}, err
	}

	return InitiativeToProjectList{
		Associations: page.Items,
		HasNextPage:  page.HasNextPage,
		EndCursor:    page.EndCursor,
	}, nil
}

// GetInitiativeToProjectByID returns one Initiative-to-Project association by Linear id.
func GetInitiativeToProjectByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeToProjectSummary, error) {
	result, err := gql.XInitiativeToProject(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeToProjectSummary{}, fmt.Errorf("get initiative to project %s: %w", id, err)
	}

	return initiativeToProjectSummary(result.InitiativeToProject.InitiativeToProjectSummaryFields), nil
}

func initiativeToProjectsNodeSummary(association initiativeToProjectsNode) InitiativeToProjectSummary {
	return initiativeToProjectSummary(association.InitiativeToProjectSummaryFields)
}

func initiativeToProjectSummary(association gql.InitiativeToProjectSummaryFields) InitiativeToProjectSummary {
	return InitiativeToProjectSummary{
		ID:             association.Id,
		InitiativeID:   association.Initiative.Id,
		InitiativeName: association.Initiative.Name,
		ProjectID:      association.Project.Id,
		ProjectName:    association.Project.Name,
		ProjectSlugID:  association.Project.SlugId,
		ProjectURL:     association.Project.Url,
		SortOrder:      association.SortOrder,
		CreatedAt:      association.CreatedAt,
		UpdatedAt:      association.UpdatedAt,
		ArchivedAt:     stringValue(association.ArchivedAt),
	}
}
