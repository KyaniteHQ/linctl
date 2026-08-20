package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

//nolint:lll
type issueAttachmentsNode = gql.IssueAttachmentsProjectionAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type issueChildrenNode = gql.IssueChildrenProjectionChildrenIssueConnectionNodesIssue

//nolint:lll
type issueDocumentsNode = gql.IssueDocumentsProjectionDocumentsDocumentConnectionNodesDocument

//nolint:lll
type issueFormerAttachmentsNode = gql.IssueFormerAttachmentsProjectionFormerAttachmentsAttachmentConnectionNodesAttachment

//nolint:lll
type issueHistoryNode = gql.IssueHistoryProjectionHistoryIssueHistoryConnectionNodesIssueHistory

//nolint:lll
type issueInverseRelationsNode = gql.IssueInverseRelationsProjectionInverseRelationsIssueRelationConnectionNodesIssueRelation

//nolint:lll
type issueLabelsNode = gql.IssueLabelsProjectionLabelsIssueLabelConnectionNodesIssueLabel

//nolint:lll
type issueRelationsForIssueNode = gql.IssueRelationsProjectionRelationsIssueRelationConnectionNodesIssueRelation

//nolint:lll
type issueReleasesNode = gql.IssueReleasesProjectionReleasesReleaseConnectionNodesRelease

//nolint:lll
type issueStateHistoryNode = gql.IssueStateHistoryProjectionStateHistoryIssueStateSpanConnectionNodesIssueStateSpan

//nolint:lll
type issueSubscribersNode = gql.IssueSubscribersProjectionSubscribersUserConnectionNodesUser

type issueChildQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// issueChildParent is the connection parent metadata issueChildQuery reads out of
// every page. Linear repeats it per page, so the last page wins.
type issueChildParent struct {
	issueID    string
	identifier string
}

// ListIssueAttachments returns attachments associated with one issue.
func ListIssueAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue attachments "+id, limit, defaultListPageSize,
		query.attachments,
		issueAttachmentNodeSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, Page: page.Page}, nil
}

// GetIssueBotActor returns the bot actor that created an issue, when present.
func GetIssueBotActor(ctx context.Context, graphqlClient graphql.Client, id string) (IssueBotActor, error) {
	result, err := gql.XIssue_botActor(ctx, graphqlClient, id)
	if err != nil {
		return IssueBotActor{}, fmt.Errorf("get issue bot actor %s: %w", id, err)
	}

	return IssueBotActor{
		IssueID: result.Issue.Id,
		Bot:     issueActorBotSummary(result.Issue.BotActor),
	}, nil
}

// ListIssueChildren returns child issues for one issue.
func ListIssueChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue children "+id, limit, defaultListPageSize,
		query.children,
		issueChildNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListIssueDocuments returns documents associated with one issue.
func ListIssueDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue documents "+id, limit, defaultListPageSize,
		query.documents,
		issueDocumentNodeSummary,
	)
	if err != nil {
		return DocumentList{}, err
	}

	return DocumentList{Documents: page.Items, Page: page.Page}, nil
}

// ListIssueFormerAttachments returns attachments formerly associated with one issue.
func ListIssueFormerAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue former attachments "+id, limit, defaultListPageSize,
		query.formerAttachments,
		issueFormerAttachmentNodeSummary,
	)
	if err != nil {
		return AttachmentList{}, err
	}

	return AttachmentList{Attachments: page.Items, Page: page.Page}, nil
}

// ListIssueHistory returns compact history metadata for one issue.
func ListIssueHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueHistoryList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue history "+id, limit, defaultListPageSize,
		query.history,
		issueHistorySummary,
	)
	if err != nil {
		return IssueHistoryList{}, err
	}

	return IssueHistoryList{History: page.Items, Page: page.Page}, nil
}

// ListIssueInverseRelations returns inverse relations associated with one issue.
func ListIssueInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue inverse relations "+id, limit, defaultListPageSize,
		query.inverseRelations,
		issueInverseRelationNodeSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, Page: page.Page}, nil
}

// ListIssueLabels returns labels associated with one issue.
func ListIssueLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (LabelList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue labels "+id, limit, defaultListPageSize,
		query.labels,
		issueLabelNodeSummary,
	)
	if err != nil {
		return LabelList{}, err
	}

	return LabelList{Labels: page.Items, Page: page.Page}, nil
}

// ListIssueRelationsForIssue returns relations associated with one issue.
func ListIssueRelationsForIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue relations "+id, limit, defaultListPageSize,
		query.relations,
		issueRelationForIssueNodeSummary,
	)
	if err != nil {
		return IssueRelationList{}, err
	}

	return IssueRelationList{Relations: page.Items, Page: page.Page}, nil
}

// ListIssueReleases returns releases associated with one issue.
func ListIssueReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue releases "+id, limit, defaultListPageSize,
		query.releases,
		issueReleaseNodeSummary,
	)
	if err != nil {
		return ReleaseList{}, err
	}

	return ReleaseList{Releases: page.Items, Page: page.Page}, nil
}

// ListIssueStateHistory returns workflow-state spans for one issue.
func ListIssueStateHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueStateHistoryList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list issue state history "+id, limit, defaultListPageSize,
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

// ListIssueSubscribers returns users subscribed to one issue.
func ListIssueSubscribers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (UserList, error) {
	query := &issueChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue subscribers "+id, limit, defaultListPageSize,
		query.subscribers,
		issueSubscriberNodeSummary,
	)
	if err != nil {
		return UserList{}, err
	}

	return UserList{Users: page.Items, Page: page.Page}, nil
}

