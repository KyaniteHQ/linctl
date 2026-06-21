package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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
	result, err := attachments(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachments: %w", err)
	}

	summaries := make([]AttachmentSummary, 0, len(result.Attachments.Nodes))
	for _, node := range result.Attachments.Nodes {
		summaries = append(summaries, attachmentSummary(node.AttachmentSummaryFields))
	}

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
	result, err := attachmentsForURL(ctx, graphqlClient, url, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachments for url %s: %w", url, err)
	}

	summaries := make([]AttachmentSummary, 0, len(result.AttachmentsForURL.Nodes))
	for _, node := range result.AttachmentsForURL.Nodes {
		summaries = append(summaries, attachmentSummary(node.AttachmentSummaryFields))
	}

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
	result, err := attachment(ctx, graphqlClient, id)
	if err != nil {
		return AttachmentSummary{}, fmt.Errorf("get attachment %s: %w", id, err)
	}

	return attachmentSummary(result.Attachment.AttachmentSummaryFields), nil
}

// GetAttachmentIssue returns the issue associated with one attachment.
func GetAttachmentIssue(ctx context.Context, graphqlClient graphql.Client, id string) (IssueSummary, error) {
	result, err := attachmentIssue(ctx, graphqlClient, id)
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
	result, err := attachmentIssue_attachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachment issue attachments %s: %w", id, err)
	}

	attachments := make([]AttachmentSummary, 0, len(result.AttachmentIssue.Attachments.Nodes))
	for _, attachment := range result.AttachmentIssue.Attachments.Nodes {
		attachments = append(attachments, attachmentSummary(attachment.AttachmentSummaryFields))
	}

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.AttachmentIssue.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Attachments.PageInfo.EndCursor,
	}, nil
}

// GetAttachmentIssueBotActor returns the issue bot actor associated with one attachment.
func GetAttachmentIssueBotActor(ctx context.Context, graphqlClient graphql.Client, id string) (IssueBotActor, error) {
	result, err := attachmentIssue_botActor(ctx, graphqlClient, id)
	if err != nil {
		return IssueBotActor{}, fmt.Errorf("get attachment issue bot actor %s: %w", id, err)
	}

	return IssueBotActor{
		IssueID: result.AttachmentIssue.Id,
		Bot:     attachmentIssueActorBotSummary(result.AttachmentIssue.BotActor),
	}, nil
}

// ListAttachmentIssueChildren returns child issues for the issue associated with one attachment.
func ListAttachmentIssueChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := attachmentIssue_children(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list attachment issue children %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.AttachmentIssue.Children.Nodes))
	for _, issue := range result.AttachmentIssue.Children.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

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
	result, err := attachmentIssue_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list attachment issue documents %s: %w", id, err)
	}

	documents := make([]DocumentSummary, 0, len(result.AttachmentIssue.Documents.Nodes))
	for _, document := range result.AttachmentIssue.Documents.Nodes {
		documents = append(documents, documentSummary(document.DocumentSummaryFields))
	}

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
	result, err := attachmentIssue_formerAttachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachment issue former attachments %s: %w", id, err)
	}

	attachments := make([]AttachmentSummary, 0, len(result.AttachmentIssue.FormerAttachments.Nodes))
	for _, attachment := range result.AttachmentIssue.FormerAttachments.Nodes {
		attachments = append(attachments, attachmentSummary(attachment.AttachmentSummaryFields))
	}

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.AttachmentIssue.FormerAttachments.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.FormerAttachments.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueHistory returns history metadata for the issue associated with one attachment.
func ListAttachmentIssueHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueHistoryList, error) {
	result, err := attachmentIssue_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueHistoryList{}, fmt.Errorf("list attachment issue history %s: %w", id, err)
	}

	history := make([]IssueHistorySummary, 0, len(result.AttachmentIssue.History.Nodes))
	for _, node := range result.AttachmentIssue.History.Nodes {
		history = append(history, attachmentIssueHistorySummary(node))
	}

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
	result, err := attachmentIssue_inverseRelations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list attachment issue inverse relations %s: %w", id, err)
	}

	relations := make([]IssueRelationSummary, 0, len(result.AttachmentIssue.InverseRelations.Nodes))
	for _, relation := range result.AttachmentIssue.InverseRelations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

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
	result, err := attachmentIssue_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list attachment issue labels %s: %w", id, err)
	}

	labels := make([]LabelSummary, 0, len(result.AttachmentIssue.Labels.Nodes))
	for _, label := range result.AttachmentIssue.Labels.Nodes {
		labels = append(labels, labelSummary(label.IssueLabelSummaryFields))
	}

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
	result, err := attachmentIssue_relations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list attachment issue relations %s: %w", id, err)
	}

	relations := make([]IssueRelationSummary, 0, len(result.AttachmentIssue.Relations.Nodes))
	for _, relation := range result.AttachmentIssue.Relations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

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
	result, err := attachmentIssue_releases(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list attachment issue releases %s: %w", id, err)
	}

	releases := make([]ReleaseSummary, 0, len(result.AttachmentIssue.Releases.Nodes))
	for _, release := range result.AttachmentIssue.Releases.Nodes {
		releases = append(releases, releaseSummary(release.ReleaseSummaryFields))
	}

	return ReleaseList{
		Releases:    releases,
		HasNextPage: result.AttachmentIssue.Releases.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Releases.PageInfo.EndCursor,
	}, nil
}

// ListAttachmentIssueStateHistory returns workflow-state spans for the issue associated with one attachment.
func ListAttachmentIssueStateHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueStateHistoryList, error) {
	result, err := attachmentIssue_stateHistory(ctx, graphqlClient, id, intPtr(limit), nil)
	if err != nil {
		return IssueStateHistoryList{}, fmt.Errorf("list attachment issue state history %s: %w", id, err)
	}

	spans := make([]IssueStateSpanSummary, 0, len(result.AttachmentIssue.StateHistory.Nodes))
	for _, node := range result.AttachmentIssue.StateHistory.Nodes {
		spans = append(spans, attachmentIssueStateSpanSummary(node))
	}

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
	result, err := attachmentIssue_subscribers(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list attachment issue subscribers %s: %w", id, err)
	}

	users := make([]UserSummary, 0, len(result.AttachmentIssue.Subscribers.Nodes))
	for _, node := range result.AttachmentIssue.Subscribers.Nodes {
		users = append(users, userSummary(node.UserSummaryFields))
	}

	return UserList{
		Users:       users,
		HasNextPage: result.AttachmentIssue.Subscribers.PageInfo.HasNextPage,
		EndCursor:   result.AttachmentIssue.Subscribers.PageInfo.EndCursor,
	}, nil
}

func attachmentSummary(fields AttachmentSummaryFields) AttachmentSummary {
	return AttachmentSummary{
		ID:         fields.Id,
		Title:      fields.Title,
		Subtitle:   stringValue(fields.Subtitle),
		URL:        fields.Url,
		SourceType: stringValue(fields.SourceType),
	}
}

func attachmentIssueActorBotSummary(bot *attachmentIssue_botActorAttachmentIssueBotActorActorBot) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return actorBotSummary(&bot.ActorBotSummaryFields)
}

func attachmentIssueHistorySummary(
	history attachmentIssue_historyAttachmentIssueHistoryIssueHistoryConnectionNodesIssueHistory,
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

func attachmentIssueStateSpanSummary(
	span attachmentIssue_stateHistoryAttachmentIssueStateHistoryIssueStateSpanConnectionNodesIssueStateSpan,
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
