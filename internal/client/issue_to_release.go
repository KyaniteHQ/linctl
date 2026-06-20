//nolint:dupl // Minimal association read glue is intentionally uniform across release-association domains.
package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListIssueToReleases returns visible Issue-to-Release associations.
func ListIssueToReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (IssueToReleaseList, error) {
	result, err := issueToReleases(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueToReleaseList{}, fmt.Errorf("list issue to releases: %w", err)
	}

	associations := make([]IssueToReleaseSummary, 0, len(result.IssueToReleases.Nodes))
	for _, association := range result.IssueToReleases.Nodes {
		associations = append(
			associations,
			issueToReleaseSummary(association.IssueToReleaseSummaryFields),
		)
	}

	return IssueToReleaseList{
		Associations: associations,
		HasNextPage:  result.IssueToReleases.PageInfo.HasNextPage,
		EndCursor:    result.IssueToReleases.PageInfo.EndCursor,
	}, nil
}

// GetIssueToReleaseByID returns one Issue-to-Release association by Linear id.
func GetIssueToReleaseByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueToReleaseSummary, error) {
	result, err := issueToRelease(ctx, graphqlClient, id)
	if err != nil {
		return IssueToReleaseSummary{}, fmt.Errorf("get issue to release %s: %w", id, err)
	}

	return issueToReleaseSummary(result.IssueToRelease.IssueToReleaseSummaryFields), nil
}

func issueToReleaseSummary(association IssueToReleaseSummaryFields) IssueToReleaseSummary {
	return IssueToReleaseSummary{
		ID:         association.Id,
		IssueID:    association.Issue.Id,
		ReleaseID:  association.Release.Id,
		CreatedAt:  association.CreatedAt,
		UpdatedAt:  association.UpdatedAt,
		ArchivedAt: stringValue(association.ArchivedAt),
	}
}
