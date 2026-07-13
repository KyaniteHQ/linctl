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
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment":       `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":         relationIssueRead(),
		"CommentUpdate": `{"commentUpdate":{"success":true,"comment":` + commentFieldsJSON("issue-id", "updated body") + `}}`,
	})}

	comment, err := UpdateComment(context.Background(), recorder, matchingTarget(), CommentUpdateRequest{
		ID:   "comment-id",
		Body: "updated body",
	})

	require.NoError(t, err)
	require.Equal(t, "comment-id", comment.ID)
	require.Equal(t, "updated body", comment.Body)
	require.JSONEq(t, `{
		"id": "comment-id",
		"input": {
			"body": "updated body"
		}
	}`, string(recorder.variablesFor(t, "CommentUpdate")))
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

func Test_ResolveComment_resolves_comment_when_target_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment":        `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":          relationIssueRead(),
		"CommentResolve": `{"commentResolve":{"success":true,"comment":` + commentFieldsJSON("issue-id", "existing body") + `}}`,
	})}

	comment, err := ResolveComment(context.Background(), recorder, matchingTarget(), "child-comment-id")

	require.NoError(t, err)
	require.Equal(t, "comment-id", comment.ID)
	require.Equal(t, "existing body", comment.Body)
	require.JSONEq(t, `{"id":"child-comment-id"}`, string(recorder.variablesFor(t, "CommentResolve")))
}

func Test_ResolveComment_requires_id(t *testing.T) {
	_, err := ResolveComment(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "")

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_ResolveComment_refuses_when_target_unresolved(t *testing.T) {
	_, err := ResolveComment(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, "comment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_ResolveComment_refuses_comment_without_an_issue_without_mutating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("", "existing body") + `}`,
	})}

	_, err := ResolveComment(context.Background(), recorder, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.Zero(t, recorder.countOf("CommentResolve"))
}

func Test_ResolveComment_refuses_when_issue_team_differs_without_mutating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueReadWrongTeam(),
	})}

	_, err := ResolveComment(context.Background(), recorder, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Zero(t, recorder.countOf("CommentResolve"))
}

func Test_ResolveComment_wraps_comment_read_error(t *testing.T) {
	_, err := ResolveComment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "comment-id",
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_ResolveComment_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueRead(),
	})

	_, err := ResolveComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.ErrorContains(t, err, "resolve comment comment-id")
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_ResolveComment_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment":        `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":          relationIssueRead(),
		"CommentResolve": `{"commentResolve":{"success":false,"comment":` + commentFieldsJSON("issue-id", "existing body") + `}}`,
	})

	_, err := ResolveComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UnresolveComment_unresolves_comment_when_target_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment":          `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":            relationIssueRead(),
		"CommentUnresolve": `{"commentUnresolve":{"success":true,"comment":` + commentFieldsJSON("issue-id", "existing body") + `}}`,
	})}

	comment, err := UnresolveComment(context.Background(), recorder, matchingTarget(), "child-comment-id")

	require.NoError(t, err)
	require.Equal(t, "comment-id", comment.ID)
	require.Equal(t, "existing body", comment.Body)
	require.JSONEq(t, `{"id":"child-comment-id"}`, string(recorder.variablesFor(t, "CommentUnresolve")))
}

func Test_UnresolveComment_requires_id(t *testing.T) {
	_, err := UnresolveComment(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "")

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UnresolveComment_refuses_when_target_unresolved(t *testing.T) {
	_, err := UnresolveComment(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, "comment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UnresolveComment_refuses_comment_without_an_issue_without_mutating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("", "existing body") + `}`,
	})}

	_, err := UnresolveComment(context.Background(), recorder, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.Zero(t, recorder.countOf("CommentUnresolve"))
}

func Test_UnresolveComment_refuses_when_issue_team_differs_without_mutating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueReadWrongTeam(),
	})}

	_, err := UnresolveComment(context.Background(), recorder, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Zero(t, recorder.countOf("CommentUnresolve"))
}

func Test_UnresolveComment_wraps_comment_read_error(t *testing.T) {
	_, err := UnresolveComment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "comment-id",
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UnresolveComment_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment": `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":   relationIssueRead(),
	})

	_, err := UnresolveComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.ErrorContains(t, err, "unresolve comment comment-id")
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UnresolveComment_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"comment":          `{"comment":` + commentFieldsJSON("issue-id", "existing body") + `}`,
		"issue":            relationIssueRead(),
		"CommentUnresolve": `{"commentUnresolve":{"success":false,"comment":` + commentFieldsJSON("issue-id", "existing body") + `}}`,
	})

	_, err := UnresolveComment(context.Background(), graphqlClient, matchingTarget(), "comment-id")

	require.ErrorIs(t, err, ErrMutationFailed)
}
