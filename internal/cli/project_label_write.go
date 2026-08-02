package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

const orgWideProjectLabelHelp = "required: project labels have no team scope; confirms this write " +
	"affects every team and project in the organization"

func addProjectLabelCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectLabelCreateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectLabelSummary]{
		Use:   "create",
		Short: "Create a project label with --org-wide, which changes every team and project in the organization",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "project label name")
			command.Flags().StringVar(&request.Color, "color", "", "project label color")
			command.Flags().StringVar(&request.Description, "description", "", "project label description")
			command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideProjectLabelHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.ProjectLabelSummary, error) {
			return client.CreateProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProjectLabel,
	})
}

func addProjectLabelUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectLabelUpdateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectLabelSummary]{
		Use:   "update PROJECT_LABEL_ID",
		Short: "Update a project label with --org-wide, which changes every team and project in the organization",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "new project label name")
			command.Flags().StringVar(&request.Color, "color", "", "new project label color")
			command.Flags().StringVar(&request.Description, "description", "", "new project label description")
			command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideProjectLabelHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectLabelSummary, error) {
			request.ID = args[0]

			return client.UpdateProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProjectLabel,
	})
}

func addProjectLabelRetireCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectLabelSummary]{
		Use:   "retire PROJECT_LABEL_ID",
		Short: "Retire a project label with --org-wide, which changes every team and project in the organization",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideProjectLabelHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectLabelSummary, error) {
			return client.RetireProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide)
		},
		Write: writeProjectLabel,
	})
}

func addProjectLabelRestoreCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectLabelSummary]{
		Use:   "restore PROJECT_LABEL_ID",
		Short: "Restore a retired project label with --org-wide, which changes every team and project",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideProjectLabelHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectLabelSummary, error) {
			return client.RestoreProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide)
		},
		Write: writeProjectLabel,
	})
}
