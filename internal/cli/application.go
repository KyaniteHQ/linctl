package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addApplicationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "application",
		Short: "Read Linear OAuth application metadata",
	}
	command.AddCommand(&cobra.Command{
		Use:   "info CLIENT_ID",
		Short: "Get public OAuth application metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			application, err := client.GetApplicationInfo(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeApplicationInfo(command, options, application)
		},
	})
	root.AddCommand(command)
}

func writeApplicationInfo(command *cobra.Command, options *rootOptions, application client.ApplicationInfo) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, application)
	}
	if wrote, err := writeIDOnly(command, options, application.ID); wrote || err != nil {
		return err
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s by %s",
		application.ID,
		application.Name,
		emptyDash(application.Developer),
	)
}
