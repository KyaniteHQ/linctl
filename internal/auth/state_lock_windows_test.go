//go:build windows

package auth

import (
	"context"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_lockStateFile_reports_context_and_non_contention_errors(t *testing.T) {
	file := tempLockFile(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.ErrorIs(t, lockStateFile(ctx, file), context.Canceled)

	want := syscall.Errno(6)
	withLockFileExCall(t, func(...uintptr) (uintptr, uintptr, error) {
		return 0, 0, want
	})
	require.ErrorIs(t, lockStateFile(context.Background(), file), want)
}

func Test_lockStateFile_retries_lock_violation_then_succeeds(t *testing.T) {
	file := tempLockFile(t)

	calls := 0
	withLockFileExCall(t, func(...uintptr) (uintptr, uintptr, error) {
		calls++
		if calls == 1 {
			return 0, 0, lockFileLockViolation
		}
		return 1, 0, nil
	})

	require.NoError(t, lockStateFile(context.Background(), file))
	require.Equal(t, 2, calls)
}

func Test_lockStateFile_cancels_while_retrying_lock_violation(t *testing.T) {
	file := tempLockFile(t)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	waiting := make(chan struct{})
	var reported bool
	withLockFileExCall(t, func(...uintptr) (uintptr, uintptr, error) {
		if !reported {
			reported = true
			close(waiting)
		}
		return 0, 0, lockFileLockViolation
	})

	result := make(chan error, 1)
	go func() {
		result <- lockStateFile(ctx, file)
	}()

	<-waiting
	cancel()
	require.ErrorIs(t, <-result, context.Canceled)
}

func Test_unlockStateFile_reports_syscall_error(t *testing.T) {
	file := tempLockFile(t)

	withUnlockFileExCall(t, func(...uintptr) (uintptr, uintptr, error) {
		return 1, 0, nil
	})
	require.NoError(t, unlockStateFile(file))

	want := syscall.Errno(6)
	withUnlockFileExCall(t, func(...uintptr) (uintptr, uintptr, error) {
		return 0, 0, want
	})
	require.ErrorIs(t, unlockStateFile(file), want)
}

func tempLockFile(t *testing.T) *os.File {
	t.Helper()
	file, err := os.CreateTemp(t.TempDir(), "lock")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, file.Close())
	})
	return file
}

func withLockFileExCall(t *testing.T, fn func(...uintptr) (uintptr, uintptr, error)) {
	t.Helper()
	original := lockFileExCall
	lockFileExCall = fn
	t.Cleanup(func() {
		lockFileExCall = original
	})
}

func withUnlockFileExCall(t *testing.T, fn func(...uintptr) (uintptr, uintptr, error)) {
	t.Helper()
	original := unlockFileExCall
	unlockFileExCall = fn
	t.Cleanup(func() {
		unlockFileExCall = original
	})
}
