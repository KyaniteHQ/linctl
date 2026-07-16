package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/gitctx"
)

func addCurrentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "current",
		Short: "Resolve the Linear issue for the current checkout",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			issueID, err := gitctx.CurrentIssueIdentifierForTeam(ctx, ".", pinnedTeamKeyHint(ctx, options))
			if err != nil {
				return err
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.GetIssueByID(ctx, runtime.graphqlClient, issueID)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	}
	addCommandWithSafety(root, CommandSafetyRead, command)
}

// pinnedTeamKeyHint returns the Pinned Target's team key as a best-effort
// parse filter for Current Issue resolution. It never errors: when config is
// absent, unreadable, or has no team key configured, it returns "", which
// keeps Current Issue resolution unfiltered. Any real config problem still
// surfaces moments later from buildCommandRuntime for commands that need it.
func pinnedTeamKeyHint(ctx context.Context, options *rootOptions) string {
	resolvedConfig, err := resolveConfig(ctx, options)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(resolvedConfig.Target.TeamKey)
}

func addDoneCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyWrite, &cobra.Command{
		Use:   "done",
		Short: "Close the current checkout issue",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			issueID, err := gitctx.CurrentIssueIdentifierForTeam(ctx, ".", pinnedTeamKeyHint(ctx, options))
			if err != nil {
				return err
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.CloseIssue(ctx, runtime.graphqlClient, runtime.config.Target, issueID)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}
