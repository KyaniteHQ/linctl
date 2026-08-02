package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addInitiativeRelationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeRelationList, client.InitiativeRelationSummary]{
			Use:       "initiative-relation",
			Short:     "Read Linear initiative relations",
			ListShort: "List visible initiative relations",
			LimitHelp: "maximum initiative relations to return",
			GetUse:    "get INITIATIVE_RELATION_ID",
			GetShort:  "Get one initiative relation by id",
			LoadList:  clientList(client.ListInitiativeRelations),
			LoadGet:   clientGet(client.GetInitiativeRelationByID),
			WriteItem: writeInitiativeRelation,
		},
	)
}

func writeInitiativeRelation(
	command *cobra.Command,
	options *rootOptions,
	relation client.InitiativeRelationSummary,
) error {
	return writeItemLine(
		command, options, relation, relation.ID,
		"%s %s -> %s order %.2f",
		relation.ID,
		relation.ParentInitiativeName,
		relation.RelatedInitiativeName,
		relation.SortOrder,
	)
}
