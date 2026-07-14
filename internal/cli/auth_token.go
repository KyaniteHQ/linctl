package cli

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

// stampAndRequireScopes applies the one post-grant rule every OAuth exchange
// shares: stamp the actor and grant type on the token state, then require the
// granted scopes before the token is used or persisted.
func stampAndRequireScopes(
	token auth.TokenState,
	actor string,
	grantType string,
	scopes []string,
) (auth.TokenState, error) {
	token.Actor = actor
	token.GrantType = grantType
	if err := requireScopes(token.Scopes, scopes); err != nil {
		return auth.TokenState{}, err
	}

	return token, nil
}

// finalizeRefreshedToken normalizes a refreshed authorization-code token:
// Linear may omit the refresh token when it does not rotate it, so the
// previous one must survive, and the actor carries over from the previous
// token state.
func finalizeRefreshedToken(
	refreshed auth.TokenState,
	previous auth.TokenState,
	scopes []string,
) (auth.TokenState, error) {
	if refreshed.RefreshToken == "" {
		refreshed.RefreshToken = previous.RefreshToken
	}

	return stampAndRequireScopes(
		refreshed,
		firstNonEmptyString(previous.Actor, appActor),
		authGrantAuthorizationCode,
		scopes,
	)
}

func currentOrRefreshedAuthToken(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	token auth.TokenState,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	if token.AccessToken != "" || token.RefreshToken != "" {
		return refreshPersistedAuthToken(ctx, authContext, app, timeout)
	}
	if app.ClientID == "" || app.ClientSecret == "" {
		return auth.TokenState{}, authReadinessReport{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"run linctl auth configure and linctl auth app",
		)
	}

	return acquirePersistedClientCredentialsToken(ctx, authContext, app, timeout)
}

func acquirePersistedClientCredentialsToken(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	var readiness authReadinessReport
	readinessChecked := false
	token, err := authContext.store.TransactTokenState(
		ctx,
		authContext.profile,
		func(current auth.TokenState) (auth.TokenState, error) {
			if tokenUsable(current, authNow()) {
				return current, nil
			}

			acquired, checkedReadiness, err := acquireClientCredentialsToken(ctx, authContext, app, timeout)
			if err != nil {
				return auth.TokenState{}, err
			}
			readiness = checkedReadiness
			readinessChecked = true

			return acquired, nil
		},
	)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	if readinessChecked {
		return token, readiness, nil
	}

	readiness, err = refreshedAuthReadiness(ctx, authContext, app, token, token.Actor, timeout)
	return token, readiness, err
}

