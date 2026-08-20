//nolint:dupl // Customer status and tier reads intentionally mirror Linear's parallel schema types.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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

//nolint:lll
type customerTiersNode = gql.XCustomerTiersCustomerTiersCustomerTierConnectionNodesCustomerTier

type customerTiersQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListCustomerTiers returns organization customer tiers.
func ListCustomerTiers(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerTierList, error) {
	query := customerTiersQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list customer tiers", limit, defaultListPageSize,
		query.page,
		customerTiersNodeSummary,
	)
	if err != nil {
		return CustomerTierList{}, err
	}

	return CustomerTierList{Tiers: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetCustomerTierByID returns one customer tier by id.
func GetCustomerTierByID(ctx context.Context, graphqlClient graphql.Client, id string) (CustomerTierSummary, error) {
	result, err := gql.XCustomerTier(ctx, graphqlClient, id)
	if err != nil {
		return CustomerTierSummary{}, fmt.Errorf("get customer tier %s: %w", id, err)
	}

	return customerTierSummary(result.CustomerTier.CustomerTierSummaryFields), nil
}

func (query customerTiersQuery) page(pageSize int, after *string) ([]customerTiersNode, bool, *string, error) {
	result, err := gql.XCustomerTiers(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.CustomerTiers.Nodes,
		result.CustomerTiers.PageInfo.HasNextPage,
		result.CustomerTiers.PageInfo.EndCursor,
		nil
}

func customerTiersNodeSummary(node customerTiersNode) CustomerTierSummary {
	return customerTierSummary(node.CustomerTierSummaryFields)
}

func customerTierSummary(fields gql.CustomerTierSummaryFields) CustomerTierSummary {
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
