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

// commentResolver is the Command Port the comment resolve command depends on.
type commentResolver interface {
	ResolveComment(ctx context.Context, commentID string) (client.CommentSummary, error)
}

// commentUnresolver is the Command Port the comment unresolve command depends on.
type commentUnresolver interface {
	UnresolveComment(ctx context.Context, commentID string) (client.CommentSummary, error)
}

var (
	_ commentUpdater    = commandClientAdapter{}
	_ commentDeleter    = commandClientAdapter{}
	_ commentResolver   = commandClientAdapter{}
	_ commentUnresolver = commandClientAdapter{}
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

func (adapter commandClientAdapter) ResolveComment(
	ctx context.Context,
	commentID string,
) (client.CommentSummary, error) {
	return client.ResolveComment(ctx, adapter.graphqlClient, adapter.target, commentID)
}

func (adapter commandClientAdapter) UnresolveComment(
	ctx context.Context,
	commentID string,
) (client.CommentSummary, error) {
	return client.UnresolveComment(ctx, adapter.graphqlClient, adapter.target, commentID)
}
