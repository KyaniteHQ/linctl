package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

func Test_AuthConfigure_saves_oauth_app_config_without_token_state(t *testing.T) {
	paths := cliAuthTestPaths(t)
	restore := useAuthPaths(t, paths)
	defer restore()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, []string{
		"--json",
		"auth",
		"configure",
		"--client-id", "client-id",
		"--client-secret", "client-secret",
		"--redirect-uri", "http://127.0.0.1:8484/callback",
		"--scopes", "read,write,issues:create,comments:create",
	})

	require.NoError(t, err)
	require.NotContains(t, stdout.String(), "client-secret")
	require.NotContains(t, stderr.String(), "client-secret")
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		Scopes:       []string{"read", "write", "issues:create", "comments:create"},
	}, got.App)
	require.Empty(t, got.Token)
	_, err = os.Stat(paths.TokenPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func Test_AuthApp_obtains_app_actor_token_after_live_readiness(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read", "write"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"oauth-access-token",
		"",
		"Bearer",
		time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC),
		[]string{"read", "write"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("app")}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, fakeReadiness)
	defer restore()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"--project", "project-id",
		"auth",
		"app",
	})

	require.NoError(t, err)
	require.Equal(t, "oauth-access-token", fakeReadiness.accessToken)
	require.Equal(t, "app", fakeReadiness.expectedActor)
	require.Equal(t, []string{"read", "write"}, fakeReadiness.requiredScopes)
	require.NotContains(t, stdout.String(), "oauth-access-token")
	require.NotContains(t, stdout.String(), "client-secret")
	require.NotContains(t, stderr.String(), "oauth-access-token")
	require.NotContains(t, stderr.String(), "client-secret")
	var report authStatusReport
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &report))
	require.Equal(t, "app", report.Actor)
	require.Equal(t, "ready", report.Target.Status)

	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "oauth-access-token", got.Token.AccessToken)
	require.Equal(t, "app", got.Token.Actor)
	require.Equal(t, "client_credentials", got.Token.GrantType)
}

func Test_AuthApp_does_not_save_token_when_readiness_fails(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read", "write"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"oauth-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read", "write"},
	)}
	restore := useAuthCommandHooks(
		t,
		paths,
		fakeOAuth,
		&fakeAuthReadinessChecker{err: client.ErrTargetMismatch},
	)
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"app",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), `"error_code":"AUTH_TARGET_MISMATCH"`)
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got.Token)
}

func Test_AuthApp_reports_missing_scope_without_saving_token(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read", "write"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"oauth-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"app",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), `"error_code":"MISSING_SCOPE"`)
	require.Contains(t, stderr.String(), "linctl auth configure --scopes read,write")
	require.Contains(t, stderr.String(), "linctl auth app or linctl auth login")
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got.Token)
}

func Test_AuthStatus_reacquires_expired_client_credentials_token_and_checks_readiness(t *testing.T) {
	paths := cliAuthTestPaths(t)
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Token: auth.TokenState{
			AccessToken: "expired-access-token",
			TokenType:   "Bearer",
			Scopes:      []string{"read"},
			ExpiresAt:   &expiredAt,
			Actor:       "app",
			GrantType:   "client_credentials",
		},
	}))
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-access-token",
		"",
		"Bearer",
		expiresAt,
		[]string{"read"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("app")}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, fakeReadiness)
	defer restore()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"status",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, "fresh-access-token", fakeReadiness.accessToken)
	require.NotContains(t, stdout.String(), "fresh-access-token")
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "fresh-access-token", got.Token.AccessToken)
}

