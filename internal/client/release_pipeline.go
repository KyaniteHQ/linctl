package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Page
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
	Page
}

//nolint:lll
type releasePipelinesNode = gql.XReleasePipelinesReleasePipelinesReleasePipelineConnectionNodesReleasePipeline

//nolint:lll
type releasePipelineReleasesNode = gql.XReleasePipeline_releasesReleasePipelineReleasesReleaseConnectionNodesRelease

//nolint:lll
type releasePipelineStagesNode = gql.XReleasePipeline_stagesReleasePipelineStagesReleaseStageConnectionNodesReleaseStage

//nolint:lll
type releasePipelineTeamsNode = gql.XReleasePipeline_teamsReleasePipelineTeamsTeamConnectionNodesTeam

//nolint:lll
type releaseStagesNode = gql.XReleaseStagesReleaseStagesReleaseStageConnectionNodesReleaseStage

//nolint:lll
type releaseStageReleasesNode = gql.XReleaseStage_releasesReleaseStageReleasesReleaseConnectionNodesRelease

type releasePipelinesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type releasePipelineScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

type releaseStagesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type releaseStageScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// ListReleasePipelines returns visible Linear release pipelines.
func ListReleasePipelines(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleasePipelineList, error) {
	query := releasePipelinesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list release pipelines", limit, defaultListPageSize,
		query.page,
		releasePipelinesNodeSummary,
	)
	if err != nil {
		return ReleasePipelineList{}, err
	}

	return ReleasePipelineList{
		ReleasePipelines: page.Items,
		Page:             page.Page,
	}, nil
}

// GetReleasePipelineByID returns one Linear release pipeline by id.
func GetReleasePipelineByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ReleasePipelineSummary, error) {
	result, err := gql.XReleasePipeline(ctx, graphqlClient, id)
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
	query := releasePipelineScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release pipeline releases "+id, limit, defaultListPageSize,
		query.releases,
		releasePipelineReleasesNodeSummary,
	)
	if err != nil {
		return ReleaseList{}, err
	}

	return ReleaseList{Releases: page.Items, Page: page.Page}, nil
}

// ListReleasePipelineStages returns stages associated with one Linear release pipeline.
func ListReleasePipelineStages(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseStageList, error) {
	query := releasePipelineScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release pipeline stages "+id, limit, defaultListPageSize,
		query.stages,
		releasePipelineStagesNodeSummary,
	)
	if err != nil {
		return ReleaseStageList{}, err
	}

	return ReleaseStageList{ReleaseStages: page.Items, Page: page.Page}, nil
}

// ListReleasePipelineTeams returns teams associated with one Linear release pipeline.
func ListReleasePipelineTeams(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (TeamList, error) {
	query := releasePipelineScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release pipeline teams "+id, limit, defaultListPageSize,
		query.teams,
		releasePipelineTeamsNodeSummary,
	)
	if err != nil {
		return TeamList{}, err
	}

	return TeamList{Teams: page.Items, Page: page.Page}, nil
}

// ListReleaseStages returns visible Linear release stages.
func ListReleaseStages(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleaseStageList, error) {
	query := releaseStagesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list release stages", limit, defaultListPageSize,
		query.page,
		releaseStagesNodeSummary,
	)
	if err != nil {
		return ReleaseStageList{}, err
	}

	return ReleaseStageList{ReleaseStages: page.Items, Page: page.Page}, nil
}

// GetReleaseStageByID returns one Linear release stage by id.
func GetReleaseStageByID(ctx context.Context, graphqlClient graphql.Client, id string) (ReleaseStageSummary, error) {
	result, err := gql.XReleaseStage(ctx, graphqlClient, id)
	if err != nil {
		return ReleaseStageSummary{}, fmt.Errorf("get release stage %s: %w", id, err)
	}

	return releaseStageSummary(result.ReleaseStage.ReleaseStageSummaryFields), nil
}

