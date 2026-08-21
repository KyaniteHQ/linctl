package cli

import (
	"fmt"
	"os"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func readJSONObjectFile(path string) ([]byte, error) {
	if path == "" {
		return nil, requiredDataFileError()
	}
	if path == "-" {
		return nil, fmt.Errorf("%w: --data-file must be a local file, not stdin", client.ErrWriteInvalid)
	}
	//nolint:gosec // The path is an explicit CLI input for reading template JSON.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: read data file %s: %w", client.ErrWriteInvalid, path, err)
	}

	return data, nil
}

func requiredDataFileError() error {
	return fmt.Errorf("%w: data-file is required", client.ErrWriteInvalid)
}
