package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/stretchr/testify/require"
)

func Test_CommandRuntime_builds_graphql_client_with_oauth_bearer_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "test-token")
	t.Setenv("LINCTL_TOKEN", "legacy-token")
	t.Setenv("LINEAR_API_KEY", "legacy-token")
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`), 0o600))

	runtime, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.NoError(t, err)
	require.Equal(t, "Bearer test-token", runtimeGraphQLAuthorizationHeader(t, runtime))
}

func Test_CommandRuntime_reads_oauth_token_from_local_auth_state(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	authPaths := configureTestAuthEnvironment(t)
	require.NoError(t, auth.NewStore(authPaths).Save(context.Background(), auth.State{
		Token: auth.TokenState{
			AccessToken: "local-oauth-token",
		},
	}))
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`), 0o600))

	runtime, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.NoError(t, err)
	require.Equal(t, "Bearer local-oauth-token", runtimeGraphQLAuthorizationHeader(t, runtime))
}

func Test_CommandRuntime_reports_auth_default_paths_error(t *testing.T) {
	restore := useAuthPathsError(t, errors.New("paths unavailable"))
	defer restore()

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.Error(t, err)
	require.Contains(t, err.Error(), "paths unavailable")
}

func Test_CommandRuntime_reports_config_load_error(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)
	require.NoError(t, os.WriteFile(".linctl.toml", []byte("["), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.Error(t, err)
	require.Contains(t, err.Error(), "parse config")
}

func Test_CommandRuntime_requires_oauth_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeNotConfigured), errorCode(err))
}

func Test_CommandRuntime_reports_local_auth_state_load_error_after_env_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	root := t.TempDir()
	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "env-oauth-token")
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(root, "state"))
	paths, err := auth.DefaultPaths(nil)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o700))
	require.NoError(t, os.Mkdir(paths.AppConfigPath, 0o700))

	_, err = newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read auth app config")
}

func Test_NewRecoveringGraphQLClient_uses_defaults_and_empty_authorization(t *testing.T) {
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{})

	require.Empty(t, runtimeClient.authorizationHeader())
	require.NotNil(t, runtimeClient.client)
	require.NotNil(t, runtimeClient.oauthClient)
}

