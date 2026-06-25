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
	addLabelChildrenCommand(ctx, labelCommand, options)
	addLabelIssuesCommand(ctx, labelCommand, options)
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

func addLabelChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "children LABEL_ID",
		Short: "List child labels under one label group",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadLabelChildren,
				labelChildrenPageWithItems,
				writeLabel,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum labels to return")
	root.AddCommand(command)
}

func addLabelIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "issues LABEL_ID",
		Short: "List Issues associated with one label",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadLabelIssues,
				labelIssuesPageWithItems,
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum Issues to return")
	root.AddCommand(command)
}

func writeLabel(command *cobra.Command, options *rootOptions, label client.LabelSummary) error {
	return writeItem(command, options, label, label.ID,
		func(command *cobra.Command, _ *rootOptions, label client.LabelSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s %s", label.ID, label.Name, label.Color)
		})
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

func loadLabelChildren(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.LabelChildList, []client.LabelSummary, error) {
	labels, err := client.ListLabelChildren(ctx, runtime.graphqlClient, args[0], limit)
	return labels, labels.Labels, err
}

func labelChildrenPageWithItems(page client.LabelChildList, labels []client.LabelSummary) client.LabelChildList {
	page.Labels = labels
	return page
}

func loadLabelIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.LabelIssueList, []client.IssueSummary, error) {
	issues, err := client.ListLabelIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func labelIssuesPageWithItems(page client.LabelIssueList, issues []client.IssueSummary) client.LabelIssueList {
	page.Issues = issues
	return page
}
