package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// ErrTargetMismatch marks a resolved target that does not match the pinned target.
var ErrTargetMismatch = errors.New("target mismatch")

// ResolvedTarget is the token-resolved Linear write target.
type ResolvedTarget struct {
	Viewer    TargetViewer     `json:"viewer"`
	Org       TargetOrg        `json:"org"`
	Team      TargetTeam       `json:"team"`
	Project   *ResolvedProject `json:"project,omitempty"`
	Expected  config.Target    `json:"expected"`
	Resolved  config.Target    `json:"resolved"`
	Confirmed bool             `json:"confirmed"`
}

// TargetViewer is the authenticated Linear user.
type TargetViewer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// TargetOrg is the resolved Linear organization.
type TargetOrg struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URLKey string `json:"url_key"`
}

// TargetTeam is the resolved Linear team.
type TargetTeam struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ResolvedProject is the resolved Linear project.
type ResolvedProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ResolveTarget resolves viewer, organization, team, and optional project from the token.
func ResolveTarget(ctx context.Context, graphqlClient graphql.Client, expected config.Target) (ResolvedTarget, error) {
	if expected.OrgID == "" || expected.TeamID == "" || expected.TeamKey == "" {
		return ResolvedTarget{}, fmt.Errorf(
			"%w: expected org_id, team_id, and team_key are required",
			ErrTargetMismatch,
		)
	}

	viewer, err := Viewer(ctx, graphqlClient)
	if err != nil {
		return ResolvedTarget{}, fmt.Errorf("resolve viewer: %w", err)
	}
	teams, err := Teams(ctx, graphqlClient, intPtr(250), nil, boolPtr(true))
	if err != nil {
		return ResolvedTarget{}, fmt.Errorf("resolve teams: %w", err)
	}

	resolvedTeam, ok := findResolvedTeam(teams.Teams.Nodes, expected)
	if !ok {
		return ResolvedTarget{}, fmt.Errorf(
			"%w: expected team_id=%s team_key=%s",
			ErrTargetMismatch,
			expected.TeamID,
			expected.TeamKey,
		)
	}

	project, hasProject, err := resolveProject(ctx, graphqlClient, expected, resolvedTeam)
	if err != nil {
		return ResolvedTarget{}, err
	}

	resolved := config.Target{
		OrgID:     viewer.Viewer.Organization.Id,
		TeamKey:   resolvedTeam.Key,
		TeamID:    resolvedTeam.Id,
		ProjectID: expected.ProjectID,
	}
	if hasProject {
		resolved.ProjectID = project.ID
	}
	if resolved.OrgID != expected.OrgID || resolved.TeamID != expected.TeamID || resolved.TeamKey != expected.TeamKey {
		return ResolvedTarget{}, fmt.Errorf("%w: expected=%+v resolved=%+v", ErrTargetMismatch, expected, resolved)
	}

	return newResolvedTarget(viewer.Viewer, resolvedTeam, project, hasProject, expected, resolved), nil
}

func newResolvedTarget(
	viewer ViewerViewerUser,
	team TeamsTeamsTeamConnectionNodesTeam,
	project ResolvedProject,
	hasProject bool,
	expected config.Target,
	resolved config.Target,
) ResolvedTarget {
	return ResolvedTarget{
		Viewer: TargetViewer{
			ID:          viewer.Id,
			Name:        viewer.Name,
			DisplayName: viewer.DisplayName,
			Email:       viewer.Email,
		},
		Org: TargetOrg{
			ID:     viewer.Organization.Id,
			Name:   viewer.Organization.Name,
			URLKey: viewer.Organization.UrlKey,
		},
		Team: TargetTeam{
			ID:   team.Id,
			Key:  team.Key,
			Name: team.Name,
		},
		Project:   optionalProject(project, hasProject),
		Expected:  expected,
		Resolved:  resolved,
		Confirmed: true,
	}
}

func findResolvedTeam(
	teams []TeamsTeamsTeamConnectionNodesTeam,
	expected config.Target,
) (TeamsTeamsTeamConnectionNodesTeam, bool) {
	for _, team := range teams {
		if team.Id == expected.TeamID && team.Key == expected.TeamKey && team.Organization.Id == expected.OrgID {
			return team, true
		}
	}

	return TeamsTeamsTeamConnectionNodesTeam{}, false
}

func resolveProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	team TeamsTeamsTeamConnectionNodesTeam,
) (ResolvedProject, bool, error) {
	if expected.ProjectID == "" {
		return ResolvedProject{}, false, nil
	}

	project, err := TargetProject(ctx, graphqlClient, expected.ProjectID)
	if err != nil {
		return ResolvedProject{}, false, fmt.Errorf("resolve project: %w", err)
	}
	for _, projectTeam := range project.Project.Teams.Nodes {
		if projectTeam.Id == team.Id {
			return ResolvedProject{
				ID:   project.Project.Id,
				Name: project.Project.Name,
			}, true, nil
		}
	}

	return ResolvedProject{}, false, fmt.Errorf(
		"%w: project_id=%s not attached to team_id=%s",
		ErrTargetMismatch,
		expected.ProjectID,
		team.Id,
	)
}

func optionalProject(project ResolvedProject, ok bool) *ResolvedProject {
	if !ok {
		return nil
	}

	return &project
}

func intPtr(value int) *int {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}
