package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	labelCommand := &cobra.Command{
		Use:   "label",
		Short: "Read Linear issue labels",
	}
	addLabelListCommand(ctx, labelCommand, options)
	addLabelGetCommand(ctx, labelCommand, options)
	root.AddCommand(labelCommand)
}

func addLabelListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List visible labels",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(ctx, command, nil, options, limit, loadLabelList, labelPageWithItems, writeLabel)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum labels to return")
	root.AddCommand(command)
}

func addLabelGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get LABEL_ID",
		Short: "Get one label by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			label, err := client.GetLabelByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeLabel(command, options, label)
		},
	})
}

func writeLabel(command *cobra.Command, options *rootOptions, label client.LabelSummary) error {
	if wrote, err := writeIDOnly(command, options, label.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, label)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s %s", label.ID, label.Name, label.Color)
}

func loadLabelList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.LabelList, []client.LabelSummary, error) {
	labels, err := client.ListLabels(ctx, runtime.graphqlClient, limit)
	return labels, labels.Labels, err
}

func labelPageWithItems(page client.LabelList, labels []client.LabelSummary) client.LabelList {
	page.Labels = labels
	return page
}
