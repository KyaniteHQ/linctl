package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

//nolint:lll
type issueVCSBranchAttachmentsNode = gql.IssueAttachmentsProjectionAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type issueVCSBranchChildrenNode = gql.IssueChildrenProjectionChildrenIssueConnectionNodesIssue

//nolint:lll
type issueVCSBranchDocumentsNode = gql.IssueDocumentsProjectionDocumentsDocumentConnectionNodesDocument

//nolint:lll
type issueVCSBranchFormerAttachmentsNode = gql.IssueFormerAttachmentsProjectionFormerAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type issueVCSBranchHistoryNode = gql.IssueHistoryProjectionHistoryIssueHistoryConnectionNodesIssueHistory

//nolint:lll
type issueVCSBranchInverseRelationsNode = gql.IssueInverseRelationsProjectionInverseRelationsIssueRelationConnectionNodesIssueRelation

//nolint:lll
type issueVCSBranchLabelsNode = gql.IssueLabelsProjectionLabelsIssueLabelConnectionNodesIssueLabel

//nolint:lll
type issueVCSBranchRelationsNode = gql.IssueRelationsProjectionRelationsIssueRelationConnectionNodesIssueRelation

//nolint:lll
type issueVCSBranchReleasesNode = gql.IssueReleasesProjectionReleasesReleaseConnectionNodesRelease

//nolint:lll
type issueVCSBranchStateHistoryNode = gql.IssueStateHistoryProjectionStateHistoryIssueStateSpanConnectionNodesIssueStateSpan

//nolint:lll
type issueVCSBranchSubscribersNode = gql.IssueSubscribersProjectionSubscribersUserConnectionNodesUser

type issueVCSBranchQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	branchName    string
}

// GetIssueByVCSBranch returns an issue by VCS branch name.
func GetIssueByVCSBranch(ctx context.Context, graphqlClient graphql.Client, branchName string) (IssueSummary, error) {
	result, err := gql.XIssueVcsBranchSearch(ctx, graphqlClient, branchName)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("get issue by vcs branch %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueSummary{}, notFoundError("get issue by vcs branch %s", branchName)
	}

	return issueSummaryFromFields(result.IssueVcsBranchSearch.IssueSummaryFields), nil
}

// ListIssueVCSBranchAttachments returns attachments for the issue matched by a VCS branch.
func ListIssueVCSBranchAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (AttachmentList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch attachments "+branchName, limit, defaultListPageSize,
		query.attachments,
		issueVCSBranchAttachmentNodeSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, Page: page.Page}, nil
}

