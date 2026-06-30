package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_DefaultCheckAuthReadiness_requires_access_token(t *testing.T) {
	_, err := defaultCheckAuthReadiness(context.Background(), authReadinessRequest{})

	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeNotConfigured), errorCode(err))
}

func Test_AuthReadiness_and_status_helpers(t *testing.T) {
	t.Run("default readiness client constructor", func(t *testing.T) {
		require.NotNil(t, newAuthReadinessGraphQLClient("access-token", time.Second))
	})

	t.Run("actor mismatch", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{},
			&fakeAuthReadinessChecker{report: readyAuthReport("user")},
		)
		defer restore()

		_, err := requireAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken:   "access-token",
			ExpectedActor: "app",
		})

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeActorMismatch), errorCode(err))
	})

	t.Run("expected actor without proven actor fails closed", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{},
			&fakeAuthReadinessChecker{report: authReadinessReport{}},
		)
		defer restore()

		_, err := requireAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken:   "access-token",
			ExpectedActor: "app",
		})

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeActorMismatch), errorCode(err))
	})

	t.Run("missing required scope fails readiness", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{},
			&fakeAuthReadinessChecker{report: readyAuthReport("app")},
		)
		defer restore()

		_, err := requireAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken:    "access-token",
			TokenActor:     "app",
			TokenScopes:    []string{"read"},
			ExpectedActor:  "app",
			RequiredScopes: []string{"read", "write"},
		})

		require.Error(t, err)
		require.Equal(t, string(auth.ErrorCodeMissingScope), errorCode(err))
	})

	t.Run("readiness error mapping", func(t *testing.T) {
		authErr := auth.NewError(auth.ErrorCodeNotConfigured, "missing")
		require.Same(t, authErr, mapAuthReadinessError(authErr))

		tokenErr := auth.NewTokenEndpointError(auth.ErrorCodeRefreshFailed, 401, "invalid_grant")
		require.Same(t, tokenErr, mapAuthReadinessError(tokenErr))

		targetErr := mapAuthReadinessError(client.ErrTargetNotConfigured)
		require.Equal(t, string(auth.ErrorCodeTargetMismatch), errorCode(targetErr))

		plainErr := errors.New("network down")
		require.Same(t, plainErr, mapAuthReadinessError(plainErr))
	})

	t.Run("logged readiness accepts nil logger and optional project", func(t *testing.T) {
		restore := useAuthCommandHooks(
			t,
			cliAuthTestPaths(t),
			&fakeOAuthTokenClient{},
			&fakeAuthReadinessChecker{
				report: authReadinessReport{
					Actor: "app",
					Target: client.ResolvedTarget{
						Org:       client.TargetOrg{ID: "org-id"},
						Team:      client.TargetTeam{ID: "team-id", Key: "LIT"},
						Project:   &client.ResolvedProject{ID: "project-id"},
						Confirmed: true,
					},
				},
			},
		)
		defer restore()

		readiness, err := requireLoggedAuthReadiness(context.Background(), nil, authReadinessRequest{
			AccessToken:   "access-token",
			ExpectedActor: "app",
		})

		require.NoError(t, err)
		require.Equal(t, "project-id", readiness.Target.Project.ID)
	})

	t.Run("revoke accepts nil logger", func(t *testing.T) {
		fakeOAuth := &fakeOAuthTokenClient{}

		revoked, failed := revokeTokenState(context.Background(), nil, fakeOAuth, auth.TokenState{
			AccessToken: "access-token",
		})

		require.False(t, failed)
		require.Equal(t, []string{"access_token"}, revoked)
	})

	t.Run("default readiness success and resolve error", func(t *testing.T) {
		restore := useAuthReadinessGraphQLClient(t, authReadinessFakeGraphQLClient{
			"Viewer": `{
				"viewer": {
					"id": "user-id",
					"name": "Omer",
					"displayName": "Omer",
					"email": "omer@example.com",
					"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
				}
			}`,
			"Teams": `{
				"teams": {
					"nodes": [{
						"id": "team-id",
						"key": "LIT",
						"name": "linctl-it",
						"organization": {"id": "org-id", "name": "Kyanite", "urlKey": "kyanite"}
					}],
					"pageInfo": {"hasNextPage": false, "endCursor": null}
				}
			}`,
		})
		defer restore()

		report, err := defaultCheckAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken:   "access-token",
			TokenActor:    "app",
			TokenScopes:   []string{"read"},
			ExpectedActor: "app",
			ExpectedTarget: config.Target{
				OrgID:   "org-id",
				TeamKey: "LIT",
				TeamID:  "team-id",
			},
		})
		require.NoError(t, err)
		require.Equal(t, "app", report.Actor)
		require.Equal(t, "team-id", report.Target.Team.ID)

		_, err = defaultCheckAuthReadiness(context.Background(), authReadinessRequest{
			AccessToken: "access-token",
		})
		require.Error(t, err)
		require.ErrorIs(t, err, client.ErrTargetNotConfigured)
	})

	t.Run("write status quiet", func(t *testing.T) {
		var stdout bytes.Buffer
		command := &cobra.Command{}
		command.SetOut(&stdout)

		err := writeAuthStatus(command, &rootOptions{quiet: true}, authStatusReport{
			Token: authTokenReport{Status: "set"},
		})

		require.NoError(t, err)
		require.Empty(t, stdout.String())
	})

	require.Equal(t, "fallback", firstNonEmptyString("", "fallback"))
}
