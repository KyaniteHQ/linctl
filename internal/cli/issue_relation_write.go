package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueRelateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	relationType := "related"
	var allowedProjects []string
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueRelationWriteResult]{
		Use:   "relate ISSUE_ID RELATED_ISSUE_ID",
		Short: "Relate two issues after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(
				&relationType, "type", relationType,
				"relation type: blocks, duplicate, related, or similar",
			)
			command.Flags().StringArrayVar(
				&allowedProjects, "allowed-project", nil,
				"project id both issues may occupy; repeat the flag; required to relate across projects",
			)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueRelationWriteResult, error) {
			return client.CreateIssueRelation(
				ctx, runtime.graphqlClient, runtime.config.Target,
				client.IssueRelationCreateRequest{
					IssueID:           args[0],
					RelatedIssueID:    args[1],
					Type:              relationType,
					AllowedProjectIDs: allowedProjects,
				},
			)
		},
		Write: writeIssueRelationResult,
	})
}

func writeIssueRelationResult(
	command *cobra.Command,
	options *rootOptions,
	result client.IssueRelationWriteResult,
) error {
	return writeIssueRelation(command, options, result.IssueRelationSummary)
}

func addIssueUnrelateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[string]{
		Use:   "unrelate ISSUE_RELATION_ID",
		Short: "Delete an issue relation after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string) (string, error) {
			return client.DeleteIssueRelation(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeDeletion,
	})
}