func Test_AuthStatus_refreshes_expired_authorization_code_token_and_checks_readiness(t *testing.T) {
	paths := cliAuthTestPaths(t)
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	freshExpiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID: "client-id",
			Scopes:   []string{"read", "write"},
		},
		Token: auth.TokenState{
			AccessToken:  "expired-access-token",
			RefreshToken: "old-refresh-token",
			TokenType:    "Bearer",
			Scopes:       []string{"read", "write"},
			ExpiresAt:    &expiredAt,
			Actor:        "app",
			GrantType:    authGrantAuthorizationCode,
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-access-token",
		"fresh-refresh-token",
		"Bearer",
		freshExpiresAt,
		[]string{"read", "write"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("app")}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, fakeReadiness)
	defer restore()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"status",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.refreshTokenCalls)
	require.Equal(t, "old-refresh-token", fakeOAuth.refreshTokenRequest.RefreshToken)
	require.Equal(t, "fresh-access-token", fakeReadiness.accessToken)
	require.NotContains(t, stdout.String(), "fresh-access-token")
	require.NotContains(t, stdout.String(), "fresh-refresh-token")
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "fresh-access-token", got.Token.AccessToken)
	require.Equal(t, "fresh-refresh-token", got.Token.RefreshToken)
	require.Equal(t, authGrantAuthorizationCode, got.Token.GrantType)
}

func Test_AuthRefresh_rotates_authorization_code_token_and_checks_readiness(t *testing.T) {
	paths := cliAuthTestPaths(t)
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read", "write"},
		},
		Token: auth.TokenState{
			AccessToken:  "old-access-token",
			RefreshToken: "old-refresh-token",
			TokenType:    "Bearer",
			Scopes:       []string{"read", "write"},
			ExpiresAt:    &expiresAt,
			Actor:        "user",
			GrantType:    authGrantAuthorizationCode,
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"new-access-token",
		"new-refresh-token",
		"Bearer",
		expiresAt,
		[]string{"read", "write"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("user")}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, fakeReadiness)
	defer restore()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"refresh",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.refreshTokenCalls)
	require.Equal(t, "old-refresh-token", fakeOAuth.refreshTokenRequest.RefreshToken)
	require.Equal(t, "client-id", fakeOAuth.refreshTokenRequest.ClientID)
	require.Equal(t, "client-secret", fakeOAuth.refreshTokenRequest.ClientSecret)
	require.Equal(t, "new-access-token", fakeReadiness.accessToken)
	require.Equal(t, "user", fakeReadiness.expectedActor)
	require.NotContains(t, stdout.String(), "new-access-token")
	require.NotContains(t, stdout.String(), "new-refresh-token")
	require.NotContains(t, stdout.String(), "client-secret")
	var report authStatusReport
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &report))
	require.Equal(t, "user", report.Actor)

	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "new-access-token", got.Token.AccessToken)
	require.Equal(t, "new-refresh-token", got.Token.RefreshToken)
	require.Equal(t, "user", got.Token.Actor)
	require.Equal(t, authGrantAuthorizationCode, got.Token.GrantType)
}

func Test_AuthRefresh_reacquires_client_credentials_token(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Token: auth.TokenState{
			AccessToken: "old-app-token",
			TokenType:   "Bearer",
			Scopes:      []string{"read"},
			Actor:       "app",
			GrantType:   authGrantClientCredentials,
		},
	}))
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"new-app-token",
		"",
		"Bearer",
		expiresAt,
		[]string{"read"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("app")}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, fakeReadiness)
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"refresh",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, "new-app-token", fakeReadiness.accessToken)
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "new-app-token", got.Token.AccessToken)
	require.Equal(t, authGrantClientCredentials, got.Token.GrantType)
}

func Test_AuthLogout_revokes_tokens_and_clears_token_state_while_keeping_app(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Token: auth.TokenState{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			TokenType:    "Bearer",
			Scopes:       []string{"read"},
			Actor:        "user",
			GrantType:    authGrantAuthorizationCode,
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{})
	defer restore()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"auth",
		"logout",
	})

	require.NoError(t, err)
	require.Equal(t, []oauth.RevocationRequest{
		{Token: "refresh-token", TokenTypeHint: "refresh_token"},
		{Token: "access-token", TokenTypeHint: "access_token"},
	}, fakeOAuth.revokeTokenRequests)
	require.NotContains(t, stdout.String(), "access-token")
	require.NotContains(t, stdout.String(), "refresh-token")
	require.NotContains(t, stdout.String(), "client-secret")
	var report authLogoutReport
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &report))
	require.Equal(t, "removed", report.Token)
	require.Equal(t, "kept", report.App)
	require.Equal(t, []string{"refresh_token", "access_token"}, report.Revoked)
	require.False(t, report.RevocationFailed)

	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "client-id", got.App.ClientID)
	require.Empty(t, got.Token)
}

