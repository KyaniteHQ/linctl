package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addApplicationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := newGroupCommand("application", "Read Linear OAuth application metadata")
	addReadGetCommand(ctx, command, options, readGetSpec[client.ApplicationInfo]{
		Use:   "info CLIENT_ID",
		Short: "Get public OAuth application metadata",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.ApplicationInfo, error) {
			return client.GetApplicationInfo(ctx, runtime.graphqlClient, id)
		},
		Write: writeApplicationInfo,
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
