package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CommandRuntime_builds_graphql_client_with_oauth_bearer_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)
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

	linear := useLinearGraphQL(t, false, true)

	runtime, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.NoError(t, err)
	require.NoError(t, runtime.graphqlClient.MakeRequest(
		context.Background(),
		&graphql.Request{OpName: "Viewer", Query: "query Viewer { viewer { id } }"},
		&graphql.Response{},
	))
	require.True(t, slices.Equal([]string{"Bearer test-token"}, linear.authorizations), "unexpected Authorization header")
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

	linear := useLinearGraphQL(t, false, true)

	runtime, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.NoError(t, err)
	require.NoError(t, runtime.graphqlClient.MakeRequest(
		context.Background(),
		&graphql.Request{OpName: "Viewer", Query: "query Viewer { viewer { id } }"},
		&graphql.Response{},
	))
	require.True(t, slices.Equal([]string{"Bearer local-oauth-token"}, linear.authorizations), "unexpected Authorization header")
}

func Test_CommandRuntime_reports_auth_default_paths_error(t *testing.T) {
	restore := useAuthPathsError(t, errors.New("paths unavailable"))
	defer restore()

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.ErrorContains(t, err, "paths unavailable")
}

func Test_CommandRuntime_reports_config_load_error(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)
	require.NoError(t, os.WriteFile(".linctl.toml", []byte("["), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.ErrorContains(t, err, "parse config")
}

func Test_ResolveConfig_uses_explicit_repo_config_outside_its_directory(t *testing.T) {
	workingDir := t.TempDir()
	t.Chdir(workingDir)
	configPath := filepath.Join(t.TempDir(), "case-scraper.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[target]
org_id = "org-id"
team_key = "NIN"
team_id = "team-id"
project_id = "case-scraper-project"
`), 0o600))

	resolved, source, err := resolveConfigWithSource(
		context.Background(),
		&rootOptions{configPath: configPath, configPathExplicit: true},
	)

	require.NoError(t, err)
	require.Equal(t, "case-scraper-project", resolved.Target.ProjectID)
	require.Equal(t, configPath, source.Path)
	require.Equal(t, repoConfigLoaded, source.Status)
}

func Test_SelectRepoConfig_reports_missing_default_in_current_directory(t *testing.T) {
	workingDir := t.TempDir()
	t.Chdir(workingDir)

	source, err := selectRepoConfig(&rootOptions{})

	require.NoError(t, err)
	require.Equal(t, filepath.Join(workingDir, ".linctl.toml"), source.Path)
	require.Equal(t, repoConfigMissing, source.Status)
}

func Test_Execute_reports_a_removed_working_directory(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("the removed working directory error is pinned on Linux getcwd")
	}
	configureTestAuthEnvironment(t)
	dir := filepath.Join(t.TempDir(), "removed")
	require.NoError(t, os.Mkdir(dir, 0o700))
	t.Chdir(dir)
	require.NoError(t, os.Remove(dir))

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{"whoami"})

	require.ErrorContains(t, err, "resolve repo config path .linctl.toml")
}

func Test_Execute_reports_a_config_path_under_a_regular_file(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows reports a path under a regular file as missing")
	}
	t.Chdir(t.TempDir())
	configureTestAuthEnvironment(t)
	require.NoError(t, os.WriteFile("plain-file", nil, 0o600))

	err := execute(
		context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{},
		[]string{"--config", filepath.Join("plain-file", "config.toml"), "whoami"},
	)

	require.ErrorContains(t, err, "inspect repo config")
	require.ErrorContains(t, err, "not a directory")
}

func Test_CommandRuntime_requires_oauth_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeNotConfigured), errorCode(err))
}

func Test_CommandRuntime_does_not_load_local_auth_state_for_env_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	paths := configureTestAuthEnvironment(t)
	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "env-oauth-token")
	require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o700))
	require.NoError(t, os.Mkdir(paths.AppConfigPath, 0o700))

	linear := useLinearGraphQL(t, false, true)

	runtime, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.NoError(t, err)
	require.NoError(t, runtime.graphqlClient.MakeRequest(
		context.Background(),
		&graphql.Request{OpName: "Viewer", Query: "query Viewer { viewer { id } }"},
		&graphql.Response{},
	))
	require.True(t, slices.Equal([]string{"Bearer env-oauth-token"}, linear.authorizations), "unexpected Authorization header")
}

func Test_CommandRuntime_reports_local_auth_state_load_error_without_env_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	paths := configureTestAuthEnvironment(t)
	require.NoError(t, os.MkdirAll(filepath.Dir(paths.AppConfigPath), 0o700))
	require.NoError(t, os.Mkdir(paths.AppConfigPath, 0o700))

	_, err := newCommandRuntime(context.Background(), &rootOptions{timeout: time.Second})

	require.ErrorContains(t, err, "read auth app config")
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"rotated-access-token",
		"rotated-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            app,
		Store:          store,
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            app,
		Store:          store,
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	boom := errors.New("linear unavailable")
	factory := &recordingRuntimeClientFactory{errors: []error{boom}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            app,
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
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
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            app,
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
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
		CredentialKind: auth.CredentialKindLocalAccessToken,
		Token:          auth.TokenState{AccessToken: "access-token"},
		App:            auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret"},
		Store:          auth.NewStore(cliAuthTestPaths(t)),
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorIs(t, err, boom)
	require.Equal(t, []string{"access-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_surfaces_rejected_injected_token_without_identity_substitution(t *testing.T) {
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"substituted-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:          auth.TokenState{AccessToken: "injected-access-token"},
		CredentialKind: auth.CredentialKindInjectedAccessToken,
		App: auth.AppConfig{
			ClientID:     "local-client-id",
			ClientSecret: "local-client-secret",
		},
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(
		context.Background(),
		&graphql.Request{Query: "query Test { viewer { id } }"},
		&graphql.Response{},
	)

	require.ErrorIs(t, err, client.ErrAuthFailed)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Zero(t, fakeOAuth.clientCredentialsCalls)
	require.Zero(t, fakeOAuth.refreshTokenCalls)
	require.Equal(t, []string{"injected-access-token"}, factory.tokens)
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_reports_token_persist_error_after_recovery(t *testing.T) {
	root := t.TempDir()
	tokenPath := filepath.Join(root, "auth-token-dir")
	require.NoError(t, os.Mkdir(tokenPath, 0o700))
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
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
		Logger:      discardLogger,
		OAuthClient: fakeOAuth,
		NewClient:   (&recordingRuntimeClientFactory{}).newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorContains(t, err, "read auth token state")
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
	require.NoError(t, store.SaveTokenState(context.Background(), "", token))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed, nil}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            app,
		Store:          store,
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.clientCredentialsCalls)
	require.Equal(t, []string{"stale-app-token", "fresh-app-token"}, factory.tokens)
	require.Equal(t, 2, factory.requestCalls)
}

func Test_CommandRuntime_logs_auth_failure_token_recovery_without_secrets(t *testing.T) {
	store := auth.NewStore(cliAuthTestPaths(t))
	app := auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken: "stale-app-token",
		Scopes:      []string{"read"},
		Actor:       "app",
		GrantType:   authGrantClientCredentials,
	}
	require.NoError(t, store.SaveTokenState(context.Background(), "", token))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed, nil}}
	var logs bytes.Buffer
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            app,
		Store:          store,
		Logger:         newDiagnosticLogger(true, false, &logs),
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	output := logs.String()
	require.Contains(t, output, `msg="auth token recovery started"`)
	require.Contains(t, output, "reason=auth_failed")
	require.Contains(t, output, "grant_type=client_credentials")
	require.Contains(t, output, `msg="auth token recovery succeeded"`)
	require.NotContains(t, output, "stale-app-token")
	require.NotContains(t, output, "fresh-app-token")
	require.NotContains(t, output, "client-secret")
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
	require.NoError(t, store.SaveTokenState(context.Background(), "", token))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"fresh-app-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour).UTC().Truncate(time.Second),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{errors: []error{client.ErrAuthFailed, client.ErrAuthFailed}}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            app,
		Store:          store,
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
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
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            auth.AppConfig{ClientID: "client-id"},
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeRefreshFailed), errorCode(err))
	require.Equal(t, 1, fakeOAuth.refreshTokenCalls)
	require.Equal(t, 0, factory.requestCalls)
}

func Test_CommandRuntime_logs_expired_token_recovery_without_secrets(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read"},
	}
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		Scopes:       []string{"read"},
		ExpiresAt:    &expiredAt,
		Actor:        "app",
		GrantType:    authGrantAuthorizationCode,
	}
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"rotated-access-token",
		"rotated-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	var logs bytes.Buffer
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            app,
		Store:          persistedRuntimeStore(t, "work", token),
		Profile:        "work",
		Logger:         newDiagnosticLogger(true, false, &logs),
		OAuthClient:    fakeOAuth,
		NewClient:      (&recordingRuntimeClientFactory{}).newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	output := logs.String()
	require.Contains(t, output, `msg="auth token recovery started"`)
	require.Contains(t, output, "reason=expired")
	require.Contains(t, output, "grant_type=authorization_code")
	require.Contains(t, output, "actor=app")
	require.Contains(t, output, "profile=work")
	require.Contains(t, output, `msg="auth token recovery succeeded"`)
	require.NotContains(t, output, "expired-access-token")
	require.NotContains(t, output, "rotated-access-token")
	require.NotContains(t, output, "old-refresh-token")
	require.NotContains(t, output, "rotated-refresh-token")
	require.NotContains(t, output, "client-secret")
}

func Test_CommandRuntime_logs_token_recovery_failure_without_error_details(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		GrantType:    authGrantAuthorizationCode,
	}
	fakeOAuth := &fakeOAuthTokenClient{err: errors.New("token endpoint unavailable")}
	var logs bytes.Buffer
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            auth.AppConfig{ClientID: "client-id"},
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         newDiagnosticLogger(true, false, &logs),
		OAuthClient:    fakeOAuth,
		NewClient:      (&recordingRuntimeClientFactory{}).newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeRefreshFailed), errorCode(err))
	output := logs.String()
	require.Contains(t, output, `msg="auth token recovery started"`)
	require.Contains(t, output, `msg="auth token recovery failed"`)
	require.Contains(t, output, "reason=expired")
	require.Contains(t, output, "phase=exchange")
	require.Contains(t, output, "error_code=AUTH_REFRESH_FAILED")
	require.NotContains(t, output, "token endpoint unavailable")
	require.NotContains(t, output, "expired-access-token")
	require.NotContains(t, output, "old-refresh-token")
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
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		App:            auth.AppConfig{ClientID: "client-id", ClientSecret: "client-secret"},
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
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
		CredentialKind: auth.CredentialKindClientCredentials,
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
		Logger:      discardLogger,
		OAuthClient: fakeOAuth,
		NewClient:   factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Equal(t, 1, factory.requestCalls)
}

func Test_CommandRuntime_requires_login_when_an_expired_authorization_code_token_has_no_refresh_token(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken: "expired-access-token",
		ExpiresAt:   &expiredAt,
		GrantType:   authGrantAuthorizationCode,
	}
	fakeOAuth := &fakeOAuthTokenClient{}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            auth.AppConfig{ClientID: "client-id"},
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorContains(t, err, "run linctl auth login")
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Zero(t, fakeOAuth.refreshTokenCalls)
	require.Zero(t, factory.requestCalls)
}

func Test_CommandRuntime_keeps_the_refresh_token_when_the_endpoint_omits_one(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		Actor:        "user",
		GrantType:    authGrantAuthorizationCode,
	}
	store := persistedRuntimeStore(t, "", token)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"new-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}},
		Store:          store,
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.NoError(t, err)
	require.Equal(t, []string{"expired-access-token", "new-access-token"}, factory.tokens)
	got, loadErr := store.Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "new-access-token", got.Token.AccessToken)
	require.Equal(t, "old-refresh-token", got.Token.RefreshToken)
	require.Equal(t, "user", got.Token.Actor)
}

func Test_CommandRuntime_reports_a_refresh_that_drops_a_required_scope(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		GrantType:    authGrantAuthorizationCode,
	}
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"new-access-token",
		"new-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{},
	)}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          token,
		App:            auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}},
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeMissingScope), errorCode(err))
	require.Zero(t, factory.requestCalls)
}

func Test_CommandRuntime_requires_app_config_to_reacquire_an_expired_app_token(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	token := auth.TokenState{
		AccessToken: "expired-app-token",
		ExpiresAt:   &expiredAt,
		GrantType:   authGrantClientCredentials,
	}
	fakeOAuth := &fakeOAuthTokenClient{}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindClientCredentials,
		Token:          token,
		Store:          persistedRuntimeStore(t, "", token),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorContains(t, err, "run linctl auth app")
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Zero(t, fakeOAuth.clientCredentialsCalls)
	require.Zero(t, factory.requestCalls)
}

func Test_CommandRuntime_refuses_a_persisted_token_without_a_recovery_grant(t *testing.T) {
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	observed := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		GrantType:    authGrantAuthorizationCode,
	}
	persisted := auth.TokenState{
		AccessToken: "replaced-access-token",
		ExpiresAt:   &expiredAt,
	}
	fakeOAuth := &fakeOAuthTokenClient{}
	factory := &recordingRuntimeClientFactory{}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		CredentialKind: auth.CredentialKindAuthorizationCode,
		Token:          observed,
		App:            auth.AppConfig{ClientID: "client-id"},
		Store:          persistedRuntimeStore(t, "", persisted),
		Logger:         discardLogger,
		OAuthClient:    fakeOAuth,
		NewClient:      factory.newClient,
	})

	err := runtimeClient.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &graphql.Response{})

	require.ErrorContains(t, err, "persisted OAuth token is not recoverable")
	require.Equal(t, string(auth.ErrorCodeReauthRequired), errorCode(err))
	require.Zero(t, fakeOAuth.refreshTokenCalls)
	require.Zero(t, factory.requestCalls)
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

func persistedRuntimeStore(t *testing.T, profile string, token auth.TokenState) auth.Store {
	t.Helper()
	store := auth.NewStore(cliAuthTestPaths(t))
	require.NoError(t, store.SaveTokenState(context.Background(), profile, token))

	return store
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
