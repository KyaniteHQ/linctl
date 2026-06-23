//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomerTierCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.CustomerTierList, client.CustomerTierSummary]{
		Use:           "customer-tier",
		Short:         "Read Linear customer tiers",
		ListShort:     "List organization customer tiers",
		LimitHelp:     "maximum customer tiers to return",
		GetUse:        "get CUSTOMER_TIER_ID",
		GetShort:      "Get one customer tier by id",
		LoadList:      loadCustomerTierList,
		PageWithItems: customerTierPageWithItems,
		LoadGet:       loadCustomerTier,
		WriteItem:     writeCustomerTier,
	})
}

func writeCustomerTier(
	command *cobra.Command,
	options *rootOptions,
	tier client.CustomerTierSummary,
) error {
	if wrote, err := writeIDOnly(command, options, tier.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, tier)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %.0f",
		tier.ID,
		tier.DisplayName,
		tier.Color,
		tier.Position,
	)
}

func loadCustomerTierList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CustomerTierList, []client.CustomerTierSummary, error) {
	tiers, err := client.ListCustomerTiers(ctx, runtime.graphqlClient, limit)
	return tiers, tiers.Tiers, err
}

func loadCustomerTier(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CustomerTierSummary, error) {
	return client.GetCustomerTierByID(ctx, runtime.graphqlClient, id)
}

func customerTierPageWithItems(
	page client.CustomerTierList,
	tiers []client.CustomerTierSummary,
) client.CustomerTierList {
	page.Tiers = tiers
	return page
}
