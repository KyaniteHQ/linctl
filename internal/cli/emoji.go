package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addEmojiCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.EmojiList, client.EmojiSummary]{
		Use:       "emoji",
		Short:     "Read Linear custom emojis",
		ListShort: "List organization custom emojis",
		LimitHelp: "maximum emojis to return",
		GetUse:    "get EMOJI_ID",
		GetShort:  "Get one custom emoji by id or name",
		LoadList:  clientList(client.ListEmojis),
		LoadGet:   clientGet(client.GetEmojiByID),
		WriteItem: writeEmoji,
	})
}

func writeEmoji(command *cobra.Command, options *rootOptions, emoji client.EmojiSummary) error {
	return writeItemLine(command, options, emoji, emoji.ID, "%s %s [%s]", emoji.ID, emoji.Name, emoji.Source)
}
