//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addExternalUserCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.ExternalUserList, client.ExternalUserSummary]{
		Use:           "external-user",
		Short:         "Read Linear ExternalUsers",
		ListShort:     "List Linear ExternalUsers",
		LimitHelp:     "maximum ExternalUsers to return",
		GetUse:        "get EXTERNAL_USER_ID",
		GetShort:      "Get one ExternalUser by id",
		LoadList:      loadExternalUserList,
		PageWithItems: externalUserPageWithItems,
		LoadGet:       loadExternalUser,
		WriteItem:     writeExternalUser,
	})
}

func writeExternalUser(command *cobra.Command, options *rootOptions, user client.ExternalUserSummary) error {
	return writeItem(command, options, user, user.ID,
		func(command *cobra.Command, _ *rootOptions, user client.ExternalUserSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s last_seen %s",
				user.ID,
				user.Name,
				user.DisplayName,
				emptyDash(user.LastSeen),
			)
		})
}

func loadExternalUserList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ExternalUserList, []client.ExternalUserSummary, error) {
	users, err := client.ListExternalUsers(ctx, runtime.graphqlClient, limit)
	return users, users.ExternalUsers, err
}

func loadExternalUser(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ExternalUserSummary, error) {
	return client.GetExternalUserByID(ctx, runtime.graphqlClient, id)
}

func externalUserPageWithItems(
	page client.ExternalUserList,
	users []client.ExternalUserSummary,
) client.ExternalUserList {
	page.ExternalUsers = users
	return page
}
