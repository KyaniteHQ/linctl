package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// TeamSummary is the compact Team model used by team commands.
type TeamSummary struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ArchivedAt  string `json:"archived_at,omitempty"`
	OrgID       string `json:"org_id"`
	OrgName     string `json:"org_name"`
	OrgURLKey   string `json:"org_url_key"`
}

// TeamList is a page of teams.
type TeamList struct {
	Teams       []TeamSummary `json:"teams"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// TeamMemberList is a page of team members.
type TeamMemberList struct {
	TeamID      string        `json:"team_id"`
	TeamKey     string        `json:"team_key"`
	TeamName    string        `json:"team_name"`
	Members     []UserSummary `json:"members"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// ListTeams returns visible teams.
func ListTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	teams, err := Teams(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamList{}, fmt.Errorf("list teams: %w", err)
	}

	summaries := make([]TeamSummary, 0, len(teams.Teams.Nodes))
	for _, team := range teams.Teams.Nodes {
		summaries = append(summaries, teamSummaryFromConnection(team))
	}

	return TeamList{
		Teams:       summaries,
		HasNextPage: teams.Teams.PageInfo.HasNextPage,
		EndCursor:   teams.Teams.PageInfo.EndCursor,
	}, nil
}

// GetTeamByID returns one Team by id.
func GetTeamByID(ctx context.Context, graphqlClient graphql.Client, id string) (TeamSummary, error) {
	team, err := TeamByID(ctx, graphqlClient, id)
	if err != nil {
		return TeamSummary{}, fmt.Errorf("get team %s: %w", id, err)
	}

	return teamSummary(team.Team.TeamSummaryFields), nil
}

// ListTeamMembers returns visible members for one Team.
func ListTeamMembers(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TeamMemberList, error) {
	team, err := TeamMembers(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamMemberList{}, fmt.Errorf("list team members %s: %w", id, err)
	}

	members := make([]UserSummary, 0, len(team.Team.Members.Nodes))
	for _, member := range team.Team.Members.Nodes {
		members = append(members, userSummary(member.UserSummaryFields))
	}

	return TeamMemberList{
		TeamID:      team.Team.Id,
		TeamKey:     team.Team.Key,
		TeamName:    team.Team.Name,
		Members:     members,
		HasNextPage: team.Team.Members.PageInfo.HasNextPage,
		EndCursor:   team.Team.Members.PageInfo.EndCursor,
	}, nil
}

func teamSummary(team TeamSummaryFields) TeamSummary {
	description := ""
	if team.Description != nil {
		description = *team.Description
	}
	archivedAt := ""
	if team.ArchivedAt != nil {
		archivedAt = *team.ArchivedAt
	}

	return TeamSummary{
		ID:          team.Id,
		Key:         team.Key,
		Name:        team.Name,
		Description: description,
		ArchivedAt:  archivedAt,
		OrgID:       team.Organization.Id,
		OrgName:     team.Organization.Name,
		OrgURLKey:   team.Organization.UrlKey,
	}
}

func teamSummaryFromConnection(team TeamsTeamsTeamConnectionNodesTeam) TeamSummary {
	return TeamSummary{
		ID:        team.Id,
		Key:       team.Key,
		Name:      team.Name,
		OrgID:     team.Organization.Id,
		OrgName:   team.Organization.Name,
		OrgURLKey: team.Organization.UrlKey,
	}
}
