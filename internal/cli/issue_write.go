package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCreateRequest{}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create an issue in the pinned target",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.CreateIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "issue title")
	command.Flags().StringVar(&request.Description, "description", "", "issue description")
	root.AddCommand(command)
}

func addIssueUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueUpdateRequest{}
	command := &cobra.Command{
		Use:   "update ISSUE_ID",
		Short: "Update an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			issue, err := client.UpdateIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "new issue title")
	command.Flags().StringVar(&request.Description, "description", "", "new issue description")
	root.AddCommand(command)
}

func addIssueCommentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	command := &cobra.Command{
		Use:   "comment ISSUE_ID",
		Short: "Comment on an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			comment, err := client.CommentOnIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}
			if options.json {
				return render.WriteJSON(command.OutOrStdout(), comment)
			}

			return render.WriteLine(command.OutOrStdout(), "comment %s on %s", comment.ID, comment.Issue.Identifier)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "comment body")
	root.AddCommand(command)
}

func addIssueCloseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "close ISSUE_ID",
		Short: "Move an issue to the completed workflow state",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := newCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.CloseIssue(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}

func writeIssue(command *cobra.Command, options *rootOptions, issue client.IssueSummary) error {
	if options.json {
		return render.WriteJSON(command.OutOrStdout(), issue)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", issue.Identifier, issue.Title, issue.State)
}
