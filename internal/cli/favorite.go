//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addFavoriteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.FavoriteList, client.FavoriteSummary]{
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
	})
}

func writeFavorite(
	command *cobra.Command,
	options *rootOptions,
	favorite client.FavoriteSummary,
) error {
	if wrote, err := writeIDOnly(command, options, favorite.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, favorite)
	}

	return render.WriteLine(command.OutOrStdout(), "%s [%s] %s", favorite.ID, favorite.Type, favorite.URL)
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

func favoritePageWithItems(
	page client.FavoriteList,
	favorites []client.FavoriteSummary,
) client.FavoriteList {
	page.Favorites = favorites
	return page
}
