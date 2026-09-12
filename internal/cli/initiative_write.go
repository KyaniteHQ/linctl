package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

const orgWideInitiativeHelp = "required: an Initiative is organization-owned and sits above every team, so a " +
	"create cannot land inside the pinned team; confirms this write adds an initiative to the organization"

func addInitiativeCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.InitiativeCreateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.InitiativeSummary]{
		Use:   "create",
		Short: "Create an initiative with --org-wide, because an initiative is always outside the pinned team",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "initiative name")
			command.Flags().StringVar(&request.Description, "description", "", "initiative description")
			command.Flags().StringVar(&request.Status, "status", "",
				"initiative status: Proposed, Planned, Active, Completed, or Canceled; Linear defaults it when unset")
			command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideInitiativeHelp)
			registerFlagCompletion(command, "status",
				cobra.FixedCompletions(client.InitiativeStatuses, cobra.ShellCompDirectiveNoFileComp))
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.InitiativeSummary, error) {
			return client.CreateInitiative(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeInitiative,
	})
}
