//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeRelationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand[
		client.InitiativeRelationList,
		client.InitiativeRelationSummary,
	](ctx, root, options, readListGetSpec[client.InitiativeRelationList, client.InitiativeRelationSummary]{
		Use:           "initiative-relation",
		Short:         "Read Linear initiative relations",
		ListShort:     "List visible initiative relations",
		LimitHelp:     "maximum initiative relations to return",
		GetUse:        "get INITIATIVE_RELATION_ID",
		GetShort:      "Get one initiative relation by id",
		LoadList:      loadInitiativeRelationList,
		PageWithItems: initiativeRelationPageWithItems,
		LoadGet:       loadInitiativeRelation,
		WriteItem:     writeInitiativeRelation,
	})
}

func writeInitiativeRelation(
	command *cobra.Command,
	options *rootOptions,
	relation client.InitiativeRelationSummary,
) error {
	if wrote, err := writeIDOnly(command, options, relation.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, relation)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s -> %s order %.2f",
		relation.ID,
		relation.ParentInitiativeName,
		relation.RelatedInitiativeName,
		relation.SortOrder,
	)
}

func loadInitiativeRelationList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeRelationList, []client.InitiativeRelationSummary, error) {
	relations, err := client.ListInitiativeRelations(ctx, runtime.graphqlClient, limit)
	return relations, relations.Relations, err
}

func loadInitiativeRelation(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeRelationSummary, error) {
	return client.GetInitiativeRelationByID(ctx, runtime.graphqlClient, id)
}

func initiativeRelationPageWithItems(
	page client.InitiativeRelationList,
	relations []client.InitiativeRelationSummary,
) client.InitiativeRelationList {
	page.Relations = relations
	return page
}
