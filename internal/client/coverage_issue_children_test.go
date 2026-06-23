package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientReadScenarios_return_issue_child_metadata_reads(t *testing.T) {
	// Given
	endCursor := "cursor-1"
	comments := `{"nodes":[` + commentMetadataJSON("issue-id", "", "user-id") +
		`],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}`
	needs := `{"nodes":[` + customerNeedJSON() +
		`],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}`
	sharedAccess := `{"isShared":true,"viewerHasOnlySharedAccess":false,` +
		`"sharedWithCount":2,"disallowedIssueFields":["description","priority"]}`
	graphqlClient := fakeGraphQLClient{
		"issue_needs":                       `{"issue":{"id":"issue-id","identifier":"LIT-1","needs":` + needs + `}}`,
		"issue_formerNeeds":                 `{"issue":{"id":"issue-id","identifier":"LIT-1","formerNeeds":` + needs + `}}`,
		"issue_sharedAccess":                `{"issue":{"id":"issue-id","identifier":"LIT-1","sharedAccess":` + sharedAccess + `}}`,
		"issueVcsBranchSearch_comments":     `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","comments":` + comments + `}}`,
		"issueVcsBranchSearch_needs":        `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","needs":` + needs + `}}`,
		"issueVcsBranchSearch_formerNeeds":  `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","formerNeeds":` + needs + `}}`,
		"issueVcsBranchSearch_sharedAccess": `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","sharedAccess":` + sharedAccess + `}}`,
		"attachmentIssue_comments":          `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","comments":` + comments + `}}`,
		"attachmentIssue_needs":             `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","needs":` + needs + `}}`,
		"attachmentIssue_formerNeeds":       `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","formerNeeds":` + needs + `}}`,
		"attachmentIssue_sharedAccess":      `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","sharedAccess":` + sharedAccess + `}}`,
	}

	// When
	issueNeeds, err := ListIssueNeeds(context.Background(), graphqlClient, "LIT-1", 1)
	require.NoError(t, err)
	issueFormerNeeds, err := ListIssueFormerNeeds(context.Background(), graphqlClient, "LIT-1", 1)
	require.NoError(t, err)
	issueSharedAccess, err := GetIssueSharedAccess(context.Background(), graphqlClient, "LIT-1")
	require.NoError(t, err)
	branchComments, err := ListIssueVCSBranchComments(context.Background(), graphqlClient, "omer/branch", 1)
	require.NoError(t, err)
	branchNeeds, err := ListIssueVCSBranchNeeds(context.Background(), graphqlClient, "omer/branch", 1)
	require.NoError(t, err)
	branchFormerNeeds, err := ListIssueVCSBranchFormerNeeds(context.Background(), graphqlClient, "omer/branch", 1)
	require.NoError(t, err)
	branchSharedAccess, err := GetIssueVCSBranchSharedAccess(context.Background(), graphqlClient, "omer/branch")
	require.NoError(t, err)
	attachmentComments, err := ListAttachmentIssueComments(context.Background(), graphqlClient, "attachment-id", 1)
	require.NoError(t, err)
	attachmentNeeds, err := ListAttachmentIssueNeeds(context.Background(), graphqlClient, "attachment-id", 1)
	require.NoError(t, err)
	attachmentFormerNeeds, err := ListAttachmentIssueFormerNeeds(context.Background(), graphqlClient, "attachment-id", 1)
	require.NoError(t, err)
	attachmentSharedAccess, err := GetAttachmentIssueSharedAccess(context.Background(), graphqlClient, "attachment-id")
	require.NoError(t, err)

	// Then
	require.True(t, issueNeeds.HasNextPage)
	require.Equal(t, &endCursor, issueNeeds.EndCursor)
	require.Equal(t, "customer-need-id", issueNeeds.Needs[0].ID)
	require.Equal(t, "Acme", issueNeeds.Needs[0].CustomerName)
	require.Equal(t, "customer-need-id", issueFormerNeeds.Needs[0].ID)
	require.Equal(t, "comment-id", branchComments.Comments[0].ID)
	require.Equal(t, "Omer", branchComments.Comments[0].DisplayName)
	require.Equal(t, "customer-need-id", branchNeeds.Needs[0].ID)
	require.Equal(t, "customer-need-id", branchFormerNeeds.Needs[0].ID)
	require.Equal(t, "comment-id", attachmentComments.Comments[0].ID)
	require.Equal(t, "customer-need-id", attachmentNeeds.Needs[0].ID)
	require.Equal(t, "customer-need-id", attachmentFormerNeeds.Needs[0].ID)
	for _, access := range []IssueSharedAccessSummary{
		issueSharedAccess,
		branchSharedAccess,
		attachmentSharedAccess,
	} {
		require.Equal(t, "issue-id", access.IssueID)
		require.Equal(t, "LIT-1", access.Identifier)
		require.True(t, access.IsShared)
		require.Equal(t, 2, access.SharedWithCount)
		require.Equal(t, []string{"description", "priority"}, access.DisallowedIssueFields)
	}
}

func Test_ClientFailureScenarios_wrap_issue_child_metadata_errors(t *testing.T) {
	errorClient := errorGraphQLClient{err: errors.New("linear unavailable")}

	_, err := ListIssueNeeds(context.Background(), errorClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue customer needs LIT-1")

	_, err = ListIssueFormerNeeds(context.Background(), errorClient, "LIT-1", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue former customer needs LIT-1")

	_, err = GetIssueSharedAccess(context.Background(), errorClient, "LIT-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue shared access LIT-1")

	_, err = ListIssueVCSBranchComments(context.Background(), errorClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch comments omer/branch")

	_, err = ListIssueVCSBranchNeeds(context.Background(), errorClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch customer needs omer/branch")

	_, err = ListIssueVCSBranchFormerNeeds(context.Background(), errorClient, "omer/branch", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list issue vcs branch former customer needs omer/branch")

	_, err = GetIssueVCSBranchSharedAccess(context.Background(), errorClient, "omer/branch")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get issue vcs branch shared access omer/branch")

	_, err = ListAttachmentIssueComments(context.Background(), errorClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue comments attachment-id")

	_, err = ListAttachmentIssueNeeds(context.Background(), errorClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue customer needs attachment-id")

	_, err = ListAttachmentIssueFormerNeeds(context.Background(), errorClient, "attachment-id", 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list attachment issue former customer needs attachment-id")

	_, err = GetAttachmentIssueSharedAccess(context.Background(), errorClient, "attachment-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "get attachment issue shared access attachment-id")

	nilBranchClient := fakeGraphQLClient{
		"issueVcsBranchSearch_comments":     `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_needs":        `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_formerNeeds":  `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_sharedAccess": `{"issueVcsBranchSearch":null}`,
	}

	_, err = ListIssueVCSBranchComments(context.Background(), nilBranchClient, "omer/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")

	_, err = ListIssueVCSBranchNeeds(context.Background(), nilBranchClient, "omer/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")

	_, err = ListIssueVCSBranchFormerNeeds(context.Background(), nilBranchClient, "omer/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")

	_, err = GetIssueVCSBranchSharedAccess(context.Background(), nilBranchClient, "omer/branch")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.Contains(t, err.Error(), "not found")
}
