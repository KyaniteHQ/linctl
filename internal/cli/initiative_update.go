package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
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
			LoadList:  loadInitiativeUpdateList,
			LoadGet:   loadInitiativeUpdate,
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
		Short: "Post a status update to an initiative after organization comparison " +
			"(Resource-Scoped; no --org-wide flag)",
		Args: cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Body, "body", "", "update body as markdown; use - to read stdin")
			command.Flags().StringVar(&bodyFile, "body-file", "", "read update body from file")
			command.Flags().StringVar(
				&health, "health", "", "initiative health: on-track, at-risk, or off-track (Linear enum aliases)",
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
		if err := resolveFileFlag(command, &request.Body, *bodyFile, "body"); err != nil {
			return client.InitiativeUpdateSummary{}, err
		}
		if err := resolveBodyFlag(command, &request.Body); err != nil {
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
	limit := 50
	command := &cobra.Command{
		Use:   "comments INITIATIVE_UPDATE_ID",
		Short: "List initiative update comments without body content",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeUpdateCommentList,
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum comments to return")
	root.AddCommand(preflightReadListCommand(command, loadInitiativeUpdateCommentList))
}

func writeInitiativeUpdate(command *cobra.Command, options *rootOptions, update client.InitiativeUpdateSummary) error {
	return writeItem(command, options, update, update.ID,
		func(command *cobra.Command, _ *rootOptions, update client.InitiativeUpdateSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %s",
				update.ID,
				update.Health,
				update.DisplayName,
				update.Body,
			)
		})
}

func loadInitiativeUpdateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeUpdateList, []client.InitiativeUpdateSummary, error) {
	updates, err := client.ListInitiativeUpdates(ctx, runtime.graphqlClient, limit)
	return updates, updates.Updates, err
}

func loadInitiativeUpdate(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeUpdateSummary, error) {
	return client.GetInitiativeUpdateByID(ctx, runtime.graphqlClient, id)
}

func loadInitiativeUpdateCommentList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.InitiativeUpdateCommentList, []client.CommentMetadataSummary, error) {
	comments, err := client.ListInitiativeUpdateComments(ctx, runtime.graphqlClient, args[0], limit)
	return comments, comments.Comments, err
}
