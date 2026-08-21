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
			addAllowedProjectFlag(ctx, command, options, &allowedProjects)
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

func addAllowedProjectFlag(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	dest *[]string,
) {
	command.Flags().StringArrayVar(
		dest, "allowed-project", nil,
		"project id both issues may occupy; repeat the flag; required for cross-project relations",
	)
	registerFlagCompletion(command, "allowed-project", flagCompletion(ctx, options, projectIDCandidates))
}

func writeIssueRelationResult(
	command *cobra.Command,
	options *rootOptions,
	result client.IssueRelationWriteResult,
) error {
	return writeItemLine(
		command, options, result, result.ID,
		"%s %s %s -> %s",
		result.ID, result.Type, result.IssueIdentifier, result.RelatedIssueIdentifier,
	)
}

func addIssueUnrelateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var allowedProjects []string
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[string]{
		Use:   "unrelate ISSUE_RELATION_ID",
		Short: "Delete an issue relation after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			addAllowedProjectFlag(ctx, command, options, &allowedProjects)
		},
		Run: func(ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string) (string, error) {
			return client.DeleteIssueRelation(
				ctx, runtime.graphqlClient, runtime.config.Target,
				client.IssueRelationDeleteRequest{
					RelationID:        args[0],
					AllowedProjectIDs: allowedProjects,
				},
			)
		},
		Write: writeDeletion,
	})
}
