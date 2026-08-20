package client

import (
	"context"
	"fmt"
)

func issueCreateStateSelector(request IssueCreateRequest) string {
	if request.StateSelector != "" {
		return request.StateSelector
	}

	return request.StateType
}

func issueUpdateStateSelector(request IssueUpdateRequest) string {
	if request.StateSelector != "" {
		return request.StateSelector
	}

	return request.StateType
}

func (guard *guardedClient) resolveRequestStateID(
	ctx context.Context,
	teamID string,
	selector string,
) (stateID string, set bool, err error) {
	if selector == "" {
		return "", false, nil
	}
	stateID, err = guard.resolveStateID(ctx, teamID, selector)
	if err != nil {
		return "", false, err
	}

	return stateID, true, nil
}

func (guard *guardedClient) finishStateWrite(
	ctx context.Context,
	issueID string,
	wantStateID string,
	writeErr error,
) (IssueSummary, error) {
	observed, err := GetIssueDetail(ctx, guard.graphqlClient, issueID)
	if err != nil {
		if writeErr != nil {
			return IssueSummary{}, writeErr
		}

		return IssueSummary{}, err
	}
	if observed.Summary.StateID == wantStateID {
		return observed.Summary, nil
	}
	if writeErr != nil {
		return IssueSummary{}, writeErr
	}

	return IssueSummary{}, fmt.Errorf(
		"%w: expected state_id=%s resolved state_id=%s name=%q",
		ErrStateMismatch,
		wantStateID,
		observed.Summary.StateID,
		observed.Summary.State,
	)
}
