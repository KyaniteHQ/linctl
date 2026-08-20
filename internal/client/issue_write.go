package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

const (
	defaultIssueCreateBatchConcurrency = 4
	maxIssueCreateBatchConcurrency     = 8
)

// ErrViewerNotAssignable marks a start request whose authenticated viewer
// cannot own or receive delegated issue work.
var ErrViewerNotAssignable = errors.New("viewer is not assignable")

// IssueCreateRequest describes a guarded issue create.
type IssueCreateRequest struct {
	Title              string
	Description        string
	StateSelector      string
	StateType          string
	Priority           string
	AssigneeID         string
	LabelIDs           []string
	DueDate            string
	Estimate           *int
	ParentID           string
	ProjectMilestoneID string
}

// IssueUpdateRequest describes a guarded issue update.
type IssueUpdateRequest struct {
	ID                 string
	Title              string
	Description        string
	Append             string
	StateSelector      string
	StateType          string
	Priority           string
	AssigneeID         string
	LabelIDs           []string
	DueDate            string
	ClearDueDate       bool
	Estimate           *int
	ClearEstimate      bool
	ProjectMilestoneID string
	ClearMilestone     bool
}

// IssueCommentRequest describes a guarded issue comment.
type IssueCommentRequest struct {
	ID       string
	Body     string
	ParentID string
}

// IssueCommentResult is the created comment plus its issue.
type IssueCommentResult struct {
	ID    string       `json:"id"`
	Body  string       `json:"body"`
	URL   string       `json:"url"`
	Issue IssueSummary `json:"issue"`
}

// CreateIssue creates an issue after resolving and comparing the pinned write target.
func CreateIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueCreateRequest,
) (IssueSummary, error) {
	if request.Title == "" {
		return IssueSummary{}, requiredFieldError("title")
	}
	if err := validateDueDate(request.DueDate); err != nil {
		return IssueSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.createIssue(ctx, request)
}

// createIssue runs one issue create through an already-resolved guarded client.
// CreateIssue and CreateIssues both call it, so a single-row create and a batch
// row go through identical input assembly and guard validation.
func (guard *guardedClient) createIssue(
	ctx context.Context,
	request IssueCreateRequest,
) (IssueSummary, error) {
	if err := guard.validateEstimate(ctx, guard.target.Team.ID, request.Estimate); err != nil {
		return IssueSummary{}, err
	}
	if request.ParentID != "" {
		if _, err := guard.requireIssue(ctx, request.ParentID); err != nil {
			return IssueSummary{}, err
		}
	}
	if err := guard.requireAttachableLabels(ctx, request.LabelIDs); err != nil {
		return IssueSummary{}, err
	}
	input := LinearIssueCreateInput{
		Title:              stringPtr(request.Title),
		Description:        optionalString(request.Description),
		TeamID:             guard.target.Team.ID,
		AssigneeID:         optionalString(request.AssigneeID),
		LabelIDs:           request.LabelIDs,
		DueDate:            optionalString(request.DueDate),
		Estimate:           request.Estimate,
		ParentID:           optionalString(request.ParentID),
		ProjectMilestoneID: optionalString(request.ProjectMilestoneID),
	}
	if guard.target.Project != nil {
		input.ProjectID = stringPtr(guard.target.Project.ID)
	}
	if err := guard.requireCreateMilestone(ctx, request.ProjectMilestoneID); err != nil {
		return IssueSummary{}, err
	}
	stateID, stateSet, err := guard.resolveRequestStateID(ctx, guard.target.Team.ID, issueCreateStateSelector(request))
	if err != nil {
		return IssueSummary{}, err
	}
	if stateSet {
		input.StateID = stringPtr(stateID)
	}
	priority, err := parsePriority(request.Priority)
	if err != nil {
		return IssueSummary{}, err
	}
	input.Priority = priority
	created, err := gql.IssueCreate(ctx, guard.graphqlClient, input)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("create issue: %w", err)
	}
	if !created.IssueCreate.Success || created.IssueCreate.Issue == nil {
		return IssueSummary{}, fmt.Errorf("%w: issueCreate returned no issue", ErrMutationFailed)
	}
	if stateSet {
		return guard.finishStateWrite(ctx, created.IssueCreate.Issue.Id, stateID, nil)
	}

	return issueSummaryFromFields(created.IssueCreate.Issue.IssueSummaryFields), nil
}

// IssueCreateOutcome is one row's result from a CreateIssues batch call,
// indexed to match its position in the request slice so a caller can report
// row order and per-row failures without losing that order under concurrency.
type IssueCreateOutcome struct {
	Index int
	Issue IssueSummary
	Err   error
}

