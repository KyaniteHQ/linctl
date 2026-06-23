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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
		summaries = append(summaries, issueSummaryFromFields(issue.IssueSummaryFields))
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
			summaries = append(summaries, issueSummaryFromFields(relation.RelatedIssue.IssueSummaryFields))
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

// GetIssueDetail returns an issue by Linear id or identifier with mutable text fields.
func GetIssueDetail(ctx context.Context, graphqlClient graphql.Client, id string) (IssueDetail, error) {
	issueResult, err := issue(ctx, graphqlClient, id)
	if err != nil {
		return IssueDetail{}, fmt.Errorf("get issue %s: %w", id, err)
	}

	return detailIssue(issueResult.Issue), nil
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