// ListReleaseStageReleases returns releases associated with one Linear release stage.
func ListReleaseStageReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseList, error) {
	query := releaseStageScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release stage releases "+id, limit, defaultListPageSize,
		query.releases,
		releaseStageReleasesNodeSummary,
	)
	if err != nil {
		return ReleaseList{}, err
	}

	return ReleaseList{Releases: page.Items, Page: page.Page}, nil
}

func (query releasePipelinesQuery) page(
	pageSize int,
	after *string,
) ([]releasePipelinesNode, bool, *string, error) {
	result, err := gql.XReleasePipelines(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleasePipelines.Nodes,
		result.ReleasePipelines.PageInfo.HasNextPage,
		result.ReleasePipelines.PageInfo.EndCursor,
		nil
}

func (query releasePipelineScopedQuery) releases(
	pageSize int,
	after *string,
) ([]releasePipelineReleasesNode, bool, *string, error) {
	result, err := gql.XReleasePipeline_releases(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleasePipeline.Releases.Nodes,
		result.ReleasePipeline.Releases.PageInfo.HasNextPage,
		result.ReleasePipeline.Releases.PageInfo.EndCursor,
		nil
}

func (query releasePipelineScopedQuery) stages(
	pageSize int,
	after *string,
) ([]releasePipelineStagesNode, bool, *string, error) {
	result, err := gql.XReleasePipeline_stages(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleasePipeline.Stages.Nodes,
		result.ReleasePipeline.Stages.PageInfo.HasNextPage,
		result.ReleasePipeline.Stages.PageInfo.EndCursor,
		nil
}

func (query releasePipelineScopedQuery) teams(
	pageSize int,
	after *string,
) ([]releasePipelineTeamsNode, bool, *string, error) {
	result, err := gql.XReleasePipeline_teams(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleasePipeline.Teams.Nodes,
		result.ReleasePipeline.Teams.PageInfo.HasNextPage,
		result.ReleasePipeline.Teams.PageInfo.EndCursor,
		nil
}

func (query releaseStagesQuery) page(pageSize int, after *string) ([]releaseStagesNode, bool, *string, error) {
	result, err := gql.XReleaseStages(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleaseStages.Nodes,
		result.ReleaseStages.PageInfo.HasNextPage,
		result.ReleaseStages.PageInfo.EndCursor,
		nil
}

func (query releaseStageScopedQuery) releases(
	pageSize int,
	after *string,
) ([]releaseStageReleasesNode, bool, *string, error) {
	result, err := gql.XReleaseStage_releases(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleaseStage.Releases.Nodes,
		result.ReleaseStage.Releases.PageInfo.HasNextPage,
		result.ReleaseStage.Releases.PageInfo.EndCursor,
		nil
}

func releasePipelinesNodeSummary(node releasePipelinesNode) ReleasePipelineSummary {
	return releasePipelineSummary(node.ReleasePipelineSummaryFields)
}

func releasePipelineReleasesNodeSummary(node releasePipelineReleasesNode) ReleaseSummary {
	return releaseSummary(node.ReleaseSummaryFields)
}

func releasePipelineStagesNodeSummary(node releasePipelineStagesNode) ReleaseStageSummary {
	return releaseStageSummary(node.ReleaseStageSummaryFields)
}

func releasePipelineTeamsNodeSummary(node releasePipelineTeamsNode) TeamSummary {
	return teamSummary(node.TeamSummaryFields)
}

func releaseStagesNodeSummary(node releaseStagesNode) ReleaseStageSummary {
	return releaseStageSummary(node.ReleaseStageSummaryFields)
}

func releaseStageReleasesNodeSummary(node releaseStageReleasesNode) ReleaseSummary {
	return releaseSummary(node.ReleaseSummaryFields)
}

func releasePipelineSummary(fields gql.ReleasePipelineSummaryFields) ReleasePipelineSummary {
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

func releaseStageSummary(fields gql.ReleaseStageSummaryFields) ReleaseStageSummary {
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
