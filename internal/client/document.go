package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// DocumentSummary is the compact Document model used by document commands.
type DocumentSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	SlugID     string `json:"slug_id"`
	ArchivedAt string `json:"archived_at,omitempty"`
	ParentType string `json:"parent_type,omitempty"`
	ParentID   string `json:"parent_id,omitempty"`
	ParentName string `json:"parent_name,omitempty"`
}

// DocumentList is a page of Documents.
type DocumentList struct {
	Documents []DocumentSummary `json:"documents"`
	Page
}

// DocumentCommentList is a page of body-free Comments associated with one Document.
type DocumentCommentList struct {
	DocumentID string                   `json:"document_id"`
	Comments   []CommentMetadataSummary `json:"comments"`
	Page
}

//nolint:lll
type documentsNode = gql.DocumentsDocumentsDocumentConnectionNodesDocument

//nolint:lll
type documentCommentsNode = gql.XDocument_commentsDocumentCommentsCommentConnectionNodesComment

type documentsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type documentCommentsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// ListDocuments returns visible Documents.
func ListDocuments(ctx context.Context, graphqlClient graphql.Client, limit int) (DocumentList, error) {
	query := documentsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list documents", limit, defaultListPageSize,
		query.page,
		documentsNodeSummary,
	)
	if err != nil {
		return DocumentList{}, err
	}

	return DocumentList{Documents: page.Items, Page: page.Page}, nil
}

// GetDocumentByID returns one Document by id or slug.
func GetDocumentByID(ctx context.Context, graphqlClient graphql.Client, id string) (DocumentSummary, error) {
	document, err := gql.XDocument(ctx, graphqlClient, id)
	if err != nil {
		return DocumentSummary{}, fmt.Errorf("get document %s: %w", id, err)
	}

	return documentSummary(document.Document.DocumentSummaryFields), nil
}

// ListDocumentComments returns body-free comments associated with one Document.
func ListDocumentComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentCommentList, error) {
	query := &documentCommentsQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list document comments "+id, limit, defaultListPageSize,
		query.comments,
		documentCommentsNodeSummary,
	)
	if err != nil {
		return DocumentCommentList{}, err
	}

	return DocumentCommentList{
		DocumentID: parent,
		Comments:   page.Items,
		Page:       page.Page,
	}, nil
}

func (query documentsQuery) page(pageSize int, after *string) ([]documentsNode, bool, *string, error) {
	result, err := gql.Documents(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Documents.Nodes,
		result.Documents.PageInfo.HasNextPage,
		result.Documents.PageInfo.EndCursor,
		nil
}

func (query *documentCommentsQuery) comments(
	pageSize int,
	after *string,
) ([]documentCommentsNode, string, bool, *string, error) {
	result, err := gql.XDocument_comments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, "", false, nil, err
	}

	return result.Document.Comments.Nodes,
		result.Document.Id,
		result.Document.Comments.PageInfo.HasNextPage,
		result.Document.Comments.PageInfo.EndCursor,
		nil
}

func documentsNodeSummary(document documentsNode) DocumentSummary {
	return documentSummary(document.DocumentSummaryFields)
}

func documentCommentsNodeSummary(node documentCommentsNode) CommentMetadataSummary {
	return commentMetadataSummary(node.CommentMetadataFields)
}

func documentSummary(document gql.DocumentSummaryFields) DocumentSummary {
	summary := DocumentSummary{
		ID:     document.Id,
		Title:  document.Title,
		SlugID: document.SlugId,
	}
	if document.ArchivedAt != nil {
		summary.ArchivedAt = *document.ArchivedAt
	}
	if document.Project != nil {
		summary.ParentType = "project"
		summary.ParentID = document.Project.Id
		summary.ParentName = document.Project.Name
	}
	if document.Team != nil {
		summary.ParentType = "team"
		summary.ParentID = document.Team.Id
		summary.ParentName = document.Team.Name
	}
	if document.Issue != nil {
		summary.ParentType = "issue"
		summary.ParentID = document.Issue.Id
		summary.ParentName = document.Issue.Identifier
	}
	if document.Cycle != nil {
		summary.ParentType = "cycle"
		summary.ParentID = document.Cycle.Id
		summary.ParentName = fmt.Sprintf("Cycle %.0f", document.Cycle.Number)
		if document.Cycle.Name != nil && *document.Cycle.Name != "" {
			summary.ParentName = *document.Cycle.Name
		}
	}

	return summary
}
