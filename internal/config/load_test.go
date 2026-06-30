package config

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/stretchr/testify/require"
)

type staticEnv map[string]string

func (env staticEnv) Lookup(key string) (string, bool) {
	value, ok := env[key]
	return value, ok
}

func Test_Load_resolves_repo_profile_and_oauth_env_token_when_present(t *testing.T) {
	// Given
	root := t.TempDir()
	globalPath := filepath.Join(root, "global.toml")
	repoPath := filepath.Join(root, "repo.toml")
	require.NoError(t, os.WriteFile(globalPath, []byte(`
profile = "global"

[profiles.global]
token = "global-token"

[profiles.global.target]
org_id = "global-org"
team_key = "GLOBAL"
team_id = "global-team"
`), 0o600))
	require.NoError(t, os.WriteFile(repoPath, []byte(`
profile = "repo"

[profiles.repo]
token = "repo-token"

[profiles.repo.target]
org_id = "repo-org"
team_key = "REPO"
team_id = "repo-team"
project_id = "repo-project"
`), 0o600))

	// When
	resolved, err := Load(context.Background(), LoadRequest{
		GlobalPath: globalPath,
		RepoPath:   repoPath,
		Env: staticEnv{
			"LINCTL_OAUTH_ACCESS_TOKEN": "env-token",
		},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "repo", resolved.Profile)
	require.Equal(t, "env-token", resolved.Token)
	require.Equal(t, "env-token", resolved.Auth.AccessToken)
	require.Equal(t, Target{
		OrgID:     "repo-org",
		TeamKey:   "REPO",
		TeamID:    "repo-team",
		ProjectID: "repo-project",
	}, resolved.Target)
}

func Test_Load_applies_explicit_profile_and_target_overrides_when_present(t *testing.T) {
	// Given
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
profile = "default"

[profiles.default]
token = "default-token"

[profiles.default.target]
org_id = "default-org"
team_key = "DEF"
team_id = "default-team"

[profiles.other]
token = "other-token"

[profiles.other.target]
org_id = "other-org"
team_key = "OTH"
team_id = "other-team"
`), 0o600))

	// When
	resolved, err := Load(context.Background(), LoadRequest{
		GlobalPath:      configPath,
		ProfileOverride: "other",
		TargetOverride: Target{
			ProjectID: "override-project",
		},
		Env: staticEnv{},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "other", resolved.Profile)
	require.Empty(t, resolved.Token)
	require.Empty(t, resolved.Auth.AccessToken)
	require.Equal(t, Target{
		OrgID:     "other-org",
		TeamKey:   "OTH",
		TeamID:    "other-team",
		ProjectID: "override-project",
	}, resolved.Target)
}

func Test_Load_keeps_profile_targets_separate_when_multiple_targets_exist(t *testing.T) {
	// Given
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[profiles.personal]
token = "personal-token"

[profiles.personal.target]
org_id = "personal-org"
team_key = "PER"
team_id = "personal-team"

[profiles.work]
token = "work-token"

[profiles.work.target]
org_id = "work-org"
team_key = "WRK"
team_id = "work-team"
project_id = "work-project"
`), 0o600))

	// When
	resolved, err := Load(context.Background(), LoadRequest{
		GlobalPath:      configPath,
		ProfileOverride: "work",
		Env: staticEnv{
			"LINCTL_OAUTH_ACCESS_TOKEN": "env-token",
		},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "work", resolved.Profile)
	require.Equal(t, "env-token", resolved.Token)
	require.Equal(t, "env-token", resolved.Auth.AccessToken)
	require.Equal(t, Target{
		OrgID:     "work-org",
		TeamKey:   "WRK",
		TeamID:    "work-team",
		ProjectID: "work-project",
	}, resolved.Target)
}

func Test_Load_refuses_unknown_explicit_profile(t *testing.T) {
	// Given
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[profiles.default]
token = "default-token"
`), 0o600))

	// When
	_, err := Load(context.Background(), LoadRequest{
		GlobalPath:      configPath,
		ProfileOverride: "missing",
		Env:             staticEnv{},
	})

	// Then
	require.ErrorIs(t, err, ErrProfileNotFound)
}

func Test_Load_ignores_legacy_config_token_without_permission_gate(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`token = "file-token"`), 0o644))

	resolved, err := Load(context.Background(), LoadRequest{
		Env:      staticEnv{},
		RepoPath: configPath,
	})

	require.NoError(t, err)
	require.Empty(t, resolved.Token)
	require.Empty(t, resolved.Auth.AccessToken)
}

func Test_Load_ignores_legacy_api_key_sources_for_product_auth(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
token = "file-token"

[profiles.default]
token = "profile-token"
`), 0o600))

	resolved, err := Load(context.Background(), LoadRequest{
		Env: staticEnv{
			"LINCTL_TOKEN":   "linctl-token",
			"LINEAR_API_KEY": "linear-token",
		},
		RepoPath:        configPath,
		ProfileOverride: "default",
	})

	require.NoError(t, err)
	require.Empty(t, resolved.Token)
	require.Empty(t, resolved.Auth.AccessToken)
}

