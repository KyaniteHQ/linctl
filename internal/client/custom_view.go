package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// CustomViewSummary is the compact custom view model used by read-only commands.
type CustomViewSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ModelName   string `json:"model_name"`
	Shared      bool   `json:"shared"`
	Color       string `json:"color,omitempty"`
	SlugID      string `json:"slug_id"`
}

// CustomViewList is a page of custom views.
type CustomViewList struct {
	CustomViews []CustomViewSummary `json:"custom_views"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// CustomViewSubscriberStatus reports whether a custom view has active subscribers.
type CustomViewSubscriberStatus struct {
	ID             string `json:"id"`
	HasSubscribers bool   `json:"has_subscribers"`
}

// ListCustomViews returns visible custom views.
func ListCustomViews(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomViewList, error) {
	result, err := customViews(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CustomViewList{}, fmt.Errorf("list custom views: %w", err)
	}

	summaries := make([]CustomViewSummary, 0, len(result.CustomViews.Nodes))
	for _, node := range result.CustomViews.Nodes {
		summaries = append(summaries, customViewSummary(node.CustomViewSummaryFields))
	}

	return CustomViewList{
		CustomViews: summaries,
		HasNextPage: result.CustomViews.PageInfo.HasNextPage,
		EndCursor:   result.CustomViews.PageInfo.EndCursor,
	}, nil
}

// GetCustomViewByID returns one custom view by Linear id or slug.
func GetCustomViewByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewSummary, error) {
	result, err := customView(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewSummary{}, fmt.Errorf("get custom view %s: %w", id, err)
	}

	return customViewSummary(result.CustomView.CustomViewSummaryFields), nil
}

// GetCustomViewSubscriberStatus returns whether a custom view has active subscribers.
func GetCustomViewSubscriberStatus(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewSubscriberStatus, error) {
	result, err := customViewHasSubscribers(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewSubscriberStatus{}, fmt.Errorf("get custom view subscribers %s: %w", id, err)
	}

	return CustomViewSubscriberStatus{
		ID:             id,
		HasSubscribers: result.CustomViewHasSubscribers.HasSubscribers,
	}, nil
}

func customViewSummary(fields CustomViewSummaryFields) CustomViewSummary {
	return CustomViewSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		Description: stringValue(fields.Description),
		ModelName:   fields.ModelName,
		Shared:      fields.Shared,
		Color:       stringValue(fields.Color),
		SlugID:      fields.SlugId,
	}
}
