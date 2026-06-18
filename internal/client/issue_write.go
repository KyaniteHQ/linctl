package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// IssueCreateRequest describes a guarded issue create.
type IssueCreateRequest struct {
	Title       string
	Description string
}

// IssueUpdateRequest describes a guarded issue update.
type IssueUpdateRequest struct {
	ID          string
	Title       string
	Description string
}

// IssueCommentRequest describes a guarded issue comment.
type IssueCommentRequest struct {
	ID   string
	Body string
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
}

// LinearIssueUpdateInput is the sparse Linear issueUpdate payload linctl supports.
type LinearIssueUpdateInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	StateID     *string `json:"stateId,omitempty"`
}

// LinearCommentCreateInput is the sparse Linear commentCreate payload linctl supports.
type LinearCommentCreateInput struct {
	Body    *string `json:"body,omitempty"`
	IssueID *string `json:"issueId,omitempty"`
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
	target, err := ResolveTarget(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	input := LinearIssueCreateInput{
		Title:       stringPtr(request.Title),
		Description: optionalString(request.Description),
		TeamID:      target.Team.ID,
	}
	if target.Project != nil {
		input.ProjectID = stringPtr(target.Project.ID)
	}
	created, err := IssueCreate(ctx, graphqlClient, input)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("create issue: %w", err)
	}
	if !created.IssueCreate.Success || created.IssueCreate.Issue == nil {
		return IssueSummary{}, fmt.Errorf("%w: issueCreate returned no issue", ErrMutationFailed)
	}

	return issueSummaryFromFields(created.IssueCreate.Issue.IssueSummaryFields), nil
}

// UpdateIssue updates an issue after resolving and comparing the pinned write target.
func UpdateIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueUpdateRequest,
) (IssueSummary, error) {
	if request.ID == "" {
		return IssueSummary{}, fmt.Errorf("%w: issue id is required", ErrWriteInvalid)
	}
	if request.Title == "" && request.Description == "" {
		return IssueSummary{}, fmt.Errorf("%w: title or description is required", ErrWriteInvalid)
	}
	if _, err := guardIssueWrite(ctx, graphqlClient, expected, request.ID); err != nil {
		return IssueSummary{}, err
	}

	updated, err := IssueUpdate(ctx, graphqlClient, request.ID, LinearIssueUpdateInput{
		Title:       optionalString(request.Title),
		Description: optionalString(request.Description),
	})
	if err != nil {
		return IssueSummary{}, fmt.Errorf("update issue %s: %w", request.ID, err)
	}
	if !updated.IssueUpdate.Success || updated.IssueUpdate.Issue == nil {
		return IssueSummary{}, fmt.Errorf("%w: issueUpdate returned no issue", ErrMutationFailed)
	}

	return issueSummaryFromFields(updated.IssueUpdate.Issue.IssueSummaryFields), nil
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
	if _, err := guardIssueWrite(ctx, graphqlClient, expected, request.ID); err != nil {
		return IssueCommentResult{}, err
	}

	comment, err := IssueCommentCreate(ctx, graphqlClient, LinearCommentCreateInput{
		Body:    stringPtr(request.Body),
		IssueID: stringPtr(request.ID),
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

// CloseIssue moves an issue to the team's completed workflow state after target comparison.
func CloseIssue(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (IssueSummary, error) {
	issue, err := guardIssueWrite(ctx, graphqlClient, expected, issueID)
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
}

// ArchiveIssue archives an issue for integration cleanup after the write surface is verified.
func ArchiveIssue(ctx context.Context, graphqlClient graphql.Client, issueID string) (IssueSummary, error) {
	archived, err := IssueArchive(ctx, graphqlClient, issueID, boolPtr(false))
	if err != nil {
		return IssueSummary{}, fmt.Errorf("archive issue %s: %w", issueID, err)
	}
	if !archived.IssueArchive.Success || archived.IssueArchive.Entity == nil {
		return IssueSummary{}, fmt.Errorf("%w: issueArchive returned no issue", ErrMutationFailed)
	}

	return issueSummaryFromFields(archived.IssueArchive.Entity.IssueSummaryFields), nil
}

func guardIssueWrite(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	issueID string,
) (IssueSummary, error) {
	target, err := ResolveTarget(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}
	issue, err := GetIssueByID(ctx, graphqlClient, issueID)
	if err != nil {
		return IssueSummary{}, err
	}
	if issue.TeamID != target.Team.ID || issue.Team != target.Team.Key {
		return IssueSummary{}, fmt.Errorf(
			"%w: expected team_id=%s team_key=%s resolved issue team_id=%s team_key=%s",
			ErrTargetMismatch,
			target.Team.ID,
			target.Team.Key,
			issue.TeamID,
			issue.Team,
		)
	}
	if target.Project != nil && issue.ProjectID != target.Project.ID {
		return IssueSummary{}, fmt.Errorf(
			"%w: expected project_id=%s resolved issue project_id=%s",
			ErrTargetMismatch,
			target.Project.ID,
			issue.ProjectID,
		)
	}

	return issue, nil
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

func optionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func stringPtr(value string) *string {
	return &value
}
