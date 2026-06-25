//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addAgentActivityCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.AgentActivityList, client.AgentActivitySummary]{
		Use:           "agent-activity",
		Short:         "Read Linear AgentActivities",
		ListShort:     "List Linear AgentActivities",
		LimitHelp:     "maximum AgentActivities to return",
		GetUse:        "get AGENT_ACTIVITY_ID",
		GetShort:      "Get one AgentActivity by id",
		LoadList:      loadAgentActivityList,
		PageWithItems: agentActivityPageWithItems,
		LoadGet:       loadAgentActivity,
		WriteItem:     writeAgentActivity,
	})
}

func writeAgentActivity(command *cobra.Command, options *rootOptions, activity client.AgentActivitySummary) error {
	return writeItem(command, options, activity, activity.ID,
		func(command *cobra.Command, _ *rootOptions, activity client.AgentActivitySummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s session %s [%s] signal %s",
				activity.ID,
				activity.AgentSessionID,
				activity.ContentType,
				emptyDash(activity.Signal),
			)
		})
}

func loadAgentActivityList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.AgentActivityList, []client.AgentActivitySummary, error) {
	activities, err := client.ListAgentActivities(ctx, runtime.graphqlClient, limit)
	return activities, activities.AgentActivities, err
}

func loadAgentActivity(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.AgentActivitySummary, error) {
	return client.GetAgentActivityByID(ctx, runtime.graphqlClient, id)
}

func agentActivityPageWithItems(
	page client.AgentActivityList,
	activities []client.AgentActivitySummary,
) client.AgentActivityList {
	page.AgentActivities = activities
	return page
}
