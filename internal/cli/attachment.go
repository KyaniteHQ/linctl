//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addAttachmentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	attachmentCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.AttachmentList, client.AttachmentSummary]{
			Use:           "attachment",
			Short:         "Read Linear attachments",
			ListShort:     "List visible issue attachments",
			LimitHelp:     "maximum attachments to return",
			GetUse:        "get ATTACHMENT_ID",
			GetShort:      "Get one attachment by id",
			LoadList:      loadAttachmentList,
			PageWithItems: attachmentPageWithItems,
			LoadGet:       loadAttachment,
			WriteItem:     writeAttachment,
		},
	)
	addAttachmentURLCommand(ctx, attachmentCommand, options)
}

func addAttachmentURLCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "url URL",
		Short: "List visible issue attachments for a URL",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadAttachmentURLList,
				attachmentPageWithItems,
				writeAttachment,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum attachments to return")
	root.AddCommand(command)
}

func writeAttachment(
	command *cobra.Command,
	options *rootOptions,
	attachment client.AttachmentSummary,
) error {
	if wrote, err := writeIDOnly(command, options, attachment.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, attachment)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", attachment.ID, attachment.Title, attachment.SourceType)
}

func loadAttachmentList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.AttachmentList, []client.AttachmentSummary, error) {
	attachments, err := client.ListAttachments(ctx, runtime.graphqlClient, limit)
	return attachments, attachments.Attachments, err
}

func loadAttachment(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.AttachmentSummary, error) {
	return client.GetAttachmentByID(ctx, runtime.graphqlClient, id)
}

func loadAttachmentURLList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.AttachmentList, []client.AttachmentSummary, error) {
	attachments, err := client.ListAttachmentsForURL(ctx, runtime.graphqlClient, args[0], limit)
	return attachments, attachments.Attachments, err
}

func attachmentPageWithItems(
	page client.AttachmentList,
	attachments []client.AttachmentSummary,
) client.AttachmentList {
	page.Attachments = attachments
	return page
}
