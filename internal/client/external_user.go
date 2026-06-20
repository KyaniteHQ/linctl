package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ExternalUserSummary is the compact ExternalUser model used by read-only commands.
type ExternalUserSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	LastSeen    string `json:"last_seen,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ArchivedAt  string `json:"archived_at,omitempty"`
}

// ExternalUserList is a page of ExternalUsers.
type ExternalUserList struct {
	ExternalUsers []ExternalUserSummary `json:"external_users"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

// ListExternalUsers returns ExternalUsers visible to the authenticated user.
func ListExternalUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (ExternalUserList, error) {
	result, err := externalUsers(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ExternalUserList{}, fmt.Errorf("list external users: %w", err)
	}

	summaries := make([]ExternalUserSummary, 0, len(result.ExternalUsers.Nodes))
	for _, node := range result.ExternalUsers.Nodes {
		summaries = append(summaries, externalUserSummary(node.ExternalUserSummaryFields))
	}

	return ExternalUserList{
		ExternalUsers: summaries,
		HasNextPage:   result.ExternalUsers.PageInfo.HasNextPage,
		EndCursor:     result.ExternalUsers.PageInfo.EndCursor,
	}, nil
}

// GetExternalUserByID returns one ExternalUser by id.
func GetExternalUserByID(ctx context.Context, graphqlClient graphql.Client, id string) (ExternalUserSummary, error) {
	result, err := externalUser(ctx, graphqlClient, id)
	if err != nil {
		return ExternalUserSummary{}, fmt.Errorf("get external user %s: %w", id, err)
	}

	return externalUserSummary(result.ExternalUser.ExternalUserSummaryFields), nil
}

func externalUserSummary(fields ExternalUserSummaryFields) ExternalUserSummary {
	return ExternalUserSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		DisplayName: fields.DisplayName,
		AvatarURL:   stringValue(fields.AvatarUrl),
		LastSeen:    stringValue(fields.LastSeen),
		CreatedAt:   fields.CreatedAt,
		UpdatedAt:   fields.UpdatedAt,
		ArchivedAt:  stringValue(fields.ArchivedAt),
	}
}
