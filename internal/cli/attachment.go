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
				attachmentPageWithItems,
				writeAttachment,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum attachments to return")
	root.AddCommand(command)
}

func addAttachmentIssueCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	issueCommand := &cobra.Command{
		Use:   "issue",
		Short: "Read the issue associated with an attachment",
	}
	root.AddCommand(issueCommand)

	issueCommand.AddCommand(&cobra.Command{
		Use:   "get ATTACHMENT_ID",
		Short: "Get the issue associated with an attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.GetAttachmentIssue(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
	addAttachmentIssueAttachmentsCommand(ctx, issueCommand, options)
	addAttachmentIssueBotActorCommand(ctx, issueCommand, options)
	addAttachmentIssueChildrenCommand(ctx, issueCommand, options)
	addAttachmentIssueCommentsCommand(ctx, issueCommand, options)
	addAttachmentIssueDocumentsCommand(ctx, issueCommand, options)
	addAttachmentIssueFormerAttachmentsCommand(ctx, issueCommand, options)
	addAttachmentIssueFormerNeedsCommand(ctx, issueCommand, options)
	addAttachmentIssueHistoryCommand(ctx, issueCommand, options)
	addAttachmentIssueInverseRelationsCommand(ctx, issueCommand, options)
	addAttachmentIssueLabelsCommand(ctx, issueCommand, options)
	addAttachmentIssueNeedsCommand(ctx, issueCommand, options)
	addAttachmentIssueRelationsCommand(ctx, issueCommand, options)
	addAttachmentIssueReleasesCommand(ctx, issueCommand, options)
	addAttachmentIssueSharedAccessCommand(ctx, issueCommand, options)
	addAttachmentIssueStateHistoryCommand(ctx, issueCommand, options)
	addAttachmentIssueSubscribersCommand(ctx, issueCommand, options)
}

func addAttachmentIssueCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCommentMetadataListCommand(
		ctx,
		root,
		options,
		"comments ATTACHMENT_ID",
		"List body-free comments for the issue associated with an attachment",
		"comments",
		func(runtime commandRuntime, attachmentID string, limit int) (client.IssueCommentMetadataList, error) {
			return client.ListAttachmentIssueComments(ctx, runtime.graphqlClient, attachmentID, limit)
		},
	)
}

func addAttachmentIssueAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"attachments ATTACHMENT_ID",
		"List attachments for the issue associated with an attachment",
		"attachments",
		func(runtime commandRuntime, attachmentID string, limit int) (client.AttachmentList, error) {
			return client.ListAttachmentIssueAttachments(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		writeAttachment,
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addAttachmentIssueBotActorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "bot-actor ATTACHMENT_ID",
		Short: "Show bot actor metadata for the issue associated with an attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			actor, err := client.GetAttachmentIssueBotActor(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssueBotActor(command, options, actor)
		},
	})
}

func addAttachmentIssueChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children ATTACHMENT_ID",
		"List child issues for the issue associated with an attachment",
		"child issues",
		func(runtime commandRuntime, attachmentID string, limit int) (client.IssueList, error) {
			return client.ListAttachmentIssueChildren(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.IssueList) int {
			return len(list.Issues)
		},
		func(list client.IssueList) (client.IssueList, error) {
			items, err := sortByJSONField(list.Issues, options.sortField, options.sortOrder)
			list.Issues = items
			return list, err
		},
		writeIssue,
		func(list client.IssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addAttachmentIssueDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"documents ATTACHMENT_ID",
		"List documents for the issue associated with an attachment",
		"documents",
		func(runtime commandRuntime, attachmentID string, limit int) (client.DocumentList, error) {
			return client.ListAttachmentIssueDocuments(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.DocumentList) int {
			return len(list.Documents)
		},
		func(list client.DocumentList) (client.DocumentList, error) {
			items, err := sortByJSONField(list.Documents, options.sortField, options.sortOrder)
			list.Documents = items
			return list, err
		},
		writeDocument,
		func(list client.DocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addAttachmentIssueFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"former-attachments ATTACHMENT_ID",
		"List former attachments for the issue associated with an attachment",
		"former attachments",
		func(runtime commandRuntime, attachmentID string, limit int) (client.AttachmentList, error) {
			return client.ListAttachmentIssueFormerAttachments(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		writeAttachment,
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addAttachmentIssueFormerNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCustomerNeedMetadataListCommand(
		ctx,
		root,
		options,
		"former-needs ATTACHMENT_ID",
		"List body-free former customer needs for the issue associated with an attachment",
		"former customer needs",
		func(runtime commandRuntime, attachmentID string, limit int) (client.IssueCustomerNeedMetadataList, error) {
			return client.ListAttachmentIssueFormerNeeds(ctx, runtime.graphqlClient, attachmentID, limit)
		},
	)
}

func addAttachmentIssueHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"history ATTACHMENT_ID",
		"List history metadata for the issue associated with an attachment",
		"history entries",
		func(runtime commandRuntime, attachmentID string, limit int) (client.IssueHistoryList, error) {
			return client.ListAttachmentIssueHistory(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.IssueHistoryList) int {
			return len(list.History)
		},
		func(list client.IssueHistoryList) (client.IssueHistoryList, error) {
			items, err := sortByJSONField(list.History, options.sortField, options.sortOrder)
			list.History = items
			return list, err
		},
		writeIssueHistory,
		func(list client.IssueHistoryList) []client.IssueHistorySummary {
			return list.History
		},
	)
}

func addAttachmentIssueInverseRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"inverse-relations ATTACHMENT_ID",
		"List inverse relations for the issue associated with an attachment",
		"inverse relations",
		func(
			ctx context.Context,
			runtime commandRuntime,
			attachmentID string,
			limit int,
		) (client.IssueRelationList, error) {
			return client.ListAttachmentIssueInverseRelations(ctx, runtime.graphqlClient, attachmentID, limit)
		},
	)
}

func addAttachmentIssueLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"labels ATTACHMENT_ID",
		"List labels for the issue associated with an attachment",
		"labels",
		func(runtime commandRuntime, attachmentID string, limit int) (client.LabelList, error) {
			return client.ListAttachmentIssueLabels(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.LabelList) int {
			return len(list.Labels)
		},
		func(list client.LabelList) (client.LabelList, error) {
			items, err := sortByJSONField(list.Labels, options.sortField, options.sortOrder)
			list.Labels = items
			return list, err
		},
		writeLabel,
		func(list client.LabelList) []client.LabelSummary {
			return list.Labels
		},
	)
}

func addAttachmentIssueNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCustomerNeedMetadataListCommand(
		ctx,
		root,
		options,
		"needs ATTACHMENT_ID",
		"List body-free customer needs for the issue associated with an attachment",
		"customer needs",
		func(runtime commandRuntime, attachmentID string, limit int) (client.IssueCustomerNeedMetadataList, error) {
			return client.ListAttachmentIssueNeeds(ctx, runtime.graphqlClient, attachmentID, limit)
		},
	)
}

func addAttachmentIssueRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"relations ATTACHMENT_ID",
		"List relations for the issue associated with an attachment",
		"relations",
		func(
			ctx context.Context,
			runtime commandRuntime,
			attachmentID string,
			limit int,
		) (client.IssueRelationList, error) {
			return client.ListAttachmentIssueRelations(ctx, runtime.graphqlClient, attachmentID, limit)
		},
	)
}

func addAttachmentIssueReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"releases ATTACHMENT_ID",
		"List releases for the issue associated with an attachment",
		"releases",
		func(runtime commandRuntime, attachmentID string, limit int) (client.ReleaseList, error) {
			return client.ListAttachmentIssueReleases(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.ReleaseList) int {
			return len(list.Releases)
		},
		func(list client.ReleaseList) (client.ReleaseList, error) {
			items, err := sortByJSONField(list.Releases, options.sortField, options.sortOrder)
			list.Releases = items
			return list, err
		},
		writeRelease,
		func(list client.ReleaseList) []client.ReleaseSummary {
			return list.Releases
		},
	)
}

func addAttachmentIssueSharedAccessCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "shared-access ATTACHMENT_ID",
		Short: "Show shared-access metadata for the issue associated with an attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			access, err := client.GetAttachmentIssueSharedAccess(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssueSharedAccess(command, options, access)
		},
	})
}

func addAttachmentIssueStateHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"state-history ATTACHMENT_ID",
		"List workflow state history for the issue associated with an attachment",
		"state spans",
		func(runtime commandRuntime, attachmentID string, limit int) (client.IssueStateHistoryList, error) {
			return client.ListAttachmentIssueStateHistory(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.IssueStateHistoryList) int {
			return len(list.Spans)
		},
		func(list client.IssueStateHistoryList) (client.IssueStateHistoryList, error) {
			items, err := sortByJSONField(list.Spans, options.sortField, options.sortOrder)
			list.Spans = items
			return list, err
		},
		writeIssueStateSpan,
		func(list client.IssueStateHistoryList) []client.IssueStateSpanSummary {
			return list.Spans
		},
	)
}

func addAttachmentIssueSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"subscribers ATTACHMENT_ID",
		"List subscribers for the issue associated with an attachment",
		"subscribers",
		func(runtime commandRuntime, attachmentID string, limit int) (client.UserList, error) {
			return client.ListAttachmentIssueSubscribers(ctx, runtime.graphqlClient, attachmentID, limit)
		},
		func(list client.UserList) int {
			return len(list.Users)
		},
		func(list client.UserList) (client.UserList, error) {
			items, err := sortByJSONField(list.Users, options.sortField, options.sortOrder)
			list.Users = items
			return list, err
		},
		writeUser,
		func(list client.UserList) []client.UserSummary {
			return list.Users
		},
	)
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

func attachmentPageWithItems(
	page client.AttachmentList,
	attachments []client.AttachmentSummary,
) client.AttachmentList {
	page.Attachments = attachments
	return page
}
