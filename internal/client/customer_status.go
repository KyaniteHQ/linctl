//nolint:dupl // Customer status and tier reads intentionally mirror Linear's parallel schema types.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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
	Statuses    []CustomerStatusSummary `json:"customer_statuses"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// ListCustomerStatuses returns workspace customer lifecycle statuses.
func ListCustomerStatuses(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerStatusList, error) {
	result, err := customerStatuses(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CustomerStatusList{}, fmt.Errorf("list customer statuses: %w", err)
	}

	summaries := make([]CustomerStatusSummary, 0, len(result.CustomerStatuses.Nodes))
	for _, node := range result.CustomerStatuses.Nodes {
		summaries = append(summaries, customerStatusSummary(node.CustomerStatusSummaryFields))
	}

	return CustomerStatusList{
		Statuses:    summaries,
		HasNextPage: result.CustomerStatuses.PageInfo.HasNextPage,
		EndCursor:   result.CustomerStatuses.PageInfo.EndCursor,
	}, nil
}

// GetCustomerStatusByID returns one customer lifecycle status by id.
func GetCustomerStatusByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomerStatusSummary, error) {
	result, err := customerStatus(ctx, graphqlClient, id)
	if err != nil {
		return CustomerStatusSummary{}, fmt.Errorf("get customer status %s: %w", id, err)
	}

	return customerStatusSummary(result.CustomerStatus.CustomerStatusSummaryFields), nil
}

func customerStatusSummary(fields CustomerStatusSummaryFields) CustomerStatusSummary {
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
