//nolint:dupl // Label-association command glue is intentionally uniform across issue and project targets.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueAddLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "add-label ISSUE_ID LABEL_ID",
		Short: "Attach a label to an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			return client.AddIssueLabel(
				ctx, runtime.graphqlClient, runtime.config.Target,
				client.IssueLabelAssociationRequest{IssueID: args[0], LabelID: args[1]},
			)
		},
		Write: writeIssue,
	})
}

func addIssueRemoveLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "remove-label ISSUE_ID LABEL_ID",
		Short: "Detach a label from an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			return client.RemoveIssueLabel(
				ctx, runtime.graphqlClient, runtime.config.Target,
				client.IssueLabelAssociationRequest{IssueID: args[0], LabelID: args[1]},
			)
		},
		Write: writeIssue,
	})
}
