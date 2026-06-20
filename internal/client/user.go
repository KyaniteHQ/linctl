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
