package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// WorkflowStateSummary is the compact workflow state model used by read-only commands.
type WorkflowStateSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Color       string  `json:"color"`
	Description string  `json:"description,omitempty"`
	Position    float64 `json:"position"`
	TeamID      string  `json:"team_id"`
	TeamKey     string  `json:"team_key"`
	TeamName    string  `json:"team_name"`
}

// WorkflowStateList is a page of workflow states.
type WorkflowStateList struct {
	WorkflowStates []WorkflowStateSummary `json:"workflow_states"`
	Page
}

// WorkflowStateIssueList is a page of Issues currently associated with one WorkflowState.
type WorkflowStateIssueList struct {
	WorkflowStateID   string         `json:"workflow_state_id"`
	WorkflowStateName string         `json:"workflow_state_name"`
	Issues            []IssueSummary `json:"issues"`
	Page
}

//nolint:lll
type workflowStatesNode = gql.XWorkflowStatesWorkflowStatesWorkflowStateConnectionNodesWorkflowState

//nolint:lll
type workflowStateIssuesNode = gql.XWorkflowState_issuesWorkflowStateIssuesIssueConnectionNodesIssue

type workflowStatesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type workflowStateIssuesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// workflowStateIssuesParent is the connection parent metadata workflowStateIssuesQuery reads out of
// every page. Linear repeats it per page, so the last page wins.
type workflowStateIssuesParent struct {
	workflowStateID   string
	workflowStateName string
}

// ListWorkflowStates returns visible workflow states.
func ListWorkflowStates(ctx context.Context, graphqlClient graphql.Client, limit int) (WorkflowStateList, error) {
	query := workflowStatesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list workflow states", limit, defaultListPageSize,
		query.page,
		workflowStatesNodeSummary,
	)
	if err != nil {
		return WorkflowStateList{}, err
	}

	return WorkflowStateList{
		WorkflowStates: page.Items,
		Page:           page.Page,
	}, nil
}

// GetWorkflowStateByID returns one workflow state by Linear id.
func GetWorkflowStateByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (WorkflowStateSummary, error) {
	state, err := gql.XWorkflowState(ctx, graphqlClient, id)
	if err != nil {
		return WorkflowStateSummary{}, fmt.Errorf("get workflow state %s: %w", id, err)
	}

	return workflowStateSummary(state.WorkflowState.WorkflowStateSummaryFields), nil
}

// ListWorkflowStateIssues returns Issues currently associated with one WorkflowState.
func ListWorkflowStateIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (WorkflowStateIssueList, error) {
	query := &workflowStateIssuesQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list workflow state issues "+id, limit, defaultListPageSize,
		query.issues,
		workflowStateIssuesNodeSummary,
	)
	if err != nil {
		return WorkflowStateIssueList{}, err
	}

	return WorkflowStateIssueList{
		WorkflowStateID:   parent.workflowStateID,
		WorkflowStateName: parent.workflowStateName,
		Issues:            page.Items,
		Page:              page.Page,
	}, nil
}

func (query workflowStatesQuery) page(pageSize int, after *string) ([]workflowStatesNode, bool, *string, error) {
	result, err := gql.XWorkflowStates(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.WorkflowStates.Nodes,
		result.WorkflowStates.PageInfo.HasNextPage,
		result.WorkflowStates.PageInfo.EndCursor,
		nil
}

func (query *workflowStateIssuesQuery) issues(
	pageSize int,
	after *string,
) ([]workflowStateIssuesNode, workflowStateIssuesParent, bool, *string, error) {
	result, err := gql.XWorkflowState_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, workflowStateIssuesParent{}, false, nil, err
	}

	return result.WorkflowState.Issues.Nodes,
		workflowStateIssuesParent{
			workflowStateID:   result.WorkflowState.Id,
			workflowStateName: result.WorkflowState.Name,
		},
		result.WorkflowState.Issues.PageInfo.HasNextPage,
		result.WorkflowState.Issues.PageInfo.EndCursor,
		nil
}

func workflowStatesNodeSummary(state workflowStatesNode) WorkflowStateSummary {
	return workflowStateSummary(state.WorkflowStateSummaryFields)
}

func workflowStateIssuesNodeSummary(issue workflowStateIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func workflowStateSummary(state gql.WorkflowStateSummaryFields) WorkflowStateSummary {
	return WorkflowStateSummary{
		ID:          state.Id,
		Name:        state.Name,
		Type:        state.Type,
		Color:       state.Color,
		Description: stringValue(state.Description),
		Position:    state.Position,
		TeamID:      state.Team.Id,
		TeamKey:     state.Team.Key,
		TeamName:    state.Team.Name,
	}
}
