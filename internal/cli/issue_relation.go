package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueRelationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.IssueRelationList, client.IssueRelationSummary]{
			Use:       "issue-relation",
			Short:     "Read Linear issue relations",
			ListShort: "List visible issue relations",
			LimitHelp: "maximum issue relations to return",
			GetUse:    "get ISSUE_RELATION_ID",
			GetShort:  "Get one issue relation by id",
			LoadList:  clientList(client.ListIssueRelations),
			LoadGet:   clientGet(client.GetIssueRelationByID),
			WriteItem: writeIssueRelation,
		},
	)
}

func writeIssueRelation(command *cobra.Command, options *rootOptions, relation client.IssueRelationSummary) error {
	return writeItemLine(
		command, options, relation, relation.ID,
		"%s %s %s -> %s", relation.ID, relation.Type, relation.IssueIdentifier, relation.RelatedIssueIdentifier,
	)
}
