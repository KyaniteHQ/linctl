package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addTriageResponsibilityCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.TriageResponsibilityList, client.TriageResponsibilitySummary]{
			Use:       "triage-responsibility",
			Short:     "Read Linear triage responsibilities",
			ListShort: "List Linear triage responsibilities",
			LimitHelp: "maximum triage responsibilities to return",
			GetUse:    "get TRIAGE_RESPONSIBILITY_ID",
			GetShort:  "Get one triage responsibility by id",
			LoadList:  clientList(client.ListTriageResponsibilities),
			LoadGet:   clientGet(client.GetTriageResponsibilityByID),
			WriteItem: writeTriageResponsibility,
		},
	)
	addTriageResponsibilityManualSelectionCommand(ctx, command, options)
}

func addTriageResponsibilityManualSelectionCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.TriageResponsibilityManualSelection]{
		Use:   "manual-selection TRIAGE_RESPONSIBILITY_ID",
		Short: "Read manual user selection for one triage responsibility",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.TriageResponsibilityManualSelection, error) {
			return client.GetTriageResponsibilityManualSelection(ctx, runtime.graphqlClient, id)
		},
		Write: writeTriageResponsibilityManualSelection,
	})
}

func writeTriageResponsibility(
	command *cobra.Command,
	options *rootOptions,
	responsibility client.TriageResponsibilitySummary,
) error {
	return writeItemLine(
		command, options, responsibility, responsibility.ID,
		"%s team %s action %s current %s",
		responsibility.ID,
		responsibility.TeamKey,
		responsibility.Action,
		emptyDash(responsibility.CurrentUserName),
	)
}

func writeTriageResponsibilityManualSelection(
	command *cobra.Command,
	options *rootOptions,
	selection client.TriageResponsibilityManualSelection,
) error {
	return writeItemLine(
		command, options, selection, selection.ID,
		"%s manual users %s", selection.ID, emptyDash(strings.Join(selection.UserIDs, ",")),
	)
}
