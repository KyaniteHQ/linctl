//go:build linux

package auth

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Store_reports_lock_permission_errors(t *testing.T) {
	t.Parallel()

	t.Run("lock directory that refuses chmod", func(t *testing.T) {
		t.Parallel()
		linkedDir := filepath.Join(t.TempDir(), "proc-self")
		require.NoError(t, os.Symlink("/proc/self", linkedDir))

		err := NewStore(Paths{TokenPath: filepath.Join(linkedDir, tokenFileName)}).
			SaveTokenState(context.Background(), "", TokenState{AccessToken: "oauth-access-token"})

		require.ErrorIs(t, err, os.ErrPermission)
		require.ErrorContains(t, err, "secure auth token state lock directory")
	})

	t.Run("lock file that refuses chmod", func(t *testing.T) {
		t.Parallel()
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(filepath.Dir(paths.TokenPath), 0o700))
		require.NoError(t, os.Symlink("/proc/self/comm", paths.TokenPath+".lock"))

		err := NewStore(paths).SaveTokenState(context.Background(), "", TokenState{AccessToken: "oauth-access-token"})

		require.ErrorIs(t, err, os.ErrPermission)
		require.ErrorContains(t, err, "secure auth token state lock "+paths.TokenPath+".lock")
		require.NoFileExists(t, paths.TokenPath)
	})
}
