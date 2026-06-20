package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// AttachmentSummary is the compact attachment model used by read-only commands.
type AttachmentSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle,omitempty"`
	URL        string `json:"url"`
	SourceType string `json:"source_type,omitempty"`
}

// AttachmentList is a page of attachments.
type AttachmentList struct {
	Attachments []AttachmentSummary `json:"attachments"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ListAttachments returns visible issue attachments.
func ListAttachments(ctx context.Context, graphqlClient graphql.Client, limit int) (AttachmentList, error) {
	result, err := attachments(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return AttachmentList{}, fmt.Errorf("list attachments: %w", err)
	}

	summaries := make([]AttachmentSummary, 0, len(result.Attachments.Nodes))
	for _, node := range result.Attachments.Nodes {
		summaries = append(summaries, attachmentSummary(node.AttachmentSummaryFields))
	}

	return AttachmentList{
		Attachments: summaries,
		HasNextPage: result.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.Attachments.PageInfo.EndCursor,
	}, nil
}

// GetAttachmentByID returns one attachment by Linear id.
func GetAttachmentByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (AttachmentSummary, error) {
	result, err := attachment(ctx, graphqlClient, id)
	if err != nil {
		return AttachmentSummary{}, fmt.Errorf("get attachment %s: %w", id, err)
	}

	return attachmentSummary(result.Attachment.AttachmentSummaryFields), nil
}

func attachmentSummary(fields AttachmentSummaryFields) AttachmentSummary {
	return AttachmentSummary{
		ID:         fields.Id,
		Title:      fields.Title,
		Subtitle:   stringValue(fields.Subtitle),
		URL:        fields.Url,
		SourceType: stringValue(fields.SourceType),
	}
}
