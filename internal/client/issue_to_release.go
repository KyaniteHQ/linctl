package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// IssueToReleaseSummary is one issue association under a release.
type IssueToReleaseSummary struct {
	ID         string `json:"id"`
	IssueID    string `json:"issue_id"`
	ReleaseID  string `json:"release_id"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	ArchivedAt string `json:"archived_at,omitempty"`
}

// IssueToReleaseList is a page of issue-to-release associations.
type IssueToReleaseList struct {
	Associations []IssueToReleaseSummary `json:"associations"`
	HasNextPage  bool                    `json:"has_next_page"`
	EndCursor    *string                 `json:"end_cursor,omitempty"`
}

//nolint:lll
type issueToReleaseNode = gql.XIssueToReleasesIssueToReleasesIssueToReleaseConnectionNodesIssueToRelease

// ListIssueToReleases returns visible Issue-to-Release associations.
func ListIssueToReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (IssueToReleaseList, error) {
	page, err := listConnection(
		"list issue to releases", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]issueToReleaseNode, bool, *string, error) {
			result, err := gql.XIssueToReleases(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.IssueToReleases.Nodes,
				result.IssueToReleases.PageInfo.HasNextPage,
				result.IssueToReleases.PageInfo.EndCursor,
				nil
		},
		issueToReleaseNodeSummary,
	)
	if err != nil {
		return IssueToReleaseList{}, err
	}

	return IssueToReleaseList{
		Associations: page.Items,
		HasNextPage:  page.HasNextPage,
		EndCursor:    page.EndCursor,
	}, nil
}

// GetIssueToReleaseByID returns one Issue-to-Release association by Linear id.
func GetIssueToReleaseByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueToReleaseSummary, error) {
	result, err := gql.XIssueToRelease(ctx, graphqlClient, id)
	if err != nil {
		return IssueToReleaseSummary{}, fmt.Errorf("get issue to release %s: %w", id, err)
	}

	return issueToReleaseSummary(result.IssueToRelease.IssueToReleaseSummaryFields), nil
}

func issueToReleaseNodeSummary(association issueToReleaseNode) IssueToReleaseSummary {
	return issueToReleaseSummary(association.IssueToReleaseSummaryFields)
}

func issueToReleaseSummary(association gql.IssueToReleaseSummaryFields) IssueToReleaseSummary {
	return IssueToReleaseSummary{
		ID:         association.Id,
		IssueID:    association.Issue.Id,
		ReleaseID:  association.Release.Id,
		CreatedAt:  association.CreatedAt,
		UpdatedAt:  association.UpdatedAt,
		ArchivedAt: stringValue(association.ArchivedAt),
	}
}
