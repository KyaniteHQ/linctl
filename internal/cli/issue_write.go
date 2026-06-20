package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCreateRequest{}
	descriptionFile := ""
	command := &cobra.Command{
		Use:   "create",
		Short: "Create an issue in the pinned target",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			if err := resolveFileFlag(&request.Description, descriptionFile, "description"); err != nil {
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
	command.Flags().StringVar(&descriptionFile, "description-file", "", "read issue description from file")
	root.AddCommand(command)
}

func addIssueUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueUpdateRequest{}
	descriptionFile := ""
	appendFile := ""
	command := &cobra.Command{
		Use:   "update ISSUE_ID",
		Short: "Update an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			if err := resolveFileFlag(&request.Description, descriptionFile, "description"); err != nil {
				return err
			}
			if err := resolveFileFlag(&request.Append, appendFile, "append"); err != nil {
				return err
			}
			issue, err := client.UpdateIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "new issue title")
	command.Flags().StringVar(&request.Description, "description", "", "new issue description")
	command.Flags().StringVar(&descriptionFile, "description-file", "", "read new issue description from file")
	command.Flags().StringVar(&request.Append, "append", "", "text to append to the issue description")
	command.Flags().StringVar(&appendFile, "append-file", "", "read text to append from file")
	root.AddCommand(command)
}

func addIssueStartCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "start ISSUE_ID",
		Short: "Assign and start an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.StartIssue(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}

func addIssueCommentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	bodyFile := ""
	command := &cobra.Command{
		Use:   "comment ISSUE_ID",
		Short: "Comment on an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			if err := resolveBodyFlag(command, &request.Body); err != nil {
				return err
			}
			if err := resolveFileFlag(&request.Body, bodyFile, "body"); err != nil {
				return err
			}
			comment, err := client.CommentOnIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, comment)
			}

			return render.WriteLine(command.OutOrStdout(), "comment %s on %s", comment.ID, comment.Issue.Identifier)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "comment body")
	command.Flags().StringVar(&bodyFile, "body-file", "", "read comment body from file")
	root.AddCommand(command)
}

func addIssueReplyCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	bodyFile := ""
	command := &cobra.Command{
		Use:   "reply ISSUE_ID COMMENT_ID",
		Short: "Reply to an issue comment after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			request.ParentID = args[1]
			if err := resolveBodyFlag(command, &request.Body); err != nil {
				return err
			}
			if err := resolveFileFlag(&request.Body, bodyFile, "body"); err != nil {
				return err
			}
			comment, err := client.CommentOnIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, comment)
			}

			return render.WriteLine(command.OutOrStdout(), "comment %s on %s", comment.ID, comment.Issue.Identifier)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "reply body")
	command.Flags().StringVar(&bodyFile, "body-file", "", "read reply body from file")
	root.AddCommand(command)
}

func resolveBodyFlag(command *cobra.Command, body *string) error {
	if *body != "-" {
		return nil
	}
	data, err := io.ReadAll(command.InOrStdin())
	if err != nil {
		return fmt.Errorf("read body from stdin: %w", err)
	}
	*body = string(data)

	return nil
}

func resolveFileFlag(value *string, path string, label string) error {
	if path == "" {
		return nil
	}
	if *value != "" {
		return fmt.Errorf("%s and %s-file are mutually exclusive", label, label)
	}

	//nolint:gosec // The path is an explicit CLI input for reading issue text.
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s from file %s: %w", label, path, err)
	}
	*value = string(data)

	return nil
}

func addIssueCloseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "close ISSUE_ID",
		Short: "Move an issue to the completed workflow state",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
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
	if wrote, err := writeIDOnly(command, options, issue.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, issue)
	}

	format, err := normalizedHumanFormat(options)
	if err != nil {
		return err
	}
	if format == "minimal" {
		return render.WriteLine(command.OutOrStdout(), "%s", issue.Identifier)
	}
	if format == "full" {
		return render.WriteLine(
			command.OutOrStdout(),
			"%s %s [%s] project=%s url=%s",
			issue.Identifier,
			issue.Title,
			issue.State,
			emptyDash(issue.Project),
			issue.URL,
		)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", issue.Identifier, issue.Title, issue.State)
}

func emptyDash(value string) string {
	if value == "" {
		return "-"
	}

	return value
}
