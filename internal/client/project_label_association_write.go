package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectLabelAssociationRequest describes a guarded ProjectLabel attach/detach on a project.
type ProjectLabelAssociationRequest struct {
	ProjectID string
	LabelID   string
}

// AddProjectLabel attaches a ProjectLabel to a project after resolving and
// comparing both the project and the label against the pinned target. A
// pinned project_id does not block this taxonomy write; ProjectLabel is
// organization-owned so only the organization is compared for the label.
func AddProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectLabelAssociationRequest,
) (ProjectSummary, error) {
	if err := validateProjectLabelAssociationRequest(request); err != nil {
		return ProjectSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectSummary, error) {
		if err := resolveProjectLabelAssociation(ctx, graphqlClient, guard, request); err != nil {
			return ProjectSummary{}, err
		}

		updated, err := ProjectAddLabel(ctx, graphqlClient, request.ProjectID, request.LabelID)
		if err != nil {
			return ProjectSummary{}, fmt.Errorf("add label to project %s: %w", request.ProjectID, err)
		}
		if !updated.ProjectAddLabel.Success || updated.ProjectAddLabel.Project == nil {
			return ProjectSummary{}, fmt.Errorf("%w: projectAddLabel reported no success", ErrMutationFailed)
		}

		return projectSummaryFromFields(updated.ProjectAddLabel.Project.ProjectSummaryFields), nil
	})
}

// RemoveProjectLabel detaches a ProjectLabel from a project after resolving
// and comparing both the project and the label against the pinned target.
func RemoveProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectLabelAssociationRequest,
) (ProjectSummary, error) {
	if err := validateProjectLabelAssociationRequest(request); err != nil {
		return ProjectSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (ProjectSummary, error) {
		if err := resolveProjectLabelAssociation(ctx, graphqlClient, guard, request); err != nil {
			return ProjectSummary{}, err
		}

		updated, err := ProjectRemoveLabel(ctx, graphqlClient, request.ProjectID, request.LabelID)
		if err != nil {
			return ProjectSummary{}, fmt.Errorf("remove label from project %s: %w", request.ProjectID, err)
		}
		if !updated.ProjectRemoveLabel.Success || updated.ProjectRemoveLabel.Project == nil {
			return ProjectSummary{}, fmt.Errorf("%w: projectRemoveLabel reported no success", ErrMutationFailed)
		}

		return projectSummaryFromFields(updated.ProjectRemoveLabel.Project.ProjectSummaryFields), nil
	})
}

// resolveProjectLabelAssociation resolves the project and confirms the
// ProjectLabel belongs to the resolved organization, shared by
// AddProjectLabel and RemoveProjectLabel before either mutation is sent.
func resolveProjectLabelAssociation(
	ctx context.Context,
	graphqlClient graphql.Client,
	guard writeGuard,
	request ProjectLabelAssociationRequest,
) error {
	if err := guard.requireProject(ctx, graphqlClient, request.ProjectID); err != nil {
		return err
	}

	return guard.requireProjectLabel(ctx, graphqlClient, request.LabelID)
}

func validateProjectLabelAssociationRequest(request ProjectLabelAssociationRequest) error {
	if request.ProjectID == "" || request.LabelID == "" {
		return fmt.Errorf("%w: project id and label id are required", ErrWriteInvalid)
	}

	return nil
}
