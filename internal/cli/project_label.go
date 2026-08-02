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
			Use:       "project-label",
			Short:     "Read and write Linear project labels",
			ListShort: "List visible Linear project labels",
			LimitHelp: "maximum project labels to return",
			GetUse:    "get PROJECT_LABEL_ID",
			GetShort:  "Get one project label by id",
			LoadList:  clientList(client.ListProjectLabels),
			LoadGet:   clientGet(client.GetProjectLabelByID),
			WriteItem: writeProjectLabel,
		},
	)
	addProjectLabelChildrenCommand(ctx, parentCommand, options)
	addProjectLabelProjectsCommand(ctx, parentCommand, options)
	addProjectLabelCreateCommand(ctx, parentCommand, options)
	addProjectLabelUpdateCommand(ctx, parentCommand, options)
	addProjectLabelRetireCommand(ctx, parentCommand, options)
	addProjectLabelRestoreCommand(ctx, parentCommand, options)
}

func addProjectLabelChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children PROJECT_LABEL_ID",
		"List child labels for one project label",
		"child project labels",
		client.ListProjectLabelChildren,
		writeProjectLabel,
	)
}

func addProjectLabelProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"projects PROJECT_LABEL_ID",
		"List projects associated with one project label",
		"projects",
		client.ListProjectLabelProjects,
		writeProject,
	)
}

func writeProjectLabel(command *cobra.Command, options *rootOptions, label client.ProjectLabelSummary) error {
	return writeItem(command, options, label, label.ID,
		func(command *cobra.Command, options *rootOptions, label client.ProjectLabelSummary) error {
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
