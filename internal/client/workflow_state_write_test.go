package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

const testWorkflowStateID = "550e8400-e29b-41d4-a716-446655440000"

func workflowStateEntityJSON(
	name string,
	stateType string,
	teamID string,
	teamKey string,
	description string,
	position float64,
) string {
	return fmt.Sprintf(
		`{"id":%q,"name":%q,"type":%q,"color":%q,"description":%q,"position":%s,`+
			`"team":{"id":%q,"key":%q,"name":"linctl"}}`,
		testWorkflowStateID, name, stateType, "#f2c94c", description,
		strconv.FormatFloat(position, 'f', -1, 64),
		teamID, teamKey,
	)
}

func matchingWorkflowStateJSON(name string, stateType string, description string, position float64) string {
	return workflowStateEntityJSON(name, stateType, "team-id", "LIT", description, position)
}

func workflowStateCreateResponse() string {
	return `{"workflowStateCreate":{"success":true}}`
}

func workflowStateUpdateResponse() string {
	return `{"workflowStateUpdate":{"success":true}}`
}

func workflowStateGetResponse(name string, stateType string, description string, position float64) string {
	return `{"workflowState":` + matchingWorkflowStateJSON(name, stateType, description, position) + `}`
}

func Test_CreateWorkflowState_sends_caller_id_and_resolved_team(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": workflowStateCreateResponse(),
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "ready items", 3),
	})}
	description := "ready items"
	position := 3.0

	state, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID:          testWorkflowStateID,
		Name:        "Ready",
		Type:        "unstarted",
		Color:       "#f2c94c",
		Description: &description,
		Position:    &position,
	})

	require.NoError(t, err)
	require.Equal(t, testWorkflowStateID, state.ID)
	require.Equal(t, "Ready", state.Name)
	require.Equal(t, "unstarted", state.Type)
	require.Equal(t, "#f2c94c", state.Color)
	require.Equal(t, "ready items", state.Description)
	require.InDelta(t, 3.0, state.Position, 0)
	require.Equal(t, "team-id", state.TeamID)
	require.Equal(t, "LIT", state.TeamKey)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
	require.Equal(t, 1, recorder.countOf("workflowState"))
	require.JSONEq(t, `{
		"input": {
			"id": "`+testWorkflowStateID+`",
			"name": "Ready",
			"type": "unstarted",
			"color": "#f2c94c",
			"teamId": "team-id",
			"description": "ready items",
			"position": 3
		}
	}`, string(recorder.variablesFor(t, "WorkflowStateCreate")))
}

func Test_CreateWorkflowState_omits_optional_fields_and_keeps_observed_defaults(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": workflowStateCreateResponse(),
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "default desc", 8),
	})}

	state, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID:    testWorkflowStateID,
		Name:  "Ready",
		Type:  "unstarted",
		Color: "#f2c94c",
	})

	require.NoError(t, err)
	require.Equal(t, "default desc", state.Description)
	require.InDelta(t, 8.0, state.Position, 0)
	require.JSONEq(t, `{
		"input": {
			"id": "`+testWorkflowStateID+`",
			"name": "Ready",
			"type": "unstarted",
			"color": "#f2c94c",
			"teamId": "team-id"
		}
	}`, string(recorder.variablesFor(t, "WorkflowStateCreate")))
}

func Test_CreateWorkflowState_accepts_supported_types(t *testing.T) {
	for _, stateType := range []string{"backlog", "unstarted", "started", "completed", "canceled"} {
		t.Run(stateType, func(t *testing.T) {
			graphqlClient := issueWriteFakeClient(map[string]string{
				"WorkflowStateCreate": workflowStateCreateResponse(),
				"workflowState":       workflowStateGetResponse("Ready", stateType, "", 1),
			})

			state, err := CreateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateCreateRequest{
				ID:    testWorkflowStateID,
				Name:  "Ready",
				Type:  stateType,
				Color: "#f2c94c",
			})

			require.NoError(t, err)
			require.Equal(t, stateType, state.Type)
		})
	}
}

func Test_CreateWorkflowState_rejects_disallowed_types(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}
	for _, stateType := range []string{"triage", "duplicate", "in-progress", "todo", "unknown"} {
		t.Run(stateType, func(t *testing.T) {
			_, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
				ID:    testWorkflowStateID,
				Name:  "Ready",
				Type:  stateType,
				Color: "#f2c94c",
			})

			require.ErrorIs(t, err, ErrWriteInvalid)
			require.False(t, recorder.sentOperation("WorkflowStateCreate"))
		})
	}
}

