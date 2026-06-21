package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/Khan/genqlient/graphql"
)

// ErrWriteInvalid marks a malformed write request.
var ErrWriteInvalid = errors.New("invalid write")

// IssueSummary is the compact issue model used by read-only commands.
type IssueSummary struct {
	ID            string  `json:"id"`
	Identifier    string  `json:"identifier"`
	Title         string  `json:"title"`
	BranchName    string  `json:"branch_name"`
	URL           string  `json:"url"`
	Priority      float64 `json:"priority"`
	PriorityLabel string  `json:"priority_label"`
	TeamID        string  `json:"team_id"`
	Team          string  `json:"team"`
	StateID       string  `json:"state_id"`
	State         string  `json:"state"`
	StateType     string  `json:"state_type"`
	Assignee      string  `json:"assignee,omitempty"`
	ProjectID     string  `json:"project_id,omitempty"`
	Project       string  `json:"project,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UnblocksCount int     `json:"unblocks_count,omitempty"`
}

// IssueDetail is one issue with fields needed for safe read-modify-write operations.
type IssueDetail struct {
	Summary     IssueSummary
	Description string
}

// IssueList is a page of read-only issues.
type IssueList struct {
	Issues      []IssueSummary `json:"issues"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// IssuePriorityValue is a Linear issue priority number and label.
type IssuePriorityValue struct {
	Priority int    `json:"priority"`
	Label    string `json:"label"`
}

// IssueFilterSuggestion is an AI-generated issue filter suggestion.
type IssueFilterSuggestion struct {
	Filter json.RawMessage `json:"filter,omitempty"`
	LogID  string          `json:"log_id,omitempty"`
}

// IssueTitleSuggestion is an AI-generated issue title suggestion.
type IssueTitleSuggestion struct {
	Title string `json:"title"`
	LogID string `json:"log_id,omitempty"`
}

// IssueHistorySummary is compact issue history metadata without raw change payloads.
type IssueHistorySummary struct {
	ID                 string `json:"id"`
	IssueID            string `json:"issue_id"`
	ActorID            string `json:"actor_id,omitempty"`
	UpdatedDescription bool   `json:"updated_description,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	ArchivedAt         string `json:"archived_at,omitempty"`
}

// IssueHistoryList is a page of issue history metadata.
type IssueHistoryList struct {
	History     []IssueHistorySummary `json:"history"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// IssueBotActor is the optional bot actor attached to an issue.
type IssueBotActor struct {
	IssueID string           `json:"issue_id"`
	Bot     *ActorBotSummary `json:"bot,omitempty"`
}

// IssueStateSpanSummary is compact workflow-state span metadata for one issue.
type IssueStateSpanSummary struct {
	ID        string `json:"id"`
	StateID   string `json:"state_id"`
	StateName string `json:"state_name,omitempty"`
	StateType string `json:"state_type,omitempty"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at,omitempty"`
}

// IssueStateHistoryList is a page of workflow-state spans for one issue.
type IssueStateHistoryList struct {
	IssueID     string                  `json:"issue_id"`
	Spans       []IssueStateSpanSummary `json:"spans"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// IssueDependencyGraph is the compact dependency graph for one issue.
type IssueDependencyGraph struct {
	ID          string         `json:"id"`
	Identifier  string         `json:"identifier"`
	Parent      *IssueSummary  `json:"parent,omitempty"`
	Children    []IssueSummary `json:"children"`
	Blocks      []IssueSummary `json:"blocks"`
	BlockedBy   []IssueSummary `json:"blocked_by"`
	HasNextPage bool           `json:"has_next_page"`
}

// IssueListFilters scopes read-only issue listing.
type IssueListFilters struct {
	StateType     string
	ProjectID     string
	AssigneeID    string
	LabelID       string
	CycleID       string
	CreatedAfter  string
	CreatedBefore string
	HasBlockers   bool
	Blocks        bool
	BlockedBy     string
}

// ListIssues returns issues across every visible Linear team for broad read-only inspection.
func ListIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	issuePage, err := issues(ctx, graphqlClient, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issuePage.Issues.Nodes))
	for _, issue := range issuePage.Issues.Nodes {
		summaries = append(summaries, allTeamIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issuePage.Issues.PageInfo.HasNextPage,
		EndCursor:   issuePage.Issues.PageInfo.EndCursor,
	}, nil
}

// ListIssuesByTeam returns issues scoped to a resolved team.
func ListIssuesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	filters IssueListFilters,
) (IssueList, error) {
	if filters.StateType != "" {
		return listIssuesByTeamState(ctx, graphqlClient, teamID, limit, filters.StateType)
	}
	if filters.ProjectID != "" {
		return listIssuesByTeamProject(ctx, graphqlClient, teamID, limit, filters.ProjectID)
	}
	if filters.AssigneeID != "" {
		return listIssuesByTeamAssignee(ctx, graphqlClient, teamID, limit, filters.AssigneeID)
	}
	if filters.LabelID != "" {
		return listIssuesByTeamLabel(ctx, graphqlClient, teamID, limit, filters.LabelID)
	}
	if filters.CycleID != "" {
		return listIssuesByTeamCycle(ctx, graphqlClient, teamID, limit, filters.CycleID)
	}
	if filters.CreatedAfter != "" {
		return listIssuesByTeamCreatedAfter(ctx, graphqlClient, teamID, limit, filters.CreatedAfter)
	}
	if filters.CreatedBefore != "" {
		return listIssuesByTeamCreatedBefore(ctx, graphqlClient, teamID, limit, filters.CreatedBefore)
	}
	if filters.HasBlockers {
		return listIssuesByTeamHasBlockers(ctx, graphqlClient, teamID, limit)
	}
	if filters.Blocks {
		return listIssuesByTeamBlocks(ctx, graphqlClient, teamID, limit)
	}
	if filters.BlockedBy != "" {
		return listIssuesBlockedByIssue(ctx, graphqlClient, teamID, limit, filters.BlockedBy)
	}

	issues, err := IssuesByTeam(ctx, graphqlClient, teamID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}
	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, listIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamState(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	stateType string,
) (IssueList, error) {
	issues, err := IssuesByTeamState(ctx, graphqlClient, teamID, stateType, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, filteredIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	projectID string,
) (IssueList, error) {
	issues, err := IssuesByTeamProject(ctx, graphqlClient, teamID, projectID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, projectIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamAssignee(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	assigneeID string,
) (IssueList, error) {
	issues, err := IssuesByTeamAssignee(ctx, graphqlClient, teamID, assigneeID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, assigneeIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	labelID string,
) (IssueList, error) {
	issues, err := IssuesByTeamLabel(ctx, graphqlClient, teamID, labelID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, labelIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamCycle(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	cycleID string,
) (IssueList, error) {
	issues, err := IssuesByTeamCycle(ctx, graphqlClient, teamID, cycleID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, cycleIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamCreatedAfter(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	createdAfter string,
) (IssueList, error) {
	issues, err := IssuesByTeamCreatedAfter(ctx, graphqlClient, teamID, createdAfter, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, createdAfterIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamCreatedBefore(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	createdBefore string,
) (IssueList, error) {
	issues, err := IssuesByTeamCreatedBefore(ctx, graphqlClient, teamID, createdBefore, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, createdBeforeIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamHasBlockers(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (IssueList, error) {
	issues, err := IssuesByTeamHasBlockers(ctx, graphqlClient, teamID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, hasBlockersIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesByTeamBlocks(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (IssueList, error) {
	issues, err := IssuesByTeamBlocks(ctx, graphqlClient, teamID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, blocksIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

// ListNextIssuesByTeam returns unstarted issues that are not blocked by another issue.
func ListNextIssuesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (IssueList, error) {
	issues, err := NextIssuesByTeam(ctx, graphqlClient, teamID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list next issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.Issues.Nodes))
	for _, issue := range issues.Issues.Nodes {
		summaries = append(summaries, nextIssueSummary(issue))
	}
	sortNextIssueCandidates(summaries)

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.Issues.PageInfo.HasNextPage,
		EndCursor:   issues.Issues.PageInfo.EndCursor,
	}, nil
}

func listIssuesBlockedByIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	blockerIssueID string,
) (IssueList, error) {
	issue, err := IssueBlockedIssues(ctx, graphqlClient, blockerIssueID, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issue.Issue.Relations.Nodes))
	for _, relation := range issue.Issue.Relations.Nodes {
		if relation.Type == "blocks" && relation.RelatedIssue.Team.Id == teamID {
			summaries = append(summaries, blockedByIssueSummary(relation.RelatedIssue))
		}
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issue.Issue.Relations.PageInfo.HasNextPage,
		EndCursor:   issue.Issue.Relations.PageInfo.EndCursor,
	}, nil
}

// GetIssueDependencies returns parent, child, and blocking relationships for one issue.
func GetIssueDependencies(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueDependencyGraph, error) {
	dependencies, err := IssueDependencies(ctx, graphqlClient, id, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueDependencyGraph{}, fmt.Errorf("get issue dependencies %s: %w", id, err)
	}

	issue := dependencies.Issue
	return IssueDependencyGraph{
		ID:          issue.Id,
		Identifier:  issue.Identifier,
		Parent:      issueDependencyParent(issue.Parent),
		Children:    issueDependencyChildren(issue.Children.Nodes),
		Blocks:      issueDependencyBlocks(issue.Relations.Nodes),
		BlockedBy:   issueDependencyBlockedBy(issue.InverseRelations.Nodes),
		HasNextPage: issueDependencyHasNextPage(issue),
	}, nil
}

// SearchIssuesByTeam searches issue content scoped to a resolved team.
func SearchIssuesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	query string,
	limit int,
) (IssueList, error) {
	issues, err := issueSearch(ctx, graphqlClient, teamID, query, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("search issues: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.IssueSearch.Nodes))
	for _, issue := range issues.IssueSearch.Nodes {
		summaries = append(summaries, searchIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.IssueSearch.PageInfo.HasNextPage,
		EndCursor:   issues.IssueSearch.PageInfo.EndCursor,
	}, nil
}

// SearchIssuesByFigmaFileKey searches issues associated with a Figma file key.
func SearchIssuesByFigmaFileKey(
	ctx context.Context,
	graphqlClient graphql.Client,
	fileKey string,
	limit int,
) (IssueList, error) {
	issues, err := issueFigmaFileKeySearch(ctx, graphqlClient, fileKey, &limit, nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("search issues by Figma file key: %w", err)
	}

	summaries := make([]IssueSummary, 0, len(issues.IssueFigmaFileKeySearch.Nodes))
	for _, issue := range issues.IssueFigmaFileKeySearch.Nodes {
		summaries = append(summaries, figmaFileKeyIssueSummary(issue))
	}

	return IssueList{
		Issues:      summaries,
		HasNextPage: issues.IssueFigmaFileKeySearch.PageInfo.HasNextPage,
		EndCursor:   issues.IssueFigmaFileKeySearch.PageInfo.EndCursor,
	}, nil
}

// ListIssuePriorityValues returns Linear issue priority labels.
func ListIssuePriorityValues(ctx context.Context, graphqlClient graphql.Client) ([]IssuePriorityValue, error) {
	result, err := issuePriorityValues(ctx, graphqlClient)
	if err != nil {
		return nil, fmt.Errorf("list issue priority values: %w", err)
	}

	values := make([]IssuePriorityValue, 0, len(result.IssuePriorityValues))
	for _, value := range result.IssuePriorityValues {
		values = append(values, IssuePriorityValue(value))
	}

	return values, nil
}

// GetIssueFilterSuggestion returns a JSON issue filter suggestion for a prompt.
func GetIssueFilterSuggestion(
	ctx context.Context,
	graphqlClient graphql.Client,
	prompt string,
	teamID string,
	projectID string,
) (IssueFilterSuggestion, error) {
	suggestion, err := issueFilterSuggestion(
		ctx,
		graphqlClient,
		prompt,
		optionalString(teamID),
		optionalString(projectID),
	)
	if err != nil {
		return IssueFilterSuggestion{}, fmt.Errorf("get issue filter suggestion: %w", err)
	}

	filter := json.RawMessage(nil)
	if suggestion.IssueFilterSuggestion.Filter != nil {
		filter = *suggestion.IssueFilterSuggestion.Filter
	}

	return IssueFilterSuggestion{
		Filter: filter,
		LogID:  stringValue(suggestion.IssueFilterSuggestion.LogId),
	}, nil
}

// GetIssueTitleSuggestionFromCustomerRequest returns a title suggestion for customer request text.
func GetIssueTitleSuggestionFromCustomerRequest(
	ctx context.Context,
	graphqlClient graphql.Client,
	request string,
) (IssueTitleSuggestion, error) {
	suggestion, err := issueTitleSuggestionFromCustomerRequest(ctx, graphqlClient, request)
	if err != nil {
		return IssueTitleSuggestion{}, fmt.Errorf("get issue title suggestion: %w", err)
	}

	return IssueTitleSuggestion{
		Title: suggestion.IssueTitleSuggestionFromCustomerRequest.Title,
		LogID: stringValue(suggestion.IssueTitleSuggestionFromCustomerRequest.LogId),
	}, nil
}

// GetIssueByID returns a read-only issue by Linear id or identifier.
func GetIssueByID(ctx context.Context, graphqlClient graphql.Client, id string) (IssueSummary, error) {
	issue, err := GetIssueDetail(ctx, graphqlClient, id)
	if err != nil {
		return IssueSummary{}, err
	}

	return issue.Summary, nil
}

// ListIssueAttachments returns attachments associated with one issue.
func ListIssueAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (AttachmentList, error) {
	result, err := issue_attachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list issue attachments %s: %w", id, err)
	}

	attachments := make([]AttachmentSummary, 0, len(result.Issue.Attachments.Nodes))
	for _, attachment := range result.Issue.Attachments.Nodes {
		attachments = append(attachments, attachmentSummary(attachment.AttachmentSummaryFields))
	}

	return AttachmentList{
		Attachments: attachments,
		HasNextPage: result.Issue.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Attachments.PageInfo.EndCursor,
	}, nil
}

// GetIssueBotActor returns the bot actor that created an issue, when present.
func GetIssueBotActor(ctx context.Context, graphqlClient graphql.Client, id string) (IssueBotActor, error) {
	result, err := issue_botActor(ctx, graphqlClient, id)
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
	result, err := issue_children(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list issue children %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.Issue.Children.Nodes))
	for _, issue := range result.Issue.Children.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

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
	result, err := issue_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return DocumentList{}, fmt.Errorf("list issue documents %s: %w", id, err)
	}

	documents := make([]DocumentSummary, 0, len(result.Issue.Documents.Nodes))
	for _, document := range result.Issue.Documents.Nodes {
		documents = append(documents, documentSummary(document.DocumentSummaryFields))
	}

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
	result, err := issue_formerAttachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list issue former attachments %s: %w", id, err)
	}

	attachments := make([]AttachmentSummary, 0, len(result.Issue.FormerAttachments.Nodes))
	for _, attachment := range result.Issue.FormerAttachments.Nodes {
		attachments = append(attachments, attachmentSummary(attachment.AttachmentSummaryFields))
	}

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
	result, err := issue_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueHistoryList{}, fmt.Errorf("list issue history %s: %w", id, err)
	}

	history := make([]IssueHistorySummary, 0, len(result.Issue.History.Nodes))
	for _, node := range result.Issue.History.Nodes {
		history = append(history, issueHistorySummary(node))
	}

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
	result, err := issue_inverseRelations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue inverse relations %s: %w", id, err)
	}

	relations := make([]IssueRelationSummary, 0, len(result.Issue.InverseRelations.Nodes))
	for _, relation := range result.Issue.InverseRelations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

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
	result, err := issue_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list issue labels %s: %w", id, err)
	}

	labels := make([]LabelSummary, 0, len(result.Issue.Labels.Nodes))
	for _, label := range result.Issue.Labels.Nodes {
		labels = append(labels, labelSummary(label.IssueLabelSummaryFields))
	}

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
	result, err := issue_relations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueRelationList{}, fmt.Errorf("list issue relations %s: %w", id, err)
	}

	relations := make([]IssueRelationSummary, 0, len(result.Issue.Relations.Nodes))
	for _, relation := range result.Issue.Relations.Nodes {
		relations = append(relations, issueRelationSummary(relation.IssueRelationSummaryFields))
	}

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
	result, err := issue_releases(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ReleaseList{}, fmt.Errorf("list issue releases %s: %w", id, err)
	}

	releases := make([]ReleaseSummary, 0, len(result.Issue.Releases.Nodes))
	for _, release := range result.Issue.Releases.Nodes {
		releases = append(releases, releaseSummary(release.ReleaseSummaryFields))
	}

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
	result, err := issue_stateHistory(ctx, graphqlClient, id, intPtr(limit), nil)
	if err != nil {
		return IssueStateHistoryList{}, fmt.Errorf("list issue state history %s: %w", id, err)
	}

	spans := make([]IssueStateSpanSummary, 0, len(result.Issue.StateHistory.Nodes))
	for _, node := range result.Issue.StateHistory.Nodes {
		spans = append(spans, issueStateSpanSummary(node))
	}

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
	result, err := issue_subscribers(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list issue subscribers %s: %w", id, err)
	}

	users := make([]UserSummary, 0, len(result.Issue.Subscribers.Nodes))
	for _, node := range result.Issue.Subscribers.Nodes {
		users = append(users, userSummary(node.UserSummaryFields))
	}

	return UserList{
		Users:       users,
		HasNextPage: result.Issue.Subscribers.PageInfo.HasNextPage,
		EndCursor:   result.Issue.Subscribers.PageInfo.EndCursor,
	}, nil
}

// GetIssueByVCSBranch returns an issue by VCS branch name.
func GetIssueByVCSBranch(ctx context.Context, graphqlClient graphql.Client, branchName string) (IssueSummary, error) {
	result, err := issueVcsBranchSearch(ctx, graphqlClient, branchName)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("get issue by vcs branch %s: %w", branchName, err)
	}
	if result.IssueVcsBranchSearch == nil {
		return IssueSummary{}, fmt.Errorf("get issue by vcs branch %s: not found", branchName)
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
		return AttachmentList{}, fmt.Errorf("list issue vcs branch attachments %s: not found", branchName)
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
		return IssueBotActor{}, fmt.Errorf("get issue vcs branch bot actor %s: not found", branchName)
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
		return IssueList{}, fmt.Errorf("list issue vcs branch children %s: not found", branchName)
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
		return DocumentList{}, fmt.Errorf("list issue vcs branch documents %s: not found", branchName)
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
		return AttachmentList{}, fmt.Errorf("list issue vcs branch former attachments %s: not found", branchName)
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
		return IssueHistoryList{}, fmt.Errorf("list issue vcs branch history %s: not found", branchName)
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
		return IssueRelationList{}, fmt.Errorf("list issue vcs branch inverse relations %s: not found", branchName)
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
		return LabelList{}, fmt.Errorf("list issue vcs branch labels %s: not found", branchName)
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
		return IssueRelationList{}, fmt.Errorf("list issue vcs branch relations %s: not found", branchName)
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
		return ReleaseList{}, fmt.Errorf("list issue vcs branch releases %s: not found", branchName)
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
		return IssueStateHistoryList{}, fmt.Errorf("list issue vcs branch state history %s: not found", branchName)
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
		return UserList{}, fmt.Errorf("list issue vcs branch subscribers %s: not found", branchName)
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

// GetIssueDetail returns an issue by Linear id or identifier with mutable text fields.
func GetIssueDetail(ctx context.Context, graphqlClient graphql.Client, id string) (IssueDetail, error) {
	issueResult, err := issue(ctx, graphqlClient, id)
	if err != nil {
		return IssueDetail{}, fmt.Errorf("get issue %s: %w", id, err)
	}

	return detailIssue(issueResult.Issue), nil
}

func allTeamIssueSummary(issue issuesIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func listIssueSummary(issue IssuesByTeamIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func filteredIssueSummary(issue IssuesByTeamStateIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func projectIssueSummary(issue IssuesByTeamProjectIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func assigneeIssueSummary(issue IssuesByTeamAssigneeIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func labelIssueSummary(issue IssuesByTeamLabelIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func cycleIssueSummary(issue IssuesByTeamCycleIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func createdAfterIssueSummary(issue IssuesByTeamCreatedAfterIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func createdBeforeIssueSummary(issue IssuesByTeamCreatedBeforeIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func hasBlockersIssueSummary(issue IssuesByTeamHasBlockersIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func blocksIssueSummary(issue IssuesByTeamBlocksIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func nextIssueSummary(issue NextIssuesByTeamIssuesIssueConnectionNodesIssue) IssueSummary {
	summary := issueSummaryFromFields(issue.IssueSummaryFields)
	summary.CreatedAt = issue.CreatedAt
	summary.UnblocksCount = nextIssueUnblocksCount(issue.Relations.Nodes)

	return summary
}

//nolint:lll // The aliased name is generated by genqlient from the GraphQL selection path.
type nextIssueRelation = NextIssuesByTeamIssuesIssueConnectionNodesIssueRelationsIssueRelationConnectionNodesIssueRelation

func nextIssueUnblocksCount(relations []nextIssueRelation) int {
	count := 0
	for _, relation := range relations {
		if relation.Type == "blocks" && isActiveIssueStateType(relation.RelatedIssue.State.Type) {
			count++
		}
	}

	return count
}

func isActiveIssueStateType(stateType string) bool {
	return stateType != "completed" && stateType != "canceled"
}

func sortNextIssueCandidates(issues []IssueSummary) {
	sort.SliceStable(issues, func(leftIndex int, rightIndex int) bool {
		left := issues[leftIndex]
		right := issues[rightIndex]
		if left.UnblocksCount != right.UnblocksCount {
			return left.UnblocksCount > right.UnblocksCount
		}
		if left.Priority != right.Priority {
			return linearPriorityRank(left.Priority) > linearPriorityRank(right.Priority)
		}

		return left.CreatedAt < right.CreatedAt
	})
}

func linearPriorityRank(priority float64) float64 {
	if priority == 0 {
		return -1
	}

	return 5 - priority
}

func blockedByIssueSummary(
	issue IssueBlockedIssuesIssueRelationsIssueRelationConnectionNodesIssueRelationRelatedIssue,
) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func issueDependencyParent(issue *IssueDependenciesIssueParentIssue) *IssueSummary {
	if issue == nil {
		return nil
	}

	summary := issueSummaryFromFields(issue.IssueSummaryFields)
	return &summary
}

func issueDependencyChildren(issues []IssueDependenciesIssueChildrenIssueConnectionNodesIssue) []IssueSummary {
	summaries := make([]IssueSummary, 0, len(issues))
	for _, issue := range issues {
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return summaries
}

func issueDependencyBlocks(
	relations []IssueDependenciesIssueRelationsIssueRelationConnectionNodesIssueRelation,
) []IssueSummary {
	summaries := make([]IssueSummary, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == "blocks" {
			summaries = append(summaries, issueSummaryFromFields(relation.RelatedIssue.IssueSummaryFields))
		}
	}

	return summaries
}

func issueDependencyBlockedBy(
	relations []IssueDependenciesIssueInverseRelationsIssueRelationConnectionNodesIssueRelation,
) []IssueSummary {
	summaries := make([]IssueSummary, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == "blocks" {
			summaries = append(summaries, issueSummaryFromFields(relation.Issue.IssueSummaryFields))
		}
	}

	return summaries
}

func issueDependencyHasNextPage(issue IssueDependenciesIssue) bool {
	return issue.Children.PageInfo.HasNextPage ||
		issue.Relations.PageInfo.HasNextPage ||
		issue.InverseRelations.PageInfo.HasNextPage
}

func searchIssueSummary(issue issueSearchIssueSearchIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func figmaFileKeyIssueSummary(
	issue issueFigmaFileKeySearchIssueFigmaFileKeySearchIssueConnectionNodesIssue,
) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func detailIssueSummary(issue issueIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func detailIssue(issue issueIssue) IssueDetail {
	description := ""
	if issue.Description != nil {
		description = *issue.Description
	}

	return IssueDetail{
		Summary:     detailIssueSummary(issue),
		Description: description,
	}
}

func issueHistorySummary(history issue_historyIssueHistoryIssueHistoryConnectionNodesIssueHistory) IssueHistorySummary {
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

func issueActorBotSummary(bot *issue_botActorIssueBotActorActorBot) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return actorBotSummary(&bot.ActorBotSummaryFields)
}

func issueVCSBranchActorBotSummary(
	bot *issueVcsBranchSearch_botActorIssueVcsBranchSearchIssueBotActorActorBot,
) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return actorBotSummary(&bot.ActorBotSummaryFields)
}

func issueStateSpanSummary(
	span issue_stateHistoryIssueStateHistoryIssueStateSpanConnectionNodesIssueStateSpan,
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

func issueSummaryFromFields(issue IssueSummaryFields) IssueSummary {
	assignee := ""
	if issue.Assignee != nil {
		assignee = issue.Assignee.DisplayName
	}
	projectID := ""
	project := ""
	if issue.Project != nil {
		projectID = issue.Project.Id
		project = issue.Project.Name
	}

	return IssueSummary{
		ID:            issue.Id,
		Identifier:    issue.Identifier,
		Title:         issue.Title,
		BranchName:    issue.BranchName,
		URL:           issue.Url,
		Priority:      issue.Priority,
		PriorityLabel: issue.PriorityLabel,
		TeamID:        issue.Team.Id,
		Team:          issue.Team.Key,
		StateID:       issue.State.Id,
		State:         issue.State.Name,
		StateType:     issue.State.Type,
		Assignee:      assignee,
		ProjectID:     projectID,
		Project:       project,
	}
}
