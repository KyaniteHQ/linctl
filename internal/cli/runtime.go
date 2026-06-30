package cli

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

type commandRuntime struct {
	config        config.Resolved
	fileClient    httpDoer
	graphqlClient graphql.Client
	logger        *slog.Logger
}

var buildCommandRuntime = newCommandRuntime

func newCommandRuntime(ctx context.Context, options *rootOptions) (commandRuntime, error) {
	logger := newDiagnosticLogger(options.debug, os.Getenv("LINCTL_DEBUG_JSON") == "1", os.Stderr)
	override := targetOverride(options)
	authStatePaths, err := authDefaultPaths(nil)
	if err != nil {
		return commandRuntime{}, err
	}
	resolvedConfig, err := config.Load(ctx, config.LoadRequest{
		GlobalPath:      defaultGlobalConfigPath(),
		RepoPath:        ".linctl.toml",
		ProfileOverride: options.profile,
		TargetOverride:  override,
	})
	if err != nil {
		return commandRuntime{}, err
	}
	applyTargetOverrideFlagSemantics(&resolvedConfig, options)
	authStore := auth.NewStore(authStatePaths)
	authSession, err := auth.SelectSession(ctx, auth.SessionRequest{
		Store:   authStore,
		Profile: resolvedConfig.Profile,
	})
	if err != nil {
		return commandRuntime{}, err
	}
	if authSession.Token.AccessToken == "" {
		return commandRuntime{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing Linear OAuth access token: run linctl auth configure, then linctl auth app or linctl auth login",
		)
	}

	logger.Debug(
		"runtime ready",
		"profile", resolvedConfig.Profile,
		"org", resolvedConfig.Target.OrgID,
		"team_key", resolvedConfig.Target.TeamKey,
		"team_id", resolvedConfig.Target.TeamID,
		"project", resolvedConfig.Target.ProjectID,
		"timeout", options.timeout.String(),
	)

	return commandRuntime{
		config:     resolvedConfig,
		fileClient: &http.Client{Timeout: options.timeout},
		logger:     logger,
		graphqlClient: newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token:       authSession.Token,
			App:         authSession.App,
			Store:       authStore,
			Profile:     resolvedConfig.Profile,
			Timeout:     options.timeout,
			Persist:     authSession.PersistentToken,
			Logger:      logger,
			OAuthClient: newAuthOAuthClient(),
			NewClient: func(accessToken string) graphql.Client {
				return client.NewTransport(client.TransportConfig{
					Token:            client.OAuthAccessToken(accessToken),
					Timeout:          options.timeout,
					DiagnosticWriter: newTransportDiagnosticWriter(logger, options.debug),
				})
			},
		}),
	}, nil
}

type recoveringGraphQLClientConfig struct {
	Token       auth.TokenState
	App         auth.AppConfig
	Store       auth.Store
	Profile     string
	Timeout     time.Duration
	Persist     bool
	Logger      *slog.Logger
	OAuthClient authOAuthClient
	NewClient   func(accessToken string) graphql.Client
}

type recoveringGraphQLClient struct {
	token       auth.TokenState
	app         auth.AppConfig
	store       auth.Store
	profile     string
	timeout     time.Duration
	persist     bool
	logger      *slog.Logger
	oauthClient authOAuthClient
	newClient   func(accessToken string) graphql.Client
	client      graphql.Client
}

func newRecoveringGraphQLClient(config recoveringGraphQLClientConfig) *recoveringGraphQLClient {
	newClient := config.NewClient
	if newClient == nil {
		newClient = func(accessToken string) graphql.Client {
			return client.NewTransport(client.TransportConfig{
				Token:   client.OAuthAccessToken(accessToken),
				Timeout: config.Timeout,
			})
		}
	}
	oauthClient := config.OAuthClient
	if oauthClient == nil {
		oauthClient = newAuthOAuthClient()
	}
	logger := config.Logger
	if logger == nil {
		logger = discardLogger
	}
	recovering := &recoveringGraphQLClient{
		token:       config.Token,
		app:         config.App,
		store:       config.Store,
		profile:     config.Profile,
		timeout:     config.Timeout,
		persist:     config.Persist,
		logger:      logger,
		oauthClient: oauthClient,
		newClient:   newClient,
	}
	recovering.client = newClient(config.Token.AccessToken)

	return recovering
}

func (recovering *recoveringGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	recovered := false
	if tokenExpired(recovering.token, authNow()) {
		if err := recovering.recoverToken(ctx, "expired"); err != nil {
			return err
		}
		recovered = true
	}
	err := recovering.client.MakeRequest(ctx, request, response)
	if err == nil {
		return nil
	}
	if recovered {
		if errors.Is(err, client.ErrAuthFailed) {
			return auth.WrapError(auth.ErrorCodeReauthRequired, "OAuth token rejected after recovery", err)
		}

		return err
	}
	if !errors.Is(err, client.ErrAuthFailed) {
		return err
	}
	if err := recovering.recoverToken(ctx, "auth_failed"); err != nil {
		return err
	}

	err = recovering.client.MakeRequest(ctx, request, response)
	if errors.Is(err, client.ErrAuthFailed) {
		return auth.WrapError(auth.ErrorCodeReauthRequired, "OAuth token rejected after recovery", err)
	}

	return err
}