// CreateIssues creates each request against ONE resolved target, preserving
// per-request guard validation (estimate, parent, milestone) for every row.
// concurrency is clamped to [defaultIssueCreateBatchConcurrency,
// maxIssueCreateBatchConcurrency]; a value <= 0 uses the default.
func CreateIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	requests []IssueCreateRequest,
	concurrency int,
) ([]IssueCreateOutcome, error) {
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return nil, err
	}

	outcomes := make([]IssueCreateOutcome, len(requests))
	tokens := make(chan struct{}, clampIssueCreateBatchConcurrency(concurrency))
	var waitGroup sync.WaitGroup
	for index, request := range requests {
		waitGroup.Add(1)
		tokens <- struct{}{}
		go func(index int, request IssueCreateRequest) {
			defer waitGroup.Done()
			defer func() { <-tokens }()
			issue, rowErr := guard.createIssueForBatchRow(ctx, request)
			outcomes[index] = IssueCreateOutcome{Index: index, Issue: issue, Err: rowErr}
		}(index, request)
	}
	waitGroup.Wait()

	return outcomes, nil
}

// createIssueForBatchRow applies CreateIssue's pre-guard input validation
// (title required, due date well-formed) to one batch row before running it
// through the shared guard-scoped create, so a malformed row fails as a row
// error instead of aborting the whole batch.
func (guard *guardedClient) createIssueForBatchRow(
	ctx context.Context,
	request IssueCreateRequest,
) (IssueSummary, error) {
	if request.Title == "" {
		return IssueSummary{}, requiredFieldError("title")
	}
	if err := validateDueDate(request.DueDate); err != nil {
		return IssueSummary{}, err
	}

	return guard.createIssue(ctx, request)
}

func clampIssueCreateBatchConcurrency(concurrency int) int {
	if concurrency <= 0 {
		return defaultIssueCreateBatchConcurrency
	}
	if concurrency > maxIssueCreateBatchConcurrency {
		return maxIssueCreateBatchConcurrency
	}

	return concurrency
}

// requireCreateMilestone verifies a requested ProjectMilestone assignment on
// create: it requires a pinned project and a milestone inside the pinned target.
func (guard *guardedClient) requireCreateMilestone(
	ctx context.Context,
	projectMilestoneID string,
) error {
	if projectMilestoneID == "" {
		return nil
	}
	if guard.target.Project == nil {
		return fmt.Errorf("%w: --milestone requires a pinned project", ErrWriteInvalid)
	}

	return guard.requireProjectMilestone(ctx, projectMilestoneID)
}

// requireUpdateMilestone validates a milestone assignment on update the same way
// create does; a clear needs no validation because it removes the assignment
// rather than pointing at one that must exist in the pinned project.
func (guard *guardedClient) requireUpdateMilestone(
	ctx context.Context,
	request IssueUpdateRequest,
) error {
	if request.ClearMilestone {
		return nil
	}

	return guard.requireCreateMilestone(ctx, request.ProjectMilestoneID)
}

