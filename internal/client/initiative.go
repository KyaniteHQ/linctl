package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	result, err := gql.XInitiatives(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeList{}, fmt.Errorf("list initiatives: %w", err)
	}

	summaries := mapNodes(result.Initiatives.Nodes, func(
		node gql.XInitiativesInitiativesInitiativeConnectionNodesInitiative,
	) InitiativeSummary {
		return initiativeSummary(node.InitiativeSummaryFields)
	})

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
	result, err := gql.XInitiative(ctx, graphqlClient, id)
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
	result, err := gql.XInitiative_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeHistoryList{}, fmt.Errorf("list initiative history %s: %w", id, err)
	}

	history := mapNodes(result.Initiative.History.Nodes, func(
		node gql.XInitiative_historyInitiativeHistoryInitiativeHistoryConnectionNodesInitiativeHistory,
	) InitiativeHistorySummary {
		return initiativeHistorySummary(node.InitiativeHistorySummaryFields)
	})

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
	result, err := gql.XInitiative_links(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return EntityExternalLinkList{}, fmt.Errorf("list initiative links %s: %w", id, err)
	}

	links := mapNodes(result.Initiative.Links.Nodes, func(
		node gql.XInitiative_linksInitiativeLinksEntityExternalLinkConnectionNodesEntityExternalLink,
	) EntityExternalLinkSummary {
		return entityExternalLinkSummary(node.EntityExternalLinkSummaryFields)
	})

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
	result, err := gql.XInitiative_subInitiatives(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeList{}, fmt.Errorf("list initiative sub-initiatives %s: %w", id, err)
	}

	initiatives := mapNodes(result.Initiative.SubInitiatives.Nodes, func(
		node gql.XInitiative_subInitiativesInitiativeSubInitiativesInitiativeConnectionNodesInitiative,
	) InitiativeSummary {
		return initiativeSummary(node.InitiativeSummaryFields)
	})

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
	result, err := gql.XInitiative_initiativeUpdates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeUpdateList{}, fmt.Errorf("list initiative updates %s: %w", id, err)
	}

	updates := mapNodes(result.Initiative.InitiativeUpdates.Nodes, func(
		node gql.
			XInitiative_initiativeUpdatesInitiativeInitiativeUpdatesInitiativeUpdateConnectionNodesInitiativeUpdate,
	) InitiativeUpdateSummary {
		return initiativeUpdateSummary(node.InitiativeUpdateSummaryFields)
	})

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
	result, err := gql.XInitiative_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list initiative documents %s: %w", id, err)
	}

	documents := mapNodes(result.Initiative.Documents.Nodes, func(
		node gql.XInitiative_documentsInitiativeDocumentsDocumentConnectionNodesDocument,
	) DocumentSummary {
		return documentSummary(node.DocumentSummaryFields)
	})

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
	result, err := gql.XInitiative_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true), boolPtr(false))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list initiative projects %s: %w", id, err)
	}

	projects := mapNodes(result.Initiative.Projects.Nodes, func(
		node gql.XInitiative_projectsInitiativeProjectsProjectConnectionNodesProject,
	) ProjectSummary {
		return projectSummaryFromFields(node.ProjectSummaryFields)
	})

	return ProjectList{
		Projects:    projects,
		HasNextPage: result.Initiative.Projects.PageInfo.HasNextPage,
		EndCursor:   result.Initiative.Projects.PageInfo.EndCursor,
	}, nil
}

func initiativeSummary(fields gql.InitiativeSummaryFields) InitiativeSummary {
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

func initiativeHistorySummary(fields gql.InitiativeHistorySummaryFields) InitiativeHistorySummary {
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
