package client

import (
	"errors"
	"fmt"
)

// ErrNotFound marks an expected missing Linear entity.
var ErrNotFound = errors.New("not found")

// ErrWriteInvalid marks a malformed write request. It is shared by every
// guarded-write surface (issues, cycles, comments, documents, attachments,
// files), so it lives with the other cross-cutting sentinels rather than on
// any one entity file.
var ErrWriteInvalid = errors.New("invalid write")

func notFoundError(format string, args ...any) error {
	args = append(args, ErrNotFound)

	return fmt.Errorf(format+": %w", args...)
}
