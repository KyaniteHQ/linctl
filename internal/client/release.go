package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ReleaseSummary is the compact release model used by read-only commands.
type ReleaseSummary struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	SlugID           string `json:"slug_id"`
	Version          string `json:"version,omitempty"`
	Description      string `json:"description,omitempty"`
	CommitSHA        string `json:"commit_sha,omitempty"`
	IssueCount       int    `json:"issue_count"`
	ReleaseNoteCount int    `json:"release_note_count"`
	Trashed          bool   `json:"trashed"`
	URL              string `json:"url"`
	StartDate        string `json:"start_date,omitempty"`
	TargetDate       string `json:"target_date,omitempty"`
	StartedAt        string `json:"started_at,omitempty"`
	CompletedAt      string `json:"completed_at,omitempty"`
	CanceledAt       string `json:"canceled_at,omitempty"`
	AutoArchivedAt   string `json:"auto_archived_at,omitempty"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	ArchivedAt       string `json:"archived_at,omitempty"`
	PipelineID       string `json:"pipeline_id"`
	PipelineName     string `json:"pipeline_name"`
	PipelineSlugID   string `json:"pipeline_slug_id"`
	StageID          string `json:"stage_id"`
	StageName        string `json:"stage_name"`
	StageType        string `json:"stage_type"`
	CreatorID        string `json:"creator_id,omitempty"`
	CreatorName      string `json:"creator_name,omitempty"`
}

// ReleaseList is a page of Linear releases.
type ReleaseList struct {
	Releases    []ReleaseSummary `json:"releases"`
	HasNextPage bool             `json:"has_next_page"`
	EndCursor   *string          `json:"end_cursor,omitempty"`
}

// ReleaseHistorySummary is the compact release history model used by read-only commands.
type ReleaseHistorySummary struct {
	ID         string          `json:"id"`
	ReleaseID  string          `json:"release_id"`
	EntryCount int             `json:"entry_count"`
	Entries    json.RawMessage `json:"entries"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
	ArchivedAt string          `json:"archived_at,omitempty"`
}

// ReleaseHistoryList is a page of Linear release history records.
type ReleaseHistoryList struct {
	History     []ReleaseHistorySummary `json:"history"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// EntityExternalLinkSummary is the compact external link model used by read-only commands.
type EntityExternalLinkSummary struct {
	ID             string  `json:"id"`
	Label          string  `json:"label"`
	URL            string  `json:"url"`
	SortOrder      float64 `json:"sort_order"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	ArchivedAt     string  `json:"archived_at,omitempty"`
	CreatorID      string  `json:"creator_id,omitempty"`
	CreatorName    string  `json:"creator_name,omitempty"`
	InitiativeID   string  `json:"initiative_id,omitempty"`
	InitiativeName string  `json:"initiative_name,omitempty"`
	ProjectID      string  `json:"project_id,omitempty"`
	ProjectName    string  `json:"project_name,omitempty"`
}

// EntityExternalLinkList is a page of Linear external links.
type EntityExternalLinkList struct {
	Links       []EntityExternalLinkSummary `json:"links"`
	HasNextPage bool                        `json:"has_next_page"`
	EndCursor   *string                     `json:"end_cursor,omitempty"`
}

