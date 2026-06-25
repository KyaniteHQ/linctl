package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// commentUpdater is the Command Port the comment update command depends on.
type commentUpdater interface {
	UpdateComment(ctx context.Context, request client.CommentUpdateRequest) (client.CommentSummary, error)
}

// commentDeleter is the Command Port the comment delete command depends on.
type commentDeleter interface {
	DeleteComment(ctx context.Context, commentID string) (string, error)
}

type commentClientAdapter struct {
	graphqlClient graphql.Client
	target        config.Target
}

func commentAdapterFor(runtime commandRuntime) commentClientAdapter {
	return commentClientAdapter{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (adapter commentClientAdapter) UpdateComment(
	ctx context.Context,
	request client.CommentUpdateRequest,
) (client.CommentSummary, error) {
	return client.UpdateComment(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commentClientAdapter) DeleteComment(ctx context.Context, commentID string) (string, error) {
	return client.DeleteComment(ctx, adapter.graphqlClient, adapter.target, commentID)
}
