package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addUserCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	userCommand := &cobra.Command{
		Use:   "user",
		Short: "Read Linear users",
	}
	addUserListCommand(ctx, userCommand, options)
	addUserGetCommand(ctx, userCommand, options)
	addUserMeCommand(ctx, userCommand, options)
	root.AddCommand(userCommand)
}

func addUserListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List visible users",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(ctx, command, nil, options, limit, loadUserList, userPageWithItems, writeUser)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum users to return")
	root.AddCommand(command)
}

func addUserGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get USER_ID",
		Short: "Get one User by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			user, err := client.GetUserByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeUser(command, options, user)
		},
	})
}

func addUserMeCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "me",
		Short: "Get the authenticated User",
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

func writeUser(command *cobra.Command, options *rootOptions, user client.UserSummary) error {
	if wrote, err := writeIDOnly(command, options, user.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, user)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s <%s>", user.ID, user.DisplayName, user.Email)
}

func loadUserList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.UserList, []client.UserSummary, error) {
	users, err := client.ListUsers(ctx, runtime.graphqlClient, limit)
	return users, users.Users, err
}

func userPageWithItems(page client.UserList, users []client.UserSummary) client.UserList {
	page.Users = users
	return page
}
