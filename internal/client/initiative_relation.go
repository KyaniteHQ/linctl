package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// InitiativeRelationSummary is one parent-child relation between initiatives.
type InitiativeRelationSummary struct {
	ID                    string  `json:"id"`
	ParentInitiativeID    string  `json:"parent_initiative_id"`
	ParentInitiativeName  string  `json:"parent_initiative_name"`
	RelatedInitiativeID   string  `json:"related_initiative_id"`
	RelatedInitiativeName string  `json:"related_initiative_name"`
	SortOrder             float64 `json:"sort_order"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
	ArchivedAt            string  `json:"archived_at,omitempty"`
	UserID                string  `json:"user_id,omitempty"`
	Name                  string  `json:"name,omitempty"`
	DisplayName           string  `json:"display_name,omitempty"`
}

// InitiativeRelationList is a page of initiative relations.
type InitiativeRelationList struct {
	Relations   []InitiativeRelationSummary `json:"relations"`
	HasNextPage bool                        `json:"has_next_page"`
	EndCursor   *string                     `json:"end_cursor,omitempty"`
}

// ListInitiativeRelations returns visible parent-child relations between initiatives.
func ListInitiativeRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (InitiativeRelationList, error) {
	result, err := initiativeRelations(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeRelationList{}, fmt.Errorf("list initiative relations: %w", err)
	}

	relations := make([]InitiativeRelationSummary, 0, len(result.InitiativeRelations.Nodes))
	for _, relation := range result.InitiativeRelations.Nodes {
		relations = append(relations, initiativeRelationSummary(relation.InitiativeRelationSummaryFields))
	}

	return InitiativeRelationList{
		Relations:   relations,
		HasNextPage: result.InitiativeRelations.PageInfo.HasNextPage,
		EndCursor:   result.InitiativeRelations.PageInfo.EndCursor,
	}, nil
}

// GetInitiativeRelationByID returns one initiative relation by Linear id.
func GetInitiativeRelationByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeRelationSummary, error) {
	result, err := initiativeRelation(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeRelationSummary{}, fmt.Errorf("get initiative relation %s: %w", id, err)
	}

	return initiativeRelationSummary(result.InitiativeRelation.InitiativeRelationSummaryFields), nil
}

func initiativeRelationSummary(relation InitiativeRelationSummaryFields) InitiativeRelationSummary {
	summary := InitiativeRelationSummary{
		ID:                    relation.Id,
		ParentInitiativeID:    relation.Initiative.Id,
		ParentInitiativeName:  relation.Initiative.Name,
		RelatedInitiativeID:   relation.RelatedInitiative.Id,
		RelatedInitiativeName: relation.RelatedInitiative.Name,
		SortOrder:             relation.SortOrder,
		CreatedAt:             relation.CreatedAt,
		UpdatedAt:             relation.UpdatedAt,
		ArchivedAt:            stringValue(relation.ArchivedAt),
	}
	if relation.User != nil {
		summary.UserID = relation.User.Id
		summary.Name = relation.User.Name
		summary.DisplayName = relation.User.DisplayName
	}

	return summary
}
