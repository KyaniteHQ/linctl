//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectRelationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand[
		client.ProjectRelationList,
		client.ProjectRelationSummary,
	](ctx, root, options, readListGetSpec[client.ProjectRelationList, client.ProjectRelationSummary]{
		Use:           "project-relation",
		Short:         "Read Linear project relations",
		ListShort:     "List visible project relations",
		LimitHelp:     "maximum project relations to return",
		GetUse:        "get PROJECT_RELATION_ID",
		GetShort:      "Get one project relation by id",
		LoadList:      loadProjectRelationList,
		PageWithItems: projectRelationPageWithItems,
		LoadGet:       loadProjectRelation,
		WriteItem:     writeProjectRelation,
	})
}

func writeProjectRelation(
	command *cobra.Command,
	options *rootOptions,
	relation client.ProjectRelationSummary,
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
		relation.ProjectName,
		relation.RelatedProjectName,
	)
}

func loadProjectRelationList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectRelationList, []client.ProjectRelationSummary, error) {
	relations, err := client.ListProjectRelations(ctx, runtime.graphqlClient, limit)
	return relations, relations.Relations, err
}

func loadProjectRelation(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ProjectRelationSummary, error) {
	return client.GetProjectRelationByID(ctx, runtime.graphqlClient, id)
}

func projectRelationPageWithItems(
	page client.ProjectRelationList,
	relations []client.ProjectRelationSummary,
) client.ProjectRelationList {
	page.Relations = relations
	return page
}
