package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectCreateRequest{}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a project in the pinned team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runProjectWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.ProjectSummary,
				error,
			) {
				return client.CreateProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
			})
		},
	}
	command.Flags().StringVar(&request.Name, "name", "", "project name")
	command.Flags().StringVar(&request.Description, "description", "", "project description")
	root.AddCommand(command)
}

func addProjectUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectUpdateRequest{}
	command := &cobra.Command{
		Use:   "update PROJECT_ID",
		Short: "Update a project after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			request.ID = args[0]
			return runProjectWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.ProjectSummary,
				error,
			) {
				return client.UpdateProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
			})
		},
	}
	command.Flags().StringVar(&request.Name, "name", "", "new project name")
	command.Flags().StringVar(&request.Description, "description", "", "new project description")
	root.AddCommand(command)
}

func addProjectArchiveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "archive PROJECT_ID",
		Short: "Archive a project after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runProjectWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.ProjectSummary,
				error,
			) {
				return client.ArchiveProject(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			})
		},
	})
}

func runProjectWriteCommand(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	write func(commandRuntime) (client.ProjectSummary, error),
) error {
	return runGuardedWrite(ctx, command, options, write, writeProject)
}
