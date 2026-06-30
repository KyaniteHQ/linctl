package cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
)

func Test_AuthApp_passes_root_timeout_to_oauth_client_factory(t *testing.T) {
	paths := cliAuthTestPaths(t)
	require.NoError(t, auth.NewStore(paths).SaveAppConfig(context.Background(), "", auth.AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read"},
	}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"oauth-access-token",
		"",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	fakeReadiness := &fakeAuthReadinessChecker{report: readyAuthReport("app")}
	restorePaths := useAuthPaths(t, paths)
	defer restorePaths()
	originalOAuthClient := newAuthOAuthClient
	originalReadiness := checkAuthReadiness
	var gotTimeout time.Duration
	newAuthOAuthClient = func(timeout time.Duration) authOAuthClient {
		gotTimeout = timeout
		return fakeOAuth
	}
	checkAuthReadiness = fakeReadiness.check
	defer func() {
		checkAuthReadiness = originalReadiness
		newAuthOAuthClient = originalOAuthClient
	}()

	err := execute(context.Background(), BuildInfo{}, nil, &bytes.Buffer{}, &bytes.Buffer{}, []string{
		"--timeout", "7s",
		"--org", "org-id",
		"--team", "LIT",
		"--team-id", "team-id",
		"auth",
		"app",
	})

	require.NoError(t, err)
	require.Equal(t, 7*time.Second, gotTimeout)
}
