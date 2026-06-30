package client

import (
	"bytes"
	"context"
	"errors"
	"io"
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

func Test_Transport_sends_oauth_access_token_as_bearer_authorization(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("authorization header = %q", got)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`{"data":{"viewer":{"id":"user-id"}}}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint: server.URL,
		Token:    OAuthAccessToken("test-token"),
		Timeout:  2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.NoError(t, err)
}

func Test_Transport_omits_authorization_header_for_empty_oauth_token(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "" {
			t.Errorf("authorization header = %q", got)
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`{"data":{"viewer":{"id":"user-id"}}}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint: server.URL,
		Token:    OAuthAccessToken(""),
		Timeout:  2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.NoError(t, err)
}

func Test_Transport_returns_graphql_error_when_response_contains_errors(t *testing.T) {
	// Given
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer test-token" {
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
		Token:    OAuthAccessToken("test-token"),
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

func Test_Transport_returns_typed_auth_error_for_401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "Bearer expired-token" {
			t.Errorf("authorization header = %q", got)
		}
		writer.WriteHeader(http.StatusUnauthorized)
		_, err := writer.Write([]byte(`{"errors":[{"message":"secret auth failure","extensions":{"code":"UNAUTHENTICATED"}}]}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		Endpoint: server.URL,
		Token:    OAuthAccessToken("expired-token"),
		Timeout:  2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrAuthFailed)
	require.Contains(t, err.Error(), "graphql http status 401 code=UNAUTHENTICATED")
	require.NotContains(t, err.Error(), "secret auth failure")
	require.NotContains(t, err.Error(), "expired-token")
}

func Test_Transport_retries_429_with_retry_after_when_present(t *testing.T) {
	// Given
	logs := bytes.Buffer{}
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
		DiagnosticWriter: &logs,
		Endpoint:         server.URL,
		Token:            OAuthAccessToken("test-token"),
		Timeout:          2 * time.Second,
		MaxRetries:       1,
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
	require.Contains(t, logs.String(), "graphql_request attempt=1")
	require.Contains(t, logs.String(), "graphql_response attempt=1 status=429")
	require.Contains(t, logs.String(), "graphql_retry attempt=1 status=429")
	require.Contains(t, logs.String(), "graphql_request_ok attempt=2 status=200")
	require.NotContains(t, logs.String(), "test-token")
	require.NotContains(t, logs.String(), "user-id")
}

func Test_Transport_retries_ratelimited_400_then_succeeds(t *testing.T) {
	// Given Linear signals a rate limit with HTTP 400 + RATELIMITED, not 429.
	logs := bytes.Buffer{}
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		requests++
		writer.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			writer.WriteHeader(http.StatusBadRequest)
			_, err := writer.Write([]byte(`{"errors":[{"message":"rate limit exceeded","extensions":{"code":"RATELIMITED"}}]}`))
			if err != nil {
				t.Errorf("write response: %v", err)
			}
			return
		}
		_, err := writer.Write([]byte(`{"data":{"viewer":{"id":"user-id"}}}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		DiagnosticWriter: &logs,
		Endpoint:         server.URL,
		Token:            OAuthAccessToken("test-token"),
		Timeout:          2 * time.Second,
		MaxRetries:       1,
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
	require.Contains(t, logs.String(), "graphql_retry attempt=1 status=400")
	require.Contains(t, logs.String(), "graphql_request_ok attempt=2 status=200")
}

func Test_Transport_returns_rate_limited_error_after_exhausting_retries(t *testing.T) {
	// Given a server that always returns the RATELIMITED 400.
	logs := bytes.Buffer{}
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		requests++
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write([]byte(`{"errors":[{"message":"rate limit exceeded","extensions":{"code":"RATELIMITED"}}]}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		DiagnosticWriter: &logs,
		Endpoint:         server.URL,
		Token:            OAuthAccessToken("test-token"),
		Timeout:          2 * time.Second,
		MaxRetries:       1,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.Error(t, err)
	require.ErrorIs(t, err, ErrRateLimited)
	require.Contains(t, err.Error(), "graphql http status 400 code=RATELIMITED")
	require.NotContains(t, err.Error(), "rate limit exceeded")
	require.Equal(t, 2, requests)
	require.Contains(t, logs.String(), "graphql_rate_limited attempt=2 status=400")
	require.NotContains(t, logs.String(), "test-token")
}

func Test_Transport_returns_error_when_context_timeout_expires(t *testing.T) {
	// Given
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
		// Hold the response open so the nanosecond client timeout is what fails
		// the request. Also release on cleanup: Windows does not reliably cancel
		// the server-side request context when the client disconnects, which
		// would otherwise hang server.Close until the test deadline.
		select {
		case <-request.Context().Done():
		case <-release:
		}
	}))
	defer server.Close()
	defer close(release)
	transport := NewTransport(TransportConfig{
		Endpoint: server.URL,
		Token:    OAuthAccessToken("test-token"),
		Timeout:  time.Nanosecond,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "request failed")
}

func Test_Transport_logs_decode_failures_without_response_body(t *testing.T) {
	// Given
	logs := bytes.Buffer{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`{"errors":[{"message":"sensitive failure detail"}]}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		DiagnosticWriter: &logs,
		Endpoint:         server.URL,
		Token:            OAuthAccessToken("test-token"),
		Timeout:          2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.Error(t, err)
	require.Contains(t, logs.String(), "graphql_decode_failed attempt=1 status=200")
	require.NotContains(t, logs.String(), "sensitive failure detail")
	require.NotContains(t, logs.String(), "test-token")
}

func Test_Transport_ignores_diagnostic_writer_failures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`{"data":{"viewer":{"id":"user-id"}}}`))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		DiagnosticWriter: failingDiagnosticWriter{},
		Endpoint:         server.URL,
		Token:            OAuthAccessToken("test-token"),
		Timeout:          2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	require.NoError(t, err)
}

func Test_Transport_returns_errors_for_request_and_body_failures(t *testing.T) {
	t.Run("unmarshalable variables", func(t *testing.T) {
		logs := bytes.Buffer{}
		transport := NewTransport(TransportConfig{
			DiagnosticWriter: &logs,
			Timeout:          2 * time.Second,
		})
		response := graphql.Response{Data: &testGraphQLData{}}

		err := transport.MakeRequest(context.Background(), &graphql.Request{
			Query:     "query Test { viewer { id } }",
			Variables: map[string]any{"bad": make(chan int)},
		}, &response)

		require.Error(t, err)
		require.Contains(t, err.Error(), "encode graphql request")
		require.Contains(t, logs.String(), "graphql_encode_failed")
		require.NotContains(t, logs.String(), "query Test")
	})

	t.Run("invalid endpoint", func(t *testing.T) {
		transport := NewTransport(TransportConfig{
			Endpoint: "://bad-url",
			Timeout:  2 * time.Second,
		})
		response := graphql.Response{Data: &testGraphQLData{}}

		err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

		require.Error(t, err)
		require.Contains(t, err.Error(), "create graphql request")
	})

	t.Run("body read error", func(t *testing.T) {
		transport := NewTransport(TransportConfig{
			Client: &http.Client{Transport: bodyFailureTransport{
				body: failingBody{readErr: errors.New("read failed")},
			}},
			Timeout: 2 * time.Second,
		})
		response := graphql.Response{Data: &testGraphQLData{}}

		err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

		require.Error(t, err)
		require.Contains(t, err.Error(), "read response body")
	})

	t.Run("body close error", func(t *testing.T) {
		transport := NewTransport(TransportConfig{
			Client: &http.Client{Transport: bodyFailureTransport{
				body: failingBody{closeErr: errors.New("close failed")},
			}},
			Timeout: 2 * time.Second,
		})
		response := graphql.Response{Data: &testGraphQLData{}}

		err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

		require.Error(t, err)
		require.Contains(t, err.Error(), "close response body")
	})

	t.Run("http error body is redacted", func(t *testing.T) {
		transport := NewTransport(TransportConfig{
			Client: &http.Client{Transport: bodyFailureTransport{
				statusCode: http.StatusBadGateway,
				body:       io.NopCloser(bytes.NewReader([]byte("sensitive upstream body"))),
			}},
			Timeout: 2 * time.Second,
		})
		response := graphql.Response{Data: &testGraphQLData{}}

		err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

		require.Error(t, err)
		require.Contains(t, err.Error(), "graphql http status 502")
		require.NotContains(t, err.Error(), "sensitive upstream body")
	})
}

func Test_Transport_logs_terminal_http_failures_without_response_body(t *testing.T) {
	// Given
	logs := bytes.Buffer{}
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("authorization header = %q", request.Header.Get("Authorization"))
		}
		writer.WriteHeader(http.StatusInternalServerError)
		_, err := writer.Write([]byte("sensitive outage detail"))
		if err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()
	transport := NewTransport(TransportConfig{
		DiagnosticWriter: &logs,
		Endpoint:         server.URL,
		Token:            OAuthAccessToken("test-token"),
		Timeout:          2 * time.Second,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	// When
	err := transport.MakeRequest(context.Background(), &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	// Then
	require.Error(t, err)
	require.NotContains(t, err.Error(), "sensitive outage detail")
	require.Contains(t, logs.String(), "graphql_decode_failed attempt=1 status=500")
	require.NotContains(t, logs.String(), "sensitive outage detail")
	require.NotContains(t, logs.String(), "test-token")
}

func Test_NewTransport_uses_project_owned_default_client(t *testing.T) {
	transport := NewTransport(TransportConfig{Timeout: 2 * time.Second})

	require.NotSame(t, http.DefaultClient, transport.httpClient)
	require.Equal(t, 2*time.Second, transport.httpClient.Timeout)
}

func Test_defaultHTTPClient_falls_back_when_default_transport_is_custom(t *testing.T) {
	originalTransport := http.DefaultTransport
	customTransport := staticResponseTransport{}
	http.DefaultTransport = customTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	client := defaultHTTPClient(2 * time.Second)

	require.Equal(t, 2*time.Second, client.Timeout)
	require.Equal(t, customTransport, client.Transport)
}

func Test_firstGraphQLErrorCode_returns_empty_when_errors_have_no_code(t *testing.T) {
	code := firstGraphQLErrorCode([]byte(`{"errors":[{"message":"body stays private"}]}`))

	require.Empty(t, code)
}

func Test_Transport_returns_error_when_retry_wait_is_canceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	logs := bytes.Buffer{}
	transport := NewTransport(TransportConfig{
		Client: &http.Client{Transport: retryCancelTransport{
			body: cancelAfterReadBody{cancel: cancel},
		}},
		DiagnosticWriter: &logs,
		Timeout:          2 * time.Second,
		MaxRetries:       1,
	})
	response := graphql.Response{Data: &testGraphQLData{}}

	err := transport.MakeRequest(ctx, &graphql.Request{Query: "query Test { viewer { id } }"}, &response)

	require.Error(t, err)
	require.Contains(t, err.Error(), "wait for retry")
	require.Contains(t, logs.String(), "graphql_retry_failed attempt=1")
}

func Benchmark_Transport_make_request_diagnostics(b *testing.B) {
	cases := []struct {
		name             string
		diagnosticWriter io.Writer
	}{
		{name: "disabled", diagnosticWriter: nil},
		{name: "enabled", diagnosticWriter: io.Discard},
	}

	for _, testCase := range cases {
		b.Run(testCase.name, func(b *testing.B) {
			transport := NewTransport(TransportConfig{
				Client:           &http.Client{Transport: staticResponseTransport{}},
				DiagnosticWriter: testCase.diagnosticWriter,
				Timeout:          2 * time.Second,
			})

			b.ReportAllocs()
			for range b.N {
				response := graphql.Response{Data: &testGraphQLData{}}
				err := transport.MakeRequest(
					context.Background(),
					&graphql.Request{Query: "query Test { viewer { id } }"},
					&response,
				)
				if err != nil {
					b.Fatalf("make request: %v", err)
				}
			}
		})
	}
}

type retryCancelTransport struct {
	body io.ReadCloser
}

func (transport retryCancelTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"1"}},
		Body:       transport.body,
	}, nil
}

type cancelAfterReadBody struct {
	cancel context.CancelFunc
}

func (body cancelAfterReadBody) Read(_ []byte) (int, error) {
	body.cancel()
	return 0, io.EOF
}

func (body cancelAfterReadBody) Close() error {
	return nil
}

type failingDiagnosticWriter struct{}

func (writer failingDiagnosticWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("diagnostic sink closed")
}

type staticResponseTransport struct{}

func (transport staticResponseTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"data":{"viewer":{"id":"user-id"}}}`))),
	}, nil
}

type bodyFailureTransport struct {
	statusCode int
	body       io.ReadCloser
}

func (transport bodyFailureTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	statusCode := transport.statusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{},
		Body:       transport.body,
	}, nil
}

type failingBody struct {
	readErr  error
	closeErr error
}

func (body failingBody) Read(_ []byte) (int, error) {
	if body.readErr != nil {
		return 0, body.readErr
	}

	return 0, io.EOF
}

func (body failingBody) Close() error {
	return body.closeErr
}
