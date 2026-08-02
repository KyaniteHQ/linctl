package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	labelCommand := newGroupCommand("label", "Read and write Linear issue labels")
	addLabelListCommand(ctx, labelCommand, options)
	addLabelGetCommand(ctx, labelCommand, options)
	addLabelChildrenCommand(ctx, labelCommand, options)
	addLabelIssuesCommand(ctx, labelCommand, options)
	addLabelCreateCommand(ctx, labelCommand, options)
	addLabelUpdateCommand(ctx, labelCommand, options)
	addLabelRetireCommand(ctx, labelCommand, options)
	addLabelRestoreCommand(ctx, labelCommand, options)
	root.AddCommand(labelCommand)
}

func addLabelListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.LabelList, client.LabelSummary]{
		Use:       "list",
		Short:     "List visible labels",
		LimitHelp: "labels",
		Args:      cobra.NoArgs,
		Load:      loadLabelList,
		WriteItem: writeLabel,
	})
}

func addLabelGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.LabelSummary]{
		Use:   "get LABEL_ID",
		Short: "Get one label by id",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.LabelSummary, error) {
			return client.GetLabelByID(ctx, runtime.graphqlClient, id)
		},
		Write: writeLabel,
	})
}

func addLabelChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children LABEL_ID",
		"List child labels under one label group",
		"labels",
		client.ListLabelChildren,
		writeLabel,
	)
}

func addLabelIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"issues LABEL_ID",
		"List Issues associated with one label",
		"Issues",
		client.ListLabelIssues,
		writeIssue,
	)
}

func writeLabel(command *cobra.Command, options *rootOptions, label client.LabelSummary) error {
	return writeItemLine(command, options, label, label.ID, "%s %s %s", label.ID, label.Name, label.Color)
}

func loadLabelList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.LabelList, error) {
	labels, err := client.ListLabels(ctx, runtime.graphqlClient, limit)
	return labels, err
}
