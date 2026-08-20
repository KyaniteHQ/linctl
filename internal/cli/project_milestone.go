package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectMilestoneCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectMilestoneCommand := newGroupCommand("project-milestone", "Read and write Linear project milestones")
	addProjectMilestoneAllCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneListCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneGetCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneIssuesCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneCreateCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneUpdateCommand(ctx, projectMilestoneCommand, options)
	addProjectMilestoneDeleteCommand(ctx, projectMilestoneCommand, options)
	root.AddCommand(projectMilestoneCommand)
}

func addProjectMilestoneAllCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectMilestoneList, client.ProjectMilestoneSummary]{
		Use:       "all",
		Short:     "List visible ProjectMilestones across the organization",
		LimitHelp: "project milestones",
		Args:      cobra.NoArgs,
		Load:      loadAllProjectMilestones,
		WriteItem: writeProjectMilestone,
	})
}

func addProjectMilestoneIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"issues PROJECT_MILESTONE_ID",
		"List issues for one ProjectMilestone",
		"issues",
		client.ListProjectMilestoneIssues,
		writeIssue,
	)
}

func addProjectMilestoneListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"list PROJECT_ID",
		"List ProjectMilestones for one project",
		"project milestones",
		client.ListProjectMilestones,
		writeProjectMilestone,
	)
}

func addProjectMilestoneGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.ProjectMilestoneSummary]{
		Use:   "get PROJECT_MILESTONE_ID",
		Short: "Get one ProjectMilestone by id",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.ProjectMilestoneSummary, error) {
			return client.GetProjectMilestoneByID(ctx, runtime.graphqlClient, id)
		},
		Write: writeProjectMilestone,
	})
}

func writeProjectMilestone(
	command *cobra.Command,
	options *rootOptions,
	milestone client.ProjectMilestoneSummary,
) error {
	return writeItem(command, options, milestone, milestone.ID,
		func(command *cobra.Command, options *rootOptions, milestone client.ProjectMilestoneSummary) error {
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
		})
}

func loadAllProjectMilestones(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectMilestoneList, error) {
	milestones, err := client.ListAllProjectMilestones(ctx, runtime.graphqlClient, limit)
	return milestones, err
}
