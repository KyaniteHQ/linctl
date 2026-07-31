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

// Test_InitiativeLabelFamilyCommandFlows_report_runtime_errors covers the
// buildCommandRuntime error branch inside retire/restore RunE closures.
func Test_InitiativeLabelFamilyCommandFlows_report_runtime_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "initiative-label retire", args: []string{"initiative-label", "retire", "initiative-label-id", "--org-wide"}},
		{name: "initiative-label restore", args: []string{"initiative-label", "restore", "initiative-label-id", "--org-wide"}},
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

// Test_InitiativeLabelFamilyCommandFlows_report_writer_errors covers the
// render step: the guarded write succeeds but stdout fails.
func Test_InitiativeLabelFamilyCommandFlows_report_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "initiative-label retire", args: []string{"initiative-label", "retire", "initiative-label-id", "--org-wide"}},
		{name: "initiative-label restore", args: []string{"initiative-label", "restore", "initiative-label-id", "--org-wide"}},
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

// Test_InitiativeLabelFamilyCommands_require_org_wide proves InitiativeLabel
// taxonomy writes refuse as invalid writes when --org-wide is omitted.
func Test_InitiativeLabelFamilyCommands_require_org_wide(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "initiative-label retire", args: []string{"initiative-label", "retire", "initiative-label-id"}},
		{name: "initiative-label restore", args: []string{"initiative-label", "restore", "initiative-label-id"}},
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

// Test_InitiativeLabelRetire_emits_stable_target_mismatch_output proves a
// retire refuses closed with TARGET_MISMATCH when the label belongs to a
// different organization.
func Test_InitiativeLabelRetire_emits_stable_target_mismatch_output(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{otherOrgInitiativeLabel: true})
	defer restore()

	var stdout, stderr bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&stderr,
		[]string{"initiative-label", "retire", "initiative-label-id", "--org-wide"},
	)

	require.ErrorIs(t, err, client.ErrTargetMismatch)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), `"error_code":"TARGET_MISMATCH"`)
}

// Test_InitiativeLabelFamilyCommands_have_no_bypass_flags proves --org-wide
// is never a force/confirm bypass.
func Test_InitiativeLabelFamilyCommands_have_no_bypass_flags(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})
	for _, path := range []string{
		"initiative-label retire", "initiative-label restore",
	} {
		command, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		require.Nil(t, command.Flags().Lookup("force"))
		require.Nil(t, command.Flags().Lookup("confirm"))
		require.NotNil(t, command.Flags().Lookup("org-wide"))
	}
}
