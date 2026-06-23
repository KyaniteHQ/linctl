package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// IssueCreateRequest describes a guarded issue create.
type IssueCreateRequest struct {
	Title       string
	Description string
	StateType   string
	Priority    string
}

// IssueUpdateRequest describes a guarded issue update.
type IssueUpdateRequest struct {
	ID          string
	Title       string
	Description string
	Append      string
	StateType   string
	Priority    string
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

// LinearIssueCreateInput is the sparse Linear issueCreate payload linctl supports.
type LinearIssueCreateInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	TeamID      string  `json:"teamId"`
	ProjectID   *string `json:"projectId,omitempty"`
	StateID     *string `json:"stateId,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

// LinearIssueUpdateInput is the sparse Linear issueUpdate payload linctl supports.
type LinearIssueUpdateInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	AssigneeID  *string `json:"assigneeId,omitempty"`
	StateID     *string `json:"stateId,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

// LinearCommentCreateInput is the sparse Linear commentCreate payload linctl supports.
type LinearCommentCreateInput struct {
	Body     *string `json:"body,omitempty"`
	IssueID  *string `json:"issueId,omitempty"`
	ParentID *string `json:"parentId,omitempty"`
}

// CreateIssue creates an issue after resolving and comparing the pinned write target.
func CreateIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueCreateRequest,
) (IssueSummary, error) {
	if request.Title == "" {
		return IssueSummary{}, fmt.Errorf("%w: title is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueSummary, error) {
		input := LinearIssueCreateInput{
			Title:       stringPtr(request.Title),
			Description: optionalString(request.Description),
			TeamID:      guard.target.Team.ID,
		}
		if guard.target.Project != nil {
			input.ProjectID = stringPtr(guard.target.Project.ID)
		}
		if request.StateType != "" {
			stateID, stateErr := firstStateIDOfType(ctx, graphqlClient, guard.target.Team.ID, request.StateType)
			if stateErr != nil {
				return IssueSummary{}, stateErr
			}
			input.StateID = stringPtr(stateID)
		}
		priority, err := parsePriority(request.Priority)
		if err != nil {
			return IssueSummary{}, err
		}
		input.Priority = priority
		created, err := IssueCreate(ctx, graphqlClient, input)
		if err != nil {
			return IssueSummary{}, fmt.Errorf("create issue: %w", err)
		}
		if !created.IssueCreate.Success || created.IssueCreate.Issue == nil {
			return IssueSummary{}, fmt.Errorf("%w: issueCreate returned no issue", ErrMutationFailed)
		}

		return issueSummaryFromFields(created.IssueCreate.Issue.IssueSummaryFields), nil
	})
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueSummary, error) {
		issue, err := guard.requireIssueDetail(ctx, graphqlClient, request.ID)
		if err != nil {
			return IssueSummary{}, err
		}
		description := request.Description
		if request.Append != "" {
			description = appendIssueDescription(issue.Description, request.Append)
		}

		updateInput, err := buildIssueUpdateInput(ctx, graphqlClient, request, issue.Summary.TeamID, description)
		if err != nil {
			return IssueSummary{}, err
		}
		updated, err := IssueUpdate(ctx, graphqlClient, request.ID, updateInput)
		if err != nil {
			return IssueSummary{}, fmt.Errorf("update issue %s: %w", request.ID, err)
		}
		if !updated.IssueUpdate.Success || updated.IssueUpdate.Issue == nil {
			return IssueSummary{}, fmt.Errorf("%w: issueUpdate returned no issue", ErrMutationFailed)
		}

		return issueSummaryFromFields(updated.IssueUpdate.Issue.IssueSummaryFields), nil
	})
}

func validateIssueUpdateRequest(request IssueUpdateRequest) error {
	if request.ID == "" {
		return fmt.Errorf("%w: issue id is required", ErrWriteInvalid)
	}
	if request.Title == "" && request.Description == "" && request.Append == "" &&
		request.StateType == "" && request.Priority == "" {
		return fmt.Errorf("%w: title, description, state, or priority is required", ErrWriteInvalid)
	}
	if request.Description != "" && request.Append != "" {
		return fmt.Errorf("%w: description and append are mutually exclusive", ErrWriteInvalid)
	}

	return nil
}

