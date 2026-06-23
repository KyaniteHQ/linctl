//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addEmojiCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.EmojiList, client.EmojiSummary]{
		Use:           "emoji",
		Short:         "Read Linear custom emojis",
		ListShort:     "List organization custom emojis",
		LimitHelp:     "maximum emojis to return",
		GetUse:        "get EMOJI_ID",
		GetShort:      "Get one custom emoji by id or name",
		LoadList:      loadEmojiList,
		PageWithItems: emojiPageWithItems,
		LoadGet:       loadEmoji,
		WriteItem:     writeEmoji,
	})
}

func writeEmoji(
	command *cobra.Command,
	options *rootOptions,
	emoji client.EmojiSummary,
) error {
	if wrote, err := writeIDOnly(command, options, emoji.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, emoji)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", emoji.ID, emoji.Name, emoji.Source)
}

func loadEmojiList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.EmojiList, []client.EmojiSummary, error) {
	emojis, err := client.ListEmojis(ctx, runtime.graphqlClient, limit)
	return emojis, emojis.Emojis, err
}

func loadEmoji(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.EmojiSummary, error) {
	return client.GetEmojiByID(ctx, runtime.graphqlClient, id)
}

func emojiPageWithItems(
	page client.EmojiList,
	emojis []client.EmojiSummary,
) client.EmojiList {
	page.Emojis = emojis
	return page
}