func Test_CreateWorkflowState_requires_id_name_type_and_color(t *testing.T) {
	tests := []WorkflowStateCreateRequest{
		{Name: "Ready", Type: "unstarted", Color: "#f2c94c"},
		{ID: testWorkflowStateID, Type: "unstarted", Color: "#f2c94c"},
		{ID: testWorkflowStateID, Name: "Ready", Color: "#f2c94c"},
		{ID: testWorkflowStateID, Name: "Ready", Type: "unstarted"},
		{ID: "not-a-uuid", Name: "Ready", Type: "unstarted", Color: "#f2c94c"},
		{ID: "11111111-1111-1111-1111-111111111111", Name: "Ready", Type: "unstarted", Color: "#f2c94c"},
	}
	for _, request := range tests {
		_, err := CreateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), request)

		require.ErrorIs(t, err, ErrWriteInvalid)
	}
}

func Test_CreateWorkflowState_refuses_when_target_mismatches(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{})}

	_, err := CreateWorkflowState(context.Background(), recorder, config.Target{
		OrgID: "other-org", TeamKey: "LIT", TeamID: "team-id", ProjectID: "project-id",
	}, WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("WorkflowStateCreate"))
}

func Test_CreateWorkflowState_refuses_when_pinned_project_is_missing(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"TargetProject": "",
	})}

	_, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.Error(t, err)
	require.False(t, recorder.sentOperation("WorkflowStateCreate"))
}

func Test_CreateWorkflowState_allows_team_owned_write_with_valid_project_pin(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": workflowStateCreateResponse(),
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "", 1),
	})

	_, err := CreateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.NoError(t, err)
}

func Test_CreateWorkflowState_allows_team_only_pin(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": workflowStateCreateResponse(),
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "", 1),
	})

	_, err := CreateWorkflowState(context.Background(), graphqlClient, teamOnlyTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.NoError(t, err)
}

func Test_CreateWorkflowState_reconciles_ambiguous_create_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": "",
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "ready items", 3),
	}).withError(errors.New("timeout"))}
	description := "ready items"
	position := 3.0

	state, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID:          testWorkflowStateID,
		Name:        "Ready",
		Type:        "unstarted",
		Color:       "#f2c94c",
		Description: &description,
		Position:    &position,
	})

	require.NoError(t, err)
	require.Equal(t, testWorkflowStateID, state.ID)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
	require.Equal(t, 1, recorder.countOf("workflowState"))
}

func Test_CreateWorkflowState_returns_conflict_when_reconciled_fields_differ(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": "",
		"workflowState":       workflowStateGetResponse("Other", "started", "ready items", 3),
	}).withError(errors.New("timeout"))}

	_, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
	require.Equal(t, 1, recorder.countOf("workflowState"))
}

func Test_CreateWorkflowState_returns_target_mismatch_when_reconciled_team_differs(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": "",
		"workflowState": `{"workflowState":` + workflowStateEntityJSON(
			"Ready", "unstarted", "other-team", "OTHER", "", 1,
		) + `}`,
	}).withError(errors.New("timeout"))}

	_, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
}

func Test_CreateWorkflowState_returns_original_error_when_id_is_absent(t *testing.T) {
	writeErr := errors.New("timeout")
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": "",
		"workflowState":       `{"workflowState":null}`,
	}).withError(writeErr)}

	_, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, writeErr)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
	require.Equal(t, 1, recorder.countOf("workflowState"))
}

func Test_CreateWorkflowState_returns_readback_error_when_reconciliation_fails(t *testing.T) {
	writeErr := errors.New("timeout")
	readErr := errors.New("read failed")
	recorder := &recordingGraphQLClient{inner: &perOpErrorClient{
		inner: issueWriteFakeClient(map[string]string{}),
		errs: map[string]error{
			"WorkflowStateCreate": writeErr,
			"workflowState":       readErr,
		},
	}}

	_, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, readErr)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
	require.Equal(t, 1, recorder.countOf("workflowState"))
}

func Test_CreateWorkflowState_returns_conflict_when_successful_mutation_readback_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": workflowStateCreateResponse(),
		"workflowState":       workflowStateGetResponse("Other", "unstarted", "", 1),
	})

	_, err := CreateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, ErrWriteConflict)
}

func Test_CreateWorkflowState_returns_readback_error_after_successful_mutation(t *testing.T) {
	readErr := errors.New("read failed")
	graphqlClient := &perOpErrorClient{
		inner: issueWriteFakeClient(map[string]string{
			"WorkflowStateCreate": workflowStateCreateResponse(),
		}),
		errs: map[string]error{
			"workflowState": readErr,
		},
	}

	_, err := CreateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.ErrorIs(t, err, readErr)
}

