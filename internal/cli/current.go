package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/gitctx"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCurrentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "current",
		Short: "Resolve the Linear issue for the current checkout",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			issueID, err := gitctx.CurrentIssueIdentifier(ctx, ".")
			if err != nil {
				return err
			}
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.GetIssueByID(ctx, runtime.graphqlClient, issueID)
			if err != nil {
				return err
			}
			if options.json {
				return render.WriteJSON(command.OutOrStdout(), issue)
			}

			return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", issue.Identifier, issue.Title, issue.State)
		},
	})
}
