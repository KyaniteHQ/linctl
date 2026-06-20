//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomerNeedCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CustomerNeedList, client.CustomerNeedSummary]{
		Use:           "customer-need",
		Short:         "Read Linear customer needs",
		ListShort:     "List visible Linear customer needs",
		LimitHelp:     "maximum customer needs to return",
		GetUse:        "get CUSTOMER_NEED_ID",
		GetShort:      "Get one customer need by id",
		LoadList:      loadCustomerNeedList,
		PageWithItems: customerNeedPageWithItems,
		LoadGet:       loadCustomerNeed,
		WriteItem:     writeCustomerNeed,
	})
}

func writeCustomerNeed(
	command *cobra.Command,
	options *rootOptions,
	need client.CustomerNeedSummary,
) error {
	if wrote, err := writeIDOnly(command, options, need.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, need)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s priority %.0f",
		need.ID,
		emptyDash(need.CustomerName),
		emptyDash(need.Issue),
		need.Priority,
	)
}

func loadCustomerNeedList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CustomerNeedList, []client.CustomerNeedSummary, error) {
	needs, err := client.ListCustomerNeeds(ctx, runtime.graphqlClient, limit)
	return needs, needs.Needs, err
}

func loadCustomerNeed(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CustomerNeedSummary, error) {
	return client.GetCustomerNeedByID(ctx, runtime.graphqlClient, id)
}

func customerNeedPageWithItems(
	page client.CustomerNeedList,
	needs []client.CustomerNeedSummary,
) client.CustomerNeedList {
	page.Needs = needs
	return page
}
