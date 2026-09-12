package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// InitiativeCreateRequest describes a guarded Initiative create. An Initiative is
// organization-owned and sits above every team, so a create necessarily lands
// outside the Pinned Target's team: OrgWide must be true and the request is
// refused otherwise.
type InitiativeCreateRequest struct {
	Name        string
	Description string
	Status      string
	OrgWide     bool
}

// InitiativeStatuses lists the InitiativeStatus values Linear accepts on create.
var InitiativeStatuses = []string{"Proposed", "Planned", "Active", "Completed", "Canceled"}

// CreateInitiative creates an Initiative after confirming --org-wide was passed,
// then confirms the created Initiative belongs to the Resolved Target's
// organization. This is an Org-Scoped Write: an Initiative has no team to compare
// against, so organization membership is the whole check, and it is made against
// the created Initiative rather than against target resolution alone.
func CreateInitiative(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request InitiativeCreateRequest,
) (InitiativeSummary, error) {
	if request.Name == "" {
		return InitiativeSummary{}, requiredFieldError("name")
	}
	if request.Status != "" && !isInitiativeStatus(request.Status) {
		return InitiativeSummary{}, fmt.Errorf(
			"%w: status must be one of %s", ErrWriteInvalid, strings.Join(InitiativeStatuses, ", "),
		)
	}
	if !request.OrgWide {
		return InitiativeSummary{}, fmt.Errorf(
			"%w: an Initiative is organization-owned and cannot belong to the pinned team; "+
				"pass --org-wide to create one",
			ErrTargetMismatch,
		)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return InitiativeSummary{}, err
	}

	return guard.createInitiative(ctx, request)
}

func isInitiativeStatus(status string) bool {
	for _, candidate := range InitiativeStatuses {
		if candidate == status {
			return true
		}
	}
	return false
}

func (guard *guardedClient) createInitiative(
	ctx context.Context,
	request InitiativeCreateRequest,
) (InitiativeSummary, error) {
	input := LinearInitiativeCreateInput{
		Name:        request.Name,
		Description: optionalString(request.Description),
		Status:      optionalString(request.Status),
	}

	created, err := gql.InitiativeCreate(ctx, guard.graphqlClient, input)
	if err != nil {
		return InitiativeSummary{}, fmt.Errorf("create initiative: %w", err)
	}
	if !created.InitiativeCreate.Success {
		return InitiativeSummary{}, fmt.Errorf("%w: initiativeCreate failed", ErrMutationFailed)
	}

	summary := initiativeSummary(created.InitiativeCreate.Initiative.InitiativeSummaryFields)
	if err := guard.requireOrganization(summary.OrgID); err != nil {
		return InitiativeSummary{}, err
	}

	return summary, nil
}
