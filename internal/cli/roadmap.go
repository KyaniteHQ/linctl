package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addRoadmapCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	roadmapCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.RoadmapList, client.RoadmapSummary]{
			Use:       "roadmap",
			Short:     "Read legacy Linear roadmaps, which the initiative commands replace for new planning",
			ListShort: "List visible legacy Linear roadmaps",
			LimitHelp: "maximum legacy roadmaps to return",
			GetUse:    "get ROADMAP_ID",
			GetShort:  "Get one legacy roadmap by id",
			LoadList:  clientList(client.ListRoadmaps),
			LoadGet:   clientGet(client.GetRoadmapByID),
			WriteItem: writeRoadmap,
		},
	)
	roadmapCommand.Long = "Roadmap is Linear's deprecated planning surface. " +
		"These reads remain for compatibility; use `linctl initiative` for new planning workflows."
	addRoadmapProjectsCommand(ctx, roadmapCommand, options)
}

func addRoadmapProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.RoadmapProjectList, client.ProjectSummary]{
		Use:   "projects ROADMAP_ID",
		Short: "List projects associated with one legacy roadmap",
		Long: "List projects associated with one legacy roadmap. " +
			"Roadmap is Linear's deprecated planning surface. " +
			"Use `linctl initiative projects` for new planning workflows.",
		LimitHelp: "projects",
		Args:      cobra.ExactArgs(1),
		Load:      loadRoadmapProjects,
		WriteItem: writeProject,
	})
}

func writeRoadmap(command *cobra.Command, options *rootOptions, roadmap client.RoadmapSummary) error {
	return writeItemLine(
		command, options, roadmap, roadmap.ID,
		"%s %s %s [legacy]", roadmap.ID, roadmap.Name, roadmap.SlugID,
	)
}

func loadRoadmapProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.RoadmapProjectList, error) {
	projects, err := client.ListRoadmapProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, err
}
