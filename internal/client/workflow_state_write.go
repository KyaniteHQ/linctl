package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

var workflowStateCreateTypes = map[string]bool{
	"backlog":   true,
	"unstarted": true,
	"started":   true,
	"completed": true,
	"canceled":  true,
}

// WorkflowStateCreateRequest describes a guarded WorkflowState create.
type WorkflowStateCreateRequest struct {
	ID          string
	Name        string
	Type        string
	Color       string
	Description *string
	Position    *float64
}

// WorkflowStateUpdateRequest describes a guarded WorkflowState update.
type WorkflowStateUpdateRequest struct {
	ID          string
	Name        *string
	Color       *string
	Description *string
	Position    *float64
}

// CreateWorkflowState creates a WorkflowState in the pinned team after target comparison.
func CreateWorkflowState(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request WorkflowStateCreateRequest,
) (WorkflowStateSummary, error) {
	if err := validateWorkflowStateCreateRequest(request); err != nil {
		return WorkflowStateSummary{}, err
	}
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return WorkflowStateSummary{}, err
	}

	return guard.createWorkflowState(ctx, request)
}

func (guard *guardedClient) createWorkflowState(
	ctx context.Context,
	request WorkflowStateCreateRequest,
) (WorkflowStateSummary, error) {
	input := workflowStateCreateInput(guard.target.Team.ID, request)
	created, err := gql.WorkflowStateCreate(ctx, guard.graphqlClient, input)

	return guard.finishWorkflowStateCreate(ctx, request, workflowStateCreateWriteError(err, created))
}

func workflowStateCreateInput(teamID string, request WorkflowStateCreateRequest) LinearWorkflowStateCreateInput {
	return LinearWorkflowStateCreateInput{
		ID:          stringPtr(request.ID),
		Name:        request.Name,
		Color:       request.Color,
		Type:        request.Type,
		TeamID:      teamID,
		Description: request.Description,
		Position:    request.Position,
	}
}

func workflowStateCreateWriteError(err error, created *gql.WorkflowStateCreateResponse) error {
	if err != nil {
		return fmt.Errorf("create workflow state: %w", err)
	}
	if created != nil && created.WorkflowStateCreate.Success {
		return nil
	}

	return fmt.Errorf("%w: workflowStateCreate reported no success", ErrMutationFailed)
}

func (guard *guardedClient) finishWorkflowStateCreate(
	ctx context.Context,
	request WorkflowStateCreateRequest,
	writeErr error,
) (WorkflowStateSummary, error) {
	observed, readErr := guard.readWorkflowState(ctx, request.ID)
	teamErr := error(nil)
	if readErr == nil {
		teamErr = guard.workflowStateTeamError(observed)
	}
	matches := workflowStateMatchesCreate(observed, request)

	return finishReconciledWrite(
		WorkflowStateWriteRetryClass(),
		observed,
		readErr,
		writeErr,
		teamErr,
		matches,
		writeConflictError("workflow state", request.ID),
	)
}

func workflowStateMatchesCreate(observed WorkflowStateSummary, request WorkflowStateCreateRequest) bool {
	if observed.ID != request.ID ||
		observed.Name != request.Name ||
		observed.Type != request.Type ||
		observed.Color != request.Color {
		return false
	}
	if request.Description != nil && observed.Description != *request.Description {
		return false
	}
	if request.Position != nil && observed.Position != *request.Position {
		return false
	}

	return true
}

// UpdateWorkflowState updates a WorkflowState after resolving and comparing its team.
func UpdateWorkflowState(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request WorkflowStateUpdateRequest,
) (WorkflowStateSummary, error) {
	if err := validateWorkflowStateUpdateRequest(request); err != nil {
		return WorkflowStateSummary{}, err
	}
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return WorkflowStateSummary{}, err
	}

	return guard.updateWorkflowState(ctx, request)
}

func (guard *guardedClient) updateWorkflowState(
	ctx context.Context,
	request WorkflowStateUpdateRequest,
) (WorkflowStateSummary, error) {
	existing, err := guard.requireWorkflowState(ctx, request.ID)
	if err != nil {
		return WorkflowStateSummary{}, err
	}
	updated, err := gql.WorkflowStateUpdate(
		ctx, guard.graphqlClient, request.ID, workflowStateUpdateInput(request),
	)

	return guard.finishWorkflowStateUpdate(
		ctx, request, existing.Type, workflowStateUpdateWriteError(request.ID, err, updated),
	)
}

