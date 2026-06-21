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

// CommentSummary is the compact read model for top-level comment reads.
type CommentSummary struct {
	ID                 string  `json:"id"`
	Body               string  `json:"body"`
	URL                string  `json:"url"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	EditedAt           *string `json:"edited_at,omitempty"`
	ResolvedAt         *string `json:"resolved_at,omitempty"`
	ParentID           string  `json:"parent_id,omitempty"`
	IssueID            string  `json:"issue_id,omitempty"`
	ProjectID          string  `json:"project_id,omitempty"`
	ProjectUpdateID    string  `json:"project_update_id,omitempty"`
	InitiativeID       string  `json:"initiative_id,omitempty"`
	InitiativeUpdateID string  `json:"initiative_update_id,omitempty"`
	DocumentContentID  string  `json:"document_content_id,omitempty"`
	UserID             string  `json:"user_id,omitempty"`
	UserName           string  `json:"user_name,omitempty"`
	DisplayName        string  `json:"display_name,omitempty"`
}

// CommentMetadataSummary is a body-free comment read model for parent-scoped comment lists.
type CommentMetadataSummary struct {
	ID                 string  `json:"id"`
	URL                string  `json:"url"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	EditedAt           *string `json:"edited_at,omitempty"`
	ResolvedAt         *string `json:"resolved_at,omitempty"`
	ParentID           string  `json:"parent_id,omitempty"`
	IssueID            string  `json:"issue_id,omitempty"`
	ProjectID          string  `json:"project_id,omitempty"`
	ProjectUpdateID    string  `json:"project_update_id,omitempty"`
	InitiativeID       string  `json:"initiative_id,omitempty"`
	InitiativeUpdateID string  `json:"initiative_update_id,omitempty"`
	DocumentContentID  string  `json:"document_content_id,omitempty"`
	UserID             string  `json:"user_id,omitempty"`
	UserName           string  `json:"user_name,omitempty"`
	DisplayName        string  `json:"display_name,omitempty"`
}

// ActorBotSummary is compact bot actor metadata without external payload details.
type ActorBotSummary struct {
	ID              string `json:"id,omitempty"`
	Type            string `json:"type"`
	SubType         string `json:"sub_type,omitempty"`
	Name            string `json:"name,omitempty"`
	UserDisplayName string `json:"user_display_name,omitempty"`
	AvatarURL       string `json:"avatar_url,omitempty"`
}

// CommentBotActor is the optional bot actor attached to a comment.
type CommentBotActor struct {
	CommentID string           `json:"comment_id"`
	Bot       *ActorBotSummary `json:"bot,omitempty"`
}

// IssueCommentList is a page of comments for one issue.
type IssueCommentList struct {
	IssueID     string                `json:"issue_id"`
	Identifier  string                `json:"identifier"`
	Comments    []IssueCommentSummary `json:"comments"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// CommentList is a page of comments visible to the authenticated user.
type CommentList struct {
	Comments    []CommentSummary `json:"comments"`
	HasNextPage bool             `json:"has_next_page"`
	EndCursor   *string          `json:"end_cursor,omitempty"`
}

// CommentChildList is a page of body-free child comment metadata.
type CommentChildList struct {
	CommentID   string                   `json:"comment_id"`
	Comments    []CommentMetadataSummary `json:"comments"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ListComments returns visible comments across parent entity types.
func ListComments(ctx context.Context, graphqlClient graphql.Client, limit int) (CommentList, error) {
	commentsPage, err := comments(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CommentList{}, fmt.Errorf("list comments: %w", err)
	}

	summaries := make([]CommentSummary, 0, len(commentsPage.Comments.Nodes))
	for _, node := range commentsPage.Comments.Nodes {
		summaries = append(summaries, topLevelCommentSummary(node.TopLevelCommentSummaryFields))
	}

	return CommentList{
		Comments:    summaries,
		HasNextPage: commentsPage.Comments.PageInfo.HasNextPage,
		EndCursor:   commentsPage.Comments.PageInfo.EndCursor,
	}, nil
}

// GetCommentByID returns one comment by Linear id.
func GetCommentByID(ctx context.Context, graphqlClient graphql.Client, id string) (CommentSummary, error) {
	commentResponse, err := comment(ctx, graphqlClient, stringPtr(id), nil)
	if err != nil {
		return CommentSummary{}, fmt.Errorf("get comment %s: %w", id, err)
	}

	return topLevelCommentSummary(commentResponse.Comment.TopLevelCommentSummaryFields), nil
}

// GetCommentBotActor returns the bot actor that created a comment, when present.
func GetCommentBotActor(ctx context.Context, graphqlClient graphql.Client, id string) (CommentBotActor, error) {
	result, err := comment_botActor(ctx, graphqlClient, stringPtr(id), nil)
	if err != nil {
		return CommentBotActor{}, fmt.Errorf("get comment bot actor %s: %w", id, err)
	}

	return CommentBotActor{
		CommentID: result.Comment.Id,
		Bot:       actorBotSummary(result.Comment.BotActor),
	}, nil
}

// ListCommentChildren returns child comments without body content.
func ListCommentChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (CommentChildList, error) {
	result, err := comment_children(ctx, graphqlClient, stringPtr(id), nil, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CommentChildList{}, fmt.Errorf("list comment children %s: %w", id, err)
	}

	comments := make([]CommentMetadataSummary, 0, len(result.Comment.Children.Nodes))
	for _, comment := range result.Comment.Children.Nodes {
		comments = append(comments, commentMetadataSummary(comment.CommentMetadataFields))
	}

	return CommentChildList{
		CommentID:   result.Comment.Id,
		Comments:    comments,
		HasNextPage: result.Comment.Children.PageInfo.HasNextPage,
		EndCursor:   result.Comment.Children.PageInfo.EndCursor,
	}, nil
}

// ListCommentCreatedIssues returns issues created from a comment.
func ListCommentCreatedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := comment_createdIssues(ctx, graphqlClient, stringPtr(id), nil, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list comment created issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.Comment.CreatedIssues.Nodes))
	for _, issue := range result.Comment.CreatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Comment.CreatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Comment.CreatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListIssueComments returns comments for one issue by Linear id or identifier.
func ListIssueComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCommentList, error) {
	comments, err := issue_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
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

func issueCommentSummary(comment issue_commentsIssueCommentsCommentConnectionNodesComment) IssueCommentSummary {
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

func topLevelCommentSummary(comment TopLevelCommentSummaryFields) CommentSummary {
	userID := ""
	userName := ""
	displayName := ""
	if comment.User != nil {
		userID = comment.User.Id
		userName = comment.User.Name
		displayName = comment.User.DisplayName
	}

	return CommentSummary{
		ID:                 comment.Id,
		Body:               comment.Body,
		URL:                comment.Url,
		CreatedAt:          comment.CreatedAt,
		UpdatedAt:          comment.UpdatedAt,
		EditedAt:           comment.EditedAt,
		ResolvedAt:         comment.ResolvedAt,
		ParentID:           stringValue(comment.ParentId),
		IssueID:            stringValue(comment.IssueId),
		ProjectID:          stringValue(comment.ProjectId),
		ProjectUpdateID:    stringValue(comment.ProjectUpdateId),
		InitiativeID:       stringValue(comment.InitiativeId),
		InitiativeUpdateID: stringValue(comment.InitiativeUpdateId),
		DocumentContentID:  stringValue(comment.DocumentContentId),
		UserID:             userID,
		UserName:           userName,
		DisplayName:        displayName,
	}
}

func commentMetadataSummary(comment CommentMetadataFields) CommentMetadataSummary {
	userID := ""
	userName := ""
	displayName := ""
	if comment.User != nil {
		userID = comment.User.Id
		userName = comment.User.Name
		displayName = comment.User.DisplayName
	}

	return CommentMetadataSummary{
		ID:                 comment.Id,
		URL:                comment.Url,
		CreatedAt:          comment.CreatedAt,
		UpdatedAt:          comment.UpdatedAt,
		EditedAt:           comment.EditedAt,
		ResolvedAt:         comment.ResolvedAt,
		ParentID:           stringValue(comment.ParentId),
		IssueID:            stringValue(comment.IssueId),
		ProjectID:          stringValue(comment.ProjectId),
		ProjectUpdateID:    stringValue(comment.ProjectUpdateId),
		InitiativeID:       stringValue(comment.InitiativeId),
		InitiativeUpdateID: stringValue(comment.InitiativeUpdateId),
		DocumentContentID:  stringValue(comment.DocumentContentId),
		UserID:             userID,
		UserName:           userName,
		DisplayName:        displayName,
	}
}

func actorBotSummary(bot *comment_botActorCommentBotActorActorBot) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return &ActorBotSummary{
		ID:              stringValue(bot.Id),
		Type:            bot.Type,
		SubType:         stringValue(bot.SubType),
		Name:            stringValue(bot.Name),
		UserDisplayName: stringValue(bot.UserDisplayName),
		AvatarURL:       stringValue(bot.AvatarUrl),
	}
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
