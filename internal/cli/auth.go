package cli

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/oauth"
	"github.com/KyaniteHQ/linctl/internal/render"
)

const (
	appActor                   = "app"
	authGrantAuthorizationCode = "authorization_code"
	authGrantClientCredentials = "client_credentials"
)

var (
	defaultOAuthScopes            = []string{"read", "write", "issues:create", "comments:create"}
	authDefaultPaths              = auth.DefaultPaths
	authNow                       = time.Now
	checkAuthReadiness            = defaultCheckAuthReadiness
	newAuthReadinessGraphQLClient = func(accessToken string, timeout time.Duration) graphql.Client {
		return client.NewTransport(client.TransportConfig{
			Token:   client.OAuthAccessToken(accessToken),
			Timeout: timeout,
		})
	}
	newAuthOAuthClient = func() authOAuthClient {
		return oauth.NewClient(oauth.ClientConfig{})
	}
)

type authOAuthClient interface {
	ClientCredentials(context.Context, oauth.ClientCredentialsRequest) (auth.TokenGrant, error)
	ExchangeAuthorizationCode(context.Context, oauth.AuthorizationCodeRequest) (auth.TokenGrant, error)
	RefreshToken(context.Context, oauth.RefreshTokenRequest) (auth.TokenGrant, error)
	RevokeToken(context.Context, oauth.RevocationRequest) error
}

type authCommandContext struct {
	paths   auth.Paths
	store   auth.Store
	state   auth.State
	profile string
	target  config.Target
}

type authConfigureFlags struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scopes       []string
}

type authAppFlags struct {
	clientID     string
	clientSecret string
	scopes       []string
}

type authLogoutFlags struct {
	forgetApp bool
}

type authReadinessRequest struct {
	AccessToken    string
	ExpectedTarget config.Target
	ExpectedActor  string
	RequiredScopes []string
	Timeout        time.Duration
}

type authReadinessReport struct {
	Actor  string                `json:"actor"`
	Target client.ResolvedTarget `json:"target"`
}

type authConfigReport struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURI  string   `json:"redirect_uri,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

type authTokenReport struct {
	Status    string     `json:"status"`
	Type      string     `json:"type,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Scopes    []string   `json:"scopes,omitempty"`
}

type authTargetStatusReport struct {
	Status   string            `json:"status"`
	Expected map[string]string `json:"expected,omitempty"`
	Resolved map[string]string `json:"resolved,omitempty"`
}

