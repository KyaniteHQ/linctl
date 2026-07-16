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
			Use:       "comment",
			Short:     "Read Linear comments",
			ListShort: "List visible comments",
			LimitHelp: "maximum comments to return",
			GetUse:    "get COMMENT_ID",
			GetShort:  "Get one comment by id",
			LoadList:  loadCommentList,
			LoadGet:   loadComment,
			WriteItem: writeComment,
		},
	)
	addCommentBotActorCommand(ctx, commentCommand, options)
	addCommentChildrenCommand(ctx, commentCommand, options)
	addCommentCreatedIssuesCommand(ctx, commentCommand, options)
	addCommentUpdateCommand(ctx, commentCommand, options)
	addCommentDeleteCommand(ctx, commentCommand, options)
	addCommentResolveCommand(ctx, commentCommand, options)
	addCommentUnresolveCommand(ctx, commentCommand, options)
}

func addCommentUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.CommentUpdateRequest{}
	bodyFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.CommentSummary]{
		Use:   "update COMMENT_ID",
		Short: "Edit a comment after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Body, "body", "", "new comment body as markdown; use - to read stdin")
			command.Flags().StringVar(&bodyFile, "body-file", "", "read new comment body from file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.CommentSummary, error) {
			request.ID = args[0]
			if err := resolveBodyFlag(command, &request.Body); err != nil {
				return client.CommentSummary{}, err
			}
			if err := resolveFileFlag(command, &request.Body, bodyFile, "body"); err != nil {
				return client.CommentSummary{}, err
			}

			return client.UpdateComment(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeComment,
	})
}

func addCommentDeleteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[string]{
		Use:   "delete COMMENT_ID",
		Short: "Delete a comment after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string) (string, error) {
			return client.DeleteComment(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeDeletion,
	})
}

func addCommentResolveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.CommentSummary]{
		Use:   "resolve COMMENT_ID",
		Short: "Resolve a comment thread after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.CommentSummary, error) {
			return client.ResolveComment(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeComment,
	})
}

func addCommentUnresolveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.CommentSummary]{
		Use:   "unresolve COMMENT_ID",
		Short: "Reopen a comment thread after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.CommentSummary, error) {
			return client.UnresolveComment(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeComment,
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
	addReadGetCommand(ctx, root, options, readGetSpec[client.CommentBotActor]{
		Use:   "bot-actor COMMENT_ID",
		Short: "Show comment bot actor metadata",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.CommentBotActor, error) {
			return client.GetCommentBotActor(ctx, runtime.graphqlClient, id)
		},
		Write: writeCommentBotActor,
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
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum child comments to return")
	root.AddCommand(preflightReadListCommand(command, loadCommentChildren))
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
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(preflightReadListCommand(command, loadCommentCreatedIssues))
}

func writeCommentBotActor(command *cobra.Command, options *rootOptions, actor client.CommentBotActor) error {
	return writeBotActor(command, options, actor, actor.CommentID, actor.Bot)
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