func acquireClientCredentialsToken(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	scopes := requiredScopes(app)
	token, err := exchangeClientCredentialsToken(ctx, newAuthOAuthClient(timeout), app)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	readiness, err := requireLoggedAuthReadiness(ctx, authContext.log(), authReadinessRequest{
		AccessToken:    token.AccessToken,
		TokenActor:     token.Actor,
		TokenScopes:    token.Scopes,
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
	refreshed, err := refreshAuthToken(ctx, app, token, timeout)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}

	readiness, err := refreshedAuthReadiness(ctx, authContext, app, refreshed, token.Actor, timeout)
	return refreshed, readiness, err
}

func refreshPersistedAuthToken(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	timeout time.Duration,
) (auth.TokenState, authReadinessReport, error) {
	var readiness authReadinessReport
	readinessChecked := false
	token, err := authContext.store.TransactTokenState(
		ctx,
		authContext.profile,
		func(current auth.TokenState) (auth.TokenState, error) {
			if current.AccessToken == "" && current.RefreshToken == "" {
				return auth.TokenState{}, missingPersistedTokenError(authContext.credentialKind)
			}
			if !current.Equal(authContext.localToken) && tokenUsable(current, authNow()) {
				return current, nil
			}

			refreshed, checkedReadiness, err := refreshAuthTokenState(ctx, authContext, app, current, timeout)
			if err != nil {
				return auth.TokenState{}, err
			}
			readiness = checkedReadiness
			readinessChecked = true

			return refreshed, nil
		},
	)
	if err != nil {
		return auth.TokenState{}, authReadinessReport{}, err
	}
	if readinessChecked {
		return token, readiness, nil
	}

	readiness, err = refreshedAuthReadiness(ctx, authContext, app, token, token.Actor, timeout)
	return token, readiness, err
}

func refreshAuthToken(
	ctx context.Context,
	app auth.AppConfig,
	token auth.TokenState,
	timeout time.Duration,
) (auth.TokenState, error) {
	if token.AccessToken == "" && token.RefreshToken == "" {
		return auth.TokenState{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing OAuth token state: run linctl auth login or linctl auth app",
		)
	}
	if token.GrantType == authGrantClientCredentials {
		if app.ClientID == "" || app.ClientSecret == "" {
			return auth.TokenState{}, auth.NewError(
				auth.ErrorCodeNotConfigured,
				"missing OAuth app client credentials: run linctl auth configure",
			)
		}

		return exchangeClientCredentialsToken(ctx, newAuthOAuthClient(timeout), app)
	}
	if token.RefreshToken == "" {
		return auth.TokenState{}, auth.NewError(
			auth.ErrorCodeReauthRequired,
			"missing OAuth refresh token: run linctl auth login or linctl auth app",
		)
	}
	if app.ClientID == "" {
		return auth.TokenState{}, auth.NewError(
			auth.ErrorCodeNotConfigured,
			"missing OAuth client id: run linctl auth configure",
		)
	}

	scopes := requiredScopes(app)
	refreshed, err := refreshAuthorizationCodeToken(ctx, newAuthOAuthClient(timeout), app, token, scopes)
	if err != nil {
		return auth.TokenState{}, err
	}

	return refreshed, nil
}

func refreshedAuthReadiness(
	ctx context.Context,
	authContext authCommandContext,
	app auth.AppConfig,
	token auth.TokenState,
	expectedActor string,
	timeout time.Duration,
) (authReadinessReport, error) {
	return requireLoggedAuthReadiness(ctx, authContext.log(), authReadinessRequest{
		AccessToken:    token.AccessToken,
		TokenActor:     token.Actor,
		TokenScopes:    token.Scopes,
		ExpectedTarget: authContext.target,
		ExpectedActor:  firstNonEmptyString(expectedActor, appActor),
		RequiredScopes: requiredScopes(app),
		Timeout:        timeout,
	})
}

func missingPersistedTokenError(kind auth.CredentialKind) error {
	if kind == auth.CredentialKindInjectedAccessToken {
		return auth.NewError(
			auth.ErrorCodeReauthRequired,
			"running on an injected LINCTL_OAUTH_ACCESS_TOKEN; unset it to manage the local session",
		)
	}

	return auth.NewError(
		auth.ErrorCodeNotConfigured,
		"missing OAuth token state: run linctl auth login or linctl auth app",
	)
}

func refreshAuthorizationCodeToken(
	ctx context.Context,
	oauthClient authOAuthClient,
	app auth.AppConfig,
	token auth.TokenState,
	scopes []string,
) (auth.TokenState, error) {
	refreshed, err := oauthClient.RefreshToken(ctx, oauth.RefreshTokenRequest{
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

	return finalizeRefreshedToken(refreshed, token, scopes)
}

func revokeTokenState(
	ctx context.Context,
	logger *slog.Logger,
	oauthClient authOAuthClient,
	token auth.TokenState,
) ([]string, bool) {
	if logger == nil {
		logger = discardLogger
	}
	revoked := []string{}
	failed := false
	for _, request := range []oauth.RevocationRequest{
		{Token: token.RefreshToken, TokenTypeHint: "refresh_token"},
		{Token: token.AccessToken, TokenTypeHint: "access_token"},
	} {
		if request.Token == "" {
			continue
		}
		logger.DebugContext(
			ctx,
			"auth token revoke started",
			"token_type", request.TokenTypeHint,
		)
		if err := oauthClient.RevokeToken(ctx, request); err != nil {
			failed = true
			logger.DebugContext(
				ctx,
				"auth token revoke failed",
				"token_type", request.TokenTypeHint,
				"error_code", errorCode(err),
			)
			continue
		}
		logger.DebugContext(
			ctx,
			"auth token revoke succeeded",
			"token_type", request.TokenTypeHint,
		)
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
	token, err := oauthClient.ClientCredentials(ctx, oauth.ClientCredentialsRequest{
		ClientID:     app.ClientID,
		ClientSecret: app.ClientSecret,
		Scopes:       scopes,
	})
	if err != nil {
		return auth.TokenState{}, err
	}

	return stampAndRequireScopes(token, appActor, authGrantClientCredentials, scopes)
}

func newOAuthHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}

func tokenExpired(token auth.TokenState, now time.Time) bool {
	return token.ExpiresAt != nil && !token.ExpiresAt.After(now)
}
