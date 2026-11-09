package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addAgentSessionCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var issue string
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.AgentSessionSummary]{
		Use:   "create",
		Short: "Start an AgentSession on an issue that is on the pinned target",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&issue, "issue", "", "issue id or identifier to start the session on")
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.AgentSessionSummary, error) {
			return client.CreateAgentSessionOnIssue(ctx, runtime.graphqlClient, runtime.config.Target, issue)
		},
		Write: writeAgentSession,
	})
}
