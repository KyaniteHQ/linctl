package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// IssueCommentSummary is the compact read model for issue comments.
type IssueCommentSummary struct {
	ID          string `json:"id"`
	Body        string `json:"body"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	ParentID    string `json:"parent_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// IssueCommentList is a page of comments for one issue.
type IssueCommentList struct {
	IssueID     string                `json:"issue_id"`
	Identifier  string                `json:"identifier"`
	Comments    []IssueCommentSummary `json:"comments"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// ListIssueComments returns comments for one issue by Linear id or identifier.
func ListIssueComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCommentList, error) {
	comments, err := IssueComments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCommentList{}, fmt.Errorf("list issue comments %s: %w", id, err)
	}

	summaries := make([]IssueCommentSummary, 0, len(comments.Issue.Comments.Nodes))
	for _, comment := range comments.Issue.Comments.Nodes {
		summaries = append(summaries, issueCommentSummary(comment))
	}

	return IssueCommentList{
		IssueID:     comments.Issue.Id,
		Identifier:  comments.Issue.Identifier,
		Comments:    summaries,
		HasNextPage: comments.Issue.Comments.PageInfo.HasNextPage,
		EndCursor:   comments.Issue.Comments.PageInfo.EndCursor,
	}, nil
}

func issueCommentSummary(comment IssueCommentsIssueCommentsCommentConnectionNodesComment) IssueCommentSummary {
	userID := ""
	userName := ""
	displayName := ""
	if comment.User != nil {
		userID = comment.User.Id
		userName = comment.User.Name
		displayName = comment.User.DisplayName
	}
	parentID := ""
	if comment.ParentId != nil {
		parentID = *comment.ParentId
	}

	return IssueCommentSummary{
		ID:          comment.Id,
		Body:        comment.Body,
		URL:         comment.Url,
		CreatedAt:   comment.CreatedAt,
		ParentID:    parentID,
		UserID:      userID,
		UserName:    userName,
		DisplayName: displayName,
	}
}
