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
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List visible teams",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(ctx, command, nil, options, limit, loadTeamList, teamPageWithItems, writeTeam)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum teams to return")
	root.AddCommand(command)
}

func addTeamGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get TEAM_ID",
		Short: "Get one Team by id",
		Args:  cobra.ExactArgs(1),
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
	limit := 50
	command := &cobra.Command{
		Use:   "members TEAM_ID",
		Short: "List team members",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadTeamMemberList,
				teamMemberPageWithItems,
				writeUser,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum members to return")
	root.AddCommand(command)
}

func addTeamCyclesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "cycles TEAM_ID",
		Short: "List team Cycles",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadTeamCycles, cyclePageWithItems, writeCycle,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum Cycles to return")
	root.AddCommand(command)
}

func addTeamIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "issues TEAM_ID",
		Short: "List team issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadTeamIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addTeamLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "labels TEAM_ID",
		Short: "List team labels",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadTeamLabels, labelPageWithItems, writeLabel,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum labels to return")
	root.AddCommand(command)
}

func addTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "memberships TEAM_ID",
		Short: "List team memberships",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadTeamMemberships,
				teamMembershipPageWithItems,
				writeTeamMembership,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum memberships to return")
	root.AddCommand(command)
}

func addTeamProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "projects TEAM_ID",
		Short: "List team projects",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadTeamProjects, projectPageWithItems, writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func addTeamReleasePipelinesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "release-pipelines TEAM_ID",
		Short: "List team release pipelines",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadTeamReleasePipelines,
				releasePipelinePageWithItems,
				writeReleasePipeline,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum release pipelines to return")
	root.AddCommand(command)
}

func addTeamStatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "states TEAM_ID",
		Short: "List team workflow states",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadTeamStates,
				workflowStatePageWithItems,
				writeWorkflowState,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum workflow states to return")
	root.AddCommand(command)
}

func addTeamGitAutomationStatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "git-automation-states TEAM_ID",
		Short: "List team Git automation states",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadTeamGitAutomationStates,
				gitAutomationStatePageWithItems,
				writeGitAutomationState,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum Git automation states to return")
	root.AddCommand(command)
}

func addTeamTemplatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "templates TEAM_ID",
		Short: "List team templates",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadTeamTemplates, templatePageWithItems, writeTemplate,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum templates to return")
	root.AddCommand(command)
}

func writeTeam(command *cobra.Command, options *rootOptions, team client.TeamSummary) error {
	if wrote, err := writeIDOnly(command, options, team.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, team)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s %s", team.ID, team.Key, team.Name)
}

func writeGitAutomationState(
	command *cobra.Command,
	options *rootOptions,
	state client.GitAutomationStateSummary,
) error {
	if wrote, err := writeIDOnly(command, options, state.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, state)
	}

	return render.WriteLine(
		command.OutOrStdout(),
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
) (client.TeamList, []client.TeamSummary, error) {
	teams, err := client.ListTeams(ctx, runtime.graphqlClient, limit)
	return teams, teams.Teams, err
}

func teamPageWithItems(page client.TeamList, teams []client.TeamSummary) client.TeamList {
	page.Teams = teams
	return page
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

func teamMemberPageWithItems(page client.TeamMemberList, members []client.UserSummary) client.TeamMemberList {
	page.Members = members
	return page
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

func cyclePageWithItems(page client.CycleList, cycles []client.CycleSummary) client.CycleList {
	page.Cycles = cycles
	return page
}

func issuePageWithItems(page client.IssueList, issues []client.IssueSummary) client.IssueList {
	page.Issues = issues
	return page
}

func gitAutomationStatePageWithItems(
	page client.GitAutomationStateList,
	states []client.GitAutomationStateSummary,
) client.GitAutomationStateList {
	page.States = states
	return page
}
