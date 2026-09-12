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

// TeamDeleteRequest describes a guarded Team delete. Linear's teamDelete archives
// the Team and schedules its data for deletion, and linctl has no restore path,
// so the command is irreversible. A Team is what a pin names: OrgWide must be
// true, the Team must belong to the Resolved Target's organization, and the
// pinned Team itself is refused because deleting it would kill the pin.
type TeamDeleteRequest struct {
	ID      string
	OrgWide bool
}

// DeleteTeam archives a Team and schedules its deletion after confirming
// --org-wide was passed, resolving the Team, comparing its organization against
// the Resolved Target's organization, and refusing the pinned Team.
func DeleteTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request TeamDeleteRequest,
) (string, error) {
	if request.ID == "" {
		return "", requiredFieldError("team id")
	}
	if !request.OrgWide {
		return "", fmt.Errorf(
			"%w: a Team is organization-owned and is what a pin names; pass --org-wide to delete one",
			ErrTargetMismatch,
		)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}

	return guard.deleteTeam(ctx, request.ID)
}

func (guard *guardedClient) deleteTeam(ctx context.Context, teamID string) (string, error) {
	team, err := GetTeamByID(ctx, guard.graphqlClient, teamID)
	if err != nil {
		return "", err
	}
	if err := guard.requireOrganization(team.OrgID); err != nil {
		return "", err
	}
	if team.ID == guard.target.Team.ID {
		return "", fmt.Errorf(
			"%w: team %s is the pinned team %s and cannot be deleted through its own pin",
			ErrTargetMismatch, team.ID, guard.target.Team.Key,
		)
	}

	deleted, err := gql.TeamDelete(ctx, guard.graphqlClient, teamID)
	if err != nil {
		return "", fmt.Errorf("delete team %s: %w", teamID, err)
	}
	if err := mutationSuccess(deleted.TeamDelete.Success, "teamDelete"); err != nil {
		return "", err
	}

	return deleted.TeamDelete.EntityId, nil
}
