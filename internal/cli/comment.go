package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCommentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	commentCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.CommentList, client.CommentSummary]{
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
	addCommentBotActorCommand(ctx, commentCommand, options)
	addCommentChildrenCommand(ctx, commentCommand, options)
	addCommentCreatedIssuesCommand(ctx, commentCommand, options)
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

func addCommentBotActorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "bot-actor COMMENT_ID",
		Short: "Show comment bot actor metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			actor, err := client.GetCommentBotActor(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCommentBotActor(command, options, actor)
		},
	})
}

func addCommentChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "children COMMENT_ID",
		Short: "List child comments without body content",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadCommentChildren,
				commentChildPageWithItems,
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum child comments to return")
	root.AddCommand(command)
}

func addCommentCreatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "created-issues COMMENT_ID",
		Short: "List issues created from a comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadCommentCreatedIssues,
				issuePageWithItems,
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func writeCommentBotActor(command *cobra.Command, options *rootOptions, actor client.CommentBotActor) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, actor)
	}
	if actor.Bot == nil {
		return render.WriteLine(command.OutOrStdout(), "%s bot -", actor.CommentID)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s bot %s %s [%s]",
		actor.CommentID,
		emptyDash(actor.Bot.ID),
		emptyDash(actor.Bot.Name),
		actor.Bot.Type,
	)
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

func loadCommentChildren(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.CommentChildList, []client.CommentMetadataSummary, error) {
	comments, err := client.ListCommentChildren(ctx, runtime.graphqlClient, args[0], limit)
	return comments, comments.Comments, err
}

func loadCommentCreatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListCommentCreatedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func commentPageWithItems(page client.CommentList, comments []client.CommentSummary) client.CommentList {
	page.Comments = comments
	return page
}

func commentChildPageWithItems(
	page client.CommentChildList,
	comments []client.CommentMetadataSummary,
) client.CommentChildList {
	page.Comments = comments
	return page
}
