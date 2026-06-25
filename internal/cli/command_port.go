package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

type commandClientAdapter struct {
	graphqlClient graphql.Client
	target        config.Target
}

func commandAdapterFor(runtime commandRuntime) commandClientAdapter {
	return commandClientAdapter{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (adapter commandClientAdapter) CreateCycle(
	ctx context.Context,
	request client.CycleCreateRequest,
) (client.CycleSummary, error) {
	return client.CreateCycle(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) UpdateCycle(
	ctx context.Context,
	request client.CycleUpdateRequest,
) (client.CycleSummary, error) {
	return client.UpdateCycle(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) ArchiveCycle(ctx context.Context, cycleID string) (client.CycleSummary, error) {
	return client.ArchiveCycle(ctx, adapter.graphqlClient, adapter.target, cycleID)
}

func (adapter commandClientAdapter) CreateProject(
	ctx context.Context,
	request client.ProjectCreateRequest,
) (client.ProjectSummary, error) {
	return client.CreateProject(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) UpdateProject(
	ctx context.Context,
	request client.ProjectUpdateRequest,
) (client.ProjectSummary, error) {
	return client.UpdateProject(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) ArchiveProject(
	ctx context.Context,
	projectID string,
) (client.ProjectSummary, error) {
	return client.ArchiveProject(ctx, adapter.graphqlClient, adapter.target, projectID)
}
