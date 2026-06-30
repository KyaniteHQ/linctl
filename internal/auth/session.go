package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	oauthAccessTokenEnv  = "LINCTL_OAUTH_ACCESS_TOKEN" //nolint:gosec // Environment variable name, not a secret.
	oauthClientIDEnv     = "LINCTL_OAUTH_CLIENT_ID"
	oauthClientSecretEnv = "LINCTL_OAUTH_CLIENT_SECRET" //nolint:gosec // Environment variable name, not a secret.
	oauthRedirectURIEnv  = "LINCTL_OAUTH_REDIRECT_URI"
	oauthScopesEnv       = "LINCTL_OAUTH_SCOPES"
)

// Session is the selected OAuth state for one linctl profile.
type Session struct {
	State           State
	Profile         string
	App             AppConfig
	Token           TokenState
	TokenSource     string
	PersistentToken bool
}

// SessionRequest describes the sources used to select OAuth state.
type SessionRequest struct {
	Env     Env
	Store   Store
	Profile string
}

// SelectSession loads local auth state once and overlays process OAuth material.
func SelectSession(ctx context.Context, request SessionRequest) (Session, error) {
	if err := ctx.Err(); err != nil {
		return Session{}, fmt.Errorf("select auth session context: %w", err)
	}

	activeEnv := activeEnv(request.Env)
	state, err := request.Store.Load(ctx)
	if err != nil {
		return Session{}, err
	}
	profileState := state.Profile(request.Profile)
	app := MergeAppConfig(profileState.App, AppConfigFromEnv(activeEnv))
	token := profileState.Token
	persistentToken := token.AccessToken != "" || token.RefreshToken != ""
	tokenSource := "local"
	if !persistentToken {
		tokenSource = "missing"
	}

	if envToken, ok := activeEnv.Lookup(oauthAccessTokenEnv); ok && envToken != "" {
		if isPersonalAPIKeyShape(envToken) {
			return Session{}, fmt.Errorf("%s contains personal API key-shaped material", oauthAccessTokenEnv)
		}
		token = TokenState{AccessToken: envToken}
		tokenSource = "env"
		persistentToken = false
	} else if isPersonalAPIKeyShape(token.AccessToken) {
		return Session{}, errors.New("local auth state contains personal API key-shaped material")
	}

	return Session{
		State:           state,
		Profile:         request.Profile,
		App:             app,
		Token:           token,
		TokenSource:     tokenSource,
		PersistentToken: persistentToken,
	}, nil
}

// AppConfigFromEnv resolves OAuth app material from environment variables.
func AppConfigFromEnv(env Env) AppConfig {
	activeEnv := activeEnv(env)
	clientID, _ := activeEnv.Lookup(oauthClientIDEnv)
	clientSecret, _ := activeEnv.Lookup(oauthClientSecretEnv)
	redirectURI, _ := activeEnv.Lookup(oauthRedirectURIEnv)
	scopeText, _ := activeEnv.Lookup(oauthScopesEnv)

	return AppConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       SplitScopes(scopeText),
	}
}

// MergeAppConfig overlays explicitly supplied OAuth app material onto a base config.
func MergeAppConfig(base AppConfig, override AppConfig) AppConfig {
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

func isPersonalAPIKeyShape(value string) bool {
	return strings.HasPrefix(value, "lin_api_")
}
