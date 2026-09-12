package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addAgentActivityCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.AgentActivityCreateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.AgentActivitySummary]{
		Use:   "create AGENT_SESSION_ID",
		Short: "Emit an AgentActivity into an agent session whose issue is on the pinned target",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Type, "type", "",
				"activity type: "+strings.Join(client.AgentActivityTypes, ", "))
			command.Flags().StringVar(&request.Body, "body", "",
				"markdown body for "+strings.Join(client.AgentActivityBodyTypes, ", ")+" activities")
			command.Flags().StringVar(&request.Action, "action", "", "action name for action activities")
			command.Flags().StringVar(&request.Parameter, "parameter", "", "action parameter for action activities")
			command.Flags().StringVar(&request.Result, "result", "", "action result for action activities")
			command.Flags().StringVar(&request.Signal, "signal", "",
				"optional signal: "+strings.Join(client.AgentActivitySignals, ", "))
			command.Flags().BoolVar(&request.Ephemeral, "ephemeral", false,
				"mark the activity ephemeral so it disappears after the next activity")
			registerFlagCompletion(command, "type",
				cobra.FixedCompletions(client.AgentActivityTypes, cobra.ShellCompDirectiveNoFileComp))
			registerFlagCompletion(command, "signal",
				cobra.FixedCompletions(client.AgentActivitySignals, cobra.ShellCompDirectiveNoFileComp))
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.AgentActivitySummary, error) {
			request.AgentSessionID = args[0]

			return client.CreateAgentActivity(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeAgentActivity,
	})
}
