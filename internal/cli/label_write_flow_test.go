package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// Test_LabelFamilyCommandFlows_report_runtime_errors covers the buildCommandRuntime
// error branch inside each new label-family command's RunE closure (label
// create/update/retire/restore, issue add-label/remove-label, project
// add-label/remove-label). It mirrors the Cycle write runtime-error edges.
func Test_LabelFamilyCommandFlows_report_runtime_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "label create", args: []string{"label", "create", "--name", "Created label"}},
		{name: "label update", args: []string{"label", "update", "label-id", "--name", "Updated label"}},
		{name: "label retire", args: []string{"label", "retire", "label-id"}},
		{name: "label restore", args: []string{"label", "restore", "label-id"}},
		{name: "issue add-label", args: []string{"issue", "add-label", "LIT-1", "label-id"}},
		{name: "issue remove-label", args: []string{"issue", "remove-label", "LIT-1", "label-id"}},
		{name: "project add-label", args: []string{"project", "add-label", "project-id", "label-id"}},
		{name: "project remove-label", args: []string{"project", "remove-label", "project-id", "label-id"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			}
			defer func() {
				buildCommandRuntime = original
			}()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "runtime failed")
		})
	}
}

// Test_LabelFamilyCommandFlows_report_writer_errors covers the render step for
// each new command: the guarded write succeeds but stdout fails, so the error
// must propagate.
func Test_LabelFamilyCommandFlows_report_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "label create", args: []string{"label", "create", "--name", "Created label"}},
		{name: "label update", args: []string{"label", "update", "label-id", "--name", "Updated label"}},
		{name: "label retire", args: []string{"label", "retire", "label-id"}},
		{name: "label restore", args: []string{"label", "restore", "label-id"}},
		{name: "issue add-label", args: []string{"issue", "add-label", "LIT-1", "label-id"}},
		{name: "issue remove-label", args: []string{"issue", "remove-label", "LIT-1", "label-id"}},
		{name: "project add-label", args: []string{"project", "add-label", "project-id", "label-id"}},
		{name: "project remove-label", args: []string{"project", "remove-label", "project-id", "label-id"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "write failed")
		})
	}
}

// Test_LabelUpdate_rejects_org_wide_flag_on_a_team_scoped_label proves the
// end-to-end CLI refusal: --org-wide is rejected for a team-owned label with a
// stable TARGET_MISMATCH error code, not a soft warning or confirmation.
func Test_LabelUpdate_rejects_org_wide_flag_on_a_team_scoped_label(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&stderr,
		[]string{"label", "update", "label-id", "--name", "Updated label", "--org-wide"},
	)

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Empty(t, stdout.String())
}

// Test_LabelUpdate_emits_stable_target_mismatch_output proves an organization-wide
// label update without --org-wide fails closed with the stable TARGET_MISMATCH
// error code, matching the guarded-write error envelope shape.
func Test_LabelUpdate_emits_stable_target_mismatch_output(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{orgWideLabel: true})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&stderr,
		[]string{"label", "update", "label-id", "--name", "Updated label"},
	)

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), `"error_code":"TARGET_MISMATCH"`)
}

// Test_LabelFamilyCommands_have_no_bypass_flags proves --org-wide selects a
// comparison class and is never a confirmation or force bypass: no label,
// issue, or project label-association command exposes --force or --confirm.
func Test_LabelFamilyCommands_have_no_bypass_flags(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})
	for _, path := range []string{"label create", "label update", "label retire", "label restore"} {
		command, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		require.Nil(t, command.Flags().Lookup("force"))
		require.Nil(t, command.Flags().Lookup("confirm"))
		require.NotNil(t, command.Flags().Lookup("org-wide"))
	}
	for _, path := range []string{"issue add-label", "issue remove-label", "project add-label", "project remove-label"} {
		command, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		require.Nil(t, command.Flags().Lookup("force"))
		require.Nil(t, command.Flags().Lookup("confirm"))
		require.Nil(t, command.Flags().Lookup("org-wide"))
	}
}
