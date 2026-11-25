package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

func Test_AuthLogin_builds_authorization_url_with_app_actor_by_default(t *testing.T) {
	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "injected-access-token")
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"auth",
		"login",
	})

	require.NoError(t, err)
	report := decodeAuthLoginStartReport(t, stdout.Bytes())
	authorizeURL := authLoginAuthorizeURL(t, report.AuthorizeURL)
	query := authorizeURL.Query()
	require.Equal(t, "https", authorizeURL.Scheme)
	require.Equal(t, "linear.app", authorizeURL.Host)
	require.Equal(t, "/oauth/authorize", authorizeURL.Path)
	require.Equal(t, "code", query.Get("response_type"))
	require.Equal(t, "client-id", query.Get("client_id"))
	require.Equal(t, "http://127.0.0.1:8484/callback", query.Get("redirect_uri"))
	require.Equal(t, "read,write", query.Get("scope"))
	require.Equal(t, "state-123", query.Get("state"))
	require.Equal(t, "challenge-123", query.Get("code_challenge"))
	require.Equal(t, "S256", query.Get("code_challenge_method"))
	require.Equal(t, "app", query.Get("actor"))
	require.Equal(t, "consent", query.Get("prompt"))
}

func Test_AuthLogin_builds_authorization_url_with_user_actor(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"--json",
		"auth",
		"login",
		"--actor", "user",
	})

	require.NoError(t, err)
	query := authLoginAuthorizeURL(t, decodeAuthLoginStartReport(t, stdout.Bytes()).AuthorizeURL).Query()
	require.Equal(t, "user", query.Get("actor"))
}

func Test_AuthLogin_requests_initiative_scopes_only_for_the_app_actor_by_default(t *testing.T) {
	tests := []struct {
		actor                 string
		wantsInitiativeScopes bool
	}{
		{actor: "app", wantsInitiativeScopes: true},
		{actor: "user", wantsInitiativeScopes: false},
	}
	for _, test := range tests {
		t.Run(test.actor, func(t *testing.T) {
			paths := cliAuthTestPaths(t)
			require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
				ClientID:    "client-id",
				RedirectURI: "http://127.0.0.1:8484/callback",
			}))
			restorePaths := useAuthPaths(t, paths)
			defer restorePaths()
			restoreLogin := useAuthLoginGeneratedValues(t)
			defer restoreLogin()
			var stdout bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
				"--json",
				"auth",
				"login",
				"--actor", test.actor,
			})

			require.NoError(t, err)
			query := authLoginAuthorizeURL(t, decodeAuthLoginStartReport(t, stdout.Bytes()).AuthorizeURL).Query()
			scopes := strings.Split(query.Get("scope"), ",")
			require.Contains(t, scopes, "read")
			require.Equal(t, test.wantsInitiativeScopes, slices.Contains(scopes, "initiative:read"))
			require.Equal(t, test.wantsInitiativeScopes, slices.Contains(scopes, "initiative:write"))
		})
	}
}

func Test_AuthLogin_quiet_start_prints_nothing(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, []string{
		"--quiet",
		"auth",
		"login",
	})

	require.NoError(t, err)
	require.Empty(t, stdout.String())
	require.Empty(t, stderr.String())
}

func Test_AuthLogin_human_start_prints_authorization_url(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stdout bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{
		"auth",
		"login",
	})

	require.NoError(t, err)
	require.Contains(t, stdout.String(), "https://linear.app/oauth/authorize")
	require.Contains(t, stdout.String(), "actor=app")
}

func Test_AuthLogin_rejects_invalid_actor(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"auth",
		"login",
		"--actor", "service-account",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), `"error_code":"AUTH_ACTOR_MISMATCH"`)
}

