//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addReleasePipelineCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleasePipelineList, client.ReleasePipelineSummary]{
			Use:           "release-pipeline",
			Short:         "Read Linear release pipelines",
			ListShort:     "List visible Linear release pipelines",
			LimitHelp:     "maximum release pipelines to return",
			GetUse:        "get RELEASE_PIPELINE_ID",
			GetShort:      "Get one release pipeline by id",
			LoadList:      loadReleasePipelineList,
			PageWithItems: releasePipelinePageWithItems,
			LoadGet:       loadReleasePipeline,
			WriteItem:     writeReleasePipeline,
		},
	)
	addReleasePipelineReleasesCommand(ctx, command, options)
	addReleasePipelineStagesCommand(ctx, command, options)
}

func addReleasePipelineReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "releases RELEASE_PIPELINE_ID",
		Short: "List releases associated with one Linear release pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadReleasePipelineReleases,
				releasePageWithItems,
				writeRelease,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum releases to return")
	root.AddCommand(command)
}

func addReleasePipelineStagesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "stages RELEASE_PIPELINE_ID",
		Short: "List stages associated with one Linear release pipeline",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadReleasePipelineStages,
				releaseStagePageWithItems,
				writeReleaseStage,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum release stages to return")
	root.AddCommand(command)
}

func addReleaseStageCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleaseStageList, client.ReleaseStageSummary]{
			Use:           "release-stage",
			Short:         "Read Linear release stages",
			ListShort:     "List visible Linear release stages",
			LimitHelp:     "maximum release stages to return",
			GetUse:        "get RELEASE_STAGE_ID",
			GetShort:      "Get one release stage by id",
			LoadList:      loadReleaseStageList,
			PageWithItems: releaseStagePageWithItems,
			LoadGet:       loadReleaseStage,
			WriteItem:     writeReleaseStage,
		},
	)
	addReleaseStageReleasesCommand(ctx, command, options)
}

func addReleaseStageReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "releases RELEASE_STAGE_ID",
		Short: "List releases associated with one Linear release stage",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadReleaseStageReleases,
				releasePageWithItems,
				writeRelease,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum releases to return")
	root.AddCommand(command)
}

func writeReleasePipeline(
	command *cobra.Command,
	options *rootOptions,
	pipeline client.ReleasePipelineSummary,
) error {
	if wrote, err := writeIDOnly(command, options, pipeline.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, pipeline)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s releases %d",
		pipeline.ID,
		pipeline.Name,
		pipeline.SlugID,
		pipeline.ApproximateReleaseCount,
	)
}

func writeReleaseStage(
	command *cobra.Command,
	options *rootOptions,
	stage client.ReleaseStageSummary,
) error {
	if wrote, err := writeIDOnly(command, options, stage.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, stage)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s [%s] pipeline %s",
		stage.ID,
		stage.Name,
		stage.Type,
		stage.PipelineName,
	)
}

func loadReleasePipelineList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ReleasePipelineList, []client.ReleasePipelineSummary, error) {
	pipelines, err := client.ListReleasePipelines(ctx, runtime.graphqlClient, limit)
	return pipelines, pipelines.ReleasePipelines, err
}

func loadReleasePipeline(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ReleasePipelineSummary, error) {
	return client.GetReleasePipelineByID(ctx, runtime.graphqlClient, id)
}

func releasePipelinePageWithItems(
	page client.ReleasePipelineList,
	pipelines []client.ReleasePipelineSummary,
) client.ReleasePipelineList {
	page.ReleasePipelines = pipelines
	return page
}

func loadReleasePipelineReleases(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ReleaseList, []client.ReleaseSummary, error) {
	releases, err := client.ListReleasePipelineReleases(ctx, runtime.graphqlClient, args[0], limit)
	return releases, releases.Releases, err
}

func loadReleasePipelineStages(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ReleaseStageList, []client.ReleaseStageSummary, error) {
	stages, err := client.ListReleasePipelineStages(ctx, runtime.graphqlClient, args[0], limit)
	return stages, stages.ReleaseStages, err
}

func loadReleaseStageList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ReleaseStageList, []client.ReleaseStageSummary, error) {
	stages, err := client.ListReleaseStages(ctx, runtime.graphqlClient, limit)
	return stages, stages.ReleaseStages, err
}

func loadReleaseStage(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ReleaseStageSummary, error) {
	return client.GetReleaseStageByID(ctx, runtime.graphqlClient, id)
}

func releaseStagePageWithItems(
	page client.ReleaseStageList,
	stages []client.ReleaseStageSummary,
) client.ReleaseStageList {
	page.ReleaseStages = stages
	return page
}

func loadReleaseStageReleases(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ReleaseList, []client.ReleaseSummary, error) {
	releases, err := client.ListReleaseStageReleases(ctx, runtime.graphqlClient, args[0], limit)
	return releases, releases.Releases, err
}
