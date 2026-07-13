package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func issueLabelAssociationIssueRead() string {
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

func Test_AddIssueLabel_attaches_team_scoped_label_when_team_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueAddLabel": `{"issueAddLabel":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "First",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}}`,
	})

	issue, err := AddIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_AddIssueLabel_attaches_organization_wide_label(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "", "") + `}`,
		"IssueAddLabel": `{"issueAddLabel":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "First",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}}`,
	})

	issue, err := AddIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_AddIssueLabel_refuses_when_label_team_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "other-team", "OTHER") + `}`,
	})}

	_, err := AddIssueLabel(context.Background(), recorder, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueAddLabel"))
}

func Test_AddIssueLabel_refuses_when_issue_team_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue": `{"issue":{
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
		}}`,
	})}

	_, err := AddIssueLabel(context.Background(), recorder, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueAddLabel"))
}

func Test_AddIssueLabel_requires_issue_and_label_ids(t *testing.T) {
	_, err := AddIssueLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		IssueLabelAssociationRequest{IssueID: "LIT-1"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_AddIssueLabel_refuses_when_target_unresolved(t *testing.T) {
	_, err := AddIssueLabel(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, IssueLabelAssociationRequest{IssueID: "LIT-1", LabelID: "label-id"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddIssueLabel_wraps_label_resolution_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue": issueLabelAssociationIssueRead(),
	})

	_, err := AddIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddIssueLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
	})

	_, err := AddIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_AddIssueLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":         issueLabelAssociationIssueRead(),
		"issueLabel":    `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueAddLabel": `{"issueAddLabel":{"success":false,"issue":null}}`,
	})

	_, err := AddIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RemoveIssueLabel_detaches_label_when_team_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueRemoveLabel": `{"issueRemoveLabel":{"success":true,"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-1",
			Title:      "First",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "state-id",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `}}`,
	})

	issue, err := RemoveIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.NoError(t, err)
	require.Equal(t, "LIT-1", issue.Identifier)
}

func Test_RemoveIssueLabel_requires_issue_and_label_ids(t *testing.T) {
	_, err := RemoveIssueLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		IssueLabelAssociationRequest{LabelID: "label-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RemoveIssueLabel_wraps_issue_resolution_error(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := RemoveIssueLabel(context.Background(), recorder, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueRemoveLabel"))
}

func Test_RemoveIssueLabel_refuses_when_label_team_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "other-team", "OTHER") + `}`,
	})}

	_, err := RemoveIssueLabel(context.Background(), recorder, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueRemoveLabel"))
}

func Test_RemoveIssueLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":      issueLabelAssociationIssueRead(),
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
	})

	_, err := RemoveIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RemoveIssueLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issue":            issueLabelAssociationIssueRead(),
		"issueLabel":       `{"issueLabel":` + issueLabelJSON("label-id", "bug", "team-id", "LIT") + `}`,
		"IssueRemoveLabel": `{"issueRemoveLabel":{"success":false,"issue":null}}`,
	})

	_, err := RemoveIssueLabel(context.Background(), graphqlClient, matchingTarget(), IssueLabelAssociationRequest{
		IssueID: "LIT-1", LabelID: "label-id",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}
