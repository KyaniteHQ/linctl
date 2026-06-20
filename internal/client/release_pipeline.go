package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ReleasePipelineSummary is the compact release pipeline model used by read-only commands.
type ReleasePipelineSummary struct {
	ID                                   string   `json:"id"`
	Name                                 string   `json:"name"`
	SlugID                               string   `json:"slug_id"`
	Type                                 string   `json:"type"`
	IsProduction                         bool     `json:"is_production"`
	AutoGenerateReleaseNotesOnCompletion bool     `json:"auto_generate_release_notes_on_completion"`
	IncludePathPatterns                  []string `json:"include_path_patterns,omitempty"`
	ApproximateReleaseCount              int      `json:"approximate_release_count"`
	Trashed                              bool     `json:"trashed,omitempty"`
	ReleaseNoteTemplateID                string   `json:"release_note_template_id,omitempty"`
	LatestReleaseNoteID                  string   `json:"latest_release_note_id,omitempty"`
	URL                                  string   `json:"url"`
	CreatedAt                            string   `json:"created_at"`
	UpdatedAt                            string   `json:"updated_at"`
	ArchivedAt                           string   `json:"archived_at,omitempty"`
}

// ReleasePipelineList is a page of Linear release pipelines.
type ReleasePipelineList struct {
	ReleasePipelines []ReleasePipelineSummary `json:"release_pipelines"`
	HasNextPage      bool                     `json:"has_next_page"`
	EndCursor        *string                  `json:"end_cursor,omitempty"`
}

// ReleaseStageSummary is the compact release stage model used by read-only commands.
type ReleaseStageSummary struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Color          string  `json:"color"`
	Type           string  `json:"type"`
	Position       float64 `json:"position"`
	Frozen         bool    `json:"frozen"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	ArchivedAt     string  `json:"archived_at,omitempty"`
	PipelineID     string  `json:"pipeline_id"`
	PipelineName   string  `json:"pipeline_name"`
	PipelineSlugID string  `json:"pipeline_slug_id"`
}

// ReleaseStageList is a page of Linear release stages.
type ReleaseStageList struct {
	ReleaseStages []ReleaseStageSummary `json:"release_stages"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

// ListReleasePipelines returns visible Linear release pipelines.
func ListReleasePipelines(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleasePipelineList, error) {
	result, err := releasePipelines(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleasePipelineList{}, fmt.Errorf("list release pipelines: %w", err)
	}

	summaries := make([]ReleasePipelineSummary, 0, len(result.ReleasePipelines.Nodes))
	for _, node := range result.ReleasePipelines.Nodes {
		summaries = append(summaries, releasePipelineSummary(node.ReleasePipelineSummaryFields))
	}

	return ReleasePipelineList{
		ReleasePipelines: summaries,
		HasNextPage:      result.ReleasePipelines.PageInfo.HasNextPage,
		EndCursor:        result.ReleasePipelines.PageInfo.EndCursor,
	}, nil
}

// GetReleasePipelineByID returns one Linear release pipeline by id.
func GetReleasePipelineByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ReleasePipelineSummary, error) {
	result, err := releasePipeline(ctx, graphqlClient, id)
	if err != nil {
		return ReleasePipelineSummary{}, fmt.Errorf("get release pipeline %s: %w", id, err)
	}

	return releasePipelineSummary(result.ReleasePipeline.ReleasePipelineSummaryFields), nil
}

// ListReleasePipelineReleases returns releases associated with one Linear release pipeline.
func ListReleasePipelineReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseList, error) {
	result, err := releasePipeline_releases(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list release pipeline releases %s: %w", id, err)
	}

	summaries := make([]ReleaseSummary, 0, len(result.ReleasePipeline.Releases.Nodes))
	for _, node := range result.ReleasePipeline.Releases.Nodes {
		summaries = append(summaries, releaseSummary(node.ReleaseSummaryFields))
	}

	return ReleaseList{
		Releases:    summaries,
		HasNextPage: result.ReleasePipeline.Releases.PageInfo.HasNextPage,
		EndCursor:   result.ReleasePipeline.Releases.PageInfo.EndCursor,
	}, nil
}

