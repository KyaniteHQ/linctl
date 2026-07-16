package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectMilestoneDeleteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[string]{
		Use:   "delete PROJECT_MILESTONE_ID",
		Short: "Hard delete a ProjectMilestone after pinned-target comparison; cannot be undone via linctl",
		Args:  cobra.ExactArgs(1),
		Run: func(ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string) (string, error) {
			return client.DeleteProjectMilestone(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeProjectMilestoneDeletion,
	})
}

// writeProjectMilestoneDeletion overrides the human deletion line with an
// explicit irreversibility warning: ProjectMilestone delete is linctl's one
// approved hard delete, and there is no restore path via linctl.
func writeProjectMilestoneDeletion(command *cobra.Command, options *rootOptions, id string) error {
	return writeDeletionMessage(
		command, options, id,
		"hard deleted ProjectMilestone "+id+": cannot be undone via linctl",
	)
}

func addProjectMilestoneCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := projectMilestoneWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectMilestoneSummary]{
		Use:       "create PROJECT_ID",
		Short:     "Create a ProjectMilestone in a pinned project",
		Args:      cobra.ExactArgs(1),
		Configure: bindProjectMilestoneWriteFlags(&flags, ""),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectMilestoneSummary, error) {
			request := client.ProjectMilestoneCreateRequest{
				ProjectID:   args[0],
				Name:        flags.Name,
				Description: flags.Description,
				TargetDate:  flags.TargetDate,
			}

			return client.CreateProjectMilestone(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProjectMilestone,
	})
}

func addProjectMilestoneUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := projectMilestoneWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectMilestoneSummary]{
		Use:       "update PROJECT_MILESTONE_ID",
		Short:     "Update a ProjectMilestone after pinned-target comparison",
		Args:      cobra.ExactArgs(1),
		Configure: bindProjectMilestoneWriteFlags(&flags, "new "),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectMilestoneSummary, error) {
			request := client.ProjectMilestoneUpdateRequest{
				ID:          args[0],
				Name:        flags.Name,
				Description: flags.Description,
				TargetDate:  flags.TargetDate,
			}

			return client.UpdateProjectMilestone(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProjectMilestone,
	})
}

type projectMilestoneWriteFlags struct {
	Name        string
	Description string
	TargetDate  string
}

func bindProjectMilestoneWriteFlags(flags *projectMilestoneWriteFlags, helpPrefix string) func(*cobra.Command) {
	return func(command *cobra.Command) {
		command.Flags().StringVar(&flags.Name, "name", "", helpPrefix+"ProjectMilestone name")
		command.Flags().StringVar(&flags.Description, "description", "", helpPrefix+"ProjectMilestone description")
		command.Flags().StringVar(&flags.TargetDate, "target-date", "", helpPrefix+"ProjectMilestone target date")
	}
}
