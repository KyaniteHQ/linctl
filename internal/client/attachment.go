package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// AttachmentSummary is the compact attachment model used by read-only commands.
type AttachmentSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle,omitempty"`
	URL        string `json:"url"`
	SourceType string `json:"source_type,omitempty"`
}

// AttachmentList is a page of attachments.
type AttachmentList struct {
	Attachments []AttachmentSummary `json:"attachments"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

//nolint:lll
type attachmentsNode = gql.XAttachmentsAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type attachmentsForURLNode = gql.XAttachmentsForURLAttachmentsForURLAttachmentConnectionNodesAttachment

//nolint:lll
type attachmentIssueAttachmentsNode = gql.IssueAttachmentsProjectionAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type attachmentIssueChildrenNode = gql.IssueChildrenProjectionChildrenIssueConnectionNodesIssue

//nolint:lll
type attachmentIssueDocumentsNode = gql.IssueDocumentsProjectionDocumentsDocumentConnectionNodesDocument

//nolint:lll
type attachmentIssueFormerAttachmentsNode = gql.IssueFormerAttachmentsProjectionFormerAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type attachmentIssueCommentsNode = gql.IssueCommentMetadataProjectionCommentsCommentConnectionNodesComment

//nolint:lll
type attachmentIssueNeedsNode = gql.IssueNeedsProjectionNeedsCustomerNeedConnectionNodesCustomerNeed

//nolint:lll
type attachmentIssueFormerNeedsNode = gql.IssueFormerNeedsProjectionFormerNeedsCustomerNeedConnectionNodesCustomerNeed

//nolint:lll
type attachmentIssueHistoryNode = gql.IssueHistoryProjectionHistoryIssueHistoryConnectionNodesIssueHistory

//nolint:lll
type attachmentIssueInverseRelationsNode = gql.IssueInverseRelationsProjectionInverseRelationsIssueRelationConnectionNodesIssueRelation

//nolint:lll
type attachmentIssueLabelsNode = gql.IssueLabelsProjectionLabelsIssueLabelConnectionNodesIssueLabel

//nolint:lll
type attachmentIssueRelationsNode = gql.IssueRelationsProjectionRelationsIssueRelationConnectionNodesIssueRelation

//nolint:lll
type attachmentIssueReleasesNode = gql.IssueReleasesProjectionReleasesReleaseConnectionNodesRelease

//nolint:lll
type attachmentIssueStateHistoryNode = gql.IssueStateHistoryProjectionStateHistoryIssueStateSpanConnectionNodesIssueStateSpan

//nolint:lll
type attachmentIssueSubscribersNode = gql.IssueSubscribersProjectionSubscribersUserConnectionNodesUser

type attachmentsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type attachmentsForURLQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	url           string
}

type attachmentIssueQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	issueID       string
	identifier    string
}

// ListAttachments returns visible issue attachments.
func ListAttachments(ctx context.Context, graphqlClient graphql.Client, limit int) (AttachmentList, error) {
	query := attachmentsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list attachments", limit, defaultListPageSize,
		query.page,
		attachmentNodeSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentsForURL returns visible issue attachments linked to a URL.
func ListAttachmentsForURL(
	ctx context.Context,
	graphqlClient graphql.Client,
	url string,
	limit int,
) (AttachmentList, error) {
	query := attachmentsForURLQuery{ctx: ctx, graphqlClient: graphqlClient, url: url}
	page, err := listConnection(
		"list attachments for url "+url, limit, defaultListPageSize,
		query.page,
		attachmentForURLNodeSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetAttachmentByID returns one attachment by Linear id.
func GetAttachmentByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (AttachmentSummary, error) {
	result, err := gql.XAttachment(ctx, graphqlClient, id)
	if err != nil {
		return AttachmentSummary{}, fmt.Errorf("get attachment %s: %w", id, err)
	}

	return attachmentSummary(result.Attachment.AttachmentSummaryFields), nil
}

// GetAttachmentIssue returns the issue associated with one attachment.
func GetAttachmentIssue(ctx context.Context, graphqlClient graphql.Client, id string) (IssueSummary, error) {
	result, err := gql.XAttachmentIssue(ctx, graphqlClient, id)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("get attachment issue %s: %w", id, err)
	}

	return issueSummaryFromFields(result.AttachmentIssue.IssueSummaryFields), nil
}

// ListAttachmentIssueAttachments returns attachments for the issue associated with one attachment.
func ListAttachmentIssueAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue attachments "+id, limit, defaultListPageSize,
		query.attachments,
		attachmentIssueAttachmentSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetAttachmentIssueBotActor returns the issue bot actor associated with one attachment.
func GetAttachmentIssueBotActor(ctx context.Context, graphqlClient graphql.Client, id string) (IssueBotActor, error) {
	result, err := gql.XAttachmentIssue_botActor(ctx, graphqlClient, id)
	if err != nil {
		return IssueBotActor{}, fmt.Errorf("get attachment issue bot actor %s: %w", id, err)
	}

	return IssueBotActor{
		IssueID: result.AttachmentIssue.Id,
		Bot:     issueActorBotSummary(result.AttachmentIssue.BotActor),
	}, nil
}

// ListAttachmentIssueChildren returns child issues for the issue associated with one attachment.
func ListAttachmentIssueChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue children "+id, limit, defaultListPageSize,
		query.children,
		attachmentIssueChildSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueDocuments returns documents for the issue associated with one attachment.
func ListAttachmentIssueDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue documents "+id, limit, defaultListPageSize,
		query.documents,
		attachmentIssueDocumentSummary,
	)
	if err != nil {
		return DocumentList{}, err
	}

	return DocumentList{Documents: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueFormerAttachments returns former attachments for the issue associated with one attachment.
func ListAttachmentIssueFormerAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue former attachments "+id, limit, defaultListPageSize,
		query.formerAttachments,
		attachmentIssueFormerAttachmentSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueComments returns body-free comments for the issue associated with one attachment.
func ListAttachmentIssueComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCommentMetadataList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue comments "+id, limit, defaultListPageSize,
		query.comments,
		attachmentIssueCommentSummary,
	)
	if err != nil {
		return IssueCommentMetadataList{}, err
	}

	return IssueCommentMetadataList{
		IssueID:     query.issueID,
		Identifier:  query.identifier,
		Comments:    page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListAttachmentIssueNeeds returns body-free customer needs for the issue associated with one attachment.
func ListAttachmentIssueNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue customer needs "+id, limit, defaultListPageSize,
		query.needs,
		attachmentIssueNeedSummary,
	)
	if err != nil {
		return IssueCustomerNeedMetadataList{}, err
	}

	return IssueCustomerNeedMetadataList{
		IssueID:     query.issueID,
		Identifier:  query.identifier,
		Needs:       page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListAttachmentIssueFormerNeeds returns body-free former customer needs for the issue associated with one attachment.
func ListAttachmentIssueFormerNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue former customer needs "+id, limit, defaultListPageSize,
		query.formerNeeds,
		attachmentIssueFormerNeedSummary,
	)
	if err != nil {
		return IssueCustomerNeedMetadataList{}, err
	}

	return IssueCustomerNeedMetadataList{
		IssueID:     query.issueID,
		Identifier:  query.identifier,
		Needs:       page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListAttachmentIssueHistory returns history metadata for the issue associated with one attachment.
func ListAttachmentIssueHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueHistoryList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue history "+id, limit, defaultListPageSize,
		query.history,
		issueHistorySummary,
	)
	if err != nil {
		return IssueHistoryList{}, err
	}

	return IssueHistoryList{History: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueInverseRelations returns inverse relations for the issue associated with one attachment.
func ListAttachmentIssueInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue inverse relations "+id, limit, defaultListPageSize,
		query.inverseRelations,
		attachmentIssueInverseRelationSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueLabels returns labels for the issue associated with one attachment.
func ListAttachmentIssueLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (LabelList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue labels "+id, limit, defaultListPageSize,
		query.labels,
		attachmentIssueLabelSummary,
	)
	if err != nil {
		return LabelList{}, err
	}

	return LabelList{Labels: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueRelations returns relations for the issue associated with one attachment.
func ListAttachmentIssueRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue relations "+id, limit, defaultListPageSize,
		query.relations,
		attachmentIssueRelationSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListAttachmentIssueReleases returns releases for the issue associated with one attachment.
func ListAttachmentIssueReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue releases "+id, limit, defaultListPageSize,
		query.releases,
		attachmentIssueReleaseSummary,
	)
	if err != nil {
		return ReleaseList{}, err
	}

	return ReleaseList{Releases: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetAttachmentIssueSharedAccess returns compact shared-access metadata for the issue associated with one attachment.
func GetAttachmentIssueSharedAccess(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (IssueSharedAccessSummary, error) {
	result, err := gql.XAttachmentIssue_sharedAccess(ctx, graphqlClient, id)
	if err != nil {
		return IssueSharedAccessSummary{}, fmt.Errorf("get attachment issue shared access %s: %w", id, err)
	}

	return issueSharedAccessSummary(
		result.AttachmentIssue.Id,
		result.AttachmentIssue.Identifier,
		result.AttachmentIssue.SharedAccess.IssueSharedAccessFields,
	), nil
}

// ListAttachmentIssueStateHistory returns workflow-state spans for the issue associated with one attachment.
func ListAttachmentIssueStateHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueStateHistoryList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue state history "+id, limit, defaultListPageSize,
		query.stateHistory,
		issueStateSpanSummary,
	)
	if err != nil {
		return IssueStateHistoryList{}, err
	}

	return IssueStateHistoryList{
		IssueID:     query.issueID,
		Spans:       page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListAttachmentIssueSubscribers returns subscribers for the issue associated with one attachment.
func ListAttachmentIssueSubscribers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (UserList, error) {
	query := &attachmentIssueQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list attachment issue subscribers "+id, limit, defaultListPageSize,
		query.subscribers,
		attachmentIssueSubscriberSummary,
	)
	if err != nil {
		return UserList{}, err
	}

	return UserList{Users: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

func (query attachmentsQuery) page(pageSize int, after *string) ([]attachmentsNode, bool, *string, error) {
	result, err := gql.XAttachments(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Attachments.Nodes, result.Attachments.PageInfo.HasNextPage, result.Attachments.PageInfo.EndCursor, nil
}

func (query attachmentsForURLQuery) page(
	pageSize int,
	after *string,
) ([]attachmentsForURLNode, bool, *string, error) {
	result, err := gql.XAttachmentsForURL(
		query.ctx, query.graphqlClient, query.url, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentsForURL.Nodes,
		result.AttachmentsForURL.PageInfo.HasNextPage,
		result.AttachmentsForURL.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) attachments(
	pageSize int,
	after *string,
) ([]attachmentIssueAttachmentsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_attachments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Attachments.Nodes,
		result.AttachmentIssue.Attachments.PageInfo.HasNextPage,
		result.AttachmentIssue.Attachments.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) children(
	pageSize int,
	after *string,
) ([]attachmentIssueChildrenNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_children(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Children.Nodes,
		result.AttachmentIssue.Children.PageInfo.HasNextPage,
		result.AttachmentIssue.Children.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) documents(
	pageSize int,
	after *string,
) ([]attachmentIssueDocumentsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_documents(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Documents.Nodes,
		result.AttachmentIssue.Documents.PageInfo.HasNextPage,
		result.AttachmentIssue.Documents.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) formerAttachments(
	pageSize int,
	after *string,
) ([]attachmentIssueFormerAttachmentsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_formerAttachments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.FormerAttachments.Nodes,
		result.AttachmentIssue.FormerAttachments.PageInfo.HasNextPage,
		result.AttachmentIssue.FormerAttachments.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) comments(
	pageSize int,
	after *string,
) ([]attachmentIssueCommentsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_comments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.issueID = result.AttachmentIssue.Id
	query.identifier = result.AttachmentIssue.Identifier

	return result.AttachmentIssue.Comments.Nodes,
		result.AttachmentIssue.Comments.PageInfo.HasNextPage,
		result.AttachmentIssue.Comments.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) needs(
	pageSize int,
	after *string,
) ([]attachmentIssueNeedsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_needs(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.issueID = result.AttachmentIssue.Id
	query.identifier = result.AttachmentIssue.Identifier

	return result.AttachmentIssue.Needs.Nodes,
		result.AttachmentIssue.Needs.PageInfo.HasNextPage,
		result.AttachmentIssue.Needs.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) formerNeeds(
	pageSize int,
	after *string,
) ([]attachmentIssueFormerNeedsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_formerNeeds(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.issueID = result.AttachmentIssue.Id
	query.identifier = result.AttachmentIssue.Identifier

	return result.AttachmentIssue.FormerNeeds.Nodes,
		result.AttachmentIssue.FormerNeeds.PageInfo.HasNextPage,
		result.AttachmentIssue.FormerNeeds.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) history(
	pageSize int,
	after *string,
) ([]attachmentIssueHistoryNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_history(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.History.Nodes,
		result.AttachmentIssue.History.PageInfo.HasNextPage,
		result.AttachmentIssue.History.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) inverseRelations(
	pageSize int,
	after *string,
) ([]attachmentIssueInverseRelationsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_inverseRelations(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.InverseRelations.Nodes,
		result.AttachmentIssue.InverseRelations.PageInfo.HasNextPage,
		result.AttachmentIssue.InverseRelations.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) labels(
	pageSize int,
	after *string,
) ([]attachmentIssueLabelsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_labels(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Labels.Nodes,
		result.AttachmentIssue.Labels.PageInfo.HasNextPage,
		result.AttachmentIssue.Labels.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) relations(
	pageSize int,
	after *string,
) ([]attachmentIssueRelationsNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_relations(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Relations.Nodes,
		result.AttachmentIssue.Relations.PageInfo.HasNextPage,
		result.AttachmentIssue.Relations.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) releases(
	pageSize int,
	after *string,
) ([]attachmentIssueReleasesNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_releases(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Releases.Nodes,
		result.AttachmentIssue.Releases.PageInfo.HasNextPage,
		result.AttachmentIssue.Releases.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) stateHistory(
	pageSize int,
	after *string,
) ([]attachmentIssueStateHistoryNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_stateHistory(query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after)
	if err != nil {
		return nil, false, nil, err
	}

	query.issueID = result.AttachmentIssue.Id

	return result.AttachmentIssue.StateHistory.Nodes,
		result.AttachmentIssue.StateHistory.PageInfo.HasNextPage,
		result.AttachmentIssue.StateHistory.PageInfo.EndCursor,
		nil
}

func (query *attachmentIssueQuery) subscribers(
	pageSize int,
	after *string,
) ([]attachmentIssueSubscribersNode, bool, *string, error) {
	result, err := gql.XAttachmentIssue_subscribers(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.AttachmentIssue.Subscribers.Nodes,
		result.AttachmentIssue.Subscribers.PageInfo.HasNextPage,
		result.AttachmentIssue.Subscribers.PageInfo.EndCursor,
		nil
}

func attachmentNodeSummary(node attachmentsNode) AttachmentSummary {
	return attachmentSummary(node.AttachmentSummaryFields)
}

func attachmentForURLNodeSummary(node attachmentsForURLNode) AttachmentSummary {
	return attachmentSummary(node.AttachmentSummaryFields)
}

func attachmentIssueAttachmentSummary(attachment attachmentIssueAttachmentsNode) AttachmentSummary {
	return attachmentSummary(attachment.AttachmentSummaryFields)
}

func attachmentIssueChildSummary(issue attachmentIssueChildrenNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func attachmentIssueDocumentSummary(document attachmentIssueDocumentsNode) DocumentSummary {
	return documentSummary(document.DocumentSummaryFields)
}

func attachmentIssueFormerAttachmentSummary(attachment attachmentIssueFormerAttachmentsNode) AttachmentSummary {
	return attachmentSummary(attachment.AttachmentSummaryFields)
}

func attachmentIssueCommentSummary(comment attachmentIssueCommentsNode) CommentMetadataSummary {
	return commentMetadataSummary(comment.CommentMetadataFields)
}

func attachmentIssueNeedSummary(need attachmentIssueNeedsNode) CustomerNeedMetadataSummary {
	return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
}

func attachmentIssueFormerNeedSummary(need attachmentIssueFormerNeedsNode) CustomerNeedMetadataSummary {
	return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
}

func attachmentIssueInverseRelationSummary(node attachmentIssueInverseRelationsNode) IssueRelationSummary {
	return issueRelationSummary(node.IssueRelationSummaryFields)
}

func attachmentIssueLabelSummary(label attachmentIssueLabelsNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func attachmentIssueRelationSummary(relation attachmentIssueRelationsNode) IssueRelationSummary {
	return issueRelationSummary(relation.IssueRelationSummaryFields)
}

func attachmentIssueReleaseSummary(release attachmentIssueReleasesNode) ReleaseSummary {
	return releaseSummary(release.ReleaseSummaryFields)
}

func attachmentIssueSubscriberSummary(node attachmentIssueSubscribersNode) UserSummary {
	return userSummary(node.UserSummaryFields)
}

func attachmentSummary(fields gql.AttachmentSummaryFields) AttachmentSummary {
	return AttachmentSummary{
		ID:         fields.Id,
		Title:      fields.Title,
		Subtitle:   stringValue(fields.Subtitle),
		URL:        fields.Url,
		SourceType: stringValue(fields.SourceType),
	}
}
