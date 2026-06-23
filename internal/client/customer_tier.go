//nolint:dupl // Customer status and tier reads intentionally mirror Linear's parallel schema types.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// CustomerTierSummary is the compact customer tier model used by read-only commands.
type CustomerTierSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Color       string  `json:"color"`
	Description string  `json:"description,omitempty"`
	Position    float64 `json:"position"`
	ArchivedAt  string  `json:"archived_at,omitempty"`
}

// CustomerTierList is a page of Linear customer tiers.
type CustomerTierList struct {
	Tiers       []CustomerTierSummary `json:"customer_tiers"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// ListCustomerTiers returns organization customer tiers.
func ListCustomerTiers(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerTierList, error) {
	result, err := customerTiers(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CustomerTierList{}, fmt.Errorf("list customer tiers: %w", err)
	}

	summaries := make([]CustomerTierSummary, 0, len(result.CustomerTiers.Nodes))
	for _, node := range result.CustomerTiers.Nodes {
		summaries = append(summaries, customerTierSummary(node.CustomerTierSummaryFields))
	}

	return CustomerTierList{
		Tiers:       summaries,
		HasNextPage: result.CustomerTiers.PageInfo.HasNextPage,
		EndCursor:   result.CustomerTiers.PageInfo.EndCursor,
	}, nil
}

// GetCustomerTierByID returns one customer tier by id.
func GetCustomerTierByID(ctx context.Context, graphqlClient graphql.Client, id string) (CustomerTierSummary, error) {
	result, err := customerTier(ctx, graphqlClient, id)
	if err != nil {
		return CustomerTierSummary{}, fmt.Errorf("get customer tier %s: %w", id, err)
	}

	return customerTierSummary(result.CustomerTier.CustomerTierSummaryFields), nil
}

func customerTierSummary(fields CustomerTierSummaryFields) CustomerTierSummary {
	return CustomerTierSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		DisplayName: fields.DisplayName,
		Color:       fields.Color,
		Description: stringValue(fields.Description),
		Position:    fields.Position,
		ArchivedAt:  stringValue(fields.ArchivedAt),
	}
}
