//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTeamMembershipCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.TeamMembershipList, client.TeamMembershipSummary]{
			Use:           "team-membership",
			Short:         "Read Linear team memberships",
			ListShort:     "List visible team memberships",
			LimitHelp:     "maximum team memberships to return",
			GetUse:        "get TEAM_MEMBERSHIP_ID",
			GetShort:      "Get one team membership by id",
			LoadList:      loadTeamMembershipList,
			PageWithItems: teamMembershipPageWithItems,
			LoadGet:       loadTeamMembership,
			WriteItem:     writeTeamMembership,
		},
	)
}

func writeTeamMembership(command *cobra.Command, options *rootOptions, membership client.TeamMembershipSummary) error {
	return writeItem(command, options, membership, membership.ID,
		func(command *cobra.Command, _ *rootOptions, membership client.TeamMembershipSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s owner %t order %.2f",
				membership.ID,
				membership.TeamKey,
				membership.DisplayName,
				membership.Owner,
				membership.SortOrder,
			)
		})
}

func loadTeamMembershipList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamMembershipList, []client.TeamMembershipSummary, error) {
	memberships, err := client.ListTeamMemberships(ctx, runtime.graphqlClient, limit)
	return memberships, memberships.Memberships, err
}

func loadTeamMembership(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.TeamMembershipSummary, error) {
	return client.GetTeamMembershipByID(ctx, runtime.graphqlClient, id)
}

func teamMembershipPageWithItems(
	page client.TeamMembershipList,
	memberships []client.TeamMembershipSummary,
) client.TeamMembershipList {
	page.Memberships = memberships
	return page
}