// ListReleasePipelineStages returns stages associated with one Linear release pipeline.
func ListReleasePipelineStages(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseStageList, error) {
	result, err := releasePipeline_stages(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseStageList{}, fmt.Errorf("list release pipeline stages %s: %w", id, err)
	}

	summaries := make([]ReleaseStageSummary, 0, len(result.ReleasePipeline.Stages.Nodes))
	for _, node := range result.ReleasePipeline.Stages.Nodes {
		summaries = append(summaries, releaseStageSummary(node.ReleaseStageSummaryFields))
	}

	return ReleaseStageList{
		ReleaseStages: summaries,
		HasNextPage:   result.ReleasePipeline.Stages.PageInfo.HasNextPage,
		EndCursor:     result.ReleasePipeline.Stages.PageInfo.EndCursor,
	}, nil
}

// ListReleaseStages returns visible Linear release stages.
func ListReleaseStages(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleaseStageList, error) {
	result, err := releaseStages(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseStageList{}, fmt.Errorf("list release stages: %w", err)
	}

	summaries := make([]ReleaseStageSummary, 0, len(result.ReleaseStages.Nodes))
	for _, node := range result.ReleaseStages.Nodes {
		summaries = append(summaries, releaseStageSummary(node.ReleaseStageSummaryFields))
	}

	return ReleaseStageList{
		ReleaseStages: summaries,
		HasNextPage:   result.ReleaseStages.PageInfo.HasNextPage,
		EndCursor:     result.ReleaseStages.PageInfo.EndCursor,
	}, nil
}

// GetReleaseStageByID returns one Linear release stage by id.
func GetReleaseStageByID(ctx context.Context, graphqlClient graphql.Client, id string) (ReleaseStageSummary, error) {
	result, err := releaseStage(ctx, graphqlClient, id)
	if err != nil {
		return ReleaseStageSummary{}, fmt.Errorf("get release stage %s: %w", id, err)
	}

	return releaseStageSummary(result.ReleaseStage.ReleaseStageSummaryFields), nil
}

func releasePipelineSummary(fields ReleasePipelineSummaryFields) ReleasePipelineSummary {
	summary := ReleasePipelineSummary{
		ID:                                   fields.Id,
		Name:                                 fields.Name,
		SlugID:                               fields.SlugId,
		Type:                                 string(fields.Type),
		IsProduction:                         fields.IsProduction,
		AutoGenerateReleaseNotesOnCompletion: fields.AutoGenerateReleaseNotesOnCompletion,
		IncludePathPatterns:                  fields.IncludePathPatterns,
		ApproximateReleaseCount:              fields.ApproximateReleaseCount,
		Trashed:                              boolValue(fields.Trashed),
		URL:                                  fields.Url,
		CreatedAt:                            fields.CreatedAt,
		UpdatedAt:                            fields.UpdatedAt,
		ArchivedAt:                           stringValue(fields.ArchivedAt),
	}
	if fields.ReleaseNoteTemplate != nil {
		summary.ReleaseNoteTemplateID = fields.ReleaseNoteTemplate.Id
	}
	if fields.LatestReleaseNote != nil {
		summary.LatestReleaseNoteID = fields.LatestReleaseNote.Id
	}

	return summary
}

func boolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func releaseStageSummary(fields ReleaseStageSummaryFields) ReleaseStageSummary {
	return ReleaseStageSummary{
		ID:             fields.Id,
		Name:           fields.Name,
		Color:          fields.Color,
		Type:           string(fields.Type),
		Position:       fields.Position,
		Frozen:         fields.Frozen,
		CreatedAt:      fields.CreatedAt,
		UpdatedAt:      fields.UpdatedAt,
		ArchivedAt:     stringValue(fields.ArchivedAt),
		PipelineID:     fields.Pipeline.Id,
		PipelineName:   fields.Pipeline.Name,
		PipelineSlugID: fields.Pipeline.SlugId,
	}
}
