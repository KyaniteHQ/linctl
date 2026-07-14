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

// ListTeams returns visible teams.
func ListTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	teams, err := gql.XTeams_list(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamList{}, fmt.Errorf("list teams: %w", err)
	}

	summaries := mapNodes(teams.Teams.Nodes, func(team gql.XTeams_listTeamsTeamConnectionNodesTeam) TeamSummary {
		return teamSummary(team.TeamSummaryFields)
	})

	return TeamList{
		Teams:       summaries,
		HasNextPage: teams.Teams.PageInfo.HasNextPage,
		EndCursor:   teams.Teams.PageInfo.EndCursor,
	}, nil
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
	team, err := gql.XTeam_members(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamMemberList{}, fmt.Errorf("list team members %s: %w", id, err)
	}

	members := mapNodes(team.Team.Members.Nodes, func(
		member gql.XTeam_membersTeamMembersUserConnectionNodesUser,
	) UserSummary {
		return userSummary(member.UserSummaryFields)
	})

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
	team, err := gql.XTeam_cycles(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CycleList{}, fmt.Errorf("list team cycles %s: %w", id, err)
	}

	cycles := mapNodes(team.Team.Cycles.Nodes, func(
		cycle gql.XTeam_cyclesTeamCyclesCycleConnectionNodesCycle,
	) CycleSummary {
		return cycleSummary(cycle.CycleSummaryFields)
	})

	return CycleList{
		Cycles:      cycles,
		HasNextPage: team.Team.Cycles.PageInfo.HasNextPage,
		EndCursor:   team.Team.Cycles.PageInfo.EndCursor,
	}, nil
}

// ListTeamIssues returns issues associated with one Team.
func ListTeamIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (IssueList, error) {
	team, err := gql.XTeam_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list team issues %s: %w", id, err)
	}

	issues := mapNodes(team.Team.Issues.Nodes, func(
		issue gql.XTeam_issuesTeamIssuesIssueConnectionNodesIssue,
	) IssueSummary {
		return issueSummaryFromFields(issue.IssueSummaryFields)
	})

	return IssueList{
		Issues:      issues,
		HasNextPage: team.Team.Issues.PageInfo.HasNextPage,
		EndCursor:   team.Team.Issues.PageInfo.EndCursor,
	}, nil
}

// ListTeamLabels returns IssueLabels associated with one Team.
func ListTeamLabels(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (LabelList, error) {
	team, err := gql.XTeam_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list team labels %s: %w", id, err)
	}

	labels := mapNodes(team.Team.Labels.Nodes, func(
		label gql.XTeam_labelsTeamLabelsIssueLabelConnectionNodesIssueLabel,
	) LabelSummary {
		return labelSummary(label.IssueLabelSummaryFields)
	})

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
	team, err := gql.XTeam_memberships(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list team memberships %s: %w", id, err)
	}

	memberships := mapNodes(team.Team.Memberships.Nodes, func(
		membership gql.XTeam_membershipsTeamMembershipsTeamMembershipConnectionNodesTeamMembership,
	) TeamMembershipSummary {
		return teamMembershipSummary(membership.TeamMembershipSummaryFields)
	})

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: team.Team.Memberships.PageInfo.HasNextPage,
		EndCursor:   team.Team.Memberships.PageInfo.EndCursor,
	}, nil
}

// ListTeamProjects returns Projects associated with one Team.
func ListTeamProjects(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (ProjectList, error) {
	team, err := gql.XTeam_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list team projects %s: %w", id, err)
	}

	projects := mapNodes(team.Team.Projects.Nodes, func(
		project gql.XTeam_projectsTeamProjectsProjectConnectionNodesProject,
	) ProjectSummary {
		return projectSummaryFromFields(project.ProjectSummaryFields)
	})

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
	team, err := gql.XTeam_releasePipelines(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleasePipelineList{}, fmt.Errorf("list team release pipelines %s: %w", id, err)
	}

	pipelines := mapNodes(team.Team.ReleasePipelines.Nodes, func(
		pipeline gql.XTeam_releasePipelinesTeamReleasePipelinesReleasePipelineConnectionNodesReleasePipeline,
	) ReleasePipelineSummary {
		return releasePipelineSummary(pipeline.ReleasePipelineSummaryFields)
	})

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
	team, err := gql.XTeam_states(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return WorkflowStateList{}, fmt.Errorf("list team states %s: %w", id, err)
	}

	states := mapNodes(team.Team.States.Nodes, func(
		state gql.XTeam_statesTeamStatesWorkflowStateConnectionNodesWorkflowState,
	) WorkflowStateSummary {
		return workflowStateSummary(state.WorkflowStateSummaryFields)
	})

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
	team, err := gql.XTeam_gitAutomationStates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return GitAutomationStateList{}, fmt.Errorf("list team git automation states %s: %w", id, err)
	}

	states := mapNodes(team.Team.GitAutomationStates.Nodes, func(
		state gql.XTeam_gitAutomationStatesTeamGitAutomationStatesGitAutomationStateConnectionNodesGitAutomationState,
	) GitAutomationStateSummary {
		return gitAutomationStateSummary(state.GitAutomationStateSummaryFields)
	})

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
	team, err := gql.XTeam_templates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TemplateList{}, fmt.Errorf("list team templates %s: %w", id, err)
	}

	templates := mapNodes(team.Team.Templates.Nodes, func(
		template gql.XTeam_templatesTeamTemplatesTemplateConnectionNodesTemplate,
	) TemplateSummary {
		return templateSummary(template.TemplateSummaryFields)
	})

	return TemplateList{
		Templates:   templates,
		TotalCount:  len(templates),
		HasNextPage: team.Team.Templates.PageInfo.HasNextPage,
		EndCursor:   team.Team.Templates.PageInfo.EndCursor,
	}, nil
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
