package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

const testTemplateID = "7c9e6679-7425-40de-944b-e07fc1f90ae7"

func issueTemplateJSON(name string, teamID string, teamKey string, pipelineID string, data string) string {
	team := "null"
	if teamID != "" {
		team = fmt.Sprintf(`{"id":%q,"key":%q,"name":"linctl"}`, teamID, teamKey)
	}
	pipeline := "null"
	if pipelineID != "" {
		pipeline = fmt.Sprintf(`{"id":%q}`, pipelineID)
	}

	return fmt.Sprintf(
		`{"id":%q,"name":%q,"type":"issue","templateData":%s,"team":%s,"pipeline":%s}`,
		testTemplateID, name, data, team, pipeline,
	)
}

func matchingIssueTemplateJSON(name string, data string) string {
	return issueTemplateJSON(name, "team-id", "LIT", "", data)
}

func templateContentResponse(name string, data string) string {
	return `{"template":` + matchingIssueTemplateJSON(name, data) + `}`
}

func templateCreateResponse(name string, data string) string {
	return `{"templateCreate":{"success":true,"template":` + matchingIssueTemplateJSON(name, data) + `}}`
}

func templateUpdateResponse(name string) string {
	return `{"templateUpdate":{"success":true,"template":` + matchingIssueTemplateJSON(name, `{"title":"Bug"}`) + `}}`
}

func Test_GetTemplateDetail_returns_canonical_data_and_scope(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"templateContent": templateContentResponse("Bug report", `{"title":"Bug","labels":["a"]}`),
	}

	detail, err := GetTemplateDetail(context.Background(), graphqlClient, testTemplateID)

	require.NoError(t, err)
	require.Equal(t, testTemplateID, detail.ID)
	require.Equal(t, "Bug report", detail.Name)
	require.Equal(t, "issue", detail.Type)
	require.Equal(t, "team-id", detail.TeamID)
	require.Equal(t, "LIT", detail.TeamKey)
	require.Empty(t, detail.PipelineID)
	require.JSONEq(t, `{"labels":["a"],"title":"Bug"}`, string(detail.Data))
}

func Test_GetTemplateDetail_unwraps_encoded_object_layer(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"templateContent": templateContentResponse("Bug report", `"{\"title\":\"Encoded\"}"`),
	}

	detail, err := GetTemplateDetail(context.Background(), graphqlClient, testTemplateID)

	require.NoError(t, err)
	require.JSONEq(t, `{"title":"Encoded"}`, string(detail.Data))
}

func Test_GetTemplateDetail_requires_template_id(t *testing.T) {
	_, err := GetTemplateDetail(context.Background(), fakeGraphQLClient{}, "")
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_GetTemplateDetail_wraps_read_error(t *testing.T) {
	_, err := GetTemplateDetail(context.Background(), fakeGraphQLClient{}, testTemplateID)
	require.ErrorContains(t, err, "get template content")
}

func Test_GetTemplateDetail_keeps_pipeline_scope_on_read(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"templateContent": `{"template":` + issueTemplateJSON(
			"Bug report", "team-id", "LIT", "pipeline-id", `{"title":"Bug"}`,
		) + `}`,
	}

	detail, err := GetTemplateDetail(context.Background(), graphqlClient, testTemplateID)

	require.NoError(t, err)
	require.Equal(t, "pipeline-id", detail.PipelineID)
}

func Test_GetTemplateDetail_rejects_non_object_data(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"templateContent": templateContentResponse("Bug report", `[]`),
	}

	_, err := GetTemplateDetail(context.Background(), graphqlClient, testTemplateID)

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.NotContains(t, err.Error(), "[")
}

