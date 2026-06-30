package auth

import (
	"slices"
	"strings"
	"time"
	"unicode"
)

// NewTokenState maps OAuth token endpoint fields into local auth state.
func NewTokenState(
	accessToken string,
	refreshToken string,
	tokenType string,
	expiresAt time.Time,
	scopes []string,
) TokenState {
	var expiresAtPtr *time.Time
	if !expiresAt.IsZero() {
		normalized := expiresAt
		expiresAtPtr = &normalized
	}

	return TokenState{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		Scopes:       slices.Clone(scopes),
		ExpiresAt:    expiresAtPtr,
	}
}

// SplitScopes parses comma- or whitespace-delimited OAuth scopes.
func SplitScopes(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})
}

// AppConfigEmpty reports whether OAuth app configuration is unset.
func AppConfigEmpty(app AppConfig) bool {
	return app.ClientID == "" &&
		app.ClientSecret == "" &&
		app.RedirectURI == "" &&
		len(app.Scopes) == 0
}
