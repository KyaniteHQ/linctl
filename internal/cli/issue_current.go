package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/gitctx"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueCurrentCommands(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "id",
		Short: "Print the Current Issue identifier",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			identifier, err := gitctx.CurrentIssueIdentifier(ctx, ".")
			if err != nil {
				return err
			}

			return writeScalar(command, options, "identifier", identifier)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "title [ISSUE_ID]",
		Short: "Print an issue title",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			issue, err := resolveIssueArgument(ctx, options, args)
			if err != nil {
				return err
			}

			return writeScalar(command, options, "title", issue.Title)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "url [ISSUE_ID]",
		Short: "Print an issue URL",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			issue, err := resolveIssueArgument(ctx, options, args)
			if err != nil {
				return err
			}

			return writeScalar(command, options, "url", issue.URL)
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "branch ISSUE_ID",
		Short: "Print the issue branch name",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			issue, err := resolveIssueArgument(ctx, options, args)
			if err != nil {
				return err
			}

			return writeScalar(command, options, "branch_name", issue.BranchName)
		},
	})
}

func writeScalar(command *cobra.Command, options *rootOptions, key string, value string) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, map[string]string{key: value})
	}

	return render.WriteLine(command.OutOrStdout(), "%s", value)
}

func resolveIssueArgument(ctx context.Context, options *rootOptions, args []string) (client.IssueSummary, error) {
	issueID, err := issueArgumentOrCurrent(ctx, args)
	if err != nil {
		return client.IssueSummary{}, err
	}
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return client.IssueSummary{}, err
	}

	return resolveIssueArgumentWithReader(ctx, issueAdapterFor(runtime), issueID)
}

func resolveIssueArgumentWithReader(
	ctx context.Context,
	reader currentIssueReader,
	issueID string,
) (client.IssueSummary, error) {
	issue, err := reader.GetIssueByID(ctx, issueID)
	if err != nil {
		return client.IssueSummary{}, err
	}

	return issue, nil
}

func issueArgumentOrCurrent(ctx context.Context, args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	return gitctx.CurrentIssueIdentifier(ctx, ".")
}
