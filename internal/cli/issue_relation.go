//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueRelationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand[
		client.IssueRelationList,
		client.IssueRelationSummary,
	](ctx, root, options, readListGetSpec[client.IssueRelationList, client.IssueRelationSummary]{
		Use:           "issue-relation",
		Short:         "Read Linear issue relations",
		ListShort:     "List visible issue relations",
		LimitHelp:     "maximum issue relations to return",
		GetUse:        "get ISSUE_RELATION_ID",
		GetShort:      "Get one issue relation by id",
		LoadList:      loadIssueRelationList,
		PageWithItems: issueRelationPageWithItems,
		LoadGet:       loadIssueRelation,
		WriteItem:     writeIssueRelation,
	})
}

func writeIssueRelation(
	command *cobra.Command,
	options *rootOptions,
	relation client.IssueRelationSummary,
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
		"%s %s %s -> %s",
		relation.ID,
		relation.Type,
		relation.IssueIdentifier,
		relation.RelatedIssueIdentifier,
	)
}

func loadIssueRelationList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueRelationList, []client.IssueRelationSummary, error) {
	relations, err := client.ListIssueRelations(ctx, runtime.graphqlClient, limit)
	return relations, relations.Relations, err
}

func loadIssueRelation(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.IssueRelationSummary, error) {
	return client.GetIssueRelationByID(ctx, runtime.graphqlClient, id)
}

func issueRelationPageWithItems(
	page client.IssueRelationList,
	relations []client.IssueRelationSummary,
) client.IssueRelationList {
	page.Relations = relations
	return page
}
