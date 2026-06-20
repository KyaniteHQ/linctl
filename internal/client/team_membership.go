package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// TeamMembershipSummary is one user's membership in a Linear team.
type TeamMembershipSummary struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Email       string  `json:"email,omitempty"`
	Active      bool    `json:"active"`
	Guest       bool    `json:"guest"`
	Admin       bool    `json:"admin"`
	TeamID      string  `json:"team_id"`
	TeamKey     string  `json:"team_key"`
	TeamName    string  `json:"team_name"`
	Owner       bool    `json:"owner"`
	SortOrder   float64 `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	ArchivedAt  string  `json:"archived_at,omitempty"`
}

// TeamMembershipList is a page of team memberships.
type TeamMembershipList struct {
	Memberships []TeamMembershipSummary `json:"memberships"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// ListTeamMemberships returns visible team memberships in the workspace.
func ListTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (TeamMembershipList, error) {
	result, err := teamMemberships(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list team memberships: %w", err)
	}

	memberships := make([]TeamMembershipSummary, 0, len(result.TeamMemberships.Nodes))
	for _, membership := range result.TeamMemberships.Nodes {
		memberships = append(
			memberships,
			teamMembershipSummary(membership.TeamMembershipSummaryFields),
		)
	}

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: result.TeamMemberships.PageInfo.HasNextPage,
		EndCursor:   result.TeamMemberships.PageInfo.EndCursor,
	}, nil
}

// GetTeamMembershipByID returns one team membership by Linear id.
func GetTeamMembershipByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (TeamMembershipSummary, error) {
	result, err := teamMembership(ctx, graphqlClient, id)
	if err != nil {
		return TeamMembershipSummary{}, fmt.Errorf("get team membership %s: %w", id, err)
	}

	return teamMembershipSummary(result.TeamMembership.TeamMembershipSummaryFields), nil
}

func teamMembershipSummary(membership TeamMembershipSummaryFields) TeamMembershipSummary {
	return TeamMembershipSummary{
		ID:          membership.Id,
		UserID:      membership.User.Id,
		Name:        membership.User.Name,
		DisplayName: membership.User.DisplayName,
		Email:       membership.User.Email,
		Active:      membership.User.Active,
		Guest:       membership.User.Guest,
		Admin:       membership.User.Admin,
		TeamID:      membership.Team.Id,
		TeamKey:     membership.Team.Key,
		TeamName:    membership.Team.Name,
		Owner:       membership.Owner,
		SortOrder:   membership.SortOrder,
		CreatedAt:   membership.CreatedAt,
		UpdatedAt:   membership.UpdatedAt,
		ArchivedAt:  stringValue(membership.ArchivedAt),
	}
}
