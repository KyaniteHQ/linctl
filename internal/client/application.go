package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ApplicationInfo is public OAuth application metadata returned by Linear.
type ApplicationInfo struct {
	ID           string `json:"id"`
	ClientID     string `json:"client_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Developer    string `json:"developer"`
	DeveloperURL string `json:"developer_url"`
	ImageURL     string `json:"image_url,omitempty"`
}

// GetApplicationInfo returns public OAuth application metadata for a client id.
func GetApplicationInfo(ctx context.Context, graphqlClient graphql.Client, clientID string) (ApplicationInfo, error) {
	result, err := gql.XApplicationInfo(ctx, graphqlClient, clientID)
	if err != nil {
		return ApplicationInfo{}, fmt.Errorf("get application info %s: %w", clientID, err)
	}

	return applicationInfoSummary(result.ApplicationInfo.ApplicationInfoFields), nil
}

func applicationInfoSummary(fields gql.ApplicationInfoFields) ApplicationInfo {
	return ApplicationInfo{
		ID:           fields.Id,
		ClientID:     fields.ClientId,
		Name:         fields.Name,
		Description:  stringValue(fields.Description),
		Developer:    fields.Developer,
		DeveloperURL: fields.DeveloperUrl,
		ImageURL:     stringValue(fields.ImageUrl),
	}
}
