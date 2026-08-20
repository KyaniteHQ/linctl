package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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

//nolint:lll
type customerNeedsNode = gql.XCustomerNeedsCustomerNeedsCustomerNeedConnectionNodesCustomerNeed

type customerNeedsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListCustomerNeeds returns visible Linear customer needs.
func ListCustomerNeeds(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomerNeedList, error) {
	query := customerNeedsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list customer needs", limit, defaultListPageSize,
		query.page,
		customerNeedsNodeSummary,
	)
	if err != nil {
		return CustomerNeedList{}, err
	}

	return CustomerNeedList{Needs: page.Items, HasNextPage: page.HasNextPage, EndCursor: page.EndCursor}, nil
}

// GetCustomerNeedByID returns one Linear customer need by id.
func GetCustomerNeedByID(ctx context.Context, graphqlClient graphql.Client, id string) (CustomerNeedSummary, error) {
	result, err := gql.XCustomerNeed(ctx, graphqlClient, &id)
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
	result, err := gql.XCustomerNeed_projectAttachment(ctx, graphqlClient, &id)
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

func (query customerNeedsQuery) page(pageSize int, after *string) ([]customerNeedsNode, bool, *string, error) {
	result, err := gql.XCustomerNeeds(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.CustomerNeeds.Nodes,
		result.CustomerNeeds.PageInfo.HasNextPage,
		result.CustomerNeeds.PageInfo.EndCursor,
		nil
}

func customerNeedsNodeSummary(node customerNeedsNode) CustomerNeedSummary {
	return customerNeedSummary(node.CustomerNeedSummaryFields)
}

func customerNeedSummary(fields gql.CustomerNeedSummaryFields) CustomerNeedSummary {
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
