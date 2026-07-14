package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CheckOrganizationExists_returns_operation_errors(t *testing.T) {
	_, err := CheckOrganizationExists(context.Background(), fakeGraphQLClient{}, "missing")

	require.ErrorContains(t, err, "missing fake response for organizationExists")
}

func Test_GetRateLimitStatus_returns_operation_errors(t *testing.T) {
	_, err := GetRateLimitStatus(context.Background(), fakeGraphQLClient{})

	require.ErrorContains(t, err, "missing fake response for rateLimitStatus")
}

func Test_ClientReadHelpers_cover_nil_actor_bot_summaries(t *testing.T) {
	require.Nil(t, actorBotSummary(nil))
	require.Nil(t, commentActorBotSummary(nil))
	require.Nil(t, issueActorBotSummary(nil))
}

func Test_ClientReadScenarios_return_not_found_for_null_vcs_branch_issue(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"issueVcsBranchSearch":                   `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_attachments":       `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_botActor":          `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_children":          `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_documents":         `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_formerAttachments": `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_history":           `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_inverseRelations":  `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_labels":            `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_relations":         `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_releases":          `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_stateHistory":      `{"issueVcsBranchSearch":null}`,
		"issueVcsBranchSearch_subscribers":       `{"issueVcsBranchSearch":null}`,
	}

	_, err := GetIssueByVCSBranch(context.Background(), graphqlClient, "missing/branch")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchAttachments(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = GetIssueVCSBranchBotActor(context.Background(), graphqlClient, "missing/branch")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchChildren(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchDocuments(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchFormerAttachments(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchHistory(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchInverseRelations(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchLabels(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchRelations(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchReleases(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchStateHistory(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
	_, err = ListIssueVCSBranchSubscribers(context.Background(), graphqlClient, "missing/branch", 1)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNotFound)
	require.ErrorIs(t, err, ErrNotFound)
}

func Test_ClientReadScenarios_rank_next_issues(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"NextIssuesByTeam": `{"issues":{"nodes":[` +
			nextIssueJSON("LIT-31", "Low priority standalone", 4, "Low", "2026-01-01T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-32", "Urgent standalone", 1, "Urgent", "2026-02-01T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-33", "Older high standalone", 2, "High", "2026-01-15T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-34", "Newer high standalone", 2, "High", "2026-02-15T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-35", "No priority standalone", 0, "No priority", "2026-01-01T00:00:00Z", []string{}) + `,` +
			nextIssueJSON("LIT-36", "Unblocks active work", 3, "Normal", "2026-03-01T00:00:00Z", []string{
				`{"type":"blocks","relatedIssue":{"id":"active-1","state":{"type":"started"}}}`,
				`{"type":"blocks","relatedIssue":{"id":"done-1","state":{"type":"completed"}}}`,
				`{"type":"relates","relatedIssue":{"id":"active-2","state":{"type":"unstarted"}}}`,
			}) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
	}

	issues, err := ListNextIssuesByTeam(context.Background(), graphqlClient, "team-id", 6)

	require.NoError(t, err)
	require.Equal(t, []string{"LIT-36", "LIT-32", "LIT-33", "LIT-34", "LIT-31", "LIT-35"}, issueIdentifiers(issues.Issues))
	require.Equal(t, 1, issues.Issues[0].UnblocksCount)
}
