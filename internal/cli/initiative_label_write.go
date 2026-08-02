package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

const orgWideInitiativeLabelHelp = "required: initiative labels have no team scope; confirms this write " +
	"affects every initiative in the organization"

func addInitiativeLabelRetireCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.InitiativeLabelSummary]{
		Use:   "retire INITIATIVE_LABEL_ID",
		Short: "Retire an initiative label with --org-wide, which changes every initiative in the organization",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideInitiativeLabelHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.InitiativeLabelSummary, error) {
			return client.RetireInitiativeLabel(ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide)
		},
		Write: writeInitiativeLabel,
	})
}

func addInitiativeLabelRestoreCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.InitiativeLabelSummary]{
		Use:   "restore INITIATIVE_LABEL_ID",
		Short: "Restore a retired initiative label with --org-wide, which changes every initiative",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(&orgWide, "org-wide", false, orgWideInitiativeLabelHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.InitiativeLabelSummary, error) {
			return client.RestoreInitiativeLabel(ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide)
		},
		Write: writeInitiativeLabel,
	})
}
