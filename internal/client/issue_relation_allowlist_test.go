package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func relationIssueJSON(id string, identifier string, projectID string, projectName string) string {
	return issueJSON(issueFixture{
		ID:         id,
		Identifier: identifier,
		Title:      identifier,
		ProjectID:  projectID,
		Project:    projectName,
		StateID:    "state-id",
		State:      "Todo",
		StateType:  "unstarted",
	})
}

func Test_CreateIssueRelation_links_issues_in_explicit_allowed_projects(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := relationIssueJSON("related-issue-id", "LIT-2", "other-project", "other")
	inner := issueWriteFakeClient(map[string]string{
		"issue_relations": emptyIssueRelationsJSON(),
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
	})
	projects := newSequentialOpClient(inner)
	projects.payloads["project"] = []string{
		`{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		`{"project":` + projectJSON(projectFixture{ID: "other-project", Name: "other", Status: "Backlog"}) + `}`,
	}
	graphqlClient := issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: projects,
	}

	result, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"project-id", "other-project"},
	})

	require.NoError(t, err)
	require.Equal(t, "relation-id", result.ID)
	require.Equal(t, "project-id", result.Issue.ProjectID)
	require.Equal(t, "other-project", result.RelatedIssue.ProjectID)
}

func Test_CreateIssueRelation_refuses_project_outside_allowlist_before_mutation(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := relationIssueJSON("related-issue-id", "LIT-2", "other-project", "other")
	recorder := &mutationRecordingClient{inner: issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: issueWriteFakeClient(map[string]string{
			"project": `{"project":` + projectJSON(projectFixture{
				ID: "project-id", Name: "fixture", Status: "Backlog",
			}) + `}`,
		}),
	}}

	_, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"project-id"},
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_refuses_cross_organization_relation(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := `{
		"id":"related-issue-id",
		"identifier":"LIT-2",
		"title":"LIT-2",
		"url":"https://linear.app/kyanite/issue/LIT-2",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"other-team","key":"OTHER","name":"other","organization":{"id":"other-org"}},
		"state":{"id":"state-id","name":"Todo","type":"unstarted"},
		"assignee":null,
		"project":{"id":"other-project","name":"other"}
	}`
	recorder := &mutationRecordingClient{inner: issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: issueWriteFakeClient(map[string]string{}),
	}}

	_, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrCrossOrganizationRelation)
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.ErrorIs(t, ErrCrossOrganizationRelation, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_reconciles_existing_relation_without_create(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{
		"issue_relations": `{"issue":{"relations":{"nodes":[` +
			relationWriteJSON("related") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	})
	recorder := &mutationRecordingClient{inner: graphqlClient}

	result, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.NoError(t, err)
	require.Equal(t, "relation-id", result.ID)
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_reconciles_ambiguous_create_without_replay(t *testing.T) {
	inner := relationIssuePairFake(map[string]string{
		"IssueRelationCreate": "",
	})
	sequenced := newSequentialOpClient(inner)
	sequenced.failAt["IssueRelationCreate"] = 1
	sequenced.payloads["issue_relations"] = []string{
		emptyIssueRelationsJSON(),
		`{"issue":{"relations":{"nodes":[` + relationWriteJSON("related") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	result, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.NoError(t, err)
	require.Equal(t, "relation-id", result.ID)
	require.Equal(t, 1, sequenced.calls["IssueRelationCreate"])
}

func Test_CreateIssueRelation_refuses_empty_allowed_project_id(t *testing.T) {
	_, err := CreateIssueRelation(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{""},
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateIssueRelation_refuses_allowed_project_on_another_team(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSONWithTeam(
			projectFixture{ID: "other-project", Name: "other", Status: "Backlog"},
			"ops-team-id",
			"OPS",
		) + `}`,
	})}

	_, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"other-project"},
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_refuses_truncated_relation_scan(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{
		"issue_relations": `{"issue":{"relations":{"nodes":[],"pageInfo":{"hasNextPage":true,"endCursor":"c"}}}}`,
	})

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "cannot reconcile")
}

func Test_CreateIssueRelation_refuses_when_readback_moves_a_project(t *testing.T) {
	moved := issueFixture{
		ID: "issue-id", Identifier: "LIT-1", Title: "First",
		ProjectID: "other-project", Project: "other",
		StateID: "state-id", State: "Todo", StateType: "unstarted",
	}
	graphqlClient := withIssueAfterWrite(relationIssuePairFake(map[string]string{
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
	}), moved)

	_, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "relation changed an issue project")
}

func Test_CreateIssueRelation_keeps_duplicate_allowed_project_ids(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := relationIssueJSON("related-issue-id", "LIT-2", "project-id", "fixture")
	graphqlClient := issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: issueWriteFakeClient(map[string]string{
			"project": `{"project":` + projectJSON(projectFixture{
				ID: "project-id", Name: "fixture", Status: "Backlog",
			}) + `}`,
			"issue_relations": emptyIssueRelationsJSON(),
			"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
				relationWriteJSON("related") + `}}`,
			"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
		}),
	}

	result, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"project-id", "project-id"},
	})

	require.NoError(t, err)
	require.Equal(t, "project-id", result.Issue.ProjectID)
}

func Test_CreateIssueRelation_refuses_unreadable_allowed_project(t *testing.T) {
	_, err := CreateIssueRelation(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"missing-project"},
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "get project")
}

func Test_CreateIssueRelation_refuses_pinned_project_mismatch(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := relationIssueJSON("related-issue-id", "LIT-2", "other-project", "other")
	recorder := &mutationRecordingClient{inner: issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: issueWriteFakeClient(map[string]string{}),
	}}

	_, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "--allowed-project")
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_reconcile_returns_write_error_when_relation_missing(t *testing.T) {
	sequenced := newSequentialOpClient(relationIssuePairFake(map[string]string{
		"IssueRelationCreate": "",
	}))
	sequenced.failAt["IssueRelationCreate"] = 1

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "injected IssueRelationCreate failure")
}

func Test_requireIssueOnTeam_fails_closed_when_org_is_absent(t *testing.T) {
	issueBody := `{
		"id":"issue-id",
		"identifier":"LIT-1",
		"title":"First",
		"url":"https://linear.app/kyanite/issue/LIT-1",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"team-id","key":"LIT","name":"linctl-it"},
		"state":{"id":"state-id","name":"Todo","type":"unstarted"},
		"assignee":null,
		"project":{"id":"project-id","name":"fixture"}
	}`
	guard, err := newGuardedClient(context.Background(), issueWriteFakeClient(map[string]string{
		"issue": `{"issue":` + issueBody + `}`,
	}), matchingTarget())
	require.NoError(t, err)

	_, err = guard.requireIssueOnTeam(context.Background(), "LIT-1")

	require.ErrorIs(t, err, ErrCrossOrganizationRelation)
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateIssueRelation_returns_readback_error_when_relation_get_fails(t *testing.T) {
	sequenced := newSequentialOpClient(relationIssuePairFake(map[string]string{
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
	}))
	sequenced.failAt["issueRelation"] = 1

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "injected issueRelation failure")
}

func Test_CreateIssueRelation_wraps_relation_scan_error(t *testing.T) {
	sequenced := newSequentialOpClient(relationIssuePairFake(map[string]string{}))
	sequenced.failAt["issue_relations"] = 1

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "injected issue_relations failure")
}

func Test_CreateIssueRelation_reconcile_returns_write_error_when_readback_fails(t *testing.T) {
	inner := relationIssuePairFake(map[string]string{})
	sequenced := newSequentialOpClient(inner)
	sequenced.failAt["IssueRelationCreate"] = 1
	sequenced.payloads["issue_relations"] = []string{
		emptyIssueRelationsJSON(),
		`{"issue":{"relations":{"nodes":[` + relationWriteJSON("related") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}
	sequenced.failAt["issueRelation"] = 1

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "injected issueRelation failure")
}

func Test_CreateIssueRelation_refuses_unreadable_allowed_issue(t *testing.T) {
	_, err := CreateIssueRelation(context.Background(), issueWriteFakeClient(map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{
			ID: "project-id", Name: "fixture", Status: "Backlog",
		}) + `}`,
	}), matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"project-id"},
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "get issue")
}

func Test_CreateIssueRelation_returns_error_when_endpoint_readback_fails(t *testing.T) {
	sequenced := newSequentialOpClient(relationIssuePairFake(map[string]string{
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
	}))
	sequenced.failAt["issue"] = 3

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "injected issue failure")
}

func Test_CreateIssueRelation_returns_error_when_related_readback_fails(t *testing.T) {
	sequenced := newSequentialOpClient(relationIssuePairFake(map[string]string{
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
	}))
	sequenced.failAt["issue"] = 4

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.Error(t, err)
	require.ErrorContains(t, err, "injected issue failure")
}

func Test_CreateIssueRelation_reconcile_refuses_moved_project(t *testing.T) {
	moved := issueFixture{
		ID: "issue-id", Identifier: "LIT-1", Title: "First",
		ProjectID: "other-project", Project: "other",
		StateID: "state-id", State: "Todo", StateType: "unstarted",
	}
	inner := relationIssuePairFake(map[string]string{})
	sequenced := newSequentialOpClient(inner)
	sequenced.failAt["IssueRelationCreate"] = 1
	sequenced.payloads["issue_relations"] = []string{
		emptyIssueRelationsJSON(),
		`{"issue":{"relations":{"nodes":[` + relationWriteJSON("related") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`,
	}

	_, err := CreateIssueRelation(
		context.Background(),
		withIssueAfterWrite(sequenced, moved),
		matchingTarget(),
		IssueRelationCreateRequest{IssueID: "LIT-1", RelatedIssueID: "LIT-2", Type: "related"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "relation changed an issue project")
}

func Test_CreateIssueRelation_keeps_pinned_project_in_allowlist_union(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := relationIssueJSON("related-issue-id", "LIT-2", "other-project", "other")
	inner := issueWriteFakeClient(map[string]string{
		"issue_relations": emptyIssueRelationsJSON(),
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
		"issueRelation": `{"issueRelation":` + relationWriteJSON("related") + `}`,
	})
	projects := newSequentialOpClient(inner)
	projects.payloads["project"] = []string{
		`{"project":` + projectJSON(projectFixture{ID: "other-project", Name: "other", Status: "Backlog"}) + `}`,
	}
	graphqlClient := issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: projects,
	}

	result, err := CreateIssueRelation(context.Background(), graphqlClient, matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"other-project"},
	})

	require.NoError(t, err)
	require.Equal(t, "project-id", result.Issue.ProjectID)
	require.Equal(t, "other-project", result.RelatedIssue.ProjectID)
}

func Test_CreateIssueRelation_refuses_cross_project_on_team_only_pin_without_allowlist(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "project-id", "fixture")
	second := relationIssueJSON("related-issue-id", "LIT-2", "other-project", "other")
	recorder := &mutationRecordingClient{inner: issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: issueWriteFakeClient(map[string]string{}),
	}}

	_, err := CreateIssueRelation(context.Background(), recorder, teamOnlyTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "--allowed-project")
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_refuses_first_issue_outside_allowlist(t *testing.T) {
	first := relationIssueJSON("issue-id", "LIT-1", "other-project", "other")
	second := relationIssueJSON("related-issue-id", "LIT-2", "project-id", "fixture")
	recorder := &mutationRecordingClient{inner: issueLookupFake{
		byID: map[string]string{
			"LIT-1": first, "issue-id": first,
			"LIT-2": second, "related-issue-id": second,
		},
		inner: issueWriteFakeClient(map[string]string{
			"project": `{"project":` + projectJSON(projectFixture{
				ID: "project-id", Name: "fixture", Status: "Backlog",
			}) + `}`,
		}),
	}}

	_, err := CreateIssueRelation(context.Background(), recorder, matchingTarget(), IssueRelationCreateRequest{
		IssueID:           "LIT-1",
		RelatedIssueID:    "LIT-2",
		Type:              "related",
		AllowedProjectIDs: []string{"project-id"},
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueRelationCreate"))
}

func Test_CreateIssueRelation_allows_same_project_on_team_only_pin(t *testing.T) {
	graphqlClient := relationIssuePairFake(map[string]string{
		"IssueRelationCreate": `{"issueRelationCreate":{"success":true,"issueRelation":` +
			relationWriteJSON("related") + `}}`,
	})

	result, err := CreateIssueRelation(context.Background(), graphqlClient, teamOnlyTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.NoError(t, err)
	require.Equal(t, "relation-id", result.ID)
}

func Test_CreateIssueRelation_reconcile_returns_error_when_relation_scan_fails(t *testing.T) {
	sequenced := newSequentialOpClient(relationIssuePairFake(map[string]string{
		"IssueRelationCreate": "",
	}))
	sequenced.failAt["IssueRelationCreate"] = 1
	sequenced.failAt["issue_relations"] = 2

	_, err := CreateIssueRelation(context.Background(), sequenced, matchingTarget(), IssueRelationCreateRequest{
		IssueID:        "LIT-1",
		RelatedIssueID: "LIT-2",
		Type:           "related",
	})

	require.ErrorContains(t, err, "injected issue_relations failure")
}
