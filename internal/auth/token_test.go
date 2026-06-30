package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NewTokenState_maps_endpoint_fields_to_token_state(t *testing.T) {
	t.Parallel()
	expiresAt := time.Date(2026, 6, 30, 12, 0, 0, 0, time.UTC)
	token := NewTokenState(
		"rotated-access-token",
		"rotated-refresh-token",
		"Bearer",
		expiresAt,
		[]string{"read", "write"},
	)

	require.Equal(t, TokenState{
		AccessToken:  "rotated-access-token",
		RefreshToken: "rotated-refresh-token",
		TokenType:    "Bearer",
		Scopes:       []string{"read", "write"},
		ExpiresAt:    &expiresAt,
	}, token)
}

func Test_SplitScopes_accepts_commas_and_whitespace(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		[]string{"read", "write", "issues:create", "comments:create"},
		SplitScopes("read, write\nissues:create\tcomments:create"),
	)
}

func Test_AppConfigEmpty_reports_unset_app_material(t *testing.T) {
	t.Parallel()

	require.True(t, AppConfigEmpty(AppConfig{}))
	require.False(t, AppConfigEmpty(AppConfig{ClientID: "client-id"}))
	require.False(t, AppConfigEmpty(AppConfig{Scopes: []string{"read"}}))
}
