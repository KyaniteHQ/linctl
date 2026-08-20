package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListAttachments_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"attachments": `{"attachments":{"nodes":[` +
			attachmentSummaryJSON("attachment-1", "First") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"attachment-cursor-1"}}}`,
		"attachments:attachment-cursor-1": `{"attachments":{"nodes":[` +
			attachmentSummaryJSON("attachment-2", "Second") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	attachments, err := ListAttachments(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"attachment-1", "attachment-2"}, []string{
		attachments.Attachments[0].ID,
		attachments.Attachments[1].ID,
	})
	require.False(t, attachments.HasNextPage)
}

func Test_ListAttachmentsForURL_pages_across_two_responses(t *testing.T) {
	url := "https://example.com/spec"
	graphqlClient := fakeGraphQLClient{
		"attachmentsForURL": `{"attachmentsForURL":{"nodes":[` +
			attachmentSummaryJSON("attachment-url-1", "First URL") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"url-cursor-1"}}}`,
		"attachmentsForURL:url-cursor-1": `{"attachmentsForURL":{"nodes":[` +
			attachmentSummaryJSON("attachment-url-2", "Second URL") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	attachments, err := ListAttachmentsForURL(context.Background(), graphqlClient, url, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"attachment-url-1", "attachment-url-2"}, []string{
		attachments.Attachments[0].ID,
		attachments.Attachments[1].ID,
	})
	require.False(t, attachments.HasNextPage)
}

func Test_ListAttachmentIssueAttachments_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"attachmentIssue_attachments": `{"attachmentIssue":{"attachments":{"nodes":[` +
			attachmentSummaryJSON("issue-attachment-1", "First issue attachment") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"issue-attachment-cursor-1"}}}}`,
		"attachmentIssue_attachments:issue-attachment-cursor-1": `{"attachmentIssue":{"attachments":{"nodes":[` +
			attachmentSummaryJSON("issue-attachment-2", "Second issue attachment") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	attachments, err := ListAttachmentIssueAttachments(context.Background(), graphqlClient, "attachment-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"issue-attachment-1", "issue-attachment-2"}, []string{
		attachments.Attachments[0].ID,
		attachments.Attachments[1].ID,
	})
	require.False(t, attachments.HasNextPage)
}

func Test_ListAttachmentIssueComments_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"attachmentIssue_comments": `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commentMetadataJSONWithID("comment-1", "", "", "", "user-id") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"comment-cursor-1"}}}}`,
		"attachmentIssue_comments:comment-cursor-1": `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commentMetadataJSONWithID("comment-2", "", "", "", "user-id") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	comments, err := ListAttachmentIssueComments(context.Background(), graphqlClient, "attachment-id", 100)

	require.NoError(t, err)
	require.Equal(t, "issue-id", comments.IssueID)
	require.Equal(t, "LIT-1", comments.Identifier)
	require.Equal(t, []string{"comment-1", "comment-2"}, []string{
		comments.Comments[0].ID,
		comments.Comments[1].ID,
	})
	require.False(t, comments.HasNextPage)
}

func attachmentSummaryJSON(id string, title string) string {
	return `{"id":"` + id + `","title":"` + title +
		`","subtitle":"feat","url":"https://example.com/` + id + `","sourceType":"github"}`
}
