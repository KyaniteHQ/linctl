package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectInitiativeLinkCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var sortOrder float64
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.InitiativeToProjectSummary]{
		Use:   "initiative-link PROJECT_ID INITIATIVE_ID",
		Short: "Add a project to an Initiative after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Configure: func(command *cobra.Command) {
			command.Flags().Float64Var(
				&sortOrder, "sort-order", 0, "position of the project within the Initiative; Linear picks one when unset",
			)
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.InitiativeToProjectSummary, error) {
			request := client.ProjectInitiativeLinkRequest{ProjectID: args[0], InitiativeID: args[1]}
			if command.Flags().Changed("sort-order") {
				request.SortOrder = &sortOrder
			}

			return client.LinkProjectInitiative(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeInitiativeToProject,
	})
}

func addProjectInitiativeUnlinkCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.InitiativeToProjectSummary]{
		Use:   "initiative-unlink PROJECT_ID INITIATIVE_ID",
		Short: "Remove a project from an Initiative after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.InitiativeToProjectSummary, error) {
			request := client.ProjectInitiativeLinkRequest{ProjectID: args[0], InitiativeID: args[1]}

			return client.UnlinkProjectInitiative(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeInitiativeToProject,
	})
}
