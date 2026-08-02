package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addAgentSkillCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.AgentSkillList, client.AgentSkillSummary]{
			Use:       "agent-skill",
			Short:     "Read Linear AgentSkills",
			ListShort: "List Linear AgentSkills",
			LimitHelp: "maximum AgentSkills to return",
			GetUse:    "get AGENT_SKILL_ID",
			GetShort:  "Get one AgentSkill by id",
			LoadList:  clientList(client.ListAgentSkills),
			LoadGet:   clientGet(client.GetAgentSkillByID),
			WriteItem: writeAgentSkill,
		},
	)
}

func writeAgentSkill(command *cobra.Command, options *rootOptions, skill client.AgentSkillSummary) error {
	return writeItemLine(
		command, options, skill, skill.ID,
		"%s %s shared %t recent %.0f", skill.ID, skill.Title, skill.Shared, skill.RecentUsageCount,
	)
}
