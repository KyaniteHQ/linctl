package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addInitiativeCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeList, client.InitiativeSummary]{
			Use:       "initiative",
			Short:     "Read Linear initiatives",
			ListShort: "List visible initiatives",
			LimitHelp: "maximum initiatives to return",
			GetUse:    "get INITIATIVE_ID",
			GetShort:  "Get one initiative by id or slug",
			LoadList:  clientList(client.ListInitiatives),
			LoadGet:   clientGet(client.GetInitiativeByID),
			WriteItem: writeInitiative,
		},
	)
	addInitiativeHistoryCommand(ctx, command, options)
	addInitiativeLinksCommand(ctx, command, options)
	addSubInitiativesCommand(ctx, command, options)
	addInitiativeScopedUpdatesCommand(ctx, command, options)
	addInitiativeDocumentsCommand(ctx, command, options)
	addInitiativeProjectsCommand(ctx, command, options)
}

func addInitiativeHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"history INITIATIVE_ID",
		"List history records associated with one Linear initiative",
		"history records",
		client.ListInitiativeHistory,
		writeInitiativeHistory,
	)
}

func addInitiativeLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"links INITIATIVE_ID",
		"List external links associated with one Linear initiative",
		"links",
		client.ListInitiativeLinks,
		writeEntityExternalLink,
	)
}

func addSubInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"sub-initiatives INITIATIVE_ID",
		"List sub-initiatives associated with one Linear initiative",
		"sub-initiatives",
		client.ListSubInitiatives,
		writeInitiative,
	)
}

func addInitiativeScopedUpdatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"updates INITIATIVE_ID",
		"List status updates associated with one Linear initiative",
		"initiative updates",
		client.ListInitiativeUpdatesForInitiative,
		writeInitiativeUpdate,
	)
}

func addInitiativeDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"documents INITIATIVE_ID",
		"List documents associated with one Linear initiative",
		"documents",
		client.ListInitiativeDocuments,
		writeDocument,
	)
}

func addInitiativeProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"projects INITIATIVE_ID",
		"List projects directly associated with one Linear initiative",
		"projects",
		client.ListInitiativeProjects,
		writeProject,
	)
}

func writeInitiative(command *cobra.Command, options *rootOptions, initiative client.InitiativeSummary) error {
	return writeItemLine(
		command, options, initiative, initiative.ID,
		"%s %s [%s]", initiative.ID, initiative.Name, initiative.Status,
	)
}

func writeInitiativeHistory(
	command *cobra.Command,
	options *rootOptions,
	history client.InitiativeHistorySummary,
) error {
	return writeItemLine(
		command, options, history, history.ID,
		"%s initiative %s entries %d", history.ID, history.InitiativeID, history.EntryCount,
	)
}