// UpdateIssue updates an issue after resolving and comparing the pinned write target.
func UpdateIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueUpdateRequest,
) (IssueSummary, error) {
	if err := validateIssueUpdateRequest(request); err != nil {
		return IssueSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.updateIssue(ctx, request)
}

func (guard *guardedClient) updateIssue(ctx context.Context, request IssueUpdateRequest) (IssueSummary, error) {
	issue, err := guard.requireIssueDetail(ctx, request.ID)
	if err != nil {
		return IssueSummary{}, err
	}
	if err = guard.validateEstimate(ctx, issue.Summary.TeamID, request.Estimate); err != nil {
		return IssueSummary{}, err
	}
	if err = guard.requireAttachableLabels(ctx, request.LabelIDs); err != nil {
		return IssueSummary{}, err
	}
	if err = guard.requireUpdateMilestone(ctx, request); err != nil {
		return IssueSummary{}, err
	}
	description := request.Description
	if request.Append != "" {
		description = appendIssueDescription(issue.Description, request.Append)
	}
	stateID, stateSet, err := guard.resolveRequestStateID(ctx, issue.Summary.TeamID, issueUpdateStateSelector(request))
	if err != nil {
		return IssueSummary{}, err
	}
	if stateSet && issue.Summary.StateID == stateID && issueUpdateOnlyState(request) {
		return issue.Summary, nil
	}

	return guard.applyIssueUpdate(ctx, request, description, stateID, stateSet)
}

func (guard *guardedClient) applyIssueUpdate(
	ctx context.Context,
	request IssueUpdateRequest,
	description string,
	stateID string,
	stateSet bool,
) (IssueSummary, error) {
	var resolvedStateID *string
	if stateSet {
		resolvedStateID = stringPtr(stateID)
	}
	updateInput, err := guard.buildIssueUpdateInput(ctx, request, description, resolvedStateID)
	if err != nil {
		return IssueSummary{}, err
	}
	updated, err := gql.IssueUpdate(ctx, guard.graphqlClient, request.ID, updateInput)
	writeErr := err
	if writeErr != nil {
		writeErr = fmt.Errorf("update issue %s: %w", request.ID, writeErr)
	} else if !updated.IssueUpdate.Success || updated.IssueUpdate.Issue == nil {
		writeErr = fmt.Errorf("%w: issueUpdate returned no issue", ErrMutationFailed)
	}
	if stateSet {
		return guard.finishStateWrite(ctx, request.ID, stateID, writeErr)
	}
	if writeErr != nil {
		return IssueSummary{}, writeErr
	}

	return issueSummaryFromFields(updated.IssueUpdate.Issue.IssueSummaryFields), nil
}

func validateIssueUpdateRequest(request IssueUpdateRequest) error {
	if request.ID == "" {
		return requiredFieldError("issue id")
	}
	if issueUpdateHasNoFields(request) {
		return requiredFieldError("title, description, state, priority, assignee, label, due date, or estimate")
	}
	if request.Description != "" && request.Append != "" {
		return fmt.Errorf("%w: description and append are mutually exclusive", ErrWriteInvalid)
	}
	if request.DueDate != "" && request.ClearDueDate {
		return fmt.Errorf("%w: due-date and clear-due-date are mutually exclusive", ErrWriteInvalid)
	}
	if request.Estimate != nil && request.ClearEstimate {
		return fmt.Errorf("%w: estimate and clear-estimate are mutually exclusive", ErrWriteInvalid)
	}
	if request.ProjectMilestoneID != "" && request.ClearMilestone {
		return fmt.Errorf("%w: milestone and clear-milestone are mutually exclusive", ErrWriteInvalid)
	}

	return validateDueDate(request.DueDate)
}

func issueUpdateHasNoFields(request IssueUpdateRequest) bool {
	return request.Title == "" && request.Description == "" && request.Append == "" &&
		issueUpdateStateSelector(request) == "" && request.Priority == "" && request.AssigneeID == "" &&
		len(request.LabelIDs) == 0 && request.DueDate == "" && !request.ClearDueDate &&
		request.Estimate == nil && !request.ClearEstimate &&
		request.ProjectMilestoneID == "" && !request.ClearMilestone
}

func issueUpdateOnlyState(request IssueUpdateRequest) bool {
	return request.Title == "" && request.Description == "" && request.Append == "" &&
		request.Priority == "" && request.AssigneeID == "" &&
		len(request.LabelIDs) == 0 && request.DueDate == "" && !request.ClearDueDate &&
		request.Estimate == nil && !request.ClearEstimate &&
		request.ProjectMilestoneID == "" && !request.ClearMilestone
}

func (guard *guardedClient) buildIssueUpdateInput(
	_ context.Context,
	request IssueUpdateRequest,
	description string,
	stateID *string,
) (LinearIssueUpdateInput, error) {
	input := LinearIssueUpdateInput{
		Title:              optionalString(request.Title),
		Description:        optionalString(description),
		AssigneeID:         optionalString(request.AssigneeID),
		LabelIDs:           request.LabelIDs,
		DueDate:            dueDateUpdateJSON(request),
		Estimate:           estimateUpdateJSON(request),
		ProjectMilestoneID: projectMilestoneUpdateJSON(request),
		StateID:            stateID,
	}
	priority, err := parsePriority(request.Priority)
	if err != nil {
		return LinearIssueUpdateInput{}, err
	}
	input.Priority = priority

	return input, nil
}

func appendIssueDescription(description string, note string) string {
	if strings.TrimSpace(description) == "" {
		return note
	}

	return strings.TrimRight(description, "\n") + "\n\n" + note
}

// StartIssue assigns a human viewer, delegates to an app viewer, and moves the
// issue to the team's started workflow state.
func StartIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (IssueSummary, error) {
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.startIssue(ctx, issueID)
}

func (guard *guardedClient) startIssue(ctx context.Context, issueID string) (IssueSummary, error) {
	if guard.target.Viewer.App && !guard.target.Viewer.Assignable {
		return IssueSummary{}, fmt.Errorf(
			"%w: the authenticated app requires the app:assignable scope",
			ErrViewerNotAssignable,
		)
	}

	issue, err := guard.requireIssue(ctx, issueID)
	if err != nil {
		return IssueSummary{}, err
	}
	stateID, err := guard.resolveStateID(ctx, issue.TeamID, "started")
	if err != nil {
		return IssueSummary{}, err
	}

	input := LinearIssueUpdateInput{StateID: stringPtr(stateID)}
	if guard.target.Viewer.App {
		input.DelegateID = stringPtr(guard.target.Viewer.ID)
	} else {
		input.AssigneeID = stringPtr(guard.target.Viewer.ID)
	}

	started, err := gql.IssueUpdate(ctx, guard.graphqlClient, issueID, input)
	writeErr := err
	if writeErr != nil {
		writeErr = fmt.Errorf("start issue %s: %w", issueID, writeErr)
	} else if !started.IssueUpdate.Success || started.IssueUpdate.Issue == nil {
		writeErr = fmt.Errorf("%w: issue start returned no issue", ErrMutationFailed)
	}

	return guard.finishStateWrite(ctx, issueID, stateID, writeErr)
}

// CommentOnIssue adds a comment after resolving and comparing the pinned write target.
func CommentOnIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueCommentRequest,
) (IssueCommentResult, error) {
	if request.ID == "" {
		return IssueCommentResult{}, requiredFieldError("issue id")
	}
	if request.Body == "" {
		return IssueCommentResult{}, requiredFieldError("body")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueCommentResult{}, err
	}

	return guard.commentOnIssue(ctx, request)
}