func Test_CommandRuntime_refreshes_expired_authorization_code_token_and_retries_once(t *testing.T) {
	paths := configureTestAuthEnvironment(t)
	store := auth.NewStore(paths)
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		TokenType:    "Bearer",
		Scopes:       []string{"read"},
		ExpiresAt:    &expiredAt,
		Actor:        "app",
		GrantType:    authGrantAuthorizationCode,
	}
	require.NoError(t, store.Save(context.Background(), auth.State{App: app, Token: token}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"rotated-access-token",
		"rotated-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         app,
		Store:       store,
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.refreshTokenCalls)
	require.Equal(t, "old-refresh-token", fakeOAuth.refreshTokenRequest.RefreshToken)
	require.Equal(t, []string{"expired-access-token", "rotated-access-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
	got, loadErr := store.Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "rotated-access-token", got.Token.AccessToken)
	require.Equal(t, "rotated-refresh-token", got.Token.RefreshToken)
	require.Equal(t, authGrantAuthorizationCode, got.Token.GrantType)
}

func Test_CommandRuntime_reacquires_expired_client_credentials_token_and_retries_once(t *testing.T) {
	paths := configureTestAuthEnvironment(t)
	store := auth.NewStore(paths)
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken: "expired-app-token",
		TokenType:   "Bearer",
		Scopes:      []string{"read"},
		ExpiresAt:   &expiredAt,
		Actor:       "app",
		GrantType:   authGrantClientCredentials,
	}
	require.NoError(t, store.Save(context.Background(), auth.State{App: app, Token: token}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         app,
		Store:       store,
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, []string{"expired-app-token", "fresh-app-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
	got, loadErr := store.Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "fresh-app-token", got.Token.AccessToken)
	require.Equal(t, authGrantClientCredentials, got.Token.GrantType)
}

func Test_CommandRuntime_returns_non_auth_error_after_pre_request_recovery(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken: "expired-app-token",
		Scopes:      []string{"read"},
		ExpiresAt:   &expiredAt,
		GrantType:   authGrantClientCredentials,
	}
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	boom := errors.New("linear unavailable")
	factory := &recordingRuntimeClientFactory{errors: []error{boom}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         app,
		Store:       auth.NewStore(cliAuthTestPaths(t)),
		Persist:     false,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorIs(t, err, boom)
	require.Equal(t, []string{"expired-app-token", "fresh-app-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_wraps_auth_failure_after_pre_request_recovery(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken: "expired-app-token",
		Scopes:      []string{"read"},
		ExpiresAt:   &expiredAt,
		GrantType:   authGrantClientCredentials,
	}
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         app,
		Store:       auth.NewStore(cliAuthTestPaths(t)),
		Persist:     false,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Equal(t, []string{"expired-app-token", "fresh-app-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_returns_non_auth_error_without_recovery(t *testing.T) {
	boom := errors.New("linear unavailable")
	factory := &recordingRuntimeClientFactory{errors: []error{boom}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:     auth.TokenState{AccessToken: "access-token"},
		App:       auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret"},
		Store:     auth.NewStore(cliAuthTestPaths(t)),
		Persist:   true,
		NewClient: factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorIs(t, err, boom)
	require.Equal(t, []string{"access-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_reports_token_persist_error_after_recovery(t *testing.T) {
	root := t.TempDir()
	tokenPath := filepath.Join(root, "auth-token-dir")
	require.NoError(t, os.Mkdir(tokenPath, 0o700))
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token: auth.TokenState{
			AccessToken: "expired-app-token",
			Scopes:      []string{"read"},
			ExpiresAt:   &expiredAt,
			GrantType:   authGrantClientCredentials,
		},
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Store: auth.NewStore(auth.Paths{
			AppConfigPath: filepath.Join(root, "auth-app.json"),
			TokenPath:     tokenPath,
		}),
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   (&recordingRuntimeClientFactory{}).newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read auth token state")
}

func Test_CommandRuntime_reacquires_client_credentials_token_after_401_once(t *testing.T) {
	store := auth.NewStore(cliAuthTestPaths(t))
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken: "stale-app-token",
		TokenType:   "Bearer",
		Scopes:      []string{"read"},
		Actor:       "app",
		GrantType:   authGrantClientCredentials,
	}
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed, nil}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         app,
		Store:       store,
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, []string{"stale-app-token", "fresh-app-token"}, factory.tokens)
	require.Equal(t, 2, factory.requestCalls)
}

func Test_CommandRuntime_returns_AUTH_REAUTH_REQUIRED_when_retried_token_is_rejected(t *testing.T) {
	store := auth.NewStore(cliAuthTestPaths(t))
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken: "stale-app-token",
		TokenType:   "Bearer",
		Scopes:      []string{"read"},
		Actor:       "app",
		GrantType:   authGrantClientCredentials,
	}
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed, client.ErrAuthFailed}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         app,
		Store:       store,
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, 2, factory.requestCalls)
}

func Test_CommandRuntime_returns_AUTH_REFRESH_FAILED_when_refresh_fails_without_retry_loop(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		GrantType:    authGrantAuthorizationCode,
	}
	fakeOAuth := &fakeOAuthTokenClient{err: errors.New("token endpoint unavailable")}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         auth.AppConfig{ClientID: "client-id"},
		Store:       auth.NewStore(cliAuthTestPaths(t)),
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeRefreshFailed), errorCode(err))
	require.Equal(t, 1, fakeOAuth.refreshTokenCalls)
	require.Equal(t, 0, factory.requestCalls)
}

func Test_CommandRuntime_returns_AUTH_REAUTH_REQUIRED_when_app_reacquire_fails_without_retry_loop(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken: "expired-app-token",
		ExpiresAt:   &expiredAt,
		GrantType:   authGrantClientCredentials,
	}
	fakeOAuth := &fakeOAuthTokenClient{err: errors.New("token endpoint unavailable")}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:       token,
		App:         auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret"},
		Store:       auth.NewStore(cliAuthTestPaths(t)),
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, 0, factory.requestCalls)
}

