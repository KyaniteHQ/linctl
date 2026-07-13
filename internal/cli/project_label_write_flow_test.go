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

// Test_ProjectLabelFamilyCommandFlows_report_runtime_errors covers the
// buildCommandRuntime error branch inside each new project-label command's
// RunE closure (create/update/retire/restore). It mirrors the Label write
// runtime-error edges.
func Test_ProjectLabelFamilyCommandFlows_report_runtime_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "project-label create", args: []string{"project-label", "create", "--name", "Created", "--org-wide"}},
		{
			name: "project-label update",
			args: []string{"project-label", "update", "project-label-id", "--name", "Updated", "--org-wide"},
		},
		{name: "project-label retire", args: []string{"project-label", "retire", "project-label-id", "--org-wide"}},
		{name: "project-label restore", args: []string{"project-label", "restore", "project-label-id", "--org-wide"}},
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

// Test_ProjectLabelFamilyCommandFlows_report_writer_errors covers the render
// step for each new command: the guarded write succeeds but stdout fails, so
// the error must propagate.
func Test_ProjectLabelFamilyCommandFlows_report_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "project-label create", args: []string{"project-label", "create", "--name", "Created", "--org-wide"}},
		{
			name: "project-label update",
			args: []string{"project-label", "update", "project-label-id", "--name", "Updated", "--org-wide"},
		},
		{name: "project-label retire", args: []string{"project-label", "retire", "project-label-id", "--org-wide"}},
		{name: "project-label restore", args: []string{"project-label", "restore", "project-label-id", "--org-wide"}},
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

// Test_ProjectLabelFamilyCommands_require_org_wide proves the end-to-end CLI
// refusal: ProjectLabel has no team scope, so every taxonomy write is
// refused as an invalid write (not a soft warning) when --org-wide is
// omitted, and the error names the flag.
func Test_ProjectLabelFamilyCommands_require_org_wide(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "project-label create", args: []string{"project-label", "create", "--name", "Created"}},
		{name: "project-label update", args: []string{"project-label", "update", "project-label-id", "--name", "Updated"}},
		{name: "project-label retire", args: []string{"project-label", "retire", "project-label-id"}},
		{name: "project-label restore", args: []string{"project-label", "restore", "project-label-id"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()

			var stdout, stderr bytes.Buffer
			err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &stderr, test.args)

			require.ErrorIs(t, err, client.ErrWriteInvalid)
			require.Empty(t, stdout.String())
			require.Contains(t, stderr.String(), "org-wide")
		})
	}
}

// Test_ProjectLabelUpdate_emits_stable_target_mismatch_output proves a
// ProjectLabel update refuses closed with the stable TARGET_MISMATCH error
// code when the resolved label belongs to a different organization, matching
// the guarded-write error envelope shape.
func Test_ProjectLabelUpdate_emits_stable_target_mismatch_output(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{otherOrgProjectLabel: true})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&stderr,
		[]string{"project-label", "update", "project-label-id", "--name", "Updated", "--org-wide"},
	)

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), `"error_code":"TARGET_MISMATCH"`)
}

// Test_ProjectLabelFamilyCommands_have_no_bypass_flags proves --org-wide
// selects the Org-Scoped Write comparison class and is never a confirmation
// or force bypass: no project-label command exposes --force or --confirm.
func Test_ProjectLabelFamilyCommands_have_no_bypass_flags(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})
	for _, path := range []string{
		"project-label create", "project-label update", "project-label retire", "project-label restore",
	} {
		command, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		require.Nil(t, command.Flags().Lookup("force"))
		require.Nil(t, command.Flags().Lookup("confirm"))
		require.NotNil(t, command.Flags().Lookup("org-wide"))
	}
}
