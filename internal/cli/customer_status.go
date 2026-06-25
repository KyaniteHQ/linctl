//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomerStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CustomerStatusList, client.CustomerStatusSummary]{
		Use:           "customer-status",
		Short:         "Read Linear customer statuses",
		ListShort:     "List organization customer statuses",
		LimitHelp:     "maximum customer statuses to return",
		GetUse:        "get CUSTOMER_STATUS_ID",
		GetShort:      "Get one customer status by id",
		LoadList:      loadCustomerStatusList,
		PageWithItems: customerStatusPageWithItems,
		LoadGet:       loadCustomerStatus,
		WriteItem:     writeCustomerStatus,
	})
}

func writeCustomerStatus(command *cobra.Command, options *rootOptions, status client.CustomerStatusSummary) error {
	return writeItem(command, options, status, status.ID,
		func(command *cobra.Command, _ *rootOptions, status client.CustomerStatusSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %.0f",
				status.ID,
				status.DisplayName,
				status.Color,
				status.Position,
			)
		})
}

func loadCustomerStatusList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CustomerStatusList, []client.CustomerStatusSummary, error) {
	statuses, err := client.ListCustomerStatuses(ctx, runtime.graphqlClient, limit)
	return statuses, statuses.Statuses, err
}

func loadCustomerStatus(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CustomerStatusSummary, error) {
	return client.GetCustomerStatusByID(ctx, runtime.graphqlClient, id)
}

func customerStatusPageWithItems(
	page client.CustomerStatusList,
	statuses []client.CustomerStatusSummary,
) client.CustomerStatusList {
	page.Statuses = statuses
	return page
}
