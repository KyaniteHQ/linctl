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

//nolint:lll
type commentsNode = gql.XCommentsCommentsCommentConnectionNodesComment

//nolint:lll
type commentChildrenNode = gql.XComment_childrenCommentChildrenCommentConnectionNodesComment

//nolint:lll
type commentCreatedIssuesNode = gql.XComment_createdIssuesCommentCreatedIssuesIssueConnectionNodesIssue

//nolint:lll
type issueCommentsNode = gql.XIssue_commentsIssueCommentsCommentConnectionNodesComment

type commentsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type commentScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	commentID     string
}

type issueCommentsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
	issueID       string
	identifier    string
}

// ListComments returns visible comments across parent entity types.
func ListComments(ctx context.Context, graphqlClient graphql.Client, limit int) (CommentList, error) {
	query := commentsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list comments", limit, defaultListPageSize,
		query.page,
		commentsNodeSummary,
	)
	if err != nil {
		return CommentList{}, err
	}

	return CommentList{Comments: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
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
	query := &commentScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list comment children "+id, limit, defaultListPageSize,
		query.children,
		commentChildrenNodeSummary,
	)
	if err != nil {
		return CommentChildList{}, err
	}

	return CommentChildList{
		CommentID:   query.commentID,
		Comments:    page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

// ListCommentCreatedIssues returns issues created from a comment.
func ListCommentCreatedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	query := commentScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list comment created issues "+id, limit, defaultListPageSize,
		query.createdIssues,
		commentCreatedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// ListIssueComments returns comments for one issue by Linear id or identifier.
func ListIssueComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueCommentList, error) {
	query := &issueCommentsQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list issue comments "+id, limit, defaultListPageSize,
		query.comments,
		issueCommentSummary,
	)
	if err != nil {
		return IssueCommentList{}, err
	}

	return IssueCommentList{
		IssueID:     query.issueID,
		Identifier:  query.identifier,
		Comments:    page.Items,
		HasNextPage: page.HasNextPage,
		EndCursor:   page.EndCursor,
	}, nil
}

func (query commentsQuery) page(pageSize int, after *string) ([]commentsNode, bool, *string, error) {
	result, err := gql.XComments(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Comments.Nodes,
		result.Comments.PageInfo.HasNextPage,
		result.Comments.PageInfo.EndCursor,
		nil
}

func (query *commentScopedQuery) children(
	pageSize int,
	after *string,
) ([]commentChildrenNode, bool, *string, error) {
	result, err := gql.XComment_children(
		query.ctx, query.graphqlClient, stringPtr(query.id), nil, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.commentID = result.Comment.Id

	return result.Comment.Children.Nodes,
		result.Comment.Children.PageInfo.HasNextPage,
		result.Comment.Children.PageInfo.EndCursor,
		nil
}

func (query commentScopedQuery) createdIssues(
	pageSize int,
	after *string,
) ([]commentCreatedIssuesNode, bool, *string, error) {
	result, err := gql.XComment_createdIssues(
		query.ctx, query.graphqlClient, stringPtr(query.id), nil, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Comment.CreatedIssues.Nodes,
		result.Comment.CreatedIssues.PageInfo.HasNextPage,
		result.Comment.CreatedIssues.PageInfo.EndCursor,
		nil
}

func (query *issueCommentsQuery) comments(
	pageSize int,
	after *string,
) ([]issueCommentsNode, bool, *string, error) {
	result, err := gql.XIssue_comments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	query.issueID = result.Issue.Id
	query.identifier = result.Issue.Identifier

	return result.Issue.Comments.Nodes,
		result.Issue.Comments.PageInfo.HasNextPage,
		result.Issue.Comments.PageInfo.EndCursor,
		nil
}

func commentsNodeSummary(node commentsNode) CommentSummary {
	return topLevelCommentSummary(node.TopLevelCommentSummaryFields)
}

func commentChildrenNodeSummary(comment commentChildrenNode) CommentMetadataSummary {
	return commentMetadataSummary(comment.CommentMetadataFields)
}

func commentCreatedIssuesNodeSummary(issue commentCreatedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func issueCommentSummary(comment issueCommentsNode) IssueCommentSummary {
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
