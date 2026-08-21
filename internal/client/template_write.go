package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

const issueTemplateType = "issue"

// TemplateCreateRequest describes a guarded issue-template create.
type TemplateCreateRequest struct {
	ID   string
	Name string
	Type string
	Data json.RawMessage
}

// TemplateUpdateRequest describes a guarded issue-template update.
type TemplateUpdateRequest struct {
	ID   string
	Name *string
	Data json.RawMessage
}

// CreateTemplate creates a team-scoped issue template after target comparison.
func CreateTemplate(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request TemplateCreateRequest,
) (TemplateDetail, error) {
	prepared, err := prepareTemplateCreateRequest(request)
	if err != nil {
		return TemplateDetail{}, err
	}
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return TemplateDetail{}, err
	}

	return guard.createTemplate(ctx, prepared)
}

func prepareTemplateCreateRequest(request TemplateCreateRequest) (TemplateCreateRequest, error) {
	if err := requireUUIDv4(request.ID, "id"); err != nil {
		return TemplateCreateRequest{}, err
	}
	if request.Name == "" {
		return TemplateCreateRequest{}, requiredFieldError("name")
	}
	if request.Type != issueTemplateType {
		return TemplateCreateRequest{}, fmt.Errorf("%w: type must be issue", ErrWriteInvalid)
	}
	data, err := CanonicalTemplateData(request.Data)
	if err != nil {
		return TemplateCreateRequest{}, err
	}
	request.Data = data

	return request, nil
}

func (guard *guardedClient) createTemplate(
	ctx context.Context,
	request TemplateCreateRequest,
) (TemplateDetail, error) {
	created, err := gql.TemplateCreate(ctx, guard.graphqlClient, LinearTemplateCreateInput{
		ID:           stringPtr(request.ID),
		Name:         request.Name,
		Type:         issueTemplateType,
		TeamID:       stringPtr(guard.target.Team.ID),
		TemplateData: request.Data,
	})

	return guard.finishTemplateWrite(ctx, request.ID, templateCreateWriteError(err, created),
		func(observed TemplateDetail) bool {
			return templateMatchesCreate(observed, request)
		})
}

func templateCreateWriteError(err error, created *gql.TemplateCreateResponse) error {
	if err != nil {
		return fmt.Errorf("create template: %w", err)
	}

	return mutationSuccess(created != nil && created.TemplateCreate.Success, "templateCreate")
}

func templateMatchesCreate(observed TemplateDetail, request TemplateCreateRequest) bool {
	return observed.ID == request.ID &&
		observed.Name == request.Name &&
		observed.Type == issueTemplateType &&
		bytes.Equal(observed.Data, request.Data)
}

// UpdateTemplate updates a team-scoped issue template after resolving its scope.
func UpdateTemplate(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request TemplateUpdateRequest,
) (TemplateDetail, error) {
	prepared, err := prepareTemplateUpdateRequest(request)
	if err != nil {
		return TemplateDetail{}, err
	}
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return TemplateDetail{}, err
	}

	return guard.updateTemplate(ctx, prepared)
}

func prepareTemplateUpdateRequest(request TemplateUpdateRequest) (TemplateUpdateRequest, error) {
	if request.ID == "" {
		return TemplateUpdateRequest{}, requiredFieldError("template id")
	}
	if request.Name == nil && request.Data == nil {
		return TemplateUpdateRequest{}, requiredFieldError("name or data")
	}
	if request.Name != nil && *request.Name == "" {
		return TemplateUpdateRequest{}, fmt.Errorf("%w: name must not be empty", ErrWriteInvalid)
	}
	if request.Data == nil {
		return request, nil
	}
	data, err := CanonicalTemplateData(request.Data)
	if err != nil {
		return TemplateUpdateRequest{}, err
	}
	request.Data = data

	return request, nil
}

func (guard *guardedClient) updateTemplate(
	ctx context.Context,
	request TemplateUpdateRequest,
) (TemplateDetail, error) {
	if _, err := guard.requireWritableIssueTemplate(ctx, request.ID); err != nil {
		return TemplateDetail{}, err
	}
	updated, err := gql.TemplateUpdate(ctx, guard.graphqlClient, request.ID, LinearTemplateUpdateInput{
		Name:         request.Name,
		TemplateData: request.Data,
	})

	return guard.finishTemplateWrite(ctx, request.ID, templateUpdateWriteError(request.ID, err, updated),
		func(observed TemplateDetail) bool {
			return templateMatchesUpdate(observed, request)
		})
}

func templateUpdateWriteError(id string, err error, updated *gql.TemplateUpdateResponse) error {
	if err != nil {
		return fmt.Errorf("update template %s: %w", id, err)
	}

	return mutationSuccess(updated != nil && updated.TemplateUpdate.Success, "templateUpdate")
}

func (guard *guardedClient) finishTemplateWrite(
	ctx context.Context,
	id string,
	writeErr error,
	matches func(TemplateDetail) bool,
) (TemplateDetail, error) {
	observed, readErr := GetTemplateDetail(ctx, guard.graphqlClient, id)
	scopeErr := error(nil)
	if readErr == nil {
		scopeErr = guard.templateScopeError(observed)
	}

	return finishReconciledWrite(
		TemplateWriteRetryClass(),
		observed,
		readErr,
		writeErr,
		scopeErr,
		matches(observed),
		writeConflictError("template", id),
	)
}

func templateMatchesUpdate(observed TemplateDetail, request TemplateUpdateRequest) bool {
	if observed.ID != request.ID {
		return false
	}
	if observed.Type != issueTemplateType {
		return false
	}
	if request.Name != nil && observed.Name != *request.Name {
		return false
	}
	if request.Data != nil && !bytes.Equal(observed.Data, request.Data) {
		return false
	}

	return true
}

func (guard *guardedClient) requireWritableIssueTemplate(
	ctx context.Context,
	id string,
) (TemplateDetail, error) {
	detail, err := GetTemplateDetail(ctx, guard.graphqlClient, id)
	if err != nil {
		return TemplateDetail{}, err
	}
	if err := guard.writableIssueTemplateError(detail); err != nil {
		return TemplateDetail{}, err
	}

	return detail, nil
}

func (guard *guardedClient) templateScopeError(detail TemplateDetail) error {
	if detail.TeamID == "" {
		return fmt.Errorf("%w: template %s has no team", ErrTargetMismatch, detail.ID)
	}
	if detail.TeamID != guard.target.Team.ID || detail.TeamKey != guard.target.Team.Key {
		return guard.teamMismatchError("template", detail.TeamID, detail.TeamKey)
	}
	if detail.PipelineID != "" {
		return fmt.Errorf("%w: template %s is pipeline-scoped", ErrTargetMismatch, detail.ID)
	}

	return nil
}

func (guard *guardedClient) writableIssueTemplateError(detail TemplateDetail) error {
	if err := guard.templateScopeError(detail); err != nil {
		return err
	}
	if detail.Type != issueTemplateType {
		return fmt.Errorf("%w: template type must be issue", ErrWriteInvalid)
	}

	return nil
}
