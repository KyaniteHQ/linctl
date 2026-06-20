// Package render writes human and JSON command output.
package render

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format selects the output renderer.
type Format string

const (
	// FormatHuman writes compact human-readable text.
	FormatHuman Format = "human"
	// FormatJSON writes indented JSON.
	FormatJSON Format = "json"
	// FormatCompact writes one-line compact text.
	FormatCompact Format = "compact"
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
