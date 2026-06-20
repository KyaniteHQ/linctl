//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.ProjectStatusList, client.ProjectStatusSummary]{
		Use:           "project-status",
		Short:         "Read Linear project statuses",
		ListShort:     "List visible Linear project statuses",
		LimitHelp:     "maximum project statuses to return",
		GetUse:        "get PROJECT_STATUS_ID",
		GetShort:      "Get one project status by id",
		LoadList:      loadProjectStatusList,
		PageWithItems: projectStatusPageWithItems,
		LoadGet:       loadProjectStatus,
		WriteItem:     writeProjectStatus,
	})
}

func writeProjectStatus(
	command *cobra.Command,
	options *rootOptions,
	status client.ProjectStatusSummary,
) error {
	if wrote, err := writeIDOnly(command, options, status.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, status)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s [%s] %s",
		status.ID,
		status.Name,
		status.Type,
		status.Color,
	)
}

func loadProjectStatusList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectStatusList, []client.ProjectStatusSummary, error) {
	statuses, err := client.ListProjectStatuses(ctx, runtime.graphqlClient, limit)
	return statuses, statuses.ProjectStatuses, err
}

func loadProjectStatus(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ProjectStatusSummary, error) {
	return client.GetProjectStatusByID(ctx, runtime.graphqlClient, id)
}

func projectStatusPageWithItems(
	page client.ProjectStatusList,
	statuses []client.ProjectStatusSummary,
) client.ProjectStatusList {
	page.ProjectStatuses = statuses
	return page
}
