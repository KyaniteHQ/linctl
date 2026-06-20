//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.InitiativeList, client.InitiativeSummary]{
		Use:           "initiative",
		Short:         "Read Linear initiatives",
		ListShort:     "List visible initiatives",
		LimitHelp:     "maximum initiatives to return",
		GetUse:        "get INITIATIVE_ID",
		GetShort:      "Get one initiative by id or slug",
		LoadList:      loadInitiativeList,
		PageWithItems: initiativePageWithItems,
		LoadGet:       loadInitiative,
		WriteItem:     writeInitiative,
	})
}

func writeInitiative(
	command *cobra.Command,
	options *rootOptions,
	initiative client.InitiativeSummary,
) error {
	if wrote, err := writeIDOnly(command, options, initiative.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, initiative)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", initiative.ID, initiative.Name, initiative.Status)
}

func loadInitiativeList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeList, []client.InitiativeSummary, error) {
	initiatives, err := client.ListInitiatives(ctx, runtime.graphqlClient, limit)
	return initiatives, initiatives.Initiatives, err
}

func loadInitiative(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeSummary, error) {
	return client.GetInitiativeByID(ctx, runtime.graphqlClient, id)
}

func initiativePageWithItems(
	page client.InitiativeList,
	initiatives []client.InitiativeSummary,
) client.InitiativeList {
	page.Initiatives = initiatives
	return page
}
