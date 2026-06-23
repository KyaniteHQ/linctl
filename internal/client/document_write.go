package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

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

// LinearDocumentCreateInput is the sparse Linear documentCreate payload linctl supports.
type LinearDocumentCreateInput struct {
	Title     string  `json:"title"`
	Content   *string `json:"content,omitempty"`
	TeamID    *string `json:"teamId,omitempty"`
	ProjectID *string `json:"projectId,omitempty"`
}

// LinearDocumentUpdateInput is the sparse Linear documentUpdate payload linctl supports.
type LinearDocumentUpdateInput struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (DocumentSummary, error) {
		input := LinearDocumentCreateInput{
			Title:   request.Title,
			Content: optionalString(request.Content),
			TeamID:  stringPtr(guard.target.Team.ID),
		}
		if guard.target.Project != nil {
			input.ProjectID = stringPtr(guard.target.Project.ID)
		}
		created, err := DocumentCreate(ctx, graphqlClient, input)
		if err != nil {
			return DocumentSummary{}, fmt.Errorf("create document: %w", err)
		}
		if !created.DocumentCreate.Success {
			return DocumentSummary{}, fmt.Errorf("%w: documentCreate returned no document", ErrMutationFailed)
		}

		return documentSummary(created.DocumentCreate.Document.DocumentSummaryFields), nil
	})
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (DocumentSummary, error) {
		if err := guardDocumentTarget(ctx, graphqlClient, guard, request.ID); err != nil {
			return DocumentSummary{}, err
		}

		updated, err := DocumentUpdate(ctx, graphqlClient, request.ID, LinearDocumentUpdateInput{
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
	})
}

// guardDocumentTarget fails closed unless the document belongs to the pinned team
// (and pinned project when configured), mirroring the resource-scoped issue guard.
func guardDocumentTarget(
	ctx context.Context,
	graphqlClient graphql.Client,
	guard writeGuard,
	documentID string,
) error {
	result, err := document(ctx, graphqlClient, documentID)
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
