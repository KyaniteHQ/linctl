package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
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
			LoadList:  clientList(client.ListComments),
			LoadGet:   clientGet(client.GetCommentByID),
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
		Short: "Update a comment after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(
				&request.Body, "body", "",
				"new comment body as markdown, or - to read the body from stdin",
			)
			command.Flags().StringVar(&bodyFile, "body-file", "", "read new comment body from file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.CommentSummary, error) {
			request.ID = args[0]
			if err := resolveBodyOrFileFlag(command, &request.Body, bodyFile, "body"); err != nil {
				return client.CommentSummary{}, err
			}

			return client.UpdateComment(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeComment,
	})
}

func addCommentDeleteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[string]{
		Use:          "delete COMMENT_ID",
		Short:        "Delete a comment after pinned-target comparison",
		Args:         cobra.ExactArgs(1),
		Irreversible: true,
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
		Short: "Unresolve a comment thread after pinned-target comparison",
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
	return writeItemLine(
		command, options, comment, comment.ID,
		"%s %s %s", comment.ID, emptyDash(comment.DisplayName), comment.Body,
	)
}

func addCommentBotActorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.CommentBotActor]{
		Use:   "bot-actor COMMENT_ID",
		Short: "Show the bot actor metadata of a comment",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.CommentBotActor, error) {
			return client.GetCommentBotActor(ctx, runtime.graphqlClient, id)
		},
		Write: writeCommentBotActor,
	})
}

func addCommentChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children COMMENT_ID",
		"List child comments without body content",
		"child comments",
		client.ListCommentChildren,
		writeCommentMetadata,
	)
}

func addCommentCreatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"created-issues COMMENT_ID",
		"List issues created from a comment",
		"issues",
		client.ListCommentCreatedIssues,
		writeIssue,
	)
}

func writeCommentBotActor(command *cobra.Command, options *rootOptions, actor client.CommentBotActor) error {
	return writeBotActor(command, options, actor, actor.CommentID, actor.Bot)
}
