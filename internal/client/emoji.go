package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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
	Emojis      []EmojiSummary `json:"emojis"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// ListEmojis returns the workspace custom emojis.
func ListEmojis(ctx context.Context, graphqlClient graphql.Client, limit int) (EmojiList, error) {
	result, err := emojis(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return EmojiList{}, fmt.Errorf("list emojis: %w", err)
	}

	summaries := make([]EmojiSummary, 0, len(result.Emojis.Nodes))
	for _, node := range result.Emojis.Nodes {
		summaries = append(summaries, emojiSummary(node.EmojiSummaryFields))
	}

	return EmojiList{
		Emojis:      summaries,
		HasNextPage: result.Emojis.PageInfo.HasNextPage,
		EndCursor:   result.Emojis.PageInfo.EndCursor,
	}, nil
}

// GetEmojiByID returns one custom emoji by Linear id or name.
func GetEmojiByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (EmojiSummary, error) {
	result, err := emoji(ctx, graphqlClient, id)
	if err != nil {
		return EmojiSummary{}, fmt.Errorf("get emoji %s: %w", id, err)
	}

	return emojiSummary(result.Emoji.EmojiSummaryFields), nil
}

func emojiSummary(fields EmojiSummaryFields) EmojiSummary {
	return EmojiSummary{
		ID:     fields.Id,
		Name:   fields.Name,
		URL:    fields.Url,
		Source: fields.Source,
	}
}
