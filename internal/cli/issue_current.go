package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/gitctx"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueCurrentCommands(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyLocal, &cobra.Command{
		Use:   "id",
		Short: "Show the identifier of the Current Issue",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			identifier, err := gitctx.CurrentIssueIdentifierForTeam(ctx, ".", pinnedTeamKeyHint(ctx, options))
			if err != nil {
				return err
			}

			return writeScalar(command, options, "identifier", identifier)
		},
	})
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
		Use:   "title [ISSUE_ID]",
		Short: "Show the title of an issue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			issue, err := resolveIssueArgument(ctx, options, args)
			if err != nil {
				return err
			}

			return writeScalar(command, options, "title", issue.Title)
		},
	})
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
		Use:   "url [ISSUE_ID]",
		Short: "Show the URL of an issue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			issue, err := resolveIssueArgument(ctx, options, args)
			if err != nil {
				return err
			}

			return writeScalar(command, options, "url", issue.URL)
		},
	})
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
		Use:   "branch ISSUE_ID",
		Short: "Show the branch name of an issue",
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
	issueID, err := issueArgumentOrCurrent(ctx, options, args)
	if err != nil {
		return client.IssueSummary{}, err
	}
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return client.IssueSummary{}, err
	}

	return client.GetIssueByID(ctx, runtime.graphqlClient, issueID)
}

func issueArgumentOrCurrent(ctx context.Context, options *rootOptions, args []string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	return gitctx.CurrentIssueIdentifierForTeam(ctx, ".", pinnedTeamKeyHint(ctx, options))
}
