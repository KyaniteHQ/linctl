package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Relations []InitiativeRelationSummary `json:"relations"`
	Page
}

//nolint:lll
type initiativeRelationNode = gql.XInitiativeRelationsInitiativeRelationsInitiativeRelationConnectionNodesInitiativeRelation

// ListInitiativeRelations returns visible parent-child relations between initiatives.
func ListInitiativeRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (InitiativeRelationList, error) {
	page, err := listConnection(
		"list initiative relations", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]initiativeRelationNode, bool, *string, error) {
			result, err := gql.XInitiativeRelations(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.InitiativeRelations.Nodes,
				result.InitiativeRelations.PageInfo.HasNextPage,
				result.InitiativeRelations.PageInfo.EndCursor,
				nil
		},
		initiativeRelationNodeSummary,
	)
	if err != nil {
		return InitiativeRelationList{}, err
	}

	return InitiativeRelationList{Relations: page.Items, Page: page.Page}, nil
}

// GetInitiativeRelationByID returns one initiative relation by Linear id.
func GetInitiativeRelationByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeRelationSummary, error) {
	result, err := gql.XInitiativeRelation(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeRelationSummary{}, fmt.Errorf("get initiative relation %s: %w", id, err)
	}

	return initiativeRelationSummary(result.InitiativeRelation.InitiativeRelationSummaryFields), nil
}

func initiativeRelationNodeSummary(relation initiativeRelationNode) InitiativeRelationSummary {
	return initiativeRelationSummary(relation.InitiativeRelationSummaryFields)
}

func initiativeRelationSummary(relation gql.InitiativeRelationSummaryFields) InitiativeRelationSummary {
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
