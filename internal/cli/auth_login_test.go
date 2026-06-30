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

func Test_GenerateOAuthState_returns_url_safe_state(t *testing.T) {
	state, err := generateOAuthState()

	require.NoError(t, err)
	require.NotEmpty(t, state)
	require.NotContains(t, state, "=")
}

func Test_GenerateOAuthState_reports_random_reader_error(t *testing.T) {
	original := authLoginRandomReader
	authLoginRandomReader = failingReader{err: errors.New("entropy unavailable")}
	defer func() {
		authLoginRandomReader = original
	}()

	_, err := generateOAuthState()

	require.Error(t, err)
	require.Contains(t, err.Error(), "generate oauth state")
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

func Test_AuthLogin_reports_generated_proof_errors(t *testing.T) {
	t.Run("pkce", func(t *testing.T) {
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

		require.Error(t, err)
		require.Contains(t, err.Error(), "pkce failed")
	})

	t.Run("state", func(t *testing.T) {
		paths := cliAuthTestPaths(t)
		saveAuthLoginApp(t, paths)
		restorePaths := useAuthPaths(t, paths)
		defer restorePaths()
		originalPKCE := generateAuthLoginPKCE
		originalState := generateAuthLoginState
		generateAuthLoginPKCE = func() (oauth.PKCE, error) {
			return oauth.PKCE{CodeVerifier: "verifier", CodeChallenge: "challenge", CodeChallengeMethod: "S256"}, nil
		}
		generateAuthLoginState = func() (string, error) {
			return "", errors.New("state failed")
		}
		defer func() {
			generateAuthLoginState = originalState
			generateAuthLoginPKCE = originalPKCE
		}()

		err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
			"auth",
			"login",
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "state failed")
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

	require.Error(t, err)
	require.Contains(t, err.Error(), "read oauth callback")
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

func Test_AuthLogin_callback_writer_error_is_reported(t *testing.T) {
	command := (cobraCommandWithIO{err: failingWriter{err: errors.New("stderr closed")}}).command()

	_, err := authLoginCallback(command, &rootOptions{}, "-", authLoginStartReport{
		AuthorizeURL: "https://linear.app/oauth/authorize",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "stderr closed")
}

func Test_AuthLogin_reports_completion_errors(t *testing.T) {
	baseRequest := authLoginCompletionRequest{
		App: auth.AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			RedirectURI:  "http://127.0.0.1:8484/callback",
		},
		Actor:         "app",
		Scopes:        []string{"read"},
		Callback:      "http://127.0.0.1:8484/callback?code=code-123&state=state-123",
		ExpectedState: "state-123",
		PKCE:          oauth.PKCE{CodeVerifier: "verifier-123"},
	}

	t.Run("exchange error", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{err: errors.New("exchange failed")},
			&fakeAuthReadinessChecker{},
		)
		defer restore()

		_, _, err := completeAuthLogin(context.Background(), authCommandContext{}, baseRequest)

		require.Error(t, err)
		require.Contains(t, err.Error(), "exchange failed")
	})

	t.Run("missing scope", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{grant: auth.NewTokenState(
				"oauth-access-token",
				"oauth-refresh-token",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{},
			)},
			&fakeAuthReadinessChecker{},
		)
		defer restore()

		_, _, err := completeAuthLogin(context.Background(), authCommandContext{}, baseRequest)

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeMissingScope), errorCode(err))
	})

	t.Run("readiness error", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{grant: auth.NewTokenState(
				"oauth-access-token",
				"oauth-refresh-token",
				"Bearer",
				time.Now().Add(time.Hour),
				[]string{"read"},
			)},
			&fakeAuthReadinessChecker{err: client.ErrTargetMismatch},
		)
		defer restore()

		_, _, err := completeAuthLogin(context.Background(), authCommandContext{}, baseRequest)

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeTargetMismatch), errorCode(err))
	})
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

	require.Error(t, err)
	require.Contains(t, err.Error(), "read auth token state")
}

func Test_AuthorizationCodeFromCallback_reports_missing_code_shapes(t *testing.T) {
	tests := []struct {
		name     string
		callback string
		want     string
	}{
		{
			name: "empty",
			want: "missing OAuth callback code",
		},
		{
			name:     "url missing code",
			callback: "http://127.0.0.1:8484/callback?state=state-123",
			want:     "OAuth callback is missing code",
		},
		{
			name:     "url-shaped without query",
			callback: "http://127.0.0.1:8484/callback",
			want:     "OAuth callback is missing code",
		},
		{
			name:     "query-shaped without code",
			callback: "?state=state-123",
			want:     "OAuth callback is missing code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := authorizationCodeFromCallback(tt.callback, "state-123")

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.want)
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
	generateAuthLoginState = func() (string, error) {
		return "state-123", nil
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
