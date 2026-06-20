package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addOrganizationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "organization",
		Short: "Read Linear organization metadata",
	}
	command.AddCommand(&cobra.Command{
		Use:   "exists URL_KEY",
		Short: "Check whether a Linear organization URL key exists",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			status, err := client.CheckOrganizationExists(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeOrganizationExists(command, options, status)
		},
	})
	root.AddCommand(command)
}

func writeOrganizationExists(
	command *cobra.Command,
	options *rootOptions,
	status client.OrganizationExistsStatus,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, status)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s exists %t success %t",
		status.URLKey,
		status.Exists,
		status.Success,
	)
}
