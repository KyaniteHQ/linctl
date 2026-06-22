package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// CommentUpdateRequest describes a guarded comment edit.
type CommentUpdateRequest struct {
	ID   string
	Body string
}

// LinearCommentUpdateInput is the sparse Linear commentUpdate payload linctl supports.
type LinearCommentUpdateInput struct {
	Body *string `json:"body,omitempty"`
}

// UpdateComment edits a comment after resolving the comment and comparing the
// pinned target through its parent issue. Only issue-attached comments are
// guarded; comments on other entities are refused.
func UpdateComment(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request CommentUpdateRequest,
) (CommentSummary, error) {
	if request.ID == "" {
		return CommentSummary{}, fmt.Errorf("%w: comment id is required", ErrWriteInvalid)
	}
	if request.Body == "" {
		return CommentSummary{}, fmt.Errorf("%w: body is required", ErrWriteInvalid)
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return CommentSummary{}, err
	}
	if err := guardCommentTarget(ctx, graphqlClient, guard, request.ID); err != nil {
		return CommentSummary{}, err
	}

	updated, err := CommentUpdate(ctx, graphqlClient, request.ID, LinearCommentUpdateInput{
		Body: stringPtr(request.Body),
	})
	if err != nil {
		return CommentSummary{}, fmt.Errorf("update comment %s: %w", request.ID, err)
	}
	if !updated.CommentUpdate.Success {
		return CommentSummary{}, fmt.Errorf("%w: commentUpdate reported no success", ErrMutationFailed)
	}

	return topLevelCommentSummary(updated.CommentUpdate.Comment.TopLevelCommentSummaryFields), nil
}

// DeleteComment removes a comment after resolving the comment and comparing the
// pinned target through its parent issue. Comment delete is the one approved
// delete and is restricted to issue-attached comments.
func DeleteComment(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	commentID string,
) (string, error) {
	if commentID == "" {
		return "", fmt.Errorf("%w: comment id is required", ErrWriteInvalid)
	}
	guard, err := newWriteGuard(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}
	if err := guardCommentTarget(ctx, graphqlClient, guard, commentID); err != nil {
		return "", err
	}

	deleted, err := CommentDelete(ctx, graphqlClient, commentID)
	if err != nil {
		return "", fmt.Errorf("delete comment %s: %w", commentID, err)
	}
	if !deleted.CommentDelete.Success {
		return "", fmt.Errorf("%w: commentDelete reported no success", ErrMutationFailed)
	}

	return commentID, nil
}

// guardCommentTarget resolves a comment and confirms its parent issue belongs to
// the resolved team. Comments not attached to an issue are refused because the
// issue guard cannot prove their target.
func guardCommentTarget(
	ctx context.Context,
	graphqlClient graphql.Client,
	guard writeGuard,
	commentID string,
) error {
	comment, err := GetCommentByID(ctx, graphqlClient, commentID)
	if err != nil {
		return err
	}
	if comment.IssueID == "" {
		return fmt.Errorf(
			"%w: comment %s is not attached to an issue; only issue comments are guarded",
			ErrWriteInvalid,
			commentID,
		)
	}
	_, err = guard.requireIssue(ctx, graphqlClient, comment.IssueID)

	return err
}
