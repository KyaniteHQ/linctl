package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

const cliWorkflowStateID = "550e8400-e29b-41d4-a716-446655440000"

func cliWorkflowStateJSON(name string, description string, position float64) string {
	return `{
		"id":"` + cliWorkflowStateID + `",
		"name":"` + name + `",
		"type":"unstarted",
		"color":"#f2c94c",
		"description":"` + description + `",
		"position":` + jsonNumber(position) + `,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func jsonNumber(value float64) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "0"
	}

	return string(encoded)
}

func workflowStateWriteFake() graphqlPayloadOverride {
	entity := cliWorkflowStateJSON("Ready", "ready items", 3)

	return graphqlPayloadOverride{
		inner: commandFlowFakeClient{},
		payloads: map[string]string{
			"WorkflowStateCreate": `{"workflowStateCreate":{"success":true,"workflowState":` + entity + `}}`,
			"WorkflowStateUpdate": `{"workflowStateUpdate":{"success":true,"workflowState":` + entity + `}}`,
			"workflowState":       `{"workflowState":` + entity + `}`,
		},
	}
}

func Test_WorkflowStateWriteCommands_have_no_bypass_or_team_flags(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})
	for _, path := range []string{"workflow-state create", "workflow-state update"} {
		command, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		require.Nil(t, command.Flags().Lookup("force"))
		require.Nil(t, command.Flags().Lookup("confirm"))
		require.Nil(t, command.Flags().Lookup("team-id"))
		require.Nil(t, command.Flags().Lookup("org-wide"))
	}
	update, _, err := root.Find([]string{"workflow-state", "update"})
	require.NoError(t, err)
	require.Nil(t, update.Flags().Lookup("type"))
	require.Nil(t, update.Flags().Lookup("id"))
}

func Test_WorkflowStateWriteCommandFlows_cover_output_modes(t *testing.T) {
	createArgs := []string{
		"workflow-state", "create",
		"--id", cliWorkflowStateID,
		"--name", "Ready",
		"--type", "unstarted",
		"--color", "#f2c94c",
		"--description", "ready items",
		"--position", "3",
	}
	updateArgs := []string{"workflow-state", "update", cliWorkflowStateID, "--name", "Ready"}
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "create human", args: createArgs, want: cliWorkflowStateID + " Ready [unstarted]\n"},
		{name: "create json", args: append([]string{"--json", "--compact"}, createArgs...), want: `"id":"` + cliWorkflowStateID + `"`},
		{name: "create id-only", args: append([]string{"--id-only"}, createArgs...), want: cliWorkflowStateID + "\n"},
		{name: "create quiet", args: append([]string{"--quiet"}, createArgs...), want: ""},
		{name: "update human", args: updateArgs, want: cliWorkflowStateID + " Ready [unstarted]\n"},
		{
			name: "update optional flags",
			args: []string{
				"workflow-state", "update", cliWorkflowStateID,
				"--name", "Ready", "--color", "#f2c94c",
				"--description", "ready items", "--position", "3",
			},
			want: cliWorkflowStateID + " Ready [unstarted]\n",
		},
		{name: "update fields", args: append([]string{"--json", "--fields", "id,name,type"}, updateArgs...), want: `"name"`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, workflowStateWriteFake())
			defer restore()
			var stdout bytes.Buffer
			err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &bytes.Buffer{}, test.args)

			require.NoError(t, err)
			require.Contains(t, stdout.String(), test.want)
		})
	}
}

func Test_WorkflowStateWriteCommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	commands := [][]string{
		{"workflow-state", "create", "--id", cliWorkflowStateID, "--name", "Ready", "--type", "unstarted", "--color", "#f2c94c"},
		{"workflow-state", "update", cliWorkflowStateID, "--name", "Ready"},
	}
	for _, args := range commands {
		t.Run("runtime "+strings.Join(args[:2], " "), func(t *testing.T) {
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			}
			defer func() { buildCommandRuntime = original }()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "runtime failed")
		})
		t.Run("writer "+strings.Join(args[:2], " "), func(t *testing.T) {
			restore := useCommandRuntime(t, workflowStateWriteFake())
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "write failed")
		})
	}
}

func Test_WorkflowStateUpdate_emits_stable_target_mismatch_output(t *testing.T) {
	restore := useCommandRuntime(t, graphqlPayloadOverride{
		inner: commandFlowFakeClient{},
		payloads: map[string]string{
			"workflowState": `{"workflowState":` + `{
				"id":"` + cliWorkflowStateID + `",
				"name":"Ready",
				"type":"unstarted",
				"color":"#f2c94c",
				"description":"",
				"position":1,
				"team":{"id":"other-team","key":"OTHER","name":"other"}
			}` + `}`,
		},
	})
	defer restore()
	var stdout, stderr bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&stderr,
		[]string{"workflow-state", "update", cliWorkflowStateID, "--name", "Ready"},
	)

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), `"error_code":"TARGET_MISMATCH"`)
}

func Test_WorkflowStateCreate_emits_conflict_without_replaying(t *testing.T) {
	restore := useCommandRuntime(t, graphqlPayloadOverride{
		inner: commandFlowFakeClient{},
		payloads: map[string]string{
			"WorkflowStateCreate": `{"workflowStateCreate":{"success":true,"workflowState":` +
				cliWorkflowStateJSON("Ready", "ready items", 3) + `}}`,
			"workflowState": `{"workflowState":` + cliWorkflowStateJSON("Other", "ready items", 3) + `}`,
		},
	})
	defer restore()
	var stdout, stderr bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&stderr,
		[]string{
			"workflow-state", "create",
			"--id", cliWorkflowStateID,
			"--name", "Ready",
			"--type", "unstarted",
			"--color", "#f2c94c",
		},
	)

	require.ErrorIs(t, err, client.ErrWriteConflict)
	require.Equal(t, "CONFLICT", errorCode(err))
	require.Contains(t, stderr.String(), `"error_code":"CONFLICT"`)
}
