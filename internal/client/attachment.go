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

// ListAttachments returns visible issue attachments.
func ListAttachments(ctx context.Context, graphqlClient graphql.Client, limit int) (AttachmentList, error) {
	result, err := gql.XAttachments(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachments: %w", err)
	}

	summaries := mapNodes(result.Attachments.Nodes, func(
		node gql.XAttachmentsAttachmentsAttachmentConnectionNodesAttachment,
	) AttachmentSummary {
		return attachmentSummary(node.AttachmentSummaryFields)
	})

	return AttachmentList{
		Attachments: summaries,
		HasNextPage: result.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.Attachments.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentsForURL returns visible issue attachments linked to a URL.
func ListAttachmentsForURL(
	ctx context.Context,
	graphqlClient graphql.Client,
	url string,
	limit int,
) (AttachmentList, error) {
	result, err := gql.XAttachmentsForURL(ctx, graphqlClient, url, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachments for url %s: %w", url, err)
	}

	summaries := mapNodes(result.AttachmentsForURL.Nodes, func(
		node gql.XAttachmentsForURLAttachmentsForURLAttachmentConnectionNodesAttachment,
	) AttachmentSummary {
		return attachmentSummary(node.AttachmentSummaryFields)
	})

	return AttachmentList{
		Attachments: summaries,
		HasNextPage: result.AttachmentsForURL.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentsForURL.PageInfo.EndCursor,
	}, nil
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
	result, err := gql.XAttachmentIssue_attachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachment issue attachments %s: %w", id, err)
	}

	attachments := mapNodes(result.AttachmentIssue.Attachments.Nodes, func(
		attachment gql.IssueAttachmentsProjectionAttachmentsAttachmentConnectionNodesAttachment,
	) AttachmentSummary {
		return attachmentSummary(attachment.AttachmentSummaryFields)
	})

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.AttachmentIssue.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Attachments.PageInfo.EndCursor,
	}, nil
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
	result, err := gql.XAttachmentIssue_children(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list attachment issue children %s: %w", id, err)
	}

	issues := mapNodes(result.AttachmentIssue.Children.Nodes, func(
		issue gql.IssueChildrenProjectionChildrenIssueConnectionNodesIssue,
	) IssueSummary {
		return issueSummaryFromFields(issue.IssueSummaryFields)
	})

	return IssueList{
		Issues:      issues,
		HasNextPage: result.AttachmentIssue.Children.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Children.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueDocuments returns documents for the issue associated with one attachment.
func ListAttachmentIssueDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (DocumentList, error) {
	result, err := gql.XAttachmentIssue_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list attachment issue documents %s: %w", id, err)
	}

	documents := mapNodes(result.AttachmentIssue.Documents.Nodes, func(
		document gql.IssueDocumentsProjectionDocumentsDocumentConnectionNodesDocument,
	) DocumentSummary {
		return documentSummary(document.DocumentSummaryFields)
	})

	return DocumentList{
		Documents:   documents,
		HasNextPage: result.AttachmentIssue.Documents.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Documents.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueFormerAttachments returns former attachments for the issue associated with one attachment.
func ListAttachmentIssueFormerAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	result, err := gql.XAttachmentIssue_formerAttachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachment issue former attachments %s: %w", id, err)
	}

	attachments := mapNodes(result.AttachmentIssue.FormerAttachments.Nodes, func(
		attachment gql.IssueFormerAttachmentsProjectionFormerAttachmentsAttachmentConnectionNodesAttachment,
	) AttachmentSummary {
		return attachmentSummary(attachment.AttachmentSummaryFields)
	})

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.AttachmentIssue.FormerAttachments.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.FormerAttachments.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueComments returns body-free comments for the issue associated with one attachment.
func ListAttachmentIssueComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCommentMetadataList, error) {
	result, err := gql.XAttachmentIssue_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCommentMetadataList{}, fmt.Errorf("list attachment issue comments %s: %w", id, err)
	}

	comments := mapNodes(result.AttachmentIssue.Comments.Nodes, func(
		comment gql.IssueCommentMetadataProjectionCommentsCommentConnectionNodesComment,
	) CommentMetadataSummary {
		return commentMetadataSummary(comment.CommentMetadataFields)
	})

	return IssueCommentMetadataList{
		IssueID:     result.AttachmentIssue.Id,
		Identifier:  result.AttachmentIssue.Identifier,
		Comments:    comments,
		HasNextPage: result.AttachmentIssue.Comments.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Comments.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueNeeds returns body-free customer needs for the issue associated with one attachment.
func ListAttachmentIssueNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	result, err := gql.XAttachmentIssue_needs(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf("list attachment issue customer needs %s: %w", id, err)
	}

	needs := mapNodes(result.AttachmentIssue.Needs.Nodes, func(
		need gql.IssueNeedsProjectionNeedsCustomerNeedConnectionNodesCustomerNeed,
	) CustomerNeedMetadataSummary {
		return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
	})

	return IssueCustomerNeedMetadataList{
		IssueID:     result.AttachmentIssue.Id,
		Identifier:  result.AttachmentIssue.Identifier,
		Needs:       needs,
		HasNextPage: result.AttachmentIssue.Needs.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Needs.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueFormerNeeds returns body-free former customer needs for the issue associated with one attachment.
func ListAttachmentIssueFormerNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCustomerNeedMetadataList, error) {
	result, err := gql.XAttachmentIssue_formerNeeds(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCustomerNeedMetadataList{}, fmt.Errorf(
			"list attachment issue former customer needs %s: %w",
			id,
			err,
		)
	}

	needs := mapNodes(result.AttachmentIssue.FormerNeeds.Nodes, func(
		need gql.IssueFormerNeedsProjectionFormerNeedsCustomerNeedConnectionNodesCustomerNeed,
	) CustomerNeedMetadataSummary {
		return customerNeedMetadataSummary(need.CustomerNeedMetadataFields)
	})

	return IssueCustomerNeedMetadataList{
		IssueID:     result.AttachmentIssue.Id,
		Identifier:  result.AttachmentIssue.Identifier,
		Needs:       needs,
		HasNextPage: result.AttachmentIssue.FormerNeeds.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.FormerNeeds.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueHistory returns history metadata for the issue associated with one attachment.
func ListAttachmentIssueHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueHistoryList, error) {
	result, err := gql.XAttachmentIssue_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueHistoryList{}, fmt.Errorf("list attachment issue history %s: %w", id, err)
	}

	history := mapNodes(result.AttachmentIssue.History.Nodes, issueHistorySummary)

	return IssueHistoryList{
		History:     history,
		HasNextPage: result.AttachmentIssue.History.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.History.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueInverseRelations returns inverse relations for the issue associated with one attachment.
func ListAttachmentIssueInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	result, err := gql.XAttachmentIssue_inverseRelations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list attachment issue inverse relations %s: %w", id, err)
	}

	relations := mapNodes(result.AttachmentIssue.InverseRelations.Nodes, func(
		node gql.IssueInverseRelationsProjectionInverseRelationsIssueRelationConnectionNodesIssueRelation,
	) IssueRelationSummary {
		return issueRelationSummary(node.IssueRelationSummaryFields)
	})

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.AttachmentIssue.InverseRelations.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.InverseRelations.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueLabels returns labels for the issue associated with one attachment.
func ListAttachmentIssueLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (LabelList, error) {
	result, err := gql.XAttachmentIssue_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list attachment issue labels %s: %w", id, err)
	}

	labels := mapNodes(result.AttachmentIssue.Labels.Nodes, func(
		label gql.IssueLabelsProjectionLabelsIssueLabelConnectionNodesIssueLabel,
	) LabelSummary {
		return labelSummary(label.IssueLabelSummaryFields)
	})

	return LabelList{
		Labels:      labels,
		HasNextPage: result.AttachmentIssue.Labels.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Labels.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueRelations returns relations for the issue associated with one attachment.
func ListAttachmentIssueRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueRelationList, error) {
	result, err := gql.XAttachmentIssue_relations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list attachment issue relations %s: %w", id, err)
	}

	relations := mapNodes(result.AttachmentIssue.Relations.Nodes, func(
		relation gql.IssueRelationsProjectionRelationsIssueRelationConnectionNodesIssueRelation,
	) IssueRelationSummary {
		return issueRelationSummary(relation.IssueRelationSummaryFields)
	})

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.AttachmentIssue.Relations.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Relations.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueReleases returns releases for the issue associated with one attachment.
func ListAttachmentIssueReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ReleaseList, error) {
	result, err := gql.XAttachmentIssue_releases(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list attachment issue releases %s: %w", id, err)
	}

	releases := mapNodes(result.AttachmentIssue.Releases.Nodes, func(
		release gql.IssueReleasesProjectionReleasesReleaseConnectionNodesRelease,
	) ReleaseSummary {
		return releaseSummary(release.ReleaseSummaryFields)
	})

	return ReleaseList{
		Releases:    releases,
		HasNextPage: result.AttachmentIssue.Releases.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Releases.PageInfo.EndCursor,
	}, nil
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
	result, err := gql.XAttachmentIssue_stateHistory(ctx, graphqlClient, id, intPtr(limit), nil)
	if err != nil {
		return IssueStateHistoryList{}, fmt.Errorf("list attachment issue state history %s: %w", id, err)
	}

	spans := mapNodes(result.AttachmentIssue.StateHistory.Nodes, issueStateSpanSummary)

	return IssueStateHistoryList{
		IssueID:     result.AttachmentIssue.Id,
		Spans:       spans,
		HasNextPage: result.AttachmentIssue.StateHistory.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.StateHistory.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueSubscribers returns subscribers for the issue associated with one attachment.
func ListAttachmentIssueSubscribers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (UserList, error) {
	result, err := gql.XAttachmentIssue_subscribers(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list attachment issue subscribers %s: %w", id, err)
	}

	users := mapNodes(result.AttachmentIssue.Subscribers.Nodes, func(
		node gql.IssueSubscribersProjectionSubscribersUserConnectionNodesUser,
	) UserSummary {
		return userSummary(node.UserSummaryFields)
	})

	return UserList{
		Users:       users,
		HasNextPage: result.AttachmentIssue.Subscribers.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Subscribers.PageInfo.EndCursor,
	}, nil
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
