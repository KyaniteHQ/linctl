package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_resolveProjectContent covers the --content / --content-file resolution:
// a file loads into content, an empty file flag is a no-op, and passing both
// content and a file is rejected.
func Test_resolveProjectContent(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "body.md")
	require.NoError(t, os.WriteFile(file, []byte("# body\n"), 0o600))

	t.Run("no file leaves content untouched", func(t *testing.T) {
		content := "inline"
		require.NoError(t, resolveProjectContent(&content, ""))
		require.Equal(t, "inline", content)
	})

	t.Run("file loads into empty content", func(t *testing.T) {
		content := ""
		require.NoError(t, resolveProjectContent(&content, file))
		require.Equal(t, "# body\n", content)
	})

	t.Run("both content and file is an error", func(t *testing.T) {
		content := "inline"
		require.Error(t, resolveProjectContent(&content, file))
	})

	t.Run("missing file is an error", func(t *testing.T) {
		content := ""
		require.Error(t, resolveProjectContent(&content, filepath.Join(dir, "missing.md")))
	})
}
