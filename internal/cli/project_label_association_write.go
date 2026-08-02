//nolint:dupl // Label-association command glue is intentionally uniform across issue and project targets.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectAddLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectSummary]{
		Use:   "add-label PROJECT_ID LABEL_ID",
		Short: "Attach a ProjectLabel to a project after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectSummary, error) {
			request := client.ProjectLabelAssociationRequest{ProjectID: args[0], LabelID: args[1]}

			return client.AddProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProject,
	})
}

func addProjectRemoveLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectSummary]{
		Use:   "remove-label PROJECT_ID LABEL_ID",
		Short: "Remove a ProjectLabel from a project after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectSummary, error) {
			request := client.ProjectLabelAssociationRequest{ProjectID: args[0], LabelID: args[1]}

			return client.RemoveProjectLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProject,
	})
}
