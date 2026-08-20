package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// CustomerSummary is the compact customer model used by read-only commands.
type CustomerSummary struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Domains              []string `json:"domains"`
	ExternalIDs          []string `json:"external_ids"`
	SlackChannelID       string   `json:"slack_channel_id,omitempty"`
	StatusID             string   `json:"status_id"`
	StatusName           string   `json:"status_name"`
	TierID               string   `json:"tier_id,omitempty"`
	TierName             string   `json:"tier_name,omitempty"`
	OwnerID              string   `json:"owner_id,omitempty"`
	OwnerDisplayName     string   `json:"owner_display_name,omitempty"`
	Revenue              *int     `json:"revenue,omitempty"`
	Size                 *float64 `json:"size,omitempty"`
	ApproximateNeedCount float64  `json:"approximate_need_count"`
	SlugID               string   `json:"slug_id"`
	URL                  string   `json:"url"`
}

// CustomerList is a page of Linear customers.
type CustomerList struct {
	Customers   []CustomerSummary `json:"customers"`
	HasNextPage bool              `json:"has_next_page"`
	EndCursor   *string           `json:"end_cursor,omitempty"`
}

//nolint:lll
type customersNode = gql.XCustomersCustomersCustomerConnectionNodesCustomer

type customersQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListCustomers returns visible Linear customers.
func ListCustomers(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerList, error) {
	query := customersQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list customers", limit, defaultListPageSize,
		query.page,
		customersNodeSummary,
	)
	if err != nil {
		return CustomerList{}, err
	}

	return CustomerList{Customers: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetCustomerByID returns one Linear customer by id or slug.
func GetCustomerByID(ctx context.Context, graphqlClient graphql.Client, id string) (CustomerSummary, error) {
	result, err := gql.XCustomer(ctx, graphqlClient, id)
	if err != nil {
		return CustomerSummary{}, fmt.Errorf("get customer %s: %w", id, err)
	}

	return customerSummary(result.Customer.CustomerSummaryFields), nil
}

func (query customersQuery) page(pageSize int, after *string) ([]customersNode, bool, *string, error) {
	result, err := gql.XCustomers(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Customers.Nodes,
		result.Customers.PageInfo.HasNextPage,
		result.Customers.PageInfo.EndCursor,
		nil
}

func customersNodeSummary(node customersNode) CustomerSummary {
	return customerSummary(node.CustomerSummaryFields)
}

func customerSummary(fields gql.CustomerSummaryFields) CustomerSummary {
	summary := CustomerSummary{
		ID:                   fields.Id,
		Name:                 fields.Name,
		Domains:              fields.Domains,
		ExternalIDs:          fields.ExternalIds,
		SlackChannelID:       stringValue(fields.SlackChannelId),
		StatusID:             fields.Status.Id,
		StatusName:           fields.Status.Name,
		Revenue:              fields.Revenue,
		Size:                 fields.Size,
		ApproximateNeedCount: fields.ApproximateNeedCount,
		SlugID:               fields.SlugId,
		URL:                  fields.Url,
	}
	if fields.Tier != nil {
		summary.TierID = fields.Tier.Id
		summary.TierName = fields.Tier.Name
	}
	if fields.Owner != nil {
		summary.OwnerID = fields.Owner.Id
		summary.OwnerDisplayName = fields.Owner.DisplayName
	}

	return summary
}
