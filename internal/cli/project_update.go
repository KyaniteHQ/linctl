package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectUpdateReadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectUpdateCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ProjectUpdateList, client.ProjectUpdateSummary]{
			Use:       "project-update",
			Short:     "Read Linear project updates",
			ListShort: "List visible project updates",
			LimitHelp: "maximum project updates to return",
			GetUse:    "get PROJECT_UPDATE_ID",
			GetShort:  "Get one project update by id",
			LoadList:  clientList(client.ListAllProjectUpdates),
			LoadGet:   clientGet(client.GetProjectUpdateByID),
			WriteItem: writeProjectUpdate,
		},
	)
	addProjectUpdateCommentsCommand(ctx, projectUpdateCommand, options)
	addProjectUpdateCreateCommand(ctx, projectUpdateCommand, options)
}

func addProjectUpdateCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectUpdateCreateRequest{}
	health := ""
	bodyFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.ProjectUpdateSummary]{
		Use:   "create PROJECT_ID",
		Short: "Create a status update on a project after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(
				&request.Body, "body", "",
				"update body as markdown, or - to read the body from stdin",
			)
			command.Flags().StringVar(&bodyFile, "body-file", "", "read update body from file")
			command.Flags().StringVar(&health, "health", "", "project health: on-track, at-risk, or off-track")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.ProjectUpdateSummary, error) {
			request.ProjectID = args[0]

			if err := resolveBodyOrFileFlag(command, &request.Body, bodyFile, "body"); err != nil {
				return client.ProjectUpdateSummary{}, err
			}
			normalizedHealth, err := normalizeAndNote(command, "health", health, normalizedHealthValue)
			if err != nil {
				return client.ProjectUpdateSummary{}, err
			}
			request.Health = normalizedHealth

			return client.CreateProjectUpdate(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeProjectUpdate,
	})
}

func writeProjectUpdate(command *cobra.Command, options *rootOptions, update client.ProjectUpdateSummary) error {
	return writeItemLine(
		command, options, update, update.ID,
		"%s %s %s %s", update.ID, update.Health, update.DisplayName, update.Body,
	)
}

func addProjectUpdateCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"comments PROJECT_UPDATE_ID",
		"List project update comments without body content",
		"comments",
		client.ListProjectUpdateComments,
		writeCommentMetadata,
	)
}
