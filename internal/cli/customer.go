//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomerCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CustomerList, client.CustomerSummary]{
		Use:           "customer",
		Short:         "Read Linear customers",
		ListShort:     "List visible Linear customers",
		LimitHelp:     "maximum customers to return",
		GetUse:        "get CUSTOMER_ID",
		GetShort:      "Get one customer by id or slug",
		LoadList:      loadCustomerList,
		PageWithItems: customerPageWithItems,
		LoadGet:       loadCustomer,
		WriteItem:     writeCustomer,
	})
}

func writeCustomer(
	command *cobra.Command,
	options *rootOptions,
	customer client.CustomerSummary,
) error {
	if wrote, err := writeIDOnly(command, options, customer.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, customer)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s [%s] needs %.0f",
		customer.ID,
		customer.Name,
		customer.StatusName,
		customer.ApproximateNeedCount,
	)
}

func loadCustomerList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CustomerList, []client.CustomerSummary, error) {
	customers, err := client.ListCustomers(ctx, runtime.graphqlClient, limit)
	return customers, customers.Customers, err
}

func loadCustomer(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CustomerSummary, error) {
	return client.GetCustomerByID(ctx, runtime.graphqlClient, id)
}

func customerPageWithItems(
	page client.CustomerList,
	customers []client.CustomerSummary,
) client.CustomerList {
	page.Customers = customers
	return page
}
