package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// mutationRecordingClient wraps a graphql.Client and records every operation
// name sent, so a guard refusal can be proven to send no mutation.
type mutationRecordingClient struct {
	inner graphql.Client
	sent  []string
}

func (client *mutationRecordingClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	client.sent = append(client.sent, request.OpName)

	return client.inner.MakeRequest(ctx, request, response)
}

func (client *mutationRecordingClient) sentOperation(name string) bool {
	for _, operation := range client.sent {
		if operation == name {
			return true
		}
	}

	return false
}

// variableCapturingClient wraps a graphql.Client and captures the JSON
// variables of each request by operation name, so outbound mutation payloads
// (resolved teamId, omitted teamId, replaceTeamLabels) can be asserted.
type variableCapturingClient struct {
	inner     graphql.Client
	variables map[string]string
}

func (client *variableCapturingClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if client.variables == nil {
		client.variables = map[string]string{}
	}
	if request.Variables != nil {
		if data, err := json.Marshal(request.Variables); err == nil {
			client.variables[request.OpName] = string(data)
		}
	}

	return client.inner.MakeRequest(ctx, request, response)
}

func issueLabelJSON(id string, name string, teamID string, teamKey string) string {
	team := "null"
	if teamID != "" {
		team = `{"id":"` + teamID + `","key":"` + teamKey + `","name":"linctl-it"}`
	}

	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":"",
		"color":"#000000",
		"isGroup":false,
		"team":` + team + `
	}`
}

func Test_CreateLabel_creates_team_scoped_label_with_resolved_team_id(t *testing.T) {
	capture := &variableCapturingClient{inner: issueWriteFakeClient(map[string]string{
		"IssueLabelCreate": `{"issueLabelCreate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "created", "team-id", "LIT") + `}}`,
	})}

	label, err := CreateLabel(
		context.Background(), capture, matchingTarget(),
		LabelCreateRequest{Name: "created", Color: "#123456"},
	)

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
	require.Equal(t, "team-id", label.TeamID)
	require.Contains(t, capture.variables["IssueLabelCreate"], `"teamId":"team-id"`)
	require.Contains(t, capture.variables["IssueLabelCreate"], `"replaceTeamLabels":false`)
}

func Test_CreateLabel_org_wide_omits_team_id(t *testing.T) {
	capture := &variableCapturingClient{inner: issueWriteFakeClient(map[string]string{
		"IssueLabelCreate": `{"issueLabelCreate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "created", "", "") + `}}`,
	})}

	label, err := CreateLabel(
		context.Background(), capture, matchingTarget(), LabelCreateRequest{Name: "created", OrgWide: true},
	)

	require.NoError(t, err)
	require.Empty(t, label.TeamID)
	require.NotContains(t, capture.variables["IssueLabelCreate"], `"teamId"`)
	require.Contains(t, capture.variables["IssueLabelCreate"], `"replaceTeamLabels":false`)
}

func Test_CreateLabel_requires_name(t *testing.T) {
	_, err := CreateLabel(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), LabelCreateRequest{})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateLabel_rejects_invalid_color(t *testing.T) {
	_, err := CreateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelCreateRequest{Name: "created", Color: "notacolor"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "color must be #RRGGBB")
}

func Test_validateLabelColor_accepts_empty_and_six_digit_hex(t *testing.T) {
	require.NoError(t, validateLabelColor(""))
	require.NoError(t, validateLabelColor("#123456"))
	require.NoError(t, validateLabelColor("#AbCdEf"))
	require.ErrorIs(t, validateLabelColor("#fff"), ErrWriteInvalid)
	require.ErrorIs(t, validateLabelColor("123456"), ErrWriteInvalid)
}

func Test_CreateLabel_refuses_when_target_unresolved(t *testing.T) {
	_, err := CreateLabel(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong-id",
	}, LabelCreateRequest{Name: "created"})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateLabel_requires_parent_in_same_team_by_default(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("parent-id", "parent", "other-team", "OTHER") + `}`,
	})}

	_, err := CreateLabel(context.Background(), recorder, matchingTarget(), LabelCreateRequest{
		Name: "created", ParentID: "parent-id",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelCreate"))
}

func Test_CreateLabel_allows_parent_in_same_team(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("parent-id", "parent", "team-id", "LIT") + `}`,
		"IssueLabelCreate": `{"issueLabelCreate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "created", "team-id", "LIT") + `}}`,
	})

	_, err := CreateLabel(context.Background(), graphqlClient, matchingTarget(), LabelCreateRequest{
		Name: "created", ParentID: "parent-id",
	})

	require.NoError(t, err)
}

