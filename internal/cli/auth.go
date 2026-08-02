package cli

import (
	"context"
	"log/slog"
	"os"
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
	authGrantAuthorizationCode = auth.GrantTypeAuthorizationCode
	authGrantClientCredentials = auth.GrantTypeClientCredentials
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
	newAuthOAuthClient = func(timeout time.Duration) authOAuthClient {
		return oauth.NewClient(oauth.ClientConfig{HTTPClient: newOAuthHTTPClient(timeout)})
	}
)

type authOAuthClient interface {
	ClientCredentials(context.Context, oauth.ClientCredentialsRequest) (auth.TokenState, error)
	ExchangeAuthorizationCode(context.Context, oauth.AuthorizationCodeRequest) (auth.TokenState, error)
	RefreshToken(context.Context, oauth.RefreshTokenRequest) (auth.TokenState, error)
	RevokeToken(context.Context, oauth.RevocationRequest) error
}

type authCommandContext struct {
	paths          auth.Paths
	store          auth.Store
	profile        string
	target         config.Target
	app            auth.AppConfig
	token          auth.TokenState
	localToken     auth.TokenState
	credentialKind auth.CredentialKind
	logger         *slog.Logger
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

func addAuthCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	authCommand := newGroupCommand("auth", "Manage linctl OAuth authentication")
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
			authContext, err := loadAuthProfileContext(ctx, command, options)
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
			if err := authContext.store.SaveAppConfig(ctx, authContext.profile, app); err != nil {
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
			authContext, err := loadAuthCommandContext(ctx, command, options)
			if err != nil {
				return err
			}
			app := auth.MergeAppConfig(authContext.app, auth.AppConfig{
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
			authContext, err := loadAuthCommandContext(ctx, command, options)
			if err != nil {
				return err
			}
			app := authContext.app
			token := authContext.token
			if token.AccessToken == "" && token.RefreshToken == "" {
				return writeCurrentOrRefreshedAuthStatus(ctx, command, options, authContext, app, token)
			}
			if tokenExpired(token, authNow()) {
				refreshed, readiness, err := currentOrRefreshedAuthToken(
					ctx,
					authContext,
					app,
					token,
					options.timeout,
				)
				if err != nil {
					return err
				}

				return writeAuthStatus(command, options, newAuthStatusReport(app, refreshed, readiness))
			}
			expectedActor := firstNonEmptyString(token.Actor, appActor)
			requiredTokenScopes := requiredScopes(app)
			if authContext.credentialKind == auth.CredentialKindInjectedAccessToken {
				expectedActor = ""
				requiredTokenScopes = nil
			}

			readiness, err := requireLoggedAuthReadiness(ctx, authContext.log(), authReadinessRequest{
				AccessToken:    token.AccessToken,
				TokenActor:     token.Actor,
				TokenScopes:    token.Scopes,
				ExpectedTarget: authContext.target,
				ExpectedActor:  expectedActor,
				RequiredScopes: requiredTokenScopes,
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
			authContext, err := loadAuthCommandContext(ctx, command, options)
			if err != nil {
				return err
			}
			app := authContext.app
			token, readiness, err := refreshPersistedAuthToken(ctx, authContext, app, options.timeout)
			if err != nil {
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
			authContext, err := loadAuthCommandContext(ctx, command, options)
			if err != nil {
				return err
			}
			var revoked []string
			var revocationFailed bool
			_, err = authContext.store.TransactTokenState(
				ctx,
				authContext.profile,
				func(current auth.TokenState) (auth.TokenState, error) {
					revoked, revocationFailed = revokeTokenState(
						ctx,
						authContext.log(),
						newAuthOAuthClient(options.timeout),
						current,
					)
					return auth.TokenState{}, nil
				},
			)
			if err != nil {
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
			authContext.log().DebugContext(
				ctx,
				"auth logout completed",
				"profile", authContext.profile,
				"app", report.App,
				"revoked_count", len(report.Revoked),
				"revocation_failed", report.RevocationFailed,
			)
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
	command.Flags().BoolVar(&flags.forgetApp, "forget-app", false, "also remove the saved OAuth app configuration")
	root.AddCommand(command)
}

func loadAuthCommandContext(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
) (authCommandContext, error) {
	authContext, err := loadAuthProfileContext(ctx, command, options)
	if err != nil {
		return authCommandContext{}, err
	}
	request := auth.SessionRequest{
		Store:   authContext.store,
		Profile: authContext.profile,
	}
	localSession, err := auth.SelectLocalSession(ctx, request)
	if err != nil {
		return authCommandContext{}, err
	}
	effectiveSession, err := auth.SelectSession(ctx, request)
	if err != nil {
		return authCommandContext{}, err
	}
	authContext.app = localSession.App
	authContext.token = effectiveSession.Token
	authContext.localToken = localSession.Token
	authContext.credentialKind = effectiveSession.CredentialKind

	return authContext, nil
}

func loadAuthProfileContext(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
) (authCommandContext, error) {
	paths, err := authDefaultPaths(nil)
	if err != nil {
		return authCommandContext{}, err
	}
	resolvedConfig, err := resolveConfig(ctx, options)
	if err != nil {
		return authCommandContext{}, err
	}

	store := auth.NewStore(paths)

	return authCommandContext{
		paths:   paths,
		store:   store,
		profile: resolvedConfig.Profile,
		target:  resolvedConfig.Target,
		logger:  authDiagnosticLogger(command, options),
	}, nil
}

func (authContext authCommandContext) log() *slog.Logger {
	if authContext.logger == nil {
		return discardLogger
	}

	return authContext.logger
}

func authDiagnosticLogger(command *cobra.Command, options *rootOptions) *slog.Logger {
	return newDiagnosticLogger(options.debug, os.Getenv("LINCTL_DEBUG_JSON") == "1", command.ErrOrStderr())
}

func firstNonEmptyString(primary string, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}
