//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	initiativeUpdateCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeUpdateList, client.InitiativeUpdateSummary]{
			Use:           "initiative-update",
			Short:         "Read Linear initiative updates",
			ListShort:     "List visible initiative updates",
			LimitHelp:     "maximum initiative updates to return",
			GetUse:        "get INITIATIVE_UPDATE_ID",
			GetShort:      "Get one initiative update by id",
			LoadList:      loadInitiativeUpdateList,
			PageWithItems: initiativeUpdatePageWithItems,
			LoadGet:       loadInitiativeUpdate,
			WriteItem:     writeInitiativeUpdate,
		},
	)
	addInitiativeUpdateCommentsCommand(ctx, initiativeUpdateCommand, options)
}

func addInitiativeUpdateCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "comments INITIATIVE_UPDATE_ID",
		Short: "List initiative update comments without body content",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeUpdateCommentList,
				initiativeUpdateCommentPageWithItems,
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum comments to return")
	root.AddCommand(command)
}

func writeInitiativeUpdate(command *cobra.Command, options *rootOptions, update client.InitiativeUpdateSummary) error {
	return writeItem(command, options, update, update.ID,
		func(command *cobra.Command, _ *rootOptions, update client.InitiativeUpdateSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %s",
				update.ID,
				update.Health,
				update.DisplayName,
				update.Body,
			)
		})
}

func loadInitiativeUpdateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeUpdateList, []client.InitiativeUpdateSummary, error) {
	updates, err := client.ListInitiativeUpdates(ctx, runtime.graphqlClient, limit)
	return updates, updates.Updates, err
}

func loadInitiativeUpdate(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeUpdateSummary, error) {
	return client.GetInitiativeUpdateByID(ctx, runtime.graphqlClient, id)
}

func initiativeUpdatePageWithItems(
	page client.InitiativeUpdateList,
	updates []client.InitiativeUpdateSummary,
) client.InitiativeUpdateList {
	page.Updates = updates
	return page
}

func loadInitiativeUpdateCommentList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.InitiativeUpdateCommentList, []client.CommentMetadataSummary, error) {
	comments, err := client.ListInitiativeUpdateComments(ctx, runtime.graphqlClient, args[0], limit)
	return comments, comments.Comments, err
}

func initiativeUpdateCommentPageWithItems(
	page client.InitiativeUpdateCommentList,
	comments []client.CommentMetadataSummary,
) client.InitiativeUpdateCommentList {
	page.Comments = comments
	return page
}
