package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

const orgWideTeamHelp = "required: a Team is organization-owned and is what a pin names, so a " +
	"create cannot land inside the pinned team; confirms this write adds a team to the organization"

func addTeamCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.TeamCreateRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.TeamSummary]{
		Use:   "create",
		Short: "Create a team with --org-wide, because a new team is always outside the pinned team",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Name, "name", "", "team name")
			command.Flags().StringVar(&request.Key, "key", "",
				"team key, which Linear derives from the name when you do not set it")
			command.Flags().StringVar(&request.Description, "description", "", "team description")
			command.Flags().BoolVar(&request.Private, "private", false, "create the team as a private team")
			command.Flags().StringVar(&request.ParentID, "parent", "",
				"parent team id, which makes the new team a sub-team")
			command.Flags().BoolVar(&request.Inherit, "inherit-workflow-states", false,
				"inherit workflow states from the parent team; needs --parent")
			command.Flags().BoolVar(&request.Triage, "triage", false, "enable triage on the new team")
			command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideTeamHelp)
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.TeamSummary, error) {
			return client.CreateTeam(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeTeam,
	})
}

const orgWideTeamDeleteHelp = "required: a Team is organization-owned and is what a pin names; confirms this " +
	"write archives a team in the organization and schedules its data for deletion"

func addTeamDeleteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.TeamDeleteRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[string]{
		Use:          "delete TEAM_ID",
		Short:        "Archive a team and schedule its deletion with --org-wide, which linctl cannot undo",
		Args:         cobra.ExactArgs(1),
		Irreversible: true,
		Configure: func(command *cobra.Command) {
			command.Flags().BoolVar(&request.OrgWide, "org-wide", false, orgWideTeamDeleteHelp)
		},
		Run: func(ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string) (string, error) {
			request.ID = args[0]

			return client.DeleteTeam(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeTeamDeletion,
	})
}

// writeTeamDeletion overrides the human deletion line with an explicit
// irreversibility warning: Linear archives the team and schedules its data for
// deletion, and there is no restore path via linctl.
func writeTeamDeletion(command *cobra.Command, options *rootOptions, id string) error {
	return writeDeletionMessage(
		command, options, id,
		"archived team "+id+" and scheduled its deletion: cannot be undone via linctl",
	)
}
