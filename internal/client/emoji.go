package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// EmojiSummary is the compact custom emoji model used by read-only commands.
type EmojiSummary struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Source string `json:"source"`
}

// EmojiList is a page of custom emojis.
type EmojiList struct {
	Emojis []EmojiSummary `json:"emojis"`
	Page
}

//nolint:lll
type emojisNode = gql.XEmojisEmojisEmojiConnectionNodesEmoji

type emojisQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListEmojis returns custom emojis visible to the authenticated user.
func ListEmojis(ctx context.Context, graphqlClient graphql.Client, limit int) (EmojiList, error) {
	query := emojisQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list emojis", limit, defaultListPageSize,
		query.page,
		emojisNodeSummary,
	)
	if err != nil {
		return EmojiList{}, err
	}

	return EmojiList{Emojis: page.Items, Page: page.Page}, nil
}

// GetEmojiByID returns one custom emoji by Linear id or name.
func GetEmojiByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (EmojiSummary, error) {
	result, err := gql.XEmoji(ctx, graphqlClient, id)
	if err != nil {
		return EmojiSummary{}, fmt.Errorf("get emoji %s: %w", id, err)
	}

	return emojiSummary(result.Emoji.EmojiSummaryFields), nil
}

func (query emojisQuery) page(pageSize int, after *string) ([]emojisNode, bool, *string, error) {
	result, err := gql.XEmojis(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Emojis.Nodes,
		result.Emojis.PageInfo.HasNextPage,
		result.Emojis.PageInfo.EndCursor,
		nil
}

func emojisNodeSummary(node emojisNode) EmojiSummary {
	return emojiSummary(node.EmojiSummaryFields)
}

func emojiSummary(fields gql.EmojiSummaryFields) EmojiSummary {
	return EmojiSummary{
		ID:     fields.Id,
		Name:   fields.Name,
		URL:    fields.Url,
		Source: fields.Source,
	}
}
