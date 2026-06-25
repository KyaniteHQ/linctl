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
