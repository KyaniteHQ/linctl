package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	commentsPage, err := gql.XComments(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CommentList{}, fmt.Errorf("list comments: %w", err)
	}

	summaries := mapNodes(commentsPage.Comments.Nodes, func(
		node gql.XCommentsCommentsCommentConnectionNodesComment,
	) CommentSummary {
		return topLevelCommentSummary(node.TopLevelCommentSummaryFields)
	})

	return CommentList{
		Comments:    summaries,
		HasNextPage: commentsPage.Comments.PageInfo.HasNextPage,
		EndCursor:   commentsPage.Comments.PageInfo.EndCursor,
	}, nil
}

// GetCommentByID returns one comment by Linear id.
func GetCommentByID(ctx context.Context, graphqlClient graphql.Client, id string) (CommentSummary, error) {
	commentResponse, err := gql.XComment(ctx, graphqlClient, stringPtr(id), nil)
	if err != nil {
		return CommentSummary{}, fmt.Errorf("get comment %s: %w", id, err)
	}

	return topLevelCommentSummary(commentResponse.Comment.TopLevelCommentSummaryFields), nil
}

// GetCommentBotActor returns the bot actor that created a comment, when present.
func GetCommentBotActor(ctx context.Context, graphqlClient graphql.Client, id string) (CommentBotActor, error) {
	result, err := gql.XComment_botActor(ctx, graphqlClient, stringPtr(id), nil)
	if err != nil {
		return CommentBotActor{}, fmt.Errorf("get comment bot actor %s: %w", id, err)
	}

	return CommentBotActor{
		CommentID: result.Comment.Id,
		Bot:       commentActorBotSummary(result.Comment.BotActor),
	}, nil
}

// ListCommentChildren returns child comments without body content.
func ListCommentChildren(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (CommentChildList, error) {
	result, err := gql.XComment_children(ctx, graphqlClient, stringPtr(id), nil, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CommentChildList{}, fmt.Errorf("list comment children %s: %w", id, err)
	}

	comments := mapNodes(result.Comment.Children.Nodes, func(
		comment gql.XComment_childrenCommentChildrenCommentConnectionNodesComment,
	) CommentMetadataSummary {
		return commentMetadataSummary(comment.CommentMetadataFields)
	})

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
	result, err := gql.XComment_createdIssues(ctx, graphqlClient, stringPtr(id), nil, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueList{}, fmt.Errorf("list comment created issues %s: %w", id, err)
	}

	issues := mapNodes(result.Comment.CreatedIssues.Nodes, func(
		issue gql.XComment_createdIssuesCommentCreatedIssuesIssueConnectionNodesIssue,
	) IssueSummary {
		return issueSummaryFromFields(issue.IssueSummaryFields)
	})

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
	comments, err := gql.XIssue_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return IssueCommentList{}, fmt.Errorf("list issue comments %s: %w", id, err)
	}

	summaries := mapNodes(comments.Issue.Comments.Nodes, issueCommentSummary)

	return IssueCommentList{
		IssueID:     comments.Issue.Id,
		Identifier:  comments.Issue.Identifier,
		Comments:    summaries,
		HasNextPage: comments.Issue.Comments.PageInfo.HasNextPage,
		EndCursor:   comments.Issue.Comments.PageInfo.EndCursor,
	}, nil
}

func issueCommentSummary(comment gql.XIssue_commentsIssueCommentsCommentConnectionNodesComment) IssueCommentSummary {
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

func topLevelCommentSummary(comment gql.TopLevelCommentSummaryFields) CommentSummary {
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

func commentMetadataSummary(comment gql.CommentMetadataFields) CommentMetadataSummary {
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

func actorBotSummary(fields *gql.ActorBotSummaryFields) *ActorBotSummary {
	if fields == nil {
		return nil
	}

	return &ActorBotSummary{
		ID:              stringValue(fields.Id),
		Type:            fields.Type,
		SubType:         stringValue(fields.SubType),
		Name:            stringValue(fields.Name),
		UserDisplayName: stringValue(fields.UserDisplayName),
		AvatarURL:       stringValue(fields.AvatarUrl),
	}
}

func commentActorBotSummary(bot *gql.XComment_botActorCommentBotActorActorBot) *ActorBotSummary {
	if bot == nil {
		return nil
	}

	return actorBotSummary(&bot.ActorBotSummaryFields)
}