func workflowStateUpdateInput(request WorkflowStateUpdateRequest) LinearWorkflowStateUpdateInput {
	return LinearWorkflowStateUpdateInput{
		Name:        request.Name,
		Color:       request.Color,
		Description: request.Description,
		Position:    request.Position,
	}
}

func workflowStateUpdateWriteError(id string, err error, updated *gql.WorkflowStateUpdateResponse) error {
	if err != nil {
		return fmt.Errorf("update workflow state %s: %w", id, err)
	}
	if updated != nil && updated.WorkflowStateUpdate.Success {
		return nil
	}

	return fmt.Errorf("%w: workflowStateUpdate reported no success", ErrMutationFailed)
}

func (guard *guardedClient) finishWorkflowStateUpdate(
	ctx context.Context,
	request WorkflowStateUpdateRequest,
	stateType string,
	writeErr error,
) (WorkflowStateSummary, error) {
	observed, readErr := guard.readWorkflowState(ctx, request.ID)
	teamErr := error(nil)
	if readErr == nil {
		teamErr = guard.workflowStateTeamError(observed)
	}
	matches := workflowStateMatchesUpdate(observed, request, stateType)

	return finishReconciledWrite(
		WorkflowStateWriteRetryClass(),
		observed,
		readErr,
		writeErr,
		teamErr,
		matches,
		writeConflictError("workflow state", request.ID),
	)
}

func workflowStateMatchesUpdate(
	observed WorkflowStateSummary,
	request WorkflowStateUpdateRequest,
	stateType string,
) bool {
	if observed.Type != stateType {
		return false
	}
	if request.Name != nil && observed.Name != *request.Name {
		return false
	}
	if request.Color != nil && observed.Color != *request.Color {
		return false
	}
	if request.Description != nil && observed.Description != *request.Description {
		return false
	}
	if request.Position != nil && observed.Position != *request.Position {
		return false
	}

	return true
}

func (guard *guardedClient) requireWorkflowState(
	ctx context.Context,
	id string,
) (WorkflowStateSummary, error) {
	state, err := guard.readWorkflowState(ctx, id)
	if err != nil {
		return WorkflowStateSummary{}, err
	}
	if err := guard.workflowStateTeamError(state); err != nil {
		return WorkflowStateSummary{}, err
	}

	return state, nil
}

func (guard *guardedClient) readWorkflowState(
	ctx context.Context,
	id string,
) (WorkflowStateSummary, error) {
	state, err := GetWorkflowStateByID(ctx, guard.graphqlClient, id)
	if err != nil {
		return WorkflowStateSummary{}, err
	}
	if state.ID == "" {
		return WorkflowStateSummary{}, notFoundError("workflow state %s", id)
	}

	return state, nil
}

func (guard *guardedClient) workflowStateTeamError(state WorkflowStateSummary) error {
	if state.TeamID != guard.target.Team.ID || state.TeamKey != guard.target.Team.Key {
		return guard.teamMismatchError("workflow state", state.TeamID, state.TeamKey)
	}

	return nil
}

func validateWorkflowStateCreateRequest(request WorkflowStateCreateRequest) error {
	if err := requireUUIDv4(request.ID, "id"); err != nil {
		return err
	}
	if request.Name == "" {
		return requiredFieldError("name")
	}
	if request.Color == "" {
		return requiredFieldError("color")
	}

	return validateWorkflowStateCreateType(request.Type)
}

func validateWorkflowStateCreateType(stateType string) error {
	if workflowStateCreateTypes[stateType] {
		return nil
	}

	return fmt.Errorf(
		"%w: type must be backlog, unstarted, started, completed, or canceled",
		ErrWriteInvalid,
	)
}

func validateWorkflowStateUpdateRequest(request WorkflowStateUpdateRequest) error {
	if request.ID == "" {
		return requiredFieldError("workflow state id")
	}
	if request.Name == nil && request.Color == nil && request.Description == nil && request.Position == nil {
		return requiredFieldError("name, color, description, or position")
	}
	if request.Name != nil && *request.Name == "" {
		return fmt.Errorf("%w: name must not be empty", ErrWriteInvalid)
	}
	if request.Color != nil && *request.Color == "" {
		return fmt.Errorf("%w: color must not be empty", ErrWriteInvalid)
	}

	return nil
}
