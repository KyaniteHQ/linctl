package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// RateLimitStatus is the authenticated client's current Linear rate-limit state.
type RateLimitStatus struct {
	Identifier string      `json:"identifier,omitempty"`
	Kind       string      `json:"kind"`
	Limits     []RateLimit `json:"limits"`
}

// RateLimit is one quota bucket inside Linear's current rate-limit state.
type RateLimit struct {
	Type            string  `json:"type"`
	RequestedAmount float64 `json:"requested_amount"`
	AllowedAmount   float64 `json:"allowed_amount"`
	Period          float64 `json:"period"`
	RemainingAmount float64 `json:"remaining_amount"`
	Reset           float64 `json:"reset"`
}

// GetRateLimitStatus returns the authenticated client's current Linear quota state.
func GetRateLimitStatus(ctx context.Context, graphqlClient graphql.Client) (RateLimitStatus, error) {
	result, err := gql.XRateLimitStatus(ctx, graphqlClient)
	if err != nil {
		return RateLimitStatus{}, err
	}

	limits := mapNodes(result.RateLimitStatus.Limits, func(
		limit gql.XRateLimitStatusRateLimitStatusRateLimitPayloadLimitsRateLimitResultPayload,
	) RateLimit {
		return RateLimit(limit)
	})

	return RateLimitStatus{
		Identifier: stringValue(result.RateLimitStatus.Identifier),
		Kind:       result.RateLimitStatus.Kind,
		Limits:     limits,
	}, nil
}
