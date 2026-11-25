package oauth_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

type tokenRequestSnapshot struct {
	Method      string
	ContentType string
	Auth        string
	Form        url.Values
	ParseErr    error
	HasBasic    bool
	BasicUser   string
	BasicPass   string
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

type responseBody struct {
	reader   *strings.Reader
	readErr  error
	closeErr error
}

func (body *responseBody) Read(data []byte) (int, error) {
	if body.readErr != nil {
		return 0, body.readErr
	}

	return body.reader.Read(data)
}

func (body *responseBody) Close() error {
	return body.closeErr
}

func writeResponseBody(writer http.ResponseWriter, body string) error {
	_, err := writer.Write([]byte(body))
	return err
}

func Test_Client_exchanges_authorization_code_with_pkce_form(t *testing.T) {
	t.Parallel()
	requests := make(chan tokenRequestSnapshot, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		parseErr := request.ParseForm()
		basicUser, basicPass, hasBasic := request.BasicAuth()
		requests <- tokenRequestSnapshot{
			Method:      request.Method,
			ContentType: request.Header.Get("Content-Type"),
			Auth:        request.Header.Get("Authorization"),
			Form:        request.PostForm,
			ParseErr:    parseErr,
			HasBasic:    hasBasic,
			BasicUser:   basicUser,
			BasicPass:   basicPass,
		}
		writer.Header().Set("Content-Type", "application/json")
		if err := writeResponseBody(writer, `{
				"access_token": "oauth-access-token",
				"refresh_token": "rotated-refresh-token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"scope": "read,write"
			}`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	now := time.Date(2026, 6, 29, 12, 0, 0, 0, time.UTC)
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
		Now: func() time.Time {
			return now
		},
	})

	token, err := client.ExchangeAuthorizationCode(context.Background(), oauth.AuthorizationCodeRequest{
		Code:         "authorization-code",
		RedirectURI:  "http://127.0.0.1:8080/callback",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		CodeVerifier: "pkce-code-verifier",
	})

	require.NoError(t, err)
	request := <-requests
	require.Equal(t, http.MethodPost, request.Method)
	require.NoError(t, request.ParseErr)
	require.Equal(t, "application/x-www-form-urlencoded", request.ContentType)
	require.Empty(t, request.Auth)
	require.Equal(t, "authorization_code", request.Form.Get("grant_type"))
	require.Equal(t, "authorization-code", request.Form.Get("code"))
	require.Equal(t, "http://127.0.0.1:8080/callback", request.Form.Get("redirect_uri"))
	require.Equal(t, "client-id", request.Form.Get("client_id"))
	require.Equal(t, "pkce-code-verifier", request.Form.Get("code_verifier"))
	require.Equal(t, "client-secret", request.Form.Get("client_secret"))
	require.Equal(t, "oauth-access-token", token.AccessToken)
	require.Equal(t, "rotated-refresh-token", token.RefreshToken)
	require.Equal(t, "Bearer", token.TokenType)
	require.Equal(t, []string{"read", "write"}, token.Scopes)
	require.NotNil(t, token.ExpiresAt)
	require.Equal(t, now.Add(time.Hour), *token.ExpiresAt)
}

func Test_Client_refreshes_rotated_token_state_with_form_client_auth(t *testing.T) {
	t.Parallel()
	requests := make(chan tokenRequestSnapshot, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		parseErr := request.ParseForm()
		basicUser, basicPass, hasBasic := request.BasicAuth()
		requests <- tokenRequestSnapshot{
			Method:      request.Method,
			ContentType: request.Header.Get("Content-Type"),
			Auth:        request.Header.Get("Authorization"),
			Form:        request.PostForm,
			ParseErr:    parseErr,
			HasBasic:    hasBasic,
			BasicUser:   basicUser,
			BasicPass:   basicPass,
		}
		writer.Header().Set("Content-Type", "application/json")
		if err := writeResponseBody(writer, `{
				"access_token": "new-access-token",
				"refresh_token": "new-refresh-token",
				"token_type": "Bearer",
				"expires_in": 1800,
				"scope": ["read", "write"]
			}`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	now := time.Date(2026, 6, 29, 13, 0, 0, 0, time.UTC)
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
		Now: func() time.Time {
			return now
		},
	})

	token, err := client.RefreshToken(context.Background(), oauth.RefreshTokenRequest{
		RefreshToken: "old-refresh-token",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	})

	require.NoError(t, err)
	request := <-requests
	require.Equal(t, http.MethodPost, request.Method)
	require.NoError(t, request.ParseErr)
	require.Equal(t, "application/x-www-form-urlencoded", request.ContentType)
	require.False(t, request.HasBasic)
	require.Equal(t, "refresh_token", request.Form.Get("grant_type"))
	require.Equal(t, "old-refresh-token", request.Form.Get("refresh_token"))
	require.Equal(t, "client-id", request.Form.Get("client_id"))
	require.Equal(t, "client-secret", request.Form.Get("client_secret"))
	require.Equal(t, "new-access-token", token.AccessToken)
	require.Equal(t, "new-refresh-token", token.RefreshToken)
	require.Equal(t, "Bearer", token.TokenType)
	require.Equal(t, []string{"read", "write"}, token.Scopes)
	require.NotNil(t, token.ExpiresAt)
	require.Equal(t, now.Add(30*time.Minute), *token.ExpiresAt)
}

func Test_Client_obtains_client_credentials_app_actor_token(t *testing.T) {
	t.Parallel()
	requests := make(chan tokenRequestSnapshot, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		parseErr := request.ParseForm()
		basicUser, basicPass, hasBasic := request.BasicAuth()
		requests <- tokenRequestSnapshot{
			Method:      request.Method,
			ContentType: request.Header.Get("Content-Type"),
			Auth:        request.Header.Get("Authorization"),
			Form:        request.PostForm,
			ParseErr:    parseErr,
			HasBasic:    hasBasic,
			BasicUser:   basicUser,
			BasicPass:   basicPass,
		}
		writer.Header().Set("Content-Type", "application/json")
		if err := writeResponseBody(writer, `{
				"access_token": "app-actor-access-token",
				"token_type": "Bearer",
				"expires_in": 7200,
				"scope": "read write"
			}`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	now := time.Date(2026, 6, 29, 14, 0, 0, 0, time.UTC)
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
		Now: func() time.Time {
			return now
		},
	})

	token, err := client.ClientCredentials(context.Background(), oauth.ClientCredentialsRequest{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"read", "write"},
	})

	require.NoError(t, err)
	request := <-requests
	require.Equal(t, http.MethodPost, request.Method)
	require.NoError(t, request.ParseErr)
	require.Equal(t, "application/x-www-form-urlencoded", request.ContentType)
	require.False(t, request.HasBasic)
	require.Empty(t, request.Auth)
	require.Equal(t, "client_credentials", request.Form.Get("grant_type"))
	require.Equal(t, "read,write", request.Form.Get("scope"))
	require.Equal(t, "client-id", request.Form.Get("client_id"))
	require.Equal(t, "client-secret", request.Form.Get("client_secret"))
	require.Equal(t, "app-actor-access-token", token.AccessToken)
	require.Empty(t, token.RefreshToken)
	require.Equal(t, "Bearer", token.TokenType)
	require.Equal(t, []string{"read", "write"}, token.Scopes)
	require.NotNil(t, token.ExpiresAt)
	require.Equal(t, now.Add(2*time.Hour), *token.ExpiresAt)
}

