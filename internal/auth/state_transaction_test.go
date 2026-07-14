package auth

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Store_token_transaction_updates_adopts_clears_and_reports_errors(t *testing.T) {
	store := NewStore(testPaths(t))
	initial := TokenState{AccessToken: "initial-access-token"}
	require.NoError(t, store.SaveTokenState(context.Background(), "", initial))

	updated, err := store.TransactTokenState(context.Background(), "", func(current TokenState) (TokenState, error) {
		require.Equal(t, initial, current)
		return TokenState{AccessToken: "updated-access-token"}, nil
	})
	require.NoError(t, err)
	require.Equal(t, "updated-access-token", updated.AccessToken)

	adopted, err := store.TransactTokenState(context.Background(), "", func(current TokenState) (TokenState, error) {
		return current, nil
	})
	require.NoError(t, err)
	require.Equal(t, updated, adopted)

	_, err = store.TransactTokenState(context.Background(), "", func(TokenState) (TokenState, error) {
		return TokenState{}, errors.New("transaction failed")
	})
	require.ErrorContains(t, err, "transaction failed")

	named, err := store.TransactTokenState(context.Background(), "work", func(current TokenState) (TokenState, error) {
		require.Empty(t, current)
		return TokenState{AccessToken: "work-access-token"}, nil
	})
	require.NoError(t, err)
	require.Equal(t, "work-access-token", named.AccessToken)
	_, err = store.TransactTokenState(context.Background(), "work", func(TokenState) (TokenState, error) {
		return TokenState{}, nil
	})
	require.NoError(t, err)

	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, updated, state.Token)
	require.Empty(t, state.Profile("work").Token)
}

func Test_Store_token_transaction_reports_context_and_parse_errors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewStore(testPaths(t)).TransactTokenState(ctx, "", func(current TokenState) (TokenState, error) {
		return current, nil
	})
	require.ErrorContains(t, err, "transact auth token state context")

	paths := testPaths(t)
	require.NoError(t, os.MkdirAll(filepath.Dir(paths.TokenPath), 0o700))
	require.NoError(t, os.WriteFile(paths.TokenPath, []byte("{"), 0o600))
	_, err = NewStore(paths).TransactTokenState(context.Background(), "", func(current TokenState) (TokenState, error) {
		return current, nil
	})
	require.ErrorContains(t, err, "parse auth token state")
}

func Test_Store_concurrent_profile_updates_preserve_every_profile(t *testing.T) {
	t.Parallel()
	const writers = 32
	store := NewStore(testPaths(t))
	start := make(chan struct{})
	errors := make(chan error, writers)
	var ready sync.WaitGroup
	ready.Add(writers)

	for index := range writers {
		go func() {
			ready.Done()
			<-start
			errors <- store.SaveTokenState(context.Background(), strconv.Itoa(index), TokenState{
				AccessToken: "oauth-access-token",
			})
		}()
	}
	ready.Wait()
	close(start)
	for range writers {
		require.NoError(t, <-errors)
	}

	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Len(t, state.Profiles, writers)
}

func Test_Store_cross_process_profile_updates_preserve_every_profile(t *testing.T) {
	const helperEnv = "LINCTL_AUTH_CONCURRENT_WRITE_HELPER"
	if os.Getenv(helperEnv) == "1" {
		runConcurrentWriteHelper(t)
		return
	}

	const writers = 16
	paths := testPaths(t)
	root := filepath.Dir(paths.TokenPath)
	startPath := filepath.Join(root, "start")
	require.NoError(t, os.MkdirAll(root, 0o700))
	commands := make([]*exec.Cmd, 0, writers)
	for index := range writers {
		profile := strconv.Itoa(index)
		command := exec.Command(os.Args[0], "-test.run=^Test_Store_cross_process_profile_updates_preserve_every_profile$")
		command.Env = append(
			os.Environ(),
			helperEnv+"=1",
			"LINCTL_AUTH_TEST_APP_PATH="+paths.AppConfigPath,
			"LINCTL_AUTH_TEST_TOKEN_PATH="+paths.TokenPath,
			"LINCTL_AUTH_TEST_PROFILE="+profile,
			"LINCTL_AUTH_TEST_START_PATH="+startPath,
			"LINCTL_AUTH_TEST_READY_PATH="+filepath.Join(root, "ready-"+profile),
		)
		require.NoError(t, command.Start())
		commands = append(commands, command)
	}
	require.Eventually(t, func() bool {
		ready, err := filepath.Glob(filepath.Join(root, "ready-*"))
		return err == nil && len(ready) == writers
	}, 10*time.Second, 10*time.Millisecond)
	require.NoError(t, os.WriteFile(startPath, []byte("start"), 0o600))
	for _, command := range commands {
		require.NoError(t, command.Wait())
	}

	state, err := NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Len(t, state.Profiles, writers)
}