func Test_AuthLogin_generates_fresh_state_and_pkce_per_run(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	start := func() url.Values {
		var stdout bytes.Buffer
		err := execute(context.Background(), BuildInfo{}, nil, &stdout, &bytes.Buffer{}, []string{"--json", "auth", "login"})
		require.NoError(t, err)

		return authLoginAuthorizeURL(t, decodeAuthLoginStartReport(t, stdout.Bytes()).AuthorizeURL).Query()
	}

	first := start()
	second := start()

	require.NotEmpty(t, first.Get("state"))
	require.NotEmpty(t, first.Get("code_challenge"))
	require.Equal(t, "S256", first.Get("code_challenge_method"))
	require.NotEqual(t, first.Get("state"), second.Get("state"))
	require.NotEqual(t, first.Get("code_challenge"), second.Get("code_challenge"))
}

func Test_AuthLogin_reports_app_configuration_errors(t *testing.T) {
	t.Run("missing client id", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		restorePaths := useAuthPaths(t, paths)
		defer restorePaths()
		var stderr bytes.Buffer

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
			"auth",
			"login",
		})

		require.Error(t, err)
		require.Contains(t, stderr.String(), "missing OAuth client id")
	})

	t.Run("missing redirect uri", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
			ClientID: "client-id",
		}))
		restorePaths := useAuthPaths(t, paths)
		defer restorePaths()
		var stderr bytes.Buffer

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
			"auth",
			"login",
		})

		require.Error(t, err)
		require.Contains(t, stderr.String(), "missing OAuth redirect URI")
	})
}

func Test_AuthLogin_reports_pkce_generation_error(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	originalPKCE := generateAuthLoginPKCE
	generateAuthLoginPKCE = func() (oauth.PKCE, error) {
		return oauth.PKCE{}, errors.New("pkce failed")
	}
	defer func() {
		generateAuthLoginPKCE = originalPKCE
	}()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"login",
	})

	require.ErrorContains(t, err, "pkce failed")
}

func Test_AuthLogin_state_mismatch_refuses_exchange_and_save(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"oauth-access-token",
		"oauth-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read", "write"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
	defer restore()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
		"auth",
		"login",
		"--callback", "http://127.0.0.1:8484/callback?code=code-123&state=wrong-state",
	})

	require.Error(t, err)
	require.Contains(t, stderr.String(), `"error_code":"AUTH_REAUTH_REQUIRED"`)
	require.Equal(t, 0, fakeOAuth.authorizationCodeCalls)
	got, loadErr := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, loadErr)
	require.Empty(t, got.Token)
}

func Test_AuthLogin_callback_url_exchanges_and_saves_after_readiness(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"oauth-access-token",
		"oauth-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read", "write"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{
		report: readyAuthReport("app"),
		beforeReturn: func() {
			got, loadErr := auth.NewStore(paths).Load(context.Background())
			require.NoError(t, loadErr)
			require.Empty(t, got.Token)
		},
	}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, fakeReadiness)
	defer restore()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, nil, &stdout, &stderr, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"login",
		"--callback", "http://127.0.0.1:8484/callback?code=code-123&state=state-123",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.authorizationCodeCalls)
	require.Equal(t, oauth.AuthorizationCodeRequest{
		Code:         "code-123",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CodeVerifier: "verifier-123",
	}, fakeOAuth.authorizationCodeRequest)
	require.Equal(t, "oauth-access-token", fakeReadiness.accessToken)
	require.Equal(t, "app", fakeReadiness.expectedActor)
	require.Equal(t, []string{"read", "write"}, fakeReadiness.requiredScopes)
	require.NotContains(t, stdout.String(), "oauth-access-token")
	require.NotContains(t, stdout.String(), "oauth-refresh-token")
	require.NotContains(t, stdout.String(), "client-secret")
	require.NotContains(t, stderr.String(), "oauth-access-token")
	require.NotContains(t, stderr.String(), "client-secret")

	got, loadErr := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "oauth-access-token", got.Token.AccessToken)
	require.Equal(t, "oauth-refresh-token", got.Token.RefreshToken)
	require.Equal(t, "app", got.Token.Actor)
	require.Equal(t, "authorization_code", got.Token.GrantType)
}

