package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueVCSBranchSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	branchCommand := newGroupCommand("vcs-branch-search", "Read the issue matched by a VCS branch")
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
		commentMetadataListItems,
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
