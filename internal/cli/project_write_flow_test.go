package cli

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// Test_ProjectCommandFlows_report_project_write_writer_errors covers the render
// step for project create/update/archive: the guarded write succeeds but stdout
// fails, so the error must propagate. It mirrors the Cycle and ProjectMilestone
// writer-error tests so the project surface is verified directly.
func Test_ProjectCommandFlows_report_project_write_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "create", args: []string{"project", "create", "--name", "Created project"}},
		{name: "update", args: []string{"project", "update", "project-id", "--name", "Updated project"}},
		{name: "archive", args: []string{"project", "archive", "project-id"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{expectedProjectCreateName: "Created project"})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "write failed")
		})
	}
}

// Test_ProjectWriteCommands_reject_content_and_content_file_together covers the
// --content/--content-file mutual exclusion on both project write commands.
func Test_ProjectWriteCommands_reject_content_and_content_file_together(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "create", args: []string{
			"project", "create", "--name", "Created project",
			"--content", "inline", "--content-file", "body.md",
		}},
		{name: "update", args: []string{
			"project", "update", "project-id",
			"--content", "inline", "--content-file", "body.md",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorIs(t, err, client.ErrWriteInvalid)
		})
	}
}
