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
	issueCommand := newGroupCommand("issue", "Read the issue associated with an attachment")
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
	addListCommand(ctx, root, options, listCommandSpec[client.AttachmentList, client.AttachmentSummary]{
		Use:       "attachments ATTACHMENT_ID",
		Short:     "List attachments for the issue associated with an attachment",
		LimitHelp: "attachments",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.AttachmentList, []client.AttachmentSummary, error) {
			list, err := client.ListAttachmentIssueAttachments(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Attachments, err
		},
		PageWithItems: attachmentPageWithItems,
		WriteItem:     writeAttachment,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "children ATTACHMENT_ID",
		Short:     "List child issues for the issue associated with an attachment",
		LimitHelp: "child issues",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueList, []client.IssueSummary, error) {
			list, err := client.ListAttachmentIssueChildren(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Issues, err
		},
		PageWithItems: issuePageWithItems,
		WriteItem:     writeIssue,
	})
}

func addAttachmentIssueDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.DocumentList, client.DocumentSummary]{
		Use:       "documents ATTACHMENT_ID",
		Short:     "List documents for the issue associated with an attachment",
		LimitHelp: "documents",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.DocumentList, []client.DocumentSummary, error) {
			list, err := client.ListAttachmentIssueDocuments(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Documents, err
		},
		PageWithItems: documentPageWithItems,
		WriteItem:     writeDocument,
	})
}

func addAttachmentIssueFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.AttachmentList, client.AttachmentSummary]{
		Use:       "former-attachments ATTACHMENT_ID",
		Short:     "List former attachments for the issue associated with an attachment",
		LimitHelp: "former attachments",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.AttachmentList, []client.AttachmentSummary, error) {
			list, err := client.ListAttachmentIssueFormerAttachments(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Attachments, err
		},
		PageWithItems: attachmentPageWithItems,
		WriteItem:     writeAttachment,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.IssueHistoryList, client.IssueHistorySummary]{
		Use:       "history ATTACHMENT_ID",
		Short:     "List history metadata for the issue associated with an attachment",
		LimitHelp: "history entries",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueHistoryList, []client.IssueHistorySummary, error) {
			list, err := client.ListAttachmentIssueHistory(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.History, err
		},
		PageWithItems: issueHistoryPageWithItems,
		WriteItem:     writeIssueHistory,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.LabelList, client.LabelSummary]{
		Use:       "labels ATTACHMENT_ID",
		Short:     "List labels for the issue associated with an attachment",
		LimitHelp: "labels",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.LabelList, []client.LabelSummary, error) {
			list, err := client.ListAttachmentIssueLabels(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Labels, err
		},
		PageWithItems: labelPageWithItems,
		WriteItem:     writeLabel,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.ReleaseList, client.ReleaseSummary]{
		Use:       "releases ATTACHMENT_ID",
		Short:     "List releases for the issue associated with an attachment",
		LimitHelp: "releases",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ReleaseList, []client.ReleaseSummary, error) {
			list, err := client.ListAttachmentIssueReleases(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Releases, err
		},
		PageWithItems: releasePageWithItems,
		WriteItem:     writeRelease,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.IssueStateHistoryList, client.IssueStateSpanSummary]{
		Use:       "state-history ATTACHMENT_ID",
		Short:     "List workflow state history for the issue associated with an attachment",
		LimitHelp: "state spans",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueStateHistoryList, []client.IssueStateSpanSummary, error) {
			list, err := client.ListAttachmentIssueStateHistory(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Spans, err
		},
		PageWithItems: issueStateSpanPageWithItems,
		WriteItem:     writeIssueStateSpan,
	})
}

func addAttachmentIssueSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.UserList, client.UserSummary]{
		Use:       "subscribers ATTACHMENT_ID",
		Short:     "List subscribers for the issue associated with an attachment",
		LimitHelp: "subscribers",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.UserList, []client.UserSummary, error) {
			list, err := client.ListAttachmentIssueSubscribers(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Users, err
		},
		PageWithItems: userPageWithItems,
		WriteItem:     writeUser,
	})
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