func (recovering *recoveringGraphQLClient) recoverToken(ctx context.Context, reason string) error {
	recoveryGrantType := recovering.recoveryGrantType()
	recovering.logger.DebugContext(
		ctx,
		"auth token recovery started",
		"reason", reason,
		"grant_type", recoveryGrantType,
		"actor", recovering.token.Actor,
		"persist", recovering.persist,
		"profile", recovering.profile,
	)

	var token auth.TokenState
	var err error
	if recovering.token.GrantType == authGrantClientCredentials || recovering.token.RefreshToken == "" {
		token, err = recovering.reacquireClientCredentials(ctx)
	} else {
		token, err = recovering.refreshAuthorizationCode(ctx)
	}
	if err != nil {
		recovering.logTokenRecoveryFailed(ctx, reason, recoveryGrantType, "exchange", err)
		return err
	}
	if recovering.persist {
		if err := recovering.store.SaveTokenState(ctx, recovering.profile, token); err != nil {
			recovering.logTokenRecoveryFailed(ctx, reason, recoveryGrantType, "persist", err)
			return err
		}
	}
	recovering.token = token
	recovering.client = recovering.newClient(token.AccessToken)
	recovering.logger.DebugContext(
		ctx,
		"auth token recovery succeeded",
		"reason", reason,
		"grant_type", recoveryGrantType,
		"actor", token.Actor,
		"persist", recovering.persist,
		"profile", recovering.profile,
	)

	return nil
}

func (recovering *recoveringGraphQLClient) recoveryGrantType() string {
	if recovering.token.GrantType == authGrantClientCredentials || recovering.token.RefreshToken == "" {
		return authGrantClientCredentials
	}

	return authGrantAuthorizationCode
}

func (recovering *recoveringGraphQLClient) logTokenRecoveryFailed(
	ctx context.Context,
	reason string,
	grantType string,
	phase string,
	err error,
) {
	recovering.logger.DebugContext(
		ctx,
		"auth token recovery failed",
		"reason", reason,
		"grant_type", grantType,
		"phase", phase,
		"error_code", errorCode(err),
		"actor", recovering.token.Actor,
		"persist", recovering.persist,
		"profile", recovering.profile,
	)
}

func (recovering *recoveringGraphQLClient) refreshAuthorizationCode(ctx context.Context) (auth.TokenState, error) {
	if recovering.token.RefreshToken == "" || recovering.app.ClientID == "" {
		return auth.TokenState{}, auth.NewError(auth.ErrorCodeReauthRequired, "run linctl auth login")
	}
	token, err := recovering.oauthClient.RefreshToken(ctx, oauth.RefreshTokenRequest{
		RefreshToken: recovering.token.RefreshToken,
		ClientID:     recovering.app.ClientID,
		ClientSecret: recovering.app.ClientSecret,
	})
	if err != nil {
		return auth.TokenState{}, auth.WrapError(
			auth.ErrorCodeRefreshFailed,
			"refresh OAuth token: run linctl auth login",
			err,
		)
	}
	if token.RefreshToken == "" {
		token.RefreshToken = recovering.token.RefreshToken
	}
	token.Actor = recovering.token.Actor
	token.GrantType = authGrantAuthorizationCode
	if err := requireScopes(token.Scopes, requiredScopes(recovering.app)); err != nil {
		return auth.TokenState{}, err
	}

	return token, nil
}

func (recovering *recoveringGraphQLClient) reacquireClientCredentials(ctx context.Context) (auth.TokenState, error) {
	if recovering.app.ClientID == "" || recovering.app.ClientSecret == "" {
		return auth.TokenState{}, auth.NewError(auth.ErrorCodeReauthRequired, "run linctl auth app")
	}
	token, err := exchangeClientCredentialsToken(ctx, recovering.oauthClient, recovering.app)
	if err != nil {
		return auth.TokenState{}, auth.WrapError(
			auth.ErrorCodeReauthRequired,
			"reacquire OAuth app token: run linctl auth app",
			err,
		)
	}

	return token, nil
}

func (recovering *recoveringGraphQLClient) authorizationHeader() string {
	if recovering.token.AccessToken == "" {
		return ""
	}

	return "Bearer " + recovering.token.AccessToken
}

func (runtime commandRuntime) resolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	target, err := client.ResolveTarget(ctx, runtime.graphqlClient, runtime.config.Target)
	logTargetResolution(runtime.log(), target, err)

	return target, err
}

func defaultGlobalConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".config", "linctl", "config.toml")
}

func targetOverride(options *rootOptions) config.Target {
	return config.Target{
		OrgID:     options.orgID,
		TeamKey:   options.team,
		TeamID:    options.teamID,
		ProjectID: options.project,
	}
}

func applyTargetOverrideFlagSemantics(resolved *config.Resolved, options *rootOptions) {
	if options.orgID == "" && options.team == "" && options.teamID == "" {
		return
	}

	resolved.Target.TeamID = options.teamID
}
