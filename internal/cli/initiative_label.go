//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeLabelList, client.InitiativeLabelSummary]{
			Use:       "initiative-label",
			Short:     "Read Linear initiative labels",
			ListShort: "List visible Linear initiative labels",
			LimitHelp: "maximum initiative labels to return",
			GetUse:    "get INITIATIVE_LABEL_ID",
			GetShort:  "Get one initiative label by id",
			LoadList:  loadInitiativeLabelList,
			LoadGet:   loadInitiativeLabel,
			WriteItem: writeInitiativeLabel,
		},
	)
}

func writeInitiativeLabel(command *cobra.Command, options *rootOptions, label client.InitiativeLabelSummary) error {
	return writeItem(command, options, label, label.ID,
		func(command *cobra.Command, options *rootOptions, label client.InitiativeLabelSummary) error {
			format, err := normalizedHumanFormat(options)
			if err != nil {
				return err
			}
			if format == "minimal" {
				return render.WriteLine(command.OutOrStdout(), "%s", label.ID)
			}
			if format == "full" {
				return render.WriteLine(
					command.OutOrStdout(),
					"%s %s %s group=%t parent=%s",
					label.ID,
					label.Name,
					label.Color,
					label.IsGroup,
					emptyDash(label.ParentName),
				)
			}

			return render.WriteLine(command.OutOrStdout(), "%s %s %s", label.ID, label.Name, label.Color)
		})
}

func loadInitiativeLabelList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeLabelList, []client.InitiativeLabelSummary, error) {
	labels, err := client.ListInitiativeLabels(ctx, runtime.graphqlClient, limit)
	return labels, labels.InitiativeLabels, err
}

func loadInitiativeLabel(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeLabelSummary, error) {
	return client.GetInitiativeLabelByID(ctx, runtime.graphqlClient, id)
}
