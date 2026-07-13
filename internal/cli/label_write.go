package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addLabelCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.LabelCreateRequest{}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a label in the pinned team, or organization-wide with --org-wide",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runLabelCreate(ctx, command, options, commandAdapterFor(runtime), request)
		},
	}
	command.Flags().StringVar(&request.Name, "name", "", "label name")
	command.Flags().StringVar(&request.Color, "color", "", "label color")
	command.Flags().StringVar(&request.Description, "description", "", "label description")
	command.Flags().StringVar(&request.ParentID, "parent", "", "parent label id")
	command.Flags().BoolVar(
		&request.OrgWide, "org-wide", false,
		"create an organization-wide label instead of a team-scoped label",
	)
	root.AddCommand(command)
}

func addLabelUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.LabelUpdateRequest{}
	command := &cobra.Command{
		Use:   "update LABEL_ID",
		Short: "Update a label after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]

			return runLabelUpdate(ctx, command, options, commandAdapterFor(runtime), request)
		},
	}
	command.Flags().StringVar(&request.Name, "name", "", "new label name")
	command.Flags().StringVar(&request.Color, "color", "", "new label color")
	command.Flags().StringVar(&request.Description, "description", "", "new label description")
	command.Flags().BoolVar(
		&request.OrgWide, "org-wide", false,
		"act on an organization-wide label instead of a team-scoped label",
	)
	root.AddCommand(command)
}

func addLabelRetireCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	command := &cobra.Command{
		Use:   "retire LABEL_ID",
		Short: "Retire a label after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runLabelRetire(ctx, command, options, commandAdapterFor(runtime), args[0], orgWide)
		},
	}
	command.Flags().BoolVar(
		&orgWide, "org-wide", false,
		"act on an organization-wide label instead of a team-scoped label",
	)
	root.AddCommand(command)
}

func addLabelRestoreCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	command := &cobra.Command{
		Use:   "restore LABEL_ID",
		Short: "Restore a retired label after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runLabelRestore(ctx, command, options, commandAdapterFor(runtime), args[0], orgWide)
		},
	}
	command.Flags().BoolVar(
		&orgWide, "org-wide", false,
		"act on an organization-wide label instead of a team-scoped label",
	)
	root.AddCommand(command)
}

func runLabelCreate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator labelCreator,
	request client.LabelCreateRequest,
) error {
	label, err := creator.CreateLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeLabel(command, options, label)
}

func runLabelUpdate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	updater labelUpdater,
	request client.LabelUpdateRequest,
) error {
	label, err := updater.UpdateLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeLabel(command, options, label)
}

func runLabelRetire(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	retirer labelRetirer,
	id string,
	orgWide bool,
) error {
	label, err := retirer.RetireLabel(ctx, id, orgWide)
	if err != nil {
		return err
	}

	return writeLabel(command, options, label)
}

func runLabelRestore(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	restorer labelRestorer,
	id string,
	orgWide bool,
) error {
	label, err := restorer.RestoreLabel(ctx, id, orgWide)
	if err != nil {
		return err
	}

	return writeLabel(command, options, label)
}
