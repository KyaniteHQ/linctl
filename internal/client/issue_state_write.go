package client

import (
	"context"
	"fmt"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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

// finishStateWrite confirms a state write with one read. written is the issue the
// mutation returned, or nil when it returned none. Linear's read path can trail an
// accepted write: a live move from Triage to Needs Plan read back Triage while the
// mutation and the data change webhook already carried Needs Plan. A successful
// write whose own payload shows the wanted state is therefore confirmed by that
// payload when the read still shows another state.
func (guard *guardedClient) finishStateWrite(
	ctx context.Context,
	issueID string,
	wantStateID string,
	written *gql.IssueSummaryFields,
	writeErr error,
) (IssueSummary, error) {
	landed := writeErr == nil && written != nil && written.State.Id == wantStateID
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
	if landed {
		return issueSummaryFromFields(*written), nil
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
