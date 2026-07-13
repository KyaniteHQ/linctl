package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

const targetResolutionPageSize = 250

// ErrTargetMismatch marks a resolved target that does not match the pinned target.
var ErrTargetMismatch = errors.New("target mismatch")

// ErrTargetNotConfigured marks a missing or incomplete pinned target (no
// org_id/team_key/team_id). It is distinct from ErrTargetMismatch, which is a
// auth credential that resolves to a target other than the one pinned.
var ErrTargetNotConfigured = errors.New("target not configured")

// ResolvedTarget is the auth-resolved Linear write target.
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

// OrganizationExistsStatus reports whether a Linear organization URL key exists.
type OrganizationExistsStatus struct {
	URLKey  string `json:"url_key"`
	Success bool   `json:"success"`
	Exists  bool   `json:"exists"`
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

// CheckOrganizationExists checks whether a Linear organization URL key exists.
func CheckOrganizationExists(
	ctx context.Context,
	graphqlClient graphql.Client,
	urlKey string,
) (OrganizationExistsStatus, error) {
	result, err := organizationExists(ctx, graphqlClient, urlKey)
	if err != nil {
		return OrganizationExistsStatus{}, err
	}

	return OrganizationExistsStatus{
		URLKey:  urlKey,
		Success: result.OrganizationExists.Success,
		Exists:  result.OrganizationExists.Exists,
	}, nil
}

// ResolveTarget resolves viewer, organization, team, and optional project from auth.
func ResolveTarget(ctx context.Context, graphqlClient graphql.Client, expected config.Target) (ResolvedTarget, error) {
	if err := requireExpectedTarget(expected); err != nil {
		return ResolvedTarget{}, err
	}

	viewer, err := Viewer(ctx, graphqlClient)
	if err != nil {
		return ResolvedTarget{}, fmt.Errorf("resolve viewer: %w", err)
	}
	resolvedTeam, ok, err := resolveTeam(ctx, graphqlClient, expected)
	if err != nil {
		return ResolvedTarget{}, err
	}
	if !ok {
		return ResolvedTarget{}, fmt.Errorf(
			"%w: the active credential cannot reach %s",
			ErrTargetMismatch,
			describeExpectedTeam(expected),
		)
	}

	project, hasProject, err := resolveProject(ctx, graphqlClient, expected, resolvedTeam)
	if err != nil {
		return ResolvedTarget{}, err
	}

	resolved := resolvedTargetConfig(viewer.Viewer.Organization.Id, resolvedTeam, project, hasProject, expected)
	if err := requireTargetMatch(expected, resolved); err != nil {
		return ResolvedTarget{}, err
	}

	return newResolvedTarget(viewer.Viewer, resolvedTeam, project, hasProject, expected, resolved), nil
}

// requireExpectedTarget rejects a target that names no team at all. team_id is
// optional: a team key is unique inside an organization, so org_id plus team_key
// already identifies the team, and the id is resolved from the active credential.
// Leaving it optional is what lets `--team KEY` retarget a write on its own and
// still fail closed when the credential cannot reach KEY.
func requireExpectedTarget(expected config.Target) error {
	if expected.OrgID == "" || expected.TeamKey == "" {
		return fmt.Errorf(
			"%w: set org_id and team_key in .linctl.toml",
			ErrTargetNotConfigured,
		)
	}

	return nil
}

func resolveTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
) (TeamsTeamsTeamConnectionNodesTeam, bool, error) {
	// Fast path: the pinned target already names the exact team id, so try a
	// direct lookup before paging the org's full team list. Any error or
	// mismatch falls through to the scan below, which remains the semantic
	// authority for failures (including ErrTargetMismatch).
	if resolvedTeam, ok := resolveTeamDirect(ctx, graphqlClient, expected); ok {
		return resolvedTeam, true, nil
	}

	var after *string
	for {
		teams, err := Teams(ctx, graphqlClient, intPtr(targetResolutionPageSize), after, boolPtr(true))
		if err != nil {
			return TeamsTeamsTeamConnectionNodesTeam{}, false, fmt.Errorf("resolve teams: %w", err)
		}
		resolvedTeam, ok := findResolvedTeam(teams.Teams.Nodes, expected)
		if ok {
			return resolvedTeam, true, nil
		}
		next, ok, err := nextTargetPageCursor(
			teams.Teams.PageInfo.HasNextPage,
			teams.Teams.PageInfo.EndCursor,
			"teams",
		)
		if err != nil || !ok {
			return TeamsTeamsTeamConnectionNodesTeam{}, false, err
		}
		after = next
	}
}

func resolvedTargetConfig(
	orgID string,
	team TeamsTeamsTeamConnectionNodesTeam,
	project ResolvedProject,
	hasProject bool,
	expected config.Target,
) config.Target {
	resolved := config.Target{
		OrgID:     orgID,
		TeamKey:   team.Key,
		TeamID:    team.Id,
		ProjectID: expected.ProjectID,
	}
	if hasProject {
		resolved.ProjectID = project.ID
	}

	return resolved
}

