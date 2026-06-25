//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addAgentSkillCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.AgentSkillList, client.AgentSkillSummary]{
		Use:           "agent-skill",
		Short:         "Read Linear AgentSkills",
		ListShort:     "List Linear AgentSkills",
		LimitHelp:     "maximum AgentSkills to return",
		GetUse:        "get AGENT_SKILL_ID",
		GetShort:      "Get one AgentSkill by id",
		LoadList:      loadAgentSkillList,
		PageWithItems: agentSkillPageWithItems,
		LoadGet:       loadAgentSkill,
		WriteItem:     writeAgentSkill,
	})
}

func writeAgentSkill(command *cobra.Command, options *rootOptions, skill client.AgentSkillSummary) error {
	return writeItem(command, options, skill, skill.ID,
		func(command *cobra.Command, _ *rootOptions, skill client.AgentSkillSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s shared %t recent %.0f",
				skill.ID,
				skill.Title,
				skill.Shared,
				skill.RecentUsageCount,
			)
		})
}

func loadAgentSkillList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.AgentSkillList, []client.AgentSkillSummary, error) {
	skills, err := client.ListAgentSkills(ctx, runtime.graphqlClient, limit)
	return skills, skills.AgentSkills, err
}

func loadAgentSkill(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.AgentSkillSummary, error) {
	return client.GetAgentSkillByID(ctx, runtime.graphqlClient, id)
}

func agentSkillPageWithItems(
	page client.AgentSkillList,
	skills []client.AgentSkillSummary,
) client.AgentSkillList {
	page.AgentSkills = skills
	return page
}