func Test_AuthLogin_raw_code_fallback_uses_authorization_code_exchange(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"oauth-access-token",
		"oauth-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read", "write"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
	defer restore()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"login",
		"--callback", "manual-code-123",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.authorizationCodeCalls)
	require.Equal(t, "manual-code-123", fakeOAuth.authorizationCodeRequest.Code)
	got, loadErr := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "oauth-access-token", got.Token.AccessToken)
	require.Equal(t, "authorization_code", got.Token.GrantType)
}

func Test_AuthLogin_reads_manual_callback_from_stdin_with_same_state(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"oauth-access-token",
		"oauth-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read", "write"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
	defer restore()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	stdin := strings.NewReader("http://127.0.0.1:8484/callback?code=stdin-code&state=state-123\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, stdin, &stdout, &stderr, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"login",
		"--callback", "-",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.authorizationCodeCalls)
	require.Equal(t, "stdin-code", fakeOAuth.authorizationCodeRequest.Code)
	require.Equal(t, "verifier-123", fakeOAuth.authorizationCodeRequest.CodeVerifier)
	require.Contains(t, stderr.String(), "state-123")
	require.NotContains(t, stdout.String(), "oauth-access-token")
	require.NotContains(t, stdout.String(), "oauth-refresh-token")
	require.NotContains(t, stderr.String(), "client-secret")
	got, loadErr := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, loadErr)
	require.Equal(t, "oauth-access-token", got.Token.AccessToken)
}

func Test_AuthLogin_reports_callback_read_errors(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()

	err := execute(
		context.Background(),
		BuildInfo{},
		failingReader{err: errors.New("stdin closed")},
		&bytes.Buffer{},
		&bytes.Buffer{},
		[]string{"auth", "login", "--callback", "-"},
	)

	require.ErrorContains(t, err, "read oauth callback")
}

func Test_AuthLogin_reports_empty_callback_from_stdin(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()
	var stderr bytes.Buffer

	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(" \n"),
		&bytes.Buffer{},
		&stderr,
		[]string{"auth", "login", "--callback", "-"},
	)

	require.Error(t, err)
	require.Contains(t, stderr.String(), "missing OAuth callback code")
}

func Test_AuthLogin_reports_authorization_url_write_error_before_reading_stdin(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()

	err := execute(
		context.Background(),
		BuildInfo{},
		failingReader{err: errors.New("stdin must not be read")},
		&bytes.Buffer{},
		failingWriter{err: errors.New("stderr closed")},
		[]string{"auth", "login", "--callback", "-"},
	)

	require.ErrorContains(t, err, "stderr closed")
}

func Test_AuthLogin_callback_failures_save_no_token(t *testing.T) {
	tests := []struct {
		name          string
		oauth         *fakeOAuthTokenClient
		readiness     *fakeAuthReadinessChecker
		wantError     string
		wantErrorCode string
	}{
		{
			name:      "token exchange fails",
			oauth:     &fakeOAuthTokenClient{err: errors.New("exchange failed")},
			readiness: &fakeAuthReadinessChecker{report: readyAuthReport("app")},
			wantError: "exchange failed",
		},
		{
			name: "granted scopes miss a required scope",
			oauth: &fakeOAuthTokenClient{grant: auth.NewTokenState(
				"oauth-access-token",
				"oauth-refresh-token",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{"read"},
			)},
			readiness:     &fakeAuthReadinessChecker{report: readyAuthReport("app")},
			wantErrorCode: "MISSING_SCOPE",
		},
		{
			name: "authorization does not match the pinned target",
			oauth: &fakeOAuthTokenClient{grant: auth.NewTokenState(
				"oauth-access-token",
				"oauth-refresh-token",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{"read", "write"},
			)},
			readiness:     &fakeAuthReadinessChecker{err: client.ErrTargetMismatch},
			wantErrorCode: "AUTH_TARGET_MISMATCH",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paths := cliAuthTestPaths(t)
			saveAuthLoginApp(t, paths)
			restore := useAuthCommandHooks(t, paths, test.oauth, test.readiness)
			defer restore()
			restoreLogin := useAuthLoginGeneratedValues(t)
			defer restoreLogin()

			err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
				"--org", "org-id",
				"--team", "LIT",
				"--team-id", "team-id",
				"auth",
				"login",
				"--callback", "http://127.0.0.1:8484/callback?code=code-123&state=state-123",
			})

			require.Error(t, err)
			if test.wantError != "" {
				require.ErrorContains(t, err, test.wantError)
			}
			if test.wantErrorCode != "" {
				require.Equal(t, test.wantErrorCode, errorCode(err))
			}
			got, loadErr := auth.NewStore(paths).Load(context.Background())
			require.NoError(t, loadErr)
			require.Empty(t, got.Token)
		})
	}
}

