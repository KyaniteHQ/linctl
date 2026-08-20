package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

//nolint:lll
type projectUpdatesNode = gql.XProjectUpdatesProjectUpdatesProjectUpdateConnectionNodesProjectUpdate

//nolint:lll
type projectScopedUpdatesNode = gql.XProject_projectUpdatesProjectProjectUpdatesProjectUpdateConnectionNodesProjectUpdate

//nolint:lll
type projectUpdateCommentsNode = gql.XProjectUpdate_commentsProjectUpdateCommentsCommentConnectionNodesComment

type projectUpdatesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type projectScopedUpdatesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

type projectUpdateCommentsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// ListProjectUpdates returns status updates for one project.
func ListProjectUpdates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectUpdateList, error) {
	query := &projectScopedUpdatesQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project updates "+id, limit, defaultListPageSize,
		query.page,
		projectScopedProjectUpdateSummary,
	)
	if err != nil {
		return ProjectUpdateList{}, err
	}

	return ProjectUpdateList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Updates:     page.Items,
		Page:        page.Page,
	}, nil
}

// ListAllProjectUpdates returns visible project status updates across projects.
func ListAllProjectUpdates(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectUpdateList, error) {
	query := projectUpdatesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list project updates", limit, defaultListPageSize,
		query.page,
		projectUpdateNodeSummary,
	)
	if err != nil {
		return ProjectUpdateList{}, err
	}

	return ProjectUpdateList{Updates: page.Items, Page: page.Page}, nil
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
	query := &projectUpdateCommentsQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project update comments "+id, limit, defaultListPageSize,
		query.page,
		projectUpdateCommentNodeSummary,
	)
	if err != nil {
		return ProjectUpdateCommentList{}, err
	}

	return ProjectUpdateCommentList{
		ProjectUpdateID: parent,
		Comments:        page.Items,
		Page:            page.Page,
	}, nil
}

func (query *projectScopedUpdatesQuery) page(
	pageSize int,
	after *string,
) ([]projectScopedUpdatesNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_projectUpdates(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.ProjectUpdates.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.ProjectUpdates.PageInfo.HasNextPage,
		result.Project.ProjectUpdates.PageInfo.EndCursor,
		nil
}

func (query projectUpdatesQuery) page(
	pageSize int,
	after *string,
) ([]projectUpdatesNode, bool, *string, error) {
	result, err := gql.XProjectUpdates(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ProjectUpdates.Nodes,
		result.ProjectUpdates.PageInfo.HasNextPage,
		result.ProjectUpdates.PageInfo.EndCursor,
		nil
}

func (query *projectUpdateCommentsQuery) page(
	pageSize int,
	after *string,
) ([]projectUpdateCommentsNode, string, bool, *string, error) {
	result, err := gql.XProjectUpdate_comments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, "", false, nil, err
	}

	return result.ProjectUpdate.Comments.Nodes,
		result.ProjectUpdate.Id,
		result.ProjectUpdate.Comments.PageInfo.HasNextPage,
		result.ProjectUpdate.Comments.PageInfo.EndCursor,
		nil
}

func projectUpdateNodeSummary(update projectUpdatesNode) ProjectUpdateSummary {
	return projectUpdateSummary(update.TopLevelProjectUpdateSummaryFields)
}

func projectUpdateCommentNodeSummary(node projectUpdateCommentsNode) CommentMetadataSummary {
	return commentMetadataSummary(node.CommentMetadataFields)
}

func projectScopedProjectUpdateSummary(update projectScopedUpdatesNode) ProjectUpdateSummary {
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
