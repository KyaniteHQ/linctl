//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand[
		client.InitiativeUpdateList,
		client.InitiativeUpdateSummary,
	](ctx, root, options, readListGetSpec[client.InitiativeUpdateList, client.InitiativeUpdateSummary]{
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
	})
}

func writeInitiativeUpdate(
	command *cobra.Command,
	options *rootOptions,
	update client.InitiativeUpdateSummary,
) error {
	if wrote, err := writeIDOnly(command, options, update.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, update)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %s",
		update.ID,
		update.Health,
		update.DisplayName,
		update.Body,
	)
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
