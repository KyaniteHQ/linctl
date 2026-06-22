package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func commentFieldsJSON(issueID string, body string) string {
	issue := `null`
	if issueID != "" {
		issue = `"` + issueID + `"`
	}

	return `{
		"id":"comment-id",
		"body":"` + body + `",
		"url":"https://linear.app/comment/comment-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:02:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":null,
		"issueId":` + issue + `,
		"projectId":null,
		"projectUpdateId":null,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func Test_UpdateComment_edits_comment_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment":       `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":         relationIssueRead(),
		"CommentUpdate": `{"commentUpdate":{"success":true,"comment":` + commentFieldsJSON("issue-id", "updated body") + `}}`,
	})

	comment, err := UpdateComment(context.Background(), graphqlClient, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "updated body",
	})

	require.NoError(t, err)
	require.Equal(t, "comment-id", comment.ID)
	require.Equal(t, "updated body", comment.Body)
}

func Test_UpdateComment_requires_id(t *testing.T) {
	_, err := UpdateComment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		CommentUpdateRequest{Body: "x"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateComment_requires_body(t *testing.T) {
	_, err := UpdateComment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		CommentUpdateRequest{ID: "comment-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateComment_refuses_when_target_unresolved(t *testing.T) {
	_, err := UpdateComment(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, CommentUpdateRequest{ID: "comment-id", Body: "x"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateComment_wraps_comment_read_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := UpdateComment(context.Background(), graphqlClient, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "x",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateComment_refuses_comment_without_an_issue(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("", "existing body") + `}`,
	})

	_, err := UpdateComment(context.Background(), graphqlClient, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "x",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateComment_refuses_when_issue_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueReadWrongTeam(),
	})

	_, err := UpdateComment(context.Background(), graphqlClient, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "x",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateComment_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueRead(),
	})

	_, err := UpdateComment(context.Background(), graphqlClient, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "x",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateComment_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment":       `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":         relationIssueRead(),
		"CommentUpdate": `{"commentUpdate":{"success":false,"comment":` + commentFieldsJSON("issue-id", "x") + `}}`,
	})

	_, err := UpdateComment(context.Background(), graphqlClient, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "x",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_DeleteComment_removes_comment_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment":       `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":         relationIssueRead(),
		"CommentDelete": `{"commentDelete":{"success":true,"entityId":"comment-id"}}`,
	})

	id, err := DeleteComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.NoError(t, err)
	require.Equal(t, "comment-id", id)
}

func Test_DeleteComment_requires_id(t *testing.T) {
	_, err := DeleteComment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "",
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_DeleteComment_refuses_when_target_unresolved(t *testing.T) {
	_, err := DeleteComment(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, "comment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteComment_refuses_comment_without_an_issue(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("", "existing body") + `}`,
	})

	_, err := DeleteComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_DeleteComment_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueRead(),
	})

	_, err := DeleteComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteComment_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment":       `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":         relationIssueRead(),
		"CommentDelete": `{"commentDelete":{"success":false,"entityId":"comment-id"}}`,
	})

	_, err := DeleteComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrMutationFailed)
}
