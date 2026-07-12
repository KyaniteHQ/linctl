package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// resolveBodyFlag replaces a literal "-" body value with the full stdin
// contents, so agents can pipe long text without shell-quoting it.
func resolveBodyFlag(command *cobra.Command, body *string) error {
	if *body != "-" {
		return nil
	}
	data, err := io.ReadAll(command.InOrStdin())
	if err != nil {
		return fmt.Errorf("read body from stdin: %w", err)
	}
	*body = string(data)

	return nil
}

// resolveFileFlag loads value from path when set, guarding against passing
// both the inline flag and its -file companion.
func resolveFileFlag(value *string, path string, label string) error {
	if path == "" {
		return nil
	}
	if *value != "" {
		return fmt.Errorf("%w: %s and %s-file are mutually exclusive", client.ErrWriteInvalid, label, label)
	}

	//nolint:gosec // The path is an explicit CLI input for reading issue text.
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s from file %s: %w", label, path, err)
	}
	*value = string(data)

	return nil
}
