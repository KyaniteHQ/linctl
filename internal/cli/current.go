package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/gitctx"
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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.GetIssueByID(ctx, runtime.graphqlClient, issueID)
			if err != nil {
				return err
			}
			return writeIssue(command, options, issue)
		},
	})
}

func addDoneCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "done",
		Short: "Close the current checkout issue",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			issueID, err := gitctx.CurrentIssueIdentifier(ctx, ".")
			if err != nil {
				return err
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.CloseIssue(ctx, runtime.graphqlClient, runtime.config.Target, issueID)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}
