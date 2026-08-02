package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectRelationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ProjectRelationList, client.ProjectRelationSummary]{
			Use:       "project-relation",
			Short:     "Read Linear project relations",
			ListShort: "List visible project relations",
			LimitHelp: "maximum project relations to return",
			GetUse:    "get PROJECT_RELATION_ID",
			GetShort:  "Get one project relation by id",
			LoadList:  clientList(client.ListProjectRelations),
			LoadGet:   clientGet(client.GetProjectRelationByID),
			WriteItem: writeProjectRelation,
		},
	)
}

func writeProjectRelation(command *cobra.Command, options *rootOptions, relation client.ProjectRelationSummary) error {
	return writeItemLine(
		command, options, relation, relation.ID,
		"%s %s %s -> %s", relation.ID, relation.Type, relation.ProjectName, relation.RelatedProjectName,
	)
}
