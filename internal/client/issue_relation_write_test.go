package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func relationWriteJSON(relationType string) string {
	return `{
		"id":"relation-id",
		"type":"` + relationType + `",
		"createdAt":"2026-06-20T00:00:00Z",
		"updatedAt":"2026-06-20T00:00:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id","identifier":"LIT-1","title":"First"},
		"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Second"}
	}`
}

// relationIssueRead resolves any relation endpoint to an issue inside the pinned target.
func relationIssueRead() string {
	return `{"issue":` + issueJSON(issueFixture{
		Identifier: "LIT-1",
		Title:      "First",
		ProjectID:  "project-id",
		Project:    "fixture",
		StateID:    "state-id",
		State:      "Todo",
		StateType:  "unstarted",
	}) + `}`
}

// relationIssueReadWrongTeam resolves an endpoint to an issue owned by a different team.
func relationIssueReadWrongTeam() string {
	return `{"issue":{
		"id":"issue-id",
		"identifier":"LIT-1",
		"title":"First",
		"url":"https://linear.app/kyanite/issue/LIT-1",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"other-team","key":"OTHER","name":"other"},
		"state":{"id":"state-id","name":"Todo","type":"unstarted"},
		"assignee":null,
		"project":{"id":"project-id","name":"fixture"}
	}}`
}

// issueRelationDepsJSON builds an IssueDependencies response; blockedBy adds an
// inverse blocks relation whose blocker resolves to the shared issue-id.
func issueRelationDepsJSON(blockedBy bool) string {
	inverse := `[]`
	if blockedBy {
		inverse = `[{"id":"blocked-by-relation","type":"blocks","issue":` + issueJSON(issueFixture{
			Identifier: "LIT-2",
			Title:      "blocker",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}]`
	}

	return `{"issue":{
		"id":"issue-id",
		"identifier":"LIT-1",
		"parent":null,
		"children":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}},
		"relations":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}},
		"inverseRelations":{"nodes":` + inverse + `,"pageInfo":{"hasNextPage":false,"endCursor":null}}
	}}`
}

func Test_CreateIssueRelation_links_issues_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueRead(),
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
	})

	relation, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.NoError(t, err)
	require.Equal(t, "relation-id", relation.ID)
	require.Equal(t, "related", relation.Type)
}

func Test_CreateIssueRelation_allows_blocks_without_a_cycle(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":             relationIssueRead(),
		"IssueDependencies": issueRelationDepsJSON(false),
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("blocks") + `}}`,
	})

	relation, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "blocks",
	})

	require.NoError(t, err)
	require.Equal(t, "blocks", relation.Type)
}

func Test_CreateIssueRelation_refuses_blocks_that_close_a_cycle(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":             relationIssueRead(),
		"IssueDependencies": issueRelationDepsJSON(true),
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "blocks",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.Contains(t, err.Error(), "create a cycle")
}

func Test_CreateIssueRelation_wraps_dependency_read_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueRead(),
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "blocks",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssueRelation_requires_both_ids(t *testing.T) {
	_, err := CreateIssueRelation(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		IssueRelationCreateRequest{IssueID: "LIT-1", Type: "related"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssueRelation_rejects_self_relation(t *testing.T) {
	_, err := CreateIssueRelation(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		IssueRelationCreateRequest{IssueID: "LIT-1", RelatedIssueID: "LIT-1", Type: "related"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssueRelation_rejects_unknown_type(t *testing.T) {
	_, err := CreateIssueRelation(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		IssueRelationCreateRequest{IssueID: "LIT-1", RelatedIssueID: "LIT-2", Type: "mentions"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssueRelation_refuses_when_target_unresolved(t *testing.T) {
	_, err := CreateIssueRelation(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, IssueRelationCreateRequest{IssueID: "LIT-1", RelatedIssueID: "LIT-2", Type: "related"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssueRelation_refuses_when_issue_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueReadWrongTeam(),
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssueRelation_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueRead(),
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssueRelation_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": relationIssueRead(),
		"IssueRelationCreate": `{"issueRelationCreate":{"success":false,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_DeleteIssueRelation_removes_relation_when_target_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueRelation":       `{"issueRelation":` + relationWriteJSON("related") + `}`,
		"issue":               relationIssueRead(),
		"IssueRelationDelete": `{"issueRelationDelete":{"success":true,"entityId":"relation-id"}}`,
	})

	id, err := DeleteIssueRelation(context.Background(), graphqlClient, matchingTarget(), "relation-id")

	require.NoError(t, err)
	require.Equal(t, "relation-id", id)
}

func Test_DeleteIssueRelation_requires_id(t *testing.T) {
	_, err := DeleteIssueRelation(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "",
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_DeleteIssueRelation_refuses_when_target_unresolved(t *testing.T) {
	_, err := DeleteIssueRelation(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "WRONG",
		TeamID:  "wrong-id",
	}, "relation-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteIssueRelation_wraps_relation_read_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{})

	_, err := DeleteIssueRelation(context.Background(), graphqlClient, matchingTarget(), "relation-id")

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteIssueRelation_refuses_when_issue_team_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
		"issue":         relationIssueReadWrongTeam(),
	})

	_, err := DeleteIssueRelation(context.Background(), graphqlClient, matchingTarget(), "relation-id")

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteIssueRelation_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
		"issue":         relationIssueRead(),
	})

	_, err := DeleteIssueRelation(context.Background(), graphqlClient, matchingTarget(), "relation-id")

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_DeleteIssueRelation_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueRelation":       `{"issueRelation":` + relationWriteJSON("related") + `}`,
		"issue":               relationIssueRead(),
		"IssueRelationDelete": `{"issueRelationDelete":{"success":false,"entityId":"relation-id"}}`,
	})

	_, err := DeleteIssueRelation(context.Background(), graphqlClient, matchingTarget(), "relation-id")

	require.ErrorIs(t, err, ErrMutationFailed)
}