func Test_CreateLabel_org_wide_requires_org_wide_parent(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("parent-id", "parent", "team-id", "LIT") + `}`,
	})}

	_, err := CreateLabel(context.Background(), recorder, matchingTarget(), LabelCreateRequest{
		Name: "created", ParentID: "parent-id", OrgWide: true,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelCreate"))
}

func Test_CreateLabel_org_wide_allows_org_wide_parent(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("parent-id", "parent", "", "") + `}`,
		"IssueLabelCreate": `{"issueLabelCreate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "created", "", "") + `}}`,
	})

	_, err := CreateLabel(context.Background(), graphqlClient, matchingTarget(), LabelCreateRequest{
		Name: "created", ParentID: "parent-id", OrgWide: true,
	})

	require.NoError(t, err)
}

func Test_CreateLabel_wraps_parent_resolution_error(t *testing.T) {
	_, err := CreateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelCreateRequest{Name: "created", ParentID: "missing-parent"},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateLabel_wraps_mutation_error(t *testing.T) {
	_, err := CreateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelCreateRequest{Name: "created"},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	_, err := CreateLabel(context.Background(), issueWriteFakeClient(map[string]string{
		"IssueLabelCreate": `{"issueLabelCreate":{"success":false,"issueLabel":` +
			issueLabelJSON("label-id", "created", "team-id", "LIT") + `}}`,
	}), matchingTarget(), LabelCreateRequest{Name: "created"})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_UpdateLabel_updates_team_scoped_label_when_team_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "team-id", "LIT") + `}`,
		"IssueLabelUpdate": `{"issueLabelUpdate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "renamed", "team-id", "LIT") + `}}`,
	})

	label, err := UpdateLabel(context.Background(), graphqlClient, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.NoError(t, err)
	require.Equal(t, "renamed", label.Name)
}

func Test_UpdateLabel_refuses_when_team_id_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "other-team", "LIT") + `}`,
	})}

	_, err := UpdateLabel(context.Background(), recorder, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelUpdate"))
}

func Test_UpdateLabel_refuses_when_team_key_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "team-id", "OTHER") + `}`,
	})}

	_, err := UpdateLabel(context.Background(), recorder, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelUpdate"))
}

func Test_UpdateLabel_refuses_organization_wide_label_without_org_wide_flag(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "", "") + `}`,
	})}

	_, err := UpdateLabel(context.Background(), recorder, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelUpdate"))
}

func Test_UpdateLabel_updates_organization_wide_label_with_org_wide_flag(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "", "") + `}`,
		"IssueLabelUpdate": `{"issueLabelUpdate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "renamed", "", "") + `}}`,
	})

	label, err := UpdateLabel(context.Background(), graphqlClient, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed", OrgWide: true,
	})

	require.NoError(t, err)
	require.Equal(t, "renamed", label.Name)
}

func Test_UpdateLabel_refuses_org_wide_flag_on_team_scoped_label(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "team-id", "LIT") + `}`,
	})}

	_, err := UpdateLabel(context.Background(), recorder, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed", OrgWide: true,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("IssueLabelUpdate"))
}

func Test_UpdateLabel_proceeds_when_pinned_project_present(t *testing.T) {
	// matchingTarget() pins project-id; taxonomy writes compare org+team only,
	// so a pinned project must not block this write (103 decision memo ruling).
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "team-id", "LIT") + `}`,
		"IssueLabelUpdate": `{"issueLabelUpdate":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "renamed", "team-id", "LIT") + `}}`,
	})

	_, err := UpdateLabel(context.Background(), graphqlClient, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.NoError(t, err)
}

