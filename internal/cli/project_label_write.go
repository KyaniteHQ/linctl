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

			return runProjectLabelCreate(ctx, command, options, commandAdapterFor(runtime), request)
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

			return runProjectLabelUpdate(ctx, command, options, commandAdapterFor(runtime), request)
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

			return runProjectLabelRetire(ctx, command, options, commandAdapterFor(runtime), args[0], orgWide)
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

			return runProjectLabelRestore(ctx, command, options, commandAdapterFor(runtime), args[0], orgWide)
		},
	}
	command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideProjectLabelHelp)
	root.AddCommand(command)
}

func runProjectLabelCreate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator projectLabelCreator,
	request client.ProjectLabelCreateRequest,
) error {
	label, err := creator.CreateProjectLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeProjectLabel(command, options, label)
}

func runProjectLabelUpdate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	updater projectLabelUpdater,
	request client.ProjectLabelUpdateRequest,
) error {
	label, err := updater.UpdateProjectLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeProjectLabel(command, options, label)
}

func runProjectLabelRetire(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	retirer projectLabelRetirer,
	id string,
	orgWide bool,
) error {
	label, err := retirer.RetireProjectLabel(ctx, id, orgWide)
	if err != nil {
		return err
	}

	return writeProjectLabel(command, options, label)
}

func runProjectLabelRestore(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	restorer projectLabelRestorer,
	id string,
	orgWide bool,
) error {
	label, err := restorer.RestoreProjectLabel(ctx, id, orgWide)
	if err != nil {
		return err
	}

	return writeProjectLabel(command, options, label)
}
