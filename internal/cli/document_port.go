package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// documentCreator is the Command Port the document create command depends on.
type documentCreator interface {
	CreateDocument(ctx context.Context, request client.DocumentCreateRequest) (client.DocumentSummary, error)
}

// documentUpdater is the Command Port the document update command depends on.
type documentUpdater interface {
	UpdateDocument(ctx context.Context, request client.DocumentUpdateRequest) (client.DocumentSummary, error)
}

type documentClientAdapter struct {
	graphqlClient graphql.Client
	target        config.Target
}

func documentAdapterFor(runtime commandRuntime) documentClientAdapter {
	return documentClientAdapter{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (adapter documentClientAdapter) CreateDocument(
	ctx context.Context,
	request client.DocumentCreateRequest,
) (client.DocumentSummary, error) {
	return client.CreateDocument(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter documentClientAdapter) UpdateDocument(
	ctx context.Context,
	request client.DocumentUpdateRequest,
) (client.DocumentSummary, error) {
	return client.UpdateDocument(ctx, adapter.graphqlClient, adapter.target, request)
}
