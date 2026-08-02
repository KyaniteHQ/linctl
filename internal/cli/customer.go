package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addCustomerCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CustomerList, client.CustomerSummary]{
		Use:       "customer",
		Short:     "Read Linear customers",
		ListShort: "List visible Linear customers",
		LimitHelp: "maximum customers to return",
		GetUse:    "get CUSTOMER_ID",
		GetShort:  "Get one customer by id or slug",
		LoadList:  clientList(client.ListCustomers),
		LoadGet:   clientGet(client.GetCustomerByID),
		WriteItem: writeCustomer,
	})
}

func writeCustomer(command *cobra.Command, options *rootOptions, customer client.CustomerSummary) error {
	return writeItemLine(
		command, options, customer, customer.ID,
		"%s %s [%s] needs %.0f", customer.ID, customer.Name, customer.StatusName, customer.ApproximateNeedCount,
	)
}
