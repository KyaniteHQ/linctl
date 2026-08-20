package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// TriageResponsibilitySummary is the compact triage responsibility model used by read-only commands.
type TriageResponsibilitySummary struct {
	ID               string   `json:"id"`
	Action           string   `json:"action"`
	TeamID           string   `json:"team_id"`
	TeamKey          string   `json:"team_key"`
	TeamName         string   `json:"team_name"`
	TimeScheduleID   string   `json:"time_schedule_id,omitempty"`
	TimeScheduleName string   `json:"time_schedule_name,omitempty"`
	CurrentUserID    string   `json:"current_user_id,omitempty"`
	CurrentUserName  string   `json:"current_user_name,omitempty"`
	ManualUserIDs    []string `json:"manual_user_ids,omitempty"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	ArchivedAt       string   `json:"archived_at,omitempty"`
}

// TriageResponsibilityList is a page of triage responsibilities.
type TriageResponsibilityList struct {
	TriageResponsibilities []TriageResponsibilitySummary `json:"triage_responsibilities"`
	Page
}

// TriageResponsibilityManualSelection is the manual user selection for one triage responsibility.
type TriageResponsibilityManualSelection struct {
	ID      string   `json:"id"`
	UserIDs []string `json:"user_ids"`
}

//nolint:lll
type triageResponsibilitiesNode = gql.XTriageResponsibilitiesTriageResponsibilitiesTriageResponsibilityConnectionNodesTriageResponsibility

type triageResponsibilitiesQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListTriageResponsibilities returns visible triage responsibility configs.
func ListTriageResponsibilities(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (TriageResponsibilityList, error) {
	query := triageResponsibilitiesQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list triage responsibilities", limit, defaultListPageSize,
		query.page,
		triageResponsibilitiesNodeSummary,
	)
	if err != nil {
		return TriageResponsibilityList{}, err
	}

	return TriageResponsibilityList{
		TriageResponsibilities: page.Items,
		Page:                   page.Page,
	}, nil
}

// GetTriageResponsibilityByID returns one triage responsibility by id.
func GetTriageResponsibilityByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (TriageResponsibilitySummary, error) {
	result, err := gql.XTriageResponsibility(ctx, graphqlClient, id)
	if err != nil {
		return TriageResponsibilitySummary{}, fmt.Errorf("get triage responsibility %s: %w", id, err)
	}

	return triageResponsibilitySummary(result.TriageResponsibility.TriageResponsibilitySummaryFields), nil
}

// GetTriageResponsibilityManualSelection returns manual user ids for one triage responsibility.
func GetTriageResponsibilityManualSelection(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (TriageResponsibilityManualSelection, error) {
	result, err := gql.XTriageResponsibility_manualSelection(ctx, graphqlClient, id)
	if err != nil {
		return TriageResponsibilityManualSelection{}, fmt.Errorf(
			"get triage responsibility manual selection %s: %w",
			id,
			err,
		)
	}

	selection := TriageResponsibilityManualSelection{ID: result.TriageResponsibility.Id}
	if result.TriageResponsibility.ManualSelection != nil {
		selection.UserIDs = result.TriageResponsibility.ManualSelection.UserIds
	}
	return selection, nil
}

func (query triageResponsibilitiesQuery) page(
	pageSize int,
	after *string,
) ([]triageResponsibilitiesNode, bool, *string, error) {
	result, err := gql.XTriageResponsibilities(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.TriageResponsibilities.Nodes,
		result.TriageResponsibilities.PageInfo.HasNextPage,
		result.TriageResponsibilities.PageInfo.EndCursor,
		nil
}

func triageResponsibilitiesNodeSummary(node triageResponsibilitiesNode) TriageResponsibilitySummary {
	return triageResponsibilitySummary(node.TriageResponsibilitySummaryFields)
}

func triageResponsibilitySummary(fields gql.TriageResponsibilitySummaryFields) TriageResponsibilitySummary {
	summary := TriageResponsibilitySummary{
		ID:         fields.Id,
		Action:     string(fields.Action),
		TeamID:     fields.Team.Id,
		TeamKey:    fields.Team.Key,
		TeamName:   fields.Team.Name,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
	}
	if fields.TimeSchedule != nil {
		summary.TimeScheduleID = fields.TimeSchedule.Id
		summary.TimeScheduleName = fields.TimeSchedule.Name
	}
	if fields.CurrentUser != nil {
		summary.CurrentUserID = fields.CurrentUser.Id
		summary.CurrentUserName = fields.CurrentUser.DisplayName
	}
	if fields.ManualSelection != nil {
		summary.ManualUserIDs = fields.ManualSelection.UserIds
	}

	return summary
}
