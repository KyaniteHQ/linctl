package auth

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SelectSession_resolves_local_auth_state(t *testing.T) {
	t.Parallel()
	paths := testAuthPaths(t)
	require.NoError(t, NewStore(paths).Save(context.Background(), State{
		App: AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Token: TokenState{
			AccessToken:  "local-oauth-token",
			RefreshToken: "refresh-token",
		},
	}))

	session, err := SelectSession(context.Background(), SessionRequest{
		Store: NewStore(paths),
	})

	require.NoError(t, err)
	require.Equal(t, "local-oauth-token", session.Token.AccessToken)
	require.Equal(t, "client-id", session.App.ClientID)
	require.Equal(t, CredentialKindLocalAccessToken, session.CredentialKind)
	require.Equal(t, "local", session.TokenSource)
	require.True(t, session.PersistentToken)
}

func Test_SelectSession_env_oauth_app_material_overrides_local_app(t *testing.T) {
	t.Parallel()
	paths := testAuthPaths(t)
	require.NoError(t, NewStore(paths).Save(context.Background(), State{
		Profiles: map[string]ProfileState{
			"work": {
				App: AppConfig{
					ClientID:     "local-client-id",
					ClientSecret: "local-client-secret",
					RedirectURI:  "http://127.0.0.1:8484/local",
					Scopes:       []string{"read"},
				},
				Token: TokenState{AccessToken: "local-oauth-token"},
			},
		},
	}))

	session, err := SelectSession(context.Background(), SessionRequest{
		Env: staticEnv{
			oauthClientIDEnv:     "env-client-id",
			oauthClientSecretEnv: "env-client-secret",
			oauthRedirectURIEnv:  "http://127.0.0.1:8484/env",
			oauthScopesEnv:       "read, write\ncomments:create",
		},
		Store:   NewStore(paths),
		Profile: "work",
	})

	require.NoError(t, err)
	require.Equal(t, "local-oauth-token", session.Token.AccessToken)
	require.Equal(t, AppConfig{
		ClientID:     "env-client-id",
		ClientSecret: "env-client-secret",
		RedirectURI:  "http://127.0.0.1:8484/env",
		Scopes:       []string{"read", "write", "comments:create"},
	}, session.App)
}

func Test_SelectSession_env_token_is_non_persistent(t *testing.T) {
	t.Parallel()
	paths := testAuthPaths(t)
	require.NoError(t, NewStore(paths).Save(context.Background(), State{
		App: AppConfig{
			ClientID:     "local-client-id",
			ClientSecret: "local-client-secret",
		},
		Token: TokenState{
			AccessToken:  "local-oauth-token",
			RefreshToken: "refresh-token",
		},
	}))

	session, err := SelectSession(context.Background(), SessionRequest{
		Env:   staticEnv{oauthAccessTokenEnv: "env-oauth-token"},
		Store: NewStore(paths),
	})

	require.NoError(t, err)
	require.Equal(t, "env-oauth-token", session.Token.AccessToken)
	require.Empty(t, session.Token.RefreshToken)
	require.Empty(t, session.App)
	require.Empty(t, session.State)
	require.Equal(t, CredentialKindInjectedAccessToken, session.CredentialKind)
	require.Equal(t, "env", session.TokenSource)
	require.False(t, session.PersistentToken)
}

func Test_SelectSession_reports_missing_token_source(t *testing.T) {
	t.Parallel()

	session, err := SelectSession(context.Background(), SessionRequest{
		Store: NewStore(testAuthPaths(t)),
	})

	require.NoError(t, err)
	require.Equal(t, CredentialKindMissing, session.CredentialKind)
	require.Equal(t, "missing", session.TokenSource)
	require.False(t, session.PersistentToken)
}

func Test_CredentialKindFromToken_uses_only_explicit_grant_variants(t *testing.T) {
	t.Parallel()

	require.Equal(t, CredentialKindAuthorizationCode, CredentialKindFromToken(TokenState{
		AccessToken: "oauth-access-token",
		GrantType:   "authorization_code",
	}))
	require.Equal(t, CredentialKindClientCredentials, CredentialKindFromToken(TokenState{
		AccessToken: "oauth-access-token",
		GrantType:   "client_credentials",
	}))
	require.True(t, CredentialKindAuthorizationCode.Recoverable())
	require.True(t, CredentialKindClientCredentials.Recoverable())
	require.False(t, CredentialKindInjectedAccessToken.Recoverable())
	require.False(t, CredentialKindLocalAccessToken.Recoverable())
}

