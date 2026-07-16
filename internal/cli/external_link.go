package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addExternalLinkCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	parentCommand := newGroupCommand("external-link", "Read Linear external links")

	addReadGetCommand(ctx, parentCommand, options, readGetSpec[client.EntityExternalLinkSummary]{
		Use:   "get EXTERNAL_LINK_ID",
		Short: "Get one external link by id",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.EntityExternalLinkSummary, error) {
			return client.GetEntityExternalLinkByID(ctx, runtime.graphqlClient, id)
		},
		Write: writeEntityExternalLink,
	})
	root.AddCommand(parentCommand)
}
