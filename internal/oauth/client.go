// Package oauth exchanges OAuth grants with Linear's token endpoint.
package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/KyaniteHQ/linctl/internal/auth"
)

const (
	defaultTokenEndpoint      = "https://api.linear.app/oauth/token" //nolint:gosec // endpoint URL, not a credential.
	defaultRevocationEndpoint = "https://api.linear.app/oauth/revoke"
)

// ClientConfig configures a Linear OAuth token client.
type ClientConfig struct {
	Endpoint           string
	RevocationEndpoint string
	HTTPClient         *http.Client
	Now                func() time.Time
}

// Client exchanges OAuth grants with Linear.
type Client struct {
	endpoint           string
	revocationEndpoint string
	httpClient         *http.Client
	now                func() time.Time
}

// AuthorizationCodeRequest describes a PKCE authorization-code token exchange.
type AuthorizationCodeRequest struct {
	Code         string
	RedirectURI  string
	ClientID     string
	ClientSecret string
	CodeVerifier string
}

// RefreshTokenRequest describes an OAuth refresh-token exchange.
type RefreshTokenRequest struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
}

// ClientCredentialsRequest describes an app actor token exchange.
type ClientCredentialsRequest struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
}

// RevocationRequest describes an OAuth token revocation request.
type RevocationRequest struct {
	Token         string
	TokenTypeHint string
}

// NewClient returns a Linear OAuth token client.
func NewClient(config ClientConfig) *Client {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}

	return &Client{
		endpoint:           firstNonEmpty(config.Endpoint, defaultTokenEndpoint),
		revocationEndpoint: firstNonEmpty(config.RevocationEndpoint, defaultRevocationEndpoint),
		httpClient:         httpClient,
		now:                now,
	}
}

// ExchangeAuthorizationCode exchanges a PKCE authorization code for tokens.
func (client *Client) ExchangeAuthorizationCode(
	ctx context.Context,
	request AuthorizationCodeRequest,
) (auth.TokenState, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", request.Code)
	form.Set("redirect_uri", request.RedirectURI)
	form.Set("client_id", request.ClientID)
	form.Set("code_verifier", request.CodeVerifier)
	if request.ClientSecret != "" {
		form.Set("client_secret", request.ClientSecret)
	}

	return client.exchange(ctx, form)
}

// RefreshToken exchanges a refresh token and returns Linear's rotated token state.
func (client *Client) RefreshToken(ctx context.Context, request RefreshTokenRequest) (auth.TokenState, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", request.RefreshToken)
	if request.ClientID != "" {
		form.Set("client_id", request.ClientID)
	}
	if request.ClientSecret != "" {
		form.Set("client_secret", request.ClientSecret)
	}

	return client.exchange(ctx, form)
}

// ClientCredentials obtains an app actor access token.
func (client *Client) ClientCredentials(
	ctx context.Context,
	request ClientCredentialsRequest,
) (auth.TokenState, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	if len(request.Scopes) > 0 {
		form.Set("scope", strings.Join(request.Scopes, ","))
	}
	form.Set("client_id", request.ClientID)
	form.Set("client_secret", request.ClientSecret)

	return client.exchange(ctx, form)
}

// RevokeToken revokes an OAuth access or refresh token.
func (client *Client) RevokeToken(ctx context.Context, request RevocationRequest) error {
	form := url.Values{}
	form.Set("token", request.Token)
	if request.TokenTypeHint != "" {
		form.Set("token_type_hint", request.TokenTypeHint)
	}

	_, err := client.postForm(ctx, client.revocationEndpoint, form, "revoke")

	return err
}

// postForm posts a form-encoded body to one OAuth endpoint and returns the
// response bytes. label names the step in a plumbing failure. A non-2xx status
// becomes a token endpoint error keyed by the form's grant type, which is
// absent on the revocation endpoint and maps to the same reauth-required code.
func (client *Client) postForm(
	ctx context.Context,
	endpoint string,
	form url.Values,
	label string,
) ([]byte, error) {
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("create oauth %s request: %w", label, err)
	}
	httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("request oauth %s: %w", label, err)
	}
	body, readErr := io.ReadAll(io.LimitReader(httpResponse.Body, maxTokenResponseBytes))
	closeErr := httpResponse.Body.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read oauth %s response: %w", label, readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close oauth %s response: %w", label, closeErr)
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return nil, tokenEndpointError(form.Get("grant_type"), httpResponse.StatusCode, body)
	}

	return body, nil
}

func (client *Client) exchange(ctx context.Context, form url.Values) (auth.TokenState, error) {
	body, err := client.postForm(ctx, client.endpoint, form, "token")
	if err != nil {
		return auth.TokenState{}, err
	}

	var response tokenEndpointResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return auth.TokenState{}, fmt.Errorf("decode oauth token response: %w", err)
	}

	return response.tokenState(client.now()), nil
}

const maxTokenResponseBytes = 1 << 20

type tokenEndpointResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	Scopes       scopeList `json:"scope"`
}

type tokenEndpointErrorResponse struct {
	Error string `json:"error"`
}

func tokenEndpointError(grantType string, statusCode int, body []byte) error {
	var response tokenEndpointErrorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return auth.NewTokenEndpointError(tokenEndpointErrorCode(grantType, ""), statusCode, "")
	}

	return auth.NewTokenEndpointError(tokenEndpointErrorCode(grantType, response.Error), statusCode, response.Error)
}

func tokenEndpointErrorCode(grantType string, oauthError string) auth.ErrorCode {
	switch oauthError {
	case "invalid_scope", "insufficient_scope":
		return auth.ErrorCodeMissingScope
	}
	switch grantType {
	case "refresh_token":
		return auth.ErrorCodeRefreshFailed
	default:
		return auth.ErrorCodeReauthRequired
	}
}

func (response tokenEndpointResponse) tokenState(now time.Time) auth.TokenState {
	var expiresAt time.Time
	if response.ExpiresIn > 0 {
		expiresAt = now.Add(time.Duration(response.ExpiresIn) * time.Second)
	}

	return auth.NewTokenState(
		response.AccessToken,
		response.RefreshToken,
		response.TokenType,
		expiresAt,
		response.Scopes,
	)
}

type scopeList []string

func (list *scopeList) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*list = nil
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*list = auth.SplitScopes(text)
		return nil
	}

	var values []string
	if err := json.Unmarshal(data, &values); err == nil {
		*list = values
		return nil
	}

	return errors.New("scope must be a string or array")
}

func firstNonEmpty(primary string, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}