func Test_UpdateLabel_requires_id(t *testing.T) {
	_, err := UpdateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelUpdateRequest{Name: "renamed"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateLabel_requires_a_field(t *testing.T) {
	_, err := UpdateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelUpdateRequest{ID: "label-id"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateLabel_rejects_invalid_color(t *testing.T) {
	_, err := UpdateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelUpdateRequest{ID: "label-id", Color: "notacolor"},
	)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "color must be #RRGGBB")
}

func Test_UpdateLabel_wraps_resolution_error(t *testing.T) {
	_, err := UpdateLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(),
		LabelUpdateRequest{ID: "label-id", Name: "renamed"},
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "team-id", "LIT") + `}`,
	})

	_, err := UpdateLabel(context.Background(), graphqlClient, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "old", "team-id", "LIT") + `}`,
		"IssueLabelUpdate": `{"issueLabelUpdate":{"success":false,"issueLabel":` +
			issueLabelJSON("label-id", "renamed", "team-id", "LIT") + `}}`,
	})

	_, err := UpdateLabel(context.Background(), graphqlClient, matchingTarget(), LabelUpdateRequest{
		ID: "label-id", Name: "renamed",
	})

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RetireLabel_retires_team_scoped_label_when_team_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "LIT") + `}`,
		"IssueLabelRetire": `{"issueLabelRetire":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "label", "team-id", "LIT") + `}}`,
	})

	label, err := RetireLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
}

func Test_RetireLabel_requires_id(t *testing.T) {
	_, err := RetireLabel(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "", false)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RetireLabel_refuses_organization_wide_label_without_org_wide_flag(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "", "") + `}`,
	})}

	_, err := RetireLabel(context.Background(), recorder, matchingTarget(), "label-id", false)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelRetire"))
}

func Test_RetireLabel_retires_organization_wide_label_with_org_wide_flag(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "", "") + `}`,
		"IssueLabelRetire": `{"issueLabelRetire":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "label", "", "") + `}}`,
	})

	_, err := RetireLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", true)

	require.NoError(t, err)
}

func Test_RetireLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "LIT") + `}`,
	})

	_, err := RetireLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RetireLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "LIT") + `}`,
		"IssueLabelRetire": `{"issueLabelRetire":{"success":false,"issueLabel":` +
			issueLabelJSON("label-id", "label", "team-id", "LIT") + `}}`,
	})

	_, err := RetireLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)

	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_RestoreLabel_restores_team_scoped_label_when_team_matches(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "LIT") + `}`,
		"IssueLabelRestore": `{"issueLabelRestore":{"success":true,"issueLabel":` +
			issueLabelJSON("label-id", "label", "team-id", "LIT") + `}}`,
	})

	label, err := RestoreLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)

	require.NoError(t, err)
	require.Equal(t, "label-id", label.ID)
}

func Test_RestoreLabel_requires_id(t *testing.T) {
	_, err := RestoreLabel(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "", false)

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_RestoreLabel_refuses_when_team_key_differs(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "OTHER") + `}`,
	})}

	_, err := RestoreLabel(context.Background(), recorder, matchingTarget(), "label-id", false)

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("IssueLabelRestore"))
}

func Test_RestoreLabel_wraps_resolution_error(t *testing.T) {
	_, err := RestoreLabel(
		context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), "label-id", false,
	)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RestoreLabel_wraps_mutation_error(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "LIT") + `}`,
	})

	_, err := RestoreLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_RestoreLabel_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"issueLabel": `{"issueLabel":` + issueLabelJSON("label-id", "label", "team-id", "LIT") + `}`,
		"IssueLabelRestore": `{"issueLabelRestore":{"success":false,"issueLabel":` +
			issueLabelJSON("label-id", "label", "team-id", "LIT") + `}}`,
	})

	_, err := RestoreLabel(context.Background(), graphqlClient, matchingTarget(), "label-id", false)

	require.ErrorIs(t, err, ErrMutationFailed)
}
