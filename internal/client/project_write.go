package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectCreateRequest describes a guarded project create.
type ProjectCreateRequest struct {
	Name        string
	Description string
}

// ProjectUpdateRequest describes a guarded project update.
type ProjectUpdateRequest struct {
	ID          string
	Name        string
	Description string
}

// LinearProjectCreateInput is the sparse Linear projectCreate payload linctl supports.
type LinearProjectCreateInput struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	TeamIDs     []string `json:"teamIds"`
}

// LinearProjectUpdateInput is the sparse Linear projectUpdate payload linctl supports.
type LinearProjectUpdateInput struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	TeamIDs     []string `json:"teamIds,omitempty"`
}

// CreateProject creates a team-scoped project after target comparison.
func CreateProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectCreateRequest,
) (ProjectSummary, error) {
	if request.Name == "" {
		return ProjectSummary{}, fmt.Errorf("%w: name is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectSummary, error) {
		created, err := ProjectCreate(ctx, graphqlClient, LinearProjectCreateInput{
			Name:        request.Name,
			Description: optionalString(request.Description),
			TeamIDs:     []string{guard.target.Team.ID},
		})
		if err != nil {
			return ProjectSummary{}, fmt.Errorf("create project: %w", err)
		}
		if !created.ProjectCreate.Success || created.ProjectCreate.Project == nil {
			return ProjectSummary{}, fmt.Errorf("%w: projectCreate returned no project", ErrMutationFailed)
		}

		return projectSummaryFromFields(created.ProjectCreate.Project.ProjectSummaryFields), nil
	})
}

// UpdateProject updates a resource-scoped project after target comparison.
func UpdateProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectUpdateRequest,
) (ProjectSummary, error) {
	if err := validateProjectUpdateRequest(request); err != nil {
		return ProjectSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectSummary, error) {
		if err := guard.requireProject(ctx, graphqlClient, request.ID); err != nil {
			return ProjectSummary{}, err
		}

		updated, err := ProjectUpdate(ctx, graphqlClient, request.ID, LinearProjectUpdateInput{
			Name:        optionalString(request.Name),
			Description: optionalString(request.Description),
		})
		if err != nil {
			return ProjectSummary{}, fmt.Errorf("update project %s: %w", request.ID, err)
		}
		if !updated.ProjectUpdate.Success || updated.ProjectUpdate.Project == nil {
			return ProjectSummary{}, fmt.Errorf("%w: projectUpdate returned no project", ErrMutationFailed)
		}

		return projectSummaryFromFields(updated.ProjectUpdate.Project.ProjectSummaryFields), nil
	})
}

// ArchiveProject archives a resource-scoped project after target comparison.
func ArchiveProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	projectID string,
) (ProjectSummary, error) {
	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectSummary, error) {
		if err := guard.requireProject(ctx, graphqlClient, projectID); err != nil {
			return ProjectSummary{}, err
		}

		archived, err := ProjectArchive(ctx, graphqlClient, projectID, boolPtr(false))
		if err != nil {
			return ProjectSummary{}, fmt.Errorf("archive project %s: %w", projectID, err)
		}
		if !archived.ProjectArchive.Success || archived.ProjectArchive.Entity == nil {
			return ProjectSummary{}, fmt.Errorf("%w: projectArchive returned no project", ErrMutationFailed)
		}

		return projectSummaryFromFields(archived.ProjectArchive.Entity.ProjectSummaryFields), nil
	})
}

func validateProjectUpdateRequest(request ProjectUpdateRequest) error {
	if request.ID == "" {
		return fmt.Errorf("%w: project id is required", ErrWriteInvalid)
	}
	if request.Name == "" && request.Description == "" {
		return fmt.Errorf("%w: name or description is required", ErrWriteInvalid)
	}

	return nil
}
