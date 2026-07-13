package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

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

// LinearProjectLabelCreateInput is the sparse Linear projectLabelCreate payload linctl supports.
type LinearProjectLabelCreateInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// LinearProjectLabelUpdateInput is the sparse Linear projectLabelUpdate payload linctl supports.
type LinearProjectLabelUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// CreateProjectLabel creates a ProjectLabel after confirming --org-wide was
// passed. ProjectLabel is organization-owned; there is no team scope to
// materialize, so the Org-Scoped Write comparison is provided by
// guardedMutation's ResolveTarget call alone.
func CreateProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectLabelCreateRequest,
) (ProjectLabelSummary, error) {
	if request.Name == "" {
		return ProjectLabelSummary{}, fmt.Errorf("%w: name is required", ErrWriteInvalid)
	}
	if err := requireProjectLabelOrgWide(request.OrgWide); err != nil {
		return ProjectLabelSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(_ writeGuard) (ProjectLabelSummary, error) {
		created, err := ProjectLabelCreate(ctx, graphqlClient, LinearProjectLabelCreateInput{
			Name:        request.Name,
			Description: optionalString(request.Description),
			Color:       optionalString(request.Color),
		})
		if err != nil {
			return ProjectLabelSummary{}, fmt.Errorf("create project label: %w", err)
		}
		if !created.ProjectLabelCreate.Success {
			return ProjectLabelSummary{}, fmt.Errorf("%w: projectLabelCreate reported no success", ErrMutationFailed)
		}

		return projectLabelSummary(created.ProjectLabelCreate.ProjectLabel.ProjectLabelSummaryFields), nil
	})
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectLabelSummary, error) {
		if err := guard.requireProjectLabel(ctx, graphqlClient, request.ID); err != nil {
			return ProjectLabelSummary{}, err
		}

		updated, err := ProjectLabelUpdate(ctx, graphqlClient, request.ID, LinearProjectLabelUpdateInput{
			Name:        optionalString(request.Name),
			Description: optionalString(request.Description),
			Color:       optionalString(request.Color),
		})
		if err != nil {
			return ProjectLabelSummary{}, fmt.Errorf("update project label %s: %w", request.ID, err)
		}
		if !updated.ProjectLabelUpdate.Success {
			return ProjectLabelSummary{}, fmt.Errorf("%w: projectLabelUpdate reported no success", ErrMutationFailed)
		}

		return projectLabelSummary(updated.ProjectLabelUpdate.ProjectLabel.ProjectLabelSummaryFields), nil
	})
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
		return ProjectLabelSummary{}, fmt.Errorf("%w: project label id is required", ErrWriteInvalid)
	}
	if err := requireProjectLabelOrgWide(orgWide); err != nil {
		return ProjectLabelSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectLabelSummary, error) {
		if err := guard.requireProjectLabel(ctx, graphqlClient, id); err != nil {
			return ProjectLabelSummary{}, err
		}

		retired, err := ProjectLabelRetire(ctx, graphqlClient, id)
		if err != nil {
			return ProjectLabelSummary{}, fmt.Errorf("retire project label %s: %w", id, err)
		}
		if !retired.ProjectLabelRetire.Success {
			return ProjectLabelSummary{}, fmt.Errorf("%w: projectLabelRetire reported no success", ErrMutationFailed)
		}

		return projectLabelSummary(retired.ProjectLabelRetire.ProjectLabel.ProjectLabelSummaryFields), nil
	})
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
		return ProjectLabelSummary{}, fmt.Errorf("%w: project label id is required", ErrWriteInvalid)
	}
	if err := requireProjectLabelOrgWide(orgWide); err != nil {
		return ProjectLabelSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectLabelSummary, error) {
		if err := guard.requireProjectLabel(ctx, graphqlClient, id); err != nil {
			return ProjectLabelSummary{}, err
		}

		restored, err := ProjectLabelRestore(ctx, graphqlClient, id)
		if err != nil {
			return ProjectLabelSummary{}, fmt.Errorf("restore project label %s: %w", id, err)
		}
		if !restored.ProjectLabelRestore.Success {
			return ProjectLabelSummary{}, fmt.Errorf("%w: projectLabelRestore reported no success", ErrMutationFailed)
		}

		return projectLabelSummary(restored.ProjectLabelRestore.ProjectLabel.ProjectLabelSummaryFields), nil
	})
}

func validateProjectLabelUpdateRequest(request ProjectLabelUpdateRequest) error {
	if request.ID == "" {
		return fmt.Errorf("%w: project label id is required", ErrWriteInvalid)
	}
	if request.Name == "" && request.Description == "" && request.Color == "" {
		return fmt.Errorf("%w: name, description, or color is required", ErrWriteInvalid)
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
