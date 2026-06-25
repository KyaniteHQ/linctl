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
		},
	)
	addCommentBotActorCommand(ctx, commentCommand, options)
	addCommentChildrenCommand(ctx, commentCommand, options)
	addCommentCreatedIssuesCommand(ctx, commentCommand, options)
	addCommentUpdateCommand(ctx, commentCommand, options)
	addCommentDeleteCommand(ctx, commentCommand, options)
}

func addCommentUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.CommentUpdateRequest{}
	bodyFile := ""
	command := &cobra.Command{
		Use:   "update COMMENT_ID",
		Short: "Edit a comment after pinned-target comparison",
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
			comment, err := client.UpdateComment(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeComment(command, options, comment)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "new comment body as markdown; use - to read stdin")
	command.Flags().StringVar(&bodyFile, "body-file", "", "read new comment body from file")
	root.AddCommand(command)
}

func addCommentDeleteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "delete COMMENT_ID",
		Short: "Delete a comment after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			id, err := client.DeleteComment(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			if err != nil {
				return err
			}

			return writeDeletion(command, options, id)
		},
	})
}

func writeComment(command *cobra.Command, options *rootOptions, comment client.CommentSummary) error {
	return writeItem(command, options, comment, comment.ID,
		func(command *cobra.Command, _ *rootOptions, comment client.CommentSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s",
				comment.ID,
				emptyDash(comment.DisplayName),
				comment.Body,
			)
		})
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
	return writeItem(command, options, actor, actor.CommentID,
		func(command *cobra.Command, _ *rootOptions, actor client.CommentBotActor) error {
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
		})
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
