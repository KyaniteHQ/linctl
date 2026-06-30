package auth

import (
	"slices"
	"time"
)

// TokenMetadata is non-secret OAuth token metadata safe for status output.
type TokenMetadata struct {
	TokenType string
	Scopes    []string
	ExpiresAt time.Time
}

// TokenGrant carries OAuth token material plus safe metadata for callers.
type TokenGrant struct {
	State    TokenState
	Metadata TokenMetadata
}

// NewTokenGrant maps token endpoint fields into local auth state and metadata.
func NewTokenGrant(
	accessToken string,
	refreshToken string,
	tokenType string,
	expiresAt time.Time,
	scopes []string,
) TokenGrant {
	var expiresAtPtr *time.Time
	if !expiresAt.IsZero() {
		normalized := expiresAt
		expiresAtPtr = &normalized
	}

	return TokenGrant{
		State: TokenState{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    tokenType,
			Scopes:       slices.Clone(scopes),
			ExpiresAt:    expiresAtPtr,
		},
		Metadata: TokenMetadata{
			TokenType: tokenType,
			Scopes:    slices.Clone(scopes),
			ExpiresAt: expiresAt,
		},
	}
}

// ApplyTokenGrant returns state with token material updated for a profile.
func ApplyTokenGrant(state State, profile string, grant TokenGrant) State {
	if profile == "" {
		state.Token = grant.State
		return state
	}
	if state.Profiles == nil {
		state.Profiles = map[string]ProfileState{}
	}
	profileState := state.Profiles[profile]
	profileState.Token = grant.State
	state.Profiles[profile] = profileState

	return state
}
