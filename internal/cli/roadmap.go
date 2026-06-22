//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addRoadmapCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	roadmapCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.RoadmapList, client.RoadmapSummary]{
			Use:           "roadmap",
			Short:         "Read legacy Linear roadmaps; prefer initiative for new planning",
			ListShort:     "List visible legacy Linear roadmaps",
			LimitHelp:     "maximum legacy roadmaps to return",
			GetUse:        "get ROADMAP_ID",
			GetShort:      "Get one legacy roadmap by id",
			LoadList:      loadRoadmapList,
			PageWithItems: roadmapPageWithItems,
			LoadGet:       loadRoadmap,
			WriteItem:     writeRoadmap,
		},
	)
	roadmapCommand.Long = "Roadmap is Linear's deprecated planning surface. " +
		"These reads remain for compatibility; use `linctl initiative` for new planning workflows."
	addRoadmapProjectsCommand(ctx, roadmapCommand, options)
}

func addRoadmapProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "projects ROADMAP_ID",
		Short: "List projects associated with one legacy roadmap",
		Long: "List projects associated with one legacy roadmap. " +
			"Roadmap is Linear's deprecated planning surface. " +
			"Use `linctl initiative projects` for new planning workflows.",
		Args: cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadRoadmapProjects,
				roadmapProjectPageWithItems,
				writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func writeRoadmap(
	command *cobra.Command,
	options *rootOptions,
	roadmap client.RoadmapSummary,
) error {
	if wrote, err := writeIDOnly(command, options, roadmap.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, roadmap)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s %s [legacy]", roadmap.ID, roadmap.Name, roadmap.SlugID)
}

func loadRoadmapList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.RoadmapList, []client.RoadmapSummary, error) {
	roadmaps, err := client.ListRoadmaps(ctx, runtime.graphqlClient, limit)
	return roadmaps, roadmaps.Roadmaps, err
}

func loadRoadmap(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.RoadmapSummary, error) {
	return client.GetRoadmapByID(ctx, runtime.graphqlClient, id)
}

func roadmapPageWithItems(
	page client.RoadmapList,
	roadmaps []client.RoadmapSummary,
) client.RoadmapList {
	page.Roadmaps = roadmaps
	return page
}

func loadRoadmapProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.RoadmapProjectList, []client.ProjectSummary, error) {
	projects, err := client.ListRoadmapProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, projects.Projects, err
}

func roadmapProjectPageWithItems(
	page client.RoadmapProjectList,
	projects []client.ProjectSummary,
) client.RoadmapProjectList {
	page.Projects = projects
	return page
}
