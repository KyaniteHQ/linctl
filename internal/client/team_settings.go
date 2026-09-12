package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// TeamSettingsRequest describes a guarded Team settings write. It reaches only
// the settings a rebuild needs: triage on or off, the default state for new
// issues, and whether a sub-team inherits its parent's workflow states. Name,
// key, and description stay out of reach.
type TeamSettingsRequest struct {
	ID             string
	Triage         *bool
	DefaultStateID string
	Inherit        *bool
	OrgWide        bool
}

// UpdateTeamSettings writes team settings after confirming --org-wide was
// passed, resolving the Team, comparing its organization against the Resolved
// Target's organization, and checking that a default state belongs to the Team
// or to the parent it inherits from.
func UpdateTeamSettings(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request TeamSettingsRequest,
) (TeamSummary, error) {
	if request.ID == "" {
		return TeamSummary{}, requiredFieldError("team id")
	}
	if request.Triage == nil && request.DefaultStateID == "" && request.Inherit == nil {
		return TeamSummary{}, requiredFieldError("triage, default state, or inherit-workflow-states")
	}
	if !request.OrgWide {
		return TeamSummary{}, fmt.Errorf(
			"%w: a Team is organization-owned and is what a pin names; pass --org-wide to change its settings",
			ErrTargetMismatch,
		)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return TeamSummary{}, err
	}

	return guard.updateTeamSettings(ctx, request)
}

func (guard *guardedClient) updateTeamSettings(
	ctx context.Context,
	request TeamSettingsRequest,
) (TeamSummary, error) {
	team, err := GetTeamByID(ctx, guard.graphqlClient, request.ID)
	if err != nil {
		return TeamSummary{}, err
	}
	if err := guard.requireOrganization(team.OrgID); err != nil {
		return TeamSummary{}, err
	}
	if request.DefaultStateID != "" {
		if err := guard.requireDefaultState(ctx, team, request.DefaultStateID); err != nil {
			return TeamSummary{}, err
		}
	}

	updated, err := gql.TeamSettingsUpdate(ctx, guard.graphqlClient, request.ID, LinearTeamSettingsInput{
		TriageEnabled:           request.Triage,
		DefaultIssueStateID:     optionalString(request.DefaultStateID),
		InheritWorkflowStatuses: request.Inherit,
	})
	if err != nil {
		return TeamSummary{}, fmt.Errorf("update team settings %s: %w", request.ID, err)
	}
	if !updated.TeamUpdate.Success || updated.TeamUpdate.Team == nil {
		return TeamSummary{}, fmt.Errorf("%w: teamUpdate failed", ErrMutationFailed)
	}

	summary := teamSummary(updated.TeamUpdate.Team.TeamSummaryFields)
	if err := guard.requireOrganization(summary.OrgID); err != nil {
		return TeamSummary{}, err
	}

	return summary, nil
}

// requireDefaultState resolves the state named as a team default and refuses
// one that belongs to another team. A sub-team that inherits states owns none
// of them, so its parent's states are accepted too.
func (guard *guardedClient) requireDefaultState(
	ctx context.Context,
	team TeamSummary,
	stateID string,
) error {
	state, err := GetWorkflowStateByID(ctx, guard.graphqlClient, stateID)
	if err != nil {
		return err
	}
	if state.TeamID == team.ID || (team.ParentID != "" && state.TeamID == team.ParentID) {
		return nil
	}

	return fmt.Errorf(
		"%w: workflow state %s belongs to team %s, not to %s or its parent",
		ErrTargetMismatch, stateID, state.TeamKey, team.Key,
	)
}
