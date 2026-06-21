package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// CustomerNeedSummary is the compact customer need model used by read-only commands.
type CustomerNeedSummary struct {
	ID           string  `json:"id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	ArchivedAt   string  `json:"archived_at,omitempty"`
	Priority     float64 `json:"priority"`
	Body         string  `json:"body,omitempty"`
	Content      string  `json:"content,omitempty"`
	URL          string  `json:"url,omitempty"`
	CustomerID   string  `json:"customer_id,omitempty"`
	CustomerName string  `json:"customer_name,omitempty"`
	IssueID      string  `json:"issue_id,omitempty"`
	Issue        string  `json:"issue,omitempty"`
	IssueTitle   string  `json:"issue_title,omitempty"`
	ProjectID    string  `json:"project_id,omitempty"`
	ProjectName  string  `json:"project_name,omitempty"`
}

// CustomerNeedList is a page of Linear customer needs.
type CustomerNeedList struct {
	Needs       []CustomerNeedSummary `json:"customer_needs"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// CustomerNeedProjectAttachment is the metadata-only ProjectAttachment linked to a customer need.
type CustomerNeedProjectAttachment struct {
	CustomerNeedID string             `json:"customer_need_id"`
	Attachment     *AttachmentSummary `json:"attachment,omitempty"`
}

// ListCustomerNeeds returns visible Linear customer needs.
func ListCustomerNeeds(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerNeedList, error) {
	result, err := customerNeeds(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CustomerNeedList{}, fmt.Errorf("list customer needs: %w", err)
	}

	summaries := make([]CustomerNeedSummary, 0, len(result.CustomerNeeds.Nodes))
	for _, node := range result.CustomerNeeds.Nodes {
		summaries = append(summaries, customerNeedSummary(node.CustomerNeedSummaryFields))
	}

	return CustomerNeedList{
		Needs:       summaries,
		HasNextPage: result.CustomerNeeds.PageInfo.HasNextPage,
		EndCursor:   result.CustomerNeeds.PageInfo.EndCursor,
	}, nil
}

// GetCustomerNeedByID returns one Linear customer need by id.
func GetCustomerNeedByID(ctx context.Context, graphqlClient graphql.Client, id string) (CustomerNeedSummary, error) {
	result, err := customerNeed(ctx, graphqlClient, &id)
	if err != nil {
		return CustomerNeedSummary{}, fmt.Errorf("get customer need %s: %w", id, err)
	}

	return customerNeedSummary(result.CustomerNeed.CustomerNeedSummaryFields), nil
}

// GetCustomerNeedProjectAttachment returns the metadata-only ProjectAttachment linked to one customer need.
func GetCustomerNeedProjectAttachment(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomerNeedProjectAttachment, error) {
	result, err := customerNeed_projectAttachment(ctx, graphqlClient, &id)
	if err != nil {
		return CustomerNeedProjectAttachment{}, fmt.Errorf("get customer need project attachment %s: %w", id, err)
	}

	attachment := (*AttachmentSummary)(nil)
	if result.CustomerNeed.ProjectAttachment != nil {
		summary := projectAttachmentSummary(result.CustomerNeed.ProjectAttachment.ProjectAttachmentSummaryFields)
		attachment = &summary
	}

	return CustomerNeedProjectAttachment{
		CustomerNeedID: result.CustomerNeed.Id,
		Attachment:     attachment,
	}, nil
}

func customerNeedSummary(fields CustomerNeedSummaryFields) CustomerNeedSummary {
	summary := CustomerNeedSummary{
		ID:         fields.Id,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
		Priority:   fields.Priority,
		Body:       stringValue(fields.Body),
		Content:    stringValue(fields.Content),
		URL:        stringValue(fields.Url),
	}
	if fields.Customer != nil {
		summary.CustomerID = fields.Customer.Id
		summary.CustomerName = fields.Customer.Name
	}
	if fields.Issue != nil {
		summary.IssueID = fields.Issue.Id
		summary.Issue = fields.Issue.Identifier
		summary.IssueTitle = fields.Issue.Title
	}
	if fields.Project != nil {
		summary.ProjectID = fields.Project.Id
		summary.ProjectName = fields.Project.Name
	}

	return summary
}
