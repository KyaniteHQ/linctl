package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// RoadmapSummary is the compact deprecated roadmap model used by read-only commands.
type RoadmapSummary struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description,omitempty"`
	Color              string  `json:"color,omitempty"`
	SlugID             string  `json:"slug_id"`
	SortOrder          float64 `json:"sort_order"`
	ArchivedAt         string  `json:"archived_at,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	URL                string  `json:"url"`
	CreatorID          string  `json:"creator_id"`
	CreatorDisplayName string  `json:"creator_display_name"`
	OwnerID            string  `json:"owner_id,omitempty"`
	OwnerDisplayName   string  `json:"owner_display_name,omitempty"`
}

// RoadmapList is a page of deprecated Linear roadmaps.
type RoadmapList struct {
	Roadmaps    []RoadmapSummary `json:"roadmaps"`
	HasNextPage bool             `json:"has_next_page"`
	EndCursor   *string          `json:"end_cursor,omitempty"`
}

// RoadmapProjectList is a page of Projects associated with one Roadmap.
type RoadmapProjectList struct {
	RoadmapID   string           `json:"roadmap_id"`
	RoadmapName string           `json:"roadmap_name"`
	Projects    []ProjectSummary `json:"projects"`
	HasNextPage bool             `json:"has_next_page"`
	EndCursor   *string          `json:"end_cursor,omitempty"`
}

// ListRoadmaps returns visible Linear roadmaps.
func ListRoadmaps(ctx context.Context, graphqlClient graphql.Client, limit int) (RoadmapList, error) {
	result, err := roadmaps(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return RoadmapList{}, fmt.Errorf("list roadmaps: %w", err)
	}

	summaries := make([]RoadmapSummary, 0, len(result.Roadmaps.Nodes))
	for _, node := range result.Roadmaps.Nodes {
		summaries = append(summaries, roadmapSummary(node.RoadmapSummaryFields))
	}

	return RoadmapList{
		Roadmaps:    summaries,
		HasNextPage: result.Roadmaps.PageInfo.HasNextPage,
		EndCursor:   result.Roadmaps.PageInfo.EndCursor,
	}, nil
}

// GetRoadmapByID returns one deprecated Linear roadmap by id.
func GetRoadmapByID(ctx context.Context, graphqlClient graphql.Client, id string) (RoadmapSummary, error) {
	result, err := roadmap(ctx, graphqlClient, id)
	if err != nil {
		return RoadmapSummary{}, fmt.Errorf("get roadmap %s: %w", id, err)
	}

	return roadmapSummary(result.Roadmap.RoadmapSummaryFields), nil
}

// ListRoadmapProjects returns Projects associated with one Roadmap.
func ListRoadmapProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (RoadmapProjectList, error) {
	result, err := roadmap_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return RoadmapProjectList{}, fmt.Errorf("list roadmap projects %s: %w", id, err)
	}

	projects := make([]ProjectSummary, 0, len(result.Roadmap.Projects.Nodes))
	for _, project := range result.Roadmap.Projects.Nodes {
		projects = append(projects, projectSummaryFromFields(project.ProjectSummaryFields))
	}

	return RoadmapProjectList{
		RoadmapID:   result.Roadmap.Id,
		RoadmapName: result.Roadmap.Name,
		Projects:    projects,
		HasNextPage: result.Roadmap.Projects.PageInfo.HasNextPage,
		EndCursor:   result.Roadmap.Projects.PageInfo.EndCursor,
	}, nil
}

func roadmapSummary(fields RoadmapSummaryFields) RoadmapSummary {
	summary := RoadmapSummary{
		ID:                 fields.Id,
		Name:               fields.Name,
		Description:        stringValue(fields.Description),
		Color:              stringValue(fields.Color),
		SlugID:             fields.SlugId,
		SortOrder:          fields.SortOrder,
		ArchivedAt:         stringValue(fields.ArchivedAt),
		CreatedAt:          fields.CreatedAt,
		UpdatedAt:          fields.UpdatedAt,
		URL:                fields.Url,
		CreatorID:          fields.Creator.Id,
		CreatorDisplayName: fields.Creator.DisplayName,
	}
	if fields.Owner != nil {
		summary.OwnerID = fields.Owner.Id
		summary.OwnerDisplayName = fields.Owner.DisplayName
	}

	return summary
}
