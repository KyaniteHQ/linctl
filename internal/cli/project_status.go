package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectStatusCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ProjectStatusList, client.ProjectStatusSummary]{
			Use:       "project-status",
			Short:     "Read Linear project statuses",
			ListShort: "List visible Linear project statuses",
			LimitHelp: "maximum project statuses to return",
			GetUse:    "get PROJECT_STATUS_ID",
			GetShort:  "Get one project status by id",
			LoadList:  clientList(client.ListProjectStatuses),
			LoadGet:   clientGet(client.GetProjectStatusByID),
			WriteItem: writeProjectStatus,
		},
	)
	addProjectStatusProjectCountCommand(ctx, projectStatusCommand, options)
}

func writeProjectStatus(command *cobra.Command, options *rootOptions, status client.ProjectStatusSummary) error {
	return writeItemLine(
		command, options, status, status.ID,
		"%s %s [%s] %s", status.ID, status.Name, status.Type, status.Color,
	)
}

func addProjectStatusProjectCountCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.ProjectStatusProjectCount]{
		Use:   "project-count PROJECT_STATUS_ID",
		Short: "Show the project counts of one project status",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.ProjectStatusProjectCount, error) {
			return client.GetProjectStatusProjectCount(ctx, runtime.graphqlClient, id)
		},
		Write: writeProjectStatusProjectCount,
	})
}

func writeProjectStatusProjectCount(
	command *cobra.Command,
	options *rootOptions,
	count client.ProjectStatusProjectCount,
) error {
	return writeItemLine(
		command, options, count, count.ProjectStatusID,
		"%s count %.0f private %.0f archived_team %.0f",
		count.ProjectStatusID,
		count.Count,
		count.PrivateCount,
		count.ArchivedTeamCount,
	)
}
