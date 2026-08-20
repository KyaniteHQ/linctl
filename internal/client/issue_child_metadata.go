package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	IssueID    string                   `json:"issue_id"`
	Identifier string                   `json:"identifier"`
	Comments   []CommentMetadataSummary `json:"comments"`
	Page
}

// IssueCustomerNeedMetadataList is a page of body-free customer needs for one issue-like root.
type IssueCustomerNeedMetadataList struct {
	IssueID    string                        `json:"issue_id"`
	Identifier string                        `json:"identifier"`
	Needs      []CustomerNeedMetadataSummary `json:"customer_needs"`
	Page
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

//nolint:lll
type issueNeedsNode = gql.IssueNeedsProjectionNeedsCustomerNeedConnectionNodesCustomerNeed

//nolint:lll
type issueFormerNeedsNode = gql.IssueFormerNeedsProjectionFormerNeedsCustomerNeedConnectionNodesCustomerNeed

//nolint:lll
type issueVCSBranchCommentsNode = gql.IssueCommentMetadataProjectionCommentsCommentConnectionNodesComment

//nolint:lll
type issueVCSBranchNeedsNode = gql.IssueNeedsProjectionNeedsCustomerNeedConnectionNodesCustomerNeed

//nolint:lll
type issueVCSBranchFormerNeedsNode = gql.IssueFormerNeedsProjectionFormerNeedsCustomerNeedConnectionNodesCustomerNeed

// ListIssueNeeds returns body-free customer needs associated with one issue.
func ListIssueNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list issue customer needs "+id, limit, defaultListPageSize,
		query.needs,
		issueNeedNodeSummary,
	)
	if err != nil {
		return IssueCustomerNeedMetadataList{}, err
	}

	return IssueCustomerNeedMetadataList{
		IssueID:    parent.issueID,
		Identifier: parent.identifier,
		Needs:      page.Items,
		Page:       page.Page,
	}, nil
}

// ListIssueFormerNeeds returns body-free customer needs formerly associated with one issue.
func ListIssueFormerNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list issue former customer needs "+id, limit, defaultListPageSize,
		query.formerNeeds,
		issueFormerNeedNodeSummary,
	)
	if err != nil {
		return IssueCustomerNeedMetadataList{}, err
	}

	return IssueCustomerNeedMetadataList{
		IssueID:    parent.issueID,
		Identifier: parent.identifier,
		Needs:      page.Items,
		Page:       page.Page,
	}, nil
}

// GetIssueSharedAccess returns compact shared-access metadata for one issue.
func GetIssueSharedAccess(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueSharedAccessSummary, error) {
	result, err := gql.XIssue_sharedAccess(ctx, graphqlClient, id)
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
func ListIssueVCSBranchComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueCommentMetadataList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, parent, err := listConnectionWithParent(
		"list issue vcs branch comments "+branchName, limit, defaultListPageSize,
		query.comments,
		issueVCSBranchCommentNodeSummary,
	)
	if err != nil {
		return IssueCommentMetadataList{}, err
	}

	return IssueCommentMetadataList{
		IssueID:    parent.issueID,
		Identifier: parent.identifier,
		Comments:   page.Items,
		Page:       page.Page,
	}, nil
}

// ListIssueVCSBranchNeeds returns body-free customer needs for the issue matched by a VCS branch.
func ListIssueVCSBranchNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, parent, err := listConnectionWithParent(
		"list issue vcs branch customer needs "+branchName, limit, defaultListPageSize,
		query.needs,
		issueVCSBranchNeedNodeSummary,
	)
	if err != nil {
		return IssueCustomerNeedMetadataList{}, err
	}

	return IssueCustomerNeedMetadataList{
		IssueID:    parent.issueID,
		Identifier: parent.identifier,
		Needs:      page.Items,
		Page:       page.Page,
	}, nil
}

// ListIssueVCSBranchFormerNeeds returns body-free former customer needs for the issue matched by a VCS branch.
func ListIssueVCSBranchFormerNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, parent, err := listConnectionWithParent(
		"list issue vcs branch former customer needs "+branchName, limit, defaultListPageSize,
		query.formerNeeds,
		issueVCSBranchFormerNeedNodeSummary,
	)
	if err != nil {
		return IssueCustomerNeedMetadataList{}, err
	}

	return IssueCustomerNeedMetadataList{
		IssueID:    parent.issueID,
		Identifier: parent.identifier,
		Needs:      page.Items,
		Page:       page.Page,
	}, nil
}

