package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// InitiativeUpdateCreateRequest describes a guarded initiative status-update create.
type InitiativeUpdateCreateRequest struct {
	InitiativeID string
	Body         string
	Health       string
}

// CreateInitiativeUpdate posts a status update to an initiative after resolving
// the initiative and comparing its organization to the pinned target.
func CreateInitiativeUpdate(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request InitiativeUpdateCreateRequest,
) (InitiativeUpdateSummary, error) {
	if err := validateInitiativeUpdateCreate(request); err != nil {
		return InitiativeUpdateSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return InitiativeUpdateSummary{}, err
	}

	return guard.createInitiativeUpdate(ctx, request)
}

func validateInitiativeUpdateCreate(request InitiativeUpdateCreateRequest) error {
	if request.InitiativeID == "" {
		return fmt.Errorf("%w: initiative id is required", ErrWriteInvalid)
	}
	if request.Body == "" && request.Health == "" {
		return fmt.Errorf("%w: body or health is required", ErrWriteInvalid)
	}

	return nil
}

func (guard *guardedClient) createInitiativeUpdate(
	ctx context.Context,
	request InitiativeUpdateCreateRequest,
) (InitiativeUpdateSummary, error) {
	if err := guard.requireInitiative(ctx, request.InitiativeID); err != nil {
		return InitiativeUpdateSummary{}, err
	}

	created, err := gql.InitiativeUpdateCreate(ctx, guard.graphqlClient, LinearInitiativeUpdateCreateInput{
		InitiativeID: request.InitiativeID,
		Body:         optionalString(request.Body),
		Health:       optionalString(request.Health),
	})
	if err != nil {
		return InitiativeUpdateSummary{}, fmt.Errorf("create initiative update: %w", err)
	}
	if !created.InitiativeUpdateCreate.Success {
		return InitiativeUpdateSummary{}, fmt.Errorf(
			"%w: initiativeUpdateCreate returned no update",
			ErrMutationFailed,
		)
	}

	return initiativeUpdateSummary(
		created.InitiativeUpdateCreate.InitiativeUpdate.InitiativeUpdateSummaryFields,
	), nil
}
