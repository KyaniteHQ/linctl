package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListFavorites returns the authenticated user's favorites.
func ListFavorites(ctx context.Context, graphqlClient graphql.Client, limit int) (FavoriteList, error) {
	result, err := favorites(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return FavoriteList{}, fmt.Errorf("list favorites: %w", err)
	}

	summaries := make([]FavoriteSummary, 0, len(result.Favorites.Nodes))
	for _, node := range result.Favorites.Nodes {
		summaries = append(summaries, favoriteSummary(node.FavoriteSummaryFields))
	}

	return FavoriteList{
		Favorites:   summaries,
		HasNextPage: result.Favorites.PageInfo.HasNextPage,
		EndCursor:   result.Favorites.PageInfo.EndCursor,
	}, nil
}

// GetFavoriteByID returns one favorite by Linear id.
func GetFavoriteByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (FavoriteSummary, error) {
	result, err := favorite(ctx, graphqlClient, id)
	if err != nil {
		return FavoriteSummary{}, fmt.Errorf("get favorite %s: %w", id, err)
	}

	return favoriteSummary(result.Favorite.FavoriteSummaryFields), nil
}

func favoriteSummary(fields FavoriteSummaryFields) FavoriteSummary {
	return FavoriteSummary{
		ID:         fields.Id,
		Type:       fields.Type,
		FolderName: stringValue(fields.FolderName),
		URL:        stringValue(fields.Url),
	}
}
