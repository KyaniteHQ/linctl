package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// UserSummary is the compact User model used by user commands.
type UserSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Active      bool   `json:"active"`
	Guest       bool   `json:"guest"`
	Admin       bool   `json:"admin"`
}

// UserList is a page of users.
type UserList struct {
	Users       []UserSummary `json:"users"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// DraftSummary is the compact saved draft model used by viewer-scoped draft reads.
type DraftSummary struct {
	ID          string `json:"id"`
	ParentType  string `json:"parent_type"`
	ParentID    string `json:"parent_id"`
	ParentKey   string `json:"parent_key,omitempty"`
	ParentTitle string `json:"parent_title,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ArchivedAt  string `json:"archived_at,omitempty"`
}

// DraftList is a page of the authenticated user's saved drafts.
type DraftList struct {
	Drafts      []DraftSummary `json:"drafts"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// ListUsers returns visible users.
func ListUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	userPage, err := users(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true), boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list users: %w", err)
	}

	summaries := make([]UserSummary, 0, len(userPage.Users.Nodes))
	for _, user := range userPage.Users.Nodes {
		summaries = append(summaries, userSummary(user.UserSummaryFields))
	}

	return UserList{
		Users:       summaries,
		HasNextPage: userPage.Users.PageInfo.HasNextPage,
		EndCursor:   userPage.Users.PageInfo.EndCursor,
	}, nil
}

// GetUserByID returns one User by id.
func GetUserByID(ctx context.Context, graphqlClient graphql.Client, id string) (UserSummary, error) {
	userResult, err := user(ctx, graphqlClient, id)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get user %s: %w", id, err)
	}

	return userSummary(userResult.User.UserSummaryFields), nil
}

// GetViewerUser returns the authenticated User.
func GetViewerUser(ctx context.Context, graphqlClient graphql.Client) (UserSummary, error) {
	userResult, err := viewer(ctx, graphqlClient)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get viewer user: %w", err)
	}

	return userSummary(userResult.Viewer.UserSummaryFields), nil
}

// ListViewerDrafts returns the authenticated user's saved draft metadata.
func ListViewerDrafts(ctx context.Context, graphqlClient graphql.Client, limit int) (DraftList, error) {
	draftPage, err := viewer_drafts(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return DraftList{}, fmt.Errorf("list viewer drafts: %w", err)
	}

	summaries := make([]DraftSummary, 0, len(draftPage.Viewer.Drafts.Nodes))
	for _, draft := range draftPage.Viewer.Drafts.Nodes {
		summaries = append(summaries, draftSummary(draft.DraftSummaryFields))
	}

	return DraftList{
		Drafts:      summaries,
		HasNextPage: draftPage.Viewer.Drafts.PageInfo.HasNextPage,
		EndCursor:   draftPage.Viewer.Drafts.PageInfo.EndCursor,
	}, nil
}

func userSummary(user UserSummaryFields) UserSummary {
	return UserSummary{
		ID:          user.Id,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Active:      user.Active,
		Guest:       user.Guest,
		Admin:       user.Admin,
	}
}

func draftSummary(draft DraftSummaryFields) DraftSummary {
	summary := DraftSummary{
		ID:         draft.Id,
		CreatedAt:  draft.CreatedAt,
		UpdatedAt:  draft.UpdatedAt,
		ArchivedAt: stringValue(draft.ArchivedAt),
	}
	switch {
	case draft.Issue != nil:
		summary.ParentType = "issue"
		summary.ParentID = draft.Issue.Id
		summary.ParentKey = draft.Issue.Identifier
		summary.ParentTitle = draft.Issue.Title
	case draft.Project != nil:
		summary.ParentType = "project"
		summary.ParentID = draft.Project.Id
		summary.ParentTitle = draft.Project.Name
	case draft.ProjectUpdate != nil:
		summary.ParentType = "project_update"
		summary.ParentID = draft.ProjectUpdate.Id
	case draft.Initiative != nil:
		summary.ParentType = "initiative"
		summary.ParentID = draft.Initiative.Id
		summary.ParentTitle = draft.Initiative.Name
	case draft.InitiativeUpdate != nil:
		summary.ParentType = "initiative_update"
		summary.ParentID = draft.InitiativeUpdate.Id
	case draft.ParentComment != nil:
		summary.ParentType = "comment"
		summary.ParentID = draft.ParentComment.Id
	case draft.CustomerNeed != nil:
		summary.ParentType = "customer_need"
		summary.ParentID = draft.CustomerNeed.Id
	case draft.Team != nil:
		summary.ParentType = "team"
		summary.ParentID = draft.Team.Id
		summary.ParentKey = draft.Team.Key
		summary.ParentTitle = draft.Team.Name
	default:
		summary.ParentType = "unknown"
	}

	return summary
}
