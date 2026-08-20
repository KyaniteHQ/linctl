package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// InitiativeUpdateSummary is one initiative status update.
type InitiativeUpdateSummary struct {
	ID             string `json:"id"`
	Body           string `json:"body"`
	Health         string `json:"health"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	URL            string `json:"url"`
	SlugID         string `json:"slug_id"`
	CommentCount   int    `json:"comment_count"`
	InitiativeID   string `json:"initiative_id"`
	InitiativeName string `json:"initiative_name"`
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
}

// InitiativeUpdateList is a page of initiative status updates.
type InitiativeUpdateList struct {
	Updates     []InitiativeUpdateSummary `json:"updates"`
	HasNextPage bool                      `json:"has_next_page"`
	EndCursor   *string                   `json:"end_cursor,omitempty"`
}

// InitiativeUpdateCommentList is a page of body-free Comments associated with one InitiativeUpdate.
type InitiativeUpdateCommentList struct {
	InitiativeUpdateID string                   `json:"initiative_update_id"`
	Comments           []CommentMetadataSummary `json:"comments"`
	HasNextPage        bool                     `json:"has_next_page"`
	EndCursor          *string                  `json:"end_cursor,omitempty"`
}

//nolint:lll
type initiativeUpdatesNode = gql.XInitiativeUpdatesInitiativeUpdatesInitiativeUpdateConnectionNodesInitiativeUpdate

//nolint:lll
type initiativeUpdateCommentsNode = gql.XInitiativeUpdate_commentsInitiativeUpdateCommentsCommentConnectionNodesComment

type initiativeUpdatesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type initiativeUpdateCommentsQuery struct {
	ctx                context.Context
	graphqlClient      graphql.Client
	id                 string
	initiativeUpdateID string
}

// ListInitiativeUpdates returns visible initiative status updates.
func ListInitiativeUpdates(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeUpdateList, error) {
	query := initiativeUpdatesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list initiative updates", limit, defaultListPageSize,
		query.page,
		initiativeUpdatesNodeSummary,
	)
	if err != nil {
		return InitiativeUpdateList{}, err
	}

	return InitiativeUpdateList{Updates: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetInitiativeUpdateByID returns one initiative update by Linear id.
func GetInitiativeUpdateByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeUpdateSummary, error) {
	result, err := gql.XInitiativeUpdate(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeUpdateSummary{}, fmt.Errorf("get initiative update %s: %w", id, err)
	}

	return initiativeUpdateSummary(result.InitiativeUpdate.InitiativeUpdateSummaryFields), nil
}

// ListInitiativeUpdateComments returns body-free comments associated with one InitiativeUpdate.
func ListInitiativeUpdateComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeUpdateCommentList, error) {
	query := &initiativeUpdateCommentsQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative update comments "+id, limit, defaultListPageSize,
		query.comments,
		initiativeUpdateCommentsNodeSummary,
	)
	if err != nil {
		return InitiativeUpdateCommentList{}, err
	}

	return InitiativeUpdateCommentList{
		InitiativeUpdateID: query.initiativeUpdateID,
		Comments:           page.Items,
		HasNextPage:        page.HasNextPage,
		EndCursor:          page.EndCursor,
	}, nil
}

func (query initiativeUpdatesQuery) page(
	pageSize int,
	after *string,
) ([]initiativeUpdatesNode, bool, *string, error) {
	result, err := gql.XInitiativeUpdates(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.InitiativeUpdates.Nodes,
		result.InitiativeUpdates.PageInfo.HasNextPage,
		result.InitiativeUpdates.PageInfo.EndCursor,
		nil
}

func (query *initiativeUpdateCommentsQuery) comments(
	pageSize int,
	after *string,
) ([]initiativeUpdateCommentsNode, bool, *string, error) {
	result, err := gql.XInitiativeUpdate_comments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.initiativeUpdateID = result.InitiativeUpdate.Id

	return result.InitiativeUpdate.Comments.Nodes,
		result.InitiativeUpdate.Comments.PageInfo.HasNextPage,
		result.InitiativeUpdate.Comments.PageInfo.EndCursor,
		nil
}

func initiativeUpdatesNodeSummary(update initiativeUpdatesNode) InitiativeUpdateSummary {
	return initiativeUpdateSummary(update.InitiativeUpdateSummaryFields)
}

func initiativeUpdateCommentsNodeSummary(node initiativeUpdateCommentsNode) CommentMetadataSummary {
	return commentMetadataSummary(node.CommentMetadataFields)
}

func initiativeUpdateSummary(update gql.InitiativeUpdateSummaryFields) InitiativeUpdateSummary {
	return InitiativeUpdateSummary{
		ID:             update.Id,
		Body:           update.Body,
		Health:         string(update.Health),
		CreatedAt:      update.CreatedAt,
		UpdatedAt:      update.UpdatedAt,
		URL:            update.Url,
		SlugID:         update.SlugId,
		CommentCount:   update.CommentCount,
		InitiativeID:   update.Initiative.Id,
		InitiativeName: update.Initiative.Name,
		UserID:         update.User.Id,
		Name:           update.User.Name,
		DisplayName:    update.User.DisplayName,
	}
}
