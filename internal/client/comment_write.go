package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// CommentUpdateRequest describes a guarded comment edit.
type CommentUpdateRequest struct {
	ID   string
	Body string
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

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return CommentSummary{}, err
	}

	return guard.updateComment(ctx, request)
}

func (guard *guardedClient) updateComment(ctx context.Context, request CommentUpdateRequest) (CommentSummary, error) {
	if err := guard.requireCommentTarget(ctx, request.ID); err != nil {
		return CommentSummary{}, err
	}

	updated, err := gql.CommentUpdate(ctx, guard.graphqlClient, request.ID, LinearCommentUpdateInput{
		Body: stringPtr(request.Body),
	})
	if err != nil {
		return CommentSummary{}, fmt.Errorf("update comment %s: %w", request.ID, err)
	}
	if err := mutationSuccess(updated.CommentUpdate.Success, "commentUpdate"); err != nil {
		return CommentSummary{}, err
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

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}

	return guard.deleteComment(ctx, commentID)
}

func (guard *guardedClient) deleteComment(ctx context.Context, commentID string) (string, error) {
	if err := guard.requireCommentTarget(ctx, commentID); err != nil {
		return "", err
	}

	deleted, err := gql.CommentDelete(ctx, guard.graphqlClient, commentID)
	if err != nil {
		return "", fmt.Errorf("delete comment %s: %w", commentID, err)
	}
	if err := mutationSuccess(deleted.CommentDelete.Success, "commentDelete"); err != nil {
		return "", err
	}

	return commentID, nil
}

// ResolveComment resolves a comment thread after comparing the pinned target
// through the comment's parent issue.
func ResolveComment(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	commentID string,
) (CommentSummary, error) {
	if commentID == "" {
		return CommentSummary{}, fmt.Errorf("%w: comment id is required", ErrWriteInvalid)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return CommentSummary{}, err
	}

	return guard.resolveComment(ctx, commentID)
}

func (guard *guardedClient) resolveComment(ctx context.Context, commentID string) (CommentSummary, error) {
	if err := guard.requireCommentTarget(ctx, commentID); err != nil {
		return CommentSummary{}, err
	}

	resolved, err := gql.CommentResolve(ctx, guard.graphqlClient, commentID)
	if err != nil {
		return CommentSummary{}, fmt.Errorf("resolve comment %s: %w", commentID, err)
	}
	if err := mutationSuccess(resolved.CommentResolve.Success, "commentResolve"); err != nil {
		return CommentSummary{}, err
	}

	return topLevelCommentSummary(resolved.CommentResolve.Comment.TopLevelCommentSummaryFields), nil
}

// UnresolveComment reopens a comment thread after comparing the pinned target
// through the comment's parent issue.
func UnresolveComment(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	commentID string,
) (CommentSummary, error) {
	if commentID == "" {
		return CommentSummary{}, fmt.Errorf("%w: comment id is required", ErrWriteInvalid)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return CommentSummary{}, err
	}

	return guard.unresolveComment(ctx, commentID)
}

func (guard *guardedClient) unresolveComment(ctx context.Context, commentID string) (CommentSummary, error) {
	if err := guard.requireCommentTarget(ctx, commentID); err != nil {
		return CommentSummary{}, err
	}

	unresolved, err := gql.CommentUnresolve(ctx, guard.graphqlClient, commentID)
	if err != nil {
		return CommentSummary{}, fmt.Errorf("unresolve comment %s: %w", commentID, err)
	}
	if err := mutationSuccess(unresolved.CommentUnresolve.Success, "commentUnresolve"); err != nil {
		return CommentSummary{}, err
	}

	return topLevelCommentSummary(unresolved.CommentUnresolve.Comment.TopLevelCommentSummaryFields), nil
}

// guardCommentTarget resolves a comment and confirms its parent issue belongs to
// the resolved team. Comments not attached to an issue are refused because the
// issue guard cannot prove their target.
func (guard *guardedClient) requireCommentTarget(
	ctx context.Context,
	commentID string,
) error {
	comment, err := GetCommentByID(ctx, guard.graphqlClient, commentID)
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
	_, err = guard.requireIssue(ctx, comment.IssueID)

	return err
}
