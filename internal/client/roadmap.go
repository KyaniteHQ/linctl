package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Roadmaps []RoadmapSummary `json:"roadmaps"`
	Page
}

// RoadmapProjectList is a page of Projects associated with one Roadmap.
type RoadmapProjectList struct {
	RoadmapID   string           `json:"roadmap_id"`
	RoadmapName string           `json:"roadmap_name"`
	Projects    []ProjectSummary `json:"projects"`
	Page
}

//nolint:lll
type roadmapsNode = gql.XRoadmapsRoadmapsRoadmapConnectionNodesRoadmap

//nolint:lll
type roadmapProjectsNode = gql.XRoadmap_projectsRoadmapProjectsProjectConnectionNodesProject

type roadmapsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type roadmapProjectsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// roadmapProjectsParent is the connection parent metadata roadmapProjectsQuery reads out of
// every page. Linear repeats it per page, so the last page wins.
type roadmapProjectsParent struct {
	roadmapID   string
	roadmapName string
}

// ListRoadmaps returns visible Linear roadmaps.
func ListRoadmaps(ctx context.Context, graphqlClient graphql.Client, limit int) (RoadmapList, error) {
	query := roadmapsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list roadmaps", limit, defaultListPageSize,
		query.page,
		roadmapsNodeSummary,
	)
	if err != nil {
		return RoadmapList{}, err
	}

	return RoadmapList{Roadmaps: page.Items, Page: page.Page}, nil
}

// GetRoadmapByID returns one deprecated Linear roadmap by id.
func GetRoadmapByID(ctx context.Context, graphqlClient graphql.Client, id string) (RoadmapSummary, error) {
	result, err := gql.XRoadmap(ctx, graphqlClient, id)
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
	query := &roadmapProjectsQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list roadmap projects "+id, limit, defaultListPageSize,
		query.projects,
		roadmapProjectsNodeSummary,
	)
	if err != nil {
		return RoadmapProjectList{}, err
	}

	return RoadmapProjectList{
		RoadmapID:   parent.roadmapID,
		RoadmapName: parent.roadmapName,
		Projects:    page.Items,
		Page:        page.Page,
	}, nil
}

func (query roadmapsQuery) page(pageSize int, after *string) ([]roadmapsNode, bool, *string, error) {
	result, err := gql.XRoadmaps(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Roadmaps.Nodes,
		result.Roadmaps.PageInfo.HasNextPage,
		result.Roadmaps.PageInfo.EndCursor,
		nil
}

func (query *roadmapProjectsQuery) projects(
	pageSize int,
	after *string,
) ([]roadmapProjectsNode, roadmapProjectsParent, bool, *string, error) {
	result, err := gql.XRoadmap_projects(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, roadmapProjectsParent{}, false, nil, err
	}

	return result.Roadmap.Projects.Nodes,
		roadmapProjectsParent{roadmapID: result.Roadmap.Id, roadmapName: result.Roadmap.Name},
		result.Roadmap.Projects.PageInfo.HasNextPage,
		result.Roadmap.Projects.PageInfo.EndCursor,
		nil
}

func roadmapsNodeSummary(node roadmapsNode) RoadmapSummary {
	return roadmapSummary(node.RoadmapSummaryFields)
}

func roadmapProjectsNodeSummary(project roadmapProjectsNode) ProjectSummary {
	return projectSummaryFromFields(project.ProjectSummaryFields)
}

func roadmapSummary(fields gql.RoadmapSummaryFields) RoadmapSummary {
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
