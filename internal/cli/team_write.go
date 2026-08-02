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
