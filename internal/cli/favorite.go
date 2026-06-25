//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addFavoriteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	favoriteCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.FavoriteList, client.FavoriteSummary]{
			Use:           "favorite",
			Short:         "Read Linear favorites",
			ListShort:     "List the authenticated user's favorites",
			LimitHelp:     "maximum favorites to return",
			GetUse:        "get FAVORITE_ID",
			GetShort:      "Get one favorite by id",
			LoadList:      loadFavoriteList,
			PageWithItems: favoritePageWithItems,
			LoadGet:       loadFavorite,
			WriteItem:     writeFavorite,
		},
	)
	addFavoriteChildrenCommand(ctx, favoriteCommand, options)
}

func addFavoriteChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "children FAVORITE_ID",
		Short: "List children of a folder favorite",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadFavoriteChildren,
				favoritePageWithItems,
				writeFavorite,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum favorites to return")
	root.AddCommand(command)
}

func writeFavorite(command *cobra.Command, options *rootOptions, favorite client.FavoriteSummary) error {
	return writeItem(command, options, favorite, favorite.ID,
		func(command *cobra.Command, _ *rootOptions, favorite client.FavoriteSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s [%s] %s", favorite.ID, favorite.Type, favorite.URL)
		})
}

func loadFavoriteList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.FavoriteList, []client.FavoriteSummary, error) {
	favorites, err := client.ListFavorites(ctx, runtime.graphqlClient, limit)
	return favorites, favorites.Favorites, err
}

func loadFavorite(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.FavoriteSummary, error) {
	return client.GetFavoriteByID(ctx, runtime.graphqlClient, id)
}

func loadFavoriteChildren(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.FavoriteList, []client.FavoriteSummary, error) {
	favorites, err := client.ListFavoriteChildren(ctx, runtime.graphqlClient, args[0], limit)
	return favorites, favorites.Favorites, err
}

func favoritePageWithItems(
	page client.FavoriteList,
	favorites []client.FavoriteSummary,
) client.FavoriteList {
	page.Favorites = favorites
	return page
}
