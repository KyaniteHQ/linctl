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

// requiredFieldError names a missing write-request field with the shared
// ErrWriteInvalid sentinel, so every guarded-write surface words a missing
// field the same way.
func requiredFieldError(field string) error {
	return fmt.Errorf("%w: %s is required", ErrWriteInvalid, field)
}

func notFoundError(format string, args ...any) error {
	args = append(args, ErrNotFound)

	return fmt.Errorf(format+": %w", args...)
}
