package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addExternalUserCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ExternalUserList, client.ExternalUserSummary]{
			Use:       "external-user",
			Short:     "Read Linear ExternalUsers",
			ListShort: "List Linear ExternalUsers",
			LimitHelp: "maximum ExternalUsers to return",
			GetUse:    "get EXTERNAL_USER_ID",
			GetShort:  "Get one ExternalUser by id",
			LoadList:  clientList(client.ListExternalUsers),
			LoadGet:   clientGet(client.GetExternalUserByID),
			WriteItem: writeExternalUser,
		},
	)
}

func writeExternalUser(command *cobra.Command, options *rootOptions, user client.ExternalUserSummary) error {
	return writeItemLine(
		command, options, user, user.ID,
		"%s %s %s last_seen %s", user.ID, user.Name, user.DisplayName, emptyDash(user.LastSeen),
	)
}
