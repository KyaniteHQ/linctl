package cli

import (
	"context"

	"github.com/spf13/cobra"

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
			return runCurrentIssueRead(ctx, command, options, issueAdapterFor(runtime), issueID)
		},
	})
}

func runCurrentIssueRead(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	reader currentIssueReader,
	issueID string,
) error {
	issue, err := reader.GetIssueByID(ctx, issueID)
	if err != nil {
		return err
	}

	return writeIssue(command, options, issue)
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
			issue, err := issueAdapterFor(runtime).CloseIssue(ctx, issueID)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}
