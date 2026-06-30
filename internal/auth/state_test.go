package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type staticEnv map[string]string

func (env staticEnv) Lookup(key string) (string, bool) {
	value, ok := env[key]
	return value, ok
}

func Test_Store_persists_local_auth_state(t *testing.T) {
	t.Parallel()
	paths := testPaths(t)
	store := NewStore(paths)
	want := State{
		App: AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Token: TokenState{
			AccessToken:  "oauth-access-token",
			RefreshToken: "refresh-token",
		},
		Profiles: map[string]ProfileState{
			"work": {
				App: AppConfig{
					ClientID: "work-client-id",
				},
				Token: TokenState{
					AccessToken: "work-oauth-access-token",
				},
			},
		},
	}

	require.NoError(t, store.Save(context.Background(), want))

	got, err := store.Load(context.Background())

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func Test_Store_save_tightens_local_auth_file_permissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows file modes do not expose unix group/world permission bits")
	}
	t.Parallel()
	root := t.TempDir()
	paths := Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", appConfigFileName),
		TokenPath:     filepath.Join(root, "state", "linctl", tokenFileName),
	}
	require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Dir(paths.TokenPath), 0o755))
	require.NoError(t, os.WriteFile(paths.AppConfigPath, []byte("{}\n"), 0o644))
	require.NoError(t, os.WriteFile(paths.TokenPath, []byte("{}\n"), 0o644))
	store := NewStore(paths)

	require.NoError(t, store.Save(context.Background(), State{
		App:   AppConfig{ClientID: "client-id"},
		Token: TokenState{AccessToken: "oauth-access-token"},
	}))

	assertPrivateDirAndFile(t, filepath.Dir(paths.AppConfigPath), paths.AppConfigPath)
	assertPrivateDirAndFile(t, filepath.Dir(paths.TokenPath), paths.TokenPath)
}

func Test_Store_clear_token_state_preserves_oauth_app_config(t *testing.T) {
	t.Parallel()
	store := NewStore(testPaths(t))
	require.NoError(t, store.Save(context.Background(), State{
		App: AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Token: TokenState{
			AccessToken:  "oauth-access-token",
			RefreshToken: "refresh-token",
		},
		Profiles: map[string]ProfileState{
			"work": {
				App:   AppConfig{ClientID: "work-client-id"},
				Token: TokenState{AccessToken: "work-oauth-access-token"},
			},
		},
	}))

	require.NoError(t, store.ClearTokenState(context.Background(), ""))
	require.NoError(t, store.ClearTokenState(context.Background(), ""))
	require.NoError(t, store.ClearTokenState(context.Background(), "work"))
	require.NoError(t, store.ClearTokenState(context.Background(), "work"))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}, got.App)
	require.Empty(t, got.Token)
	require.Equal(t, AppConfig{ClientID: "work-client-id"}, got.Profiles["work"].App)
	require.Empty(t, got.Profiles["work"].Token)
}

func Test_Store_clear_app_config_preserves_oauth_token_state(t *testing.T) {
	t.Parallel()
	store := NewStore(testPaths(t))
	require.NoError(t, store.Save(context.Background(), State{
		App: AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Token: TokenState{
			AccessToken:  "oauth-access-token",
			RefreshToken: "refresh-token",
		},
		Profiles: map[string]ProfileState{
			"work": {
				App:   AppConfig{ClientID: "work-client-id"},
				Token: TokenState{AccessToken: "work-oauth-access-token"},
			},
		},
	}))

	require.NoError(t, store.ClearAppConfig(context.Background(), ""))
	require.NoError(t, store.ClearAppConfig(context.Background(), ""))
	require.NoError(t, store.ClearAppConfig(context.Background(), "work"))
	require.NoError(t, store.ClearAppConfig(context.Background(), "work"))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got.App)
	require.Equal(t, TokenState{
		AccessToken:  "oauth-access-token",
		RefreshToken: "refresh-token",
	}, got.Token)
	require.Empty(t, got.Profiles["work"].App)
	require.Equal(t, TokenState{AccessToken: "work-oauth-access-token"}, got.Profiles["work"].Token)
}

