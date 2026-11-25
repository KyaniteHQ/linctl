package cli

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/oauth"
	"github.com/KyaniteHQ/linctl/internal/render"
)

const (
	linearAuthorizeEndpoint = "https://linear.app/oauth/authorize"
	userActor               = "user"
)

const authLoginPrompt = "Open this URL in a browser, approve linctl, " +
	"then paste the redirect URL or the authorization code here and press Enter:"

func addAuthLoginCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var actor string
	command := &cobra.Command{
		Use:   "login",
		Short: "Authorize with OAuth browser login",
		Long: "Print a Linear authorization URL on stderr, then read the redirect URL or the " +
			"authorization code as one line from stdin and save the token after a live readiness check.",
		Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runAuthLogin(ctx, command, options, actor)
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	command.Flags().StringVar(&actor, "actor", appActor, "OAuth actor: app or user")
	root.AddCommand(command)
}

func runAuthLogin(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	requestedActor string,
) error {
	authContext, err := loadAuthCommandContext(ctx, command, options)
	if err != nil {
		return err
	}
	app, err := authLoginApp(authContext)
	if err != nil {
		return err
	}
	actor, err := normalizedAuthLoginActor(requestedActor)
	if err != nil {
		return err
	}
	pkce := oauth.GeneratePKCE()
	state := rand.Text()
	scopes := requiredScopes(app, actor)
	authorizeURL := buildAuthLoginAuthorizeURL(app, actor, scopes, state, pkce)
	if err := render.WriteLine(command.ErrOrStderr(), "%s\n%s", authLoginPrompt, authorizeURL); err != nil {
		return err
	}
	callback, err := readAuthLoginCallback(ctx, command.InOrStdin())
	if err != nil {
		return err
	}

	token, readiness, err := completeAuthLogin(ctx, authContext, authLoginCompletionRequest{
		App:           app,
		Actor:         actor,
		Scopes:        scopes,
		Callback:      callback,
		ExpectedState: state,
		PKCE:          pkce,
		Timeout:       options.timeout,
	})
	if err != nil {
		return err
	}
	if err := authContext.store.SaveTokenState(ctx, authContext.profile, token); err != nil {
		return err
	}

	return writeAuthStatus(command, options, newAuthStatusReport(app, token, readiness))
}

func authLoginApp(authContext authCommandContext) (auth.AppConfig, error) {
	app := authContext.app
	if strings.TrimSpace(app.ClientID) == "" {
		return auth.AppConfig{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing OAuth client id: run linctl auth configure",
		)
	}
	if strings.TrimSpace(app.RedirectURI) == "" {
		return auth.AppConfig{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing OAuth redirect URI: run linctl auth configure",
		)
	}

	return app, nil
}

type authLoginCompletionRequest struct {
	App           auth.AppConfig
	Actor         string
	Scopes        []string
	Callback      string
	ExpectedState string
	PKCE          oauth.PKCE
	Timeout       time.Duration
}

func completeAuthLogin(
	ctx context.Context,
	authContext authCommandContext,
	request authLoginCompletionRequest,
) (auth.TokenState, authReadinessReport, error) {
	code, err := authorizationCodeFromCallback(request.Callback, request.ExpectedState)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	token, err := newAuthOAuthClient(request.Timeout).ExchangeAuthorizationCode(ctx, oauth.AuthorizationCodeRequest{
		Code:         code,
		RedirectURI:  request.App.RedirectURI,
		ClientID:     request.App.ClientID,
		ClientSecret: request.App.ClientSecret,
		CodeVerifier: request.PKCE.CodeVerifier,
	})
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	token, err = stampAndRequireScopes(token, request.Actor, authGrantAuthorizationCode, request.Scopes)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	readiness, err := requireAuthReadiness(ctx, authReadinessRequest{
		AccessToken:    token.AccessToken,
		TokenActor:     token.Actor,
		TokenScopes:    token.Scopes,
		ExpectedTarget: authContext.target,
		ExpectedActor:  request.Actor,
		RequiredScopes: request.Scopes,
		Timeout:        request.Timeout,
	})
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}

	return token, readiness, nil
}

func buildAuthLoginAuthorizeURL(
	app auth.AppConfig,
	actor string,
	scopes []string,
	state string,
	pkce oauth.PKCE,
) string {
	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", app.ClientID)
	values.Set("redirect_uri", app.RedirectURI)
	values.Set("scope", strings.Join(scopes, ","))
	values.Set("state", state)
	values.Set("code_challenge", pkce.CodeChallenge)
	values.Set("code_challenge_method", pkce.CodeChallengeMethod)
	values.Set("actor", actor)
	values.Set("prompt", "consent")

	return linearAuthorizeEndpoint + "?" + values.Encode()
}

func authorizationCodeFromCallback(callback string, expectedState string) (string, error) {
	callback = strings.TrimSpace(callback)
	parsed, err := url.Parse(callback)
	if err == nil && parsed.RawQuery != "" {
		query := parsed.Query()
		code := strings.TrimSpace(query.Get("code"))
		if code == "" {
			return "", auth.NewError(auth.ErrorCodeReauthRequired, "OAuth callback is missing code")
		}
		if state := query.Get("state"); state != expectedState {
			return "", auth.NewError(auth.ErrorCodeReauthRequired, "OAuth callback state mismatch")
		}

		return code, nil
	}
	if strings.Contains(callback, "://") || strings.HasPrefix(callback, "?") {
		return "", auth.NewError(auth.ErrorCodeReauthRequired, "OAuth callback is missing code")
	}

	return callback, nil
}

func normalizedAuthLoginActor(actor string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(actor)) {
	case "", appActor:
		return appActor, nil
	case userActor:
		return userActor, nil
	default:
		return "", auth.NewError(auth.ErrorCodeActorMismatch, "OAuth actor must be app or user")
	}
}

func readAuthLoginCallback(ctx context.Context, stdin io.Reader) (string, error) {
	type lineResult struct {
		line string
		err  error
	}
	read := make(chan lineResult, 1)
	go func() {
		line, err := readLineWithoutReadAhead(stdin)
		read <- lineResult{line: line, err: err}
	}()

	var result lineResult
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("read oauth callback: %w", ctx.Err())
	case result = <-read:
	}
	if result.err != nil {
		return "", fmt.Errorf("read oauth callback: %w", result.err)
	}
	callback := strings.TrimSpace(result.line)
	if callback == "" {
		return "", auth.NewError(auth.ErrorCodeReauthRequired, "missing OAuth callback code")
	}

	return callback, nil
}

func readLineWithoutReadAhead(reader io.Reader) (string, error) {
	var line []byte
	var next [1]byte
	for {
		count, err := reader.Read(next[:])
		if count == 1 {
			if next[0] == '\n' {
				return string(line), nil
			}
			line = append(line, next[0])
		}
		if errors.Is(err, io.EOF) {
			return string(line), nil
		}
		if err != nil {
			return "", err
		}
	}
}
