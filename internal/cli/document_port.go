package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// documentCreator is the Command Port the document create command depends on.
type documentCreator interface {
	CreateDocument(ctx context.Context, request client.DocumentCreateRequest) (client.DocumentSummary, error)
}

// documentUpdater is the Command Port the document update command depends on.
type documentUpdater interface {
	UpdateDocument(ctx context.Context, request client.DocumentUpdateRequest) (client.DocumentSummary, error)
}

var (
	_ documentCreator = commandClientAdapter{}
	_ documentUpdater = commandClientAdapter{}
)

func (adapter commandClientAdapter) CreateDocument(
	ctx context.Context,
	request client.DocumentCreateRequest,
) (client.DocumentSummary, error) {
	return client.CreateDocument(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) UpdateDocument(
	ctx context.Context,
	request client.DocumentUpdateRequest,
) (client.DocumentSummary, error) {
	return client.UpdateDocument(ctx, adapter.graphqlClient, adapter.target, request)
}
