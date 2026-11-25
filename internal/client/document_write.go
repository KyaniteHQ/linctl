package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// DocumentCreateRequest describes a guarded document create.
type DocumentCreateRequest struct {
	Title   string
	Content string
}

// DocumentUpdateRequest describes a guarded document update.
type DocumentUpdateRequest struct {
	ID      string
	Title   string
	Content string
}

// CreateDocument creates a document in the pinned project after target comparison.
func CreateDocument(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request DocumentCreateRequest,
) (DocumentSummary, error) {
	if request.Title == "" {
		return DocumentSummary{}, requiredFieldError("title")
	}
	if err := requireDocumentProjectPin(expected); err != nil {
		return DocumentSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return DocumentSummary{}, err
	}

	return guard.createDocument(ctx, request)
}

func (guard *guardedClient) createDocument(
	ctx context.Context,
	request DocumentCreateRequest,
) (DocumentSummary, error) {
	input := LinearDocumentCreateInput{
		Title:     request.Title,
		Content:   optionalString(request.Content),
		ProjectID: stringPtr(guard.target.Project.ID),
	}
	created, err := gql.DocumentCreate(ctx, guard.graphqlClient, input)
	if err != nil {
		return DocumentSummary{}, fmt.Errorf("create document: %w", err)
	}
	if err := mutationSuccess(created.DocumentCreate.Success, "documentCreate"); err != nil {
		return DocumentSummary{}, err
	}

	return documentSummary(created.DocumentCreate.Document.DocumentSummaryFields), nil
}

// UpdateDocument updates an existing document after resolving and comparing the pinned target.
func UpdateDocument(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request DocumentUpdateRequest,
) (DocumentSummary, error) {
	if request.ID == "" {
		return DocumentSummary{}, requiredFieldError("document id")
	}
	if request.Title == "" && request.Content == "" {
		return DocumentSummary{}, requiredFieldError("title or content")
	}
	if err := requireDocumentProjectPin(expected); err != nil {
		return DocumentSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return DocumentSummary{}, err
	}

	return guard.updateDocument(ctx, request)
}

func (guard *guardedClient) updateDocument(
	ctx context.Context,
	request DocumentUpdateRequest,
) (DocumentSummary, error) {
	if err := guard.requireDocumentTarget(ctx, request.ID); err != nil {
		return DocumentSummary{}, err
	}

	updated, err := gql.DocumentUpdate(ctx, guard.graphqlClient, request.ID, LinearDocumentUpdateInput{
		Title:   optionalString(request.Title),
		Content: optionalString(request.Content),
	})
	if err != nil {
		return DocumentSummary{}, fmt.Errorf("update document %s: %w", request.ID, err)
	}
	if err := mutationSuccess(updated.DocumentUpdate.Success, "documentUpdate"); err != nil {
		return DocumentSummary{}, err
	}

	return documentSummary(updated.DocumentUpdate.Document.DocumentSummaryFields), nil
}

// requireDocumentProjectPin refuses document writes without a pinned project:
// Linear's public schema exposes no team on a document, so the pinned project
// is the only parent the guard can prove.
func requireDocumentProjectPin(expected config.Target) error {
	if expected.ProjectID == "" {
		return fmt.Errorf("%w: document writes require a pinned project (set project_id in .linctl.toml)", ErrWriteInvalid)
	}

	return nil
}

// requireDocumentTarget fails closed unless the document belongs to the pinned project.
func (guard *guardedClient) requireDocumentTarget(
	ctx context.Context,
	documentID string,
) error {
	result, err := gql.XDocument(ctx, guard.graphqlClient, documentID)
	if err != nil {
		return fmt.Errorf("get document %s: %w", documentID, err)
	}
	project := result.Document.Project
	if project == nil || project.Id != guard.target.Project.ID {
		return fmt.Errorf("%w: expected project_id=%s", ErrTargetMismatch, guard.target.Project.ID)
	}

	return nil
}
