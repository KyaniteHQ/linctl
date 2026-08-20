package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// SearchIssuesByTeam searches issue content scoped to a resolved team.
func SearchIssuesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	query string,
	limit int,
) (IssueList, error) {
	issues, err := gql.XIssueSearch(ctx, graphqlClient, teamID, query, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("search issues: %w", err)
	}

	summaries := mapNodes(issues.IssueSearch.Nodes, searchIssueSummary)

	return IssueList{
		Issues: summaries,
		Page: Page{
			HasNextPage: issues.IssueSearch.PageInfo.HasNextPage,
			EndCursor:   issues.IssueSearch.PageInfo.EndCursor,
		},
	}, nil
}

// SearchIssuesByFigmaFileKey searches issues associated with a Figma file key.
func SearchIssuesByFigmaFileKey(
	ctx context.Context,
	graphqlClient graphql.Client,
	fileKey string,
	limit int,
) (IssueList, error) {
	issues, err := gql.XIssueFigmaFileKeySearch(ctx, graphqlClient, fileKey, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("search issues by Figma file key: %w", err)
	}

	summaries := mapNodes(issues.IssueFigmaFileKeySearch.Nodes, figmaFileKeyIssueSummary)

	return IssueList{
		Issues: summaries,
		Page: Page{
			HasNextPage: issues.IssueFigmaFileKeySearch.PageInfo.HasNextPage,
			EndCursor:   issues.IssueFigmaFileKeySearch.PageInfo.EndCursor,
		},
	}, nil
}

// GetIssueFilterSuggestion returns a JSON issue filter suggestion for a prompt.
func GetIssueFilterSuggestion(
	ctx context.Context,
	graphqlClient graphql.Client,
	prompt string,
	teamID string,
	projectID string,
) (IssueFilterSuggestion, error) {
	suggestion, err := gql.XIssueFilterSuggestion(
		ctx,
		graphqlClient,
		prompt,
		optionalString(teamID),
		optionalString(projectID),
	)
	if err != nil {
		return IssueFilterSuggestion{}, fmt.Errorf("get issue filter suggestion: %w", err)
	}

	filter := json.RawMessage(nil)
	if suggestion.IssueFilterSuggestion.Filter != nil {
		filter = *suggestion.IssueFilterSuggestion.Filter
	}

	return IssueFilterSuggestion{
		Filter: filter,
		LogID:  stringValue(suggestion.IssueFilterSuggestion.LogId),
	}, nil
}

// GetIssueTitleSuggestionFromCustomerRequest returns a title suggestion for customer request text.
func GetIssueTitleSuggestionFromCustomerRequest(
	ctx context.Context,
	graphqlClient graphql.Client,
	request string,
) (IssueTitleSuggestion, error) {
	suggestion, err := gql.XIssueTitleSuggestionFromCustomerRequest(ctx, graphqlClient, request)
	if err != nil {
		return IssueTitleSuggestion{}, fmt.Errorf("get issue title suggestion: %w", err)
	}

	return IssueTitleSuggestion{
		Title: suggestion.IssueTitleSuggestionFromCustomerRequest.Title,
		LogID: stringValue(suggestion.IssueTitleSuggestionFromCustomerRequest.LogId),
	}, nil
}

func searchIssueSummary(issue gql.XIssueSearchIssueSearchIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func figmaFileKeyIssueSummary(
	issue gql.XIssueFigmaFileKeySearchIssueFigmaFileKeySearchIssueConnectionNodesIssue,
) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}