func Test_AuthLogout_forgets_app_and_clears_state_when_revocation_fails(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Token: auth.TokenState{
			AccessToken: "access-token",
			TokenType:   "Bearer",
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{revokeTokenErr: errors.New("revoked already")}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{})
	defer restore()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"auth",
		"logout",
		"--forget-app",
	})

	require.NoError(t, err)
	var report authLogoutReport
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &report))
	require.Equal(t, "removed", report.Token)
	require.Equal(t, "forgotten", report.App)
	require.True(t, report.RevocationFailed)
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got.App)
	require.Empty(t, got.Token)
}

func Test_DefaultCheckAuthReadiness_requires_access_token(t *testing.T) {
	_, err := defaultCheckAuthReadiness(context.Background(), authReadinessRequest{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeNotConfigured), errorCode(err))
}

func Test_AuthCommands_report_default_path_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "configure",
			args: []string{"auth", "configure", "--client-id", "client-id"},
		},
		{
			name: "app",
			args: []string{"auth", "app"},
		},
		{
			name: "status",
			args: []string{"auth", "status"},
		},
		{
			name: "login",
			args: []string{"auth", "login"},
		},
		{
			name: "refresh",
			args: []string{"auth", "refresh"},
		},
		{
			name: "logout",
			args: []string{"auth", "logout"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := useAuthPathsError(t, errors.New("paths unavailable"))
			defer restore()

			err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, tt.args)

			require.Error(t, err)
			require.Contains(t, err.Error(), "paths unavailable")
		})
	}
}

func Test_AuthConfigure_reports_missing_client_id(t *testing.T) {
	paths := cliAuthTestPaths(t)
	restore := useAuthPaths(t, paths)
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"auth",
		"configure",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), `"error_code":"AUTH_NOT_CONFIGURED"`)
	require.Contains(t, stderr.String(), "missing --client-id")
}

func Test_AuthConfigure_quiet_and_human_output(t *testing.T) {
	t.Run("quiet", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		restore := useAuthPaths(t, paths)
		defer restore()
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, []string{
			"--quiet",
			"auth",
			"configure",
			"--client-id", "client-id",
		})

		require.NoError(t, err)
		require.Empty(t, stdout.String())
		require.Empty(t, stderr.String())
	})

	t.Run("human", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		restore := useAuthPaths(t, paths)
		defer restore()
		var stdout bytes.Buffer

		err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
			"auth",
			"configure",
			"--client-id", "client-id",
		})

		require.NoError(t, err)
		require.Contains(t, stdout.String(), "OAuth app configured")
	})
}

func Test_AuthConfigure_reports_save_error(t *testing.T) {
	root := t.TempDir()
	appConfigPath := filepath.Join(root, "auth-app-dir")
	require.NoError(t, os.Mkdir(appConfigPath, 0o700))
	paths := auth.Paths{
		AppConfigPath: appConfigPath,
		TokenPath:     filepath.Join(root, "auth-token.json"),
	}
	restore := useAuthPaths(t, paths)
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"configure",
		"--client-id", "client-id",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "auth app config")
}

func Test_AuthCommandContext_reports_config_and_state_load_errors(t *testing.T) {
	t.Run("config", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		paths := cliAuthTestPaths(t)
		restore := useAuthPaths(t, paths)
		defer restore()
		require.NoError(t, os.WriteFile(".linctl.toml", []byte("["), 0o600))

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
			"auth",
			"status",
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "parse config")
	})

	t.Run("state", func(t *testing.T) {
		root := t.TempDir()
		appConfigPath := filepath.Join(root, "auth-app-dir")
		require.NoError(t, os.Mkdir(appConfigPath, 0o700))
		paths := auth.Paths{
			AppConfigPath: appConfigPath,
			TokenPath:     filepath.Join(root, "auth-token.json"),
		}
		restore := useAuthPaths(t, paths)
		defer restore()

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
			"auth",
			"status",
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "read auth app config")
	})
}

