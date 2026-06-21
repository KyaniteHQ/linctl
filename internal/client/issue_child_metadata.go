package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// CustomerNeedMetadataSummary is body-free customer need metadata for issue child reads.
type CustomerNeedMetadataSummary struct {
	ID           string  `json:"id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	ArchivedAt   string  `json:"archived_at,omitempty"`
	Priority     float64 `json:"priority"`
	URL          string  `json:"url,omitempty"`
	CustomerID   string  `json:"customer_id,omitempty"`
	CustomerName string  `json:"customer_name,omitempty"`
	IssueID      string  `json:"issue_id,omitempty"`
	Issue        string  `json:"issue,omitempty"`
	IssueTitle   string  `json:"issue_title,omitempty"`
	ProjectID    string  `json:"project_id,omitempty"`
	ProjectName  string  `json:"project_name,omitempty"`
}

// IssueCommentMetadataList is a page of body-free comments for one issue-like root.
type IssueCommentMetadataList struct {
	IssueID     string                   `json:"issue_id"`
	Identifier  string                   `json:"identifier"`
	Comments    []CommentMetadataSummary `json:"comments"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// IssueCustomerNeedMetadataList is a page of body-free customer needs for one issue-like root.
type IssueCustomerNeedMetadataList struct {
	IssueID     string                        `json:"issue_id"`
	Identifier  string                        `json:"identifier"`
	Needs       []CustomerNeedMetadataSummary `json:"customer_needs"`
	HasNextPage bool                          `json:"has_next_page"`
	EndCursor   *string                       `json:"end_cursor,omitempty"`
}

// IssueSharedAccessSummary is compact shared-access metadata without shared user details.
type IssueSharedAccessSummary struct {
	IssueID                   string   `json:"issue_id"`
	Identifier                string   `json:"identifier"`
	IsShared                  bool     `json:"is_shared"`
	ViewerHasOnlySharedAccess bool     `json:"viewer_has_only_shared_access"`
	SharedWithCount           int      `json:"shared_with_count"`
	DisallowedIssueFields     []string `json:"disallowed_issue_fields,omitempty"`
}

// ListIssueNeeds returns body-free customer needs associated with one issue.
func ListIssueNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	result, err := issue_needs(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf("list issue customer needs %s: %w", id, err)
	}

	needs := make([]CustomerNeedMetadataSummary, 0, len(result.Issue.Needs.Nodes))
	for _, need := range result.Issue.Needs.Nodes {
		needs = append(needs, customerNeedMetadataSummary(need.CustomerNeedMetadataFields))
	}

	return IssueCustomerNeedMetadataList{
		IssueID:     result.Issue.Id,
		Identifier:  result.Issue.Identifier,
		Needs:       needs,
		HasNextPage: result.Issue.Needs.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Needs.PageInfo.EndCursor,
	}, nil
}

// ListIssueFormerNeeds returns body-free customer needs formerly associated with one issue.
func ListIssueFormerNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	result, err := issue_formerNeeds(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf("list issue former customer needs %s: %w", id, err)
	}

	needs := make([]CustomerNeedMetadataSummary, 0, len(result.Issue.FormerNeeds.Nodes))
	for _, need := range result.Issue.FormerNeeds.Nodes {
		needs = append(needs, customerNeedMetadataSummary(need.CustomerNeedMetadataFields))
	}

	return IssueCustomerNeedMetadataList{
		IssueID:     result.Issue.Id,
		Identifier:  result.Issue.Identifier,
		Needs:       needs,
		HasNextPage: result.Issue.FormerNeeds.PageInfo.HasNextPage,
		EndCursor:   result.Issue.FormerNeeds.PageInfo.EndCursor,
	}, nil
}

// GetIssueSharedAccess returns compact shared-access metadata for one issue.
func GetIssueSharedAccess(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueSharedAccessSummary, error) {
	result, err := issue_sharedAccess(ctx, graphqlClient, id)
	if err != nil {
		return IssueSharedAccessSummary{}, fmt.Errorf("get issue shared access %s: %w", id, err)
	}

	return issueSharedAccessSummary(
		result.Issue.Id,
		result.Issue.Identifier,
		result.Issue.SharedAccess.IssueSharedAccessFields,
	), nil
}

