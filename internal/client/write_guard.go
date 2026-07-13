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

	return guard.requireProjectTeam(project)
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

	return guard.requireProjectTeam(milestone.Project)
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

// requireOrganization compares a resolved organization id (read from an
// entity such as a ProjectLabel) against the guard's resolved organization.
// It is the Org-Scoped Write hard stop: organization-owned entities have no
// team to compare, so organization membership is the whole check.
func (guard writeGuard) requireOrganization(orgID string) error {
	if orgID != guard.target.Org.ID {
		return guard.organizationMismatchError(orgID)
	}

	return nil
}

// organizationMismatchError reports a Target Mismatch between the pinned
// organization and the organization resolved from an existing entity.
func (guard writeGuard) organizationMismatchError(orgID string) error {
	return fmt.Errorf(
		"%w: expected org_id=%s resolved org_id=%s",
		ErrTargetMismatch,
		guard.target.Org.ID,
		orgID,
	)
}

// requireIssueLabel resolves an IssueLabel for a taxonomy mutation (update,
// retire, restore). A team-scoped label must match the resolved team and must
// not be combined with orgWide. An organization-wide label (null team)
// requires orgWide and fails closed otherwise.
func (guard writeGuard) requireIssueLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	labelID string,
	orgWide bool,
) error {
	label, err := GetLabelByID(ctx, graphqlClient, labelID)
	if err != nil {
		return err
	}
	if label.TeamID == "" {
		if !orgWide {
			return fmt.Errorf(
				"%w: label %s is organization-wide; pass --org-wide to act on it",
				ErrTargetMismatch,
				labelID,
			)
		}

		return nil
	}
	if orgWide {
		return fmt.Errorf(
			"%w: label %s is team-scoped; --org-wide only applies to organization-wide labels",
			ErrWriteInvalid,
			labelID,
		)
	}
	if label.TeamID != guard.target.Team.ID || label.TeamKey != guard.target.Team.Key {
		return guard.teamMismatchError("label", label.TeamID, label.TeamKey)
	}

	return nil
}

// requireLabelParentScope resolves a candidate parent IssueLabel and confirms
// it shares the effective scope of the label being created: the resolved team
// by default, or organization-wide when orgWide is set.
func (guard writeGuard) requireLabelParentScope(
	ctx context.Context,
	graphqlClient graphql.Client,
	parentID string,
	orgWide bool,
) error {
	parent, err := GetLabelByID(ctx, graphqlClient, parentID)
	if err != nil {
		return err
	}
	if orgWide {
		if parent.TeamID != "" {
			return fmt.Errorf(
				"%w: parent label %s is team-scoped; an organization-wide label requires an organization-wide parent",
				ErrTargetMismatch,
				parentID,
			)
		}

		return nil
	}
	if parent.TeamID != guard.target.Team.ID || parent.TeamKey != guard.target.Team.Key {
		return guard.teamMismatchError("parent label", parent.TeamID, parent.TeamKey)
	}

	return nil
}

// requireAttachableLabel resolves an IssueLabel for attaching to or removing
// from an issue: a team-scoped label must match the resolved team, and an
// organization-wide label (null team) is always attachable within the
// resolved organization. There is no --org-wide flag for association writes.
func (guard writeGuard) requireAttachableLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	labelID string,
) error {
	label, err := GetLabelByID(ctx, graphqlClient, labelID)
	if err != nil {
		return err
	}
	if label.TeamID == "" {
		return nil
	}
	if label.TeamID != guard.target.Team.ID || label.TeamKey != guard.target.Team.Key {
		return guard.teamMismatchError("label", label.TeamID, label.TeamKey)
	}

	return nil
}

// requireProjectLabel resolves a ProjectLabel and confirms it belongs to the
// resolved organization. ProjectLabel is organization-owned; there is no team
// scope to compare.
func (guard writeGuard) requireProjectLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	labelID string,
) error {
	label, err := GetProjectLabelByID(ctx, graphqlClient, labelID)
	if err != nil {
		return err
	}

	return guard.requireOrganization(label.OrgID)
}

// requireProjectTeam compares a resolved project's teams against the pinned
// team. An unmatched result on a truncated team page fails closed instead of
// silently trusting the first 50 teams.
func (guard writeGuard) requireProjectTeam(project ProjectSummary) error {
	if projectHasTeam(project, guard.target.Team.ID, guard.target.Team.Key) {
		return nil
	}
	if project.TeamsTruncated {
		return fmt.Errorf(
			"%w: project %s has more than 50 teams; cannot verify pinned team membership",
			ErrTargetMismatch,
			project.ID,
		)
	}

	return guard.teamNotAttachedError()
}

func projectHasTeam(project ProjectSummary, teamID string, teamKey string) bool {
	for _, team := range project.Teams {
		if team.ID == teamID && team.Key == teamKey {
			return true
		}
	}

	return false
}
