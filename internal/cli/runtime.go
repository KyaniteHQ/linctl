package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

type commandRuntime struct {
	config        config.Resolved
	repoConfig    repoConfigSelection
	fileClient    httpDoer
	graphqlClient graphql.Client
	logger        *slog.Logger
}

var buildCommandRuntime = newCommandRuntime

const (
	repoConfigLoaded  = "loaded"
	repoConfigMissing = "missing"
)

type repoConfigSelection struct {
	Path   string
	Status string
}

// resolveConfig loads the layered configuration: the one config-resolution
// path every command uses.
func resolveConfig(ctx context.Context, options *rootOptions) (config.Resolved, error) {
	resolved, _, err := resolveConfigWithSource(ctx, options)

	return resolved, err
}

func resolveConfigWithSource(
	ctx context.Context,
	options *rootOptions,
) (config.Resolved, repoConfigSelection, error) {
	repoConfig, err := selectRepoConfig(options)
	if err != nil {
		return config.Resolved{}, repoConfigSelection{}, err
	}
	resolved, err := config.Load(ctx, config.LoadRequest{
		GlobalPath:      defaultGlobalConfigPath(),
		RepoPath:        repoConfig.Path,
		ProfileOverride: options.profile,
		TargetOverride:  targetOverride(options),
	})
	if err != nil {
		return config.Resolved{}, repoConfigSelection{}, err
	}

	return resolved, repoConfig, nil
}

func selectRepoConfig(options *rootOptions) (repoConfigSelection, error) {
	return selectRepoConfigWithFS(options, filepath.Abs, os.Stat)
}

func selectRepoConfigWithFS(
	options *rootOptions,
	absolutePath func(string) (string, error),
	stat func(string) (os.FileInfo, error),
) (repoConfigSelection, error) {
	path := options.configPath
	explicit := options.configPathExplicit
	if path == "" {
		path = ".linctl.toml"
	}
	resolvedPath, err := absolutePath(path)
	if err != nil {
		return repoConfigSelection{}, fmt.Errorf("resolve repo config path %s: %w", path, err)
	}
	_, err = stat(resolvedPath)
	if err == nil {
		return repoConfigSelection{Path: resolvedPath, Status: repoConfigLoaded}, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		if explicit {
			return repoConfigSelection{}, fmt.Errorf("read explicit repo config %s: %w", resolvedPath, err)
		}

		return repoConfigSelection{Path: resolvedPath, Status: repoConfigMissing}, nil
	}

	return repoConfigSelection{}, fmt.Errorf("inspect repo config %s: %w", resolvedPath, err)
}

func newCommandRuntime(ctx context.Context, options *rootOptions) (commandRuntime, error) {
	stderr := options.stderr
	if stderr == nil {
		stderr = os.Stderr
	}
	logger := newDiagnosticLogger(options.debug, os.Getenv("LINCTL_DEBUG_JSON") == "1", stderr)
	authStatePaths, err := authDefaultPaths(nil)
	if err != nil {
		return commandRuntime{}, err
	}
	resolvedConfig, repoConfig, err := resolveConfigWithSource(ctx, options)
	if err != nil {
		return commandRuntime{}, err
	}
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
		repoConfig: repoConfig,
		fileClient: &http.Client{Timeout: options.timeout},
		logger:     logger,
		graphqlClient: newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token:          authSession.Token,
			CredentialKind: authSession.CredentialKind,
			App:            authSession.App,
			Store:          authStore,
			Profile:        resolvedConfig.Profile,
			Timeout:        options.timeout,
			Logger:         logger,
			OAuthClient:    newAuthOAuthClient(options.timeout),
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
	Token          auth.TokenState
	CredentialKind auth.CredentialKind
	App            auth.AppConfig
	Store          auth.Store
	Profile        string
	Timeout        time.Duration
	Logger         *slog.Logger
	OAuthClient    authOAuthClient
	NewClient      func(accessToken string) graphql.Client
}

type recoveringGraphQLClientState struct {
	token          auth.TokenState
	credentialKind auth.CredentialKind
	client         graphql.Client
}

type recoveringGraphQLClient struct {
	mu          sync.RWMutex
	state       recoveringGraphQLClientState
	app         auth.AppConfig
	store       auth.Store
	profile     string
	timeout     time.Duration
	logger      *slog.Logger
	oauthClient authOAuthClient
	newClient   func(accessToken string) graphql.Client
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
		oauthClient = newAuthOAuthClient(config.Timeout)
	}
	logger := config.Logger
	if logger == nil {
		logger = discardLogger
	}
	credentialKind := config.CredentialKind
	if credentialKind == auth.CredentialKindMissing {
		credentialKind = auth.CredentialKindFromToken(config.Token)
	}
	recovering := &recoveringGraphQLClient{
		state: recoveringGraphQLClientState{
			token:          config.Token,
			credentialKind: credentialKind,
		},
		app:         config.App,
		store:       config.Store,
		profile:     config.Profile,
		timeout:     config.Timeout,
		logger:      logger,
		oauthClient: oauthClient,
		newClient:   newClient,
	}
	recovering.state.client = newClient(config.Token.AccessToken)

	return recovering
}