type authStatusReport struct {
	App       authConfigReport       `json:"app"`
	Token     authTokenReport        `json:"token"`
	Actor     string                 `json:"actor,omitempty"`
	Scopes    []string               `json:"scopes,omitempty"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
	TokenType string                 `json:"token_type,omitempty"`
	Target    authTargetStatusReport `json:"target"`
}

type authLogoutReport struct {
	Token            string   `json:"token"`
	App              string   `json:"app"`
	Revoked          []string `json:"revoked,omitempty"`
	RevocationFailed bool     `json:"revocation_failed,omitempty"`
}

func addAuthCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	authCommand := &cobra.Command{
		Use:   "auth",
		Short: "Manage linctl OAuth authentication",
	}
	annotateCommand(authCommand, commandSafetyAnnotation, string(CommandSafetyLocal))
	addAuthConfigureCommand(ctx, authCommand, options)
	addAuthLoginCommand(ctx, authCommand, options)
	addAuthAppCommand(ctx, authCommand, options)
	addAuthStatusCommand(ctx, authCommand, options)
	addAuthRefreshCommand(ctx, authCommand, options)
	addAuthLogoutCommand(ctx, authCommand, options)
	root.AddCommand(authCommand)
}

func addAuthConfigureCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var flags authConfigureFlags
	command := &cobra.Command{
		Use:   "configure",
		Short: "Save OAuth app client configuration",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			paths, err := authDefaultPaths(nil)
			if err != nil {
				return err
			}
			if strings.TrimSpace(flags.clientID) == "" {
				return auth.NewError(auth.ErrorCodeNotConfigured, "missing --client-id")
			}
			app := auth.AppConfig{
				ClientID:     strings.TrimSpace(flags.clientID),
				ClientSecret: flags.clientSecret,
				RedirectURI:  strings.TrimSpace(flags.redirectURI),
				Scopes:       normalizedScopes(flags.scopes),
			}
			if err := auth.NewStore(paths).SaveAppConfig(ctx, options.profile, app); err != nil {
				return err
			}
			if options.quiet {
				return nil
			}
			report := redactedAppConfigReport(app)
			if options.json {
				return writeJSONValue(command, options, report)
			}

			return render.WriteLine(command.OutOrStdout(), "OAuth app configured")
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	command.Flags().StringVar(&flags.clientID, "client-id", "", "OAuth app client id")
	command.Flags().StringVar(&flags.clientSecret, "client-secret", "", "OAuth app client secret")
	command.Flags().StringVar(&flags.redirectURI, "redirect-uri", "", "OAuth redirect URI")
	command.Flags().StringSliceVar(&flags.scopes, "scopes", nil, "OAuth scopes")
	root.AddCommand(command)
}

func addAuthAppCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var flags authAppFlags
	command := &cobra.Command{
		Use:   "app",
		Short: "Authorize with OAuth client credentials as the app actor",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			authContext, err := loadAuthCommandContext(ctx, options)
			if err != nil {
				return err
			}
			profile := authContext.state.Profile(authContext.profile)
			app := mergeAppConfig(profile.App, auth.AppConfig{
				ClientID:     strings.TrimSpace(flags.clientID),
				ClientSecret: flags.clientSecret,
				Scopes:       normalizedScopes(flags.scopes),
			})
			if app.ClientID == "" {
				return auth.NewError(auth.ErrorCodeNotConfigured, "missing OAuth client id: run linctl auth configure")
			}
			if app.ClientSecret == "" {
				return auth.NewError(
					auth.ErrorCodeNotConfigured,
					"missing OAuth client secret: run linctl auth configure",
				)
			}

			token, readiness, err := acquireClientCredentialsToken(ctx, authContext, app, options.timeout)
			if err != nil {
				return err
			}
			if err := authContext.store.SaveTokenState(ctx, authContext.profile, token); err != nil {
				return err
			}
			if options.quiet {
				return nil
			}
			status := newAuthStatusReport(app, token, readiness)
			if options.json {
				return writeJSONValue(command, options, status)
			}

			return writeAuthStatusHuman(command, status)
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	command.Flags().StringVar(&flags.clientID, "client-id", "", "OAuth app client id")
	command.Flags().StringVar(&flags.clientSecret, "client-secret", "", "OAuth app client secret")
	command.Flags().StringSliceVar(&flags.scopes, "scopes", nil, "OAuth scopes")
	root.AddCommand(command)
}

func addAuthStatusCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "status",
		Short: "Check OAuth token, actor, scopes, and target readiness",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			authContext, err := loadAuthCommandContext(ctx, options)
			if err != nil {
				return err
			}
			profile := authContext.state.Profile(authContext.profile)
			app := profile.App
			token := profile.Token
			if token.AccessToken == "" || tokenExpired(token, authNow()) {
				return writeCurrentOrRefreshedAuthStatus(ctx, command, options, authContext, app, token)
			}

			readiness, err := requireAuthReadiness(ctx, authReadinessRequest{
				AccessToken:    token.AccessToken,
				ExpectedTarget: authContext.target,
				ExpectedActor:  firstNonEmptyString(token.Actor, appActor),
				RequiredScopes: requiredScopes(app),
				Timeout:        options.timeout,
			})
			if err != nil {
				return err
			}

			return writeAuthStatus(command, options, newAuthStatusReport(app, token, readiness))
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	root.AddCommand(command)
}

func addAuthRefreshCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh OAuth token state",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			authContext, err := loadAuthCommandContext(ctx, options)
			if err != nil {
				return err
			}
			profile := authContext.state.Profile(authContext.profile)
			app := profile.App
			token, readiness, err := refreshAuthTokenState(ctx, authContext, app, profile.Token, options.timeout)
			if err != nil {
				return err
			}
			if err := authContext.store.SaveTokenState(ctx, authContext.profile, token); err != nil {
				return err
			}

			return writeAuthStatus(command, options, newAuthStatusReport(app, token, readiness))
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	root.AddCommand(command)
}

func addAuthLogoutCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var flags authLogoutFlags
	command := &cobra.Command{
		Use:   "logout",
		Short: "Revoke OAuth tokens and remove local token state",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			authContext, err := loadAuthCommandContext(ctx, options)
			if err != nil {
				return err
			}
			profile := authContext.state.Profile(authContext.profile)
			revoked, revocationFailed := revokeTokenState(ctx, newAuthOAuthClient(), profile.Token)
			if err := authContext.store.ClearTokenState(ctx, authContext.profile); err != nil {
				return err
			}
			appStatus := "kept"
			if flags.forgetApp {
				if err := authContext.store.ClearAppConfig(ctx, authContext.profile); err != nil {
					return err
				}
				appStatus = "forgotten"
			}
			report := authLogoutReport{
				Token:            "removed",
				App:              appStatus,
				Revoked:          revoked,
				RevocationFailed: revocationFailed,
			}
			if options.quiet {
				return nil
			}
			if options.json {
				return writeJSONValue(command, options, report)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"auth logout token %s app %s revoked %s revocation_failed %t",
				report.Token,
				report.App,
				strings.Join(report.Revoked, ","),
				report.RevocationFailed,
			)
		},
	}
	annotateCommand(command, commandSafetyAnnotation, string(CommandSafetyLocal))
	command.Flags().BoolVar(&flags.forgetApp, "forget-app", false, "also remove saved OAuth app configuration")
	root.AddCommand(command)
}

func loadAuthCommandContext(ctx context.Context, options *rootOptions) (authCommandContext, error) {
	paths, err := authDefaultPaths(nil)
	if err != nil {
		return authCommandContext{}, err
	}
	resolvedConfig, err := config.Load(ctx, config.LoadRequest{
		GlobalPath:      defaultGlobalConfigPath(),
		RepoPath:        ".linctl.toml",
		AuthStatePaths:  paths,
		ProfileOverride: options.profile,
		TargetOverride:  targetOverride(options),
	})
	if err != nil {
		return authCommandContext{}, err
	}
	applyTargetOverrideFlagSemantics(&resolvedConfig, options)

	store := auth.NewStore(paths)
	state, err := store.Load(ctx)
	if err != nil {
		return authCommandContext{}, err
	}
	profileState := state.Profile(resolvedConfig.Profile)
	mergeResolvedAuthAppConfig(&state, resolvedConfig.Profile, profileState, resolvedConfig.Auth.App)

	return authCommandContext{
		paths:   paths,
		store:   store,
		state:   state,
		profile: resolvedConfig.Profile,
		target:  resolvedConfig.Target,
	}, nil
}

func writeCurrentOrRefreshedAuthStatus(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	authContext authCommandContext,
	app auth.AppConfig,
	token auth.TokenState,
) error {
	token, readiness, err := currentOrRefreshedAuthToken(ctx, authContext, app, token, options.timeout)
	if err != nil {
		return err
	}
	if err := authContext.store.SaveTokenState(ctx, authContext.profile, token); err != nil {
		return err
	}

	return writeAuthStatus(command, options, newAuthStatusReport(app, token, readiness))
}

func currentOrRefreshedAuthToken(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	token auth.TokenState,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	if token.AccessToken != "" || token.RefreshToken != "" {
		return refreshAuthTokenState(ctx, authContext, app, token, timeout)
	}
	if app.ClientID == "" || app.ClientSecret == "" {
		return auth.TokenState{}, authReadinessReport{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"run linctl auth configure and linctl auth app",
		)
	}

	return acquireClientCredentialsToken(ctx, authContext, app, timeout)
}

func mergeResolvedAuthAppConfig(
	state *auth.State,
	profileName string,
	profileState auth.ProfileState,
	app auth.AppConfig,
) {
	if authAppConfigEmpty(app) {
		return
	}
	profileState.App = mergeAppConfig(profileState.App, app)
	if profileName == "" {
		state.App = profileState.App
		return
	}
	if state.Profiles == nil {
		state.Profiles = map[string]auth.ProfileState{}
	}
	state.Profiles[profileName] = profileState
}

func authAppConfigEmpty(app auth.AppConfig) bool {
	return app.ClientID == "" &&
		app.ClientSecret == "" &&
		app.RedirectURI == "" &&
		len(app.Scopes) == 0
}

func acquireClientCredentialsToken(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	scopes := requiredScopes(app)
	token, err := exchangeClientCredentialsToken(ctx, newAuthOAuthClient(), app)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	readiness, err := requireAuthReadiness(ctx, authReadinessRequest{
		AccessToken:    token.AccessToken,
		ExpectedTarget: authContext.target,
		ExpectedActor:  appActor,
		RequiredScopes: scopes,
		Timeout:        timeout,
	})
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}

	return token, readiness, nil
}

func refreshAuthTokenState(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	token auth.TokenState,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	if token.AccessToken == "" && token.RefreshToken == "" {
		return auth.TokenState{}, authReadinessReport{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing OAuth token state: run linctl auth login or linctl auth app",
		)
	}
	if token.GrantType == authGrantClientCredentials {
		if app.ClientID == "" || app.ClientSecret == "" {
			return auth.TokenState{}, authReadinessReport{}, auth.NewError(
				auth.ErrorCodeNotConfigured,
				"missing OAuth app client credentials: run linctl auth configure",
			)
		}

		return acquireClientCredentialsToken(ctx, authContext, app, timeout)
	}
	if token.RefreshToken == "" {
		return auth.TokenState{}, authReadinessReport{}, auth.NewError(
			auth.ErrorCodeReauthRequired,
			"missing OAuth refresh token: run linctl auth login or linctl auth app",
		)
	}
	if app.ClientID == "" {
		return auth.TokenState{}, authReadinessReport{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing OAuth client id: run linctl auth configure",
		)
	}

	scopes := requiredScopes(app)
	refreshed, err := refreshAuthorizationCodeToken(ctx, newAuthOAuthClient(), app, token, scopes)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	readiness, err := requireAuthReadiness(ctx, authReadinessRequest{
		AccessToken:    refreshed.AccessToken,
		ExpectedTarget: authContext.target,
		ExpectedActor:  firstNonEmptyString(token.Actor, appActor),
		RequiredScopes: scopes,
		Timeout:        timeout,
	})
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}

	return refreshed, readiness, nil
}

func refreshAuthorizationCodeToken(
	ctx context.Context,
	oauthClient authOAuthClient,
	app auth.AppConfig,
	token auth.TokenState,
	scopes []string,
) (auth.TokenState, error) {
	grant, err := oauthClient.RefreshToken(ctx, oauth.RefreshTokenRequest{
		RefreshToken: token.RefreshToken,
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
	})
	if err != nil {
		return auth.TokenState{}, auth.WrapError(
			auth.ErrorCodeRefreshFailed,
			"refresh OAuth token: run linctl auth login",
			err,
		)
	}
	refreshed := grant.State
	refreshed.Actor = firstNonEmptyString(token.Actor, appActor)
	refreshed.GrantType = authGrantAuthorizationCode
	if err := requireScopes(refreshed.Scopes, scopes); err != nil {
		return auth.TokenState{}, err
	}

	return refreshed, nil
}

func revokeTokenState(
	ctx context.Context,
	oauthClient authOAuthClient,
	token auth.TokenState,
) ([]string, bool) {
	revoked := []string{}
	failed := false
	for _, request := range []oauth.RevocationRequest{
		{Token: token.RefreshToken, TokenTypeHint: "refresh_token"},
		{Token: token.AccessToken, TokenTypeHint: "access_token"},
	} {
		if request.Token == "" {
			continue
		}
		if err := oauthClient.RevokeToken(ctx, request); err != nil {
			failed = true
			continue
		}
		revoked = append(revoked, request.TokenTypeHint)
	}

	return revoked, failed
}

func exchangeClientCredentialsToken(
	ctx context.Context,
	oauthClient authOAuthClient,
	app auth.AppConfig,
) (auth.TokenState, error) {
	scopes := requiredScopes(app)
	grant, err := oauthClient.ClientCredentials(ctx, oauth.ClientCredentialsRequest{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Scopes:       scopes,
	})
	if err != nil {
		return auth.TokenState{}, err
	}
	token := grant.State
	token.Actor = appActor
	token.GrantType = authGrantClientCredentials
	if err := requireScopes(token.Scopes, scopes); err != nil {
		return auth.TokenState{}, err
	}

	return token, nil
}

func requireAuthReadiness(ctx context.Context, request authReadinessRequest) (authReadinessReport, error) {
	readiness, err := checkAuthReadiness(ctx, request)
	if err != nil {
		return authReadinessReport{}, mapAuthReadinessError(err)
	}
	if request.ExpectedActor != "" && readiness.Actor != "" && readiness.Actor != request.ExpectedActor {
		return authReadinessReport{}, auth.NewError(
			auth.ErrorCodeActorMismatch,
			fmt.Sprintf("expected actor %q but resolved %q", request.ExpectedActor, readiness.Actor),
		)
	}

	return readiness, nil
}

func defaultCheckAuthReadiness(ctx context.Context, request authReadinessRequest) (authReadinessReport, error) {
	if request.AccessToken == "" {
		return authReadinessReport{}, auth.NewError(auth.ErrorCodeNotConfigured, "missing OAuth access token")
	}
	graphqlClient := newAuthReadinessGraphQLClient(request.AccessToken, request.Timeout)
	target, err := client.ResolveTarget(ctx, graphqlClient, request.ExpectedTarget)
	if err != nil {
		return authReadinessReport{}, err
	}

	return authReadinessReport{Actor: request.ExpectedActor, Target: target}, nil
}

func mapAuthReadinessError(err error) error {
	var authErr *auth.AuthError
	var tokenErr *auth.TokenEndpointError
	switch {
	case errors.As(err, &authErr):
		return err
	case errors.As(err, &tokenErr):
		return err
	case errors.Is(err, client.ErrTargetMismatch), errors.Is(err, client.ErrTargetNotConfigured):
		return auth.WrapError(
			auth.ErrorCodeTargetMismatch,
			"OAuth authorization does not match the pinned target",
			err,
		)
	default:
		return err
	}
}

func requireScopes(actual []string, required []string) error {
	missing := missingScopes(actual, required)
	if len(missing) == 0 {
		return nil
	}

	return auth.NewError(
		auth.ErrorCodeMissingScope,
		"missing OAuth scopes: "+strings.Join(missing, ",")+
			"; run linctl auth configure --scopes "+strings.Join(required, ",")+
			" then linctl auth app or linctl auth login",
	)
}

func missingScopes(actual []string, required []string) []string {
	actualSet := map[string]bool{}
	for _, scope := range actual {
		actualSet[scope] = true
	}
	missing := []string{}
	for _, scope := range required {
		if !actualSet[scope] {
			missing = append(missing, scope)
		}
	}

	return missing
}

func mergeAppConfig(base auth.AppConfig, override auth.AppConfig) auth.AppConfig {
	merged := base
	if override.ClientID != "" {
		merged.ClientID = override.ClientID
	}
	if override.ClientSecret != "" {
		merged.ClientSecret = override.ClientSecret
	}
	if override.RedirectURI != "" {
		merged.RedirectURI = override.RedirectURI
	}
	if len(override.Scopes) > 0 {
		merged.Scopes = slices.Clone(override.Scopes)
	}

	return merged
}

func requiredScopes(app auth.AppConfig) []string {
	if len(app.Scopes) > 0 {
		return slices.Clone(app.Scopes)
	}

	return slices.Clone(defaultOAuthScopes)
}

func normalizedScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}
	normalized := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		for _, part := range strings.FieldsFunc(scope, func(r rune) bool {
			return r == ',' || r == ' ' || r == '\t' || r == '\n'
		}) {
			part = strings.TrimSpace(part)
			if part != "" && !slices.Contains(normalized, part) {
				normalized = append(normalized, part)
			}
		}
	}

	return normalized
}

func tokenExpired(token auth.TokenState, now time.Time) bool {
	return token.ExpiresAt != nil && !token.ExpiresAt.After(now)
}

func newAuthStatusReport(
	app auth.AppConfig,
	token auth.TokenState,
	readiness authReadinessReport,
) authStatusReport {
	return authStatusReport{
		App: redactedAppConfigReport(app),
		Token: authTokenReport{
			Status:    presence(token.AccessToken),
			Type:      token.TokenType,
			ExpiresAt: token.ExpiresAt,
			Scopes:    slices.Clone(token.Scopes),
		},
		Actor:     firstNonEmptyString(readiness.Actor, token.Actor),
		Scopes:    slices.Clone(token.Scopes),
		ExpiresAt: token.ExpiresAt,
		TokenType: token.TokenType,
		Target: authTargetStatusReport{
			Status:   "ready",
			Expected: targetMap(readiness.Target.Expected),
			Resolved: targetMap(readiness.Target.Resolved),
		},
	}
}

func redactedAppConfigReport(app auth.AppConfig) authConfigReport {
	return authConfigReport{
		ClientID:     presence(app.ClientID),
		ClientSecret: presence(app.ClientSecret),
		RedirectURI:  app.RedirectURI,
		Scopes:       slices.Clone(app.Scopes),
	}
}

func presence(value string) string {
	if value == "" {
		return "missing"
	}

	return "set"
}

func firstNonEmptyString(primary string, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}

func writeAuthStatus(command *cobra.Command, options *rootOptions, status authStatusReport) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, status)
	}

	return writeAuthStatusHuman(command, status)
}

func writeAuthStatusHuman(command *cobra.Command, status authStatusReport) error {
	return render.WriteLine(
		command.OutOrStdout(),
		"auth %s actor %s scopes %s target %s",
		status.Token.Status,
		defaultString(status.Actor, "unknown"),
		strings.Join(status.Scopes, ","),
		status.Target.Status,
	)
}
