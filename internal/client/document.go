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
	Documents   []DocumentSummary `json:"documents"`
	HasNextPage bool              `json:"has_next_page"`
	EndCursor   *string           `json:"end_cursor,omitempty"`
}

// DocumentCommentList is a page of body-free Comments associated with one Document.
type DocumentCommentList struct {
	DocumentID  string                   `json:"document_id"`
	Comments    []CommentMetadataSummary `json:"comments"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ListDocuments returns visible Documents.
func ListDocuments(ctx context.Context, graphqlClient graphql.Client, limit int) (DocumentList, error) {
	documents, err := gql.Documents(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list documents: %w", err)
	}

	summaries := mapNodes(documents.Documents.Nodes, func(
		document gql.DocumentsDocumentsDocumentConnectionNodesDocument,
	) DocumentSummary {
		return documentSummary(document.DocumentSummaryFields)
	})

	return DocumentList{
		Documents:   summaries,
		HasNextPage: documents.Documents.PageInfo.HasNextPage,
		EndCursor:   documents.Documents.PageInfo.EndCursor,
	}, nil
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
	result, err := gql.XDocument_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentCommentList{}, fmt.Errorf("list document comments %s: %w", id, err)
	}

	comments := mapNodes(result.Document.Comments.Nodes, func(
		node gql.XDocument_commentsDocumentCommentsCommentConnectionNodesComment,
	) CommentMetadataSummary {
		return commentMetadataSummary(node.CommentMetadataFields)
	})

	return DocumentCommentList{
		DocumentID:  result.Document.Id,
		Comments:    comments,
		HasNextPage: result.Document.Comments.PageInfo.HasNextPage,
		EndCursor:   result.Document.Comments.PageInfo.EndCursor,
	}, nil
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
