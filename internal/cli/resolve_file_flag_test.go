package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// Test_resolveFileFlag covers the --content / --content-file resolution:
// a file loads into the value, an empty file flag is a no-op, passing both
// content and a file is rejected with ErrWriteInvalid, and a missing file errors.
func Test_resolveFileFlag(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "body.md")
	require.NoError(t, os.WriteFile(file, []byte("# body\n"), 0o600))

	t.Run("no file leaves value untouched", func(t *testing.T) {
		value := "inline"
		require.NoError(t, resolveFileFlag(&value, "", "content"))
		require.Equal(t, "inline", value)
	})

	t.Run("file loads into empty value", func(t *testing.T) {
		value := ""
		require.NoError(t, resolveFileFlag(&value, file, "content"))
		require.Equal(t, "# body\n", value)
	})

	t.Run("both value and file is ErrWriteInvalid", func(t *testing.T) {
		value := "inline"
		require.ErrorIs(t, resolveFileFlag(&value, file, "content"), client.ErrWriteInvalid)
	})

	t.Run("missing file is an error", func(t *testing.T) {
		value := ""
		require.Error(t, resolveFileFlag(&value, filepath.Join(dir, "missing.md"), "content"))
	})
}
