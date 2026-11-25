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
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

func Test_AuthLogin_builds_authorization_url_with_app_actor_by_default(t *testing.T) {
	paths := cliAuthTestPaths(t)
	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "injected-access-token")
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()

	authorizeURL := startAuthLogin(t, "auth", "login")

	query := authorizeURL.Query()
	require.Equal(t, "https", authorizeURL.Scheme)
	require.Equal(t, "linear.app", authorizeURL.Host)
	require.Equal(t, "/oauth/authorize", authorizeURL.Path)
	require.Equal(t, "code", query.Get("response_type"))
	require.Equal(t, "client-id", query.Get("client_id"))
	require.Equal(t, "http://127.0.0.1:8484/callback", query.Get("redirect_uri"))
	require.Equal(t, "read,write", query.Get("scope"))
	require.NotEmpty(t, query.Get("state"))
	require.Regexp(t, `^[A-Za-z0-9_-]{43}$`, query.Get("code_challenge"))
	require.Equal(t, "S256", query.Get("code_challenge_method"))
	require.Equal(t, "app", query.Get("actor"))
	require.Equal(t, "consent", query.Get("prompt"))
}

func Test_AuthLogin_builds_authorization_url_with_user_actor(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()

	query := startAuthLogin(t, "auth", "login", "--actor", "user").Query()

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

			query := startAuthLogin(t, "auth", "login", "--actor", test.actor).Query()

			scopes := strings.Split(query.Get("scope"), ",")
			require.Contains(t, scopes, "read")
			require.Equal(t, test.wantsInitiativeScopes, slices.Contains(scopes, "initiative:read"))
			require.Equal(t, test.wantsInitiativeScopes, slices.Contains(scopes, "initiative:write"))
		})
	}
}

func Test_AuthLogin_prints_the_authorization_url_even_when_quiet(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()

	authorizeURL := startAuthLogin(t, "--quiet", "auth", "login")

	require.Equal(t, "linear.app", authorizeURL.Host)
}

func Test_AuthLogin_rejects_invalid_actor(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
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
		return startAuthLogin(t, "auth", "login").Query()
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
	var stderr bytes.Buffer
	stdin := strings.NewReader("http://127.0.0.1:8484/callback?code=code-123&state=wrong-state\n")

	err := execute(context.Background(), BuildInfo{}, stdin, &bytes.Buffer{}, &stderr, []string{"auth", "login"})

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
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, answerAuthLoginPrompt(t, &stderr, "code-123"), &stdout, &stderr, []string{
		"--json",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"login",
	})

	require.NoError(t, err)
	require.Equal(t, 1, fakeOAuth.authorizationCodeCalls)
	exchange := fakeOAuth.authorizationCodeRequest
	require.Equal(t,
		printedAuthorizeURL(t, stderr.String()).Query().Get("code_challenge"),
		oauth.CodeChallengeS256(exchange.CodeVerifier))
	exchange.CodeVerifier = ""
	require.Equal(t, oauth.AuthorizationCodeRequest{
		Code:         "code-123",
		RedirectURI:  "http://127.0.0.1:8484/callback",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}, exchange)
	require.Equal(t, "oauth-access-token", fakeReadiness.accessToken)
	require.Equal(t, "app", fakeReadiness.expectedActor)
	require.Equal(t, []string{"read", "write"}, fakeReadiness.requiredScopes)
	var status struct {
		Actor string `json:"actor"`
		Token struct {
			Status string `json:"status"`
		} `json:"token"`
		Target struct {
			Status string `json:"status"`
		} `json:"target"`
	}
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &status))
	require.Equal(t, "app", status.Actor)
	require.Equal(t, "set", status.Token.Status)
	require.Equal(t, "ready", status.Target.Status)
	promptLines := strings.Split(strings.TrimSuffix(stderr.String(), "\n"), "\n")
	require.Len(t, promptLines, 2)
	require.Equal(t,
		"Open this URL in a browser, approve linctl, then paste the redirect URL or the authorization code here and press Enter:",
		promptLines[0])
	require.True(t, strings.HasPrefix(promptLines[1], "https://linear.app/oauth/authorize?"))
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

func Test_AuthLogin_exchanges_a_raw_code_read_from_the_first_stdin_line(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantRemaining string
	}{
		{name: "leaves later input unread", input: "manual-code-123\nsecond line\n", wantRemaining: "second line\n"},
		{name: "accepts a last line without newline", input: "manual-code-123"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
			stdin := strings.NewReader(test.input)

			err := execute(context.Background(), BuildInfo{}, stdin, &bytes.Buffer{}, &bytes.Buffer{}, []string{
				"--org", "org-id",
				"--team", "LIT",
				"--team-id", "team-id",
				"auth",
				"login",
			})

			require.NoError(t, err)
			require.Equal(t, 1, fakeOAuth.authorizationCodeCalls)
			require.Equal(t, "manual-code-123", fakeOAuth.authorizationCodeRequest.Code)
			remaining, readErr := io.ReadAll(stdin)
			require.NoError(t, readErr)
			require.Equal(t, test.wantRemaining, string(remaining))
			got, loadErr := auth.NewStore(paths).Load(context.Background())
			require.NoError(t, loadErr)
			require.Equal(t, "oauth-access-token", got.Token.AccessToken)
			require.Equal(t, "authorization_code", got.Token.GrantType)
		})
	}
}

