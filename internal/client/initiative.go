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
	OrgID       string `json:"org_id"`
}

// InitiativeList is a page of initiatives.
type InitiativeList struct {
	Initiatives []InitiativeSummary `json:"initiatives"`
	Page
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
	History []InitiativeHistorySummary `json:"history"`
	Page
}

//nolint:lll
type initiativesNode = gql.XInitiativesInitiativesInitiativeConnectionNodesInitiative

//nolint:lll
type initiativeHistoryNode = gql.XInitiative_historyInitiativeHistoryInitiativeHistoryConnectionNodesInitiativeHistory

//nolint:lll
type initiativeLinksNode = gql.XInitiative_linksInitiativeLinksEntityExternalLinkConnectionNodesEntityExternalLink

//nolint:lll
type subInitiativesNode = gql.XInitiative_subInitiativesInitiativeSubInitiativesInitiativeConnectionNodesInitiative

//nolint:lll
type initiativeUpdatesForInitiativeNode = gql.XInitiative_initiativeUpdatesInitiativeInitiativeUpdatesInitiativeUpdateConnectionNodesInitiativeUpdate

//nolint:lll
type initiativeDocumentsNode = gql.XInitiative_documentsInitiativeDocumentsDocumentConnectionNodesDocument

//nolint:lll
type initiativeProjectsNode = gql.XInitiative_projectsInitiativeProjectsProjectConnectionNodesProject

type initiativesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type initiativeScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// ListInitiatives returns visible initiatives.
func ListInitiatives(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeList, error) {
	query := initiativesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list initiatives", limit, defaultListPageSize,
		query.page,
		initiativesNodeSummary,
	)
	if err != nil {
		return InitiativeList{}, err
	}

	return InitiativeList{Initiatives: page.Items, Page: page.Page}, nil
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
	query := initiativeScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative history "+id, limit, defaultListPageSize,
		query.history,
		initiativeHistoryNodeSummary,
	)
	if err != nil {
		return InitiativeHistoryList{}, err
	}

	return InitiativeHistoryList{History: page.Items, Page: page.Page}, nil
}

// ListInitiativeLinks returns external links associated with one initiative.
func ListInitiativeLinks(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (EntityExternalLinkList, error) {
	query := initiativeScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative links "+id, limit, defaultListPageSize,
		query.links,
		initiativeLinksNodeSummary,
	)
	if err != nil {
		return EntityExternalLinkList{}, err
	}

	return EntityExternalLinkList{Links: page.Items, Page: page.Page}, nil
}

// ListSubInitiatives returns child initiatives associated with one initiative.
func ListSubInitiatives(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeList, error) {
	query := initiativeScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative sub-initiatives "+id, limit, defaultListPageSize,
		query.subInitiatives,
		subInitiativesNodeSummary,
	)
	if err != nil {
		return InitiativeList{}, err
	}

	return InitiativeList{Initiatives: page.Items, Page: page.Page}, nil
}

// ListInitiativeUpdatesForInitiative returns status updates associated with one initiative.
func ListInitiativeUpdatesForInitiative(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeUpdateList, error) {
	query := initiativeScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative updates "+id, limit, defaultListPageSize,
		query.updates,
		initiativeUpdatesForInitiativeNodeSummary,
	)
	if err != nil {
		return InitiativeUpdateList{}, err
	}

	return InitiativeUpdateList{Updates: page.Items, Page: page.Page}, nil
}

// ListInitiativeDocuments returns Documents associated with one initiative.
func ListInitiativeDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	query := initiativeScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative documents "+id, limit, defaultListPageSize,
		query.documents,
		initiativeDocumentsNodeSummary,
	)
	if err != nil {
		return DocumentList{}, err
	}

	return DocumentList{Documents: page.Items, Page: page.Page}, nil
}

// ListInitiativeProjects returns Projects directly associated with one initiative.
func ListInitiativeProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectList, error) {
	query := initiativeScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list initiative projects "+id, limit, defaultListPageSize,
		query.projects,
		initiativeProjectsNodeSummary,
	)
	if err != nil {
		return ProjectList{}, err
	}

	return ProjectList{Projects: page.Items, Page: page.Page}, nil
}