func Test_Load_rejects_personal_api_key_shaped_oauth_env_token(t *testing.T) {
	_, err := Load(context.Background(), LoadRequest{
		Env: staticEnv{
			"LINCTL_OAUTH_ACCESS_TOKEN": "lin_api_personal_key",
		},
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "personal API key")
}

func Test_Load_accepts_oauth_app_material_from_env(t *testing.T) {
	resolved, err := Load(context.Background(), LoadRequest{
		Env: staticEnv{
			"LINCTL_OAUTH_CLIENT_ID":     "client-id",
			"LINCTL_OAUTH_CLIENT_SECRET": "client-secret",
			"LINCTL_TOKEN":               "linctl-token",
			"LINEAR_API_KEY":             "linear-token",
		},
	})

	require.NoError(t, err)
	require.Empty(t, resolved.Auth.AccessToken)
	require.Equal(t, "client-id", resolved.Auth.App.ClientID)
	require.Equal(t, "client-secret", resolved.Auth.App.ClientSecret)
}

func Test_Load_env_oauth_app_material_overrides_local_auth_app(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`[profiles.work]`), 0o600))
	authPaths := auth.Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", "auth-app.json"),
		TokenPath:     filepath.Join(root, "state", "linctl", "auth-token.json"),
	}
	require.NoError(t, auth.NewStore(authPaths).Save(context.Background(), auth.State{
		Profiles: map[string]auth.ProfileState{
			"work": {
				App: auth.AppConfig{
					ClientID:     "local-client-id",
					ClientSecret: "local-client-secret",
					RedirectURI:  "http://127.0.0.1:8484/local",
					Scopes:       []string{"read"},
				},
				Token: auth.TokenState{AccessToken: "local-oauth-token"},
			},
		},
	}))

	resolved, err := Load(context.Background(), LoadRequest{
		Env: staticEnv{
			"LINCTL_OAUTH_CLIENT_ID":     "env-client-id",
			"LINCTL_OAUTH_CLIENT_SECRET": "env-client-secret",
			"LINCTL_OAUTH_REDIRECT_URI":  "http://127.0.0.1:8484/env",
			"LINCTL_OAUTH_SCOPES":        "read, write\ncomments:create",
		},
		GlobalPath:      configPath,
		AuthStatePaths:  authPaths,
		ProfileOverride: "work",
	})

	require.NoError(t, err)
	require.Equal(t, "local-oauth-token", resolved.Auth.AccessToken)
	require.Equal(t, auth.AppConfig{
		ClientID:     "env-client-id",
		ClientSecret: "env-client-secret",
		RedirectURI:  "http://127.0.0.1:8484/env",
		Scopes:       []string{"read", "write", "comments:create"},
	}, resolved.Auth.App)
}

func Test_Load_resolves_oauth_token_from_local_auth_state(t *testing.T) {
	root := t.TempDir()
	authPaths := auth.Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", "auth-app.json"),
		TokenPath:     filepath.Join(root, "state", "linctl", "auth-token.json"),
	}
	require.NoError(t, auth.NewStore(authPaths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Token: auth.TokenState{
			AccessToken:  "local-oauth-token",
			RefreshToken: "refresh-token",
		},
	}))

	resolved, err := Load(context.Background(), LoadRequest{
		Env:            staticEnv{},
		AuthStatePaths: authPaths,
	})

	require.NoError(t, err)
	require.Equal(t, "local-oauth-token", resolved.Token)
	require.Equal(t, "local-oauth-token", resolved.Auth.AccessToken)
	require.Equal(t, "client-id", resolved.Auth.App.ClientID)
}

func Test_Load_rejects_personal_api_key_shaped_local_auth_state_token(t *testing.T) {
	root := t.TempDir()
	authPaths := auth.Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", "auth-app.json"),
		TokenPath:     filepath.Join(root, "state", "linctl", "auth-token.json"),
	}
	require.NoError(t, auth.NewStore(authPaths).Save(context.Background(), auth.State{
		Token: auth.TokenState{
			AccessToken: "lin_api_personal_key",
		},
	}))

	_, err := Load(context.Background(), LoadRequest{
		Env:            staticEnv{},
		AuthStatePaths: authPaths,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "personal API key")
}

func Test_Load_reports_auth_state_read_error(t *testing.T) {
	root := t.TempDir()
	appConfigPath := filepath.Join(root, "auth-app-dir")
	require.NoError(t, os.Mkdir(appConfigPath, 0o700))

	_, err := Load(context.Background(), LoadRequest{
		Env: staticEnv{},
		AuthStatePaths: auth.Paths{
			AppConfigPath: appConfigPath,
			TokenPath:     filepath.Join(root, "auth-token.json"),
		},
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read auth app config")
}

func Test_OAuthTokenSource_resolves_access_token(t *testing.T) {
	token, err := (OAuthTokenSource{AccessToken: "test-token"}).Resolve(context.Background())

	require.NoError(t, err)
	require.Equal(t, "test-token", token)
}

func Test_OAuthTokenSource_respects_context_cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := OAuthTokenSource{AccessToken: "test-token"}.Resolve(ctx)

	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve oauth token source")
}

func Test_Load_allows_broad_config_permissions_without_tokens(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows file modes do not expose unix group/world permission bits")
	}
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[target]
team_key = "ENG"
`), 0o644))

	resolved, err := Load(context.Background(), LoadRequest{
		Env:      staticEnv{},
		RepoPath: configPath,
	})

	require.NoError(t, err)
	require.Equal(t, "ENG", resolved.Target.TeamKey)
}

func Test_Load_reports_read_error_after_config_stat_succeeds(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config-dir")
	require.NoError(t, os.Mkdir(configPath, 0o700))

	_, err := Load(context.Background(), LoadRequest{
		Env:      staticEnv{},
		RepoPath: configPath,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read config")
}

func Test_Load_reports_config_stat_error(t *testing.T) {
	_, err := Load(context.Background(), LoadRequest{
		Env:      staticEnv{},
		RepoPath: "bad\x00path",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read config")
}
