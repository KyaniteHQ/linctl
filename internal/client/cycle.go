package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Cycles []CycleSummary `json:"cycles"`
	Page
}

// SprintReport is a read-only Cycle report with its assigned issues.
type SprintReport struct {
	Cycle  CycleSummary   `json:"cycle"`
	Issues []IssueSummary `json:"issues"`
	Page
}

// CycleIssueList is a page of Issues associated with one Cycle.
type CycleIssueList struct {
	Cycle  CycleSummary   `json:"cycle"`
	Issues []IssueSummary `json:"issues"`
	Page
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

//nolint:lll
type cyclesNode = gql.XCyclesCyclesCycleConnectionNodesCycle

//nolint:lll
type cycleIssuesNode = gql.XCycle_issuesCycleIssuesIssueConnectionNodesIssue

//nolint:lll
type cycleUncompletedIssuesNode = gql.XCycle_uncompletedIssuesUponCloseCycleUncompletedIssuesUponCloseIssueConnectionNodesIssue

type cyclesByTeamQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	teamID        string
}

type cycleScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// cycleScopedParent is the connection parent metadata cycleScopedQuery reads out of
// every page. Linear repeats it per page, so the last page wins.
type cycleScopedParent struct {
	cycle CycleSummary
}

// ListCyclesByTeam returns Cycles scoped to a resolved team.
func ListCyclesByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (CycleList, error) {
	query := cyclesByTeamQuery{ctx: ctx, graphqlClient: graphqlClient, teamID: teamID}
	page, err := listConnection(
		"list cycles", limit, defaultListPageSize,
		query.page,
		cyclesNodeSummary,
	)
	if err != nil {
		return CycleList{}, err
	}

	return CycleList{Cycles: page.Items, Page: page.Page}, nil
}

// GetCycleByID returns a Cycle by Linear id or slug.
func GetCycleByID(ctx context.Context, graphqlClient graphql.Client, id string) (CycleSummary, error) {
	cycle, err := gql.XCycle(ctx, graphqlClient, id)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("get cycle %s: %w", id, err)
	}

	return cycleSummary(cycle.Cycle.CycleSummaryFields), nil
}

// CurrentCycleByTeam returns the active Cycle for a team, using a server-side
// isActive filter so history beyond one page never hides the active Cycle.
func CurrentCycleByTeam(ctx context.Context, graphqlClient graphql.Client, teamID string) (CycleSummary, error) {
	cycle, err := currentActiveCycle(ctx, graphqlClient, teamID)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("current sprint: %w", err)
	}

	return cycle, nil
}

func currentActiveCycle(ctx context.Context, graphqlClient graphql.Client, teamID string) (CycleSummary, error) {
	page, err := gql.XActiveCyclesByTeam(ctx, graphqlClient, teamID)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("list cycles: %w", err)
	}

	switch len(page.Cycles.Nodes) {
	case 0:
		return CycleSummary{}, fmt.Errorf("no active Cycle for team %s", teamID)
	case 1:
		return cycleSummary(page.Cycles.Nodes[0].CycleSummaryFields), nil
	default:
		return CycleSummary{}, fmt.Errorf("multiple active Cycles for team %s", teamID)
	}
}

// GetSprintReport returns one Cycle and its assigned issues.
func GetSprintReport(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (SprintReport, error) {
	report, err := gql.CycleReport(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return SprintReport{}, fmt.Errorf("sprint report %s: %w", id, err)
	}

	issues := mapNodes(report.Cycle.Issues.Nodes, func(
		issue gql.CycleReportCycleIssuesIssueConnectionNodesIssue,
	) IssueSummary {
		return issueSummaryFromFields(issue.IssueSummaryFields)
	})

	return SprintReport{
		Cycle:  cycleSummary(report.Cycle.CycleSummaryFields),
		Issues: issues,
		Page: Page{
			HasNextPage: report.Cycle.Issues.PageInfo.HasNextPage,
			EndCursor:   report.Cycle.Issues.PageInfo.EndCursor,
		},
	}, nil
}

// ListCycleIssues returns Issues assigned to one Cycle.
func ListCycleIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (CycleIssueList, error) {
	query := &cycleScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list cycle issues "+id, limit, defaultListPageSize,
		query.issues,
		cycleIssuesNodeSummary,
	)
	if err != nil {
		return CycleIssueList{}, err
	}

	return CycleIssueList{
		Cycle:  parent.cycle,
		Issues: page.Items,
		Page:   page.Page,
	}, nil
}

// ListCycleUncompletedIssuesUponClose returns Issues left open when one Cycle closed.
func ListCycleUncompletedIssuesUponClose(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (CycleIssueList, error) {
	query := &cycleScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list cycle uncompleted issues "+id, limit, defaultListPageSize,
		query.uncompletedIssues,
		cycleUncompletedIssuesNodeSummary,
	)
	if err != nil {
		return CycleIssueList{}, err
	}

	return CycleIssueList{
		Cycle:  parent.cycle,
		Issues: page.Items,
		Page:   page.Page,
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
		return CycleSummary{}, requiredFieldError("starts at")
	}
	if request.EndsAt == "" {
		return CycleSummary{}, requiredFieldError("ends at")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return CycleSummary{}, err
	}

	return guard.createCycle(ctx, request)
}