func Test_CommandRuntime_returns_recovery_error_after_401(t *testing.T) {
	fakeOAuth := &fakeOAuthTokenClient{err: errors.New("token endpoint unavailable")}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token: auth.TokenState{
			AccessToken: "stale-app-token",
			Scopes:      []string{"read"},
			GrantType:   authGrantClientCredentials,
		},
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       []string{"read"},
		},
		Store:       auth.NewStore(cliAuthTestPaths(t)),
		Persist:     true,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_refresh_authorization_code_edge_cases(t *testing.T) {
	t.Run("missing refresh state", func(t *testing.T) {
		runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token: auth.TokenState{},
			App:   auth.AppConfig{ClientID: "client-id"},
			Store: auth.NewStore(cliAuthTestPaths(t)),
		})

		_, err := runtimeClient.refreshAuthorizationCode(context.Background())

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	})

	t.Run("keeps old refresh token when endpoint omits one", func(t *testing.T) {
		fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
			"new-access-token",
			"",
			"Bearer",
			time.Now().Add(time.Hour),
			[]string{"read"},
		)}
		runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token: auth.TokenState{
				RefreshToken: "old-refresh-token",
				Actor:        "user",
			},
			App: auth.AppConfig{
				ClientID: "client-id",
				Scopes:   []string{"read"},
			},
			Store:       auth.NewStore(cliAuthTestPaths(t)),
			OAuthClient: fakeOAuth,
		})

		token, err := runtimeClient.refreshAuthorizationCode(context.Background())

		require.NoError(t, err)
		require.Equal(t, "old-refresh-token", token.RefreshToken)
		require.Equal(t, "user", token.Actor)
	})

	t.Run("missing scope", func(t *testing.T) {
		fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenGrant(
			"new-access-token",
			"new-refresh-token",
			"Bearer",
			time.Now().Add(time.Hour),
			[]string{},
		)}
		runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token: auth.TokenState{RefreshToken: "old-refresh-token"},
			App: auth.AppConfig{
				ClientID: "client-id",
				Scopes:   []string{"read"},
			},
			Store:       auth.NewStore(cliAuthTestPaths(t)),
			OAuthClient: fakeOAuth,
		})

		_, err := runtimeClient.refreshAuthorizationCode(context.Background())

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeMissingScope), errorCode(err))
	})
}

func Test_CommandRuntime_reacquire_client_credentials_requires_app_config(t *testing.T) {
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token: auth.TokenState{GrantType: authGrantClientCredentials},
		Store: auth.NewStore(cliAuthTestPaths(t)),
	})

	_, err := runtimeClient.reacquireClientCredentials(context.Background())

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
}

func runtimeGraphQLAuthorizationHeader(t *testing.T, runtime commandRuntime) string {
	t.Helper()

	if client, ok := runtime.graphqlClient.(*recoveringGraphQLClient); ok {
		return client.authorizationHeader()
	}
	value := reflect.ValueOf(runtime.graphqlClient)
	require.Equal(t, "*client.Transport", value.Type().String())

	return value.Elem().FieldByName("token").FieldByName("authorization").String()
}

func configureTestAuthEnvironment(t *testing.T) auth.Paths {
	t.Helper()

	root := t.TempDir()
	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "")
	t.Setenv("LINCTL_TOKEN", "")
	t.Setenv("LINEAR_API_KEY", "")
	switch runtime.GOOS {
	case "windows":
		t.Setenv("HOME", root)
		t.Setenv("USERPROFILE", root)
		t.Setenv("APPDATA", filepath.Join(root, "appdata"))
		t.Setenv("LOCALAPPDATA", filepath.Join(root, "localappdata"))
	case "darwin":
		t.Setenv("HOME", root)
	default:
		t.Setenv("HOME", root)
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
		t.Setenv("XDG_STATE_HOME", filepath.Join(root, "state"))
	}

	paths, err := auth.DefaultPaths(nil)
	require.NoError(t, err)

	return paths
}

type recordingRuntimeClientFactory struct {
	tokens       []string
	errors       []error
	requestCalls int
}

func (factory *recordingRuntimeClientFactory) newClient(accessToken string) graphql.Client {
	factory.tokens = append(factory.tokens, accessToken)
	return recordingRuntimeGraphQLClient{factory: factory}
}

type recordingRuntimeGraphQLClient struct {
	factory *recordingRuntimeClientFactory
}

func (client recordingRuntimeGraphQLClient) MakeRequest(
	_ context.Context,
	_ *graphql.Request,
	_ *graphql.Response,
) error {
	client.factory.requestCalls++
	if len(client.factory.errors) == 0 {
		return nil
	}
	err := client.factory.errors[0]
	client.factory.errors = client.factory.errors[1:]

	return err
}
