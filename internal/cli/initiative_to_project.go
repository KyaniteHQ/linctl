//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeToProjectCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand[
		client.InitiativeToProjectList,
		client.InitiativeToProjectSummary,
	](ctx, root, options, readListGetSpec[client.InitiativeToProjectList, client.InitiativeToProjectSummary]{
		Use:           "initiative-to-project",
		Short:         "Read Linear Initiative-to-Project associations",
		ListShort:     "List visible Initiative-to-Project associations",
		LimitHelp:     "maximum Initiative-to-Project associations to return",
		GetUse:        "get INITIATIVE_TO_PROJECT_ID",
		GetShort:      "Get one Initiative-to-Project association by id",
		LoadList:      loadInitiativeToProjectList,
		PageWithItems: initiativeToProjectPageWithItems,
		LoadGet:       loadInitiativeToProject,
		WriteItem:     writeInitiativeToProject,
	})
}

func writeInitiativeToProject(
	command *cobra.Command,
	options *rootOptions,
	association client.InitiativeToProjectSummary,
) error {
	if wrote, err := writeIDOnly(command, options, association.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, association)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s -> %s order %s",
		association.ID,
		association.InitiativeName,
		association.ProjectName,
		association.SortOrder,
	)
}

func loadInitiativeToProjectList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeToProjectList, []client.InitiativeToProjectSummary, error) {
	associations, err := client.ListInitiativeToProjects(ctx, runtime.graphqlClient, limit)
	return associations, associations.Associations, err
}

func loadInitiativeToProject(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeToProjectSummary, error) {
	return client.GetInitiativeToProjectByID(ctx, runtime.graphqlClient, id)
}

func initiativeToProjectPageWithItems(
	page client.InitiativeToProjectList,
	associations []client.InitiativeToProjectSummary,
) client.InitiativeToProjectList {
	page.Associations = associations
	return page
}
