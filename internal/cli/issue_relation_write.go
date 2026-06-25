package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueRelateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	relationType := "related"
	command := &cobra.Command{
		Use:   "relate ISSUE_ID RELATED_ISSUE_ID",
		Short: "Relate two issues after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			relation, err := issueAdapterFor(runtime).CreateIssueRelation(
				ctx,
				client.IssueRelationCreateRequest{
					IssueID:        args[0],
					RelatedIssueID: args[1],
					Type:           relationType,
				},
			)
			if err != nil {
				return err
			}

			return writeIssueRelation(command, options, relation)
		},
	}
	command.Flags().StringVar(
		&relationType, "type", relationType,
		"relation type: blocks, duplicate, related, or similar",
	)
	root.AddCommand(command)
}

func addIssueUnrelateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "unrelate ISSUE_RELATION_ID",
		Short: "Delete an issue relation after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			id, err := issueAdapterFor(runtime).DeleteIssueRelation(ctx, args[0])
			if err != nil {
				return err
			}

			return writeDeletion(command, options, id)
		},
	})
}
