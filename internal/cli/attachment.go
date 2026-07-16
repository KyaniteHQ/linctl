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
			Use:       "attachment",
			Short:     "Read Linear attachments",
			ListShort: "List visible issue attachments",
			LimitHelp: "maximum attachments to return",
			GetUse:    "get ATTACHMENT_ID",
			GetShort:  "Get one attachment by id",
			LoadList:  loadAttachmentList,
			LoadGet:   loadAttachment,
			WriteItem: writeAttachment,
		},
	)
	addAttachmentURLCommand(ctx, attachmentCommand, options)
	addAttachmentIssueCommand(ctx, attachmentCommand, options)
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
				writeAttachment,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum attachments to return")
	root.AddCommand(preflightReadListCommand(command, loadAttachmentURLList))
}

func addAttachmentIssueCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	issueCommand := newGroupCommand("issue", "Read the issue associated with an attachment")
	root.AddCommand(issueCommand)

	addReadGetCommand(ctx, issueCommand, options, readGetSpec[client.IssueSummary]{
		Use:   "get ATTACHMENT_ID",
		Short: "Get the issue associated with an attachment",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.IssueSummary, error) {
			return client.GetAttachmentIssue(ctx, runtime.graphqlClient, id)
		},
		Write: writeIssue,
	})
	addIssueChildCommands(ctx, issueCommand, options, issueChildCommandBundleForAttachment())
	addAttachmentIssueCommentsCommand(ctx, issueCommand, options)
}

func addAttachmentIssueCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"comments ATTACHMENT_ID",
		"List body-free comments for the issue associated with an attachment",
		"comments",
		client.ListAttachmentIssueComments,
		commentMetadataListItems,
		writeCommentMetadata,
	)
}

func issueChildCommandBundleForAttachment() issueChildCommandBundle {
	return issueChildCommandBundle{
		Argument:          "ATTACHMENT_ID",
		Text:              scopedIssueChildCommandText("the issue associated with an attachment"),
		Attachments:       client.ListAttachmentIssueAttachments,
		BotActor:          client.GetAttachmentIssueBotActor,
		Children:          client.ListAttachmentIssueChildren,
		Documents:         client.ListAttachmentIssueDocuments,
		FormerAttachments: client.ListAttachmentIssueFormerAttachments,
		FormerNeeds:       client.ListAttachmentIssueFormerNeeds,
		History:           client.ListAttachmentIssueHistory,
		InverseRelations:  client.ListAttachmentIssueInverseRelations,
		Labels:            client.ListAttachmentIssueLabels,
		Needs:             client.ListAttachmentIssueNeeds,
		Relations:         client.ListAttachmentIssueRelations,
		Releases:          client.ListAttachmentIssueReleases,
		SharedAccess:      client.GetAttachmentIssueSharedAccess,
		StateHistory:      client.ListAttachmentIssueStateHistory,
		Subscribers:       client.ListAttachmentIssueSubscribers,
	}
}

func writeAttachment(command *cobra.Command, options *rootOptions, attachment client.AttachmentSummary) error {
	return writeItem(command, options, attachment, attachment.ID,
		func(command *cobra.Command, _ *rootOptions, attachment client.AttachmentSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s]",
				attachment.ID,
				attachment.Title,
				attachment.SourceType,
			)
		})
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
