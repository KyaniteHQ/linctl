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

// resolveBodyOrFileFlag resolves a text value that a command accepts through
// an inline flag or its -file companion. The file flag runs first, so passing
// both flags always reports the mutual-exclusion error before either one
// consumes stdin.
func resolveBodyOrFileFlag(command *cobra.Command, value *string, path string, label string) error {
	if err := resolveFileFlag(command, value, path, label); err != nil {
		return err
	}

	return resolveBodyFlag(command, value)
}

// resolveFileFlag loads value from path when set, guarding against passing
// both the inline flag and its -file companion. A path of "-" reads stdin
// instead of a literal file named "-", matching the inline flag's own "-"
// convention.
func resolveFileFlag(command *cobra.Command, value *string, path string, label string) error {
	if path == "" {
		return nil
	}
	if *value != "" {
		return fmt.Errorf("%w: %s and %s-file are mutually exclusive", client.ErrWriteInvalid, label, label)
	}
	if path == "-" {
		data, err := io.ReadAll(command.InOrStdin())
		if err != nil {
			return fmt.Errorf("read %s from stdin: %w", label, err)
		}
		*value = string(data)

		return nil
	}

	//nolint:gosec // The path is an explicit CLI input for reading issue text.
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s from file %s: %w", label, path, err)
	}
	*value = string(data)

	return nil
}
