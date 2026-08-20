package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// LabelSummary is the compact IssueLabel model used by label commands.
type LabelSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color"`
	IsGroup     bool   `json:"is_group"`
	TeamID      string `json:"team_id,omitempty"`
	TeamKey     string `json:"team_key,omitempty"`
	TeamName    string `json:"team_name,omitempty"`
}

// LabelList is a page of labels.
type LabelList struct {
	Labels      []LabelSummary `json:"labels"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// LabelChildList is a page of child labels for one IssueLabel group.
type LabelChildList struct {
	LabelID     string         `json:"label_id"`
	LabelName   string         `json:"label_name"`
	Labels      []LabelSummary `json:"labels"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// LabelIssueList is a page of issues associated with one IssueLabel.
type LabelIssueList struct {
	LabelID     string         `json:"label_id"`
	LabelName   string         `json:"label_name"`
	Issues      []IssueSummary `json:"issues"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

//nolint:lll
type labelsNode = gql.IssueLabelsIssueLabelsIssueLabelConnectionNodesIssueLabel

//nolint:lll
type labelChildrenNode = gql.XIssueLabel_childrenIssueLabelChildrenIssueLabelConnectionNodesIssueLabel

//nolint:lll
type labelIssuesNode = gql.XIssueLabel_issuesIssueLabelIssuesIssueConnectionNodesIssue

type labelsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type labelScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	labelID       string
	labelName     string
}

// ListLabels returns visible IssueLabels.
func ListLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (LabelList, error) {
	query := labelsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list labels", limit, defaultListPageSize,
		query.page,
		labelsNodeSummary,
	)
	if err != nil {
		return LabelList{}, err
	}

	return LabelList{Labels: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetLabelByID returns one IssueLabel by id.
func GetLabelByID(ctx context.Context, graphqlClient graphql.Client, id string) (LabelSummary, error) {
	label, err := gql.XIssueLabel(ctx, graphqlClient, id)
	if err != nil {
		return LabelSummary{}, fmt.Errorf("get label %s: %w", id, err)
	}

	return labelSummary(label.IssueLabel.IssueLabelSummaryFields), nil
}

// ListLabelChildren returns child labels under one IssueLabel group.
func ListLabelChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (LabelChildList, error) {
	query := &labelScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list label children "+id, limit, defaultListPageSize,
		query.children,
		labelChildrenNodeSummary,
	)
	if err != nil {
		return LabelChildList{}, err
	}

	return LabelChildList{
		LabelID:     query.labelID,
		LabelName:   query.labelName,
		Labels:      page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListLabelIssues returns issues associated with one IssueLabel.
func ListLabelIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (LabelIssueList, error) {
	query := &labelScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list label issues "+id, limit, defaultListPageSize,
		query.issues,
		labelIssuesNodeSummary,
	)
	if err != nil {
		return LabelIssueList{}, err
	}

	return LabelIssueList{
		LabelID:     query.labelID,
		LabelName:   query.labelName,
		Issues:      page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

func (query labelsQuery) page(pageSize int, after *string) ([]labelsNode, bool, *string, error) {
	result, err := gql.IssueLabels(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.IssueLabels.Nodes,
		result.IssueLabels.PageInfo.HasNextPage,
		result.IssueLabels.PageInfo.EndCursor,
		nil
}

func (query *labelScopedQuery) children(pageSize int, after *string) ([]labelChildrenNode, bool, *string, error) {
	result, err := gql.XIssueLabel_children(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.labelID = result.IssueLabel.Id
	query.labelName = result.IssueLabel.Name

	return result.IssueLabel.Children.Nodes,
		result.IssueLabel.Children.PageInfo.HasNextPage,
		result.IssueLabel.Children.PageInfo.EndCursor,
		nil
}

func (query *labelScopedQuery) issues(pageSize int, after *string) ([]labelIssuesNode, bool, *string, error) {
	result, err := gql.XIssueLabel_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.labelID = result.IssueLabel.Id
	query.labelName = result.IssueLabel.Name

	return result.IssueLabel.Issues.Nodes,
		result.IssueLabel.Issues.PageInfo.HasNextPage,
		result.IssueLabel.Issues.PageInfo.EndCursor,
		nil
}

func labelsNodeSummary(label labelsNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func labelChildrenNodeSummary(label labelChildrenNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func labelIssuesNodeSummary(issue labelIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func labelSummary(label gql.IssueLabelSummaryFields) LabelSummary {
	summary := LabelSummary{
		ID:          label.Id,
		Name:        label.Name,
		Description: stringValue(label.Description),
		Color:       label.Color,
		IsGroup:     label.IsGroup,
	}
	if label.Team != nil {
		summary.TeamID = label.Team.Id
		summary.TeamKey = label.Team.Key
		summary.TeamName = label.Team.Name
	}

	return summary
}
