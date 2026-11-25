package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_readJSONObjectFile_returns_file_bytes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "template.json")
	raw := []byte(`{ "b": 1, "a": true }`)
	require.NoError(t, os.WriteFile(path, raw, 0o600))

	data, err := readJSONObjectFile(path)

	require.NoError(t, err)
	require.Equal(t, string(raw), string(data))
}