// describeExpectedTeam names the pinned team for an error message, omitting the
// team id when the target identifies the team by key alone.
func describeExpectedTeam(expected config.Target) string {
	if expected.TeamID == "" {
		return "team_key=" + expected.TeamKey
	}

	return fmt.Sprintf("team_key=%s team_id=%s", expected.TeamKey, expected.TeamID)
}

func requireTargetMatch(expected config.Target, resolved config.Target) error {
	orgMatches := resolved.OrgID == expected.OrgID
	keyMatches := resolved.TeamKey == expected.TeamKey
	// An empty pinned team id means the target named the team by key alone, so the
	// id resolved from the credential is the answer rather than something to check.
	idMatches := expected.TeamID == "" || resolved.TeamID == expected.TeamID

	if orgMatches && keyMatches && idMatches {
		return nil
	}

	return fmt.Errorf("%w: expected=%+v resolved=%+v", ErrTargetMismatch, expected, resolved)
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

func resolveTeamDirect(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
) (TeamsTeamsTeamConnectionNodesTeam, bool) {
	// Without a pinned team id there is nothing to look up directly; the scan
	// below resolves the team from its key instead.
	if expected.TeamID == "" {
		return TeamsTeamsTeamConnectionNodesTeam{}, false
	}

	result, err := team(ctx, graphqlClient, expected.TeamID)
	if err != nil {
		return TeamsTeamsTeamConnectionNodesTeam{}, false
	}

	found := result.Team
	if found.Id != expected.TeamID || found.Key != expected.TeamKey || found.Organization.Id != expected.OrgID {
		return TeamsTeamsTeamConnectionNodesTeam{}, false
	}

	return TeamsTeamsTeamConnectionNodesTeam{
		Id:   found.Id,
		Key:  found.Key,
		Name: found.Name,
		Organization: TeamsTeamsTeamConnectionNodesTeamOrganization{
			Id:     found.Organization.Id,
			Name:   found.Organization.Name,
			UrlKey: found.Organization.UrlKey,
		},
	}, true
}

func findResolvedTeam(
	teams []TeamsTeamsTeamConnectionNodesTeam,
	expected config.Target,
) (TeamsTeamsTeamConnectionNodesTeam, bool) {
	for _, team := range teams {
		if team.Key != expected.TeamKey || team.Organization.Id != expected.OrgID {
			continue
		}
		// A pinned team id is still compared exactly. It is only skipped when the
		// target names the team by key alone.
		if expected.TeamID != "" && team.Id != expected.TeamID {
			continue
		}

		return team, true
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

	project, err := TargetProject(
		ctx,
		graphqlClient,
		expected.ProjectID,
		intPtr(targetResolutionPageSize),
		nil,
	)
	if err != nil {
		return ResolvedProject{}, false, fmt.Errorf("resolve project: %w", err)
	}
	projectID := project.Project.Id
	projectName := project.Project.Name
	projectTeamPages := []TargetProjectProjectTeamsTeamConnection{project.Project.Teams}
	after := project.Project.Teams.PageInfo.EndCursor
	for {
		for _, projectTeam := range projectTeamPageNodes(projectTeamPages) {
			if projectTeam.Id == team.Id {
				return ResolvedProject{
					ID:   projectID,
					Name: projectName,
				}, true, nil
			}
		}
		next, ok, err := nextTargetPageCursor(
			projectTeamPages[len(projectTeamPages)-1].PageInfo.HasNextPage,
			after,
			"project teams",
		)
		if err != nil {
			return ResolvedProject{}, false, err
		}
		if !ok {
			break
		}
		project, err := TargetProject(
			ctx,
			graphqlClient,
			expected.ProjectID,
			intPtr(targetResolutionPageSize),
			next,
		)
		if err != nil {
			return ResolvedProject{}, false, fmt.Errorf("resolve project: %w", err)
		}
		projectTeamPages = []TargetProjectProjectTeamsTeamConnection{project.Project.Teams}
		after = project.Project.Teams.PageInfo.EndCursor
	}

	return ResolvedProject{}, false, fmt.Errorf(
		"%w: project_id=%s not attached to team_id=%s",
		ErrTargetMismatch,
		expected.ProjectID,
		team.Id,
	)
}

func projectTeamPageNodes(
	pages []TargetProjectProjectTeamsTeamConnection,
) []TargetProjectProjectTeamsTeamConnectionNodesTeam {
	teams := []TargetProjectProjectTeamsTeamConnectionNodesTeam{}
	for _, page := range pages {
		teams = append(teams, page.Nodes...)
	}

	return teams
}

func nextTargetPageCursor(hasNextPage bool, endCursor *string, collection string) (*string, bool, error) {
	if !hasNextPage {
		return nil, false, nil
	}
	if endCursor == nil || *endCursor == "" {
		return nil, false, fmt.Errorf(
			"%w: %s pageInfo.hasNextPage without endCursor",
			ErrTargetMismatch,
			collection,
		)
	}

	return stringPtr(*endCursor), true, nil
}

func optionalProject(project ResolvedProject, ok bool) *ResolvedProject {
	if !ok {
		return nil
	}

	return &project
}
