package cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
)

func Test_AuthRefresh_keeps_the_saved_refresh_token_when_linear_does_not_rotate_it(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).Save(context.Background(), auth.State{
		App: auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}},
		Token: auth.TokenState{
			AccessToken:  "old-access-token",
			RefreshToken: "old-refresh-token",
			Scopes:       []string{"read"},
			Actor:        "user",
			GrantType:    authGrantAuthorizationCode,
		},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"new-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	restore := useAuthCommandHooks(t, paths, fakeOAuth, &fakeAuthReadinessChecker{report: readyAuthReport("user")})
	defer restore()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"refresh",
	})

	require.NoError(t, err)
	got, err := auth.NewStore(paths).Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "new-access-token", got.Token.AccessToken)
	require.Equal(t, "old-refresh-token", got.Token.RefreshToken)
	require.Equal(t, "user", got.Token.Actor)
	require.Equal(t, authGrantAuthorizationCode, got.Token.GrantType)
}
