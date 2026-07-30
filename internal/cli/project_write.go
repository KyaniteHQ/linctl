package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectCreateRequest{}
	var contentFile string
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectSummary]{
		Use:   "create",
		Short: "Create a project in the pinned team",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "project name")
			command.Flags().StringVar(&request.Description, "description", "", "project description")
			command.Flags().StringVar(&request.Content, "content", "", "project content as markdown")
			command.Flags().StringVar(&contentFile, "content-file", "", "read project content markdown from a file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, _ []string,
		) (client.ProjectSummary, error) {
			if err := resolveFileFlag(command, &request.Content, contentFile, "content"); err != nil {
				return client.ProjectSummary{}, err
			}

			return client.CreateProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProject,
	})
}

func addProjectUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectUpdateRequest{}
	var contentFile string
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectSummary]{
		Use:   "update PROJECT_ID",
		Short: "Update a project after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "new project name")
			command.Flags().StringVar(&request.Description, "description", "", "new project description")
			command.Flags().StringVar(&request.Content, "content", "", "new project content as markdown")
			command.Flags().StringVar(&contentFile, "content-file", "", "read new project content markdown from a file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectSummary, error) {
			request.ID = args[0]
			if err := resolveFileFlag(command, &request.Content, contentFile, "content"); err != nil {
				return client.ProjectSummary{}, err
			}

			return client.UpdateProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProject,
	})
}

func addProjectArchiveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectSummary]{
		Use:   "archive PROJECT_ID",
		Short: "Archive a project after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectSummary, error) {
			return client.ArchiveProject(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeProject,
	})
}

func addProjectAddTeamCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectAddTeamRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectSummary]{
		Use:   "add-team PROJECT_ID",
		Short: "Attach a team to a project without dropping existing teams",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.TeamKey, "to-team", "", "team key to attach")
			command.Flags().StringVar(&request.TeamID, "to-team-id", "", "team id to attach")
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectSummary, error) {
			request.ProjectID = args[0]

			return client.AddProjectTeam(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProject,
	})
}
