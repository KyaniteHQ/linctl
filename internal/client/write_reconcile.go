package client

import "errors"

func finishReconciledWrite[T any](
	class MutationRetryClass,
	observed T,
	readErr error,
	writeErr error,
	teamErr error,
	matches bool,
	conflict error,
) (T, error) {
	var zero T
	if readErr != nil {
		if writeErr != nil && errors.Is(readErr, ErrNotFound) {
			return zero, writeErr
		}

		return zero, withMutationError(readErr, writeErr)
	}
	if teamErr != nil {
		return zero, withMutationError(teamErr, writeErr)
	}
	if !matches {
		return zero, withMutationError(conflict, writeErr)
	}
	if writeErr != nil {
		return applyMutationRetryClass(class, observed, true, writeErr)
	}

	return observed, nil
}

func withMutationError(outcome error, writeErr error) error {
	if writeErr == nil {
		return outcome
	}

	return &ReconcileError{Outcome: outcome, MutationErr: writeErr}
}
