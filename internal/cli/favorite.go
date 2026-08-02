//nolint:dupl // Declarative registration only; addReadListGetCommand and addChildListCommand own the behavior.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addFavoriteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	favoriteCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.FavoriteList, client.FavoriteSummary]{
			Use:       "favorite",
			Short:     "Read Linear favorites",
			ListShort: "List the favorites of the authenticated user",
			LimitHelp: "maximum favorites to return",
			GetUse:    "get FAVORITE_ID",
			GetShort:  "Get one favorite by id",
			LoadList:  clientList(client.ListFavorites),
			LoadGet:   clientGet(client.GetFavoriteByID),
			WriteItem: writeFavorite,
		},
	)
	addFavoriteChildrenCommand(ctx, favoriteCommand, options)
}

func addFavoriteChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children FAVORITE_ID",
		"List children of a folder favorite",
		"favorites",
		client.ListFavoriteChildren,
		writeFavorite,
	)
}

func writeFavorite(command *cobra.Command, options *rootOptions, favorite client.FavoriteSummary) error {
	return writeItemLine(
		command, options, favorite, favorite.ID,
		"%s [%s] %s", favorite.ID, favorite.Type, favorite.URL,
	)
}
