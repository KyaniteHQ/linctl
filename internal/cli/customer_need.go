//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomerNeedCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	customerNeedCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.CustomerNeedList, client.CustomerNeedSummary]{
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
		},
	)
	addCustomerNeedProjectAttachmentCommand(ctx, customerNeedCommand, options)
}

func addCustomerNeedProjectAttachmentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "project-attachment CUSTOMER_NEED_ID",
		Short: "Get the project attachment linked to one customer need",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			attachment, err := client.GetCustomerNeedProjectAttachment(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomerNeedProjectAttachment(command, options, attachment)
		},
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

func writeCustomerNeedProjectAttachment(
	command *cobra.Command,
	options *rootOptions,
	attachment client.CustomerNeedProjectAttachment,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, attachment)
	}
	if attachment.Attachment == nil {
		return render.WriteLine(command.OutOrStdout(), "%s project_attachment -", attachment.CustomerNeedID)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s project_attachment %s %s [%s]",
		attachment.CustomerNeedID,
		attachment.Attachment.ID,
		attachment.Attachment.Title,
		attachment.Attachment.SourceType,
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
