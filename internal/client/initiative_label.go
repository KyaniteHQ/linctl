package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// InitiativeLabelSummary is the compact initiative label model used by read-only commands.
type InitiativeLabelSummary struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Color         string `json:"color"`
	IsGroup       bool   `json:"is_group"`
	OrgID         string `json:"org_id,omitempty"`
	ParentID      string `json:"parent_id,omitempty"`
	ParentName    string `json:"parent_name,omitempty"`
	ParentColor   string `json:"parent_color,omitempty"`
	LastAppliedAt string `json:"last_applied_at,omitempty"`
	RetiredAt     string `json:"retired_at,omitempty"`
	ArchivedAt    string `json:"archived_at,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// InitiativeLabelList is a page of Linear initiative labels.
type InitiativeLabelList struct {
	InitiativeLabels []InitiativeLabelSummary `json:"initiative_labels"`
	Page
}

//nolint:lll
type initiativeLabelNode = gql.XInitiativeLabelsInitiativeLabelsInitiativeLabelConnectionNodesInitiativeLabel

// ListInitiativeLabels returns visible Linear initiative labels.
func ListInitiativeLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeLabelList, error) {
	page, err := listConnection(
		"list initiative labels", limit, defaultListPageSize,
		func(pageSize int, after *string) ([]initiativeLabelNode, bool, *string, error) {
			result, err := gql.XInitiativeLabels(ctx, graphqlClient, intPtr(pageSize), after, boolPtr(true))
			if err != nil {
				return nil, false, nil, err
			}

			return result.InitiativeLabels.Nodes,
				result.InitiativeLabels.PageInfo.HasNextPage,
				result.InitiativeLabels.PageInfo.EndCursor,
				nil
		},
		initiativeLabelNodeSummary,
	)
	if err != nil {
		return InitiativeLabelList{}, err
	}

	return InitiativeLabelList{
		InitiativeLabels: page.Items,
		Page:             page.Page,
	}, nil
}

// GetInitiativeLabelByID returns one Linear initiative label by id.
func GetInitiativeLabelByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (InitiativeLabelSummary, error) {
	result, err := gql.XInitiativeLabel(ctx, graphqlClient, id)
	if err != nil {
		return InitiativeLabelSummary{}, fmt.Errorf("get initiative label %s: %w", id, err)
	}

	return initiativeLabelSummary(result.InitiativeLabel.InitiativeLabelSummaryFields), nil
}

func initiativeLabelNodeSummary(label initiativeLabelNode) InitiativeLabelSummary {
	return initiativeLabelSummary(label.InitiativeLabelSummaryFields)
}

func initiativeLabelSummary(fields gql.InitiativeLabelSummaryFields) InitiativeLabelSummary {
	label := InitiativeLabelSummary{
		ID:            fields.Id,
		Name:          fields.Name,
		Description:   stringValue(fields.Description),
		Color:         fields.Color,
		IsGroup:       fields.IsGroup,
		OrgID:         fields.Organization.Id,
		LastAppliedAt: stringValue(fields.LastAppliedAt),
		RetiredAt:     stringValue(fields.RetiredAt),
		ArchivedAt:    stringValue(fields.ArchivedAt),
		CreatedAt:     fields.CreatedAt,
		UpdatedAt:     fields.UpdatedAt,
	}
	if fields.Parent != nil {
		label.ParentID = fields.Parent.Id
		label.ParentName = fields.Parent.Name
		label.ParentColor = fields.Parent.Color
	}

	return label
}
