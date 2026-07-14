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

			request := client.ProjectLabelAssociationRequest{ProjectID: args[0], LabelID: args[1]}
			project, err := client.AddProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeProject(command, options, project)
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

			request := client.ProjectLabelAssociationRequest{ProjectID: args[0], LabelID: args[1]}
			project, err := client.RemoveProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeProject(command, options, project)
		},
	})
}
