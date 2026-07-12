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

func guardedMutation[T any](
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	mutate func(writeGuard) (T, error),
) (T, error) {
	var zero T
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return zero, err
	}

	return mutate(guard)
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
		return IssueDetail{}, guard.teamMismatchError("issue", issue.Summary.TeamID, issue.Summary.Team)
	}
	if guard.target.Project != nil && issue.Summary.ProjectID != guard.target.Project.ID {
		return IssueDetail{}, guard.projectMismatchError("issue project_id", issue.Summary.ProjectID)
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
		return guard.projectMismatchError("project_id", project.ID)
	}
	if !projectHasTeam(project, guard.target.Team.ID, guard.target.Team.Key) {
		return guard.teamNotAttachedError()
	}

	return nil
}

func (guard writeGuard) requireProjectMilestone(
	ctx context.Context,
	graphqlClient graphql.Client,
	projectMilestoneID string,
) error {
	milestone, err := GetProjectMilestoneDetail(ctx, graphqlClient, projectMilestoneID)
	if err != nil {
		return err
	}
	if guard.target.Project != nil && milestone.Project.ID != guard.target.Project.ID {
		return guard.projectMismatchError("project_id", milestone.Project.ID)
	}
	if !projectHasTeam(milestone.Project, guard.target.Team.ID, guard.target.Team.Key) {
		return guard.teamNotAttachedError()
	}

	return nil
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
		return guard.teamMismatchError("cycle", cycle.TeamID, cycle.TeamKey)
	}

	return nil
}

// teamMismatchError reports a Target Mismatch between the pinned team and the
// team resolved from an existing entity. It is a hard stop for guarded writes.
func (guard writeGuard) teamMismatchError(entity string, teamID string, teamKey string) error {
	return fmt.Errorf(
		"%w: expected team_id=%s team_key=%s resolved %s team_id=%s team_key=%s",
		ErrTargetMismatch,
		guard.target.Team.ID,
		guard.target.Team.Key,
		entity,
		teamID,
		teamKey,
	)
}

// projectMismatchError reports a Target Mismatch between the pinned project
// and the project resolved from an existing entity.
func (guard writeGuard) projectMismatchError(label string, projectID string) error {
	return fmt.Errorf(
		"%w: expected project_id=%s resolved %s=%s",
		ErrTargetMismatch,
		guard.target.Project.ID,
		label,
		projectID,
	)
}

// teamNotAttachedError reports a Target Mismatch when a resolved project is
// not attached to the pinned team.
func (guard writeGuard) teamNotAttachedError() error {
	return fmt.Errorf(
		"%w: expected team_id=%s team_key=%s",
		ErrTargetMismatch,
		guard.target.Team.ID,
		guard.target.Team.Key,
	)
}

func projectHasTeam(project ProjectSummary, teamID string, teamKey string) bool {
	for _, team := range project.Teams {
		if team.ID == teamID && team.Key == teamKey {
			return true
		}
	}

	return false
}
