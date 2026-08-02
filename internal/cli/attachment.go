package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
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
			LoadList:  clientList(client.ListAttachments),
			LoadGet:   clientGet(client.GetAttachmentByID),
			WriteItem: writeAttachment,
		},
	)
	addAttachmentURLCommand(ctx, attachmentCommand, options)
	addAttachmentIssueCommand(ctx, attachmentCommand, options)
}

func addAttachmentURLCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"url URL",
		"List visible issue attachments for a URL",
		"attachments",
		client.ListAttachmentsForURL,
		writeAttachment,
	)
}

func addAttachmentIssueCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	issueCommand := newGroupCommand("issue", "Read the issue associated with an attachment")
	root.AddCommand(issueCommand)

	addReadGetCommand(ctx, issueCommand, options, readGetSpec[client.IssueSummary]{
		Use:   "get ATTACHMENT_ID",
		Short: "Get the issue attached to an attachment",
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
	return writeItemLine(
		command, options, attachment, attachment.ID,
		"%s %s [%s]", attachment.ID, attachment.Title, attachment.SourceType,
	)
}
