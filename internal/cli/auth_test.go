package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
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
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, []string{
		"--debug",
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
	logs := stderr.String()
	require.Contains(t, logs, `msg="auth token revoke failed"`)
	require.Contains(t, logs, `token_type=access_token`)
	require.Contains(t, logs, `error_code=INTERNAL`)
	require.Contains(t, logs, `msg="auth logout completed"`)
	require.Contains(t, logs, `revocation_failed=true`)
	require.NotContains(t, logs, "access-token")
	require.NotContains(t, logs, "client-secret")
	require.NotContains(t, logs, "revoked already")
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, got.App)
	require.Empty(t, got.Token)
}
