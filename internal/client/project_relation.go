package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ProjectRelationSummary is one dependency relation between Linear projects.
type ProjectRelationSummary struct {
	ID                          string `json:"id"`
	Type                        string `json:"type"`
	AnchorType                  string `json:"anchor_type"`
	RelatedAnchorType           string `json:"related_anchor_type"`
	ProjectID                   string `json:"project_id"`
	ProjectName                 string `json:"project_name"`
	ProjectMilestoneID          string `json:"project_milestone_id,omitempty"`
	ProjectMilestoneName        string `json:"project_milestone_name,omitempty"`
	RelatedProjectID            string `json:"related_project_id"`
	RelatedProjectName          string `json:"related_project_name"`
	RelatedProjectMilestoneID   string `json:"related_project_milestone_id,omitempty"`
	RelatedProjectMilestoneName string `json:"related_project_milestone_name,omitempty"`
	CreatedAt                   string `json:"created_at"`
	UpdatedAt                   string `json:"updated_at"`
	ArchivedAt                  string `json:"archived_at,omitempty"`
	UserID                      string `json:"user_id,omitempty"`
	Name                        string `json:"name,omitempty"`
	DisplayName                 string `json:"display_name,omitempty"`
}

// ProjectRelationList is a page of project dependency relations.
type ProjectRelationList struct {
	Relations   []ProjectRelationSummary `json:"relations"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ListProjectRelations returns visible dependency relations between projects.
func ListProjectRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (ProjectRelationList, error) {
	result, err := projectRelations(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectRelationList{}, fmt.Errorf("list project relations: %w", err)
	}

	relations := make([]ProjectRelationSummary, 0, len(result.ProjectRelations.Nodes))
	for _, relation := range result.ProjectRelations.Nodes {
		relations = append(relations, projectRelationSummary(relation.ProjectRelationSummaryFields))
	}

	return ProjectRelationList{
		Relations:   relations,
		HasNextPage: result.ProjectRelations.PageInfo.HasNextPage,
		EndCursor:   result.ProjectRelations.PageInfo.EndCursor,
	}, nil
}

// GetProjectRelationByID returns one project relation by Linear id.
func GetProjectRelationByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectRelationSummary, error) {
	result, err := projectRelation(ctx, graphqlClient, id)
	if err != nil {
		return ProjectRelationSummary{}, fmt.Errorf("get project relation %s: %w", id, err)
	}

	return projectRelationSummary(result.ProjectRelation.ProjectRelationSummaryFields), nil
}

func projectRelationSummary(relation ProjectRelationSummaryFields) ProjectRelationSummary {
	summary := ProjectRelationSummary{
		ID:                 relation.Id,
		Type:               relation.Type,
		AnchorType:         relation.AnchorType,
		RelatedAnchorType:  relation.RelatedAnchorType,
		ProjectID:          relation.Project.Id,
		ProjectName:        relation.Project.Name,
		RelatedProjectID:   relation.RelatedProject.Id,
		RelatedProjectName: relation.RelatedProject.Name,
		CreatedAt:          relation.CreatedAt,
		UpdatedAt:          relation.UpdatedAt,
		ArchivedAt:         stringValue(relation.ArchivedAt),
	}
	if relation.ProjectMilestone != nil {
		summary.ProjectMilestoneID = relation.ProjectMilestone.Id
		summary.ProjectMilestoneName = relation.ProjectMilestone.Name
	}
	if relation.RelatedProjectMilestone != nil {
		summary.RelatedProjectMilestoneID = relation.RelatedProjectMilestone.Id
		summary.RelatedProjectMilestoneName = relation.RelatedProjectMilestone.Name
	}
	if relation.User != nil {
		summary.UserID = relation.User.Id
		summary.Name = relation.User.Name
		summary.DisplayName = relation.User.DisplayName
	}

	return summary
}
