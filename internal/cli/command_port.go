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

func (adapter commandClientAdapter) CreateProjectMilestone(
	ctx context.Context,
	request client.ProjectMilestoneCreateRequest,
) (client.ProjectMilestoneSummary, error) {
	return client.CreateProjectMilestone(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) UpdateProjectMilestone(
	ctx context.Context,
	request client.ProjectMilestoneUpdateRequest,
) (client.ProjectMilestoneSummary, error) {
	return client.UpdateProjectMilestone(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) DeleteProjectMilestone(
	ctx context.Context,
	projectMilestoneID string,
) (string, error) {
	return client.DeleteProjectMilestone(ctx, adapter.graphqlClient, adapter.target, projectMilestoneID)
}

func (adapter commandClientAdapter) CreateLabel(
	ctx context.Context,
	request client.LabelCreateRequest,
) (client.LabelSummary, error) {
	return client.CreateLabel(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) UpdateLabel(
	ctx context.Context,
	request client.LabelUpdateRequest,
) (client.LabelSummary, error) {
	return client.UpdateLabel(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) RetireLabel(
	ctx context.Context,
	id string,
	orgWide bool,
) (client.LabelSummary, error) {
	return client.RetireLabel(ctx, adapter.graphqlClient, adapter.target, id, orgWide)
}

func (adapter commandClientAdapter) RestoreLabel(
	ctx context.Context,
	id string,
	orgWide bool,
) (client.LabelSummary, error) {
	return client.RestoreLabel(ctx, adapter.graphqlClient, adapter.target, id, orgWide)
}

func (adapter commandClientAdapter) AddProjectLabel(
	ctx context.Context,
	request client.ProjectLabelAssociationRequest,
) (client.ProjectSummary, error) {
	return client.AddProjectLabel(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) RemoveProjectLabel(
	ctx context.Context,
	request client.ProjectLabelAssociationRequest,
) (client.ProjectSummary, error) {
	return client.RemoveProjectLabel(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) CreateProjectLabel(
	ctx context.Context,
	request client.ProjectLabelCreateRequest,
) (client.ProjectLabelSummary, error) {
	return client.CreateProjectLabel(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) UpdateProjectLabel(
	ctx context.Context,
	request client.ProjectLabelUpdateRequest,
) (client.ProjectLabelSummary, error) {
	return client.UpdateProjectLabel(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter commandClientAdapter) RetireProjectLabel(
	ctx context.Context,
	id string,
	orgWide bool,
) (client.ProjectLabelSummary, error) {
	return client.RetireProjectLabel(ctx, adapter.graphqlClient, adapter.target, id, orgWide)
}

func (adapter commandClientAdapter) RestoreProjectLabel(
	ctx context.Context,
	id string,
	orgWide bool,
) (client.ProjectLabelSummary, error) {
	return client.RestoreProjectLabel(ctx, adapter.graphqlClient, adapter.target, id, orgWide)
}
