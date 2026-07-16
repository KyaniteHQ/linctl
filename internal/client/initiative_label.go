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
	HasNextPage      bool                     `json:"has_next_page"`
	EndCursor        *string                  `json:"end_cursor,omitempty"`
}

//nolint:lll
type initiativeLabelNode = gql.XInitiativeLabelsInitiativeLabelsInitiativeLabelConnectionNodesInitiativeLabel

// ListInitiativeLabels returns visible Linear initiative labels.
func ListInitiativeLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (InitiativeLabelList, error) {
	result, err := gql.XInitiativeLabels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeLabelList{}, fmt.Errorf("list initiative labels: %w", err)
	}

	labels := mapNodes(result.InitiativeLabels.Nodes, func(label initiativeLabelNode) InitiativeLabelSummary {
		return initiativeLabelSummary(label.InitiativeLabelSummaryFields)
	})

	return InitiativeLabelList{
		InitiativeLabels: labels,
		HasNextPage:      result.InitiativeLabels.PageInfo.HasNextPage,
		EndCursor:        result.InitiativeLabels.PageInfo.EndCursor,
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

func initiativeLabelSummary(fields gql.InitiativeLabelSummaryFields) InitiativeLabelSummary {
	label := InitiativeLabelSummary{
		ID:            fields.Id,
		Name:          fields.Name,
		Description:   stringValue(fields.Description),
		Color:         fields.Color,
		IsGroup:       fields.IsGroup,
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
