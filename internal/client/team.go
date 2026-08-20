package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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

//nolint:lll
type teamMembersNode = gql.XTeam_membersTeamMembersUserConnectionNodesUser

//nolint:lll
type teamCyclesNode = gql.XTeam_cyclesTeamCyclesCycleConnectionNodesCycle

//nolint:lll
type teamIssuesNode = gql.XTeam_issuesTeamIssuesIssueConnectionNodesIssue

//nolint:lll
type teamLabelsNode = gql.XTeam_labelsTeamLabelsIssueLabelConnectionNodesIssueLabel

//nolint:lll
type teamMembershipsForTeamNode = gql.XTeam_membershipsTeamMembershipsTeamMembershipConnectionNodesTeamMembership

//nolint:lll
type teamProjectsNode = gql.XTeam_projectsTeamProjectsProjectConnectionNodesProject

//nolint:lll
type teamReleasePipelinesNode = gql.XTeam_releasePipelinesTeamReleasePipelinesReleasePipelineConnectionNodesReleasePipeline

//nolint:lll
type teamWorkflowStatesNode = gql.XTeam_statesTeamStatesWorkflowStateConnectionNodesWorkflowState

//nolint:lll
type teamGitAutomationStatesNode = gql.XTeam_gitAutomationStatesTeamGitAutomationStatesGitAutomationStateConnectionNodesGitAutomationState

//nolint:lll
type teamTemplatesNode = gql.XTeam_templatesTeamTemplatesTemplateConnectionNodesTemplate

type teamScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	teamID        string
	teamKey       string
	teamName      string
}