func (query *issueChildQuery) attachments(
	pageSize int,
	after *string,
) ([]issueAttachmentsNode, bool, *string, error) {
	result, err := gql.XIssue_attachments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Attachments.Nodes,
		result.Issue.Attachments.PageInfo.HasNextPage,
		result.Issue.Attachments.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) children(pageSize int, after *string) ([]issueChildrenNode, bool, *string, error) {
	result, err := gql.XIssue_children(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Children.Nodes,
		result.Issue.Children.PageInfo.HasNextPage,
		result.Issue.Children.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) documents(pageSize int, after *string) ([]issueDocumentsNode, bool, *string, error) {
	result, err := gql.XIssue_documents(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Documents.Nodes,
		result.Issue.Documents.PageInfo.HasNextPage,
		result.Issue.Documents.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) formerAttachments(
	pageSize int,
	after *string,
) ([]issueFormerAttachmentsNode, bool, *string, error) {
	result, err := gql.XIssue_formerAttachments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.FormerAttachments.Nodes,
		result.Issue.FormerAttachments.PageInfo.HasNextPage,
		result.Issue.FormerAttachments.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) history(pageSize int, after *string) ([]issueHistoryNode, bool, *string, error) {
	result, err := gql.XIssue_history(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.History.Nodes,
		result.Issue.History.PageInfo.HasNextPage,
		result.Issue.History.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) inverseRelations(
	pageSize int,
	after *string,
) ([]issueInverseRelationsNode, bool, *string, error) {
	result, err := gql.XIssue_inverseRelations(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.InverseRelations.Nodes,
		result.Issue.InverseRelations.PageInfo.HasNextPage,
		result.Issue.InverseRelations.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) labels(pageSize int, after *string) ([]issueLabelsNode, bool, *string, error) {
	result, err := gql.XIssue_labels(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Labels.Nodes,
		result.Issue.Labels.PageInfo.HasNextPage,
		result.Issue.Labels.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) relations(
	pageSize int,
	after *string,
) ([]issueRelationsForIssueNode, bool, *string, error) {
	result, err := gql.XIssue_relations(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Relations.Nodes,
		result.Issue.Relations.PageInfo.HasNextPage,
		result.Issue.Relations.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) releases(pageSize int, after *string) ([]issueReleasesNode, bool, *string, error) {
	result, err := gql.XIssue_releases(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Releases.Nodes,
		result.Issue.Releases.PageInfo.HasNextPage,
		result.Issue.Releases.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) stateHistory(
	pageSize int,
	after *string,
) ([]issueStateHistoryNode, issueChildParent, bool, *string, error) {
	result, err := gql.XIssue_stateHistory(query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after)
	if err != nil {
		return nil, issueChildParent{}, false, nil, err
	}

	return result.Issue.StateHistory.Nodes,
		issueChildParent{issueID: result.Issue.Id},
		result.Issue.StateHistory.PageInfo.HasNextPage,
		result.Issue.StateHistory.PageInfo.EndCursor,
		nil
}

func (query *issueChildQuery) subscribers(pageSize int, after *string) ([]issueSubscribersNode, bool, *string, error) {
	result, err := gql.XIssue_subscribers(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Issue.Subscribers.Nodes,
		result.Issue.Subscribers.PageInfo.HasNextPage,
		result.Issue.Subscribers.PageInfo.EndCursor,
		nil
}

func issueAttachmentNodeSummary(attachment issueAttachmentsNode) AttachmentSummary {
	return attachmentSummary(attachment.AttachmentSummaryFields)
}

func issueChildNodeSummary(issue issueChildrenNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func issueDocumentNodeSummary(document issueDocumentsNode) DocumentSummary {
	return documentSummary(document.DocumentSummaryFields)
}

func issueFormerAttachmentNodeSummary(attachment issueFormerAttachmentsNode) AttachmentSummary {
	return attachmentSummary(attachment.AttachmentSummaryFields)
}

func issueInverseRelationNodeSummary(relation issueInverseRelationsNode) IssueRelationSummary {
	return issueRelationSummary(relation.IssueRelationSummaryFields)
}

func issueLabelNodeSummary(label issueLabelsNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func issueRelationForIssueNodeSummary(relation issueRelationsForIssueNode) IssueRelationSummary {
	return issueRelationSummary(relation.IssueRelationSummaryFields)
}

func issueReleaseNodeSummary(release issueReleasesNode) ReleaseSummary {
	return releaseSummary(release.ReleaseSummaryFields)
}

func issueSubscriberNodeSummary(node issueSubscribersNode) UserSummary {
	return userSummary(node.UserSummaryFields)
}

func issueHistorySummary(
	history gql.IssueHistoryProjectionHistoryIssueHistoryConnectionNodesIssueHistory,
) IssueHistorySummary {
	return IssueHistorySummary{
		ID:                 history.Id,
		IssueID:            history.Issue.Id,
		ActorID:            stringValue(history.ActorId),
		UpdatedDescription: boolValue(history.UpdatedDescription),
		CreatedAt:          history.CreatedAt,
		UpdatedAt:          history.UpdatedAt,
		ArchivedAt:         stringValue(history.ArchivedAt),
	}
}

func issueActorBotSummary(bot *gql.IssueBotActorProjectionBotActorActorBot) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return actorBotSummary(&bot.ActorBotSummaryFields)
}

func issueStateSpanSummary(
	span gql.IssueStateHistoryProjectionStateHistoryIssueStateSpanConnectionNodesIssueStateSpan,
) IssueStateSpanSummary {
	stateName := ""
	stateType := ""
	if span.State != nil {
		stateName = span.State.Name
		stateType = span.State.Type
	}

	return IssueStateSpanSummary{
		ID:        span.Id,
		StateID:   span.StateId,
		StateName: stateName,
		StateType: stateType,
		StartedAt: span.StartedAt,
		EndedAt:   stringValue(span.EndedAt),
	}
}
