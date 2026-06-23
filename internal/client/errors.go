package client

import (
	"errors"
	"fmt"
)

// ErrNotFound marks an expected missing Linear entity.
var ErrNotFound = errors.New("not found")

func notFoundError(format string, args ...any) error {
	args = append(args, ErrNotFound)

	return fmt.Errorf(format+": %w", args...)
}
