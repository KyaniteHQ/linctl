package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// IssueRelationSummary is one directional relation between Linear issues.
type IssueRelationSummary struct {
	ID                     string `json:"id"`
	Type                   string `json:"type"`
	IssueID                string `json:"issue_id"`
	IssueIdentifier        string `json:"issue_identifier"`
	IssueTitle             string `json:"issue_title"`
	RelatedIssueID         string `json:"related_issue_id"`
	RelatedIssueIdentifier string `json:"related_issue_identifier"`
	RelatedIssueTitle      string `json:"related_issue_title"`
	CreatedAt              string `json:"created_at"`
	UpdatedAt              string `json:"updated_at"`
	ArchivedAt             string `json:"archived_at,omitempty"`
}

// IssueRelationList is a page of issue relations.
type IssueRelationList struct {
	Relations []IssueRelationSummary `json:"relations"`
	Page
}

//nolint:lll
type issueRelationsNode = gql.XIssueRelationsIssueRelationsIssueRelationConnectionNodesIssueRelation

// ListIssueRelations returns visible relations between issues.
func ListIssueRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (IssueRelationList, error) {
	page, err := listConnection(
		"list issue relations", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]issueRelationsNode, bool, *string, error) {
			result, err := gql.XIssueRelations(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.IssueRelations.Nodes,
				result.IssueRelations.PageInfo.HasNextPage,
				result.IssueRelations.PageInfo.EndCursor,
				nil
		},
		issueRelationNodeSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, Page: page.Page}, nil
}

// GetIssueRelationByID returns one issue relation by Linear id.
func GetIssueRelationByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueRelationSummary, error) {
	result, err := gql.XIssueRelation(ctx, graphqlClient, id)
	if err != nil {
		return IssueRelationSummary{}, fmt.Errorf("get issue relation %s: %w", id, err)
	}

	return issueRelationSummary(result.IssueRelation.IssueRelationSummaryFields), nil
}

func issueRelationNodeSummary(relation issueRelationsNode) IssueRelationSummary {
	return issueRelationSummary(relation.IssueRelationSummaryFields)
}

func issueRelationSummary(relation gql.IssueRelationSummaryFields) IssueRelationSummary {
	return IssueRelationSummary{
		ID:                     relation.Id,
		Type:                   relation.Type,
		IssueID:                relation.Issue.Id,
		IssueIdentifier:        relation.Issue.Identifier,
		IssueTitle:             relation.Issue.Title,
		RelatedIssueID:         relation.RelatedIssue.Id,
		RelatedIssueIdentifier: relation.RelatedIssue.Identifier,
		RelatedIssueTitle:      relation.RelatedIssue.Title,
		CreatedAt:              relation.CreatedAt,
		UpdatedAt:              relation.UpdatedAt,
		ArchivedAt:             stringValue(relation.ArchivedAt),
	}
}