func runConcurrentWriteHelper(t *testing.T) {
	t.Helper()
	readyPath := os.Getenv("LINCTL_AUTH_TEST_READY_PATH")
	startPath := os.Getenv("LINCTL_AUTH_TEST_START_PATH")
	require.NoError(t, os.WriteFile(readyPath, []byte("ready"), 0o600))
	require.Eventually(t, func() bool {
		_, err := os.Stat(startPath)
		return err == nil
	}, 10*time.Second, time.Millisecond)

	store := NewStore(Paths{
		AppConfigPath: os.Getenv("LINCTL_AUTH_TEST_APP_PATH"),
		TokenPath:     os.Getenv("LINCTL_AUTH_TEST_TOKEN_PATH"),
	})
	require.NoError(t, store.SaveTokenState(
		context.Background(),
		os.Getenv("LINCTL_AUTH_TEST_PROFILE"),
		TokenState{AccessToken: "oauth-access-token"},
	))
}

func Test_Store_expired_save_waiting_for_file_lock_does_not_write(t *testing.T) {
	paths := testPaths(t)
	store := NewStore(paths)
	initial := TokenState{AccessToken: "initial-access-token"}
	require.NoError(t, store.SaveTokenState(context.Background(), "", initial))

	lockStarted := make(chan struct{})
	releaseLock := make(chan struct{})
	withLockAuthFile(t, func(ctx context.Context, _ *os.File) error {
		close(lockStarted)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-releaseLock:
			return nil
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		result <- store.SaveTokenState(ctx, "", TokenState{AccessToken: "late-access-token"})
	}()

	<-lockStarted
	cancel()
	require.ErrorIs(t, <-result, context.Canceled)
	close(releaseLock)

	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, initial, state.Token)
}

func Test_Store_expired_save_after_process_lock_does_not_write(t *testing.T) {
	paths := testPaths(t)
	store := NewStore(paths)
	initial := TokenState{AccessToken: "initial-access-token"}
	require.NoError(t, store.SaveTokenState(context.Background(), "", initial))

	ctx, cancel := context.WithCancel(context.Background())
	withLockAuthFile(t, func(context.Context, *os.File) error {
		cancel()
		return nil
	})
	err := store.SaveTokenState(ctx, "", TokenState{AccessToken: "late-access-token"})
	require.ErrorIs(t, err, context.Canceled)

	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, initial, state.Token)
}

func Test_withAuthFileLock_cancels_while_waiting_for_process_lock_without_leaking(t *testing.T) {
	paths := testPaths(t)
	processLock := authFileLock(paths.TokenPath)
	require.NoError(t, processLock.acquire(context.Background()))
	held := true
	t.Cleanup(func() {
		if held {
			processLock.release()
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	result := make(chan error, 1)
	var actionCalled atomic.Bool
	go func() {
		close(started)
		result <- withAuthFileLock(ctx, paths.TokenPath, "auth token state", "test process wait", func() error {
			actionCalled.Store(true)
			return nil
		})
	}()

	<-started
	require.Never(t, func() bool {
		select {
		case <-result:
			return true
		default:
			return false
		}
	}, 20*time.Millisecond, time.Millisecond)
	cancel()
	require.ErrorIs(t, <-result, context.Canceled)
	require.False(t, actionCalled.Load())

	processLock.release()
	held = false
	require.NoError(t, withAuthFileLock(
		context.Background(),
		paths.TokenPath,
		"auth token state",
		"test process retry",
		func() error {
			actionCalled.Store(true)
			return nil
		},
	))
	require.True(t, actionCalled.Load())
}

func Test_withAuthFileLock_reports_lock_setup_errors(t *testing.T) {
	t.Run("secure lock directory", func(t *testing.T) {
		withRuntimeGOOS(t, "linux")
		withChmodFile(t, func(string, os.FileMode) error {
			return errors.New("chmod failed")
		})

		err := NewStore(testPaths(t)).SaveTokenState(context.Background(), "", TokenState{})

		require.ErrorContains(t, err, "secure auth token state lock directory")
	})

	t.Run("open lock file", func(t *testing.T) {
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(paths.TokenPath+".lock", 0o700))

		err := NewStore(paths).SaveTokenState(context.Background(), "", TokenState{})

		require.ErrorContains(t, err, "open auth token state lock")
	})

	t.Run("secure lock file", func(t *testing.T) {
		withRuntimeGOOS(t, "linux")
		calls := 0
		withChmodFile(t, func(string, os.FileMode) error {
			calls++
			if calls == 2 {
				return errors.New("chmod failed")
			}
			return nil
		})

		err := NewStore(testPaths(t)).SaveTokenState(context.Background(), "", TokenState{})

		require.ErrorContains(t, err, "secure auth token state lock")
	})

	t.Run("acquire lock", func(t *testing.T) {
		unlockCalled := false
		withLockAuthFile(t, func(context.Context, *os.File) error {
			return errors.New("lock failed")
		})
		withUnlockAuthFile(t, func(*os.File) error {
			unlockCalled = true
			return errors.New("bogus unlock")
		})

		err := NewStore(testPaths(t)).SaveTokenState(context.Background(), "", TokenState{})

		require.ErrorContains(t, err, "acquire auth token state lock")
		require.NotContains(t, err.Error(), "bogus unlock")
		require.False(t, unlockCalled)
	})

	t.Run("release lock", func(t *testing.T) {
		withUnlockAuthFile(t, func(*os.File) error {
			return errors.New("unlock failed")
		})

		err := NewStore(testPaths(t)).SaveTokenState(context.Background(), "", TokenState{})

		require.ErrorContains(t, err, "unlock failed")
	})
}
