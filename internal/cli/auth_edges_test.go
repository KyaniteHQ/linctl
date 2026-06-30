package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
)

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
			fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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

func Test_AuthStatus_debug_logs_readiness_mismatch_without_secrets(t *testing.T) {
	paths := cliAuthTestPaths(t)
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Token: auth.TokenState{
			AccessToken: "current-access-token",
			TokenType:   "Bearer",
			Scopes:      []string{"read"},
			ExpiresAt:   &expiresAt,
			Actor:       "app",
		},
	}))
	restore := useAuthCommandHooks(t, paths, &fakeOAuthTokenClient{}, &fakeAuthReadinessChecker{
		err: client.ErrTargetMismatch,
	})
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"--debug",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"status",
	})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeTargetMismatch), errorCode(err))
	output := stderr.String()
	require.Contains(t, output, `msg="auth readiness check started"`)
	require.Contains(t, output, `msg="auth readiness check failed"`)
	require.Contains(t, output, `error_code=AUTH_TARGET_MISMATCH`)
	require.NotContains(t, output, "current-access-token")
	require.NotContains(t, output, "client-secret")
}

func Test_AuthStatus_acquires_client_credentials_when_token_state_is_missing(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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

func Test_AuthRefresh_debug_logs_readiness_mismatch_without_secrets(t *testing.T) {
	paths := cliAuthTestPaths(t)
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
			Actor:        "app",
			GrantType:    authGrantAuthorizationCode,
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"new-access-token",
		"new-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{
		err: client.ErrTargetMismatch,
	})
	defer restore()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"--debug",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"refresh",
	})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeTargetMismatch), errorCode(err))
	output := stderr.String()
	require.Contains(t, output, `msg="auth readiness check started"`)
	require.Contains(t, output, `msg="auth readiness check failed"`)
	require.Contains(t, output, `error_code=AUTH_TARGET_MISMATCH`)
	require.NotContains(t, output, "old-access-token")
	require.NotContains(t, output, "old-refresh-token")
	require.NotContains(t, output, "new-access-token")
	require.NotContains(t, output, "new-refresh-token")
	require.NotContains(t, output, "client-secret")
	got, loadErr := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "old-access-token", got.Token.AccessToken)
	require.Equal(t, "old-refresh-token", got.Token.RefreshToken)
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
			&fakeOAuthTokenClient{grant: auth.NewTokenState(
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
			&fakeOAuthTokenClient{grant: auth.NewTokenState(
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
