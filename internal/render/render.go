// Package render writes human and JSON command output.
package render

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteJSON writes a JSON value.
func WriteJSON(writer io.Writer, value any, compact bool) error {
	encoder := json.NewEncoder(writer)
	if !compact {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("write json: %w", err)
	}

	return nil
}

// WriteLine writes a single line.
func WriteLine(writer io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(writer, format+"\n", args...)
	if err != nil {
		return fmt.Errorf("write line: %w", err)
	}

	return nil
}
