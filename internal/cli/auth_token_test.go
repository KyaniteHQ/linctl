package cli

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
)

// Test_refreshAuthorizationCodeToken_preserves_refresh_token_when_endpoint_omits_one
// pins the same rotation guarantee the recovering runtime client has: when the
// OAuth refresh response carries no refresh_token, the stored one must survive
// so a later `linctl auth refresh` does not strand the profile.
func Test_refreshAuthorizationCodeToken_preserves_refresh_token_when_endpoint_omits_one(t *testing.T) {
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"new-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	app := auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}}
	previous := auth.TokenState{RefreshToken: "old-refresh-token", Actor: "user"}

	token, err := refreshAuthorizationCodeToken(context.Background(), fakeOAuth, app, previous, requiredScopes(app))

	require.NoError(t, err)
	require.Equal(t, "old-refresh-token", token.RefreshToken)
	require.Equal(t, "user", token.Actor)
	require.Equal(t, authGrantAuthorizationCode, token.GrantType)
}
