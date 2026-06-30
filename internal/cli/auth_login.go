package cli

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"slices"
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

var (
	generateAuthLoginPKCE  = oauth.GeneratePKCE
	generateAuthLoginState = generateOAuthState
	authLoginRandomReader  = rand.Reader
)

type authLoginFlags struct {
	actor    string
	callback string
}

type authLoginStartReport struct {
	AuthorizeURL string   `json:"authorize_url"`
	Actor        string   `json:"actor"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
}

func addAuthLoginCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var flags authLoginFlags
	command := &cobra.Command{
		Use:   "login",
		Short: "Authorize with OAuth browser login",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runAuthLogin(ctx, command, options, flags)
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	command.Flags().StringVar(&flags.actor, "actor", appActor, "OAuth actor: app or user")
	command.Flags().StringVar(
		&flags.callback,
		"callback",
		"",
		"OAuth callback URL, authorization code, or '-' to read from stdin",
	)
	root.AddCommand(command)
}

func runAuthLogin(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	flags authLoginFlags,
) error {
	authContext, err := loadAuthCommandContext(ctx, command, options)
	if err != nil {
		return err
	}
	app, err := authLoginApp(authContext)
	if err != nil {
		return err
	}
	actor, err := normalizedAuthLoginActor(flags.actor)
	if err != nil {
		return err
	}
	pkce, state, err := newAuthLoginProof()
	if err != nil {
		return err
	}
	scopes := requiredScopes(app)
	report := newAuthLoginStartReport(app, actor, scopes, state, pkce)
	callback, err := authLoginCallback(command, options, strings.TrimSpace(flags.callback), report)
	if err != nil {
		return err
	}
	if callback == "" {
		return writeAuthLoginStart(command, options, report)
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

func newAuthLoginProof() (oauth.PKCE, string, error) {
	pkce, err := generateAuthLoginPKCE()
	if err != nil {
		return oauth.PKCE{}, "", err
	}
	state, err := generateAuthLoginState()
	if err != nil {
		return oauth.PKCE{}, "", err
	}

	return pkce, state, nil
}

func newAuthLoginStartReport(
	app auth.AppConfig,
	actor string,
	scopes []string,
	state string,
	pkce oauth.PKCE,
) authLoginStartReport {
	return authLoginStartReport{
		AuthorizeURL: buildAuthLoginAuthorizeURL(app, actor, scopes, state, pkce),
		Actor:        actor,
		RedirectURI:  app.RedirectURI,
		Scopes:       slices.Clone(scopes),
	}
}

func authLoginCallback(
	command *cobra.Command,
	options *rootOptions,
	callback string,
	report authLoginStartReport,
) (string, error) {
	if callback != "-" {
		return callback, nil
	}
	if !options.quiet {
		if err := render.WriteLine(command.ErrOrStderr(), "%s", report.AuthorizeURL); err != nil {
			return "", err
		}
	}

	return readAuthLoginCallback(command)
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
	token.Actor = request.Actor
	token.GrantType = authGrantAuthorizationCode
	if err := requireScopes(token.Scopes, request.Scopes); err != nil {
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
	if callback == "" {
		return "", auth.NewError(auth.ErrorCodeReauthRequired, "missing OAuth callback code")
	}
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

func generateOAuthState() (string, error) {
	random := make([]byte, 32)
	if _, err := io.ReadFull(authLoginRandomReader, random); err != nil {
		return "", fmt.Errorf("generate oauth state: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(random), nil
}

func writeAuthLoginStart(command *cobra.Command, options *rootOptions, report authLoginStartReport) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, report)
	}

	return render.WriteLine(command.OutOrStdout(), "%s", report.AuthorizeURL)
}

func readAuthLoginCallback(command *cobra.Command) (string, error) {
	data, err := io.ReadAll(command.InOrStdin())
	if err != nil {
		return "", fmt.Errorf("read oauth callback: %w", err)
	}
	callback := strings.TrimSpace(string(data))
	if callback == "" {
		return "", auth.NewError(auth.ErrorCodeReauthRequired, "missing OAuth callback code")
	}

	return callback, nil
}
