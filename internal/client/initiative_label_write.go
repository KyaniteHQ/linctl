package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// RetireInitiativeLabel retires an InitiativeLabel after confirming --org-wide
// was passed and that the label belongs to the resolved organization.
func RetireInitiativeLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
	orgWide bool,
) (InitiativeLabelSummary, error) {
	if id == "" {
		return InitiativeLabelSummary{}, fmt.Errorf("%w: initiative label id is required", ErrWriteInvalid)
	}
	if err := requireInitiativeLabelOrgWide(orgWide); err != nil {
		return InitiativeLabelSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return InitiativeLabelSummary{}, err
	}

	return guard.retireInitiativeLabel(ctx, id)
}

func (guard *guardedClient) retireInitiativeLabel(
	ctx context.Context,
	id string,
) (InitiativeLabelSummary, error) {
	if err := guard.requireInitiativeLabel(ctx, id); err != nil {
		return InitiativeLabelSummary{}, err
	}

	retired, err := gql.InitiativeLabelRetire(ctx, guard.graphqlClient, id)
	if err != nil {
		return InitiativeLabelSummary{}, fmt.Errorf("retire initiative label %s: %w", id, err)
	}
	if err := mutationSuccess(retired.InitiativeLabelRetire.Success, "initiativeLabelRetire"); err != nil {
		return InitiativeLabelSummary{}, err
	}

	return initiativeLabelSummary(retired.InitiativeLabelRetire.InitiativeLabel.InitiativeLabelSummaryFields), nil
}

// RestoreInitiativeLabel restores a previously retired InitiativeLabel after
// confirming --org-wide was passed and that the label belongs to the
// resolved organization.
func RestoreInitiativeLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
	orgWide bool,
) (InitiativeLabelSummary, error) {
	if id == "" {
		return InitiativeLabelSummary{}, fmt.Errorf("%w: initiative label id is required", ErrWriteInvalid)
	}
	if err := requireInitiativeLabelOrgWide(orgWide); err != nil {
		return InitiativeLabelSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return InitiativeLabelSummary{}, err
	}

	return guard.restoreInitiativeLabel(ctx, id)
}

func (guard *guardedClient) restoreInitiativeLabel(
	ctx context.Context,
	id string,
) (InitiativeLabelSummary, error) {
	if err := guard.requireInitiativeLabel(ctx, id); err != nil {
		return InitiativeLabelSummary{}, err
	}

	restored, err := gql.InitiativeLabelRestore(ctx, guard.graphqlClient, id)
	if err != nil {
		return InitiativeLabelSummary{}, fmt.Errorf("restore initiative label %s: %w", id, err)
	}
	if err := mutationSuccess(restored.InitiativeLabelRestore.Success, "initiativeLabelRestore"); err != nil {
		return InitiativeLabelSummary{}, err
	}

	return initiativeLabelSummary(restored.InitiativeLabelRestore.InitiativeLabel.InitiativeLabelSummaryFields), nil
}

// requireInitiativeLabelOrgWide refuses an InitiativeLabel taxonomy write that
// omits --org-wide. InitiativeLabel has no team scope, so --org-wide is the
// only available Org-Scoped Write path and is required rather than optional;
// the blast radius is every initiative in the organization.
func requireInitiativeLabelOrgWide(orgWide bool) error {
	if !orgWide {
		return fmt.Errorf(
			"%w: initiative labels have no team scope; pass --org-wide to confirm this affects "+
				"every initiative in the organization",
			ErrWriteInvalid,
		)
	}

	return nil
}
