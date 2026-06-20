package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectMilestoneCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectMilestoneCommand := &cobra.Command{
		Use:   "project-milestone",
		Short: "Read and write Linear project milestones",
	}
	addProjectMilestoneListCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneGetCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneCreateCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneUpdateCommand(ctx, projectMilestoneCommand, options)
	root.AddCommand(projectMilestoneCommand)
}

func addProjectMilestoneListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list PROJECT_ID",
		Short: "List milestones for one project",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadProjectMilestoneList,
				projectMilestonePageWithItems,
				writeProjectMilestone,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum project milestones to return")
	root.AddCommand(command)
}

func addProjectMilestoneGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get PROJECT_MILESTONE_ID",
		Short: "Get one ProjectMilestone by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			milestone, err := client.GetProjectMilestoneByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeProjectMilestone(command, options, milestone)
		},
	})
}

func writeProjectMilestone(
	command *cobra.Command,
	options *rootOptions,
	milestone client.ProjectMilestoneSummary,
) error {
	if wrote, err := writeIDOnly(command, options, milestone.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, milestone)
	}

	format, err := normalizedHumanFormat(options)
	if err != nil {
		return err
	}
	if format == "minimal" {
		return render.WriteLine(command.OutOrStdout(), "%s", milestone.ID)
	}
	if format == "full" {
		return render.WriteLine(
			command.OutOrStdout(),
			"%s %s [%s] target_date=%s progress=%0.2f",
			milestone.ID,
			milestone.Name,
			milestone.Status,
			emptyDash(milestone.TargetDate),
			milestone.Progress,
		)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", milestone.ID, milestone.Name, milestone.Status)
}

func loadProjectMilestoneList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectMilestoneList, []client.ProjectMilestoneSummary, error) {
	milestones, err := client.ListProjectMilestones(ctx, runtime.graphqlClient, args[0], limit)
	return milestones, milestones.Milestones, err
}

func projectMilestonePageWithItems(
	page client.ProjectMilestoneList,
	milestones []client.ProjectMilestoneSummary,
) client.ProjectMilestoneList {
	page.Milestones = milestones
	return page
}
