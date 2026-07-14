package cli

import (
	"context"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
)

func Test_CommandRuntime_logout_wins_over_in_flight_refresh(t *testing.T) {
	store := auth.NewStore(cliAuthTestPaths(t))
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		GrantType:    authGrantAuthorizationCode,
	}
	require.NoError(t, store.Save(context.Background(), auth.State{App: app, Token: token}))
	exchangeStarted := make(chan struct{})
	releaseExchange := make(chan struct{})
	fakeOAuth := &fakeOAuthTokenClient{
		grant: auth.NewTokenState(
			"refreshed-access-token",
			"rotated-refresh-token",
			"Bearer",
			time.Now().Add(time.Hour),
			[]string{"read"},
		),
		beforeRefresh: func() {
			close(exchangeStarted)
			<-releaseExchange
		},
	}
	runtimeClient := newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
		Token:          token,
		CredentialKind: auth.CredentialKindAuthorizationCode,
		App:            app,
		Store:          store,
		OAuthClient:    fakeOAuth,
		NewClient:      (&recordingRuntimeClientFactory{}).newClient,
	})
	requestDone := make(chan error, 1)
	go func() {
		requestDone <- runtimeClient.MakeRequest(
			context.Background(),
			&graphql.Request{Query: "query Test { viewer { id } }"},
			&graphql.Response{},
		)
	}()

	<-exchangeStarted
	clearDone := make(chan error, 1)
	go func() {
		clearDone <- store.ClearTokenState(context.Background(), "")
	}()
	clearFinished := false
	select {
	case err := <-clearDone:
		require.NoError(t, err)
		clearFinished = true
	case <-time.After(100 * time.Millisecond):
	}
	close(releaseExchange)
	require.NoError(t, <-requestDone)
	if !clearFinished {
		require.NoError(t, <-clearDone)
	}

	finalState, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Empty(t, finalState.Token, "logout must win over an in-flight refresh")
}

func Test_CommandRuntime_concurrent_refreshes_exchange_once_and_adopt(t *testing.T) {
	store := auth.NewStore(cliAuthTestPaths(t))
	expiredAt := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	app := auth.AppConfig{ClientID: "client-id", Scopes: []string{"read"}}
	token := auth.TokenState{
		AccessToken:  "expired-access-token",
		RefreshToken: "old-refresh-token",
		ExpiresAt:    &expiredAt,
		GrantType:    authGrantAuthorizationCode,
	}
	require.NoError(t, store.Save(context.Background(), auth.State{App: app, Token: token}))
	fakeOAuth := &fakeOAuthTokenClient{grant: auth.NewTokenState(
		"refreshed-access-token",
		"rotated-refresh-token",
		"Bearer",
		time.Now().Add(time.Hour),
		[]string{"read"},
	)}
	clients := []*recoveringGraphQLClient{
		newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token:          token,
			CredentialKind: auth.CredentialKindAuthorizationCode,
			App:            app,
			Store:          store,
			OAuthClient:    fakeOAuth,
			NewClient:      (&recordingRuntimeClientFactory{}).newClient,
		}),
		newRecoveringGraphQLClient(recoveringGraphQLClientConfig{
			Token:          token,
			CredentialKind: auth.CredentialKindAuthorizationCode,
			App:            app,
			Store:          store,
			OAuthClient:    fakeOAuth,
			NewClient:      (&recordingRuntimeClientFactory{}).newClient,
		}),
	}
	start := make(chan struct{})
	results := make(chan error, len(clients))
	for _, runtimeClient := range clients {
		go func() {
			<-start
			results <- runtimeClient.MakeRequest(
				context.Background(),
				&graphql.Request{Query: "query Test { viewer { id } }"},
				&graphql.Response{},
			)
		}()
	}

	close(start)
	for range clients {
		require.NoError(t, <-results)
	}
	require.Equal(t, 1, fakeOAuth.refreshTokenCalls)
	state, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, "rotated-refresh-token", state.Token.RefreshToken)
}