// GetIssueVCSBranchSharedAccess returns compact shared-access metadata for an issue matched by a VCS branch.
func GetIssueVCSBranchSharedAccess(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
) (IssueSharedAccessSummary, error) {
	result, err := gql.XIssueVcsBranchSearch_sharedAccess(ctx, graphqlClient, branchName)
	if err != nil {
		return IssueSharedAccessSummary{}, fmt.Errorf("get issue vcs branch shared access %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueSharedAccessSummary{}, notFoundError("get issue vcs branch shared access %s", branchName)
	}

	return issueSharedAccessSummary(
		result.IssueVcsBranchSearch.Id,
		result.IssueVcsBranchSearch.Identifier,
		result.IssueVcsBranchSearch.SharedAccess.IssueSharedAccessFields,
	), nil
}

func (query *issueChildQuery) needs(
	pageSize int,
	after *string,
) ([]issueNeedsNode, issueChildParent, bool, *string, error) {
	result, err := gql.XIssue_needs(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, issueChildParent{}, false, nil, err
	}

	return result.Issue.Needs.Nodes,
		issueChildParent{issueID: result.Issue.Id, identifier: result.Issue.Identifier},
		result.Issue.Needs.PageInfo.HasNextPage,
		result.Issue.Needs.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) formerNeeds(
	pageSize int,
	after *string,
) ([]issueFormerNeedsNode, issueChildParent, bool, *string, error) {
	result, err := gql.XIssue_formerNeeds(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, issueChildParent{}, false, nil, err
	}

	return result.Issue.FormerNeeds.Nodes,
		issueChildParent{issueID: result.Issue.Id, identifier: result.Issue.Identifier},
		result.Issue.FormerNeeds.PageInfo.HasNextPage,
		result.Issue.FormerNeeds.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) comments(
	pageSize int,
	after *string,
) ([]issueVCSBranchCommentsNode, issueVCSBranchParent, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_comments(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, issueVCSBranchParent{}, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, issueVCSBranchParent{}, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Comments.Nodes,
		issueVCSBranchParent{
			issueID:    result.IssueVcsBranchSearch.Id,
			identifier: result.IssueVcsBranchSearch.Identifier,
		},
		result.IssueVcsBranchSearch.Comments.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Comments.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) needs(
	pageSize int,
	after *string,
) ([]issueVCSBranchNeedsNode, issueVCSBranchParent, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_needs(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, issueVCSBranchParent{}, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, issueVCSBranchParent{}, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Needs.Nodes,
		issueVCSBranchParent{
			issueID:    result.IssueVcsBranchSearch.Id,
			identifier: result.IssueVcsBranchSearch.Identifier,
		},
		result.IssueVcsBranchSearch.Needs.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Needs.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) formerNeeds(
	pageSize int,
	after *string,
) ([]issueVCSBranchFormerNeedsNode, issueVCSBranchParent, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_formerNeeds(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, issueVCSBranchParent{}, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, issueVCSBranchParent{}, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.FormerNeeds.Nodes,
		issueVCSBranchParent{
			issueID:    result.IssueVcsBranchSearch.Id,
			identifier: result.IssueVcsBranchSearch.Identifier,
		},
		result.IssueVcsBranchSearch.FormerNeeds.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.FormerNeeds.PageInfo.EndCursor,
		nil
}

func issueNeedNodeSummary(need issueNeedsNode) CustomerNeedMetadataSummary {
	return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
}

func issueFormerNeedNodeSummary(need issueFormerNeedsNode) CustomerNeedMetadataSummary {
	return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
}

func issueVCSBranchCommentNodeSummary(comment issueVCSBranchCommentsNode) CommentMetadataSummary {
	return commentMetadataSummary(comment.CommentMetadataFields)
}

func issueVCSBranchNeedNodeSummary(need issueVCSBranchNeedsNode) CustomerNeedMetadataSummary {
	return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
}

func issueVCSBranchFormerNeedNodeSummary(need issueVCSBranchFormerNeedsNode) CustomerNeedMetadataSummary {
	return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
}

func customerNeedMetadataSummary(fields gql.CustomerNeedMetadataFields) CustomerNeedMetadataSummary {
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
	fields gql.IssueSharedAccessFields,
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

func issueSharedAccessDisallowedFields(fields []gql.IssueSharedAccessDisallowedField) []string {
	values := mapNodes(fields, func(field gql.IssueSharedAccessDisallowedField) string {
		return string(field)
	})

	return values
}
