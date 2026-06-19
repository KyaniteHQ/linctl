package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type testGraphQLData struct {
	Viewer struct {
		ID string `json:"id"`
	} `json:"viewer"`
}

func Test_Transport_returns_graphql_error_when_response_contains_errors(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "test-token" {
			t.Errorf("authorization header = %q", request.Header.Get("Authorization"))
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`{"errors":[{"message":"bad query","extensions":{"code":"BAD_USER_INPUT"}}]}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint: server.URL,
		Token:    PersonalAPIToken("test-token"),
		Timeout:  2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrGraphQL)
	require.Contains(t, err.Error(), "bad query")
	require.Contains(t, err.Error(), "BAD_USER_INPUT")
}

func Test_Transport_retries_429_with_retry_after_when_present(t *testing.T) {
	// Given
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		requests++
		if requests == 1 {
			writer.Header().Set("Retry-After", "0")
			writer.WriteHeader(http.StatusTooManyRequests)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`{"data":{"viewer":{"id":"user-id"}}}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint:   server.URL,
		Token:      PersonalAPIToken("test-token"),
		Timeout:    2 * time.Second,
		MaxRetries: 1,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.NoError(t, err)
	require.Equal(t, 2, requests)
	data, ok := response.Data.(*testGraphQLData)
	require.True(t, ok)
	require.Equal(t, "user-id", data.Viewer.ID)
}

func Test_Transport_returns_error_when_context_timeout_expires(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint: server.URL,
		Token:    PersonalAPIToken("test-token"),
		Timeout:  time.Nanosecond,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "request failed")
}