func Test_AuthApp_reports_missing_client_configuration(t *testing.T) {
	t.Run("client id", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{})
		defer restore()
		var stderr bytes.Buffer

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
			"auth",
			"app",
		})

		require.Error(t, err)
		require.Contains(t, stderr.String(), "missing OAuth client id")
	})

	t.Run("client secret", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
			ClientID: "client-id",
		}))
		restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{})
		defer restore()
		var stderr bytes.Buffer

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
			"auth",
			"app",
		})

		require.Error(t, err)
		require.Contains(t, stderr.String(), "missing OAuth client secret")
	})
}

func Test_AuthApp_quiet_and_human_output(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{
			name: "quiet",
			args: []string{"--quiet", "auth", "app"},
		},
		{
			name:       "human",
			args:       []string{"auth", "app"},
			wantOutput: "auth set actor app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := cliAuthTestPaths(t)
			require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
				ClientID:     "client-id",
				ClientSecret: "client-secret",
				Scopes:       []string{"read"},
			}))
			fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
				"oauth-access-token",
				"",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{"read"},
			)}
			restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
			defer restore()
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, tt.args)

			require.NoError(t, err)
			require.NotContains(t, stdout.String(), "oauth-access-token")
			require.NotContains(t, stderr.String(), "oauth-access-token")
			if tt.wantOutput == "" {
				require.Empty(t, stdout.String())
			} else {
				require.Contains(t, stdout.String(), tt.wantOutput)
			}
		})
	}
}

func Test_AuthApp_reports_token_state_save_error(t *testing.T) {
	root := t.TempDir()
	paths := auth.Paths{
		AppConfigPath: filepath.Join(root, "auth-app.json"),
		TokenPath:     filepath.Join(root, "auth-token.json"),
	}
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"oauth-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{
		report: readyAuthReport("app"),
		beforeReturn: func() {
			require.NoError(t, os.Mkdir(paths.TokenPath, 0o700))
		},
	})
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"app",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "auth token state")
}

func Test_AuthStatus_checks_current_token_readiness(t *testing.T) {
	paths := cliAuthTestPaths(t)
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}},
		Token: auth.TokenState{
			AccessToken: "current-access-token",
			TokenType:   "Bearer",
			Scopes:      []string{"read"},
			ExpiresAt:   &expiresAt,
			Actor:       "app",
		},
	}))
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("app")}
	restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, fakeReadiness)
	defer restore()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"auth",
		"status",
	})

	require.NoError(t, err)
	require.Equal(t, "current-access-token", fakeReadiness.accessToken)
	require.Equal(t, "app", fakeReadiness.expectedActor)
	require.Contains(t, stdout.String(), "auth set actor app")
	require.NotContains(t, stdout.String(), "current-access-token")
}

func Test_AuthStatus_reports_current_token_readiness_error(t *testing.T) {
	paths := cliAuthTestPaths(t)
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}},
		Token: auth.TokenState{
			AccessToken: "current-access-token",
			TokenType:   "Bearer",
			Scopes:      []string{"read"},
			ExpiresAt:   &expiresAt,
			Actor:       "app",
		},
	}))
	restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{
		err: auth.NewError(auth.ErrorCodeActorMismatch, "wrong actor"),
	})
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"auth",
		"status",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), `"error_code":"AUTH_ACTOR_MISMATCH"`)
}

func Test_AuthStatus_acquires_client_credentials_when_token_state_is_missing(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"status",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	got, loadErr := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "fresh-access-token", got.Token.AccessToken)
}

