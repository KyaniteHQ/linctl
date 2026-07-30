package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/client/internal/gqlmodel"
)

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
	UpdatedAfter  string
	UpdatedBefore string
	HasBlockers   bool
	Blocks        bool
	BlockedBy     string
}

// issueListPageSize matches Linear's per-request cap so accumulated issue
// listing keeps single-request behavior for every limit up to the cap.
const issueListPageSize = 250

// ListIssues returns issues across every visible Linear team for broad read-only inspection.
func ListIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	page, err := collectNodePages(
		"list issues", limit, issueListPageSize,
		func(pageSize int, after *string) (nodePage[IssueSummary], error) {
			issuePage, err := gql.XIssues(ctx, graphqlClient, &pageSize, after, boolPtr(true))
			if err != nil {
				return nodePage[IssueSummary]{}, err
			}

			return issueSummaryNodePage(
				mapNodes(issuePage.Issues.Nodes, func(issue gql.XIssuesIssuesIssueConnectionNodesIssue) IssueSummary {
					return issueSummaryFromFields(issue.IssueSummaryFields)
				}),
				issuePage.Issues.PageInfo.HasNextPage,
				issuePage.Issues.PageInfo.EndCursor,
			), nil
		},
	)
	if err != nil {
		return IssueList{}, err
	}

	return issueListFromNodePage(page), nil
}

func issueSummaryNodePage(items []IssueSummary, hasNextPage bool, endCursor *string) nodePage[IssueSummary] {
	return nodePage[IssueSummary]{Items: items, HasNextPage: hasNextPage, EndCursor: endCursor}
}

func issueListFromNodePage(page nodePage[IssueSummary]) IssueList {
	return IssueList{Issues: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}
}

// ListIssuesByTeam returns issues scoped to a resolved team, composing every
// set filter into one IssueFilter. The blocked-by filter stays on its own
// relations-based path because IssueFilter cannot express a relation target.
func ListIssuesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	filters IssueListFilters,
) (IssueList, error) {
	if filters.BlockedBy != "" {
		return listIssuesBlockedByIssue(ctx, graphqlClient, teamID, limit, filters.BlockedBy)
	}

	filter := buildIssueFilter(teamID, filters)
	page, err := collectNodePages(
		"list issues", limit, issueListPageSize,
		func(pageSize int, after *string) (nodePage[IssueSummary], error) {
			issues, err := gql.IssuesByTeamFiltered(ctx, graphqlClient, filter, &pageSize, after, boolPtr(true))
			if err != nil {
				return nodePage[IssueSummary]{}, err
			}

			return issueSummaryNodePage(
				mapNodes(issues.Issues.Nodes, func(
					issue gql.IssuesByTeamFilteredIssuesIssueConnectionNodesIssue,
				) IssueSummary {
					return issueSummaryFromFields(issue.IssueSummaryFields)
				}),
				issues.Issues.PageInfo.HasNextPage,
				issues.Issues.PageInfo.EndCursor,
			), nil
		},
	)
	if err != nil {
		return IssueList{}, err
	}

	return issueListFromNodePage(page), nil
}

func buildIssueFilter(teamID string, filters IssueListFilters) gqlmodel.LinearIssueFilter {
	filter := gqlmodel.LinearIssueFilter{
		Team: &gqlmodel.LinearIDFilter{ID: gqlmodel.LinearIDComparator{Eq: teamID}},
	}
	if filters.StateType != "" {
		filter.State = &gqlmodel.LinearWorkflowStateTypeFilter{
			Type: gqlmodel.LinearStringComparator{Eq: filters.StateType},
		}
	}
	if filters.ProjectID != "" {
		filter.Project = &gqlmodel.LinearIDFilter{ID: gqlmodel.LinearIDComparator{Eq: filters.ProjectID}}
	}
	if filters.AssigneeID != "" {
		filter.Assignee = &gqlmodel.LinearIDFilter{ID: gqlmodel.LinearIDComparator{Eq: filters.AssigneeID}}
	}
	if filters.LabelID != "" {
		filter.Labels = &gqlmodel.LinearLabelCollectionFilter{
			Some: gqlmodel.LinearIDFilter{ID: gqlmodel.LinearIDComparator{Eq: filters.LabelID}},
		}
	}
	if filters.CycleID != "" {
		filter.Cycle = &gqlmodel.LinearIDFilter{ID: gqlmodel.LinearIDComparator{Eq: filters.CycleID}}
	}
	if filters.CreatedAfter != "" || filters.CreatedBefore != "" {
		filter.CreatedAt = &gqlmodel.LinearDateComparator{
			Gte: optionalString(filters.CreatedAfter),
			Lte: optionalString(filters.CreatedBefore),
		}
	}
	if filters.UpdatedAfter != "" || filters.UpdatedBefore != "" {
		filter.UpdatedAt = &gqlmodel.LinearDateComparator{
			Gte: optionalString(filters.UpdatedAfter),
			Lte: optionalString(filters.UpdatedBefore),
		}
	}
	if filters.HasBlockers {
		filter.HasBlockedByRelations = &gqlmodel.LinearRelationExistsComparator{Eq: true}
	}
	if filters.Blocks {
		filter.HasBlockingRelations = &gqlmodel.LinearRelationExistsComparator{Eq: true}
	}

	return filter
}

func listIssuesBlockedByIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
	blockerIssueID string,
) (IssueList, error) {
	issue, err := gql.IssueBlockedIssues(ctx, graphqlClient, blockerIssueID, &limit, nil, boolPtr(true))
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

// ListIssuePriorityValues returns Linear issue priority labels.
func ListIssuePriorityValues(ctx context.Context, graphqlClient graphql.Client) ([]IssuePriorityValue, error) {
	result, err := gql.XIssuePriorityValues(ctx, graphqlClient)
	if err != nil {
		return nil, fmt.Errorf("list issue priority values: %w", err)
	}

	values := mapNodes(result.IssuePriorityValues, func(
		value gql.XIssuePriorityValuesIssuePriorityValuesIssuePriorityValue,
	) IssuePriorityValue {
		return IssuePriorityValue(value)
	})

	return values, nil
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
	issueResult, err := gql.XIssue(ctx, graphqlClient, id)
	if err != nil {
		return IssueDetail{}, fmt.Errorf("get issue %s: %w", id, err)
	}

	return detailIssue(issueResult.Issue), nil
}

func detailIssueSummary(issue gql.XIssueIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func detailIssue(issue gql.XIssueIssue) IssueDetail {
	description := ""
	if issue.Description != nil {
		description = *issue.Description
	}

	return IssueDetail{
		Summary:     detailIssueSummary(issue),
		Description: description,
	}
}

func issueSummaryFromFields(issue gql.IssueSummaryFields) IssueSummary {
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
