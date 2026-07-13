package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// Test_resolveFileFlag covers the --content / --content-file resolution:
// a file loads into the value, an empty file flag is a no-op, passing both
// content and a file is rejected with ErrWriteInvalid, a missing file errors,
// and a path of "-" reads stdin instead of a literal file named "-".
func Test_resolveFileFlag(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "body.md")
	require.NoError(t, os.WriteFile(file, []byte("# body\n"), 0o600))

	t.Run("no file leaves value untouched", func(t *testing.T) {
		value := "inline"
		require.NoError(t, resolveFileFlag(&cobra.Command{}, &value, "", "content"))
		require.Equal(t, "inline", value)
	})

	t.Run("file loads into empty value", func(t *testing.T) {
		value := ""
		require.NoError(t, resolveFileFlag(&cobra.Command{}, &value, file, "content"))
		require.Equal(t, "# body\n", value)
	})

	t.Run("both value and file is ErrWriteInvalid", func(t *testing.T) {
		value := "inline"
		require.ErrorIs(t, resolveFileFlag(&cobra.Command{}, &value, file, "content"), client.ErrWriteInvalid)
	})

	t.Run("missing file is an error", func(t *testing.T) {
		value := ""
		require.Error(t, resolveFileFlag(&cobra.Command{}, &value, filepath.Join(dir, "missing.md"), "content"))
	})

	t.Run("dash path reads stdin", func(t *testing.T) {
		value := ""
		command := &cobra.Command{}
		command.SetIn(strings.NewReader("stdin content"))
		require.NoError(t, resolveFileFlag(command, &value, "-", "content"))
		require.Equal(t, "stdin content", value)
	})

	t.Run("dash path with value set is ErrWriteInvalid without reading stdin", func(t *testing.T) {
		value := "inline"
		command := &cobra.Command{}
		command.SetIn(failingReader{err: errors.New("stdin must not be read when the conflict check fails first")})
		require.ErrorIs(t, resolveFileFlag(command, &value, "-", "content"), client.ErrWriteInvalid)
	})

	t.Run("dash path stdin read error is reported", func(t *testing.T) {
		value := ""
		command := &cobra.Command{}
		command.SetIn(failingReader{err: errors.New("stdin broken")})
		require.ErrorContains(t, resolveFileFlag(command, &value, "-", "content"), "read content from stdin")
	})
}