func Test_AuthStatus_reports_save_error_after_refresh_or_acquire(t *testing.T) {
	root := t.TempDir()
	paths := auth.Paths{
		AppConfigPath: filepath.Join(root, "auth-app.json"),
		TokenPath:     filepath.Join(root, "auth-token.json"),
	}
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{
		report: readyAuthReport("app"),
		beforeReturn: func() {
			require.NoError(t, os.Mkdir(paths.TokenPath, 0o700))
		},
	})
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"status",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "auth token state")
}

func Test_AuthStatus_reports_missing_app_for_implicit_token_acquire(t *testing.T) {
	paths := cliAuthTestPaths(t)
	restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{})
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"auth",
		"status",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), "run linctl auth configure and linctl auth app")
}

func Test_AuthRefresh_reports_token_state_configuration_errors(t *testing.T) {
	tests := []struct {
		name  string
		state auth.State
		want  string
	}{
		{
			name: "missing token state",
			want: "missing OAuth token state",
		},
		{
			name: "client credentials missing app",
			state: auth.State{
				Token: auth.TokenState{
					AccessToken: "access-token",
					GrantType:   authGrantClientCredentials,
				},
			},
			want: "missing OAuth app client credentials",
		},
		{
			name: "authorization code missing refresh token",
			state: auth.State{
				App: auth.AppConfig{ClientID: "client-id"},
				Token: auth.TokenState{
					AccessToken: "access-token",
					GrantType:   authGrantAuthorizationCode,
				},
			},
			want: "missing OAuth refresh token",
		},
		{
			name: "authorization code missing client id",
			state: auth.State{
				Token: auth.TokenState{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					GrantType:    authGrantAuthorizationCode,
				},
			},
			want: "missing OAuth client id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := cliAuthTestPaths(t)
			require.NoError(t, auth.NewStore(paths).Save(context.Background(), tt.state))
			restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{})
			defer restore()
			var stderr bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
				"auth",
				"refresh",
			})

			require.Error(t, err)
			require.Contains(t, stderr.String(), tt.want)
		})
	}
}

func Test_AuthRefresh_reports_save_error(t *testing.T) {
	root := t.TempDir()
	paths := auth.Paths{
		AppConfigPath: filepath.Join(root, "auth-app.json"),
		TokenPath:     filepath.Join(root, "auth-token.json"),
	}
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Token: auth.TokenState{
			AccessToken:  "old-access-token",
			RefreshToken: "old-refresh-token",
			Scopes:       []string{"read"},
			GrantType:    authGrantAuthorizationCode,
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"new-access-token",
		"new-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{
		report: readyAuthReport("app"),
		beforeReturn: func() {
			require.NoError(t, os.Remove(paths.TokenPath))
			require.NoError(t, os.Mkdir(paths.TokenPath, 0o700))
		},
	})
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"refresh",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "auth token state")
}

func Test_AuthLogout_quiet_and_human_output(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
	}{
		{
			name: "quiet",
			args: []string{"--quiet", "auth", "logout"},
		},
		{
			name:       "human",
			args:       []string{"auth", "logout"},
			wantOutput: "auth logout token removed app kept",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := cliAuthTestPaths(t)
			require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
				App:   auth.AppConfig{ClientID: "client-id"},
				Token: auth.TokenState{AccessToken: "access-token"},
			}))
			restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{})
			defer restore()
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, tt.args)

			require.NoError(t, err)
			require.NotContains(t, stdout.String(), "access-token")
			require.NotContains(t, stderr.String(), "access-token")
			if tt.wantOutput == "" {
				require.Empty(t, stdout.String())
			} else {
				require.Contains(t, stdout.String(), tt.wantOutput)
			}
		})
	}
}

