package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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

// ListTimeSchedules returns visible Linear time schedules.
func ListTimeSchedules(ctx context.Context, graphqlClient graphql.Client, limit int) (TimeScheduleList, error) {
	result, err := timeSchedules(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return TimeScheduleList{}, fmt.Errorf("list time schedules: %w", err)
	}

	summaries := make([]TimeScheduleSummary, 0, len(result.TimeSchedules.Nodes))
	for _, node := range result.TimeSchedules.Nodes {
		summaries = append(summaries, timeScheduleSummary(node.TimeScheduleSummaryFields))
	}

	return TimeScheduleList{
		TimeSchedules: summaries,
		HasNextPage:   result.TimeSchedules.PageInfo.HasNextPage,
		EndCursor:     result.TimeSchedules.PageInfo.EndCursor,
	}, nil
}

// GetTimeScheduleByID returns one Linear time schedule by id.
func GetTimeScheduleByID(ctx context.Context, graphqlClient graphql.Client, id string) (TimeScheduleSummary, error) {
	result, err := timeSchedule(ctx, graphqlClient, id)
	if err != nil {
		return TimeScheduleSummary{}, fmt.Errorf("get time schedule %s: %w", id, err)
	}

	return timeScheduleSummary(result.TimeSchedule.TimeScheduleSummaryFields), nil
}

func timeScheduleSummary(fields TimeScheduleSummaryFields) TimeScheduleSummary {
	entries := make([]TimeScheduleEntrySummary, 0, len(fields.Entries))
	for _, entry := range fields.Entries {
		entries = append(entries, TimeScheduleEntrySummary{
			StartsAt:  entry.StartsAt,
			EndsAt:    entry.EndsAt,
			UserID:    stringValue(entry.UserId),
			UserEmail: stringValue(entry.UserEmail),
		})
	}

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
