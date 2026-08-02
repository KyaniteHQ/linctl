package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectCreateRequest describes a guarded project create.
type ProjectCreateRequest struct {
	Name        string
	Description string
	Content     string
}

// ProjectUpdateRequest describes a guarded project update.
type ProjectUpdateRequest struct {
	ID          string
	Name        string
	Description string
	Content     string
}

// CreateProject creates a team-scoped project after target comparison.
func CreateProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectCreateRequest,
) (ProjectSummary, error) {
	if request.Name == "" {
		return ProjectSummary{}, requiredFieldError("name")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectSummary{}, err
	}

	return guard.createProject(ctx, request)
}

func (guard *guardedClient) createProject(ctx context.Context, request ProjectCreateRequest) (ProjectSummary, error) {
	created, err := gql.ProjectCreate(ctx, guard.graphqlClient, LinearProjectCreateInput{
		Name:        request.Name,
		Description: optionalString(request.Description),
		Content:     optionalString(request.Content),
		TeamIDs:     []string{guard.target.Team.ID},
	})
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("create project: %w", err)
	}
	if !created.ProjectCreate.Success || created.ProjectCreate.Project == nil {
		return ProjectSummary{}, fmt.Errorf("%w: projectCreate returned no project", ErrMutationFailed)
	}

	return projectSummaryFromFields(created.ProjectCreate.Project.ProjectSummaryFields), nil
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

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectSummary{}, err
	}

	return guard.updateProject(ctx, request)
}

func (guard *guardedClient) updateProject(ctx context.Context, request ProjectUpdateRequest) (ProjectSummary, error) {
	if err := guard.requireProject(ctx, request.ID); err != nil {
		return ProjectSummary{}, err
	}

	updated, err := gql.ProjectUpdate(ctx, guard.graphqlClient, request.ID, LinearProjectUpdateInput{
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		Content:     optionalString(request.Content),
	})
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("update project %s: %w", request.ID, err)
	}
	if !updated.ProjectUpdate.Success || updated.ProjectUpdate.Project == nil {
		return ProjectSummary{}, fmt.Errorf("%w: projectUpdate returned no project", ErrMutationFailed)
	}

	return projectSummaryFromFields(updated.ProjectUpdate.Project.ProjectSummaryFields), nil
}

// ArchiveProject archives a resource-scoped project after target comparison.
func ArchiveProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	projectID string,
) (ProjectSummary, error) {
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectSummary{}, err
	}

	return guard.archiveProject(ctx, projectID)
}

func (guard *guardedClient) archiveProject(ctx context.Context, projectID string) (ProjectSummary, error) {
	if err := guard.requireProject(ctx, projectID); err != nil {
		return ProjectSummary{}, err
	}

	archived, err := gql.ProjectArchive(ctx, guard.graphqlClient, projectID, boolPtr(false))
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("archive project %s: %w", projectID, err)
	}
	if !archived.ProjectArchive.Success || archived.ProjectArchive.Entity == nil {
		return ProjectSummary{}, fmt.Errorf("%w: projectArchive returned no project", ErrMutationFailed)
	}

	return projectSummaryFromFields(archived.ProjectArchive.Entity.ProjectSummaryFields), nil
}

func validateProjectUpdateRequest(request ProjectUpdateRequest) error {
	if request.ID == "" {
		return requiredFieldError("project id")
	}
	if request.Name == "" && request.Description == "" && request.Content == "" {
		return requiredFieldError("name, description, or content")
	}

	return nil
}
