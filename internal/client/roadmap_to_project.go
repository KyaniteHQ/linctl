//nolint:dupl // Minimal association read glue is intentionally uniform across project-association domains.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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
	HasNextPage  bool                      `json:"has_next_page"`
	EndCursor    *string                   `json:"end_cursor,omitempty"`
}

// ListRoadmapToProjects returns visible Roadmap-to-Project associations.
func ListRoadmapToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (RoadmapToProjectList, error) {
	result, err := roadmapToProjects(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return RoadmapToProjectList{}, fmt.Errorf("list roadmap to projects: %w", err)
	}

	associations := make([]RoadmapToProjectSummary, 0, len(result.RoadmapToProjects.Nodes))
	for _, association := range result.RoadmapToProjects.Nodes {
		associations = append(
			associations,
			roadmapToProjectSummary(association.RoadmapToProjectSummaryFields),
		)
	}

	return RoadmapToProjectList{
		Associations: associations,
		HasNextPage:  result.RoadmapToProjects.PageInfo.HasNextPage,
		EndCursor:    result.RoadmapToProjects.PageInfo.EndCursor,
	}, nil
}

// GetRoadmapToProjectByID returns one Roadmap-to-Project association by Linear id.
func GetRoadmapToProjectByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (RoadmapToProjectSummary, error) {
	result, err := roadmapToProject(ctx, graphqlClient, id)
	if err != nil {
		return RoadmapToProjectSummary{}, fmt.Errorf("get roadmap to project %s: %w", id, err)
	}

	return roadmapToProjectSummary(result.RoadmapToProject.RoadmapToProjectSummaryFields), nil
}

func roadmapToProjectSummary(association RoadmapToProjectSummaryFields) RoadmapToProjectSummary {
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
