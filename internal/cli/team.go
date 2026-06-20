package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTeamCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	teamCommand := &cobra.Command{
		Use:   "team",
		Short: "Read Linear teams",
	}
	addTeamListCommand(ctx, teamCommand, options)
	addTeamGetCommand(ctx, teamCommand, options)
	addTeamMembersCommand(ctx, teamCommand, options)
	root.AddCommand(teamCommand)
}

func addTeamListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List visible teams",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(ctx, command, nil, options, limit, loadTeamList, teamPageWithItems, writeTeam)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum teams to return")
	root.AddCommand(command)
}

func addTeamGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get TEAM_ID",
		Short: "Get one Team by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			team, err := client.GetTeamByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeTeam(command, options, team)
		},
	})
}

func addTeamMembersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "members TEAM_ID",
		Short: "List team members",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadTeamMemberList,
				teamMemberPageWithItems,
				writeUser,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum members to return")
	root.AddCommand(command)
}

func writeTeam(command *cobra.Command, options *rootOptions, team client.TeamSummary) error {
	if wrote, err := writeIDOnly(command, options, team.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, team)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s %s", team.ID, team.Key, team.Name)
}

func loadTeamList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, []client.TeamSummary, error) {
	teams, err := client.ListTeams(ctx, runtime.graphqlClient, limit)
	return teams, teams.Teams, err
}

func teamPageWithItems(page client.TeamList, teams []client.TeamSummary) client.TeamList {
	page.Teams = teams
	return page
}

func loadTeamMemberList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMemberList, []client.UserSummary, error) {
	members, err := client.ListTeamMembers(ctx, runtime.graphqlClient, args[0], limit)
	return members, members.Members, err
}

func teamMemberPageWithItems(page client.TeamMemberList, members []client.UserSummary) client.TeamMemberList {
	page.Members = members
	return page
}
