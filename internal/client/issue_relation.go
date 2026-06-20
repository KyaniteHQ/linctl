package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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
	Relations   []IssueRelationSummary `json:"relations"`
	HasNextPage bool                   `json:"has_next_page"`
	EndCursor   *string                `json:"end_cursor,omitempty"`
}

// ListIssueRelations returns visible relations between issues.
func ListIssueRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (IssueRelationList, error) {
	result, err := issueRelations(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue relations: %w", err)
	}

	relations := make([]IssueRelationSummary, 0, len(result.IssueRelations.Nodes))
	for _, relation := range result.IssueRelations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.IssueRelations.PageInfo.HasNextPage,
		EndCursor:   result.IssueRelations.PageInfo.EndCursor,
	}, nil
}

// GetIssueRelationByID returns one issue relation by Linear id.
func GetIssueRelationByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueRelationSummary, error) {
	result, err := issueRelation(ctx, graphqlClient, id)
	if err != nil {
		return IssueRelationSummary{}, fmt.Errorf("get issue relation %s: %w", id, err)
	}

	return issueRelationSummary(result.IssueRelation.IssueRelationSummaryFields), nil
}

func issueRelationSummary(relation IssueRelationSummaryFields) IssueRelationSummary {
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
