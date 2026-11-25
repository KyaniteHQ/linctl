package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
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
}

func writeProjectStatus(command *cobra.Command, options *rootOptions, status client.ProjectStatusSummary) error {
	return writeItemLine(
		command, options, status, status.ID,
		"%s %s [%s] %s", status.ID, status.Name, status.Type, status.Color,
	)
}