func (guard *guardedClient) createCycle(ctx context.Context, request CycleCreateRequest) (CycleSummary, error) {
	created, err := gql.CycleCreate(ctx, guard.graphqlClient, LinearCycleCreateInput{
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
	if !created.CycleCreate.Success || created.CycleCreate.Cycle == nil {
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

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return CycleSummary{}, err
	}

	return guard.updateCycle(ctx, request)
}

func (guard *guardedClient) updateCycle(ctx context.Context, request CycleUpdateRequest) (CycleSummary, error) {
	if err := guard.requireCycle(ctx, request.ID); err != nil {
		return CycleSummary{}, err
	}

	updated, err := gql.CycleUpdate(ctx, guard.graphqlClient, request.ID, LinearCycleUpdateInput{
		Name:        optionalString(request.Name),
		Description: optionalString(request.Description),
		StartsAt:    optionalString(request.StartsAt),
		EndsAt:      optionalString(request.EndsAt),
		CompletedAt: optionalString(request.CompletedAt),
	})
	if err != nil {
		return CycleSummary{}, fmt.Errorf("update cycle %s: %w", request.ID, err)
	}
	if !updated.CycleUpdate.Success || updated.CycleUpdate.Cycle == nil {
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
		return CycleSummary{}, requiredFieldError("cycle id")
	}
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return CycleSummary{}, err
	}

	return guard.archiveCycle(ctx, id)
}

func (guard *guardedClient) archiveCycle(ctx context.Context, id string) (CycleSummary, error) {
	if err := guard.requireCycle(ctx, id); err != nil {
		return CycleSummary{}, err
	}

	archived, err := gql.CycleArchive(ctx, guard.graphqlClient, id)
	if err != nil {
		return CycleSummary{}, fmt.Errorf("archive cycle %s: %w", id, err)
	}
	if !archived.CycleArchive.Success || archived.CycleArchive.Entity == nil {
		return CycleSummary{}, fmt.Errorf("%w: cycleArchive failed", ErrMutationFailed)
	}

	return cycleSummary(archived.CycleArchive.Entity.CycleSummaryFields), nil
}

func validateCycleUpdateRequest(request CycleUpdateRequest) error {
	if request.ID == "" {
		return requiredFieldError("cycle id")
	}
	if request.Name == "" &&
		request.Description == "" &&
		request.StartsAt == "" &&
		request.EndsAt == "" &&
		request.CompletedAt == "" {
		return requiredFieldError("name, description, starts at, ends at, or completed at")
	}

	return nil
}

func (query cyclesByTeamQuery) page(pageSize int, after *string) ([]cyclesNode, bool, *string, error) {
	result, err := gql.XCycles(
		query.ctx, query.graphqlClient, query.teamID, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Cycles.Nodes,
		result.Cycles.PageInfo.HasNextPage,
		result.Cycles.PageInfo.EndCursor,
		nil
}

func (query *cycleScopedQuery) issues(
	pageSize int,
	after *string,
) ([]cycleIssuesNode, cycleScopedParent, bool, *string, error) {
	result, err := gql.XCycle_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, cycleScopedParent{}, false, nil, err
	}

	return result.Cycle.Issues.Nodes,
		cycleScopedParent{cycle: cycleSummary(result.Cycle.CycleSummaryFields)},
		result.Cycle.Issues.PageInfo.HasNextPage,
		result.Cycle.Issues.PageInfo.EndCursor,
		nil
}

func (query *cycleScopedQuery) uncompletedIssues(
	pageSize int,
	after *string,
) ([]cycleUncompletedIssuesNode, cycleScopedParent, bool, *string, error) {
	result, err := gql.XCycle_uncompletedIssuesUponClose(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, cycleScopedParent{}, false, nil, err
	}

	return result.Cycle.UncompletedIssuesUponClose.Nodes,
		cycleScopedParent{cycle: cycleSummary(result.Cycle.CycleSummaryFields)},
		result.Cycle.UncompletedIssuesUponClose.PageInfo.HasNextPage,
		result.Cycle.UncompletedIssuesUponClose.PageInfo.EndCursor,
		nil
}

func cyclesNodeSummary(cycle cyclesNode) CycleSummary {
	return cycleSummary(cycle.CycleSummaryFields)
}

func cycleIssuesNodeSummary(issue cycleIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func cycleUncompletedIssuesNodeSummary(issue cycleUncompletedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func cycleSummary(cycle gql.CycleSummaryFields) CycleSummary {
	name := fmt.Sprintf("Cycle %.0f", cycle.Number)
	if cycle.Name != nil && *cycle.Name != "" {
		name = *cycle.Name
	}
	description := stringValue(cycle.Description)
	completedAt := stringValue(cycle.CompletedAt)

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
