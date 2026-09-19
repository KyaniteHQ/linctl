package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func attachmentLinkURLJSON(success bool) string {
	flag := "false"
	if success {
		flag = "true"
	}

	return `{"attachmentCreate":{"success":` + flag + `,"attachment":{
		"id":"attachment-id",
		"title":"Linked PR",
		"subtitle":"PR #1",
		"url":"https://example.com/pr/1",
		"sourceType":"github"
	}}}`
}

func attachmentIssueRead() string {
	return `{"attachmentIssue":` + issueJSON(issueFixture{
		Identifier: "LIT-1",
		Title:      "First",
		ProjectID:  "project-id",
		Project:    "fixture",
		StateID:    "state-id",
		State:      "Todo",
		StateType:  "unstarted",
	}) + `}`
}

func emptyAttachmentIssueRead() string {
	return `{"attachmentIssue":{
		"id":"",
		"identifier":"",
		"title":"",
		"url":"",
		"priority":0,
		"priorityLabel":"",
		"team":{"id":"team-id","key":"LIT","name":"linctl-it","organization":{"id":"org-id"}},
		"state":{"id":"state-id","name":"Todo","type":"unstarted"},
		"assignee":null,
		"project":null
	}}`
}

func Test_DeleteAttachment_removes_attachment_when_target_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"attachmentIssue":  attachmentIssueRead(),
		"issue":            relationIssueRead(),
		"AttachmentDelete": `{"attachmentDelete":{"success":true,"entityId":"attachment-id"}}`,
	})}

	id, err := DeleteAttachment(context.Background(), recorder, matchingTarget(), "attachment-id")

	require.NoError(t, err)
	require.Equal(t, "attachment-id", id)
	require.JSONEq(t, `{"id":"attachment-id"}`, string(recorder.variablesFor(t, "AttachmentDelete")))
}

func Test_DeleteAttachment_requires_id(t *testing.T) {
	_, err := DeleteAttachment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "",
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_DeleteAttachment_refuses_when_target_unresolved(t *testing.T) {
	_, err := DeleteAttachment(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, "attachment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteAttachment_refuses_attachment_without_an_issue_without_mutating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"attachmentIssue": emptyAttachmentIssueRead(),
	})}

	_, err := DeleteAttachment(context.Background(), recorder, matchingTarget(), "attachment-id")

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.Zero(t, recorder.countOf("AttachmentDelete"))
}

func Test_DeleteAttachment_refuses_when_issue_team_differs_without_mutating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"attachmentIssue": attachmentIssueRead(),
		"issue":           relationIssueReadWrongTeam(),
	})}

	_, err := DeleteAttachment(context.Background(), recorder, matchingTarget(), "attachment-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Zero(t, recorder.countOf("AttachmentDelete"))
}

func Test_DeleteAttachment_wraps_issue_read_error(t *testing.T) {
	_, err := DeleteAttachment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "attachment-id",
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteAttachment_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"attachmentIssue": attachmentIssueRead(),
		"issue":           relationIssueRead(),
	})

	_, err := DeleteAttachment(context.Background(), graphqlClient, matchingTarget(), "attachment-id")

	require.ErrorContains(t, err, "delete attachment attachment-id")
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteAttachment_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"attachmentIssue":  attachmentIssueRead(),
		"issue":            relationIssueRead(),
		"AttachmentDelete": `{"attachmentDelete":{"success":false,"entityId":"attachment-id"}}`,
	})

	_, err := DeleteAttachment(context.Background(), graphqlClient, matchingTarget(), "attachment-id")

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_LinkIssueAttachment_attaches_url_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":             relationIssueRead(),
		"AttachmentLinkURL": attachmentLinkURLJSON(true),
	})

	attachment, err := LinkIssueAttachment(context.Background(), graphqlClient, matchingTarget(), AttachmentLinkRequest{
		IssueID:  "LIT-1",
		URL:      "https://example.com/pr/1",
		Title:    "Linked PR",
		Subtitle: "PR #1",
	})

	require.NoError(t, err)
	require.Equal(t, "attachment-id", attachment.ID)
	require.Equal(t, "https://example.com/pr/1", attachment.URL)
	require.Equal(t, "Linked PR", attachment.Title)
	require.Equal(t, "PR #1", attachment.Subtitle)
}

func Test_LinkIssueAttachment_requires_issue_id(t *testing.T) {
	_, err := LinkIssueAttachment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		AttachmentLinkRequest{URL: "https://example.com/pr/1"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_LinkIssueAttachment_requires_url(t *testing.T) {
	_, err := LinkIssueAttachment(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		AttachmentLinkRequest{IssueID: "LIT-1"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_LinkIssueAttachment_refuses_when_target_unresolved(t *testing.T) {
	_, err := LinkIssueAttachment(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, AttachmentLinkRequest{IssueID: "LIT-1", URL: "https://example.com/pr/1"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_LinkIssueAttachment_refuses_when_issue_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueReadWrongTeam(),
	})

	_, err := LinkIssueAttachment(context.Background(), graphqlClient, matchingTarget(), AttachmentLinkRequest{
		IssueID: "LIT-1",
		URL:     "https://example.com/pr/1",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_LinkIssueAttachment_wraps_issue_read_error(t *testing.T) {
	_, err := LinkIssueAttachment(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		AttachmentLinkRequest{IssueID: "LIT-1", URL: "https://example.com/pr/1"})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_LinkIssueAttachment_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueRead(),
	})

	_, err := LinkIssueAttachment(context.Background(), graphqlClient, matchingTarget(), AttachmentLinkRequest{
		IssueID: "LIT-1",
		URL:     "https://example.com/pr/1",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_LinkIssueAttachment_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":             relationIssueRead(),
		"AttachmentLinkURL": attachmentLinkURLJSON(false),
	})

	_, err := LinkIssueAttachment(context.Background(), graphqlClient, matchingTarget(), AttachmentLinkRequest{
		IssueID: "LIT-1",
		URL:     "https://example.com/pr/1",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}