// ListTeams returns visible teams.
func ListTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	page, err := listConnection(
		"list teams", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]gql.XTeams_listTeamsTeamConnectionNodesTeam, bool, *string, error) {
			result, err := gql.XTeams_list(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.Teams.Nodes, result.Teams.PageInfo.HasNextPage, result.Teams.PageInfo.EndCursor, nil
		},
		func(team gql.XTeams_listTeamsTeamConnectionNodesTeam) TeamSummary {
			return teamSummary(team.TeamSummaryFields)
		},
	)
	if err != nil {
		return TeamList{}, err
	}

	return TeamList{Teams: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetTeamByID returns one Team by id.
func GetTeamByID(ctx context.Context, graphqlClient graphql.Client, id string) (TeamSummary, error) {
	teamResult, err := gql.XTeam(ctx, graphqlClient, id)
	if err != nil {
		return TeamSummary{}, fmt.Errorf("get team %s: %w", id, err)
	}

	return teamSummary(teamResult.Team.TeamSummaryFields), nil
}

// ListTeamMembers returns visible members for one Team.
func ListTeamMembers(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TeamMemberList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team members "+id, limit, defaultListPageSize,
		query.members,
		teamMembersNodeSummary,
	)
	if err != nil {
		return TeamMemberList{}, err
	}

	return TeamMemberList{
		TeamID:      query.teamID,
		TeamKey:     query.teamKey,
		TeamName:    query.teamName,
		Members:     page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListTeamCycles returns Cycles associated with one Team.
func ListTeamCycles(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (CycleList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team cycles "+id, limit, defaultListPageSize,
		query.cycles,
		teamCyclesNodeSummary,
	)
	if err != nil {
		return CycleList{}, err
	}

	return CycleList{Cycles: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListTeamIssues returns issues associated with one Team.
func ListTeamIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (IssueList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team issues "+id, limit, defaultListPageSize,
		query.issues,
		teamIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListTeamLabels returns IssueLabels associated with one Team.
func ListTeamLabels(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (LabelList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team labels "+id, limit, defaultListPageSize,
		query.labels,
		teamLabelsNodeSummary,
	)
	if err != nil {
		return LabelList{}, err
	}

	return LabelList{Labels: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListTeamMembershipsForTeam returns TeamMemberships associated with one Team.
func ListTeamMembershipsForTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (TeamMembershipList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team memberships "+id, limit, defaultListPageSize,
		query.memberships,
		teamMembershipsForTeamNodeSummary,
	)
	if err != nil {
		return TeamMembershipList{}, err
	}

	return TeamMembershipList{Memberships: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListTeamProjects returns Projects associated with one Team.
func ListTeamProjects(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (ProjectList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team projects "+id, limit, defaultListPageSize,
		query.projects,
		teamProjectsNodeSummary,
	)
	if err != nil {
		return ProjectList{}, err
	}

	return ProjectList{Projects: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListTeamReleasePipelines returns ReleasePipelines associated with one Team.
func ListTeamReleasePipelines(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleasePipelineList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team release pipelines "+id, limit, defaultListPageSize,
		query.releasePipelines,
		teamReleasePipelinesNodeSummary,
	)
	if err != nil {
		return ReleasePipelineList{}, err
	}

	return ReleasePipelineList{
		ReleasePipelines: page.Items,
		HasNextPage:      page.HasNextPage,
		EndCursor:        page.EndCursor,
	}, nil
}

// ListTeamWorkflowStates returns workflow states associated with one Team.
func ListTeamWorkflowStates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (WorkflowStateList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team states "+id, limit, defaultListPageSize,
		query.states,
		teamWorkflowStatesNodeSummary,
	)
	if err != nil {
		return WorkflowStateList{}, err
	}

	return WorkflowStateList{
		WorkflowStates: page.Items,
		HasNextPage:    page.HasNextPage,
		EndCursor:      page.EndCursor,
	}, nil
}

// ListTeamGitAutomationStates returns Git automation rules associated with one Team.
func ListTeamGitAutomationStates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (GitAutomationStateList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team git automation states "+id, limit, defaultListPageSize,
		query.gitAutomationStates,
		teamGitAutomationStatesNodeSummary,
	)
	if err != nil {
		return GitAutomationStateList{}, err
	}

	return GitAutomationStateList{
		TeamID:      query.teamID,
		TeamKey:     query.teamKey,
		TeamName:    query.teamName,
		States:      page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListTeamTemplates returns Templates associated with one Team.
func ListTeamTemplates(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TemplateList, error) {
	query := &teamScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list team templates "+id, limit, defaultListPageSize,
		query.templates,
		teamTemplatesNodeSummary,
	)
	if err != nil {
		return TemplateList{}, err
	}

	return TemplateList{
		Templates:   page.Items,
		TotalCount:  len(page.Items),
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

func (query *teamScopedQuery) members(pageSize int, after *string) ([]teamMembersNode, bool, *string, error) {
	result, err := gql.XTeam_members(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.teamID = result.Team.Id
	query.teamKey = result.Team.Key
	query.teamName = result.Team.Name

	return result.Team.Members.Nodes,
		result.Team.Members.PageInfo.HasNextPage,
		result.Team.Members.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) cycles(pageSize int, after *string) ([]teamCyclesNode, bool, *string, error) {
	result, err := gql.XTeam_cycles(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.Cycles.Nodes,
		result.Team.Cycles.PageInfo.HasNextPage,
		result.Team.Cycles.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) issues(pageSize int, after *string) ([]teamIssuesNode, bool, *string, error) {
	result, err := gql.XTeam_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.Issues.Nodes,
		result.Team.Issues.PageInfo.HasNextPage,
		result.Team.Issues.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) labels(pageSize int, after *string) ([]teamLabelsNode, bool, *string, error) {
	result, err := gql.XTeam_labels(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.Labels.Nodes,
		result.Team.Labels.PageInfo.HasNextPage,
		result.Team.Labels.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) memberships(
	pageSize int,
	after *string,
) ([]teamMembershipsForTeamNode, bool, *string, error) {
	result, err := gql.XTeam_memberships(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.Memberships.Nodes,
		result.Team.Memberships.PageInfo.HasNextPage,
		result.Team.Memberships.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) projects(pageSize int, after *string) ([]teamProjectsNode, bool, *string, error) {
	result, err := gql.XTeam_projects(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.Projects.Nodes,
		result.Team.Projects.PageInfo.HasNextPage,
		result.Team.Projects.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) releasePipelines(
	pageSize int,
	after *string,
) ([]teamReleasePipelinesNode, bool, *string, error) {
	result, err := gql.XTeam_releasePipelines(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.ReleasePipelines.Nodes,
		result.Team.ReleasePipelines.PageInfo.HasNextPage,
		result.Team.ReleasePipelines.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) states(pageSize int, after *string) ([]teamWorkflowStatesNode, bool, *string, error) {
	result, err := gql.XTeam_states(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.States.Nodes,
		result.Team.States.PageInfo.HasNextPage,
		result.Team.States.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) gitAutomationStates(
	pageSize int,
	after *string,
) ([]teamGitAutomationStatesNode, bool, *string, error) {
	result, err := gql.XTeam_gitAutomationStates(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.teamID = result.Team.Id
	query.teamKey = result.Team.Key
	query.teamName = result.Team.Name

	return result.Team.GitAutomationStates.Nodes,
		result.Team.GitAutomationStates.PageInfo.HasNextPage,
		result.Team.GitAutomationStates.PageInfo.EndCursor,
		nil
}

func (query *teamScopedQuery) templates(pageSize int, after *string) ([]teamTemplatesNode, bool, *string, error) {
	result, err := gql.XTeam_templates(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Team.Templates.Nodes,
		result.Team.Templates.PageInfo.HasNextPage,
		result.Team.Templates.PageInfo.EndCursor,
		nil
}

func teamMembersNodeSummary(member teamMembersNode) UserSummary {
	return userSummary(member.UserSummaryFields)
}

func teamCyclesNodeSummary(cycle teamCyclesNode) CycleSummary {
	return cycleSummary(cycle.CycleSummaryFields)
}

func teamIssuesNodeSummary(issue teamIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func teamLabelsNodeSummary(label teamLabelsNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func teamMembershipsForTeamNodeSummary(membership teamMembershipsForTeamNode) TeamMembershipSummary {
	return teamMembershipSummary(membership.TeamMembershipSummaryFields)
}

func teamProjectsNodeSummary(project teamProjectsNode) ProjectSummary {
	return projectSummaryFromFields(project.ProjectSummaryFields)
}

func teamReleasePipelinesNodeSummary(pipeline teamReleasePipelinesNode) ReleasePipelineSummary {
	return releasePipelineSummary(pipeline.ReleasePipelineSummaryFields)
}

func teamWorkflowStatesNodeSummary(state teamWorkflowStatesNode) WorkflowStateSummary {
	return workflowStateSummary(state.WorkflowStateSummaryFields)
}

func teamGitAutomationStatesNodeSummary(state teamGitAutomationStatesNode) GitAutomationStateSummary {
	return gitAutomationStateSummary(state.GitAutomationStateSummaryFields)
}

func teamTemplatesNodeSummary(template teamTemplatesNode) TemplateSummary {
	return templateSummary(template.TemplateSummaryFields)
}

func teamSummary(team gql.TeamSummaryFields) TeamSummary {
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

func gitAutomationStateSummary(fields gql.GitAutomationStateSummaryFields) GitAutomationStateSummary {
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
