package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Page
}

//nolint:lll
type externalUsersNode = gql.XExternalUsersExternalUsersExternalUserConnectionNodesExternalUser

type externalUsersQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListExternalUsers returns ExternalUsers visible to the authenticated user.
func ListExternalUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (ExternalUserList, error) {
	query := externalUsersQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list external users", limit, defaultListPageSize,
		query.page,
		externalUsersNodeSummary,
	)
	if err != nil {
		return ExternalUserList{}, err
	}

	return ExternalUserList{
		ExternalUsers: page.Items,
		Page:          page.Page,
	}, nil
}

// GetExternalUserByID returns one ExternalUser by id.
func GetExternalUserByID(ctx context.Context, graphqlClient graphql.Client, id string) (ExternalUserSummary, error) {
	result, err := gql.XExternalUser(ctx, graphqlClient, id)
	if err != nil {
		return ExternalUserSummary{}, fmt.Errorf("get external user %s: %w", id, err)
	}

	return externalUserSummary(result.ExternalUser.ExternalUserSummaryFields), nil
}

func (query externalUsersQuery) page(pageSize int, after *string) ([]externalUsersNode, bool, *string, error) {
	result, err := gql.XExternalUsers(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.ExternalUsers.Nodes,
		result.ExternalUsers.PageInfo.HasNextPage,
		result.ExternalUsers.PageInfo.EndCursor,
		nil
}

func externalUsersNodeSummary(node externalUsersNode) ExternalUserSummary {
	return externalUserSummary(node.ExternalUserSummaryFields)
}

func externalUserSummary(fields gql.ExternalUserSummaryFields) ExternalUserSummary {
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