func Test_AuthLogin_reports_token_state_save_error_after_callback(t *testing.T) {
	root := t.TempDir()
	paths := auth.Paths{
		AppConfigPath: filepath.Join(root, "auth-app.json"),
		TokenPath:     filepath.Join(root, "auth-token.json"),
	}
	saveAuthLoginApp(t, paths)
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"oauth-access-token",
		"oauth-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read", "write"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{
		report: readyAuthReport("app"),
		beforeReturn: func() {
			require.NoError(t, os.Mkdir(paths.TokenPath, 0o700))
		},
	})
	defer restore()
	restoreLogin := useAuthLoginGeneratedValues(t)
	defer restoreLogin()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"auth",
		"login",
		"--callback", "http://127.0.0.1:8484/callback?code=code-123&state=state-123",
	})

	require.ErrorContains(t, err, "read auth token state")
}

func Test_AuthLogin_refuses_callbacks_without_a_code(t *testing.T) {
	tests := []struct {
		name     string
		callback string
	}{
		{name: "url missing code", callback: "http://127.0.0.1:8484/callback?state=state-123"},
		{name: "url-shaped without query", callback: "http://127.0.0.1:8484/callback"},
		{name: "query-shaped without code", callback: "?state=state-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := cliAuthTestPaths(t)
			saveAuthLoginApp(t, paths)
			fakeOAuth := &fakeOAuthTokenClient{}
			restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
			defer restore()
			restoreLogin := useAuthLoginGeneratedValues(t)
			defer restoreLogin()
			var stderr bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &stderr, []string{
				"auth",
				"login",
				"--callback", tt.callback,
			})

			require.Error(t, err)
			require.Equal(t, "AUTH_REAUTH_REQUIRED", errorCode(err))
			require.Contains(t, stderr.String(), "OAuth callback is missing code")
			require.Zero(t, fakeOAuth.authorizationCodeCalls)
		})
	}
}

func saveAuthLoginApp(t *testing.T, paths auth.Paths) {
	t.Helper()
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		Scopes:       []string{"read", "write"},
	}))
}

func useAuthLoginGeneratedValues(t *testing.T) func() {
	t.Helper()
	originalPKCE := generateAuthLoginPKCE
	originalState := generateAuthLoginState
	generateAuthLoginPKCE = func() (oauth.PKCE, error) {
		return oauth.PKCE{
			CodeVerifier:        "verifier-123",
			CodeChallenge:       "challenge-123",
			CodeChallengeMethod: "S256",
		}, nil
	}
	generateAuthLoginState = func() string {
		return "state-123"
	}

	return func() {
		generateAuthLoginState = originalState
		generateAuthLoginPKCE = originalPKCE
	}
}

func decodeAuthLoginStartReport(t *testing.T, data []byte) authLoginStartReport {
	t.Helper()
	var report authLoginStartReport
	require.NoError(t, json.Unmarshal(data, &report))

	return report
}

func authLoginAuthorizeURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	parsed, err := url.Parse(rawURL)
	require.NoError(t, err)

	return parsed
}

type failingReader struct {
	err error
}

func (reader failingReader) Read(_ []byte) (int, error) {
	return 0, reader.err
}

type failingWriter struct {
	err error
}

func (writer failingWriter) Write(_ []byte) (int, error) {
	return 0, writer.err
}

type cobraCommandWithIO struct {
	err io.Writer
}

func (commandWithIO cobraCommandWithIO) command() *cobra.Command {
	command := &cobra.Command{}
	if commandWithIO.err != nil {
		command.SetErr(commandWithIO.err)
	}

	return command
}
