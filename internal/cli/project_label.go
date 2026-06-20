//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectLabelCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	parentCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ProjectLabelList, client.ProjectLabelSummary]{
			Use:           "project-label",
			Short:         "Read Linear project labels",
			ListShort:     "List visible Linear project labels",
			LimitHelp:     "maximum project labels to return",
			GetUse:        "get PROJECT_LABEL_ID",
			GetShort:      "Get one project label by id",
			LoadList:      loadProjectLabelList,
			PageWithItems: projectLabelPageWithItems,
			LoadGet:       loadProjectLabel,
			WriteItem:     writeProjectLabel,
		},
	)
	addProjectLabelChildrenCommand(ctx, parentCommand, options)
	addProjectLabelProjectsCommand(ctx, parentCommand, options)
}

func addProjectLabelChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "children PROJECT_LABEL_ID",
		Short: "List child labels for one project label",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadProjectLabelChildrenList,
				projectLabelChildrenPageWithItems,
				writeProjectLabel,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum child project labels to return")
	root.AddCommand(command)
}

func addProjectLabelProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "projects PROJECT_LABEL_ID",
		Short: "List projects associated with one project label",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadProjectLabelProjectsList,
				projectLabelProjectsPageWithItems,
				writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func writeProjectLabel(command *cobra.Command, options *rootOptions, label client.ProjectLabelSummary) error {
	if wrote, err := writeIDOnly(command, options, label.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, label)
	}

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
}

func loadProjectLabelList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectLabelList, []client.ProjectLabelSummary, error) {
	labels, err := client.ListProjectLabels(ctx, runtime.graphqlClient, limit)
	return labels, labels.ProjectLabels, err
}

func loadProjectLabel(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ProjectLabelSummary, error) {
	return client.GetProjectLabelByID(ctx, runtime.graphqlClient, id)
}

func loadProjectLabelChildrenList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectLabelChildrenList, []client.ProjectLabelSummary, error) {
	labels, err := client.ListProjectLabelChildren(ctx, runtime.graphqlClient, args[0], limit)
	return labels, labels.ProjectLabels, err
}

func loadProjectLabelProjectsList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectLabelProjectsList, []client.ProjectSummary, error) {
	projects, err := client.ListProjectLabelProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, projects.Projects, err
}

func projectLabelPageWithItems(
	page client.ProjectLabelList,
	labels []client.ProjectLabelSummary,
) client.ProjectLabelList {
	page.ProjectLabels = labels
	return page
}

func projectLabelChildrenPageWithItems(
	page client.ProjectLabelChildrenList,
	labels []client.ProjectLabelSummary,
) client.ProjectLabelChildrenList {
	page.ProjectLabels = labels
	return page
}

func projectLabelProjectsPageWithItems(
	page client.ProjectLabelProjectsList,
	projects []client.ProjectSummary,
) client.ProjectLabelProjectsList {
	page.Projects = projects
	return page
}