func Test_Client_revokes_token_with_token_form_field(t *testing.T) {
	t.Parallel()
	requests := make(chan tokenRequestSnapshot, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		parseErr := request.ParseForm()
		basicUser, basicPass, hasBasic := request.BasicAuth()
		requests <- tokenRequestSnapshot{
			Method:      request.Method,
			ContentType: request.Header.Get("Content-Type"),
			Auth:        request.Header.Get("Authorization"),
			Form:        request.PostForm,
			ParseErr:    parseErr,
			HasBasic:    hasBasic,
			BasicUser:   basicUser,
			BasicPass:   basicPass,
		}
		writer.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	client := oauth.NewClient(oauth.ClientConfig{
		RevocationEndpoint: server.URL,
		HTTPClient:         server.Client(),
	})

	err := client.RevokeToken(context.Background(), oauth.RevocationRequest{
		Token:         "refresh-token",
		TokenTypeHint: "refresh_token",
	})

	require.NoError(t, err)
	request := <-requests
	require.Equal(t, http.MethodPost, request.Method)
	require.NoError(t, request.ParseErr)
	require.Equal(t, "application/x-www-form-urlencoded", request.ContentType)
	require.False(t, request.HasBasic)
	require.Empty(t, request.Auth)
	require.Equal(t, "refresh-token", request.Form.Get("token"))
	require.Equal(t, "refresh_token", request.Form.Get("token_type_hint"))
	require.Empty(t, request.Form.Get("access_token"))
	require.Empty(t, request.Form.Get("refresh_token"))
}

func Test_Client_revoke_error_redacts_secret_values(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		if err := writeResponseBody(writer, `{
				"error": "invalid_grant",
				"error_description": "token=refresh-token client_secret=client-secret"
			}`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	client := oauth.NewClient(oauth.ClientConfig{
		RevocationEndpoint: server.URL,
		HTTPClient:         server.Client(),
	})

	err := client.RevokeToken(context.Background(), oauth.RevocationRequest{
		Token:         "refresh-token",
		TokenTypeHint: "refresh_token",
	})

	require.Error(t, err)
	var tokenErr *auth.TokenEndpointError
	require.ErrorAs(t, err, &tokenErr)
	require.Equal(t, auth.ErrorCodeReauthRequired, tokenErr.Code)
	require.NotContains(t, err.Error(), "refresh-token")
	require.NotContains(t, err.Error(), "client-secret")
}

func Test_Client_token_endpoint_error_has_code_and_redacts_secrets(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		if err := writeResponseBody(writer, `{
				"error": "invalid_grant",
				"error_description": "client_secret=client-secret access_token=oauth-access-token refresh_token=rotated-refresh-token secret-from-response"
			}`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
	})

	_, err := client.RefreshToken(context.Background(), oauth.RefreshTokenRequest{
		RefreshToken: "old-refresh-token",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	})

	require.Error(t, err)
	var tokenErr *auth.TokenEndpointError
	require.ErrorAs(t, err, &tokenErr)
	require.Equal(t, auth.ErrorCodeRefreshFailed, tokenErr.Code)
	errorText := err.Error()
	require.Contains(t, errorText, "AUTH_REFRESH_FAILED")
	require.Contains(t, errorText, "http status 400")
	for _, forbidden := range []string{
		"oauth-access-token",
		"old-refresh-token",
		"rotated-refresh-token",
		"client-secret",
		"secret-from-response",
		"access_token",
		"refresh_token",
		"client_secret",
	} {
		require.NotContains(t, errorText, forbidden)
	}
}

func Test_Client_token_endpoint_invalid_scope_maps_missing_scope(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		if err := writeResponseBody(writer, `{"error":"invalid_scope","error_description":"scope denied"}`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
	})

	_, err := client.ClientCredentials(context.Background(), oauth.ClientCredentialsRequest{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		Scopes:       []string{"admin"},
	})

	require.Error(t, err)
	var tokenErr *auth.TokenEndpointError
	require.ErrorAs(t, err, &tokenErr)
	require.Equal(t, auth.ErrorCodeMissingScope, tokenErr.Code)
	require.Equal(t, "invalid_scope", tokenErr.OAuthError)
	require.NotContains(t, err.Error(), "scope denied")
}

