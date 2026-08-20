package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Page
}

//nolint:lll
type teamMembershipsNode = gql.XTeamMembershipsTeamMembershipsTeamMembershipConnectionNodesTeamMembership

type teamMembershipsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListTeamMemberships returns team memberships visible to the authenticated user.
func ListTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (TeamMembershipList, error) {
	query := teamMembershipsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list team memberships", limit, defaultListPageSize,
		query.page,
		teamMembershipsNodeSummary,
	)
	if err != nil {
		return TeamMembershipList{}, err
	}

	return TeamMembershipList{Memberships: page.Items, Page: page.Page}, nil
}

// GetTeamMembershipByID returns one team membership by Linear id.
func GetTeamMembershipByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (TeamMembershipSummary, error) {
	result, err := gql.XTeamMembership(ctx, graphqlClient, id)
	if err != nil {
		return TeamMembershipSummary{}, fmt.Errorf("get team membership %s: %w", id, err)
	}

	return teamMembershipSummary(result.TeamMembership.TeamMembershipSummaryFields), nil
}

func (query teamMembershipsQuery) page(pageSize int, after *string) ([]teamMembershipsNode, bool, *string, error) {
	result, err := gql.XTeamMemberships(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.TeamMemberships.Nodes,
		result.TeamMemberships.PageInfo.HasNextPage,
		result.TeamMemberships.PageInfo.EndCursor,
		nil
}

func teamMembershipsNodeSummary(membership teamMembershipsNode) TeamMembershipSummary {
	return teamMembershipSummary(membership.TeamMembershipSummaryFields)
}

func teamMembershipSummary(membership gql.TeamMembershipSummaryFields) TeamMembershipSummary {
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
