package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTeamCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	teamCommand := &cobra.Command{
		Use:   "team",
		Short: "Read Linear teams",
	}
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
		Use:           "list",
		Short:         "List visible teams",
		LimitHelp:     "teams",
		Args:          cobra.NoArgs,
		Load:          loadTeamList,
		PageWithItems: teamPageWithItems,
		WriteItem:     writeTeam,
	})
}

func addTeamGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:               "get TEAM_ID",
		Short:             "Get one Team by id",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: firstArgCompletion(ctx, options, teamKeyCandidates),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			team, err := client.GetTeamByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeTeam(command, options, team)
		},
	})
}

func addTeamMembersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamMemberList, client.UserSummary]{
		Use:           "members TEAM_ID",
		Short:         "List team members",
		LimitHelp:     "members",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamMemberList,
		PageWithItems: teamMemberPageWithItems,
		WriteItem:     writeUser,
	})
}

func addTeamCyclesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.CycleList, client.CycleSummary]{
		Use:           "cycles TEAM_ID",
		Short:         "List team Cycles",
		LimitHelp:     "Cycles",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamCycles,
		PageWithItems: cyclePageWithItems,
		WriteItem:     writeCycle,
	})
}

func addTeamIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:           "issues TEAM_ID",
		Short:         "List team issues",
		LimitHelp:     "issues",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamIssues,
		PageWithItems: issuePageWithItems,
		WriteItem:     writeIssue,
	})
}

func addTeamLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.LabelList, client.LabelSummary]{
		Use:           "labels TEAM_ID",
		Short:         "List team labels",
		LimitHelp:     "labels",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamLabels,
		PageWithItems: labelPageWithItems,
		WriteItem:     writeLabel,
	})
}

func addTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamMembershipList, client.TeamMembershipSummary]{
		Use:           "memberships TEAM_ID",
		Short:         "List team memberships",
		LimitHelp:     "memberships",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamMemberships,
		PageWithItems: teamMembershipPageWithItems,
		WriteItem:     writeTeamMembership,
	})
}

func addTeamProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectList, client.ProjectSummary]{
		Use:           "projects TEAM_ID",
		Short:         "List team projects",
		LimitHelp:     "projects",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamProjects,
		PageWithItems: projectPageWithItems,
		WriteItem:     writeProject,
	})
}

func addTeamReleasePipelinesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ReleasePipelineList, client.ReleasePipelineSummary]{
		Use:           "release-pipelines TEAM_ID",
		Short:         "List team release pipelines",
		LimitHelp:     "release pipelines",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamReleasePipelines,
		PageWithItems: releasePipelinePageWithItems,
		WriteItem:     writeReleasePipeline,
	})
}

func addTeamStatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.WorkflowStateList, client.WorkflowStateSummary]{
		Use:           "states TEAM_ID",
		Short:         "List team workflow states",
		LimitHelp:     "workflow states",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamStates,
		PageWithItems: workflowStatePageWithItems,
		WriteItem:     writeWorkflowState,
	})
}

func addTeamGitAutomationStatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.GitAutomationStateList, client.GitAutomationStateSummary]{
		Use:           "git-automation-states TEAM_ID",
		Short:         "List team Git automation states",
		LimitHelp:     "Git automation states",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamGitAutomationStates,
		PageWithItems: gitAutomationStatePageWithItems,
		WriteItem:     writeGitAutomationState,
	})
}

func addTeamTemplatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TemplateList, client.TemplateSummary]{
		Use:           "templates TEAM_ID",
		Short:         "List team templates",
		LimitHelp:     "templates",
		Args:          cobra.ExactArgs(1),
		Load:          loadTeamTemplates,
		PageWithItems: templatePageWithItems,
		WriteItem:     writeTemplate,
	})
}

func writeTeam(command *cobra.Command, options *rootOptions, team client.TeamSummary) error {
	return writeItem(command, options, team, team.ID,
		func(command *cobra.Command, _ *rootOptions, team client.TeamSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s %s", team.ID, team.Key, team.Name)
		})
}

func writeGitAutomationState(
	command *cobra.Command,
	options *rootOptions,
	state client.GitAutomationStateSummary,
) error {
	return writeItem(command, options, state, state.ID,
		func(command *cobra.Command, _ *rootOptions, state client.GitAutomationStateSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s state %s target %s",
				state.ID,
				state.Event,
				emptyDash(state.StateName),
				emptyDash(state.TargetBranchPattern),
			)
		})
}

func loadTeamList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, []client.TeamSummary, error) {
	teams, err := client.ListTeams(ctx, runtime.graphqlClient, limit)
	return teams, teams.Teams, err
}

func loadTeamMemberList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMemberList, []client.UserSummary, error) {
	members, err := client.ListTeamMembers(ctx, runtime.graphqlClient, args[0], limit)
	return members, members.Members, err
}

func loadTeamCycles(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.CycleList, []client.CycleSummary, error) {
	cycles, err := client.ListTeamCycles(ctx, runtime.graphqlClient, args[0], limit)
	return cycles, cycles.Cycles, err
}

func loadTeamIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListTeamIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func loadTeamLabels(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.LabelList, []client.LabelSummary, error) {
	labels, err := client.ListTeamLabels(ctx, runtime.graphqlClient, args[0], limit)
	return labels, labels.Labels, err
}

func loadTeamMemberships(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMembershipList, []client.TeamMembershipSummary, error) {
	memberships, err := client.ListTeamMembershipsForTeam(ctx, runtime.graphqlClient, args[0], limit)
	return memberships, memberships.Memberships, err
}

func loadTeamProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectList, []client.ProjectSummary, error) {
	projects, err := client.ListTeamProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, projects.Projects, err
}

func loadTeamReleasePipelines(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ReleasePipelineList, []client.ReleasePipelineSummary, error) {
	pipelines, err := client.ListTeamReleasePipelines(ctx, runtime.graphqlClient, args[0], limit)
	return pipelines, pipelines.ReleasePipelines, err
}

func loadTeamStates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.WorkflowStateList, []client.WorkflowStateSummary, error) {
	states, err := client.ListTeamWorkflowStates(ctx, runtime.graphqlClient, args[0], limit)
	return states, states.WorkflowStates, err
}

func loadTeamGitAutomationStates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.GitAutomationStateList, []client.GitAutomationStateSummary, error) {
	states, err := client.ListTeamGitAutomationStates(ctx, runtime.graphqlClient, args[0], limit)
	return states, states.States, err
}

func loadTeamTemplates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TemplateList, []client.TemplateSummary, error) {
	templates, err := client.ListTeamTemplates(ctx, runtime.graphqlClient, args[0], limit)
	return templates, templates.Templates, err
}
