package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueVCSBranchSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	branchCommand := &cobra.Command{
		Use:   "vcs-branch-search",
		Short: "Read the issue matched by a VCS branch",
	}
	root.AddCommand(branchCommand)

	branchCommand.AddCommand(&cobra.Command{
		Use:   "get BRANCH_NAME",
		Short: "Get the issue matched by a VCS branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := issueAdapterFor(runtime).GetIssueByVCSBranch(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
	addIssueVCSBranchAttachmentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchBotActorCommand(ctx, branchCommand, options)
	addIssueVCSBranchChildrenCommand(ctx, branchCommand, options)
	addIssueVCSBranchDocumentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchFormerAttachmentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchCommentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchFormerNeedsCommand(ctx, branchCommand, options)
	addIssueVCSBranchHistoryCommand(ctx, branchCommand, options)
	addIssueVCSBranchInverseRelationsCommand(ctx, branchCommand, options)
	addIssueVCSBranchLabelsCommand(ctx, branchCommand, options)
	addIssueVCSBranchNeedsCommand(ctx, branchCommand, options)
	addIssueVCSBranchRelationsCommand(ctx, branchCommand, options)
	addIssueVCSBranchReleasesCommand(ctx, branchCommand, options)
	addIssueVCSBranchSharedAccessCommand(ctx, branchCommand, options)
	addIssueVCSBranchStateHistoryCommand(ctx, branchCommand, options)
	addIssueVCSBranchSubscribersCommand(ctx, branchCommand, options)
}

func addIssueVCSBranchCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCommentMetadataListCommand(
		ctx,
		root,
		options,
		"comments BRANCH_NAME",
		"List body-free comments for the issue matched by a VCS branch",
		"comments",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueCommentMetadataList, error) {
			return issueAdapterFor(runtime).ListIssueVCSBranchComments(ctx, branchName, limit)
		},
	)
}

func addIssueVCSBranchFormerNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCustomerNeedMetadataListCommand(
		ctx,
		root,
		options,
		"former-needs BRANCH_NAME",
		"List body-free former customer needs for the issue matched by a VCS branch",
		"former customer needs",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueCustomerNeedMetadataList, error) {
			return issueAdapterFor(runtime).ListIssueVCSBranchFormerNeeds(ctx, branchName, limit)
		},
	)
}

func addIssueVCSBranchAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.AttachmentList, client.AttachmentSummary]{
		Use:       "attachments BRANCH_NAME",
		Short:     "List attachments for the issue matched by a VCS branch",
		LimitHelp: "attachments",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.AttachmentList, []client.AttachmentSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchAttachments(ctx, args[0], limit)
			return list, list.Attachments, err
		},
		PageWithItems: attachmentPageWithItems,
		WriteItem:     writeAttachment,
	})
}

func addIssueVCSBranchBotActorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "bot-actor BRANCH_NAME",
		Short: "Show bot actor metadata for the issue matched by a VCS branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			actor, err := issueAdapterFor(runtime).GetIssueVCSBranchBotActor(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssueBotActor(command, options, actor)
		},
	})
}

func addIssueVCSBranchChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "children BRANCH_NAME",
		Short:     "List child issues for the issue matched by a VCS branch",
		LimitHelp: "child issues",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueList, []client.IssueSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchChildren(ctx, args[0], limit)
			return list, list.Issues, err
		},
		PageWithItems: issuePageWithItems,
		WriteItem:     writeIssue,
	})
}

func addIssueVCSBranchDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.DocumentList, client.DocumentSummary]{
		Use:       "documents BRANCH_NAME",
		Short:     "List documents for the issue matched by a VCS branch",
		LimitHelp: "documents",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.DocumentList, []client.DocumentSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchDocuments(ctx, args[0], limit)
			return list, list.Documents, err
		},
		PageWithItems: documentPageWithItems,
		WriteItem:     writeDocument,
	})
}

func addIssueVCSBranchFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.AttachmentList, client.AttachmentSummary]{
		Use:       "former-attachments BRANCH_NAME",
		Short:     "List former attachments for the issue matched by a VCS branch",
		LimitHelp: "former attachments",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.AttachmentList, []client.AttachmentSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchFormerAttachments(ctx, args[0], limit)
			return list, list.Attachments, err
		},
		PageWithItems: attachmentPageWithItems,
		WriteItem:     writeAttachment,
	})
}

func addIssueVCSBranchHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueHistoryList, client.IssueHistorySummary]{
		Use:       "history BRANCH_NAME",
		Short:     "List history metadata for the issue matched by a VCS branch",
		LimitHelp: "history entries",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueHistoryList, []client.IssueHistorySummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchHistory(ctx, args[0], limit)
			return list, list.History, err
		},
		PageWithItems: issueHistoryPageWithItems,
		WriteItem:     writeIssueHistory,
	})
}

func addIssueVCSBranchInverseRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"inverse-relations BRANCH_NAME",
		"List inverse relations for the issue matched by a VCS branch",
		"inverse relations",
		func(
			ctx context.Context,
			runtime commandRuntime,
			branchName string,
			limit int,
		) (client.IssueRelationList, error) {
			return issueAdapterFor(runtime).ListIssueVCSBranchInverseRelations(ctx, branchName, limit)
		},
	)
}

func addIssueVCSBranchLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.LabelList, client.LabelSummary]{
		Use:       "labels BRANCH_NAME",
		Short:     "List labels for the issue matched by a VCS branch",
		LimitHelp: "labels",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.LabelList, []client.LabelSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchLabels(ctx, args[0], limit)
			return list, list.Labels, err
		},
		PageWithItems: labelPageWithItems,
		WriteItem:     writeLabel,
	})
}

func addIssueVCSBranchNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCustomerNeedMetadataListCommand(
		ctx,
		root,
		options,
		"needs BRANCH_NAME",
		"List body-free customer needs for the issue matched by a VCS branch",
		"customer needs",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueCustomerNeedMetadataList, error) {
			return issueAdapterFor(runtime).ListIssueVCSBranchNeeds(ctx, branchName, limit)
		},
	)
}

func addIssueVCSBranchRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"relations BRANCH_NAME",
		"List relations for the issue matched by a VCS branch",
		"relations",
		func(
			ctx context.Context,
			runtime commandRuntime,
			branchName string,
			limit int,
		) (client.IssueRelationList, error) {
			return issueAdapterFor(runtime).ListIssueVCSBranchRelations(ctx, branchName, limit)
		},
	)
}

func addIssueVCSBranchReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ReleaseList, client.ReleaseSummary]{
		Use:       "releases BRANCH_NAME",
		Short:     "List releases for the issue matched by a VCS branch",
		LimitHelp: "releases",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ReleaseList, []client.ReleaseSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchReleases(ctx, args[0], limit)
			return list, list.Releases, err
		},
		PageWithItems: releasePageWithItems,
		WriteItem:     writeRelease,
	})
}

func addIssueVCSBranchSharedAccessCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "shared-access BRANCH_NAME",
		Short: "Show shared-access metadata for the issue matched by a VCS branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			access, err := issueAdapterFor(runtime).GetIssueVCSBranchSharedAccess(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssueSharedAccess(command, options, access)
		},
	})
}

func addIssueVCSBranchStateHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueStateHistoryList, client.IssueStateSpanSummary]{
		Use:       "state-history BRANCH_NAME",
		Short:     "List workflow state history for the issue matched by a VCS branch",
		LimitHelp: "state spans",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueStateHistoryList, []client.IssueStateSpanSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchStateHistory(ctx, args[0], limit)
			return list, list.Spans, err
		},
		PageWithItems: issueStateSpanPageWithItems,
		WriteItem:     writeIssueStateSpan,
	})
}

func addIssueVCSBranchSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.UserList, client.UserSummary]{
		Use:       "subscribers BRANCH_NAME",
		Short:     "List subscribers for the issue matched by a VCS branch",
		LimitHelp: "subscribers",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.UserList, []client.UserSummary, error) {
			list, err := issueAdapterFor(runtime).ListIssueVCSBranchSubscribers(ctx, args[0], limit)
			return list, list.Users, err
		},
		PageWithItems: userPageWithItems,
		WriteItem:     writeUser,
	})
}
