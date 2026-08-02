package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// TeamCreateRequest describes a guarded Team create. A Team is organization-owned
// and is the pin's own subject, so a create necessarily lands outside the Pinned
// Target's team: OrgWide must be true and the request is refused otherwise.
type TeamCreateRequest struct {
	Name        string
	Key         string
	Description string
	Private     bool
	OrgWide     bool
}

// CreateTeam creates a Team after confirming --org-wide was passed, then confirms the
// created Team belongs to the Resolved Target's organization. This is an Org-Scoped
// Write: a Team has no team to compare against, and it cannot compare against the
// pinned team because it is the kind of thing a pin names. Organization membership is
// therefore the whole check, and it is made against the created Team rather than
// against target resolution alone, so the guard reports what Linear actually did.
func CreateTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request TeamCreateRequest,
) (TeamSummary, error) {
	if request.Name == "" {
		return TeamSummary{}, requiredFieldError("name")
	}
	if !request.OrgWide {
		return TeamSummary{}, fmt.Errorf(
			"%w: a Team is organization-owned and cannot belong to the pinned team; pass --org-wide to create one",
			ErrTargetMismatch,
		)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return TeamSummary{}, err
	}

	return guard.createTeam(ctx, request)
}

func (guard *guardedClient) createTeam(
	ctx context.Context,
	request TeamCreateRequest,
) (TeamSummary, error) {
	input := LinearTeamCreateInput{
		Name:        request.Name,
		Key:         optionalString(request.Key),
		Description: optionalString(request.Description),
	}
	if request.Private {
		input.Private = &request.Private
	}

	created, err := gql.TeamCreate(ctx, guard.graphqlClient, input)
	if err != nil {
		return TeamSummary{}, fmt.Errorf("create team: %w", err)
	}
	if !created.TeamCreate.Success || created.TeamCreate.Team == nil {
		return TeamSummary{}, fmt.Errorf("%w: teamCreate failed", ErrMutationFailed)
	}

	summary := teamSummary(created.TeamCreate.Team.TeamSummaryFields)
	if err := guard.requireOrganization(summary.OrgID); err != nil {
		return TeamSummary{}, err
	}

	return summary, nil
}
