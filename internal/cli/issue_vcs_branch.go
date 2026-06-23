//nolint:dupl // Issue child read commands intentionally share the same list-command shape.
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
			issue, err := client.GetIssueByVCSBranch(ctx, runtime.graphqlClient, args[0])
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
			return client.ListIssueVCSBranchComments(ctx, runtime.graphqlClient, branchName, limit)
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
			return client.ListIssueVCSBranchFormerNeeds(ctx, runtime.graphqlClient, branchName, limit)
		},
	)
}

func addIssueVCSBranchAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"attachments BRANCH_NAME",
		"List attachments for the issue matched by a VCS branch",
		"attachments",
		func(runtime commandRuntime, branchName string, limit int) (client.AttachmentList, error) {
			return client.ListIssueVCSBranchAttachments(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
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
			actor, err := client.GetIssueVCSBranchBotActor(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssueBotActor(command, options, actor)
		},
	})
}

func addIssueVCSBranchChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children BRANCH_NAME",
		"List child issues for the issue matched by a VCS branch",
		"child issues",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueList, error) {
			return client.ListIssueVCSBranchChildren(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.IssueList) int {
			return len(list.Issues)
		},
		func(list client.IssueList) (client.IssueList, error) {
			items, err := sortByJSONField(list.Issues, options.sortField, options.sortOrder)
			list.Issues = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueSummary) error {
			return writeIssue(command, options, item)
		},
		func(list client.IssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addIssueVCSBranchDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"documents BRANCH_NAME",
		"List documents for the issue matched by a VCS branch",
		"documents",
		func(runtime commandRuntime, branchName string, limit int) (client.DocumentList, error) {
			return client.ListIssueVCSBranchDocuments(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.DocumentList) int {
			return len(list.Documents)
		},
		func(list client.DocumentList) (client.DocumentList, error) {
			items, err := sortByJSONField(list.Documents, options.sortField, options.sortOrder)
			list.Documents = items
			return list, err
		},
		func(command *cobra.Command, item client.DocumentSummary) error {
			return writeDocument(command, options, item)
		},
		func(list client.DocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addIssueVCSBranchFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"former-attachments BRANCH_NAME",
		"List former attachments for the issue matched by a VCS branch",
		"former attachments",
		func(runtime commandRuntime, branchName string, limit int) (client.AttachmentList, error) {
			return client.ListIssueVCSBranchFormerAttachments(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addIssueVCSBranchHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"history BRANCH_NAME",
		"List history metadata for the issue matched by a VCS branch",
		"history entries",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueHistoryList, error) {
			return client.ListIssueVCSBranchHistory(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.IssueHistoryList) int {
			return len(list.History)
		},
		func(list client.IssueHistoryList) (client.IssueHistoryList, error) {
			items, err := sortByJSONField(list.History, options.sortField, options.sortOrder)
			list.History = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueHistorySummary) error {
			return writeIssueHistory(command, options, item)
		},
		func(list client.IssueHistoryList) []client.IssueHistorySummary {
			return list.History
		},
	)
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
			return client.ListIssueVCSBranchInverseRelations(ctx, runtime.graphqlClient, branchName, limit)
		},
	)
}

func addIssueVCSBranchLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"labels BRANCH_NAME",
		"List labels for the issue matched by a VCS branch",
		"labels",
		func(runtime commandRuntime, branchName string, limit int) (client.LabelList, error) {
			return client.ListIssueVCSBranchLabels(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.LabelList) int {
			return len(list.Labels)
		},
		func(list client.LabelList) (client.LabelList, error) {
			items, err := sortByJSONField(list.Labels, options.sortField, options.sortOrder)
			list.Labels = items
			return list, err
		},
		func(command *cobra.Command, item client.LabelSummary) error {
			return writeLabel(command, options, item)
		},
		func(list client.LabelList) []client.LabelSummary {
			return list.Labels
		},
	)
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
			return client.ListIssueVCSBranchNeeds(ctx, runtime.graphqlClient, branchName, limit)
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
			return client.ListIssueVCSBranchRelations(ctx, runtime.graphqlClient, branchName, limit)
		},
	)
}

func addIssueVCSBranchReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"releases BRANCH_NAME",
		"List releases for the issue matched by a VCS branch",
		"releases",
		func(runtime commandRuntime, branchName string, limit int) (client.ReleaseList, error) {
			return client.ListIssueVCSBranchReleases(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.ReleaseList) int {
			return len(list.Releases)
		},
		func(list client.ReleaseList) (client.ReleaseList, error) {
			items, err := sortByJSONField(list.Releases, options.sortField, options.sortOrder)
			list.Releases = items
			return list, err
		},
		func(command *cobra.Command, item client.ReleaseSummary) error {
			return writeRelease(command, options, item)
		},
		func(list client.ReleaseList) []client.ReleaseSummary {
			return list.Releases
		},
	)
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
			access, err := client.GetIssueVCSBranchSharedAccess(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssueSharedAccess(command, options, access)
		},
	})
}

func addIssueVCSBranchStateHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"state-history BRANCH_NAME",
		"List workflow state history for the issue matched by a VCS branch",
		"state spans",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueStateHistoryList, error) {
			return client.ListIssueVCSBranchStateHistory(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.IssueStateHistoryList) int {
			return len(list.Spans)
		},
		func(list client.IssueStateHistoryList) (client.IssueStateHistoryList, error) {
			items, err := sortByJSONField(list.Spans, options.sortField, options.sortOrder)
			list.Spans = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueStateSpanSummary) error {
			return writeIssueStateSpan(command, options, item)
		},
		func(list client.IssueStateHistoryList) []client.IssueStateSpanSummary {
			return list.Spans
		},
	)
}

func addIssueVCSBranchSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"subscribers BRANCH_NAME",
		"List subscribers for the issue matched by a VCS branch",
		"subscribers",
		func(runtime commandRuntime, branchName string, limit int) (client.UserList, error) {
			return client.ListIssueVCSBranchSubscribers(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.UserList) int {
			return len(list.Users)
		},
		func(list client.UserList) (client.UserList, error) {
			items, err := sortByJSONField(list.Users, options.sortField, options.sortOrder)
			list.Users = items
			return list, err
		},
		func(command *cobra.Command, item client.UserSummary) error {
			return writeUser(command, options, item)
		},
		func(list client.UserList) []client.UserSummary {
			return list.Users
		},
	)
}
