package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCommentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CommentList, client.CommentSummary]{
		Use:           "comment",
		Short:         "Read Linear comments",
		ListShort:     "List visible comments",
		LimitHelp:     "maximum comments to return",
		GetUse:        "get COMMENT_ID",
		GetShort:      "Get one comment by id",
		LoadList:      loadCommentList,
		PageWithItems: commentPageWithItems,
		LoadGet:       loadComment,
		WriteItem:     writeComment,
	})
}

func writeComment(command *cobra.Command, options *rootOptions, comment client.CommentSummary) error {
	if wrote, err := writeIDOnly(command, options, comment.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, comment)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s %s", comment.ID, emptyDash(comment.DisplayName), comment.Body)
}

func loadCommentList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CommentList, []client.CommentSummary, error) {
	comments, err := client.ListComments(ctx, runtime.graphqlClient, limit)
	return comments, comments.Comments, err
}

func loadComment(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CommentSummary, error) {
	return client.GetCommentByID(ctx, runtime.graphqlClient, id)
}

func commentPageWithItems(page client.CommentList, comments []client.CommentSummary) client.CommentList {
	page.Comments = comments
	return page
}
