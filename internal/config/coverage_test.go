package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_LoadScenarios_resolve_sources_and_report_config_errors(t *testing.T) {
	t.Run("context cancellation is returned before file reads", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := Load(ctx, LoadRequest{})

		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("missing files resolve empty config", func(t *testing.T) {
		config, err := Load(context.Background(), LoadRequest{
			GlobalPath: filepath.Join(t.TempDir(), "missing-global.toml"),
			RepoPath:   filepath.Join(t.TempDir(), "missing-repo.toml"),
		})

		require.NoError(t, err)
		require.Empty(t, config.Profile)
		require.Empty(t, config.Target)
	})

	t.Run("parse errors include the config path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "broken.toml")
		require.NoError(t, os.WriteFile(path, []byte("[target\n"), 0o600))

		_, err := Load(context.Background(), LoadRequest{RepoPath: path})

		require.ErrorContains(t, err, path)
	})

	t.Run("read errors include the config path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config-dir")
		require.NoError(t, os.Mkdir(path, 0o700))

		_, err := Load(context.Background(), LoadRequest{RepoPath: path})

		require.ErrorContains(t, err, "read config")
		require.ErrorContains(t, err, path)
	})

	t.Run("global config read errors stop loading", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "global-dir")
		require.NoError(t, os.Mkdir(path, 0o700))

		_, err := Load(context.Background(), LoadRequest{GlobalPath: path})

		require.ErrorContains(t, err, "read config")
		require.ErrorContains(t, err, path)
	})
}