func Test_Store_saves_app_config_without_creating_token_state(t *testing.T) {
	t.Parallel()
	paths := testPaths(t)
	store := NewStore(paths)

	require.NoError(t, store.SaveAppConfig(context.Background(), "", AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		Scopes:       []string{"read", "write"},
	}))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		Scopes:       []string{"read", "write"},
	}, got.App)
	require.Empty(t, got.Token)
	_, err = os.Stat(paths.TokenPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func Test_Store_saves_token_state_without_rewriting_app_config(t *testing.T) {
	t.Parallel()
	store := NewStore(testPaths(t))
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, store.SaveAppConfig(context.Background(), "", AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}))

	require.NoError(t, store.SaveTokenState(context.Background(), "", TokenState{
		AccessToken: "oauth-access-token",
		TokenType:   "Bearer",
		Scopes:      []string{"read"},
		ExpiresAt:   &expiresAt,
		Actor:       "app",
		GrantType:   "client_credentials",
	}))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}, got.App)
	require.Equal(t, TokenState{
		AccessToken: "oauth-access-token",
		TokenType:   "Bearer",
		Scopes:      []string{"read"},
		ExpiresAt:   &expiresAt,
		Actor:       "app",
		GrantType:   "client_credentials",
	}, got.Token)
}

func Test_Store_saves_named_profile_app_config_and_token_state(t *testing.T) {
	t.Parallel()
	store := NewStore(testPaths(t))

	require.NoError(t, store.SaveAppConfig(context.Background(), "work", AppConfig{
		ClientID:     "work-client-id",
		ClientSecret: "work-client-secret",
		Scopes:       []string{"read"},
	}))
	require.NoError(t, store.SaveTokenState(context.Background(), "work", TokenState{
		AccessToken: "work-oauth-access-token",
		TokenType:   "Bearer",
	}))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got.App)
	require.Empty(t, got.Token)
	require.Equal(t, AppConfig{
		ClientID:     "work-client-id",
		ClientSecret: "work-client-secret",
		Scopes:       []string{"read"},
	}, got.Profiles["work"].App)
	require.Equal(t, TokenState{
		AccessToken: "work-oauth-access-token",
		TokenType:   "Bearer",
	}, got.Profiles["work"].Token)
}

func Test_Store_empty_paths_are_noop(t *testing.T) {
	t.Parallel()
	store := NewStore(Paths{})

	require.NoError(t, store.Save(context.Background(), State{
		App:   AppConfig{ClientID: "client-id"},
		Token: TokenState{AccessToken: "oauth-access-token"},
	}))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got)
}