func (guard *guardedClient) commentOnIssue(
	ctx context.Context,
	request IssueCommentRequest,
) (IssueCommentResult, error) {
	issue, err := guard.requireIssue(ctx, request.ID)
	if err != nil {
		return IssueCommentResult{}, err
	}
	if err := guard.requireParentCommentOnIssue(ctx, request.ParentID, issue); err != nil {
		return IssueCommentResult{}, err
	}

	comment, err := gql.IssueCommentCreate(ctx, guard.graphqlClient, LinearCommentCreateInput{
		Body:     stringPtr(request.Body),
		IssueID:  stringPtr(request.ID),
		ParentID: optionalString(request.ParentID),
	})
	if err != nil {
		return IssueCommentResult{}, fmt.Errorf("comment on issue %s: %w", request.ID, err)
	}
	if !comment.CommentCreate.Success || comment.CommentCreate.Comment.Issue == nil {
		return IssueCommentResult{}, fmt.Errorf("%w: commentCreate returned no issue", ErrMutationFailed)
	}

	return IssueCommentResult{
		ID:    comment.CommentCreate.Comment.Id,
		Body:  comment.CommentCreate.Comment.Body,
		URL:   comment.CommentCreate.Comment.Url,
		Issue: issueSummaryFromFields(comment.CommentCreate.Comment.Issue.IssueSummaryFields),
	}, nil
}

// requireParentCommentOnIssue verifies a reply's parent comment is attached to
// the guarded issue itself, so a reply cannot land in a different thread than
// the issue the target comparison approved.
func (guard *guardedClient) requireParentCommentOnIssue(
	ctx context.Context,
	parentID string,
	issue IssueSummary,
) error {
	if parentID == "" {
		return nil
	}
	parent, err := GetCommentByID(ctx, guard.graphqlClient, parentID)
	if err != nil {
		return err
	}
	if parent.IssueID != issue.ID {
		return fmt.Errorf(
			"%w: parent comment %s is attached to issue id %q, not issue %s",
			ErrWriteInvalid, parentID, parent.IssueID, issue.Identifier,
		)
	}

	return nil
}

// CloseIssue moves an issue to the team's completed workflow state after target comparison.
func CloseIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (IssueSummary, error) {
	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.closeIssue(ctx, issueID)
}

func (guard *guardedClient) closeIssue(ctx context.Context, issueID string) (IssueSummary, error) {
	issue, err := guard.requireIssue(ctx, issueID)
	if err != nil {
		return IssueSummary{}, err
	}
	stateID, err := guard.resolveStateID(ctx, issue.TeamID, "completed")
	if err != nil {
		return IssueSummary{}, err
	}

	closed, err := gql.IssueClose(ctx, guard.graphqlClient, issueID, LinearIssueUpdateInput{
		StateID: stringPtr(stateID),
	})
	writeErr := err
	if writeErr != nil {
		writeErr = fmt.Errorf("close issue %s: %w", issueID, writeErr)
	} else if !closed.IssueUpdate.Success || closed.IssueUpdate.Issue == nil {
		writeErr = fmt.Errorf("%w: issue close returned no issue", ErrMutationFailed)
	}

	return guard.finishStateWrite(ctx, issueID, stateID, writeErr)
}

