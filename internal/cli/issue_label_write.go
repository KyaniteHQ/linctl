//nolint:dupl // Label-association command glue is intentionally uniform across issue and project targets.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueAddLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyWrite, &cobra.Command{
		Use:   "add-label ISSUE_ID LABEL_ID",
		Short: "Attach a label to an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runIssueAddLabel(
				ctx, command, options, issueAdapterFor(runtime),
				client.IssueLabelAssociationRequest{IssueID: args[0], LabelID: args[1]},
			)
		},
	})
}

func addIssueRemoveLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyWrite, &cobra.Command{
		Use:   "remove-label ISSUE_ID LABEL_ID",
		Short: "Detach a label from an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runIssueRemoveLabel(
				ctx, command, options, issueAdapterFor(runtime),
				client.IssueLabelAssociationRequest{IssueID: args[0], LabelID: args[1]},
			)
		},
	})
}

func runIssueAddLabel(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	adder issueLabelAdder,
	request client.IssueLabelAssociationRequest,
) error {
	issue, err := adder.AddIssueLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeIssue(command, options, issue)
}

func runIssueRemoveLabel(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	remover issueLabelRemover,
	request client.IssueLabelAssociationRequest,
) error {
	issue, err := remover.RemoveIssueLabel(ctx, request)
	if err != nil {
		return err
	}

	return writeIssue(command, options, issue)
}
