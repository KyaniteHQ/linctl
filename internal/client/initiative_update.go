package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListInitiativeUpdates returns visible initiative status updates.
func ListInitiativeUpdates(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeUpdateList, error) {
	result, err := initiativeUpdates(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeUpdateList{}, fmt.Errorf("list initiative updates: %w", err)
	}

	updates := make([]InitiativeUpdateSummary, 0, len(result.InitiativeUpdates.Nodes))
	for _, update := range result.InitiativeUpdates.Nodes {
		updates = append(updates, initiativeUpdateSummary(update.InitiativeUpdateSummaryFields))
	}

	return InitiativeUpdateList{
		Updates:     updates,
		HasNextPage: result.InitiativeUpdates.PageInfo.HasNextPage,
		EndCursor:   result.InitiativeUpdates.PageInfo.EndCursor,
	}, nil
}

// GetInitiativeUpdateByID returns one initiative update by Linear id.
func GetInitiativeUpdateByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeUpdateSummary, error) {
	result, err := initiativeUpdate(ctx, graphqlClient, id)
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
	result, err := initiativeUpdate_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeUpdateCommentList{}, fmt.Errorf("list initiative update comments %s: %w", id, err)
	}

	comments := make([]CommentMetadataSummary, 0, len(result.InitiativeUpdate.Comments.Nodes))
	for _, node := range result.InitiativeUpdate.Comments.Nodes {
		comments = append(comments, commentMetadataSummary(node.CommentMetadataFields))
	}

	return InitiativeUpdateCommentList{
		InitiativeUpdateID: result.InitiativeUpdate.Id,
		Comments:           comments,
		HasNextPage:        result.InitiativeUpdate.Comments.PageInfo.HasNextPage,
		EndCursor:          result.InitiativeUpdate.Comments.PageInfo.EndCursor,
	}, nil
}

func initiativeUpdateSummary(update InitiativeUpdateSummaryFields) InitiativeUpdateSummary {
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