func buildIssueUpdateInput(
	ctx context.Context,
	graphqlClient graphql.Client,
	request IssueUpdateRequest,
	teamID string,
	description string,
) (LinearIssueUpdateInput, error) {
	input := LinearIssueUpdateInput{
		Title:       optionalString(request.Title),
		Description: optionalString(description),
	}
	if request.StateType != "" {
		stateID, err := firstStateIDOfType(ctx, graphqlClient, teamID, request.StateType)
		if err != nil {
			return LinearIssueUpdateInput{}, err
		}
		input.StateID = stringPtr(stateID)
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

// StartIssue assigns an issue to the viewer and moves it to the team's started workflow state.
func StartIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (IssueSummary, error) {
	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueSummary, error) {
		issue, err := guard.requireIssue(ctx, graphqlClient, issueID)
		if err != nil {
			return IssueSummary{}, err
		}
		stateID, err := firstStartedStateID(ctx, graphqlClient, issue.TeamID)
		if err != nil {
			return IssueSummary{}, err
		}

		started, err := IssueUpdate(ctx, graphqlClient, issueID, LinearIssueUpdateInput{
			AssigneeID: stringPtr(guard.target.Viewer.ID),
			StateID:    stringPtr(stateID),
		})
		if err != nil {
			return IssueSummary{}, fmt.Errorf("start issue %s: %w", issueID, err)
		}
		if !started.IssueUpdate.Success || started.IssueUpdate.Issue == nil {
			return IssueSummary{}, fmt.Errorf("%w: issue start returned no issue", ErrMutationFailed)
		}

		return issueSummaryFromFields(started.IssueUpdate.Issue.IssueSummaryFields), nil
	})
}

// CommentOnIssue adds a comment after resolving and comparing the pinned write target.
func CommentOnIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueCommentRequest,
) (IssueCommentResult, error) {
	if request.ID == "" {
		return IssueCommentResult{}, fmt.Errorf("%w: issue id is required", ErrWriteInvalid)
	}
	if request.Body == "" {
		return IssueCommentResult{}, fmt.Errorf("%w: body is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueCommentResult, error) {
		if _, err := guard.requireIssue(ctx, graphqlClient, request.ID); err != nil {
			return IssueCommentResult{}, err
		}

		comment, err := IssueCommentCreate(ctx, graphqlClient, LinearCommentCreateInput{
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
	})
}

// CloseIssue moves an issue to the team's completed workflow state after target comparison.
func CloseIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (IssueSummary, error) {
	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueSummary, error) {
		issue, err := guard.requireIssue(ctx, graphqlClient, issueID)
		if err != nil {
			return IssueSummary{}, err
		}
		stateID, err := firstCompletedStateID(ctx, graphqlClient, issue.TeamID)
		if err != nil {
			return IssueSummary{}, err
		}

		closed, err := IssueClose(ctx, graphqlClient, issueID, LinearIssueUpdateInput{
			StateID: stringPtr(stateID),
		})
		if err != nil {
			return IssueSummary{}, fmt.Errorf("close issue %s: %w", issueID, err)
		}
		if !closed.IssueUpdate.Success || closed.IssueUpdate.Issue == nil {
			return IssueSummary{}, fmt.Errorf("%w: issue close returned no issue", ErrMutationFailed)
		}

		return issueSummaryFromFields(closed.IssueUpdate.Issue.IssueSummaryFields), nil
	})
}

func firstCompletedStateID(ctx context.Context, graphqlClient graphql.Client, teamID string) (string, error) {
	states, err := CompletedWorkflowStates(ctx, graphqlClient, teamID, intPtr(50))
	if err != nil {
		return "", fmt.Errorf("list completed workflow states: %w", err)
	}
	if len(states.WorkflowStates.Nodes) == 0 {
		return "", fmt.Errorf("%w: completed workflow state missing for team_id=%s", ErrWriteInvalid, teamID)
	}

	state := states.WorkflowStates.Nodes[0]
	for _, candidate := range states.WorkflowStates.Nodes[1:] {
		if candidate.Position < state.Position {
			state = candidate
		}
	}

	return state.Id, nil
}

func firstStartedStateID(ctx context.Context, graphqlClient graphql.Client, teamID string) (string, error) {
	states, err := StartedWorkflowStates(ctx, graphqlClient, teamID, intPtr(50))
	if err != nil {
		return "", fmt.Errorf("list started workflow states: %w", err)
	}
	if len(states.WorkflowStates.Nodes) == 0 {
		return "", fmt.Errorf("%w: started workflow state missing for team_id=%s", ErrWriteInvalid, teamID)
	}

	state := states.WorkflowStates.Nodes[0]
	for _, candidate := range states.WorkflowStates.Nodes[1:] {
		if candidate.Position < state.Position {
			state = candidate
		}
	}

	return state.Id, nil
}

func firstStateIDOfType(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	stateType string,
) (string, error) {
	states, err := WorkflowStatesByType(ctx, graphqlClient, teamID, stateType, intPtr(50))
	if err != nil {
		return "", fmt.Errorf("list %s workflow states: %w", stateType, err)
	}
	if len(states.WorkflowStates.Nodes) == 0 {
		return "", fmt.Errorf("%w: %s workflow state missing for team_id=%s", ErrWriteInvalid, stateType, teamID)
	}

	state := states.WorkflowStates.Nodes[0]
	for _, candidate := range states.WorkflowStates.Nodes[1:] {
		if candidate.Position < state.Position {
			state = candidate
		}
	}

	return state.Id, nil
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