func Test_SelectSession_rejects_personal_api_key_shapes(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name  string
		env   Env
		state State
	}{
		{
			name: "env token",
			env:  staticEnv{oauthAccessTokenEnv: "lin_api_personal_key"},
		},
		{
			name:  "local token",
			state: State{Token: TokenState{AccessToken: "lin_api_personal_key"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			paths := testAuthPaths(t)
			require.NoError(t, NewStore(paths).Save(context.Background(), test.state))

			_, err := SelectSession(context.Background(), SessionRequest{
				Env:   test.env,
				Store: NewStore(paths),
			})

			require.ErrorContains(t, err, "personal API key")
		})
	}
}

func Test_SelectSession_reports_context_and_store_errors(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := SelectSession(ctx, SessionRequest{})

		require.ErrorContains(t, err, "select auth session context")
	})

	t.Run("store read error", func(t *testing.T) {
		root := t.TempDir()
		appConfigPath := filepath.Join(root, "auth-app-dir")
		require.NoError(t, os.Mkdir(appConfigPath, 0o700))

		_, err := SelectSession(context.Background(), SessionRequest{
			Store: NewStore(Paths{
				AppConfigPath: appConfigPath,
				TokenPath:     filepath.Join(root, "auth-token.json"),
			}),
		})

		require.ErrorContains(t, err, "read auth app config")
	})
}

func Test_SelectLocalSession_loads_local_state_despite_injected_token(t *testing.T) {
	paths := testAuthPaths(t)
	require.NoError(t, NewStore(paths).Save(context.Background(), State{
		App:   AppConfig{ClientID: "local-client-id"},
		Token: TokenState{AccessToken: "local-access-token"},
	}))

	session, err := SelectLocalSession(context.Background(), SessionRequest{
		Env:   staticEnv{oauthAccessTokenEnv: "injected-access-token"},
		Store: NewStore(paths),
	})

	require.NoError(t, err)
	require.Equal(t, "local-client-id", session.App.ClientID)
	require.Equal(t, "local-access-token", session.Token.AccessToken)
	require.Equal(t, CredentialKindLocalAccessToken, session.CredentialKind)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = SelectLocalSession(ctx, SessionRequest{})
	require.ErrorContains(t, err, "select local auth session context")
}

func Test_AppConfigFromEnv_resolves_oauth_material(t *testing.T) {
	t.Parallel()

	app := AppConfigFromEnv(staticEnv{
		oauthClientIDEnv:     "client-id",
		oauthClientSecretEnv: "client-secret",
		oauthRedirectURIEnv:  "http://127.0.0.1:8484/callback",
		oauthScopesEnv:       "read,write comments:create",
		"LINCTL_TOKEN":       "legacy-token",
		"LINEAR_API_KEY":     "legacy-token",
	})

	require.Equal(t, AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		Scopes:       []string{"read", "write", "comments:create"},
	}, app)
}

func Test_MergeAppConfig_overlays_explicit_fields(t *testing.T) {
	t.Parallel()

	got := MergeAppConfig(
		AppConfig{
			ClientID:     "local-client-id",
			ClientSecret: "local-client-secret",
			RedirectURI:  "http://127.0.0.1:8484/local",
			Scopes:       []string{"read"},
		},
		AppConfig{
			ClientID: "env-client-id",
			Scopes:   []string{"read", "write"},
		},
	)

	require.Equal(t, AppConfig{
		ClientID:     "env-client-id",
		ClientSecret: "local-client-secret",
		RedirectURI:  "http://127.0.0.1:8484/local",
		Scopes:       []string{"read", "write"},
	}, got)
}

func testAuthPaths(t *testing.T) Paths {
	t.Helper()
	root := t.TempDir()
	return Paths{
		AppConfigPath: filepath.Join(root, "config", "linctl", "auth-app.json"),
		TokenPath:     filepath.Join(root, "state", "linctl", "auth-token.json"),
	}
}
