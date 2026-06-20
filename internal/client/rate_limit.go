package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"
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
	result, err := rateLimitStatus(ctx, graphqlClient)
	if err != nil {
		return RateLimitStatus{}, err
	}

	limits := make([]RateLimit, 0, len(result.RateLimitStatus.Limits))
	for _, limit := range result.RateLimitStatus.Limits {
		limits = append(limits, RateLimit(limit))
	}

	return RateLimitStatus{
		Identifier: stringValue(result.RateLimitStatus.Identifier),
		Kind:       result.RateLimitStatus.Kind,
		Limits:     limits,
	}, nil
}
