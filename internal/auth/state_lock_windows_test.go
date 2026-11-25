//go:build windows

package auth

import (
	"context"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	errorInvalidHandle = syscall.Errno(6)
	errorNotLocked     = syscall.Errno(158)
)

func Test_lockStateFile_reports_context_and_non_contention_errors(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "lock"))
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, lockStateFile(ctx, file), context.Canceled)

	require.NoError(t, file.Close())
	require.ErrorIs(t, lockStateFile(context.Background(), file), errorInvalidHandle)
}

func Test_lockStateFile_waits_for_a_held_lock_until_release_or_cancel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")
	holder := openLockHandle(t, path)
	waiter := openLockHandle(t, path)
	require.NoError(t, lockStateFile(context.Background(), holder))

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	require.ErrorIs(t, lockStateFile(ctx, waiter), context.DeadlineExceeded)

	released := make(chan error, 1)
	time.AfterFunc(50*time.Millisecond, func() {
		released <- unlockStateFile(holder)
	})
	require.NoError(t, lockStateFile(context.Background(), waiter))
	require.NoError(t, <-released)
	require.NoError(t, unlockStateFile(waiter))
}

func Test_unlockStateFile_reports_unlocking_an_unlocked_file(t *testing.T) {
	file := openLockHandle(t, filepath.Join(t.TempDir(), "lock"))

	require.ErrorIs(t, unlockStateFile(file), errorNotLocked)
}

func openLockHandle(t *testing.T, path string) *os.File {
	t.Helper()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, file.Close())
	})
	return file
}
