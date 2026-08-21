package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_readJSONObjectFile_returns_file_bytes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "template.json")
	raw := []byte(`{ "b": 1, "a": true }`)
	require.NoError(t, os.WriteFile(path, raw, 0o600))

	data, err := readJSONObjectFile(path)

	require.NoError(t, err)
	require.Equal(t, string(raw), string(data))
}

func Test_readJSONObjectFile_rejects_stdin_and_missing(t *testing.T) {
	_, err := readJSONObjectFile("-")
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "not stdin")

	_, err = readJSONObjectFile("")
	require.ErrorIs(t, err, client.ErrWriteInvalid)

	_, err = readJSONObjectFile(filepath.Join(t.TempDir(), "missing.json"))
	require.ErrorIs(t, err, client.ErrWriteInvalid)
}
