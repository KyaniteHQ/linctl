package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// projectCreator is the Command Port the project create command depends on.
type projectCreator interface {
	CreateProject(ctx context.Context, request client.ProjectCreateRequest) (client.ProjectSummary, error)
}

// projectUpdater is the Command Port the project update command depends on.
type projectUpdater interface {
	UpdateProject(ctx context.Context, request client.ProjectUpdateRequest) (client.ProjectSummary, error)
}

// projectArchiver is the Command Port the project archive command depends on.
type projectArchiver interface {
	ArchiveProject(ctx context.Context, projectID string) (client.ProjectSummary, error)
}

type projectClientAdapter struct {
	graphqlClient graphql.Client
	target        config.Target
}

func projectAdapterFor(runtime commandRuntime) projectClientAdapter {
	return projectClientAdapter{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (adapter projectClientAdapter) CreateProject(
	ctx context.Context,
	request client.ProjectCreateRequest,
) (client.ProjectSummary, error) {
	return client.CreateProject(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter projectClientAdapter) UpdateProject(
	ctx context.Context,
	request client.ProjectUpdateRequest,
) (client.ProjectSummary, error) {
	return client.UpdateProject(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter projectClientAdapter) ArchiveProject(
	ctx context.Context,
	projectID string,
) (client.ProjectSummary, error) {
	return client.ArchiveProject(ctx, adapter.graphqlClient, adapter.target, projectID)
}
