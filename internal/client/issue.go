package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ErrWriteInvalid marks a malformed write request.
var ErrWriteInvalid = errors.New("invalid write")

// IssueSummary is the compact issue model used by read-only commands.
type IssueSummary struct {
	ID            string  `json:"id"`
	Identifier    string  `json:"identifier"`
	Title         string  `json:"title"`
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
}

// IssueList is a page of read-only issues.
type IssueList struct {
	Issues      []IssueSummary `json:"issues"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// ListIssuesByTeam returns issues scoped to a resolved team.
func ListIssuesByTeam(ctx context.Context, graphqlClient graphql.Client, teamID string, limit int) (IssueList, error) {
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

// GetIssueByID returns a read-only issue by Linear id or identifier.
func GetIssueByID(ctx context.Context, graphqlClient graphql.Client, id string) (IssueSummary, error) {
	issue, err := IssueByID(ctx, graphqlClient, id)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("get issue %s: %w", id, err)
	}

	return detailIssueSummary(issue.Issue), nil
}

func listIssueSummary(issue IssuesByTeamIssuesIssueConnectionNodesIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func detailIssueSummary(issue IssueByIDIssue) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
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