func Test_AuthLogout_reports_clear_state_errors(t *testing.T) {
	t.Run("token state", func(t *testing.T) {
		root := t.TempDir()
		paths := auth.Paths{
			AppConfigPath: filepath.Join(root, "auth-app.json"),
			TokenPath:     filepath.Join(root, "auth-token.json"),
		}
		require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
			App:   auth.AppConfig{ClientID: "client-id"},
			Token: auth.TokenState{AccessToken: "access-token"},
		}))
		restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{
			beforeRevoke: func() {
				require.NoError(t, os.Remove(paths.TokenPath))
				require.NoError(t, os.Mkdir(paths.TokenPath, 0o700))
			},
		}, &fakeAuthReadinessChecker{})
		defer restore()

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
			"auth",
			"logout",
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "read auth token state")
	})

	t.Run("app config", func(t *testing.T) {
		root := t.TempDir()
		paths := auth.Paths{
			AppConfigPath: filepath.Join(root, "auth-app.json"),
			TokenPath:     filepath.Join(root, "auth-token.json"),
		}
		require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
			App:   auth.AppConfig{ClientID: "client-id"},
			Token: auth.TokenState{AccessToken: "access-token"},
		}))
		restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{
			beforeRevoke: func() {
				require.NoError(t, os.Remove(paths.AppConfigPath))
				require.NoError(t, os.Mkdir(paths.AppConfigPath, 0o700))
			},
		}, &fakeAuthReadinessChecker{})
		defer restore()

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
			"auth",
			"logout",
			"--forget-app",
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "read auth app config")
	})
}

func Test_AuthHelper_refresh_errors(t *testing.T) {
	authContext := authCommandContext{
		store: auth.NewStore(cliAuthTestPaths(t)),
	}

	t.Run("refresh token endpoint failure", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{err: errors.New("token endpoint down")},
			&fakeAuthReadinessChecker{},
		)
		defer restore()

		_, _, err := refreshAuthTokenState(context.Background(), authContext, auth.AppConfig{
			ClientID: "client-id",
			Scopes:   []string{"read"},
		}, auth.TokenState{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			GrantType:    authGrantAuthorizationCode,
		}, time.Second)

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeRefreshFailed), errorCode(err))
	})

	t.Run("refresh readiness failure", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{grant: auth.NewTokenGrant(
				"new-access-token",
				"new-refresh-token",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{"read"},
			)},
			&fakeAuthReadinessChecker{err: client.ErrTargetMismatch},
		)
		defer restore()

		_, _, err := refreshAuthTokenState(context.Background(), authContext, auth.AppConfig{
			ClientID: "client-id",
			Scopes:   []string{"read"},
		}, auth.TokenState{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			Actor:        "app",
			GrantType:    authGrantAuthorizationCode,
		}, time.Second)

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeTargetMismatch), errorCode(err))
	})

	t.Run("refresh missing scope", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{grant: auth.NewTokenGrant(
				"new-access-token",
				"new-refresh-token",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{},
			)},
			&fakeAuthReadinessChecker{},
		)
		defer restore()

		_, _, err := refreshAuthTokenState(context.Background(), authContext, auth.AppConfig{
			ClientID: "client-id",
			Scopes:   []string{"read"},
		}, auth.TokenState{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			GrantType:    authGrantAuthorizationCode,
		}, time.Second)

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeMissingScope), errorCode(err))
	})
}

