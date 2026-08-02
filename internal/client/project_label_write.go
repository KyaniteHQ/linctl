package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectLabelCreateRequest describes a guarded ProjectLabel create.
// ProjectLabel has no team scope, so OrgWide must be true; the request is
// refused otherwise.
type ProjectLabelCreateRequest struct {
	Name        string
	Color       string
	Description string
	OrgWide     bool
}

// ProjectLabelUpdateRequest describes a guarded ProjectLabel update. OrgWide
// must be true for the same reason as ProjectLabelCreateRequest.
type ProjectLabelUpdateRequest struct {
	ID          string
	Name        string
	Color       string
	Description string
	OrgWide     bool
}

// CreateProjectLabel creates a ProjectLabel after confirming --org-wide was
// passed. ProjectLabel is organization-owned; there is no team scope to
// materialize, so the Org-Scoped Write comparison is provided by
// guarded client's target resolution alone.
func CreateProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectLabelCreateRequest,
) (ProjectLabelSummary, error) {
	if request.Name == "" {
		return ProjectLabelSummary{}, requiredFieldError("name")
	}
	if err := requireProjectLabelOrgWide(request.OrgWide); err != nil {
		return ProjectLabelSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectLabelSummary{}, err
	}

	return guard.createProjectLabel(ctx, request)
}

func (guard *guardedClient) createProjectLabel(
	ctx context.Context,
	request ProjectLabelCreateRequest,
) (ProjectLabelSummary, error) {
	created, err := gql.ProjectLabelCreate(ctx, guard.graphqlClient, LinearProjectLabelCreateInput{
		Name:        request.Name,
		Description: optionalString(request.Description),
		Color:       optionalString(request.Color),
	})
	if err != nil {
		return ProjectLabelSummary{}, fmt.Errorf("create project label: %w", err)
	}
	if err := mutationSuccess(created.ProjectLabelCreate.Success, "projectLabelCreate"); err != nil {
		return ProjectLabelSummary{}, err
	}

	return projectLabelSummary(created.ProjectLabelCreate.ProjectLabel.ProjectLabelSummaryFields), nil
}

// UpdateProjectLabel updates a ProjectLabel after confirming --org-wide was
// passed and that the label belongs to the resolved organization.
//
//nolint:dupl // Mirrors UpdateProjectMilestone's resolve-then-mutate shape; the guard target and mutation differ.
func UpdateProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectLabelUpdateRequest,
) (ProjectLabelSummary, error) {
	if err := validateProjectLabelUpdateRequest(request); err != nil {
		return ProjectLabelSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectLabelSummary{}, err
	}

	return guard.updateProjectLabel(ctx, request)
}

func (guard *guardedClient) updateProjectLabel(
	ctx context.Context,
	request ProjectLabelUpdateRequest,
) (ProjectLabelSummary, error) {
	if err := guard.requireProjectLabel(ctx, request.ID); err != nil {
		return ProjectLabelSummary{}, err
	}

	updated, err := gql.ProjectLabelUpdate(ctx, guard.graphqlClient, request.ID, LinearProjectLabelUpdateInput{
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		Color:       optionalString(request.Color),
	})
	if err != nil {
		return ProjectLabelSummary{}, fmt.Errorf("update project label %s: %w", request.ID, err)
	}
	if err := mutationSuccess(updated.ProjectLabelUpdate.Success, "projectLabelUpdate"); err != nil {
		return ProjectLabelSummary{}, err
	}

	return projectLabelSummary(updated.ProjectLabelUpdate.ProjectLabel.ProjectLabelSummaryFields), nil
}

// RetireProjectLabel retires a ProjectLabel after confirming --org-wide was
// passed and that the label belongs to the resolved organization.
func RetireProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
	orgWide bool,
) (ProjectLabelSummary, error) {
	if id == "" {
		return ProjectLabelSummary{}, requiredFieldError("project label id")
	}
	if err := requireProjectLabelOrgWide(orgWide); err != nil {
		return ProjectLabelSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectLabelSummary{}, err
	}

	return guard.retireProjectLabel(ctx, id)
}

func (guard *guardedClient) retireProjectLabel(ctx context.Context, id string) (ProjectLabelSummary, error) {
	if err := guard.requireProjectLabel(ctx, id); err != nil {
		return ProjectLabelSummary{}, err
	}

	retired, err := gql.ProjectLabelRetire(ctx, guard.graphqlClient, id)
	if err != nil {
		return ProjectLabelSummary{}, fmt.Errorf("retire project label %s: %w", id, err)
	}
	if err := mutationSuccess(retired.ProjectLabelRetire.Success, "projectLabelRetire"); err != nil {
		return ProjectLabelSummary{}, err
	}

	return projectLabelSummary(retired.ProjectLabelRetire.ProjectLabel.ProjectLabelSummaryFields), nil
}

// RestoreProjectLabel restores a previously retired ProjectLabel after
// confirming --org-wide was passed and that the label belongs to the
// resolved organization.
func RestoreProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
	orgWide bool,
) (ProjectLabelSummary, error) {
	if id == "" {
		return ProjectLabelSummary{}, requiredFieldError("project label id")
	}
	if err := requireProjectLabelOrgWide(orgWide); err != nil {
		return ProjectLabelSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectLabelSummary{}, err
	}

	return guard.restoreProjectLabel(ctx, id)
}

func (guard *guardedClient) restoreProjectLabel(ctx context.Context, id string) (ProjectLabelSummary, error) {
	if err := guard.requireProjectLabel(ctx, id); err != nil {
		return ProjectLabelSummary{}, err
	}

	restored, err := gql.ProjectLabelRestore(ctx, guard.graphqlClient, id)
	if err != nil {
		return ProjectLabelSummary{}, fmt.Errorf("restore project label %s: %w", id, err)
	}
	if err := mutationSuccess(restored.ProjectLabelRestore.Success, "projectLabelRestore"); err != nil {
		return ProjectLabelSummary{}, err
	}

	return projectLabelSummary(restored.ProjectLabelRestore.ProjectLabel.ProjectLabelSummaryFields), nil
}

func validateProjectLabelUpdateRequest(request ProjectLabelUpdateRequest) error {
	if request.ID == "" {
		return requiredFieldError("project label id")
	}
	if request.Name == "" && request.Description == "" && request.Color == "" {
		return requiredFieldError("name, description, or color")
	}

	return requireProjectLabelOrgWide(request.OrgWide)
}

// requireProjectLabelOrgWide refuses a ProjectLabel taxonomy write that omits
// --org-wide. ProjectLabel has no team scope, so --org-wide is the only
// available Org-Scoped Write path and is required rather than optional; the
// blast radius is every team and project in the organization.
func requireProjectLabelOrgWide(orgWide bool) error {
	if !orgWide {
		return fmt.Errorf(
			"%w: project labels have no team scope; pass --org-wide to confirm this affects "+
				"every team and project in the organization",
			ErrWriteInvalid,
		)
	}

	return nil
}