func Test_Store_reports_canceled_context(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	store := NewStore(testPaths(t))

	tests := []struct {
		name string
		run  func() error
		want string
	}{
		{
			name: "load",
			run: func() error {
				_, err := store.Load(ctx)
				return err
			},
			want: "load auth state context",
		},
		{
			name: "save",
			run: func() error {
				return store.Save(ctx, State{})
			},
			want: "save auth state context",
		},
		{
			name: "save app config",
			run: func() error {
				return store.SaveAppConfig(ctx, "", AppConfig{})
			},
			want: "save auth app config context",
		},
		{
			name: "save token state",
			run: func() error {
				return store.SaveTokenState(ctx, "", TokenState{})
			},
			want: "save auth token state context",
		},
		{
			name: "clear token state",
			run: func() error {
				return store.ClearTokenState(ctx, "")
			},
			want: "clear auth token state context",
		},
		{
			name: "clear app config",
			run: func() error {
				return store.ClearAppConfig(ctx, "")
			},
			want: "clear auth app config context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			require.ErrorIs(t, err, context.Canceled)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func Test_Store_reports_read_parse_and_write_errors(t *testing.T) {
	t.Parallel()

	t.Run("read app config directory", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		appConfigPath := filepath.Join(root, "auth-app-dir")
		require.NoError(t, os.Mkdir(appConfigPath, 0o700))
		store := NewStore(Paths{
			AppConfigPath: appConfigPath,
			TokenPath:     filepath.Join(root, tokenFileName),
		})

		_, err := store.Load(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "read auth app config")
		require.Contains(t, err.Error(), appConfigPath)
	})

	t.Run("read token state directory", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		tokenPath := filepath.Join(root, "auth-token-dir")
		require.NoError(t, os.Mkdir(tokenPath, 0o700))
		store := NewStore(Paths{
			AppConfigPath: filepath.Join(root, appConfigFileName),
			TokenPath:     tokenPath,
		})

		_, err := store.Load(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "read auth token state")
		require.Contains(t, err.Error(), tokenPath)
	})

	t.Run("parse app config", func(t *testing.T) {
		t.Parallel()
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o700))
		require.NoError(t, os.WriteFile(paths.AppConfigPath, []byte("{"), 0o600))
		store := NewStore(paths)

		_, err := store.Load(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "parse auth app config")
		require.Contains(t, err.Error(), paths.AppConfigPath)
	})

	t.Run("create app config directory", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		parentFile := filepath.Join(root, "not-a-directory")
		require.NoError(t, os.WriteFile(parentFile, []byte("file"), 0o600))
		store := NewStore(Paths{
			AppConfigPath: filepath.Join(parentFile, appConfigFileName),
			TokenPath:     filepath.Join(root, tokenFileName),
		})

		err := store.Save(context.Background(), State{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "create auth app config directory")
	})

	t.Run("write app config file", func(t *testing.T) {
		t.Parallel()
		root := t.TempDir()
		appConfigPath := filepath.Join(root, "auth-app-dir")
		require.NoError(t, os.Mkdir(appConfigPath, 0o700))
		store := NewStore(Paths{
			AppConfigPath: appConfigPath,
			TokenPath:     filepath.Join(root, tokenFileName),
		})

		err := store.Save(context.Background(), State{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write auth app config")
		require.Contains(t, err.Error(), appConfigPath)
	})
}

func Test_Store_reports_operation_read_errors(t *testing.T) {
	t.Parallel()

	t.Run("save app config parse error", func(t *testing.T) {
		t.Parallel()
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o700))
		require.NoError(t, os.WriteFile(paths.AppConfigPath, []byte("{"), 0o600))
		store := NewStore(paths)

		err := store.SaveAppConfig(context.Background(), "", AppConfig{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "parse auth app config")
	})

	t.Run("save token state parse error", func(t *testing.T) {
		t.Parallel()
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(filepath.Dir(paths.TokenPath), 0o700))
		require.NoError(t, os.WriteFile(paths.TokenPath, []byte("{"), 0o600))
		store := NewStore(paths)

		err := store.SaveTokenState(context.Background(), "", TokenState{})

		require.Error(t, err)
		require.Contains(t, err.Error(), "parse auth token state")
	})

	t.Run("clear token state parse error", func(t *testing.T) {
		t.Parallel()
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(filepath.Dir(paths.TokenPath), 0o700))
		require.NoError(t, os.WriteFile(paths.TokenPath, []byte("{"), 0o600))
		store := NewStore(paths)

		err := store.ClearTokenState(context.Background(), "")

		require.Error(t, err)
		require.Contains(t, err.Error(), "parse auth token state")
	})

	t.Run("clear app config parse error", func(t *testing.T) {
		t.Parallel()
		paths := testPaths(t)
		require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o700))
		require.NoError(t, os.WriteFile(paths.AppConfigPath, []byte("{"), 0o600))
		store := NewStore(paths)

		err := store.ClearAppConfig(context.Background(), "")

		require.Error(t, err)
		require.Contains(t, err.Error(), "parse auth app config")
	})
}

func Test_writeJSON_reports_encode_and_permission_errors(t *testing.T) {
	t.Run("encode error", func(t *testing.T) {
		err := writeJSON(filepath.Join(t.TempDir(), "auth.json"), func() {}, "auth app config")

		require.Error(t, err)
		require.Contains(t, err.Error(), "encode auth app config")
	})

	t.Run("secure directory error", func(t *testing.T) {
		withRuntimeGOOS(t, "linux")
		withChmodFile(t, func(string, os.FileMode) error {
			return errors.New("chmod failed")
		})

		err := writeJSON(filepath.Join(t.TempDir(), "dir", "auth.json"), struct{}{}, "auth app config")

		require.Error(t, err)
		require.Contains(t, err.Error(), "secure auth app config directory")
	})

	t.Run("secure file error", func(t *testing.T) {
		withRuntimeGOOS(t, "linux")
		calls := 0
		withChmodFile(t, func(string, os.FileMode) error {
			calls++
			if calls == 2 {
				return errors.New("chmod failed")
			}

			return nil
		})

		err := writeJSON(filepath.Join(t.TempDir(), "dir", "auth.json"), struct{}{}, "auth app config")

		require.Error(t, err)
		require.Contains(t, err.Error(), "secure auth app config")
		require.Equal(t, 2, calls)
	})
}

func Test_State_profile_selects_default_named_and_missing_state(t *testing.T) {
	t.Parallel()
	state := State{
		App:   AppConfig{ClientID: "client-id"},
		Token: TokenState{AccessToken: "oauth-access-token"},
		Profiles: map[string]ProfileState{
			"work": {
				App:   AppConfig{ClientID: "work-client-id"},
				Token: TokenState{AccessToken: "work-oauth-access-token"},
			},
		},
	}

	require.Equal(t, ProfileState{
		App:   AppConfig{ClientID: "client-id"},
		Token: TokenState{AccessToken: "oauth-access-token"},
	}, state.Profile(""))
	require.Equal(t, ProfileState{
		App:   AppConfig{ClientID: "work-client-id"},
		Token: TokenState{AccessToken: "work-oauth-access-token"},
	}, state.Profile("work"))
	require.Empty(t, state.Profile("missing"))
}

