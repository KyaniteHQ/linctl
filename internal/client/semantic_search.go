package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// SemanticSearchResultSummary is a compact reference returned by semantic search.
type SemanticSearchResultSummary struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Key   string `json:"key,omitempty"`
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
}

// SemanticSearchList is the semantic search result set.
type SemanticSearchList struct {
	Results []SemanticSearchResultSummary `json:"results"`
}

// SearchSemantic returns compact references from Linear semantic search.
func SearchSemantic(
	ctx context.Context,
	graphqlClient graphql.Client,
	query string,
	limit int,
) (SemanticSearchList, error) {
	result, err := semanticSearch(ctx, graphqlClient, query, intPtr(limit), boolPtr(false))
	if err != nil {
		return SemanticSearchList{}, fmt.Errorf("semantic search: %w", err)
	}

	results := make([]SemanticSearchResultSummary, 0, len(result.SemanticSearch.Results))
	for _, searchResult := range result.SemanticSearch.Results {
		results = append(results, semanticSearchResultSummary(searchResult.SemanticSearchResultSummaryFields))
	}

	return SemanticSearchList{Results: results}, nil
}

func semanticSearchResultSummary(fields SemanticSearchResultSummaryFields) SemanticSearchResultSummary {
	summary := SemanticSearchResultSummary{
		Type: string(fields.Type),
		ID:   fields.Id,
	}
	if fields.Issue != nil {
		summary.Key = fields.Issue.Identifier
		summary.Title = fields.Issue.Title
		summary.URL = fields.Issue.Url
		return summary
	}
	if fields.Project != nil {
		summary.Title = fields.Project.Name
		summary.URL = fields.Project.Url
		return summary
	}
	if fields.Initiative != nil {
		summary.Title = fields.Initiative.Name
		summary.URL = fields.Initiative.Url
		return summary
	}
	if fields.Document != nil {
		summary.Title = fields.Document.Title
		summary.URL = fields.Document.Url
		return summary
	}

	return summary
}
