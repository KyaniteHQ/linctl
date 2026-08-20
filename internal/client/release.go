package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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

//nolint:lll
type releasesNode = gql.XReleasesReleasesReleaseConnectionNodesRelease

//nolint:lll
type releaseHistoryNode = gql.XRelease_historyReleaseHistoryReleaseHistoryConnectionNodesReleaseHistory

//nolint:lll
type releaseDocumentsNode = gql.XRelease_documentsReleaseDocumentsDocumentConnectionNodesDocument

//nolint:lll
type releaseIssuesNode = gql.XRelease_issuesReleaseIssuesIssueConnectionNodesIssue

//nolint:lll
type releaseLinksNode = gql.XRelease_linksReleaseLinksEntityExternalLinkConnectionNodesEntityExternalLink

//nolint:lll
type releaseNotesNode = gql.XReleaseNotesReleaseNotesReleaseNoteConnectionNodesReleaseNote

type releasesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type releaseScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

type releaseNotesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListReleases returns visible Linear releases.
func ListReleases(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleaseList, error) {
	query := releasesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list releases", limit, defaultListPageSize,
		query.page,
		releasesNodeSummary,
	)
	if err != nil {
		return ReleaseList{}, err
	}

	return ReleaseList{Releases: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetReleaseByID returns one Linear release by id.
func GetReleaseByID(ctx context.Context, graphqlClient graphql.Client, id string) (ReleaseSummary, error) {
	result, err := gql.XRelease(ctx, graphqlClient, id)
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
	query := releaseScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release history "+id, limit, defaultListPageSize,
		query.history,
		releaseHistoryNodeSummary,
	)
	if err != nil {
		return ReleaseHistoryList{}, err
	}

	return ReleaseHistoryList{History: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListReleaseDocuments returns documents associated with one Linear release.
func ListReleaseDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	query := releaseScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release documents "+id, limit, defaultListPageSize,
		query.documents,
		releaseDocumentsNodeSummary,
	)
	if err != nil {
		return DocumentList{}, err
	}

	return DocumentList{Documents: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListReleaseIssues returns issues associated with one Linear release.
func ListReleaseIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	query := releaseScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release issues "+id, limit, defaultListPageSize,
		query.issues,
		releaseIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListReleaseLinks returns external links associated with one Linear release.
func ListReleaseLinks(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (EntityExternalLinkList, error) {
	query := releaseScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list release links "+id, limit, defaultListPageSize,
		query.links,
		releaseLinksNodeSummary,
	)
	if err != nil {
		return EntityExternalLinkList{}, err
	}

	return EntityExternalLinkList{Links: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetEntityExternalLinkByID returns one Linear external link by id.
func GetEntityExternalLinkByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (EntityExternalLinkSummary, error) {
	result, err := gql.XEntityExternalLink(ctx, graphqlClient, id)
	if err != nil {
		return EntityExternalLinkSummary{}, fmt.Errorf("get external link %s: %w", id, err)
	}

	return entityExternalLinkSummary(result.EntityExternalLink.EntityExternalLinkSummaryFields), nil
}

// SearchReleases returns Linear releases matching a term.
func SearchReleases(ctx context.Context, graphqlClient graphql.Client, term string, limit int) (ReleaseList, error) {
	result, err := gql.XReleaseSearch(ctx, graphqlClient, stringPtr(term), intPtr(limit))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("search releases: %w", err)
	}

	summaries := mapNodes(result.ReleaseSearch, func(node gql.XReleaseSearchReleaseSearchRelease) ReleaseSummary {
		return releaseSummary(node.ReleaseSummaryFields)
	})

	return ReleaseList{Releases: summaries}, nil
}

// ListReleaseNotes returns visible Linear release notes.
func ListReleaseNotes(ctx context.Context, graphqlClient graphql.Client, limit int) (ReleaseNoteList, error) {
	query := releaseNotesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list release notes", limit, defaultListPageSize,
		query.page,
		releaseNotesNodeSummary,
	)
	if err != nil {
		return ReleaseNoteList{}, err
	}

	return ReleaseNoteList{ReleaseNotes: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetReleaseNoteByID returns one Linear release note by id.
func GetReleaseNoteByID(ctx context.Context, graphqlClient graphql.Client, id string) (ReleaseNoteSummary, error) {
	result, err := gql.XReleaseNote(ctx, graphqlClient, id)
	if err != nil {
		return ReleaseNoteSummary{}, fmt.Errorf("get release note %s: %w", id, err)
	}

	return releaseNoteSummary(result.ReleaseNote.ReleaseNoteSummaryFields), nil
}

func (query releasesQuery) page(pageSize int, after *string) ([]releasesNode, bool, *string, error) {
	result, err := gql.XReleases(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Releases.Nodes,
		result.Releases.PageInfo.HasNextPage,
		result.Releases.PageInfo.EndCursor,
		nil
}

func (query releaseScopedQuery) history(
	pageSize int,
	after *string,
) ([]releaseHistoryNode, bool, *string, error) {
	result, err := gql.XRelease_history(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Release.History.Nodes,
		result.Release.History.PageInfo.HasNextPage,
		result.Release.History.PageInfo.EndCursor,
		nil
}

func (query releaseScopedQuery) documents(
	pageSize int,
	after *string,
) ([]releaseDocumentsNode, bool, *string, error) {
	result, err := gql.XRelease_documents(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Release.Documents.Nodes,
		result.Release.Documents.PageInfo.HasNextPage,
		result.Release.Documents.PageInfo.EndCursor,
		nil
}

func (query releaseScopedQuery) issues(pageSize int, after *string) ([]releaseIssuesNode, bool, *string, error) {
	result, err := gql.XRelease_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Release.Issues.Nodes,
		result.Release.Issues.PageInfo.HasNextPage,
		result.Release.Issues.PageInfo.EndCursor,
		nil
}

func (query releaseScopedQuery) links(pageSize int, after *string) ([]releaseLinksNode, bool, *string, error) {
	result, err := gql.XRelease_links(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Release.Links.Nodes,
		result.Release.Links.PageInfo.HasNextPage,
		result.Release.Links.PageInfo.EndCursor,
		nil
}

func (query releaseNotesQuery) page(pageSize int, after *string) ([]releaseNotesNode, bool, *string, error) {
	result, err := gql.XReleaseNotes(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ReleaseNotes.Nodes,
		result.ReleaseNotes.PageInfo.HasNextPage,
		result.ReleaseNotes.PageInfo.EndCursor,
		nil
}

func releasesNodeSummary(node releasesNode) ReleaseSummary {
	return releaseSummary(node.ReleaseSummaryFields)
}

func releaseHistoryNodeSummary(node releaseHistoryNode) ReleaseHistorySummary {
	return releaseHistorySummary(node.ReleaseHistorySummaryFields)
}

func releaseDocumentsNodeSummary(node releaseDocumentsNode) DocumentSummary {
	return documentSummary(node.DocumentSummaryFields)
}

func releaseIssuesNodeSummary(node releaseIssuesNode) IssueSummary {
	return issueSummaryFromFields(node.IssueSummaryFields)
}

func releaseLinksNodeSummary(node releaseLinksNode) EntityExternalLinkSummary {
	return entityExternalLinkSummary(node.EntityExternalLinkSummaryFields)
}

func releaseNotesNodeSummary(node releaseNotesNode) ReleaseNoteSummary {
	return releaseNoteSummary(node.ReleaseNoteSummaryFields)
}

func releaseSummary(fields gql.ReleaseSummaryFields) ReleaseSummary {
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

func releaseHistorySummary(fields gql.ReleaseHistorySummaryFields) ReleaseHistorySummary {
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

func entityExternalLinkSummary(fields gql.EntityExternalLinkSummaryFields) EntityExternalLinkSummary {
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

func releaseNoteSummary(fields gql.ReleaseNoteSummaryFields) ReleaseNoteSummary {
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

func releaseNoteGenerationStatus(status *gql.ReleaseNoteGenerationStatus) string {
	if status == nil {
		return ""
	}

	return string(*status)
}