// GetIssueVCSBranchBotActor returns bot actor metadata for the issue matched by a VCS branch.
func GetIssueVCSBranchBotActor(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
) (IssueBotActor, error) {
	result, err := gql.XIssueVcsBranchSearch_botActor(ctx, graphqlClient, branchName)
	if err != nil {
		return IssueBotActor{}, fmt.Errorf("get issue vcs branch bot actor %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueBotActor{}, notFoundError("get issue vcs branch bot actor %s", branchName)
	}

	return IssueBotActor{
		IssueID: result.IssueVcsBranchSearch.Id,
		Bot:     issueActorBotSummary(result.IssueVcsBranchSearch.BotActor),
	}, nil
}

// ListIssueVCSBranchChildren returns child issues for the issue matched by a VCS branch.
func ListIssueVCSBranchChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch children "+branchName, limit, defaultListPageSize,
		query.children,
		issueVCSBranchChildNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchDocuments returns documents for the issue matched by a VCS branch.
func ListIssueVCSBranchDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (DocumentList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch documents "+branchName, limit, defaultListPageSize,
		query.documents,
		issueVCSBranchDocumentNodeSummary,
	)
	if err != nil {
		return DocumentList{}, err
	}

	return DocumentList{Documents: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchFormerAttachments returns former attachments for the issue matched by a VCS branch.
func ListIssueVCSBranchFormerAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (AttachmentList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch former attachments "+branchName, limit, defaultListPageSize,
		query.formerAttachments,
		issueVCSBranchFormerAttachmentNodeSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchHistory returns history metadata for the issue matched by a VCS branch.
func ListIssueVCSBranchHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueHistoryList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch history "+branchName, limit, defaultListPageSize,
		query.history,
		issueHistorySummary,
	)
	if err != nil {
		return IssueHistoryList{}, err
	}

	return IssueHistoryList{History: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchInverseRelations returns inverse relations for the issue matched by a VCS branch.
func ListIssueVCSBranchInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueRelationList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch inverse relations "+branchName, limit, defaultListPageSize,
		query.inverseRelations,
		issueVCSBranchInverseRelationNodeSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchLabels returns labels for the issue matched by a VCS branch.
func ListIssueVCSBranchLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (LabelList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch labels "+branchName, limit, defaultListPageSize,
		query.labels,
		issueVCSBranchLabelNodeSummary,
	)
	if err != nil {
		return LabelList{}, err
	}

	return LabelList{Labels: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchRelations returns relations for the issue matched by a VCS branch.
func ListIssueVCSBranchRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueRelationList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch relations "+branchName, limit, defaultListPageSize,
		query.relations,
		issueVCSBranchRelationNodeSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchReleases returns releases for the issue matched by a VCS branch.
func ListIssueVCSBranchReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (ReleaseList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch releases "+branchName, limit, defaultListPageSize,
		query.releases,
		issueVCSBranchReleaseNodeSummary,
	)
	if err != nil {
		return ReleaseList{}, err
	}

	return ReleaseList{Releases: page.Items, Page: page.Page}, nil
}

// ListIssueVCSBranchStateHistory returns workflow-state spans for the issue matched by a VCS branch.
func ListIssueVCSBranchStateHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueStateHistoryList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, parent, err := listConnectionWithParent(
		"list issue vcs branch state history "+branchName, limit, defaultListPageSize,
		query.stateHistory,
		issueStateSpanSummary,
	)
	if err != nil {
		return IssueStateHistoryList{}, err
	}

	return IssueStateHistoryList{
		IssueID: parent.issueID,
		Spans:   page.Items,
		Page:    page.Page,
	}, nil
}

// ListIssueVCSBranchSubscribers returns subscribers for the issue matched by a VCS branch.
func ListIssueVCSBranchSubscribers(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (UserList, error) {
	query := &issueVCSBranchQuery{ctx: ctx, graphqlClient: graphqlClient, branchName: branchName}
	page, err := listConnection(
		"list issue vcs branch subscribers "+branchName, limit, defaultListPageSize,
		query.subscribers,
		issueVCSBranchSubscriberNodeSummary,
	)
	if err != nil {
		return UserList{}, err
	}

	return UserList{Users: page.Items, Page: page.Page}, nil
}

func (query *issueVCSBranchQuery) attachments(
	pageSize int,
	after *string,
) ([]issueVCSBranchAttachmentsNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_attachments(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Attachments.Nodes,
		result.IssueVcsBranchSearch.Attachments.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Attachments.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) children(
	pageSize int,
	after *string,
) ([]issueVCSBranchChildrenNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_children(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Children.Nodes,
		result.IssueVcsBranchSearch.Children.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Children.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) documents(
	pageSize int,
	after *string,
) ([]issueVCSBranchDocumentsNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_documents(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Documents.Nodes,
		result.IssueVcsBranchSearch.Documents.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Documents.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) formerAttachments(
	pageSize int,
	after *string,
) ([]issueVCSBranchFormerAttachmentsNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_formerAttachments(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.FormerAttachments.Nodes,
		result.IssueVcsBranchSearch.FormerAttachments.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.FormerAttachments.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) history(
	pageSize int,
	after *string,
) ([]issueVCSBranchHistoryNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_history(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.History.Nodes,
		result.IssueVcsBranchSearch.History.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.History.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) inverseRelations(
	pageSize int,
	after *string,
) ([]issueVCSBranchInverseRelationsNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_inverseRelations(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.InverseRelations.Nodes,
		result.IssueVcsBranchSearch.InverseRelations.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.InverseRelations.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) labels(
	pageSize int,
	after *string,
) ([]issueVCSBranchLabelsNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_labels(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Labels.Nodes,
		result.IssueVcsBranchSearch.Labels.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Labels.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) relations(
	pageSize int,
	after *string,
) ([]issueVCSBranchRelationsNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_relations(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Relations.Nodes,
		result.IssueVcsBranchSearch.Relations.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Relations.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) releases(
	pageSize int,
	after *string,
) ([]issueVCSBranchReleasesNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_releases(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Releases.Nodes,
		result.IssueVcsBranchSearch.Releases.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Releases.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) stateHistory(
	pageSize int,
	after *string,
) ([]issueVCSBranchStateHistoryNode, issueParent, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_stateHistory(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after,
	)
	if err != nil {
		return nil, issueParent{}, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, issueParent{}, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.StateHistory.Nodes,
		issueParent{issueID: result.IssueVcsBranchSearch.Id},
		result.IssueVcsBranchSearch.StateHistory.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.StateHistory.PageInfo.EndCursor,
		nil
}

func (query *issueVCSBranchQuery) subscribers(
	pageSize int,
	after *string,
) ([]issueVCSBranchSubscribersNode, bool, *string, error) {
	result, err := gql.XIssueVcsBranchSearch_subscribers(
		query.ctx, query.graphqlClient, query.branchName, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}
	if result.IssueVcsBranchSearch == nil {
		return nil, false, nil, ErrNotFound
	}

	return result.IssueVcsBranchSearch.Subscribers.Nodes,
		result.IssueVcsBranchSearch.Subscribers.PageInfo.HasNextPage,
		result.IssueVcsBranchSearch.Subscribers.PageInfo.EndCursor,
		nil
}

func issueVCSBranchAttachmentNodeSummary(node issueVCSBranchAttachmentsNode) AttachmentSummary {
	return attachmentSummary(node.AttachmentSummaryFields)
}

func issueVCSBranchChildNodeSummary(issue issueVCSBranchChildrenNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func issueVCSBranchDocumentNodeSummary(document issueVCSBranchDocumentsNode) DocumentSummary {
	return documentSummary(document.DocumentSummaryFields)
}

func issueVCSBranchFormerAttachmentNodeSummary(attachment issueVCSBranchFormerAttachmentsNode) AttachmentSummary {
	return attachmentSummary(attachment.AttachmentSummaryFields)
}

func issueVCSBranchInverseRelationNodeSummary(relation issueVCSBranchInverseRelationsNode) IssueRelationSummary {
	return issueRelationSummary(relation.IssueRelationSummaryFields)
}

func issueVCSBranchLabelNodeSummary(label issueVCSBranchLabelsNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func issueVCSBranchRelationNodeSummary(node issueVCSBranchRelationsNode) IssueRelationSummary {
	return issueRelationSummary(node.IssueRelationSummaryFields)
}

func issueVCSBranchReleaseNodeSummary(release issueVCSBranchReleasesNode) ReleaseSummary {
	return releaseSummary(release.ReleaseSummaryFields)
}

func issueVCSBranchSubscriberNodeSummary(node issueVCSBranchSubscribersNode) UserSummary {
	return userSummary(node.UserSummaryFields)
}
