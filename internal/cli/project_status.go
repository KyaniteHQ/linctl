//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectStatusCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ProjectStatusList, client.ProjectStatusSummary]{
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
		},
	)
	addProjectStatusProjectCountCommand(ctx, projectStatusCommand, options)
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

func addProjectStatusProjectCountCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "project-count PROJECT_STATUS_ID",
		Short: "Show project counts for one project status",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			count, err := client.GetProjectStatusProjectCount(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeProjectStatusProjectCount(command, options, count)
		},
	})
}

func writeProjectStatusProjectCount(
	command *cobra.Command,
	options *rootOptions,
	count client.ProjectStatusProjectCount,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, count)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s count %.0f private %.0f archived_team %.0f",
		count.ProjectStatusID,
		count.Count,
		count.PrivateCount,
		count.ArchivedTeamCount,
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
