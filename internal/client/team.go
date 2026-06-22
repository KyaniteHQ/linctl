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

// GitAutomationStateSummary is the compact Git automation rule model used by read-only commands.
type GitAutomationStateSummary struct {
	ID                  string `json:"id"`
	Event               string `json:"event"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	ArchivedAt          string `json:"archived_at,omitempty"`
	StateID             string `json:"state_id,omitempty"`
	StateName           string `json:"state_name,omitempty"`
	StateType           string `json:"state_type,omitempty"`
	TargetBranchID      string `json:"target_branch_id,omitempty"`
	TargetBranchPattern string `json:"target_branch_pattern,omitempty"`
	TargetBranchIsRegex bool   `json:"target_branch_is_regex"`
}

// GitAutomationStateList is a page of Git automation rules associated with one Team.
type GitAutomationStateList struct {
	TeamID      string                      `json:"team_id"`
	TeamKey     string                      `json:"team_key"`
	TeamName    string                      `json:"team_name"`
	States      []GitAutomationStateSummary `json:"git_automation_states"`
	HasNextPage bool                        `json:"has_next_page"`
	EndCursor   *string                     `json:"end_cursor,omitempty"`
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
	teamResult, err := team(ctx, graphqlClient, id)
	if err != nil {
		return TeamSummary{}, fmt.Errorf("get team %s: %w", id, err)
	}

	return teamSummary(teamResult.Team.TeamSummaryFields), nil
}

// ListTeamMembers returns visible members for one Team.
func ListTeamMembers(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TeamMemberList, error) {
	team, err := team_members(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
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

// ListTeamCycles returns Cycles associated with one Team.
func ListTeamCycles(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (CycleList, error) {
	team, err := team_cycles(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CycleList{}, fmt.Errorf("list team cycles %s: %w", id, err)
	}

	cycles := make([]CycleSummary, 0, len(team.Team.Cycles.Nodes))
	for _, cycle := range team.Team.Cycles.Nodes {
		cycles = append(cycles, cycleSummary(cycle.CycleSummaryFields))
	}

	return CycleList{
		Cycles:      cycles,
		HasNextPage: team.Team.Cycles.PageInfo.HasNextPage,
		EndCursor:   team.Team.Cycles.PageInfo.EndCursor,
	}, nil
}

// ListTeamIssues returns issues associated with one Team.
func ListTeamIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (IssueList, error) {
	team, err := team_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list team issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(team.Team.Issues.Nodes))
	for _, issue := range team.Team.Issues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: team.Team.Issues.PageInfo.HasNextPage,
		EndCursor:   team.Team.Issues.PageInfo.EndCursor,
	}, nil
}

// ListTeamLabels returns IssueLabels associated with one Team.
func ListTeamLabels(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (LabelList, error) {
	team, err := team_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list team labels %s: %w", id, err)
	}

	labels := make([]LabelSummary, 0, len(team.Team.Labels.Nodes))
	for _, label := range team.Team.Labels.Nodes {
		labels = append(labels, labelSummary(label.IssueLabelSummaryFields))
	}

	return LabelList{
		Labels:      labels,
		HasNextPage: team.Team.Labels.PageInfo.HasNextPage,
		EndCursor:   team.Team.Labels.PageInfo.EndCursor,
	}, nil
}

// ListTeamMembershipsForTeam returns TeamMemberships associated with one Team.
func ListTeamMembershipsForTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (TeamMembershipList, error) {
	team, err := team_memberships(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list team memberships %s: %w", id, err)
	}

	memberships := make([]TeamMembershipSummary, 0, len(team.Team.Memberships.Nodes))
	for _, membership := range team.Team.Memberships.Nodes {
		memberships = append(memberships, teamMembershipSummary(membership.TeamMembershipSummaryFields))
	}

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: team.Team.Memberships.PageInfo.HasNextPage,
		EndCursor:   team.Team.Memberships.PageInfo.EndCursor,
	}, nil
}

// ListTeamProjects returns Projects associated with one Team.
func ListTeamProjects(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (ProjectList, error) {
	team, err := team_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list team projects %s: %w", id, err)
	}

	projects := make([]ProjectSummary, 0, len(team.Team.Projects.Nodes))
	for _, project := range team.Team.Projects.Nodes {
		projects = append(projects, projectSummaryFromFields(project.ProjectSummaryFields))
	}

	return ProjectList{
		Projects:    projects,
		HasNextPage: team.Team.Projects.PageInfo.HasNextPage,
		EndCursor:   team.Team.Projects.PageInfo.EndCursor,
	}, nil
}

// ListTeamReleasePipelines returns ReleasePipelines associated with one Team.
func ListTeamReleasePipelines(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleasePipelineList, error) {
	team, err := team_releasePipelines(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleasePipelineList{}, fmt.Errorf("list team release pipelines %s: %w", id, err)
	}

	pipelines := make([]ReleasePipelineSummary, 0, len(team.Team.ReleasePipelines.Nodes))
	for _, pipeline := range team.Team.ReleasePipelines.Nodes {
		pipelines = append(pipelines, releasePipelineSummary(pipeline.ReleasePipelineSummaryFields))
	}

	return ReleasePipelineList{
		ReleasePipelines: pipelines,
		HasNextPage:      team.Team.ReleasePipelines.PageInfo.HasNextPage,
		EndCursor:        team.Team.ReleasePipelines.PageInfo.EndCursor,
	}, nil
}

// ListTeamWorkflowStates returns workflow states associated with one Team.
func ListTeamWorkflowStates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (WorkflowStateList, error) {
	team, err := team_states(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return WorkflowStateList{}, fmt.Errorf("list team states %s: %w", id, err)
	}

	states := make([]WorkflowStateSummary, 0, len(team.Team.States.Nodes))
	for _, state := range team.Team.States.Nodes {
		states = append(states, workflowStateSummary(state.WorkflowStateSummaryFields))
	}

	return WorkflowStateList{
		WorkflowStates: states,
		HasNextPage:    team.Team.States.PageInfo.HasNextPage,
		EndCursor:      team.Team.States.PageInfo.EndCursor,
	}, nil
}

// ListTeamGitAutomationStates returns Git automation rules associated with one Team.
func ListTeamGitAutomationStates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (GitAutomationStateList, error) {
	team, err := team_gitAutomationStates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return GitAutomationStateList{}, fmt.Errorf("list team git automation states %s: %w", id, err)
	}

	states := make([]GitAutomationStateSummary, 0, len(team.Team.GitAutomationStates.Nodes))
	for _, state := range team.Team.GitAutomationStates.Nodes {
		states = append(states, gitAutomationStateSummary(state.GitAutomationStateSummaryFields))
	}

	return GitAutomationStateList{
		TeamID:      team.Team.Id,
		TeamKey:     team.Team.Key,
		TeamName:    team.Team.Name,
		States:      states,
		HasNextPage: team.Team.GitAutomationStates.PageInfo.HasNextPage,
		EndCursor:   team.Team.GitAutomationStates.PageInfo.EndCursor,
	}, nil
}

// ListTeamTemplates returns Templates associated with one Team.
func ListTeamTemplates(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TemplateList, error) {
	team, err := team_templates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TemplateList{}, fmt.Errorf("list team templates %s: %w", id, err)
	}

	templates := make([]TemplateSummary, 0, len(team.Team.Templates.Nodes))
	for _, template := range team.Team.Templates.Nodes {
		templates = append(templates, templateSummary(template.TemplateSummaryFields))
	}

	return TemplateList{
		Templates:   templates,
		TotalCount:  len(templates),
		HasNextPage: team.Team.Templates.PageInfo.HasNextPage,
		EndCursor:   team.Team.Templates.PageInfo.EndCursor,
	}, nil
}

func teamSummary(team TeamSummaryFields) TeamSummary {
	return TeamSummary{
		ID:          team.Id,
		Key:         team.Key,
		Name:        team.Name,
		Description: stringValue(team.Description),
		ArchivedAt:  stringValue(team.ArchivedAt),
		OrgID:       team.Organization.Id,
		OrgName:     team.Organization.Name,
		OrgURLKey:   team.Organization.UrlKey,
	}
}

func gitAutomationStateSummary(fields GitAutomationStateSummaryFields) GitAutomationStateSummary {
	summary := GitAutomationStateSummary{
		ID:         fields.Id,
		Event:      string(fields.Event),
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
	}
	if fields.State != nil {
		summary.StateID = fields.State.Id
		summary.StateName = fields.State.Name
		summary.StateType = fields.State.Type
	}
	if fields.TargetBranch != nil {
		summary.TargetBranchID = fields.TargetBranch.Id
		summary.TargetBranchPattern = fields.TargetBranch.BranchPattern
		summary.TargetBranchIsRegex = fields.TargetBranch.IsRegex
	}

	return summary
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