func Test_CreateTemplate_sends_caller_id_issue_type_and_team(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate":  templateCreateResponse("Bug report", `{"title":"Bug"}`),
		"templateContent": templateContentResponse("Bug report", `{ "title": "Bug" }`),
	})}

	detail, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID:   testTemplateID,
		Name: "Bug report",
		Type: "issue",
		Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.NoError(t, err)
	require.Equal(t, testTemplateID, detail.ID)
	require.Equal(t, "issue", detail.Type)
	require.Equal(t, "team-id", detail.TeamID)
	require.JSONEq(t, `{"title":"Bug"}`, string(detail.Data))
	require.Equal(t, 1, recorder.countOf("TemplateCreate"))
	require.Equal(t, 1, recorder.countOf("templateContent"))
	require.JSONEq(t, `{
		"input": {
			"id": "`+testTemplateID+`",
			"name": "Bug report",
			"type": "issue",
			"teamId": "team-id",
			"templateData": {"title":"Bug"}
		}
	}`, string(recorder.variablesFor(t, "TemplateCreate")))
}

func Test_CreateTemplate_requires_uuid_v4(t *testing.T) {
	_, err := CreateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateCreateRequest{
		ID:   "not-a-uuid",
		Name: "Bug report",
		Type: "issue",
		Data: json.RawMessage(`{"title":"Bug"}`),
	})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateTemplate_requires_name(t *testing.T) {
	_, err := CreateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateCreateRequest{
		ID:   testTemplateID,
		Type: "issue",
		Data: json.RawMessage(`{"title":"Bug"}`),
	})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateTemplate_reconciles_success_false_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate": `{"templateCreate":{"success":false,"template":` +
			matchingIssueTemplateJSON("Bug report", `{"title":"Bug"}`) + `}}`,
		"templateContent": templateContentResponse("Bug report", `{"title":"Bug"}`),
	})}

	detail, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.NoError(t, err)
	require.Equal(t, testTemplateID, detail.ID)
	require.Equal(t, 1, recorder.countOf("TemplateCreate"))
}

func Test_CreateTemplate_rejects_non_issue_types(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}
	for _, templateType := range []string{"project", "document", "releaseNote", ""} {
		t.Run(templateType, func(t *testing.T) {
			_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
				ID:   testTemplateID,
				Name: "Bug report",
				Type: templateType,
				Data: json.RawMessage(`{"title":"Bug"}`),
			})

			require.ErrorIs(t, err, ErrWriteInvalid)
			require.False(t, recorder.sentOperation("TemplateCreate"))
		})
	}
}

func Test_CreateTemplate_rejects_non_object_data_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID:   testTemplateID,
		Name: "Bug report",
		Type: "issue",
		Data: json.RawMessage(`[]`),
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("TemplateCreate"))
}

func Test_CreateTemplate_rejects_encoded_object_layer_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID:   testTemplateID,
		Name: "Bug report",
		Type: "issue",
		Data: json.RawMessage(`"{\"title\":\"Bug\"}"`),
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("TemplateCreate"))
}

func Test_CreateTemplate_refuses_when_target_mismatches(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateTemplate(context.Background(), recorder, config.Target{
		OrgID: "other-org", TeamKey: "LIT", TeamID: "team-id", ProjectID: "project-id",
	}, TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TemplateCreate"))
}

func Test_CreateTemplate_reconciles_ambiguous_create_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate":  "",
		"templateContent": templateContentResponse("Bug report", `{"title":"Bug"}`),
	}).withError(errors.New("timeout"))}

	detail, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.NoError(t, err)
	require.Equal(t, testTemplateID, detail.ID)
	require.Equal(t, 1, recorder.countOf("TemplateCreate"))
	require.Equal(t, 1, recorder.countOf("templateContent"))
}

func Test_CreateTemplate_returns_conflict_when_data_differs(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate":  "",
		"templateContent": templateContentResponse("Bug report", `{"title":"Other"}`),
	}).withError(errors.New("timeout"))}

	_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.NotContains(t, err.Error(), "Other")
	require.Equal(t, 1, recorder.countOf("TemplateCreate"))
}

