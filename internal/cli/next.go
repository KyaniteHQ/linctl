package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addNextCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	dryRun := false
	limit := 20
	command := &cobra.Command{
		Use:   "next --dry-run",
		Short: "Preview the next unblocked issue",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if !dryRun {
				return errors.New("next requires --dry-run because checkout/worktree creation is not implemented")
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			target, err := runtime.resolveTarget(ctx)
			if err != nil {
				return err
			}
			issues, err := client.ListNextIssuesByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(issues.Issues)); err != nil {
				return err
			}
			if len(issues.Issues) == 0 {
				return errors.New("next issue not found")
			}

			return writeIssue(command, options, issues.Issues[0])
		},
	}
	command.Flags().BoolVar(&dryRun, "dry-run", dryRun, "preview without creating a checkout or worktree")
	command.Flags().IntVar(&limit, "limit", limit, "maximum candidate issues to inspect")
	root.AddCommand(command)
}