func Test_DefaultPaths_use_os_config_and_state_locations(t *testing.T) {
	root := t.TempDir()
	env := staticEnv{}
	var wantAppConfigPath string
	var wantTokenPath string
	switch runtime.GOOS {
	case "darwin":
		env["HOME"] = root
		base := filepath.Join(root, "Library", "Application Support", "linctl")
		wantAppConfigPath = filepath.Join(base, appConfigFileName)
		wantTokenPath = filepath.Join(base, tokenFileName)
	case "windows":
		env["APPDATA"] = root
		env["LOCALAPPDATA"] = filepath.Join(root, "local")
		wantAppConfigPath = filepath.Join(root, "linctl", appConfigFileName)
		wantTokenPath = filepath.Join(root, "local", "linctl", tokenFileName)
	default:
		env["XDG_CONFIG_HOME"] = filepath.Join(root, "config")
		env["XDG_STATE_HOME"] = filepath.Join(root, "state")
		wantAppConfigPath = filepath.Join(root, "config", "linctl", appConfigFileName)
		wantTokenPath = filepath.Join(root, "state", "linctl", tokenFileName)
	}

	got, err := DefaultPaths(env)

	require.NoError(t, err)
	require.Equal(t, wantAppConfigPath, got.AppConfigPath)
	require.Equal(t, wantTokenPath, got.TokenPath)
	require.NotContains(t, got.AppConfigPath, ".linctl.toml")
	require.NotContains(t, got.TokenPath, ".linctl.toml")
}

func Test_DefaultPaths_cover_supported_os_locations(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name              string
		goos              string
		env               staticEnv
		wantAppConfigPath string
		wantTokenPath     string
	}{
		{
			name: "windows local app data",
			goos: "windows",
			env: staticEnv{
				"APPDATA":      filepath.Join(root, "appdata"),
				"LOCALAPPDATA": filepath.Join(root, "localappdata"),
			},
			wantAppConfigPath: filepath.Join(root, "appdata", "linctl", appConfigFileName),
			wantTokenPath:     filepath.Join(root, "localappdata", "linctl", tokenFileName),
		},
		{
			name: "windows app data fallback",
			goos: "windows",
			env: staticEnv{
				"APPDATA": filepath.Join(root, "appdata"),
			},
			wantAppConfigPath: filepath.Join(root, "appdata", "linctl", appConfigFileName),
			wantTokenPath:     filepath.Join(root, "appdata", "linctl", tokenFileName),
		},
		{
			name: "darwin application support",
			goos: "darwin",
			env:  staticEnv{"HOME": root},
			wantAppConfigPath: filepath.Join(
				root,
				"Library",
				"Application Support",
				"linctl",
				appConfigFileName,
			),
			wantTokenPath: filepath.Join(root, "Library", "Application Support", "linctl", tokenFileName),
		},
		{
			name: "linux xdg",
			goos: "linux",
			env: staticEnv{
				"XDG_CONFIG_HOME": filepath.Join(root, "config"),
				"XDG_STATE_HOME":  filepath.Join(root, "state"),
			},
			wantAppConfigPath: filepath.Join(root, "config", "linctl", appConfigFileName),
			wantTokenPath:     filepath.Join(root, "state", "linctl", tokenFileName),
		},
		{
			name:              "linux home fallback",
			goos:              "linux",
			env:               staticEnv{"HOME": root},
			wantAppConfigPath: filepath.Join(root, ".config", "linctl", appConfigFileName),
			wantTokenPath:     filepath.Join(root, ".local", "state", "linctl", tokenFileName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRuntimeGOOS(t, tt.goos)

			got, err := DefaultPaths(tt.env)

			require.NoError(t, err)
			require.Equal(t, tt.wantAppConfigPath, got.AppConfigPath)
			require.Equal(t, tt.wantTokenPath, got.TokenPath)
		})
	}
}

