package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// FavoriteSummary is the compact favorite model used by read-only commands.
type FavoriteSummary struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	FolderName string `json:"folder_name,omitempty"`
	URL        string `json:"url,omitempty"`
}

// FavoriteList is a page of favorites.
type FavoriteList struct {
	Favorites   []FavoriteSummary `json:"favorites"`
	HasNextPage bool              `json:"has_next_page"`
	EndCursor   *string           `json:"end_cursor,omitempty"`
}

//nolint:lll
type favoritesNode = gql.XFavoritesFavoritesFavoriteConnectionNodesFavorite

//nolint:lll
type favoriteChildrenNode = gql.XFavorite_childrenFavoriteChildrenFavoriteConnectionNodesFavorite

type favoritesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type favoriteChildrenQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// ListFavorites returns the authenticated user's favorites.
func ListFavorites(ctx context.Context, graphqlClient graphql.Client, limit int) (FavoriteList, error) {
	query := favoritesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list favorites", limit, defaultListPageSize,
		query.page,
		favoritesNodeSummary,
	)
	if err != nil {
		return FavoriteList{}, err
	}

	return FavoriteList{Favorites: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListFavoriteChildren returns child favorites under a folder favorite.
func ListFavoriteChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (FavoriteList, error) {
	query := favoriteChildrenQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list favorite children "+id, limit, defaultListPageSize,
		query.children,
		favoriteChildrenNodeSummary,
	)
	if err != nil {
		return FavoriteList{}, err
	}

	return FavoriteList{Favorites: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetFavoriteByID returns one favorite by Linear id.
func GetFavoriteByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (FavoriteSummary, error) {
	result, err := gql.XFavorite(ctx, graphqlClient, id)
	if err != nil {
		return FavoriteSummary{}, fmt.Errorf("get favorite %s: %w", id, err)
	}

	return favoriteSummary(result.Favorite.FavoriteSummaryFields), nil
}

func (query favoritesQuery) page(pageSize int, after *string) ([]favoritesNode, bool, *string, error) {
	result, err := gql.XFavorites(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Favorites.Nodes,
		result.Favorites.PageInfo.HasNextPage,
		result.Favorites.PageInfo.EndCursor,
		nil
}

func (query favoriteChildrenQuery) children(
	pageSize int,
	after *string,
) ([]favoriteChildrenNode, bool, *string, error) {
	result, err := gql.XFavorite_children(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Favorite.Children.Nodes,
		result.Favorite.Children.PageInfo.HasNextPage,
		result.Favorite.Children.PageInfo.EndCursor,
		nil
}

func favoritesNodeSummary(node favoritesNode) FavoriteSummary {
	return favoriteSummary(node.FavoriteSummaryFields)
}

func favoriteChildrenNodeSummary(node favoriteChildrenNode) FavoriteSummary {
	return favoriteSummary(node.FavoriteSummaryFields)
}

func favoriteSummary(fields gql.FavoriteSummaryFields) FavoriteSummary {
	return FavoriteSummary{
		ID:         fields.Id,
		Type:       fields.Type,
		FolderName: stringValue(fields.FolderName),
		URL:        stringValue(fields.Url),
	}
}
