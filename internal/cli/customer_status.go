package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addCustomerStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.CustomerStatusList, client.CustomerStatusSummary]{
			Use:       "customer-status",
			Short:     "Read Linear customer statuses",
			ListShort: "List organization customer statuses",
			LimitHelp: "maximum customer statuses to return",
			GetUse:    "get CUSTOMER_STATUS_ID",
			GetShort:  "Get one customer status by id",
			LoadList:  clientList(client.ListCustomerStatuses),
			LoadGet:   clientGet(client.GetCustomerStatusByID),
			WriteItem: writeCustomerStatus,
		},
	)
}

func writeCustomerStatus(command *cobra.Command, options *rootOptions, status client.CustomerStatusSummary) error {
	return writeItemLine(
		command, options, status, status.ID,
		"%s %s %s %.0f", status.ID, status.DisplayName, status.Color, status.Position,
	)
}
