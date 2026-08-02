package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addUserCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	userCommand := newGroupCommand("user", "Read Linear users")
	addUserListCommand(ctx, userCommand, options)
	addUserGetCommand(ctx, userCommand, options)
	addUserMeCommand(ctx, userCommand, options)
	addUserDraftsCommand(ctx, userCommand, options)
	addUserSettingsCommand(ctx, userCommand, options)
	addUserAssignedIssuesCommand(ctx, userCommand, options)
	addUserCreatedIssuesCommand(ctx, userCommand, options)
	addUserDelegatedIssuesCommand(ctx, userCommand, options)
	addUserTeamMembershipsCommand(ctx, userCommand, options)
	addUserTeamsCommand(ctx, userCommand, options)
	addViewerAssignedIssuesCommand(ctx, userCommand, options)
	addViewerCreatedIssuesCommand(ctx, userCommand, options)
	addViewerDelegatedIssuesCommand(ctx, userCommand, options)
	addViewerTeamMembershipsCommand(ctx, userCommand, options)
	addViewerTeamsCommand(ctx, userCommand, options)
	root.AddCommand(userCommand)
}

func addUserListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.UserList, client.UserSummary]{
		Use:       "list",
		Short:     "List visible users",
		LimitHelp: "users",
		Args:      cobra.NoArgs,
		Load:      loadUserList,
		WriteItem: writeUser,
	})
}

func addUserGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.UserSummary]{
		Use:   "get USER_ID",
		Short: "Get one user by id",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.UserSummary, error) {
			return client.GetUserByID(ctx, runtime.graphqlClient, id)
		},
		Write: writeUser,
	})
}

func addUserMeCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
		Use:   "me",
		Short: "Get the authenticated user",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			user, err := client.GetViewerUser(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUser(command, options, user)
		},
	})
}

func addUserDraftsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.DraftList, client.DraftSummary]{
		Use:       "drafts",
		Short:     "List the saved draft metadata of the authenticated user",
		LimitHelp: "drafts",
		Args:      cobra.NoArgs,
		Load:      loadViewerDraftList,
		WriteItem: writeDraft,
	})
}

func addUserAssignedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "assigned-issues USER_ID",
		Short:     "List the issues assigned to a user",
		LimitHelp: "issues",
		Args:      cobra.ExactArgs(1),
		Load:      loadUserAssignedIssues,
		WriteItem: writeIssue,
	})
}

func addUserCreatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "created-issues USER_ID",
		Short:     "List the issues created by a user",
		LimitHelp: "issues",
		Args:      cobra.ExactArgs(1),
		Load:      loadUserCreatedIssues,
		WriteItem: writeIssue,
	})
}

func addUserDelegatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "delegated-issues USER_ID",
		Short:     "List the issues delegated to a user",
		LimitHelp: "issues",
		Args:      cobra.ExactArgs(1),
		Load:      loadUserDelegatedIssues,
		WriteItem: writeIssue,
	})
}

func addUserTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamMembershipList, client.TeamMembershipSummary]{
		Use:       "team-memberships USER_ID",
		Short:     "List the team memberships of a user",
		LimitHelp: "memberships",
		Args:      cobra.ExactArgs(1),
		Load:      loadUserTeamMemberships,
		WriteItem: writeTeamMembership,
	})
}

func addUserTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamList, client.TeamSummary]{
		Use:       "teams USER_ID",
		Short:     "List the teams of a user",
		LimitHelp: "teams",
		Args:      cobra.ExactArgs(1),
		Load:      loadUserTeams,
		WriteItem: writeTeam,
	})
}

func addViewerAssignedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "my-assigned-issues",
		Short:     "List the issues assigned to the authenticated user",
		LimitHelp: "issues",
		Args:      cobra.NoArgs,
		Load:      loadViewerAssignedIssues,
		WriteItem: writeIssue,
	})
}

func addViewerCreatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "my-created-issues",
		Short:     "List the issues created by the authenticated user",
		LimitHelp: "issues",
		Args:      cobra.NoArgs,
		Load:      loadViewerCreatedIssues,
		WriteItem: writeIssue,
	})
}

func addViewerDelegatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "my-delegated-issues",
		Short:     "List the issues delegated to the authenticated user",
		LimitHelp: "issues",
		Args:      cobra.NoArgs,
		Load:      loadViewerDelegatedIssues,
		WriteItem: writeIssue,
	})
}

func addViewerTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamMembershipList, client.TeamMembershipSummary]{
		Use:       "my-team-memberships",
		Short:     "List the team memberships of the authenticated user",
		LimitHelp: "memberships",
		Args:      cobra.NoArgs,
		Load:      loadViewerTeamMemberships,
		WriteItem: writeTeamMembership,
	})
}

func addViewerTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamList, client.TeamSummary]{
		Use:       "my-teams",
		Short:     "List the teams of the authenticated user",
		LimitHelp: "teams",
		Args:      cobra.NoArgs,
		Load:      loadViewerTeams,
		WriteItem: writeTeam,
	})
}

func writeUser(command *cobra.Command, options *rootOptions, user client.UserSummary) error {
	return writeItemLine(command, options, user, user.ID, "%s %s <%s>", user.ID, user.DisplayName, user.Email)
}

func writeDraft(command *cobra.Command, options *rootOptions, draft client.DraftSummary) error {
	return writeItem(command, options, draft, draft.ID,
		func(command *cobra.Command, _ *rootOptions, draft client.DraftSummary) error {
			parentKey := defaultString(draft.ParentKey, "-")
			parentTitle := defaultString(draft.ParentTitle, "-")
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %s",
				draft.ID,
				draft.ParentType,
				parentKey,
				parentTitle,
			)
		})
}

func loadUserList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.UserList, error) {
	users, err := client.ListUsers(ctx, runtime.graphqlClient, limit)
	return users, err
}

func loadViewerDraftList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.DraftList, error) {
	drafts, err := client.ListViewerDrafts(ctx, runtime.graphqlClient, limit)
	return drafts, err
}

func loadUserAssignedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListUserAssignedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, err
}

func loadUserCreatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListUserCreatedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, err
}

func loadUserDelegatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListUserDelegatedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, err
}

func loadUserTeamMemberships(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMembershipList, error) {
	memberships, err := client.ListUserTeamMemberships(ctx, runtime.graphqlClient, args[0], limit)
	return memberships, err
}

func loadUserTeams(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamList, error) {
	teams, err := client.ListUserTeams(ctx, runtime.graphqlClient, args[0], limit)
	return teams, err
}

func loadViewerAssignedIssues(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListViewerAssignedIssues(ctx, runtime.graphqlClient, limit)
	return issues, err
}

func loadViewerCreatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListViewerCreatedIssues(ctx, runtime.graphqlClient, limit)
	return issues, err
}

func loadViewerDelegatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.ListViewerDelegatedIssues(ctx, runtime.graphqlClient, limit)
	return issues, err
}

func loadViewerTeamMemberships(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamMembershipList, error) {
	memberships, err := client.ListViewerTeamMemberships(ctx, runtime.graphqlClient, limit)
	return memberships, err
}

func loadViewerTeams(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, error) {
	teams, err := client.ListViewerTeams(ctx, runtime.graphqlClient, limit)
	return teams, err
}
