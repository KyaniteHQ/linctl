package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addReleasePipelineCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleasePipelineList, client.ReleasePipelineSummary]{
			Use:       "release-pipeline",
			Short:     "Read Linear release pipelines",
			ListShort: "List visible Linear release pipelines",
			LimitHelp: "maximum release pipelines to return",
			GetUse:    "get RELEASE_PIPELINE_ID",
			GetShort:  "Get one release pipeline by id",
			LoadList:  clientList(client.ListReleasePipelines),
			LoadGet:   clientGet(client.GetReleasePipelineByID),
			WriteItem: writeReleasePipeline,
		},
	)
	addReleasePipelineReleasesCommand(ctx, command, options)
	addReleasePipelineStagesCommand(ctx, command, options)
	addReleasePipelineTeamsCommand(ctx, command, options)
}

func addReleasePipelineReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"releases RELEASE_PIPELINE_ID",
		"List releases associated with one Linear release pipeline",
		"releases",
		client.ListReleasePipelineReleases,
		writeRelease,
	)
}

func addReleasePipelineStagesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"stages RELEASE_PIPELINE_ID",
		"List stages associated with one Linear release pipeline",
		"release stages",
		client.ListReleasePipelineStages,
		writeReleaseStage,
	)
}

func addReleasePipelineTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"teams RELEASE_PIPELINE_ID",
		"List teams associated with one Linear release pipeline",
		"teams",
		client.ListReleasePipelineTeams,
		writeTeam,
	)
}

func addReleaseStageCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleaseStageList, client.ReleaseStageSummary]{
			Use:       "release-stage",
			Short:     "Read Linear release stages",
			ListShort: "List visible Linear release stages",
			LimitHelp: "maximum release stages to return",
			GetUse:    "get RELEASE_STAGE_ID",
			GetShort:  "Get one release stage by id",
			LoadList:  clientList(client.ListReleaseStages),
			LoadGet:   clientGet(client.GetReleaseStageByID),
			WriteItem: writeReleaseStage,
		},
	)
	addReleaseStageReleasesCommand(ctx, command, options)
}

func addReleaseStageReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"releases RELEASE_STAGE_ID",
		"List releases associated with one Linear release stage",
		"releases",
		client.ListReleaseStageReleases,
		writeRelease,
	)
}

func writeReleasePipeline(command *cobra.Command, options *rootOptions, pipeline client.ReleasePipelineSummary) error {
	return writeItemLine(
		command, options, pipeline, pipeline.ID,
		"%s %s %s releases %d", pipeline.ID, pipeline.Name, pipeline.SlugID, pipeline.ApproximateReleaseCount,
	)
}

func writeReleaseStage(command *cobra.Command, options *rootOptions, stage client.ReleaseStageSummary) error {
	return writeItemLine(
		command, options, stage, stage.ID,
		"%s %s [%s] pipeline %s", stage.ID, stage.Name, stage.Type, stage.PipelineName,
	)
}
