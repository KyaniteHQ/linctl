package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// workflowStateTypeAliases maps human state words to Linear WorkflowState.type
// values from the vendored schema. Exact state names are resolved separately.
var workflowStateTypeAliases = map[string]string{
	"triage":      "triage",
	"backlog":     "backlog",
	"unstarted":   "unstarted",
	"started":     "started",
	"completed":   "completed",
	"canceled":    "canceled",
	"todo":        "unstarted",
	"to do":       "unstarted",
	"in progress": "started",
	"in-progress": "started",
	"done":        "completed",
	"complete":    "completed",
	"closed":      "completed",
	"cancelled":   "canceled",
	"wont do":     "canceled",
	"wont-do":     "canceled",
	"won't do":    "canceled",
}

type workflowStateCandidate struct {
	ID       string
	Name     string
	Type     string
	Position float64
}

const workflowStatePageSize = 50

// CanonicalWorkflowStateType maps a human state word to a Linear workflow
// state type. It does not resolve a team state id.
func CanonicalWorkflowStateType(raw string) (string, bool) {
	canonical, ok := workflowStateTypeAliases[strings.ToLower(strings.TrimSpace(raw))]

	return canonical, ok
}

func (guard *guardedClient) resolveStateID(
	ctx context.Context,
	teamID string,
	selector string,
) (string, error) {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return "", requiredFieldError("state")
	}
	states, err := guard.teamWorkflowStates(ctx, teamID)
	if err != nil {
		return "", err
	}

	return selectWorkflowStateID(states, teamID, selector)
}

func (guard *guardedClient) teamWorkflowStates(
	ctx context.Context,
	teamID string,
) ([]workflowStateCandidate, error) {
	guard.stateIDs.mu.Lock()
	defer guard.stateIDs.mu.Unlock()
	if states, ok := guard.stateIDs.lists[teamID]; ok {
		return states, nil
	}

	states, err := listTeamWorkflowStates(ctx, guard.graphqlClient, teamID)
	if err != nil {
		return nil, err
	}
	guard.stateIDs.lists[teamID] = states

	return states, nil
}

func listTeamWorkflowStates(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
) ([]workflowStateCandidate, error) {
	result, err := gql.WorkflowStatesByTeam(ctx, graphqlClient, teamID, intPtr(workflowStatePageSize))
	if err != nil {
		return nil, fmt.Errorf("list workflow states for team_id=%s: %w", teamID, err)
	}
	if result.WorkflowStates.PageInfo.HasNextPage {
		return nil, fmt.Errorf(
			"%w: team team_id=%s has more than %d workflow states; cannot resolve a state uniquely",
			ErrWriteInvalid,
			teamID,
			workflowStatePageSize,
		)
	}

	states := make([]workflowStateCandidate, 0, len(result.WorkflowStates.Nodes))
	for _, node := range result.WorkflowStates.Nodes {
		states = append(states, workflowStateCandidate{
			ID:       node.Id,
			Name:     node.Name,
			Type:     node.Type,
			Position: node.Position,
		})
	}

	return states, nil
}

func selectWorkflowStateID(
	states []workflowStateCandidate,
	teamID string,
	selector string,
) (string, error) {
	if id, ok, err := uniqueStateNameID(states, selector); err != nil || ok {
		return id, err
	}
	stateType, ok := CanonicalWorkflowStateType(selector)
	if !ok {
		return "", fmt.Errorf("%w: unknown workflow state %q", ErrWriteInvalid, selector)
	}

	return firstStateIDOfCandidates(states, teamID, stateType)
}

func uniqueStateNameID(
	states []workflowStateCandidate,
	selector string,
) (string, bool, error) {
	want := strings.ToLower(selector)
	var matches []workflowStateCandidate
	for _, state := range states {
		if strings.ToLower(state.Name) == want {
			matches = append(matches, state)
		}
	}
	if len(matches) == 0 {
		return "", false, nil
	}
	if len(matches) > 1 {
		return "", false, fmt.Errorf(
			"%w: workflow state name %q is ambiguous",
			ErrWriteInvalid,
			selector,
		)
	}

	return matches[0].ID, true, nil
}

func firstStateIDOfCandidates(
	states []workflowStateCandidate,
	teamID string,
	stateType string,
) (string, error) {
	var selected *workflowStateCandidate
	for index := range states {
		candidate := states[index]
		if candidate.Type != stateType {
			continue
		}
		if selected == nil || candidate.Position < selected.Position {
			selected = &states[index]
		}
	}
	if selected == nil {
		return "", fmt.Errorf(
			"%w: %s workflow state missing for team_id=%s",
			ErrWriteInvalid,
			stateType,
			teamID,
		)
	}

	return selected.ID, nil
}