func Test_Client_exchange_request_creation_error(t *testing.T) {
	t.Parallel()
	client := oauth.NewClient(oauth.ClientConfig{Endpoint: "http://[::1"})

	_, err := client.ExchangeAuthorizationCode(context.Background(), oauth.AuthorizationCodeRequest{})

	require.ErrorContains(t, err, "create oauth token request")
}

func Test_Client_exchange_http_error(t *testing.T) {
	t.Parallel()
	requestErr := errors.New("transport unavailable")
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint: "https://oauth.test/token",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, requestErr
		})},
	})

	_, err := client.RefreshToken(context.Background(), oauth.RefreshTokenRequest{RefreshToken: "refresh-token"})

	require.ErrorIs(t, err, requestErr)
	require.ErrorContains(t, err, "request oauth token")
}

func Test_Client_exchange_response_body_errors(t *testing.T) {
	t.Parallel()
	readErr := errors.New("read failed")
	closeErr := errors.New("close failed")
	for _, test := range []struct {
		name      string
		body      io.ReadCloser
		wantError error
		wantText  string
	}{
		{
			name:      "read error",
			body:      &responseBody{readErr: readErr},
			wantError: readErr,
			wantText:  "read oauth token response",
		},
		{
			name:      "close error",
			body:      &responseBody{reader: strings.NewReader(`{"access_token":"access-token"}`), closeErr: closeErr},
			wantError: closeErr,
			wantText:  "close oauth token response",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := oauth.NewClient(oauth.ClientConfig{
				Endpoint: "https://oauth.test/token",
				HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       test.body,
					}, nil
				})},
			})

			_, err := client.ClientCredentials(context.Background(), oauth.ClientCredentialsRequest{})

			require.ErrorIs(t, err, test.wantError)
			require.ErrorContains(t, err, test.wantText)
		})
	}
}

