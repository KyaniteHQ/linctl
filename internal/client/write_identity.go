package client

import (
	"fmt"
	"regexp"
)

var uuidV4Pattern = regexp.MustCompile(
	`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func requireUUIDv4(id string, field string) error {
	if !uuidV4Pattern.MatchString(id) {
		return fmt.Errorf("%w: %s must be a UUID v4", ErrWriteInvalid, field)
	}

	return nil
}

func writeConflictError(entity string, id string) error {
	return fmt.Errorf("%w: %s %s does not match the requested fields", ErrWriteConflict, entity, id)
}
