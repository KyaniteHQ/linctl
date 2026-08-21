package client

import (
	"context"
	"fmt"
)

func (guard *guardedClient) resolveNamedOrTypedState(
	ctx context.Context,
	teamID string,
	name string,
	stateType string,
) (stateID string, set bool, err error) {
	if name != "" {
		stateID, err = guard.resolveStateID(ctx, teamID, name)

		return stateID, true, err
	}
	if stateType != "" {
		stateID, err = guard.resolveStateTypeID(ctx, teamID, stateType)

		return stateID, true, err
	}

	return "", false, nil
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
		if writeErr != nil {
			return applyMutationRetryClass(
				IssueStateWriteRetryClass(), observed.Summary, true, writeErr,
			)
		}

		return observed.Summary, nil
	}

	mismatch := fmt.Errorf(
		"%w: expected state_id=%s resolved state_id=%s name=%q",
		ErrStateMismatch,
		wantStateID,
		observed.Summary.StateID,
		observed.Summary.State,
	)
	if writeErr != nil {
		return IssueSummary{}, fmt.Errorf("%w: %w", mismatch, writeErr)
	}

	return IssueSummary{}, mismatch
}
