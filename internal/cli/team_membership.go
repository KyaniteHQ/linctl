package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addTeamMembershipCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.TeamMembershipList, client.TeamMembershipSummary]{
			Use:       "team-membership",
			Short:     "Read Linear team memberships",
			ListShort: "List visible team memberships",
			LimitHelp: "maximum team memberships to return",
			GetUse:    "get TEAM_MEMBERSHIP_ID",
			GetShort:  "Get one team membership by id",
			LoadList:  clientList(client.ListTeamMemberships),
			LoadGet:   clientGet(client.GetTeamMembershipByID),
			WriteItem: writeTeamMembership,
		},
	)
}

func writeTeamMembership(command *cobra.Command, options *rootOptions, membership client.TeamMembershipSummary) error {
	return writeItemLine(
		command, options, membership, membership.ID,
		"%s %s %s owner %t order %.2f",
		membership.ID,
		membership.TeamKey,
		membership.DisplayName,
		membership.Owner,
		membership.SortOrder,
	)
}
