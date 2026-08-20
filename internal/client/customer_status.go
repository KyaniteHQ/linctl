//nolint:dupl // Customer status and tier reads intentionally mirror Linear's parallel schema types.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// CustomerStatusSummary is the compact customer status model used by read-only commands.
type CustomerStatusSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Color       string  `json:"color"`
	Description string  `json:"description,omitempty"`
	Position    float64 `json:"position"`
	ArchivedAt  string  `json:"archived_at,omitempty"`
}

// CustomerStatusList is a page of Linear customer statuses.
type CustomerStatusList struct {
	Statuses []CustomerStatusSummary `json:"customer_statuses"`
	Page
}

//nolint:lll
type customerStatusesNode = gql.XCustomerStatusesCustomerStatusesCustomerStatusConnectionNodesCustomerStatus

type customerStatusesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListCustomerStatuses returns organization customer lifecycle statuses.
func ListCustomerStatuses(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerStatusList, error) {
	query := customerStatusesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list customer statuses", limit, defaultListPageSize,
		query.page,
		customerStatusesNodeSummary,
	)
	if err != nil {
		return CustomerStatusList{}, err
	}

	return CustomerStatusList{Statuses: page.Items, Page: page.Page}, nil
}

// GetCustomerStatusByID returns one customer lifecycle status by id.
func GetCustomerStatusByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomerStatusSummary, error) {
	result, err := gql.XCustomerStatus(ctx, graphqlClient, id)
	if err != nil {
		return CustomerStatusSummary{}, fmt.Errorf("get customer status %s: %w", id, err)
	}

	return customerStatusSummary(result.CustomerStatus.CustomerStatusSummaryFields), nil
}

func (query customerStatusesQuery) page(
	pageSize int,
	after *string,
) ([]customerStatusesNode, bool, *string, error) {
	result, err := gql.XCustomerStatuses(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.CustomerStatuses.Nodes,
		result.CustomerStatuses.PageInfo.HasNextPage,
		result.CustomerStatuses.PageInfo.EndCursor,
		nil
}

func customerStatusesNodeSummary(node customerStatusesNode) CustomerStatusSummary {
	return customerStatusSummary(node.CustomerStatusSummaryFields)
}

func customerStatusSummary(fields gql.CustomerStatusSummaryFields) CustomerStatusSummary {
	return CustomerStatusSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		DisplayName: fields.DisplayName,
		Color:       fields.Color,
		Description: stringValue(fields.Description),
		Position:    fields.Position,
		ArchivedAt:  stringValue(fields.ArchivedAt),
	}
}