func Test_AuthLogin_stops_waiting_for_the_callback_when_cancelled(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	fakeOAuth := &fakeOAuthTokenClient{}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("app")})
	defer restore()
	stdin, stdinWriter := io.Pipe()
	defer func() { require.NoError(t, stdinWriter.Close()) }()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	prompted := &signalingWriter{written: make(chan struct{})}
	result := make(chan error, 1)

	go func() {
		result <- execute(ctx, BuildInfo{}, stdin, &bytes.Buffer{}, prompted, []string{"auth", "login"})
	}()
	<-prompted.written
	cancel()

	select {
	case err := <-result:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "auth login kept waiting for stdin after cancellation")
	}
	require.Zero(t, fakeOAuth.authorizationCodeCalls)
}

func Test_AuthLogin_reports_callback_read_errors(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()

	err := execute(
		context.Background(),
		BuildInfo{},
		failingReader{err: errors.New("stdin closed")},
		&bytes.Buffer{},
		&bytes.Buffer{},
		[]string{"auth", "login"},
	)

	require.ErrorContains(t, err, "read oauth callback")
}

func Test_AuthLogin_reports_empty_callback_from_stdin(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	var stderr bytes.Buffer

	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(" \n"),
		&bytes.Buffer{},
		&stderr,
		[]string{"auth", "login"},
	)

	require.Error(t, err)
	require.Contains(t, stderr.String(), "missing OAuth callback code")
}

func Test_AuthLogin_reports_authorization_url_write_error_before_reading_stdin(t *testing.T) {
	paths := cliAuthTestPaths(t)
	saveAuthLoginApp(t, paths)
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()

	err := execute(
		context.Background(),
		BuildInfo{},
		failingReader{err: errors.New("stdin must not be read")},
		&bytes.Buffer{},
		failingWriter{err: errors.New("stderr closed")},
		[]string{"auth", "login"},
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

			var stderr bytes.Buffer

			err := execute(context.Background(), BuildInfo{}, answerAuthLoginPrompt(t, &stderr, "code-123"), &bytes.Buffer{}, &stderr, []string{
				"--org", "org-id",
				"--team", "LIT",
				"--team-id", "team-id",
				"auth",
				"login",
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

	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, answerAuthLoginPrompt(t, &stderr, "code-123"), &bytes.Buffer{}, &stderr, []string{
		"auth",
		"login",
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
			var stderr bytes.Buffer

			stdin := strings.NewReader(tt.callback + "\n")

			err := execute(context.Background(), BuildInfo{}, stdin, &bytes.Buffer{}, &stderr, []string{"auth", "login"})

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

func startAuthLogin(t *testing.T, args ...string) *url.URL {
	t.Helper()
	var stderr bytes.Buffer

	err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &bytes.Buffer{}, &stderr, args)

	require.ErrorContains(t, err, "missing OAuth callback code")

	return printedAuthorizeURL(t, stderr.String())
}

func printedAuthorizeURL(t *testing.T, stderr string) *url.URL {
	t.Helper()
	for line := range strings.Lines(stderr) {
		if strings.HasPrefix(line, linearAuthorizeEndpoint) {
			parsed, err := url.Parse(strings.TrimSpace(line))
			require.NoError(t, err)

			return parsed
		}
	}
	require.FailNow(t, "auth login printed no authorization URL", stderr)

	return nil
}

func answerAuthLoginPrompt(t *testing.T, stderr *bytes.Buffer, code string) io.Reader {
	t.Helper()

	return &authLoginPromptAnswer{t: t, stderr: stderr, code: code}
}

type authLoginPromptAnswer struct {
	t      *testing.T
	stderr *bytes.Buffer
	code   string
	answer io.Reader
}

func (prompt *authLoginPromptAnswer) Read(data []byte) (int, error) {
	if prompt.answer == nil {
		state := printedAuthorizeURL(prompt.t, prompt.stderr.String()).Query().Get("state")
		callback := url.URL{
			Scheme:   "http",
			Host:     "127.0.0.1:8484",
			Path:     "/callback",
			RawQuery: url.Values{"code": {prompt.code}, "state": {state}}.Encode(),
		}
		prompt.answer = strings.NewReader(callback.String() + "\n")
	}

	return prompt.answer.Read(data)
}

type failingReader struct {
	err error
}

func (reader failingReader) Read(_ []byte) (int, error) {
	return 0, reader.err
}

type signalingWriter struct {
	once    sync.Once
	written chan struct{}
}

func (writer *signalingWriter) Write(data []byte) (int, error) {
	writer.once.Do(func() { close(writer.written) })

	return len(data), nil
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