func Test_Client_exchange_rejects_malformed_token_json(t *testing.T) {
	t.Parallel()
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint: "https://oauth.test/token",
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":`)),
			}, nil
		})},
	})

	_, err := client.ClientCredentials(context.Background(), oauth.ClientCredentialsRequest{})

	require.ErrorContains(t, err, "decode oauth token response")
}

func Test_Client_token_endpoint_malformed_error_body_uses_empty_oauth_error(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		if err := writeResponseBody(writer, `{"error":`); err != nil {
			return
		}
	}))
	t.Cleanup(server.Close)
	client := oauth.NewClient(oauth.ClientConfig{
		Endpoint:   server.URL,
		HTTPClient: server.Client(),
	})

	_, err := client.RefreshToken(context.Background(), oauth.RefreshTokenRequest{RefreshToken: "refresh-token"})

	require.Error(t, err)
	var tokenErr *auth.TokenEndpointError
	require.ErrorAs(t, err, &tokenErr)
	require.Equal(t, auth.ErrorCodeRefreshFailed, tokenErr.Code)
	require.Empty(t, tokenErr.OAuthError)
}

func Test_NewClient_defaults_to_linear_token_and_revocation_endpoints(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := oauth.NewClient(oauth.ClientConfig{})

	_, exchangeErr := client.ExchangeAuthorizationCode(ctx, oauth.AuthorizationCodeRequest{})
	revokeErr := client.RevokeToken(ctx, oauth.RevocationRequest{Token: "refresh-token"})

	var exchangeURLErr *url.Error
	require.ErrorAs(t, exchangeErr, &exchangeURLErr)
	require.Equal(t, "https://api.linear.app/oauth/token", exchangeURLErr.URL)
	require.ErrorIs(t, exchangeErr, context.Canceled)
	var revokeURLErr *url.Error
	require.ErrorAs(t, revokeErr, &revokeURLErr)
	require.Equal(t, "https://api.linear.app/oauth/revoke", revokeURLErr.URL)
	require.ErrorIs(t, revokeErr, context.Canceled)
}

func Test_Client_without_clock_expires_tokens_from_the_current_time(t *testing.T) {
	t.Parallel()
	client := oauth.NewClient(oauth.ClientConfig{
		HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"access_token":"access-token","expires_in":3600}`)),
			}, nil
		})},
	})
	before := time.Now()

	token, err := client.ClientCredentials(context.Background(), oauth.ClientCredentialsRequest{})

	after := time.Now()
	require.NoError(t, err)
	require.NotNil(t, token.ExpiresAt)
	require.WithinRange(t, *token.ExpiresAt, before.Add(time.Hour), after.Add(time.Hour))
}

func Test_Client_token_response_scope_null_and_non_string_values(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name      string
		body      string
		wantError string
	}{
		{
			name: "null scope leaves scopes unset",
			body: `{"access_token":"access-token","scope":null}`,
		},
		{
			name:      "numeric scope is rejected",
			body:      `{"access_token":"access-token","scope":42}`,
			wantError: "decode oauth token response: scope must be a string or array",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			client := oauth.NewClient(oauth.ClientConfig{
				HTTPClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(test.body)),
					}, nil
				})},
			})

			token, err := client.ClientCredentials(context.Background(), oauth.ClientCredentialsRequest{})

			if test.wantError != "" {
				require.EqualError(t, err, test.wantError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, "access-token", token.AccessToken)
			require.Nil(t, token.Scopes)
		})
	}
}
