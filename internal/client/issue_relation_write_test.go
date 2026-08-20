package client

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/Khan/genqlient/graphql"
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
	return issueRelationDepsJSONWithNextPage(blockedBy, false)
}

// issueRelationDepsJSONWithNextPage additionally controls the inverse-relations
// pageInfo, to exercise the fail-closed truncated-scan boundary.
func issueRelationDepsJSONWithNextPage(blockedBy bool, hasNextPage bool) string {
	inverse := `[]`
	if blockedBy {
		inverse = `[{"id":"blocked-by-relation","type":"blocks","issue":` + issueJSON(issueFixture{
			ID:         "related-issue-id",
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
		"inverseRelations":{"nodes":` + inverse + `,"pageInfo":{"hasNextPage":` + strconv.FormatBool(hasNextPage) + `,"endCursor":null}}
	}}`
}

func graphqlRequestID(request *graphql.Request) string {
	if request.Variables == nil {
		return ""
	}
	data, err := json.Marshal(request.Variables)
	if err != nil {
		return ""
	}
	var variables map[string]any
	if err := json.Unmarshal(data, &variables); err != nil {
		return ""
	}
	id, ok := variables["id"].(string)
	if !ok {
		return ""
	}

	return id
}

type issueLookupFake struct {
	byID  map[string]string
	inner graphql.Client
}

func (fake issueLookupFake) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "issue" {
		if payload, ok := fake.byID[graphqlRequestID(request)]; ok {
			return fakeGraphQLClient(map[string]string{"issue": `{"issue":` + payload + `}`}).
				MakeRequest(ctx, request, response)
		}
	}

	return fake.inner.MakeRequest(ctx, request, response)
}

func relationPinnedIssueJSON(id string, identifier string, title string) string {
	return issueJSON(issueFixture{
		ID:         id,
		Identifier: identifier,
		Title:      title,
		ProjectID:  "project-id",
		Project:    "fixture",
		StateID:    "state-id",
		State:      "Todo",
		StateType:  "unstarted",
	})
}

func relationIssuePairFake(extra map[string]string) graphql.Client {
	first := relationPinnedIssueJSON("issue-id", "LIT-1", "First")
	second := relationPinnedIssueJSON("related-issue-id", "LIT-2", "Second")

	return issueLookupFake{
		byID: map[string]string{
			"LIT-1":            first,
			"issue-id":         first,
			"LIT-2":            second,
			"related-issue-id": second,
		},
		inner: issueWriteFakeClient(extra),
	}
}

func Test_CreateIssueRelation_links_issues_when_target_matches(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{
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
	graphqlClient := relationIssuePairFake(map[string]string{
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
	graphqlClient := relationIssuePairFake(map[string]string{
		"IssueDependencies": issueRelationDepsJSON(true),
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "blocks",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "create a cycle")
}

func Test_CreateIssueRelation_wraps_dependency_read_error(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "blocks",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

// Test_CreateIssueRelation_refuses_blocks_when_dependency_scan_is_truncated
// proves a truncated blocked-by scan (hasNextPage=true) with no match found in
// the fetched page fails closed instead of assuming no cycle exists.
func Test_CreateIssueRelation_refuses_blocks_when_dependency_scan_is_truncated(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{
		"IssueDependencies": issueRelationDepsJSONWithNextPage(false, true),
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "blocks",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "more than 50 relations")
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
		context.Background(), relationIssuePairFake(map[string]string{}), matchingTarget(),
		IssueRelationCreateRequest{IssueID: "LIT-1", RelatedIssueID: "LIT-1", Type: "related"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "cannot relate to itself")
}

func Test_CreateIssueRelation_rejects_self_relation_when_identifier_and_uuid_name_the_same_issue(t *testing.T) {
	same := relationPinnedIssueJSON("issue-uuid-42", "LIT-42", "Same issue")
	recorder := &mutationRecordingClient{inner: issueLookupFake{
		byID: map[string]string{
			"LIT-42":        same,
			"issue-uuid-42": same,
		},
		inner: issueWriteFakeClient(map[string]string{
			"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
				relationWriteJSON("blocks") + `}}`,
		}),
	}}

	_, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-42",
		RelatedIssueID: "issue-uuid-42",
		Type:           "blocks",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "cannot relate to itself")
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
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
	graphqlClient := relationIssuePairFake(map[string]string{})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssueRelation_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{
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
