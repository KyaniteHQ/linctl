package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// CycleSummary is the compact Cycle model used by cycle commands.
type CycleSummary struct {
	ID          string  `json:"id"`
	Number      float64 `json:"number"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	StartsAt    string  `json:"starts_at"`
	EndsAt      string  `json:"ends_at"`
	CompletedAt string  `json:"completed_at,omitempty"`
	Progress    float64 `json:"progress"`
	Status      string  `json:"status"`
	TeamID      string  `json:"team_id"`
	TeamKey     string  `json:"team_key"`
	TeamName    string  `json:"team_name"`
}

// CycleList is a page of Cycles scoped to a team.
type CycleList struct {
	Cycles      []CycleSummary `json:"cycles"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// SprintReport is a read-only Cycle report with its assigned issues.
type SprintReport struct {
	Cycle       CycleSummary   `json:"cycle"`
	Issues      []IssueSummary `json:"issues"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// CycleCreateRequest describes a guarded Cycle create in the pinned team.
type CycleCreateRequest struct {
	Name        string
	Description string
	StartsAt    string
	EndsAt      string
	CompletedAt string
}

// CycleUpdateRequest describes a guarded Cycle update.
type CycleUpdateRequest struct {
	ID          string
	Name        string
	Description string
	StartsAt    string
	EndsAt      string
	CompletedAt string
}

// LinearCycleCreateInput is the sparse Linear cycleCreate payload linctl supports.
type LinearCycleCreateInput struct {
	TeamID      string  `json:"teamId"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	StartsAt    string  `json:"startsAt"`
	EndsAt      string  `json:"endsAt"`
	CompletedAt *string `json:"completedAt,omitempty"`
}

// LinearCycleUpdateInput is the sparse Linear cycleUpdate payload linctl supports.
type LinearCycleUpdateInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	StartsAt    *string `json:"startsAt,omitempty"`
	EndsAt      *string `json:"endsAt,omitempty"`
	CompletedAt *string `json:"completedAt,omitempty"`
}

// ListCyclesByTeam returns Cycles scoped to a resolved team.
func ListCyclesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (CycleList, error) {
	cyclePage, err := cycles(ctx, graphqlClient, teamID, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CycleList{}, fmt.Errorf("list cycles: %w", err)
	}

	summaries := make([]CycleSummary, 0, len(cyclePage.Cycles.Nodes))
	for _, cycle := range cyclePage.Cycles.Nodes {
		summaries = append(summaries, cycleSummary(cycle.CycleSummaryFields))
	}

	return CycleList{
		Cycles:      summaries,
		HasNextPage: cyclePage.Cycles.PageInfo.HasNextPage,
		EndCursor:   cyclePage.Cycles.PageInfo.EndCursor,
	}, nil
}

// GetCycleByID returns a Cycle by Linear id or slug.
func GetCycleByID(ctx context.Context, graphqlClient graphql.Client, id string) (CycleSummary, error) {
	cycle, err := cycle(ctx, graphqlClient, id)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("get cycle %s: %w", id, err)
	}

	return cycleSummary(cycle.Cycle.CycleSummaryFields), nil
}

// CurrentCycleByTeam returns the active Cycle for a team.
func CurrentCycleByTeam(ctx context.Context, graphqlClient graphql.Client, teamID string) (CycleSummary, error) {
	cycles, err := ListCyclesByTeam(ctx, graphqlClient, teamID, 50)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("current sprint: %w", err)
	}
	for _, cycle := range cycles.Cycles {
		if cycle.Status == "active" {
			return cycle, nil
		}
	}

	return CycleSummary{}, fmt.Errorf("current sprint: no active Cycle for team %s", teamID)
}

