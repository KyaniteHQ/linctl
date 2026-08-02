package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

const projectTeamsPageSize = 50

// ProjectAddTeamRequest describes a resource-scoped write that attaches one more
// team to an existing project without dropping the teams already on it.
type ProjectAddTeamRequest struct {
	ProjectID string
	TeamID    string
	TeamKey   string
}

// IssueMoveTeamRequest describes a resource-scoped write that moves an issue from
// the pinned team to another team in the same organization. The issue must start
// on the pin; after the move it deliberately no longer matches the pin.
type IssueMoveTeamRequest struct {
	IssueID string
	TeamID  string
	TeamKey string
}

// AddProjectTeam attaches a team to a project after comparing the project against
// the pinned target. Linear's projectUpdate teamIds field replaces the full set,
// so the write merges the destination onto the complete current membership and
// refuses when that membership cannot be fully enumerated. The destination team
// must belong to the same organization as the pin. Adding a team already present
// is a no-op that returns the current project.
func AddProjectTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectAddTeamRequest,
) (ProjectSummary, error) {
	if request.ProjectID == "" {
		return ProjectSummary{}, requiredFieldError("project id")
	}
	if request.TeamID == "" && request.TeamKey == "" {
		return ProjectSummary{}, requiredFieldError("--to-team or --to-team-id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return ProjectSummary{}, err
	}

	return guard.addProjectTeam(ctx, request)
}

func (guard *guardedClient) addProjectTeam(
	ctx context.Context,
	request ProjectAddTeamRequest,
) (ProjectSummary, error) {
	project, err := GetProjectByID(ctx, guard.graphqlClient, request.ProjectID)
	if err != nil {
		return ProjectSummary{}, err
	}
	if guard.target.Project != nil && project.ID != guard.target.Project.ID {
		return ProjectSummary{}, guard.projectMismatchError("project_id", project.ID)
	}
	if err := guard.requireProjectTeam(project); err != nil {
		return ProjectSummary{}, err
	}

	destination, err := guard.resolveDestinationTeam(ctx, request.TeamID, request.TeamKey)
	if err != nil {
		return ProjectSummary{}, err
	}
	if projectHasTeam(project, destination.ID, destination.Key) {
		return project, nil
	}

	teamIDs, err := guard.listAllProjectTeamIDs(ctx, request.ProjectID, project)
	if err != nil {
		return ProjectSummary{}, err
	}
	teamIDs = append(teamIDs, destination.ID)

	updated, err := gql.ProjectUpdate(ctx, guard.graphqlClient, request.ProjectID, LinearProjectUpdateInput{
		TeamIDs: teamIDs,
	})
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("add team to project %s: %w", request.ProjectID, err)
	}
	if !updated.ProjectUpdate.Success || updated.ProjectUpdate.Project == nil {
		return ProjectSummary{}, fmt.Errorf("%w: projectUpdate returned no project", ErrMutationFailed)
	}

	summary := projectSummaryFromFields(updated.ProjectUpdate.Project.ProjectSummaryFields)
	if !projectHasTeam(summary, destination.ID, destination.Key) {
		return ProjectSummary{}, fmt.Errorf(
			"%w: project %s is missing destination team_id=%s team_key=%s after update",
			ErrMutationFailed,
			summary.ID,
			destination.ID,
			destination.Key,
		)
	}
	if err := guard.requireProjectTeam(summary); err != nil {
		return ProjectSummary{}, err
	}

	return summary, nil
}

