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
	// GrantTypeAuthorizationCode is the persisted authorization-code OAuth grant variant.
	GrantTypeAuthorizationCode = "authorization_code"
	// GrantTypeClientCredentials is the persisted client-credentials OAuth grant variant.
	GrantTypeClientCredentials = "client_credentials"
)

// CredentialKind identifies the provenance and recovery contract of an OAuth credential.
type CredentialKind uint8

const (
	// CredentialKindMissing identifies an absent OAuth credential.
	CredentialKindMissing CredentialKind = iota
	// CredentialKindInjectedAccessToken identifies a non-recoverable process override.
	CredentialKindInjectedAccessToken
	// CredentialKindLocalAccessToken identifies local token material without an explicit recovery grant.
	CredentialKindLocalAccessToken
	// CredentialKindAuthorizationCode identifies a refreshable authorization-code credential.
	CredentialKindAuthorizationCode
	// CredentialKindClientCredentials identifies a reacquirable OAuth app credential.
	CredentialKindClientCredentials
)

// Recoverable reports whether the credential kind permits an OAuth token exchange.
func (kind CredentialKind) Recoverable() bool {
	return kind == CredentialKindAuthorizationCode || kind == CredentialKindClientCredentials
}

// Session is the selected OAuth state for one linctl profile.
type Session struct {
	State           State
	Profile         string
	App             AppConfig
	Token           TokenState
	CredentialKind  CredentialKind
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
	if envToken, ok := activeEnv.Lookup(oauthAccessTokenEnv); ok && envToken != "" {
		if isPersonalAPIKeyShape(envToken) {
			return Session{}, fmt.Errorf("%s contains personal API key-shaped material", oauthAccessTokenEnv)
		}

		return Session{
			Profile:         request.Profile,
			Token:           TokenState{AccessToken: envToken},
			CredentialKind:  CredentialKindInjectedAccessToken,
			TokenSource:     "env",
			PersistentToken: false,
		}, nil
	}

	return selectLocalSession(ctx, request, activeEnv)
}

// SelectLocalSession loads profile state and process app configuration without applying an injected access token.
func SelectLocalSession(ctx context.Context, request SessionRequest) (Session, error) {
	if err := ctx.Err(); err != nil {
		return Session{}, fmt.Errorf("select local auth session context: %w", err)
	}

	return selectLocalSession(ctx, request, activeEnv(request.Env))
}

func selectLocalSession(ctx context.Context, request SessionRequest, env Env) (Session, error) {
	state, err := request.Store.Load(ctx)
	if err != nil {
		return Session{}, err
	}
	profileState := state.Profile(request.Profile)
	app := MergeAppConfig(profileState.App, AppConfigFromEnv(env))
	token := profileState.Token
	persistentToken := token.AccessToken != "" || token.RefreshToken != ""
	credentialKind := CredentialKindFromToken(token)
	tokenSource := "local"
	if !persistentToken {
		tokenSource = "missing"
	}

	if isPersonalAPIKeyShape(token.AccessToken) {
		return Session{}, errors.New("local auth state contains personal API key-shaped material")
	}

	return Session{
		State:           state,
		Profile:         request.Profile,
		App:             app,
		Token:           token,
		CredentialKind:  credentialKind,
		TokenSource:     tokenSource,
		PersistentToken: persistentToken,
	}, nil
}

// CredentialKindFromToken returns the explicit recovery variant persisted with a local token.
func CredentialKindFromToken(token TokenState) CredentialKind {
	if token.AccessToken == "" && token.RefreshToken == "" {
		return CredentialKindMissing
	}

	switch token.GrantType {
	case GrantTypeAuthorizationCode:
		return CredentialKindAuthorizationCode
	case GrantTypeClientCredentials:
		return CredentialKindClientCredentials
	default:
		return CredentialKindLocalAccessToken
	}
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