func Test_CreateWorkflowState_reconciles_success_false_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": `{"workflowStateCreate":{"success":false}}`,
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "", 1),
	})}

	state, err := CreateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateCreateRequest{
		ID: testWorkflowStateID, Name: "Ready", Type: "unstarted", Color: "#f2c94c",
	})

	require.NoError(t, err)
	require.Equal(t, testWorkflowStateID, state.ID)
	require.Equal(t, 1, recorder.countOf("WorkflowStateCreate"))
}

func Test_UpdateWorkflowState_updates_changed_fields_with_presence(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"workflowState":       workflowStateGetResponse("Next", "unstarted", "", 0),
		"WorkflowStateUpdate": workflowStateUpdateResponse(),
	})}
	name := "Next"
	color := "#f2c94c"
	description := ""
	position := 0.0

	state, err := UpdateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateUpdateRequest{
		ID:          testWorkflowStateID,
		Name:        &name,
		Color:       &color,
		Description: &description,
		Position:    &position,
	})

	require.NoError(t, err)
	require.Equal(t, "Next", state.Name)
	require.Empty(t, state.Description)
	require.InDelta(t, 0.0, state.Position, 0)
	require.Equal(t, 1, recorder.countOf("WorkflowStateUpdate"))
	require.JSONEq(t, `{
		"id": "`+testWorkflowStateID+`",
		"input": {
			"name": "Next",
			"color": "#f2c94c",
			"description": "",
			"position": 0
		}
	}`, string(recorder.variablesFor(t, "WorkflowStateUpdate")))
}

func Test_UpdateWorkflowState_omits_unchanged_fields(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"workflowState":       workflowStateGetResponse("Next", "unstarted", "old", 3),
		"WorkflowStateUpdate": workflowStateUpdateResponse(),
	})}
	name := "Next"

	_, err := UpdateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.NoError(t, err)
	require.JSONEq(t, `{
		"id": "`+testWorkflowStateID+`",
		"input": {"name": "Next"}
	}`, string(recorder.variablesFor(t, "WorkflowStateUpdate")))
}

func Test_UpdateWorkflowState_refuses_wrong_team_before_mutation(t *testing.T) {
	recorder := &mutationRecordingClient{inner: issueWriteFakeClient(map[string]string{
		"workflowState": `{"workflowState":` + workflowStateEntityJSON(
			"Ready", "unstarted", "other-team", "OTHER", "", 1,
		) + `}`,
	})}
	name := "Next"

	_, err := UpdateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("WorkflowStateUpdate"))
}

func Test_UpdateWorkflowState_refuses_when_only_team_id_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"workflowState": `{"workflowState":` + workflowStateEntityJSON(
			"Ready", "unstarted", "other-team", "LIT", "", 1,
		) + `}`,
	})
	name := "Next"

	_, err := UpdateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateWorkflowState_refuses_when_only_team_key_differs(t *testing.T) {
	graphqlClient := issueWriteFakeClient(map[string]string{
		"workflowState": `{"workflowState":` + workflowStateEntityJSON(
			"Ready", "unstarted", "team-id", "OTHER", "", 1,
		) + `}`,
	})
	name := "Next"

	_, err := UpdateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateWorkflowState_requires_id(t *testing.T) {
	name := "Next"
	_, err := UpdateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), WorkflowStateUpdateRequest{
		Name: &name,
	})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateWorkflowState_refuses_when_target_mismatches(t *testing.T) {
	name := "Next"
	_, err := UpdateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID: "other-org", TeamKey: "LIT", TeamID: "team-id", ProjectID: "project-id",
	}, WorkflowStateUpdateRequest{ID: testWorkflowStateID, Name: &name})

	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_UpdateWorkflowState_wraps_missing_state_read(t *testing.T) {
	name := "Next"
	_, err := UpdateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})
	require.ErrorContains(t, err, "get workflow state")
}

func Test_UpdateWorkflowState_reconciles_success_false_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"workflowState":       workflowStateGetResponse("Next", "unstarted", "", 1),
		"WorkflowStateUpdate": `{"workflowStateUpdate":{"success":false}}`,
	})}
	name := "Next"

	state, err := UpdateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.NoError(t, err)
	require.Equal(t, "Next", state.Name)
}

