package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTargetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "target",
		Short: "Print the resolved Linear target",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			target, err := client.ResolveTarget(ctx, runtime.graphqlClient, runtime.config.Target)
			if err != nil {
				return err
			}
			if options.json {
				return render.WriteJSON(command.OutOrStdout(), target)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"org %s team %s/%s project %s confirmed %t",
				target.Org.ID,
				target.Team.Key,
				target.Team.ID,
				projectID(target.Project),
				target.Confirmed,
			)
		},
	})
}

func addWhoamiCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "whoami",
		Short: "Print the authenticated Linear user",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			target, err := client.ResolveTarget(ctx, runtime.graphqlClient, runtime.config.Target)
			if err != nil {
				return err
			}
			if options.json {
				return render.WriteJSON(command.OutOrStdout(), target.Viewer)
			}

			return render.WriteLine(command.OutOrStdout(), "%s <%s>", target.Viewer.Name, target.Viewer.Email)
		},
	})
}

func projectID(project *client.ResolvedProject) string {
	if project == nil {
		return ""
	}

	return project.ID
}
