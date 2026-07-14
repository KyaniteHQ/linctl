//go:build !windows

package auth

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"

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

	originalFlock := flockFile
	waiting := make(chan struct{})
	var reported bool
	flockFile = func(fd int, how int) error {
		err := originalFlock(fd, how)
		if !reported && (errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN)) {
			reported = true
			close(waiting)
		}
		return err
	}
	t.Cleanup(func() {
		flockFile = originalFlock
	})

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		result <- store.SaveTokenState(ctx, "", TokenState{AccessToken: "late-access-token"})
	}()

	<-waiting
	cancel()
	require.ErrorIs(t, <-result, context.Canceled)
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