func Test_DefaultPaths_reports_supported_os_environment_errors(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name string
		goos string
		env  staticEnv
		want string
	}{
		{
			name: "windows missing app data",
			goos: "windows",
			env:  staticEnv{},
			want: "APPDATA is not set",
		},
		{
			name: "darwin missing home",
			goos: "darwin",
			env:  staticEnv{},
			want: "HOME is not set",
		},
		{
			name: "linux relative config home",
			goos: "linux",
			env: staticEnv{
				"HOME":            root,
				"XDG_CONFIG_HOME": "relative",
			},
			want: "XDG_CONFIG_HOME must be absolute",
		},
		{
			name: "linux relative state home",
			goos: "linux",
			env: staticEnv{
				"HOME":            root,
				"XDG_CONFIG_HOME": filepath.Join(root, "config"),
				"XDG_STATE_HOME":  "relative",
			},
			want: "XDG_STATE_HOME must be absolute",
		},
		{
			name: "linux missing state home fallback",
			goos: "linux",
			env: staticEnv{
				"XDG_CONFIG_HOME": filepath.Join(root, "config"),
			},
			want: "HOME is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRuntimeGOOS(t, tt.goos)

			_, err := DefaultPaths(tt.env)

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func Test_chmodIfSupported_skips_windows_permissions(t *testing.T) {
	withRuntimeGOOS(t, "windows")
	withChmodFile(t, func(string, os.FileMode) error {
		return errors.New("chmod should not be called")
	})

	require.NoError(t, chmodIfSupported("ignored", 0o600))
}

func Test_DefaultPaths_uses_xdg_state_fallback_on_linux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("XDG state fallback is a Linux path contract")
	}
	root := t.TempDir()

	got, err := DefaultPaths(staticEnv{
		"HOME":            root,
		"XDG_CONFIG_HOME": filepath.Join(root, "config"),
	})

	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, ".local", "state", "linctl", tokenFileName), got.TokenPath)
}

func Test_DefaultPaths_uses_process_environment_when_env_is_nil(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("nil env fallback test uses XDG paths")
	}
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(root, "state"))

	got, err := DefaultPaths(nil)

	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, "config", "linctl", appConfigFileName), got.AppConfigPath)
	require.Equal(t, filepath.Join(root, "state", "linctl", tokenFileName), got.TokenPath)
}

func Test_DefaultPaths_reports_environment_errors(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name string
		env  staticEnv
		want string
	}{
		{
			name: "relative XDG_CONFIG_HOME",
			env: staticEnv{
				"HOME":            root,
				"XDG_CONFIG_HOME": "relative",
				"XDG_STATE_HOME":  filepath.Join(root, "state"),
			},
			want: "XDG_CONFIG_HOME must be absolute",
		},
		{
			name: "relative XDG_STATE_HOME",
			env: staticEnv{
				"HOME":            root,
				"XDG_CONFIG_HOME": filepath.Join(root, "config"),
				"XDG_STATE_HOME":  "relative",
			},
			want: "XDG_STATE_HOME must be absolute",
		},
		{
			name: "missing HOME",
			env:  staticEnv{},
			want: "HOME is not set",
		},
	}
	if runtime.GOOS == "windows" {
		tests = []struct {
			name string
			env  staticEnv
			want string
		}{
			{
				name: "missing APPDATA",
				env:  staticEnv{},
				want: "APPDATA is not set",
			},
		}
	}
	if runtime.GOOS == "darwin" {
		tests = []struct {
			name string
			env  staticEnv
			want string
		}{
			{
				name: "missing HOME",
				env:  staticEnv{},
				want: "HOME is not set",
			},
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DefaultPaths(tt.env)
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
		})
	}
}

func Test_osEnv_lookup_reads_process_environment(t *testing.T) {
	const key = "LINCTL_AUTH_TEST_LOOKUP"
	t.Setenv(key, "set")

	value, ok := osEnv{}.Lookup(key)
	require.True(t, ok)
	require.Equal(t, "set", value)

	value, ok = osEnv{}.Lookup(key + "_MISSING")
	require.False(t, ok)
	require.Empty(t, value)
}

func testPaths(t *testing.T) Paths {
	t.Helper()
	root := t.TempDir()

	return Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", appConfigFileName),
		TokenPath:     filepath.Join(root, "state", "linctl", tokenFileName),
	}
}

func assertPrivateDirAndFile(t *testing.T, dir string, path string) {
	t.Helper()
	dirInfo, err := os.Stat(dir)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o700), dirInfo.Mode().Perm())
	fileInfo, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), fileInfo.Mode().Perm())
}

func withRuntimeGOOS(t *testing.T, goos string) {
	t.Helper()
	original := runtimeGOOS
	runtimeGOOS = goos
	t.Cleanup(func() {
		runtimeGOOS = original
	})
}

func withChmodFile(t *testing.T, fn func(string, os.FileMode) error) {
	t.Helper()
	original := chmodFile
	chmodFile = fn
	t.Cleanup(func() {
		chmodFile = original
	})
}
