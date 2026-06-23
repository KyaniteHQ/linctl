package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

func Test_TransportScenarios_return_actionable_errors(t *testing.T) {
	require.Equal(t, "fallback", firstNonEmpty("", "fallback"))
	require.Equal(t, "primary", firstNonEmpty("primary", "fallback"))
	require.Equal(t, 3*time.Second, defaultDuration(3*time.Second, time.Second))
	require.Equal(t, time.Second, defaultDuration(0, time.Second))
	require.Equal(t, 200*time.Millisecond, retryDelay(http.Header{}, 1))
	require.Equal(t, 100*time.Millisecond, retryDelay(http.Header{"Retry-After": []string{"not-a-number"}}, 0))
	require.Equal(t, 2*time.Second, retryDelay(http.Header{"Retry-After": []string{"2"}}, 0))
	require.Equal(t, maxRetryDelay, retryDelay(http.Header{"Retry-After": []string{"120"}}, 0))

	require.True(t, isRateLimited(http.StatusTooManyRequests, nil))
	require.True(t, isRateLimited(http.StatusBadRequest, []byte(`{"errors":[{"extensions":{"code":"RATELIMITED"}}]}`)))
	require.False(t, isRateLimited(http.StatusBadRequest, []byte(`{"errors":[{"extensions":{"code":"BAD_USER_INPUT"}}]}`)))
	require.False(t, isRateLimited(http.StatusBadRequest, []byte("not json")))
	require.False(t, isRateLimited(http.StatusOK, nil))
	rateLimitedBody := []byte(`{"errors":[{"extensions":{"code":"RATELIMITED"}}]}`)
	require.False(t, isRateLimited(http.StatusInternalServerError, rateLimitedBody))
	require.ErrorIs(t, rateLimitError(http.StatusTooManyRequests, []byte("slow down")), ErrRateLimited)

	response := graphql.Response{}
	err := decodeGraphQLResponse([]byte("not json"), http.StatusOK, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode graphql response")

	err = decodeGraphQLResponse([]byte("server down"), http.StatusBadGateway, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "graphql http status 502")
}

func Test_CustomViewPreferenceReads_return_empty_values_when_organization_defaults_are_absent(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient{
		"customView_organizationViewPreferences":             `{"customView":{"organizationViewPreferences":null}}`,
		"customView_organizationViewPreferences_preferences": `{"customView":{"organizationViewPreferences":null}}`,
	}

	// When
	preferences, err := GetCustomViewOrganizationPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	values, err := GetCustomViewOrganizationPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)

	// Then
	require.Equal(t, "custom-view-id", preferences.CustomViewID)
	require.Empty(t, preferences.ID)
	require.Equal(t, "custom-view-id", values.CustomViewID)
	require.False(t, values.HasOrganizationPreferences)
}

func Test_CustomViewPreferenceReads_return_empty_values_when_user_preferences_are_absent(t *testing.T) {
	// Given
	graphqlClient := fakeGraphQLClient{
		"customView_userViewPreferences":             `{"customView":{"userViewPreferences":null}}`,
		"customView_userViewPreferences_preferences": `{"customView":{"userViewPreferences":null}}`,
	}

	// When
	preferences, err := GetCustomViewUserPreferences(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)
	values, err := GetCustomViewUserPreferenceValues(context.Background(), graphqlClient, "custom-view-id")
	require.NoError(t, err)

	// Then
	require.Equal(t, "custom-view-id", preferences.CustomViewID)
	require.Empty(t, preferences.ID)
	require.Equal(t, "custom-view-id", values.CustomViewID)
	require.False(t, values.HasUserPreferences)
}
