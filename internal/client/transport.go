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
	Client     *http.Client
	Endpoint   string
	Token      AuthToken
	Timeout    time.Duration
	MaxRetries int
}

// Transport implements genqlient's GraphQL client interface.
type Transport struct {
	httpClient *http.Client
	endpoint   string
	token      AuthToken
	timeout    time.Duration
	maxRetries int
}

// PersonalAPIToken sends a raw Linear personal API key.
func PersonalAPIToken(value string) AuthToken {
	return AuthToken{
		authorization: value,
	}
}

// OAuthToken sends a Bearer OAuth token.
func OAuthToken(value string) AuthToken {
	return AuthToken{
		authorization: "Bearer " + value,
	}
}

// NewTransport creates a Linear GraphQL transport.
func NewTransport(config TransportConfig) *Transport {
	httpClient := config.Client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Transport{
		httpClient: httpClient,
		endpoint:   firstNonEmpty(config.Endpoint, "https://api.linear.app/graphql"),
		token:      config.Token,
		timeout:    defaultDuration(config.Timeout, 30*time.Second),
		maxRetries: defaultRetries(config.MaxRetries),
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
		return fmt.Errorf("encode graphql request: %w", err)
	}

	for attempt := range transport.maxRetries + 1 {
		body, statusCode, header, err := transport.send(ctx, payload)
		if err != nil {
			return err
		}
		if statusCode == http.StatusTooManyRequests && attempt < transport.maxRetries {
			if err := waitForRetry(ctx, retryDelay(header.Get("Retry-After"), attempt)); err != nil {
				return err
			}
			continue
		}
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

	return fmt.Errorf("graphql retry loop exhausted: %w", ErrGraphQL)
}

// AssertMutationSuccess checks a mutation payload for success and returned entity id.
func AssertMutationSuccess(reader io.Reader, mutationName string, entityPath string) error {
	var root map[string]json.RawMessage
	if err := json.NewDecoder(reader).Decode(&root); err != nil {
		return fmt.Errorf("decode mutation payload: %w", err)
	}

	mutationPayload, ok := root[mutationName]
	if !ok {
		return fmt.Errorf("%w: mutation %s missing", ErrMutationFailed, mutationName)
	}

	var mutation map[string]json.RawMessage
	if err := json.Unmarshal(mutationPayload, &mutation); err != nil {
		return fmt.Errorf("decode mutation %s: %w", mutationName, err)
	}
	if !jsonBool(mutation["success"]) {
		return fmt.Errorf("%w: mutation %s success false or missing", ErrMutationFailed, mutationName)
	}
	if entityPath != "" && !jsonPathStringExists(mutationPayload, entityPath) {
		return fmt.Errorf("%w: mutation %s entity %s missing", ErrMutationFailed, mutationName, entityPath)
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

func jsonBool(raw json.RawMessage) bool {
	var value bool
	if err := json.Unmarshal(raw, &value); err != nil {
		return false
	}

	return value
}

func jsonPathStringExists(raw json.RawMessage, path string) bool {
	current := raw
	for _, pathPart := range strings.Split(path, ".") {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(current, &object); err != nil {
			return false
		}
		next, ok := object[pathPart]
		if !ok {
			return false
		}
		current = next
	}

	var value string
	if err := json.Unmarshal(current, &value); err != nil {
		return false
	}

	return value != ""
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
