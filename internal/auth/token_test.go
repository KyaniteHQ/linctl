package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ApplyTokenGrant_persists_rotated_token_state_without_deleting_app_config(t *testing.T) {
	store := NewStore(testPaths(t))
	state := State{
		App: AppConfig{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Profiles: map[string]ProfileState{
			"work": {
				App: AppConfig{ClientID: "work-client-id"},
			},
		},
	}
	expiresAt := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	grant := NewTokenGrant(
		"rotated-access-token",
		"rotated-refresh-token",
		"Bearer",
		expiresAt,
		[]string{"read", "write"},
	)

	require.NoError(t, store.Save(context.Background(), ApplyTokenGrant(state, "work", grant)))

	got, err := store.Load(context.Background())
	require.NoError(t, err)
	require.Equal(t, AppConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}, got.App)
	require.Equal(t, AppConfig{ClientID: "work-client-id"}, got.Profiles["work"].App)
	require.Equal(t, TokenState{
		AccessToken:  "rotated-access-token",
		RefreshToken: "rotated-refresh-token",
		TokenType:    "Bearer",
		Scopes:       []string{"read", "write"},
		ExpiresAt:    grant.State.ExpiresAt,
	}, got.Profiles["work"].Token)
}
