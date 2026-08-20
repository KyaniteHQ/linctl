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

// ErrStateMismatch marks a state write whose readback did not land in the
// selected workflow state. It is a hard stop: linctl does not replay the
// mutation.
var ErrStateMismatch = errors.New("state mismatch")

// ErrCrossOrganizationRelation marks a relation whose endpoints are not in the
// pinned organization. It wraps ErrTargetMismatch. It is a hard stop, not a
// retry or bypass path. The JSON code CROSS_ORGANIZATION_RELATION names this
// one boundary; ADR 0001 keeps every other guarded-write failure as
// TARGET_MISMATCH.
var ErrCrossOrganizationRelation = errors.New("cross-organization relation")

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
