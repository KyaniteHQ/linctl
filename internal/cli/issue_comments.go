package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "comments ISSUE_ID",
		Short: "List issue comments",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			comments, err := issueAdapterFor(runtime).ListIssueComments(ctx, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(comments.Comments)); err != nil {
				return err
			}
			comments.Comments, err = sortByJSONField(comments.Comments, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, comments)
			}

			return writeIssueComments(command, comments.Comments)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum comments to return")
	root.AddCommand(command)
}

func writeIssueComments(command *cobra.Command, comments []client.IssueCommentSummary) error {
	for _, comment := range comments {
		if err := render.WriteLine(
			command.OutOrStdout(),
			"%s %s %s",
			comment.ID,
			emptyDash(comment.DisplayName),
			comment.Body,
		); err != nil {
			return err
		}
	}

	return nil
}
