package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addTeamCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	teamCommand := newGroupCommand("team", "Read Linear teams and create one")
	addTeamCreateCommand(ctx, teamCommand, options)
	addTeamListCommand(ctx, teamCommand, options)
	addTeamGetCommand(ctx, teamCommand, options)
	addTeamCyclesCommand(ctx, teamCommand, options)
	addTeamIssuesCommand(ctx, teamCommand, options)
	addTeamLabelsCommand(ctx, teamCommand, options)
	addTeamMembersCommand(ctx, teamCommand, options)
	addTeamMembershipsCommand(ctx, teamCommand, options)
	addTeamProjectsCommand(ctx, teamCommand, options)
	addTeamReleasePipelinesCommand(ctx, teamCommand, options)
	addTeamStatesCommand(ctx, teamCommand, options)
	addTeamGitAutomationStatesCommand(ctx, teamCommand, options)
	addTeamTemplatesCommand(ctx, teamCommand, options)
	root.AddCommand(teamCommand)
}

func addTeamListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamList, client.TeamSummary]{
		Use:       "list",
		Short:     "List visible teams",
		LimitHelp: "teams",
		Args:      cobra.NoArgs,
		Load:      loadTeamList,
		WriteItem: writeTeam,
	})
}

func addTeamGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.TeamSummary]{
		Use:   "get TEAM_ID",
		Short: "Get one team by id",
		Configure: func(command *cobra.Command) {
			command.ValidArgsFunction = firstArgCompletion(ctx, options, teamKeyCandidates)
		},
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.TeamSummary, error) {
			return client.GetTeamByID(ctx, runtime.graphqlClient, id)
		},
		Write: writeTeam,
	})
}

func addTeamMembersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamMemberList, client.UserSummary]{
		Use:       "members TEAM_ID",
		Short:     "List team members",
		LimitHelp: "members",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamMemberList,
		WriteItem: writeUser,
	})
}

func addTeamCyclesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.CycleList, client.CycleSummary]{
		Use:       "cycles TEAM_ID",
		Short:     "List team Cycles",
		LimitHelp: "Cycles",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamCycles,
		WriteItem: writeCycle,
	})
}

func addTeamIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "issues TEAM_ID",
		Short:     "List team issues",
		LimitHelp: "issues",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamIssues,
		WriteItem: writeIssue,
	})
}

func addTeamLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.LabelList, client.LabelSummary]{
		Use:       "labels TEAM_ID",
		Short:     "List team labels",
		LimitHelp: "labels",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamLabels,
		WriteItem: writeLabel,
	})
}

func addTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamMembershipList, client.TeamMembershipSummary]{
		Use:       "memberships TEAM_ID",
		Short:     "List team memberships",
		LimitHelp: "memberships",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamMemberships,
		WriteItem: writeTeamMembership,
	})
}

func addTeamProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectList, client.ProjectSummary]{
		Use:       "projects TEAM_ID",
		Short:     "List team projects",
		LimitHelp: "projects",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamProjects,
		WriteItem: writeProject,
	})
}

func addTeamReleasePipelinesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ReleasePipelineList, client.ReleasePipelineSummary]{
		Use:       "release-pipelines TEAM_ID",
		Short:     "List team release pipelines",
		LimitHelp: "release pipelines",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamReleasePipelines,
		WriteItem: writeReleasePipeline,
	})
}

func addTeamStatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.WorkflowStateList, client.WorkflowStateSummary]{
		Use:       "states TEAM_ID",
		Short:     "List team workflow states",
		LimitHelp: "workflow states",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamStates,
		WriteItem: writeWorkflowState,
	})
}

func addTeamGitAutomationStatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.GitAutomationStateList, client.GitAutomationStateSummary]{
		Use:       "git-automation-states TEAM_ID",
		Short:     "List team Git automation states",
		LimitHelp: "Git automation states",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamGitAutomationStates,
		WriteItem: writeGitAutomationState,
	})
}

func addTeamTemplatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TemplateList, client.TemplateSummary]{
		Use:       "templates TEAM_ID",
		Short:     "List team templates",
		LimitHelp: "templates",
		Args:      cobra.ExactArgs(1),
		Load:      loadTeamTemplates,
		WriteItem: writeTemplate,
	})
}

func writeTeam(command *cobra.Command, options *rootOptions, team client.TeamSummary) error {
	return writeItemLine(command, options, team, team.ID, "%s %s %s", team.ID, team.Key, team.Name)
}

func writeGitAutomationState(
	command *cobra.Command,
	options *rootOptions,
	state client.GitAutomationStateSummary,
) error {
	return writeItemLine(
		command, options, state, state.ID,
		"%s %s state %s target %s",
		state.ID,
		state.Event,
		emptyDash(state.StateName),
		emptyDash(state.TargetBranchPattern),
	)
}

func loadTeamList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, error) {
	teams, err := client.ListTeams(ctx, runtime.graphqlClient, limit)
	return teams, err
}

func loadTeamMemberList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMemberList, error) {
	members, err := client.ListTeamMembers(ctx, runtime.graphqlClient, args[0], limit)
	return members, err
}

func loadTeamCycles(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.CycleList, error) {
	cycles, err := client.ListTeamCycles(ctx, runtime.graphqlClient, args[0], limit)
	return cycles, err
}

func loadTeamIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListTeamIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, err
}

func loadTeamLabels(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.LabelList, error) {
	labels, err := client.ListTeamLabels(ctx, runtime.graphqlClient, args[0], limit)
	return labels, err
}

func loadTeamMemberships(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMembershipList, error) {
	memberships, err := client.ListTeamMembershipsForTeam(ctx, runtime.graphqlClient, args[0], limit)
	return memberships, err
}

func loadTeamProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectList, error) {
	projects, err := client.ListTeamProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, err
}

func loadTeamReleasePipelines(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ReleasePipelineList, error) {
	pipelines, err := client.ListTeamReleasePipelines(ctx, runtime.graphqlClient, args[0], limit)
	return pipelines, err
}

func loadTeamStates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.WorkflowStateList, error) {
	states, err := client.ListTeamWorkflowStates(ctx, runtime.graphqlClient, args[0], limit)
	return states, err
}

func loadTeamGitAutomationStates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.GitAutomationStateList, error) {
	states, err := client.ListTeamGitAutomationStates(ctx, runtime.graphqlClient, args[0], limit)
	return states, err
}

func loadTeamTemplates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TemplateList, error) {
	templates, err := client.ListTeamTemplates(ctx, runtime.graphqlClient, args[0], limit)
	return templates, err
}
