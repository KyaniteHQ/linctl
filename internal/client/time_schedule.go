package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// TimeScheduleEntrySummary is one compact entry in a Linear time schedule.
type TimeScheduleEntrySummary struct {
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
	UserID    string `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

// TimeScheduleSummary is the compact time schedule model used by read-only commands.
type TimeScheduleSummary struct {
	ID            string                     `json:"id"`
	Name          string                     `json:"name"`
	CreatedAt     string                     `json:"created_at"`
	UpdatedAt     string                     `json:"updated_at"`
	ArchivedAt    string                     `json:"archived_at,omitempty"`
	ExternalID    string                     `json:"external_id,omitempty"`
	ExternalURL   string                     `json:"external_url,omitempty"`
	IntegrationID string                     `json:"integration_id,omitempty"`
	EntryCount    int                        `json:"entry_count"`
	Entries       []TimeScheduleEntrySummary `json:"entries"`
}

// TimeScheduleList is a page of Linear time schedules.
type TimeScheduleList struct {
	TimeSchedules []TimeScheduleSummary `json:"time_schedules"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

//nolint:lll
type timeSchedulesNode = gql.XTimeSchedulesTimeSchedulesTimeScheduleConnectionNodesTimeSchedule

type timeSchedulesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListTimeSchedules returns visible Linear time schedules.
func ListTimeSchedules(ctx context.Context, graphqlClient graphql.Client, limit int) (TimeScheduleList, error) {
	query := timeSchedulesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list time schedules", limit, defaultListPageSize,
		query.page,
		timeSchedulesNodeSummary,
	)
	if err != nil {
		return TimeScheduleList{}, err
	}

	return TimeScheduleList{
		TimeSchedules: page.Items,
		HasNextPage:   page.HasNextPage,
		EndCursor:     page.EndCursor,
	}, nil
}

// GetTimeScheduleByID returns one Linear time schedule by id.
func GetTimeScheduleByID(ctx context.Context, graphqlClient graphql.Client, id string) (TimeScheduleSummary, error) {
	result, err := gql.XTimeSchedule(ctx, graphqlClient, id)
	if err != nil {
		return TimeScheduleSummary{}, fmt.Errorf("get time schedule %s: %w", id, err)
	}

	return timeScheduleSummary(result.TimeSchedule.TimeScheduleSummaryFields), nil
}

func (query timeSchedulesQuery) page(pageSize int, after *string) ([]timeSchedulesNode, bool, *string, error) {
	result, err := gql.XTimeSchedules(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.TimeSchedules.Nodes,
		result.TimeSchedules.PageInfo.HasNextPage,
		result.TimeSchedules.PageInfo.EndCursor,
		nil
}

func timeSchedulesNodeSummary(node timeSchedulesNode) TimeScheduleSummary {
	return timeScheduleSummary(node.TimeScheduleSummaryFields)
}

func timeScheduleSummary(fields gql.TimeScheduleSummaryFields) TimeScheduleSummary {
	entries := mapNodes(fields.Entries, func(
		entry gql.TimeScheduleSummaryFieldsEntriesTimeScheduleEntry,
	) TimeScheduleEntrySummary {
		return TimeScheduleEntrySummary{
			StartsAt:  entry.StartsAt,
			EndsAt:    entry.EndsAt,
			UserID:    stringValue(entry.UserId),
			UserEmail: stringValue(entry.UserEmail),
		}
	})

	summary := TimeScheduleSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		CreatedAt:   fields.CreatedAt,
		UpdatedAt:   fields.UpdatedAt,
		ArchivedAt:  stringValue(fields.ArchivedAt),
		ExternalID:  stringValue(fields.ExternalId),
		ExternalURL: stringValue(fields.ExternalUrl),
		EntryCount:  len(entries),
		Entries:     entries,
	}
	if fields.Integration != nil {
		summary.IntegrationID = fields.Integration.Id
	}

	return summary
}