// ListIssueVCSBranchComments returns body-free comments for the issue matched by a VCS branch.
//
//nolint:dupl // VCS branch child-read methods intentionally mirror generated GraphQL connection shapes.
func ListIssueVCSBranchComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueCommentMetadataList, error) {
	result, err := issueVcsBranchSearch_comments(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCommentMetadataList{}, fmt.Errorf("list issue vcs branch comments %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueCommentMetadataList{}, fmt.Errorf("list issue vcs branch comments %s: not found", branchName)
	}

	comments := make([]CommentMetadataSummary, 0, len(result.IssueVcsBranchSearch.Comments.Nodes))
	for _, comment := range result.IssueVcsBranchSearch.Comments.Nodes {
		comments = append(comments, commentMetadataSummary(comment.CommentMetadataFields))
	}

	return IssueCommentMetadataList{
		IssueID:     result.IssueVcsBranchSearch.Id,
		Identifier:  result.IssueVcsBranchSearch.Identifier,
		Comments:    comments,
		HasNextPage: result.IssueVcsBranchSearch.Comments.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Comments.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchNeeds returns body-free customer needs for the issue matched by a VCS branch.
//
//nolint:dupl // VCS branch child-read methods intentionally mirror generated GraphQL connection shapes.
func ListIssueVCSBranchNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	result, err := issueVcsBranchSearch_needs(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf(
			"list issue vcs branch customer needs %s: %w",
			branchName,
			err,
		)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf(
			"list issue vcs branch customer needs %s: not found",
			branchName,
		)
	}

	needs := make([]CustomerNeedMetadataSummary, 0, len(result.IssueVcsBranchSearch.Needs.Nodes))
	for _, need := range result.IssueVcsBranchSearch.Needs.Nodes {
		needs = append(needs, customerNeedMetadataSummary(need.CustomerNeedMetadataFields))
	}

	return IssueCustomerNeedMetadataList{
		IssueID:     result.IssueVcsBranchSearch.Id,
		Identifier:  result.IssueVcsBranchSearch.Identifier,
		Needs:       needs,
		HasNextPage: result.IssueVcsBranchSearch.Needs.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Needs.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchFormerNeeds returns body-free former customer needs for the issue matched by a VCS branch.
//
//nolint:dupl // VCS branch child-read methods intentionally mirror generated GraphQL connection shapes.
func ListIssueVCSBranchFormerNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	result, err := issueVcsBranchSearch_formerNeeds(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf(
			"list issue vcs branch former customer needs %s: %w",
			branchName,
			err,
		)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf(
			"list issue vcs branch former customer needs %s: not found",
			branchName,
		)
	}

	needs := make([]CustomerNeedMetadataSummary, 0, len(result.IssueVcsBranchSearch.FormerNeeds.Nodes))
	for _, need := range result.IssueVcsBranchSearch.FormerNeeds.Nodes {
		needs = append(needs, customerNeedMetadataSummary(need.CustomerNeedMetadataFields))
	}

	return IssueCustomerNeedMetadataList{
		IssueID:     result.IssueVcsBranchSearch.Id,
		Identifier:  result.IssueVcsBranchSearch.Identifier,
		Needs:       needs,
		HasNextPage: result.IssueVcsBranchSearch.FormerNeeds.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.FormerNeeds.PageInfo.EndCursor,
	}, nil
}

// GetIssueVCSBranchSharedAccess returns compact shared-access metadata for an issue matched by a VCS branch.
func GetIssueVCSBranchSharedAccess(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
) (IssueSharedAccessSummary, error) {
	result, err := issueVcsBranchSearch_sharedAccess(ctx, graphqlClient, branchName)
	if err != nil {
		return IssueSharedAccessSummary{}, fmt.Errorf("get issue vcs branch shared access %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueSharedAccessSummary{}, fmt.Errorf("get issue vcs branch shared access %s: not found", branchName)
	}

	return issueSharedAccessSummary(
		result.IssueVcsBranchSearch.Id,
		result.IssueVcsBranchSearch.Identifier,
		result.IssueVcsBranchSearch.SharedAccess.IssueSharedAccessFields,
	), nil
}

func customerNeedMetadataSummary(fields CustomerNeedMetadataFields) CustomerNeedMetadataSummary {
	summary := CustomerNeedMetadataSummary{
		ID:         fields.Id,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
		Priority:   fields.Priority,
		URL:        stringValue(fields.Url),
	}
	if fields.Customer != nil {
		summary.CustomerID = fields.Customer.Id
		summary.CustomerName = fields.Customer.Name
	}
	if fields.Issue != nil {
		summary.IssueID = fields.Issue.Id
		summary.Issue = fields.Issue.Identifier
		summary.IssueTitle = fields.Issue.Title
	}
	if fields.Project != nil {
		summary.ProjectID = fields.Project.Id
		summary.ProjectName = fields.Project.Name
	}

	return summary
}

func issueSharedAccessSummary(
	issueID string,
	identifier string,
	fields IssueSharedAccessFields,
) IssueSharedAccessSummary {
	return IssueSharedAccessSummary{
		IssueID:                   issueID,
		Identifier:                identifier,
		IsShared:                  fields.IsShared,
		ViewerHasOnlySharedAccess: fields.ViewerHasOnlySharedAccess,
		SharedWithCount:           fields.SharedWithCount,
		DisallowedIssueFields:     issueSharedAccessDisallowedFields(fields.DisallowedIssueFields),
	}
}

func issueSharedAccessDisallowedFields(fields []IssueSharedAccessDisallowedField) []string {
	values := make([]string, 0, len(fields))
	for _, field := range fields {
		values = append(values, string(field))
	}

	return values
}
