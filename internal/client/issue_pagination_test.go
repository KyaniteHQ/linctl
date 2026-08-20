package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ListIssueAttachments_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issue_attachments": `{"issue":{"attachments":{"nodes":[` +
			issueAttachmentPageJSON("issue-attachment-1", "First") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"issue-attachment-cursor-1"}}}}`,
		"issue_attachments:issue-attachment-cursor-1": `{"issue":{"attachments":{"nodes":[` +
			issueAttachmentPageJSON("issue-attachment-2", "Second") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	attachments, err := ListIssueAttachments(context.Background(), graphqlClient, "LIT-1", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"issue-attachment-1", "issue-attachment-2"}, []string{
		attachments.Attachments[0].ID,
		attachments.Attachments[1].ID,
	})
	require.False(t, attachments.HasNextPage)
}

func Test_ListIssueNeeds_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issue_needs": `{"issue":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			issueNeedPageJSON("need-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"need-cursor-1"}}}}`,
		"issue_needs:need-cursor-1": `{"issue":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			issueNeedPageJSON("need-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	needs, err := ListIssueNeeds(context.Background(), graphqlClient, "LIT-1", 100)

	require.NoError(t, err)
	require.Equal(t, "issue-id", needs.IssueID)
	require.Equal(t, "LIT-1", needs.Identifier)
	require.Equal(t, []string{"need-1", "need-2"}, []string{
		needs.Needs[0].ID,
		needs.Needs[1].ID,
	})
	require.False(t, needs.HasNextPage)
}

func Test_ListIssueRelations_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issueRelations": `{"issueRelations":{"nodes":[` +
			issueRelationPageJSON("relation-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"relation-cursor-1"}}}`,
		"issueRelations:relation-cursor-1": `{"issueRelations":{"nodes":[` +
			issueRelationPageJSON("relation-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	relations, err := ListIssueRelations(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"relation-1", "relation-2"}, []string{
		relations.Relations[0].ID,
		relations.Relations[1].ID,
	})
	require.False(t, relations.HasNextPage)
}

func Test_ListIssueToReleases_pages_across_two_responses(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issueToReleases": `{"issueToReleases":{"nodes":[` +
			issueToReleasePageJSON("issue-to-release-1") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"issue-to-release-cursor-1"}}}`,
		"issueToReleases:issue-to-release-cursor-1": `{"issueToReleases":{"nodes":[` +
			issueToReleasePageJSON("issue-to-release-2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	associations, err := ListIssueToReleases(context.Background(), graphqlClient, 100)

	require.NoError(t, err)
	require.Equal(t, []string{"issue-to-release-1", "issue-to-release-2"}, []string{
		associations.Associations[0].ID,
		associations.Associations[1].ID,
	})
	require.False(t, associations.HasNextPage)
}

func Test_ListIssueVCSBranchComments_pages_and_keeps_parent_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issueVcsBranchSearch_comments": `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commentMetadataJSONWithID("comment-1", "", "", "", "user-id") +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"comment-cursor-1"}}}}`,
		"issueVcsBranchSearch_comments:comment-cursor-1": `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commentMetadataJSONWithID("comment-2", "", "", "", "user-id") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	comments, err := ListIssueVCSBranchComments(context.Background(), graphqlClient, "omer/branch", 100)

	require.NoError(t, err)
	require.Equal(t, "issue-id", comments.IssueID)
	require.Equal(t, "LIT-1", comments.Identifier)
	require.Equal(t, []string{"comment-1", "comment-2"}, []string{
		comments.Comments[0].ID,
		comments.Comments[1].ID,
	})
	require.False(t, comments.HasNextPage)
}

func Test_ListNextIssuesByTeam_pages_then_ranks_across_the_window(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"NextIssuesByTeam": `{"issues":{"nodes":[` +
			nextIssueJSON("LIT-LOW", "Older low priority", 4, "Low", "2026-01-01T00:00:00Z", []string{}) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"next-cursor-1"}}}`,
		"NextIssuesByTeam:next-cursor-1": `{"issues":{"nodes":[` +
			nextIssueJSON("LIT-BEST", "Later unblocking issue", 3, "Normal", "2026-06-01T00:00:00Z", []string{
				`{"type":"blocks","relatedIssue":{"id":"active-1","state":{"type":"started"}}}`,
				`{"type":"blocks","relatedIssue":{"id":"active-2","state":{"type":"unstarted"}}}`,
			}) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	issues, err := ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 100)

	require.NoError(t, err)
	require.Equal(t, []string{"LIT-BEST", "LIT-LOW"}, issueIdentifiers(issues.Issues))
	require.Equal(t, 2, issues.Issues[0].UnblocksCount)
	require.False(t, issues.HasNextPage)
}

func Test_ListNextIssuesByTeam_keeps_has_next_page_when_limit_stops_the_window(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"NextIssuesByTeam": `{"issues":{"nodes":[` +
			nextIssueJSON("LIT-LOW", "Older low priority", 4, "Low", "2026-01-01T00:00:00Z", []string{}) +
			`],"pageInfo":{"hasNextPage":true,"endCursor":"next-cursor-1"}}}`,
	}

	issues, err := ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 1)

	require.NoError(t, err)
	require.Equal(t, []string{"LIT-LOW"}, issueIdentifiers(issues.Issues))
	require.True(t, issues.HasNextPage)
}

func issueAttachmentPageJSON(id string, title string) string {
	return `{"id":"` + id + `","title":"` + title +
		`","subtitle":"feat","url":"https://example.com/` + id + `","sourceType":"github"}`
}

func issueNeedPageJSON(id string) string {
	return `{"id":"` + id + `","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:01:00Z","archivedAt":null,` +
		`"priority":1,"url":"https://example.com/need","customer":{"id":"customer-id","name":"Acme"},` +
		`"issue":{"id":"issue-id","identifier":"LIT-1","title":"Need issue"},` +
		`"project":{"id":"project-id","name":"Customer project"}}`
}

func issueRelationPageJSON(id string) string {
	return `{"id":"` + id + `","type":"blocks","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z",` +
		`"archivedAt":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},` +
		`"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}`
}

func issueToReleasePageJSON(id string) string {
	return `{"id":"` + id + `","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,` +
		`"issue":{"id":"issue-id"},"release":{"id":"release-id"}}`
}
