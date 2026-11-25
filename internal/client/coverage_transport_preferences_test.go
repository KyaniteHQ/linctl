package client

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Transport_rejects_a_success_status_with_a_body_that_is_not_json(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writeTestBody(t, writer, "not json")
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{Endpoint: server.URL, Timeout: 2 * time.Second})

	err := requestViewer(transport, &testGraphQLData{})

	require.ErrorContains(t, err, "decode graphql response")
}

func Test_Transport_does_not_retry_responses_that_are_not_rate_limits(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		want       string
	}{
		{
			name:       "400 with another error code",
			statusCode: http.StatusBadRequest,
			body:       `{"errors":[{"message":"bad input","extensions":{"code":"BAD_USER_INPUT"}}]}`,
			want:       "graphql http status 400 code=BAD_USER_INPUT",
		},
		{
			name:       "400 with a body that is not json",
			statusCode: http.StatusBadRequest,
			body:       "not json",
			want:       "graphql http status 400",
		},
		{
			name:       "400 with valid json and no code",
			statusCode: http.StatusBadRequest,
			body:       `{"errors":[{"message":"private detail"}]}`,
			want:       "graphql http status 400",
		},
		{
			name:       "500 with a RATELIMITED code",
			statusCode: http.StatusInternalServerError,
			body:       `{"errors":[{"message":"rate limit exceeded","extensions":{"code":"RATELIMITED"}}]}`,
			want:       "graphql http status 500 code=RATELIMITED",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				requests.Add(1)
				writer.WriteHeader(test.statusCode)
				writeTestBody(t, writer, test.body)
			}))
			defer server.Close()
			transport := NewTransport(TransportConfig{Endpoint: server.URL, Timeout: 2 * time.Second, MaxRetries: 1})

			err := requestViewer(transport, &testGraphQLData{})

			require.EqualError(t, err, test.want)
			require.NotContains(t, err.Error(), "private detail")
			require.NotErrorIs(t, err, ErrRateLimited)
			require.Equal(t, int32(1), requests.Load())
		})
	}
}

func Test_Transport_does_not_retry_at_once_when_retry_after_overflows_a_duration(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) == 1 {
			writer.Header().Set("Retry-After", strconv.Itoa(math.MaxInt))
			writer.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writeTestBody(t, writer, `{"data":{"viewer":{"id":"user-id"}}}`)
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint:   server.URL,
		Timeout:    300 * time.Millisecond,
		MaxRetries: 1,
	})

	err := requestViewer(transport, &testGraphQLData{})

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.ErrorContains(t, err, "wait for retry")
	require.Equal(t, int32(1), requests.Load())
}

func Test_NewTransport_without_a_timeout_still_completes_requests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writeTestBody(t, writer, `{"data":{"viewer":{"id":"user-id"}}}`)
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{Endpoint: server.URL})
	data := &testGraphQLData{}

	err := requestViewer(transport, data)

	require.NoError(t, err)
	require.Equal(t, "user-id", data.Viewer.ID)
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
