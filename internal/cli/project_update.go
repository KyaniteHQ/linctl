//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectUpdateReadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectUpdateCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ProjectUpdateList, client.ProjectUpdateSummary]{
			Use:           "project-update",
			Short:         "Read Linear project updates",
			ListShort:     "List visible project updates",
			LimitHelp:     "maximum project updates to return",
			GetUse:        "get PROJECT_UPDATE_ID",
			GetShort:      "Get one project update by id",
			LoadList:      loadProjectUpdateList,
			PageWithItems: projectUpdatePageWithItems,
			LoadGet:       loadProjectUpdate,
			WriteItem:     writeProjectUpdate,
		},
	)
	addProjectUpdateCommentsCommand(ctx, projectUpdateCommand, options)
}

func writeProjectUpdate(command *cobra.Command, options *rootOptions, update client.ProjectUpdateSummary) error {
	if wrote, err := writeIDOnly(command, options, update.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, update)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %s",
		update.ID,
		update.Health,
		update.DisplayName,
		update.Body,
	)
}

func loadProjectUpdateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectUpdateList, []client.ProjectUpdateSummary, error) {
	updates, err := client.ListAllProjectUpdates(ctx, runtime.graphqlClient, limit)
	return updates, updates.Updates, err
}

func loadProjectUpdate(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ProjectUpdateSummary, error) {
	return client.GetProjectUpdateByID(ctx, runtime.graphqlClient, id)
}

func projectUpdatePageWithItems(
	page client.ProjectUpdateList,
	updates []client.ProjectUpdateSummary,
) client.ProjectUpdateList {
	page.Updates = updates
	return page
}

func addProjectUpdateCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "comments PROJECT_UPDATE_ID",
		Short: "List project update comments without body content",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadProjectUpdateCommentList,
				projectUpdateCommentPageWithItems,
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum comments to return")
	root.AddCommand(command)
}

func loadProjectUpdateCommentList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectUpdateCommentList, []client.CommentMetadataSummary, error) {
	comments, err := client.ListProjectUpdateComments(ctx, runtime.graphqlClient, args[0], limit)
	return comments, comments.Comments, err
}

func projectUpdateCommentPageWithItems(
	page client.ProjectUpdateCommentList,
	comments []client.CommentMetadataSummary,
) client.ProjectUpdateCommentList {
	page.Comments = comments
	return page
}
