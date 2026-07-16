package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// GetIssueDependencies returns parent, child, and blocking relationships for one issue.
func GetIssueDependencies(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueDependencyGraph, error) {
	dependencies, err := gql.IssueDependencies(ctx, graphqlClient, id, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueDependencyGraph{}, fmt.Errorf("get issue dependencies %s: %w", id, err)
	}

	issue := dependencies.Issue
	return IssueDependencyGraph{
		ID:          issue.Id,
		Identifier:  issue.Identifier,
		Parent:      issueDependencyParent(issue.Parent),
		Children:    issueDependencyChildren(issue.Children.Nodes),
		Blocks:      issueDependencyBlocks(issue.Relations.Nodes),
		BlockedBy:   issueDependencyBlockedBy(issue.InverseRelations.Nodes),
		HasNextPage: issueDependencyHasNextPage(issue),
	}, nil
}

func issueDependencyParent(issue *gql.IssueDependenciesIssueParentIssue) *IssueSummary {
	if issue == nil {
		return nil
	}

	summary := issueSummaryFromFields(issue.IssueSummaryFields)
	return &summary
}

func issueDependencyChildren(issues []gql.IssueDependenciesIssueChildrenIssueConnectionNodesIssue) []IssueSummary {
	summaries := mapNodes(issues, func(issue gql.IssueDependenciesIssueChildrenIssueConnectionNodesIssue) IssueSummary {
		return issueSummaryFromFields(issue.IssueSummaryFields)
	})

	return summaries
}

func issueDependencyBlocks(
	relations []gql.IssueDependenciesIssueRelationsIssueRelationConnectionNodesIssueRelation,
) []IssueSummary {
	summaries := make([]IssueSummary, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == "blocks" {
			summaries = append(summaries, issueSummaryFromFields(relation.RelatedIssue.IssueSummaryFields))
		}
	}

	return summaries
}

func issueDependencyBlockedBy(
	relations []gql.IssueDependenciesIssueInverseRelationsIssueRelationConnectionNodesIssueRelation,
) []IssueSummary {
	summaries := make([]IssueSummary, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == "blocks" {
			summaries = append(summaries, issueSummaryFromFields(relation.Issue.IssueSummaryFields))
		}
	}

	return summaries
}

func issueDependencyHasNextPage(issue gql.IssueDependenciesIssue) bool {
	return issue.Children.PageInfo.HasNextPage ||
		issue.Relations.PageInfo.HasNextPage ||
		issue.InverseRelations.PageInfo.HasNextPage
}