func parsePriority(raw string) (*int, error) {
	if raw == "" {
		return nil, nil //nolint:nilnil // nil *int is the intentional "no priority" signal
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: priority must be a number (0-4), got %q", ErrWriteInvalid, raw)
	}
	if value < 0 || value > 4 {
		return nil, fmt.Errorf("%w: priority must be 0-4, got %d", ErrWriteInvalid, value)
	}

	return &value, nil
}

// validateDueDate ensures a non-empty due date is a calendar date (YYYY-MM-DD).
func validateDueDate(raw string) error {
	if raw == "" {
		return nil
	}
	if _, err := time.Parse("2006-01-02", raw); err != nil {
		return fmt.Errorf("%w: due date must be YYYY-MM-DD, got %q", ErrWriteInvalid, raw)
	}

	return nil
}

// dueDateUpdateJSON renders the issueUpdate dueDate field: an explicit null to
// clear it, a quoted date to set it, or nil to leave it untouched.
func dueDateUpdateJSON(request IssueUpdateRequest) json.RawMessage {
	if request.ClearDueDate {
		return json.RawMessage("null")
	}
	if request.DueDate == "" {
		return nil
	}

	return json.RawMessage(strconv.Quote(request.DueDate))
}

// teamEstimateConfig is the team-level estimate configuration linctl can
// verify before a mutation.
type teamEstimateConfig struct {
	estimationType string
	allowZero      bool
}

// estimateConfigForTeam reads a team's estimate configuration, memoized per
// guard so batch writes read each team configuration once.
func (guard *guardedClient) estimateConfigForTeam(
	ctx context.Context,
	teamID string,
) (teamEstimateConfig, error) {
	guard.estimateConfigs.mu.Lock()
	defer guard.estimateConfigs.mu.Unlock()
	if config, ok := guard.estimateConfigs.configs[teamID]; ok {
		return config, nil
	}

	result, err := gql.XTeamEstimateConfig(ctx, guard.graphqlClient, teamID)
	if err != nil {
		return teamEstimateConfig{}, fmt.Errorf("read team estimate config for team_id=%s: %w", teamID, err)
	}
	config := teamEstimateConfig{
		estimationType: result.Team.IssueEstimationType,
		allowZero:      result.Team.IssueEstimationAllowZero,
	}
	guard.estimateConfigs.configs[teamID] = config

	return config, nil
}

// validateEstimate fails closed before a mutation when the team cannot accept
// the requested estimate. It performs a free read of the team's estimate
// configuration and enforces the team-level constraints linctl can verify from
// that configuration: estimates must be enabled, and a zero estimate is only
// accepted when the team allows it. Linear remains authoritative for the exact
// point scale of each estimation type.
func (guard *guardedClient) validateEstimate(ctx context.Context, teamID string, estimate *int) error {
	if estimate == nil {
		return nil
	}
	if *estimate < 0 {
		return fmt.Errorf("%w: estimate must not be negative, got %d", ErrWriteInvalid, *estimate)
	}
	config, err := guard.estimateConfigForTeam(ctx, teamID)
	if err != nil {
		return err
	}
	if config.estimationType == "" || config.estimationType == "notUsed" {
		return fmt.Errorf("%w: team team_id=%s has estimates disabled", ErrWriteInvalid, teamID)
	}
	if *estimate == 0 && !config.allowZero {
		return fmt.Errorf("%w: team team_id=%s does not allow a zero estimate", ErrWriteInvalid, teamID)
	}

	return nil
}

// estimateUpdateJSON renders the issueUpdate estimate field: an explicit null to
// clear it, the integer to set it, or nil to leave it untouched.
func estimateUpdateJSON(request IssueUpdateRequest) json.RawMessage {
	if request.ClearEstimate {
		return json.RawMessage("null")
	}
	if request.Estimate == nil {
		return nil
	}

	return json.RawMessage(strconv.Itoa(*request.Estimate))
}

// projectMilestoneUpdateJSON renders the issueUpdate projectMilestoneId field: an
// explicit null to clear it, a quoted id to set it, or nil to leave it untouched.
func projectMilestoneUpdateJSON(request IssueUpdateRequest) json.RawMessage {
	if request.ClearMilestone {
		return json.RawMessage("null")
	}
	if request.ProjectMilestoneID == "" {
		return nil
	}

	return json.RawMessage(strconv.Quote(request.ProjectMilestoneID))
}
