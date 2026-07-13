package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// LabelCreateRequest describes a guarded IssueLabel create in the pinned
// team, or organization-wide when OrgWide is set.
type LabelCreateRequest struct {
	Name        string
	Color       string
	Description string
	ParentID    string
	OrgWide     bool
}

// LabelUpdateRequest describes a guarded IssueLabel update.
type LabelUpdateRequest struct {
	ID          string
	Name        string
	Color       string
	Description string
	OrgWide     bool
}

// LinearIssueLabelCreateInput is the sparse Linear issueLabelCreate payload linctl supports.
type LinearIssueLabelCreateInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	TeamID      *string `json:"teamId,omitempty"`
}

// LinearIssueLabelUpdateInput is the sparse Linear issueLabelUpdate payload linctl supports.
type LinearIssueLabelUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Color       *string `json:"color,omitempty"`
}

// CreateLabel creates an IssueLabel in the pinned team, or organization-wide
// when OrgWide is set, after target comparison. replaceTeamLabels is never
// sent true.
func CreateLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request LabelCreateRequest,
) (LabelSummary, error) {
	if request.Name == "" {
		return LabelSummary{}, fmt.Errorf("%w: name is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (LabelSummary, error) {
		var teamID *string
		if !request.OrgWide {
			teamID = stringPtr(guard.target.Team.ID)
		}
		if request.ParentID != "" {
			if err := guard.requireLabelParentScope(ctx, graphqlClient, request.ParentID, request.OrgWide); err != nil {
				return LabelSummary{}, err
			}
		}

		created, err := IssueLabelCreate(ctx, graphqlClient, boolPtr(false), LinearIssueLabelCreateInput{
			Name:        request.Name,
			Description: optionalString(request.Description),
			Color:       optionalString(request.Color),
			ParentID:    optionalString(request.ParentID),
			TeamID:      teamID,
		})
		if err != nil {
			return LabelSummary{}, fmt.Errorf("create label: %w", err)
		}
		if !created.IssueLabelCreate.Success {
			return LabelSummary{}, fmt.Errorf("%w: issueLabelCreate reported no success", ErrMutationFailed)
		}

		return labelSummary(created.IssueLabelCreate.IssueLabel.IssueLabelSummaryFields), nil
	})
}

// UpdateLabel updates an IssueLabel after resolving and comparing its scope:
// a team-scoped label compares team id+key, and an organization-wide label
// (null team) requires OrgWide and fails closed otherwise.
func UpdateLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request LabelUpdateRequest,
) (LabelSummary, error) {
	if err := validateLabelUpdateRequest(request); err != nil {
		return LabelSummary{}, err
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (LabelSummary, error) {
		if err := guard.requireIssueLabel(ctx, graphqlClient, request.ID, request.OrgWide); err != nil {
			return LabelSummary{}, err
		}

		updated, err := IssueLabelUpdate(ctx, graphqlClient, request.ID, boolPtr(false), LinearIssueLabelUpdateInput{
			Name:        optionalString(request.Name),
			Description: optionalString(request.Description),
			Color:       optionalString(request.Color),
		})
		if err != nil {
			return LabelSummary{}, fmt.Errorf("update label %s: %w", request.ID, err)
		}
		if !updated.IssueLabelUpdate.Success {
			return LabelSummary{}, fmt.Errorf("%w: issueLabelUpdate reported no success", ErrMutationFailed)
		}

		return labelSummary(updated.IssueLabelUpdate.IssueLabel.IssueLabelSummaryFields), nil
	})
}

// RetireLabel retires an IssueLabel after resolving and comparing its scope.
func RetireLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
	orgWide bool,
) (LabelSummary, error) {
	if id == "" {
		return LabelSummary{}, fmt.Errorf("%w: label id is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (LabelSummary, error) {
		if err := guard.requireIssueLabel(ctx, graphqlClient, id, orgWide); err != nil {
			return LabelSummary{}, err
		}

		retired, err := IssueLabelRetire(ctx, graphqlClient, id)
		if err != nil {
			return LabelSummary{}, fmt.Errorf("retire label %s: %w", id, err)
		}
		if !retired.IssueLabelRetire.Success {
			return LabelSummary{}, fmt.Errorf("%w: issueLabelRetire reported no success", ErrMutationFailed)
		}

		return labelSummary(retired.IssueLabelRetire.IssueLabel.IssueLabelSummaryFields), nil
	})
}

// RestoreLabel restores a previously retired IssueLabel after resolving and
// comparing its scope.
func RestoreLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
	orgWide bool,
) (LabelSummary, error) {
	if id == "" {
		return LabelSummary{}, fmt.Errorf("%w: label id is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (LabelSummary, error) {
		if err := guard.requireIssueLabel(ctx, graphqlClient, id, orgWide); err != nil {
			return LabelSummary{}, err
		}

		restored, err := IssueLabelRestore(ctx, graphqlClient, id)
		if err != nil {
			return LabelSummary{}, fmt.Errorf("restore label %s: %w", id, err)
		}
		if !restored.IssueLabelRestore.Success {
			return LabelSummary{}, fmt.Errorf("%w: issueLabelRestore reported no success", ErrMutationFailed)
		}

		return labelSummary(restored.IssueLabelRestore.IssueLabel.IssueLabelSummaryFields), nil
	})
}

func validateLabelUpdateRequest(request LabelUpdateRequest) error {
	if request.ID == "" {
		return fmt.Errorf("%w: label id is required", ErrWriteInvalid)
	}
	if request.Name == "" && request.Description == "" && request.Color == "" {
		return fmt.Errorf("%w: name, description, or color is required", ErrWriteInvalid)
	}

	return nil
}
