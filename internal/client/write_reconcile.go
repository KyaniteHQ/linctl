package client

import (
	"errors"
	"fmt"
)

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

		return zero, readErr
	}
	if teamErr != nil {
		return zero, teamErr
	}
	if !matches {
		if writeErr != nil {
			return zero, fmt.Errorf("%w; the mutation returned: %w", conflict, writeErr)
		}

		return zero, conflict
	}
	if writeErr != nil {
		return applyMutationRetryClass(class, observed, true, writeErr)
	}

	return observed, nil
}
