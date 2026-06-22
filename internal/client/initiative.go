package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// InitiativeSummary is the compact initiative model used by read-only commands.
type InitiativeSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	TargetDate  string `json:"target_date,omitempty"`
	SlugID      string `json:"slug_id"`
	URL         string `json:"url"`
}

// InitiativeList is a page of initiatives.
type InitiativeList struct {
	Initiatives []InitiativeSummary `json:"initiatives"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// InitiativeHistorySummary is the compact initiative history model used by read-only commands.
type InitiativeHistorySummary struct {
	ID           string          `json:"id"`
	InitiativeID string          `json:"initiative_id"`
	EntryCount   int             `json:"entry_count"`
	Entries      json.RawMessage `json:"entries"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
	ArchivedAt   string          `json:"archived_at,omitempty"`
}

// InitiativeHistoryList is a page of Linear initiative history records.
type InitiativeHistoryList struct {
	History     []InitiativeHistorySummary `json:"history"`
	HasNextPage bool                       `json:"has_next_page"`
	EndCursor   *string                    `json:"end_cursor,omitempty"`
}

// ListInitiatives returns visible initiatives.
func ListInitiatives(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeList, error) {
	result, err := initiatives(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeList{}, fmt.Errorf("list initiatives: %w", err)
	}

	summaries := make([]InitiativeSummary, 0, len(result.Initiatives.Nodes))
	for _, node := range result.Initiatives.Nodes {
		summaries = append(summaries, initiativeSummary(node.InitiativeSummaryFields))
	}

	return InitiativeList{
		Initiatives: summaries,
		HasNextPage: result.Initiatives.PageInfo.HasNextPage,
		EndCursor:   result.Initiatives.PageInfo.EndCursor,
	}, nil
}

// GetInitiativeByID returns one initiative by Linear id or slug.
func GetInitiativeByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeSummary, error) {
	result, err := initiative(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeSummary{}, fmt.Errorf("get initiative %s: %w", id, err)
	}

	return initiativeSummary(result.Initiative.InitiativeSummaryFields), nil
}

// ListInitiativeHistory returns history records associated with one initiative.
func ListInitiativeHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeHistoryList, error) {
	result, err := initiative_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeHistoryList{}, fmt.Errorf("list initiative history %s: %w", id, err)
	}

	history := make([]InitiativeHistorySummary, 0, len(result.Initiative.History.Nodes))
	for _, node := range result.Initiative.History.Nodes {
		history = append(history, initiativeHistorySummary(node.InitiativeHistorySummaryFields))
	}

	return InitiativeHistoryList{
		History:     history,
		HasNextPage: result.Initiative.History.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.History.PageInfo.EndCursor,
	}, nil
}

// ListInitiativeLinks returns external links associated with one initiative.
func ListInitiativeLinks(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (EntityExternalLinkList, error) {
	result, err := initiative_links(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return EntityExternalLinkList{}, fmt.Errorf("list initiative links %s: %w", id, err)
	}

	links := make([]EntityExternalLinkSummary, 0, len(result.Initiative.Links.Nodes))
	for _, node := range result.Initiative.Links.Nodes {
		links = append(links, entityExternalLinkSummary(node.EntityExternalLinkSummaryFields))
	}

	return EntityExternalLinkList{
		Links:       links,
		HasNextPage: result.Initiative.Links.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.Links.PageInfo.EndCursor,
	}, nil
}

// ListSubInitiatives returns child initiatives associated with one initiative.
func ListSubInitiatives(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeList, error) {
	result, err := initiative_subInitiatives(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeList{}, fmt.Errorf("list initiative sub-initiatives %s: %w", id, err)
	}

	initiatives := make([]InitiativeSummary, 0, len(result.Initiative.SubInitiatives.Nodes))
	for _, node := range result.Initiative.SubInitiatives.Nodes {
		initiatives = append(initiatives, initiativeSummary(node.InitiativeSummaryFields))
	}

	return InitiativeList{
		Initiatives: initiatives,
		HasNextPage: result.Initiative.SubInitiatives.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.SubInitiatives.PageInfo.EndCursor,
	}, nil
}

// ListInitiativeUpdatesForInitiative returns status updates associated with one initiative.
func ListInitiativeUpdatesForInitiative(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeUpdateList, error) {
	result, err := initiative_initiativeUpdates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeUpdateList{}, fmt.Errorf("list initiative updates %s: %w", id, err)
	}

	updates := make([]InitiativeUpdateSummary, 0, len(result.Initiative.InitiativeUpdates.Nodes))
	for _, node := range result.Initiative.InitiativeUpdates.Nodes {
		updates = append(updates, initiativeUpdateSummary(node.InitiativeUpdateSummaryFields))
	}

	return InitiativeUpdateList{
		Updates:     updates,
		HasNextPage: result.Initiative.InitiativeUpdates.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.InitiativeUpdates.PageInfo.EndCursor,
	}, nil
}

// ListInitiativeDocuments returns Documents associated with one initiative.
func ListInitiativeDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	result, err := initiative_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list initiative documents %s: %w", id, err)
	}

	documents := make([]DocumentSummary, 0, len(result.Initiative.Documents.Nodes))
	for _, node := range result.Initiative.Documents.Nodes {
		documents = append(documents, documentSummary(node.DocumentSummaryFields))
	}

	return DocumentList{
		Documents:   documents,
		HasNextPage: result.Initiative.Documents.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.Documents.PageInfo.EndCursor,
	}, nil
}

// ListInitiativeProjects returns Projects directly associated with one initiative.
func ListInitiativeProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectList, error) {
	result, err := initiative_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true), boolPtr(false))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list initiative projects %s: %w", id, err)
	}

	projects := make([]ProjectSummary, 0, len(result.Initiative.Projects.Nodes))
	for _, node := range result.Initiative.Projects.Nodes {
		projects = append(projects, projectSummaryFromFields(node.ProjectSummaryFields))
	}

	return ProjectList{
		Projects:    projects,
		HasNextPage: result.Initiative.Projects.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.Projects.PageInfo.EndCursor,
	}, nil
}

func initiativeSummary(fields InitiativeSummaryFields) InitiativeSummary {
	return InitiativeSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		Description: stringValue(fields.Description),
		Status:      string(fields.Status),
		Priority:    fields.Priority,
		TargetDate:  stringValue(fields.TargetDate),
		SlugID:      fields.SlugId,
		URL:         fields.Url,
	}
}

func initiativeHistorySummary(fields InitiativeHistorySummaryFields) InitiativeHistorySummary {
	return InitiativeHistorySummary{
		ID:           fields.Id,
		InitiativeID: fields.Initiative.Id,
		EntryCount:   countJSONArrayEntries(fields.Entries),
		Entries:      fields.Entries,
		CreatedAt:    fields.CreatedAt,
		UpdatedAt:    fields.UpdatedAt,
		ArchivedAt:   stringValue(fields.ArchivedAt),
	}
}
