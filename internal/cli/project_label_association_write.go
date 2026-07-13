//nolint:dupl // Label-association command glue is intentionally uniform across issue and project targets.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectAddLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyWrite, &cobra.Command{
		Use:   "add-label PROJECT_ID LABEL_ID",
		Short: "Attach a ProjectLabel to a project after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runProjectAddLabel(
				ctx, command, options, commandAdapterFor(runtime),
				client.ProjectLabelAssociationRequest{ProjectID: args[0], LabelID: args[1]},
			)
		},
	})
}

func addProjectRemoveLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyWrite, &cobra.Command{
		Use:   "remove-label PROJECT_ID LABEL_ID",
		Short: "Detach a ProjectLabel from a project after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runProjectRemoveLabel(
				ctx, command, options, commandAdapterFor(runtime),
				client.ProjectLabelAssociationRequest{ProjectID: args[0], LabelID: args[1]},
			)
		},
	})
}

func runProjectAddLabel(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	adder projectLabelAdder,
	request client.ProjectLabelAssociationRequest,
) error {
	project, err := adder.AddProjectLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeProject(command, options, project)
}

func runProjectRemoveLabel(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	remover projectLabelRemover,
	request client.ProjectLabelAssociationRequest,
) error {
	project, err := remover.RemoveProjectLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeProject(command, options, project)
}