func Test_AuthReadiness_and_status_helpers(t *testing.T) {
	t.Run("default readiness client constructor", func(t *testing.T) {
		require.NotNil(t, newAuthReadinessGraphQLClient("access-token", time.Second))
	})

	t.Run("actor mismatch", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{},
			&fakeAuthReadinessChecker{report: readyAuthReport("user")},
		)
		defer restore()

		_, err := requireAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken:   "access-token",
			ExpectedActor: "app",
		})

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeActorMismatch), errorCode(err))
	})

	t.Run("readiness error mapping", func(t *testing.T) {
		authErr := auth.NewError(auth.ErrorCodeNotConfigured, "missing")
		require.Same(t, authErr, mapAuthReadinessError(authErr))

		tokenErr := auth.NewTokenEndpointError(auth.ErrorCodeRefreshFailed, 401, "invalid_grant")
		require.Same(t, tokenErr, mapAuthReadinessError(tokenErr))

		targetErr := mapAuthReadinessError(client.ErrTargetNotConfigured)
		require.Equal(t, string(auth.ErrorCodeTargetMismatch), errorCode(targetErr))

		plainErr := errors.New("network down")
		require.Same(t, plainErr, mapAuthReadinessError(plainErr))
	})

	t.Run("default readiness success and resolve error", func(t *testing.T) {
		restore := useAuthReadinessGraphQLClient(t, authReadinessFakeGraphQLClient{
			"Viewer": `{
				"viewer": {
					"id": "user-id",
					"name": "Omer",
					"displayName": "Omer",
					"email": "omer@example.com",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}
			}`,
			"Teams": `{
				"teams": {
					"nodes": [{
						"id": "team-id",
						"key": "LIT",
						"name": "linctl-it",
						"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
					}],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}`,
		})
		defer restore()

		report, err := defaultCheckAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken:   "access-token",
			ExpectedActor: "app",
			ExpectedTarget: config.Target{
				OrgID:   "org-id",
				TeamKey: "LIT",
				TeamID:  "team-id",
			},
		})
		require.NoError(t, err)
		require.Equal(t, "app", report.Actor)
		require.Equal(t, "team-id", report.Target.Team.ID)

		_, err = defaultCheckAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken: "access-token",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, client.ErrTargetNotConfigured)
	})

	t.Run("write status quiet", func(t *testing.T) {
		var stdout bytes.Buffer
		command := &cobra.Command{}
		command.SetOut(&stdout)

		err := writeAuthStatus(command, &rootOptions{quiet: true}, authStatusReport{
			Token: authTokenReport{Status: "set"},
		})

		require.NoError(t, err)
		require.Empty(t, stdout.String())
	})

	require.Equal(t, "fallback", firstNonEmptyString("", "fallback"))
}

func Test_MergeResolvedAuthAppConfig_branches(t *testing.T) {
	t.Run("empty app", func(t *testing.T) {
		state := auth.State{App: auth.AppConfig{ClientID: "local-client-id"}}

		mergeResolvedAuthAppConfig(&state, "", auth.ProfileState{}, auth.AppConfig{})

		require.Equal(t, "local-client-id", state.App.ClientID)
	})

	t.Run("default profile", func(t *testing.T) {
		var state auth.State

		mergeResolvedAuthAppConfig(&state, "", auth.ProfileState{}, auth.AppConfig{ClientID: "env-client-id"})

		require.Equal(t, "env-client-id", state.App.ClientID)
	})

	t.Run("named profile", func(t *testing.T) {
		var state auth.State

		mergeResolvedAuthAppConfig(&state, "work", auth.ProfileState{}, auth.AppConfig{ClientID: "env-client-id"})

		require.Equal(t, "env-client-id", state.Profiles["work"].App.ClientID)
	})
}

func Test_ValidateCommandFlags_reports_limit_parse_error(t *testing.T) {
	command := &cobra.Command{}
	command.Flags().String("limit", "not-an-int", "")

	err := validateCommandFlags(command)

	require.Error(t, err)
	require.Contains(t, err.Error(), "read --limit")
}

func cliAuthTestPaths(t *testing.T) auth.Paths {
	t.Helper()
	root := t.TempDir()
	return auth.Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", "auth-app.json"),
		TokenPath:     filepath.Join(root, "state", "linctl", "auth-token.json"),
	}
}

func useAuthPaths(t *testing.T, paths auth.Paths) func() {
	t.Helper()
	original := authDefaultPaths
	authDefaultPaths = func(auth.Env) (auth.Paths, error) {
		return paths, nil
	}
	return func() {
		authDefaultPaths = original
	}
}

func useAuthPathsError(t *testing.T, err error) func() {
	t.Helper()
	original := authDefaultPaths
	authDefaultPaths = func(auth.Env) (auth.Paths, error) {
		return auth.Paths{}, err
	}
	return func() {
		authDefaultPaths = original
	}
}

