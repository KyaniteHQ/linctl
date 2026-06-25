//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addAgentSessionCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.AgentSessionList, client.AgentSessionSummary]{
		Use:           "agent-session",
		Short:         "Read Linear AgentSessions",
		ListShort:     "List Linear AgentSessions",
		LimitHelp:     "maximum AgentSessions to return",
		GetUse:        "get AGENT_SESSION_ID",
		GetShort:      "Get one AgentSession by id",
		LoadList:      loadAgentSessionList,
		PageWithItems: agentSessionPageWithItems,
		LoadGet:       loadAgentSession,
		WriteItem:     writeAgentSession,
	})
}

func writeAgentSession(
	command *cobra.Command,
	options *rootOptions,
	session client.AgentSessionSummary,
) error {
	if wrote, err := writeIDOnly(command, options, session.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, session)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s [%s] %s",
		session.ID,
		session.SlugID,
		session.Status,
		emptyDash(session.IssueIdentifier),
	)
}

func loadAgentSessionList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.AgentSessionList, []client.AgentSessionSummary, error) {
	sessions, err := client.ListAgentSessions(ctx, runtime.graphqlClient, limit)
	return sessions, sessions.AgentSessions, err
}

func loadAgentSession(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.AgentSessionSummary, error) {
	return client.GetAgentSessionByID(ctx, runtime.graphqlClient, id)
}

func agentSessionPageWithItems(
	page client.AgentSessionList,
	sessions []client.AgentSessionSummary,
) client.AgentSessionList {
	page.AgentSessions = sessions
	return page
}
