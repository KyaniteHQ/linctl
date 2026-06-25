//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addRoadmapToProjectCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand[
		client.RoadmapToProjectList,
		client.RoadmapToProjectSummary,
	](ctx, root, options, readListGetSpec[client.RoadmapToProjectList, client.RoadmapToProjectSummary]{
		Use:           "roadmap-to-project",
		Short:         "Read legacy Linear Roadmap-to-Project associations",
		ListShort:     "List visible legacy Roadmap-to-Project associations",
		LimitHelp:     "maximum legacy Roadmap-to-Project associations to return",
		GetUse:        "get ROADMAP_TO_PROJECT_ID",
		GetShort:      "Get one legacy Roadmap-to-Project association by id",
		LoadList:      loadRoadmapToProjectList,
		PageWithItems: roadmapToProjectPageWithItems,
		LoadGet:       loadRoadmapToProject,
		WriteItem:     writeRoadmapToProject,
	})
	command.Long = "RoadmapToProject is Linear's deprecated planning association surface. " +
		"These reads remain for compatibility; use `linctl initiative-to-project` for new workflows."
}

func writeRoadmapToProject(
	command *cobra.Command,
	options *rootOptions,
	association client.RoadmapToProjectSummary,
) error {
	return writeItem(command, options, association, association.ID,
		func(command *cobra.Command, _ *rootOptions, association client.RoadmapToProjectSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s -> %s order %s [legacy]",
				association.ID,
				association.RoadmapName,
				association.ProjectName,
				association.SortOrder,
			)
		})
}

func loadRoadmapToProjectList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.RoadmapToProjectList, []client.RoadmapToProjectSummary, error) {
	associations, err := client.ListRoadmapToProjects(ctx, runtime.graphqlClient, limit)
	return associations, associations.Associations, err
}

func loadRoadmapToProject(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.RoadmapToProjectSummary, error) {
	return client.GetRoadmapToProjectByID(ctx, runtime.graphqlClient, id)
}

func roadmapToProjectPageWithItems(
	page client.RoadmapToProjectList,
	associations []client.RoadmapToProjectSummary,
) client.RoadmapToProjectList {
	page.Associations = associations
	return page
}
