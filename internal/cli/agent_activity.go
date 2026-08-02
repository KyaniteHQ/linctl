package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addAgentActivityCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.AgentActivityList, client.AgentActivitySummary]{
			Use:       "agent-activity",
			Short:     "Read Linear AgentActivities",
			ListShort: "List Linear AgentActivities",
			LimitHelp: "maximum AgentActivities to return",
			GetUse:    "get AGENT_ACTIVITY_ID",
			GetShort:  "Get one AgentActivity by id",
			LoadList:  clientList(client.ListAgentActivities),
			LoadGet:   clientGet(client.GetAgentActivityByID),
			WriteItem: writeAgentActivity,
		},
	)
}

func writeAgentActivity(command *cobra.Command, options *rootOptions, activity client.AgentActivitySummary) error {
	return writeItemLine(
		command, options, activity, activity.ID,
		"%s session %s [%s] signal %s",
		activity.ID,
		activity.AgentSessionID,
		activity.ContentType,
		emptyDash(activity.Signal),
	)
}