func (query initiativesQuery) page(pageSize int, after *string) ([]initiativesNode, bool, *string, error) {
	result, err := gql.XInitiatives(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiatives.Nodes,
		result.Initiatives.PageInfo.HasNextPage,
		result.Initiatives.PageInfo.EndCursor,
		nil
}

func (query initiativeScopedQuery) history(
	pageSize int,
	after *string,
) ([]initiativeHistoryNode, bool, *string, error) {
	result, err := gql.XInitiative_history(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiative.History.Nodes,
		result.Initiative.History.PageInfo.HasNextPage,
		result.Initiative.History.PageInfo.EndCursor,
		nil
}

func (query initiativeScopedQuery) links(
	pageSize int,
	after *string,
) ([]initiativeLinksNode, bool, *string, error) {
	result, err := gql.XInitiative_links(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiative.Links.Nodes,
		result.Initiative.Links.PageInfo.HasNextPage,
		result.Initiative.Links.PageInfo.EndCursor,
		nil
}

func (query initiativeScopedQuery) subInitiatives(
	pageSize int,
	after *string,
) ([]subInitiativesNode, bool, *string, error) {
	result, err := gql.XInitiative_subInitiatives(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiative.SubInitiatives.Nodes,
		result.Initiative.SubInitiatives.PageInfo.HasNextPage,
		result.Initiative.SubInitiatives.PageInfo.EndCursor,
		nil
}

func (query initiativeScopedQuery) updates(
	pageSize int,
	after *string,
) ([]initiativeUpdatesForInitiativeNode, bool, *string, error) {
	result, err := gql.XInitiative_initiativeUpdates(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiative.InitiativeUpdates.Nodes,
		result.Initiative.InitiativeUpdates.PageInfo.HasNextPage,
		result.Initiative.InitiativeUpdates.PageInfo.EndCursor,
		nil
}

func (query initiativeScopedQuery) documents(
	pageSize int,
	after *string,
) ([]initiativeDocumentsNode, bool, *string, error) {
	result, err := gql.XInitiative_documents(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiative.Documents.Nodes,
		result.Initiative.Documents.PageInfo.HasNextPage,
		result.Initiative.Documents.PageInfo.EndCursor,
		nil
}

func (query initiativeScopedQuery) projects(
	pageSize int,
	after *string,
) ([]initiativeProjectsNode, bool, *string, error) {
	result, err := gql.XInitiative_projects(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true), boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Initiative.Projects.Nodes,
		result.Initiative.Projects.PageInfo.HasNextPage,
		result.Initiative.Projects.PageInfo.EndCursor,
		nil
}

func initiativesNodeSummary(node initiativesNode) InitiativeSummary {
	return initiativeSummary(node.InitiativeSummaryFields)
}

func initiativeHistoryNodeSummary(node initiativeHistoryNode) InitiativeHistorySummary {
	return initiativeHistorySummary(node.InitiativeHistorySummaryFields)
}

func initiativeLinksNodeSummary(node initiativeLinksNode) EntityExternalLinkSummary {
	return entityExternalLinkSummary(node.EntityExternalLinkSummaryFields)
}

func subInitiativesNodeSummary(node subInitiativesNode) InitiativeSummary {
	return initiativeSummary(node.InitiativeSummaryFields)
}

func initiativeUpdatesForInitiativeNodeSummary(node initiativeUpdatesForInitiativeNode) InitiativeUpdateSummary {
	return initiativeUpdateSummary(node.InitiativeUpdateSummaryFields)
}

func initiativeDocumentsNodeSummary(node initiativeDocumentsNode) DocumentSummary {
	return documentSummary(node.DocumentSummaryFields)
}

func initiativeProjectsNodeSummary(node initiativeProjectsNode) ProjectSummary {
	return projectSummaryFromFields(node.ProjectSummaryFields)
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
		OrgID:       fields.Organization.Id,
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
