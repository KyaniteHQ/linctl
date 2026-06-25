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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runProjectCreate(ctx, command, options, projectAdapterFor(runtime), request)
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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]

			return runProjectUpdate(ctx, command, options, projectAdapterFor(runtime), request)
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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runProjectArchive(ctx, command, options, projectAdapterFor(runtime), args[0])
		},
	})
}

func runProjectCreate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator projectCreator,
	request client.ProjectCreateRequest,
) error {
	project, err := creator.CreateProject(ctx, request)
	if err != nil {
		return err
	}

	return writeProject(command, options, project)
}

func runProjectUpdate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	updater projectUpdater,
	request client.ProjectUpdateRequest,
) error {
	project, err := updater.UpdateProject(ctx, request)
	if err != nil {
		return err
	}

	return writeProject(command, options, project)
}

func runProjectArchive(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	archiver projectArchiver,
	projectID string,
) error {
	project, err := archiver.ArchiveProject(ctx, projectID)
	if err != nil {
		return err
	}

	return writeProject(command, options, project)
}
