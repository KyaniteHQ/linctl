package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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

type pinWriter interface {
	io.Writer
	Close() error
}

var (
	openPinFile        = openPinFileDefault
	marshalPin         = toml.Marshal
	openPinFileDefault = func(path string, flag int, perm os.FileMode) (pinWriter, error) {
		//nolint:gosec // pin paths are explicit local destinations chosen by the caller
		return os.OpenFile(path, flag, perm)
	}
)

// WritePin writes a target-only `.linctl.toml` pin to path with mode 0o644.
// It refuses to overwrite an existing path (O_EXCL).
func WritePin(path string, target Target) error {
	if path == "" {
		return errors.New("write pin: path is required")
	}
	if target.OrgID == "" || target.TeamKey == "" || target.TeamID == "" {
		return errors.New("write pin: org_id, team_key, and team_id are required")
	}

	payload, err := marshalPin(pinFile{Target: target})
	if err != nil {
		return fmt.Errorf("write pin: encode: %w", err)
	}

	file, err := openPinFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
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
	defer file.Close() //nolint:errcheck // close error checked after write

	if _, err := file.Write(payload); err != nil {
		return fmt.Errorf("write pin %s: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("write pin %s: %w", path, err)
	}

	return nil
}

// PinFileKeys documents the only TOML keys WritePin may emit under [target].
// A drift test fails if Target tags or this allowlist diverge.
func PinFileKeys() []string {
	return []string{"org_id", "team_key", "team_id", "project_id"}
}

// SerializedPinKeys returns the key names present under [target] in a pin file body.
func SerializedPinKeys(body []byte) ([]string, error) {
	var decoded pinFile
	if err := toml.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}
	keys := make([]string, 0, 4)
	if decoded.Target.OrgID != "" {
		keys = append(keys, "org_id")
	}
	if decoded.Target.TeamKey != "" {
		keys = append(keys, "team_key")
	}
	if decoded.Target.TeamID != "" {
		keys = append(keys, "team_id")
	}
	if decoded.Target.ProjectID != "" {
		keys = append(keys, "project_id")
	}
	if strings.Contains(string(body), "token") || strings.Contains(string(body), "secret") {
		return nil, errors.New("serialized pin must not contain auth material")
	}

	return keys, nil
}
