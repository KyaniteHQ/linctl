// Package client contains Linear GraphQL client primitives.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// ErrGraphQL marks a GraphQL errors[] response.
var ErrGraphQL = errors.New("graphql error")

// ErrMutationFailed marks a mutation payload without success and entity id.
var ErrMutationFailed = errors.New("mutation failed")

// AuthToken formats the Linear Authorization header.
type AuthToken struct {
	authorization string
}

// TransportConfig configures the Linear GraphQL transport.
type TransportConfig struct {
	Client           *http.Client
	DiagnosticWriter io.Writer
	Endpoint         string
	Token            AuthToken
	Timeout          time.Duration
	MaxRetries       int
}

// Transport implements genqlient's GraphQL client interface.
type Transport struct {
	httpClient       *http.Client
	diagnosticWriter io.Writer
	endpoint         string
	token            AuthToken
	timeout          time.Duration
	maxRetries       int
}

// PersonalAPIToken sends a raw Linear personal API key.
func PersonalAPIToken(value string) AuthToken {
	return AuthToken{
		authorization: value,
	}
}

// NewTransport creates a Linear GraphQL transport.
func NewTransport(config TransportConfig) *Transport {
	httpClient := config.Client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Transport{
		httpClient:       httpClient,
		diagnosticWriter: config.DiagnosticWriter,
		endpoint:         firstNonEmpty(config.Endpoint, "https://api.linear.app/graphql"),
		token:            config.Token,
		timeout:          defaultDuration(config.Timeout, 30*time.Second),
		maxRetries:       defaultRetries(config.MaxRetries),
	}
}

// MakeRequest sends a GraphQL request and decodes a GraphQL response.
func (transport *Transport) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	payload, err := json.Marshal(request)
	if err != nil {
		transport.log("graphql_encode_failed error=%q", err.Error())
		return fmt.Errorf("encode graphql request: %w", err)
	}

	for attempt := 0; ; attempt++ {
		transport.log("graphql_request attempt=%d", attempt+1)
		body, statusCode, header, err := transport.send(ctx, payload)
		if err != nil {
			transport.log("graphql_request_failed attempt=%d error=%q", attempt+1, err.Error())
			return err
		}
		transport.log("graphql_response attempt=%d status=%d", attempt+1, statusCode)
		retry, err := waitForRateLimitRetry(ctx, statusCode, header, attempt, transport.maxRetries)
		if err != nil {
			transport.log("graphql_retry_failed attempt=%d error=%q", attempt+1, err.Error())
			return err
		}
		if retry {
			transport.log("graphql_retry attempt=%d status=%d", attempt+1, statusCode)
			continue
		}
		if err := decodeGraphQLResponse(body, statusCode, response); err != nil {
			transport.log("graphql_decode_failed attempt=%d status=%d", attempt+1, statusCode)
			return err
		}
		transport.log("graphql_request_ok attempt=%d status=%d", attempt+1, statusCode)
		return nil
	}
}

func (transport *Transport) log(format string, args ...any) {
	if transport.diagnosticWriter == nil {
		return
	}

	_, err := fmt.Fprintf(transport.diagnosticWriter, format+"\n", args...)
	if err != nil {
		return
	}
}

func waitForRateLimitRetry(
	ctx context.Context,
	statusCode int,
	header http.Header,
	attempt int,
	maxRetries int,
) (bool, error) {
	if statusCode != http.StatusTooManyRequests || attempt >= maxRetries {
		return false, nil
	}
	if err := waitForRetry(ctx, retryDelay(header.Get("Retry-After"), attempt)); err != nil {
		return false, err
	}

	return true, nil
}

func decodeGraphQLResponse(body []byte, statusCode int, response *graphql.Response) error {
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("graphql http status %d: %s", statusCode, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("decode graphql response: %w", err)
	}
	if len(response.Errors) > 0 {
		return formatGraphQLErrors(response.Errors)
	}

	return nil
}

func (transport *Transport) send(ctx context.Context, payload []byte) ([]byte, int, http.Header, error) {
	requestCtx, cancel := context.WithTimeout(ctx, transport.timeout)
	defer cancel()

	httpRequest, err := http.NewRequestWithContext(
		requestCtx,
		http.MethodPost,
		transport.endpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("create graphql request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	if transport.token.authorization != "" {
		httpRequest.Header.Set("Authorization", transport.token.authorization)
	}

	httpResponse, err := transport.httpClient.Do(httpRequest)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("request failed: %w", err)
	}
	body, readErr := io.ReadAll(httpResponse.Body)
	closeErr := httpResponse.Body.Close()
	if readErr != nil {
		return nil, 0, nil, fmt.Errorf("read response body: %w", readErr)
	}
	if closeErr != nil {
		return nil, 0, nil, fmt.Errorf("close response body: %w", closeErr)
	}

	return body, httpResponse.StatusCode, httpResponse.Header, nil
}

func retryDelay(retryAfter string, attempt int) time.Duration {
	if retryAfter != "" {
		seconds, err := strconv.Atoi(retryAfter)
		if err == nil {
			return time.Duration(seconds) * time.Second
		}
	}

	return time.Duration(attempt+1) * 100 * time.Millisecond
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return fmt.Errorf("wait for retry: %w", ctx.Err())
	case <-timer.C:
		return nil
	}
}

func formatGraphQLErrors(graphqlErrors gqlerror.List) error {
	messages := make([]string, 0, len(graphqlErrors))
	for _, graphqlError := range graphqlErrors {
		code, ok := graphqlError.Extensions["code"]
		if ok {
			messages = append(messages, fmt.Sprintf("%s (%s)", graphqlError.Error(), code))
			continue
		}
		messages = append(messages, graphqlError.Error())
	}

	return fmt.Errorf("%w: %s", ErrGraphQL, strings.Join(messages, "; "))
}

func firstNonEmpty(primary string, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}

func defaultDuration(value time.Duration, fallback time.Duration) time.Duration {
	if value > 0 {
		return value
	}

	return fallback
}

func defaultRetries(value int) int {
	if value > 0 {
		return value
	}

	return 3
}
