package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// WorkflowStateSummary is the compact workflow state model used by read-only commands.
type WorkflowStateSummary struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Color    string  `json:"color"`
	Position float64 `json:"position"`
	TeamID   string  `json:"team_id"`
	TeamKey  string  `json:"team_key"`
	TeamName string  `json:"team_name"`
}

// WorkflowStateList is a page of workflow states.
type WorkflowStateList struct {
	WorkflowStates []WorkflowStateSummary `json:"workflow_states"`
	HasNextPage    bool                   `json:"has_next_page"`
	EndCursor      *string                `json:"end_cursor,omitempty"`
}

// ListWorkflowStates returns visible workflow states.
func ListWorkflowStates(ctx context.Context, graphqlClient graphql.Client, limit int) (WorkflowStateList, error) {
	states, err := workflowStates(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return WorkflowStateList{}, fmt.Errorf("list workflow states: %w", err)
	}

	summaries := make([]WorkflowStateSummary, 0, len(states.WorkflowStates.Nodes))
	for _, state := range states.WorkflowStates.Nodes {
		summaries = append(summaries, workflowStateSummary(state.WorkflowStateSummaryFields))
	}

	return WorkflowStateList{
		WorkflowStates: summaries,
		HasNextPage:    states.WorkflowStates.PageInfo.HasNextPage,
		EndCursor:      states.WorkflowStates.PageInfo.EndCursor,
	}, nil
}

// GetWorkflowStateByID returns one workflow state by Linear id.
func GetWorkflowStateByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (WorkflowStateSummary, error) {
	state, err := workflowState(ctx, graphqlClient, id)
	if err != nil {
		return WorkflowStateSummary{}, fmt.Errorf("get workflow state %s: %w", id, err)
	}

	return workflowStateSummary(state.WorkflowState.WorkflowStateSummaryFields), nil
}

func workflowStateSummary(state WorkflowStateSummaryFields) WorkflowStateSummary {
	return WorkflowStateSummary{
		ID:       state.Id,
		Name:     state.Name,
		Type:     state.Type,
		Color:    state.Color,
		Position: state.Position,
		TeamID:   state.Team.Id,
		TeamKey:  state.Team.Key,
		TeamName: state.Team.Name,
	}
}