// GetSprintReport returns one Cycle and its assigned issues.
func GetSprintReport(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (SprintReport, error) {
	report, err := CycleReport(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return SprintReport{}, fmt.Errorf("sprint report %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(report.Cycle.Issues.Nodes))
	for _, issue := range report.Cycle.Issues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return SprintReport{
		Cycle:       cycleSummary(report.Cycle.CycleSummaryFields),
		Issues:      issues,
		HasNextPage: report.Cycle.Issues.PageInfo.HasNextPage,
		EndCursor:   report.Cycle.Issues.PageInfo.EndCursor,
	}, nil
}

// CreateCycle creates a Cycle in the pinned team after target comparison.
func CreateCycle(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request CycleCreateRequest,
) (CycleSummary, error) {
	if request.StartsAt == "" {
		return CycleSummary{}, fmt.Errorf("%w: starts at is required", ErrWriteInvalid)
	}
	if request.EndsAt == "" {
		return CycleSummary{}, fmt.Errorf("%w: ends at is required", ErrWriteInvalid)
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return CycleSummary{}, err
	}

	created, err := CycleCreate(ctx, graphqlClient, LinearCycleCreateInput{
		TeamID:      guard.target.Team.ID,
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		StartsAt:    request.StartsAt,
		EndsAt:      request.EndsAt,
		CompletedAt: optionalString(request.CompletedAt),
	})
	if err != nil {
		return CycleSummary{}, fmt.Errorf("create cycle: %w", err)
	}
	if !created.CycleCreate.Success {
		return CycleSummary{}, fmt.Errorf("%w: cycleCreate failed", ErrMutationFailed)
	}

	return cycleSummary(created.CycleCreate.Cycle.CycleSummaryFields), nil
}

// UpdateCycle updates a Cycle after resolving and comparing its team.
func UpdateCycle(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request CycleUpdateRequest,
) (CycleSummary, error) {
	if err := validateCycleUpdateRequest(request); err != nil {
		return CycleSummary{}, err
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return CycleSummary{}, err
	}
	if err := guard.requireCycle(ctx, graphqlClient, request.ID); err != nil {
		return CycleSummary{}, err
	}

	updated, err := CycleUpdate(ctx, graphqlClient, request.ID, LinearCycleUpdateInput{
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		StartsAt:    optionalString(request.StartsAt),
		EndsAt:      optionalString(request.EndsAt),
		CompletedAt: optionalString(request.CompletedAt),
	})
	if err != nil {
		return CycleSummary{}, fmt.Errorf("update cycle %s: %w", request.ID, err)
	}
	if !updated.CycleUpdate.Success {
		return CycleSummary{}, fmt.Errorf("%w: cycleUpdate failed", ErrMutationFailed)
	}

	return cycleSummary(updated.CycleUpdate.Cycle.CycleSummaryFields), nil
}

// ArchiveCycle archives a Cycle after resolving and comparing its team.
func ArchiveCycle(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	id string,
) (CycleSummary, error) {
	if id == "" {
		return CycleSummary{}, fmt.Errorf("%w: cycle id is required", ErrWriteInvalid)
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return CycleSummary{}, err
	}
	if err := guard.requireCycle(ctx, graphqlClient, id); err != nil {
		return CycleSummary{}, err
	}

	archived, err := CycleArchive(ctx, graphqlClient, id)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("archive cycle %s: %w", id, err)
	}
	if !archived.CycleArchive.Success {
		return CycleSummary{}, fmt.Errorf("%w: cycleArchive failed", ErrMutationFailed)
	}

	return cycleSummary(archived.CycleArchive.Entity.CycleSummaryFields), nil
}

func validateCycleUpdateRequest(request CycleUpdateRequest) error {
	if request.ID == "" {
		return fmt.Errorf("%w: cycle id is required", ErrWriteInvalid)
	}
	if request.Name == "" &&
		request.Description == "" &&
		request.StartsAt == "" &&
		request.EndsAt == "" &&
		request.CompletedAt == "" {
		return fmt.Errorf("%w: name, description, starts at, ends at, or completed at is required", ErrWriteInvalid)
	}

	return nil
}

func cycleSummary(cycle CycleSummaryFields) CycleSummary {
	name := fmt.Sprintf("Cycle %.0f", cycle.Number)
	if cycle.Name != nil && *cycle.Name != "" {
		name = *cycle.Name
	}
	description := ""
	if cycle.Description != nil {
		description = *cycle.Description
	}
	completedAt := ""
	if cycle.CompletedAt != nil {
		completedAt = *cycle.CompletedAt
	}

	return CycleSummary{
		ID:          cycle.Id,
		Number:      cycle.Number,
		Name:        name,
		Description: description,
		StartsAt:    cycle.StartsAt,
		EndsAt:      cycle.EndsAt,
		CompletedAt: completedAt,
		Progress:    cycle.Progress,
		Status:      cycleStatus(cycle.StartsAt, cycle.EndsAt, completedAt),
		TeamID:      cycle.Team.Id,
		TeamKey:     cycle.Team.Key,
		TeamName:    cycle.Team.Name,
	}
}

func cycleStatus(startsAt string, endsAt string, completedAt string) string {
	if completedAt != "" {
		return "completed"
	}
	now := time.Now().UTC()
	start, startErr := time.Parse(time.RFC3339, startsAt)
	end, endErr := time.Parse(time.RFC3339, endsAt)
	if startErr != nil || endErr != nil {
		return "unknown"
	}
	if now.Before(start) {
		return "future"
	}
	if now.After(end) {
		return "past"
	}
	return "active"
}
