package client

import (
	"context"
	"fmt"
	"time"
)

// stateReadbackDelays are the waits before each re-read of a state write that
// still shows the old state. Linear's read path can trail an accepted
// issueUpdate: a live move from Triage to Needs Plan read back Triage while the
// data change webhook already carried Needs Plan.
var stateReadbackDelays = []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second}

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
	observed, err := guard.readBackState(ctx, issueID, wantStateID)
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

// readBackState reads the issue, and re-reads it after each delay while the
// state is not the wanted one yet.
func (guard *guardedClient) readBackState(
	ctx context.Context,
	issueID string,
	wantStateID string,
) (IssueDetail, error) {
	observed, err := GetIssueDetail(ctx, guard.graphqlClient, issueID)
	for _, delay := range stateReadbackDelays {
		if err != nil || observed.Summary.StateID == wantStateID {
			break
		}
		select {
		case <-ctx.Done():
			return observed, nil
		case <-time.After(delay):
		}
		observed, err = GetIssueDetail(ctx, guard.graphqlClient, issueID)
	}

	return observed, err
}
