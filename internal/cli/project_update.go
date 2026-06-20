//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectUpdateReadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.ProjectUpdateList, client.ProjectUpdateSummary]{
		Use:           "project-update",
		Short:         "Read Linear project updates",
		ListShort:     "List visible project updates",
		LimitHelp:     "maximum project updates to return",
		GetUse:        "get PROJECT_UPDATE_ID",
		GetShort:      "Get one project update by id",
		LoadList:      loadProjectUpdateList,
		PageWithItems: projectUpdatePageWithItems,
		LoadGet:       loadProjectUpdate,
		WriteItem:     writeProjectUpdate,
	})
}

func writeProjectUpdate(command *cobra.Command, options *rootOptions, update client.ProjectUpdateSummary) error {
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

func loadProjectUpdateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectUpdateList, []client.ProjectUpdateSummary, error) {
	updates, err := client.ListAllProjectUpdates(ctx, runtime.graphqlClient, limit)
	return updates, updates.Updates, err
}

func loadProjectUpdate(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ProjectUpdateSummary, error) {
	return client.GetProjectUpdateByID(ctx, runtime.graphqlClient, id)
}

func projectUpdatePageWithItems(
	page client.ProjectUpdateList,
	updates []client.ProjectUpdateSummary,
) client.ProjectUpdateList {
	page.Updates = updates
	return page
}
