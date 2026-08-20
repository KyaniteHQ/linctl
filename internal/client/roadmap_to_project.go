//nolint:dupl // Minimal association read glue is intentionally uniform across project-association domains.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// RoadmapToProjectSummary is one project association under a Roadmap.
type RoadmapToProjectSummary struct {
	ID            string `json:"id"`
	RoadmapID     string `json:"roadmap_id"`
	RoadmapName   string `json:"roadmap_name"`
	ProjectID     string `json:"project_id"`
	ProjectName   string `json:"project_name"`
	ProjectSlugID string `json:"project_slug_id"`
	ProjectURL    string `json:"project_url"`
	SortOrder     string `json:"sort_order"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ArchivedAt    string `json:"archived_at,omitempty"`
}

// RoadmapToProjectList is a page of Roadmap-to-Project associations.
type RoadmapToProjectList struct {
	Associations []RoadmapToProjectSummary `json:"associations"`
	Page
}

//nolint:lll
type roadmapToProjectsNode = gql.XRoadmapToProjectsRoadmapToProjectsRoadmapToProjectConnectionNodesRoadmapToProject

// ListRoadmapToProjects returns visible Roadmap-to-Project associations.
func ListRoadmapToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (RoadmapToProjectList, error) {
	page, err := listConnection(
		"list roadmap to projects", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]roadmapToProjectsNode, bool, *string, error) {
			result, err := gql.XRoadmapToProjects(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.RoadmapToProjects.Nodes,
				result.RoadmapToProjects.PageInfo.HasNextPage,
				result.RoadmapToProjects.PageInfo.EndCursor,
				nil
		},
		roadmapToProjectsNodeSummary,
	)
	if err != nil {
		return RoadmapToProjectList{}, err
	}

	return RoadmapToProjectList{
		Associations: page.Items,
		Page:         page.Page,
	}, nil
}

// GetRoadmapToProjectByID returns one Roadmap-to-Project association by Linear id.
func GetRoadmapToProjectByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (RoadmapToProjectSummary, error) {
	result, err := gql.XRoadmapToProject(ctx, graphqlClient, id)
	if err != nil {
		return RoadmapToProjectSummary{}, fmt.Errorf("get roadmap to project %s: %w", id, err)
	}

	return roadmapToProjectSummary(result.RoadmapToProject.RoadmapToProjectSummaryFields), nil
}

func roadmapToProjectsNodeSummary(association roadmapToProjectsNode) RoadmapToProjectSummary {
	return roadmapToProjectSummary(association.RoadmapToProjectSummaryFields)
}

func roadmapToProjectSummary(association gql.RoadmapToProjectSummaryFields) RoadmapToProjectSummary {
	return RoadmapToProjectSummary{
		ID:            association.Id,
		RoadmapID:     association.Roadmap.Id,
		RoadmapName:   association.Roadmap.Name,
		ProjectID:     association.Project.Id,
		ProjectName:   association.Project.Name,
		ProjectSlugID: association.Project.SlugId,
		ProjectURL:    association.Project.Url,
		SortOrder:     association.SortOrder,
		CreatedAt:     association.CreatedAt,
		UpdatedAt:     association.UpdatedAt,
		ArchivedAt:    stringValue(association.ArchivedAt),
	}
}
