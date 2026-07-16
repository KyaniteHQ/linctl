package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addLabelCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.LabelCreateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.LabelSummary]{
		Use:   "create",
		Short: "Create a label in the pinned team, or organization-wide with --org-wide",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "label name")
			command.Flags().StringVar(&request.Color, "color", "", "label color")
			command.Flags().StringVar(&request.Description, "description", "", "label description")
			command.Flags().StringVar(&request.ParentID, "parent", "", "parent label id")
			command.Flags().BoolVar(
				&request.OrgWide, "org-wide", false,
				"create an organization-wide label instead of a team-scoped label",
			)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.LabelSummary, error) {
			return client.CreateLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeLabel,
	})
}

func addLabelUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.LabelUpdateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.LabelSummary]{
		Use:   "update LABEL_ID",
		Short: "Update a label after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "new label name")
			command.Flags().StringVar(&request.Color, "color", "", "new label color")
			command.Flags().StringVar(&request.Description, "description", "", "new label description")
			command.Flags().BoolVar(
				&request.OrgWide, "org-wide", false,
				"act on an organization-wide label instead of a team-scoped label",
			)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.LabelSummary, error) {
			request.ID = args[0]

			return client.UpdateLabel(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeLabel,
	})
}

func addLabelRetireCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.LabelSummary]{
		Use:   "retire LABEL_ID",
		Short: "Retire a label after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(
				&orgWide, "org-wide", false,
				"act on an organization-wide label instead of a team-scoped label",
			)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.LabelSummary, error) {
			return client.RetireLabel(ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide)
		},
		Write: writeLabel,
	})
}

func addLabelRestoreCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	orgWide := false
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.LabelSummary]{
		Use:   "restore LABEL_ID",
		Short: "Restore a retired label after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(
				&orgWide, "org-wide", false,
				"act on an organization-wide label instead of a team-scoped label",
			)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.LabelSummary, error) {
			return client.RestoreLabel(ctx, runtime.graphqlClient, runtime.config.Target, args[0], orgWide)
		},
		Write: writeLabel,
	})
}
