//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addTriageResponsibilityCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.TriageResponsibilityList, client.TriageResponsibilitySummary]{
			Use:           "triage-responsibility",
			Short:         "Read Linear triage responsibilities",
			ListShort:     "List Linear triage responsibilities",
			LimitHelp:     "maximum triage responsibilities to return",
			GetUse:        "get TRIAGE_RESPONSIBILITY_ID",
			GetShort:      "Get one triage responsibility by id",
			LoadList:      loadTriageResponsibilityList,
			PageWithItems: triageResponsibilityPageWithItems,
			LoadGet:       loadTriageResponsibility,
			WriteItem:     writeTriageResponsibility,
		},
	)
	addTriageResponsibilityManualSelectionCommand(ctx, command, options)
}

func addTriageResponsibilityManualSelectionCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
) {
	root.AddCommand(&cobra.Command{
		Use:   "manual-selection TRIAGE_RESPONSIBILITY_ID",
		Short: "Read manual user selection for one triage responsibility",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			selection, err := client.GetTriageResponsibilityManualSelection(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeTriageResponsibilityManualSelection(command, options, selection)
		},
	})
}

func writeTriageResponsibility(
	command *cobra.Command,
	options *rootOptions,
	responsibility client.TriageResponsibilitySummary,
) error {
	return writeItem(command, options, responsibility, responsibility.ID,
		func(command *cobra.Command, _ *rootOptions, responsibility client.TriageResponsibilitySummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s team %s action %s current %s",
				responsibility.ID,
				responsibility.TeamKey,
				responsibility.Action,
				emptyDash(responsibility.CurrentUserName),
			)
		})
}

func writeTriageResponsibilityManualSelection(
	command *cobra.Command,
	options *rootOptions,
	selection client.TriageResponsibilityManualSelection,
) error {
	return writeItem(command, options, selection, selection.ID,
		func(command *cobra.Command, _ *rootOptions, selection client.TriageResponsibilityManualSelection) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s manual users %s",
				selection.ID,
				emptyDash(strings.Join(selection.UserIDs, ",")),
			)
		})
}

func loadTriageResponsibilityList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TriageResponsibilityList, []client.TriageResponsibilitySummary, error) {
	responsibilities, err := client.ListTriageResponsibilities(ctx, runtime.graphqlClient, limit)
	return responsibilities, responsibilities.TriageResponsibilities, err
}

func loadTriageResponsibility(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.TriageResponsibilitySummary, error) {
	return client.GetTriageResponsibilityByID(ctx, runtime.graphqlClient, id)
}

func triageResponsibilityPageWithItems(
	page client.TriageResponsibilityList,
	responsibilities []client.TriageResponsibilitySummary,
) client.TriageResponsibilityList {
	page.TriageResponsibilities = responsibilities
	return page
}
