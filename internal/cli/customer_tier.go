package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addCustomerTierCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.CustomerTierList, client.CustomerTierSummary]{
			Use:       "customer-tier",
			Short:     "Read Linear customer tiers",
			ListShort: "List organization customer tiers",
			LimitHelp: "maximum customer tiers to return",
			GetUse:    "get CUSTOMER_TIER_ID",
			GetShort:  "Get one customer tier by id",
			LoadList:  clientList(client.ListCustomerTiers),
			LoadGet:   clientGet(client.GetCustomerTierByID),
			WriteItem: writeCustomerTier,
		},
	)
}

func writeCustomerTier(command *cobra.Command, options *rootOptions, tier client.CustomerTierSummary) error {
	return writeItemLine(
		command, options, tier, tier.ID,
		"%s %s %s %.0f", tier.ID, tier.DisplayName, tier.Color, tier.Position,
	)
}