func Test_UpdateWorkflowState_returns_conflict_when_optional_fields_differ(t *testing.T) {
	color := "#00ff00"
	description := "new"
	position := 9.0
	graphqlClient := issueWriteFakeClient(map[string]string{
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "old", 1),
		"WorkflowStateUpdate": "",
	})

	_, err := UpdateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateUpdateRequest{
		ID:    testWorkflowStateID,
		Color: &color,
	})
	require.ErrorIs(t, err, ErrWriteConflict)

	_, err = UpdateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateUpdateRequest{
		ID:          testWorkflowStateID,
		Description: &description,
	})
	require.ErrorIs(t, err, ErrWriteConflict)

	_, err = UpdateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateUpdateRequest{
		ID:       testWorkflowStateID,
		Position: &position,
	})
	require.ErrorIs(t, err, ErrWriteConflict)
}

func Test_CreateWorkflowState_returns_conflict_when_optional_fields_differ(t *testing.T) {
	description := "wanted"
	position := 4.0
	graphqlClient := issueWriteFakeClient(map[string]string{
		"WorkflowStateCreate": workflowStateCreateResponse(),
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "other", 1),
	})

	_, err := CreateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateCreateRequest{
		ID:          testWorkflowStateID,
		Name:        "Ready",
		Type:        "unstarted",
		Color:       "#f2c94c",
		Description: &description,
	})
	require.ErrorIs(t, err, ErrWriteConflict)

	_, err = CreateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateCreateRequest{
		ID:       testWorkflowStateID,
		Name:     "Ready",
		Type:     "unstarted",
		Color:    "#f2c94c",
		Position: &position,
	})
	require.ErrorIs(t, err, ErrWriteConflict)
}

func Test_UpdateWorkflowState_returns_conflict_when_type_changes(t *testing.T) {
	name := "Ready"
	graphqlClient := &sequentialPayloadClient{
		inner: issueWriteFakeClient(map[string]string{
			"WorkflowStateUpdate": workflowStateUpdateResponse(),
		}),
		payloads: map[string][]string{
			"workflowState": {
				workflowStateGetResponse("Ready", "unstarted", "", 1),
				workflowStateGetResponse("Ready", "started", "", 1),
			},
		},
	}

	_, err := UpdateWorkflowState(context.Background(), graphqlClient, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrWriteConflict)
}

func Test_UpdateWorkflowState_requires_a_changed_field(t *testing.T) {
	_, err := UpdateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), WorkflowStateUpdateRequest{
		ID: testWorkflowStateID,
	})

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateWorkflowState_rejects_empty_name_and_color(t *testing.T) {
	empty := ""
	_, err := UpdateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &empty,
	})
	require.ErrorIs(t, err, ErrWriteInvalid)

	_, err = UpdateWorkflowState(context.Background(), issueWriteFakeClient(map[string]string{}), matchingTarget(), WorkflowStateUpdateRequest{
		ID:    testWorkflowStateID,
		Color: &empty,
	})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_UpdateWorkflowState_reconciles_ambiguous_update_without_replay(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"workflowState":       workflowStateGetResponse("Next", "unstarted", "", 1),
		"WorkflowStateUpdate": "",
	}).withError(errors.New("timeout"))}
	name := "Next"

	state, err := UpdateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.NoError(t, err)
	require.Equal(t, "Next", state.Name)
	require.Equal(t, 1, recorder.countOf("WorkflowStateUpdate"))
	require.Equal(t, 2, recorder.countOf("workflowState"))
}

func Test_UpdateWorkflowState_returns_conflict_when_reconciled_fields_differ(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: issueWriteFakeClient(map[string]string{
		"workflowState":       workflowStateGetResponse("Ready", "unstarted", "", 1),
		"WorkflowStateUpdate": "",
	}).withError(errors.New("timeout"))}
	name := "Next"

	_, err := UpdateWorkflowState(context.Background(), recorder, matchingTarget(), WorkflowStateUpdateRequest{
		ID:   testWorkflowStateID,
		Name: &name,
	})

	require.ErrorIs(t, err, ErrWriteConflict)
	require.Equal(t, 1, recorder.countOf("WorkflowStateUpdate"))
}

type perOpErrorClient struct {
	inner graphql.Client
	errs  map[string]error
}

func (client *perOpErrorClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := client.errs[request.OpName]; err != nil {
		return err
	}

	return client.inner.MakeRequest(ctx, request, response)
}

type sequentialPayloadClient struct {
	inner    graphql.Client
	payloads map[string][]string
	seen     map[string]int
}

func (client *sequentialPayloadClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if client.seen == nil {
		client.seen = map[string]int{}
	}
	payloads := client.payloads[request.OpName]
	index := client.seen[request.OpName]
	if index < len(payloads) {
		client.seen[request.OpName]++
		wrapped := []byte(`{"data":` + payloads[index] + `}`)

		return json.Unmarshal(wrapped, response)
	}

	return client.inner.MakeRequest(ctx, request, response)
}