// ReleaseNoteSummary is the compact release note model used by read-only commands.
type ReleaseNoteSummary struct {
	ID                  string `json:"id"`
	Title               string `json:"title,omitempty"`
	SlugID              string `json:"slug_id"`
	GenerationStatus    string `json:"generation_status,omitempty"`
	ReleaseCount        int    `json:"release_count"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	ArchivedAt          string `json:"archived_at,omitempty"`
	PipelineID          string `json:"pipeline_id"`
	PipelineName        string `json:"pipeline_name"`
	PipelineSlugID      string `json:"pipeline_slug_id"`
	FirstReleaseID      string `json:"first_release_id,omitempty"`
	FirstReleaseName    string `json:"first_release_name,omitempty"`
	FirstReleaseVersion string `json:"first_release_version,omitempty"`
	LastReleaseID       string `json:"last_release_id,omitempty"`
	LastReleaseName     string `json:"last_release_name,omitempty"`
	LastReleaseVersion  string `json:"last_release_version,omitempty"`
}

// ReleaseNoteList is a page of Linear release notes.
type ReleaseNoteList struct {
	ReleaseNotes []ReleaseNoteSummary `json:"release_notes"`
	HasNextPage  bool                 `json:"has_next_page"`
	EndCursor    *string              `json:"end_cursor,omitempty"`
}

// ListReleases returns visible Linear releases.
func ListReleases(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleaseList, error) {
	result, err := releases(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list releases: %w", err)
	}

	summaries := make([]ReleaseSummary, 0, len(result.Releases.Nodes))
	for _, node := range result.Releases.Nodes {
		summaries = append(summaries, releaseSummary(node.ReleaseSummaryFields))
	}

	return ReleaseList{
		Releases:    summaries,
		HasNextPage: result.Releases.PageInfo.HasNextPage,
		EndCursor:   result.Releases.PageInfo.EndCursor,
	}, nil
}

// GetReleaseByID returns one Linear release by id.
func GetReleaseByID(ctx context.Context, graphqlClient graphql.Client, id string) (ReleaseSummary, error) {
	result, err := release(ctx, graphqlClient, id)
	if err != nil {
		return ReleaseSummary{}, fmt.Errorf("get release %s: %w", id, err)
	}

	return releaseSummary(result.Release.ReleaseSummaryFields), nil
}

// ListReleaseHistory returns history records associated with one Linear release.
func ListReleaseHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseHistoryList, error) {
	result, err := release_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseHistoryList{}, fmt.Errorf("list release history %s: %w", id, err)
	}

	history := make([]ReleaseHistorySummary, 0, len(result.Release.History.Nodes))
	for _, node := range result.Release.History.Nodes {
		history = append(history, releaseHistorySummary(node.ReleaseHistorySummaryFields))
	}

	return ReleaseHistoryList{
		History:     history,
		HasNextPage: result.Release.History.PageInfo.HasNextPage,
		EndCursor:   result.Release.History.PageInfo.EndCursor,
	}, nil
}

// ListReleaseDocuments returns documents associated with one Linear release.
func ListReleaseDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	result, err := release_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list release documents %s: %w", id, err)
	}

	documents := make([]DocumentSummary, 0, len(result.Release.Documents.Nodes))
	for _, node := range result.Release.Documents.Nodes {
		documents = append(documents, documentSummary(node.DocumentSummaryFields))
	}

	return DocumentList{
		Documents:   documents,
		HasNextPage: result.Release.Documents.PageInfo.HasNextPage,
		EndCursor:   result.Release.Documents.PageInfo.EndCursor,
	}, nil
}

// ListReleaseIssues returns issues associated with one Linear release.
func ListReleaseIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := release_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list release issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.Release.Issues.Nodes))
	for _, node := range result.Release.Issues.Nodes {
		issues = append(issues, issueSummaryFromFields(node.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Release.Issues.PageInfo.HasNextPage,
		EndCursor:   result.Release.Issues.PageInfo.EndCursor,
	}, nil
}

// ListReleaseLinks returns external links associated with one Linear release.
func ListReleaseLinks(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (EntityExternalLinkList, error) {
	result, err := release_links(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return EntityExternalLinkList{}, fmt.Errorf("list release links %s: %w", id, err)
	}

	links := make([]EntityExternalLinkSummary, 0, len(result.Release.Links.Nodes))
	for _, node := range result.Release.Links.Nodes {
		links = append(links, entityExternalLinkSummary(node.EntityExternalLinkSummaryFields))
	}

	return EntityExternalLinkList{
		Links:       links,
		HasNextPage: result.Release.Links.PageInfo.HasNextPage,
		EndCursor:   result.Release.Links.PageInfo.EndCursor,
	}, nil
}

// GetEntityExternalLinkByID returns one Linear external link by id.
func GetEntityExternalLinkByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (EntityExternalLinkSummary, error) {
	result, err := entityExternalLink(ctx, graphqlClient, id)
	if err != nil {
		return EntityExternalLinkSummary{}, fmt.Errorf("get external link %s: %w", id, err)
	}

	return entityExternalLinkSummary(result.EntityExternalLink.EntityExternalLinkSummaryFields), nil
}

// SearchReleases returns Linear releases matching a term.
func SearchReleases(ctx context.Context, graphqlClient graphql.Client, term string, limit int) (ReleaseList, error) {
	result, err := releaseSearch(ctx, graphqlClient, stringPtr(term), intPtr(limit))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("search releases: %w", err)
	}

	summaries := make([]ReleaseSummary, 0, len(result.ReleaseSearch))
	for _, node := range result.ReleaseSearch {
		summaries = append(summaries, releaseSummary(node.ReleaseSummaryFields))
	}

	return ReleaseList{Releases: summaries}, nil
}

// ListReleaseNotes returns visible Linear release notes.
func ListReleaseNotes(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleaseNoteList, error) {
	result, err := releaseNotes(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseNoteList{}, fmt.Errorf("list release notes: %w", err)
	}

	summaries := make([]ReleaseNoteSummary, 0, len(result.ReleaseNotes.Nodes))
	for _, node := range result.ReleaseNotes.Nodes {
		summaries = append(summaries, releaseNoteSummary(node.ReleaseNoteSummaryFields))
	}

	return ReleaseNoteList{
		ReleaseNotes: summaries,
		HasNextPage:  result.ReleaseNotes.PageInfo.HasNextPage,
		EndCursor:    result.ReleaseNotes.PageInfo.EndCursor,
	}, nil
}

// GetReleaseNoteByID returns one Linear release note by id.
func GetReleaseNoteByID(ctx context.Context, graphqlClient graphql.Client, id string) (ReleaseNoteSummary, error) {
	result, err := releaseNote(ctx, graphqlClient, id)
	if err != nil {
		return ReleaseNoteSummary{}, fmt.Errorf("get release note %s: %w", id, err)
	}

	return releaseNoteSummary(result.ReleaseNote.ReleaseNoteSummaryFields), nil
}

func releaseSummary(fields ReleaseSummaryFields) ReleaseSummary {
	summary := ReleaseSummary{
		ID:               fields.Id,
		Name:             fields.Name,
		SlugID:           fields.SlugId,
		Version:          stringValue(fields.Version),
		Description:      stringValue(fields.Description),
		CommitSHA:        stringValue(fields.CommitSha),
		IssueCount:       fields.IssueCount,
		ReleaseNoteCount: len(fields.ReleaseNotes),
		Trashed:          boolValue(fields.Trashed),
		URL:              fields.Url,
		StartDate:        stringValue(fields.StartDate),
		TargetDate:       stringValue(fields.TargetDate),
		StartedAt:        stringValue(fields.StartedAt),
		CompletedAt:      stringValue(fields.CompletedAt),
		CanceledAt:       stringValue(fields.CanceledAt),
		AutoArchivedAt:   stringValue(fields.AutoArchivedAt),
		CreatedAt:        fields.CreatedAt,
		UpdatedAt:        fields.UpdatedAt,
		ArchivedAt:       stringValue(fields.ArchivedAt),
		PipelineID:       fields.Pipeline.Id,
		PipelineName:     fields.Pipeline.Name,
		PipelineSlugID:   fields.Pipeline.SlugId,
		StageID:          fields.Stage.Id,
		StageName:        fields.Stage.Name,
		StageType:        string(fields.Stage.Type),
	}
	if fields.Creator != nil {
		summary.CreatorID = fields.Creator.Id
		summary.CreatorName = fields.Creator.DisplayName
	}

	return summary
}

func releaseHistorySummary(fields ReleaseHistorySummaryFields) ReleaseHistorySummary {
	return ReleaseHistorySummary{
		ID:         fields.Id,
		ReleaseID:  fields.Release.Id,
		EntryCount: countJSONArrayEntries(fields.Entries),
		Entries:    fields.Entries,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
	}
}

func entityExternalLinkSummary(fields EntityExternalLinkSummaryFields) EntityExternalLinkSummary {
	summary := EntityExternalLinkSummary{
		ID:         fields.Id,
		Label:      fields.Label,
		URL:        fields.Url,
		SortOrder:  fields.SortOrder,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
	}
	if fields.Creator != nil {
		summary.CreatorID = fields.Creator.Id
		summary.CreatorName = fields.Creator.DisplayName
	}
	if fields.Initiative != nil {
		summary.InitiativeID = fields.Initiative.Id
		summary.InitiativeName = fields.Initiative.Name
	}
	if fields.Project != nil {
		summary.ProjectID = fields.Project.Id
		summary.ProjectName = fields.Project.Name
	}

	return summary
}

func releaseNoteSummary(fields ReleaseNoteSummaryFields) ReleaseNoteSummary {
	summary := ReleaseNoteSummary{
		ID:               fields.Id,
		Title:            stringValue(fields.Title),
		SlugID:           fields.SlugId,
		GenerationStatus: releaseNoteGenerationStatus(fields.GenerationStatus),
		ReleaseCount:     fields.ReleaseCount,
		CreatedAt:        fields.CreatedAt,
		UpdatedAt:        fields.UpdatedAt,
		ArchivedAt:       stringValue(fields.ArchivedAt),
		PipelineID:       fields.Pipeline.Id,
		PipelineName:     fields.Pipeline.Name,
		PipelineSlugID:   fields.Pipeline.SlugId,
	}
	if fields.FirstRelease != nil {
		summary.FirstReleaseID = fields.FirstRelease.Id
		summary.FirstReleaseName = fields.FirstRelease.Name
		summary.FirstReleaseVersion = stringValue(fields.FirstRelease.Version)
	}
	if fields.LastRelease != nil {
		summary.LastReleaseID = fields.LastRelease.Id
		summary.LastReleaseName = fields.LastRelease.Name
		summary.LastReleaseVersion = stringValue(fields.LastRelease.Version)
	}

	return summary
}

func releaseNoteGenerationStatus(status *ReleaseNoteGenerationStatus) string {
	if status == nil {
		return ""
	}

	return string(*status)
}
