//go:build !windows

package auth

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Store_cancels_while_waiting_for_OS_lock_without_writing_or_leaking(t *testing.T) {
	paths := testPaths(t)
	store := NewStore(paths)
	initial := TokenState{AccessToken: "initial-access-token"}
	require.NoError(t, store.SaveTokenState(context.Background(), "", initial))

	lockFile, err := os.OpenFile(paths.TokenPath+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	require.NoError(t, err)
	require.NoError(t, syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB))
	held := true
	t.Cleanup(func() {
		if held {
			require.NoError(t, syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN))
		}
		require.NoError(t, lockFile.Close())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err = store.SaveTokenState(ctx, "", TokenState{AccessToken: "late-access-token"})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, initial, state.Token)

	require.NoError(t, syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN))
	held = false
	replacement := TokenState{AccessToken: "replacement-access-token"}
	require.NoError(t, store.SaveTokenState(context.Background(), "", replacement))
	state, err = store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, replacement, state.Token)
}

func Test_lockStateFile_reports_context_and_non_contention_errors(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "lock")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, lockStateFile(ctx, file), context.Canceled)

	require.NoError(t, file.Close())
	require.ErrorIs(t, lockStateFile(context.Background(), file), syscall.EBADF)
}

func Test_Store_failed_write_keeps_existing_state_and_leaves_no_temp_file(t *testing.T) {
	paths := testPaths(t)
	store := NewStore(paths)
	initial := TokenState{AccessToken: "initial-access-token"}
	require.NoError(t, store.SaveTokenState(context.Background(), "", initial))
	var original syscall.Rlimit
	require.NoError(t, syscall.Getrlimit(syscall.RLIMIT_FSIZE, &original))
	limited := original
	limited.Cur = 0
	require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &limited))
	restore := func() { require.NoError(t, syscall.Setrlimit(syscall.RLIMIT_FSIZE, &original)) }
	t.Cleanup(restore)

	err := store.SaveTokenState(context.Background(), "", TokenState{AccessToken: "new-access-token"})
	restore()

	require.ErrorIs(t, err, syscall.EFBIG)
	require.ErrorContains(t, err, "write auth token state "+paths.TokenPath)
	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, initial, state.Token)
	entries, err := os.ReadDir(filepath.Dir(paths.TokenPath))
	require.NoError(t, err)
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	require.Equal(t, []string{tokenFileName, tokenFileName + ".lock"}, names)
}

func Test_Store_reports_temp_file_creation_error(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	appConfigPath := filepath.Join(root, strings.Repeat("a", 250))

	err := NewStore(Paths{AppConfigPath: appConfigPath}).Save(context.Background(), State{})

	require.ErrorIs(t, err, syscall.ENAMETOOLONG)
	require.ErrorContains(t, err, "write auth app config "+appConfigPath)
	require.NoFileExists(t, appConfigPath)
}

func Test_Store_saves_readable_private_files_under_a_restrictive_umask(t *testing.T) {
	paths := testPaths(t)
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Dir(paths.AppConfigPath)), 0o700))
	require.NoError(t, os.MkdirAll(filepath.Dir(filepath.Dir(paths.TokenPath)), 0o700))
	previous := syscall.Umask(0o477)
	t.Cleanup(func() { syscall.Umask(previous) })
	store := NewStore(paths)
	saved := State{
		App:   AppConfig{ClientID: "client-id"},
		Token: TokenState{AccessToken: "oauth-access-token"},
	}

	require.NoError(t, store.Save(context.Background(), saved))

	loaded, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, saved.App, loaded.App)
	require.Equal(t, saved.Token, loaded.Token)
	assertPrivateDirAndFile(t, filepath.Dir(paths.AppConfigPath), paths.AppConfigPath)
	assertPrivateDirAndFile(t, filepath.Dir(paths.TokenPath), paths.TokenPath)
}
