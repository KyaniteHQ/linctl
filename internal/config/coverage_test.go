package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type mapEnv map[string]string

func (env mapEnv) Lookup(key string) (string, bool) {
	value, ok := env[key]
	return value, ok
}

func Test_LoadScenarios_resolve_sources_and_report_config_errors(t *testing.T) {
	t.Run("context cancellation is returned before file reads", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := Load(ctx, LoadRequest{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "load config context")
	})

	t.Run("missing files resolve empty config", func(t *testing.T) {
		config, err := Load(context.Background(), LoadRequest{
			Env:        mapEnv{"LINEAR_API_KEY": "linear-token"},
			GlobalPath: filepath.Join(t.TempDir(), "missing-global.toml"),
			RepoPath:   filepath.Join(t.TempDir(), "missing-repo.toml"),
		})

		require.NoError(t, err)
		require.Equal(t, "linear-token", config.Token)
	})

	t.Run("parse errors include the config path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "broken.toml")
		require.NoError(t, os.WriteFile(path, []byte("[target\n"), 0o600))

		_, err := Load(context.Background(), LoadRequest{RepoPath: path})

		require.Error(t, err)
		require.Contains(t, err.Error(), path)
	})

	t.Run("read errors include the config path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config-dir")
		require.NoError(t, os.Mkdir(path, 0o700))

		_, err := Load(context.Background(), LoadRequest{RepoPath: path})

		require.Error(t, err)
		require.Contains(t, err.Error(), "read config")
		require.Contains(t, err.Error(), path)
	})

	t.Run("global config read errors stop loading", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "global-dir")
		require.NoError(t, os.Mkdir(path, 0o700))

		_, err := Load(context.Background(), LoadRequest{GlobalPath: path})

		require.Error(t, err)
		require.Contains(t, err.Error(), "read config")
		require.Contains(t, err.Error(), path)
	})

	t.Run("config without profile table initializes empty profiles", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.toml")
		require.NoError(t, os.WriteFile(path, []byte(`token = "file-token"`), 0o600))

		config, err := Load(context.Background(), LoadRequest{
			Env:      mapEnv{},
			RepoPath: path,
		})

		require.NoError(t, err)
		require.Equal(t, "file-token", config.Token)
	})

	t.Run("LINCTL_TOKEN wins over LINEAR_API_KEY and profile token", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "config.toml")
		require.NoError(t, os.WriteFile(path, []byte(`
profile = "daily"

[profiles.daily]
token = "profile-token"
`), 0o600))

		config, err := Load(context.Background(), LoadRequest{
			Env:      mapEnv{"LINCTL_TOKEN": "linctl-token", "LINEAR_API_KEY": "linear-token"},
			RepoPath: path,
		})

		require.NoError(t, err)
		require.Equal(t, "linctl-token", config.Token)
	})

	t.Run("nil env uses process environment", func(t *testing.T) {
		t.Setenv("LINCTL_TOKEN", "process-token")

		config, err := Load(context.Background(), LoadRequest{})

		require.NoError(t, err)
		require.Equal(t, "process-token", config.Token)
	})
}