func (recovering *recoveringGraphQLClient) snapshot() recoveringGraphQLClientState {
	recovering.mu.RLock()
	defer recovering.mu.RUnlock()

	return recovering.state
}

func (recovering *recoveringGraphQLClient) setState(state recoveringGraphQLClientState) {
	recovering.mu.Lock()
	recovering.state = state
	recovering.mu.Unlock()
}

func (recovering *recoveringGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	state := recovering.snapshot()
	recovered := false
	if state.credentialKind.Recoverable() && tokenExpired(state.token, authNow()) {
		next, err := recovering.recoverToken(ctx, "expired", state)
		if err != nil {
			return err
		}
		state = next
		recovered = true
	}
	err := state.client.MakeRequest(ctx, request, response)
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
	if !state.credentialKind.Recoverable() {
		return auth.WrapError(
			auth.ErrorCodeReauthRequired,
			"OAuth token rejected and credential is not recoverable: run linctl auth login or linctl auth app",
			err,
		)
	}
	next, err := recovering.recoverToken(ctx, "auth_failed", state)
	if err != nil {
		return err
	}

	err = next.client.MakeRequest(ctx, request, response)
	if errors.Is(err, client.ErrAuthFailed) {
		return auth.WrapError(auth.ErrorCodeReauthRequired, "OAuth token rejected after recovery", err)
	}

	return err
}

func (recovering *recoveringGraphQLClient) recoverToken(
	ctx context.Context,
	reason string,
	observed recoveringGraphQLClientState,
) (recoveringGraphQLClientState, error) {
	recoveryGrantType := recoveryGrantTypeFor(observed.credentialKind)
	recovering.logger.DebugContext(
		ctx,
		"auth token recovery started",
		"reason", reason,
		"grant_type", recoveryGrantType,
		"actor", observed.token.Actor,
		"profile", recovering.profile,
	)

	token, _, err := transactRecoveredToken(
		ctx,
		recovering.store,
		recovering.profile,
		auth.NewError(auth.ErrorCodeReauthRequired, "persisted OAuth session was removed"),
		observed.token,
		func(current auth.TokenState) (auth.TokenState, error) {
			return recovering.exchangeCurrentToken(ctx, current)
		},
	)
	if err != nil {
		recovering.logTokenRecoveryFailed(ctx, reason, recoveryGrantType, "exchange", err, observed.token.Actor)
		return recoveringGraphQLClientState{}, err
	}
	next := recoveringGraphQLClientState{
		token:          token,
		credentialKind: auth.CredentialKindFromToken(token),
		client:         recovering.newClient(token.AccessToken),
	}
	recovering.setState(next)
	recovering.logger.DebugContext(
		ctx,
		"auth token recovery succeeded",
		"reason", reason,
		"grant_type", recoveryGrantType,
		"actor", token.Actor,
		"profile", recovering.profile,
	)

	return next, nil
}

func (recovering *recoveringGraphQLClient) exchangeCurrentToken(
	ctx context.Context,
	current auth.TokenState,
) (auth.TokenState, error) {
	credentialKind := auth.CredentialKindFromToken(current)
	if credentialKind == auth.CredentialKindClientCredentials {
		return recovering.reacquireClientCredentials(ctx)
	}
	if credentialKind == auth.CredentialKindAuthorizationCode {
		return recovering.refreshAuthorizationCode(ctx, current)
	}

	return auth.TokenState{}, auth.NewError(
		auth.ErrorCodeReauthRequired,
		"persisted OAuth token is not recoverable",
	)
}

func recoveryGrantTypeFor(kind auth.CredentialKind) string {
	if kind == auth.CredentialKindClientCredentials {
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
	actor string,
) {
	recovering.logger.DebugContext(
		ctx,
		"auth token recovery failed",
		"reason", reason,
		"grant_type", grantType,
		"phase", phase,
		"error_code", errorCode(err),
		"actor", actor,
		"profile", recovering.profile,
	)
}

func (recovering *recoveringGraphQLClient) refreshAuthorizationCode(
	ctx context.Context,
	token auth.TokenState,
) (auth.TokenState, error) {
	if token.RefreshToken == "" || recovering.app.ClientID == "" {
		return auth.TokenState{}, auth.NewError(auth.ErrorCodeReauthRequired, "run linctl auth login")
	}

	return refreshAuthorizationCodeToken(
		ctx,
		recovering.oauthClient,
		recovering.app,
		token,
		requiredScopes(recovering.app),
	)
}

func tokenUsable(token auth.TokenState, now time.Time) bool {
	return token.AccessToken != "" && !tokenExpired(token, now)
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
