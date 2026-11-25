package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ErrPinExists marks a refuse-to-overwrite failure when the destination already
// has a pin file.
var ErrPinExists = errors.New("pin file already exists")

// pinFile is the only shape WritePin serializes. Auth material has no field
// here, so secrets cannot reach the written file by construction.
type pinFile struct {
	Target Target `toml:"target"`
}

// WritePin writes a target-only `.linctl.toml` pin to path with mode 0o644.
// It refuses to overwrite an existing path (O_EXCL).
func WritePin(path string, target Target) error {
	if path == "" {
		return errors.New("write pin: path is required")
	}
	if target.OrgID == "" || target.TeamKey == "" || target.TeamID == "" {
		return errors.New("write pin: org_id, team_key, and team_id are required")
	}

	//nolint:gosec // pin paths are explicit local destinations chosen by the caller
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf(
				"%w: %s already exists; edit or remove it, then re-run (no --force)",
				ErrPinExists,
				path,
			)
		}

		return fmt.Errorf("write pin %s: %w", path, err)
	}

	encodeErr := toml.NewEncoder(file).Encode(pinFile{Target: target})
	if err := errors.Join(encodeErr, file.Close()); err != nil {
		_ = os.Remove(path) //nolint:errcheck // a partial pin must not block the next init.
		return fmt.Errorf("write pin %s: %w", path, err)
	}

	return nil
}
