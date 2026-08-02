package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx, root, options, "comments ISSUE_ID", "List the comments of an issue", "comments",
		client.ListIssueComments, writeIssueCommentSummary,
	)
}

func writeIssueCommentSummary(
	command *cobra.Command,
	options *rootOptions,
	comment client.IssueCommentSummary,
) error {
	return writeItemLine(
		command, options, comment, comment.ID,
		"%s %s %s",
		comment.ID,
		emptyDash(comment.DisplayName),
		comment.Body,
	)
}
