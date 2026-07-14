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
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a project label; requires --org-wide (affects every team and project in the organization)",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			label, err := client.CreateProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeProjectLabel(command, options, label)
		},
	}
	command.Flags().StringVar(&request.Name, "name", "", "project label name")
	command.Flags().StringVar(&request.Color, "color", "", "project label color")
	command.Flags().StringVar(&request.Description, "description", "", "project label description")
	command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideProjectLabelHelp)
	root.AddCommand(command)
}

func addProjectLabelUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectLabelUpdateRequest{}
	command := &cobra.Command{
		Use:   "update PROJECT_LABEL_ID",
		Short: "Update a project label; requires --org-wide (affects every team and project in the organization)",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]

			label, err := client.UpdateProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeProjectLabel(command, options, label)
		},
	}
	command.Flags().StringVar(&request.Name, "name", "", "new project label name")
	command.Flags().StringVar(&request.Color, "color", "", "new project label color")
	command.Flags().StringVar(&request.Description, "description", "", "new project label description")
	command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideProjectLabelHelp)
	root.AddCommand(command)
}

func addProjectLabelRetireCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	command := &cobra.Command{
		Use:   "retire PROJECT_LABEL_ID",
		Short: "Retire a project label; requires --org-wide (affects every team and project in the organization)",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			label, err := client.RetireProjectLabel(
				ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide,
			)
			if err != nil {
				return err
			}

			return writeProjectLabel(command, options, label)
		},
	}
	command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideProjectLabelHelp)
	root.AddCommand(command)
}

func addProjectLabelRestoreCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	command := &cobra.Command{
		Use:   "restore PROJECT_LABEL_ID",
		Short: "Restore a retired project label; requires --org-wide (affects every team and project)",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			label, err := client.RestoreProjectLabel(
				ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide,
			)
			if err != nil {
				return err
			}

			return writeProjectLabel(command, options, label)
		},
	}
	command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideProjectLabelHelp)
	root.AddCommand(command)
}
