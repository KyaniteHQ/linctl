package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// commentUpdater is the Command Port the comment update command depends on.
type commentUpdater interface {
	UpdateComment(ctx context.Context, request client.CommentUpdateRequest) (client.CommentSummary, error)
}

// commentDeleter is the Command Port the comment delete command depends on.
type commentDeleter interface {
	DeleteComment(ctx context.Context, commentID string) (string, error)
}

var (
	_ commentUpdater = commandClientAdapter{}
	_ commentDeleter = commandClientAdapter{}
)

func (adapter commandClientAdapter) UpdateComment(
	ctx context.Context,
	request client.CommentUpdateRequest,
) (client.CommentSummary, error) {
	return client.UpdateComment(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) DeleteComment(ctx context.Context, commentID string) (string, error) {
	return client.DeleteComment(ctx, adapter.graphqlClient, adapter.target, commentID)
}
