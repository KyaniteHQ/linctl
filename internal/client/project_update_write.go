package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectUpdateCreateRequest describes a guarded project status-update create.
type ProjectUpdateCreateRequest struct {
	ProjectID string
	Body      string
	Health    string
}

// CreateProjectUpdate posts a status update to a project after resolving and
// comparing the pinned target. It is resource-scoped: the target project must
// match the pinned project (when configured) and belong to the resolved team.
func CreateProjectUpdate(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectUpdateCreateRequest,
) (ProjectUpdateSummary, error) {
	if request.ProjectID == "" {
		return ProjectUpdateSummary{}, fmt.Errorf("%w: project id is required", ErrWriteInvalid)
	}
	if request.Body == "" && request.Health == "" {
		return ProjectUpdateSummary{}, fmt.Errorf("%w: body or health is required", ErrWriteInvalid)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectUpdateSummary{}, err
	}

	return guard.createProjectUpdate(ctx, request)
}

func (guard *guardedClient) createProjectUpdate(
	ctx context.Context,
	request ProjectUpdateCreateRequest,
) (ProjectUpdateSummary, error) {
	if err := guard.requireProject(ctx, request.ProjectID); err != nil {
		return ProjectUpdateSummary{}, err
	}

	created, err := gql.ProjectUpdateCreate(ctx, guard.graphqlClient, LinearProjectUpdateCreateInput{
		ProjectID: request.ProjectID,
		Body:      optionalString(request.Body),
		Health:    optionalString(request.Health),
	})
	if err != nil {
		return ProjectUpdateSummary{}, fmt.Errorf("create project update: %w", err)
	}
	if !created.ProjectUpdateCreate.Success {
		return ProjectUpdateSummary{}, fmt.Errorf("%w: projectUpdateCreate returned no update", ErrMutationFailed)
	}

	return projectUpdateSummary(created.ProjectUpdateCreate.ProjectUpdate.TopLevelProjectUpdateSummaryFields), nil
}
