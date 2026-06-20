package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListInitiativeToProjects returns visible Initiative-to-Project associations.
func ListInitiativeToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (InitiativeToProjectList, error) {
	result, err := initiativeToProjects(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeToProjectList{}, fmt.Errorf("list initiative to projects: %w", err)
	}

	associations := make([]InitiativeToProjectSummary, 0, len(result.InitiativeToProjects.Nodes))
	for _, association := range result.InitiativeToProjects.Nodes {
		associations = append(
			associations,
			initiativeToProjectSummary(association.InitiativeToProjectSummaryFields),
		)
	}

	return InitiativeToProjectList{
		Associations: associations,
		HasNextPage:  result.InitiativeToProjects.PageInfo.HasNextPage,
		EndCursor:    result.InitiativeToProjects.PageInfo.EndCursor,
	}, nil
}

// GetInitiativeToProjectByID returns one Initiative-to-Project association by Linear id.
func GetInitiativeToProjectByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeToProjectSummary, error) {
	result, err := initiativeToProject(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeToProjectSummary{}, fmt.Errorf("get initiative to project %s: %w", id, err)
	}

	return initiativeToProjectSummary(result.InitiativeToProject.InitiativeToProjectSummaryFields), nil
}

func initiativeToProjectSummary(association InitiativeToProjectSummaryFields) InitiativeToProjectSummary {
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
