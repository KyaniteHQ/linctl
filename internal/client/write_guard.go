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
	issue, err := guard.requireIssueDetail(ctx, graphqlClient, issueID)
	if err != nil {
		return IssueSummary{}, err
	}

	return issue.Summary, nil
}

func (guard writeGuard) requireIssueDetail(
	ctx context.Context,
	graphqlClient graphql.Client,
	issueID string,
) (IssueDetail, error) {
	issue, err := GetIssueDetail(ctx, graphqlClient, issueID)
	if err != nil {
		return IssueDetail{}, err
	}
	if issue.Summary.TeamID != guard.target.Team.ID || issue.Summary.Team != guard.target.Team.Key {
		return IssueDetail{}, fmt.Errorf(
			"%w: expected team_id=%s team_key=%s resolved issue team_id=%s team_key=%s",
			ErrTargetMismatch,
			guard.target.Team.ID,
			guard.target.Team.Key,
			issue.Summary.TeamID,
			issue.Summary.Team,
		)
	}
	if guard.target.Project != nil && issue.Summary.ProjectID != guard.target.Project.ID {
		return IssueDetail{}, fmt.Errorf(
			"%w: expected project_id=%s resolved issue project_id=%s",
			ErrTargetMismatch,
			guard.target.Project.ID,
			issue.Summary.ProjectID,
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

func (guard writeGuard) requireProjectMilestone(
	ctx context.Context,
	graphqlClient graphql.Client,
	projectMilestoneID string,
) (ProjectMilestoneDetail, error) {
	milestone, err := GetProjectMilestoneDetail(ctx, graphqlClient, projectMilestoneID)
	if err != nil {
		return ProjectMilestoneDetail{}, err
	}
	if guard.target.Project != nil && milestone.Project.ID != guard.target.Project.ID {
		return ProjectMilestoneDetail{}, fmt.Errorf(
			"%w: expected project_id=%s resolved project_id=%s",
			ErrTargetMismatch,
			guard.target.Project.ID,
			milestone.Project.ID,
		)
	}
	if !projectHasTeam(milestone.Project, guard.target.Team.ID, guard.target.Team.Key) {
		return ProjectMilestoneDetail{}, fmt.Errorf(
			"%w: expected team_id=%s team_key=%s",
			ErrTargetMismatch,
			guard.target.Team.ID,
			guard.target.Team.Key,
		)
	}

	return milestone, nil
}

func (guard writeGuard) requireCycle(
	ctx context.Context,
	graphqlClient graphql.Client,
	cycleID string,
) error {
	cycle, err := GetCycleByID(ctx, graphqlClient, cycleID)
	if err != nil {
		return err
	}
	if cycle.TeamID != guard.target.Team.ID || cycle.TeamKey != guard.target.Team.Key {
		return fmt.Errorf(
			"%w: expected team_id=%s team_key=%s resolved cycle team_id=%s team_key=%s",
			ErrTargetMismatch,
			guard.target.Team.ID,
			guard.target.Team.Key,
			cycle.TeamID,
			cycle.TeamKey,
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
