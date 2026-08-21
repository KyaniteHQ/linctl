package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_readJSONObjectFile_reads_canonical_object(t *testing.T) {
	path := filepath.Join(t.TempDir(), "template.json")
	require.NoError(t, os.WriteFile(path, []byte(`{ "b": 1, "a": true }`), 0o600))

	data, err := readJSONObjectFile(path)

	require.NoError(t, err)
	require.Equal(t, `{"a":true,"b":1}`, string(data))
}

func Test_readJSONObjectFile_rejects_stdin_missing_and_non_objects(t *testing.T) {
	_, err := readJSONObjectFile("-")
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "not stdin")

	_, err = readJSONObjectFile("")
	require.ErrorIs(t, err, client.ErrWriteInvalid)

	_, err = readJSONObjectFile(filepath.Join(t.TempDir(), "missing.json"))
	require.ErrorIs(t, err, client.ErrWriteInvalid)

	path := filepath.Join(t.TempDir(), "array.json")
	require.NoError(t, os.WriteFile(path, []byte(`[]`), 0o600))
	_, err = readJSONObjectFile(path)
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.NotContains(t, err.Error(), "[")
}
