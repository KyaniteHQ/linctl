package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ListIssueAttachments returns attachments associated with one issue.
func ListIssueAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	result, err := gql.XIssue_attachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list issue attachments %s: %w", id, err)
	}

	attachments := mapNodes(result.Issue.Attachments.Nodes, func(
		attachment gql.IssueAttachmentsProjectionAttachmentsAttachmentConnectionNodesAttachment,
	) AttachmentSummary {
		return attachmentSummary(attachment.AttachmentSummaryFields)
	})

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.Issue.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Attachments.PageInfo.EndCursor,
	}, nil
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
	result, err := gql.XIssue_children(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issue children %s: %w", id, err)
	}

	issues := mapNodes(result.Issue.Children.Nodes, func(
		issue gql.IssueChildrenProjectionChildrenIssueConnectionNodesIssue,
	) IssueSummary {
		return issueSummaryFromFields(issue.IssueSummaryFields)
	})

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Issue.Children.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Children.PageInfo.EndCursor,
	}, nil
}

// ListIssueDocuments returns documents associated with one issue.
func ListIssueDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	result, err := gql.XIssue_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list issue documents %s: %w", id, err)
	}

	documents := mapNodes(result.Issue.Documents.Nodes, func(
		document gql.IssueDocumentsProjectionDocumentsDocumentConnectionNodesDocument,
	) DocumentSummary {
		return documentSummary(document.DocumentSummaryFields)
	})

	return DocumentList{
		Documents:   documents,
		HasNextPage: result.Issue.Documents.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Documents.PageInfo.EndCursor,
	}, nil
}

// ListIssueFormerAttachments returns attachments formerly associated with one issue.
func ListIssueFormerAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	result, err := gql.XIssue_formerAttachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list issue former attachments %s: %w", id, err)
	}

	attachments := mapNodes(result.Issue.FormerAttachments.Nodes, func(
		attachment gql.IssueFormerAttachmentsProjectionFormerAttachmentsAttachmentConnectionNodesAttachment,
	) AttachmentSummary {
		return attachmentSummary(attachment.AttachmentSummaryFields)
	})

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.Issue.FormerAttachments.PageInfo.HasNextPage,
		EndCursor:   result.Issue.FormerAttachments.PageInfo.EndCursor,
	}, nil
}

// ListIssueHistory returns compact history metadata for one issue.
func ListIssueHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueHistoryList, error) {
	result, err := gql.XIssue_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueHistoryList{}, fmt.Errorf("list issue history %s: %w", id, err)
	}

	history := mapNodes(result.Issue.History.Nodes, issueHistorySummary)

	return IssueHistoryList{
		History:     history,
		HasNextPage: result.Issue.History.PageInfo.HasNextPage,
		EndCursor:   result.Issue.History.PageInfo.EndCursor,
	}, nil
}

// ListIssueInverseRelations returns inverse relations associated with one issue.
func ListIssueInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	result, err := gql.XIssue_inverseRelations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue inverse relations %s: %w", id, err)
	}

	relations := mapNodes(result.Issue.InverseRelations.Nodes, func(
		relation gql.IssueInverseRelationsProjectionInverseRelationsIssueRelationConnectionNodesIssueRelation,
	) IssueRelationSummary {
		return issueRelationSummary(relation.IssueRelationSummaryFields)
	})

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.Issue.InverseRelations.PageInfo.HasNextPage,
		EndCursor:   result.Issue.InverseRelations.PageInfo.EndCursor,
	}, nil
}

// ListIssueLabels returns labels associated with one issue.
func ListIssueLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (LabelList, error) {
	result, err := gql.XIssue_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list issue labels %s: %w", id, err)
	}

	labels := mapNodes(result.Issue.Labels.Nodes, func(
		label gql.IssueLabelsProjectionLabelsIssueLabelConnectionNodesIssueLabel,
	) LabelSummary {
		return labelSummary(label.IssueLabelSummaryFields)
	})

	return LabelList{
		Labels:      labels,
		HasNextPage: result.Issue.Labels.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Labels.PageInfo.EndCursor,
	}, nil
}

// ListIssueRelationsForIssue returns relations associated with one issue.
func ListIssueRelationsForIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	result, err := gql.XIssue_relations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue relations %s: %w", id, err)
	}

	relations := mapNodes(result.Issue.Relations.Nodes, func(
		relation gql.IssueRelationsProjectionRelationsIssueRelationConnectionNodesIssueRelation,
	) IssueRelationSummary {
		return issueRelationSummary(relation.IssueRelationSummaryFields)
	})

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.Issue.Relations.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Relations.PageInfo.EndCursor,
	}, nil
}

// ListIssueReleases returns releases associated with one issue.
func ListIssueReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseList, error) {
	result, err := gql.XIssue_releases(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list issue releases %s: %w", id, err)
	}

	releases := mapNodes(result.Issue.Releases.Nodes, func(
		release gql.IssueReleasesProjectionReleasesReleaseConnectionNodesRelease,
	) ReleaseSummary {
		return releaseSummary(release.ReleaseSummaryFields)
	})

	return ReleaseList{
		Releases:    releases,
		HasNextPage: result.Issue.Releases.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Releases.PageInfo.EndCursor,
	}, nil
}

// ListIssueStateHistory returns workflow-state spans for one issue.
func ListIssueStateHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueStateHistoryList, error) {
	result, err := gql.XIssue_stateHistory(ctx, graphqlClient, id, intPtr(limit), nil)
	if err != nil {
		return IssueStateHistoryList{}, fmt.Errorf("list issue state history %s: %w", id, err)
	}

	spans := mapNodes(result.Issue.StateHistory.Nodes, issueStateSpanSummary)

	return IssueStateHistoryList{
		IssueID:     result.Issue.Id,
		Spans:       spans,
		HasNextPage: result.Issue.StateHistory.PageInfo.HasNextPage,
		EndCursor:   result.Issue.StateHistory.PageInfo.EndCursor,
	}, nil
}

// ListIssueSubscribers returns users subscribed to one issue.
func ListIssueSubscribers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (UserList, error) {
	result, err := gql.XIssue_subscribers(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list issue subscribers %s: %w", id, err)
	}

	users := mapNodes(result.Issue.Subscribers.Nodes, func(
		node gql.IssueSubscribersProjectionSubscribersUserConnectionNodesUser,
	) UserSummary {
		return userSummary(node.UserSummaryFields)
	})

	return UserList{
		Users:       users,
		HasNextPage: result.Issue.Subscribers.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Subscribers.PageInfo.EndCursor,
	}, nil
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