func Test_CreateTemplate_returns_target_mismatch_when_team_is_null(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate": "",
		"templateContent": `{"template":` + issueTemplateJSON(
			"Bug report", "", "", "", `{"title":"Bug"}`,
		) + `}`,
	}).withError(errors.New("timeout"))}

	_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateTemplate_returns_target_mismatch_when_pipeline_is_set(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate": "",
		"templateContent": `{"template":` + issueTemplateJSON(
			"Bug report", "team-id", "LIT", "pipeline-id", `{"title":"Bug"}`,
		) + `}`,
	}).withError(errors.New("timeout"))}

	_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_CreateTemplate_returns_original_error_when_id_is_absent(t *testing.T) {
	writeErr := errors.New("timeout")
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"TemplateCreate":  "",
		"templateContent": `{"template":null}`,
	}).withError(writeErr)}

	_, err := CreateTemplate(context.Background(), recorder, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"title":"Bug"}`),
	})

	require.ErrorIs(t, err, writeErr)
}

func Test_UpdateTemplate_omits_team_and_type(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": templateContentResponse("Updated", `{"title":"Bug"}`),
		"TemplateUpdate":  templateUpdateResponse("Updated"),
	})}
	name := "Updated"

	detail, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.NoError(t, err)
	require.Equal(t, "Updated", detail.Name)
	require.JSONEq(t, `{
		"id": "`+testTemplateID+`",
		"input": {"name": "Updated"}
	}`, string(recorder.variablesFor(t, "TemplateUpdate")))
	require.NotContains(t, string(recorder.variablesFor(t, "TemplateUpdate")), "teamId")
	require.NotContains(t, string(recorder.variablesFor(t, "TemplateUpdate")), `"type"`)
}

func Test_UpdateTemplate_sends_canonical_data(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": templateContentResponse("Bug report", `{"title":"Bug"}`),
		"TemplateUpdate":  templateUpdateResponse("Bug report"),
	})}

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Data: json.RawMessage(`{ "title": "Bug" }`),
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"id": "`+testTemplateID+`",
		"input": {"templateData": {"title":"Bug"}}
	}`, string(recorder.variablesFor(t, "TemplateUpdate")))
}

func Test_UpdateTemplate_refuses_organization_template_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": `{"template":` + issueTemplateJSON(
			"Bug report", "", "", "", `{"title":"Bug"}`,
		) + `}`,
	})}
	name := "Updated"

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TemplateUpdate"))
}

func Test_UpdateTemplate_refuses_pipeline_template_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": `{"template":` + issueTemplateJSON(
			"Bug report", "team-id", "LIT", "pipeline-id", `{"title":"Bug"}`,
		) + `}`,
	})}
	name := "Updated"

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TemplateUpdate"))
}

func Test_UpdateTemplate_refuses_non_issue_template_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": `{"template":{"id":"` + testTemplateID + `","name":"Doc","type":"document",` +
			`"templateData":{"title":"Doc"},"team":{"id":"team-id","key":"LIT","name":"linctl"},"pipeline":null}}`,
	})}
	name := "Updated"

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.False(t, recorder.sentOperation("TemplateUpdate"))
}

func Test_UpdateTemplate_refuses_wrong_team_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": `{"template":` + issueTemplateJSON(
			"Bug report", "other-team", "OTHER", "", `{"title":"Bug"}`,
		) + `}`,
	})}
	name := "Updated"

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("TemplateUpdate"))
}

func Test_UpdateTemplate_reconciles_ambiguous_update_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": templateContentResponse("Updated", `{"title":"Bug"}`),
		"TemplateUpdate":  "",
	}).withError(errors.New("timeout"))}
	name := "Updated"

	detail, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.NoError(t, err)
	require.Equal(t, "Updated", detail.Name)
	require.Equal(t, 1, recorder.countOf("TemplateUpdate"))
	require.Equal(t, 2, recorder.countOf("templateContent"))
}

func Test_UpdateTemplate_returns_conflict_when_name_differs(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": templateContentResponse("Bug report", `{"title":"Bug"}`),
		"TemplateUpdate":  "",
	}).withError(errors.New("timeout"))}
	name := "Updated"

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.Equal(t, 1, recorder.countOf("TemplateUpdate"))
}

