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

// CreateDocument creates a document in the resolved team after target comparison.
func CreateDocument(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request DocumentCreateRequest,
) (DocumentSummary, error) {
	if request.Title == "" {
		return DocumentSummary{}, fmt.Errorf("%w: title is required", ErrWriteInvalid)
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
		Title:   request.Title,
		Content: optionalString(request.Content),
		TeamID:  stringPtr(guard.target.Team.ID),
	}
	if guard.target.Project != nil {
		input.ProjectID = stringPtr(guard.target.Project.ID)
	}
	created, err := gql.DocumentCreate(ctx, guard.graphqlClient, input)
	if err != nil {
		return DocumentSummary{}, fmt.Errorf("create document: %w", err)
	}
	if !created.DocumentCreate.Success {
		return DocumentSummary{}, fmt.Errorf("%w: documentCreate returned no document", ErrMutationFailed)
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
		return DocumentSummary{}, fmt.Errorf("%w: document id is required", ErrWriteInvalid)
	}
	if request.Title == "" && request.Content == "" {
		return DocumentSummary{}, fmt.Errorf("%w: title or content is required", ErrWriteInvalid)
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
	if !updated.DocumentUpdate.Success {
		return DocumentSummary{}, fmt.Errorf("%w: documentUpdate returned no document", ErrMutationFailed)
	}

	return documentSummary(updated.DocumentUpdate.Document.DocumentSummaryFields), nil
}

// guardDocumentTarget fails closed unless the document belongs to the pinned team
// (and pinned project when configured), mirroring the resource-scoped issue guard.
func (guard *guardedClient) requireDocumentTarget(
	ctx context.Context,
	documentID string,
) error {
	result, err := gql.XDocument(ctx, guard.graphqlClient, documentID)
	if err != nil {
		return fmt.Errorf("get document %s: %w", documentID, err)
	}
	fields := result.Document.DocumentSummaryFields
	if fields.Team == nil || fields.Team.Id != guard.target.Team.ID || fields.Team.Key != guard.target.Team.Key {
		return fmt.Errorf(
			"%w: expected team_id=%s team_key=%s",
			ErrTargetMismatch,
			guard.target.Team.ID,
			guard.target.Team.Key,
		)
	}
	if guard.target.Project != nil {
		if fields.Project == nil || fields.Project.Id != guard.target.Project.ID {
			return fmt.Errorf(
				"%w: expected project_id=%s",
				ErrTargetMismatch,
				guard.target.Project.ID,
			)
		}
	}

	return nil
}
