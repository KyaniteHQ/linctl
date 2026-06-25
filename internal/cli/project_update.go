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
	addProjectUpdateCreateCommand(ctx, projectUpdateCommand, options)
}

func addProjectUpdateCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.ProjectUpdateCreateRequest{}
	health := ""
	bodyFile := ""
	command := &cobra.Command{
		Use:   "create PROJECT_ID",
		Short: "Post a status update to a project after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ProjectID = args[0]

			return runProjectUpdateCreate(
				ctx, command, options, issueAdapterFor(runtime), request, health, bodyFile,
			)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "update body as markdown; use - to read stdin")
	command.Flags().StringVar(&bodyFile, "body-file", "", "read update body from file")
	command.Flags().StringVar(&health, "health", "", "project health: on-track, at-risk, or off-track")
	root.AddCommand(command)
}

// runProjectUpdateCreate resolves the body (stdin or file), normalizes the
// health alias, then posts the update through the Command Port. Splitting it
// from the cobra wiring makes the port the test surface: the body and health
// logic is exercised against an in-memory fake, not canned GraphQL JSON.
func runProjectUpdateCreate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator projectUpdateCreator,
	request client.ProjectUpdateCreateRequest,
	health string,
	bodyFile string,
) error {
	if err := resolveFileFlag(&request.Body, bodyFile, "body"); err != nil {
		return err
	}
	if err := resolveBodyFlag(command, &request.Body); err != nil {
		return err
	}
	normalizedHealth, err := normalizeAndNote(command, "health", health, normalizedHealthValue)
	if err != nil {
		return err
	}
	request.Health = normalizedHealth
	update, err := creator.CreateProjectUpdate(ctx, request)
	if err != nil {
		return err
	}

	return writeProjectUpdate(command, options, update)
}

func writeProjectUpdate(command *cobra.Command, options *rootOptions, update client.ProjectUpdateSummary) error {
	return writeItem(command, options, update, update.ID,
		func(command *cobra.Command, _ *rootOptions, update client.ProjectUpdateSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %s",
				update.ID,
				update.Health,
				update.DisplayName,
				update.Body,
			)
		})
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