func useAuthCommandHooks(
	t *testing.T,
	paths auth.Paths,
	oauthClient *fakeOAuthTokenClient,
	readiness *fakeAuthReadinessChecker,
) func() {
	t.Helper()
	restorePaths := useAuthPaths(t, paths)
	originalOAuthClient := newAuthOAuthClient
	originalReadiness := checkAuthReadiness
	newAuthOAuthClient = func() authOAuthClient {
		return oauthClient
	}
	checkAuthReadiness = readiness.check
	return func() {
		checkAuthReadiness = originalReadiness
		newAuthOAuthClient = originalOAuthClient
		restorePaths()
	}
}

func useAuthReadinessGraphQLClient(t *testing.T, client graphql.Client) func() {
	t.Helper()
	original := newAuthReadinessGraphQLClient
	newAuthReadinessGraphQLClient = func(string, time.Duration) graphql.Client {
		return client
	}
	return func() {
		newAuthReadinessGraphQLClient = original
	}
}

type fakeOAuthTokenClient struct {
	grant                    auth.TokenGrant
	err                      error
	clientCredentialsCalls   int
	clientCredentialsRequest oauth.ClientCredentialsRequest
	authorizationCodeCalls   int
	authorizationCodeRequest oauth.AuthorizationCodeRequest
	refreshTokenCalls        int
	refreshTokenRequest      oauth.RefreshTokenRequest
	revokeTokenCalls         int
	revokeTokenRequests      []oauth.RevocationRequest
	revokeTokenErr           error
	beforeRevoke             func()
}

func (client *fakeOAuthTokenClient) ClientCredentials(
	_ context.Context,
	request oauth.ClientCredentialsRequest,
) (auth.TokenGrant, error) {
	client.clientCredentialsCalls++
	client.clientCredentialsRequest = request
	if client.err != nil {
		return auth.TokenGrant{}, client.err
	}

	return client.grant, nil
}

func (client *fakeOAuthTokenClient) RefreshToken(
	_ context.Context,
	request oauth.RefreshTokenRequest,
) (auth.TokenGrant, error) {
	client.refreshTokenCalls++
	client.refreshTokenRequest = request
	if client.err != nil {
		return auth.TokenGrant{}, client.err
	}

	return client.grant, nil
}

func (client *fakeOAuthTokenClient) ExchangeAuthorizationCode(
	_ context.Context,
	request oauth.AuthorizationCodeRequest,
) (auth.TokenGrant, error) {
	client.authorizationCodeCalls++
	client.authorizationCodeRequest = request
	if client.err != nil {
		return auth.TokenGrant{}, client.err
	}

	return client.grant, nil
}

func (client *fakeOAuthTokenClient) RevokeToken(_ context.Context, request oauth.RevocationRequest) error {
	client.revokeTokenCalls++
	client.revokeTokenRequests = append(client.revokeTokenRequests, request)
	if client.beforeRevoke != nil {
		client.beforeRevoke()
		client.beforeRevoke = nil
	}
	if client.revokeTokenErr != nil {
		return client.revokeTokenErr
	}

	return nil
}

type fakeAuthReadinessChecker struct {
	report         authReadinessReport
	err            error
	accessToken    string
	expectedActor  string
	requiredScopes []string
	beforeReturn   func()
}

func (checker *fakeAuthReadinessChecker) check(
	_ context.Context,
	request authReadinessRequest,
) (authReadinessReport, error) {
	checker.accessToken = request.AccessToken
	checker.expectedActor = request.ExpectedActor
	checker.requiredScopes = request.RequiredScopes
	if checker.beforeReturn != nil {
		checker.beforeReturn()
	}
	if checker.err != nil {
		return authReadinessReport{}, checker.err
	}

	return checker.report, nil
}

func readyAuthReport(actor string) authReadinessReport {
	return authReadinessReport{
		Actor: actor,
		Target: client.ResolvedTarget{
			Org:      client.TargetOrg{ID: "org-id"},
			Team:     client.TargetTeam{ID: "team-id", Key: "LIT"},
			Expected: config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"},
			Resolved: config.Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"},
		},
	}
}

type authReadinessFakeGraphQLClient map[string]string

func (client authReadinessFakeGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	payload, ok := client[request.OpName]
	if !ok {
		return errors.New("missing fake response for " + request.OpName)
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

var _ = errors.New
