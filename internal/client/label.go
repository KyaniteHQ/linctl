package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListLabels returns visible IssueLabels.
func ListLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (LabelList, error) {
	labels, err := IssueLabels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return LabelList{}, fmt.Errorf("list labels: %w", err)
	}

	summaries := make([]LabelSummary, 0, len(labels.IssueLabels.Nodes))
	for _, label := range labels.IssueLabels.Nodes {
		summaries = append(summaries, labelSummary(label.IssueLabelSummaryFields))
	}

	return LabelList{
		Labels:      summaries,
		HasNextPage: labels.IssueLabels.PageInfo.HasNextPage,
		EndCursor:   labels.IssueLabels.PageInfo.EndCursor,
	}, nil
}

// GetLabelByID returns one IssueLabel by id.
func GetLabelByID(ctx context.Context, graphqlClient graphql.Client, id string) (LabelSummary, error) {
	label, err := IssueLabelByID(ctx, graphqlClient, id)
	if err != nil {
		return LabelSummary{}, fmt.Errorf("get label %s: %w", id, err)
	}

	return labelSummary(label.IssueLabel.IssueLabelSummaryFields), nil
}

func labelSummary(label IssueLabelSummaryFields) LabelSummary {
	description := ""
	if label.Description != nil {
		description = *label.Description
	}
	summary := LabelSummary{
		ID:          label.Id,
		Name:        label.Name,
		Description: description,
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
