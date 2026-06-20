package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// InitiativeSummary is the compact initiative model used by read-only commands.
type InitiativeSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	TargetDate  string `json:"target_date,omitempty"`
	SlugID      string `json:"slug_id"`
	URL         string `json:"url"`
}

// InitiativeList is a page of initiatives.
type InitiativeList struct {
	Initiatives []InitiativeSummary `json:"initiatives"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ListInitiatives returns visible initiatives.
func ListInitiatives(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeList, error) {
	result, err := initiatives(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeList{}, fmt.Errorf("list initiatives: %w", err)
	}

	summaries := make([]InitiativeSummary, 0, len(result.Initiatives.Nodes))
	for _, node := range result.Initiatives.Nodes {
		summaries = append(summaries, initiativeSummary(node.InitiativeSummaryFields))
	}

	return InitiativeList{
		Initiatives: summaries,
		HasNextPage: result.Initiatives.PageInfo.HasNextPage,
		EndCursor:   result.Initiatives.PageInfo.EndCursor,
	}, nil
}

// GetInitiativeByID returns one initiative by Linear id or slug.
func GetInitiativeByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeSummary, error) {
	result, err := initiative(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeSummary{}, fmt.Errorf("get initiative %s: %w", id, err)
	}

	return initiativeSummary(result.Initiative.InitiativeSummaryFields), nil
}

func initiativeSummary(fields InitiativeSummaryFields) InitiativeSummary {
	return InitiativeSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		Description: stringValue(fields.Description),
		Status:      string(fields.Status),
		Priority:    fields.Priority,
		TargetDate:  stringValue(fields.TargetDate),
		SlugID:      fields.SlugId,
		URL:         fields.Url,
	}
}
