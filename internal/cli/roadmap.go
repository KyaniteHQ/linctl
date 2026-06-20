//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addRoadmapCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(ctx, root, options, readListGetSpec[client.RoadmapList, client.RoadmapSummary]{
		Use:           "roadmap",
		Short:         "Read Linear roadmaps",
		ListShort:     "List visible Linear roadmaps",
		LimitHelp:     "maximum roadmaps to return",
		GetUse:        "get ROADMAP_ID",
		GetShort:      "Get one roadmap by id",
		LoadList:      loadRoadmapList,
		PageWithItems: roadmapPageWithItems,
		LoadGet:       loadRoadmap,
		WriteItem:     writeRoadmap,
	})
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

	return render.WriteLine(command.OutOrStdout(), "%s %s %s", roadmap.ID, roadmap.Name, roadmap.SlugID)
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
