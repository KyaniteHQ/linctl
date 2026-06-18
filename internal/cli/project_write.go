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
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			project, err := client.CreateProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeProject(command, options, project)
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
			return runProjectUpdateCommand(ctx, command, options, request)
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
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			project, err := client.ArchiveProject(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			if err != nil {
				return err
			}

			return writeProject(command, options, project)
		},
	})
}

func runProjectUpdateCommand(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	request client.ProjectUpdateRequest,
) error {
	runtime, err := newCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	project, err := client.UpdateProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
	if err != nil {
		return err
	}

	return writeProject(command, options, project)
}