func Test_UpdateTemplate_returns_conflict_when_readback_id_differs(t *testing.T) {
	const otherTemplateID = "11111111-1111-4111-8111-111111111111"
	name := "Updated"
	wrongIDReadback := `{"template":{"id":"` + otherTemplateID + `","name":"Updated","type":"issue",` +
		`"templateData":{"title":"Bug"},"team":{"id":"team-id","key":"LIT","name":"linctl"},"pipeline":null}}`
	recorder := &recordingGraphQLClient{inner: &sequentialPayloadClient{
		inner: issueWriteFakeClient(map[string]string{
			"TemplateUpdate": templateUpdateResponse("Updated"),
		}),
		payloads: map[string][]string{
			"templateContent": {
				templateContentResponse("Updated", `{"title":"Bug"}`),
				wrongIDReadback,
			},
		},
	}}

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.NotContains(t, err.Error(), otherTemplateID)
	require.Equal(t, 1, recorder.countOf("TemplateUpdate"))
}

func Test_UpdateTemplate_refuses_when_target_mismatches(t *testing.T) {
	name := "Updated"
	_, err := UpdateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "other-org", TeamKey: "LIT", TeamID: "team-id", ProjectID: "project-id",
	}, TemplateUpdateRequest{ID: testTemplateID, Name: &name})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateTemplate_rejects_empty_name_and_invalid_data(t *testing.T) {
	empty := ""
	_, err := UpdateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &empty,
	})
	require.ErrorIs(t, err, ErrWriteInvalid)

	_, err = UpdateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Data: json.RawMessage(`[]`),
	})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateTemplate_returns_conflict_when_data_or_type_differs(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": templateContentResponse("Bug report", `{"title":"Other"}`),
		"TemplateUpdate":  "",
	}).withError(errors.New("timeout"))}

	_, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Data: json.RawMessage(`{"title":"Bug"}`),
	})
	require.ErrorIs(t, err, ErrWriteConflict)

	name := "Bug report"
	typeClient := &sequentialPayloadClient{
		inner: issueWriteFakeClient(map[string]string{
			"TemplateUpdate": templateUpdateResponse("Bug report"),
		}),
		payloads: map[string][]string{
			"templateContent": {
				templateContentResponse("Bug report", `{"title":"Bug"}`),
				`{"template":{"id":"` + testTemplateID + `","name":"Bug report","type":"document",` +
					`"templateData":{"title":"Bug"},"team":{"id":"team-id","key":"LIT","name":"linctl"},"pipeline":null}}`,
			},
		},
	}
	_, err = UpdateTemplate(context.Background(), typeClient, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})
	require.ErrorIs(t, err, ErrWriteConflict)
}

func Test_UpdateTemplate_wraps_missing_template_read(t *testing.T) {
	name := "Updated"
	_, err := UpdateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})
	require.ErrorContains(t, err, "get template content")
}

func Test_UpdateTemplate_reconciles_success_false_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"templateContent": templateContentResponse("Updated", `{"title":"Bug"}`),
		"TemplateUpdate": `{"templateUpdate":{"success":false,"template":` +
			matchingIssueTemplateJSON("Updated", `{"title":"Bug"}`) + `}}`,
	})}
	name := "Updated"

	detail, err := UpdateTemplate(context.Background(), recorder, matchingTarget(), TemplateUpdateRequest{
		ID:   testTemplateID,
		Name: &name,
	})

	require.NoError(t, err)
	require.Equal(t, "Updated", detail.Name)
}

func Test_UpdateTemplate_requires_template_id(t *testing.T) {
	name := "Updated"
	_, err := UpdateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateUpdateRequest{
		Name: &name,
	})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateTemplate_requires_a_changed_field(t *testing.T) {
	_, err := UpdateTemplate(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), TemplateUpdateRequest{
		ID: testTemplateID,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_CreateTemplate_preserves_number_lexemes_in_readback(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"TemplateCreate":  templateCreateResponse("Bug report", `{"n":1.0}`),
		"templateContent": templateContentResponse("Bug report", `{"n":1.0}`),
	})

	detail, err := CreateTemplate(context.Background(), graphqlClient, matchingTarget(), TemplateCreateRequest{
		ID: testTemplateID, Name: "Bug report", Type: "issue", Data: json.RawMessage(`{"n":1.0}`),
	})

	require.NoError(t, err)
	require.Equal(t, `{"n":1.0}`, string(detail.Data)) //nolint:testifylint // lexeme must stay 1.0
}
