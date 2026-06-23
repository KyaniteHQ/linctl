package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// GetIssueByVCSBranch returns an issue by VCS branch name.
func GetIssueByVCSBranch(ctx context.Context, graphqlClient graphql.Client, branchName string) (IssueSummary, error) {
	result, err := issueVcsBranchSearch(ctx, graphqlClient, branchName)
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
	result, err := issueVcsBranchSearch_attachments(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list issue vcs branch attachments %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return AttachmentList{}, notFoundError("list issue vcs branch attachments %s", branchName)
	}

	attachments := make([]AttachmentSummary, 0, len(result.IssueVcsBranchSearch.Attachments.Nodes))
	for _, attachment := range result.IssueVcsBranchSearch.Attachments.Nodes {
		attachments = append(attachments, attachmentSummary(attachment.AttachmentSummaryFields))
	}

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.IssueVcsBranchSearch.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Attachments.PageInfo.EndCursor,
	}, nil
}

// GetIssueVCSBranchBotActor returns bot actor metadata for the issue matched by a VCS branch.
func GetIssueVCSBranchBotActor(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
) (IssueBotActor, error) {
	result, err := issueVcsBranchSearch_botActor(ctx, graphqlClient, branchName)
	if err != nil {
		return IssueBotActor{}, fmt.Errorf("get issue vcs branch bot actor %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueBotActor{}, notFoundError("get issue vcs branch bot actor %s", branchName)
	}

	return IssueBotActor{
		IssueID: result.IssueVcsBranchSearch.Id,
		Bot:     issueVCSBranchActorBotSummary(result.IssueVcsBranchSearch.BotActor),
	}, nil
}

// ListIssueVCSBranchChildren returns child issues for the issue matched by a VCS branch.
func ListIssueVCSBranchChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueList, error) {
	result, err := issueVcsBranchSearch_children(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issue vcs branch children %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueList{}, notFoundError("list issue vcs branch children %s", branchName)
	}

	issues := make([]IssueSummary, 0, len(result.IssueVcsBranchSearch.Children.Nodes))
	for _, issue := range result.IssueVcsBranchSearch.Children.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.IssueVcsBranchSearch.Children.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Children.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchDocuments returns documents for the issue matched by a VCS branch.
func ListIssueVCSBranchDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (DocumentList, error) {
	result, err := issueVcsBranchSearch_documents(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list issue vcs branch documents %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return DocumentList{}, notFoundError("list issue vcs branch documents %s", branchName)
	}

	documents := make([]DocumentSummary, 0, len(result.IssueVcsBranchSearch.Documents.Nodes))
	for _, document := range result.IssueVcsBranchSearch.Documents.Nodes {
		documents = append(documents, documentSummary(document.DocumentSummaryFields))
	}

	return DocumentList{
		Documents:   documents,
		HasNextPage: result.IssueVcsBranchSearch.Documents.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Documents.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchFormerAttachments returns former attachments for the issue matched by a VCS branch.
func ListIssueVCSBranchFormerAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (AttachmentList, error) {
	result, err := issueVcsBranchSearch_formerAttachments(
		ctx,
		graphqlClient,
		branchName,
		intPtr(limit),
		nil,
		boolPtr(true),
	)
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list issue vcs branch former attachments %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return AttachmentList{}, notFoundError("list issue vcs branch former attachments %s", branchName)
	}

	attachments := make([]AttachmentSummary, 0, len(result.IssueVcsBranchSearch.FormerAttachments.Nodes))
	for _, attachment := range result.IssueVcsBranchSearch.FormerAttachments.Nodes {
		attachments = append(attachments, attachmentSummary(attachment.AttachmentSummaryFields))
	}

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.IssueVcsBranchSearch.FormerAttachments.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.FormerAttachments.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchHistory returns history metadata for the issue matched by a VCS branch.
func ListIssueVCSBranchHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueHistoryList, error) {
	result, err := issueVcsBranchSearch_history(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueHistoryList{}, fmt.Errorf("list issue vcs branch history %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueHistoryList{}, notFoundError("list issue vcs branch history %s", branchName)
	}

	history := make([]IssueHistorySummary, 0, len(result.IssueVcsBranchSearch.History.Nodes))
	for _, node := range result.IssueVcsBranchSearch.History.Nodes {
		history = append(history, issueVCSBranchHistorySummary(node))
	}

	return IssueHistoryList{
		History:     history,
		HasNextPage: result.IssueVcsBranchSearch.History.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.History.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchInverseRelations returns inverse relations for the issue matched by a VCS branch.
func ListIssueVCSBranchInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueRelationList, error) {
	result, err := issueVcsBranchSearch_inverseRelations(
		ctx,
		graphqlClient,
		branchName,
		intPtr(limit),
		nil,
		boolPtr(true),
	)
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue vcs branch inverse relations %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueRelationList{}, notFoundError("list issue vcs branch inverse relations %s", branchName)
	}

	relations := make([]IssueRelationSummary, 0, len(result.IssueVcsBranchSearch.InverseRelations.Nodes))
	for _, relation := range result.IssueVcsBranchSearch.InverseRelations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.IssueVcsBranchSearch.InverseRelations.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.InverseRelations.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchLabels returns labels for the issue matched by a VCS branch.
func ListIssueVCSBranchLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (LabelList, error) {
	result, err := issueVcsBranchSearch_labels(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list issue vcs branch labels %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return LabelList{}, notFoundError("list issue vcs branch labels %s", branchName)
	}

	labels := make([]LabelSummary, 0, len(result.IssueVcsBranchSearch.Labels.Nodes))
	for _, label := range result.IssueVcsBranchSearch.Labels.Nodes {
		labels = append(labels, labelSummary(label.IssueLabelSummaryFields))
	}

	return LabelList{
		Labels:      labels,
		HasNextPage: result.IssueVcsBranchSearch.Labels.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Labels.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchRelations returns relations for the issue matched by a VCS branch.
func ListIssueVCSBranchRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueRelationList, error) {
	result, err := issueVcsBranchSearch_relations(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue vcs branch relations %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueRelationList{}, notFoundError("list issue vcs branch relations %s", branchName)
	}

	relations := make([]IssueRelationSummary, 0, len(result.IssueVcsBranchSearch.Relations.Nodes))
	for _, relation := range result.IssueVcsBranchSearch.Relations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

	return IssueRelationList{
		Relations:   relations,
		HasNextPage: result.IssueVcsBranchSearch.Relations.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Relations.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchReleases returns releases for the issue matched by a VCS branch.
func ListIssueVCSBranchReleases(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (ReleaseList, error) {
	result, err := issueVcsBranchSearch_releases(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list issue vcs branch releases %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return ReleaseList{}, notFoundError("list issue vcs branch releases %s", branchName)
	}

	releases := make([]ReleaseSummary, 0, len(result.IssueVcsBranchSearch.Releases.Nodes))
	for _, release := range result.IssueVcsBranchSearch.Releases.Nodes {
		releases = append(releases, releaseSummary(release.ReleaseSummaryFields))
	}

	return ReleaseList{
		Releases:    releases,
		HasNextPage: result.IssueVcsBranchSearch.Releases.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Releases.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchStateHistory returns workflow-state spans for the issue matched by a VCS branch.
func ListIssueVCSBranchStateHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (IssueStateHistoryList, error) {
	result, err := issueVcsBranchSearch_stateHistory(ctx, graphqlClient, branchName, intPtr(limit), nil)
	if err != nil {
		return IssueStateHistoryList{}, fmt.Errorf("list issue vcs branch state history %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueStateHistoryList{}, notFoundError("list issue vcs branch state history %s", branchName)
	}

	spans := make([]IssueStateSpanSummary, 0, len(result.IssueVcsBranchSearch.StateHistory.Nodes))
	for _, node := range result.IssueVcsBranchSearch.StateHistory.Nodes {
		spans = append(spans, issueVCSBranchStateSpanSummary(node))
	}

	return IssueStateHistoryList{
		IssueID:     result.IssueVcsBranchSearch.Id,
		Spans:       spans,
		HasNextPage: result.IssueVcsBranchSearch.StateHistory.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.StateHistory.PageInfo.EndCursor,
	}, nil
}

// ListIssueVCSBranchSubscribers returns subscribers for the issue matched by a VCS branch.
func ListIssueVCSBranchSubscribers(
	ctx context.Context,
	graphqlClient graphql.Client,
	branchName string,
	limit int,
) (UserList, error) {
	result, err := issueVcsBranchSearch_subscribers(ctx, graphqlClient, branchName, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list issue vcs branch subscribers %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return UserList{}, notFoundError("list issue vcs branch subscribers %s", branchName)
	}

	users := make([]UserSummary, 0, len(result.IssueVcsBranchSearch.Subscribers.Nodes))
	for _, node := range result.IssueVcsBranchSearch.Subscribers.Nodes {
		users = append(users, userSummary(node.UserSummaryFields))
	}

	return UserList{
		Users:       users,
		HasNextPage: result.IssueVcsBranchSearch.Subscribers.PageInfo.HasNextPage,
		EndCursor:   result.IssueVcsBranchSearch.Subscribers.PageInfo.EndCursor,
	}, nil
}

func issueVCSBranchHistorySummary(
	history issueVcsBranchSearch_historyIssueVcsBranchSearchIssueHistoryIssueHistoryConnectionNodesIssueHistory,
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

func issueVCSBranchActorBotSummary(
	bot *issueVcsBranchSearch_botActorIssueVcsBranchSearchIssueBotActorActorBot,
) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return actorBotSummary(&bot.ActorBotSummaryFields)
}

//nolint:lll // The aliased name is generated by genqlient from the GraphQL selection path.
type issueVCSBranchStateSpan = issueVcsBranchSearch_stateHistoryIssueVcsBranchSearchIssueStateHistoryIssueStateSpanConnectionNodesIssueStateSpan

func issueVCSBranchStateSpanSummary(span issueVCSBranchStateSpan) IssueStateSpanSummary {
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
