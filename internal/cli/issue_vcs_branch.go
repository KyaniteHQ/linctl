package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueVCSBranchSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	branchCommand := newGroupCommand("vcs-branch-search", "Read the issue matched by a VCS branch")
	root.AddCommand(branchCommand)

	addReadGetCommand(ctx, branchCommand, options, readGetSpec[client.IssueSummary]{
		Use:   "get BRANCH_NAME",
		Short: "Get the issue that matches a VCS branch",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.IssueSummary, error) {
			return client.GetIssueByVCSBranch(ctx, runtime.graphqlClient, id)
		},
		Write: writeIssue,
	})
	addIssueChildCommands(ctx, branchCommand, options, issueChildCommandBundleForVCSBranch())
	addIssueVCSBranchCommentsCommand(ctx, branchCommand, options)
}

func addIssueVCSBranchCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"comments BRANCH_NAME",
		"List body-free comments for the issue matched by a VCS branch",
		"comments",
		client.ListIssueVCSBranchComments,
		writeCommentMetadata,
	)
}

func issueChildCommandBundleForVCSBranch() issueChildCommandBundle {
	return issueChildCommandBundle{
		Argument:          "BRANCH_NAME",
		Text:              scopedIssueChildCommandText("the issue matched by a VCS branch"),
		Attachments:       client.ListIssueVCSBranchAttachments,
		BotActor:          client.GetIssueVCSBranchBotActor,
		Children:          client.ListIssueVCSBranchChildren,
		Documents:         client.ListIssueVCSBranchDocuments,
		FormerAttachments: client.ListIssueVCSBranchFormerAttachments,
		FormerNeeds:       client.ListIssueVCSBranchFormerNeeds,
		History:           client.ListIssueVCSBranchHistory,
		InverseRelations:  client.ListIssueVCSBranchInverseRelations,
		Labels:            client.ListIssueVCSBranchLabels,
		Needs:             client.ListIssueVCSBranchNeeds,
		Relations:         client.ListIssueVCSBranchRelations,
		Releases:          client.ListIssueVCSBranchReleases,
		SharedAccess:      client.GetIssueVCSBranchSharedAccess,
		StateHistory:      client.ListIssueVCSBranchStateHistory,
		Subscribers:       client.ListIssueVCSBranchSubscribers,
	}
}
