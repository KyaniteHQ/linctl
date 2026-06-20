package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addExternalLinkCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	parentCommand := &cobra.Command{
		Use:   "external-link",
		Short: "Read Linear external links",
	}

	getCommand := &cobra.Command{
		Use:   "get EXTERNAL_LINK_ID",
		Short: "Get one external link by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			link, err := client.GetEntityExternalLinkByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeEntityExternalLink(command, options, link)
		},
	}

	parentCommand.AddCommand(getCommand)
	root.AddCommand(parentCommand)
}
