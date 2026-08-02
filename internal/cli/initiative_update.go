package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addInitiativeUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	initiativeUpdateCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeUpdateList, client.InitiativeUpdateSummary]{
			Use:       "initiative-update",
			Short:     "Read Linear initiative updates",
			ListShort: "List visible initiative updates",
			LimitHelp: "maximum initiative updates to return",
			GetUse:    "get INITIATIVE_UPDATE_ID",
			GetShort:  "Get one initiative update by id",
			LoadList:  clientList(client.ListInitiativeUpdates),
			LoadGet:   clientGet(client.GetInitiativeUpdateByID),
			WriteItem: writeInitiativeUpdate,
		},
	)
	addInitiativeUpdateCommentsCommand(ctx, initiativeUpdateCommand, options)
	addInitiativeUpdateCreateCommand(ctx, initiativeUpdateCommand, options)
}

func addInitiativeUpdateCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.InitiativeUpdateCreateRequest{}
	health := ""
	bodyFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.InitiativeUpdateSummary]{
		Use: "create INITIATIVE_ID",
		Short: "Create a status update on an initiative after organization comparison, " +
			"which needs no --org-wide flag",
		Args: cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(
				&request.Body, "body", "",
				"update body as markdown, or - to read the body from stdin",
			)
			command.Flags().StringVar(&bodyFile, "body-file", "", "read update body from file")
			command.Flags().StringVar(
				&health, "health", "", "initiative health: on-track, at-risk, or off-track",
			)
		},
		Run:   runInitiativeUpdateCreate(&request, &health, &bodyFile),
		Write: writeInitiativeUpdate,
	})
}

func runInitiativeUpdateCreate(
	request *client.InitiativeUpdateCreateRequest,
	health *string,
	bodyFile *string,
) func(context.Context, *cobra.Command, commandRuntime, []string) (client.InitiativeUpdateSummary, error) {
	return func(
		ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
	) (client.InitiativeUpdateSummary, error) {
		request.InitiativeID = args[0]
		if err := resolveBodyOrFileFlag(command, &request.Body, *bodyFile, "body"); err != nil {
			return client.InitiativeUpdateSummary{}, err
		}
		normalizedHealth, err := normalizeAndNote(command, "health", *health, normalizedHealthValue)
		if err != nil {
			return client.InitiativeUpdateSummary{}, err
		}
		request.Health = normalizedHealth

		return client.CreateInitiativeUpdate(ctx, runtime.graphqlClient, runtime.config.Target, *request)
	}
}

func addInitiativeUpdateCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"comments INITIATIVE_UPDATE_ID",
		"List initiative update comments without body content",
		"comments",
		client.ListInitiativeUpdateComments,
		writeCommentMetadata,
	)
}

func writeInitiativeUpdate(command *cobra.Command, options *rootOptions, update client.InitiativeUpdateSummary) error {
	return writeItemLine(
		command, options, update, update.ID,
		"%s %s %s %s", update.ID, update.Health, update.DisplayName, update.Body,
	)
}
