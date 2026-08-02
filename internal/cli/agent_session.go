package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addAgentSessionCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.AgentSessionList, client.AgentSessionSummary]{
			Use:       "agent-session",
			Short:     "Read Linear AgentSessions",
			ListShort: "List Linear AgentSessions",
			LimitHelp: "maximum AgentSessions to return",
			GetUse:    "get AGENT_SESSION_ID",
			GetShort:  "Get one AgentSession by id",
			LoadList:  clientList(client.ListAgentSessions),
			LoadGet:   clientGet(client.GetAgentSessionByID),
			WriteItem: writeAgentSession,
		},
	)
}

func writeAgentSession(command *cobra.Command, options *rootOptions, session client.AgentSessionSummary) error {
	return writeItemLine(
		command, options, session, session.ID,
		"%s %s [%s] %s", session.ID, session.SlugID, session.Status, emptyDash(session.IssueIdentifier),
	)
}