// MoveIssueTeam moves an issue to another team in the same organization. The
// issue must currently match the pinned team (and project, when one is pinned).
// The destination is compared against the organization only: after the move the
// issue leaves the pin by design, so the returned issue is not re-checked for
// team match.
func MoveIssueTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueMoveTeamRequest,
) (IssueSummary, error) {
	if request.IssueID == "" {
		return IssueSummary{}, requiredFieldError("issue id")
	}
	if request.TeamID == "" && request.TeamKey == "" {
		return IssueSummary{}, requiredFieldError("--to-team or --to-team-id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.moveIssueTeam(ctx, request)
}

func (guard *guardedClient) moveIssueTeam(
	ctx context.Context,
	request IssueMoveTeamRequest,
) (IssueSummary, error) {
	issue, err := guard.requireIssueDetail(ctx, request.IssueID)
	if err != nil {
		return IssueSummary{}, err
	}

	destination, err := guard.resolveDestinationTeam(ctx, request.TeamID, request.TeamKey)
	if err != nil {
		return IssueSummary{}, err
	}
	if destination.ID == issue.Summary.TeamID || destination.Key == issue.Summary.Team {
		return IssueSummary{}, fmt.Errorf(
			"%w: issue %s is already on team_id=%s team_key=%s",
			ErrWriteInvalid,
			request.IssueID,
			issue.Summary.TeamID,
			issue.Summary.Team,
		)
	}

	updated, err := gql.IssueUpdate(ctx, guard.graphqlClient, request.IssueID, LinearIssueUpdateInput{
		TeamID: stringPtr(destination.ID),
	})
	if err != nil {
		return IssueSummary{}, fmt.Errorf("move issue %s: %w", request.IssueID, err)
	}
	if !updated.IssueUpdate.Success || updated.IssueUpdate.Issue == nil {
		return IssueSummary{}, fmt.Errorf("%w: issueUpdate returned no issue", ErrMutationFailed)
	}

	summary := issueSummaryFromFields(updated.IssueUpdate.Issue.IssueSummaryFields)
	if summary.TeamID != destination.ID || summary.Team != destination.Key {
		return IssueSummary{}, fmt.Errorf(
			"%w: issue %s landed on team_id=%s team_key=%s, expected team_id=%s team_key=%s",
			ErrMutationFailed,
			request.IssueID,
			summary.TeamID,
			summary.Team,
			destination.ID,
			destination.Key,
		)
	}

	return summary, nil
}

// resolveDestinationTeam resolves a team inside the resolved organization by id
// and/or key. When both are set they must name the same team. Organization
// membership is the hard stop: a team in another organization is a Target Mismatch.
func (guard *guardedClient) resolveDestinationTeam(
	ctx context.Context,
	teamID string,
	teamKey string,
) (TeamSummary, error) {
	if teamID != "" {
		team, err := GetTeamByID(ctx, guard.graphqlClient, teamID)
		if err != nil {
			return TeamSummary{}, err
		}
		if err := guard.requireOrganization(team.OrgID); err != nil {
			return TeamSummary{}, err
		}
		if teamKey != "" && team.Key != teamKey {
			return TeamSummary{}, fmt.Errorf(
				"%w: team_id=%s has team_key=%s, not %s",
				ErrWriteInvalid,
				teamID,
				team.Key,
				teamKey,
			)
		}

		return team, nil
	}

	team, err := guard.findTeamByKey(ctx, teamKey)
	if err != nil {
		return TeamSummary{}, err
	}
	if err := guard.requireOrganization(team.OrgID); err != nil {
		return TeamSummary{}, err
	}

	return team, nil
}

func (guard *guardedClient) findTeamByKey(ctx context.Context, teamKey string) (TeamSummary, error) {
	var after *string
	for {
		teams, err := gql.XTeams_list(ctx, guard.graphqlClient, intPtr(targetResolutionPageSize), after, boolPtr(true))
		if err != nil {
			return TeamSummary{}, fmt.Errorf("list teams: %w", err)
		}
		for _, node := range teams.Teams.Nodes {
			summary := teamSummary(node.TeamSummaryFields)
			if summary.Key == teamKey {
				return summary, nil
			}
		}
		if !teams.Teams.PageInfo.HasNextPage {
			return TeamSummary{}, fmt.Errorf("%w: no visible team with key %s", ErrWriteInvalid, teamKey)
		}
		if teams.Teams.PageInfo.EndCursor == nil {
			return TeamSummary{}, errors.New("list teams: next page has no end cursor")
		}
		after = teams.Teams.PageInfo.EndCursor
	}
}

// listAllProjectTeamIDs returns every team id currently on the project. Linear's
// projectUpdate teamIds field replaces the association set, so a truncated first
// page must be fully walked before the write; otherwise an unlisted team would
// be dropped.
func (guard *guardedClient) listAllProjectTeamIDs(
	ctx context.Context,
	projectID string,
	project ProjectSummary,
) ([]string, error) {
	if !project.TeamsTruncated {
		ids := make([]string, 0, len(project.Teams))
		for _, team := range project.Teams {
			ids = append(ids, team.ID)
		}

		return ids, nil
	}

	ids := make([]string, 0, projectTeamsPageSize)
	var after *string
	for {
		page, err := gql.XProject_teams(
			ctx, guard.graphqlClient, projectID, intPtr(projectTeamsPageSize), after, boolPtr(true),
		)
		if err != nil {
			return nil, fmt.Errorf("list project teams %s: %w", projectID, err)
		}
		for _, node := range page.Project.Teams.Nodes {
			ids = append(ids, node.Id)
		}
		if !page.Project.Teams.PageInfo.HasNextPage {
			if len(ids) == 0 {
				return nil, fmt.Errorf(
					"%w: project %s has no teams to merge before add",
					ErrMutationFailed,
					projectID,
				)
			}

			return ids, nil
		}
		if page.Project.Teams.PageInfo.EndCursor == nil {
			return nil, fmt.Errorf(
				"%w: project %s teams page is truncated without an end cursor",
				ErrTargetMismatch,
				projectID,
			)
		}
		after = page.Project.Teams.PageInfo.EndCursor
	}
}
