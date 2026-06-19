package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

type writeGuard struct {
	target ResolvedTarget
}

func newWriteGuard(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
) (writeGuard, error) {
	target, err := ResolveTarget(ctx, graphqlClient, expected)
	if err != nil {
		return writeGuard{}, err
	}

	return writeGuard{target: target}, nil
}

func (guard writeGuard) requireIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	issueID string,
) (IssueSummary, error) {
	issue, err := GetIssueByID(ctx, graphqlClient, issueID)
	if err != nil {
		return IssueSummary{}, err
	}
	if issue.TeamID != guard.target.Team.ID || issue.Team != guard.target.Team.Key {
		return IssueSummary{}, fmt.Errorf(
			"%w: expected team_id=%s team_key=%s resolved issue team_id=%s team_key=%s",
			ErrTargetMismatch,
			guard.target.Team.ID,
			guard.target.Team.Key,
			issue.TeamID,
			issue.Team,
		)
	}
	if guard.target.Project != nil && issue.ProjectID != guard.target.Project.ID {
		return IssueSummary{}, fmt.Errorf(
			"%w: expected project_id=%s resolved issue project_id=%s",
			ErrTargetMismatch,
			guard.target.Project.ID,
			issue.ProjectID,
		)
	}

	return issue, nil
}

func (guard writeGuard) requireProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	projectID string,
) error {
	project, err := GetProjectByID(ctx, graphqlClient, projectID)
	if err != nil {
		return err
	}
	if guard.target.Project != nil && project.ID != guard.target.Project.ID {
		return fmt.Errorf(
			"%w: expected project_id=%s resolved project_id=%s",
			ErrTargetMismatch,
			guard.target.Project.ID,
			project.ID,
		)
	}
	if !projectHasTeam(project, guard.target.Team.ID, guard.target.Team.Key) {
		return fmt.Errorf(
			"%w: expected team_id=%s team_key=%s",
			ErrTargetMismatch,
			guard.target.Team.ID,
			guard.target.Team.Key,
		)
	}

	return nil
}

func projectHasTeam(project ProjectSummary, teamID string, teamKey string) bool {
	for _, team := range project.Teams {
		if team.ID == teamID && team.Key == teamKey {
			return true
		}
	}

	return false
}
