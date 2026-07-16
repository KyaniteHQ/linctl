package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ListProjectUpdates returns status updates for one project.
func ListProjectUpdates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectUpdateList, error) {
	project, err := gql.XProject_projectUpdates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateList{}, fmt.Errorf("list project updates %s: %w", id, err)
	}

	updates := mapNodes(project.Project.ProjectUpdates.Nodes, projectScopedProjectUpdateSummary)

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
	updatesResponse, err := gql.XProjectUpdates(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateList{}, fmt.Errorf("list project updates: %w", err)
	}

	updates := mapNodes(updatesResponse.ProjectUpdates.Nodes, func(
		update gql.XProjectUpdatesProjectUpdatesProjectUpdateConnectionNodesProjectUpdate,
	) ProjectUpdateSummary {
		return projectUpdateSummary(update.TopLevelProjectUpdateSummaryFields)
	})

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
	update, err := gql.XProjectUpdate(ctx, graphqlClient, id)
	if err != nil {
		return ProjectUpdateSummary{}, fmt.Errorf("get project update %s: %w", id, err)
	}

	return projectUpdateSummary(update.ProjectUpdate.TopLevelProjectUpdateSummaryFields), nil
}

// ListProjectUpdateComments returns body-free comments associated with one ProjectUpdate.
func ListProjectUpdateComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectUpdateCommentList, error) {
	result, err := gql.XProjectUpdate_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateCommentList{}, fmt.Errorf("list project update comments %s: %w", id, err)
	}

	comments := mapNodes(result.ProjectUpdate.Comments.Nodes, func(
		node gql.XProjectUpdate_commentsProjectUpdateCommentsCommentConnectionNodesComment,
	) CommentMetadataSummary {
		return commentMetadataSummary(node.CommentMetadataFields)
	})

	return ProjectUpdateCommentList{
		ProjectUpdateID: result.ProjectUpdate.Id,
		Comments:        comments,
		HasNextPage:     result.ProjectUpdate.Comments.PageInfo.HasNextPage,
		EndCursor:       result.ProjectUpdate.Comments.PageInfo.EndCursor,
	}, nil
}

func projectScopedProjectUpdateSummary(
	update gql.XProject_projectUpdatesProjectProjectUpdatesProjectUpdateConnectionNodesProjectUpdate,
) ProjectUpdateSummary {
	return ProjectUpdateSummary{
		ID:          update.Id,
		Health:      string(update.Health),
		CreatedAt:   update.CreatedAt,
		UpdatedAt:   update.UpdatedAt,
		URL:         update.Url,
		UserID:      update.User.Id,
		Name:        update.User.Name,
		DisplayName: update.User.DisplayName,
	}
}

func projectUpdateSummary(update gql.TopLevelProjectUpdateSummaryFields) ProjectUpdateSummary {
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
